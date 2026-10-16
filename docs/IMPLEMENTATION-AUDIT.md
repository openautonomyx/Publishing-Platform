# Implementation audit — 2026-10-16

Baseline: `79ac2d54018989a246f21586805062421e5884c1` (`main`). This is a source and local execution audit, not production certification.

## Verified evidence and release gates

- Issue #6 remains open. Run https://github.com/openautonomyx/Publishing-Platform/actions/runs/38863532312 has `startup_failure`; its jobs API returns `total_count: 0`. No test result can be inferred from that run.
- The existing investigation at https://github.com/openautonomyx/Publishing-Platform/issues/6#issuecomment-6209826747 reports that even a no-actions probe cannot start. Its account/repository restriction hypothesis is not an administrator-confirmed root cause.
- The administrator must inspect the banner on https://github.com/openautonomyx/Publishing-Platform/actions/runs/38863157242, repository/organization Actions settings, account restrictions and billing, and escalate with the run IDs if needed. Successful executed jobs are required before closing #6. These settings were not changed here.
- `node --test workspace-api/server.test.cjs`: 2 passed on the baseline (Node 24.19.0). This tests the single-operator evaluation API only.
- `go test -race ./...` and `go vet ./...` in `src/api`: pass with Go 1.26.2. PostgreSQL integration cases skip without `DATABASE_URL`; these are not database integration passes.
- Supplying an unreachable `DATABASE_URL` now causes those tests to fail instead of skip. A negative control using port 1 confirmed failures in all nine database integration cases.
- Baseline frontend build fails: missing `frontend/tsconfig.node.json`. Existing PR #11 supplies this and removes unused imports. An isolated checkout of its head `75ae0aa` built successfully against the installed frontend dependencies. This PR does not duplicate #11; merge/validate it before expecting the workspace pipeline to pass.
- All workflow files pass `actionlint v1.7.7 -shellcheck= -pyflakes=` after this change. ShellCheck and Pyflakes were unavailable and are explicitly disabled; this is workflow/expression/action metadata validation, not shell validation or a successful hosted run.
- Compose configuration validates with a local dummy token. Container build/runtime validation remains UNVERIFIED: the Docker client exists, but no daemon socket is available. No production deployment was attempted. An isolated PostgreSQL binary was also attempted, but the host sandbox rejected its required shared-memory creation (`shmget: Operation not permitted`); database execution remains unverified.
- Root SDK `npm run build`, `npm test -- --runInBand`, and `npm run lint` fail on the baseline: missing `src/client.ts`, SDK compilation includes incomplete React pages, agent SDK references missing `setupInterceptors`, LLM typing errors, no Jest tests/configuration, duplicate package names, and no ESLint configuration. Do not suppress these checks or use `--passWithNoTests`.

## Implementation inventory

| Area | Classification | Source evidence and remaining work |
|---|---|---|
| Evaluation editorial API | IMPLEMENTED, limited | `workspace-api/server.cjs` and its HTTP tests cover draft/review/approval/internal publication with a shared token and atomic file replacement. No multi-user identity, DB transactions, external delivery or multi-process safety. |
| PostgreSQL API | PARTIALLY IMPLEMENTED / UNVERIFIED integration | `src/api/app.go`, `handlers_content.go`, `handlers_org_users.go`, `handlers_agents_keys.go` route real APIs. Database integration tests exist in `app_integration_test.go`; must run against real PostgreSQL. |
| Identity and authorization | PARTIALLY IMPLEMENTED | `src/api/logtoauth`, `rbac.go`, `app.go`: OIDC verification, tenant identity and RBAC are present. Do not classify all authentication as missing. Live issuer validation, onboarding and frontend session integration still need end-to-end validation. Workspace shared token is evaluation-only. |
| Database / migrations | PARTIALLY IMPLEMENTED | `db/schema.sql`, `db/migrations`, Go transaction helper and `FOR UPDATE` content locking exist. Migration ordering, schema compatibility, production upgrades, backup/restore and retention remain unverified. |
| Frontend | PARTIALLY IMPLEMENTED / BROKEN baseline build | `/content` uses `/api/publications` shared-token API, not the Go `/api/v1/content` API. PR #11 fixes its build. Other pages show static sample data (including “Connected” WordPress in `IntegrationsPage.tsx`) and do not prove real provider connectivity. |
| SDK | BROKEN | `src/index.ts` exports nonexistent `./client`; service URLs use `/publications` while Go routes `/content`. Root build includes unrelated incomplete UI and utilities. SDK/backend contract must be reconciled and tested. |
| Publishing orchestration | PARTIALLY IMPLEMENTED | Workspace timer publishes internally; it has no provider delivery receipt, durable queue/leases, retries, idempotency or cancellation. Several Go publishing/workflow files have `//go:build draft` and are excluded from the runnable server. |
| Provider integrations | UNVERIFIED | Adapter/prototype files and static UI are not evidence of live WordPress, LinkedIn, X, Facebook, Instagram or YouTube publishing. Need supported contracts, credentials, mock-provider tests and live verification separately. |
| AI / multilingual | PARTIALLY IMPLEMENTED / UNVERIFIED | `src/services/llm.ts` has a local-model path and explicitly unimplemented cloud provider path. Registries include internal `zip` identifier. PR #13 contains registry regression tests. No model quality or live AI claims established. |
| Analytics / collaboration | PARTIALLY IMPLEMENTED | API organizational resources exist; frontend metrics and SDK service methods do not establish provider ingestion or end-to-end team workflows. |
| Security / operations | PARTIALLY IMPLEMENTED | Token comparison, Go RBAC, security headers, rate limiting and audits exist. Tenant regression tests require PostgreSQL. Dependency remediation, secrets encryption, restore drills and monitoring validation remain release gates. |
| Containers / release | BROKEN / UNVERIFIED | Root `Dockerfile` requires an ignored lockfile, copies absent optional files with invalid COPY fallback syntax, excludes README/LICENSE via `.dockerignore`, and starts an SDK entrypoint as a web server. API image runs as root. Workspace image runtime is unverified. |
| CI | BROKEN externally; code defects repaired in this PR | Zero executed hosted jobs; invalid top-level `secrets`, invalid secret expression contexts, retired action runtimes, missing Docker build step ID and masked smoke failures were independently found. Actual build/test failures above remain visible. |

## Changes in this stabilization PR

1. Remove invalid workflow syntax and invalid secret contexts, update retired action runtimes, repair the Docker digest step ID, and add pinned actionlint validation for all workflows.
2. Disable automatic VPS deployment on push; deployment jobs require manual dispatch and reference `production`. An administrator must configure required reviewers on that environment; a YAML environment name alone does not enforce approval. No deployment was dispatched.
3. Preserve existing checks and make container smoke validation exercise the real web/API authentication and full internal editorial lifecycle. Requested Docker validation fails when tooling/daemon is unavailable. Smoke state uses an isolated Compose project and its volumes are removed afterward.
4. Commit the frontend lockfile and use `npm ci` in workspace validation/container builds for repeatability.
5. Fail database integration validation when a provided database cannot be reached.

## Ordered remaining production plan

1. Administrator resolves #6; run new workflow validation, Go integration tests, PR #11 frontend build and Compose lifecycle checks. Review PR #13 separately. Do not merge or deploy without human approval.
2. Restore a coherent SDK package and contract tests; repair root container packaging. Inventory dependency/security findings with scanner evidence.
3. Choose and document the primary application API (existing transactional Go API is the natural candidate); connect the editorial UI to authenticated tenant-scoped routes. Do not introduce a competing database.
4. Validate migrations and implement content versions/concurrency, durable publication jobs/leases/retries/cancellation and audit records. Preserve separate internal and external delivery states.
5. Add provider adapters incrementally with mock contract tests, OAuth/token security and receipts. Mark credentials/approval-dependent live validation as unverified.
6. Connect media, AI, multilingual workflow, team administration and real analytics; replace static connected/status/metric claims.
7. Validate security, dependency and container scans, restore drills, health/metrics, staging smoke and rollback. Obtain approval for a verified production rollout.

The full production plan remains incomplete. This audit deliberately records runnable behavior separately from design documents and draft code.
