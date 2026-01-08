-- Enterprise features: Organizations, Teams, SSO, IP Whitelist

-- Organizations (顶层组织)
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(50) NOT NULL UNIQUE,          -- URL-friendly identifier
    description TEXT,
    logo_url VARCHAR(500),
    -- Settings
    settings JSONB DEFAULT '{}',
    -- SSO config
    sso_enabled BOOLEAN DEFAULT false,
    sso_provider VARCHAR(50),                  -- saml, oidc
    sso_config JSONB,                          -- Provider-specific config
    -- Limits
    max_teams INT DEFAULT 10,
    max_members INT DEFAULT 100,
    max_tenants INT DEFAULT 5,
    -- Status
    status VARCHAR(20) DEFAULT 'active',       -- active, suspended, deleted
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_org_slug ON organizations(slug);

-- Teams (组织内的团队)
CREATE TABLE IF NOT EXISTS teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    -- Permissions inherited by members
    default_role VARCHAR(50) DEFAULT 'member', -- owner, admin, member, viewer
    permissions TEXT[] DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(org_id, name)
);

CREATE INDEX IF NOT EXISTS idx_team_org ON teams(org_id);

-- Organization members
CREATE TABLE IF NOT EXISTS org_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    -- Role
    role VARCHAR(50) NOT NULL DEFAULT 'member', -- owner, admin, member, viewer
    -- Team assignments
    team_ids UUID[] DEFAULT '{}',
    -- Status
    status VARCHAR(20) DEFAULT 'active',        -- active, invited, suspended
    invited_by UUID,
    invited_at TIMESTAMPTZ,
    joined_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(org_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_member_org ON org_members(org_id);
CREATE INDEX IF NOT EXISTS idx_member_user ON org_members(user_id);

-- SSO sessions
CREATE TABLE IF NOT EXISTS sso_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    session_id VARCHAR(255) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255),
    attributes JSONB,                           -- SAML attributes or OIDC claims
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sso_session ON sso_sessions(session_id);
CREATE INDEX IF NOT EXISTS idx_sso_expires ON sso_sessions(expires_at);

-- IP Whitelist
CREATE TABLE IF NOT EXISTS ip_whitelist (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Scope: org, tenant, or api_key level
    scope_type VARCHAR(20) NOT NULL,            -- org, tenant, api_key
    scope_id UUID NOT NULL,
    -- IP config
    ip_address VARCHAR(45),                     -- Single IP
    ip_cidr VARCHAR(50),                        -- CIDR range
    description VARCHAR(255),
    enabled BOOLEAN DEFAULT true,
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_whitelist_scope ON ip_whitelist(scope_type, scope_id);

-- Link tenants to organizations
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS org_id UUID REFERENCES organizations(id);
CREATE INDEX IF NOT EXISTS idx_tenant_org ON tenants(org_id);

-- Audit log for admin operations
CREATE TABLE IF NOT EXISTS admin_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID,
    user_id UUID NOT NULL,
    action VARCHAR(100) NOT NULL,               -- create_org, invite_member, update_sso, etc
    resource_type VARCHAR(50),                  -- org, team, member, sso, etc
    resource_id UUID,
    old_value JSONB,
    new_value JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_org ON admin_audit_log(org_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_user ON admin_audit_log(user_id, created_at DESC);
