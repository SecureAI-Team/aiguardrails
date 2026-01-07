package rag

import (
	"errors"
	"regexp"
	"strings"

	"aiguardrails/internal/policy"
)

var (
	ErrNamespaceNotAllowed = errors.New("namespace not allowed")
	ErrQueryInjection      = errors.New("potential query injection detected")
	ErrDocumentAccess      = errors.New("document access denied")
)

// Document represents a RAG document result.
type Document struct {
	ID          string            `json:"id"`
	Namespace   string            `json:"namespace"`
	Content     string            `json:"content"`
	Score       float64           `json:"score"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Sensitivity string            `json:"sensitivity,omitempty"` // public, internal, confidential, secret
}

// Security enforces namespace isolation, query validation, and result filtering.
type Security struct {
	policy            policy.Engine
	injectionPatterns []*regexp.Regexp
	filters           []ResultFilter
}

// NewSecurity constructs a Security layer.
func NewSecurity(p policy.Engine) *Security {
	return &Security{
		policy:            p,
		injectionPatterns: defaultInjectionPatterns(),
		filters: []ResultFilter{
			&DLPFilter{},
			&SensitivityFilter{},
		},
	}
}

// ValidateNamespace ensures queries stay within allowed namespaces.
func (s *Security) ValidateNamespace(tenantID, namespace string) error {
	allowed := s.policy.AllowedNamespaces(tenantID)
	if len(allowed) == 0 {
		return nil // open by default if not configured
	}
	for _, ns := range allowed {
		if ns == namespace || ns == "*" {
			return nil
		}
		// Support wildcard prefix matching
		if strings.HasSuffix(ns, "*") && strings.HasPrefix(namespace, strings.TrimSuffix(ns, "*")) {
			return nil
		}
	}
	return ErrNamespaceNotAllowed
}

// ValidateQuery checks query for injection patterns.
func (s *Security) ValidateQuery(tenantID, query string) error {
	queryLower := strings.ToLower(query)
	for _, pattern := range s.injectionPatterns {
		if pattern.MatchString(queryLower) {
			return ErrQueryInjection
		}
	}
	return nil
}

// FilterResults applies all filters to results.
func (s *Security) FilterResults(tenantID string, docs []Document, userLevel string) []Document {
	filtered := docs
	for _, filter := range s.filters {
		filtered = filter.Filter(tenantID, filtered, userLevel)
	}
	return filtered
}

// RedactResult performs basic masking for secrets.
func (s *Security) RedactResult(result string) string {
	customTerms := s.policy.CustomTerms("")
	for _, term := range customTerms {
		if strings.Contains(strings.ToLower(result), strings.ToLower(term)) {
			return "[REDACTED: contains sensitive term]"
		}
	}
	if strings.Contains(strings.ToLower(result), "password") {
		return "[REDACTED]"
	}
	return result
}

// AddFilter adds a custom result filter.
func (s *Security) AddFilter(f ResultFilter) {
	s.filters = append(s.filters, f)
}

// defaultInjectionPatterns returns patterns to detect injection.
func defaultInjectionPatterns() []*regexp.Regexp {
	patterns := []string{
		`(?i)ignore\s+previous\s+instructions`,
		`(?i)forget\s+all\s+(previous|prior)`,
		`(?i)disregard\s+(above|previous)`,
		`(?i)system\s*:\s*you\s+are`,
		`(?i)\[INST\]`,
		`(?i)<\|im_start\|>`,
		`(?i)###\s*instruction`,
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if re, err := regexp.Compile(p); err == nil {
			compiled = append(compiled, re)
		}
	}
	return compiled
}

// ResultFilter interface for filtering results.
type ResultFilter interface {
	Filter(tenantID string, docs []Document, userLevel string) []Document
}

// DLPFilter removes documents containing sensitive data.
type DLPFilter struct{}

func (f *DLPFilter) Filter(tenantID string, docs []Document, userLevel string) []Document {
	result := make([]Document, 0, len(docs))
	for _, doc := range docs {
		if !containsSensitiveData(doc.Content) {
			result = append(result, doc)
		} else {
			// Redact instead of removing
			doc.Content = "[REDACTED: DLP detected]"
			result = append(result, doc)
		}
	}
	return result
}

func containsSensitiveData(content string) bool {
	lc := strings.ToLower(content)
	patterns := []string{
		"password", "secret", "api_key", "apikey",
		"access_token", "private_key", "credential",
	}
	for _, p := range patterns {
		if strings.Contains(lc, p) {
			return true
		}
	}
	return false
}

// SensitivityFilter filters by document sensitivity level.
type SensitivityFilter struct{}

var sensitivityLevels = map[string]int{
	"public":       1,
	"internal":     2,
	"confidential": 3,
	"secret":       4,
}

func (f *SensitivityFilter) Filter(tenantID string, docs []Document, userLevel string) []Document {
	maxLevel := sensitivityLevels[userLevel]
	if maxLevel == 0 {
		maxLevel = 1 // default to public only
	}
	result := make([]Document, 0, len(docs))
	for _, doc := range docs {
		docLevel := sensitivityLevels[doc.Sensitivity]
		if docLevel == 0 {
			docLevel = 1 // default to public
		}
		if docLevel <= maxLevel {
			result = append(result, doc)
		}
	}
	return result
}

// NamespaceFilter filters by allowed namespaces.
type NamespaceFilter struct {
	allowedNamespaces []string
}

func NewNamespaceFilter(allowed []string) *NamespaceFilter {
	return &NamespaceFilter{allowedNamespaces: allowed}
}

func (f *NamespaceFilter) Filter(tenantID string, docs []Document, userLevel string) []Document {
	if len(f.allowedNamespaces) == 0 {
		return docs
	}
	result := make([]Document, 0, len(docs))
	for _, doc := range docs {
		for _, ns := range f.allowedNamespaces {
			if ns == doc.Namespace || ns == "*" {
				result = append(result, doc)
				break
			}
		}
	}
	return result
}
