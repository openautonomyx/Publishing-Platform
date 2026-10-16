import { useState } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Navigation from './components/Navigation';
import RightPanel from './components/RightPanel';
import HomePage from './pages/HomePage';
import BlogPage from './pages/BlogPage';
import ContentPage from './pages/EditorialWorkspace';
import IntegrationsPage from './pages/IntegrationsPage';
import FormatsPage from './pages/FormatsPage';
import AdminPage from './pages/AdminPage';
import './App.css';

function App() {
  const [activeContext, setActiveContext] = useState<string>('home');
  const [contextData, setContextData] = useState<any>({});

  return (
    <Router>
      <div className="app-container">
        <Navigation onNavigate={setActiveContext} />

        <div className="main-layout">
          <div className="main-content">
            <Routes>
              <Route path="/" element={<HomePage onContextChange={(ctx, data) => { setActiveContext(ctx); setContextData(data); }} />} />
              <Route path="/blog" element={<BlogPage onContextChange={(ctx, data) => { setActiveContext(ctx); setContextData(data); }} />} />
              <Route path="/content" element={<ContentPage onContextChange={(ctx, data) => { setActiveContext(ctx); setContextData(data); }} />} />
              <Route path="/integrations" element={<IntegrationsPage onContextChange={(ctx, data) => { setActiveContext(ctx); setContextData(data); }} />} />
              <Route path="/formats" element={<FormatsPage onContextChange={(ctx, data) => { setActiveContext(ctx); setContextData(data); }} />} />
              <Route path="/admin" element={<AdminPage onContextChange={(ctx, data) => { setActiveContext(ctx); setContextData(data); }} />} />
            </Routes>
          </div>

          <RightPanel context={activeContext} data={contextData} />
        </div>
      </div>
    </Router>
  );
}

export default App;
