package policy

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"aiguardrails/internal/types"
)

// PGEngine stores policies in Postgres.
type PGEngine struct {
	db *sql.DB
}

// NewPGEngine constructs PGEngine.
func NewPGEngine(db *sql.DB) *PGEngine {
	return &PGEngine{db: db}
}

func (e *PGEngine) CreatePolicy(p types.Policy) (types.Policy, error) {
	if p.TenantID == "" || p.Name == "" {
		return p, errors.New("tenantID and name required")
	}
	p.ID = uuid.NewString()
	p.LastModifiedAt = time.Now().UTC()
	if err := e.insert(p); err != nil {
		return p, err
	}
	_ = e.insertHistory(p)
	return p, nil
}

// UpdatePolicy updates an existing policy and records history.
func (e *PGEngine) UpdatePolicy(p types.Policy) (types.Policy, error) {
	if p.TenantID == "" || p.ID == "" {
		return p, errors.New("tenantID and id required")
	}
	p.LastModifiedAt = time.Now().UTC()
	pr, _ := json.Marshal(p.PromptRules)
	tl, _ := json.Marshal(p.ToolAllowList)
	rn, _ := json.Marshal(p.RAGNamespaces)
	of, _ := json.Marshal(p.OutputFilters)
	st, _ := json.Marshal(p.SensitiveTerms)
	_, err := e.db.Exec(`UPDATE policies SET name=$1, prompt_rules=$2, tool_allowlist=$3, rag_namespaces=$4, output_filters=$5, sensitive_terms=$6, updated_at=$7 WHERE id=$8 AND tenant_id=$9`,
		p.Name, pr, tl, rn, of, st, p.LastModifiedAt, p.ID, p.TenantID)
	if err != nil {
		return p, err
	}
	_ = e.insertHistory(p)
	return p, nil
}

func (e *PGEngine) insert(p types.Policy) error {
	pr, _ := json.Marshal(p.PromptRules)
	tl, _ := json.Marshal(p.ToolAllowList)
	rn, _ := json.Marshal(p.RAGNamespaces)
	of, _ := json.Marshal(p.OutputFilters)
	st, _ := json.Marshal(p.SensitiveTerms)
	_, err := e.db.Exec(`INSERT INTO policies (id, tenant_id, name, prompt_rules, tool_allowlist, rag_namespaces, output_filters, sensitive_terms, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		p.ID, p.TenantID, p.Name, pr, tl, rn, of, st, p.LastModifiedAt)
	return err
}

func (e *PGEngine) ListPolicies(tenantID string) ([]types.Policy, error) {
	rows, err := e.db.Query(`SELECT id, name, prompt_rules, tool_allowlist, rag_namespaces, output_filters, sensitive_terms, updated_at FROM policies WHERE tenant_id=$1`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []types.Policy
	for rows.Next() {
		var p types.Policy
		p.TenantID = tenantID
		var pr, tl, rn, of, st []byte
		if err := rows.Scan(&p.ID, &p.Name, &pr, &tl, &rn, &of, &st, &p.LastModifiedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(pr, &p.PromptRules)
		_ = json.Unmarshal(tl, &p.ToolAllowList)
		_ = json.Unmarshal(rn, &p.RAGNamespaces)
		_ = json.Unmarshal(of, &p.OutputFilters)
		_ = json.Unmarshal(st, &p.SensitiveTerms)
		out = append(out, p)
	}
	return out, nil
}

// GetPolicy returns a policy by ID.
func (e *PGEngine) GetPolicy(tenantID, policyID string) (*types.Policy, error) {
	row := e.db.QueryRow(`SELECT id, name, prompt_rules, tool_allowlist, rag_namespaces, output_filters, sensitive_terms, updated_at FROM policies WHERE tenant_id=$1 AND id=$2`, tenantID, policyID)
	var p types.Policy
	p.TenantID = tenantID
	var pr, tl, rn, of, st []byte
	if err := row.Scan(&p.ID, &p.Name, &pr, &tl, &rn, &of, &st, &p.LastModifiedAt); err != nil {
		return nil, err
	}
	_ = json.Unmarshal(pr, &p.PromptRules)
	_ = json.Unmarshal(tl, &p.ToolAllowList)
	_ = json.Unmarshal(rn, &p.RAGNamespaces)
	_ = json.Unmarshal(of, &p.OutputFilters)
	_ = json.Unmarshal(st, &p.SensitiveTerms)
	return &p, nil
}

// ListHistory returns policy history for a tenant.
func (e *PGEngine) ListHistory(tenantID string, limit int) ([]types.Policy, error) {
	rows, err := e.db.Query(`SELECT policy_id, name, prompt_rules, tool_allowlist, rag_namespaces, output_filters, sensitive_terms, updated_at FROM policy_history WHERE tenant_id=$1 ORDER BY updated_at DESC LIMIT $2`, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []types.Policy
	for rows.Next() {
		var p types.Policy
		p.TenantID = tenantID
		var pr, tl, rn, of, st []byte
		if err := rows.Scan(&p.ID, &p.Name, &pr, &tl, &rn, &of, &st, &p.LastModifiedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(pr, &p.PromptRules)
		_ = json.Unmarshal(tl, &p.ToolAllowList)
		_ = json.Unmarshal(rn, &p.RAGNamespaces)
		_ = json.Unmarshal(of, &p.OutputFilters)
		_ = json.Unmarshal(st, &p.SensitiveTerms)
		out = append(out, p)
	}
	return out, nil
}

func (e *PGEngine) insertHistory(p types.Policy) error {
	pr, _ := json.Marshal(p.PromptRules)
	tl, _ := json.Marshal(p.ToolAllowList)
	rn, _ := json.Marshal(p.RAGNamespaces)
	of, _ := json.Marshal(p.OutputFilters)
	st, _ := json.Marshal(p.SensitiveTerms)
	_, err := e.db.Exec(`INSERT INTO policy_history (policy_id, tenant_id, name, prompt_rules, tool_allowlist, rag_namespaces, output_filters, sensitive_terms, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		p.ID, p.TenantID, p.Name, pr, tl, rn, of, st, p.LastModifiedAt)
	return err
}

func (e *PGEngine) EvaluatePrompt(tenantID, prompt string) types.GuardrailResult {
	policies, err := e.ListPolicies(tenantID)
	if err != nil {
		return types.GuardrailResult{Allowed: false, Reason: "policy_load_error", Signals: []string{err.Error()}}
	}
	for _, p := range policies {
		for _, rule := range p.PromptRules {
			if strings.Contains(strings.ToLower(prompt), strings.ToLower(rule)) {
				return types.GuardrailResult{Allowed: false, Reason: "blocked_by_prompt_rule", Signals: []string{rule}}
			}
		}
	}
	return types.GuardrailResult{Allowed: true}
}

func (e *PGEngine) AllowTool(tenantID, tool string) bool {
	policies, err := e.ListPolicies(tenantID)
	if err != nil {
		return false
	}
	allowlistExists := false
	for _, p := range policies {
		if len(p.ToolAllowList) > 0 {
			allowlistExists = true
			for _, allowed := range p.ToolAllowList {
				if allowed == tool {
					return true
				}
			}
		}
	}
	if allowlistExists {
		return false
	}
	return true
}

func (e *PGEngine) AllowedNamespaces(tenantID string) []string {
	policies, err := e.ListPolicies(tenantID)
	if err != nil {
		return nil
	}
	set := map[string]struct{}{}
	for _, p := range policies {
		for _, ns := range p.RAGNamespaces {
			set[ns] = struct{}{}
		}
	}
	out := make([]string, 0, len(set))
	for ns := range set {
		out = append(out, ns)
	}
	return out
}

// CustomTerms aggregates sensitive terms.
func (e *PGEngine) CustomTerms(tenantID string) []string {
	policies, err := e.ListPolicies(tenantID)
	if err != nil {
		return nil
	}
	set := map[string]struct{}{}
	for _, p := range policies {
		for _, t := range p.SensitiveTerms {
			set[t] = struct{}{}
		}
	}
	out := make([]string, 0, len(set))
	for t := range set {
		out = append(out, t)
	}
	return out
}

// Ensure interface compliance.
var _ Engine = (*PGEngine)(nil)

