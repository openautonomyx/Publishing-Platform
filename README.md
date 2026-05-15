# Publishing Platform

**An AI-native, cloud-native, automation-native, feed-native, canvas-native, multimodal, multilingual, location-aware, time-aware, Schema.org-native, SEO-native, composable, professional publishing platform for authoritative opinions.**

Publishing Platform is built for research organizations, think tanks, universities, analysts, experts, institutions, and knowledge-driven teams that need to publish credible, high-quality, evidence-backed content at scale.

The goal is to let any organization launch a polished publishing presence, manage contributors, publish authoritative articles and reports, customize branding, localize content, contextualize by geography and time, create visual assets, syndicate through feeds, connect their own automations, optimize for discovery, and operate from one unified backend.

## Location Primitive

Location is a first-class primitive.

The platform uses **Plus Codes**, based on the open-source **Open Location Code** project, as canonical location identifiers. Plus Codes encode latitude and longitude into a text code, can be encoded and decoded offline, and are suitable where street addresses are missing or inconsistent.

Location model:

```text
Plus Code        = canonical location identity
GeoCoordinates   = actual latitude/longitude
Schema.org Place = public JSON-LD representation
SurrealDB record = internal graph/storage representation
```

This keeps location identity independent from internal database IDs.

## Stack

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

- `tenant` — an organization or publication
- `user` — an authenticated person
- `session` — authenticated session state
- `membership` — connection between a user and tenant
- `role` — permissions within a tenant
- `article` — canonical publishable content owned by a tenant
- `article_translation` — language-specific article version
- `feed` — syndication definition for RSS, Atom, or JSON Feed output
- `location` — linked geo object using Plus Code as canonical identity
- `timeline_event` — temporal context for content, reports, and events
- `seo_metadata` — search and social metadata
- `schema_entity` — structured data for JSON-LD output
- `canvas_document` — editable Fabric.js canvas source document
- `canvas_asset` — exported image or reusable visual asset
- `automation_event` — platform event emitted for automation
- `trigger` — platform event definition that can start workflows
- `webhook_endpoint` — tenant-owned automation destination
- `api_client` — tenant API client credentials and scopes
- `mcp_tool` — AI-agent accessible platform capability
- `media_asset` — uploaded images, PDFs, and files

## First Milestone

Build the publishing loop:

1. Create an Astro app
2. Add React islands
3. Connect SurrealDB
4. Add Auth.js / NextAuth authentication
5. Create tenant records
6. Create membership records
7. Create article records
8. Add language, location, time, SEO, and schema metadata
9. Add basic feed generation
10. Add basic automation events and webhook endpoints
11. Add a simple Tiptap article editor
12. Render public article pages with metadata and JSON-LD

## Development Principles

- Build the core loop first
- Use Astro as the lightweight publishing-first frontend
- Use Auth.js / NextAuth for authentication
- Store app user, session, membership, and RBAC data in SurrealDB
- Use React islands only for interactive tools
- Treat language, location, time, SEO, Schema.org, feeds, canvas, and automation as primitives
- Use Plus Codes as canonical location identities
- Expose RSS, Atom, JSON Feed, OpenAPI, MCP, triggers, and webhooks
- Let tenants connect their own automation tools
- Use Fabric.js for visual canvas composition
- Do not embed n8n as a tenant-facing builder
- Keep the backend unified
- Keep tenant isolation explicit

## Status

Early build phase.

## License

TBD
