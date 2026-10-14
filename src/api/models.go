package main

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ============================================================================
// Organization (Tenant)
// ============================================================================

type Organization struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	Name        string                 `json:"name" db:"name"`
	Slug        string                 `json:"slug" db:"slug"`
	Description string                 `json:"description" db:"description"`
	Tier        string                 `json:"tier" db:"tier"` // free, pro, enterprise
	Status      string                 `json:"status" db:"status"`
	Metadata    json.RawMessage        `json:"metadata" db:"metadata"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
	CreatedBy   uuid.UUID              `json:"created_by" db:"created_by"`
	DeletedAt   *time.Time             `json:"deleted_at" db:"deleted_at"`
}

// ============================================================================
// User (Multi-tenant)
// ============================================================================

type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	OrgID        uuid.UUID  `json:"org_id" db:"org_id"`
	Email        string     `json:"email" db:"email"`
	Name         string     `json:"name" db:"name"`
	PasswordHash string     `json:"-" db:"password_hash"`
	AvatarURL    string     `json:"avatar_url" db:"avatar_url"`
	Role         string     `json:"role" db:"role"` // admin, editor, viewer, user
	Status       string     `json:"status" db:"status"`
	LastLoginAt  *time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" db:"deleted_at"`
}

// ============================================================================
// API Key
// ============================================================================

type APIKey struct {
	ID          uuid.UUID     `json:"id" db:"id"`
	OrgID       uuid.UUID     `json:"org_id" db:"org_id"`
	UserID      uuid.UUID     `json:"user_id" db:"user_id"`
	KeyHash     string        `json:"-" db:"key_hash"`
	Name        string        `json:"name" db:"name"`
	Permissions pq.StringArray `json:"permissions" db:"permissions"`
	Status      string        `json:"status" db:"status"`
	LastUsedAt  *time.Time    `json:"last_used_at" db:"last_used_at"`
	ExpiresAt   *time.Time    `json:"expires_at" db:"expires_at"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	DeletedAt   *time.Time    `json:"deleted_at" db:"deleted_at"`
}

// ============================================================================
// Workflow
// ============================================================================

type Workflow struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	OrgID        uuid.UUID       `json:"org_id" db:"org_id"`
	Name         string          `json:"name" db:"name"`
	Description  string          `json:"description" db:"description"`
	WorkflowType string          `json:"workflow_type" db:"workflow_type"`
	Status       string          `json:"status" db:"status"`
	Config       json.RawMessage `json:"config" db:"config"`
	CreatedBy    uuid.UUID       `json:"created_by" db:"created_by"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time      `json:"deleted_at" db:"deleted_at"`
}

// ============================================================================
// Approval Stage
// ============================================================================

type ApprovalStage struct {
	ID              uuid.UUID `json:"id" db:"id"`
	WorkflowID      uuid.UUID `json:"workflow_id" db:"workflow_id"`
	Name            string    `json:"name" db:"name"`
	SequenceOrder   int       `json:"sequence_order" db:"sequence_order"`
	ApproverRole    string    `json:"approver_role" db:"approver_role"`
	RequiredCount   int       `json:"required_count" db:"required_count"`
	TimeoutHours    int       `json:"timeout_hours" db:"timeout_hours"`
	Status          string    `json:"status" db:"status"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// ============================================================================
// Content
// ============================================================================

type Content struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	OrgID        uuid.UUID       `json:"org_id" db:"org_id"`
	WorkflowID   *uuid.UUID      `json:"workflow_id" db:"workflow_id"`
	Title        string          `json:"title" db:"title"`
	Description  string          `json:"description" db:"description"`
	ContentType  string          `json:"content_type" db:"content_type"`
	Status       string          `json:"status" db:"status"`
	Metadata     json.RawMessage `json:"metadata" db:"metadata"`
	CreatedBy    uuid.UUID       `json:"created_by" db:"created_by"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
	PublishedAt  *time.Time      `json:"published_at" db:"published_at"`
	DeletedAt    *time.Time      `json:"deleted_at" db:"deleted_at"`
}

// ============================================================================
// Approval
// ============================================================================

type Approval struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ContentID  uuid.UUID `json:"content_id" db:"content_id"`
	StageID    uuid.UUID `json:"stage_id" db:"stage_id"`
	ApproverID uuid.UUID `json:"approver_id" db:"approver_id"`
	Decision   string    `json:"decision" db:"decision"`
	Comment    string    `json:"comment" db:"comment"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// ============================================================================
// Agent
// ============================================================================

type Agent struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	OrgID        uuid.UUID       `json:"org_id" db:"org_id"`
	Name         string          `json:"name" db:"name"`
	AgentType    string          `json:"agent_type" db:"agent_type"`
	Description  string          `json:"description" db:"description"`
	Config       json.RawMessage `json:"config" db:"config"`
	Status       string          `json:"status" db:"status"`
	LastRunAt    *time.Time      `json:"last_run_at" db:"last_run_at"`
	ErrorMessage string          `json:"error_message" db:"error_message"`
	CreatedBy    uuid.UUID       `json:"created_by" db:"created_by"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time      `json:"deleted_at" db:"deleted_at"`
}

// ============================================================================
// Agent Run
// ============================================================================

type AgentRun struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	AgentID       uuid.UUID       `json:"agent_id" db:"agent_id"`
	ContentID     *uuid.UUID      `json:"content_id" db:"content_id"`
	Status        string          `json:"status" db:"status"`
	InputData     json.RawMessage `json:"input_data" db:"input_data"`
	OutputData    json.RawMessage `json:"output_data" db:"output_data"`
	ErrorMessage  string          `json:"error_message" db:"error_message"`
	DurationMs    int             `json:"duration_ms" db:"duration_ms"`
	StartedAt     time.Time       `json:"started_at" db:"started_at"`
	CompletedAt   *time.Time      `json:"completed_at" db:"completed_at"`
}

// ============================================================================
// Distribution
// ============================================================================

type Distribution struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	OrgID        uuid.UUID       `json:"org_id" db:"org_id"`
	ContentID    uuid.UUID       `json:"content_id" db:"content_id"`
	Channel      string          `json:"channel" db:"channel"`
	Status       string          `json:"status" db:"status"`
	PublishedURL string          `json:"published_url" db:"published_url"`
	PublishedAt  *time.Time      `json:"published_at" db:"published_at"`
	Metadata     json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
}

// ============================================================================
// Audit Log
// ============================================================================

type AuditLog struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	OrgID        uuid.UUID       `json:"org_id" db:"org_id"`
	UserID       *uuid.UUID      `json:"user_id" db:"user_id"`
	Action       string          `json:"action" db:"action"`
	ResourceType string          `json:"resource_type" db:"resource_type"`
	ResourceID   *uuid.UUID      `json:"resource_id" db:"resource_id"`
	Changes      json.RawMessage `json:"changes" db:"changes"`
	IPAddress    string          `json:"ip_address" db:"ip_address"`
	UserAgent    string          `json:"user_agent" db:"user_agent"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
}

// ============================================================================
// Request/Response DTOs
// ============================================================================

type CreateOrganizationRequest struct {
	Name        string                 `json:"name" validate:"required,min=3,max=255"`
	Slug        string                 `json:"slug" validate:"required,alphanum,min=3,max=255"`
	Description string                 `json:"description" validate:"max=1000"`
	Tier        string                 `json:"tier" validate:"oneof=free pro enterprise"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=255"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"oneof=admin editor viewer user"`
}

type CreateContentRequest struct {
	Title       string                 `json:"title" validate:"required,min=3,max=255"`
	Description string                 `json:"description" validate:"max=5000"`
	ContentType string                 `json:"content_type" validate:"required,oneof=article video image podcast social"`
	WorkflowID  *uuid.UUID             `json:"workflow_id"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type CreateWorkflowRequest struct {
	Name         string                 `json:"name" validate:"required,min=3,max=255"`
	Description  string                 `json:"description" validate:"max=1000"`
	WorkflowType string                 `json:"workflow_type" validate:"required,oneof=content-creation approval publishing distribution"`
	Config       map[string]interface{} `json:"config" validate:"required"`
}

type CreateApprovalRequest struct {
	ContentID uuid.UUID `json:"content_id" validate:"required"`
	StageID   uuid.UUID `json:"stage_id" validate:"required"`
	Decision  string    `json:"decision" validate:"required,oneof=approved rejected requested-changes"`
	Comment   string    `json:"comment" validate:"max=2000"`
}

type CreateAgentRequest struct {
	Name        string                 `json:"name" validate:"required,min=3,max=255"`
	AgentType   string                 `json:"agent_type" validate:"required,oneof=content-creator approver publisher distributor"`
	Description string                 `json:"description" validate:"max=1000"`
	Config      map[string]interface{} `json:"config" validate:"required"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
}

type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Uptime    int64     `json:"uptime_ms"`
}

// ============================================================================
// Context Keys
// ============================================================================

type contextKey string

const (
	ContextKeyOrgID    contextKey = "org_id"
	ContextKeyUserID   contextKey = "user_id"
	ContextKeyUserRole contextKey = "user_role"
)

