# Publishing Platform

**AI-native infrastructure for professional, authoritative, multilingual publishing.**

Publishing Platform is a cloud-native, automation-native, feed-native, canvas-native, multimodal, multilingual, location-aware, time-aware, Schema.org-native, SEO-native, and composable publishing system for organizations that need to produce credible content at scale.

It is designed for research organizations, think tanks, universities, analysts, experts, institutions, and knowledge-driven teams that publish articles, reports, explainers, briefs, datasets, opinions, and structured knowledge.

## Vision

Modern publishing should be more than a blog engine.

Publishing Platform gives organizations a unified system to:

- Launch a polished publication presence
- Manage tenants, teams, contributors, roles, and permissions
- Publish authoritative articles, reports, and editorial content
- Add multilingual, geographic, temporal, SEO, and structured data context
- Create reusable visual and canvas-based assets
- Syndicate content through RSS, Atom, and JSON Feed
- Connect workflows through OpenAPI, webhooks, triggers, and MCP tools
- Operate with a clean backend that is ready for AI-assisted publishing workflows

The platform treats publishing as a structured knowledge workflow, not just page rendering.

## Core Principles

- **Publishing-first** — every feature should support credible, high-quality publishing.
- **AI-native** — the system should be easy for agents and automations to read, write, and extend.
- **Tenant-aware** — organizations, roles, permissions, and isolation are first-class concerns.
- **Schema-native** — structured data and JSON-LD should be built into the content model.
- **Location-aware** — location is modeled as a durable primitive, not a loose text field.
- **Time-aware** — events, publication dates, updates, and timelines are explicit.
- **Feed-native** — RSS, Atom, and JSON Feed are core outputs.
- **Automation-native** — triggers, webhooks, OpenAPI, and MCP are part of the platform surface.
- **Canvas-native** — visual assets and editorial graphics can be composed and exported.
- **Composable** — teams should be able to replace, extend, or integrate major capabilities.

## Location Primitive

Location is a first-class primitive.

The platform uses **Plus Codes**, based on the open-source **Open Location Code** project, as canonical location identifiers. Plus Codes encode latitude and longitude into a text code, can be encoded and decoded offline, and are useful where street addresses are missing, inconsistent, or politically sensitive.

```text
Plus Code         = canonical location identity
GeoCoordinates    = actual latitude/longitude
Schema.org Place  = public JSON-LD representation
SurrealDB record  = internal graph/storage representation
```

This keeps public location identity independent from internal database IDs.

## Technology Stack

| Layer | Technology |
| --- | --- |
| Frontend | Astro |
| Interactive Islands | React |
| Language | TypeScript |
| Styling | Tailwind CSS |
| Editor | Tiptap in React island |
| Canvas | Fabric.js in React island |
| Database | SurrealDB |
| Auth | Auth.js / NextAuth |
| Media Storage | Cloudflare R2 |
| Analytics | PostHog |
| Billing | Stripe |
| Hosting | Vercel or Cloudflare Pages |
| Automation | OpenAPI, MCP, triggers, webhooks |
| Feeds | RSS, Atom, JSON Feed |
| Schema | Schema.org / JSON-LD |
| SEO | Metadata, sitemaps, canonical URLs, hreflang, feed discovery |

## Core Data Model

The initial platform model includes:

| Entity | Purpose |
| --- | --- |
| `tenant` | Organization or publication workspace |
| `user` | Authenticated person |
| `session` | Authenticated session state |
| `membership` | Connection between a user and tenant |
| `role` | Permissions within a tenant |
| `article` | Canonical publishable content owned by a tenant |
| `article_translation` | Language-specific article version |
| `feed` | RSS, Atom, or JSON Feed definition |
| `location` | Geo object using Plus Code as canonical identity |
| `timeline_event` | Temporal context for content, reports, and events |
| `seo_metadata` | Search and social metadata |
| `schema_entity` | Structured data for JSON-LD output |
| `canvas_document` | Editable Fabric.js canvas source document |
| `canvas_asset` | Exported image or reusable visual asset |
| `automation_event` | Platform event emitted for automation |
| `trigger` | Event definition that can start workflows |
| `webhook_endpoint` | Tenant-owned automation destination |
| `api_client` | Tenant API client credentials and scopes |
| `mcp_tool` | AI-agent accessible platform capability |
| `media_asset` | Uploaded images, PDFs, and files |

## Current Implementation Status

The repository is moving from architecture planning into executable implementation.

Current state:

- Workspace package scaffold exists at the repository root.
- `feat/astro-scaffold` has been created for the first runnable web application.
- The first Astro layout file has been committed on the scaffold branch.
- Existing auth, signup, webhook, and schema work will be wired into this repository after the Astro shell is runnable.
- The next implementation step is to finish committing the `apps/web` scaffold, then open a pull request into `main`.

Pending application scaffold files:

- `apps/web/package.json`
- `apps/web/astro.config.mjs`
- `apps/web/tsconfig.json`
- `apps/web/tailwind.config.mjs`
- `apps/web/src/pages/index.astro`
- `apps/web/src/styles/global.css`
- `apps/web/src/lib/db.ts`
- `apps/web/.env.example`

After the scaffold lands, the repository should support:

```bash
pnpm install
pnpm dev
pnpm build
```

## Integration Plan

After the Astro shell is in place, the next work is to wire in existing platform capabilities rather than redesign them from scratch:

1. Signup and login flows
2. Auth/session handling
3. Tenant, user, membership, and role schema
4. RBAC-protected admin/editor routes
5. Signed webhooks and delivery logs
6. Article CRUD and publishing workflow
7. Schema.org JSON-LD and SEO metadata
8. Feed generation

## First Milestone

The first milestone is the complete publishing loop:

1. Create the Astro application
2. Add React islands for interactive surfaces
3. Connect SurrealDB
4. Add Auth.js / NextAuth authentication
5. Create tenant records
6. Create membership records
7. Create article records
8. Add language, location, time, SEO, and schema metadata
9. Generate RSS, Atom, and JSON Feed outputs
10. Emit automation events and support webhook endpoints
11. Add a basic Tiptap article editor
12. Render public article pages with metadata and JSON-LD

## Development Roadmap

### Phase 1 — Core Publishing Loop

- Astro shell
- Public article pages
- Tenant model
- User authentication
- Article CRUD
- Metadata and JSON-LD
- Feed generation

### Phase 2 — Editorial System

- Tiptap editor
- Draft and publish workflow
- Contributor roles
- Editorial review states
- Version history
- Media asset management

### Phase 3 — Structured Context

- Multilingual articles
- Plus Code location model
- Timeline events
- Schema.org entity graph
- Canonical URLs and hreflang

### Phase 4 — Automation Surface

- Webhook endpoints
- Automation events
- Trigger definitions
- OpenAPI surface
- MCP tools for agent workflows

### Phase 5 — Visual Publishing

- Fabric.js canvas editor
- Reusable visual templates
- Exported social cards
- Article graphics and report assets

## Development Principles

- Build the core loop first
- Keep the backend unified
- Keep tenant isolation explicit
- Use Astro as the lightweight publishing-first frontend
- Use React islands only for interactive tools
- Store app user, session, membership, and RBAC data in SurrealDB
- Treat language, location, time, SEO, Schema.org, feeds, canvas, and automation as primitives
- Use Plus Codes as canonical location identities
- Expose RSS, Atom, JSON Feed, OpenAPI, MCP, triggers, and webhooks
- Let tenants connect their own automation tools
- Use Fabric.js for visual canvas composition
- Do not embed n8n as a tenant-facing builder

## Repository Status

This project is in the early build phase.

The current focus is finishing the first runnable Astro web application scaffold, then wiring in existing signup, auth, schema, webhook, and publishing workflow capabilities.

## License

TBD
