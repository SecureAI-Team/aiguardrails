package rules

import (
	"time"
)

type RuleType string

const (
	RuleTypeOPA     RuleType = "opa"
	RuleTypeLLM     RuleType = "llm"
	RuleTypeKeyword RuleType = "keyword" // Legacy support
)

// Rule represents a guardrail definition.
type Rule struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        RuleType  `json:"type"`
	Content     string    `json:"content"`  // Rego code, Prompt Template, or Keywords
	Severity    string    `json:"severity"` // low, medium, high
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsSystem    bool      `json:"is_system"` // If true, cannot be deleted
}

// Store defines persistence for rules.
type Store interface {
	Add(rule Rule) error
	Get(id string) (*Rule, error)
	List() ([]Rule, error)
	Delete(id string) error
	Update(rule Rule) error
}
