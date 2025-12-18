package usage

import (
	"sync"
	"time"

	"aiguardrails/internal/types"
)

// Meter records usage for apps.
type Meter struct {
	mu     sync.Mutex
	events []types.UsageRecord
}

// NewMeter creates a Meter.
func NewMeter() *Meter {
	return &Meter{events: []types.UsageRecord{}}
}

// Record app usage count.
func (m *Meter) Record(appID string, count int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, types.UsageRecord{AppID: appID, Count: count, Timestamp: time.Now().UTC()})
}

// Summary returns aggregated counts for the app in the last hour.
func (m *Meter) Summary(appID string, since time.Time) int64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	var total int64
	for _, e := range m.events {
		if e.AppID == appID && e.Timestamp.After(since) {
			total += e.Count
		}
	}
	return total
}

