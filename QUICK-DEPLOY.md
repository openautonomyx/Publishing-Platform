# OpenAutonomyX - Quick Deployment Guide

**🚀 Fastest path to production in 4 commands**

---

## What's Being Deployed

```
📍 REACT FRONTEND
   https://app.publishing.openautonomyx.com
   Platform: Vercel (auto-deploys on git push)
   
📍 MARKETING SITE  
   https://publishing.openautonomyx.com
   Platform: GitHub Pages (live now!)
   
📍 ENTERPRISE PORTAL
   https://publishing.openautonomyx.com:8080
   Platform: Docker / Liferay DXP
```

---

## 4-Step Deployment

### Step 1: Setup Vercel (React Frontend)

```bash
# Install Vercel CLI
npm install -g vercel

# Login to Vercel account
vercel login

# Deploy production build
cd /Users/chinmaypanda/CustomApps/frontend
vercel --prod

# Configure when prompted:
# Project name: openautonomyx-platform
# Framework: Vite
# Root directory: ./
# Build command: npm run build
# Output directory: dist
```

**Result:** React frontend live at `https://openautonomyx-platform.vercel.app`

---

### Step 2: Configure Custom Domain (Vercel)

```
1. Open: https://vercel.com/dashboard
2. Click: openautonomyx-platform project
3. Go to: Settings → Domains
4. Add: app.publishing.openautonomyx.com
5. Follow DNS instructions
```

**Result:** Frontend accessible at `https://app.publishing.openautonomyx.com`

---

### Step 3: Enable GitHub Pages (Marketing Site)

```
1. Go to: https://github.com/Open-Autonomyx/Publishing-Platform/settings/pages
2. Source: Deploy from a branch
3. Branch: main
4. Folder: /docs
5. Custom domain: publishing.openautonomyx.com
6. Save
```

**Then update DNS:**

```
publishing.openautonomyx.com  CNAME  openautonomyx.github.io.
```

**Result:** Marketing site live at `https://publishing.openautonomyx.com` ✅

---

### Step 4: Start Liferay Portal (Local)

```bash
cd /Users/chinmaypanda/CustomApps

# Create environment file
cat > .env << 'EOF'
DB_USER=postgres
DB_PASS=your_secure_password
ENVIRONMENT=production
EOF

# Start Docker containers
docker-compose -f docker-compose.prod.yml up -d

# Wait for startup (30-60 seconds)
docker-compose -f docker-compose.prod.yml logs -f liferay

# Access portal
open http://localhost:8080
```

**Default credentials:**
- Username: `test@liferay.com`
- Password: `test`

---

## DNS Configuration Checklist

Add these records to your domain registrar:

```
CNAME Records:
  publishing.openautonomyx.com          →  openautonomyx.github.io
  app.publishing.openautonomyx.com      →  cname.vercel-dns.com
  api.publishing.openautonomyx.com      →  [your-vps-ip] or API domain
  blog.publishing.openautonomyx.com     →  [blog-domain]
  formats.publishing.openautonomyx.com  →  [formats-domain]
  integrations.publishing.openautonomyx.com → [integrations-domain]

A Records (if using VPS):
  @                                      →  [your-vps-ip]
  api                                    →  [your-vps-ip]
```

---

## Verify All Systems Are Live

```bash
# Marketing site
curl -I https://publishing.openautonomyx.com
# Expected: 200 OK (from GitHub Pages)

# React frontend
curl -I https://app.publishing.openautonomyx.com
# Expected: 200 OK (from Vercel)

# Liferay portal
curl -I http://localhost:8080/
# Expected: 200 OK (from Docker)
```

---

## Troubleshooting

### "Vercel deployment failed"
```bash
# Rebuild locally first
cd /Users/chinmaypanda/CustomApps/frontend
npm install
npm run build

# Check dist folder exists
ls -la dist/

# Try deployment again
vercel --prod --force
```

### "GitHub Pages not showing custom domain"
```bash
# DNS can take 24 hours to propagate
# Check DNS status:
nslookup publishing.openautonomyx.com

# Should show: openautonomyx.github.io
```

### "Liferay won't start"
```bash
# Check Docker is running
docker ps

# Check logs
docker-compose -f docker-compose.prod.yml logs liferay

# Restart services
docker-compose -f docker-compose.prod.yml restart
```

---

## What's Next?

### Immediate (Today)
- ✅ Deploy React frontend
- ✅ Enable GitHub Pages
- ✅ Start Liferay locally

### Short-term (This Week)
- [ ] Configure custom domains
- [ ] Setup SSL certificates (Let's Encrypt)
- [ ] Configure API Gateway
- [ ] Test end-to-end workflow

### Production (Next Week)
- [ ] Deploy to VPS with K3s
- [ ] Setup monitoring (Prometheus/Grafana)
- [ ] Configure load balancing
- [ ] Enable auto-scaling

---

## Automated Deployment Script

Use the interactive deployment script:

```bash
# Interactive mode (choose steps)
./DEPLOYMENT-SCRIPTS.sh

# Or run all steps at once
./DEPLOYMENT-SCRIPTS.sh all

# Or run individual steps
./DEPLOYMENT-SCRIPTS.sh react      # Deploy React to Vercel
./DEPLOYMENT-SCRIPTS.sh pages      # Enable GitHub Pages
./DEPLOYMENT-SCRIPTS.sh ssl        # Setup SSL
./DEPLOYMENT-SCRIPTS.sh liferay    # Start Liferay
```

---

## Production Deployment Summary

```
┌─────────────────────────────────────────────────────┐
│  OpenAutonomyX - Three Layers Deployed              │
├─────────────────────────────────────────────────────┤
│                                                     │
│  🎨 REACT FRONTEND (Vercel)                        │
│     https://app.publishing.openautonomyx.com        │
│     ✅ Auto-deploys on git push                     │
│     ✅ Global CDN                                   │
│     ✅ SSL included                                 │
│                                                     │
│  📖 MARKETING SITE (GitHub Pages)                  │
│     https://publishing.openautonomyx.com            │
│     ✅ Live immediately                             │
│     ✅ Updates on git push                          │
│     ✅ Free hosting                                 │
│                                                     │
│  👥 ENTERPRISE PORTAL (Liferay + Docker)           │
│     https://publishing.openautonomyx.com:8080       │
│     ✅ Full content creation                        │
│     ✅ Portlet architecture                         │
│     ✅ Multi-tenant ready                           │
│                                                     │
│  🎉 ALL SYSTEMS ONLINE                             │
│                                                     │
└─────────────────────────────────────────────────────┘
```

---

## Support

- 📖 Full docs: [DEPLOY-LANDING-PAGES.md](DEPLOY-LANDING-PAGES.md)
- 🔧 Architecture: [PLATFORM-ARCHITECTURE.md](PLATFORM-ARCHITECTURE.md)
- 🚀 Launch checklist: [FINAL-LAUNCH-CHECKLIST.md](FINAL-LAUNCH-CHECKLIST.md)
- 📍 GitHub: [Open-Autonomyx/Publishing-Platform](https://github.com/Open-Autonomyx/Publishing-Platform)

---

**🎯 Production-ready. Ship-ready. Go! 🚀**
