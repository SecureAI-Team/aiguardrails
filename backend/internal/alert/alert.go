package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// Event represents an audit/security event.
type Event struct {
	Type      string            `json:"type"`
	TenantID  string            `json:"tenant_id"`
	AppID     string            `json:"app_id,omitempty"`
	Reason    string            `json:"reason,omitempty"`
	Severity  string            `json:"severity,omitempty"`
	Signals   []string          `json:"signals,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// Alert represents a triggered alert.
type Alert struct {
	Rule      string    `json:"rule"`
	Severity  string    `json:"severity"`
	Message   string    `json:"message"`
	Event     Event     `json:"event"`
	Timestamp time.Time `json:"timestamp"`
}

// Rule defines when to trigger an alert.
type Rule struct {
	Name      string
	Severity  string // critical, high, medium, low
	Condition func(Event) bool
}

// Manager handles alert rules and notifications.
type Manager struct {
	mu         sync.RWMutex
	rules      []Rule
	webhooks   []string
	httpClient *http.Client
	alertCh    chan Alert
}

// NewManager creates an alert manager.
func NewManager(webhooks []string) *Manager {
	m := &Manager{
		rules:      defaultRules(),
		webhooks:   webhooks,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		alertCh:    make(chan Alert, 100),
	}
	return m
}

// Start begins processing alerts.
func (m *Manager) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case alert := <-m.alertCh:
				m.send(alert)
			}
		}
	}()
}

// AddRule adds a custom alert rule.
func (m *Manager) AddRule(rule Rule) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rules = append(m.rules, rule)
}

// Evaluate checks event against all rules.
func (m *Manager) Evaluate(event Event) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, rule := range m.rules {
		if rule.Condition(event) {
			alert := Alert{
				Rule:      rule.Name,
				Severity:  rule.Severity,
				Message:   event.Reason,
				Event:     event,
				Timestamp: time.Now(),
			}
			select {
			case m.alertCh <- alert:
			default:
				// Channel full, drop alert
			}
		}
	}
}

// send dispatches alert to all webhooks.
func (m *Manager) send(alert Alert) {
	body, _ := json.Marshal(alert)
	for _, webhook := range m.webhooks {
		go func(url string) {
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			if err != nil {
				return
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := m.httpClient.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()
		}(webhook)
	}
}

// defaultRules returns built-in alert rules.
func defaultRules() []Rule {
	return []Rule{
		{
			Name:     "critical_block",
			Severity: "critical",
			Condition: func(e Event) bool {
				return e.Severity == "critical" && e.Type == "blocked"
			},
		},
		{
			Name:     "industrial_safety_violation",
			Severity: "critical",
			Condition: func(e Event) bool {
				for _, sig := range e.Signals {
					if sig == "severity:critical" {
						return true
					}
				}
				return false
			},
		},
		{
			Name:     "prompt_injection_detected",
			Severity: "high",
			Condition: func(e Event) bool {
				return e.Reason == "prompt_injection_detected"
			},
		},
		{
			Name:     "opa_block",
			Severity: "high",
			Condition: func(e Event) bool {
				return e.Reason == "opa_block" ||
					e.Reason == "opa_industrial_prompt_stop" ||
					e.Reason == "opa_industrial_prompt_override"
			},
		},
		{
			Name:     "dlp_violation",
			Severity: "medium",
			Condition: func(e Event) bool {
				return e.Reason == "dlp_detected" || e.Reason == "sensitive_data_exposed"
			},
		},
	}
}
