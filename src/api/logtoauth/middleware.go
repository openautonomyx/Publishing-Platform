package logtoauth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type principalKey struct{}

// FromContext returns the caller that Require authenticated.
func FromContext(ctx context.Context) (*Principal, bool) {
	p, ok := ctx.Value(principalKey{}).(*Principal)
	return p, ok
}

// ErrorWriter writes an authentication (401) or authorization (403) failure.
type ErrorWriter func(w http.ResponseWriter, r *http.Request, status int, msg string)

// Require rejects requests without a valid Logto bearer token, a tenant, and at least
// one publishing:* scope. Errors use a {"success":false,"error":"…"} body and never
// echo verification details back to the caller.
func Require(v *Verifier) func(http.Handler) http.Handler {
	return RequireWith(v, func(w http.ResponseWriter, _ *http.Request, status int, msg string) {
		writeError(w, status, msg)
	})
}

// RequireWith is Require with a caller-supplied error format.
func RequireWith(v *Verifier, writeErr ErrorWriter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			scheme, token, ok := strings.Cut(r.Header.Get("Authorization"), " ")
			token = strings.TrimSpace(token)
			if !ok || !strings.EqualFold(scheme, "Bearer") || token == "" {
				writeErr(w, r, http.StatusUnauthorized, "Missing or malformed bearer token")
				return
			}

			p, err := v.Verify(r.Context(), token)
			if err != nil {
				log.Printf("auth rejected: %v", err)
				writeErr(w, r, http.StatusUnauthorized, "Invalid token")
				return
			}
			if p.TenantID == "" {
				writeErr(w, r, http.StatusForbidden, "Token carries no tenant and DEFAULT_TENANT_ID is not set")
				return
			}
			if p.Role() == "" {
				writeErr(w, r, http.StatusForbidden, "Token has no publishing:* scope")
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), principalKey{}, p)))
		})
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	if status == http.StatusUnauthorized {
		w.Header().Set("WWW-Authenticate", "Bearer")
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}{false, msg})
}
