package agent

import (
	"errors"
	"time"

	"aiguardrails/internal/policy"
	"aiguardrails/internal/promptfw"
	"aiguardrails/internal/types"
)

// Gateway orchestrates agent interactions with guardrails.
type Gateway struct {
	policy   policy.Engine
	firewall *promptfw.Firewall
}

// NewGateway constructs a Gateway.
func NewGateway(p policy.Engine, fw *promptfw.Firewall) *Gateway {
	return &Gateway{policy: p, firewall: fw}
}

// PlanAndAct is a stubbed planning loop that enforces tool allowlists.
func (g *Gateway) PlanAndAct(tenantID, prompt string, proposedTools []string) (types.GuardrailResult, error) {
	check := g.firewall.CheckPrompt(tenantID, prompt)
	if !check.Allowed {
		return check, errors.New("prompt rejected")
	}
	for _, t := range proposedTools {
		if !g.policy.AllowTool(tenantID, t) {
			return types.GuardrailResult{Allowed: false, Reason: "tool_not_allowed", Signals: []string{t}}, errors.New("tool blocked")
		}
	}
	// Stub: simulate execution time budget.
	time.Sleep(10 * time.Millisecond)
	return types.GuardrailResult{Allowed: true}, nil
}

