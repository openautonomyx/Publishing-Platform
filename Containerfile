# OCI-compatible lightweight container build for the Publishing Platform web app.
# Build with Podman or Buildah. Avoids Docker-specific assumptions.

FROM node:22-alpine AS deps
WORKDIR /workspace
RUN corepack enable
COPY package.json pnpm-lock.yaml* pnpm-workspace.yaml* ./
COPY apps/web/package.json apps/web/package.json
RUN pnpm install --frozen-lockfile=false

FROM node:22-alpine AS build
WORKDIR /workspace
RUN corepack enable
COPY --from=deps /workspace/node_modules ./node_modules
COPY --from=deps /workspace/apps/web/node_modules ./apps/web/node_modules
COPY . .
RUN pnpm --filter @publishing-platform/web build

FROM node:22-alpine AS runtime
WORKDIR /app
ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=4321
RUN addgroup -S app && adduser -S app -G app
COPY --from=build /workspace/apps/web/dist ./dist
COPY --from=build /workspace/apps/web/package.json ./package.json
USER app
EXPOSE 4321
CMD ["node", "./dist/server/entry.mjs"]
