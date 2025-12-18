package rag

import (
	"errors"
	"strings"

	"aiguardrails/internal/policy"
)

// Security enforces namespace isolation and basic redaction.
type Security struct {
	policy policy.Engine
}

// NewSecurity constructs a Security layer.
func NewSecurity(p policy.Engine) *Security {
	return &Security{policy: p}
}

// ValidateNamespace ensures queries stay within allowed namespaces.
func (s *Security) ValidateNamespace(tenantID, namespace string) error {
	allowed := s.policy.AllowedNamespaces(tenantID)
	if len(allowed) == 0 {
		return nil // open by default if not configured
	}
	for _, ns := range allowed {
		if ns == namespace {
			return nil
		}
	}
	return errors.New("namespace not allowed")
}

// RedactResult performs basic masking for secrets.
func (s *Security) RedactResult(result string) string {
	if strings.Contains(strings.ToLower(result), "password") {
		return "[REDACTED]"
	}
	return result
}

