#!/bin/bash

##############################################################################
# OpenAutonomyX Podman Deployment (AlmaLinux)
# Lightweight container runtime - no Docker daemon needed
##############################################################################

set -e

echo "╔════════════════════════════════════════════════════════════╗"
echo "║  OpenAutonomyX - Podman Deployment (AlmaLinux)             ║"
echo "║  Lightweight Production Setup                              ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

REPO_URL="https://github.com/openautonomyx/original.git"
DEPLOY_DIR="/opt/openautonomyx"
APP_PORT="3000"
DOMAIN="${1:-openautonomyx.com}"

echo -e "${YELLOW}📋 Configuration${NC}"
echo "Repository: $REPO_URL"
echo "Deploy Path: $DEPLOY_DIR"
echo "Domain: $DOMAIN"
echo ""

if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}❌ This script must be run as root${NC}"
   exit 1
fi

echo -e "${YELLOW}🔄 Step 1: Update System${NC}"
dnf update -y
dnf install -y curl wget git

echo -e "${YELLOW}🔄 Step 2: Install Podman${NC}"
dnf install -y podman podman-compose

echo -e "${YELLOW}🔄 Step 3: Enable Podman Socket (for rootless)${NC}"
systemctl enable podman.socket
systemctl start podman.socket

echo -e "${YELLOW}🔄 Step 4: Install Nginx${NC}"
dnf install -y nginx

echo -e "${YELLOW}🔄 Step 5: Install Certbot (SSL)${NC}"
dnf install -y certbot python3-certbot-nginx

echo -e "${YELLOW}🔄 Step 6: Disable SELinux${NC}"
setenforce 0 || true
sed -i 's/^SELINUX=enforcing$/SELINUX=disabled/' /etc/selinux/config

echo -e "${YELLOW}🔄 Step 7: Configure Firewall${NC}"
systemctl start firewalld
systemctl enable firewalld
firewall-cmd --permanent --add-service=http
firewall-cmd --permanent --add-service=https
firewall-cmd --permanent --add-port=$APP_PORT/tcp
firewall-cmd --reload

echo -e "${YELLOW}🔄 Step 8: Clone Repository${NC}"
mkdir -p $DEPLOY_DIR
cd $DEPLOY_DIR
if [ -d ".git" ]; then
    git pull origin main
else
    git clone $REPO_URL .
fi

echo -e "${YELLOW}🔄 Step 9: Create Environment File${NC}"
cat > $DEPLOY_DIR/.env << EOF
NODE_ENV=production
PORT=$APP_PORT

# Database
DATABASE_TYPE=postgresql
DATABASE_URL=postgresql://openautonomyx:secure_password@postgres:5432/openautonomyx

# Cache
REDIS_URL=redis://redis:6379

# LLM Configuration
LLM_PROVIDER=ollama
LLM_MODEL=mistral
OLLAMA_API_URL=http://ollama:11434

# Storage
STORAGE_TYPE=local
STORAGE_PATH=/data/uploads

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Security
JWT_SECRET=$(openssl rand -base64 32)
API_KEY_PREFIX=sk-openautonomyx

# Domain
DOMAIN=$DOMAIN
EOF

echo -e "${GREEN}✅ .env file created${NC}"

echo -e "${YELLOW}🔄 Step 10: Start Podman Containers${NC}"
cd $DEPLOY_DIR

# Create pod for all services
podman pod create -n openautonomyx \
    -p 5432:5432 \
    -p 6379:6379 \
    -p 11434:11434 \
    -p 5050:80 \
    -p 8081:8081 \
    -p $APP_PORT:$APP_PORT

# PostgreSQL
echo "Starting PostgreSQL..."
podman run -d \
    --name openautonomyx-postgres \
    --pod openautonomyx \
    -e POSTGRES_USER=openautonomyx \
    -e POSTGRES_PASSWORD=secure_password \
    -e POSTGRES_DB=openautonomyx \
    -v postgres_data:/var/lib/postgresql/data \
    docker.io/library/postgres:15-alpine

# Redis
echo "Starting Redis..."
podman run -d \
    --name openautonomyx-redis \
    --pod openautonomyx \
    -v redis_data:/data \
    docker.io/library/redis:7-alpine \
    redis-server --appendonly yes

# Ollama
echo "Starting Ollama..."
podman run -d \
    --name openautonomyx-ollama \
    --pod openautonomyx \
    -v ollama_data:/root/.ollama \
    docker.io/ollama/ollama:latest

# Application
echo "Starting Application..."
podman run -d \
    --name openautonomyx-app \
    --pod openautonomyx \
    -v $DEPLOY_DIR/public:/app/public \
    -v $DEPLOY_DIR/docs:/app/docs \
    --env-file .env \
    docker.io/library/node:18-alpine \
    sh -c "cd /app && npm install && npm start"

sleep 5
echo -e "${YELLOW}🔄 Step 11: Pull LLM Model${NC}"
podman exec openautonomyx-ollama ollama pull mistral || true

echo -e "${YELLOW}🔄 Step 12: Configure Nginx${NC}"
cat > /etc/nginx/conf.d/openautonomyx.conf << EOF
upstream openautonomyx_backend {
    server localhost:$APP_PORT;
}

server {
    listen 80;
    server_name $DOMAIN www.$DOMAIN;
    client_max_body_size 100M;

    location / {
        proxy_pass http://openautonomyx_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_cache_bypass \$http_upgrade;
    }
}
EOF

nginx -t
systemctl start nginx
systemctl enable nginx

echo -e "${YELLOW}🔄 Step 13: Setup SSL Certificate${NC}"
certbot --nginx -d $DOMAIN -d www.$DOMAIN --non-interactive --agree-tos -m admin@$DOMAIN || true

echo -e "${YELLOW}🔄 Step 14: Create Systemd Service${NC}"
cat > /etc/systemd/system/openautonomyx.service << EOF
[Unit]
Description=OpenAutonomyX Podman Services
After=network-online.target
Wants=network-online.target

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=$DEPLOY_DIR
ExecStart=/usr/bin/podman pod start openautonomyx
ExecStop=/usr/bin/podman pod stop openautonomyx
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable openautonomyx.service

echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo -e "${GREEN}✅ Deployment Complete!${NC}"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""
echo "🌐 Your platform is live at:"
echo -e "${GREEN}   https://$DOMAIN${NC}"
echo ""
echo "📊 Services:"
echo "   • App: http://localhost:$APP_PORT"
echo "   • PostgreSQL: localhost:5432"
echo "   • Redis: localhost:6379"
echo "   • Ollama: localhost:11434"
echo ""
echo "📋 Podman Commands:"
echo "   • View containers: podman pod ps"
echo "   • View logs: podman pod logs openautonomyx"
echo "   • Container logs: podman logs openautonomyx-app"
echo "   • Stop all: systemctl stop openautonomyx"
echo "   • Start all: systemctl start openautonomyx"
echo ""
echo "🔄 Auto-start on Boot:"
echo "   systemctl enable openautonomyx.service"
echo ""
echo "🔒 Security Checklist:"
echo "   ☐ Update .env file with secure passwords"
echo "   ☐ Verify firewall rules"
echo "   ☐ Set up backups"
echo ""
