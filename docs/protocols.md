# Protocols and Standards

The platform should respect open standards and use the right protocol for the right job.

Schema.org is important, but it is not the default representation for everything.

## Principle

```text
Use the standard that matches the use case.
```

Internal platform data may live in SurrealDB, but external interfaces should expose stable, portable, standards-based representations.

## Schema.org and the Database

Schema.org can influence the database structure, especially for public publishing objects.

Use Schema.org-shaped fields where they naturally fit:

- `Person` for authors, contributors, editors, experts, and profiles
- `Organization` for tenants, publishers, institutions, and research organizations
- `Place` for public location representation
- `CreativeWork` and subtypes for publishable content
- `Dataset` for datasets
- `MediaObject` for media assets

But do not force every internal record to be Schema.org.

For example, these are platform/application records, not Schema.org-first entities:

- sessions
- memberships
- RBAC permissions
- billing state
- webhook delivery logs
- API clients
- internal workflow events
- raw editor documents

The database model should therefore have two layers:

```text
Operational layer = SurrealDB-native records for app behavior
Semantic layer    = Schema.org-compatible fields/entities for public meaning
```

Example:

```text
article.content        = editor document / operational data
article.status         = publishing workflow state
article.schema         = Schema.org CreativeWork JSON-LD-compatible object
article.seo            = SEO metadata
article.tenant         = SurrealDB tenant relation
article.location       = linked geo object
```

## Protocol Mapping

| Capability | Primary Standard / Protocol | Purpose |
| --- | --- | --- |
| Public semantic metadata | Schema.org + JSON-LD | Search engines, knowledge graphs, structured public pages |
| API contracts | OpenAPI | Developer-facing HTTP API documentation and client generation |
| AI agent access | MCP | Tool and resource access for agents |
| Web automation | Webhooks | Event delivery to tenant-owned systems |
| Event definitions | Trigger catalog | Stable platform event names and payload shapes |
| Feeds | RSS, Atom, JSON Feed | Content syndication |
| Geo identity | Open Location Code / Plus Codes | Canonical portable location identifiers |
| Geo coordinates | WGS84 latitude/longitude | Actual coordinate representation |
| Media files | HTTP URLs + content types | Public and private asset access |
| SEO metadata | HTML meta tags + canonical URLs + sitemaps | Search discovery |
| Multilingual discovery | hreflang | Locale-aware search routing |
| Internal graph and records | SurrealDB records and relations | Storage, tenant isolation, permissions, linked data |

## Schema.org Scope

Use Schema.org for public semantic entities such as:

- Person
- Organization
- Place
- CreativeWork
- Article
- NewsArticle
- OpinionNewsArticle
- Dataset
- MediaObject
- BreadcrumbList
- WebSite

Do not use Schema.org as the only representation for:

- Internal authorization
- Sessions
- Billing state
- Webhook delivery logs
- API client secrets
- Internal workflow state
- Raw editor documents
- Operational analytics

Those should use internal application schemas and dedicated protocols.

## Feature Registry

The platform feature registry is currently exposed as plain HTTP JSON:

```text
/features.json
```

That endpoint is not Schema.org by default.

Future protocol-specific representations can include:

```text
/features.json              -> HTTP JSON
/openapi.json               -> OpenAPI
/mcp                        -> MCP tools/resources
/schema/platform.jsonld     -> optional semantic JSON-LD
```

## Location

Location uses Open Location Code / Plus Codes as canonical identifiers.

```text
Plus Code        = canonical location identity
GeoCoordinates   = WGS84 latitude/longitude
Schema.org Place = public semantic representation
SurrealDB record = internal linked graph representation
```

## Implementation Rule

Every outward-facing platform object should answer:

1. What is the internal SurrealDB record?
2. Which public standard applies?
3. Which protocol exposes it?
4. Is Schema.org appropriate here, or is another standard better?

Schema.org is a semantic layer, not a universal storage model.
