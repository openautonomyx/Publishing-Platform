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
(cd frontend && npm install --no-audit --no-fund && npm run build)

if [[ "${RUN_DOCKER_SMOKE:-0}" == "1" ]] && { ! command -v docker >/dev/null || ! docker compose version >/dev/null 2>&1; }; then
  echo 'ERROR: Docker Compose is required when RUN_DOCKER_SMOKE=1' >&2
  exit 1
fi

if command -v docker >/dev/null && docker compose version >/dev/null 2>&1; then
  step 'Docker Compose schema validation'
  PUBLISHING_API_TOKEN='ci-validation-placeholder-token-32-chars' docker compose -f compose.workspace.yml config --quiet
  if [[ "${RUN_DOCKER_SMOKE:-0}" == "1" ]]; then
    step 'Docker Compose build and local smoke'
    export PUBLISHING_API_TOKEN="${PUBLISHING_API_TOKEN:-$(node -e "console.log(require('node:crypto').randomBytes(32).toString('hex'))")}"
    cleanup() { docker compose -f compose.workspace.yml down || true; }
    trap cleanup EXIT
    docker compose -f compose.workspace.yml up --build -d
    for attempt in $(seq 1 30); do
      if curl --silent --fail http://127.0.0.1:8080/api/health | grep -q '"status":"ok"'; then
        echo 'API health passed'
        break
      fi
      if [[ "$attempt" == 30 ]]; then
        echo 'API did not become healthy' >&2
        exit 1
      fi
      sleep 2
    done
    curl --silent --fail http://127.0.0.1:8080/ >/dev/null
    auth_status=$(curl --silent --output /dev/null --write-out '%{http_code}' http://127.0.0.1:8080/api/publications)
    [[ "$auth_status" == '401' ]] || { echo "Expected 401 without token, got $auth_status" >&2; exit 1; }
    curl --silent --show-error --fail --header "Authorization: Bearer $PUBLISHING_API_TOKEN" http://127.0.0.1:8080/api/publications | grep -q '"items"'
    echo 'Web and authenticated API smoke passed'
  fi
else
  echo 'SKIP: docker compose not available (container smoke not validated)'
fi
step 'Validation commands finished successfully'
