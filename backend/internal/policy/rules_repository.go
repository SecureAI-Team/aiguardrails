package policy

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Rule represents a compliance/industry rule entry.
type Rule struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Jurisdiction string  `json:"jurisdiction"` // e.g., EU
	Regulation  string   `json:"regulation"`   // e.g., EU AI Act, GDPR
	Vendor      string   `json:"vendor"`       // e.g., Siemens, ABB
	Product     string   `json:"product"`
	Severity    string   `json:"severity"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	References  []string `json:"references"` // URLs or clause references
	Description string   `json:"description"`
	Remediation string   `json:"remediation"`
}

// RuleAttachment links rule to policy.
type RuleAttachment struct {
	PolicyID string
	RuleID   string
}

// RulesRepository loads rules from json files.
type RulesRepository struct {
	rules []Rule
}

// NewRulesRepository loads all JSON files under dir.
func NewRulesRepository(dir string) (*RulesRepository, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var rules []Rule
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if filepath.Ext(f.Name()) != ".json" {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
		var rs []Rule
		if err := json.Unmarshal(b, &rs); err != nil {
			return nil, err
		}
		rules = append(rules, rs...)
	}
	return &RulesRepository{rules: rules}, nil
}

// List returns rules with optional filters.
func (r *RulesRepository) List(filter map[string]string) []Rule {
	var out []Rule
	for _, rule := range r.rules {
		if filter["jurisdiction"] != "" && filter["jurisdiction"] != rule.Jurisdiction {
			continue
		}
		if filter["regulation"] != "" && filter["regulation"] != rule.Regulation {
			continue
		}
		if filter["vendor"] != "" && filter["vendor"] != rule.Vendor {
			continue
		}
		if filter["product"] != "" && filter["product"] != rule.Product {
			continue
		}
		out = append(out, rule)
	}
	return out
}

// Get returns rule by ID.
func (r *RulesRepository) Get(id string) (*Rule, error) {
	for _, rule := range r.rules {
		if rule.ID == id {
			return &rule, nil
		}
	}
	return nil, errors.New("rule not found")
}

