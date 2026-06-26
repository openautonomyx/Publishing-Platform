# Format Converter Service - Multiple Content Formats

Transform content into different formats: EPUB, Slides, PDF, Audio, Video, HTML

**Service:** `@publishing-platform/formats`
**Port:** `3011`

---

## Supported Formats

| Format | MIME Type | Use Case | Extension |
|--------|-----------|----------|-----------|
| **EPUB** | application/epub+zip | E-books, Kindles | .epub |
| **Slides** | application/vnd.openxmlformats-officedocument.presentationml.presentation | Presentations | .pptx |
| **PDF** | application/pdf | Documents, Print | .pdf |
| **Audio** | audio/mpeg | Podcasts, Audiobooks | .mp3 |
| **Video** | video/mp4 | Video Content | .mp4 |
| **HTML** | text/html | Web pages, Interactive | .html |
| **Markdown** | text/markdown | Docs, GitHub | .md |

---

## API Endpoints

### Convert to EPUB (E-Book)

```bash
curl -X POST http://localhost:3011/api/v1/formats/epub \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Getting Started with OpenAutonomyX",
    "author": "Your Name",
    "content": "# Chapter 1\n\nYour content in markdown...",
    "chapters": [
      { "title": "Introduction", "content": "..." },
      { "title": "Setup", "content": "..." }
    ]
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "epub-uuid",
    "format": "epub",
    "title": "Getting Started with OpenAutonomyX",
    "filename": "getting-started-with-openautonomyx.epub",
    "downloadUrl": "/api/v1/formats/epub/{id}/download",
    "previewUrl": "/api/v1/formats/epub/{id}/preview",
    "chapters": 2,
    "estimatedSize": "2.5 MB"
  }
}
```

---

### Convert to Slides (PowerPoint)

```bash
curl -X POST http://localhost:3011/api/v1/formats/slides \
  -H "Content-Type: application/json" \
  -d '{
    "title": "OpenAutonomyX Overview",
    "author": "Your Name",
    "content": "# Title Slide\n\nContent here",
    "sections": [
      {
        "title": "What is OpenAutonomyX?",
        "content": "- Vendor-neutral\n- Local LLM support\n- Multi-format publishing"
      },
      {
        "title": "Key Features",
        "content": "- EPUB conversion\n- Slide generation\n- PDF export"
      }
    ]
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "slides-uuid",
    "format": "slides",
    "title": "OpenAutonomyX Overview",
    "filename": "openautonomyx-overview.pptx",
    "downloadUrl": "/api/v1/formats/slides/{id}/download",
    "slideCount": 3,
    "estimatedSize": "1.8 MB"
  }
}
```

---

### Convert to PDF

```bash
curl -X POST http://localhost:3011/api/v1/formats/pdf \
  -H "Content-Type: application/json" \
  -d '{
    "title": "OpenAutonomyX Guide",
    "author": "Your Name",
    "content": "Full markdown content..."
  }'
```

---

### Convert to Audio (Podcast/Audiobook)

```bash
curl -X POST http://localhost:3011/api/v1/formats/audio \
  -H "Content-Type: application/json" \
  -d '{
    "title": "OpenAutonomyX Podcast",
    "author": "Your Name",
    "narrator": "narrator-voice-id",
    "content": "Your script or article text..."
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "audio-uuid",
    "format": "audio",
    "title": "OpenAutonomyX Podcast",
    "filename": "openautonomyx-podcast.mp3",
    "streamUrl": "/api/v1/formats/audio/{id}/stream",
    "duration": 1200,
    "estimatedSize": "20 MB"
  }
}
```

---

### Convert to Video

```bash
curl -X POST http://localhost:3011/api/v1/formats/video \
  -H "Content-Type: application/json" \
  -d '{
    "title": "OpenAutonomyX Tutorial",
    "author": "Your Name",
    "thumbnail": "https://...",
    "content": "Video script or description..."
  }'
```

---

### Convert to Interactive HTML

```bash
curl -X POST http://localhost:3011/api/v1/formats/html \
  -H "Content-Type: application/json" \
  -d '{
    "title": "OpenAutonomyX Web Version",
    "author": "Your Name",
    "content": "# Content\n\nMarkdown here..."
  }'
```

---

### Get Supported Formats

```bash
curl http://localhost:3011/api/v1/formats/supported
```

**Response:**
```json
{
  "success": true,
  "data": {
    "formats": [
      {
        "name": "EPUB",
        "id": "epub",
        "description": "E-book format",
        "useCase": "Digital books, e-readers"
      },
      // ... more formats
    ]
  }
}
```

---

### Convert Between Formats

```bash
curl -X POST http://localhost:3011/api/v1/formats/convert \
  -H "Content-Type: application/json" \
  -d '{
    "contentId": "content-uuid",
    "fromFormat": "markdown",
    "toFormat": "epub"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "conversionId": "conversion-uuid",
    "from": "markdown",
    "to": "epub",
    "status": "queued",
    "estimatedTime": "2-5 minutes"
  }
}
```

---

## Workflow: One Content, Multiple Formats

### Step 1: Write Once
```
Write in OpenAutonomyX Platform
↓
Content stored (Markdown + metadata)
```

### Step 2: Generate All Formats

```bash
# Create EPUB
curl -X POST http://localhost:3011/api/v1/formats/epub \
  -d '{ "title": "My Book", "author": "Me", "content": "..." }'

# Create Slides
curl -X POST http://localhost:3011/api/v1/formats/slides \
  -d '{ "title": "My Presentation", "author": "Me", "content": "..." }'

# Create PDF
curl -X POST http://localhost:3011/api/v1/formats/pdf \
  -d '{ "title": "My Document", "author": "Me", "content": "..." }'

# Create Audio
curl -X POST http://localhost:3011/api/v1/formats/audio \
  -d '{ "title": "My Podcast", "author": "Me", "content": "..." }'

# Create Video
curl -X POST http://localhost:3011/api/v1/formats/video \
  -d '{ "title": "My Video", "author": "Me", "content": "..." }'
```

### Step 3: Distribute Everywhere

```
📱 Web (HTML)
📚 Ebooks (EPUB → Kindle, Kobo, Apple Books)
📊 Presentations (Slides → PowerPoint)
📄 Documents (PDF)
🎧 Podcast/Audiobook (Audio → Spotify, Apple Podcasts)
🎬 Video Platforms (Video → YouTube, Vimeo)
```

---

## Integration with Publishing Platform

### In Content Service

```typescript
// After publishing content
const content = await contentService.create({
  title: "My Post",
  content: "Markdown content..."
});

// Auto-generate all formats
await formatService.generateAll(content);

// Now available as:
// - EPUB (for e-readers)
// - Slides (for presentations)
// - PDF (for downloads)
// - Audio (for podcasts)
// - Video (for YouTube)
// - HTML (for web)
```

---

## Docker Setup

```yaml
formats:
  build:
    context: ./services/formats
    dockerfile: Dockerfile
  container_name: pp-formats
  ports:
    - "3011:3011"
  environment:
    NODE_ENV: development
    PORT: 3011
    SERVICE_NAME: formats
  depends_on:
    - event-bus
  networks:
    - pp-network
  restart: unless-stopped
```

---

## Use Cases

###📚 Author
- Write once
- Export as EPUB for Kindle
- Export as PDF for print
- Export as HTML for website

### 🎤 Podcaster
- Write script
- Convert to audio
- Publish to Spotify, Apple Podcasts
- Include transcript as blog post

### 🎓 Educator
- Create course content
- Generate slides for teaching
- Export as PDF handouts
- Create video tutorials

### 📰 Journalist
- Write article
- Publish as HTML on website
- Export as PDF for archive
- Audio narration for accessibility

### 🎬 Content Creator
- Write blog post
- Create video from script
- Generate slides for promotion
- Audio version for podcast
- All from single markdown file!

---

## Format Conversion Matrix

```
markdown ─→ EPUB
         ├→ PDF
         ├→ HTML
         ├→ Slides
         ├→ Audio
         └→ Video

html ────→ PDF
       ├→ Slides
       └→ Audio

epub ────→ PDF
       ├→ HTML
       └→ Audio

pdf ─────→ HTML
       └→ Slides
```

---

## Platform Integration

**20 Services + Format Converter = 21 Services**

```
┌─ Infrastructure (5)
├─ Core (2)
├─ Business Modules (8)
├─ Publishing (3)
│  ├─ Blog (3009)
│  ├─ Integrations (3010)
│  └─ Formats (3011) ← NEW
└─ Monitoring (2)
```

---

## Next Steps

- ☐ Create Dockerfile for formats service
- ☐ Add to docker-compose.yml
- ☐ Integrate with content service
- ☐ Create web UI for format selector
- ☐ Add format-specific viewers
- ☐ Implement webhook notifications
- ☐ Add batch conversion API

---

**Status:** ✅ Ready for Implementation

One content → Infinite formats! 🚀
