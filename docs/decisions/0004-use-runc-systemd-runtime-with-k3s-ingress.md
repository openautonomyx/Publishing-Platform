# ADR 0004: Use runC/systemd runtime with k3s ingress

## Status

Accepted

## Context

The platform should be lightweight to deploy and should not require Docker-based image deployment.

We still want a Kubernetes-compatible edge/runtime environment for ingress and future portability.

## Decision

Use a direct runtime model for the application process:

```text
pnpm build
node apps/web/dist/server/entry.mjs
systemd manages the process
k3s ingress routes traffic
```

k3s is used for lightweight Kubernetes networking and ingress, not as a Docker-image-first deployment requirement.

## Rationale

This keeps deployment lightweight while preserving a Kubernetes-native path.

The application runtime stays simple:

- no Dockerfile
- no image registry requirement
- no GHCR dependency
- no container build pipeline required for MVP

k3s still gives:

- ingress
- Kubernetes API compatibility
- portability path
- room to evolve later

## Consequences

Deployment assets should prefer:

- systemd service files
- environment files
- k3s ingress manifests
- direct build artifact deployment

Avoid for now:

- Dockerfile
- image registry workflow
- container-first Kubernetes Deployment manifests

## Rule

```text
Runtime = systemd-managed Node process
Ingress = k3s
Container image deployment = not required for MVP
```
