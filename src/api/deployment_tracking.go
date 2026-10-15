//go:build draft

// Excluded from the default build: not routed from main.go and depends on code
// that does not exist yet. Build with -tags draft to work on it (see README).

package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DeploymentPhase string

const (
	PhaseConnect         DeploymentPhase = "connect"
	PhaseDependencies    DeploymentPhase = "dependencies"
	PhaseCloneRepo       DeploymentPhase = "clone_repo"
	PhaseSecrets         DeploymentPhase = "secrets"
	PhasePostgres        DeploymentPhase = "postgres"
	PhaseSchema          DeploymentPhase = "schema"
	PhaseRedis           DeploymentPhase = "redis"
	PhaseLiferay         DeploymentPhase = "liferay"
	PhaseAPI             DeploymentPhase = "api"
	PhaseNginx           DeploymentPhase = "nginx"
	PhaseHealthCheck     DeploymentPhase = "health_check"
	PhaseSaveInfo        DeploymentPhase = "save_info"
)

type DeploymentStatus string

const (
	StatusPending    DeploymentStatus = "PENDING"
	StatusInProgress DeploymentStatus = "IN_PROGRESS"
	StatusSuccess    DeploymentStatus = "SUCCESS"
	StatusFailed     DeploymentStatus = "FAILED"
	StatusRolledBack DeploymentStatus = "ROLLED_BACK"
)

type DeploymentEvent struct {
	ID          uuid.UUID
	DeploymentID uuid.UUID
	Phase       DeploymentPhase
	Status      DeploymentStatus
	Message     string
	Timestamp   time.Time
	Duration    *time.Duration
	Error       *string
	Metadata    map[string]interface{}
}

type Deployment struct {
	ID            uuid.UUID
	VPSHost       string
	Branch        string
	Status        DeploymentStatus
	Progress      int
	TotalSteps    int
	StartedAt     time.Time
	CompletedAt   *time.Time
	FailedPhase   *DeploymentPhase
	ErrorMessage  *string
	Events        []DeploymentEvent
	CurrentPhase  DeploymentPhase
	mu            sync.RWMutex
}

type DeploymentTracker struct {
	db         *Database
	logger     *zap.Logger
	deployments map[uuid.UUID]*Deployment
	mu         sync.RWMutex
}

func NewDeploymentTracker() *DeploymentTracker {
	return &DeploymentTracker{
		db:          db,
		logger:      logger,
		deployments: make(map[uuid.UUID]*Deployment),
	}
}

// StartDeployment creates a new deployment tracking record
func (dt *DeploymentTracker) StartDeployment(ctx context.Context, vpsHost, branch string) (*Deployment, error) {
	deploymentID := uuid.New()

	phases := []DeploymentPhase{
		PhaseConnect,
		PhaseDependencies,
		PhaseCloneRepo,
		PhaseSecrets,
		PhasePostgres,
		PhaseSchema,
		PhaseRedis,
		PhaseLiferay,
		PhaseAPI,
		PhaseNginx,
		PhaseHealthCheck,
		PhaseSaveInfo,
	}

	deployment := &Deployment{
		ID:         deploymentID,
		VPSHost:    vpsHost,
		Branch:     branch,
		Status:     StatusPending,
		Progress:   0,
		TotalSteps: len(phases),
		StartedAt:  time.Now(),
		Events:     []DeploymentEvent{},
		CurrentPhase: phases[0],
	}

	// Store in database
	query := `
		INSERT INTO deployments (id, vps_host, branch, status, progress, total_steps, started_at, current_phase)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := dt.db.ExecContext(ctx, query,
		deployment.ID, vpsHost, branch, string(StatusPending),
		0, len(phases), time.Now(), string(phases[0]),
	)

	if err != nil {
		dt.logger.Error("failed to create deployment record", zap.Error(err))
		return nil, err
	}

	// Store in memory
	dt.mu.Lock()
	dt.deployments[deploymentID] = deployment
	dt.mu.Unlock()

	// Start workflow
	go dt.startDeploymentWorkflow(context.Background(), deployment)

	dt.logger.Info("deployment started",
		zap.String("deployment_id", deploymentID.String()),
		zap.String("vps_host", vpsHost),
		zap.String("branch", branch),
	)

	return deployment, nil
}

// RecordPhaseStart records the start of a deployment phase
func (dt *DeploymentTracker) RecordPhaseStart(ctx context.Context, deploymentID uuid.UUID, phase DeploymentPhase) error {
	dt.mu.Lock()
	deployment, exists := dt.deployments[deploymentID]
	dt.mu.Unlock()

	if !exists {
		return fmt.Errorf("deployment not found: %s", deploymentID)
	}

	deployment.mu.Lock()
	deployment.Status = StatusInProgress
	deployment.CurrentPhase = phase
	deployment.mu.Unlock()

	event := DeploymentEvent{
		ID:           uuid.New(),
		DeploymentID: deploymentID,
		Phase:        phase,
		Status:       StatusInProgress,
		Message:      fmt.Sprintf("Starting phase: %s", phase),
		Timestamp:    time.Now(),
		Metadata: map[string]interface{}{
			"step": phase,
		},
	}

	// Store event
	query := `
		INSERT INTO deployment_events (id, deployment_id, phase, status, message, timestamp, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := dt.db.ExecContext(ctx, query,
		event.ID, deploymentID, string(phase), string(StatusInProgress),
		event.Message, time.Now(), event.Metadata,
	)

	deployment.mu.Lock()
	deployment.Events = append(deployment.Events, event)
	deployment.mu.Unlock()

	dt.logger.Info("phase started",
		zap.String("deployment_id", deploymentID.String()),
		zap.String("phase", string(phase)),
	)

	return err
}

// RecordPhaseComplete records the successful completion of a phase
func (dt *DeploymentTracker) RecordPhaseComplete(ctx context.Context, deploymentID uuid.UUID, phase DeploymentPhase, duration time.Duration) error {
	dt.mu.Lock()
	deployment, exists := dt.deployments[deploymentID]
	dt.mu.Unlock()

	if !exists {
		return fmt.Errorf("deployment not found: %s", deploymentID)
	}

	deployment.mu.Lock()
	deployment.Progress++
	deployment.mu.Unlock()

	event := DeploymentEvent{
		ID:           uuid.New(),
		DeploymentID: deploymentID,
		Phase:        phase,
		Status:       StatusSuccess,
		Message:      fmt.Sprintf("Completed phase: %s", phase),
		Timestamp:    time.Now(),
		Duration:     &duration,
		Metadata: map[string]interface{}{
			"phase":    string(phase),
			"duration": duration.Seconds(),
		},
	}

	// Store event
	query := `
		INSERT INTO deployment_events (id, deployment_id, phase, status, message, timestamp, duration_ms, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := dt.db.ExecContext(ctx, query,
		event.ID, deploymentID, string(phase), string(StatusSuccess),
		event.Message, time.Now(), duration.Milliseconds(), event.Metadata,
	)

	deployment.mu.Lock()
	deployment.Events = append(deployment.Events, event)
	deployment.mu.Unlock()

	dt.logger.Info("phase completed",
		zap.String("deployment_id", deploymentID.String()),
		zap.String("phase", string(phase)),
		zap.Float64("duration_seconds", duration.Seconds()),
	)

	return err
}

// RecordPhaseFailed records a phase failure
func (dt *DeploymentTracker) RecordPhaseFailed(ctx context.Context, deploymentID uuid.UUID, phase DeploymentPhase, errorMsg string) error {
	dt.mu.Lock()
	deployment, exists := dt.deployments[deploymentID]
	dt.mu.Unlock()

	if !exists {
		return fmt.Errorf("deployment not found: %s", deploymentID)
	}

	deployment.mu.Lock()
	deployment.Status = StatusFailed
	deployment.FailedPhase = &phase
	deployment.ErrorMessage = &errorMsg
	deployment.mu.Unlock()

	event := DeploymentEvent{
		ID:           uuid.New(),
		DeploymentID: deploymentID,
		Phase:        phase,
		Status:       StatusFailed,
		Message:      fmt.Sprintf("Phase failed: %s", phase),
		Timestamp:    time.Now(),
		Error:        &errorMsg,
		Metadata: map[string]interface{}{
			"error": errorMsg,
			"phase": string(phase),
		},
	}

	// Store event
	query := `
		INSERT INTO deployment_events (id, deployment_id, phase, status, message, timestamp, error, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := dt.db.ExecContext(ctx, query,
		event.ID, deploymentID, string(phase), string(StatusFailed),
		event.Message, time.Now(), errorMsg, event.Metadata,
	)

	deployment.mu.Lock()
	deployment.Events = append(deployment.Events, event)
	deployment.mu.Unlock()

	// Send alert
	go dt.sendAlert(context.Background(), deployment, errorMsg)

	dt.logger.Error("phase failed",
		zap.String("deployment_id", deploymentID.String()),
		zap.String("phase", string(phase)),
		zap.String("error", errorMsg),
	)

	return err
}

// CompleteDeployment marks deployment as complete
func (dt *DeploymentTracker) CompleteDeployment(ctx context.Context, deploymentID uuid.UUID) error {
	dt.mu.Lock()
	deployment, exists := dt.deployments[deploymentID]
	dt.mu.Unlock()

	if !exists {
		return fmt.Errorf("deployment not found: %s", deploymentID)
	}

	now := time.Now()
	deployment.mu.Lock()
	deployment.Status = StatusSuccess
	deployment.CompletedAt = &now
	deployment.mu.Unlock()

	totalDuration := time.Since(deployment.StartedAt)

	// Update database
	query := `
		UPDATE deployments
		SET status = $1, completed_at = $2, progress = $3
		WHERE id = $4
	`
	_, err := dt.db.ExecContext(ctx, query,
		string(StatusSuccess), now, deployment.Progress, deploymentID,
	)

	event := DeploymentEvent{
		ID:           uuid.New(),
		DeploymentID: deploymentID,
		Phase:        PhaseSaveInfo,
		Status:       StatusSuccess,
		Message:      "Deployment completed successfully",
		Timestamp:    now,
		Duration:     &totalDuration,
		Metadata: map[string]interface{}{
			"total_duration_seconds": totalDuration.Seconds(),
			"phases_completed":       deployment.Progress,
			"total_phases":           deployment.TotalSteps,
		},
	}

	deployment.mu.Lock()
	deployment.Events = append(deployment.Events, event)
	deployment.mu.Unlock()

	// Send success notification
	go dt.sendNotification(context.Background(), deployment, "Deployment completed successfully")

	dt.logger.Info("deployment completed",
		zap.String("deployment_id", deploymentID.String()),
		zap.Float64("total_duration_seconds", totalDuration.Seconds()),
	)

	return err
}

// GetDeployment retrieves deployment status
func (dt *DeploymentTracker) GetDeployment(ctx context.Context, deploymentID uuid.UUID) (*Deployment, error) {
	dt.mu.RLock()
	deployment, exists := dt.deployments[deploymentID]
	dt.mu.RUnlock()

	if !exists {
		// Try to load from database
		query := `
			SELECT id, vps_host, branch, status, progress, total_steps, started_at, completed_at, failed_phase, error_message, current_phase
			FROM deployments
			WHERE id = $1
		`

		var d Deployment
		var failedPhase, errorMsg, currentPhase *string
		var completedAt *time.Time

		err := dt.db.QueryRowContext(ctx, query, deploymentID).Scan(
			&d.ID, &d.VPSHost, &d.Branch, &d.Status, &d.Progress, &d.TotalSteps,
			&d.StartedAt, &completedAt, &failedPhase, &errorMsg, &currentPhase,
		)

		if err != nil {
			return nil, fmt.Errorf("deployment not found: %s", deploymentID)
		}

		d.CompletedAt = completedAt
		if failedPhase != nil {
			phase := DeploymentPhase(*failedPhase)
			d.FailedPhase = &phase
		}
		d.ErrorMessage = errorMsg
		if currentPhase != nil {
			d.CurrentPhase = DeploymentPhase(*currentPhase)
		}

		// Load events
		eventQuery := `
			SELECT id, deployment_id, phase, status, message, timestamp, duration_ms, error, metadata
			FROM deployment_events
			WHERE deployment_id = $1
			ORDER BY timestamp ASC
		`

		rows, err := dt.db.QueryContext(ctx, eventQuery, deploymentID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var event DeploymentEvent
				var durationMs, error_msg *int64
				var metadata map[string]interface{}

				rows.Scan(
					&event.ID, &event.DeploymentID, &event.Phase, &event.Status,
					&event.Message, &event.Timestamp, &durationMs, &error_msg, &metadata,
				)

				if durationMs != nil {
					duration := time.Duration(*durationMs) * time.Millisecond
					event.Duration = &duration
				}
				if error_msg != nil {
					event.Error = &*error_msg
				}
				event.Metadata = metadata

				d.Events = append(d.Events, event)
			}
		}

		return &d, nil
	}

	return deployment, nil
}

// startDeploymentWorkflow runs the deployment workflow
func (dt *DeploymentTracker) startDeploymentWorkflow(ctx context.Context, deployment *Deployment) {
	phases := []DeploymentPhase{
		PhaseConnect,
		PhaseDependencies,
		PhaseCloneRepo,
		PhaseSecrets,
		PhasePostgres,
		PhaseSchema,
		PhaseRedis,
		PhaseLiferay,
		PhaseAPI,
		PhaseNginx,
		PhaseHealthCheck,
		PhaseSaveInfo,
	}

	for _, phase := range phases {
		phaseStart := time.Now()

		if err := dt.RecordPhaseStart(ctx, deployment.ID, phase); err != nil {
			dt.logger.Error("failed to record phase start", zap.Error(err))
			continue
		}

		// Simulate phase execution (in real scenario, this calls actual deployment steps)
		phaseErr := dt.executePhase(ctx, deployment, phase)

		if phaseErr != nil {
			dt.RecordPhaseFailed(ctx, deployment.ID, phase, phaseErr.Error())
			return
		}

		duration := time.Since(phaseStart)
		if err := dt.RecordPhaseComplete(ctx, deployment.ID, phase, duration); err != nil {
			dt.logger.Error("failed to record phase complete", zap.Error(err))
		}

		// Small delay between phases
		time.Sleep(100 * time.Millisecond)
	}

	// Mark deployment as complete
	if err := dt.CompleteDeployment(ctx, deployment.ID); err != nil {
		dt.logger.Error("failed to complete deployment", zap.Error(err))
	}
}

// executePhase simulates phase execution
func (dt *DeploymentTracker) executePhase(ctx context.Context, deployment *Deployment, phase DeploymentPhase) error {
	// In production, this would execute actual deployment commands
	// For now, just simulate success
	return nil
}

// sendAlert sends deployment failure alert
func (dt *DeploymentTracker) sendAlert(ctx context.Context, deployment *Deployment, errorMsg string) {
	// Send to Slack/Email/PagerDuty
	dt.logger.Warn("deployment alert",
		zap.String("deployment_id", deployment.ID.String()),
		zap.String("error", errorMsg),
	)
}

// sendNotification sends deployment completion notification
func (dt *DeploymentTracker) sendNotification(ctx context.Context, deployment *Deployment, message string) {
	// Send to Slack/Email
	dt.logger.Info("deployment notification",
		zap.String("deployment_id", deployment.ID.String()),
		zap.String("message", message),
	)
}

// GetDeploymentProgress returns real-time progress
func (dt *DeploymentTracker) GetDeploymentProgress(deploymentID uuid.UUID) map[string]interface{} {
	dt.mu.RLock()
	deployment, exists := dt.deployments[deploymentID]
	dt.mu.RUnlock()

	if !exists {
		return nil
	}

	deployment.mu.RLock()
	defer deployment.mu.RUnlock()

	progressPercent := 0
	if deployment.TotalSteps > 0 {
		progressPercent = (deployment.Progress * 100) / deployment.TotalSteps
	}

	return map[string]interface{}{
		"deployment_id":   deployment.ID.String(),
		"status":          string(deployment.Status),
		"progress":        deployment.Progress,
		"total_steps":     deployment.TotalSteps,
		"progress_percent": progressPercent,
		"current_phase":   string(deployment.CurrentPhase),
		"started_at":      deployment.StartedAt,
		"completed_at":    deployment.CompletedAt,
		"vps_host":        deployment.VPSHost,
		"branch":          deployment.Branch,
		"events_count":    len(deployment.Events),
	}
}

// ListDeployments lists all deployments
func (dt *DeploymentTracker) ListDeployments(ctx context.Context, limit int) ([]Deployment, error) {
	query := `
		SELECT id, vps_host, branch, status, progress, total_steps, started_at, completed_at
		FROM deployments
		ORDER BY started_at DESC
		LIMIT $1
	`

	rows, err := dt.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deployments []Deployment
	for rows.Next() {
		var d Deployment
		var completedAt *time.Time

		err := rows.Scan(
			&d.ID, &d.VPSHost, &d.Branch, &d.Status, &d.Progress,
			&d.TotalSteps, &d.StartedAt, &completedAt,
		)
		if err != nil {
			continue
		}

		d.CompletedAt = completedAt
		deployments = append(deployments, d)
	}

	return deployments, nil
}
