# ⚡ Automation Setup Checklist

**Time Required:** 2 minutes  
**Result:** Fully automated deployment from code push to production

---

## Quick Setup (Copy & Paste)

### Step 1: Get Your VPS Password

You'll need the password for `almalinux` user on agennext.com.

### Step 2: Add GitHub Secret

```bash
# Open this URL in browser:
https://github.com/fractional-pm/creative-platform/settings/secrets/actions

# Click "New repository secret"
# Fill in:
#   Name:  VPS_PASSWORD
#   Value: [paste your VPS password]
# Click "Add secret"
```

### Step 3: Enable Workflows

```bash
# Open:
https://github.com/fractional-pm/creative-platform/actions

# Verify both workflows show (enabled):
# ✅ Build and Push Docker Image to GHCR
# ✅ Auto Deploy to VPS
```

### Step 4: Test It

```bash
cd ~/CustomApps/creative-platform

# Make a test change
echo "# Test deployment" >> README.md

# Commit and push
git add README.md
git commit -m "test: trigger automation"
git push origin main
```

### Step 5: Watch Deployment

Open: https://github.com/fractional-pm/creative-platform/actions

You'll see:
1. **Build and Push Docker Image to GHCR** - Running (2-3 min)
2. **Auto Deploy to VPS** - Waiting or Running (2-3 min)
3. ✅ Complete - API is live

---

## After Setup

### Deploy New Code

```bash
git add .
git commit -m "your message"
git push origin main
```

**That's it.** Automation handles the rest.

---

## Verify Deployment

After workflow completes:

```bash
curl http://agennext.com:3001/health
```

Should return:
```json
{"status":"OK","timestamp":"2026-06-25T12:34:56Z"}
```

---

## Troubleshooting

**If deployment fails:**

1. Check GitHub Actions logs: https://github.com/fractional-pm/creative-platform/actions
2. Look for error in "Deploy to VPS via SSH" step
3. Common causes:
   - VPS_PASSWORD secret not set → Add it
   - Wrong password → Update secret
   - VPS offline → Check VPS status

**If API doesn't respond:**

```bash
ssh almalinux@agennext.com
cd ~/creative-platform
docker-compose logs api --tail=50
```

---

## Done! 🎉

You now have:
- ✅ Automatic Docker builds on every push
- ✅ Automatic deployment to VPS
- ✅ Automatic health checks
- ✅ Zero manual steps

**Next:** `git push origin main` and watch your code go live.

