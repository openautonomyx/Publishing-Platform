package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type DeploymentPhase string

const (
	PhaseDev        DeploymentPhase = "dev"
	PhaseStaging    DeploymentPhase = "staging"
	PhaseProduction DeploymentPhase = "production"
)

type DeploymentStatus string

const (
	StatusPending     DeploymentStatus = "PENDING"
	StatusInProgress  DeploymentStatus = "IN_PROGRESS"
	StatusTesting     DeploymentStatus = "TESTING"
	StatusApprovalReq DeploymentStatus = "APPROVAL_REQUIRED"
	StatusHealthy     DeploymentStatus = "HEALTHY"
	StatusUnhealthy   DeploymentStatus = "UNHEALTHY"
	StatusRolledBack  DeploymentStatus = "ROLLED_BACK"
	StatusFailed      DeploymentStatus = "FAILED"
)

type Deployment struct {
	ID              string                 `json:"id"`
	Commit          string                 `json:"commit"`
	Branch          string                 `json:"branch"`
	Phase           DeploymentPhase        `json:"phase"`
	Status          DeploymentStatus       `json:"status"`
	TriggeredBy     string                 `json:"triggered_by"`
	ApprovedBy      *string                `json:"approved_by,omitempty"`
	StartedAt       time.Time              `json:"started_at"`
	CompletedAt     *time.Time             `json:"completed_at,omitempty"`
	RolledBackAt    *time.Time             `json:"rolled_back_at,omitempty"`
	Duration        *time.Duration         `json:"duration,omitempty"`
	ErrorRate       float64                `json:"error_rate"`
	LatencyP95      float64                `json:"latency_p95"`
	HealthCheckPass int                    `json:"health_check_pass"`
	HealthCheckFail int                    `json:"health_check_fail"`
	Metrics         map[string]interface{} `json:"metrics,omitempty"`
	CanRollback     bool                   `json:"can_rollback"`
}

// DeploymentService manages deployment orchestration
type DeploymentService struct {
	logger *zap.Logger
	db     *Database
}

func NewDeploymentService() *DeploymentService {
	return &DeploymentService{
		logger: logger,
		db:     db,
	}
}

// RegisterDeploymentRoutes registers all deployment endpoints
func (ds *DeploymentService) RegisterDeploymentRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/deployments", ds.ListDeploymentsHandler).Methods("GET")
	router.HandleFunc("/api/v1/deployments/{id}", ds.GetDeploymentHandler).Methods("GET")
	router.HandleFunc("/api/v1/deployments/{id}/approve", ds.ApproveDeploymentHandler).Methods("POST")
	router.HandleFunc("/api/v1/deployments/{id}/rollback", ds.RollbackDeploymentHandler).Methods("POST")
	router.HandleFunc("/api/v1/deployments/{id}/logs", ds.GetDeploymentLogsHandler).Methods("GET")
	router.HandleFunc("/api/v1/deployments/{id}/metrics", ds.GetDeploymentMetricsHandler).Methods("GET")
	router.HandleFunc("/api/v1/deployments/webhook/github", ds.GithubWebhookHandler).Methods("POST")
}

// ListDeploymentsHandler lists deployments by phase
func (ds *DeploymentService) ListDeploymentsHandler(w http.ResponseWriter, r *http.Request) {
	phase := r.URL.Query().Get("phase")
	if phase == "" {
		phase = "production"
	}

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	query := `
		SELECT id, commit, branch, phase, status, triggered_by, approved_by,
		       started_at, completed_at, rolled_back_at, error_rate, latency_p95,
		       health_check_pass, health_check_fail, can_rollback
		FROM deployments
		WHERE phase = $1
		ORDER BY started_at DESC
		LIMIT $2
	`

	rows, err := db.QueryContext(r.Context(), query, phase, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch deployments")
		return
	}
	defer rows.Close()

	var deployments []Deployment
	for rows.Next() {
		var d Deployment
		var completedAt, rolledBackAt *time.Time

		err := rows.Scan(
			&d.ID, &d.Commit, &d.Branch, &d.Phase, &d.Status, &d.TriggeredBy, &d.ApprovedBy,
			&d.StartedAt, &completedAt, &rolledBackAt, &d.ErrorRate, &d.LatencyP95,
			&d.HealthCheckPass, &d.HealthCheckFail, &d.CanRollback,
		)
		if err != nil {
			continue
		}

		d.CompletedAt = completedAt
		d.RolledBackAt = rolledBackAt

		if completedAt != nil {
			duration := completedAt.Sub(d.StartedAt)
			d.Duration = &duration
		}

		deployments = append(deployments, d)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"deployments": deployments,
		"count":       len(deployments),
	})
}

// GetDeploymentHandler gets a single deployment
func (ds *DeploymentService) GetDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	deploymentID := mux.Vars(r)["id"]

	query := `
		SELECT id, commit, branch, phase, status, triggered_by, approved_by,
		       started_at, completed_at, rolled_back_at, error_rate, latency_p95,
		       health_check_pass, health_check_fail, can_rollback
		FROM deployments
		WHERE id = $1
	`

	var d Deployment
	var completedAt, rolledBackAt *time.Time

	err := db.QueryRowContext(r.Context(), query, deploymentID).Scan(
		&d.ID, &d.Commit, &d.Branch, &d.Phase, &d.Status, &d.TriggeredBy, &d.ApprovedBy,
		&d.StartedAt, &completedAt, &rolledBackAt, &d.ErrorRate, &d.LatencyP95,
		&d.HealthCheckPass, &d.HealthCheckFail, &d.CanRollback,
	)

	if err != nil {
		respondError(w, http.StatusNotFound, "Deployment not found")
		return
	}

	d.CompletedAt = completedAt
	d.RolledBackAt = rolledBackAt

	if completedAt != nil {
		duration := completedAt.Sub(d.StartedAt)
		d.Duration = &duration
	}

	respondJSON(w, http.StatusOK, d)
}

// ApproveDeploymentHandler approves a deployment
func (ds *DeploymentService) ApproveDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	deploymentID := mux.Vars(r)["id"]
	userID := r.Context().Value("user_id").(string)

	ds.logger.Info("approving deployment",
		zap.String("deployment_id", deploymentID),
		zap.String("user_id", userID),
	)

	// Update approval
	query := `
		UPDATE deployments
		SET approved_by = $1, status = $2, updated_at = $3
		WHERE id = $4
	`

	_, err := db.ExecContext(r.Context(), query,
		userID, string(StatusApprovalReq), time.Now(), deploymentID,
	)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to approve deployment")
		return
	}

	// Start async blue-green deployment
	go ds.executeBlueGreenDeploy(context.Background(), deploymentID)

	respondJSON(w, http.StatusOK, map[string]string{
		"status": "approved",
		"message": "Deployment approved and will proceed to production",
	})
}

// RollbackDeploymentHandler rolls back a deployment
func (ds *DeploymentService) RollbackDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	deploymentID := mux.Vars(r)["id"]
	userID := r.Context().Value("user_id").(string)

	ds.logger.Info("rolling back deployment",
		zap.String("deployment_id", deploymentID),
		zap.String("user_id", userID),
	)

	// Check if rollback is allowed
	var canRollback bool
	query := `SELECT can_rollback FROM deployments WHERE id = $1`
	err := db.QueryRowContext(r.Context(), query, deploymentID).Scan(&canRollback)
	if err != nil || !canRollback {
		respondError(w, http.StatusBadRequest, "Rollback not available for this deployment")
		return
	}

	// Execute rollback
	err = ds.rollbackDeployment(r.Context(), deploymentID, userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Rollback failed")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status": "rolled_back",
		"message": "Deployment has been rolled back",
	})
}

// GetDeploymentLogsHandler retrieves deployment logs
func (ds *DeploymentService) GetDeploymentLogsHandler(w http.ResponseWriter, r *http.Request) {
	deploymentID := mux.Vars(r)["id"]

	// Stream logs (simplified - would integrate with GitHub Actions API)
	logs := []map[string]interface{}{
		{"timestamp": time.Now(), "level": "INFO", "message": "Deployment started"},
		{"timestamp": time.Now(), "level": "INFO", "message": "Building Docker image"},
		{"timestamp": time.Now(), "level": "INFO", "message": "Running tests"},
		{"timestamp": time.Now(), "level": "INFO", "message": "Deploying to blue"},
		{"timestamp": time.Now(), "level": "INFO", "message": "Health checks passed"},
		{"timestamp": time.Now(), "level": "INFO", "message": "Traffic switched to green"},
	}

	w.Header().Set("Content-Type", "application/json")
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"deployment_id": deploymentID,
		"logs": logs,
	})
}

// GetDeploymentMetricsHandler retrieves deployment metrics
func (ds *DeploymentService) GetDeploymentMetricsHandler(w http.ResponseWriter, r *http.Request) {
	deploymentID := mux.Vars(r)["id"]

	query := `
		SELECT error_rate, latency_p50, latency_p95, latency_p99,
		       cpu_usage, memory_usage, requests_per_sec
		FROM deployment_metrics
		WHERE deployment_id = $1
		ORDER BY timestamp DESC
		LIMIT 100
	`

	rows, err := db.QueryContext(r.Context(), query, deploymentID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch metrics")
		return
	}
	defer rows.Close()

	type Metric struct {
		ErrorRate     float64   `json:"error_rate"`
		LatencyP50    float64   `json:"latency_p50"`
		LatencyP95    float64   `json:"latency_p95"`
		LatencyP99    float64   `json:"latency_p99"`
		CPUUsage      float64   `json:"cpu_usage"`
		MemoryUsage   float64   `json:"memory_usage"`
		RequestsPerSec float64  `json:"requests_per_sec"`
	}

	var metrics []Metric
	for rows.Next() {
		var m Metric
		err := rows.Scan(&m.ErrorRate, &m.LatencyP50, &m.LatencyP95, &m.LatencyP99,
			&m.CPUUsage, &m.MemoryUsage, &m.RequestsPerSec)
		if err != nil {
			continue
		}
		metrics = append(metrics, m)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"deployment_id": deploymentID,
		"metrics": metrics,
	})
}

// GithubWebhookHandler handles GitHub webhook events
func (ds *DeploymentService) GithubWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// GitHub push event
	var payload struct {
		Ref    string `json:"ref"`
		Commits []struct {
			ID      string `json:"id"`
			Message string `json:"message"`
		} `json:"commits"`
		Pusher struct {
			Name string `json:"name"`
		} `json:"pusher"`
	}

	if err := parseJSON(r.Body, &payload); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	// Determine branch
	branch := payload.Ref
	if len(payload.Commits) == 0 {
		respondError(w, http.StatusBadRequest, "No commits")
		return
	}

	commit := payload.Commits[0].ID
	phase := determinePhasefromBranch(branch)

	// Create deployment record
	deploymentID := fmt.Sprintf("dep-%s-%d", commit[:7], time.Now().Unix())

	query := `
		INSERT INTO deployments (
			id, commit, branch, phase, status, triggered_by, started_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := db.ExecContext(r.Context(), query,
		deploymentID, commit, branch, string(phase), string(StatusPending),
		payload.Pusher.Name, time.Now(),
	)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create deployment")
		return
	}

	ds.logger.Info("deployment created from github webhook",
		zap.String("deployment_id", deploymentID),
		zap.String("commit", commit),
		zap.String("branch", branch),
	)

	respondJSON(w, http.StatusAccepted, map[string]string{
		"deployment_id": deploymentID,
		"status": "pending",
	})
}

// executeBlueGreenDeploy executes the blue-green deployment
func (ds *DeploymentService) executeBlueGreenDeploy(ctx context.Context, deploymentID string) {
	ds.logger.Info("executing blue-green deployment", zap.String("deployment_id", deploymentID))

	// Mark as in progress
	ds.updateDeploymentStatus(ctx, deploymentID, StatusInProgress)

	// Get deployment details
	var commit string
	query := `SELECT commit FROM deployments WHERE id = $1`
	err := db.QueryRowContext(ctx, query, deploymentID).Scan(&commit)
	if err != nil {
		ds.logger.Error("failed to get deployment", zap.Error(err))
		ds.updateDeploymentStatus(ctx, deploymentID, StatusFailed)
		return
	}

	// Deploy to green
	if err := ds.deployToGreen(ctx, commit); err != nil {
		ds.logger.Error("failed to deploy to green", zap.Error(err))
		ds.updateDeploymentStatus(ctx, deploymentID, StatusUnhealthy)
		return
	}

	// Run health checks
	healthy := ds.runHealthChecks(ctx, 3002, 30*time.Second)
	if !healthy {
		ds.logger.Error("green health checks failed")
		ds.rollbackDeployment(ctx, deploymentID, "system")
		return
	}

	// Switch traffic
	if err := ds.switchTraffic(ctx, 3001, 3002); err != nil {
		ds.logger.Error("failed to switch traffic", zap.Error(err))
		ds.rollbackDeployment(ctx, deploymentID, "system")
		return
	}

	// Monitor deployment
	if err := ds.monitorDeployment(ctx, deploymentID, 5*time.Minute); err != nil {
		ds.logger.Error("deployment metrics unhealthy", zap.Error(err))
		ds.rollbackDeployment(ctx, deploymentID, "system")
		return
	}

	// Mark as healthy
	now := time.Now()
	ds.updateDeploymentStatusWithTime(ctx, deploymentID, StatusHealthy, &now)
	ds.logger.Info("deployment completed successfully", zap.String("deployment_id", deploymentID))
}

// Helper methods

func (ds *DeploymentService) deployToGreen(ctx context.Context, commit string) error {
	ds.logger.Info("deploying to green", zap.String("commit", commit))
	// SSH to VPS and deploy
	// (implementation would SSH into production VPS and run deployment)
	return nil
}

func (ds *DeploymentService) runHealthChecks(ctx context.Context, port int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		url := fmt.Sprintf("http://localhost:%d/health", port)
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return true
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(3 * time.Second)
	}
	return false
}

func (ds *DeploymentService) switchTraffic(ctx context.Context, oldPort, newPort int) error {
	ds.logger.Info("switching traffic",
		zap.Int("from", oldPort),
		zap.Int("to", newPort),
	)
	// Update Nginx configuration and reload
	return nil
}

func (ds *DeploymentService) monitorDeployment(ctx context.Context, deploymentID string, duration time.Duration) error {
	deadline := time.Now().Add(duration)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Check metrics from monitoring system
		// If error rate > 1% or latency > 500ms, return error

		time.Sleep(10 * time.Second)
	}
	return nil
}

func (ds *DeploymentService) rollbackDeployment(ctx context.Context, deploymentID string, requestedBy string) error {
	ds.logger.Info("rolling back deployment",
		zap.String("deployment_id", deploymentID),
		zap.String("requested_by", requestedBy),
	)

	// Switch traffic back to blue
	if err := ds.switchTraffic(ctx, 3002, 3001); err != nil {
		return err
	}

	now := time.Now()
	query := `
		UPDATE deployments
		SET status = $1, rolled_back_at = $2, can_rollback = $3
		WHERE id = $4
	`
	_, err := db.ExecContext(ctx, query, string(StatusRolledBack), now, false, deploymentID)
	return err
}

func (ds *DeploymentService) updateDeploymentStatus(ctx context.Context, deploymentID string, status DeploymentStatus) {
	query := `UPDATE deployments SET status = $1, updated_at = $2 WHERE id = $3`
	db.ExecContext(ctx, query, string(status), time.Now(), deploymentID)
}

func (ds *DeploymentService) updateDeploymentStatusWithTime(ctx context.Context, deploymentID string, status DeploymentStatus, completedAt *time.Time) {
	query := `UPDATE deployments SET status = $1, completed_at = $2, updated_at = $3 WHERE id = $4`
	db.ExecContext(ctx, query, string(status), completedAt, time.Now(), deploymentID)
}

func determinePhasefromBranch(ref string) DeploymentPhase {
	if ref == "refs/heads/develop" {
		return PhaseDev
	} else if ref == "refs/heads/staging" {
		return PhaseStaging
	} else if ref == "refs/heads/main" {
		return PhaseProduction
	}
	return PhaseDev
}
