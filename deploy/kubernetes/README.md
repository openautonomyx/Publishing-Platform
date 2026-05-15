# Kubernetes Deployment

The Publishing Platform is designed to be deployment-platform neutral and Kubernetes-native.

## Components

- Publishing Platform web application
- SurrealDB
- Cloudflare R2 (external object storage)
- Auth.js-compatible authentication providers
- Optional PostHog and Stripe integrations

## Deployment Assets

- `deployment.yaml`

## Environment Variables

Configure application settings using Kubernetes ConfigMaps and Secrets.

## Recommended Additions

- Helm chart
- Kustomize overlays
- Horizontal Pod Autoscaler
- Cert-Manager
- External Secrets Operator
- GitOps via Argo CD or Flux

## Principle

```text
Application = cloud native
Deployment = Kubernetes native
Infrastructure = provider neutral
```
