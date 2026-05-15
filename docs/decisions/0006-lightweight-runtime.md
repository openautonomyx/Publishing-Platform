# ADR 0006: Lightweight runtime

## Status

Accepted

## Decision

Use a lightweight runtime architecture.

## Default

- k3s for lightweight Kubernetes ingress
- systemd-managed Node process for the app
- no Docker requirement for MVP

## Rule

- use CNCF where it helps infrastructure
- use open standards for protocols
- avoid heavy infrastructure until needed
