package server

import (
	"aiguardrails/internal/opa"
	"aiguardrails/internal/rules"
	"fmt"
)

// syncOPARules merges base policies and dynamic rules, then reloads OPA.
func (s *Server) syncOPARules() {
	if s.opaEval == nil {
		return
	}

	// 1. Load base policies from disk
	// Adjust path if needed, assuming running from /app or backend root
	modules, err := opa.LoadModules("opa/policies")
	if err != nil {
		fmt.Printf("Warning: Failed to load base OPA policies: %v. Usage OPA features might fail.\n", err)
		modules = map[string]string{}
	}

	// 2. Load dynamic rules from store
	allRules, err := s.ruleStore.List()
	if err == nil {
		for _, r := range allRules {
			if r.Type == rules.RuleTypeOPA && r.Content != "" {
				// Create a unique module name
				name := fmt.Sprintf("dynamic_%s.rego", r.ID)
				modules[name] = r.Content
			}
		}
	}

	// 3. Reload OPA
	if err := s.opaEval.ReloadFromContent(modules); err != nil {
		fmt.Printf("Error reloading OPA modules: %v\n", err)
	} else {
		fmt.Printf("OPA Engine synced: %d active modules\n", len(modules))
	}
}
