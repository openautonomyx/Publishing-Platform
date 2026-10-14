package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/fractional-pm/creative-platform/logtoauth"
)

const Version = "0.1.0"

var (
	startTime = time.Now()
	db        *sql.DB
)

func main() {
	// Load configuration from environment
	config := loadConfig()

	// Initialize database
	var err error
	db, err = initDB(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Fail closed: refuse to start without a Logto issuer to verify tokens against.
	verifier, err := logtoauth.NewVerifier(logtoauth.ConfigFromEnv(), nil)
	if err != nil {
		log.Fatalf("Failed to initialize auth: %v", err)
	}

	// Initialize router
	router := mux.NewRouter()

	// Global middleware
	router.Use(loggingMiddleware)
	router.Use(corsMiddleware)
	router.Use(requestIDMiddleware)
	router.Use(recoveryMiddleware)

	// Health check endpoint (no auth required)
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")

	// API v1 routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(authMiddleware(verifier))

	// Organizations
	api.HandleFunc("/organizations", createOrganizationHandler).Methods("POST")
	api.HandleFunc("/organizations", listOrganizationsHandler).Methods("GET")
	api.HandleFunc("/organizations/{id}", getOrganizationHandler).Methods("GET")
	api.HandleFunc("/organizations/{id}", updateOrganizationHandler).Methods("PUT")
	api.HandleFunc("/organizations/{id}", deleteOrganizationHandler).Methods("DELETE")

	// Users
	api.HandleFunc("/users", createUserHandler).Methods("POST")
	api.HandleFunc("/users", listUsersHandler).Methods("GET")
	api.HandleFunc("/users/{id}", getUserHandler).Methods("GET")
	api.HandleFunc("/users/{id}", updateUserHandler).Methods("PUT")
	api.HandleFunc("/users/{id}", deleteUserHandler).Methods("DELETE")

	// Content
	api.HandleFunc("/content", createContentHandler).Methods("POST")
	api.HandleFunc("/content", listContentHandler).Methods("GET")
	api.HandleFunc("/content/{id}", getContentHandler).Methods("GET")
	api.HandleFunc("/content/{id}", updateContentHandler).Methods("PUT")
	api.HandleFunc("/content/{id}", deleteContentHandler).Methods("DELETE")

	// Workflows
	api.HandleFunc("/workflows", createWorkflowHandler).Methods("POST")
	api.HandleFunc("/workflows", listWorkflowsHandler).Methods("GET")
	api.HandleFunc("/workflows/{id}", getWorkflowHandler).Methods("GET")
	api.HandleFunc("/workflows/{id}", updateWorkflowHandler).Methods("PUT")
	api.HandleFunc("/workflows/{id}", deleteWorkflowHandler).Methods("DELETE")

	// Approvals
	api.HandleFunc("/approvals", createApprovalHandler).Methods("POST")
	api.HandleFunc("/content/{id}/approvals", listApprovalsHandler).Methods("GET")

	// Agents
	api.HandleFunc("/agents", createAgentHandler).Methods("POST")
	api.HandleFunc("/agents", listAgentsHandler).Methods("GET")
	api.HandleFunc("/agents/{id}", getAgentHandler).Methods("GET")
	api.HandleFunc("/agents/{id}", updateAgentHandler).Methods("PUT")
	api.HandleFunc("/agents/{id}", deleteAgentHandler).Methods("DELETE")
	api.HandleFunc("/agents/{id}/run", runAgentHandler).Methods("POST")

	// Distributions
	api.HandleFunc("/distributions", createDistributionHandler).Methods("POST")
	api.HandleFunc("/content/{id}/distributions", listDistributionsHandler).Methods("GET")

	// API Keys
	api.HandleFunc("/api-keys", createAPIKeyHandler).Methods("POST")
	api.HandleFunc("/api-keys", listAPIKeysHandler).Methods("GET")
	api.HandleFunc("/api-keys/{id}", revokeAPIKeyHandler).Methods("DELETE")

	// Audit logs
	api.HandleFunc("/audit-logs", listAuditLogsHandler).Methods("GET")

	// HTTP Server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting API server on port %s (v%s)", config.Port, Version)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

func initDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Set connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")
	return db, nil
}

func loadConfig() Config {
	return Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgresql://dev:dev@localhost:5432/creative_platform"),
		Port:        getEnv("API_PORT", "3001"),
		Env:         getEnv("API_ENV", "development"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

type Config struct {
	DatabaseURL string
	Port        string
	Env         string
}

// ============================================================================
// Handlers (stubs - implemented in handlers.go)
// ============================================================================

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	uptime := int64(time.Since(startTime).Milliseconds())
	response := HealthCheckResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   Version,
		Uptime:    uptime,
	}

	encodeJSON(w, response)
}

// Organization handlers
func createOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Create Organization", http.StatusCreated)
}

func listOrganizationsHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "List Organizations", http.StatusOK)
}

func getOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Get Organization", http.StatusOK)
}

func updateOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Update Organization", http.StatusOK)
}

func deleteOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Delete Organization", http.StatusOK)
}

// User handlers
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Create User", http.StatusCreated)
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "List Users", http.StatusOK)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Get User", http.StatusOK)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Update User", http.StatusOK)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Delete User", http.StatusOK)
}

// Content handlers
func createContentHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Create Content", http.StatusCreated)
}

func listContentHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "List Content", http.StatusOK)
}

func getContentHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Get Content", http.StatusOK)
}

func updateContentHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Update Content", http.StatusOK)
}

func deleteContentHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Delete Content", http.StatusOK)
}

// Workflow handlers
func createWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Create Workflow", http.StatusCreated)
}

func listWorkflowsHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "List Workflows", http.StatusOK)
}

func getWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Get Workflow", http.StatusOK)
}

func updateWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Update Workflow", http.StatusOK)
}

func deleteWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Delete Workflow", http.StatusOK)
}

// Approval handlers
func createApprovalHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Create Approval", http.StatusCreated)
}

func listApprovalsHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "List Approvals", http.StatusOK)
}

// Agent handlers
func createAgentHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Create Agent", http.StatusCreated)
}

func listAgentsHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "List Agents", http.StatusOK)
}

func getAgentHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Get Agent", http.StatusOK)
}

func updateAgentHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Update Agent", http.StatusOK)
}

func deleteAgentHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Delete Agent", http.StatusOK)
}

func runAgentHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Run Agent", http.StatusOK)
}

// Distribution handlers
func createDistributionHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Create Distribution", http.StatusCreated)
}

func listDistributionsHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "List Distributions", http.StatusOK)
}

// API Key handlers
func createAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Create API Key", http.StatusCreated)
}

func listAPIKeysHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "List API Keys", http.StatusOK)
}

func revokeAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "Revoke API Key", http.StatusOK)
}

// Audit log handler
func listAuditLogsHandler(w http.ResponseWriter, r *http.Request) {
	handlerStub(w, "List Audit Logs", http.StatusOK)
}

func handlerStub(w http.ResponseWriter, action string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := APIResponse{
		Success: true,
		Message: action + " - Implementation in progress",
	}
	encodeJSON(w, response)
}

func encodeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("JSON encoding error: %v", err)
	}
}

// respondJSON writes v as the JSON response body with the given status.
func respondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encodeJSON(w, v)
}

// respondError writes the API's standard error envelope.
func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, APIResponse{Success: false, Error: msg})
}
