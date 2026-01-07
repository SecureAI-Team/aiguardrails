-- Policy versioning enhancements
ALTER TABLE policy_history ADD COLUMN IF NOT EXISTS version INT NOT NULL DEFAULT 1;
ALTER TABLE policy_history ADD COLUMN IF NOT EXISTS change_summary TEXT;
ALTER TABLE policy_history ADD COLUMN IF NOT EXISTS changed_by TEXT;

-- Index for version queries
CREATE INDEX IF NOT EXISTS idx_policy_history_version ON policy_history(policy_id, version DESC);
CREATE INDEX IF NOT EXISTS idx_policy_history_tenant ON policy_history(tenant_id, updated_at DESC);

-- Add version tracking to policies table
ALTER TABLE policies ADD COLUMN IF NOT EXISTS version INT NOT NULL DEFAULT 1;
ALTER TABLE policies ADD COLUMN IF NOT EXISTS change_summary TEXT;
ALTER TABLE policies ADD COLUMN IF NOT EXISTS changed_by TEXT;
