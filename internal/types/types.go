package types

import "time"

// Tenant represents a logical organization using the SaaS.
type Tenant struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

// App represents an application key pair issued to a tenant.
type App struct {
	ID         string
	TenantID   string
	Name       string
	APIKey     string
	APISecret  string
	QuotaPerHr int64
	CreatedAt  time.Time
	Revoked    bool
}

// Policy defines guardrails applied to requests.
type Policy struct {
	ID             string
	TenantID       string
	Name           string
	PromptRules    []string // e.g., regex or keywords
	ToolAllowList  []string
	RAGNamespaces  []string
	OutputFilters  []string // PII, secrets, toxicity tags
	SensitiveTerms []string // custom sensitive words
	LastModifiedAt time.Time
}

// UsageRecord tracks app usage for metering.
type UsageRecord struct {
	AppID     string
	Timestamp time.Time
	Count     int64
}

// GuardrailResult captures prompt firewall decisions.
type GuardrailResult struct {
	Allowed bool
	Reason  string
	Signals []string
}

