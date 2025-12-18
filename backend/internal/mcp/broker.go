package mcp

import "aiguardrails/internal/policy"

// Broker mediates MCP provider calls with allowlists.
type Broker struct {
	policy policy.Engine
	store  *Store
}

// NewBroker constructs a Broker.
func NewBroker(p policy.Engine, store *Store) *Broker {
	return &Broker{policy: p, store: store}
}

// AllowCapability checks if capability is allowed for tenant.
func (b *Broker) AllowCapability(tenantID, capability string) bool {
	return b.policy.AllowTool(tenantID, capability)
}

// DescribeCapability returns registry info if present.
func (b *Broker) DescribeCapability(name string) (Capability, bool) {
	if b.store == nil {
		return Capability{}, false
	}
	all, err := b.store.List("")
	if err != nil {
		return Capability{}, false
	}
	for _, c := range all {
		if c.Name == name {
			return c, true
		}
	}
	return Capability{}, false
}

