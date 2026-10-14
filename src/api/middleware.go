package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/fractional-pm/creative-platform/logtoauth"
)

// ============================================================================
// Middleware Chain
// ============================================================================

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("Response time: %v", time.Since(start))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		w.Header().Set("X-Request-ID", requestID)
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				response := APIResponse{
					Success: false,
					Error:   "Internal server error",
				}
				encodeJSON(w, response)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// authMiddleware verifies Logto access tokens (see logtoauth.Require) and copies the
// caller's tenant, user and role into the request context for handlers and RLS.
func authMiddleware(verifier *logtoauth.Verifier) mux.MiddlewareFunc {
	require := logtoauth.Require(verifier)
	return func(next http.Handler) http.Handler {
		return require(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal, _ := logtoauth.FromContext(r.Context())

			ctx := context.WithValue(r.Context(), ContextKeyOrgID, principal.TenantID)
			ctx = context.WithValue(ctx, ContextKeyUserID, principal.Subject)
			ctx = context.WithValue(ctx, ContextKeyUserRole, principal.Role())

			// Set org_id in PostgreSQL session for RLS
			if err := setPostgreSQLOrgID(ctx, principal.TenantID); err != nil {
				log.Printf("Failed to set PostgreSQL org_id: %v", err)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}))
	}
}

// ============================================================================
// Helpers
// ============================================================================

func setPostgreSQLOrgID(ctx context.Context, orgID string) error {
	// PLACEHOLDER: Week 2 will implement proper session variable setting
	// This enables row-level security in PostgreSQL
	return nil
}

// ============================================================================
// Request/Response Interceptors
// ============================================================================

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body = append(rw.body, b...)
	return rw.ResponseWriter.Write(b)
}
