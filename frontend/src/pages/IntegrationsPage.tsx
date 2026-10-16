import { useEffect } from 'react';
import './Pages.css';

interface PageProps {
  onContextChange: (context: string, data: any) => void;
}

function IntegrationsPage({ onContextChange }: PageProps) {
  useEffect(() => {
    onContextChange('integrations', {});
  }, [onContextChange]);

  return (
    <div className="page-container">
      <div className="page-header">
        <h1>🔌 Integrations</h1>
        <p>Connect with external platforms and publish everywhere</p>
        <button className="btn-action">+ Add Integration</button>
      </div>

      <div className="content-section">
        <div className="blog-post">
          <h3>📘 WordPress</h3>
          <div className="post-meta">
            <span>Status: ✅ Connected</span>
            <span>Domain: blog.publishing.openautonomyx.com</span>
          </div>
          <p>Publish directly to your WordPress blog</p>
          <button className="btn-secondary">Configure</button>
        </div>

        <div className="blog-post">
          <h3>📱 Medium</h3>
          <div className="post-meta">
            <span>Status: ⭕ Disconnected</span>
          </div>
          <p>Reach millions of readers on Medium</p>
          <button className="btn-secondary">Connect</button>
        </div>

        <div className="blog-post">
          <h3>📬 Substack</h3>
          <div className="post-meta">
            <span>Status: ⭕ Disconnected</span>
          </div>
          <p>Build your newsletter audience on Substack</p>
          <button className="btn-secondary">Connect</button>
        </div>

        <div className="blog-post">
          <h3>𝕏 Twitter</h3>
          <div className="post-meta">
            <span>Status: ⭕ Disconnected</span>
          </div>
          <p>Share your content with Twitter followers</p>
          <button className="btn-secondary">Connect</button>
        </div>

        <div className="blog-post">
          <h3>💼 LinkedIn</h3>
          <div className="post-meta">
            <span>Status: ⭕ Disconnected</span>
          </div>
          <p>Publish professional content on LinkedIn</p>
          <button className="btn-secondary">Connect</button>
        </div>

        <div className="blog-post">
          <h3>👥 Facebook</h3>
          <div className="post-meta">
            <span>Status: ⭕ Disconnected</span>
          </div>
          <p>Reach your Facebook audience</p>
          <button className="btn-secondary">Connect</button>
        </div>
      </div>
    </div>
  );
}

export default IntegrationsPage;
