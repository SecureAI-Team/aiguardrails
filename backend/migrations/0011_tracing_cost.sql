-- P1/P2 Feature tables: Tracing, Cost, Models

-- Request traces (full request lifecycle)
CREATE TABLE IF NOT EXISTS request_traces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trace_id VARCHAR(64) NOT NULL,           -- Distributed trace ID
    span_id VARCHAR(32),                     -- Span ID
    parent_span_id VARCHAR(32),
    tenant_id UUID NOT NULL,
    app_id UUID,
    -- Request info
    method VARCHAR(10),
    path VARCHAR(255),
    query_params JSONB,
    headers JSONB,                           -- Filtered headers
    request_body TEXT,                       -- Truncated if large
    -- Response info
    status_code INT,
    response_body TEXT,                      -- Truncated if large
    -- Timing
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    duration_ms INT,
    -- Processing stages
    stages JSONB,                            -- [{name, start, end, status}]
    -- Security
    blocked BOOLEAN DEFAULT false,
    block_reason VARCHAR(255),
    signals TEXT[],
    -- Tokens
    input_tokens INT DEFAULT 0,
    output_tokens INT DEFAULT 0,
    -- Error
    error TEXT,
    -- Metadata
    user_agent TEXT,
    client_ip VARCHAR(45),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_traces_tenant ON request_traces(tenant_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_traces_trace_id ON request_traces(trace_id);
CREATE INDEX IF NOT EXISTS idx_traces_blocked ON request_traces(blocked, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_traces_status ON request_traces(status_code);

-- Cost records
CREATE TABLE IF NOT EXISTS cost_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    app_id UUID,
    date DATE NOT NULL,
    -- Usage
    request_count BIGINT DEFAULT 0,
    input_tokens BIGINT DEFAULT 0,
    output_tokens BIGINT DEFAULT 0,
    -- Cost (in cents)
    input_cost_cents BIGINT DEFAULT 0,
    output_cost_cents BIGINT DEFAULT 0,
    total_cost_cents BIGINT DEFAULT 0,
    -- Pricing config at time of record
    pricing_config JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, app_id, date)
);

CREATE INDEX IF NOT EXISTS idx_cost_tenant_date ON cost_records(tenant_id, date DESC);

-- Model catalog
CREATE TABLE IF NOT EXISTS model_catalog (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider VARCHAR(50) NOT NULL,           -- openai, anthropic, qwen, etc
    model_id VARCHAR(100) NOT NULL,          -- gpt-4, claude-3, etc
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    -- Capabilities
    capabilities TEXT[],                     -- ['chat', 'completion', 'embedding']
    context_window INT,
    max_output_tokens INT,
    -- Pricing (per 1M tokens)
    input_price_per_m DECIMAL(10,4),
    output_price_per_m DECIMAL(10,4),
    -- Status
    enabled BOOLEAN DEFAULT true,
    deprecated BOOLEAN DEFAULT false,
    -- Metadata
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(provider, model_id)
);

-- Insert default models
INSERT INTO model_catalog (provider, model_id, display_name, description, capabilities, context_window, input_price_per_m, output_price_per_m) VALUES
('qwen', 'qwen-turbo', '通义千问-Turbo', '高性价比通用模型', ARRAY['chat'], 8192, 0.8, 2.0),
('qwen', 'qwen-plus', '通义千问-Plus', '增强推理能力', ARRAY['chat'], 32768, 4.0, 12.0),
('qwen', 'qwen-max', '通义千问-Max', '最强能力模型', ARRAY['chat'], 32768, 40.0, 120.0),
('openai', 'gpt-4o', 'GPT-4o', 'OpenAI最新多模态模型', ARRAY['chat', 'vision'], 128000, 5.0, 15.0),
('openai', 'gpt-4o-mini', 'GPT-4o Mini', '快速经济型模型', ARRAY['chat'], 128000, 0.15, 0.6),
('anthropic', 'claude-3-5-sonnet', 'Claude 3.5 Sonnet', 'Anthropic均衡模型', ARRAY['chat', 'vision'], 200000, 3.0, 15.0)
ON CONFLICT DO NOTHING;

-- Export jobs
CREATE TABLE IF NOT EXISTS export_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,               -- usage, traces, logs, cost
    status VARCHAR(20) DEFAULT 'pending',    -- pending, processing, completed, failed
    filters JSONB,                           -- Date range, etc
    file_path VARCHAR(500),
    file_size BIGINT,
    row_count INT,
    error TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_export_tenant ON export_jobs(tenant_id, created_at DESC);
