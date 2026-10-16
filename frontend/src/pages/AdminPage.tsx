import { useEffect } from 'react';
import './Pages.css';

interface PageProps {
  onContextChange: (context: string, data: any) => void;
}

function AdminPage({ onContextChange }: PageProps) {
  useEffect(() => {
    onContextChange('admin', {});
  }, [onContextChange]);

  return (
    <div className="page-container">
      <div className="page-header">
        <h1>⚙️ Admin Panel</h1>
        <p>Manage your platform settings and configurations</p>
      </div>

      <div className="content-section">
        <div className="blog-post">
          <h3>🔧 Settings</h3>
          <p>Configure platform-wide settings and preferences</p>
          <button className="btn-secondary">Open Settings</button>
        </div>

        <div className="blog-post">
          <h3>👥 Users</h3>
          <p>Manage user accounts and permissions</p>
          <button className="btn-secondary">Manage Users</button>
        </div>

        <div className="blog-post">
          <h3>📊 Analytics</h3>
          <p>View analytics and performance metrics</p>
          <button className="btn-secondary">View Analytics</button>
        </div>

        <div className="blog-post">
          <h3>🔐 Security</h3>
          <p>Configure security settings and API keys</p>
          <button className="btn-secondary">Security Settings</button>
        </div>

        <div className="blog-post">
          <h3>📦 Services</h3>
          <p>Manage and monitor running services</p>
          <button className="btn-secondary">Service Status</button>
        </div>
      </div>
    </div>
  );
}

export default AdminPage;
