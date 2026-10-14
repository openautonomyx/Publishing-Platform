FROM node:20-alpine
WORKDIR /app
ENV NODE_ENV=production PORT=3000 PUBLISHING_DATA_FILE=/data/publications.json
COPY server.cjs ./
USER node
EXPOSE 3000
CMD ["node","server.cjs"]
