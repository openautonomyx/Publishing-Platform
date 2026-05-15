# ADR 0005: Use CNCF or open standards by layer

## Status

Accepted

## Context

The platform should be cloud-native, portable, and standards-respecting.

Not every part of a CMS belongs in the CNCF ecosystem. CNCF is most relevant for infrastructure, runtime, deployment, observability, networking, and cloud-native operations.

Publishing, CMS modeling, schema, feeds, layout, geo identity, and API contracts should follow the appropriate open standard instead of being forced into a CNCF-only filter.

## Decision

Use this rule:

```text
Infrastructure layer -> prefer CNCF official projects where practical
Protocol/product layer -> prefer open standards
Application layer -> use best-fit libraries with clear boundaries
```

## Layer Policy

| Layer | Preference |
| --- | --- |
| Orchestration | CNCF official / Kubernetes-compatible |
| Ingress / networking | Kubernetes-compatible, CNCF-aligned where practical |
| Runtime | OCI / runC ecosystem where relevant |
| Observability later | CNCF official projects preferred |
| API contracts | OpenAPI |
| AI agent access | MCP |
| Public semantic metadata | Schema.org + JSON-LD |
| Feeds | RSS, Atom, JSON Feed |
| Geo identity | Open Location Code / Plus Codes |
| Coordinates | WGS84 latitude/longitude |
| CMS internal records | SurrealDB operational schema |
| Editor / canvas / UI | best-fit application libraries |

## Rationale

This avoids fake standardization.

CNCF should guide cloud-native infrastructure choices. Open standards should guide protocol and interoperability choices. Product/application libraries should be chosen for fit, maintainability, and replaceability.

## Consequences

We should not claim a component is CNCF-backed unless it is verified from CNCF official sources.

We should not use Wikipedia as an authority for standards or CNCF status.

We should document whether a choice is:

- CNCF official
- CNCF-aligned
- open standard
- application library
- external service
- internal implementation detail

## Rule

```text
CNCF official or open standard where applicable.
Best-fit application library where no standard applies.
```
