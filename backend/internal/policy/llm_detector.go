package policy

import (
	"context"
	"sync"
	"time"

	"aiguardrails/internal/types"
)

// LLMClient defines moderation client behavior.
type LLMClient interface {
	Moderate(ctx context.Context, text string) (types.GuardrailResult, error)
}

type job struct {
	text string
}

// LLMDetector runs async moderation with caching.
type LLMDetector struct {
	client LLMClient
	queue  chan job
	mu     sync.RWMutex
	cache  map[string]cacheEntry
	tokens <-chan time.Time
	ttl    time.Duration
	onDone func(text string, res types.GuardrailResult)
}

type cacheEntry struct {
	res     types.GuardrailResult
	expires time.Time
}

// NewLLMDetector constructs a detector with queue size, ttl, and rate limit.
func NewLLMDetector(client LLMClient, queueSize int, ttl time.Duration, rps int) *LLMDetector {
	var tokens <-chan time.Time
	if rps > 0 {
		tokens = time.NewTicker(time.Second / time.Duration(rps)).C
	}
	return &LLMDetector{
		client: client,
		queue:  make(chan job, queueSize),
		cache:  map[string]cacheEntry{},
		tokens: tokens,
		ttl:    ttl,
	}
}

// Start launches worker goroutines.
func (d *LLMDetector) Start(workers int) {
	for i := 0; i < workers; i++ {
		go d.worker()
	}
}

// worker processes moderation jobs.
func (d *LLMDetector) worker() {
	for j := range d.queue {
		if d.tokens != nil {
			<-d.tokens
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		res, err := d.client.Moderate(ctx, j.text)
		cancel()
		if err != nil {
			res = types.GuardrailResult{Allowed: true, Reason: "llm_error"}
		}
		d.mu.Lock()
		d.cache[j.text] = cacheEntry{res: res, expires: time.Now().Add(d.ttl)}
		d.mu.Unlock()
		if d.onDone != nil {
			d.onDone(j.text, res)
		}
	}
}

// Check returns cached result or enqueues for evaluation.
func (d *LLMDetector) Check(text string) types.GuardrailResult {
	d.mu.RLock()
	if entry, ok := d.cache[text]; ok && entry.expires.After(time.Now()) {
		d.mu.RUnlock()
		return entry.res
	}
	d.mu.RUnlock()
	// drop expired if existed
	d.mu.Lock()
	delete(d.cache, text)
	d.mu.Unlock()
	select {
	case d.queue <- job{text: text}:
	default:
		// queue full; do nothing
	}
	return types.GuardrailResult{Allowed: true, Reason: "llm_pending"}
}

// WithCallback sets optional completion callback.
func (d *LLMDetector) WithCallback(cb func(text string, res types.GuardrailResult)) {
	d.onDone = cb
}

