# k3s Lightweight Deployment

Publishing Platform targets a lightweight OCI-first deployment model:

```text
Containerfile -> OCI image -> k3s Deployment -> Service/Ingress
```

This project does not require a full multi-node Kubernetes platform to start. Use k3s for a small, production-friendly runtime.

## Why k3s

- Lightweight Kubernetes distribution
- Low memory footprint
- Works well on a single VPS or small cluster
- Supports standard Kubernetes manifests
- Compatible with OCI images built by Podman, Buildah, or Docker

## Recommended local/edge shape

```text
k3s node
├── publishing-platform-web Deployment
├── publishing-platform-web Service
├── Traefik Ingress Controller (bundled with k3s by default)
└── external or sidecar SurrealDB, depending on environment
```

## Build image

```bash
podman build -t ghcr.io/openautonomyx/publishing-platform:latest -f Containerfile .
podman push ghcr.io/openautonomyx/publishing-platform:latest
```

## Deploy

```bash
kubectl apply -k deploy/kubernetes
```

The manifests under `deploy/kubernetes` are intentionally minimal and k3s-compatible.

## Production notes

- Keep SurrealDB storage durable.
- Use Kubernetes Secrets or an external secret manager for credentials.
- Use the bundled k3s Traefik ingress for simple deployments.
- Move to Helm only when the deployment surface grows.
