FROM node:22-bookworm-slim

WORKDIR /app

COPY apps/web/package*.json ./
RUN npm install -g pnpm@9.12.0 && pnpm install --frozen-lockfile

COPY . .
RUN pnpm build

EXPOSE 3000
CMD ["node", "apps/web/dist/server/entry.mjs"]