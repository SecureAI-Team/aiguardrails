package promptfw

import (
	"strings"

	"aiguardrails/internal/policy"
	"aiguardrails/internal/types"
)

// Firewall wraps prompt evaluation for injection prevention.
type Firewall struct {
	policy policy.Engine
	llm    *policy.LLMDetector
	mode   string // "block" or "mark"
}

// NewFirewall constructs a Firewall.
func NewFirewall(p policy.Engine) *Firewall {
	return &Firewall{policy: p}
}

// WithLLM attaches an LLM detector.
func (f *Firewall) WithLLM(det *policy.LLMDetector, mode string) {
	f.llm = det
	if mode == "block" {
		f.mode = "block"
	} else {
		f.mode = "mark"
	}
}

// CheckPrompt runs prompt through guardrails: injection check + explicit keywords.
func (f *Firewall) CheckPrompt(tenantID, prompt string, keywords []string) types.GuardrailResult {
	if strings.Contains(strings.ToLower(prompt), "ignore previous instructions") {
		return types.GuardrailResult{Allowed: false, Reason: "prompt_injection_detected", Signals: []string{"ignore previous instructions"}}
	}

	lowPrompt := strings.ToLower(prompt)
	for _, kw := range keywords {
		if strings.Contains(lowPrompt, strings.ToLower(kw)) {
			return types.GuardrailResult{
				Allowed: false,
				Reason:  "keyword_block",
				Signals: []string{kw},
			}
		}
	}
	// Note: We no longer call policy.EvaluatePrompt here because aggregation happens upstream (server)
	return types.GuardrailResult{Allowed: true}
}

// FilterOutput applies DLP with built-in and custom terms.
func (f *Firewall) FilterOutput(tenantID, output string) types.GuardrailResult {
	custom := f.policy.CustomTerms(tenantID)
	dlp := policy.DetectDLP(output, custom)
	if dlp.Hit {
		return types.GuardrailResult{Allowed: false, Reason: dlp.Reason, Signals: dlp.Matches}
	}
	if f.llm != nil {
		res := f.llm.Check(output)
		if !res.Allowed {
			return res
		}
		if res.Reason == "llm_pending" && f.mode == "mark" {
			return types.GuardrailResult{Allowed: true, Reason: "llm_pending", Signals: res.Signals}
		}
		if f.mode == "block" && res.Reason == "llm_pending" {
			return types.GuardrailResult{Allowed: false, Reason: "llm_pending"}
		}
	}
	return types.GuardrailResult{Allowed: true}
}
