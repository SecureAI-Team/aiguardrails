package policy

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"aiguardrails/internal/types"
)

// Engine evaluates guardrail policies.
type Engine interface {
	CreatePolicy(p types.Policy) (types.Policy, error)
	UpdatePolicy(p types.Policy) (types.Policy, error)
	DeletePolicy(tenantID, policyID string) error
	ListPolicies(tenantID string) ([]types.Policy, error)
	ListHistory(tenantID string, limit int) ([]types.Policy, error)
	GetPolicy(tenantID, policyID string) (*types.Policy, error)
	EvaluatePrompt(tenantID, prompt string) types.GuardrailResult
	AllowTool(tenantID, tool string) bool
	AllowedNamespaces(tenantID string) []string
	CustomTerms(tenantID string) []string
}

// MemoryEngine stores policies in memory.
type MemoryEngine struct {
	mu       sync.RWMutex
	policies map[string][]types.Policy // tenantID -> policies
}

// NewMemoryEngine builds a MemoryEngine.
func NewMemoryEngine() *MemoryEngine {
	return &MemoryEngine{
		policies: map[string][]types.Policy{},
	}
}

// CreatePolicy adds a policy for the tenant.
func (e *MemoryEngine) CreatePolicy(p types.Policy) (types.Policy, error) {
	p.ID = uuid.NewString()
	p.LastModifiedAt = time.Now().UTC()
	e.mu.Lock()
	defer e.mu.Unlock()
	e.policies[p.TenantID] = append(e.policies[p.TenantID], p)
	return p, nil
}

// UpdatePolicy replaces a policy by ID.
func (e *MemoryEngine) UpdatePolicy(p types.Policy) (types.Policy, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	list := e.policies[p.TenantID]
	for i, policy := range list {
		if policy.ID == p.ID {
			p.LastModifiedAt = time.Now().UTC()
			list[i] = p
			e.policies[p.TenantID] = list
			return p, nil
		}
	}
	return p, errors.New("policy not found")
}

// DeletePolicy removes a policy by ID.
func (e *MemoryEngine) DeletePolicy(tenantID, policyID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	list := e.policies[tenantID]
	for i, policy := range list {
		if policy.ID == policyID {
			e.policies[tenantID] = append(list[:i], list[i+1:]...)
			return nil
		}
	}
	return errors.New("policy not found")
}

// GetPolicy returns a policy by ID.
func (e *MemoryEngine) GetPolicy(tenantID, policyID string) (*types.Policy, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, p := range e.policies[tenantID] {
		if p.ID == policyID {
			copy := p
			return &copy, nil
		}
	}
	return nil, errors.New("policy not found")
}

// ListPolicies returns policies for a tenant.
func (e *MemoryEngine) ListPolicies(tenantID string) ([]types.Policy, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return append([]types.Policy{}, e.policies[tenantID]...), nil
}

// EvaluatePrompt checks prompt against simple rules.
func (e *MemoryEngine) EvaluatePrompt(tenantID, prompt string) types.GuardrailResult {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, p := range e.policies[tenantID] {
		for _, rule := range p.PromptRules {
			if strings.Contains(strings.ToLower(prompt), strings.ToLower(rule)) {
				return types.GuardrailResult{
					Allowed: false,
					Reason:  "blocked_by_prompt_rule",
					Signals: []string{rule},
				}
			}
		}
	}
	return types.GuardrailResult{Allowed: true}
}

// AllowTool determines if a tool is permitted.
func (e *MemoryEngine) AllowTool(tenantID, tool string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, p := range e.policies[tenantID] {
		if len(p.ToolAllowList) == 0 {
			continue
		}
		for _, allowed := range p.ToolAllowList {
			if allowed == tool {
				return true
			}
		}
	}
	// If allowlist exists but did not match, deny; if none, allow.
	for _, p := range e.policies[tenantID] {
		if len(p.ToolAllowList) > 0 {
			return false
		}
	}
	return true
}

// AllowedNamespaces lists RAG namespaces permitted by policy.
func (e *MemoryEngine) AllowedNamespaces(tenantID string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	set := map[string]struct{}{}
	for _, p := range e.policies[tenantID] {
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
func (e *MemoryEngine) CustomTerms(tenantID string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	set := map[string]struct{}{}
	for _, p := range e.policies[tenantID] {
		for _, term := range p.SensitiveTerms {
			set[term] = struct{}{}
		}
	}
	out := make([]string, 0, len(set))
	for t := range set {
		out = append(out, t)
	}
	return out
}

// ListHistory returns empty for memory engine.
func (e *MemoryEngine) ListHistory(tenantID string, limit int) ([]types.Policy, error) {
	return e.ListPolicies(tenantID)
}
