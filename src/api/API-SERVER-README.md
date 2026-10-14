# Universal Creative Platform - API Server

Production-grade REST API server in Go for the universal creative platform.

## Overview

The API server is the core backend of the creative platform, handling:

- **Content Management** — Create, read, update, version content
- **Approval Workflows** — Multi-stage governance and compliance
- **Ingestion** — Accept data from any source (databases, APIs, files)
- **Distribution** — Publish to multiple destinations globally
- **Analytics** — Track metrics and engagement
- **Multi-Tenancy** — Isolated environments per enterprise
- **Security** — Authentication, authorization, audit trails
- **Agents** — API-native support for autonomous systems

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  HTTP Handler Layer (Gin)                                   │
│  ├─ Content Handlers                                        │
│  ├─ Approval Handlers                                       │
│  ├─ Ingestion Handlers                                      │
│  ├─ Analytics Handlers                                      │
│  └─ Distribution Handlers                                   │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│  Service Layer (Business Logic)                             │
│  ├─ ContentService        (CRUD, versioning, publishing)   │
│  ├─ ApprovalService       (Workflows, decisions, audit)     │
│  ├─ IngestionService      (Data ingestion, adapters)        │
│  ├─ AnalyticsService      (Metrics, reporting)              │
│  └─ AuthService           (Token validation)                │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│  Middleware Layer                                           │
│  ├─ Auth Middleware       (JWT validation)                  │
│  ├─ RBAC Middleware       (Role-based access)               │
│  ├─ Tenant Middleware     (Tenant isolation)                │
│  ├─ Trace Middleware      (Request tracing)                 │
│  └─ CORS Middleware                                         │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│  Database Layer (PostgreSQL)                                │
│  ├─ Content tables                                          │
│  ├─ Approval tables                                         │
│  ├─ Audit log tables (immutable)                            │
│  ├─ Ingestion tables                                        │
│  └─ Metrics tables                                          │
└─────────────────────────────────────────────────────────────┘
```

## File Structure

```
src/
├── api-main.go              # Server initialization, routing
├── api-handlers.go          # HTTP handlers
├── api-services.go          # Business logic layer
├── api-middleware-types.go  # Middleware + request/response types
├── go.mod                   # Go dependencies
├── Dockerfile.api           # Container image
└── .env.example             # Environment template
```

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Docker (optional)

### Local Development

```bash
# Set up environment
cp .env.example .env
# Edit .env with your settings

# Install dependencies
go mod download

# Run migrations
psql -h localhost -U dev -d creative_platform < ../01-database-schema.sql

# Run server
go run .

# Server runs on http://localhost:3001
```

### With Docker

```bash
# Build image
docker build -f Dockerfile.api -t creative-platform-api:latest .

# Run container
docker run -p 3001:3001 \
  -e DATABASE_URL=postgresql://dev:dev@postgres:5432/creative_platform \
  -e API_PORT=3001 \
  creative-platform-api:latest
```

### With Docker Compose

```bash
# Start entire stack (from root directory)
docker-compose up -d

# API runs on http://localhost:3001
```

## API Routes

### Content Management

```
POST   /api/v1/content                # Create content
GET    /api/v1/content/:id            # Get content
GET    /api/v1/content                # List content
PUT    /api/v1/content/:id            # Update content
DELETE /api/v1/content/:id            # Archive content
GET    /api/v1/content/:id/versions   # Get version history
POST   /api/v1/content/:id/versions/:vid/revert  # Revert version
```

### Approvals & Governance

```
POST   /api/v1/approvals              # Submit for approval
GET    /api/v1/approvals              # List pending approvals
GET    /api/v1/approvals/:id          # Get approval status
POST   /api/v1/approvals/:id/approve  # Approve content
POST   /api/v1/approvals/:id/reject   # Reject content
GET    /api/v1/approvals/:id/audit-trail  # Audit trail
```

### Content Types

```
GET    /api/v1/content-types          # List available types
POST   /api/v1/content-types          # Create custom type
GET    /api/v1/content-types/:id      # Get type schema
PUT    /api/v1/content-types/:id      # Update type
```

### Ingestion

```
POST   /api/v1/ingest/adapters        # Create adapter
GET    /api/v1/ingest/adapters        # List adapters
GET    /api/v1/ingest/adapters/:id    # Get adapter config
POST   /api/v1/ingest/adapters/:id/sync  # Trigger sync
GET    /api/v1/ingest/jobs/:id        # Get job status
```

### Distribution & Publishing

```
GET    /api/v1/content/:id/distribution    # Distribution status
POST   /api/v1/content/:id/publish         # Publish content
POST   /api/v1/content/:id/republish       # Republish to edges
```

### Analytics

```
GET    /api/v1/analytics/content/:id  # Content metrics
GET    /api/v1/analytics/agent        # Agent/creator metrics
GET    /api/v1/analytics/system       # System-wide metrics
GET    /api/v1/audit-logs             # Query audit logs
```

### Webhooks

```
POST   /api/v1/webhooks/events        # Incoming webhook
GET    /api/v1/webhooks/test          # Test webhook
```

## Authentication

Every `/api/v1` route requires a **Logto access token** for the API resource named by
`LOGTO_AUDIENCE`. The middleware (`logtoauth/`) checks it against the issuer's JWKS at
`$LOGTO_ENDPOINT/oidc/jwks`: signature (ES384, ES256 or RS256), issuer
`$LOGTO_ENDPOINT/oidc`, audience, and expiry. It uses the same contract as Identity Fabric,
so one Logto setup serves both services.

| From the token | Used as |
|---|---|
| `sub` | user ID (`user_id` in the request context) |
| `$TENANT_CLAIM` (default `tenant_id`), else `DEFAULT_TENANT_ID` | org/tenant ID; also sets the PostgreSQL RLS org |
| highest of `publishing:admin` / `publishing:approve` / `publishing:edit` / `publishing:read` in `scope` | role: admin / approver / editor / viewer |

The tenant always comes from the token. Clients cannot pick it with a header.

| Response | When |
|---|---|
| `401` | no bearer token, or the token fails verification |
| `403` | valid token but no tenant, or no `publishing:*` scope |

**Logto setup:** create an API resource (its indicator is `LOGTO_AUDIENCE`) with the four
`publishing:*` permissions, assign them to roles, and add a custom JWT claim `tenant_id`
for access tokens (or set `TENANT_CLAIM=organization_id` to use Logto organizations).

The API **refuses to start** if `LOGTO_ENDPOINT` or `LOGTO_AUDIENCE` is unset.

```bash
Authorization: Bearer <logto_access_token>
X-Request-ID: unique-request-id   # optional
```

### Example Request

```bash
curl -X POST http://localhost:3001/api/v1/content \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "content_type_id": "blog-post",
    "payload": {
      "title": "My Blog Post",
      "body": "Content here...",
      "tags": ["tech", "design"]
    }
  }'
```

## Response Format

### Success Response (200, 201)

```json
{
  "id": "content-123",
  "content_type_id": "blog-post",
  "status": "draft",
  "payload": { ... },
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### List Response (200)

```json
{
  "data": [ ... ],
  "total": 100,
  "limit": 20,
  "offset": 0
}
```

### Error Response (4xx, 5xx)

```json
{
  "error": "Error message",
  "request_id": "req-123"
}
```

## Environment Variables

```env
# Server
API_PORT=3001
ENVIRONMENT=development
LOG_LEVEL=info

# Database
DATABASE_URL=postgresql://dev:dev@postgres:5432/creative_platform
DATABASE_ENGINE=postgres

# Caching
REDIS_URL=redis://redis:6379

# Message Queue
KAFKA_BROKERS=kafka:29092

# Observability
JAEGER_AGENT_HOST=jaeger
JAEGER_AGENT_PORT=6831
PROMETHEUS_PORT=9091

# Security (Logto; see Authentication)
LOGTO_ENDPOINT=https://auth.example.com
LOGTO_AUDIENCE=https://api.publishing.example.com
TENANT_CLAIM=tenant_id
DEFAULT_TENANT_ID=
TENANT_ISOLATION_LEVEL=row
```

## Development

### Running Tests

```bash
go test ./...
go test -v ./...
go test -cover ./...
```

### Draft files (`-tags draft`)

These files are excluded from the default build with `//go:build draft`. None of them
is routed from `main.go`, and each depends on code that doesn't exist yet:

| File | Missing |
|---|---|
| `deployment_service.go` | `Database` type. Also conflicts with `deployment_tracking.go`: both define `Deployment`, `DeploymentPhase`, `DeploymentStatus` and the same `/api/v1/deployments` routes, so pick one design |
| `deployment_tracking.go`, `deployment_handlers.go` | `Database` type; the `deployment_service.go` conflict above |
| `onboarding_service.go` | `LiferayClient.CreateOrganization/CreateSite/CreateUser`, `sendWelcomeEmail`, `sendSetupCompleteEmail` |
| `workflows_approval.go` | `RecordAuditLog`, `AuditLogEntry` |

To finish one, remove its build tag and make `go build ./...` pass.
`go build -tags draft ./...` shows what's still missing.

### Code Organization

- **Handlers** — HTTP request/response handling (api-handlers.go)
- **Services** — Business logic and domain models (api-services.go)
- **Middleware** — Cross-cutting concerns (api-middleware-types.go)
- **Models** — Data structures (models in separate file)
- **Database** — SQL queries and migrations (in services)

### Adding New Endpoints

1. Create handler in `api-handlers.go`
2. Add service method in `api-services.go`
3. Add route in `setupRouter()` in `api-main.go`
4. Add request/response types in `api-middleware-types.go`
5. Add database schema changes if needed

### Error Handling

```go
// Standard error response
if err != nil {
  log.Printf("Error: %v", err)
  c.JSON(http.StatusInternalServerError, gin.H{
    "error": "Human-readable error message",
  })
  return
}
```

## Deployment

### Docker Compose (Development)

```bash
docker-compose up api
```

### Kubernetes (Production)

```bash
kubectl apply -k k8s/overlays/production
```

### Terraform (Infrastructure)

```bash
cd infra
tofu init
tofu apply -var-file=../environments/prod.tfvars
```

## Monitoring

### Health Check

```bash
curl http://localhost:3001/health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Metrics

Prometheus metrics available at `:9091/metrics`

```bash
curl http://localhost:9091/metrics
```

### Logging

Logs output to stdout (configurable via LOG_LEVEL env var)

### Tracing

Distributed traces sent to Jaeger at `jaeger:6831`

## Security Considerations

### Authentication

- [x] Verify Logto access tokens (JWKS signature, issuer, audience, expiry)
- [x] Token refresh and expiry handled by Logto (short-lived access tokens)
- [ ] Add rate limiting per API key

### Authorization

- [ ] Implement full RBAC from database
- [ ] Add permission checking on all endpoints
- [ ] Implement attribute-based access control

### Data Protection

- [ ] Encrypt sensitive fields at rest
- [ ] Implement TLS for all connections
- [ ] Add request/response logging (sanitized)
- [ ] Implement PII detection and masking

### Compliance

- [ ] Add signature verification for audit logs
- [ ] Implement hash chain for audit log integrity
- [ ] Add compliance event logging
- [ ] Implement GDPR data export/deletion

## Performance Optimization

### Database

- [ ] Add connection pooling (implemented)
- [ ] Add query result caching (Redis)
- [ ] Add database query indexes
- [ ] Implement prepared statements

### Caching

- [ ] Cache content types in Redis
- [ ] Cache user permissions in Redis
- [ ] Implement cache invalidation strategy

### API

- [ ] Add response compression (gzip)
- [ ] Implement pagination for list endpoints
- [ ] Add request throttling
- [ ] Implement async processing for long operations

## Troubleshooting

### Server won't start

```bash
# Check if port 3001 is already in use
lsof -i :3001

# Check database connection
psql $DATABASE_URL -c "SELECT version();"
```

### Database connection errors

```bash
# Check PostgreSQL is running
docker-compose logs postgres

# Check connection string
echo $DATABASE_URL

# Test connection
psql $DATABASE_URL -c "SELECT 1;"
```

### Request failures

```bash
# Check logs
docker-compose logs api

# Test endpoint
curl -v http://localhost:3001/health
```

## Next Steps

1. **Implement JWT validation** — Use `github.com/golang-jwt/jwt`
2. **Add request validation** — Use `github.com/go-playground/validator`
3. **Implement database query builder** — Use `sqlc` for type-safe queries
4. **Add more middleware** — Rate limiting, compression, etc.
5. **Implement async job processing** — Use Temporal or Kafka consumers
6. **Add comprehensive logging** — Use structured logging (Zap or Logrus)
7. **Add API documentation** — Swagger/OpenAPI with `swaggo`
8. **Implement caching strategy** — Redis for performance
9. **Add end-to-end tests** — Integration tests with test database
10. **Setup CI/CD pipeline** — GitHub Actions or similar

## Contributing

This is a draft implementation. Contributions welcome!

- Report issues in the GitHub issue tracker
- Submit PRs with improvements
- Test thoroughly before submitting
- Follow Go best practices and conventions

## License

Apache 2.0
