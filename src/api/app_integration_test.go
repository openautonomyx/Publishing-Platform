package main

// End-to-end tests: real HTTP router, real PostgreSQL (db/schema.sql), real token
// verification against an in-process fake Logto issuer.
//
// They need DATABASE_URL (CI's Test job sets it). Each run creates and drops its own
// schema, so they never touch existing tables.

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/fractional-pm/creative-platform/logtoauth"
)

const testAudience = "https://api.publishing.test"

type testEnv struct {
	t      *testing.T
	api    *httptest.Server
	issuer *httptest.Server
	key    *ecdsa.PrivateKey
	db     *sql.DB
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	base := os.Getenv("DATABASE_URL")
	if base == "" {
		t.Skip("DATABASE_URL not set; skipping integration tests")
	}

	suffix := make([]byte, 6)
	if _, err := rand.Read(suffix); err != nil {
		t.Fatal(err)
	}
	schema := "it_" + hex.EncodeToString(suffix)
	admin, err := sql.Open("postgres", base)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { admin.Close() })
	if err := admin.Ping(); err != nil {
		t.Fatalf("DATABASE_URL is set but database is not reachable: %v", err)
	}
	if _, err := admin.Exec(`CREATE SCHEMA ` + schema); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := admin.Exec(`DROP SCHEMA ` + schema + ` CASCADE`); err != nil {
			t.Logf("drop schema: %v", err)
		}
	})

	dsn := base + map[bool]string{true: "&", false: "?"}[strings.Contains(base, "?")] + "search_path=" + schema + ",public"
	if !strings.Contains(base, "://") {
		dsn = base + " search_path=" + schema + ",public"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	ddl, err := os.ReadFile("../../db/schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(string(ddl)); err != nil {
		t.Fatalf("load schema: %v", err)
	}

	env := &testEnv{t: t, db: db}
	env.key, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	env.issuer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := func(i interface{ FillBytes([]byte) []byte }) string {
			return base64.RawURLEncoding.EncodeToString(i.FillBytes(make([]byte, 48)))
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"keys": []map[string]string{{
			"kty": "EC", "crv": "P-384", "kid": "k1", "use": "sig", "x": b(env.key.X), "y": b(env.key.Y)}}})
	}))
	t.Cleanup(env.issuer.Close)

	verifier, err := logtoauth.NewVerifier(logtoauth.Config{Endpoint: env.issuer.URL, Audience: testAudience}, env.issuer.Client())
	if err != nil {
		t.Fatal(err)
	}
	app := &App{db: db, verifier: verifier}
	env.api = httptest.NewServer(app.Router())
	t.Cleanup(env.api.Close)
	return env
}

// caller is someone signing in through Logto.
type caller struct {
	env    *testEnv
	tenant string
	sub    string
	scope  string
	email  string
	apiKey string
}

func (e *testEnv) as(tenant, sub, scope string) *caller {
	return &caller{env: e, tenant: tenant, sub: sub, scope: scope}
}

func (c *caller) token() string {
	claims := jwt.MapClaims{
		"iss": c.env.issuer.URL + "/oidc", "aud": testAudience, "sub": c.sub,
		"tenant_id": c.tenant, "scope": c.scope, "exp": time.Now().Add(time.Hour).Unix(),
	}
	if c.email != "" {
		claims["email"] = c.email
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodES384, claims)
	tok.Header["kid"] = "k1"
	s, err := tok.SignedString(c.env.key)
	if err != nil {
		c.env.t.Fatal(err)
	}
	return s
}

type resp struct {
	Status int
	Body   map[string]interface{}
	Raw    string
}

func (r resp) data() map[string]interface{} {
	d, _ := r.Body["data"].(map[string]interface{})
	return d
}

func (r resp) list() []interface{} {
	l, _ := r.Body["data"].([]interface{})
	return l
}

func (c *caller) do(method, path string, body interface{}) resp {
	c.env.t.Helper()
	var rd io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			c.env.t.Fatal(err)
		}
		rd = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, c.env.api.URL+path, rd)
	if err != nil {
		c.env.t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	} else {
		req.Header.Set("Authorization", "Bearer "+c.token())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.env.t.Fatal(err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	out := resp{Status: res.StatusCode, Raw: string(raw)}
	_ = json.Unmarshal(raw, &out.Body)
	return out
}

func expect(t *testing.T, r resp, status int) resp {
	t.Helper()
	if r.Status != status {
		t.Fatalf("status = %d, want %d; body: %s", r.Status, status, r.Raw)
	}
	return r
}

func str(m map[string]interface{}, k string) string {
	s, _ := m[k].(string)
	return s
}

// ----------------------------------------------------------------------------

func TestProvisioningAndOrganization(t *testing.T) {
	env := newTestEnv(t)
	admin := env.as("acme", "u_admin", "publishing:admin")

	orgs := expect(t, admin.do("GET", "/api/v1/organizations", nil), 200).list()
	if len(orgs) != 1 {
		t.Fatalf("organizations = %d, want 1", len(orgs))
	}
	org := orgs[0].(map[string]interface{})
	if str(org, "external_id") != "acme" || str(org, "slug") != "acme" {
		t.Errorf("provisioned org = %v", org)
	}
	orgID := str(org, "id")

	users := expect(t, admin.do("GET", "/api/v1/users", nil), 200).list()
	if len(users) != 1 || str(users[0].(map[string]interface{}), "role") != "admin" {
		t.Fatalf("users = %v", users)
	}

	// Second request reuses the same org and user.
	expect(t, admin.do("GET", "/api/v1/users", nil), 200)
	var n int
	if err := env.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&n); err != nil || n != 1 {
		t.Fatalf("users rows = %d (%v), want 1", n, err)
	}

	expect(t, admin.do("POST", "/api/v1/organizations", map[string]string{"name": "x"}), 403)
	upd := expect(t, admin.do("PUT", "/api/v1/organizations/"+orgID,
		map[string]interface{}{"name": "Acme Inc", "metadata": map[string]string{"plan": "x"}}), 200).data()
	if str(upd, "name") != "Acme Inc" {
		t.Errorf("updated org = %v", upd)
	}
	expect(t, env.as("acme", "u_editor", "publishing:edit").do("PUT", "/api/v1/organizations/"+orgID,
		map[string]string{"name": "Nope"}), 403)

	// Deleting the org locks everyone out of it.
	expect(t, admin.do("DELETE", "/api/v1/organizations/"+orgID, nil), 200)
	expect(t, admin.do("GET", "/api/v1/users", nil), 403)
}

func TestTenantIsolation(t *testing.T) {
	env := newTestEnv(t)
	a := env.as("tenant-a", "u_a", "publishing:admin")
	b := env.as("tenant-b", "u_b", "publishing:admin")

	c := expect(t, a.do("POST", "/api/v1/content",
		map[string]string{"title": "A's post", "content_type": "article"}), 201).data()
	id := str(c, "id")
	aUser := str(expect(t, a.do("GET", "/api/v1/users", nil), 200).list()[0].(map[string]interface{}), "id")
	aOrg := str(c, "org_id")

	expect(t, b.do("GET", "/api/v1/content/"+id, nil), 404)
	expect(t, b.do("PUT", "/api/v1/content/"+id, map[string]string{"title": "hijack"}), 404)
	expect(t, b.do("DELETE", "/api/v1/content/"+id, nil), 404)
	expect(t, b.do("GET", "/api/v1/content/"+id+"/approvals", nil), 404)
	expect(t, b.do("GET", "/api/v1/users/"+aUser, nil), 404)
	expect(t, b.do("GET", "/api/v1/organizations/"+aOrg, nil), 404)
	if l := expect(t, b.do("GET", "/api/v1/content", nil), 200).list(); len(l) != 0 {
		t.Errorf("tenant B sees %d content items", len(l))
	}
	if l := expect(t, b.do("GET", "/api/v1/audit-logs", nil), 200).list(); len(l) != 0 {
		t.Errorf("tenant B sees %d audit entries", len(l))
	}
	expect(t, b.do("POST", "/api/v1/distributions",
		map[string]string{"content_id": id, "channel": "twitter"}), 404)
}

func TestRBAC(t *testing.T) {
	env := newTestEnv(t)
	viewer := env.as("t", "u_view", "publishing:read")
	editor := env.as("t", "u_edit", "publishing:edit")
	approver := env.as("t", "u_appr", "publishing:approve")

	expect(t, viewer.do("GET", "/api/v1/content", nil), 200)
	expect(t, viewer.do("POST", "/api/v1/content", map[string]string{"title": "Nope", "content_type": "article"}), 403)
	expect(t, editor.do("POST", "/api/v1/content", map[string]string{"title": "Yes", "content_type": "article"}), 201)
	for _, c := range []*caller{editor, approver} {
		expect(t, c.do("POST", "/api/v1/workflows", map[string]string{"name": "WF", "workflow_type": "approval"}), 403)
		expect(t, c.do("GET", "/api/v1/audit-logs", nil), 403)
		expect(t, c.do("GET", "/api/v1/api-keys", nil), 403)
		expect(t, c.do("POST", "/api/v1/users", map[string]string{"email": "x@y.io", "name": "X"}), 403)
	}
	expect(t, editor.do("POST", "/api/v1/approvals", map[string]string{
		"content_id": "00000000-0000-0000-0000-000000000000", "stage_id": "00000000-0000-0000-0000-000000000000", "decision": "approved"}), 403)

	r := expect(t, viewer.do("POST", "/api/v1/content", map[string]string{"title": "x"}), 403)
	if str(r.Body, "error") != "FORBIDDEN" || r.Body["success"] != false {
		t.Errorf("error body = %s", r.Raw)
	}
}

func TestErrorsAndValidation(t *testing.T) {
	env := newTestEnv(t)
	ed := env.as("t", "u_edit", "publishing:edit")

	r := expect(t, ed.do("POST", "/api/v1/content", map[string]string{"title": "x", "content_type": "novel"}), 400)
	if str(r.Body, "error") != "VALIDATION_ERROR" || len(r.Body["details"].([]interface{})) != 2 {
		t.Errorf("validation body = %s", r.Raw)
	}
	expect(t, ed.do("POST", "/api/v1/content", map[string]string{"title": "Fine", "content_type": "article", "bogus": "x"}), 400)
	expect(t, ed.do("GET", "/api/v1/content/not-a-uuid", nil), 400)
	expect(t, ed.do("GET", "/api/v1/content/00000000-0000-0000-0000-000000000000", nil), 404)
	expect(t, ed.do("GET", "/api/v1/nope", nil), 404)

	res, err := http.Get(env.api.URL + "/api/v1/content")
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != 401 || res.Header.Get("X-Request-ID") == "" {
		t.Errorf("unauthenticated: status %d, request id %q", res.StatusCode, res.Header.Get("X-Request-ID"))
	}

	// CORS preflight is answered without auth.
	req, _ := http.NewRequest("OPTIONS", env.api.URL+"/api/v1/content", nil)
	req.Header.Set("Origin", "https://app.example")
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != 204 || res.Header.Get("Access-Control-Allow-Origin") == "" {
		t.Errorf("preflight: %d %v", res.StatusCode, res.Header)
	}

	// Pagination.
	for i := 0; i < 3; i++ {
		expect(t, ed.do("POST", "/api/v1/content", map[string]string{"title": fmt.Sprintf("Post %d", i), "content_type": "article"}), 201)
	}
	p := expect(t, ed.do("GET", "/api/v1/content?limit=2&page=2", nil), 200)
	if p.Body["total"] != float64(3) || len(p.list()) != 1 {
		t.Errorf("page 2 = %s", p.Raw)
	}
}

func TestApprovalWorkflowEndToEnd(t *testing.T) {
	env := newTestEnv(t)
	admin := env.as("t", "u_admin", "publishing:admin")
	author := env.as("t", "u_author", "publishing:edit")
	rev1 := env.as("t", "u_rev1", "publishing:approve")
	rev2 := env.as("t", "u_rev2", "publishing:approve")

	wf := expect(t, admin.do("POST", "/api/v1/workflows", map[string]interface{}{
		"name": "Editorial", "workflow_type": "approval",
		"stages": []map[string]interface{}{
			{"name": "Copy edit", "approver_role": "approver"},
			{"name": "Legal", "approver_role": "admin"},
		},
	}), 201).data()
	stages := wf["stages"].([]interface{})
	copyEdit := str(stages[0].(map[string]interface{}), "id")
	legal := str(stages[1].(map[string]interface{}), "id")

	c := expect(t, author.do("POST", "/api/v1/content", map[string]interface{}{
		"title": "Launch post", "content_type": "article", "workflow_id": wf["id"]}), 201).data()
	cid := str(c, "id")
	approve := func(who *caller, stage, decision string) resp {
		return who.do("POST", "/api/v1/approvals", map[string]string{"content_id": cid, "stage_id": stage, "decision": decision})
	}

	expect(t, approve(rev1, copyEdit, "approved"), 409) // still draft
	expect(t, author.do("PUT", "/api/v1/content/"+cid, map[string]string{"status": "published"}), 409)
	expect(t, author.do("PUT", "/api/v1/content/"+cid, map[string]string{"status": "in-review"}), 200)
	expect(t, author.do("PUT", "/api/v1/content/"+cid, map[string]string{"title": "Edited in review"}), 409)

	// Distribution needs approval first.
	expect(t, author.do("POST", "/api/v1/distributions", map[string]string{"content_id": cid, "channel": "linkedin"}), 409)

	expect(t, approve(rev1, legal, "approved"), 409)    // stage order
	expect(t, approve(rev1, copyEdit, "rejected"), 201) // back to draft
	got := expect(t, author.do("GET", "/api/v1/content/"+cid, nil), 200).data()
	if str(got, "status") != "draft" {
		t.Fatalf("after rejection status = %s", str(got, "status"))
	}

	// Resubmit: the old round's decisions no longer count.
	expect(t, author.do("PUT", "/api/v1/content/"+cid, map[string]string{"status": "in-review"}), 200)
	r := expect(t, approve(rev1, copyEdit, "approved"), 201).data()
	if str(r, "content_status") != "in-review" {
		t.Errorf("after stage 1: %v", r)
	}
	expect(t, approve(rev2, copyEdit, "approved"), 409) // stage already satisfied
	expect(t, approve(rev2, legal, "approved"), 403)    // legal needs admin
	r = expect(t, approve(admin, legal, "approved"), 201).data()
	if str(r, "content_status") != "approved" {
		t.Fatalf("after all stages: %v", r)
	}

	// Authors can't approve their own content.
	own := expect(t, admin.do("POST", "/api/v1/content", map[string]interface{}{
		"title": "Admin post", "content_type": "article", "workflow_id": wf["id"]}), 201).data()
	expect(t, admin.do("PUT", "/api/v1/content/"+str(own, "id"), map[string]string{"status": "in-review"}), 200)
	expect(t, admin.do("POST", "/api/v1/approvals", map[string]string{
		"content_id": str(own, "id"), "stage_id": copyEdit, "decision": "approved"}), 403)

	if l := expect(t, author.do("GET", "/api/v1/content/"+cid+"/approvals", nil), 200).list(); len(l) != 3 {
		t.Errorf("approvals = %d, want 3", len(l))
	}

	// Stages are locked once used.
	expect(t, admin.do("PUT", "/api/v1/workflows/"+str(wf, "id"), map[string]interface{}{"stages": []interface{}{}}), 409)

	d := expect(t, author.do("POST", "/api/v1/distributions", map[string]interface{}{
		"content_id": cid, "channel": "linkedin", "scheduled_at": time.Now().Add(time.Hour)}), 201).data()
	if str(d, "status") != "scheduled" {
		t.Errorf("distribution = %v", d)
	}
	pub := expect(t, author.do("PUT", "/api/v1/content/"+cid, map[string]string{"status": "published"}), 200).data()
	if str(pub, "status") != "published" || pub["published_at"] == nil {
		t.Errorf("published = %v", pub)
	}
	if l := expect(t, author.do("GET", "/api/v1/content/"+cid+"/distributions", nil), 200).list(); len(l) != 1 {
		t.Errorf("distributions = %d", len(l))
	}

	// Content without a workflow can be published directly.
	quick := expect(t, author.do("POST", "/api/v1/content", map[string]string{"title": "Quick note", "content_type": "social"}), 201).data()
	expect(t, author.do("PUT", "/api/v1/content/"+str(quick, "id"), map[string]string{"status": "published"}), 200)
}

func TestAPIKeys(t *testing.T) {
	env := newTestEnv(t)
	admin := env.as("t", "u_admin", "publishing:admin")
	editor := env.as("t", "u_edit", "publishing:edit")
	expect(t, editor.do("GET", "/api/v1/content", nil), 200) // provision editor

	k := expect(t, admin.do("POST", "/api/v1/api-keys", map[string]interface{}{
		"name": "CI bot", "permissions": []string{"read"}}), 201).data()
	key := str(k, "key")
	if !strings.HasPrefix(key, "pk_") || str(k, "key_prefix") != key[:11] {
		t.Fatalf("created key = %v", k)
	}
	bot := &caller{env: env, apiKey: key}
	expect(t, bot.do("GET", "/api/v1/content", nil), 200)
	expect(t, bot.do("POST", "/api/v1/content", map[string]string{"title": "Bot post", "content_type": "article"}), 403)

	listed := expect(t, admin.do("GET", "/api/v1/api-keys", nil), 200)
	if strings.Contains(listed.Raw, key) || strings.Contains(listed.Raw, "key_hash") {
		t.Fatal("key listing leaks the key or its hash")
	}

	expect(t, admin.do("DELETE", "/api/v1/api-keys/"+str(k, "id"), nil), 200)
	expect(t, bot.do("GET", "/api/v1/content", nil), 401)
	expect(t, (&caller{env: env, apiKey: "pk_wrong"}).do("GET", "/api/v1/content", nil), 401)

	// A write key works, but is capped by its owner's role once the owner is demoted.
	wk := str(expect(t, admin.do("POST", "/api/v1/api-keys", map[string]interface{}{
		"name": "Writer", "permissions": []string{"write"}}), 201).data(), "key")
	writer := &caller{env: env, apiKey: wk}
	expect(t, writer.do("POST", "/api/v1/content", map[string]string{"title": "Bot post", "content_type": "article"}), 201)
	expect(t, env.as("t", "u_admin", "publishing:read").do("GET", "/api/v1/content", nil), 200) // admin's token now read-only
	expect(t, writer.do("POST", "/api/v1/content", map[string]string{"title": "Bot post 2", "content_type": "article"}), 403)

	// A key can't outrank its creator.
	expect(t, env.as("t", "u_admin", "publishing:admin").do("GET", "/api/v1/content", nil), 200)
	expect(t, admin.do("POST", "/api/v1/api-keys", map[string]interface{}{"name": "Bad", "permissions": []string{"superuser"}}), 400)
}

func TestUsersInviteLinkDeactivate(t *testing.T) {
	env := newTestEnv(t)
	admin := env.as("t", "u_admin", "publishing:admin")

	inv := expect(t, admin.do("POST", "/api/v1/users", map[string]string{"email": "Jo@Example.com", "name": "Jo"}), 201).data()
	if str(inv, "status") != "invited" {
		t.Fatalf("invited user = %v", inv)
	}
	expect(t, admin.do("POST", "/api/v1/users", map[string]string{"email": "Jo@Example.com", "name": "Jo again"}), 409)

	jo := env.as("t", "u_jo", "publishing:edit")
	jo.email = "jo@example.com"
	expect(t, jo.do("GET", "/api/v1/content", nil), 200)
	linked := expect(t, admin.do("GET", "/api/v1/users/"+str(inv, "id"), nil), 200).data()
	if str(linked, "status") != "active" || str(linked, "external_id") != "u_jo" || str(linked, "role") != "editor" {
		t.Fatalf("linked user = %v", linked)
	}

	// Self-service profile edit; no status change without admin.
	expect(t, jo.do("PUT", "/api/v1/users/"+str(inv, "id"), map[string]string{"name": "Joanna"}), 200)
	expect(t, jo.do("PUT", "/api/v1/users/"+str(inv, "id"), map[string]string{"status": "active"}), 403)

	expect(t, admin.do("PUT", "/api/v1/users/"+str(inv, "id"), map[string]string{"status": "inactive"}), 200)
	expect(t, jo.do("GET", "/api/v1/content", nil), 403)

	adminID := str(expect(t, admin.do("GET", "/api/v1/users?status=active", nil), 200).list()[0].(map[string]interface{}), "id")
	expect(t, admin.do("DELETE", "/api/v1/users/"+adminID, nil), 400)
	expect(t, admin.do("DELETE", "/api/v1/users/"+str(inv, "id"), nil), 200)
	expect(t, jo.do("GET", "/api/v1/content", nil), 403)
}

func TestAgentsAndAudit(t *testing.T) {
	env := newTestEnv(t)
	admin := env.as("t", "u_admin", "publishing:admin")

	ag := expect(t, admin.do("POST", "/api/v1/agents", map[string]interface{}{
		"name": "Summariser", "agent_type": "content-creator", "config": map[string]string{"model": "x"}}), 201).data()
	aid := str(ag, "id")
	run := expect(t, admin.do("POST", "/api/v1/agents/"+aid+"/run", map[string]interface{}{"input": map[string]int{"n": 1}}), 202).data()
	if str(run, "status") != "queued" {
		t.Errorf("run = %v", run)
	}
	if l := expect(t, admin.do("GET", "/api/v1/agents/"+aid+"/runs", nil), 200).list(); len(l) != 1 {
		t.Errorf("runs = %d", len(l))
	}
	expect(t, admin.do("PUT", "/api/v1/agents/"+aid, map[string]string{"status": "paused"}), 200)
	expect(t, admin.do("POST", "/api/v1/agents/"+aid+"/run", map[string]interface{}{}), 409)
	expect(t, admin.do("DELETE", "/api/v1/agents/"+aid, nil), 200)
	expect(t, admin.do("GET", "/api/v1/agents/"+aid, nil), 404)

	logs := expect(t, admin.do("GET", "/api/v1/audit-logs?resource_type=agent", nil), 200)
	actions := map[string]bool{}
	for _, e := range logs.list() {
		actions[str(e.(map[string]interface{}), "action")] = true
	}
	for _, want := range []string{"agent.create", "agent.run", "agent.update", "agent.delete"} {
		if !actions[want] {
			t.Errorf("audit log missing %s (have %v)", want, actions)
		}
	}
}

func TestNullJSONBColumnsScan(t *testing.T) {
	env := newTestEnv(t)
	ed := env.as("t", "u_edit", "publishing:edit")
	c := expect(t, ed.do("POST", "/api/v1/content", map[string]string{"title": "Null meta", "content_type": "article"}), 201).data()
	if _, err := env.db.Exec(`UPDATE content SET metadata = NULL WHERE id = $1`, str(c, "id")); err != nil {
		t.Fatal(err)
	}
	got := expect(t, ed.do("GET", "/api/v1/content/"+str(c, "id"), nil), 200).data()
	if v, ok := got["metadata"]; !ok || v != nil {
		t.Errorf("metadata = %v (present=%v), want null", v, ok)
	}
}
