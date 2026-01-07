package agent

import (
	"context"
	"errors"
	"strings"
	"time"
)

var (
	ErrToolNotAllowed = errors.New("tool not allowed in sandbox")
	ErrTimeout        = errors.New("sandbox execution timeout")
	ErrMemoryLimit    = errors.New("memory limit exceeded")
)

// Sandbox provides isolated execution environment for tools.
type Sandbox struct {
	timeout     time.Duration
	memLimit    int64 // bytes
	allowedCmds []string
	denyList    []string
}

// SandboxOption configures Sandbox.
type SandboxOption func(*Sandbox)

// WithDenyList sets commands that are always denied.
func WithDenyList(cmds []string) SandboxOption {
	return func(s *Sandbox) { s.denyList = cmds }
}

// NewSandbox creates a new Sandbox with constraints.
func NewSandbox(timeout time.Duration, memLimit int64, allowedCmds []string, opts ...SandboxOption) *Sandbox {
	s := &Sandbox{
		timeout:     timeout,
		memLimit:    memLimit,
		allowedCmds: allowedCmds,
		denyList: []string{
			"rm", "del", "format", "shutdown", "reboot",
			"dd", "mkfs", "fdisk", "kill", "pkill",
			"chmod", "chown", "sudo", "su",
		},
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// ValidateTool checks if a tool is allowed to execute.
func (s *Sandbox) ValidateTool(tool string) error {
	toolLower := strings.ToLower(tool)

	// Check deny list
	for _, denied := range s.denyList {
		if strings.Contains(toolLower, denied) {
			return ErrToolNotAllowed
		}
	}

	// If allowlist is set, check it
	if len(s.allowedCmds) > 0 {
		allowed := false
		for _, cmd := range s.allowedCmds {
			if strings.EqualFold(cmd, tool) {
				allowed = true
				break
			}
		}
		if !allowed {
			return ErrToolNotAllowed
		}
	}

	return nil
}

// Execute runs a tool in sandbox with constraints.
func (s *Sandbox) Execute(ctx context.Context, tool string, args map[string]interface{}) (interface{}, error) {
	// Validate tool first
	if err := s.ValidateTool(tool); err != nil {
		return nil, err
	}

	// Create timeout context
	execCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// Execute in goroutine with result channel
	resultCh := make(chan sandboxResult, 1)
	go func() {
		result, err := s.executeInternal(tool, args)
		resultCh <- sandboxResult{result: result, err: err}
	}()

	// Wait for result or timeout
	select {
	case <-execCtx.Done():
		return nil, ErrTimeout
	case res := <-resultCh:
		return res.result, res.err
	}
}

type sandboxResult struct {
	result interface{}
	err    error
}

// executeInternal performs the actual execution.
// In production, this would use process isolation or containers.
func (s *Sandbox) executeInternal(tool string, args map[string]interface{}) (interface{}, error) {
	// Stub implementation - in production:
	// 1. Fork a subprocess with resource limits (cgroups, rlimit)
	// 2. Set up seccomp filters
	// 3. Use namespaces for isolation
	// 4. Monitor memory usage

	// For now, simulate execution
	time.Sleep(10 * time.Millisecond)

	return map[string]interface{}{
		"tool":   tool,
		"args":   args,
		"status": "executed",
		"output": "Sandbox execution stub",
	}, nil
}

// SetTimeout updates the timeout.
func (s *Sandbox) SetTimeout(t time.Duration) {
	s.timeout = t
}

// SetMemLimit updates memory limit.
func (s *Sandbox) SetMemLimit(bytes int64) {
	s.memLimit = bytes
}

// AddAllowed adds tools to allowlist.
func (s *Sandbox) AddAllowed(tools ...string) {
	s.allowedCmds = append(s.allowedCmds, tools...)
}

// AddDenied adds tools to denylist.
func (s *Sandbox) AddDenied(tools ...string) {
	s.denyList = append(s.denyList, tools...)
}
