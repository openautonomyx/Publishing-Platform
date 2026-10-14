package main

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"strings"
)

// ============================================================================
// Organizations
//
// A caller belongs to exactly one organization: the one provisioned from their
// token's tenant. Organizations are created on first sign-in, not through the API.
// ============================================================================

const orgColumns = `id, name, slug, external_id, description, tier, status, metadata, created_at, updated_at, created_by`

func scanOrganization(row interface{ Scan(...interface{}) error }) (*Organization, error) {
	o := &Organization{}
	err := row.Scan(&o.ID, &o.Name, &o.Slug, &o.ExternalID, &o.Description, &o.Tier, &o.Status,
		&o.Metadata, &o.CreatedAt, &o.UpdatedAt, &o.CreatedBy)
	return o, err
}

func (a *App) createOrganization(w http.ResponseWriter, r *http.Request) {
	a.fail(w, r, NewForbiddenError("organizations are created automatically on first sign-in from a Logto tenant"))
}

func (a *App) listOrganizations(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	o, err := scanOrganization(a.db.QueryRowContext(r.Context(),
		`SELECT `+orgColumns+` FROM organizations WHERE id = $1 AND deleted_at IS NULL`, id.OrgID))
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondPage(w, []*Organization{o}, 1, page{Page: 1, Limit: 1})
}

func (a *App) getOrganization(w http.ResponseWriter, r *http.Request) {
	orgID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	if orgID != identityFrom(r.Context()).OrgID {
		a.fail(w, r, NewNotFoundError("organization"))
		return
	}
	o, err := scanOrganization(a.db.QueryRowContext(r.Context(),
		`SELECT `+orgColumns+` FROM organizations WHERE id = $1 AND deleted_at IS NULL`, orgID))
	if err != nil {
		a.fail(w, r, notFound(err, "organization"))
		return
	}
	respondData(w, http.StatusOK, o)
}

func (a *App) updateOrganization(w http.ResponseWriter, r *http.Request) {
	orgID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	if orgID != identityFrom(r.Context()).OrgID {
		a.fail(w, r, NewNotFoundError("organization"))
		return
	}
	var req UpdateOrganizationRequest
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

	var o *Organization
	err = a.inTx(r.Context(), func(tx *sql.Tx) error {
		o, err = scanOrganization(tx.QueryRowContext(r.Context(), `
			UPDATE organizations SET
				name = COALESCE($2, name), slug = COALESCE($3, slug),
				description = COALESCE($4, description), metadata = COALESCE($5, metadata)
			WHERE id = $1 AND deleted_at IS NULL
			RETURNING `+orgColumns, orgID, req.Name, req.Slug, req.Description, metadata))
		if err != nil {
			return notFound(err, "organization")
		}
		return a.audit(r.Context(), tx, r, "organization.update", "organization", orgID, req)
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondData(w, http.StatusOK, o)
}

func (a *App) deleteOrganization(w http.ResponseWriter, r *http.Request) {
	orgID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	if orgID != identityFrom(r.Context()).OrgID {
		a.fail(w, r, NewNotFoundError("organization"))
		return
	}
	err = a.inTx(r.Context(), func(tx *sql.Tx) error {
		if err := a.audit(r.Context(), tx, r, "organization.delete", "organization", orgID, nil); err != nil {
			return err
		}
		return execOne(r.Context(), tx, "organization", `
			UPDATE organizations SET status = 'deleted', deleted_at = NOW()
			WHERE id = $1 AND deleted_at IS NULL`, orgID)
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Message: "organization deleted"})
}

// ============================================================================
// Users
//
// Users sign in through Logto and are provisioned on first request. POST /users
// invites someone by email; they are linked when they first sign in with a token
// that carries that email. Roles come from the token's publishing:* scopes, so the
// role stored here mirrors the last token seen and is not editable through the API.
// ============================================================================

const userColumns = `id, org_id, external_id, email, name, avatar_url, role, status, last_login_at, created_at, updated_at`

func scanUser(row interface{ Scan(...interface{}) error }) (*User, error) {
	u := &User{}
	err := row.Scan(&u.ID, &u.OrgID, &u.ExternalID, &u.Email, &u.Name, &u.AvatarURL, &u.Role,
		&u.Status, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	var req CreateUserRequest
	if err := decode(r, &req); err != nil {
		a.fail(w, r, err)
		return
	}
	var u *User
	err := a.inTx(r.Context(), func(tx *sql.Tx) error {
		var err error
		u, err = scanUser(tx.QueryRowContext(r.Context(), `
			INSERT INTO users (org_id, email, name, role, status)
			VALUES ($1, $2, $3, 'viewer', 'invited')
			RETURNING `+userColumns, id.OrgID, strings.TrimSpace(req.Email), req.Name))
		if isUniqueViolation(err, "") {
			return NewConflictError("a user with this email already exists in the organization")
		}
		if err != nil {
			return err
		}
		return a.audit(r.Context(), tx, r, "user.invite", "user", u.ID, map[string]string{"email": req.Email})
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondData(w, http.StatusCreated, u)
}

func (a *App) listUsers(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	p := pagination(r)
	f := &queryFilter{args: []interface{}{id.OrgID}}
	if s := r.URL.Query().Get("status"); s != "" {
		f.add("status = $%d", s)
	}
	where := `FROM users WHERE org_id = $1 AND deleted_at IS NULL` + f.sql()

	var total int
	if err := a.db.QueryRowContext(r.Context(), `SELECT COUNT(*) `+where, f.args...).Scan(&total); err != nil {
		a.fail(w, r, err)
		return
	}
	users, err := queryList(r.Context(), a.db, scanUser,
		`SELECT `+userColumns+` `+where+` ORDER BY created_at LIMIT `+limitOffset(p), f.args...)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondPage(w, users, total, p)
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	userID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	u, err := scanUser(a.db.QueryRowContext(r.Context(), `
		SELECT `+userColumns+` FROM users WHERE id = $1 AND org_id = $2 AND deleted_at IS NULL`,
		userID, identityFrom(r.Context()).OrgID))
	if err != nil {
		a.fail(w, r, notFound(err, "user"))
		return
	}
	respondData(w, http.StatusOK, u)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	userID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	var req UpdateUserRequest
	if err := decode(r, &req); err != nil {
		a.fail(w, r, err)
		return
	}
	isAdmin := canPerformAction(id.Role, PermissionAdmin)
	switch {
	case userID != id.UserID && !isAdmin:
		a.fail(w, r, NewForbiddenError("only admins can update other users"))
		return
	case req.Status != nil && !isAdmin:
		a.fail(w, r, NewForbiddenError("only admins can change a user's status"))
		return
	case req.Status != nil && userID == id.UserID && *req.Status != "active":
		a.fail(w, r, NewBadRequestError("you cannot deactivate yourself"))
		return
	}

	var u *User
	err = a.inTx(r.Context(), func(tx *sql.Tx) error {
		var err error
		u, err = scanUser(tx.QueryRowContext(r.Context(), `
			UPDATE users SET name = COALESCE($3, name), avatar_url = COALESCE($4, avatar_url),
				status = COALESCE($5, status)
			WHERE id = $1 AND org_id = $2 AND deleted_at IS NULL
			RETURNING `+userColumns, userID, id.OrgID, req.Name, req.AvatarURL, req.Status))
		if err != nil {
			return notFound(err, "user")
		}
		return a.audit(r.Context(), tx, r, "user.update", "user", userID, req)
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondData(w, http.StatusOK, u)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := identityFrom(r.Context())
	userID, err := pathID(r)
	if err != nil {
		a.fail(w, r, err)
		return
	}
	if userID == id.UserID {
		a.fail(w, r, NewBadRequestError("you cannot remove yourself"))
		return
	}
	err = a.inTx(r.Context(), func(tx *sql.Tx) error {
		if err := execOne(r.Context(), tx, "user", `
			UPDATE users SET status = 'inactive', deleted_at = NOW()
			WHERE id = $1 AND org_id = $2 AND deleted_at IS NULL`, userID, id.OrgID); err != nil {
			return err
		}
		// A removed user's API keys stop working with them.
		if _, err := tx.ExecContext(r.Context(), `
			UPDATE api_keys SET status = 'revoked' WHERE user_id = $1 AND status = 'active'`, userID); err != nil {
			return err
		}
		return a.audit(r.Context(), tx, r, "user.delete", "user", userID, nil)
	})
	if err != nil {
		a.fail(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Message: "user removed"})
}

// ============================================================================
// Shared query helpers
// ============================================================================

// queryList runs a query and scans every row with scan.
func queryList[T any](ctx context.Context, q querier, scan func(interface{ Scan(...interface{}) error }) (T, error), query string, args ...interface{}) ([]T, error) {
	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []T{}
	for rows.Next() {
		v, err := scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func limitOffset(p page) string {
	return strconv.Itoa(p.Limit) + ` OFFSET ` + strconv.Itoa(p.Offset)
}

// execOne runs an UPDATE/DELETE that must affect exactly one row, else 404.
func execOne(ctx context.Context, q querier, resource, query string, args ...interface{}) error {
	res, err := q.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return NewNotFoundError(resource)
	}
	return nil
}

// notFound turns sql.ErrNoRows into a named 404.
func notFound(err error, resource string) error {
	if isNoRows(err) {
		return NewNotFoundError(resource)
	}
	return err
}
