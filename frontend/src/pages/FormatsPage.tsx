import { useEffect } from 'react';
import './Pages.css';

interface PageProps {
  onContextChange: (context: string, data: any) => void;
}

function FormatsPage({ onContextChange }: PageProps) {
  useEffect(() => {
    onContextChange('formats', {});
  }, [onContextChange]);

  return (
    <div className="page-container">
      <div className="page-header">
        <h1>🎬 Format Converter</h1>
        <p>Convert your content to multiple formats</p>
      </div>

      <div className="content-section">
        <div className="blog-post">
          <h3>📚 EPUB</h3>
          <p>Generate e-books compatible with Kindle, Kobo, Apple Books, and more</p>
          <button className="btn-secondary">Convert to EPUB</button>
        </div>

        <div className="blog-post">
          <h3>📊 Slides</h3>
          <p>Create PowerPoint presentations from your content</p>
          <button className="btn-secondary">Generate Slides</button>
        </div>

        <div className="blog-post">
          <h3>📄 PDF</h3>
          <p>Export professional PDF documents for printing and sharing</p>
          <button className="btn-secondary">Export to PDF</button>
        </div>

        <div className="blog-post">
          <h3>🎧 Audio</h3>
          <p>Create podcasts and audiobooks from your text</p>
          <button className="btn-secondary">Generate Audio</button>
        </div>

        <div className="blog-post">
          <h3>🎬 Video</h3>
          <p>Transform content into video format for YouTube and other platforms</p>
          <button className="btn-secondary">Create Video</button>
        </div>

        <div className="blog-post">
          <h3>🌐 HTML</h3>
          <p>Generate interactive web pages for your content</p>
          <button className="btn-secondary">Generate HTML</button>
        </div>
      </div>
    </div>
  );
}

export default FormatsPage;
