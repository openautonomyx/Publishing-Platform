# Integrated editorial workspace (first implementation)

This directory contains a **single-operator internal publishing API** based on Node.js built-ins. The existing React `/content` route is connected to this API by `frontend/src/pages/EditorialWorkspace.tsx`.

## Run locally

Requires Node.js 20+.

```bash
export PUBLISHING_API_TOKEN="$(node -e "console.log(require('crypto').randomBytes(32).toString('hex'))")"
node workspace-api/server.cjs
# second terminal:
cd frontend
npm install
npm run dev
```

Open http://localhost:5173/content and enter the generated token. Vite proxies `/api` to port 3000.

Data is saved atomically to `./data/publications.json` (override with `PUBLISHING_DATA_FILE`). Back it up; do not commit it. Do not expose API port 3000 publicly; terminate TLS at a reverse proxy and keep the token secret. This is not multi-tenant authentication and is **not production-ready for multiple users**.

## API

- `GET /api/health` (public liveness)
- `GET /api/publications` (list)
- `POST /api/publications` (create draft with title and content)
- `PUT /api/publications/:id` (edit draft)
- `POST /api/publications/:id/submit|approve|reject|schedule|publish|archive`
- Schedule body: `{"scheduledFor":"2026-12-01T10:00:00Z"}`

All content routes require `Authorization: Bearer <PUBLISHING_API_TOKEN>`. Scheduling and immediate publication update the **internal catalog only**; no external provider post is sent. A 15-second poll publishes due scheduled items while the process is running, and requests reconcile overdue items after restart.

## Known limitations and next implementation gates

- One token grants all editorial permissions: implement authenticated users, organization isolation and audited RBAC before deploying for a team.
- JSON file is only suitable for a single process. Use a transactional database plus migrations, optimistic concurrency, and transactional outbox for production.
- No outbound social-media adapters, OAuth callbacks, token refresh, retries, or delivery receipts.
- No webhook verification, scheduling queue, notifications, media storage, rate limiting, or analytics ingestion.
- Existing root Docker Compose references absent services. This workspace is standalone until the monorepo deployment is reconciled.
- Run `node --test workspace-api/server.test.cjs` and `cd frontend && npm run build` before merging.
