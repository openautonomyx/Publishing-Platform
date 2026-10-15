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
	"strconv"
	"strings"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/fractional-pm/creative-platform/logtoauth"
)

const Version = "0.2.0"

var startTime = time.Now()

func main() {
	config := loadConfig()

	db, err := initDB(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Fail closed: refuse to start without a Logto issuer to verify tokens against.
	verifier, err := logtoauth.NewVerifier(logtoauth.ConfigFromEnv(), nil)
	if err != nil {
		log.Fatalf("Failed to initialize auth: %v", err)
	}

	app := &App{db: db, verifier: verifier, corsOrigins: config.CORSOrigins, trustProxy: config.TrustProxy}
	if config.RateLimitRPS > 0 {
		app.limiter = NewRateLimiter(config.RateLimitRPS, config.RateLimitBurst)
		defer app.limiter.Stop()
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", config.Port),
		Handler:           app.Router(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("Starting API server on port %s (v%s)", config.Port, Version)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

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
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	log.Println("Database connected successfully")
	return db, nil
}

type Config struct {
	DatabaseURL    string
	Port           string
	CORSOrigins    map[string]bool
	TrustProxy     bool
	RateLimitRPS   int
	RateLimitBurst int
}

func loadConfig() Config {
	c := Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgresql://dev:dev@localhost:5432/creative_platform"),
		Port:           getEnv("API_PORT", "3001"),
		CORSOrigins:    map[string]bool{},
		TrustProxy:     getEnv("TRUST_PROXY", "false") == "true",
		RateLimitRPS:   getEnvInt("RATE_LIMIT_REQUESTS_PER_SECOND", 100),
		RateLimitBurst: getEnvInt("RATE_LIMIT_BURST_SIZE", 200),
	}
	for _, o := range strings.Split(getEnv("CORS_ORIGINS", ""), ",") {
		if o = strings.TrimSpace(o); o != "" && o != "*" {
			c.CORSOrigins[o] = true
		}
	}
	return c
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if v, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return v
	}
	return defaultVal
}

func encodeJSON(w http.ResponseWriter, v interface{}) {
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
