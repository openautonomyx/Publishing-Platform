# Platform ground rules

These choices are non-negotiable at the platform level.

Implementation choices are reviewable.

## Publishing and content management

- Professional publishing
- Structured content management
- Editorial workflow support
- Multimodal publishing
- Multilingual publishing
- Feed support through open feed formats
- Canvas support through browser/web fundamentals

## Taxonomy and semantics

- Schema.org support where public semantic metadata fits
- Taxonomy support for topics, tags, authors, collections, locations, and content types
- Location as a first-class concept
- Time as a first-class concept

## Infrastructure

- Cloud-native architecture
- CNCF or CNCF-aligned infrastructure where it helps
- Lightweight runtime where possible
- Open standards for externally visible protocols

## Data storage and portability

- Portable identifiers where they exist
- Open Location Code / Plus Codes for location identity
- WGS84 latitude/longitude for coordinates
- Operational data stored in the chosen database layer
- Public meaning exposed through appropriate standards

## AI and automation

- AI-native workflows
- Automation support through OpenAPI, MCP, triggers, and webhooks
- Tenant-owned automation endpoints

## Discovery and distribution

- Open search and discovery standards
- Metadata
- Canonical URLs
- Sitemaps
- robots.txt
- hreflang
- Feed discovery
- Structured data where appropriate

## Reviewable implementation choices

Current implementation choices include:

- Astro
- SurrealDB
- Auth.js / NextAuth
- Fabric.js
- k3s
- systemd-managed Node runtime
- package structure
- route implementation
- deployment scripts
- adapter code
- UI components
- helper functions
- naming details
- operational tooling

## Notes

Feed support means using open formats such as RSS, Atom, and JSON Feed.

Canvas support means using browser/web platform capabilities. Fabric.js is a current implementation choice, not a ground rule.

Ground rules stay fixed.

Implementation can improve.
