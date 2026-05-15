# Multi-Tenant Routing Contract

The platform should follow a stable SaaS routing contract that works in Astro now and can be mirrored in Next.js App Router later.

## Principle

Do not couple product URLs to a framework.

```text
Product route contract = stable
Framework implementation = replaceable
```

## Standard Tenant Routes

Public publishing routes:

```text
/[tenant]
/[tenant]/articles/[slug]
/[tenant]/feed.xml
/[tenant]/feed.json
/[tenant]/sitemap.xml
```

Publishing console routes:

```text
/[tenant]/console
/[tenant]/console/articles
/[tenant]/console/articles/new
/[tenant]/console/articles/[articleId]/edit
/[tenant]/console/media
/[tenant]/console/canvas
/[tenant]/console/feeds
/[tenant]/console/schema
/[tenant]/console/locations
/[tenant]/console/automation
/[tenant]/console/settings
```

Platform API routes:

```text
/features.json
/openapi.json
/api/webhooks
/api/triggers
/api/mcp
```

## Next.js App Router Equivalent

If the console later moves to Next.js, the equivalent shape is:

```text
app/
├── [tenant]/
│   ├── page.tsx
│   ├── articles/[slug]/page.tsx
│   ├── feed.xml/route.ts
│   ├── feed.json/route.ts
│   ├── sitemap.xml/route.ts
│   └── console/
│       ├── page.tsx
│       ├── articles/page.tsx
│       ├── articles/new/page.tsx
│       ├── articles/[articleId]/edit/page.tsx
│       ├── media/page.tsx
│       ├── canvas/page.tsx
│       ├── feeds/page.tsx
│       ├── schema/page.tsx
│       ├── locations/page.tsx
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
│   ├── articles/[slug].astro
│   ├── feed.xml.ts
│   ├── feed.json.ts
│   ├── sitemap.xml.ts
│   └── console/
│       ├── index.astro
│       ├── articles/index.astro
│       ├── articles/new.astro
│       ├── articles/[articleId]/edit.astro
│       ├── media.astro
│       ├── canvas.astro
│       ├── feeds.astro
│       ├── schema.astro
│       ├── locations.astro
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
/acme/console
/acme/articles/my-post
```

Custom domains can be added later without changing internal route names.

## Rule

Use the same route contract regardless of framework.

This lets us start with Astro and later split or migrate the console to Next.js without changing product URLs.
