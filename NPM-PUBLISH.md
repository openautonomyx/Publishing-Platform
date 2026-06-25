# NPM Package Publishing Guide

## 📦 OpenAutonomyX SDK

**Package Name:** `openautonomyx`
**Scope:** Public (@openautonomyx/openautonomyx)
**Registry:** npm.js.org

---

## 🚀 Publishing Steps

### 1. Prepare for Release

```bash
# Navigate to project root
cd /Users/chinmaypanda/CustomApps

# Install dependencies
npm install

# Build the project
npm run build

# Run tests
npm run test

# Generate documentation
npm run docs
```

### 2. Update Version

```bash
# Patch release (1.0.0 → 1.0.1)
npm version patch

# Minor release (1.0.0 → 1.1.0)
npm version minor

# Major release (1.0.0 → 2.0.0)
npm version major
```

### 3. Publish to NPM

```bash
# Dry run (see what will be published)
npm publish --dry-run

# Publish to npm registry
npm publish

# Publish as scoped package
npm publish --access public
```

### 4. Verify Publication

```bash
# Check on npm registry
npm view openautonomyx

# Install from npm
npm install openautonomyx

# Verify it works
node -e "const oax = require('openautonomyx'); console.log(oax)"
```

---

## 📋 Package Contents

When published, the npm package includes:

```
openautonomyx/
├── dist/
│   ├── index.js
│   ├── index.d.ts
│   ├── openautonomyx.js
│   ├── openautonomyx.d.ts
│   ├── client.js
│   ├── client.d.ts
│   ├── types.js
│   ├── types.d.ts
│   ├── services/
│   │   ├── content.js
│   │   ├── publishing.js
│   │   ├── analytics.js
│   │   ├── team.js
│   │   ├── templates.js
│   │   ├── ai.js
│   │   └── multilingual.js
│   └── utils/
│       ├── errors.js
│       └── validators.js
├── package.json
├── README.md
└── LICENSE
```

---

## 📊 Version Strategy

### Current: 1.0.0

**Releases:**
- `1.0.0` - Initial release (MVP)
- `1.0.1` - Bug fixes
- `1.1.0` - New features (multilingual improvements)
- `2.0.0` - Breaking changes (new API design)

**Tagging:**
```bash
git tag -a v1.0.0 -m "Initial release"
git push origin v1.0.0
```

---

## 🔐 Authentication

### Setup NPM Token

```bash
# Login to npm
npm login

# Create token (on npmjs.com/settings/tokens)
# Token with read and publish permissions

# Store in ~/.npmrc or environment
npm config set //registry.npmjs.org/:_authToken YOUR_TOKEN_HERE
```

---

## 📤 GitHub Release

After npm publishing:

```bash
# Create GitHub release
gh release create v1.0.0 \
  --title "OpenAutonomyX v1.0.0" \
  --notes "Initial release with full SDK support"

# Upload dist/ folder as artifact
gh release upload v1.0.0 dist/
```

---

## 🎯 Distribution Checklist

- [ ] All tests passing
- [ ] Code coverage > 80%
- [ ] TypeScript types generated
- [ ] README updated
- [ ] CHANGELOG updated
- [ ] Version bumped in package.json
- [ ] Built (npm run build)
- [ ] Dry run successful (npm publish --dry-run)
- [ ] Published to npm
- [ ] Verified on npmjs.com
- [ ] GitHub release created
- [ ] Tagged in git

---

## 💻 Usage After Publishing

### Installation

```bash
npm install openautonomyx
# or
yarn add openautonomyx
# or
pnpm add openautonomyx
```

### Import

```typescript
import { OpenAutonomyX } from 'openautonomyx';

const oax = new OpenAutonomyX({
  apiKey: 'your-api-key'
});
```

---

## 🐛 Publishing Troubleshooting

### Authentication Error

```bash
# Clear npm cache
npm cache clean --force

# Re-login
npm logout
npm login
```

### Version Already Exists

```bash
# Get current version
npm view openautonomyx version

# Bump version in package.json and try again
npm version minor
npm publish
```

### TypeScript Definition Issues

```bash
# Regenerate types
npm run build

# Check declaration files
ls -la dist/*.d.ts
```

---

## 📚 NPM Package Documentation

See the full README at: [README.md](README.md)

---

## 🔗 Links

- **npm Registry:** https://www.npmjs.com/package/openautonomyx
- **GitHub:** https://github.com/openautonomyx/openautonomyx
- **Documentation:** https://openautonomyx.com/docs
- **API Reference:** https://openautonomyx.com/api

---

## 📞 Support

For publishing issues:
- Email: support@openautonomyx.com
- GitHub Issues: https://github.com/openautonomyx/openautonomyx/issues
- Discord: https://discord.gg/openautonomyx

---

**Ready to publish!** 🚀
