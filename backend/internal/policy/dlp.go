package policy

import (
	"regexp"
	"strings"
)

var (
	piiPatterns = []*regexp.Regexp{
		regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`),            // SSN-like
		regexp.MustCompile(`\b4[0-9]{12}(?:[0-9]{3})?\b`),      // Visa-like
		regexp.MustCompile(`\b[13][a-km-zA-HJ-NP-Z1-9]{25,34}\b`), // crypto-ish
	}
)

// DLPResult captures detection outcome.
type DLPResult struct {
	Hit     bool
	Reason  string
	Matches []string
}

// DetectDLP combines regex, dictionary, and custom terms.
func DetectDLP(text string, customTerms []string) DLPResult {
	lower := strings.ToLower(text)
	matches := []string{}

	// Regex patterns
	for _, re := range piiPatterns {
		if loc := re.FindString(text); loc != "" {
			matches = append(matches, loc)
		}
	}

	// Built-in dictionary
	keywords := []string{
		"password", "ssn", "secret", "apikey", "credit card",
		"passport", "id card", "social security", "iban", "swift",
		"private key", "api key", "token", "secret key", "credential",
		"confidential", "proprietary",
	}
	for _, k := range keywords {
		if strings.Contains(lower, k) {
			matches = append(matches, k)
		}
	}

	// Custom terms
	for _, term := range customTerms {
		if term == "" {
			continue
		}
		if strings.Contains(lower, strings.ToLower(term)) {
			matches = append(matches, term)
		}
	}

	if len(matches) > 0 {
		return DLPResult{Hit: true, Reason: "dlp_match", Matches: matches}
	}

	// Stub for LLM-based detector
	// In production, call an LLM classifier; here we just return no-hit.
	return DLPResult{Hit: false}
}


