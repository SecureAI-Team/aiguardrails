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

// GetHistoryVersion returns a specific version of a policy.
func (e *PGEngine) GetHistoryVersion(tenantID, policyID string, version int) (*types.Policy, error) {
	row := e.db.QueryRow(`SELECT policy_id, name, prompt_rules, tool_allowlist, rag_namespaces, output_filters, sensitive_terms, updated_at, version, COALESCE(change_summary, ''), COALESCE(changed_by, '') 
		FROM policy_history WHERE tenant_id=$1 AND policy_id=$2 AND version=$3`, tenantID, policyID, version)
	var p types.Policy
	p.TenantID = tenantID
	var pr, tl, rn, of, st []byte
	if err := row.Scan(&p.ID, &p.Name, &pr, &tl, &rn, &of, &st, &p.LastModifiedAt, &p.Version, &p.ChangeSummary, &p.ChangedBy); err != nil {
		return nil, err
	}
	_ = json.Unmarshal(pr, &p.PromptRules)
	_ = json.Unmarshal(tl, &p.ToolAllowList)
	_ = json.Unmarshal(rn, &p.RAGNamespaces)
	_ = json.Unmarshal(of, &p.OutputFilters)
	_ = json.Unmarshal(st, &p.SensitiveTerms)
	return &p, nil
}

// ListPolicyVersions returns all versions of a specific policy.
func (e *PGEngine) ListPolicyVersions(tenantID, policyID string) ([]types.PolicyVersion, error) {
	rows, err := e.db.Query(`SELECT version, updated_at, COALESCE(change_summary, ''), COALESCE(changed_by, '') 
		FROM policy_history WHERE tenant_id=$1 AND policy_id=$2 ORDER BY version DESC`, tenantID, policyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var versions []types.PolicyVersion
	for rows.Next() {
		var v types.PolicyVersion
		v.PolicyID = policyID
		if err := rows.Scan(&v.Version, &v.UpdatedAt, &v.ChangeSummary, &v.ChangedBy); err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}
	return versions, nil
}

// RollbackPolicy restores a policy to a previous version.
func (e *PGEngine) RollbackPolicy(tenantID, policyID string, version int, changedBy string) (types.Policy, error) {
	// Get the historical version
	historyPolicy, err := e.GetHistoryVersion(tenantID, policyID, version)
	if err != nil {
		return types.Policy{}, err
	}

	// Get current version number
	var currentVersion int
	row := e.db.QueryRow(`SELECT COALESCE(version, 1) FROM policies WHERE tenant_id=$1 AND id=$2`, tenantID, policyID)
	_ = row.Scan(&currentVersion)

	// Update policy with historical values
	newVersion := currentVersion + 1
	historyPolicy.Version = newVersion
	historyPolicy.ChangeSummary = "Rollback to version " + string(rune(version+'0'))
	historyPolicy.ChangedBy = changedBy
	historyPolicy.LastModifiedAt = time.Now().UTC()

	pr, _ := json.Marshal(historyPolicy.PromptRules)
	tl, _ := json.Marshal(historyPolicy.ToolAllowList)
	rn, _ := json.Marshal(historyPolicy.RAGNamespaces)
	of, _ := json.Marshal(historyPolicy.OutputFilters)
	st, _ := json.Marshal(historyPolicy.SensitiveTerms)

	_, err = e.db.Exec(`UPDATE policies SET name=$1, prompt_rules=$2, tool_allowlist=$3, rag_namespaces=$4, output_filters=$5, sensitive_terms=$6, updated_at=$7, version=$8, change_summary=$9, changed_by=$10 WHERE id=$11 AND tenant_id=$12`,
		historyPolicy.Name, pr, tl, rn, of, st, historyPolicy.LastModifiedAt, newVersion, historyPolicy.ChangeSummary, changedBy, policyID, tenantID)
	if err != nil {
		return types.Policy{}, err
	}

	// Record in history
	_ = e.insertHistoryWithVersion(*historyPolicy, newVersion, historyPolicy.ChangeSummary, changedBy)
	return *historyPolicy, nil
}

// insertHistoryWithVersion inserts history with version info.
func (e *PGEngine) insertHistoryWithVersion(p types.Policy, version int, summary, changedBy string) error {
	pr, _ := json.Marshal(p.PromptRules)
	tl, _ := json.Marshal(p.ToolAllowList)
	rn, _ := json.Marshal(p.RAGNamespaces)
	of, _ := json.Marshal(p.OutputFilters)
	st, _ := json.Marshal(p.SensitiveTerms)
	_, err := e.db.Exec(`INSERT INTO policy_history (policy_id, tenant_id, name, prompt_rules, tool_allowlist, rag_namespaces, output_filters, sensitive_terms, updated_at, version, change_summary, changed_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		p.ID, p.TenantID, p.Name, pr, tl, rn, of, st, p.LastModifiedAt, version, summary, changedBy)
	return err
}

// CompareVersions returns diff between two versions.
func (e *PGEngine) CompareVersions(tenantID, policyID string, v1, v2 int) (*types.PolicyDiff, error) {
	p1, err := e.GetHistoryVersion(tenantID, policyID, v1)
	if err != nil {
		return nil, err
	}
	p2, err := e.GetHistoryVersion(tenantID, policyID, v2)
	if err != nil {
		return nil, err
	}

	diff := &types.PolicyDiff{
		PolicyID: policyID,
		FromVer:  v1,
		ToVer:    v2,
		Changes:  make(map[string]types.FieldChange),
	}

	if p1.Name != p2.Name {
		diff.Changes["name"] = types.FieldChange{Field: "name", OldValue: p1.Name, NewValue: p2.Name}
	}
	if !slicesEqual(p1.PromptRules, p2.PromptRules) {
		diff.Changes["prompt_rules"] = types.FieldChange{Field: "prompt_rules", OldValue: p1.PromptRules, NewValue: p2.PromptRules}
	}
	if !slicesEqual(p1.ToolAllowList, p2.ToolAllowList) {
		diff.Changes["tool_allowlist"] = types.FieldChange{Field: "tool_allowlist", OldValue: p1.ToolAllowList, NewValue: p2.ToolAllowList}
	}
	if !slicesEqual(p1.RAGNamespaces, p2.RAGNamespaces) {
		diff.Changes["rag_namespaces"] = types.FieldChange{Field: "rag_namespaces", OldValue: p1.RAGNamespaces, NewValue: p2.RAGNamespaces}
	}
	if !slicesEqual(p1.OutputFilters, p2.OutputFilters) {
		diff.Changes["output_filters"] = types.FieldChange{Field: "output_filters", OldValue: p1.OutputFilters, NewValue: p2.OutputFilters}
	}
	if !slicesEqual(p1.SensitiveTerms, p2.SensitiveTerms) {
		diff.Changes["sensitive_terms"] = types.FieldChange{Field: "sensitive_terms", OldValue: p1.SensitiveTerms, NewValue: p2.SensitiveTerms}
	}

	return diff, nil
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Ensure interface compliance.
var _ Engine = (*PGEngine)(nil)
