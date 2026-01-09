package audit

import (
	"database/sql"
	"encoding/json"
	"time"
	"strings"
	"strconv"
)

// Store persists audit events.
type Store struct {
	db *sql.DB
}

// NewStore constructs Store.
func NewStore(db *sql.DB) *Store { return &Store{db: db} }

// Init creates table if not exists.
func (s *Store) Init() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS audit_events (
		id SERIAL PRIMARY KEY,
		event TEXT NOT NULL,
		fields JSONB NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`)
	return err
}

// Record saves an event.
func (s *Store) Record(event string, fields map[string]string) error {
	b, err := json.Marshal(fields)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`INSERT INTO audit_events (event, fields) VALUES ($1, $2)`, event, string(b))
	return err
}

// List returns recent events with optional filters.
func (s *Store) List(limit int, eventLike, tenant string) ([]map[string]any, error) {
	query := `SELECT event, fields, created_at FROM audit_events`
	args := []any{}
	clauses := []string{}
	if eventLike != "" {
		args = append(args, "%"+eventLike+"%")
		clauses = append(clauses, `event ILIKE $`+itoa(len(args)))
	}
	if tenant != "" {
		args = append(args, tenant)
		clauses = append(clauses, `fields->>'tenant_id' = $`+itoa(len(args)))
	}
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}
	args = append(args, limit)
	query += " ORDER BY created_at DESC LIMIT $" + itoa(len(args))

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]any
	for rows.Next() {
		var event string
		var fieldsJSON []byte
		var ts time.Time
		if err := rows.Scan(&event, &fieldsJSON, &ts); err != nil {
			return nil, err
		}
		var fields map[string]any
		if err := json.Unmarshal(fieldsJSON, &fields); err != nil {
			// fallback/ignore error?
			fields = map[string]any{}
		}

		fields["event"] = event
		fields["created_at"] = ts
		out = append(out, fields)
	}
	return out, nil
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

