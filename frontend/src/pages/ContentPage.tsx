import { useEffect } from 'react';
import './Pages.css';

interface PageProps {
  onContextChange: (context: string, data: any) => void;
}

function ContentPage({ onContextChange }: PageProps) {
  useEffect(() => {
    onContextChange('content', {});
  }, [onContextChange]);

  return (
    <div className="page-container">
      <div className="page-header">
        <h1>📄 Content Management</h1>
        <p>Create and manage your content</p>
        <button className="btn-action">+ Create Content</button>
      </div>

      <div className="content-section">
        <div className="blog-post">
          <h3>Getting Started Guide</h3>
          <div className="post-meta">
            <span>📅 June 26, 2026</span>
            <span>👤 You</span>
            <span>🏷️ Guide</span>
          </div>
          <p>Learn how to use OpenAutonomyX to publish content across multiple platforms.</p>
          <button className="btn-secondary">Edit</button>
        </div>

        <div className="blog-post">
          <h3>Platform Overview</h3>
          <div className="post-meta">
            <span>📅 June 26, 2026</span>
            <span>👤 You</span>
            <span>🏷️ Documentation</span>
          </div>
          <p>Understand the architecture and capabilities of the publishing platform.</p>
          <button className="btn-secondary">Edit</button>
        </div>
      </div>
    </div>
  );
}

export default ContentPage;
