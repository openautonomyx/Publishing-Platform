#!/usr/bin/env bash
# Repeatable validation without relying on GitHub-hosted Actions.
set -Eeuo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

step() { printf '\n==> %s\n' "$*"; }
step 'Checking installed tools'
command -v node >/dev/null || { echo 'Node.js >=20 required' >&2; exit 1; }
node -e 'if (+process.versions.node.split(".")[0] < 20) process.exit(1)'
command -v npm >/dev/null || { echo 'npm required' >&2; exit 1; }

step 'Node syntax validation'
node --check workspace-api/server.cjs
node --check workspace-api/server.test.cjs

step 'Publishing API integration tests'
node --test workspace-api/server.test.cjs

step 'Frontend TypeScript and Vite production build'
(cd frontend && npm ci --no-audit --no-fund && npm run build)

if command -v docker >/dev/null && docker compose version >/dev/null 2>&1; then
  step 'Docker Compose schema validation'
  PUBLISHING_API_TOKEN='ci-validation-placeholder-token-32-chars' docker compose -f compose.workspace.yml config --quiet
  if [[ "${RUN_DOCKER_SMOKE:-0}" == "1" ]]; then
    step 'Docker Compose build and local smoke'
    export PUBLISHING_API_TOKEN="${PUBLISHING_API_TOKEN:-$(node -e "console.log(require('node:crypto').randomBytes(32).toString('hex'))")}"
    # Isolate test state from an operator's existing Compose project.
    export COMPOSE_PROJECT_NAME="publishing-smoke-$$"
    docker info >/dev/null
    cleanup() { docker compose -f compose.workspace.yml down --volumes || true; }
    trap cleanup EXIT
    docker compose -f compose.workspace.yml up --build -d
    node scripts/smoke-workspace.cjs

  fi
else
  if [[ "${RUN_DOCKER_SMOKE:-0}" == "1" ]]; then
    echo 'FAIL: Docker Compose is required for requested smoke validation' >&2
    exit 1
  fi
  echo 'SKIP: docker compose not available (container smoke not validated)'
fi
step 'Validation commands finished successfully'
