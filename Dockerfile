# Build stage
FROM node:18-alpine AS builder

WORKDIR /build

# Copy package files
COPY package*.json ./
COPY tsconfig.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY src ./src
COPY public ./public
COPY docs ./docs

# Build TypeScript
RUN npm run build

# Generate SBOM
RUN npm run generate:sbom --if-present || true

# Production stage
FROM node:18-alpine

LABEL org.opencontainers.image.title="OpenAutonomyX"
LABEL org.opencontainers.image.description="Vendor-neutral creative publishing platform with local LLM support"
LABEL org.opencontainers.image.authors="OpenAutonomyX Contributors, Claude AI"
LABEL org.opencontainers.image.vendor="OpenAutonomyX"
LABEL org.opencontainers.image.version="1.0.0"
LABEL org.opencontainers.image.source="https://github.com/openautonomyx/original"
LABEL org.opencontainers.image.documentation="https://openautonomyx.github.io/original"
LABEL org.opencontainers.image.licenses="MIT"

WORKDIR /app

# Create non-root user for security
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

# Copy package files
COPY package*.json ./

# Install production dependencies only
RUN npm ci --only=production && \
    npm cache clean --force

# Copy built application from builder
COPY --from=builder /build/dist ./dist
COPY --from=builder /build/public ./public
COPY --from=builder /build/docs ./docs

# Copy documentation and license
COPY README.md LICENSE CONTRIBUTORS.md ./
COPY DEPLOYMENT.md VENDOR-NEUTRAL-ARCHITECTURE.md ./

# Change ownership to nodejs user
RUN chown -R nodejs:nodejs /app

# Switch to nodejs user
USER nodejs

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
    CMD node -e "require('http').get('http://localhost:' + (process.env.PORT || 3000) + '/api/health', (r) => {if (r.statusCode !== 200) throw new Error(r.statusCode)})"

# Expose port
EXPOSE 3000

# Set environment
ENV NODE_ENV=production
ENV PORT=3000

# Start application
CMD ["node", "dist/index.js"]
