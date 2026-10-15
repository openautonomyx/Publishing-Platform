// Package logtoauth verifies Logto access tokens (JWT) against the issuer's JWKS.
//
// It follows the same contract as the Identity Fabric service: the same environment
// variables (LOGTO_ENDPOINT, LOGTO_AUDIENCE, TENANT_CLAIM, DEFAULT_TENANT_ID), the same
// accepted algorithms, and the same required claims, so one Logto API resource and one
// token can be used against both services.
package logtoauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Scopes (Logto API resource permissions) that grant publishing roles, highest first.
const (
	ScopeAdmin   = "publishing:admin"
	ScopeApprove = "publishing:approve"
	ScopeEdit    = "publishing:edit"
	ScopeRead    = "publishing:read"
)

// ErrNotConfigured is returned when LOGTO_ENDPOINT or LOGTO_AUDIENCE is unset.
var ErrNotConfigured = errors.New("logto is not configured (LOGTO_ENDPOINT / LOGTO_AUDIENCE)")

type Config struct {
	Endpoint        string // e.g. https://auth.example.com (no trailing slash)
	Audience        string // the Logto API resource indicator
	TenantClaim     string // claim that carries the tenant; default "tenant_id"
	DefaultTenantID string // used when the token carries no tenant claim
}

// ConfigFromEnv reads the same variables as Identity Fabric.
func ConfigFromEnv() Config {
	claim := os.Getenv("TENANT_CLAIM")
	if claim == "" {
		claim = "tenant_id"
	}
	return Config{
		Endpoint:        strings.TrimRight(os.Getenv("LOGTO_ENDPOINT"), "/"),
		Audience:        os.Getenv("LOGTO_AUDIENCE"),
		TenantClaim:     claim,
		DefaultTenantID: os.Getenv("DEFAULT_TENANT_ID"),
	}
}

// Principal is the authenticated caller.
type Principal struct {
	Subject  string // Logto user ID (or M2M app ID for client-credentials tokens)
	ClientID string
	TenantID string
	Email    string // optional "email" claim (a Logto custom JWT claim)
	Name     string // optional "name" claim
	Scopes   map[string]bool
}

// Role maps the caller's scopes onto the platform's roles: admin, approver, editor, viewer.
// A token with none of the publishing scopes gets "" and should be refused by the caller.
func (p *Principal) Role() string {
	switch {
	case p.Scopes[ScopeAdmin]:
		return "admin"
	case p.Scopes[ScopeApprove]:
		return "approver"
	case p.Scopes[ScopeEdit]:
		return "editor"
	case p.Scopes[ScopeRead]:
		return "viewer"
	}
	return ""
}

type Verifier struct {
	cfg    Config
	jwks   *jwksCache
	parser *jwt.Parser
}

// NewVerifier returns ErrNotConfigured when the endpoint or audience is missing, so the
// API fails closed instead of accepting unverified tokens.
func NewVerifier(cfg Config, client *http.Client) (*Verifier, error) {
	if cfg.Endpoint == "" || cfg.Audience == "" {
		return nil, ErrNotConfigured
	}
	if cfg.TenantClaim == "" {
		cfg.TenantClaim = "tenant_id"
	}
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	return &Verifier{
		cfg:  cfg,
		jwks: newJWKSCache(cfg.Endpoint+"/oidc/jwks", client),
		parser: jwt.NewParser(
			jwt.WithValidMethods([]string{"ES384", "ES256", "RS256"}),
			jwt.WithIssuer(cfg.Endpoint+"/oidc"),
			jwt.WithAudience(cfg.Audience),
			jwt.WithExpirationRequired(),
			jwt.WithLeeway(30*time.Second),
		),
	}, nil
}

// Verify checks signature, issuer, audience and expiry, and returns the caller.
func (v *Verifier) Verify(ctx context.Context, token string) (*Principal, error) {
	claims := jwt.MapClaims{}
	_, err := v.parser.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		kid, _ := t.Header["kid"].(string)
		return v.jwks.key(ctx, kid)
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	sub, _ := claims["sub"].(string)
	if sub == "" {
		return nil, errors.New("invalid token: missing sub")
	}
	p := &Principal{Subject: sub, Scopes: map[string]bool{}}
	p.ClientID, _ = claims["client_id"].(string)
	p.Email, _ = claims["email"].(string)
	p.Name, _ = claims["name"].(string)
	if scope, ok := claims["scope"].(string); ok {
		for _, s := range strings.Fields(scope) {
			p.Scopes[s] = true
		}
	}
	p.TenantID, _ = claims[v.cfg.TenantClaim].(string)
	if p.TenantID == "" {
		p.TenantID = v.cfg.DefaultTenantID
	}
	return p, nil
}
