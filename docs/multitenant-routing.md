# Multi-Tenant CMS Routing Contract

This is a multi-tenant CMS and professional publishing platform, not a generic SaaS dashboard.

The route contract should be CMS-first and publishing-first.

## Principle

Do not couple product URLs to a framework.

```text
CMS route contract = stable
Framework implementation = replaceable
```

The CMS should support many tenants, but the product model is content management, editorial workflow, structured publishing, feeds, schema, and distribution.

## Public Publishing Routes

These are public reader-facing routes.

```text
/[tenant]
/[tenant]/[contentType]/[slug]
/[tenant]/authors/[authorSlug]
/[tenant]/topics/[topicSlug]
/[tenant]/tags/[tagSlug]
/[tenant]/locations/[locationSlug]
/[tenant]/collections/[collectionSlug]
/[tenant]/feeds/[feedSlug].xml
/[tenant]/feed.xml
/[tenant]/feed.json
/[tenant]/sitemap.xml
```

Examples:

```text
/acme/articles/future-of-ai
/acme/reports/state-of-education-2026
/acme/opinions/why-cities-matter
/acme/datasets/global-work-index
```

## CMS Admin Routes

These are editorial and administrative routes.

```text
/[tenant]/cms
/[tenant]/cms/content
/[tenant]/cms/content/new
/[tenant]/cms/content/[entryId]
/[tenant]/cms/content/[entryId]/edit
/[tenant]/cms/content-types
/[tenant]/cms/content-types/[contentType]
/[tenant]/cms/workflow
/[tenant]/cms/revisions
/[tenant]/cms/media
/[tenant]/cms/canvas
/[tenant]/cms/taxonomy
/[tenant]/cms/authors
/[tenant]/cms/locations
/[tenant]/cms/schema
/[tenant]/cms/seo
/[tenant]/cms/feeds
/[tenant]/cms/automation
/[tenant]/cms/settings
```

## CMS Concepts

The CMS should be centered around these primitives:

- content entries
- content types
- editorial workflow
- revisions
- media library
- canvas assets
- taxonomy
- authors and contributors
- locations
- Schema.org entities
- SEO metadata
- feeds
- automation triggers
- publishing settings

## Content Types

The CMS should support configurable content types backed by Schema.org-compatible semantic mappings where appropriate.

Examples:

```text
article       -> Article / NewsArticle / OpinionNewsArticle
report        -> Report / CreativeWork
dataset       -> Dataset
web_page      -> WebPage
profile       -> ProfilePage / Person
media_object  -> MediaObject / ImageObject / VideoObject / AudioObject
```

Schema.org is used for public semantic meaning, not for every internal CMS record.

## Platform API Routes

```text
/features.json
/openapi.json
/api/webhooks
/api/triggers
/api/mcp
```

## Next.js App Router Equivalent

If the CMS app later moves to Next.js, mirror this product route contract:

```text
app/
├── [tenant]/
│   ├── page.tsx
│   ├── [contentType]/[slug]/page.tsx
│   ├── authors/[authorSlug]/page.tsx
│   ├── topics/[topicSlug]/page.tsx
│   ├── tags/[tagSlug]/page.tsx
│   ├── locations/[locationSlug]/page.tsx
│   ├── collections/[collectionSlug]/page.tsx
│   ├── feed.xml/route.ts
│   ├── feed.json/route.ts
│   ├── sitemap.xml/route.ts
│   └── cms/
│       ├── page.tsx
│       ├── content/page.tsx
│       ├── content/new/page.tsx
│       ├── content/[entryId]/page.tsx
│       ├── content/[entryId]/edit/page.tsx
│       ├── content-types/page.tsx
│       ├── workflow/page.tsx
│       ├── revisions/page.tsx
│       ├── media/page.tsx
│       ├── canvas/page.tsx
│       ├── taxonomy/page.tsx
│       ├── authors/page.tsx
│       ├── locations/page.tsx
│       ├── schema/page.tsx
│       ├── seo/page.tsx
│       ├── feeds/page.tsx
│       ├── automation/page.tsx
│       └── settings/page.tsx
├── features.json/route.ts
├── openapi.json/route.ts
└── api/
    ├── webhooks/route.ts
    ├── triggers/route.ts
    └── mcp/route.ts
```

## Astro Implementation

Astro should mirror the same route contract:

```text
src/pages/
├── [tenant]/
│   ├── index.astro
│   ├── [contentType]/[slug].astro
│   ├── authors/[authorSlug].astro
│   ├── topics/[topicSlug].astro
│   ├── tags/[tagSlug].astro
│   ├── locations/[locationSlug].astro
│   ├── collections/[collectionSlug].astro
│   ├── feed.xml.ts
│   ├── feed.json.ts
│   ├── sitemap.xml.ts
│   └── cms/
│       ├── index.astro
│       ├── content/index.astro
│       ├── content/new.astro
│       ├── content/[entryId]/index.astro
│       ├── content/[entryId]/edit.astro
│       ├── content-types/index.astro
│       ├── workflow.astro
│       ├── revisions.astro
│       ├── media.astro
│       ├── canvas.astro
│       ├── taxonomy.astro
│       ├── authors.astro
│       ├── locations.astro
│       ├── schema.astro
│       ├── seo.astro
│       ├── feeds.astro
│       ├── automation.astro
│       └── settings.astro
├── features.json.ts
├── openapi.json.ts
└── api/
    ├── webhooks.ts
    ├── triggers.ts
    └── mcp.ts
```

## Tenant Resolution

Resolution order:

```text
custom domain -> tenant slug -> 404
```

Initial MVP uses path-based tenant slugs:

```text
/acme/cms
/acme/articles/my-post
```

Custom domains can be added later without changing internal route names.

## Rule

Use CMS-first URLs, not generic SaaS-console URLs.

This lets the platform stay true to its product model while remaining framework-portable.
