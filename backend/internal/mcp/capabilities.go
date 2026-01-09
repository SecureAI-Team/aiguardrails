package mcp

import "time"

// Capability represents a tool/capability entry.
type Capability struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
}

// Registry holds allowed capabilities.
type Registry struct {
	entries map[string]Capability
}

// NewRegistry constructs a registry with optional seed.
func NewRegistry(seed []Capability) *Registry {
	r := &Registry{entries: map[string]Capability{}}
	for _, c := range seed {
		r.entries[c.Name] = c
	}
	return r
}

// Get returns capability by name.
func (r *Registry) Get(name string) (Capability, bool) {
	c, ok := r.entries[name]
	return c, ok
}

// List returns all capabilities.
func (r *Registry) List() []Capability {
	out := make([]Capability, 0, len(r.entries))
	for _, c := range r.entries {
		out = append(out, c)
	}
	return out
}
