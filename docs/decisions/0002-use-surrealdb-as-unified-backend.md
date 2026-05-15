# ADR 0002: Use SurrealDB as the unified backend

## Status

Accepted

## Context

The platform needs one backend foundation for multi-tenant content, users, memberships, publishing workflow, schema entities, geo objects, revisions, events, and relationships.

## Decision

Use SurrealDB as the unified backend and system of record for operational CMS data.

## Rationale

SurrealDB supports structured records, document-style content, graph-style relationships, permissions, and flexible modeling in one database layer.

This fits the platform because CMS entities are highly connected:

- tenants
- users
- memberships
- content entries
- revisions
- publishing events
- schema entities
- locations
- feeds
- automation events

## Consequences

SurrealDB will store operational platform state.

External systems remain responsible for their own domains:

- object storage for files
- Stripe for billing
- PostHog for product analytics
- Auth.js-compatible providers for authentication flow

## Alternatives Considered

### PostgreSQL

Strong relational foundation, but graph/content flexibility would need more application-layer modeling.

### MongoDB

Good document storage, but tenant permissions and graph relationships would need more custom structure.

### Multiple specialized databases

Rejected for now because the platform should start with one unified backend and add complexity only when needed.
