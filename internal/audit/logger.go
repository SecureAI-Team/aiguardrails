package audit

import "log"

// Logger is a simple audit sink; replace with structured logging later.
type Logger struct{}

// NewLogger constructs a Logger.
func NewLogger() *Logger { return &Logger{} }

// Record logs an audit event.
func (l *Logger) Record(event string, fields map[string]string) {
	log.Printf("[AUDIT] %s %v", event, fields)
}

// RecordStore logs and persists if store provided.
func (l *Logger) RecordStore(store *Store, event string, fields map[string]string) {
	l.Record(event, fields)
	if store != nil {
		_ = store.Record(event, fields)
	}
}

