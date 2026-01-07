package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// RequestsTotal counts all guardrail requests.
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "guardrails_requests_total",
			Help: "Total number of guardrail requests",
		},
		[]string{"endpoint", "status", "tenant_id"},
	)

	// BlockedTotal counts blocked requests.
	BlockedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "guardrails_blocked_total",
			Help: "Total number of blocked requests",
		},
		[]string{"reason", "tenant_id", "severity"},
	)

	// LatencyHistogram measures request latency.
	LatencyHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "guardrails_latency_seconds",
			Help:    "Request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)

	// OPAEvalDuration measures OPA evaluation time.
	OPAEvalDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "guardrails_opa_eval_seconds",
			Help:    "OPA evaluation duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
		},
	)

	// ActiveAgentLoops tracks currently running agent loops.
	ActiveAgentLoops = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "guardrails_active_agent_loops",
			Help: "Number of currently active agent loops",
		},
	)

	// PolicyVersionGauge tracks current policy versions by tenant.
	PolicyVersionGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "guardrails_policy_version",
			Help: "Current policy version by tenant",
		},
		[]string{"tenant_id", "policy_id"},
	)

	// OPAReloadTotal counts OPA policy reloads.
	OPAReloadTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "guardrails_opa_reload_total",
			Help: "Total number of OPA policy reloads",
		},
	)

	// AlertsFiredTotal counts alerts fired.
	AlertsFiredTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "guardrails_alerts_fired_total",
			Help: "Total number of alerts fired",
		},
		[]string{"rule", "severity"},
	)
)

// RecordRequest records a guardrail request.
func RecordRequest(endpoint, status, tenantID string) {
	RequestsTotal.WithLabelValues(endpoint, status, tenantID).Inc()
}

// RecordBlocked records a blocked request.
func RecordBlocked(reason, tenantID, severity string) {
	BlockedTotal.WithLabelValues(reason, tenantID, severity).Inc()
}

// RecordLatency records request latency.
func RecordLatency(endpoint string, duration time.Duration) {
	LatencyHistogram.WithLabelValues(endpoint).Observe(duration.Seconds())
}

// RecordOPAEval records OPA evaluation duration.
func RecordOPAEval(duration time.Duration) {
	OPAEvalDuration.Observe(duration.Seconds())
}

// AgentLoopStart increments active agent loops.
func AgentLoopStart() {
	ActiveAgentLoops.Inc()
}

// AgentLoopEnd decrements active agent loops.
func AgentLoopEnd() {
	ActiveAgentLoops.Dec()
}

// RecordAlert records an alert fired.
func RecordAlert(rule, severity string) {
	AlertsFiredTotal.WithLabelValues(rule, severity).Inc()
}

// Handler returns the prometheus HTTP handler.
func Handler() http.Handler {
	return promhttp.Handler()
}

// Middleware returns a metrics middleware for chi router.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		RecordLatency(r.URL.Path, time.Since(start))
		status := "success"
		if wrapped.statusCode >= 400 {
			status = "error"
		}
		RecordRequest(r.URL.Path, status, r.Header.Get("X-Tenant-ID"))
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
