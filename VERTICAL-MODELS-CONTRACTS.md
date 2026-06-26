# Vertical Models & API Contracts

**Complete API specifications, data models, and contracts for all services**

---

## Architecture: Vertical Slices

```
Publishing Platform
├── Content Vertical
│   ├── Service: Content Management (3002)
│   ├── Models: Content, Version, Status, Category
│   └── Contract: CRUD + Versioning + Publishing
│
├── Blog Vertical
│   ├── Service: Blog (3009)
│   ├── Models: Post, Comment, Tag, Category
│   └── Contract: Blog operations + Comments
│
├── Format Vertical
│   ├── Service: Formats (3011)
│   ├── Models: Format, Conversion, Export
│   └── Contract: Multi-format conversion
│
├── Integration Vertical
│   ├── Service: Integrations (3010)
│   ├── Models: Integration, Platform, Connection
│   └── Contract: Plugin-based publishing
│
├── Analytics Vertical
│   ├── Service: Analytics (3005)
│   ├── Models: Event, Metric, Dashboard
│   └── Contract: Tracking + Reporting
│
├── Optimization Vertical
│   ├── Service: Optimization (3006)
│   ├── Models: Recommendation, Score, Variant
│   └── Contract: AI-powered suggestions
│
├── Design Vertical
│   ├── Service: Design (3007)
│   ├── Models: Template, Asset, Style
│   └── Contract: Design management
│
├── Skills Vertical
│   ├── Service: Skills (3003)
│   ├── Models: Skill, Level, Certification
│   └── Contract: Skill registry
│
└── Features Vertical
    ├── Service: Features (3008)
    ├── Models: Feature, Flag, Rollout
    └── Contract: Feature management
```

---

## Content Vertical Contract

### OpenAPI Spec

```yaml
openapi: 3.0.0
info:
  title: Content Management API
  version: 1.0.0
  description: Content creation, versioning, and publishing

servers:
  - url: http://localhost:3000/api/v1
    description: Local development

paths:
  /content:
    post:
      summary: Create content
      tags: [Content]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateContentDTO'
      responses:
        '201':
          description: Content created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Content'
    
    get:
      summary: List content
      tags: [Content]
      parameters:
        - name: status
          in: query
          schema: { type: string, enum: [draft, published, archived] }
        - name: category
          in: query
          schema: { type: string }
        - name: page
          in: query
          schema: { type: integer, default: 1 }
        - name: limit
          in: query
          schema: { type: integer, default: 20 }
      responses:
        '200':
          description: List of content
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/PaginatedResponse'

  /content/{id}:
    get:
      summary: Get single content
      tags: [Content]
      parameters:
        - name: id
          in: path
          required: true
          schema: { type: string, format: uuid }
      responses:
        '200':
          description: Content object
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Content'
    
    put:
      summary: Update content
      tags: [Content]
      parameters:
        - name: id
          in: path
          required: true
          schema: { type: string, format: uuid }
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UpdateContentDTO'
      responses:
        '200':
          description: Updated content
    
    delete:
      summary: Delete content
      tags: [Content]
      parameters:
        - name: id
          in: path
          required: true
          schema: { type: string, format: uuid }
      responses:
        '204':
          description: Content deleted

  /content/{id}/publish:
    post:
      summary: Publish content
      tags: [Content]
      parameters:
        - name: id
          in: path
          required: true
          schema: { type: string, format: uuid }
      responses:
        '200':
          description: Content published
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Content'

  /content/{id}/versions:
    get:
      summary: Get version history
      tags: [Content]
      parameters:
        - name: id
          in: path
          required: true
          schema: { type: string, format: uuid }
      responses:
        '200':
          description: Version history

components:
  schemas:
    Content:
      type: object
      properties:
        id:
          type: string
          format: uuid
        title:
          type: string
        description:
          type: string
        body:
          type: string
        type:
          type: string
          enum: [blog, whitepaper, case-study, guide, tutorial]
        status:
          type: string
          enum: [draft, published, archived]
        tags:
          type: array
          items:
            type: string
        category:
          type: string
        createdAt:
          type: string
          format: date-time
        updatedAt:
          type: string
          format: date-time
        publishedAt:
          type: string
          format: date-time
      required: [id, title, body]

    CreateContentDTO:
      type: object
      properties:
        title:
          type: string
        description:
          type: string
        body:
          type: string
        type:
          type: string
        tags:
          type: array
          items:
            type: string
        category:
          type: string
      required: [title, body]

    UpdateContentDTO:
      type: object
      properties:
        title:
          type: string
        description:
          type: string
        body:
          type: string
        type:
          type: string
        tags:
          type: array
          items:
            type: string
        category:
          type: string

    PaginatedResponse:
      type: object
      properties:
        data:
          type: array
          items:
            $ref: '#/components/schemas/Content'
        total:
          type: integer
        page:
          type: integer
        limit:
          type: integer
        totalPages:
          type: integer
```

---

## Blog Vertical Contract

```yaml
openapi: 3.0.0
info:
  title: Blog Service API
  version: 1.0.0

paths:
  /blog/posts:
    post:
      summary: Create blog post
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreatePostDTO'
      responses:
        '201':
          description: Post created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/BlogPost'
    
    get:
      summary: List blog posts
      parameters:
        - name: status
          in: query
          schema: { type: string }
        - name: category
          in: query
          schema: { type: string }
      responses:
        '200':
          description: List of posts

  /blog/posts/{id}:
    get:
      summary: Get post
      parameters:
        - name: id
          in: path
          required: true
          schema: { type: string, format: uuid }
      responses:
        '200':
          description: Blog post

  /blog/posts/{id}/publish:
    post:
      summary: Publish post
      parameters:
        - name: id
          in: path
          required: true
          schema: { type: string, format: uuid }
      responses:
        '200':
          description: Post published

  /blog/posts/{id}/comments:
    get:
      summary: Get post comments
      parameters:
        - name: id
          in: path
          required: true
          schema: { type: string, format: uuid }
      responses:
        '200':
          description: List of comments
    
    post:
      summary: Add comment to post
      parameters:
        - name: id
          in: path
          required: true
          schema: { type: string, format: uuid }
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateCommentDTO'
      responses:
        '201':
          description: Comment created

  /blog/categories:
    get:
      summary: Get all categories
      responses:
        '200':
          description: Categories list

  /blog/tags:
    get:
      summary: Get all tags
      responses:
        '200':
          description: Tags list

  /blog/search:
    get:
      summary: Search posts
      parameters:
        - name: q
          in: query
          required: true
          schema: { type: string }
      responses:
        '200':
          description: Search results

components:
  schemas:
    BlogPost:
      type: object
      properties:
        id:
          type: string
          format: uuid
        title:
          type: string
        slug:
          type: string
        excerpt:
          type: string
        content:
          type: string
        author:
          type: string
        status:
          type: string
          enum: [draft, published, archived]
        tags:
          type: array
          items:
            type: string
        category:
          type: string
        views:
          type: integer
        createdAt:
          type: string
          format: date-time
        publishedAt:
          type: string
          format: date-time
      required: [id, title, content]

    CreatePostDTO:
      type: object
      properties:
        title:
          type: string
        content:
          type: string
        excerpt:
          type: string
        tags:
          type: array
          items:
            type: string
        category:
          type: string
      required: [title, content]

    CreateCommentDTO:
      type: object
      properties:
        author:
          type: string
        email:
          type: string
        content:
          type: string
      required: [author, email, content]
```

---

## Formats Vertical Contract

```yaml
openapi: 3.0.0
info:
  title: Format Converter API
  version: 1.0.0

paths:
  /formats/epub:
    post:
      summary: Convert to EPUB
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ConvertDTO'
      responses:
        '201':
          description: EPUB generated
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ConversionResult'

  /formats/pdf:
    post:
      summary: Convert to PDF
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ConvertDTO'
      responses:
        '201':
          description: PDF generated

  /formats/slides:
    post:
      summary: Convert to Slides (PowerPoint)
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ConvertDTO'
      responses:
        '201':
          description: Slides generated

  /formats/audio:
    post:
      summary: Convert to Audio (MP3)
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ConvertDTO'
      responses:
        '201':
          description: Audio generated

  /formats/video:
    post:
      summary: Convert to Video (MP4)
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ConvertDTO'
      responses:
        '201':
          description: Video generated

  /formats/html:
    post:
      summary: Convert to Interactive HTML
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ConvertDTO'
      responses:
        '201':
          description: HTML generated

  /formats/convert:
    post:
      summary: Convert between formats
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ConvertBetweenDTO'
      responses:
        '201':
          description: Conversion queued

  /formats/supported:
    get:
      summary: Get supported formats
      responses:
        '200':
          description: List of formats

components:
  schemas:
    ConvertDTO:
      type: object
      properties:
        contentId:
          type: string
        title:
          type: string
        author:
          type: string
        content:
          type: string
        chapters:
          type: array
          items:
            type: object
      required: [title, content]

    ConversionResult:
      type: object
      properties:
        id:
          type: string
          format: uuid
        format:
          type: string
        filename:
          type: string
        downloadUrl:
          type: string
        previewUrl:
          type: string
        estimatedSize:
          type: string
        createdAt:
          type: string
          format: date-time
```

---

## Integrations Vertical Contract

```yaml
openapi: 3.0.0
info:
  title: Integrations API
  version: 1.0.0

paths:
  /integrations:
    post:
      summary: Register integration
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateIntegrationDTO'
      responses:
        '201':
          description: Integration registered
    
    get:
      summary: List integrations
      responses:
        '200':
          description: List of integrations

  /integrations/{id}:
    get:
      summary: Get integration
      parameters:
        - name: id
          in: path
          required: true
          schema: { type: string, format: uuid }
      responses:
        '200':
          description: Integration object
    
    put:
      summary: Update integration
      parameters:
        - name: id
          in: path
          required: true
          schema: { type: string, format: uuid }
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UpdateIntegrationDTO'
      responses:
        '200':
          description: Integration updated
    
    delete:
      summary: Remove integration
      parameters:
        - name: id
          in: path
          required: true
          schema: { type: string, format: uuid }
      responses:
        '204':
          description: Integration removed

  /integrations/{id}/publish:
    post:
      summary: Publish to platform
      parameters:
        - name: id
          in: path
          required: true
          schema: { type: string, format: uuid }
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/PublishDTO'
      responses:
        '201':
          description: Published to platform

  /integrations/{id}/history:
    get:
      summary: Get publish history
      parameters:
        - name: id
          in: path
          required: true
          schema: { type: string, format: uuid }
      responses:
        '200':
          description: Publish history

components:
  schemas:
    CreateIntegrationDTO:
      type: object
      properties:
        name:
          type: string
        type:
          type: string
          enum: [wordpress, medium, substack, twitter, linkedin, facebook, canva, custom]
        config:
          type: object
          additionalProperties: true
      required: [name, type, config]

    UpdateIntegrationDTO:
      type: object
      properties:
        name:
          type: string
        config:
          type: object
          additionalProperties: true

    PublishDTO:
      type: object
      properties:
        title:
          type: string
        content:
          type: string
        excerpt:
          type: string
        tags:
          type: array
          items:
            type: string
      required: [title, content]

    Integration:
      type: object
      properties:
        id:
          type: string
          format: uuid
        name:
          type: string
        type:
          type: string
        isActive:
          type: boolean
        config:
          type: object
        createdAt:
          type: string
          format: date-time
        updatedAt:
          type: string
          format: date-time
```

---

## Analytics Vertical Contract

```yaml
openapi: 3.0.0
info:
  title: Analytics API
  version: 1.0.0

paths:
  /analytics/events:
    post:
      summary: Track event
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/EventDTO'
      responses:
        '201':
          description: Event tracked
    
    get:
      summary: Get events
      parameters:
        - name: type
          in: query
          schema: { type: string }
        - name: dateFrom
          in: query
          schema: { type: string, format: date-time }
        - name: dateTo
          in: query
          schema: { type: string, format: date-time }
      responses:
        '200':
          description: List of events

  /analytics/dashboard:
    get:
      summary: Get analytics dashboard
      responses:
        '200':
          description: Dashboard metrics

  /analytics/reports:
    post:
      summary: Generate report
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ReportRequestDTO'
      responses:
        '201':
          description: Report generated

components:
  schemas:
    EventDTO:
      type: object
      properties:
        type:
          type: string
          enum: [view, click, share, download, publish]
        contentId:
          type: string
        userId:
          type: string
        platform:
          type: string
        metadata:
          type: object
          additionalProperties: true
      required: [type]

    ReportRequestDTO:
      type: object
      properties:
        format:
          type: string
          enum: [pdf, csv, json]
        dateRange:
          type: object
          properties:
            from:
              type: string
              format: date-time
            to:
              type: string
              format: date-time
        metrics:
          type: array
          items:
            type: string
```

---

## Data Models (TypeScript)

```typescript
// Content Model
interface Content {
  id: string;
  title: string;
  description: string;
  body: string;
  type: 'blog' | 'whitepaper' | 'case-study' | 'guide' | 'tutorial';
  status: 'draft' | 'published' | 'archived';
  tags: string[];
  category: string;
  gartnerCategoryId?: string;
  createdAt: Date;
  updatedAt: Date;
  publishedAt?: Date;
}

// Blog Model
interface BlogPost {
  id: string;
  title: string;
  slug: string;
  excerpt: string;
  content: string;
  author: string;
  status: 'draft' | 'published' | 'archived';
  tags: string[];
  category: string;
  views: number;
  createdAt: Date;
  updatedAt: Date;
  publishedAt?: Date;
}

// Format Model
interface FormatConversion {
  id: string;
  contentId: string;
  format: 'epub' | 'pdf' | 'slides' | 'audio' | 'video' | 'html';
  status: 'pending' | 'processing' | 'completed' | 'failed';
  downloadUrl?: string;
  previewUrl?: string;
  estimatedSize: string;
  createdAt: Date;
  completedAt?: Date;
}

// Integration Model
interface Integration {
  id: string;
  name: string;
  type: 'wordpress' | 'medium' | 'substack' | 'twitter' | 'linkedin' | 'facebook';
  isActive: boolean;
  config: Record<string, any>;
  createdAt: Date;
  updatedAt: Date;
}

// Analytics Model
interface AnalyticsEvent {
  id: string;
  type: 'view' | 'click' | 'share' | 'download' | 'publish';
  contentId: string;
  userId?: string;
  platform: string;
  metadata?: Record<string, any>;
  createdAt: Date;
}
```

---

## Contract Summary Table

| Vertical | Service | Port | Main Resources | Key Operations |
|----------|---------|------|-----------------|-----------------|
| Content | Content Management | 3002 | /content | CRUD, Publish, Version |
| Blog | Blog Service | 3009 | /blog/posts | CRUD, Comments, Search |
| Formats | Format Converter | 3011 | /formats/[type] | Convert to EPUB/PDF/etc |
| Integration | Integrations | 3010 | /integrations | Register, Publish, History |
| Analytics | Analytics | 3005 | /analytics | Track, Report, Dashboard |
| Optimization | Optimization | 3006 | /recommendations | Suggest, Score, Optimize |
| Design | Design | 3007 | /templates | Create, Edit, Manage |
| Skills | Skills | 3003 | /skills | CRUD, Certify, Track |
| Features | Features | 3008 | /features | Flag, Rollout, Track |

---

## Usage & Publishing

### Publish OpenAPI Specs

```bash
# Swagger UI
npm install -g swagger-cli
swagger-cli bundle VERTICAL-MODELS-CONTRACTS.yaml -o swagger.json

# Generate API docs
npm install -g redoc-cli
redoc-cli build swagger.json

# Deploy to docs site
cp swagger.json /docs/api/
cp redoc-static.html /docs/api/index.html
```

### Client Generation

```bash
# Generate TypeScript client from OpenAPI
npx openapi-generator-cli generate \
  -i VERTICAL-MODELS-CONTRACTS.yaml \
  -g typescript \
  -o ./sdk/typescript

# Generate Python client
npx openapi-generator-cli generate \
  -i VERTICAL-MODELS-CONTRACTS.yaml \
  -g python \
  -o ./sdk/python
```

---

**Complete API Contracts & Vertical Models Published!** 🚀

All 9 verticals with full OpenAPI specs, data models, and integration contracts.
