import './RightPanel.css';

interface RightPanelProps {
  context: string;
  data: any;
}

function RightPanel({ context }: RightPanelProps) {
  const renderContextualContent = () => {
    switch (context) {
      case 'home':
        return (
          <div className="panel-section">
            <h3>📊 Quick Stats</h3>
            <div className="stat-item">
              <span className="stat-label">Services Running</span>
              <span className="stat-value">21</span>
            </div>
            <div className="stat-item">
              <span className="stat-label">Domains</span>
              <span className="stat-value">Multiple</span>
            </div>
            <div className="stat-item">
              <span className="stat-label">Formats</span>
              <span className="stat-value">6+</span>
            </div>
            <div className="stat-item">
              <span className="stat-label">Status</span>
              <span className="stat-value status-active">● Active</span>
            </div>
          </div>
        );

      case 'blog':
        return (
          <div className="panel-section">
            <h3>📝 Blog Info</h3>
            <div className="info-item">
              <strong>Platform:</strong>
              <span>WordPress</span>
            </div>
            <div className="info-item">
              <strong>Domain:</strong>
              <span>blog.publishing.openautonomyx.com</span>
            </div>
            <div className="info-item">
              <strong>Status:</strong>
              <span className="badge-active">Connected</span>
            </div>
            <button className="btn-secondary" style={{ width: '100%', marginTop: '10px' }}>
              🔗 Visit Blog
            </button>
          </div>
        );

      case 'content':
        return (
          <div className="panel-section">
            <h3>📄 Content Tools</h3>
            <div className="tools-list">
              <button className="tool-btn">✏️ Create Post</button>
              <button className="tool-btn">📚 My Content</button>
              <button className="tool-btn">🏷️ Tags</button>
              <button className="tool-btn">📂 Categories</button>
            </div>
          </div>
        );

      case 'formats':
        return (
          <div className="panel-section">
            <h3>🎬 Format Converter</h3>
            <div className="formats-list">
              <div className="format-item">
                <span>📚 EPUB</span>
                <small>E-books</small>
              </div>
              <div className="format-item">
                <span>📊 Slides</span>
                <small>Presentations</small>
              </div>
              <div className="format-item">
                <span>📄 PDF</span>
                <small>Documents</small>
              </div>
              <div className="format-item">
                <span>🎧 Audio</span>
                <small>Podcasts</small>
              </div>
              <div className="format-item">
                <span>🎬 Video</span>
                <small>Videos</small>
              </div>
              <div className="format-item">
                <span>🌐 HTML</span>
                <small>Web</small>
              </div>
            </div>
          </div>
        );

      case 'integrations':
        return (
          <div className="panel-section">
            <h3>🔌 Integrations</h3>
            <div className="integrations-list">
              <div className="integration-item">
                <span>📘 WordPress</span>
              </div>
              <div className="integration-item">
                <span>📱 Medium</span>
              </div>
              <div className="integration-item">
                <span>📬 Substack</span>
              </div>
              <div className="integration-item">
                <span>𝕏 Twitter</span>
              </div>
              <div className="integration-item">
                <span>💼 LinkedIn</span>
              </div>
              <div className="integration-item">
                <span>👥 Facebook</span>
              </div>
            </div>
          </div>
        );

      case 'admin':
        return (
          <div className="panel-section">
            <h3>⚙️ Admin Panel</h3>
            <div className="admin-menu">
              <button className="admin-btn">🔧 Settings</button>
              <button className="admin-btn">👥 Users</button>
              <button className="admin-btn">📊 Analytics</button>
              <button className="admin-btn">🔐 Security</button>
              <button className="admin-btn">📦 Services</button>
            </div>
          </div>
        );

      default:
        return (
          <div className="panel-section">
            <h3>ℹ️ Info</h3>
            <p>Select a section to view details</p>
          </div>
        );
    }
  };

  return (
    <div className="right-panel">
      <div className="panel-header">
        <h2>Context</h2>
      </div>

      <div className="panel-content">
        {renderContextualContent()}

        <div className="panel-section" style={{ marginTop: '30px' }}>
          <h3>🚀 Getting Started</h3>
          <ul className="getting-started">
            <li><a href="#docs">📚 Documentation</a></li>
            <li><a href="#api">🔌 API Reference</a></li>
            <li><a href="#support">💬 Support</a></li>
            <li><a href="#github">🐙 GitHub</a></li>
          </ul>
        </div>
      </div>
    </div>
  );
}

export default RightPanel;
