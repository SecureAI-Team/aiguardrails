package agent

import (
	"context"
	"errors"
	"sync"
)

var ErrExecutorNotFound = errors.New("executor not found")

// Executor defines interface for tool execution.
type Executor interface {
	// Execute runs the tool with given arguments.
	Execute(ctx context.Context, args map[string]interface{}) (interface{}, error)
	// Name returns the tool name.
	Name() string
	// Description returns tool description for LLM.
	Description() string
	// Schema returns JSON schema for arguments.
	Schema() map[string]interface{}
}

// Registry manages tool executors.
type Registry struct {
	mu        sync.RWMutex
	executors map[string]Executor
}

// NewRegistry creates an empty registry.
func NewRegistry() *Registry {
	return &Registry{
		executors: make(map[string]Executor),
	}
}

// Register adds an executor to the registry.
func (r *Registry) Register(exec Executor) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.executors[exec.Name()] = exec
	return nil
}

// Unregister removes an executor.
func (r *Registry) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.executors, name)
}

// Get returns an executor by name.
func (r *Registry) Get(name string) (Executor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	exec, ok := r.executors[name]
	if !ok {
		return nil, ErrExecutorNotFound
	}
	return exec, nil
}

// List returns all registered executor names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.executors))
	for name := range r.executors {
		names = append(names, name)
	}
	return names
}

// Describe returns all executor descriptions for LLM context.
func (r *Registry) Describe() []ToolDescription {
	r.mu.RLock()
	defer r.mu.RUnlock()
	descs := make([]ToolDescription, 0, len(r.executors))
	for _, exec := range r.executors {
		descs = append(descs, ToolDescription{
			Name:        exec.Name(),
			Description: exec.Description(),
			Schema:      exec.Schema(),
		})
	}
	return descs
}

// ToolDescription for LLM.
type ToolDescription struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema,omitempty"`
}

// NoOpExecutor is a stub executor for testing.
type NoOpExecutor struct {
	name string
	desc string
}

func NewNoOpExecutor(name, desc string) *NoOpExecutor {
	return &NoOpExecutor{name: name, desc: desc}
}

func (e *NoOpExecutor) Name() string        { return e.name }
func (e *NoOpExecutor) Description() string { return e.desc }
func (e *NoOpExecutor) Schema() map[string]interface{} {
	return map[string]interface{}{"type": "object"}
}
func (e *NoOpExecutor) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return map[string]interface{}{
		"tool":    e.name,
		"status":  "noop",
		"message": "NoOp executor stub",
	}, nil
}
