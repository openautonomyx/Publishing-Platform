# ADR 0003: Use Astro as the default web framework

## Status

Accepted

## Context

The platform needs a lightweight web framework for public publishing and the CMS interface.

## Decision

Use Astro as the default web framework.

## Rationale

Astro provides a lightweight architecture focused on content and publishing while still allowing React islands for rich interactive tools such as Tiptap and Fabric.js.

## Consequences

The platform uses:

- Astro for application and publishing routes
- React islands for interactive components
- Framework-independent route contracts

## Alternatives Considered

### Next.js

Strong option, but adds more framework weight than needed for the current stage.

### Full React SPA

Rejected because the platform is content and publishing oriented rather than client-only application oriented.
