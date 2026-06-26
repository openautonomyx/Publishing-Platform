# Multi-Domain Setup for OpenAutonomyX

Support multiple domains pointing to the same platform with separate configurations.

---

## Your Domains

Configure as many as you need:

```
┌─────────────────────────────────────────────────────────┐
│                 openautonomyx.com                       │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  publishing.openautonomyx.com                          │
│  ├─ Main platform                                      │
│  ├─ blog.publishing.openautonomyx.com (WordPress)      │
│  └─ Admin, API, Services                               │
│                                                         │
│  publish.openautonomyx.com                             │
│  ├─ Alternative endpoint                               │
│  ├─ blog.publish.openautonomyx.com (WordPress)         │
│  └─ Same services, different domain                    │
│                                                         │
│  Other subdomains                                       │
│  ├─ app.openautonomyx.com                              │
│  ├─ api.openautonomyx.com                              │
│  └─ Whatever you need                                  │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

---

## Hostinger Setup

### Domain 1: `publishing.openautonomyx.com`

**In Hostinger Dashboard → Domains:**

1. Add subdomain: `publishing`
2. Points to: OpenAutonomyX Platform
3. Create WordPress subdomain: `blog.publishing.openautonomyx.com`

### Domain 2: `publish.openautonomyx.com`

**Same process:**

1. Add subdomain: `publish`
2. Points to: Same OpenAutonomyX Platform (or separate instance)
3. Create WordPress subdomain: `blog.publish.openautonomyx.com`

### Domain 3+: Add More

Add as many subdomains as needed. Each can have:
- Own WordPress blog
- Separate configuration
- Different branding
- Independent analytics

---

## Platform Configuration

### Option A: Single Platform, Multiple Domains

All domains route to same platform, different configurations:

```typescript
// API Gateway router
app.get('/api/v1/config', (req, res) => {
  const domain = req.hostname; // publishing.openautonomyx.com or publish.openautonomyx.com
  
  const configs = {
    'publishing.openautonomyx.com': { /* config 1 */ },
    'publish.openautonomyx.com': { /* config 2 */ },
    'app.openautonomyx.com': { /* config 3 */ }
  };
  
  res.json(configs[domain] || configs['publishing.openautonomyx.com']);
});
```

### Option B: Separate Instances

Run multiple platform instances, each with own database:

```yaml
# Instance 1: publishing.openautonomyx.com
publishing-instance:
  DATABASE_URL: postgresql://...@postgres:5432/publishing_platform_1
  REDIS_URL: redis://redis:6379/0

# Instance 2: publish.openautonomyx.com
publish-instance:
  DATABASE_URL: postgresql://...@postgres:5432/publishing_platform_2
  REDIS_URL: redis://redis:6379/1
```

---

## WordPress Blogs

### Blog 1: blog.publishing.openautonomyx.com

```bash
curl -X POST http://localhost:3010/api/v1/integrations \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Publishing Blog",
    "type": "wordpress",
    "config": {
      "apiUrl": "https://blog.publishing.openautonomyx.com",
      "username": "admin",
      "password": "app-password-1"
    }
  }'
```

### Blog 2: blog.publish.openautonomyx.com

```bash
curl -X POST http://localhost:3010/api/v1/integrations \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Publish Blog",
    "type": "wordpress",
    "config": {
      "apiUrl": "https://blog.publish.openautonomyx.com",
      "username": "admin",
      "password": "app-password-2"
    }
  }'
```

---

## Publishing to Multiple Blogs

Publish same content to all blogs:

```bash
# Get all integrations
INTEGRATIONS=$(curl -s http://localhost:3010/api/v1/integrations)

# Publish to each
for INT in $(echo $INTEGRATIONS | jq -r '.data[].id'); do
  curl -X POST http://localhost:3010/api/v1/integrations/$INT/publish \
    -H "Content-Type: application/json" \
    -d '{
      "title": "My Post",
      "content": "Content...",
      "excerpt": "Summary..."
    }'
done
```

---

## Multi-Tenant Database Schema

Single database supporting multiple domains:

```sql
CREATE TABLE tenants (
  id UUID PRIMARY KEY,
  domain VARCHAR(255) UNIQUE,
  name VARCHAR(255),
  config JSONB,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE content (
  id UUID PRIMARY KEY,
  tenant_id UUID REFERENCES tenants(id),
  title VARCHAR(255),
  content TEXT,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Query by domain
SELECT * FROM content 
WHERE tenant_id = (
  SELECT id FROM tenants WHERE domain = 'publishing.openautonomyx.com'
);
```

---

## API Gateway Multi-Domain Routing

```typescript
import express from 'express';

const app = express();

// Domain-aware middleware
app.use((req, res, next) => {
  req.tenant = {
    domain: req.hostname,
    id: getTenantIdByDomain(req.hostname)
  };
  next();
});

// All routes now have req.tenant context
app.get('/api/v1/content', (req, res) => {
  // Automatically filters by tenant
  const content = db.content.find({ tenant_id: req.tenant.id });
  res.json(content);
});

// Route to appropriate integrations
app.post('/api/v1/integrations/:id/publish', (req, res) => {
  // Uses integrations for this tenant's domain
  const integration = db.integrations.findOne({
    id: req.params.id,
    tenant_id: req.tenant.id
  });
  // ... publish logic
});
```

---

## Environment Setup

Add to your `.env`:

```bash
# Primary domain
PRIMARY_DOMAIN=publishing.openautonomyx.com

# Multi-domain support
ALLOWED_DOMAINS=publishing.openautonomyx.com,publish.openautonomyx.com,app.openautonomyx.com

# Database
DATABASE_URL=postgresql://user:pass@postgres:5432/openautonomyx

# Cache (shared Redis)
REDIS_URL=redis://redis:6379/0

# Integrations service
INTEGRATIONS_URL=http://localhost:3010
```

---

## Docker Compose - Multi-Domain

```yaml
version: '3.9'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: openautonomyx
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  api-gateway:
    build: ./services/api-gateway
    ports:
      - "3000:3000"
    environment:
      DATABASE_URL: postgresql://user:password@postgres:5432/openautonomyx
      REDIS_URL: redis://redis:6379/0
      ALLOWED_DOMAINS: publishing.openautonomyx.com,publish.openautonomyx.com
    depends_on:
      - postgres
      - redis

  integrations:
    build: ./services/integrations
    ports:
      - "3010:3010"
    environment:
      EVENT_BUS_URL: http://event-bus:3001

volumes:
  postgres_data:
```

---

## Hostinger Setup Checklist

- [ ] Domain 1: `publishing.openautonomyx.com`
  - [ ] Subdomain created
  - [ ] WordPress installed on `blog.publishing.openautonomyx.com`
  - [ ] API credentials saved

- [ ] Domain 2: `publish.openautonomyx.com`
  - [ ] Subdomain created
  - [ ] WordPress installed on `blog.publish.openautonomyx.com`
  - [ ] API credentials saved

- [ ] Domain 3+: Additional domains
  - [ ] Subdomains created
  - [ ] WordPress installed
  - [ ] Credentials saved

---

## Register All Integrations

```bash
# Domain 1
curl -X POST http://localhost:3010/api/v1/integrations \
  -d '{"name":"Blog 1","type":"wordpress","config":{"apiUrl":"https://blog.publishing.openautonomyx.com",...}}'

# Domain 2
curl -X POST http://localhost:3010/api/v1/integrations \
  -d '{"name":"Blog 2","type":"wordpress","config":{"apiUrl":"https://blog.publish.openautonomyx.com",...}}'

# Get all
curl http://localhost:3010/api/v1/integrations
```

---

## Multi-Domain Publishing

Create content once, publish to all domains:

```bash
POST /api/v1/content
{
  "title": "Post for all domains",
  "content": "...",
  "publishTo": [
    "blog.publishing.openautonomyx.com",
    "blog.publish.openautonomyx.com",
    "all-integrations"
  ]
}
```

---

## Benefits of Multi-Domain Setup

✅ **Multiple entry points** - Different URLs, same platform
✅ **Separate branding** - Each domain has own look/feel
✅ **Independent blogs** - Multiple WordPress instances
✅ **Centralized management** - One admin panel
✅ **Unified analytics** - Track all domains
✅ **Easy expansion** - Add domains anytime
✅ **Flexible configuration** - Different settings per domain
✅ **Cross-domain publishing** - Publish to multiple blogs at once

---

## Summary

You can have as many domains as you need:
- `publishing.openautonomyx.com` - Main platform
- `publish.openautonomyx.com` - Alternative endpoint
- `app.openautonomyx.com` - App-specific
- `blog.openautonomyx.com` - Blog subdomain
- Any other subdomains you want!

All connected to same OpenAutonomyX platform with separate WordPress blogs for each domain.

**Ready to configure your domains!** 🌐
