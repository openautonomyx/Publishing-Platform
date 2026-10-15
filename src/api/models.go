package main

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ============================================================================
// JSONB
// ============================================================================

// JSONB is a JSON document stored in a PostgreSQL JSONB column. Unlike
// json.RawMessage it scans SQL NULL (as JSON null) and marshals as raw JSON.
type JSONB []byte

func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return []byte(j), nil
}

func (j *JSONB) Scan(src interface{}) error {
	switch v := src.(type) {
	case nil:
		*j = nil
	case []byte:
		*j = append((*j)[:0], v...)
	case string:
		*j = JSONB(v)
	default:
		return fmt.Errorf("JSONB: cannot scan %T", src)
	}
	return nil
}

func (j JSONB) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return j, nil
}

func (j *JSONB) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*j = nil
		return nil
	}
	*j = append((*j)[:0], data...)
	return nil
}

// jsonObject returns m as a JSONB object, defaulting to {}.
func jsonObject(m map[string]interface{}) (JSONB, error) {
	if m == nil {
		return JSONB("{}"), nil
	}
	b, err := json.Marshal(m)
	return JSONB(b), err
}

// ============================================================================
// Organization (Tenant)
// ============================================================================

type Organization struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	ExternalID  *string    `json:"external_id"`
	Description *string    `json:"description"`
	Tier        string     `json:"tier"` // free, pro, enterprise
	Status      string     `json:"status"`
	Metadata    JSONB      `json:"metadata"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CreatedBy   *uuid.UUID `json:"created_by"`
}

// ============================================================================
// User (Multi-tenant)
// ============================================================================

type User struct {
	ID          uuid.UUID  `json:"id"`
	OrgID       uuid.UUID  `json:"org_id"`
	ExternalID  *string    `json:"external_id"`
	Email       *string    `json:"email"`
	Name        string     `json:"name"`
	AvatarURL   *string    `json:"avatar_url"`
	Role        string     `json:"role"` // admin, approver, editor, viewer
	Status      string     `json:"status"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ============================================================================
// API Key
// ============================================================================

type APIKey struct {
	ID          uuid.UUID      `json:"id"`
	OrgID       uuid.UUID      `json:"org_id"`
	UserID      uuid.UUID      `json:"user_id"`
	KeyPrefix   string         `json:"key_prefix"`
	Name        *string        `json:"name"`
	Permissions pq.StringArray `json:"permissions"`
	Status      string         `json:"status"`
	LastUsedAt  *time.Time     `json:"last_used_at"`
	ExpiresAt   *time.Time     `json:"expires_at"`
	CreatedAt   time.Time      `json:"created_at"`
	// Key is the plaintext key, returned only by the create call.
	Key string `json:"key,omitempty"`
}

// ============================================================================
// Workflow
// ============================================================================

type Workflow struct {
	ID           uuid.UUID       `json:"id"`
	OrgID        uuid.UUID       `json:"org_id"`
	Name         string          `json:"name"`
	Description  *string         `json:"description"`
	WorkflowType string          `json:"workflow_type"`
	Status       string          `json:"status"`
	Config       JSONB           `json:"config"`
	CreatedBy    uuid.UUID       `json:"created_by"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Stages       []ApprovalStage `json:"stages,omitempty"`
}

// ============================================================================
// Approval Stage
// ============================================================================

type ApprovalStage struct {
	ID            uuid.UUID `json:"id"`
	WorkflowID    uuid.UUID `json:"workflow_id"`
	Name          string    `json:"name"`
	SequenceOrder int       `json:"sequence_order"`
	ApproverRole  string    `json:"approver_role"`
	RequiredCount int       `json:"required_count"`
	TimeoutHours  int       `json:"timeout_hours"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ============================================================================
// Content
// ============================================================================

type Content struct {
	ID          uuid.UUID  `json:"id"`
	OrgID       uuid.UUID  `json:"org_id"`
	WorkflowID  *uuid.UUID `json:"workflow_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	ContentType string     `json:"content_type"`
	Status      string     `json:"status"`
	Metadata    JSONB      `json:"metadata"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	SubmittedAt *time.Time `json:"submitted_at"`
	PublishedAt *time.Time `json:"published_at"`
}

// ============================================================================
// Approval
// ============================================================================

type Approval struct {
	ID         uuid.UUID `json:"id"`
	ContentID  uuid.UUID `json:"content_id"`
	StageID    uuid.UUID `json:"stage_id"`
	ApproverID uuid.UUID `json:"approver_id"`
	Decision   string    `json:"decision"`
	Comment    *string   `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ============================================================================
// Agent
// ============================================================================

type Agent struct {
	ID           uuid.UUID  `json:"id"`
	OrgID        uuid.UUID  `json:"org_id"`
	Name         string     `json:"name"`
	AgentType    string     `json:"agent_type"`
	Description  *string    `json:"description"`
	Config       JSONB      `json:"config"`
	Status       string     `json:"status"`
	LastRunAt    *time.Time `json:"last_run_at"`
	ErrorMessage *string    `json:"error_message"`
	CreatedBy    uuid.UUID  `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// ============================================================================
// Agent Run
// ============================================================================

type AgentRun struct {
	ID           uuid.UUID  `json:"id"`
	AgentID      uuid.UUID  `json:"agent_id"`
	ContentID    *uuid.UUID `json:"content_id"`
	Status       string     `json:"status"`
	InputData    JSONB      `json:"input_data"`
	OutputData   JSONB      `json:"output_data"`
	ErrorMessage *string    `json:"error_message"`
	DurationMs   *int       `json:"duration_ms"`
	StartedAt    time.Time  `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at"`
}

// ============================================================================
// Distribution
// ============================================================================

type Distribution struct {
	ID           uuid.UUID  `json:"id"`
	OrgID        uuid.UUID  `json:"org_id"`
	ContentID    uuid.UUID  `json:"content_id"`
	Channel      string     `json:"channel"`
	Status       string     `json:"status"`
	ScheduledAt  *time.Time `json:"scheduled_at"`
	PublishedURL *string    `json:"published_url"`
	PublishedAt  *time.Time `json:"published_at"`
	Metadata     JSONB      `json:"metadata"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// ============================================================================
// Audit Log
// ============================================================================

type AuditLog struct {
	ID           uuid.UUID  `json:"id"`
	OrgID        uuid.UUID  `json:"org_id"`
	UserID       *uuid.UUID `json:"user_id"`
	Action       string     `json:"action"`
	ResourceType *string    `json:"resource_type"`
	ResourceID   *uuid.UUID `json:"resource_id"`
	Changes      JSONB      `json:"changes"`
	IPAddress    *string    `json:"ip_address"`
	UserAgent    *string    `json:"user_agent"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ============================================================================
// Request DTOs
// ============================================================================

// Pointer fields in Update* requests are optional: nil leaves the value unchanged.

type UpdateOrganizationRequest struct {
	Name        *string                `json:"name" validate:"omitempty,min=3,max=255"`
	Slug        *string                `json:"slug" validate:"omitempty,slug,min=3,max=255"`
	Description *string                `json:"description" validate:"omitempty,max=1000"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type CreateUserRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
	Name  string `json:"name" validate:"required,min=2,max=255"`
}

type UpdateUserRequest struct {
	Name      *string `json:"name" validate:"omitempty,min=2,max=255"`
	AvatarURL *string `json:"avatar_url" validate:"omitempty,url,max=2000"`
	Status    *string `json:"status" validate:"omitempty,oneof=active inactive"`
}

type CreateContentRequest struct {
	Title       string                 `json:"title" validate:"required,min=3,max=255"`
	Description string                 `json:"description" validate:"max=5000"`
	ContentType string                 `json:"content_type" validate:"required,oneof=article video image podcast social"`
	WorkflowID  *uuid.UUID             `json:"workflow_id"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type UpdateContentRequest struct {
	Title       *string                `json:"title" validate:"omitempty,min=3,max=255"`
	Description *string                `json:"description" validate:"omitempty,max=5000"`
	WorkflowID  *uuid.UUID             `json:"workflow_id"`
	Status      *string                `json:"status" validate:"omitempty,oneof=draft in-review published archived"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type StageRequest struct {
	Name          string `json:"name" validate:"required,min=1,max=255"`
	ApproverRole  string `json:"approver_role" validate:"required,oneof=approver admin"`
	RequiredCount int    `json:"required_count" validate:"omitempty,min=1,max=20"`
	TimeoutHours  int    `json:"timeout_hours" validate:"omitempty,min=1,max=8760"`
}

type CreateWorkflowRequest struct {
	Name         string                 `json:"name" validate:"required,min=3,max=255"`
	Description  string                 `json:"description" validate:"max=1000"`
	WorkflowType string                 `json:"workflow_type" validate:"required,oneof=content-creation approval publishing distribution"`
	Config       map[string]interface{} `json:"config"`
	Stages       []StageRequest         `json:"stages" validate:"max=20,dive"`
}

type UpdateWorkflowRequest struct {
	Name        *string                `json:"name" validate:"omitempty,min=3,max=255"`
	Description *string                `json:"description" validate:"omitempty,max=1000"`
	Status      *string                `json:"status" validate:"omitempty,oneof=draft active archived"`
	Config      map[string]interface{} `json:"config"`
	// Stages, when present, replaces all stages. Refused once any approval references them.
	Stages *[]StageRequest `json:"stages" validate:"omitempty,max=20,dive"`
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
	Config      map[string]interface{} `json:"config"`
}

type UpdateAgentRequest struct {
	Name        *string                `json:"name" validate:"omitempty,min=3,max=255"`
	Description *string                `json:"description" validate:"omitempty,max=1000"`
	Status      *string                `json:"status" validate:"omitempty,oneof=active paused"`
	Config      map[string]interface{} `json:"config"`
}

type RunAgentRequest struct {
	ContentID *uuid.UUID             `json:"content_id"`
	Input     map[string]interface{} `json:"input"`
}

type CreateDistributionRequest struct {
	ContentID   uuid.UUID              `json:"content_id" validate:"required"`
	Channel     string                 `json:"channel" validate:"required,oneof=twitter linkedin medium substack website"`
	ScheduledAt *time.Time             `json:"scheduled_at"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type CreateAPIKeyRequest struct {
	Name          string   `json:"name" validate:"required,min=3,max=255"`
	Permissions   []string `json:"permissions" validate:"required,min=1,dive,oneof=read write approve admin"`
	ExpiresInDays int      `json:"expires_in_days" validate:"omitempty,min=1,max=3650"`
}

// ============================================================================
// Responses
// ============================================================================

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
