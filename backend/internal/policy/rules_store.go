package policy

import (
	"database/sql"
)

// RuleStore persists policy-rule attachments in Postgres.
type RuleStore struct {
	db *sql.DB
}

// NewRuleStore constructs RuleStore.
func NewRuleStore(db *sql.DB) *RuleStore {
	return &RuleStore{db: db}
}

// Attach links a rule to a policy.
func (s *RuleStore) Attach(policyID, ruleID string) error {
	_, err := s.db.Exec(`INSERT INTO policy_rules (policy_id, rule_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, policyID, ruleID)
	return err
}

// ListByPolicy returns rule IDs attached to a policy.
func (s *RuleStore) ListByPolicy(policyID string) ([]string, error) {
	rows, err := s.db.Query(`SELECT rule_id FROM policy_rules WHERE policy_id=$1`, policyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, nil
}

