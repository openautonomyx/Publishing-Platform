# Shared Infrastructure Boundary

This repository contains application code for the Publishing Platform.

Infrastructure is shared across multiple apps and should live outside this repository.

## This repo should include

- Astro application code
- React islands
- SurrealDB schema and seed files used by this app
- Auth, RBAC, article, feed, and webhook application logic
- A lightweight OCI `Containerfile` for building the app image
- Minimal runtime notes showing required environment variables and exposed port

## This repo should not own

- k3s cluster installation
- Kubernetes cluster lifecycle
- shared ingress controllers
- shared certificate management
- shared container registry setup
- shared observability stack
- shared secret manager
- shared SurrealDB infrastructure lifecycle

Those belong in the shared infrastructure bundle.

## Runtime contract

The application image should expose:

```text
PORT=4321
HOST=0.0.0.0
```

The shared infra layer is responsible for:

```text
image pull
service routing
tls
domain mapping
secrets injection
persistent storage
observability
scaling policy
```

## Recommended flow

```text
Publishing-Platform repo
  -> build OCI image
  -> push image to registry
  -> shared infra deploys image
```

This keeps the application repository lightweight and lets one infrastructure bundle run multiple apps consistently.
