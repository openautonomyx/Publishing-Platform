package logtoauth

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const audience = "https://api.publishing.example"

// fakeLogto serves /oidc/jwks for a set of EC P-384 keys that tests can rotate.
type fakeLogto struct {
	srv     *httptest.Server
	mu      sync.Mutex
	keys    map[string]*ecdsa.PrivateKey
	fetches atomic.Int32
}

func newFakeLogto(t *testing.T) *fakeLogto {
	f := &fakeLogto{keys: map[string]*ecdsa.PrivateKey{}}
	f.addKey(t, "k1")
	f.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oidc/jwks" {
			http.NotFound(w, r)
			return
		}
		f.fetches.Add(1)
		f.mu.Lock()
		defer f.mu.Unlock()
		var set struct {
			Keys []map[string]string `json:"keys"`
		}
		for kid, k := range f.keys {
			set.Keys = append(set.Keys, map[string]string{
				"kty": "EC", "crv": "P-384", "kid": kid, "use": "sig", "alg": "ES384",
				"x": base64.RawURLEncoding.EncodeToString(k.X.FillBytes(make([]byte, 48))),
				"y": base64.RawURLEncoding.EncodeToString(k.Y.FillBytes(make([]byte, 48))),
			})
		}
		_ = json.NewEncoder(w).Encode(set)
	}))
	t.Cleanup(f.srv.Close)
	return f
}

func (f *fakeLogto) addKey(t *testing.T, kid string) {
	k, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	f.mu.Lock()
	f.keys[kid] = k
	f.mu.Unlock()
}

func (f *fakeLogto) sign(t *testing.T, kid string, claims jwt.MapClaims) string {
	f.mu.Lock()
	k := f.keys[kid]
	f.mu.Unlock()
	tok := jwt.NewWithClaims(jwt.SigningMethodES384, claims)
	tok.Header["kid"] = kid
	s, err := tok.SignedString(k)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func (f *fakeLogto) claims() jwt.MapClaims {
	now := time.Now()
	return jwt.MapClaims{
		"iss":       f.srv.URL + "/oidc",
		"aud":       audience,
		"sub":       "u_abc123",
		"client_id": "app_xyz",
		"scope":     "publishing:edit publishing:read",
		"tenant_id": "tenant-a",
		"iat":       now.Unix(),
		"exp":       now.Add(time.Hour).Unix(),
	}
}

func newVerifier(t *testing.T, f *fakeLogto, cfg Config) *Verifier {
	cfg.Endpoint = f.srv.URL
	if cfg.Audience == "" {
		cfg.Audience = audience
	}
	v, err := NewVerifier(cfg, f.srv.Client())
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func TestVerifyValidToken(t *testing.T) {
	f := newFakeLogto(t)
	v := newVerifier(t, f, Config{})

	p, err := v.Verify(context.Background(), f.sign(t, "k1", f.claims()))
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if p.Subject != "u_abc123" || p.ClientID != "app_xyz" || p.TenantID != "tenant-a" {
		t.Errorf("unexpected principal: %+v", p)
	}
	if p.Role() != "editor" {
		t.Errorf("Role() = %q, want editor", p.Role())
	}
}

func TestVerifyRejects(t *testing.T) {
	f := newFakeLogto(t)
	v := newVerifier(t, f, Config{})

	cases := map[string]func() string{
		"wrong issuer": func() string {
			c := f.claims()
			c["iss"] = "https://evil.example/oidc"
			return f.sign(t, "k1", c)
		},
		"wrong audience": func() string {
			c := f.claims()
			c["aud"] = "https://other-api.example"
			return f.sign(t, "k1", c)
		},
		"expired": func() string {
			c := f.claims()
			c["exp"] = time.Now().Add(-time.Hour).Unix()
			return f.sign(t, "k1", c)
		},
		"missing exp": func() string {
			c := f.claims()
			delete(c, "exp")
			return f.sign(t, "k1", c)
		},
		"missing sub": func() string {
			c := f.claims()
			delete(c, "sub")
			return f.sign(t, "k1", c)
		},
		"unknown kid": func() string {
			c := f.claims()
			tok := jwt.NewWithClaims(jwt.SigningMethodES384, c)
			tok.Header["kid"] = "nope"
			other, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
			s, _ := tok.SignedString(other)
			return s
		},
		"signed by a different key with a known kid": func() string {
			tok := jwt.NewWithClaims(jwt.SigningMethodES384, f.claims())
			tok.Header["kid"] = "k1"
			other, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
			s, _ := tok.SignedString(other)
			return s
		},
		"alg none": func() string {
			tok := jwt.NewWithClaims(jwt.SigningMethodNone, f.claims())
			tok.Header["kid"] = "k1"
			s, _ := tok.SignedString(jwt.UnsafeAllowNoneSignatureType)
			return s
		},
		"HS256 key confusion": func() string {
			tok := jwt.NewWithClaims(jwt.SigningMethodHS256, f.claims())
			tok.Header["kid"] = "k1"
			s, _ := tok.SignedString([]byte("guessable"))
			return s
		},
		"garbage": func() string { return "not.a.jwt" },
	}
	for name, mk := range cases {
		t.Run(name, func(t *testing.T) {
			if p, err := v.Verify(context.Background(), mk()); err == nil {
				t.Fatalf("expected rejection, got principal %+v", p)
			}
		})
	}
}

func TestKeyRotationRefetchesJWKS(t *testing.T) {
	f := newFakeLogto(t)
	v := newVerifier(t, f, Config{})
	v.jwks.minRefresh = 0

	if _, err := v.Verify(context.Background(), f.sign(t, "k1", f.claims())); err != nil {
		t.Fatal(err)
	}
	f.addKey(t, "k2")
	if _, err := v.Verify(context.Background(), f.sign(t, "k2", f.claims())); err != nil {
		t.Fatalf("token signed with rotated key: %v", err)
	}
	if got := f.fetches.Load(); got != 2 {
		t.Errorf("jwks fetches = %d, want 2", got)
	}
}

func TestUnknownKidRefetchIsRateLimited(t *testing.T) {
	f := newFakeLogto(t)
	v := newVerifier(t, f, Config{})

	for i := 0; i < 5; i++ {
		tok := jwt.NewWithClaims(jwt.SigningMethodES384, f.claims())
		tok.Header["kid"] = "attacker-kid"
		other, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		s, _ := tok.SignedString(other)
		_, _ = v.Verify(context.Background(), s)
	}
	if got := f.fetches.Load(); got != 1 {
		t.Errorf("jwks fetches = %d, want 1 (unknown kids must not hammer the issuer)", got)
	}
}

func TestTenantClaimAndDefault(t *testing.T) {
	f := newFakeLogto(t)

	v := newVerifier(t, f, Config{TenantClaim: "organization_id"})
	c := f.claims()
	c["organization_id"] = "org-9"
	p, err := v.Verify(context.Background(), f.sign(t, "k1", c))
	if err != nil || p.TenantID != "org-9" {
		t.Fatalf("custom tenant claim: tenant=%v err=%v", p, err)
	}

	v = newVerifier(t, f, Config{DefaultTenantID: "default-tenant"})
	c = f.claims()
	delete(c, "tenant_id")
	p, err = v.Verify(context.Background(), f.sign(t, "k1", c))
	if err != nil || p.TenantID != "default-tenant" {
		t.Fatalf("default tenant: tenant=%v err=%v", p, err)
	}
}

func TestRole(t *testing.T) {
	for scope, want := range map[string]string{
		"publishing:admin publishing:read": "admin",
		"publishing:approve":               "approver",
		"publishing:edit":                  "editor",
		"publishing:read":                  "viewer",
		"identity:read":                    "",
		"":                                 "",
	} {
		p := &Principal{Scopes: map[string]bool{}}
		for _, s := range strings.Fields(scope) {
			p.Scopes[s] = true
		}
		if got := p.Role(); got != want {
			t.Errorf("scope %q: Role() = %q, want %q", scope, got, want)
		}
	}
}

func TestNewVerifierFailsClosed(t *testing.T) {
	for _, cfg := range []Config{{}, {Endpoint: "https://auth.example"}, {Audience: audience}} {
		if _, err := NewVerifier(cfg, nil); !errors.Is(err, ErrNotConfigured) {
			t.Errorf("NewVerifier(%+v) err = %v, want ErrNotConfigured", cfg, err)
		}
	}
}

func TestConfigFromEnv(t *testing.T) {
	t.Setenv("LOGTO_ENDPOINT", "https://auth.example/")
	t.Setenv("LOGTO_AUDIENCE", audience)
	t.Setenv("TENANT_CLAIM", "")
	t.Setenv("DEFAULT_TENANT_ID", "t0")
	cfg := ConfigFromEnv()
	if cfg.Endpoint != "https://auth.example" || cfg.Audience != audience ||
		cfg.TenantClaim != "tenant_id" || cfg.DefaultTenantID != "t0" {
		t.Errorf("ConfigFromEnv() = %+v", cfg)
	}
}
