package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ============================================================================
// Agents and runs
//
// POST /agents/{id}/run queues a run (status "queued"). Executing agents is the job
// of a worker outside this service, which picks up queued runs and records their
// status, output and duration.
// ============================================================================

const agentColumns = `id, org_id, name, agent_type, description, config, status, last_run_at, error_message, created_by, created_at, updated_at`
const runColumns = `id, agent_id, content_id, status, input_data, output_data, error_message, duration_ms, started_at, completed_at`

func scanAgent(row interface{ Scan(...interface{}) error }) (*Agent, error) {
	ag := &Agent{}
	err := row.Scan(&ag.ID, &ag.OrgID, &ag.Name, &ag.AgentType, &ag.Description, &ag.Config, &ag.Status,
		&ag.LastRunAt, &ag.ErrorMessage, &ag.CreatedBy, &ag.CreatedAt, &ag.UpdatedAt)
	return ag, err
}

func scanRun(row interface{ Scan(...interface{}) error }) (*AgentRun, error) {
	run := &AgentRun{}
	err := row.Scan(&run.ID, &run.AgentID, &run.ContentID, &run.Status, &run.InputData, &run.OutputData,
		&run.ErrorMessage, &run.DurationMs, &run.StartedAt, &run.CompletedAt)
	return run, err
}

func (a *App) createAgent(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	var req CreateAgentRequest
	if err := decode(r, &req); err != nil {
		a.fail(w, r, err)
		return
	}
	config, err := jsonObject(req.Config)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	var ag *Agent
	err = a.inTx(r.Context(), func(tx *sql.Tx) error {
		var err error
		ag, err = scanAgent(tx.QueryRowContext(r.Context(), `
			INSERT INTO agents (org_id, name, agent_type, description, config, status, created_by)
			VALUES ($1, $2, $3, $4, $5, 'active', $6) RETURNING `+agentColumns,
			id.OrgID, req.Name, req.AgentType, nullString(req.Description), config, id.UserID))
		if err != nil {
			return err
		}
		return a.audit(r.Context(), tx, r, "agent.create", "agent", ag.ID, map[string]string{"name": req.Name, "agent_type": req.AgentType})
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondData(w, http.StatusCreated, ag)
}

func (a *App) listAgents(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	p := pagination(r)
	f := &queryFilter{args: []interface{}{id.OrgID}}
	if v := r.URL.Query().Get("status"); v != "" {
		f.add("status = $%d", v)
	}
	if v := r.URL.Query().Get("agent_type"); v != "" {
		f.add("agent_type = $%d", v)
	}
	where := `FROM agents WHERE org_id = $1 AND deleted_at IS NULL` + f.sql()
	var total int
	if err := a.db.QueryRowContext(r.Context(), `SELECT COUNT(*) `+where, f.args...).Scan(&total); err != nil {
		a.fail(w, r, err)
		return
	}
	items, err := queryList(r.Context(), a.db, scanAgent,
		`SELECT `+agentColumns+` `+where+` ORDER BY created_at DESC LIMIT `+limitOffset(p), f.args...)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondPage(w, items, total, p)
}

func (a *App) getAgent(w http.ResponseWriter, r *http.Request) {
	agentID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	ag, err := a.agentInOrg(r, a.db, agentID, false)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondData(w, http.StatusOK, ag)
}

func (a *App) agentInOrg(r *http.Request, q querier, agentID uuid.UUID, forUpdate bool) (*Agent, error) {
	query := `SELECT ` + agentColumns + ` FROM agents WHERE id = $1 AND org_id = $2 AND deleted_at IS NULL`
	if forUpdate {
		query += ` FOR UPDATE`
	}
	ag, err := scanAgent(q.QueryRowContext(r.Context(), query, agentID, identityFrom(r.Context()).OrgID))
	return ag, notFound(err, "agent")
}

func (a *App) updateAgent(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	agentID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	var req UpdateAgentRequest
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
	var ag *Agent
	err = a.inTx(r.Context(), func(tx *sql.Tx) error {
		var err error
		ag, err = scanAgent(tx.QueryRowContext(r.Context(), `
			UPDATE agents SET name = COALESCE($3, name), description = COALESCE($4, description),
				status = COALESCE($5, status), config = COALESCE($6, config)
			WHERE id = $1 AND org_id = $2 AND deleted_at IS NULL RETURNING `+agentColumns,
			agentID, id.OrgID, req.Name, req.Description, req.Status, config))
		if err != nil {
			return notFound(err, "agent")
		}
		return a.audit(r.Context(), tx, r, "agent.update", "agent", agentID, req)
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondData(w, http.StatusOK, ag)
}

func (a *App) deleteAgent(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	agentID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	err = a.inTx(r.Context(), func(tx *sql.Tx) error {
		if err := execOne(r.Context(), tx, "agent", `
			UPDATE agents SET status = 'deleted', deleted_at = NOW()
			WHERE id = $1 AND org_id = $2 AND deleted_at IS NULL`, agentID, id.OrgID); err != nil {
			return err
		}
		return a.audit(r.Context(), tx, r, "agent.delete", "agent", agentID, nil)
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Message: "agent deleted"})
}

func (a *App) runAgent(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	agentID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	var req RunAgentRequest
	if r.ContentLength != 0 { // every field is optional, so an empty body is fine
		if err := decode(r, &req); err != nil {
			a.fail(w, r, err)
			return
		}
	}
	input, err := jsonObject(req.Input)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	var run *AgentRun
	err = a.inTx(r.Context(), func(tx *sql.Tx) error {
		ag, err := a.agentInOrg(r, tx, agentID, true)
		if err != nil {
			return err
		}
		if ag.Status != "active" {
			return NewConflictError("agent is " + ag.Status)
		}
		if req.ContentID != nil {
			if _, err := getContentRow(r.Context(), tx, id.OrgID, *req.ContentID, false); err != nil {
				return err
			}
		}
		run, err = scanRun(tx.QueryRowContext(r.Context(), `
			INSERT INTO agent_runs (agent_id, content_id, status, input_data)
			VALUES ($1, $2, 'queued', $3) RETURNING `+runColumns, agentID, req.ContentID, input))
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(r.Context(), `UPDATE agents SET last_run_at = NOW() WHERE id = $1`, agentID); err != nil {
			return err
		}
		return a.audit(r.Context(), tx, r, "agent.run", "agent", agentID, map[string]interface{}{"run_id": run.ID})
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondData(w, http.StatusAccepted, run)
}

func (a *App) listAgentRuns(w http.ResponseWriter, r *http.Request) {
	agentID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	if _, err := a.agentInOrg(r, a.db, agentID, false); err != nil {
		a.fail(w, r, err)
		return
	}
	p := pagination(r)
	var total int
	if err := a.db.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM agent_runs WHERE agent_id = $1`, agentID).Scan(&total); err != nil {
		a.fail(w, r, err)
		return
	}
	items, err := queryList(r.Context(), a.db, scanRun,
		`SELECT `+runColumns+` FROM agent_runs WHERE agent_id = $1 ORDER BY started_at DESC LIMIT `+limitOffset(p), agentID)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondPage(w, items, total, p)
}

// ============================================================================
// API keys
//
// Keys authenticate with the X-API-Key header. Only a SHA-256 hash is stored; the
// key is returned once, on creation. A key's permissions map to a role
// (read=viewer, write=editor, approve=approver, admin=admin), capped by the role of
// the user who owns it, and a key can't be created above the creator's own role.
// ============================================================================

const apiKeyColumns = `id, org_id, user_id, key_prefix, name, permissions, status, last_used_at, expires_at, created_at`

func scanAPIKey(row interface{ Scan(...interface{}) error }) (*APIKey, error) {
	k := &APIKey{}
	err := row.Scan(&k.ID, &k.OrgID, &k.UserID, &k.KeyPrefix, &k.Name, &k.Permissions, &k.Status,
		&k.LastUsedAt, &k.ExpiresAt, &k.CreatedAt)
	return k, err
}

func (a *App) createAPIKey(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	var req CreateAPIKeyRequest
	if err := decode(r, &req); err != nil {
		a.fail(w, r, err)
		return
	}
	if getRolePermissionLevel(roleForPermissions(req.Permissions)) > getRolePermissionLevel(id.Role) {
		a.fail(w, r, NewForbiddenError("a key cannot have more permissions than its creator"))
		return
	}
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		a.fail(w, r, err)
		return
	}
	key := "pk_" + base64.RawURLEncoding.EncodeToString(secret)
	var expires *time.Time
	if req.ExpiresInDays > 0 {
		t := time.Now().Add(time.Duration(req.ExpiresInDays) * 24 * time.Hour)
		expires = &t
	}

	var k *APIKey
	err := a.inTx(r.Context(), func(tx *sql.Tx) error {
		var err error
		k, err = scanAPIKey(tx.QueryRowContext(r.Context(), `
			INSERT INTO api_keys (org_id, user_id, key_hash, key_prefix, name, permissions, expires_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING `+apiKeyColumns,
			id.OrgID, id.UserID, hashAPIKey(key), key[:11], req.Name, pq.StringArray(req.Permissions), expires))
		if err != nil {
			return err
		}
		return a.audit(r.Context(), tx, r, "api_key.create", "api_key", k.ID,
			map[string]interface{}{"name": req.Name, "permissions": req.Permissions})
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	k.Key = key
	respondData(w, http.StatusCreated, k)
}

func (a *App) listAPIKeys(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	p := pagination(r)
	var total int
	if err := a.db.QueryRowContext(r.Context(),
		`SELECT COUNT(*) FROM api_keys WHERE org_id = $1 AND deleted_at IS NULL`, id.OrgID).Scan(&total); err != nil {
		a.fail(w, r, err)
		return
	}
	items, err := queryList(r.Context(), a.db, scanAPIKey, `
		SELECT `+apiKeyColumns+` FROM api_keys WHERE org_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT `+limitOffset(p), id.OrgID)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondPage(w, items, total, p)
}

func (a *App) revokeAPIKey(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	keyID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	err = a.inTx(r.Context(), func(tx *sql.Tx) error {
		if err := execOne(r.Context(), tx, "api key", `
			UPDATE api_keys SET status = 'revoked' WHERE id = $1 AND org_id = $2 AND status = 'active'`,
			keyID, id.OrgID); err != nil {
			return err
		}
		return a.audit(r.Context(), tx, r, "api_key.revoke", "api_key", keyID, nil)
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Message: "API key revoked"})
}

// ============================================================================
// Audit logs
// ============================================================================

const auditColumns = `id, org_id, user_id, action, resource_type, resource_id, changes, host(ip_address), user_agent, created_at`

func scanAudit(row interface{ Scan(...interface{}) error }) (*AuditLog, error) {
	l := &AuditLog{}
	err := row.Scan(&l.ID, &l.OrgID, &l.UserID, &l.Action, &l.ResourceType, &l.ResourceID, &l.Changes,
		&l.IPAddress, &l.UserAgent, &l.CreatedAt)
	return l, err
}

func (a *App) listAuditLogs(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	p := pagination(r)
	q := r.URL.Query()
	f := &queryFilter{args: []interface{}{id.OrgID}}
	if v := q.Get("action"); v != "" {
		f.add("action = $%d", v)
	}
	if v := q.Get("resource_type"); v != "" {
		f.add("resource_type = $%d", v)
	}
	for _, key := range []string{"user_id", "resource_id"} {
		if v := q.Get(key); v != "" {
			u, err := uuid.Parse(v)
			if err != nil {
				a.fail(w, r, NewBadRequestError("invalid "+key))
				return
			}
			f.add(key+" = $%d", u)
		}
	}
	where := `FROM audit_logs WHERE org_id = $1` + f.sql()
	var total int
	if err := a.db.QueryRowContext(r.Context(), `SELECT COUNT(*) `+where, f.args...).Scan(&total); err != nil {
		a.fail(w, r, err)
		return
	}
	items, err := queryList(r.Context(), a.db, scanAudit,
		`SELECT `+auditColumns+` `+where+` ORDER BY created_at DESC, id LIMIT `+limitOffset(p), f.args...)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondPage(w, items, total, p)
}
