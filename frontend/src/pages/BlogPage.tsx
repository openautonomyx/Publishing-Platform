import { useEffect } from 'react';
import './Pages.css';

interface PageProps {
  onContextChange: (context: string, data: any) => void;
}

function BlogPage({ onContextChange }: PageProps) {
  useEffect(() => {
    onContextChange('blog', {});
  }, [onContextChange]);

  return (
    <div className="page-container">
      <div className="page-header">
        <h1>📝 Blog</h1>
        <p>Your WordPress blog connected to OpenAutonomyX</p>
        <button className="btn-action">+ New Post</button>
      </div>

      <div className="content-section">
        <div className="blog-post">
          <h3>Welcome to OpenAutonomyX Blog</h3>
          <div className="post-meta">
            <span>📅 June 26, 2026</span>
            <span>👤 OpenAutonomyX</span>
            <span>🏷️ Platform</span>
          </div>
          <p>Your blog is now connected and ready to publish. Write once and publish to multiple platforms.</p>
          <button className="btn-secondary">Read More</button>
        </div>

        <div className="blog-post">
          <h3>Multi-Format Publishing</h3>
          <div className="post-meta">
            <span>📅 Coming Soon</span>
            <span>👤 OpenAutonomyX</span>
            <span>🏷️ Features</span>
          </div>
          <p>Transform your blog posts into EPUB, PDF, Slides, and more with one click.</p>
          <button className="btn-secondary">Read More</button>
        </div>
      </div>
    </div>
  );
}

export default BlogPage;
