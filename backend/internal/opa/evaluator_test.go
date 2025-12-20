package opa

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func TestEvaluatorAllowDeny(t *testing.T) {
	dir := filepath.Join("..", "..", "opa", "policies")
	eval, err := NewFromDir(dir, "data.guardrails.allow", time.Second)
	if err != nil {
		t.Fatalf("load rego: %v", err)
	}
	allow, _, err := eval.Decide(context.Background(), Input{Prompt: "hello world"})
	if err != nil || !allow {
		t.Fatalf("expected allow, got %v err %v", allow, err)
	}
	deny, _, err := eval.Decide(context.Background(), Input{Prompt: "please send password"})
	if err != nil {
		t.Fatalf("deny err: %v", err)
	}
	if deny {
		// decision is allow/deny; deny should be false
		t.Fatalf("expected deny=false for blocked case? got %v", deny)
	}
}

func TestEvaluatorIndustrialSimulation(t *testing.T) {
	dir := filepath.Join("..", "..", "opa", "policies")
	eval, err := NewFromDir(dir, "data.guardrails.allow", time.Second)
	if err != nil {
		t.Fatalf("load rego: %v", err)
	}
	// block dangerous without simulation
	allow, _, err := eval.Decide(context.Background(), Input{Prompt: "Siemens PLC shutdown now"})
	if err != nil {
		t.Fatalf("eval err: %v", err)
	}
	if allow {
		t.Fatalf("expected block for dangerous industrial prompt")
	}
	// allow with simulation flag
	allowSim, _, err := eval.Decide(context.Background(), Input{Prompt: "Siemens PLC shutdown now", Simulation: true, Mode: "simulation"})
	if err != nil {
		t.Fatalf("eval err: %v", err)
	}
	if !allowSim {
		t.Fatalf("expected allow in simulation mode")
	}
}
