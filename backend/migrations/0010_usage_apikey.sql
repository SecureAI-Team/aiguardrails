-- Usage statistics and API key management

-- API usage statistics (aggregated)
CREATE TABLE IF NOT EXISTS api_usage_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    app_id UUID,
    date DATE NOT NULL,
    hour SMALLINT,                       -- 0-23, NULL for daily aggregates
    -- Metrics
    request_count BIGINT DEFAULT 0,
    success_count BIGINT DEFAULT 0,
    error_count BIGINT DEFAULT 0,
    blocked_count BIGINT DEFAULT 0,
    -- Latency (milliseconds)
    latency_sum BIGINT DEFAULT 0,        -- For calculating average
    latency_p50 INT DEFAULT 0,
    latency_p99 INT DEFAULT 0,
    -- Token usage
    input_tokens BIGINT DEFAULT 0,
    output_tokens BIGINT DEFAULT 0,
    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, app_id, date, hour)
);

CREATE INDEX IF NOT EXISTS idx_usage_stats_tenant_date ON api_usage_stats(tenant_id, date DESC);
CREATE INDEX IF NOT EXISTS idx_usage_stats_app ON api_usage_stats(app_id, date DESC);

-- API Keys
CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    app_id UUID,
    name VARCHAR(100) NOT NULL,
    key_prefix VARCHAR(8) NOT NULL,      -- First 8 chars for display
    key_hash VARCHAR(64) NOT NULL,       -- SHA256 hash of the full key
    -- Permissions
    scopes TEXT[] DEFAULT '{}',          -- ['read', 'write', 'admin']
    ip_whitelist TEXT[] DEFAULT '{}',    -- Allowed IPs (empty = all)
    -- Limits
    rate_limit_rpm INT,                  -- Requests per minute
    rate_limit_rpd INT,                  -- Requests per day
    -- Status
    enabled BOOLEAN DEFAULT true,
    expires_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ,
    -- Metadata
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_api_keys_tenant ON api_keys(tenant_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_prefix ON api_keys(key_prefix);
CREATE INDEX IF NOT EXISTS idx_api_keys_hash ON api_keys(key_hash);

-- Quota configuration
CREATE TABLE IF NOT EXISTS quota_config (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL UNIQUE,
    -- Quotas
    daily_request_limit BIGINT,
    monthly_request_limit BIGINT,
    daily_token_limit BIGINT,
    monthly_token_limit BIGINT,
    -- Current usage (updated by triggers/jobs)
    current_daily_requests BIGINT DEFAULT 0,
    current_monthly_requests BIGINT DEFAULT 0,
    current_daily_tokens BIGINT DEFAULT 0,
    current_monthly_tokens BIGINT DEFAULT 0,
    -- Reset tracking
    last_daily_reset DATE,
    last_monthly_reset DATE,
    -- Alerts
    alert_threshold_percent INT DEFAULT 80,  -- Alert when usage > 80%
    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_quota_tenant ON quota_config(tenant_id);

-- API key usage log (for auditing)
CREATE TABLE IF NOT EXISTS api_key_usage_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key_id UUID NOT NULL,
    endpoint VARCHAR(255),
    method VARCHAR(10),
    status_code INT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    request_tokens INT,
    response_tokens INT,
    latency_ms INT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_key_usage_key ON api_key_usage_log(key_id, created_at DESC);
