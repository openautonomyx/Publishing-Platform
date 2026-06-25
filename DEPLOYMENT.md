# OpenAutonomyX Deployment Guide

## 🚀 Quick Deployment Options

### Option 1: GitHub Pages (Marketing Site - Free)
**Status:** ✅ Ready to deploy
**URL:** https://openautonomyx.github.io

```bash
# Push to GitHub (triggers automatic deployment)
git push -u origin main

# GitHub Pages will auto-deploy docs/ folder
# Site live in ~1 minute at: https://openautonomyx.github.io
```

### Option 2: NPM Package Registry (Free)
**Status:** ✅ Ready to publish
**URL:** https://www.npmjs.com/package/openautonomyx

```bash
# Build the package
npm run build

# Login to npm (if not already)
npm login

# Publish to npm registry
npm run publish:npm

# Install globally
npm install -g openautonomyx
```

### Option 3: VPS Deployment (Self-Hosted API)
**Status:** ✅ Ready (requires VPS)
**Recommended:** DigitalOcean, Linode, AWS EC2, or your own server

```bash
# On your VPS:
git clone https://github.com/openautonomyx/openautonomyx.git
cd openautonomyx

# Install dependencies
npm install

# Set environment variables
export NODE_ENV=production
export LLM_PROVIDER=ollama
export DATABASE_URL=postgresql://user:pass@localhost/openautonomyx
export PORT=3000

# Start server
npm start

# Or use PM2 for process management
pm2 start server.js --name openautonomyx
pm2 save
```

### Option 4: Docker (Container Deployment)
**Status:** ✅ Ready (requires Docker)

```dockerfile
# Dockerfile (create this)
FROM node:18-alpine

WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .

EXPOSE 3000
CMD ["npm", "start"]
```

```bash
# Build and push
docker build -t openautonomyx:latest .
docker push yourusername/openautonomyx:latest

# Run
docker run -e NODE_ENV=production \
  -e LLM_PROVIDER=ollama \
  -p 3000:3000 \
  openautonomyx:latest
```

### Option 5: Vercel (Serverless - Free)
**Status:** ✅ Ready
**URL:** https://openautonomyx.vercel.app

```bash
# Install Vercel CLI
npm i -g vercel

# Deploy
vercel --prod

# Follow prompts
```

### Option 6: Railway (PaaS - Simple)
**Status:** ✅ Ready
**URL:** Auto-assigned, custom domain available

```bash
# Connect GitHub repo at https://railway.app
# Or use CLI:
npm i -g @railway/cli
railway init
railway up
```

---

## 📋 Deployment Checklist

### Pre-Deployment
- [ ] All code committed and pushed to GitHub
- [ ] License (MIT) in place
- [ ] README with attribution complete
- [ ] CONTRIBUTORS.md created
- [ ] Environment variables documented
- [ ] No secrets in code
- [ ] Version bumped (if needed)

### During Deployment
- [ ] GitHub Pages enabled (Settings → Pages → main/docs)
- [ ] NPM package published (if distributing SDK)
- [ ] API server started successfully
- [ ] Health check endpoint responding
- [ ] Database migrations run (if applicable)

### Post-Deployment
- [ ] Site accessible at public URL
- [ ] API responding on correct port
- [ ] Admin panel accessible
- [ ] Logging working
- [ ] Monitoring set up
- [ ] Backup strategy in place

---

## 🔧 Environment Configuration

### Required Variables
```bash
NODE_ENV=production
PORT=3000
```

### LLM Configuration
```bash
# Local (Ollama)
LLM_PROVIDER=ollama
LLM_MODEL=mistral
OLLAMA_API_URL=http://localhost:11434

# Cloud (OpenAI)
LLM_PROVIDER=openai
LLM_MODEL=gpt-4
OPENAI_API_KEY=sk-...

# Cloud (Anthropic)
LLM_PROVIDER=anthropic
LLM_MODEL=claude-3-sonnet
ANTHROPIC_API_KEY=sk-ant-...
```

### Database Configuration
```bash
# PostgreSQL (recommended)
DATABASE_TYPE=postgresql
DATABASE_URL=postgresql://user:password@host:5432/openautonomyx

# SQLite (development)
DATABASE_TYPE=sqlite
DATABASE_PATH=/data/openautonomyx.db
```

### Storage Configuration
```bash
# Local filesystem
STORAGE_TYPE=local
STORAGE_PATH=/data/uploads

# S3-compatible
STORAGE_TYPE=s3
AWS_ACCESS_KEY_ID=...
AWS_SECRET_ACCESS_KEY=...
AWS_REGION=us-east-1
```

---

## 🌐 Domain Setup

### Custom Domain (GitHub Pages)
1. Create `CNAME` file in `docs/` folder:
```
openautonomyx.com
```

2. Point DNS to GitHub Pages:
```
A record: 185.199.108.153
A record: 185.199.109.153
A record: 185.199.110.153
A record: 185.199.111.153
```

3. Or CNAME:
```
CNAME: openautonomyx.github.io
```

### Custom Domain (VPS)
1. Point A record to VPS IP
2. Configure reverse proxy (Nginx/Apache)
3. Set up SSL (Let's Encrypt)

---

## 📊 Monitoring & Logging

### Health Check Endpoint
```bash
curl https://api.openautonomyx.com/api/health
# Response: { "status": "healthy", "timestamp": "..." }
```

### Log Levels
```bash
DEBUG=openautonomyx:* npm start    # Verbose
LOG_LEVEL=info npm start           # Standard
LOG_LEVEL=error npm start          # Errors only
```

### Recommended Monitoring
- **Uptime:** Pingdom, UptimeRobot
- **Analytics:** Google Analytics, Plausible
- **Error Tracking:** Sentry, Rollbar
- **Logs:** ELK Stack, Datadog, NewRelic

---

## 🔒 Security Checklist

- [ ] HTTPS enforced (SSL certificate)
- [ ] Rate limiting enabled
- [ ] CORS properly configured
- [ ] Authentication tokens secure
- [ ] Database encrypted
- [ ] Backups automated
- [ ] DDoS protection (Cloudflare)
- [ ] Security headers set
- [ ] Regular security audits
- [ ] Dependency updates automated

---

## 📈 Performance Tips

1. **Enable caching:**
```bash
CACHE_TTL=3600
REDIS_URL=redis://localhost:6379
```

2. **Compress responses:**
```bash
GZIP_ENABLED=true
```

3. **Use CDN for static files:**
```bash
CDN_URL=https://cdn.openautonomyx.com
```

4. **Database optimization:**
- Create indexes on frequently queried columns
- Use connection pooling
- Regular VACUUM (PostgreSQL)

---

## 🆘 Troubleshooting

### Site not live after push
- Check GitHub Actions tab for deployment errors
- Verify docs/ folder exists
- Check branch settings (main vs master)

### API not responding
- Check process is running: `pm2 list`
- Check logs: `pm2 logs openautonomyx`
- Verify port is open: `netstat -tuln | grep 3000`
- Check firewall rules

### Database connection errors
- Verify DATABASE_URL is correct
- Check database server is running
- Verify credentials have proper permissions
- Check network connectivity

### Performance issues
- Check CPU/memory usage
- Enable response caching
- Add more database connections
- Scale horizontally (multiple instances)

---

## 📞 Support

- **Documentation:** https://openautonomyx.com/docs
- **GitHub Issues:** https://github.com/openautonomyx/openautonomyx/issues
- **Email:** support@openautonomyx.com
- **Discord:** https://discord.gg/openautonomyx

---

**Last Updated:** June 25, 2026
**Status:** Production Ready
