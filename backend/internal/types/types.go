package types

import "time"

// Tenant represents a logical organization using the SaaS.
type Tenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// App represents an application key pair issued to a tenant.
type App struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenant_id"`
	Name       string    `json:"name"`
	APIKey     string    `json:"api_key,omitempty"`
	APISecret  string    `json:"api_secret,omitempty"`
	QuotaPerHr int64     `json:"quota_per_hr"`
	CreatedAt  time.Time `json:"created_at"`
	Revoked    bool      `json:"is_revoked"`
}

// Policy defines guardrails applied to requests.
type Policy struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenant_id"`
	Name           string    `json:"name"`
	PromptRules    []string  `json:"prompt_rules"` // e.g., regex or keywords
	ToolAllowList  []string  `json:"tool_allowlist"`
	RAGNamespaces  []string  `json:"rag_namespaces"`
	OutputFilters  []string  `json:"output_filters"`  // PII, secrets, toxicity tags
	SensitiveTerms []string  `json:"sensitive_terms"` // custom sensitive words
	LastModifiedAt time.Time `json:"last_modified_at"`
	Version        int       `json:"version"`        // version number
	ChangeSummary  string    `json:"change_summary"` // description of changes
	ChangedBy      string    `json:"changed_by"`     // user who made the change
}

// PolicyVersion represents a version entry in history.
type PolicyVersion struct {
	PolicyID      string    `json:"policy_id"`
	Version       int       `json:"version"`
	UpdatedAt     time.Time `json:"updated_at"`
	ChangeSummary string    `json:"change_summary"`
	ChangedBy     string    `json:"changed_by"`
}

// PolicyDiff represents differences between two policy versions.
type PolicyDiff struct {
	PolicyID string                 `json:"policy_id"`
	FromVer  int                    `json:"from_version"`
	ToVer    int                    `json:"to_version"`
	Changes  map[string]FieldChange `json:"changes"`
}

// FieldChange represents a change to a single field.
type FieldChange struct {
	Field    string      `json:"field"`
	OldValue interface{} `json:"old_value"`
	NewValue interface{} `json:"new_value"`
}

// UsageRecord tracks app usage for metering.
type UsageRecord struct {
	AppID     string    `json:"app_id"`
	Timestamp time.Time `json:"timestamp"`
	Count     int64     `json:"count"`
}

// GuardrailResult captures prompt firewall decisions.
type GuardrailResult struct {
	Allowed bool     `json:"allowed"`
	Reason  string   `json:"reason"`
	Signals []string `json:"signals"`
}
