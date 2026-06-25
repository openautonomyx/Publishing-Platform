#!/bin/bash

# Complete GitHub OAuth Setup & Deploy Script
# Run this on your local machine to push both repositories

set -e

echo ""
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║     GitHub OAuth Setup & Repository Push                       ║"
echo "║     OpenAGX Creative Platform Deployment                       ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# ============================================================================
# STEP 1: CHECK & INSTALL GITHUB CLI
# ============================================================================

echo "${BLUE}STEP 1: Setting up GitHub CLI${NC}"
echo "════════════════════════════════════════════════════════════════"
echo ""

if ! command -v gh &> /dev/null; then
    echo "GitHub CLI not found. Installing..."
    echo ""

    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo "macOS detected. Installing via Homebrew..."
        brew install gh
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "Linux detected. Installing via apt..."
        apt install gh
    else
        echo "Please install GitHub CLI from: https://cli.github.com"
        exit 1
    fi
else
    echo "✅ GitHub CLI already installed"
    gh --version
fi

echo ""
echo "✅ GitHub CLI ready"
echo ""

# ============================================================================
# STEP 2: AUTHENTICATE WITH OAUTH
# ============================================================================

echo "${BLUE}STEP 2: OAuth Authentication${NC}"
echo "════════════════════════════════════════════════════════════════"
echo ""

if gh auth status > /dev/null 2>&1; then
    echo "✅ Already authenticated:"
    gh auth status --show-token=false
    echo ""
else
    echo "⏳ Opening GitHub login in browser..."
    echo ""
    echo "Please:"
    echo "  1. Choose SSH (recommended) when prompted"
    echo "  2. Generate new SSH key when asked"
    echo "  3. Authorize in the browser window that opens"
    echo ""

    gh auth login

    if ! gh auth status > /dev/null 2>&1; then
        echo "❌ Authentication failed"
        exit 1
    fi

    echo ""
    echo "✅ OAuth authenticated successfully:"
    gh auth status --show-token=false
    echo ""
fi

# ============================================================================
# STEP 3: PUSH REPOSITORY 1 (creative-platform)
# ============================================================================

echo ""
echo "${BLUE}STEP 3: Pushing creative-platform${NC}"
echo "════════════════════════════════════════════════════════════════"
echo ""

cd /Users/chinmaypanda/CustomApps/creative-platform

echo "Repository: creative-platform"
echo "Remote: $(git remote get-url origin)"
echo "Branch: $(git rev-parse --abbrev-ref HEAD)"
echo "Commits to push: $(git rev-list --count origin/main..HEAD 2>/dev/null || echo 'unknown')"
echo ""

echo "Pushing to GitHub..."
if git push -u origin main; then
    echo ""
    echo "✅ creative-platform pushed successfully"
    echo "   https://github.com/openagx/creative-platform"
else
    echo ""
    echo "❌ Push failed. Retrying..."
    sleep 2
    git push -u origin main
fi

echo ""

# ============================================================================
# STEP 4: PUSH REPOSITORY 2 (agent-registry)
# ============================================================================

echo ""
echo "${BLUE}STEP 4: Pushing agent-registry${NC}"
echo "════════════════════════════════════════════════════════════════"
echo ""

cd /Users/chinmaypanda/CustomApps/agent-registry

echo "Repository: agent-registry"
echo "Remote: $(git remote get-url origin)"
echo "Branch: $(git rev-parse --abbrev-ref HEAD)"
echo "Commits to push: $(git rev-list --count origin/main..HEAD 2>/dev/null || echo 'unknown')"
echo ""

echo "Pushing to GitHub..."
if git push -u origin main; then
    echo ""
    echo "✅ agent-registry pushed successfully"
    echo "   https://github.com/openagx/agent-registry"
else
    echo ""
    echo "❌ Push failed. Retrying..."
    sleep 2
    git push -u origin main
fi

# ============================================================================
# STEP 5: VERIFY & SUMMARY
# ============================================================================

echo ""
echo ""
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║                    ✅ DEPLOYMENT COMPLETE                     ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

echo "${GREEN}Both repositories are now on GitHub:${NC}"
echo ""
echo "  1️⃣  creative-platform"
echo "      https://github.com/openagx/creative-platform"
echo ""
echo "  2️⃣  agent-registry"
echo "      https://github.com/openagx/agent-registry"
echo ""

echo "${YELLOW}Next Steps:${NC}"
echo ""
echo "  1. Configure GitHub Secrets"
echo "     URL: https://github.com/openagx/creative-platform/settings/secrets/actions"
echo ""
echo "     Add secrets:"
echo "       - VPS_HOST=agennext.com"
echo "       - VPS_USER=almalinux"
echo "       - VPS_PASSWORD=[your password]"
echo "       - VPS_SSH_KEY=[your SSH key]"
echo "       - DB_PASSWORD=[database password]"
echo "       - JWT_SECRET=[JWT secret]"
echo ""
echo "  2. Follow the deployment guide:"
echo "     /Users/chinmaypanda/CustomApps/MASTER-DEPLOYMENT-GUIDE.md"
echo ""
echo "  3. Deploy to VPS"
echo "     cd /Users/chinmaypanda/CustomApps/creative-platform"
echo "     VPS_PASSWORD='pw' bash deploy/vps-automation-agent.sh"
echo ""

echo "${GREEN}🎉 GitHub push complete!${NC}"
echo ""
