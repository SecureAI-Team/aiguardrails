package rules

import (
	"encoding/json"
	"os"
	"time"
)

type ValidationRule struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	// Minimal fields we care about for seed
}

// Convert JSON seed format to Rule struct
func LoadFromJSON(path string, store Store) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var seedRules []ValidationRule
	if err := json.Unmarshal(bytes, &seedRules); err != nil {
		return err
	}

	for _, sr := range seedRules {
		// Map existing seed data to Rule
		// Defaulting to OPA type for vendor rules
		r := Rule{
			ID:          sr.ID,
			Name:        sr.Name,
			Description: sr.Description,
			Type:        RuleTypeOPA, // Current seeds are OPA rules
			Content:     "",          // OPA content is currently in file system (.rego), we rely on ID linkage
			Severity:    sr.Severity,
			Category:    sr.Category,
			Tags:        sr.Tags,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsSystem:    true,
		}
		_ = store.Add(r)
	}
	return nil
}
