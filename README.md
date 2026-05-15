# Publishing Platform

**An AI-native, cloud-native, multimodal, multilingual, location-aware, time-aware, Schema.org-native, SEO-native, composable, professional publishing platform for authoritative opinions.**

Publishing Platform is built for research organizations, think tanks, universities, analysts, experts, institutions, and knowledge-driven teams that need to publish credible, high-quality, evidence-backed content at scale.

The goal is to let any organization launch a polished publishing presence, manage contributors, publish authoritative articles and reports, customize branding, localize content, contextualize by geography and time, optimize for discovery, and operate from one unified backend.

## Vision

Build the professional publishing infrastructure for authoritative voices across languages, formats, markets, locations, time, search engines, and knowledge graphs.

This platform is designed for organizations and experts who publish:

- Opinion essays
- Research reports
- Policy briefs
- Analysis pieces
- Expert commentary
- Data-backed articles
- Multilingual editions
- Region-specific perspectives
- Time-sensitive analysis
- Multimodal explainers
- Institutional knowledge

It should feel as easy as a modern SaaS tool, but powerful enough for serious publishing teams.

## Product Principles

### AI Native

AI should assist the publishing workflow from draft to distribution.

Planned AI capabilities:

- Draft assistance
- Research summarization
- Editorial suggestions
- Headline generation
- SEO metadata generation
- Schema.org metadata generation
- Source organization
- Translation assistance
- Localization workflows
- Content repurposing
- Multimodal article enrichment
- Time-aware and location-aware recommendations

### Cloud Native

The platform should be built for modern cloud deployment from day one.

Core expectations:

- Multi-tenant architecture
- Custom domains
- Serverless-friendly deployment
- CDN-backed media delivery
- Scalable storage
- API-first backend
- Secure authentication and access control

### Multimodal

Publishing should not be limited to text.

Supported content types should include:

- Articles
- Reports
- PDFs
- Charts
- Images
- Tables
- Data files
- Embeds
- Audio and video later

### Multilingual

The platform should support publishing across languages from the beginning.

Core multilingual capabilities:

- Language-specific article versions
- Translation workflow
- Locale-aware slugs
- Locale-aware SEO metadata
- RTL language support later
- Language switcher on public pages
- Tenant-level default language
- Per-publication language settings

### Location as a Primitive

Publishing often depends on geography.

Location should be modeled directly instead of treated as plain text.

Core location capabilities:

- Country, region, city, and market metadata
- Location-aware article targeting
- Location-specific editions
- Geographic filters and discovery
- Localized SEO and metadata
- Maps and regional pages later

### Time as a Primitive

Authoritative publishing depends on temporal context.

Time should be modeled directly across the platform.

Core time capabilities:

- Published time
- Updated time
- Scheduled publishing
- Event time
- Report period
- Timezone-aware display
- Historical archive views
- Time-based search and filtering
- Time-sensitive analysis pages

### Schema.org Native

Structured data should be part of the publishing model, not an afterthought.

Core Schema.org capabilities:

- Article schema
- NewsArticle schema
- OpinionNewsArticle schema
- Report schema
- Person and Organization schema
- Breadcrumb schema
- FAQ schema later
- Dataset schema later
- JSON-LD generation
- Tenant-level publisher identity
- Author and contributor metadata

### SEO Native

Discovery should be built into the core workflow.

Core SEO capabilities:

- Meta titles and descriptions
- Canonical URLs
- Open Graph metadata
- Twitter/X card metadata
- XML sitemaps
- Robots.txt
- Clean slugs
- Locale-aware routing
- hreflang support
- Structured data
- Redirect management later
- Search preview in editor later

### Composable

The system should be modular and extensible.

Core modules:

- Publishing
- Editor
- Media
- Search
- Analytics
- Billing
- Auth
- Tenant management
- AI workflows
- Translation workflows
- Location workflows
- Time workflows
- SEO workflows
- Schema.org workflows

### Professional

This is not a casual blogging tool.

It is designed for credible publishing, institutional trust, and polished presentation.

## Core Product Loop

A user can sign up, create a workspace, write an article, publish it, and view it on a public URL.

That loop comes first.

## Stack

| Layer | Technology |
| --- | --- |
| Frontend | Next.js |
| Language | TypeScript |
| Styling | Tailwind CSS |
| UI | shadcn/ui |
| Editor | Tiptap |
| Database | SurrealDB |
| Auth | Clerk |
| Media Storage | Cloudflare R2 |
| Analytics | PostHog |
| Billing | Stripe |
| Hosting | Vercel |
| AI | OpenAI-compatible model layer |
| i18n | Locale-aware routing and content model |
| Location | Structured geo metadata |
| Time | Timezone-aware temporal metadata |
| Schema | Schema.org / JSON-LD |
| SEO | Metadata, sitemaps, canonical URLs, hreflang |

## Why SurrealDB

SurrealDB is the backend foundation because it supports structured records, document-style content, graph relationships, flexible schemas, permissions, search, and real-time capabilities in one system.

That makes it a strong fit for a multi-tenant publishing platform where tenants, users, memberships, articles, translations, authors, tags, locations, timelines, schema entities, SEO metadata, media, and reports are deeply connected.

## Architecture

```text
User -> Next.js app -> Server actions/API routes -> SurrealDB
                                      |
                                      -> Cloudflare R2 for media
```

External services:

- Clerk for authentication
- Stripe for billing
- PostHog for analytics
- Vercel for hosting
- OpenAI-compatible APIs for AI workflows

## Core Data Model

- `tenant` — an organization or publication
- `user` — an authenticated person
- `membership` — connection between a user and tenant
- `role` — permissions within a tenant
- `article` — canonical publishable content owned by a tenant
- `article_translation` — language-specific article version
- `location` — structured geography for content and tenants
- `timeline_event` — temporal context for content, reports, and events
- `seo_metadata` — search and social metadata
- `schema_entity` — structured data for JSON-LD output
- `media_asset` — uploaded images, PDFs, and files

Every tenant-owned record should include a tenant reference.

## MVP Scope

### Must Have

- User authentication
- Tenant/workspace creation
- Membership model
- Article CRUD
- Rich text editor
- Draft and published states
- Public article pages
- Tenant-aware routing
- Basic branding settings
- Language metadata
- Published and updated timestamps
- Basic SEO metadata
- Article JSON-LD

### Soon After

- Media uploads
- Custom domains
- Basic analytics
- Team invitations
- Translation workflow
- Location metadata
- Scheduled publishing
- Timezone-aware publishing
- XML sitemaps
- hreflang support
- Publisher and author schema

### Later

- Approval workflows
- Reports and PDFs
- Newsletter tools
- Advanced search
- Billing plans
- Enterprise SSO
- Audit logs
- Regional editions
- Historical archives
- Maps and timeline views
- Redirect management
- Search preview
- Advanced schema types

## Suggested Repository Structure

```text
Publishing-Platform/
├── apps/
│   └── web/
├── packages/
│   ├── ui/
│   ├── db/
│   ├── editor/
│   └── shared/
├── infra/
├── docs/
└── README.md
```

## First Milestone

Build the publishing loop:

1. Create a Next.js app
2. Connect SurrealDB
3. Add Clerk auth
4. Create tenant records
5. Create membership records
6. Create article records
7. Add language, location, time, SEO, and schema metadata
8. Add a simple article editor
9. Render public article pages with metadata and JSON-LD

Success path:

```text
/chinmay/dashboard/articles/new
        -> publish
/chinmay/articles/my-first-post
```

## Development Principles

- Build the core loop first
- Treat language, location, time, SEO, and Schema.org as primitives
- Keep the backend unified
- Avoid unnecessary infrastructure
- Add complexity only when usage demands it
- Keep tenant isolation explicit
- Make migration paths possible

## Status

Early build phase.

Current focus:

- Project setup
- Core data model
- Tenant-aware publishing loop

## License

TBD
