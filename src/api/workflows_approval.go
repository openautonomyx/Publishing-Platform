package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ApprovalDecision string

const (
	APPROVED        ApprovalDecision = "APPROVED"
	REJECTED        ApprovalDecision = "REJECTED"
	REQUEST_CHANGES ApprovalDecision = "REQUEST_CHANGES"
)

// ApprovalRequest represents a content approval submission
type ApprovalRequest struct {
	ContentID    uuid.UUID
	Title        string
	AuthorID     uuid.UUID
	AuthorEmail  string
	ReviewerIDs  []uuid.UUID
	SubmittedAt  time.Time
	DueDate      time.Time
	OrgID        uuid.UUID
	Description  string
	Tags         []string
}

// ApprovalResult represents the approval decision
type ApprovalResult struct {
	Decision        ApprovalDecision
	RejectionReason string
	Feedback        string
	ReviewerID      uuid.UUID
	ReviewedAt      time.Time
}

// SubmitForApproval submits content for approval workflow
func SubmitForApproval(ctx context.Context, req ApprovalRequest) (string, error) {
	logger.Info("submitting content for approval",
		zap.String("content_id", req.ContentID.String()),
		zap.String("author", req.AuthorID.String()),
		zap.Int("reviewers", len(req.ReviewerIDs)),
	)

	// Validate request
	if len(req.ReviewerIDs) == 0 {
		return "", fmt.Errorf("at least one reviewer required")
	}

	// Create approval record in database
	approvalID := uuid.New()
	workflowID := fmt.Sprintf("approval-%s-%d", approvalID.String()[:8], time.Now().Unix())

	// Insert into content_approvals table
	query := `
		INSERT INTO content_approvals (
			id, content_id, workflow_id, org_id, author_id,
			submitted_at, due_date, current_status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	err := db.QueryRowContext(
		ctx, query,
		approvalID, req.ContentID, workflowID, req.OrgID, req.AuthorID,
		req.SubmittedAt, req.DueDate, "PENDING",
	).Err()

	if err != nil {
		logger.Error("failed to create approval record", zap.Error(err))
		return "", err
	}

	// Add reviewers to approval queue
	for _, reviewerID := range req.ReviewerIDs {
		queueQuery := `
			INSERT INTO approval_queue (approval_id, reviewer_id, status)
			VALUES ($1, $2, $3)
		`
		_, err := db.ExecContext(ctx, queueQuery, approvalID, reviewerID, "PENDING")
		if err != nil {
			logger.Error("failed to add reviewer to queue", zap.Error(err))
		}
	}

	// Record audit log
	RecordAuditLog(ctx, AuditLogEntry{
		UserID:       req.AuthorID,
		OrgID:        req.OrgID,
		Action:       "SUBMIT_FOR_APPROVAL",
		ResourceType: "CONTENT",
		ResourceID:   req.ContentID,
		Timestamp:    time.Now(),
		Details: map[string]string{
			"workflow_id": workflowID,
			"reviewers":   fmt.Sprintf("%d", len(req.ReviewerIDs)),
		},
	})

	logger.Info("approval workflow started",
		zap.String("approval_id", approvalID.String()),
		zap.String("workflow_id", workflowID),
	)

	return approvalID.String(), nil
}

// ApproveContent approves content submission
func ApproveContent(ctx context.Context, approvalID string, reviewerID uuid.UUID) error {
	logger.Info("approving content",
		zap.String("approval_id", approvalID),
		zap.String("reviewer_id", reviewerID.String()),
	)

	now := time.Now()

	// Update approval record
	query := `
		UPDATE content_approvals
		SET current_status = $1, decision = $2, reviewer_id = $3, reviewed_at = $4, updated_at = $5
		WHERE id = $6
	`

	result, err := db.ExecContext(
		ctx, query,
		"APPROVED", "APPROVED", reviewerID, now, now,
		approvalID,
	)

	if err != nil {
		logger.Error("failed to approve content", zap.Error(err))
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("approval not found: %s", approvalID)
	}

	// Update approval queue
	_, err = db.ExecContext(
		ctx,
		"UPDATE approval_queue SET status = $1, action_taken_at = $2 WHERE approval_id = $3 AND reviewer_id = $4",
		"APPROVED", now, approvalID, reviewerID,
	)

	if err != nil {
		logger.Error("failed to update approval queue", zap.Error(err))
		return err
	}

	// Get approval details for audit log
	var contentID uuid.UUID
	var orgID uuid.UUID
	err = db.QueryRowContext(
		ctx,
		"SELECT content_id, org_id FROM content_approvals WHERE id = $1",
		approvalID,
	).Scan(&contentID, &orgID)

	if err == nil {
		RecordAuditLog(ctx, AuditLogEntry{
			UserID:       reviewerID,
			OrgID:        orgID,
			Action:       "APPROVE_CONTENT",
			ResourceType: "CONTENT",
			ResourceID:   contentID,
			Timestamp:    now,
		})
	}

	return nil
}

// RejectContent rejects content with feedback
func RejectContent(ctx context.Context, approvalID string, reviewerID uuid.UUID, reason string) error {
	logger.Info("rejecting content",
		zap.String("approval_id", approvalID),
		zap.String("reviewer_id", reviewerID.String()),
		zap.String("reason", reason),
	)

	now := time.Now()

	// Update approval record
	query := `
		UPDATE content_approvals
		SET current_status = $1, decision = $2, reviewer_id = $3, reviewed_at = $4,
		    rejection_reason = $5, updated_at = $6
		WHERE id = $7
	`

	result, err := db.ExecContext(
		ctx, query,
		"REJECTED", "REJECTED", reviewerID, now, reason, now,
		approvalID,
	)

	if err != nil {
		logger.Error("failed to reject content", zap.Error(err))
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("approval not found: %s", approvalID)
	}

	// Update approval queue
	_, err = db.ExecContext(
		ctx,
		"UPDATE approval_queue SET status = $1, action_taken_at = $2 WHERE approval_id = $3 AND reviewer_id = $4",
		"REJECTED", now, approvalID, reviewerID,
	)

	if err != nil {
		logger.Error("failed to update approval queue", zap.Error(err))
		return err
	}

	// Get approval details for audit log
	var contentID uuid.UUID
	var orgID uuid.UUID
	err = db.QueryRowContext(
		ctx,
		"SELECT content_id, org_id FROM content_approvals WHERE id = $1",
		approvalID,
	).Scan(&contentID, &orgID)

	if err == nil {
		RecordAuditLog(ctx, AuditLogEntry{
			UserID:       reviewerID,
			OrgID:        orgID,
			Action:       "REJECT_CONTENT",
			ResourceType: "CONTENT",
			ResourceID:   contentID,
			Timestamp:    now,
			Details: map[string]string{
				"reason": reason,
			},
		})
	}

	return nil
}

// RequestChanges asks author to revise content
func RequestChanges(ctx context.Context, approvalID string, reviewerID uuid.UUID, feedback string) error {
	logger.Info("requesting content changes",
		zap.String("approval_id", approvalID),
		zap.String("reviewer_id", reviewerID.String()),
	)

	now := time.Now()

	// Update approval record
	query := `
		UPDATE content_approvals
		SET current_status = $1, decision = $2, reviewer_id = $3, reviewed_at = $4,
		    feedback = $5, updated_at = $6, revision_count = revision_count + 1
		WHERE id = $7
	`

	result, err := db.ExecContext(
		ctx, query,
		"REVISION_REQUESTED", "REQUEST_CHANGES", reviewerID, now, feedback, now,
		approvalID,
	)

	if err != nil {
		logger.Error("failed to request changes", zap.Error(err))
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("approval not found: %s", approvalID)
	}

	// Update approval queue
	_, err = db.ExecContext(
		ctx,
		"UPDATE approval_queue SET status = $1, action_taken_at = $2 WHERE approval_id = $3 AND reviewer_id = $4",
		"REVISION_REQUESTED", now, approvalID, reviewerID,
	)

	if err != nil {
		logger.Error("failed to update approval queue", zap.Error(err))
		return err
	}

	// Get approval details for audit log
	var contentID uuid.UUID
	var orgID uuid.UUID
	err = db.QueryRowContext(
		ctx,
		"SELECT content_id, org_id FROM content_approvals WHERE id = $1",
		approvalID,
	).Scan(&contentID, &orgID)

	if err == nil {
		RecordAuditLog(ctx, AuditLogEntry{
			UserID:       reviewerID,
			OrgID:        orgID,
			Action:       "REQUEST_CHANGES",
			ResourceType: "CONTENT",
			ResourceID:   contentID,
			Timestamp:    now,
			Details: map[string]string{
				"feedback": feedback,
			},
		})
	}

	return nil
}

// GetApprovalStatus returns current approval status
func GetApprovalStatus(ctx context.Context, contentID uuid.UUID) (map[string]interface{}, error) {
	query := `
		SELECT id, workflow_id, current_status, decision, reviewer_id,
		       reviewed_at, rejection_reason, feedback, revision_count,
		       submitted_at, due_date
		FROM content_approvals
		WHERE content_id = $1
		ORDER BY submitted_at DESC
		LIMIT 1
	`

	var (
		id              uuid.UUID
		workflowID      string
		status          string
		decision        *string
		reviewerID      *uuid.UUID
		reviewedAt      *time.Time
		rejectionReason *string
		feedback        *string
		revisionCount   int
		submittedAt     time.Time
		dueDate         *time.Time
	)

	err := db.QueryRowContext(ctx, query, contentID).Scan(
		&id, &workflowID, &status, &decision, &reviewerID,
		&reviewedAt, &rejectionReason, &feedback, &revisionCount,
		&submittedAt, &dueDate,
	)

	if err != nil {
		logger.Error("failed to get approval status", zap.Error(err))
		return nil, err
	}

	return map[string]interface{}{
		"id":                id.String(),
		"workflow_id":       workflowID,
		"status":            status,
		"decision":          decision,
		"reviewer_id":       reviewerID,
		"reviewed_at":       reviewedAt,
		"rejection_reason":  rejectionReason,
		"feedback":          feedback,
		"revision_count":    revisionCount,
		"submitted_at":      submittedAt,
		"due_date":          dueDate,
	}, nil
}

// ListPendingApprovals returns all pending approvals for a reviewer
func ListPendingApprovals(ctx context.Context, reviewerID uuid.UUID, limit int) ([]map[string]interface{}, error) {
	if limit == 0 {
		limit = 10
	}

	query := `
		SELECT
			ca.id, ca.content_id, ca.title, ca.author_id, ca.submitted_at,
			ca.due_date, ca.current_status, ca.revision_count
		FROM content_approvals ca
		INNER JOIN approval_queue aq ON ca.id = aq.approval_id
		WHERE aq.reviewer_id = $1 AND aq.status = 'PENDING'
		ORDER BY ca.due_date ASC NULLS LAST
		LIMIT $2
	`

	rows, err := db.QueryContext(ctx, query, reviewerID, limit)
	if err != nil {
		logger.Error("failed to list pending approvals", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var approvals []map[string]interface{}

	for rows.Next() {
		var (
			id           uuid.UUID
			contentID    uuid.UUID
			title        string
			authorID     uuid.UUID
			submittedAt  time.Time
			dueDate      *time.Time
			status       string
			revisionCount int
		)

		err := rows.Scan(&id, &contentID, &title, &authorID, &submittedAt, &dueDate, &status, &revisionCount)
		if err != nil {
			logger.Error("failed to scan approval row", zap.Error(err))
			continue
		}

		approvals = append(approvals, map[string]interface{}{
			"id":              id.String(),
			"content_id":      contentID.String(),
			"title":           title,
			"author_id":       authorID.String(),
			"submitted_at":    submittedAt,
			"due_date":        dueDate,
			"status":          status,
			"revision_count":  revisionCount,
		})
	}

	return approvals, nil
}

// GetApprovalHistory returns complete history of approval actions
func GetApprovalHistory(ctx context.Context, approvalID string) ([]map[string]interface{}, error) {
	query := `
		SELECT timestamp, action, user_id, details
		FROM approval_audit_log
		WHERE approval_id = $1
		ORDER BY timestamp DESC
	`

	rows, err := db.QueryContext(ctx, query, approvalID)
	if err != nil {
		logger.Error("failed to get approval history", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var history []map[string]interface{}

	for rows.Next() {
		var (
			timestamp time.Time
			action    string
			userID    uuid.UUID
			details   *string
		)

		err := rows.Scan(&timestamp, &action, &userID, &details)
		if err != nil {
			logger.Error("failed to scan history row", zap.Error(err))
			continue
		}

		history = append(history, map[string]interface{}{
			"timestamp": timestamp,
			"action":    action,
			"user_id":   userID.String(),
			"details":   details,
		})
	}

	return history, nil
}
