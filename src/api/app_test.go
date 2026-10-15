package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestClientIPIgnoresProxyHeadersUnlessTrusted(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "10.0.0.7:5555"
	r.Header.Set("X-Forwarded-For", "1.2.3.4, 10.0.0.1")
	r.Header.Set("X-Real-IP", "5.6.7.8")

	if got := clientIP(r, false); got != "10.0.0.7" {
		t.Errorf("untrusted: got %s, want the socket address", got)
	}
	if got := clientIP(r, true); got != "1.2.3.4" {
		t.Errorf("trusted: got %s, want first X-Forwarded-For entry", got)
	}
	r.Header.Del("X-Forwarded-For")
	if got := clientIP(r, true); got != "5.6.7.8" {
		t.Errorf("trusted without XFF: got %s, want X-Real-IP", got)
	}
}

func TestSlugify(t *testing.T) {
	for in, want := range map[string]string{
		"acme":            "acme",
		"Acme Corp!":      "acme-corp",
		"org_9":           "org-9",
		"ab":              "org-ab",
		"--Weird__ID--":   "weird-id",
		"logto-org-xYz12": "logto-org-xyz12",
	} {
		if got := slugify(in); got != want {
			t.Errorf("slugify(%q) = %q, want %q", in, got, want)
		}
		if !slugPattern.MatchString(slugify(in)) {
			t.Errorf("slugify(%q) = %q is not a valid slug", in, slugify(in))
		}
	}
}

func TestJSONB(t *testing.T) {
	var j JSONB
	if err := j.Scan(nil); err != nil || j != nil {
		t.Fatalf("Scan(nil) = %v, %v", j, err)
	}
	if b, _ := json.Marshal(struct{ M JSONB }{j}); string(b) != `{"M":null}` {
		t.Errorf("null marshal = %s", b)
	}
	if err := j.Scan([]byte(`{"a":1}`)); err != nil {
		t.Fatal(err)
	}
	if b, _ := json.Marshal(struct{ M JSONB }{j}); string(b) != `{"M":{"a":1}}` {
		t.Errorf("object marshal = %s", b)
	}
	if v, _ := JSONB(nil).Value(); v != nil {
		t.Errorf("nil Value = %v, want SQL NULL", v)
	}
	if err := j.Scan(42); err == nil {
		t.Error("Scan(int) should fail")
	}
}

func TestRateLimitThroughRouter(t *testing.T) {
	app := &App{limiter: NewRateLimiter(1, 2)}
	defer app.limiter.Stop()
	h := app.Router()

	codes := []int{}
	for i := 0; i < 3; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/health", nil)
		req.RemoteAddr = "192.0.2.10:1234"
		h.ServeHTTP(rec, req)
		codes = append(codes, rec.Code)
	}
	if codes[0] != 200 || codes[1] != 200 || codes[2] != 429 {
		t.Fatalf("status codes = %v, want [200 200 429]", codes)
	}

	// A spoofed X-Forwarded-For does not get a fresh bucket when proxies aren't trusted.
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	req.RemoteAddr = "192.0.2.10:1234"
	req.Header.Set("X-Forwarded-For", "203.0.113.99")
	h.ServeHTTP(rec, req)
	if rec.Code != 429 {
		t.Errorf("spoofed XFF got %d, want 429", rec.Code)
	}
}
