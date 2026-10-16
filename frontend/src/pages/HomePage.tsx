import { useEffect } from 'react';
import './HomePage.css';

interface HomePageProps {
  onContextChange: (context: string, data: any) => void;
}

function HomePage({ onContextChange }: HomePageProps) {
  useEffect(() => {
    onContextChange('home', {});
  }, [onContextChange]);

  return (
    <div className="home-page">
      {/* Hero Section */}
      <section className="hero">
        <div className="hero-content">
          <h1 className="hero-title">
            Create Once<br />
            <span>Publish Everywhere</span>
          </h1>
          <p className="hero-subtitle">
            Vendor-neutral creative publishing platform<br />
            powered by local AI, unlimited formats, infinite possibilities
          </p>
          <div className="hero-buttons">
            <button className="btn-primary">Get Started</button>
            <button className="btn-secondary">Learn More</button>
          </div>
        </div>
        <div className="hero-features">
          <div className="feature-badge">🚀 21 Services</div>
          <div className="feature-badge">📚 6+ Formats</div>
          <div className="feature-badge">🌐 Multi-Domain</div>
        </div>
      </section>

      {/* Features Section */}
      <section className="features">
        <h2>Powerful Features</h2>
        <div className="features-grid">
          <div className="feature-card">
            <div className="feature-icon">📝</div>
            <h3>Content Management</h3>
            <p>Create, edit, and manage content with markdown support and version control</p>
          </div>

          <div className="feature-card">
            <div className="feature-icon">🎬</div>
            <h3>Multi-Format Export</h3>
            <p>Generate EPUB, PDF, Slides, Audio, Video, and HTML from single content</p>
          </div>

          <div className="feature-card">
            <div className="feature-icon">🔌</div>
            <h3>Plugin Architecture</h3>
            <p>Publish to WordPress, Medium, Substack, Twitter, LinkedIn, and more</p>
          </div>

          <div className="feature-card">
            <div className="feature-icon">🤖</div>
            <h3>Local AI</h3>
            <p>Run Ollama locally or connect to OpenAI, Anthropic, Azure, and others</p>
          </div>

          <div className="feature-card">
            <div className="feature-icon">🔐</div>
            <h3>Vendor Neutral</h3>
            <p>Zero lock-in. Full control over your data. SBOM-based supply chain security</p>
          </div>

          <div className="feature-card">
            <div className="feature-icon">📊</div>
            <h3>Analytics</h3>
            <p>Track performance across all platforms with unified analytics dashboard</p>
          </div>
        </div>
      </section>

      {/* Services Section */}
      <section className="services">
        <h2>21 Microservices Architecture</h2>
        <div className="services-grid">
          <div className="service-group">
            <h4>Infrastructure (5)</h4>
            <ul>
              <li>PostgreSQL</li>
              <li>Redis</li>
              <li>Elasticsearch</li>
              <li>MinIO</li>
              <li>Ollama</li>
            </ul>
          </div>

          <div className="service-group">
            <h4>Core (2)</h4>
            <ul>
              <li>API Gateway</li>
              <li>Event Bus</li>
            </ul>
          </div>

          <div className="service-group">
            <h4>Publishing (3)</h4>
            <ul>
              <li>Blog (3009)</li>
              <li>Integrations (3010)</li>
              <li>Formats (3011)</li>
            </ul>
          </div>

          <div className="service-group">
            <h4>Business Modules (8)</h4>
            <ul>
              <li>Content</li>
              <li>Skills</li>
              <li>Tools</li>
              <li>Analytics</li>
              <li>Optimization</li>
              <li>Design</li>
              <li>Features</li>
              <li>+1 More</li>
            </ul>
          </div>
        </div>
      </section>

      {/* Tech Stack */}
      <section className="tech-stack">
        <h2>Built With Modern Tech</h2>
        <div className="stack-grid">
          <div className="stack-item">
            <span className="stack-icon">⚙️</span>
            <h4>Backend</h4>
            <p>Node.js, Express, TypeScript</p>
          </div>
          <div className="stack-item">
            <span className="stack-icon">🎨</span>
            <h4>Frontend</h4>
            <p>React, TypeScript, Vite</p>
          </div>
          <div className="stack-item">
            <span className="stack-icon">🗄️</span>
            <h4>Database</h4>
            <p>PostgreSQL, Redis, Elasticsearch</p>
          </div>
          <div className="stack-item">
            <span className="stack-icon">📦</span>
            <h4>Deployment</h4>
            <p>Docker, Kubernetes, GitHub Pages</p>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="cta">
        <h2>Ready to Transform Your Publishing?</h2>
        <p>Join the creative revolution with vendor-neutral publishing</p>
        <button className="btn-large">Start Your Free Trial</button>
      </section>
    </div>
  );
}

export default HomePage;
