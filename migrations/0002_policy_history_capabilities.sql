-- Add sensitive_terms to policies
ALTER TABLE policies ADD COLUMN IF NOT EXISTS sensitive_terms JSONB NOT NULL DEFAULT '[]'::jsonb;

-- Policy history
CREATE TABLE IF NOT EXISTS policy_history (
    id SERIAL PRIMARY KEY,
    policy_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    name TEXT NOT NULL,
    prompt_rules JSONB NOT NULL DEFAULT '[]'::jsonb,
    tool_allowlist JSONB NOT NULL DEFAULT '[]'::jsonb,
    rag_namespaces JSONB NOT NULL DEFAULT '[]'::jsonb,
    output_filters JSONB NOT NULL DEFAULT '[]'::jsonb,
    sensitive_terms JSONB NOT NULL DEFAULT '[]'::jsonb,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Capabilities registry
CREATE TABLE IF NOT EXISTS capabilities (
    id UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    tags JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

