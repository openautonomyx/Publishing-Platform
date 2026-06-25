#!/bin/bash

##############################################################################
# OpenAutonomyX Podman Deployment (Simple - No Compose)
# Direct podman commands - minimal dependencies
##############################################################################

set -e

echo "╔════════════════════════════════════════════════════════════╗"
echo "║  OpenAutonomyX - Podman Simple Deployment                  ║"
echo "║  Direct Podman Commands (No Compose)                       ║"
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
echo "Deployment Path: $DEPLOY_DIR"
echo "Domain: $DOMAIN"
echo "App Port: $APP_PORT"
echo ""

if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}❌ This script must be run as root${NC}"
   exit 1
fi

echo -e "${YELLOW}🔄 Step 1: Update System${NC}"
dnf update -y
dnf install -y curl wget git

echo -e "${YELLOW}🔄 Step 2: Install Podman${NC}"
dnf install -y podman

echo -e "${YELLOW}🔄 Step 3: Enable Podman Socket${NC}"
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
firewall-cmd --permanent --add-port=5432/tcp
firewall-cmd --permanent --add-port=6379/tcp
firewall-cmd --permanent --add-port=11434/tcp
firewall-cmd --reload

echo -e "${YELLOW}🔄 Step 8: Create Persistent Volumes${NC}"
podman volume create postgres_data || true
podman volume create redis_data || true
podman volume create ollama_data || true
podman volume create app_data || true

echo -e "${YELLOW}🔄 Step 9: Clone Repository${NC}"
mkdir -p $DEPLOY_DIR
cd $DEPLOY_DIR
if [ -d ".git" ]; then
    git pull origin main
else
    git clone $REPO_URL .
fi

echo -e "${YELLOW}🔄 Step 10: Create Environment File${NC}"
cat > $DEPLOY_DIR/.env << EOF
NODE_ENV=production
PORT=$APP_PORT
DATABASE_URL=postgresql://openautonomyx:secure_password@localhost:5432/openautonomyx
REDIS_URL=redis://localhost:6379
LLM_PROVIDER=ollama
LLM_MODEL=mistral
OLLAMA_API_URL=http://localhost:11434
STORAGE_TYPE=local
STORAGE_PATH=/app/data
LOG_LEVEL=info
JWT_SECRET=$(openssl rand -base64 32)
API_KEY_PREFIX=sk-openautonomyx
DOMAIN=$DOMAIN
EOF

echo -e "${GREEN}✅ Configuration ready${NC}"

echo -e "${YELLOW}🔄 Step 11: Create Pod${NC}"
podman pod create -n openautonomyx \
    -p 127.0.0.1:5432:5432 \
    -p 127.0.0.1:6379:6379 \
    -p 127.0.0.1:11434:11434 \
    -p 127.0.0.1:$APP_PORT:$APP_PORT \
    || echo "Pod already exists"

echo -e "${YELLOW}🔄 Step 12: Start PostgreSQL${NC}"
podman run -d \
    --name openautonomyx-postgres \
    --pod openautonomyx \
    -e POSTGRES_USER=openautonomyx \
    -e POSTGRES_PASSWORD=secure_password \
    -e POSTGRES_DB=openautonomyx \
    -v postgres_data:/var/lib/postgresql/data \
    --restart=on-failure:5 \
    postgres:15-alpine \
    || echo "PostgreSQL already running"

sleep 3

echo -e "${YELLOW}🔄 Step 13: Start Redis${NC}"
podman run -d \
    --name openautonomyx-redis \
    --pod openautonomyx \
    -v redis_data:/data \
    --restart=on-failure:5 \
    redis:7-alpine \
    redis-server --appendonly yes \
    || echo "Redis already running"

echo -e "${YELLOW}🔄 Step 14: Start Ollama${NC}"
podman run -d \
    --name openautonomyx-ollama \
    --pod openautonomyx \
    -v ollama_data:/root/.ollama \
    --restart=on-failure:5 \
    ollama/ollama:latest \
    || echo "Ollama already running"

sleep 5

echo -e "${YELLOW}🔄 Step 15: Pull LLM Model${NC}"
podman exec openautonomyx-ollama ollama pull mistral || true

echo -e "${YELLOW}🔄 Step 16: Start Application${NC}"
podman run -d \
    --name openautonomyx-app \
    --pod openautonomyx \
    -v $DEPLOY_DIR:/app \
    --env-file $DEPLOY_DIR/.env \
    --restart=on-failure:5 \
    node:18-alpine \
    sh -c "cd /app && npm install --production && npm start" \
    || echo "App already running"

sleep 3

echo -e "${YELLOW}🔄 Step 17: Configure Nginx${NC}"
cat > /etc/nginx/conf.d/openautonomyx.conf << EOF
upstream openautonomyx_backend {
    server 127.0.0.1:$APP_PORT;
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
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
}
EOF

nginx -t
systemctl start nginx
systemctl enable nginx

echo -e "${YELLOW}🔄 Step 18: Setup SSL Certificate${NC}"
certbot --nginx -d $DOMAIN -d www.$DOMAIN --non-interactive --agree-tos -m admin@$DOMAIN || true

echo -e "${YELLOW}🔄 Step 19: Create Systemd Service${NC}"
cat > /etc/systemd/system/openautonomyx.service << EOF
[Unit]
Description=OpenAutonomyX Podman Pod
After=network-online.target
Wants=network-online.target

[Service]
Type=oneshot
RemainAfterExit=yes
ExecStart=/usr/bin/podman pod start openautonomyx
ExecStop=/usr/bin/podman pod stop -t 10 openautonomyx
Restart=on-failure
RestartSec=10s

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
echo "📊 Running Services:"
podman pod ps
echo ""
podman pod top openautonomyx
echo ""
echo "📋 Useful Podman Commands:"
echo "   • View pod: podman pod ps"
echo "   • View logs: podman pod logs openautonomyx"
echo "   • App logs: podman logs openautonomyx-app"
echo "   • DB logs: podman logs openautonomyx-postgres"
echo "   • Stop all: systemctl stop openautonomyx"
echo "   • Start all: systemctl start openautonomyx"
echo "   • Restart: systemctl restart openautonomyx"
echo ""
echo "🔗 Services (local access):"
echo "   • App: http://localhost:$APP_PORT"
echo "   • PostgreSQL: localhost:5432"
echo "   • Redis: localhost:6379"
echo "   • Ollama: localhost:11434"
echo ""
echo "🔄 Auto-start on Boot:"
echo "   ✅ Enabled via systemd"
echo ""
echo "📝 Configuration:"
echo "   .env file: $DEPLOY_DIR/.env"
echo "   Update passwords in .env after first run!"
echo ""
echo "🆘 Troubleshooting:"
echo "   • Check status: systemctl status openautonomyx"
echo "   • View all: podman ps -a"
echo "   • Restart specific: podman restart openautonomyx-app"
echo ""
