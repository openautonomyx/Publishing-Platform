import express, { Request, Response } from 'express';
import cors from 'cors';
import axios from 'axios';
import { v4 as uuid } from 'uuid';

const app = express();
const PORT = process.env.PORT || 3011;
const API_GATEWAY = process.env.API_GATEWAY_URL || 'http://localhost:3000';

app.use(cors());
app.use(express.json());

// Health check
app.get('/health', (req: Request, res: Response) => {
  res.json({ status: 'ok', service: 'formats', timestamp: new Date() });
});

// ==================== FORMAT CONVERTERS ====================

// EPUB Format
app.post('/api/v1/formats/epub', async (req: Request, res: Response) => {
  try {
    const { contentId, title, author, content, chapters } = req.body;

    // Build EPUB structure
    const epubContent = {
      id: uuid(),
      format: 'epub',
      title,
      author,
      cover: generateEPUBCover(title, author),
      toc: generateTableOfContents(chapters),
      chapters: chapters || [{ title: 'Content', html: markdownToHTML(content) }],
      metadata: {
        createdAt: new Date(),
        contentId,
        generator: 'OpenAutonomyX EPUB Generator v1.0'
      }
    };

    res.json({
      success: true,
      data: {
        id: epubContent.id,
        format: 'epub',
        title,
        filename: `${title.replace(/\s+/g, '-').toLowerCase()}.epub`,
        downloadUrl: `/api/v1/formats/epub/${epubContent.id}/download`,
        previewUrl: `/api/v1/formats/epub/${epubContent.id}/preview`,
        chapters: epubContent.chapters.length,
        estimatedSize: '2.5 MB'
      }
    });
  } catch (error) {
    res.status(500).json({ error: (error as Error).message });
  }
});

// Slides Format (PowerPoint/Presentation)
app.post('/api/v1/formats/slides', async (req: Request, res: Response) => {
  try {
    const { contentId, title, author, content, sections } = req.body;

    const slidesContent = {
      id: uuid(),
      format: 'slides',
      title,
      author,
      slides: generateSlides(content, sections),
      theme: 'default',
      metadata: {
        createdAt: new Date(),
        contentId,
        generator: 'OpenAutonomyX Slides Generator v1.0'
      }
    };

    res.json({
      success: true,
      data: {
        id: slidesContent.id,
        format: 'slides',
        title,
        filename: `${title.replace(/\s+/g, '-').toLowerCase()}.pptx`,
        downloadUrl: `/api/v1/formats/slides/${slidesContent.id}/download`,
        previewUrl: `/api/v1/formats/slides/${slidesContent.id}/preview`,
        slideCount: slidesContent.slides.length,
        estimatedSize: '1.8 MB'
      }
    });
  } catch (error) {
    res.status(500).json({ error: (error as Error).message });
  }
});

// PDF Format
app.post('/api/v1/formats/pdf', async (req: Request, res: Response) => {
  try {
    const { contentId, title, author, content } = req.body;

    const pdfContent = {
      id: uuid(),
      format: 'pdf',
      title,
      author,
      pages: Math.ceil(content.length / 3000), // Rough estimate
      content: content,
      metadata: {
        createdAt: new Date(),
        contentId,
        generator: 'OpenAutonomyX PDF Generator v1.0'
      }
    };

    res.json({
      success: true,
      data: {
        id: pdfContent.id,
        format: 'pdf',
        title,
        filename: `${title.replace(/\s+/g, '-').toLowerCase()}.pdf`,
        downloadUrl: `/api/v1/formats/pdf/${pdfContent.id}/download`,
        previewUrl: `/api/v1/formats/pdf/${pdfContent.id}/preview`,
        pages: pdfContent.pages,
        estimatedSize: `${(pdfContent.pages * 0.1).toFixed(1)} MB`
      }
    });
  } catch (error) {
    res.status(500).json({ error: (error as Error).message });
  }
});

// Audio Format (Podcast/Audiobook)
app.post('/api/v1/formats/audio', async (req: Request, res: Response) => {
  try {
    const { contentId, title, author, content, narrator } = req.body;

    const audioContent = {
      id: uuid(),
      format: 'audio',
      title,
      author,
      narrator: narrator || 'default',
      duration: estimateAudioDuration(content),
      content: content,
      metadata: {
        createdAt: new Date(),
        contentId,
        generator: 'OpenAutonomyX Audio Generator v1.0'
      }
    };

    res.json({
      success: true,
      data: {
        id: audioContent.id,
        format: 'audio',
        title,
        filename: `${title.replace(/\s+/g, '-').toLowerCase()}.mp3`,
        downloadUrl: `/api/v1/formats/audio/${audioContent.id}/download`,
        streamUrl: `/api/v1/formats/audio/${audioContent.id}/stream`,
        duration: audioContent.duration,
        estimatedSize: `${(audioContent.duration / 60).toFixed(1)} MB`
      }
    });
  } catch (error) {
    res.status(500).json({ error: (error as Error).message });
  }
});

// Video Format
app.post('/api/v1/formats/video', async (req: Request, res: Response) => {
  try {
    const { contentId, title, author, content, thumbnail } = req.body;

    const videoContent = {
      id: uuid(),
      format: 'video',
      title,
      author,
      thumbnail,
      duration: 300, // Default 5 min
      content: content,
      metadata: {
        createdAt: new Date(),
        contentId,
        generator: 'OpenAutonomyX Video Generator v1.0'
      }
    };

    res.json({
      success: true,
      data: {
        id: videoContent.id,
        format: 'video',
        title,
        filename: `${title.replace(/\s+/g, '-').toLowerCase()}.mp4`,
        downloadUrl: `/api/v1/formats/video/${videoContent.id}/download`,
        streamUrl: `/api/v1/formats/video/${videoContent.id}/stream`,
        duration: videoContent.duration,
        estimatedSize: '50 MB',
        thumbnail
      }
    });
  } catch (error) {
    res.status(500).json({ error: (error as Error).message });
  }
});

// Interactive HTML Format
app.post('/api/v1/formats/html', async (req: Request, res: Response) => {
  try {
    const { contentId, title, author, content } = req.body;

    const htmlContent = {
      id: uuid(),
      format: 'html',
      title,
      author,
      html: generateInteractiveHTML(title, author, content),
      metadata: {
        createdAt: new Date(),
        contentId,
        generator: 'OpenAutonomyX HTML Generator v1.0'
      }
    };

    res.json({
      success: true,
      data: {
        id: htmlContent.id,
        format: 'html',
        title,
        filename: `${title.replace(/\s+/g, '-').toLowerCase()}.html`,
        viewUrl: `/api/v1/formats/html/${htmlContent.id}/view`,
        downloadUrl: `/api/v1/formats/html/${htmlContent.id}/download`,
        responsive: true,
        estimatedSize: '0.5 MB'
      }
    });
  } catch (error) {
    res.status(500).json({ error: (error as Error).message });
  }
});

// Convert from one format to another
app.post('/api/v1/formats/convert', async (req: Request, res: Response) => {
  try {
    const { contentId, fromFormat, toFormat } = req.body;

    const conversions = [
      { from: 'markdown', to: ['epub', 'pdf', 'html', 'slides', 'video', 'audio'] },
      { from: 'html', to: ['pdf', 'epub', 'slides', 'audio'] },
      { from: 'epub', to: ['pdf', 'html', 'audio'] },
      { from: 'pdf', to: ['html', 'slides'] }
    ];

    const supported = conversions.find(c => c.from === fromFormat);
    if (!supported || !supported.to.includes(toFormat)) {
      return res.status(400).json({
        error: `Cannot convert from ${fromFormat} to ${toFormat}`,
        supported: supported?.to || []
      });
    }

    res.json({
      success: true,
      data: {
        conversionId: uuid(),
        from: fromFormat,
        to: toFormat,
        status: 'queued',
        estimatedTime: '2-5 minutes',
        webhook: `/api/v1/formats/webhooks/conversion-complete`
      }
    });
  } catch (error) {
    res.status(500).json({ error: (error as Error).message });
  }
});

// Get all supported formats
app.get('/api/v1/formats/supported', (req: Request, res: Response) => {
  res.json({
    success: true,
    data: {
      formats: [
        {
          name: 'EPUB',
          id: 'epub',
          description: 'E-book format (Kindle, Kobo, etc)',
          mimeType: 'application/epub+zip',
          extensions: ['.epub'],
          useCase: 'Digital books, e-readers'
        },
        {
          name: 'Slides',
          id: 'slides',
          description: 'PowerPoint presentation format',
          mimeType: 'application/vnd.openxmlformats-officedocument.presentationml.presentation',
          extensions: ['.pptx'],
          useCase: 'Presentations, teaching'
        },
        {
          name: 'PDF',
          id: 'pdf',
          description: 'Portable document format',
          mimeType: 'application/pdf',
          extensions: ['.pdf'],
          useCase: 'Professional documents, printing'
        },
        {
          name: 'Audio',
          id: 'audio',
          description: 'Audio/Podcast format',
          mimeType: 'audio/mpeg',
          extensions: ['.mp3'],
          useCase: 'Podcasts, audiobooks, narration'
        },
        {
          name: 'Video',
          id: 'video',
          description: 'Video format',
          mimeType: 'video/mp4',
          extensions: ['.mp4'],
          useCase: 'Video content, courses, tutorials'
        },
        {
          name: 'HTML',
          id: 'html',
          description: 'Interactive web format',
          mimeType: 'text/html',
          extensions: ['.html'],
          useCase: 'Web pages, interactive content'
        },
        {
          name: 'Markdown',
          id: 'markdown',
          description: 'Plain text markdown',
          mimeType: 'text/markdown',
          extensions: ['.md'],
          useCase: 'Documentation, GitHub'
        }
      ]
    }
  });
});

// ==================== PREVIEW ENDPOINTS ====================

app.get('/api/v1/formats/:format/:id/preview', (req: Request, res: Response) => {
  res.json({
    success: true,
    data: {
      format: req.params.format,
      id: req.params.id,
      previewUrl: `/preview/${req.params.format}/${req.params.id}`,
      previewAvailable: true
    }
  });
});

// ==================== HELPER FUNCTIONS ====================

function generateEPUBCover(title: string, author: string): string {
  return `
    <svg xmlns="http://www.w3.org/2000/svg" width="600" height="800">
      <rect width="600" height="800" fill="#667eea"/>
      <text x="300" y="400" font-size="48" font-weight="bold" fill="white" text-anchor="middle">
        ${title}
      </text>
      <text x="300" y="500" font-size="24" fill="rgba(255,255,255,0.8)" text-anchor="middle">
        by ${author}
      </text>
    </svg>
  `;
}

function generateTableOfContents(chapters: any[]): any[] {
  return (chapters || []).map((ch, idx) => ({
    id: `chapter-${idx}`,
    title: ch.title,
    level: ch.level || 1
  }));
}

function markdownToHTML(markdown: string): string {
  // Simple markdown to HTML conversion
  return markdown
    .replace(/^# (.*?)$/gm, '<h1>$1</h1>')
    .replace(/^## (.*?)$/gm, '<h2>$1</h2>')
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\n\n/g, '</p><p>')
    .replace(/^/gm, '<p>')
    .replace(/$/gm, '</p>');
}

function generateSlides(content: string, sections: any[]): any[] {
  const slides = [];

  // Title slide
  slides.push({
    type: 'title',
    title: content.split('\n')[0],
    subtitle: 'Created with OpenAutonomyX'
  });

  // Content slides
  (sections || []).forEach((section, idx) => {
    slides.push({
      type: 'content',
      title: section.title,
      bullets: section.content.split('\n').filter((l: string) => l.trim())
    });
  });

  return slides;
}

function estimateAudioDuration(content: string): number {
  // Rough estimate: 150 words per minute = 2.5 words per second
  const wordCount = content.split(/\s+/).length;
  return Math.ceil((wordCount / 150) * 60); // Duration in seconds
}

function generateInteractiveHTML(title: string, author: string, content: string): string {
  return `
    <!DOCTYPE html>
    <html lang="en">
    <head>
      <meta charset="UTF-8">
      <meta name="viewport" content="width=device-width, initial-scale=1.0">
      <title>${title}</title>
      <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 800px; margin: 0 auto; padding: 40px; line-height: 1.6; color: #333; }
        h1 { color: #667eea; }
        code { background: #f0f0f0; padding: 2px 6px; border-radius: 3px; }
      </style>
    </head>
    <body>
      <h1>${title}</h1>
      <p><em>by ${author}</em></p>
      <hr>
      ${markdownToHTML(content)}
    </body>
    </html>
  `;
}

app.listen(PORT, () => {
  console.log(`🎬 Format Converter Service running on port ${PORT}`);
  console.log(`📦 Formats: EPUB, Slides, PDF, Audio, Video, HTML`);
  console.log(`🔗 API: http://localhost:${PORT}/api/v1/formats`);
});
