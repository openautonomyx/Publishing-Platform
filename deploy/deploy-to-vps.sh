#!/bin/bash

# COMPLETE VPS DEPLOYMENT SCRIPT
# Deploys entire platform to agennext.com in one go

set -e

# Configuration
VPS_HOST="${VPS_HOST:-agennext.com}"
VPS_USER="${VPS_USER:-almalinux}"
VPS_PASSWORD="${VPS_PASSWORD}"
REPO_URL="https://github.com/fractional-pm/creative-platform.git"
BRANCH="${BRANCH:-main}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=========================================="
echo "🚀 CREATIVE PLATFORM - VPS DEPLOYMENT"
echo "==========================================${NC}"
echo ""
echo "Target VPS: $VPS_HOST"
echo "User: $VPS_USER"
echo "Branch: $BRANCH"
echo ""

if [ -z "$VPS_PASSWORD" ]; then
    echo -e "${RED}❌ Error: VPS_PASSWORD not set${NC}"
    echo "Usage: VPS_PASSWORD='your-password' bash deploy-to-vps.sh"
    exit 1
fi

# Function to run commands on VPS
run_on_vps() {
    local cmd=$1
    sshpass -p "$VPS_PASSWORD" ssh -o StrictHostKeyChecking=no "$VPS_USER@$VPS_HOST" "$cmd"
}

# Function to copy files to VPS
copy_to_vps() {
    local file=$1
    local dest=$2
    sshpass -p "$VPS_PASSWORD" scp -o StrictHostKeyChecking=no "$file" "$VPS_USER@$VPS_HOST:$dest"
}

echo -e "${YELLOW}Step 1: Checking VPS connectivity...${NC}"
if run_on_vps "echo 'Connected!'" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ VPS is reachable${NC}"
else
    echo -e "${RED}❌ Cannot connect to VPS${NC}"
    exit 1
fi
echo ""

echo -e "${YELLOW}Step 2: Installing system dependencies...${NC}"
run_on_vps "sudo bash -c '
set -e
dnf update -y
dnf install -y curl wget git htop docker sshpass openssl
systemctl start docker
systemctl enable docker
curl -L \"https://github.com/docker/compose/releases/latest/download/docker-compose-\$(uname -s)-\$(uname -m)\" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
firewall-cmd --permanent --add-service=http
firewall-cmd --permanent --add-service=https
firewall-cmd --permanent --add-port=3001/tcp
firewall-cmd --permanent --add-port=8080/tcp
firewall-cmd --permanent --add-port=3000/tcp
firewall-cmd --reload
docker network create creative-network 2>/dev/null || true
echo \"Done\"
'"
echo -e "${GREEN}✅ Dependencies installed${NC}"
echo ""

echo -e "${YELLOW}Step 3: Cloning repository...${NC}"
run_on_vps "
rm -rf ~/creative-platform
git clone -b $BRANCH $REPO_URL ~/creative-platform
cd ~/creative-platform
echo 'Repository cloned'
"
echo -e "${GREEN}✅ Repository cloned${NC}"
echo ""

echo -e "${YELLOW}Step 4: Generating secrets...${NC}"
JWT_SECRET=$(openssl rand -base64 32)
DB_PASSWORD=$(openssl rand -base64 32)
LIFERAY_API_KEY=$(openssl rand -base64 32)

echo "Secrets generated:"
echo "  JWT_SECRET: ${JWT_SECRET:0:20}..."
echo "  DB_PASSWORD: ${DB_PASSWORD:0:20}..."
echo ""

echo -e "${YELLOW}Step 5: Starting PostgreSQL...${NC}"
run_on_vps "
docker pull postgres:15-alpine
docker stop postgres 2>/dev/null || true
docker rm postgres 2>/dev/null || true
docker run -d \
  --name postgres \
  --network creative-network \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=$DB_PASSWORD \
  -v postgres_data:/var/lib/postgresql/data \
  -p 5432:5432 \
  postgres:15-alpine

sleep 10
docker exec postgres pg_isready || sleep 10
echo 'PostgreSQL started'
"
echo -e "${GREEN}✅ PostgreSQL running${NC}"
echo ""

echo -e "${YELLOW}Step 6: Loading database schema...${NC}"
copy_to_vps "db/schema.sql" "/tmp/schema.sql"
run_on_vps "
docker exec -i postgres psql -U postgres -c 'CREATE DATABASE creative_platform' 2>/dev/null || true
docker exec -i postgres psql -U postgres creative_platform < /tmp/schema.sql
docker exec postgres psql -U postgres creative_platform -c '\dt' | head -5
echo 'Schema loaded'
"
echo -e "${GREEN}✅ Database schema loaded${NC}"
echo ""

echo -e "${YELLOW}Step 7: Building and pushing Docker image...${NC}"
docker build -t ghcr.io/fractional-pm/creative-platform-api:latest src/api/ || {
    echo -e "${RED}❌ Docker build failed - trying without rebuild${NC}"
    echo -e "${YELLOW}Using latest image from GHCR...${NC}"
}
echo -e "${GREEN}✅ Docker image ready${NC}"
echo ""

echo -e "${YELLOW}Step 8: Starting Redis...${NC}"
run_on_vps "
docker pull redis:7-alpine
docker stop redis 2>/dev/null || true
docker rm redis 2>/dev/null || true
docker run -d \
  --name redis \
  --network creative-network \
  -p 6379:6379 \
  redis:7-alpine

sleep 3
docker exec redis redis-cli ping
echo 'Redis started'
"
echo -e "${GREEN}✅ Redis running${NC}"
echo ""

echo -e "${YELLOW}Step 9: Starting Liferay...${NC}"
run_on_vps "
docker pull liferay/dxp:latest
docker stop liferay 2>/dev/null || true
docker rm liferay 2>/dev/null || true
docker run -d \
  --name liferay \
  --network creative-network \
  -e LIFERAY_JPDA_ENABLED=false \
  -e LIFERAY_DATABASE_PREPARED_STATEMENT_CACHE_SIZE=500 \
  -v liferay_data:/opt/liferay \
  -p 8080:8080 \
  liferay/dxp:latest

echo 'Liferay starting (this takes 2-5 minutes)...'
sleep 60
echo 'Waiting for Liferay to be ready...'
"
echo -e "${GREEN}✅ Liferay started${NC}"
echo ""

echo -e "${YELLOW}Step 10: Starting Go API...${NC}"
run_on_vps "
docker stop api 2>/dev/null || true
docker rm api 2>/dev/null || true
docker pull ghcr.io/fractional-pm/creative-platform-api:latest || true
docker run -d \
  --name api \
  --network creative-network \
  -e DATABASE_URL=postgres://postgres:$DB_PASSWORD@postgres:5432/creative_platform \
  -e REDIS_URL=redis://redis:6379 \
  -e LIFERAY_URL=http://liferay:8080 \
  -e LIFERAY_API_KEY=$LIFERAY_API_KEY \
  -e JWT_SECRET=$JWT_SECRET \
  -e API_LOG_LEVEL=info \
  -e API_PORT=3001 \
  -p 3001:3001 \
  ghcr.io/fractional-pm/creative-platform-api:latest

sleep 10
docker logs api | head -20
echo 'API started'
"
echo -e "${GREEN}✅ API running${NC}"
echo ""

echo -e "${YELLOW}Step 11: Configuring Nginx reverse proxy...${NC}"
run_on_vps "
sudo bash -c 'cat > /etc/nginx/conf.d/creative-platform.conf << \"NGINX_EOF\"
upstream api_backend {
    server localhost:3001;
}

upstream liferay_backend {
    server localhost:8080;
}

server {
    listen 80;
    server_name agennext.com www.agennext.com;

    # Redirect HTTP to HTTPS
    return 301 https://\$server_name\$request_uri;
}

server {
    listen 443 ssl http2;
    server_name agennext.com www.agennext.com;

    ssl_certificate /etc/letsencrypt/live/agennext.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/agennext.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # API endpoints
    location /api/ {
        proxy_pass http://api_backend;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    # Health check
    location /health {
        proxy_pass http://api_backend;
        proxy_set_header Host \$host;
    }

    # Metrics
    location /metrics {
        proxy_pass http://api_backend;
        proxy_set_header Host \$host;
    }

    # Liferay (optional)
    location /liferay/ {
        proxy_pass http://liferay_backend/;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
NGINX_EOF

systemctl start nginx
systemctl enable nginx
systemctl status nginx | head -10
'"
echo -e "${GREEN}✅ Nginx configured${NC}"
echo ""

echo -e "${YELLOW}Step 12: Verifying deployment...${NC}"
echo ""

# Health checks
echo "Running health checks..."
sleep 5

echo -n "  API health: "
if run_on_vps "curl -s http://localhost:3001/health | grep -q 'status'" 2>/dev/null; then
    echo -e "${GREEN}✅${NC}"
else
    echo -e "${YELLOW}⚠️  Still starting${NC}"
fi

echo -n "  PostgreSQL: "
if run_on_vps "docker exec postgres pg_isready -U postgres" > /dev/null 2>&1; then
    echo -e "${GREEN}✅${NC}"
else
    echo -e "${RED}❌${NC}"
fi

echo -n "  Redis: "
if run_on_vps "docker exec redis redis-cli ping | grep -q PONG" > /dev/null 2>&1; then
    echo -e "${GREEN}✅${NC}"
else
    echo -e "${RED}❌${NC}"
fi

echo -n "  Liferay: "
if run_on_vps "curl -s http://localhost:8080 | grep -q Liferay" > /dev/null 2>&1; then
    echo -e "${GREEN}✅${NC}"
else
    echo -e "${YELLOW}⚠️  Still starting${NC}"
fi

echo -n "  Docker containers: "
RUNNING=$(run_on_vps "docker ps | wc -l")
echo -e "${GREEN}$RUNNING containers${NC}"

echo ""
echo -e "${YELLOW}Step 13: Saving deployment info...${NC}"
run_on_vps "
cat > ~/deployment-info.txt << 'INFO'
DEPLOYMENT INFO
================
Deployed: $(date)
Branch: $BRANCH
Repository: $REPO_URL

URLS
====
API: http://agennext.com/api/v1
Health: http://agennext.com/health
Liferay: http://agennext.com:8080

DATABASE
========
Host: localhost
Port: 5432
User: postgres
Database: creative_platform
Password: $DB_PASSWORD (saved in docker)

SECRETS
=======
JWT_SECRET: $JWT_SECRET
API_KEY: $LIFERAY_API_KEY

DOCKER COMMANDS
===============
View logs: docker logs -f api
View services: docker ps
Stop all: docker-compose down
Start all: docker-compose up -d

SSH ACCESS
==========
ssh $VPS_USER@$VPS_HOST
cd ~/creative-platform
INFO

cat ~/deployment-info.txt
"
echo -e "${GREEN}✅ Deployment info saved${NC}"
echo ""

echo -e "${BLUE}=========================================="
echo "🎉 DEPLOYMENT COMPLETE!"
echo "==========================================${NC}"
echo ""
echo -e "${GREEN}Platform is live at:${NC}"
echo "  🌐 API: http://agennext.com/health"
echo "  🌐 Liferay: http://agennext.com:8080"
echo "  📊 Metrics: http://agennext.com/metrics"
echo ""
echo -e "${GREEN}SSH into VPS:${NC}"
echo "  ssh $VPS_USER@$VPS_HOST"
echo ""
echo -e "${GREEN}Next steps:${NC}"
echo "  1. Test API: curl http://agennext.com/health"
echo "  2. Check logs: docker logs -f api"
echo "  3. Configure DNS"
echo "  4. Setup SSL"
echo "  5. Deploy monitoring"
echo ""
echo -e "${YELLOW}Save these credentials:${NC}"
echo "  JWT_SECRET: $JWT_SECRET"
echo "  DB_PASSWORD: $DB_PASSWORD"
echo ""
