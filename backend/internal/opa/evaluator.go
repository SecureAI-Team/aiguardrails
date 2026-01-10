package opa

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
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
	Rules      []string    `json:"rules,omitempty"`
	Signals    interface{} `json:"signals,omitempty"`
}

// Evaluator wraps OPA rego evaluation with hot-reload support.
type Evaluator struct {
	mu       sync.RWMutex
	query    string
	modules  map[string]string
	timeout  time.Duration
	version  int64
	onChange func(version int64)
}

// NewFromDir loads all .rego files under dir and builds evaluator.
func NewFromDir(dir, decision string, timeout time.Duration) (*Evaluator, error) {
	modules, err := loadModules(dir)
	if err != nil {
		return nil, err
	}
	if len(modules) == 0 {
		return nil, fmt.Errorf("no rego modules found in %s", dir)
	}
	return &Evaluator{
		query:   decision,
		modules: modules,
		timeout: timeout,
		version: 1,
	}, nil
}

// loadModules reads all .rego files from directory.
func loadModules(dir string) (map[string]string, error) {
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
	return modules, nil
}

// Reload reloads modules from directory or raw content.
func (e *Evaluator) Reload(dir string) error {
	modules, err := loadModules(dir)
	if err != nil {
		return err
	}
	if len(modules) == 0 {
		return fmt.Errorf("no rego modules found in %s", dir)
	}

	e.mu.Lock()
	e.modules = modules
	e.version++
	version := e.version
	onChange := e.onChange
	e.mu.Unlock()

	if onChange != nil {
		onChange(version)
	}
	return nil
}

// ReloadFromContent reloads modules from provided content map.
func (e *Evaluator) ReloadFromContent(modules map[string]string) error {
	if len(modules) == 0 {
		return fmt.Errorf("empty modules")
	}

	e.mu.Lock()
	e.modules = modules
	e.version++
	version := e.version
	onChange := e.onChange
	e.mu.Unlock()

	if onChange != nil {
		onChange(version)
	}
	return nil
}

// Version returns current policy version.
func (e *Evaluator) Version() int64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.version
}

// SetOnChange sets callback for version changes.
func (e *Evaluator) SetOnChange(fn func(version int64)) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.onChange = fn
}

// Modules returns current module names.
func (e *Evaluator) Modules() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	names := make([]string, 0, len(e.modules))
	for name := range e.modules {
		names = append(names, name)
	}
	return names
}

// Decide returns (allow, data, error).
func (e *Evaluator) Decide(ctx context.Context, in Input) (bool, interface{}, error) {
	e.mu.RLock()
	modules := e.modules
	query := e.query
	timeout := e.timeout
	e.mu.RUnlock()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	opts := []func(*rego.Rego){
		rego.Query(query),
		rego.Input(in),
	}
	for name, mod := range modules {
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
