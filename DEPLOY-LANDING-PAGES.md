# Deploy All Landing Pages - Complete Guide

**Production deployment for React Frontend, GitHub Pages, and Liferay Portal**

---

## 1. Deploy React Frontend to Vercel

### Step 1: Build Production Bundle

```bash
cd /Users/chinmaypanda/CustomApps/frontend

# Install dependencies
npm install

# Build for production
npm run build

# Output: ./dist/ directory
ls -la dist/
```

### Step 2: Deploy to Vercel

```bash
# Install Vercel CLI
npm install -g vercel

# Login to Vercel
vercel login

# Deploy
vercel --prod

# Configure
# - Project name: openautonomyx-platform
# - Framework: Vite
# - Root directory: ./
# - Build command: npm run build
# - Output directory: dist

# Result: https://openautonomyx-platform.vercel.app
```

### Step 3: Configure Custom Domain

```bash
# Add domain in Vercel dashboard
# Settings → Domains → Add
# - Add: app.publishing.openautonomyx.com
# - Update DNS records:
#   CNAME app.publishing.openautonomyx.com → cname.vercel-dns.com
```

### Step 4: Environment Variables

```bash
# .env.production
VITE_API_GATEWAY_URL=https://api.publishing.openautonomyx.com
VITE_BLOG_URL=https://blog.publishing.openautonomyx.com
VITE_FORMATS_URL=https://formats.publishing.openautonomyx.com
VITE_INTEGRATIONS_URL=https://integrations.publishing.openautonomyx.com

# Set in Vercel dashboard: Settings → Environment Variables
```

### Result
```
✅ React Frontend Live
🌐 https://app.publishing.openautonomyx.com
📊 Auto-deploys on git push to main
```

---

## 2. Enhance GitHub Pages Marketing Site

### Step 1: Create Professional Landing Page

```bash
cat > /Users/chinmaypanda/CustomApps/docs/index.html << 'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>OpenAutonomyX - Vendor-Neutral Creative Publishing Platform</title>
  <meta name="description" content="Create once, publish everywhere. AI-powered content generation for 50+ languages and 6 social platforms.">
  <style>
    * {
      margin: 0;
      padding: 0;
      box-sizing: border-box;
    }

    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      background: #0f172a;
      color: #e2e8f0;
      line-height: 1.6;
    }

    header {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      padding: 80px 40px;
      text-align: center;
      color: white;
    }

    header h1 {
      font-size: 48px;
      margin-bottom: 20px;
      line-height: 1.2;
    }

    header p {
      font-size: 20px;
      opacity: 0.95;
      margin-bottom: 40px;
    }

    .cta-buttons {
      display: flex;
      gap: 20px;
      justify-content: center;
      flex-wrap: wrap;
    }

    button {
      padding: 16px 40px;
      border: none;
      border-radius: 8px;
      font-size: 16px;
      font-weight: 600;
      cursor: pointer;
      transition: all 0.3s;
    }

    .btn-primary {
      background: white;
      color: #667eea;
    }

    .btn-primary:hover {
      transform: translateY(-3px);
      box-shadow: 0 15px 40px rgba(0,0,0,0.3);
    }

    .btn-secondary {
      background: rgba(255,255,255,0.2);
      color: white;
      border: 2px solid white;
    }

    .btn-secondary:hover {
      background: rgba(255,255,255,0.3);
    }

    .features {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
      gap: 30px;
      padding: 80px 40px;
      background: #1a1f35;
    }

    .feature {
      background: #1e293b;
      padding: 30px;
      border-radius: 12px;
      border: 1px solid #334155;
    }

    .feature h3 {
      margin-bottom: 15px;
      color: #f1f5f9;
      font-size: 20px;
    }

    .feature p {
      color: #94a3b8;
      font-size: 14px;
    }

    .stats {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 20px;
      padding: 60px 40px;
      background: #0f172a;
      text-align: center;
    }

    .stat-item h4 {
      font-size: 32px;
      color: #667eea;
      margin-bottom: 10px;
    }

    .stat-item p {
      color: #94a3b8;
      font-size: 14px;
    }

    footer {
      background: #1e293b;
      padding: 40px;
      text-align: center;
      border-top: 1px solid #334155;
      color: #94a3b8;
    }

    footer a {
      color: #667eea;
      text-decoration: none;
      margin: 0 15px;
    }

    footer a:hover {
      text-decoration: underline;
    }

    .container {
      max-width: 1200px;
      margin: 0 auto;
    }
  </style>
</head>
<body>
  <header>
    <div class="container">
      <h1>OpenAutonomyX</h1>
      <p>Create Once. Publish Everywhere.<br>In Any Language.</p>
      <p style="font-size: 16px; opacity: 0.9;">
        Vendor-neutral creative publishing platform<br>
        powered by local AI, unlimited formats, infinite possibilities
      </p>
      <div class="cta-buttons">
        <button class="btn-primary" onclick="window.location='https://app.publishing.openautonomyx.com'">
          Get Started Free
        </button>
        <button class="btn-secondary" onclick="document.getElementById('features').scrollIntoView()">
          Learn More
        </button>
      </div>
    </div>
  </header>

  <section class="features" id="features">
    <div class="feature">
      <h3>📝 Content Creation</h3>
      <p>Create rich, markdown-based content with version control and publishing workflow.</p>
    </div>
    <div class="feature">
      <h3>🎬 Multi-Format Export</h3>
      <p>Generate EPUB, PDF, Slides, Audio, Video, and HTML from single content.</p>
    </div>
    <div class="feature">
      <h3>🌐 Omnichannel Publishing</h3>
      <p>Publish to WordPress, Medium, Substack, Twitter, LinkedIn, Facebook, and more.</p>
    </div>
    <div class="feature">
      <h3>🤖 Local AI</h3>
      <p>Run Ollama locally or connect to OpenAI, Anthropic, Azure, and 50+ providers.</p>
    </div>
    <div class="feature">
      <h3>🔐 Vendor Neutral</h3>
      <p>Zero lock-in. Full control. SBOM-based supply chain security.</p>
    </div>
    <div class="feature">
      <h3>📊 Analytics</h3>
      <p>Unified analytics across all platforms with real-time dashboards.</p>
    </div>
  </section>

  <section class="stats">
    <div class="stat-item">
      <h4>21</h4>
      <p>Microservices</p>
    </div>
    <div class="stat-item">
      <h4>50+</h4>
      <p>Languages</p>
    </div>
    <div class="stat-item">
      <h4>6+</h4>
      <p>Export Formats</p>
    </div>
    <div class="stat-item">
      <h4>∞</h4>
      <p>Platforms</p>
    </div>
  </section>

  <footer>
    <div class="container">
      <p>&copy; 2026 OpenAutonomyX. MIT License.</p>
      <div style="margin-top: 20px;">
        <a href="https://github.com/open-autonomyx">GitHub</a>
        <a href="https://docs.publishing.openautonomyx.com">Documentation</a>
        <a href="https://discord.gg/openautonomyx">Discord</a>
        <a href="mailto:support@openautonomyx.com">Support</a>
      </div>
      <p style="margin-top: 20px; font-size: 12px;">
        Built with ❤️ by Claude AI | Anthropic
      </p>
    </div>
  </footer>
</body>
</html>
EOF
```

### Step 2: Enable GitHub Pages in Repository

```bash
# In GitHub repository settings:
# Settings → Pages
# - Source: Deploy from a branch
# - Branch: main
# - Folder: /docs
# - Custom domain: publishing.openautonomyx.com

# Update DNS:
# publishing.openautonomyx.com → CNAME → openautonomyx.github.io.
```

### Step 3: Update DNS Records

```bash
# In your domain registrar:

# CNAME Records:
publishing.openautonomyx.com  CNAME  openautonomyx.github.io
app.publishing.openautonomyx.com  CNAME  cname.vercel-dns.com
api.publishing.openautonomyx.com  CNAME  [your-api-ip-or-domain]
blog.publishing.openautonomyx.com  CNAME  [your-blog-domain]

# A Records (if using VPS):
@  A  [your-vps-ip]
api  A  [your-vps-ip]
blog  A  [your-vps-ip]
```

### Result
```
✅ Marketing Site Live
🌐 https://publishing.openautonomyx.com
📚 Full documentation
```

---

## 3. Deploy Liferay Portal to Production

### Step 1: Prepare Liferay for Production

```bash
# Create production docker-compose

cat > /tmp/liferay-production.yml << 'EOF'
version: '3.9'

services:
  liferay:
    image: liferay/dxp:latest
    container_name: liferay-prod
    environment:
      LIFERAY_JDBC_ONE_URL: jdbc:postgresql://postgres:5432/liferay
      LIFERAY_JDBC_ONE_DRIVER_CLASS_NAME: org.postgresql.Driver
      LIFERAY_JDBC_ONE_USERNAME: ${DB_USER}
      LIFERAY_JDBC_ONE_PASSWORD: ${DB_PASS}
      LIFERAY_SETUP_WIZARD_ENABLED: "false"
    ports:
      - "8080:8080"
    volumes:
      - liferay_data:/opt/liferay/data
      - ./liferay-config.properties:/opt/liferay/portal-ext.properties
    depends_on:
      - postgres
    restart: unless-stopped
    networks:
      - prod-network

  postgres:
    image: postgres:15-alpine
    container_name: postgres-prod
    environment:
      POSTGRES_DB: liferay
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASS}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped
    networks:
      - prod-network

  nginx:
    image: nginx:latest
    container_name: nginx-prod
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - liferay
    restart: unless-stopped
    networks:
      - prod-network

volumes:
  liferay_data:
  postgres_data:

networks:
  prod-network:
    driver: bridge
EOF
```

### Step 2: Configure Nginx Reverse Proxy

```bash
cat > /tmp/nginx.conf << 'EOF'
events {
  worker_connections 1024;
}

http {
  upstream liferay {
    server liferay:8080;
  }

  server {
    listen 80;
    server_name publishing.openautonomyx.com;
    return 301 https://$server_name$request_uri;
  }

  server {
    listen 443 ssl http2;
    server_name publishing.openautonomyx.com;

    ssl_certificate /etc/nginx/ssl/tls.crt;
    ssl_certificate_key /etc/nginx/ssl/tls.key;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    client_max_body_size 100M;

    location / {
      proxy_pass http://liferay;
      proxy_set_header Host $host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto $scheme;
      proxy_read_timeout 600s;
      proxy_connect_timeout 600s;
    }
  }
}
EOF
```

### Step 3: Generate SSL Certificates

```bash
# Using Let's Encrypt (certbot)
sudo apt-get install -y certbot python3-certbot-nginx

# Generate certificate
sudo certbot certonly --standalone \
  -d publishing.openautonomyx.com \
  -d app.publishing.openautonomyx.com \
  --email support@openautonomyx.com

# Copy to Docker volume
sudo cp /etc/letsencrypt/live/publishing.openautonomyx.com/* ./ssl/
```

### Step 4: Deploy to Production VPS

```bash
# SSH into VPS
ssh root@your-vps-ip

# Clone repository
git clone https://github.com/Open-Autonomyx/Publishing-Platform.git
cd Publishing-Platform

# Create .env file
cat > .env << 'EOF'
DB_USER=postgres
DB_PASS=your_secure_password
ENVIRONMENT=production
EOF

# Start services
docker-compose -f liferay-production.yml up -d

# Verify
docker-compose -f liferay-production.yml ps
```

### Step 5: Configure DNS & SSL

```bash
# Update DNS for Liferay domain
publishing.openautonomyx.com  A  [your-vps-ip]

# Verify SSL
curl -I https://publishing.openautonomyx.com
```

### Result
```
✅ Liferay Portal Live
🌐 https://publishing.openautonomyx.com
👥 Enterprise portal with portlets
📊 Full content creation & publishing
```

---

## Complete Deployment Summary

```
┌─────────────────────────────────────────┐
│  OpenAutonomyX - Production Ready        │
├─────────────────────────────────────────┤
│                                         │
│  🎨 React Frontend                      │
│  https://app.publishing.openautonomyx.com
│  Platform: Vercel                       │
│                                         │
│  📖 Marketing Site                      │
│  https://publishing.openautonomyx.com   │
│  (GitHub Pages)                         │
│                                         │
│  🔌 API Gateway                         │
│  https://api.publishing.openautonomyx.com
│  Platform: VPS / Kubernetes             │
│                                         │
│  👥 Enterprise Portal (Liferay)        │
│  https://publishing.openautonomyx.com:8080
│  Platform: VPS / Docker                 │
│                                         │
│  📝 Blog                                │
│  https://blog.publishing.openautonomyx.com
│  (WordPress on Hostinger)              │
│                                         │
│  ✅ All Systems Online                  │
│                                         │
└─────────────────────────────────────────┘
```

---

## Post-Deployment Checklist

- [ ] React frontend accessible
- [ ] GitHub Pages live
- [ ] Liferay portal running
- [ ] SSL certificates valid
- [ ] DNS records propagated
- [ ] API Gateway connected
- [ ] Database backed up
- [ ] Monitoring alerts active
- [ ] Email notifications working
- [ ] Support channels ready

---

## Monitoring & Maintenance

```bash
# Monitor uptime
curl -I https://publishing.openautonomyx.com
curl -I https://app.publishing.openautonomyx.com

# Check Liferay logs
docker-compose -f liferay-production.yml logs -f liferay

# Backup database
docker-compose -f liferay-production.yml exec postgres \
  pg_dump -U postgres liferay | gzip > backup-$(date +%Y%m%d).sql.gz

# Scale services
docker-compose -f liferay-production.yml up -d --scale content=3
```

---

🚀 **Platform Live! Ready to Serve Millions!**
