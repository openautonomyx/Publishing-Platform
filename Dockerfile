FROM node:22-bookworm-slim

WORKDIR /app

# Install pnpm globally
RUN npm install -g pnpm@9.12.0

# Copy workspace files
COPY pnpm-workspace.yaml package.json pnpm-lock.yaml* ./
COPY apps ./apps
COPY packages ./packages

# Install dependencies
RUN pnpm install

# Build
RUN pnpm build

EXPOSE 3000
CMD ["pnpm", "--filter", "@publishing-platform/web", "start"]