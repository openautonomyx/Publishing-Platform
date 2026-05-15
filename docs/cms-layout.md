# CMS Layout Architecture

There is no single universal protocol for application layout like OpenAPI or RSS.

For the CMS interface, the platform should follow open web standards and established CMS information architecture patterns.

## Standards and Conventions

Use:

- Semantic HTML for structure
- WAI-ARIA only where native HTML is not enough
- Responsive CSS layout
- Keyboard navigability
- Landmark regions
- Consistent navigation hierarchy
- CMS-first information architecture

## Layout Principle

```text
CMS layout = content operations first
```

The layout should prioritize editorial and publishing work, not generic SaaS dashboards.

## Primary Regions

```text
App Shell
├── Sidebar navigation
├── Top bar
├── Workspace / tenant context
├── Main content region
├── Secondary panel / inspector
└── Status / workflow surface
```

## Semantic HTML Landmarks

```html
<body>
  <aside aria-label="CMS navigation"></aside>
  <header></header>
  <main></main>
  <section aria-label="Inspector"></section>
</body>
```

## CMS Navigation

Primary CMS navigation:

```text
Content
Content Types
Workflow
Revisions
Media
Canvas
Taxonomy
Authors
Locations
Schema
SEO
Feeds
Automation
Settings
```

## Content Editing Layout

The content editor should use a three-region layout:

```text
┌──────────────────────────────────────────────┐
│ Top bar: title, status, save, preview, publish│
├───────────────┬────────────────────┬─────────┤
│ Content tree  │ Editor canvas/body │ Inspector│
│ / sections    │                    │ Metadata │
└───────────────┴────────────────────┴─────────┘
```

## Inspector Panels

The inspector should expose publishing primitives:

- Status
- Schedule
- Author
- Language
- Location
- SEO
- Schema.org mapping
- Feeds
- Automation triggers
- Permissions

## Public Layout

Public publishing pages should prioritize:

- Article readability
- Author identity
- Organization/publisher identity
- Publication date and updated date
- Language and location context
- Schema.org JSON-LD
- SEO metadata
- Feeds and discovery links

## Framework Independence

The layout contract should be framework-independent.

```text
Layout IA = stable
Astro / React / Next.js implementation = replaceable
```

## Rule

Do not design the CMS as a generic SaaS admin panel.

Design it as an editorial operating system for structured professional publishing.
