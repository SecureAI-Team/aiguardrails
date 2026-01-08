-- Alert system tables

-- Alert rules (user-configurable)
CREATE TABLE IF NOT EXISTS alert_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID,                              -- NULL=platform-level rule
    name VARCHAR(100) NOT NULL,
    description TEXT,
    -- Trigger conditions
    event_types TEXT[] DEFAULT '{}',             -- ['blocked', 'login_failed', ...]
    severity_threshold VARCHAR(20) DEFAULT 'high', -- critical, high, medium, low
    threshold_count INT DEFAULT 1,               -- Trigger count threshold
    threshold_window_sec INT DEFAULT 60,         -- Time window in seconds
    -- Notification config
    notify_channels TEXT[] DEFAULT '{}',         -- ['sms', 'wechat', 'email', 'webhook']
    notify_recipients JSONB DEFAULT '{}',        -- Recipient configuration
    cooldown_sec INT DEFAULT 300,                -- Cooldown to prevent duplicates
    -- Status
    enabled BOOLEAN DEFAULT true,
    priority INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_alert_rules_tenant ON alert_rules(tenant_id);
CREATE INDEX IF NOT EXISTS idx_alert_rules_enabled ON alert_rules(enabled);

-- Alert history
CREATE TABLE IF NOT EXISTS alert_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_id UUID,
    rule_name VARCHAR(100),
    tenant_id UUID,
    severity VARCHAR(20) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT,
    event_data JSONB,
    -- Notification status
    notify_status JSONB DEFAULT '{}',            -- {"sms": "sent", "wechat": "failed"}
    -- Acknowledgement
    acknowledged BOOLEAN DEFAULT false,
    acknowledged_by UUID,
    acknowledged_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_alert_history_tenant ON alert_history(tenant_id);
CREATE INDEX IF NOT EXISTS idx_alert_history_severity ON alert_history(severity);
CREATE INDEX IF NOT EXISTS idx_alert_history_created ON alert_history(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_alert_history_ack ON alert_history(acknowledged);

-- Notification channels
CREATE TABLE IF NOT EXISTS notification_channels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID,
    channel_type VARCHAR(50) NOT NULL,           -- sms, wechat, wecom, email, webhook, dingtalk
    name VARCHAR(100) NOT NULL,
    config JSONB NOT NULL,                       -- Channel-specific configuration
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notification_channels_tenant ON notification_channels(tenant_id);
CREATE INDEX IF NOT EXISTS idx_notification_channels_type ON notification_channels(channel_type);

-- Insert default platform-level rules
INSERT INTO alert_rules (tenant_id, name, description, event_types, severity_threshold, threshold_count, threshold_window_sec, notify_channels, priority) VALUES
(NULL, 'prompt_injection_attack', '检测到提示注入攻击', ARRAY['prompt_injection'], 'critical', 1, 60, ARRAY['sms', 'wechat'], 100),
(NULL, 'brute_force_login', '暴力登录尝试', ARRAY['login_failed'], 'high', 5, 300, ARRAY['sms'], 90),
(NULL, 'sensitive_data_leak', '敏感数据泄露检测', ARRAY['dlp_violation', 'sensitive_data_exposed'], 'high', 1, 60, ARRAY['sms', 'wechat'], 95),
(NULL, 'opa_policy_block', 'OPA策略阻断', ARRAY['opa_block'], 'medium', 3, 60, ARRAY['webhook'], 50),
(NULL, 'rate_limit_exceeded', '请求速率超限', ARRAY['rate_limit'], 'medium', 10, 60, ARRAY['webhook'], 40),
(NULL, 'industrial_safety_violation', '工业安全违规', ARRAY['industrial_safety'], 'critical', 1, 60, ARRAY['sms', 'wechat'], 100)
ON CONFLICT DO NOTHING;
