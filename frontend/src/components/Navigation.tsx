import { Link } from 'react-router-dom';
import './Navigation.css';

interface NavigationProps {
  onNavigate: (route: string) => void;
}

function Navigation({ onNavigate }: NavigationProps) {
  return (
    <nav className="navbar">
      <div className="nav-container">
        <Link to="/" className="nav-logo" onClick={() => onNavigate('home')}>
          <span className="logo-icon">🚀</span>
          <span className="logo-text">OpenAutonomyX</span>
        </Link>

        <ul className="nav-menu">
          <li>
            <Link to="/" className="nav-link" onClick={() => onNavigate('home')}>
              Home
            </Link>
          </li>
          <li>
            <Link to="/content" className="nav-link" onClick={() => onNavigate('content')}>
              Content
            </Link>
          </li>
          <li>
            <Link to="/blog" className="nav-link" onClick={() => onNavigate('blog')}>
              Blog
            </Link>
          </li>
          <li>
            <Link to="/formats" className="nav-link" onClick={() => onNavigate('formats')}>
              Formats
            </Link>
          </li>
          <li>
            <Link to="/integrations" className="nav-link" onClick={() => onNavigate('integrations')}>
              Integrations
            </Link>
          </li>
          <li>
            <Link to="/admin" className="nav-link" onClick={() => onNavigate('admin')}>
              Admin
            </Link>
          </li>
        </ul>

        <div className="nav-auth">
          <button className="btn-login">Login</button>
          <button className="btn-signup">Sign Up</button>
        </div>
      </div>
    </nav>
  );
}

export default Navigation;
