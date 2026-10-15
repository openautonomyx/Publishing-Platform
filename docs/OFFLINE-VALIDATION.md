# Offline validation

GitHub Actions is currently returning workflow-level `startup_failure` before jobs execute (see issue #6). This script provides reproducible local checks independently of GitHub Actions.

From a trusted clone on a machine with **Node.js >=20** and npm:

```bash
bash scripts/validate-workspace.sh
```

It checks JavaScript syntax, runs the publishing API integration tests, builds the React frontend, and validates Docker Compose configuration when Docker is available. A successful run **does not** imply production security readiness.

To additionally build and smoke-test the actual Caddy + Node containers, with Docker Engine and Compose available:

```bash
RUN_DOCKER_SMOKE=1 bash scripts/validate-workspace.sh
```

The smoke test requires Docker Compose (otherwise the script exits nonzero), checks the API health endpoint, verifies that unauthenticated publishing access is denied with HTTP 401, and confirms a token-authorized list response. The smoke test uses localhost port 8080 and temporarily starts the evaluation stack. It stops the stack on exit but leaves the named data volume untouched. Don't run this against a production deployment. It does not make external social-platform API calls.

## Evidence required before release

Record exact commit SHA, Node/npm/Docker versions, command output and exit status, and link passing GitHub Checks once Actions starts executing jobs. Report skipped Docker checks as **skipped**, not passed. Never upload API tokens or unpublished content.
