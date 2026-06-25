# OpenAutonomyX SDK

**The global creative platform for creating, publishing, and analyzing content in 50+ languages across 6 platforms.**

[![npm version](https://badge.fury.io/js/openautonomyx.svg)](https://badge.fury.io/js/openautonomyx)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## 🌐 What is OpenAutonomyX?

OpenAutonomyX is a comprehensive SDK for creators, businesses, and agencies to:

- ✨ **Create content** once in any language
- 🌍 **Generate native versions** in 50+ languages automatically (not translations)
- 📱 **Publish simultaneously** to 6 social platforms
- 📊 **Track analytics** in real-time
- 🤖 **Leverage AI** for content generation and optimization
- 👥 **Collaborate** with teams

---

## 🚀 Quick Start

### Installation

```bash
npm install openautonomyx
```

### Basic Usage

```typescript
import { OpenAutonomyX } from 'openautonomyx';

// Initialize the SDK
const oax = new OpenAutonomyX({
  apiKey: process.env.OPENAUTONOMYX_API_KEY,
  apiUrl: 'https://openautonomyx.com:3001/api/v1'
});

// Create content
const publication = await oax.content.create({
  title: 'My Article',
  content: 'This is my article content...',
  languages: ['en', 'es', 'fr', 'ja']
});

// Publish to multiple platforms
const results = await oax.publishing.publish(publication.id, {
  platforms: ['instagram', 'twitter', 'linkedin'],
  scheduledFor: '2026-06-30T10:00:00Z'
});

// Get analytics
const report = await oax.analytics.getReport({
  start: '2026-01-01',
  end: '2026-06-30'
});

console.log(`📊 Total impressions: ${report.totalImpressions}`);
```

---

## 📚 Complete API Reference

### Content Management

```typescript
// Create a publication
await oax.content.create({
  title: 'Your Title',
  content: 'Your content...',
  contentType: 'blog',
  languages: ['en', 'es', 'fr'],
  tags: ['ai', 'automation']
});

// Get publication
const pub = await oax.content.get(publicationId);

// List publications
const publications = await oax.content.list({
  limit: 20,
  offset: 0,
  sortBy: 'createdAt'
});

// Update publication
await oax.content.update(publicationId, {
  title: 'Updated Title'
});

// Archive publication
await oax.content.archive(publicationId);

// Duplicate publication
await oax.content.duplicate(publicationId);
```

### Publishing

```typescript
// Publish immediately
const results = await oax.publishing.publish(publicationId, {
  platforms: ['instagram', 'twitter', 'linkedin']
});

// Schedule for later
await oax.publishing.schedule(publicationId, {
  platforms: ['tiktok', 'youtube'],
  scheduledFor: '2026-07-01T12:00:00Z'
});

// Get status
const status = await oax.publishing.getStatus(publicationId);

// Cancel scheduled post
await oax.publishing.cancel(publicationId, 'instagram');
```

### Analytics

```typescript
// Get metrics for a publication
const metrics = await oax.analytics.getMetrics(publicationId);

// Get performance report
const report = await oax.analytics.getReport({
  start: '2026-01-01',
  end: '2026-06-30'
});

console.log(`Impressions: ${report.totalImpressions}`);
console.log(`Engagement: ${report.totalEngagement}`);
console.log(`Top performers:`, report.topPerformers);
```

### Team Management

```typescript
// Get team info
const team = await oax.team.getTeam();

// Invite team member
await oax.team.invite({
  email: 'teammate@example.com',
  role: 'editor'
});

// Update member role
await oax.team.updateRole(memberId, 'admin');

// Remove member
await oax.team.removeMember(memberId);
```

### Templates

```typescript
// Search templates
const templates = await oax.templates.search({
  category: 'marketing',
  tags: ['email'],
  limit: 20
});

// Get template details
const template = await oax.templates.get(templateId);

// Clone a template
const cloned = await oax.templates.clone(templateId);
```

### AI Content Generation

```typescript
// Generate content with AI
const generated = await oax.ai.generate({
  topic: 'The future of AI',
  contentType: 'blog',
  tone: 'professional',
  language: 'en',
  keywords: ['ai', 'future', 'innovation']
});

console.log(generated.content);
console.log(`Reading time: ${generated.readingTime} minutes`);

// Improve existing content
const improved = await oax.ai.improve(yourContent);

// Summarize content
const summary = await oax.ai.summarize(longContent);

// Extract keywords
const keywords = await oax.ai.extractKeywords(content);
```

### Multilingual Publishing

```typescript
// Generate native-fluent content in multiple languages
const multilingualContent = await oax.multilingual.generateMultilingual(
  'Your topic',
  ['en', 'es', 'fr', 'ja', 'ar', 'zh']
);

// Each is native composition, not translation
for (const content of multilingualContent) {
  console.log(`${content.language}: ${content.title}`);
  console.log(`Reading time: ${content.readingTime} min`);
}

// Translate content
const translations = await oax.multilingual.translate({
  content: 'Your content to translate',
  targetLanguages: ['es', 'fr', 'de', 'it']
});

// Detect language
const detected = await oax.multilingual.detectLanguage(content);
console.log(`Detected language: ${detected.language}`);

// Get supported languages
const languages = await oax.multilingual.getAvailableLanguages();
console.log(`Supported: ${languages.languages.join(', ')}`);
```

---

## 🔑 Types & Interfaces

### Publication

```typescript
interface Publication {
  id: string;
  title: string;
  content: string;
  contentType: 'blog' | 'social' | 'email' | 'newsletter' | 'product' | 'announcement';
  languages: SupportedLanguage[];
  status: 'draft' | 'published' | 'archived';
  createdAt: string;
  updatedAt: string;
  createdBy: string;
}
```

### Supported Languages (50+)

```
English, Spanish, French, German, Italian, Portuguese, Russian, Japanese, Chinese, Korean,
Arabic, Hindi, Thai, Vietnamese, Turkish, Polish, Dutch, Swedish, Danish, Norwegian,
Finnish, Indonesian, Malay, Tagalog, Bulgarian, Czech, Croatian, Romanian, Slovak, Slovenian,
Hungarian, Greek, Ukrainian, Hebrew, Persian, Urdu, Bengali, Tamil, Telugu, Kannada, Malayalam,
Afrikaans, Albanian, Estonian, Lithuanian, Latvian, Macedonian
```

### Social Platforms

```typescript
type SocialPlatform =
  | 'instagram'
  | 'twitter'
  | 'linkedin'
  | 'facebook'
  | 'tiktok'
  | 'youtube';
```

---

## 🔒 Authentication

```typescript
const oax = new OpenAutonomyX({
  apiKey: process.env.OPENAUTONOMYX_API_KEY,
  // Optional: custom API URL (default: https://openautonomyx.com:3001/api/v1)
  apiUrl: 'https://custom-api.example.com/api/v1',
  // Optional: custom timeout in milliseconds (default: 30000)
  timeout: 60000
});
```

Get your API key at: **https://openautonomyx.com/dashboard/api**

---

## 🚀 Examples

### Example 1: Blog Post to 50 Languages

```typescript
// Create content in English
const publication = await oax.content.create({
  title: 'The Future of AI',
  content: 'AI is transforming...',
  contentType: 'blog'
});

// Generate native versions in 50 languages
const multilingual = await oax.multilingual.generateMultilingual(
  'The Future of AI',
  ['en', 'es', 'fr', 'de', 'it', 'pt', 'ru', 'ja', 'zh', 'ko', 
   'ar', 'hi', 'th', 'vi', 'tr', 'pl', 'nl', 'sv', 'da', 'no',
   'fi', 'id', 'ms', 'tl', 'bg', 'cs', 'hr', 'ro', 'sk', 'sl',
   'hu', 'el', 'uk', 'he', 'fa', 'ur', 'bn', 'ta', 'te', 'kn',
   'ml', 'af', 'sq', 'et', 'lt', 'lv', 'mk'] // ... more
);

// Publish to all platforms
for (const lang of multilingual) {
  await oax.publishing.publish(publication.id, {
    platforms: ['instagram', 'twitter', 'linkedin'],
    content: lang.content
  });
}
```

### Example 2: Social Media Scheduling

```typescript
// Create social content
const social = await oax.content.create({
  title: 'Weekly Tips',
  content: 'Check out our latest...',
  contentType: 'social',
  languages: ['en', 'es', 'fr']
});

// Schedule for tomorrow at noon (UTC)
const tomorrow = new Date();
tomorrow.setDate(tomorrow.getDate() + 1);
tomorrow.setHours(12, 0, 0, 0);

await oax.publishing.schedule(social.id, {
  platforms: ['instagram', 'twitter', 'tiktok'],
  scheduledFor: tomorrow.toISOString()
});
```

### Example 3: Content with AI Enhancement

```typescript
// Generate AI content
const generated = await oax.ai.generate({
  topic: 'Smart Home Technology',
  contentType: 'blog',
  tone: 'technical',
  targetAudience: 'tech enthusiasts',
  keywords: ['smart home', 'IoT', 'automation']
});

// Create publication from AI output
const publication = await oax.content.create({
  title: generated.title,
  content: generated.content,
  tags: generated.keywords
});

// Publish to all platforms
await oax.publishing.publish(publication.id, {
  platforms: ['instagram', 'twitter', 'linkedin', 'facebook']
});
```

---

## 🛠️ Error Handling

```typescript
import {
  OpenAutonomyXError,
  ValidationError,
  AuthenticationError,
  NotFoundError
} from 'openautonomyx';

try {
  await oax.content.create({
    title: '',  // This will fail validation
    content: 'content'
  });
} catch (error) {
  if (error instanceof ValidationError) {
    console.error('Validation failed:', error.message);
  } else if (error instanceof AuthenticationError) {
    console.error('Auth failed:', error.message);
  } else if (error instanceof OpenAutonomyXError) {
    console.error(`Error [${error.code}]: ${error.message}`);
  }
}
```

---

## 📖 Documentation

- **Full Docs:** https://openautonomyx.com/docs
- **API Reference:** https://openautonomyx.com/api
- **Examples:** https://github.com/openautonomyx/examples
- **Blog:** https://openautonomyx.com/blog

---

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md)

---

## 📄 License

MIT - See [LICENSE](LICENSE)

---

## 💬 Support

- **Email:** support@openautonomyx.com
- **Discord:** https://discord.gg/openautonomyx
- **Twitter:** @openautonomyx
- **GitHub Issues:** https://github.com/openautonomyx/openautonomyx/issues

---

## 🌟 Features at a Glance

| Feature | Status | Details |
|---------|--------|---------|
| Content Creation | ✅ | Create, edit, manage publications |
| Multi-Language | ✅ | 50+ languages with native fluency |
| Social Publishing | ✅ | 6 platforms, simultaneous posting |
| AI Generation | ✅ | Claude-powered content creation |
| Analytics | ✅ | Real-time performance tracking |
| Team Collaboration | ✅ | Invite teammates, set roles |
| Templates | ✅ | 1,000+ professional templates |
| Scheduling | ✅ | Post now or schedule for later |
| API Access | ✅ | Full programmatic access |
| Webhooks | ✅ | Real-time event notifications |

---

**Made with ❤️ by OpenAutonomyX**

*Create once. Publish everywhere. In any language.*
