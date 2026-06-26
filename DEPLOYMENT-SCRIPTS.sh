#!/bin/bash

# OpenAutonomyX Complete Deployment Scripts
# Automates React Frontend, GitHub Pages, and Liferay Portal deployment

set -e

COLOR_GREEN='\033[0;32m'
COLOR_BLUE='\033[0;34m'
COLOR_YELLOW='\033[1;33m'
COLOR_RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${COLOR_BLUE}🚀 OpenAutonomyX Deployment Suite${NC}\n"

# ===== 1. DEPLOY REACT FRONTEND TO VERCEL =====

deploy_react_vercel() {
  echo -e "${COLOR_BLUE}📦 Step 1: Deploy React Frontend to Vercel${NC}\n"

  cd /Users/chinmaypanda/CustomApps/frontend

  # Check if Vercel CLI is installed
  if ! command -v vercel &> /dev/null; then
    echo -e "${COLOR_YELLOW}Installing Vercel CLI...${NC}"
    npm install -g vercel
  fi

  # Build production bundle
  echo -e "${COLOR_YELLOW}Building production bundle...${NC}"
  npm install
  npm run build

  if [ -d "dist" ]; then
    echo -e "${COLOR_GREEN}✅ Production build complete${NC}"
    ls -lh dist/ | head -10
  else
    echo -e "${COLOR_RED}❌ Build failed${NC}"
    return 1
  fi

  # Deploy to Vercel
  echo -e "${COLOR_YELLOW}Deploying to Vercel...${NC}"
  vercel --prod \
    --name openautonomyx-platform \
    --region sfo1

  echo -e "${COLOR_GREEN}✅ React Frontend deployed to Vercel${NC}\n"
}

# ===== 2. ENABLE GITHUB PAGES =====

enable_github_pages() {
  echo -e "${COLOR_BLUE}📖 Step 2: Enable GitHub Pages${NC}\n"

  cd /Users/chinmaypanda/CustomApps

  echo -e "${COLOR_YELLOW}GitHub Pages setup:${NC}"
  echo "1. Go to: https://github.com/Open-Autonomyx/Publishing-Platform/settings/pages"
  echo "2. Under 'Source', select: Deploy from a branch"
  echo "3. Select branch: main"
  echo "4. Select folder: /docs"
  echo "5. Add custom domain: publishing.openautonomyx.com"
  echo ""
  echo -e "${COLOR_YELLOW}Update DNS records at your registrar:${NC}"
  echo "publishing.openautonomyx.com  CNAME  openautonomyx.github.io."
  echo ""
  echo "Waiting for manual setup... (Press Enter when done)"
  read

  echo -e "${COLOR_GREEN}✅ GitHub Pages configured${NC}\n"
}

# ===== 3. SETUP SSL CERTIFICATES =====

setup_ssl_certificates() {
  echo -e "${COLOR_BLUE}🔐 Step 3: Setup SSL Certificates${NC}\n"

  echo -e "${COLOR_YELLOW}Installing certbot...${NC}"
  # Check if running on macOS or Linux
  if [[ "$OSTYPE" == "darwin"* ]]; then
    brew install certbot
  else
    sudo apt-get install -y certbot python3-certbot-nginx
  fi

  echo -e "${COLOR_YELLOW}Generating SSL certificates...${NC}"
  sudo certbot certonly --standalone \
    -d publishing.openautonomyx.com \
    -d app.publishing.openautonomyx.com \
    --email support@openautonomyx.com \
    --agree-tos --no-eff-email

  # Create SSL directory
  mkdir -p /Users/chinmaypanda/CustomApps/ssl

  echo -e "${COLOR_GREEN}✅ SSL certificates generated${NC}\n"
}

# ===== 4. PREPARE LIFERAY DOCKER COMPOSE =====

prepare_liferay_docker() {
  echo -e "${COLOR_BLUE}🔌 Step 4: Prepare Liferay Docker Setup${NC}\n"

  cd /Users/chinmaypanda/CustomApps

  # Create .env file for production
  cat > .env.production << 'EOF'
# Production Environment Variables
ENVIRONMENT=production
DB_USER=liferay_user
DB_PASS=your_secure_password_here
LIFERAY_ADMIN_PASSWORD=your_admin_password_here
API_GATEWAY_URL=https://api.publishing.openautonomyx.com
EVENT_BUS_URL=https://event-bus.publishing.openautonomyx.com
EOF

  echo -e "${COLOR_YELLOW}⚠️  IMPORTANT: Update .env.production with your passwords!${NC}"
  echo "File: /Users/chinmaypanda/CustomApps/.env.production"

  # Create production docker-compose
  cat > docker-compose.prod.yml << 'EOF'
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
    depends_on:
      - postgres
    restart: unless-stopped
    networks:
      - prod-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/"]
      interval: 30s
      timeout: 10s
      retries: 3

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
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  liferay_data:
  postgres_data:

networks:
  prod-network:
    driver: bridge
EOF

  echo -e "${COLOR_GREEN}✅ Docker Compose files created${NC}\n"
}

# ===== 5. COMMIT AND PUSH TO GIT =====

commit_and_push() {
  echo -e "${COLOR_BLUE}📤 Step 5: Commit and Push to Git${NC}\n"

  cd /Users/chinmaypanda/CustomApps

  git add docs/index.html DEPLOY-LANDING-PAGES.md docker-compose.prod.yml
  git commit -m "Add: Complete production deployment - React Vercel, GitHub Pages, Liferay Portal with all configs"
  git push origin main

  echo -e "${COLOR_GREEN}✅ Changes pushed to GitHub${NC}\n"
}

# ===== 6. SUMMARY =====

deployment_summary() {
  echo -e "${COLOR_GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
  echo -e "${COLOR_GREEN}║  OpenAutonomyX - Production Deployment Complete! 🎉         ║${NC}"
  echo -e "${COLOR_GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

  echo -e "${COLOR_BLUE}📍 DEPLOYMENT CHECKLIST:${NC}\n"

  echo "✅ React Frontend:"
  echo "   https://app.publishing.openautonomyx.com"
  echo "   (Auto-deploys on git push)"
  echo ""

  echo "✅ Marketing Site (GitHub Pages):"
  echo "   https://publishing.openautonomyx.com"
  echo "   (Live immediately)"
  echo ""

  echo "✅ Liferay Portal (Docker):"
  echo "   http://localhost:8080 (local)"
  echo "   https://publishing.openautonomyx.com:8080 (production)"
  echo ""

  echo "✅ SSL Certificates:"
  echo "   Let's Encrypt (Auto-renew)"
  echo ""

  echo -e "${COLOR_YELLOW}📋 Next Steps:${NC}\n"

  echo "1. Update DNS records:"
  echo "   publishing.openautonomyx.com  CNAME  openautonomyx.github.io"
  echo "   app.publishing.openautonomyx.com  CNAME  cname.vercel-dns.com"
  echo ""

  echo "2. Configure environment variables in Vercel:"
  echo "   VITE_API_GATEWAY_URL=https://api.publishing.openautonomyx.com"
  echo "   VITE_BLOG_URL=https://blog.publishing.openautonomyx.com"
  echo ""

  echo "3. Start Liferay (local):"
  echo "   docker-compose -f docker-compose.prod.yml up -d"
  echo ""

  echo "4. Monitor deployments:"
  echo "   Vercel: https://vercel.com/dashboard"
  echo "   GitHub Pages: https://github.com/Open-Autonomyx/Publishing-Platform/deployments"
  echo ""

  echo -e "${COLOR_GREEN}🚀 All systems ready to launch!${NC}\n"
}

# ===== INTERACTIVE MENU =====

show_menu() {
  echo ""
  echo -e "${COLOR_BLUE}Select deployment step:${NC}"
  echo "1. Deploy React Frontend (Vercel)"
  echo "2. Enable GitHub Pages"
  echo "3. Setup SSL Certificates"
  echo "4. Prepare Liferay Docker"
  echo "5. Commit & Push to Git"
  echo "6. Run All Steps"
  echo "7. Show Summary"
  echo "0. Exit"
  echo ""
}

# Main execution
if [ $# -eq 0 ]; then
  # Interactive mode
  while true; do
    show_menu
    read -p "Enter choice [0-7]: " choice

    case $choice in
      1) deploy_react_vercel ;;
      2) enable_github_pages ;;
      3) setup_ssl_certificates ;;
      4) prepare_liferay_docker ;;
      5) commit_and_push ;;
      6)
        deploy_react_vercel
        enable_github_pages
        setup_ssl_certificates
        prepare_liferay_docker
        commit_and_push
        deployment_summary
        ;;
      7) deployment_summary ;;
      0)
        echo -e "${COLOR_GREEN}Goodbye!${NC}"
        exit 0
        ;;
      *)
        echo -e "${COLOR_RED}Invalid choice${NC}"
        ;;
    esac
  done
else
  # Execute specific command
  case $1 in
    react) deploy_react_vercel ;;
    pages) enable_github_pages ;;
    ssl) setup_ssl_certificates ;;
    liferay) prepare_liferay_docker ;;
    commit) commit_and_push ;;
    all)
      deploy_react_vercel
      enable_github_pages
      setup_ssl_certificates
      prepare_liferay_docker
      commit_and_push
      deployment_summary
      ;;
    summary) deployment_summary ;;
    *)
      echo "Usage: $0 [react|pages|ssl|liferay|commit|all|summary]"
      exit 1
      ;;
  esac
fi
