-- Tenant-level configurable rules
CREATE TABLE IF NOT EXISTS tenant_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    rule_type VARCHAR(50) NOT NULL,  -- 'business' | 'permission'
    name VARCHAR(255) NOT NULL,
    description TEXT,
    config JSONB NOT NULL DEFAULT '{}'::jsonb,
    enabled BOOLEAN DEFAULT true,
    priority INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by VARCHAR(255),
    UNIQUE(tenant_id, rule_type, name)
);

CREATE INDEX IF NOT EXISTS idx_tenant_rules_tenant ON tenant_rules(tenant_id);
CREATE INDEX IF NOT EXISTS idx_tenant_rules_type ON tenant_rules(rule_type);
CREATE INDEX IF NOT EXISTS idx_tenant_rules_enabled ON tenant_rules(tenant_id, enabled);

-- Rule templates (platform-provided templates for tenants)
CREATE TABLE IF NOT EXISTS rule_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    rule_type VARCHAR(50) NOT NULL,
    description TEXT,
    config_schema JSONB NOT NULL DEFAULT '{}'::jsonb,  -- JSON schema for validation
    default_config JSONB NOT NULL DEFAULT '{}'::jsonb,
    tags TEXT[] DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Insert default templates
INSERT INTO rule_templates (id, name, rule_type, description, config_schema, default_config, tags) VALUES
(gen_random_uuid(), 'vendor_exclusive', 'business', '厂商专属助手 - 限制只回答特定厂商产品', 
 '{"type":"object","properties":{"allowed_vendors":{"type":"array"},"blocked_vendors":{"type":"array"}}}',
 '{"enabled":true,"mode":"exclusive","allowed_vendors":[],"blocked_vendors":[],"responses":{"blocked":"抱歉，我无法回答这个问题。"}}',
 ARRAY['vendor', 'exclusive']),
 
(gen_random_uuid(), 'domain_boundary', 'business', '领域知识边界 - 限制回答特定领域问题',
 '{"type":"object","properties":{"allowed_topics":{"type":"array"},"blocked_topics":{"type":"array"}}}',
 '{"enabled":true,"allowed_topics":[],"blocked_topics":[],"responses":{"out_of_scope":"这个问题超出了我的专业范围。"}}',
 ARRAY['domain', 'boundary']),

(gen_random_uuid(), 'custom_roles', 'permission', '自定义角色权限 - 定义租户内角色层级',
 '{"type":"object","properties":{"roles":{"type":"object"}}}',
 '{"roles":{"admin":{"level":3,"permissions":["*"]},"user":{"level":1,"permissions":["chat"]}}}',
 ARRAY['role', 'permission']),

(gen_random_uuid(), 'tool_acl', 'permission', '工具访问控制 - 控制工具使用权限',
 '{"type":"object","properties":{"tool_permissions":{"type":"object"}}}',
 '{"tool_permissions":{"read_status":{"min_level":1},"set_parameter":{"min_level":2,"requires_confirmation":true}}}',
 ARRAY['tool', 'acl', 'permission'])
ON CONFLICT (name) DO NOTHING;
