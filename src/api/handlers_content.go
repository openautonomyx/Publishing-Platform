package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// Content
//
// Lifecycle:
//
//	draft ──submit──► in-review ──all stages approved──► approved ──► published ──► archived
//	  ▲                  │  rejected / requested-changes        │
//	  └──────────────────┴──────────────────────────────────────┘ (back to draft)
//
// Content without a workflow can go straight from draft to published. Title,
// description, metadata and workflow can only be edited in draft.
// ============================================================================

const contentColumns = `id, org_id, workflow_id, title, description, content_type, status, metadata,
	created_by, created_at, updated_at, submitted_at, published_at`

func scanContent(row interface{ Scan(...interface{}) error }) (*Content, error) {
	c := &Content{}
	err := row.Scan(&c.ID, &c.OrgID, &c.WorkflowID, &c.Title, &c.Description, &c.ContentType, &c.Status,
		&c.Metadata, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt, &c.SubmittedAt, &c.PublishedAt)
	return c, err
}

// contentTransitions lists the status changes PUT /content/{id} may make. Moving to
// "approved" happens only through approvals.
var contentTransitions = map[string]map[string]bool{
	"draft":     {"in-review": true, "published": true, "archived": true},
	"in-review": {"draft": true, "archived": true},
	"approved":  {"published": true, "draft": true, "archived": true},
	"published": {"archived": true},
	"archived":  {"draft": true},
}

func (a *App) createContent(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	var req CreateContentRequest
	if err := decode(r, &req); err != nil {
		a.fail(w, r, err)
		return
	}
	metadata, err := jsonObject(req.Metadata)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	var c *Content
	err = a.inTx(r.Context(), func(tx *sql.Tx) error {
		if req.WorkflowID != nil {
			if err := activeWorkflow(r.Context(), tx, id.OrgID, *req.WorkflowID); err != nil {
				return err
			}
		}
		var err error
		c, err = scanContent(tx.QueryRowContext(r.Context(), `
			INSERT INTO content (org_id, workflow_id, title, description, content_type, status, metadata, created_by)
			VALUES ($1, $2, $3, $4, $5, 'draft', $6, $7)
			RETURNING `+contentColumns,
			id.OrgID, req.WorkflowID, req.Title, nullString(req.Description), req.ContentType, metadata, id.UserID))
		if err != nil {
			return err
		}
		return a.audit(r.Context(), tx, r, "content.create", "content", c.ID, map[string]interface{}{"title": req.Title})
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondData(w, http.StatusCreated, c)
}

func (a *App) listContent(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	p := pagination(r)
	q := r.URL.Query()
	f := &queryFilter{args: []interface{}{id.OrgID}}
	if v := q.Get("status"); v != "" {
		f.add("status = $%d", v)
	}
	if v := q.Get("content_type"); v != "" {
		f.add("content_type = $%d", v)
	}
	if v := q.Get("workflow_id"); v != "" {
		wf, err := uuid.Parse(v)
		if err != nil {
			a.fail(w, r, NewBadRequestError("invalid workflow_id"))
			return
		}
		f.add("workflow_id = $%d", wf)
	}
	where := `FROM content WHERE org_id = $1 AND deleted_at IS NULL` + f.sql()

	var total int
	if err := a.db.QueryRowContext(r.Context(), `SELECT COUNT(*) `+where, f.args...).Scan(&total); err != nil {
		a.fail(w, r, err)
		return
	}
	items, err := queryList(r.Context(), a.db, scanContent,
		`SELECT `+contentColumns+` `+where+` ORDER BY created_at DESC LIMIT `+limitOffset(p), f.args...)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondPage(w, items, total, p)
}

func (a *App) getContent(w http.ResponseWriter, r *http.Request) {
	contentID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	c, err := getContentRow(r.Context(), a.db, identityFrom(r.Context()).OrgID, contentID, false)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondData(w, http.StatusOK, c)
}

func (a *App) updateContent(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	contentID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	var req UpdateContentRequest
	if err := decode(r, &req); err != nil {
		a.fail(w, r, err)
		return
	}
	var metadata interface{}
	if req.Metadata != nil {
		m, err := jsonObject(req.Metadata)
		if err != nil {
			a.fail(w, r, err)
			return
		}
		metadata = m
	}
	editsBody := req.Title != nil || req.Description != nil || req.Metadata != nil || req.WorkflowID != nil

	var c *Content
	err = a.inTx(r.Context(), func(tx *sql.Tx) error {
		cur, err := getContentRow(r.Context(), tx, id.OrgID, contentID, true)
		if err != nil {
			return err
		}
		if editsBody && cur.Status != "draft" {
			return NewConflictError("content can only be edited in draft; move it back to draft first")
		}
		workflowID := cur.WorkflowID
		if req.WorkflowID != nil {
			if err := activeWorkflow(r.Context(), tx, id.OrgID, *req.WorkflowID); err != nil {
				return err
			}
			workflowID = req.WorkflowID
		}

		status := cur.Status
		if req.Status != nil && *req.Status != cur.Status {
			status = *req.Status
			if !contentTransitions[cur.Status][status] {
				return NewConflictError("cannot move content from " + cur.Status + " to " + status)
			}
			switch {
			case status == "in-review" && workflowID == nil:
				return NewConflictError("content needs a workflow before it can be submitted for review")
			case status == "in-review":
				var stages int
				if err := tx.QueryRowContext(r.Context(),
					`SELECT COUNT(*) FROM approval_stages WHERE workflow_id = $1 AND status = 'active'`, *workflowID,
				).Scan(&stages); err != nil {
					return err
				}
				if stages == 0 {
					return NewConflictError("the workflow has no approval stages")
				}
			case status == "published" && cur.Status == "draft" && workflowID != nil:
				return NewConflictError("content with a workflow must be approved before publishing")
			}
		}

		c, err = scanContent(tx.QueryRowContext(r.Context(), `
			UPDATE content SET
				title = COALESCE($3, title), description = COALESCE($4, description),
				metadata = COALESCE($5, metadata), workflow_id = $6, status = $7::text,
				submitted_at = CASE WHEN $7::text = 'in-review' AND status <> 'in-review' THEN NOW() ELSE submitted_at END,
				published_at = CASE WHEN $7::text = 'published' AND status <> 'published' THEN NOW() ELSE published_at END
			WHERE id = $1 AND org_id = $2
			RETURNING `+contentColumns,
			contentID, id.OrgID, req.Title, req.Description, metadata, workflowID, status))
		if err != nil {
			return err
		}
		changes := map[string]interface{}{"fields": req}
		if status != cur.Status {
			changes["status"] = map[string]string{"from": cur.Status, "to": status}
		}
		return a.audit(r.Context(), tx, r, "content.update", "content", contentID, changes)
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondData(w, http.StatusOK, c)
}

func (a *App) deleteContent(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	contentID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	err = a.inTx(r.Context(), func(tx *sql.Tx) error {
		if err := execOne(r.Context(), tx, "content", `
			UPDATE content SET deleted_at = NOW() WHERE id = $1 AND org_id = $2 AND deleted_at IS NULL`,
			contentID, id.OrgID); err != nil {
			return err
		}
		return a.audit(r.Context(), tx, r, "content.delete", "content", contentID, nil)
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Message: "content deleted"})
}

func getContentRow(ctx context.Context, q querier, orgID, contentID uuid.UUID, forUpdate bool) (*Content, error) {
	query := `SELECT ` + contentColumns + ` FROM content WHERE id = $1 AND org_id = $2 AND deleted_at IS NULL`
	if forUpdate {
		query += ` FOR UPDATE`
	}
	c, err := scanContent(q.QueryRowContext(ctx, query, contentID, orgID))
	return c, notFound(err, "content")
}

func activeWorkflow(ctx context.Context, q querier, orgID, workflowID uuid.UUID) error {
	var status string
	err := q.QueryRowContext(ctx,
		`SELECT status FROM workflows WHERE id = $1 AND org_id = $2 AND deleted_at IS NULL`, workflowID, orgID,
	).Scan(&status)
	if isNoRows(err) {
		return NewBadRequestError("workflow not found")
	}
	if err != nil {
		return err
	}
	if status != "active" {
		return NewBadRequestError("workflow is " + status + ", not active")
	}
	return nil
}

// ============================================================================
// Workflows (with ordered approval stages)
// ============================================================================

const workflowColumns = `id, org_id, name, description, workflow_type, status, config, created_by, created_at, updated_at`
const stageColumns = `id, workflow_id, name, sequence_order, approver_role, required_count, timeout_hours, status, created_at, updated_at`

func scanWorkflow(row interface{ Scan(...interface{}) error }) (*Workflow, error) {
	wf := &Workflow{}
	err := row.Scan(&wf.ID, &wf.OrgID, &wf.Name, &wf.Description, &wf.WorkflowType, &wf.Status,
		&wf.Config, &wf.CreatedBy, &wf.CreatedAt, &wf.UpdatedAt)
	return wf, err
}

func scanStage(row interface{ Scan(...interface{}) error }) (ApprovalStage, error) {
	s := ApprovalStage{}
	err := row.Scan(&s.ID, &s.WorkflowID, &s.Name, &s.SequenceOrder, &s.ApproverRole, &s.RequiredCount,
		&s.TimeoutHours, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func loadStages(ctx context.Context, q querier, workflowID uuid.UUID) ([]ApprovalStage, error) {
	return queryList(ctx, q, scanStage,
		`SELECT `+stageColumns+` FROM approval_stages WHERE workflow_id = $1 ORDER BY sequence_order`, workflowID)
}

func insertStages(ctx context.Context, tx *sql.Tx, workflowID uuid.UUID, stages []StageRequest) error {
	for i, s := range stages {
		required, timeout := s.RequiredCount, s.TimeoutHours
		if required == 0 {
			required = 1
		}
		if timeout == 0 {
			timeout = 48
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO approval_stages (workflow_id, name, sequence_order, approver_role, required_count, timeout_hours)
			VALUES ($1, $2, $3, $4, $5, $6)`, workflowID, s.Name, i+1, s.ApproverRole, required, timeout); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) createWorkflow(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	var req CreateWorkflowRequest
	if err := decode(r, &req); err != nil {
		a.fail(w, r, err)
		return
	}
	config, err := jsonObject(req.Config)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	var wf *Workflow
	err = a.inTx(r.Context(), func(tx *sql.Tx) error {
		var err error
		wf, err = scanWorkflow(tx.QueryRowContext(r.Context(), `
			INSERT INTO workflows (org_id, name, description, workflow_type, status, config, created_by)
			VALUES ($1, $2, $3, $4, 'active', $5, $6)
			RETURNING `+workflowColumns,
			id.OrgID, req.Name, nullString(req.Description), req.WorkflowType, config, id.UserID))
		if err != nil {
			return err
		}
		if err := insertStages(r.Context(), tx, wf.ID, req.Stages); err != nil {
			return err
		}
		if wf.Stages, err = loadStages(r.Context(), tx, wf.ID); err != nil {
			return err
		}
		return a.audit(r.Context(), tx, r, "workflow.create", "workflow", wf.ID, req)
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondData(w, http.StatusCreated, wf)
}

func (a *App) listWorkflows(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	p := pagination(r)
	f := &queryFilter{args: []interface{}{id.OrgID}}
	if v := r.URL.Query().Get("status"); v != "" {
		f.add("status = $%d", v)
	}
	if v := r.URL.Query().Get("workflow_type"); v != "" {
		f.add("workflow_type = $%d", v)
	}
	where := `FROM workflows WHERE org_id = $1 AND deleted_at IS NULL` + f.sql()
	var total int
	if err := a.db.QueryRowContext(r.Context(), `SELECT COUNT(*) `+where, f.args...).Scan(&total); err != nil {
		a.fail(w, r, err)
		return
	}
	items, err := queryList(r.Context(), a.db, scanWorkflow,
		`SELECT `+workflowColumns+` `+where+` ORDER BY created_at DESC LIMIT `+limitOffset(p), f.args...)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondPage(w, items, total, p)
}

func (a *App) getWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	wf, err := scanWorkflow(a.db.QueryRowContext(r.Context(), `
		SELECT `+workflowColumns+` FROM workflows WHERE id = $1 AND org_id = $2 AND deleted_at IS NULL`,
		workflowID, identityFrom(r.Context()).OrgID))
	if err != nil {
		a.fail(w, r, notFound(err, "workflow"))
		return
	}
	if wf.Stages, err = loadStages(r.Context(), a.db, wf.ID); err != nil {
		a.fail(w, r, err)
		return
	}
	respondData(w, http.StatusOK, wf)
}

func (a *App) updateWorkflow(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	workflowID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	var req UpdateWorkflowRequest
	if err := decode(r, &req); err != nil {
		a.fail(w, r, err)
		return
	}
	var config interface{}
	if req.Config != nil {
		c, err := jsonObject(req.Config)
		if err != nil {
			a.fail(w, r, err)
			return
		}
		config = c
	}
	var wf *Workflow
	err = a.inTx(r.Context(), func(tx *sql.Tx) error {
		var err error
		wf, err = scanWorkflow(tx.QueryRowContext(r.Context(), `
			UPDATE workflows SET name = COALESCE($3, name), description = COALESCE($4, description),
				status = COALESCE($5, status), config = COALESCE($6, config)
			WHERE id = $1 AND org_id = $2 AND deleted_at IS NULL
			RETURNING `+workflowColumns, workflowID, id.OrgID, req.Name, req.Description, req.Status, config))
		if err != nil {
			return notFound(err, "workflow")
		}
		if req.Stages != nil {
			var used bool
			if err := tx.QueryRowContext(r.Context(), `
				SELECT EXISTS (SELECT 1 FROM approvals a JOIN approval_stages s ON s.id = a.stage_id
				               WHERE s.workflow_id = $1)`, workflowID).Scan(&used); err != nil {
				return err
			}
			if used {
				return NewConflictError("stages cannot be replaced once approvals reference them; create a new workflow")
			}
			if _, err := tx.ExecContext(r.Context(), `DELETE FROM approval_stages WHERE workflow_id = $1`, workflowID); err != nil {
				return err
			}
			if err := insertStages(r.Context(), tx, workflowID, *req.Stages); err != nil {
				return err
			}
		}
		if wf.Stages, err = loadStages(r.Context(), tx, workflowID); err != nil {
			return err
		}
		return a.audit(r.Context(), tx, r, "workflow.update", "workflow", workflowID, req)
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondData(w, http.StatusOK, wf)
}

func (a *App) deleteWorkflow(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	workflowID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	err = a.inTx(r.Context(), func(tx *sql.Tx) error {
		var inReview int
		if err := tx.QueryRowContext(r.Context(), `
			SELECT COUNT(*) FROM content
			WHERE workflow_id = $1 AND org_id = $2 AND status = 'in-review' AND deleted_at IS NULL`,
			workflowID, id.OrgID).Scan(&inReview); err != nil {
			return err
		}
		if inReview > 0 {
			return NewConflictError("content is in review under this workflow")
		}
		if err := execOne(r.Context(), tx, "workflow", `
			UPDATE workflows SET status = 'archived', deleted_at = NOW()
			WHERE id = $1 AND org_id = $2 AND deleted_at IS NULL`, workflowID, id.OrgID); err != nil {
			return err
		}
		return a.audit(r.Context(), tx, r, "workflow.delete", "workflow", workflowID, nil)
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Message: "workflow deleted"})
}

// ============================================================================
// Approvals
//
// Stages are approved in order. Each stage needs required_count "approved"
// decisions from distinct approvers whose role is at least the stage's
// approver_role; authors cannot approve their own content. Any "rejected" or
// "requested-changes" decision sends the content back to draft, and only decisions
// made since the latest submission (content.submitted_at) count.
// ============================================================================

const approvalColumns = `id, content_id, stage_id, approver_id, decision, comment, created_at, updated_at`

func scanApproval(row interface{ Scan(...interface{}) error }) (*Approval, error) {
	ap := &Approval{}
	err := row.Scan(&ap.ID, &ap.ContentID, &ap.StageID, &ap.ApproverID, &ap.Decision, &ap.Comment, &ap.CreatedAt, &ap.UpdatedAt)
	return ap, err
}

type approvalResult struct {
	Approval      *Approval `json:"approval"`
	ContentStatus string    `json:"content_status"`
}

func (a *App) createApproval(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	var req CreateApprovalRequest
	if err := decode(r, &req); err != nil {
		a.fail(w, r, err)
		return
	}
	var result approvalResult
	err := a.inTx(r.Context(), func(tx *sql.Tx) error {
		c, err := getContentRow(r.Context(), tx, id.OrgID, req.ContentID, true)
		if err != nil {
			return err
		}
		if c.Status != "in-review" {
			return NewConflictError("content is " + c.Status + ", not in review")
		}
		if c.CreatedBy == id.UserID {
			return NewForbiddenError("authors cannot approve their own content")
		}
		if c.WorkflowID == nil || activeWorkflow(r.Context(), tx, id.OrgID, *c.WorkflowID) != nil {
			return NewConflictError("the content's workflow is no longer active")
		}
		stages, err := loadStages(r.Context(), tx, *c.WorkflowID)
		if err != nil {
			return err
		}
		approvedCount := func(stageID uuid.UUID) (int, error) {
			var n int
			err := tx.QueryRowContext(r.Context(), `
				SELECT COUNT(DISTINCT approver_id) FROM approvals
				WHERE content_id = $1 AND stage_id = $2 AND decision = 'approved'
				  AND created_at >= COALESCE((SELECT submitted_at FROM content WHERE id = $1), '-infinity')`,
				c.ID, stageID).Scan(&n)
			return n, err
		}

		var stage *ApprovalStage
		for i := range stages {
			s := &stages[i]
			if s.ID == req.StageID {
				stage = s
				break
			}
			n, err := approvedCount(s.ID)
			if err != nil {
				return err
			}
			if n < s.RequiredCount {
				return NewConflictError("stage \"" + s.Name + "\" must be approved first")
			}
		}
		if stage == nil {
			return NewBadRequestError("stage does not belong to the content's workflow")
		}
		if getRolePermissionLevel(id.Role) < getRolePermissionLevel(stage.ApproverRole) {
			return NewForbiddenError("this stage requires the " + stage.ApproverRole + " role")
		}
		n, err := approvedCount(stage.ID)
		if err != nil {
			return err
		}
		if n >= stage.RequiredCount {
			return NewConflictError("stage \"" + stage.Name + "\" is already approved")
		}
		var already bool
		if err := tx.QueryRowContext(r.Context(), `
			SELECT EXISTS (SELECT 1 FROM approvals WHERE content_id = $1 AND stage_id = $2 AND approver_id = $3
			               AND created_at >= COALESCE((SELECT submitted_at FROM content WHERE id = $1), '-infinity'))`,
			c.ID, stage.ID, id.UserID).Scan(&already); err != nil {
			return err
		}
		if already {
			return NewConflictError("you have already decided on this stage")
		}

		result.Approval, err = scanApproval(tx.QueryRowContext(r.Context(), `
			INSERT INTO approvals (content_id, stage_id, approver_id, decision, comment)
			VALUES ($1, $2, $3, $4, $5) RETURNING `+approvalColumns,
			c.ID, stage.ID, id.UserID, req.Decision, nullString(req.Comment)))
		if err != nil {
			return err
		}

		result.ContentStatus = "in-review"
		if req.Decision != "approved" {
			result.ContentStatus = "draft"
		} else {
			done := true
			for _, s := range stages {
				n, err := approvedCount(s.ID)
				if err != nil {
					return err
				}
				if n < s.RequiredCount {
					done = false
					break
				}
			}
			if done {
				result.ContentStatus = "approved"
			}
		}
		if result.ContentStatus != c.Status {
			if _, err := tx.ExecContext(r.Context(),
				`UPDATE content SET status = $2 WHERE id = $1`, c.ID, result.ContentStatus); err != nil {
				return err
			}
		}
		return a.audit(r.Context(), tx, r, "approval.create", "content", c.ID, map[string]interface{}{
			"stage_id": stage.ID, "decision": req.Decision, "content_status": result.ContentStatus,
		})
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondData(w, http.StatusCreated, result)
}

func (a *App) listApprovals(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	contentID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	if _, err := getContentRow(r.Context(), a.db, id.OrgID, contentID, false); err != nil {
		a.fail(w, r, err)
		return
	}
	items, err := queryList(r.Context(), a.db, scanApproval,
		`SELECT `+approvalColumns+` FROM approvals WHERE content_id = $1 ORDER BY created_at`, contentID)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondPage(w, items, len(items), page{Page: 1, Limit: len(items)})
}

// ============================================================================
// Distributions
//
// A distribution records the intent to publish approved or published content to a
// channel. The API stores the request; posting to the channel is done by a worker
// outside this service, which updates status / published_url.
// ============================================================================

const distributionColumns = `id, org_id, content_id, channel, status, scheduled_at, published_url, published_at, metadata, created_at, updated_at`

func scanDistribution(row interface{ Scan(...interface{}) error }) (*Distribution, error) {
	d := &Distribution{}
	err := row.Scan(&d.ID, &d.OrgID, &d.ContentID, &d.Channel, &d.Status, &d.ScheduledAt, &d.PublishedURL,
		&d.PublishedAt, &d.Metadata, &d.CreatedAt, &d.UpdatedAt)
	return d, err
}

func (a *App) createDistribution(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	var req CreateDistributionRequest
	if err := decode(r, &req); err != nil {
		a.fail(w, r, err)
		return
	}
	metadata, err := jsonObject(req.Metadata)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	status := "pending"
	if req.ScheduledAt != nil && req.ScheduledAt.After(time.Now()) {
		status = "scheduled"
	}
	var d *Distribution
	err = a.inTx(r.Context(), func(tx *sql.Tx) error {
		c, err := getContentRow(r.Context(), tx, id.OrgID, req.ContentID, false)
		if err != nil {
			return err
		}
		if c.Status != "approved" && c.Status != "published" {
			return NewConflictError("only approved or published content can be distributed; content is " + c.Status)
		}
		d, err = scanDistribution(tx.QueryRowContext(r.Context(), `
			INSERT INTO distributions (org_id, content_id, channel, status, scheduled_at, metadata)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING `+distributionColumns,
			id.OrgID, c.ID, req.Channel, status, req.ScheduledAt, metadata))
		if err != nil {
			return err
		}
		return a.audit(r.Context(), tx, r, "distribution.create", "distribution", d.ID,
			map[string]interface{}{"content_id": c.ID, "channel": req.Channel})
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondData(w, http.StatusCreated, d)
}

func (a *App) listDistributions(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	contentID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	if _, err := getContentRow(r.Context(), a.db, id.OrgID, contentID, false); err != nil {
		a.fail(w, r, err)
		return
	}
	items, err := queryList(r.Context(), a.db, scanDistribution,
		`SELECT `+distributionColumns+` FROM distributions WHERE content_id = $1 AND org_id = $2 ORDER BY created_at`,
		contentID, id.OrgID)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondPage(w, items, len(items), page{Page: 1, Limit: len(items)})
}

func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func isNoRows(err error) bool { return errors.Is(err, sql.ErrNoRows) }
