# Publishing Platform

A multi-tenant publishing platform for research organizations, think tanks, universities, analysts, and knowledge-driven teams.

The goal is to let any organization launch a polished publishing site, manage contributors, publish articles and reports, customize branding, and operate from one unified backend.

## Vision

Build a modern SaaS publishing platform similar in spirit to research and media sites, but designed from the ground up for multiple tenants.

Each tenant should be able to:

- Create and manage a branded publication
- Publish articles, reports, pages, and datasets
- Invite team members
- Assign editorial roles
- Upload media and PDFs
- Configure SEO metadata
- Connect a custom domain
- Track basic analytics
- Manage subscription and billing

## Core Product Loop

A user can sign up, create a workspace, write an article, publish it, and view it on a public URL.

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

## Why SurrealDB

SurrealDB is the backend foundation because it supports structured records, document-style content, graph relationships, flexible schemas, permissions, search, and real-time capabilities in one system.

That makes it a strong fit for a multi-tenant publishing platform where tenants, users, memberships, articles, authors, tags, media, and reports are deeply connected.

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

## Core Data Model

- `tenant` — an organization or publication
- `user` — an authenticated person
- `membership` — connection between a user and tenant
- `role` — permissions within a tenant
- `article` — publishable content owned by a tenant

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

### Soon After

- Media uploads
- SEO metadata
- Custom domains
- Basic analytics
- Team invitations

### Later

- Approval workflows
- Scheduled publishing
- Reports and PDFs
- Newsletter tools
- Advanced search
- Billing plans
- Enterprise SSO
- Audit logs

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
7. Add a simple article editor
8. Render public article pages

Success path:

```text
/chinmay/dashboard/articles/new
        -> publish
/chinmay/articles/my-first-post
```

## Development Principles

- Build the core loop first
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
