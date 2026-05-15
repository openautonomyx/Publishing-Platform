# OCI Images

Containerisation is important for distribution.

Use OCI images without requiring Docker.

Current direction:

- build OCI images with Stacker
- run on OCI-compatible runtimes
- deploy with k3s when needed
- keep Docker out of the required workflow

## Stacker

Stacker builds OCI images from a declarative `stacker.yaml` file.

Repo: https://github.com/project-stacker/stacker

## Build

```bash
stacker build
```

## Image

```text
publishing-platform
```

## Notes

OCI images give us a standard distribution artifact.

Docker is not required for the build path.
