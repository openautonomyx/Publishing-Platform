package logtoauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequire(t *testing.T) {
	f := newFakeLogto(t)
	v := newVerifier(t, f, Config{})

	var seen *Principal
	h := Require(v)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen, _ = FromContext(r.Context())
		w.WriteHeader(http.StatusNoContent)
	}))

	noTenant := f.claims()
	delete(noTenant, "tenant_id")
	noScope := f.claims()
	noScope["scope"] = "identity:read"

	cases := []struct {
		name   string
		header string
		want   int
	}{
		{"valid", "Bearer " + f.sign(t, "k1", f.claims()), http.StatusNoContent},
		{"lowercase scheme", "bearer " + f.sign(t, "k1", f.claims()), http.StatusNoContent},
		{"no header", "", http.StatusUnauthorized},
		{"basic auth", "Basic dXNlcjpwYXNz", http.StatusUnauthorized},
		{"empty bearer", "Bearer ", http.StatusUnauthorized},
		{"any string (old stub accepted this)", "Bearer anything", http.StatusUnauthorized},
		{"no tenant", "Bearer " + f.sign(t, "k1", noTenant), http.StatusForbidden},
		{"no publishing scope", "Bearer " + f.sign(t, "k1", noScope), http.StatusForbidden},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			seen = nil
			req := httptest.NewRequest(http.MethodGet, "/api/v1/content", nil)
			if c.header != "" {
				req.Header.Set("Authorization", c.header)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != c.want {
				t.Fatalf("status = %d, want %d (body %s)", rec.Code, c.want, rec.Body)
			}
			if c.want != http.StatusNoContent {
				if seen != nil {
					t.Error("handler ran for a rejected request")
				}
				var body struct {
					Success bool   `json:"success"`
					Error   string `json:"error"`
				}
				if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.Success || body.Error == "" {
					t.Errorf("error body = %s", rec.Body)
				}
				return
			}
			if seen == nil || seen.Subject != "u_abc123" || seen.TenantID != "tenant-a" || seen.Role() != "editor" {
				t.Errorf("principal in context = %+v", seen)
			}
		})
	}
}
