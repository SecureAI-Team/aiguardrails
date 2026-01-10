package agent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"aiguardrails/internal/policy"
	"aiguardrails/internal/promptfw"
	"aiguardrails/internal/types"
)

// Step represents one iteration in the ReAct loop.
type Step struct {
	Iteration   int           `json:"iteration"`
	Thought     string        `json:"thought"`
	Action      string        `json:"action"`
	ActionInput interface{}   `json:"action_input,omitempty"`
	Observation string        `json:"observation"`
	Duration    time.Duration `json:"duration_ms"`
	Blocked     bool          `json:"blocked,omitempty"`
	BlockReason string        `json:"block_reason,omitempty"`
}

// PlanRequest is input to PlanAndAct.
type PlanRequest struct {
	TenantID      string            `json:"tenant_id"`
	Prompt        string            `json:"prompt"`
	Tools         []string          `json:"tools"`
	MaxIterations int               `json:"max_iterations,omitempty"`
	Timeout       time.Duration     `json:"timeout,omitempty"`
	Context       map[string]string `json:"context,omitempty"`
}

// PlanResponse is the result of PlanAndAct.
type PlanResponse struct {
	Allowed     bool          `json:"allowed"`
	Reason      string        `json:"reason,omitempty"`
	Steps       []Step        `json:"steps"`
	FinalResult interface{}   `json:"final_result,omitempty"`
	TotalTime   time.Duration `json:"total_time_ms"`
	Signals     []string      `json:"signals,omitempty"`
}

// Metrics tracks agent execution stats.
type Metrics struct {
	mu              sync.Mutex
	TotalExecutions int64
	TotalBlocked    int64
	AvgIterations   float64
}

func (m *Metrics) record(blocked bool, iterations int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalExecutions++
	if blocked {
		m.TotalBlocked++
	}
	m.AvgIterations = (m.AvgIterations*float64(m.TotalExecutions-1) + float64(iterations)) / float64(m.TotalExecutions)
}

// Gateway orchestrates agent interactions with guardrails.
type Gateway struct {
	policy         policy.Engine
	firewall       *promptfw.Firewall
	sandbox        *Sandbox
	registry       *Registry
	metrics        *Metrics
	defaultMaxIter int
	defaultTimeout time.Duration
}

// GatewayOption configures Gateway.
type GatewayOption func(*Gateway)

// WithSandbox sets custom sandbox.
func WithSandbox(s *Sandbox) GatewayOption {
	return func(g *Gateway) { g.sandbox = s }
}

// WithRegistry sets tool registry.
func WithRegistry(r *Registry) GatewayOption {
	return func(g *Gateway) { g.registry = r }
}

// WithDefaults sets default max iterations and timeout.
func WithDefaults(maxIter int, timeout time.Duration) GatewayOption {
	return func(g *Gateway) {
		g.defaultMaxIter = maxIter
		g.defaultTimeout = timeout
	}
}

// NewGateway constructs a Gateway with options.
func NewGateway(p policy.Engine, fw *promptfw.Firewall, opts ...GatewayOption) *Gateway {
	g := &Gateway{
		policy:         p,
		firewall:       fw,
		sandbox:        NewSandbox(30*time.Second, 128*1024*1024, nil),
		registry:       NewRegistry(),
		metrics:        &Metrics{},
		defaultMaxIter: 10,
		defaultTimeout: 60 * time.Second,
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// GetMetrics returns current metrics.
func (g *Gateway) GetMetrics() Metrics {
	g.metrics.mu.Lock()
	defer g.metrics.mu.Unlock()
	return *g.metrics
}

// PlanAndAct implements ReAct loop with guardrail enforcement.
func (g *Gateway) PlanAndAct(ctx context.Context, req PlanRequest) (*PlanResponse, error) {
	start := time.Now()
	resp := &PlanResponse{Steps: []Step{}}

	// Set defaults
	maxIter := req.MaxIterations
	if maxIter <= 0 {
		maxIter = g.defaultMaxIter
	}
	timeout := req.Timeout
	if timeout <= 0 {
		timeout = g.defaultTimeout
	}

	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Phase 1: Check prompt
	check := g.firewall.CheckPrompt(req.TenantID, req.Prompt, []string{})
	if !check.Allowed {
		resp.Allowed = false
		resp.Reason = check.Reason
		resp.Signals = check.Signals
		resp.TotalTime = time.Since(start)
		g.metrics.record(true, 0)
		return resp, errors.New("prompt rejected")
	}

	// Phase 2: Validate all proposed tools
	for _, tool := range req.Tools {
		if !g.policy.AllowTool(req.TenantID, tool) {
			resp.Allowed = false
			resp.Reason = "tool_not_allowed"
			resp.Signals = []string{tool}
			resp.TotalTime = time.Since(start)
			g.metrics.record(true, 0)
			return resp, errors.New("tool blocked by policy")
		}
		if err := g.sandbox.ValidateTool(tool); err != nil {
			resp.Allowed = false
			resp.Reason = "tool_sandbox_rejected"
			resp.Signals = []string{tool, err.Error()}
			resp.TotalTime = time.Since(start)
			g.metrics.record(true, 0)
			return resp, err
		}
	}

	// Phase 3: ReAct loop
	for i := 0; i < maxIter; i++ {
		select {
		case <-ctx.Done():
			resp.Allowed = false
			resp.Reason = "timeout"
			resp.TotalTime = time.Since(start)
			g.metrics.record(true, i)
			return resp, ctx.Err()
		default:
		}

		iterStart := time.Now()
		step := Step{Iteration: i + 1}

		// Simulate thought generation (in real implementation, call LLM)
		step.Thought = fmt.Sprintf("Analyzing step %d for prompt: %s", i+1, truncate(req.Prompt, 50))

		// Determine action (stub: use first available tool or finish)
		if i < len(req.Tools) {
			step.Action = req.Tools[i]
			step.ActionInput = map[string]interface{}{"iteration": i}

			// Execute in sandbox
			result, err := g.sandbox.Execute(ctx, step.Action, map[string]interface{}{
				"prompt":    req.Prompt,
				"iteration": i,
			})
			if err != nil {
				step.Observation = fmt.Sprintf("Error: %v", err)
				step.Blocked = true
				step.BlockReason = err.Error()
			} else {
				step.Observation = fmt.Sprintf("Result: %v", result)
			}
		} else {
			// No more tools, finish
			step.Action = "finish"
			step.Observation = "Task completed"
		}

		step.Duration = time.Since(iterStart)
		resp.Steps = append(resp.Steps, step)

		// Check if blocked
		if step.Blocked {
			resp.Allowed = false
			resp.Reason = step.BlockReason
			resp.TotalTime = time.Since(start)
			g.metrics.record(true, i+1)
			return resp, errors.New(step.BlockReason)
		}

		// Check if finished
		if step.Action == "finish" {
			break
		}

		// FilterOutput through output filter
		filterResult := g.firewall.FilterOutput(req.TenantID, step.Observation, nil)
		if !filterResult.Allowed {
			resp.Allowed = false
			resp.Reason = filterResult.Reason
			resp.Signals = filterResult.Signals
			resp.TotalTime = time.Since(start)
			g.metrics.record(true, i+1)
			return resp, errors.New("output filtered")
		}
	}

	resp.Allowed = true
	resp.FinalResult = "Agent execution completed successfully"
	resp.TotalTime = time.Since(start)
	g.metrics.record(false, len(resp.Steps))
	return resp, nil
}

// LegacyPlanAndAct provides backward compatibility with old signature.
func (g *Gateway) LegacyPlanAndAct(tenantID, prompt string, proposedTools []string) (types.GuardrailResult, error) {
	ctx := context.Background()
	resp, err := g.PlanAndAct(ctx, PlanRequest{
		TenantID: tenantID,
		Prompt:   prompt,
		Tools:    proposedTools,
	})
	if err != nil {
		return types.GuardrailResult{
			Allowed: false,
			Reason:  resp.Reason,
			Signals: resp.Signals,
		}, err
	}
	return types.GuardrailResult{Allowed: resp.Allowed}, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
