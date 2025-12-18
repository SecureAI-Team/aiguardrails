CREATE TABLE IF NOT EXISTS policy_rules (
    policy_id UUID NOT NULL REFERENCES policies(id) ON DELETE CASCADE,
    rule_id TEXT NOT NULL,
    PRIMARY KEY (policy_id, rule_id)
);

