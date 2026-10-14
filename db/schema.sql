-- Universal Creative Platform - Database Schema
-- PostgreSQL 15+, Multi-tenant with Row-Level Security, ACID compliant

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================================
-- Organizations (Tenants)
-- ============================================================================

CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    tier VARCHAR(50) DEFAULT 'free', -- free, pro, enterprise
    status VARCHAR(50) DEFAULT 'active', -- active, suspended, deleted
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_organizations_slug ON organizations(slug);
CREATE INDEX idx_organizations_status ON organizations(status);

-- ============================================================================
-- Users (Multi-tenant)
-- ============================================================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255),
    avatar_url TEXT,
    role VARCHAR(50) DEFAULT 'user', -- admin, editor, viewer, user
    status VARCHAR(50) DEFAULT 'active', -- active, inactive, deleted
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(org_id, email)
);

CREATE INDEX idx_users_org_id ON users(org_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);

-- ============================================================================
-- API Keys (for agent authentication)
-- ============================================================================

CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key_hash VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255),
    permissions TEXT[] DEFAULT ARRAY['read', 'write'],
    status VARCHAR(50) DEFAULT 'active', -- active, revoked
    last_used_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_api_keys_org_id ON api_keys(org_id);
CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_status ON api_keys(status);

-- ============================================================================
-- Workflows (Creative workflows with approval stages)
-- ============================================================================

CREATE TABLE workflows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    workflow_type VARCHAR(50) NOT NULL, -- content-creation, approval, publishing, distribution
    status VARCHAR(50) DEFAULT 'draft', -- draft, active, archived, deleted
    config JSONB NOT NULL DEFAULT '{}', -- workflow-specific configuration
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_workflows_org_id ON workflows(org_id);
CREATE INDEX idx_workflows_type ON workflows(workflow_type);
CREATE INDEX idx_workflows_status ON workflows(status);

-- ============================================================================
-- Approval Stages (Sequential stages within workflow)
-- ============================================================================

CREATE TABLE approval_stages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    sequence_order INT NOT NULL,
    approver_role VARCHAR(50) NOT NULL, -- admin, manager, editor
    required_count INT DEFAULT 1,
    timeout_hours INT DEFAULT 48,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_approval_stages_workflow_id ON approval_stages(workflow_id);
CREATE INDEX idx_approval_stages_order ON approval_stages(workflow_id, sequence_order);

-- ============================================================================
-- Content (Pieces of creative work)
-- ============================================================================

CREATE TABLE content (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    workflow_id UUID REFERENCES workflows(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    content_type VARCHAR(50) NOT NULL, -- article, video, image, podcast, social
    status VARCHAR(50) DEFAULT 'draft', -- draft, in-review, approved, published, archived
    metadata JSONB DEFAULT '{}', -- type-specific metadata
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_content_org_id ON content(org_id);
CREATE INDEX idx_content_workflow_id ON content(workflow_id);
CREATE INDEX idx_content_status ON content(status);
CREATE INDEX idx_content_type ON content(content_type);
CREATE INDEX idx_content_created_by ON content(created_by);

-- ============================================================================
-- Approvals (Approval decisions for content)
-- ============================================================================

CREATE TABLE approvals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    stage_id UUID NOT NULL REFERENCES approval_stages(id),
    approver_id UUID NOT NULL REFERENCES users(id),
    decision VARCHAR(50) NOT NULL, -- approved, rejected, requested-changes
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_approvals_content_id ON approvals(content_id);
CREATE INDEX idx_approvals_stage_id ON approvals(stage_id);
CREATE INDEX idx_approvals_approver_id ON approvals(approver_id);
CREATE INDEX idx_approvals_decision ON approvals(decision);

-- ============================================================================
-- Agents (Autonomous agents for content creation)
-- ============================================================================

CREATE TABLE agents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    agent_type VARCHAR(50) NOT NULL, -- content-creator, approver, publisher, distributor
    description TEXT,
    config JSONB NOT NULL DEFAULT '{}', -- agent-specific configuration
    status VARCHAR(50) DEFAULT 'active', -- active, paused, error, deleted
    last_run_at TIMESTAMP,
    error_message TEXT,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_agents_org_id ON agents(org_id);
CREATE INDEX idx_agents_type ON agents(agent_type);
CREATE INDEX idx_agents_status ON agents(status);

-- ============================================================================
-- Agent Runs (Execution history)
-- ============================================================================

CREATE TABLE agent_runs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    content_id UUID REFERENCES content(id) ON DELETE SET NULL,
    status VARCHAR(50) NOT NULL, -- running, success, failed, timeout
    input_data JSONB DEFAULT '{}',
    output_data JSONB DEFAULT '{}',
    error_message TEXT,
    duration_ms INT,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE INDEX idx_agent_runs_agent_id ON agent_runs(agent_id);
CREATE INDEX idx_agent_runs_content_id ON agent_runs(content_id);
CREATE INDEX idx_agent_runs_status ON agent_runs(status);
CREATE INDEX idx_agent_runs_started_at ON agent_runs(started_at DESC);

-- ============================================================================
-- Distributions (Publishing to external platforms)
-- ============================================================================

CREATE TABLE distributions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    channel VARCHAR(100) NOT NULL, -- twitter, linkedin, medium, substack, website
    status VARCHAR(50) DEFAULT 'pending', -- pending, published, failed, scheduled
    published_url TEXT,
    published_at TIMESTAMP,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_distributions_org_id ON distributions(org_id);
CREATE INDEX idx_distributions_content_id ON distributions(content_id);
CREATE INDEX idx_distributions_channel ON distributions(channel);
CREATE INDEX idx_distributions_status ON distributions(status);

-- ============================================================================
-- Audit Log (Compliance and debugging)
-- ============================================================================

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(255) NOT NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    changes JSONB DEFAULT '{}',
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_org_id ON audit_logs(org_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);

-- ============================================================================
-- Row-Level Security (Multi-tenant isolation)
-- ============================================================================

ALTER TABLE organizations ENABLE ROW LEVEL SECURITY;
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE api_keys ENABLE ROW LEVEL SECURITY;
ALTER TABLE workflows ENABLE ROW LEVEL SECURITY;
ALTER TABLE approval_stages ENABLE ROW LEVEL SECURITY;
ALTER TABLE content ENABLE ROW LEVEL SECURITY;
ALTER TABLE approvals ENABLE ROW LEVEL SECURITY;
ALTER TABLE agents ENABLE ROW LEVEL SECURITY;
ALTER TABLE agent_runs ENABLE ROW LEVEL SECURITY;
ALTER TABLE distributions ENABLE ROW LEVEL SECURITY;
ALTER TABLE audit_logs ENABLE ROW LEVEL SECURITY;

-- RLS Policies for organizations (super users only)
CREATE POLICY org_isolation ON organizations
    USING (TRUE); -- Admin-only in application layer

-- RLS Policies for users
CREATE POLICY user_org_isolation ON users
    USING (org_id = CURRENT_SETTING('app.current_org_id')::UUID);

CREATE POLICY user_insert ON users
    WITH CHECK (org_id = CURRENT_SETTING('app.current_org_id')::UUID);

-- RLS Policies for content
CREATE POLICY content_org_isolation ON content
    USING (org_id = CURRENT_SETTING('app.current_org_id')::UUID);

CREATE POLICY content_insert ON content
    WITH CHECK (org_id = CURRENT_SETTING('app.current_org_id')::UUID);

-- ============================================================================
-- Functions
-- ============================================================================

-- Update timestamp on modification
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for updated_at
CREATE TRIGGER update_organizations_updated_at BEFORE UPDATE ON organizations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_workflows_updated_at BEFORE UPDATE ON workflows
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_content_updated_at BEFORE UPDATE ON content
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_approvals_updated_at BEFORE UPDATE ON approvals
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_agents_updated_at BEFORE UPDATE ON agents
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_distributions_updated_at BEFORE UPDATE ON distributions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- Views
-- ============================================================================

-- Content with approval status
CREATE VIEW v_content_approval_status AS
SELECT
    c.id,
    c.org_id,
    c.title,
    c.status,
    COUNT(CASE WHEN a.decision = 'approved' THEN 1 END) as approvals_count,
    COUNT(CASE WHEN a.decision = 'rejected' THEN 1 END) as rejections_count,
    MAX(a.created_at) as last_approval_at
FROM content c
LEFT JOIN approvals a ON c.id = a.content_id
GROUP BY c.id, c.org_id, c.title, c.status;

-- Agent execution statistics
CREATE VIEW v_agent_statistics AS
SELECT
    a.id,
    a.name,
    a.org_id,
    COUNT(*) as total_runs,
    COUNT(CASE WHEN ar.status = 'success' THEN 1 END) as successful_runs,
    COUNT(CASE WHEN ar.status = 'failed' THEN 1 END) as failed_runs,
    AVG(ar.duration_ms) as avg_duration_ms,
    MAX(ar.started_at) as last_run_at
FROM agents a
LEFT JOIN agent_runs ar ON a.id = ar.agent_id
GROUP BY a.id, a.name, a.org_id;
