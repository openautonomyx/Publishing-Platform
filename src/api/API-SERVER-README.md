# Creative Platform API

Go REST API for multi-tenant content: content with an approval workflow, distribution
requests, agents, API keys and an audit log. PostgreSQL storage; sign-in through
[Logto](https://logto.io) with the same token contract as Identity Fabric.

## Run it

```bash
# 1. Database
createdb creative_platform
psql -d creative_platform -f ../../db/schema.sql

# 2. Configure (see Environment)
export DATABASE_URL=postgresql://dev:dev@localhost:5432/creative_platform
export LOGTO_ENDPOINT=https://auth.example.com
export LOGTO_AUDIENCE=https://api.publishing.example.com

# 3. Start
go run .            # http://localhost:3001
```

The API **refuses to start** without `LOGTO_ENDPOINT` and `LOGTO_AUDIENCE`.

Docker: `docker build -t creative-platform-api .` (from this directory).

## Identity and access

### Callers

- **People** send a Logto access token: `Authorization: Bearer <token>`. It is checked
  against `$LOGTO_ENDPOINT/oidc/jwks` (ES384/ES256/RS256), issuer `$LOGTO_ENDPOINT/oidc`,
  audience `$LOGTO_AUDIENCE` and expiry.
- **Automation** sends an API key: `X-API-Key: pk_…` (see API keys).

### Tenants and users are provisioned on first request

| Token claim | Becomes |
|---|---|
| `$TENANT_CLAIM` (default `tenant_id`), else `DEFAULT_TENANT_ID` | the caller's **organization**, matched on `organizations.external_id`; created on first sign-in |
| `sub` | the caller's **user** in that org, matched on `users.external_id`; created on first sign-in |
| `email` (optional custom claim) | links the sign-in to a user invited with `POST /users` |
| `name` (optional) | display name for new users |
| `scope` | the role (below) |

Every query is scoped to the caller's organization. Another tenant's records return
404, as if they didn't exist.

### Roles

The role comes from the highest `publishing:*` permission in the token's `scope`.
Configure these permissions on the Logto API resource and assign them to Logto roles.

| Scope | Role | Can |
|---|---|---|
| `publishing:read` | viewer | read everything in the org; edit own profile |
| `publishing:edit` | editor | + create/edit/delete content and agents, queue agent runs, request distributions |
| `publishing:approve` | approver | + decide on approval stages |
| `publishing:admin` | admin | + workflows, users, organization settings, API keys, audit log |

A token without a `publishing:*` scope gets 403. Users can be deactivated or removed
locally (`PUT`/`DELETE /users/{id}`), which blocks them even with a valid token.

### API keys

`POST /api-keys` (admin) returns the key **once**; only its SHA-256 hash is stored.
Permissions map to roles (`read`=viewer, `write`=editor, `approve`=approver,
`admin`=admin). A key can't have more permissions than its creator, and it's always
capped by its owner's current role. Removing the owner revokes their keys.

## Endpoints

All under `/api/v1`. Bodies are JSON, and unknown fields are rejected.

| Method & path | Min role | Notes |
|---|---|---|
| `GET /organizations` | viewer | your organization |
| `GET /organizations/{id}` | viewer | |
| `PUT /organizations/{id}` | admin | `name`, `slug`, `description`, `metadata` |
| `DELETE /organizations/{id}` | admin | soft delete; locks every member out |
| `POST /organizations` | — | 403: organizations come from Logto tenants |
| `GET /users` · `GET /users/{id}` | viewer | `?status=` |
| `POST /users` | admin | invite: `email`, `name` (status `invited` until they sign in) |
| `PUT /users/{id}` | viewer | self: `name`, `avatar_url`; admin: also `status` (`active`/`inactive`) |
| `DELETE /users/{id}` | admin | soft delete + revoke their API keys; not yourself |
| `GET /content` · `GET /content/{id}` | viewer | `?status=&content_type=&workflow_id=` |
| `POST /content` | editor | `title`, `content_type`, optional `workflow_id`, `description`, `metadata` |
| `PUT /content/{id}` | editor | edit fields (draft only) and/or `status` (see lifecycle) |
| `DELETE /content/{id}` | editor | soft delete |
| `GET /workflows` · `GET /workflows/{id}` | viewer | `GET /{id}` includes stages |
| `POST /workflows` | admin | `name`, `workflow_type`, `config`, ordered `stages[]` (`name`, `approver_role`, `required_count`, `timeout_hours`) |
| `PUT /workflows/{id}` | admin | `stages` replaces all stages; refused once approvals use them |
| `DELETE /workflows/{id}` | admin | refused while content is in review under it |
| `POST /approvals` | approver | `content_id`, `stage_id`, `decision` (`approved`/`rejected`/`requested-changes`), `comment` |
| `GET /content/{id}/approvals` | viewer | |
| `POST /distributions` | editor | `content_id`, `channel`, optional `scheduled_at`, `metadata`; approved/published content only |
| `GET /content/{id}/distributions` | viewer | |
| `GET /agents` · `GET /agents/{id}` | viewer | |
| `POST /agents` · `PUT /agents/{id}` · `DELETE /agents/{id}` | editor | `status`: `active`/`paused` |
| `POST /agents/{id}/run` | editor | queues a run (202); body optional: `content_id`, `input` |
| `GET /agents/{id}/runs` | viewer | |
| `POST /api-keys` · `GET /api-keys` · `DELETE /api-keys/{id}` | admin | |
| `GET /audit-logs` | admin | `?action=&resource_type=&user_id=&resource_id=` |

Unauthenticated: `GET /health` (liveness) and `GET /ready` (checks the database).

Lists take `?page=` (from 1) and `?limit=` (default 20, max 100).

### Content lifecycle

```
draft ──submit──► in-review ──all stages approved──► approved ──► published ──► archived
  ▲                  │ rejected / requested-changes        │
  └──────────────────┴─────────────────────────────────────┘ (back to draft)
```

- Submitting (`status: in-review`) needs a workflow with at least one stage.
- Stages are approved **in order**. A stage needs `required_count` approvals from
  distinct users whose role is at least its `approver_role` (`approver` or `admin`).
- Authors can't approve their own content.
- Any rejection sends the content back to draft. Only decisions made since the latest
  submission count.
- Content without a workflow can go straight from draft to published.
- Title, description, metadata and workflow can only change in draft.

### Not done by this service

Agent runs and distributions are **recorded** here (status `queued` / `pending` or
`scheduled`). Executing agents and posting to channels belongs to a worker that reads
those rows and updates their status. No such worker exists yet.

## Responses

```jsonc
// single item (200/201/202)
{"success": true, "data": {...}}
// list
{"success": true, "data": [...], "total": 42, "page": 1, "limit": 20}
// error (4xx/5xx)
{"success": false, "error": "VALIDATION_ERROR", "message": "validation failed",
 "details": ["field 'CreateContentRequest.Title' failed validation: required"],
 "request_id": "…", "timestamp": "…"}
```

Error codes: `BAD_REQUEST`, `VALIDATION_ERROR`, `UNAUTHORIZED`, `FORBIDDEN`,
`NOT_FOUND`, `CONFLICT` (e.g. a lifecycle rule), `TOO_MANY_REQUESTS`, `INTERNAL_ERROR`
(details are logged, never returned). Every response carries `X-Request-ID`.

Every create, update and delete is written to `audit_logs` in the same transaction.

## Environment

| Variable | Default | |
|---|---|---|
| `DATABASE_URL` | `postgresql://dev:dev@localhost:5432/creative_platform` | |
| `API_PORT` | `3001` | |
| `LOGTO_ENDPOINT` | — | **required**, e.g. `https://auth.example.com` |
| `LOGTO_AUDIENCE` | — | **required**, the Logto API resource indicator |
| `TENANT_CLAIM` | `tenant_id` | or `organization_id` for Logto organizations |
| `DEFAULT_TENANT_ID` | — | tenant for tokens without the claim |
| `CORS_ORIGINS` | any (`*`) | comma-separated allow-list |
| `RATE_LIMIT_REQUESTS_PER_SECOND` / `RATE_LIMIT_BURST_SIZE` | `100` / `200` | per client IP; `0` disables |
| `TRUST_PROXY` | `false` | `true` only behind a proxy that sets `X-Forwarded-For` (otherwise clients could spoof their IP) |
| `API_ENV` / `API_LOG_LEVEL` | development / `info` | `production` switches to JSON logs |

## Database

`db/schema.sql` is the schema (PostgreSQL 15+). The API enforces tenant isolation in
every query. The schema also enables row-level security, but the API connects as the
table owner, which bypasses RLS, so the policies only constrain other, non-owner roles.

## Tests

```bash
go test ./...                                   # unit tests; integration tests skip
DATABASE_URL=postgresql://… go test -race ./... # + end-to-end tests
```

The end-to-end tests (`app_integration_test.go`) run the real router against
PostgreSQL. Each run loads `db/schema.sql` into its own throwaway schema. A fake Logto
issuer signs real tokens. They cover provisioning, tenant isolation, RBAC, the approval
lifecycle, API keys, invites, agents, audit and validation. CI's Test job provides
`DATABASE_URL`.

## Code layout

| File | |
|---|---|
| `main.go` | config, database pool, server lifecycle |
| `app.go` | router, middleware, identity resolution, RBAC gate, transactions, audit, error mapping |
| `handlers_*.go` | endpoints |
| `models.go` | rows and request types (`JSONB` handles SQL NULL) |
| `rbac.go`, `validation.go`, `errors.go`, `logging.go`, `ratelimit.go` | support |
| `logtoauth/` | Logto token verification, reusable by other Go services |

### Draft files (`-tags draft`)

These files are excluded from the default build with `//go:build draft`. None is
routed, and each depends on code that doesn't exist yet:

| File | Missing / why excluded |
|---|---|
| `deployment_service.go` | `Database` type. Also a second, conflicting deployment design: same `Deployment` types and `/api/v1/deployments` routes as the tracking files, so pick one |
| `deployment_tracking.go`, `deployment_handlers.go` | `Database` type; the conflict above |
| `onboarding_service.go` | `LiferayClient.CreateOrganization/CreateSite/CreateUser`, `sendWelcomeEmail`, `sendSetupCompleteEmail` |
| `liferay_integration.go` | only used by `onboarding_service.go` |
| `workflows_approval.go` | `RecordAuditLog`, `AuditLogEntry` (approvals are implemented in `handlers_content.go`) |
| `publishing.go`, `social_features.go` | read the caller from the client-supplied `X-User-ID` header (must use the verified identity), and their tables exist only in `db/migrations/010_create_publishing_tables.sql` |

To finish one, remove its build tag and make `go build ./...` pass.
`go build -tags draft ./...` shows what's still missing.
