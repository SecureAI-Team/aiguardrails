package opa

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/open-policy-agent/opa/rego"
)

// Input defines data passed to OPA.
type Input struct {
	TenantID   string      `json:"tenantId"`
	AppID      string      `json:"appId"`
	Mode       string      `json:"mode"` // prompt_check | output_filter | tool
	Prompt     string      `json:"prompt,omitempty"`
	Output     string      `json:"output,omitempty"`
	Tool       string      `json:"tool,omitempty"`
	Simulation bool        `json:"simulation,omitempty"`
	Namespaces []string    `json:"namespaces,omitempty"`
	Signals    interface{} `json:"signals,omitempty"`
}

// Evaluator wraps OPA rego evaluation.
type Evaluator struct {
	query   string
	modules map[string]string
	timeout time.Duration
}

// NewFromDir loads all .rego files under dir and builds evaluator.
func NewFromDir(dir, decision string, timeout time.Duration) (*Evaluator, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	modules := map[string]string{}
	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != ".rego" {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
		modules[f.Name()] = string(b)
	}
	if len(modules) == 0 {
		return nil, fmt.Errorf("no rego modules found in %s", dir)
	}
	return &Evaluator{query: decision, modules: modules, timeout: timeout}, nil
}

// Decide returns (allow, reason, signals, error).
func (e *Evaluator) Decide(ctx context.Context, in Input) (bool, interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	opts := []func(*rego.Rego){
		rego.Query(e.query),
		rego.Input(in),
	}
	for name, mod := range e.modules {
		opts = append(opts, rego.Module(name, mod))
	}
	r := rego.New(opts...)
	rs, err := r.Eval(ctx)
	if err != nil {
		return false, nil, err
	}
	if len(rs) == 0 || len(rs[0].Expressions) == 0 {
		return false, nil, fmt.Errorf("empty decision")
	}
	val := rs[0].Expressions[0].Value
	// Expecting allow boolean or object with allow/bool and reason/signals
	switch v := val.(type) {
	case bool:
		return v, nil, nil
	case map[string]interface{}:
		allow, ok := v["allow"].(bool)
		if !ok {
			return false, nil, fmt.Errorf("decision missing allow")
		}
		return allow, v, nil
	default:
		return false, nil, fmt.Errorf("unexpected decision type %T", v)
	}
}
