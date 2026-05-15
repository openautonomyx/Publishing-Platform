# Design Considerations

These design considerations guide the platform.

They reflect our current understanding and are open to review.

Feedback is welcome.

## Publishing and content management

- Professional publishing
- Structured content management
- Editorial workflow support
- Multimodal publishing
- Multilingual publishing
- Feed support through open feed formats
- Canvas support through browser and web platform fundamentals

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
- Open Location Code (Plus Codes) for location identity
- WGS84 latitude and longitude for coordinates
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

Examples include:

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
