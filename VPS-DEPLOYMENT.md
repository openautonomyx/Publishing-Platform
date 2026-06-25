# VPS Deployment Guide

## 🚀 Quick Start (5 minutes)

### Step 1: Get a VPS
Choose one of these providers:

| Provider | Cost | Features |
|----------|------|----------|
| **DigitalOcean** | $5-12/mo | Simple, reliable, great for beginners |
| **Linode** | $5-10/mo | Powerful, good support |
| **AWS EC2** | Variable | Scale, enterprise features |
| **Hetzner** | €3-5/mo | Cheapest, good performance |
| **Vultr** | $2.50+/mo | Global, flexible |

**Recommended for beginners:** DigitalOcean (Droplet)
- Ubuntu 20.04 or 22.04 LTS
- 2GB RAM minimum
- 50GB SSD

### Step 2: SSH into Your VPS
```bash
ssh root@your-vps-ip-address
```

### Step 3: Download & Run Deployment Script
```bash
curl -fsSL https://raw.githubusercontent.com/openautonomyx/original/main/deploy-vps.sh -o deploy-vps.sh
sudo bash deploy-vps.sh openautonomyx.com
```

Replace `openautonomyx.com` with your domain.

### Step 4: Configure DNS
Point your domain to the VPS IP:

```
A Record: your-domain.com  →  your-vps-ip
A Record: www.your-domain.com  →  your-vps-ip
```

### Step 5: Done! ✅
Your platform is live at: `https://your-domain.com`

---

## 📋 What the Script Does

```
1. Updates system packages
2. Installs Docker & Docker Compose
3. Installs Nginx (reverse proxy)
4. Installs Certbot (SSL certificates)
5. Clones the repository
6. Creates .env configuration
7. Starts Docker containers (app, DB, cache, LLM)
8. Configures Nginx
9. Generates free SSL certificate
10. Sets up auto-renewal
```

---

## 🔧 Manual Setup (If Script Doesn't Work)

### Prerequisites
```bash
apt-get update && apt-get upgrade -y
apt-get install -y curl wget git
```

### Install Docker
```bash
curl -fsSL https://get.docker.com | sh
systemctl start docker
systemctl enable docker
```

### Install Docker Compose
```bash
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
```

### Clone & Deploy
```bash
mkdir -p /opt/openautonomyx
cd /opt/openautonomyx
git clone https://github.com/openautonomyx/original.git .

# Create .env file
nano .env
# Add configuration (see below)

# Start services
docker-compose up -d

# Pull LLM model
docker exec openautonomyx-ollama ollama pull mistral
```

### Configure Nginx
```bash
apt-get install -y nginx certbot python3-certbot-nginx

# Create nginx config
nano /etc/nginx/sites-available/openautonomyx
```

Paste this:
```nginx
upstream app {
    server localhost:3000;
}

server {
    listen 80;
    server_name your-domain.com www.your-domain.com;

    location / {
        proxy_pass http://app;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable & restart:
```bash
ln -s /etc/nginx/sites-available/openautonomyx /etc/nginx/sites-enabled/
systemctl restart nginx
```

### Setup SSL
```bash
certbot --nginx -d your-domain.com
```

---

## 🔐 Environment Configuration

Create `.env` file with:

```bash
# Core
NODE_ENV=production
PORT=3000

# Database
DATABASE_TYPE=postgresql
DATABASE_URL=postgresql://openautonomyx:password@postgres:5432/openautonomyx

# Cache
REDIS_URL=redis://redis:6379

# LLM
LLM_PROVIDER=ollama
LLM_MODEL=mistral
OLLAMA_API_URL=http://ollama:11434

# Storage
STORAGE_TYPE=local
STORAGE_PATH=/data/uploads

# Security
JWT_SECRET=generate-random-string-here
API_KEY_PREFIX=sk-openautonomyx

# Domain
DOMAIN=your-domain.com
```

---

## 📊 Service Status

Check if everything is running:

```bash
# View all containers
docker ps

# View logs
docker-compose logs -f app

# Check specific service
docker-compose logs postgres
docker-compose logs redis
docker-compose logs ollama
```

---

## 🌐 Access Your Services

| Service | URL |
|---------|-----|
| **App** | https://your-domain.com |
| **API** | https://your-domain.com/api |
| **Health** | https://your-domain.com/api/health |
| **Admin DB** | http://your-vps-ip:5050 (pgAdmin) |
| **Redis** | your-vps-ip:6379 (internal only) |
| **Ollama API** | http://your-vps-ip:11434 |

---

## 🔒 Security Checklist

- [ ] Change PostgreSQL password in `.env`
- [ ] Change pgAdmin default credentials
- [ ] Configure firewall rules:
  ```bash
  ufw allow 22/tcp  # SSH
  ufw allow 80/tcp  # HTTP
  ufw allow 443/tcp # HTTPS
  ufw enable
  ```
- [ ] Set up automated backups
- [ ] Monitor disk space: `df -h`
- [ ] Monitor memory: `free -m`
- [ ] Enable log rotation

---

## 📈 Common Commands

### View Logs
```bash
docker-compose logs -f app
docker-compose logs -f postgres
```

### Restart Services
```bash
docker-compose restart app
docker-compose restart
```

### Stop Everything
```bash
docker-compose down
```

### Start Again
```bash
docker-compose up -d
```

### SSH into Container
```bash
docker exec -it openautonomyx-app bash
docker exec -it openautonomyx-postgres psql -U openautonomyx
```

### Update Code
```bash
cd /opt/openautonomyx
git pull origin main
docker-compose restart app
```

### View Disk Usage
```bash
du -sh /opt/openautonomyx
docker system df
```

---

## 🆘 Troubleshooting

### Port already in use
```bash
lsof -i :3000
kill -9 <PID>
```

### Database connection error
```bash
docker-compose ps  # Check if postgres is running
docker-compose logs postgres
```

### SSL certificate error
```bash
certbot renew --dry-run
certbot renew
```

### App won't start
```bash
docker-compose logs app
# Check .env file for errors
```

### Out of disk space
```bash
df -h
docker system prune -a  # WARNING: Removes unused images
```

---

## 📞 Support

- **Docs:** https://openautonomyx.github.io/original
- **Issues:** https://github.com/openautonomyx/original/issues
- **Email:** support@openautonomyx.com

---

## 🎯 Next Steps

1. **Monitor:** Set up uptime monitoring (Pingdom, UptimeRobot)
2. **Backups:** Configure automated database backups
3. **Analytics:** Add Google Analytics or Plausible
4. **Scaling:** If traffic grows, add load balancer or extra servers
5. **CDN:** Use Cloudflare for global caching

---

**Your production platform is live!** 🚀
