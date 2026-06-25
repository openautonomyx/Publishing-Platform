# Publishing Platform - Modular Architecture

## 🏗️ Platform Structure

```
Publishing-Platform (Orchestrator)
├── @publishing-platform/core (Shared Library)
│   ├── Data Models
│   ├── Auth System
│   ├── Event Bus
│   └── Types
├── api-gateway (Router & Load Balancer)
├── content-management (Module 1)
├── analytics (Module 2)
├── optimization (Module 3)
├── design (Module 4)
├── feature-store (Module 5)
└── skills (Module 6)
```

---

## 📦 Shared Core Library

**@publishing-platform/core** (npm package)

```typescript
// types/index.ts
export interface GartnerCategory {
  id: string;
  marketName: string;
  quadrant: 'leader' | 'visionary' | 'challenger' | 'niche';
  industry: string;
  vertical: string;
  vendor?: string;
  score?: number;
  capabilities: string[];
  tags: string[];
}

export interface Content {
  id: string;
  title: string;
  type: 'blog' | 'whitepaper' | 'case-study' | 'guide';
  gartnerCategory: GartnerCategory;
  createdAt: Date;
  updatedAt: Date;
}

export interface Skill {
  id: string;
  name: string;
  category: GartnerCategory;
  level: 'beginner' | 'intermediate' | 'advanced' | 'expert';
  prerequisite: string[];
}

export interface User {
  id: string;
  email: string;
  name: string;
  roles: Role[];
  skills: Skill[];
  permissions: Permission[];
}

export type Role = 'admin' | 'creator' | 'analyst' | 'designer' | 'viewer';
export type Permission = string;
```

---

## 🔐 Authentication & Authorization

**Shared Auth Module:**
```typescript
// auth/index.ts
export interface JWTPayload {
  userId: string;
  email: string;
  roles: Role[];
  permissions: Permission[];
  iat: number;
  exp: number;
}

export class AuthService {
  generateToken(user: User): string
  verifyToken(token: string): JWTPayload
  validatePermission(token: JWTPayload, required: Permission): boolean
}
```

---

## 🔌 Event Bus (Service Communication)

**Pub/Sub Pattern:**
```typescript
// events/index.ts
export interface Event {
  id: string;
  type: 'content.created' | 'skill.updated' | 'analysis.complete' | ...;
  timestamp: Date;
  source: string; // which module
  data: any;
}

export class EventBus {
  subscribe(eventType: string, handler: (event: Event) => void)
  publish(event: Event)
  unsubscribe(eventType: string, handler)
}
```

---

## 🌐 API Gateway

**Central entry point for all modules:**

```typescript
// gateway/index.ts
app.use('/api/v1/content', contentModuleRouter);
app.use('/api/v1/analytics', analyticsModuleRouter);
app.use('/api/v1/design', designModuleRouter);
app.use('/api/v1/skills', skillsModuleRouter);
app.use('/api/v1/optimization', optimizationModuleRouter);
app.use('/api/v1/features', featureStoreRouter);
app.use('/api/v1/gartner', gartnerCategoriesRouter);

// Middleware
app.use(authMiddleware);
app.use(rateLimitMiddleware);
app.use(loggingMiddleware);
```

---

## 📊 Module Structure (Template)

Each module follows this pattern:

```
content-management/
├── src/
│   ├── api/
│   │   └── routes.ts
│   ├── services/
│   │   └── contentService.ts
│   ├── models/
│   │   └── contentModel.ts
│   ├── events/
│   │   └── contentEvents.ts
│   └── index.ts
├── package.json
├── Dockerfile
├── k8s-deployment.yaml
└── README.md
```

---

## 📁 Database Schema

**Shared Tables:**
```sql
-- Users
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email VARCHAR(255) UNIQUE,
  name VARCHAR(255),
  roles TEXT[],
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- Gartner Categories
CREATE TABLE gartner_categories (
  id UUID PRIMARY KEY,
  market_name VARCHAR(255),
  quadrant VARCHAR(50),
  industry VARCHAR(255),
  vendor VARCHAR(255),
  score FLOAT,
  capabilities TEXT[],
  tags TEXT[]
);

-- Module-specific tables
-- content.* (content management)
-- analytics.* (analytics)
-- skills.* (skills)
-- etc.
```

---

## 🚀 Deployment (K3s)

**Each module deploys as separate service:**

```yaml
---
# content-management deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: content-management
spec:
  replicas: 2
  template:
    spec:
      containers:
      - name: content-management
        image: ghcr.io/openautonomyx/content-management:latest
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: url
        - name: EVENT_BUS_URL
          value: "http://event-bus:3001"
---
# analytics deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: analytics
spec:
  replicas: 2
  # ... similar config
```

---

## 🔄 Module Communication Flow

```
User Request
    ↓
API Gateway (routing, auth, rate limit)
    ↓
Module API (e.g., content-management)
    ↓
Module Service Logic
    ↓
Database Query
    ↓
Publish Event (e.g., "content.created")
    ↓
Event Bus
    ↓
Other Modules Subscribe (e.g., analytics listens for content.created)
    ↓
Update Their State
```

---

## 📦 Repositories to Create

1. **@publishing-platform/core** - Shared library (npm)
2. **publishing-platform-gateway** - API Gateway
3. **publishing-platform-content** - Content Management
4. **publishing-platform-analytics** - Analytics
5. **publishing-platform-optimization** - Optimization
6. **publishing-platform-design** - Design
7. **publishing-platform-features** - Feature Store
8. **publishing-platform-skills** - Skills
9. **publishing-platform-orchestrator** - K3s deployment & docs

---

## 🔐 Security

- ✅ JWT authentication (shared core)
- ✅ Role-based access control (RBAC)
- ✅ API key for service-to-service
- ✅ Request signing (modules verify origin)
- ✅ Rate limiting (API Gateway)
- ✅ Encrypted secrets (K3s)

---

## 📈 Scaling

Each module scales independently:
```bash
# Scale content-management to 5 replicas
kubectl scale deployment content-management --replicas=5

# Scale analytics to 3 replicas
kubectl scale deployment analytics --replicas=3
```

---

## 🎯 Development Workflow

1. **Clone all repos locally**
2. **Install @publishing-platform/core locally**
3. **Run each module with `npm start`**
4. **Event Bus runs on localhost:3001**
5. **API Gateway on localhost:3000**
6. **Access platform at http://localhost:3000**

---

## 🚀 Next Steps

1. Create @publishing-platform/core repo with shared types
2. Create API Gateway repo with routing
3. Create first module (content-management)
4. Setup event bus
5. Create remaining modules
6. Deploy to K3s

Ready? 🎯
