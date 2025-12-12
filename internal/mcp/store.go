package mcp

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Store persists capabilities in Postgres.
type Store struct {
	db *sql.DB
}

// NewStore constructs Store.
func NewStore(db *sql.DB) *Store { return &Store{db: db} }

// Add inserts a capability.
func (s *Store) Add(c Capability) (Capability, error) {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now().UTC()
	}
	tags, _ := json.Marshal(c.Tags)
	_, err := s.db.Exec(`INSERT INTO capabilities (id, name, description, tags, created_at) VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (name) DO UPDATE SET description=excluded.description, tags=excluded.tags`, c.ID, c.Name, c.Description, tags, c.CreatedAt)
	return c, err
}

// List returns capabilities filtered by optional tag.
func (s *Store) List(tag string) ([]Capability, error) {
	query := `SELECT id, name, description, tags, created_at FROM capabilities`
	args := []any{}
	if tag != "" {
		query += ` WHERE tags ? $1`
		args = append(args, tag)
	}
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Capability
	for rows.Next() {
		var c Capability
		var tags []byte
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &tags, &c.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(tags, &c.Tags)
		out = append(out, c)
	}
	return out, nil
}

// FilterAllowed returns only capabilities allowed by policy.
func (s *Store) FilterAllowed(all []Capability, allowList []string) []Capability {
	if len(allowList) == 0 {
		return all
	}
	set := map[string]struct{}{}
	for _, a := range allowList {
		set[strings.ToLower(a)] = struct{}{}
	}
	var out []Capability
	for _, c := range all {
		if _, ok := set[strings.ToLower(c.Name)]; ok {
			out = append(out, c)
		}
	}
	return out
}

