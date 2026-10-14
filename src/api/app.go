package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lib/pq"

	"github.com/fractional-pm/creative-platform/logtoauth"
)

// App holds the API's dependencies. Handlers are methods on it, so tests can build
// one against a real database and a fake identity provider.
type App struct {
	db          *sql.DB
	verifier    *logtoauth.Verifier
	limiter     *RateLimiter    // nil disables rate limiting
	corsOrigins map[string]bool // empty allows any origin
	trustProxy  bool            // trust X-Forwarded-For / X-Real-IP for the client IP
}

// Identity is the authenticated caller, resolved to rows in this database.
type Identity struct {
	OrgID    uuid.UUID
	UserID   uuid.UUID
	Role     string // admin, approver, editor, viewer
	Via      string // "logto" or "api_key"
	APIKeyID *uuid.UUID
}

type ctxKey int

const (
	identityKey ctxKey = iota
	requestIDKey
)

func identityFrom(ctx context.Context) *Identity {
	id, _ := ctx.Value(identityKey).(*Identity)
	return id
}

func requestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}

// querier is satisfied by both *sql.DB and *sql.Tx.
type querier interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// ============================================================================
// Router
// ============================================================================

// Router returns the full HTTP handler. Cross-cutting middleware wraps the mux (not
// mux.Use) so it also runs for CORS preflights and unmatched routes.
func (a *App) Router() http.Handler {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.fail(w, r, NewNotFoundError("route"))
	})
	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.fail(w, r, &AppError{Code: ErrorCodeBadRequest, Message: "method not allowed", StatusCode: http.StatusMethodNotAllowed})
	})

	r.HandleFunc("/health", a.health).Methods("GET")
	r.HandleFunc("/ready", a.ready).Methods("GET")

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(a.authenticate)
	route := func(method, path string, min PermissionLevel, h http.HandlerFunc) {
		api.Handle(path, a.requireLevel(min, h)).Methods(method)
	}

	route("POST", "/organizations", PermissionAdmin, a.createOrganization)
	route("GET", "/organizations", PermissionView, a.listOrganizations)
	route("GET", "/organizations/{id}", PermissionView, a.getOrganization)
	route("PUT", "/organizations/{id}", PermissionAdmin, a.updateOrganization)
	route("DELETE", "/organizations/{id}", PermissionAdmin, a.deleteOrganization)

	route("POST", "/users", PermissionAdmin, a.createUser)
	route("GET", "/users", PermissionView, a.listUsers)
	route("GET", "/users/{id}", PermissionView, a.getUser)
	route("PUT", "/users/{id}", PermissionView, a.updateUser) // self-service; admin for others
	route("DELETE", "/users/{id}", PermissionAdmin, a.deleteUser)

	route("POST", "/content", PermissionEdit, a.createContent)
	route("GET", "/content", PermissionView, a.listContent)
	route("GET", "/content/{id}", PermissionView, a.getContent)
	route("PUT", "/content/{id}", PermissionEdit, a.updateContent)
	route("DELETE", "/content/{id}", PermissionEdit, a.deleteContent)

	route("POST", "/workflows", PermissionAdmin, a.createWorkflow)
	route("GET", "/workflows", PermissionView, a.listWorkflows)
	route("GET", "/workflows/{id}", PermissionView, a.getWorkflow)
	route("PUT", "/workflows/{id}", PermissionAdmin, a.updateWorkflow)
	route("DELETE", "/workflows/{id}", PermissionAdmin, a.deleteWorkflow)

	route("POST", "/approvals", PermissionApprove, a.createApproval)
	route("GET", "/content/{id}/approvals", PermissionView, a.listApprovals)

	route("POST", "/agents", PermissionEdit, a.createAgent)
	route("GET", "/agents", PermissionView, a.listAgents)
	route("GET", "/agents/{id}", PermissionView, a.getAgent)
	route("PUT", "/agents/{id}", PermissionEdit, a.updateAgent)
	route("DELETE", "/agents/{id}", PermissionEdit, a.deleteAgent)
	route("POST", "/agents/{id}/run", PermissionEdit, a.runAgent)
	route("GET", "/agents/{id}/runs", PermissionView, a.listAgentRuns)

	route("POST", "/distributions", PermissionEdit, a.createDistribution)
	route("GET", "/content/{id}/distributions", PermissionView, a.listDistributions)

	route("POST", "/api-keys", PermissionAdmin, a.createAPIKey)
	route("GET", "/api-keys", PermissionAdmin, a.listAPIKeys)
	route("DELETE", "/api-keys/{id}", PermissionAdmin, a.revokeAPIKey)

	route("GET", "/audit-logs", PermissionAdmin, a.listAuditLogs)

	var h http.Handler = r
	if a.limiter != nil {
		h = a.rateLimit(h)
	}
	return a.requestID(a.logRequests(a.cors(a.recoverPanics(h))))
}

func (a *App) health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, HealthCheckResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Version:   Version,
		Uptime:    time.Since(startTime).Milliseconds(),
	})
}

// ready reports whether the database is reachable (for readiness probes).
func (a *App) ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := a.db.PingContext(ctx); err != nil {
		logger.Warnw("readiness check failed", "error", err)
		a.fail(w, r, NewServiceUnavailableError("database unavailable"))
		return
	}
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Message: "ready"})
}

// ============================================================================
// Cross-cutting middleware
// ============================================================================

var requestIDPattern = regexp.MustCompile(`^[A-Za-z0-9._-]{1,128}$`)

func (a *App) requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if !requestIDPattern.MatchString(id) {
			id = uuid.NewString()
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDKey, id)))
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	if s.status == 0 {
		s.status = code
	}
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusRecorder) Write(b []byte) (int, error) {
	if s.status == 0 {
		s.status = http.StatusOK
	}
	return s.ResponseWriter.Write(b)
}

func (a *App) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(rec, r)
		logRequest(r.Method, r.URL.Path, rec.status, time.Since(start), requestIDFrom(r.Context()))
	})
}

func (a *App) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		switch {
		case len(a.corsOrigins) == 0:
			w.Header().Set("Access-Control-Allow-Origin", "*")
		case a.corsOrigins[origin]:
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key, X-Request-ID")
		w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *App) recoverPanics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				logPanic("panic in http handler", v)
				a.fail(w, r, fmt.Errorf("panic: %v", v))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (a *App) rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r, a.trustProxy)
		if !a.limiter.IsAllowed(ip) {
			logRateLimitExceeded(ip)
			w.Header().Set("Retry-After", "1")
			a.fail(w, r, NewTooManyRequestsError())
			return
		}
		next.ServeHTTP(w, r)
	})
}

// clientIP returns the caller's IP. Proxy headers are client-controlled, so they are
// only trusted when the API runs behind a proxy that sets them (TRUST_PROXY=true).
func clientIP(r *http.Request, trustProxy bool) string {
	if trustProxy {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			first := strings.TrimSpace(strings.Split(xff, ",")[0])
			if net.ParseIP(first) != nil {
				return first
			}
		}
		if real := strings.TrimSpace(r.Header.Get("X-Real-IP")); net.ParseIP(real) != nil {
			return real
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// ============================================================================
// Authentication: Logto bearer token or X-API-Key
// ============================================================================

func (a *App) authenticate(next http.Handler) http.Handler {
	bearer := logtoauth.RequireWith(a.verifier, func(w http.ResponseWriter, r *http.Request, status int, msg string) {
		code := ErrorCodeUnauthorized
		if status == http.StatusForbidden {
			code = ErrorCodeForbidden
		} else {
			w.Header().Set("WWW-Authenticate", "Bearer")
		}
		a.fail(w, r, &AppError{Code: code, Message: msg, StatusCode: status})
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, _ := logtoauth.FromContext(r.Context())
		id, err := a.resolveLogtoIdentity(r.Context(), p)
		if err != nil {
			a.fail(w, r, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), identityKey, id)))
	}))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if key := r.Header.Get("X-API-Key"); key != "" {
			id, err := a.resolveAPIKey(r.Context(), key)
			if err != nil {
				a.fail(w, r, err)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), identityKey, id)))
			return
		}
		bearer.ServeHTTP(w, r)
	})
}

func (a *App) requireLevel(min PermissionLevel, h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := identityFrom(r.Context())
		if id == nil || !canPerformAction(id.Role, min) {
			a.fail(w, r, NewForbiddenError(fmt.Sprintf("requires the %s role", levelName(min))))
			return
		}
		h(w, r)
	})
}

// resolveLogtoIdentity maps a verified token to this database's organization and user,
// creating them on first sign-in (just-in-time provisioning). An invited user (created
// with POST /users) is linked by email when the token carries an email claim.
func (a *App) resolveLogtoIdentity(ctx context.Context, p *logtoauth.Principal) (*Identity, error) {
	orgID, err := a.provisionOrganization(ctx, p.TenantID)
	if err != nil {
		return nil, err
	}
	role := p.Role()

	var (
		userID    uuid.UUID
		status    string
		curRole   string
		lastLogin sql.NullTime
		deletedAt sql.NullTime
	)
	err = a.db.QueryRowContext(ctx, `
		SELECT id, status, role, last_login_at, deleted_at FROM users
		WHERE org_id = $1 AND external_id = $2`, orgID, p.Subject,
	).Scan(&userID, &status, &curRole, &lastLogin, &deletedAt)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		userID, status, err = a.provisionUser(ctx, orgID, p, role)
		if err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
	case deletedAt.Valid:
		return nil, NewForbiddenError("user has been removed from this organization")
	case curRole != role || !lastLogin.Valid || time.Since(lastLogin.Time) > 5*time.Minute:
		if _, err := a.db.ExecContext(ctx,
			`UPDATE users SET role = $2, last_login_at = NOW() WHERE id = $1`, userID, role); err != nil {
			return nil, err
		}
	}
	if status != "active" {
		return nil, NewForbiddenError("user is " + status)
	}
	return &Identity{OrgID: orgID, UserID: userID, Role: role, Via: "logto"}, nil
}

func (a *App) provisionOrganization(ctx context.Context, tenant string) (uuid.UUID, error) {
	var (
		id        uuid.UUID
		status    string
		deletedAt sql.NullTime
	)
	lookup := func() error {
		return a.db.QueryRowContext(ctx,
			`SELECT id, status, deleted_at FROM organizations WHERE external_id = $1`, tenant,
		).Scan(&id, &status, &deletedAt)
	}
	err := lookup()
	if errors.Is(err, sql.ErrNoRows) {
		slug := slugify(tenant)
		for attempt := 0; attempt < 2; attempt++ {
			_, err = a.db.ExecContext(ctx, `
				INSERT INTO organizations (name, slug, external_id) VALUES ($1, $2, $1)
				ON CONFLICT (external_id) DO NOTHING`, tenant, slug)
			if !isUniqueViolation(err, "organizations_slug_key") {
				break
			}
			sum := sha256.Sum256([]byte(tenant))
			slug = slug + "-" + hex.EncodeToString(sum[:3])
		}
		if err != nil {
			return uuid.Nil, err
		}
		err = lookup()
	}
	if err != nil {
		return uuid.Nil, err
	}
	if deletedAt.Valid {
		return uuid.Nil, NewForbiddenError("organization has been deleted")
	}
	if status != "active" {
		return uuid.Nil, NewForbiddenError("organization is " + status)
	}
	return id, nil
}

func (a *App) provisionUser(ctx context.Context, orgID uuid.UUID, p *logtoauth.Principal, role string) (uuid.UUID, string, error) {
	var (
		id     uuid.UUID
		status string
	)
	if p.Email != "" {
		err := a.db.QueryRowContext(ctx, `
			UPDATE users SET external_id = $3, role = $4, last_login_at = NOW(),
				status = CASE WHEN status = 'invited' THEN 'active' ELSE status END
			WHERE org_id = $1 AND lower(email) = lower($2) AND external_id IS NULL AND deleted_at IS NULL
			RETURNING id, status`, orgID, p.Email, p.Subject, role,
		).Scan(&id, &status)
		if err == nil {
			return id, status, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, "", err
		}
	}

	name := firstNonEmpty(p.Name, p.Email, p.Subject)
	email := sql.NullString{String: p.Email, Valid: p.Email != ""}
	insert := func(email sql.NullString) error {
		return a.db.QueryRowContext(ctx, `
			INSERT INTO users (org_id, external_id, email, name, role, status, last_login_at)
			VALUES ($1, $2, $3, $4, $5, 'active', NOW())
			ON CONFLICT (org_id, external_id) DO UPDATE SET last_login_at = EXCLUDED.last_login_at
			RETURNING id, status`, orgID, p.Subject, email, name, role,
		).Scan(&id, &status)
	}
	err := insert(email)
	if isUniqueViolation(err, "users_org_id_email_key") {
		// The email belongs to another user in this org; provision without it.
		err = insert(sql.NullString{})
	}
	return id, status, err
}

// resolveAPIKey authenticates an X-API-Key. The key's role is capped by its owner's role.
func (a *App) resolveAPIKey(ctx context.Context, key string) (*Identity, error) {
	var (
		keyID, orgID, userID uuid.UUID
		perms                pq.StringArray
		userRole, userStatus string
		orgStatus            string
		userDeleted, orgDel  sql.NullTime
	)
	err := a.db.QueryRowContext(ctx, `
		SELECT k.id, k.org_id, k.user_id, k.permissions,
		       u.role, u.status, u.deleted_at, o.status, o.deleted_at
		FROM api_keys k
		JOIN users u ON u.id = k.user_id
		JOIN organizations o ON o.id = k.org_id
		WHERE k.key_hash = $1 AND k.status = 'active' AND k.deleted_at IS NULL
		  AND (k.expires_at IS NULL OR k.expires_at > NOW())`, hashAPIKey(key),
	).Scan(&keyID, &orgID, &userID, &perms, &userRole, &userStatus, &userDeleted, &orgStatus, &orgDel)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, NewUnauthorizedError("invalid API key")
	}
	if err != nil {
		return nil, err
	}
	if orgDel.Valid || orgStatus != "active" {
		return nil, NewForbiddenError("organization is not active")
	}
	if userDeleted.Valid || userStatus != "active" {
		return nil, NewForbiddenError("the key's owner is not active")
	}

	role := roleForPermissions(perms)
	if getRolePermissionLevel(userRole) < getRolePermissionLevel(role) {
		role = userRole
	}
	if _, err := a.db.ExecContext(ctx, `
		UPDATE api_keys SET last_used_at = NOW()
		WHERE id = $1 AND (last_used_at IS NULL OR last_used_at < NOW() - INTERVAL '1 minute')`, keyID); err != nil {
		return nil, err
	}
	return &Identity{OrgID: orgID, UserID: userID, Role: role, Via: "api_key", APIKeyID: &keyID}, nil
}

func hashAPIKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}

func roleForPermissions(perms []string) string {
	role := "viewer"
	for _, p := range perms {
		var r string
		switch p {
		case "admin":
			r = "admin"
		case "approve":
			r = "approver"
		case "write":
			r = "editor"
		default:
			r = "viewer"
		}
		if getRolePermissionLevel(r) > getRolePermissionLevel(role) {
			role = r
		}
	}
	return role
}

var slugInvalid = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	slug := strings.Trim(slugInvalid.ReplaceAllString(strings.ToLower(s), "-"), "-")
	if len(slug) > 200 {
		slug = slug[:200]
	}
	if len(slug) < 3 {
		slug = "org-" + slug
	}
	return strings.Trim(slug, "-")
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// ============================================================================
// Request helpers
// ============================================================================

const maxBodyBytes = 1 << 20

// decode reads a JSON body into v (unknown fields rejected) and validates it.
func decode(r *http.Request, v interface{}) error {
	dec := json.NewDecoder(io.LimitReader(r.Body, maxBodyBytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		if errors.Is(err, io.EOF) {
			return NewBadRequestError("request body is required")
		}
		return NewBadRequestError("invalid JSON: " + err.Error())
	}
	return validateRequest(v)
}

func pathID(r *http.Request) (uuid.UUID, error) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		return uuid.Nil, NewBadRequestError("invalid id")
	}
	return id, nil
}

type page struct{ Page, Limit, Offset int }

func pagination(r *http.Request) page {
	p := page{Page: 1, Limit: 20}
	if v, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && v > 0 {
		p.Page = v
	}
	if v, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && v > 0 {
		p.Limit = v
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	p.Offset = (p.Page - 1) * p.Limit
	return p
}

func respondPage(w http.ResponseWriter, data interface{}, total int, p page) {
	respondJSON(w, http.StatusOK, PaginatedResponse{Success: true, Data: data, Total: total, Page: p.Page, Limit: p.Limit})
}

func respondData(w http.ResponseWriter, status int, data interface{}) {
	respondJSON(w, status, APIResponse{Success: true, Data: data})
}

// queryFilter builds "AND col = $n" clauses for optional query-string filters.
type queryFilter struct {
	where []string
	args  []interface{}
}

func (f *queryFilter) add(clause string, arg interface{}) {
	f.args = append(f.args, arg)
	f.where = append(f.where, fmt.Sprintf(clause, len(f.args)))
}

func (f *queryFilter) sql() string {
	if len(f.where) == 0 {
		return ""
	}
	return " AND " + strings.Join(f.where, " AND ")
}

// ============================================================================
// Transactions and audit
// ============================================================================

func (a *App) inTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

// audit records a mutation in audit_logs, in the same transaction as the change.
func (a *App) audit(ctx context.Context, q querier, r *http.Request, action, resourceType string, resourceID uuid.UUID, changes interface{}) error {
	id := identityFrom(ctx)
	if changes == nil {
		changes = map[string]interface{}{}
	}
	body, err := json.Marshal(changes)
	if err != nil {
		return err
	}
	ip := sql.NullString{}
	if v := clientIP(r, a.trustProxy); net.ParseIP(v) != nil {
		ip = sql.NullString{String: v, Valid: true}
	}
	ua := r.UserAgent()
	if len(ua) > 500 {
		ua = ua[:500]
	}
	_, err = q.ExecContext(ctx, `
		INSERT INTO audit_logs (org_id, user_id, action, resource_type, resource_id, changes, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		id.OrgID, id.UserID, action, resourceType, resourceID, body, ip, ua)
	return err
}

// ============================================================================
// Errors
// ============================================================================

// fail writes err as the API's standard error response. Database errors are mapped
// to 4xx where the cause is the request; anything else is logged and returned as a
// generic 500 without internal details.
func (a *App) fail(w http.ResponseWriter, r *http.Request, err error) {
	reqID := requestIDFrom(r.Context())
	var appErr *AppError
	var pqErr *pq.Error
	switch {
	case errors.As(err, &appErr):
	case errors.Is(err, sql.ErrNoRows):
		appErr = NewNotFoundError("resource")
	case errors.As(err, &pqErr) && pqErr.Code == "23505":
		appErr = NewConflictError("a resource with the same unique value already exists")
	case errors.As(err, &pqErr) && pqErr.Code == "23503":
		appErr = NewBadRequestError("a referenced resource does not exist")
	case errors.As(err, &pqErr) && (pqErr.Code == "22P02" || pqErr.Code == "22001"):
		appErr = NewBadRequestError("invalid value")
	default:
		logger.Errorw("request failed", "error", err, "request_id", reqID, "path", r.URL.Path)
		appErr = NewInternalError()
	}
	writeError(w, appErr, reqID)
}

func isUniqueViolation(err error, constraint string) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505" && (constraint == "" || pqErr.Constraint == constraint)
}
