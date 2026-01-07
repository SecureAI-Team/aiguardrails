package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents an admin/local user.
type User struct {
	ID           string
	Username     string
	PasswordHash string
	Role         string
	Email        string
	DisplayName  string
	Status       string // active, inactive, suspended
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastLoginAt  *time.Time
}

// UserStore provides CRUD for users.
type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(username, password, role string) (*User, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password required")
	}
	if role == "" {
		role = "tenant_admin"
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	id := uuid.NewString()
	now := time.Now().UTC()
	_, err = s.db.Exec(`INSERT INTO users (id, username, password_hash, role, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		id, username, string(hash), role, now, now)
	if err != nil {
		return nil, err
	}
	return &User{ID: id, Username: username, PasswordHash: string(hash), Role: role, CreatedAt: now, UpdatedAt: now}, nil
}

func (s *UserStore) GetByUsername(username string) (*User, error) {
	var u User
	err := s.db.QueryRow(`SELECT id, username, password_hash, role, created_at, updated_at FROM users WHERE username=$1`, username).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *UserStore) EnsureBootUser(username, password, role string) (*User, error) {
	if username == "" || password == "" {
		return nil, errors.New("boot user requires username and password")
	}
	u, err := s.GetByUsername(username)
	if err == nil {
		return u, nil
	}
	return s.Create(username, password, role)
}

func (s *UserStore) Verify(username, password string) (*User, error) {
	u, err := s.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return u, nil
}

// List returns all users with optional filters.
func (s *UserStore) List(role, status string, limit, offset int) ([]User, error) {
	query := `SELECT id, username, password_hash, role, COALESCE(email,''), COALESCE(display_name,''), COALESCE(status,'active'), created_at, updated_at, last_login_at FROM users WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if role != "" {
		query += ` AND role=$` + string(rune('0'+argIdx))
		args = append(args, role)
		argIdx++
	}
	if status != "" {
		query += ` AND status=$` + string(rune('0'+argIdx))
		args = append(args, status)
		argIdx++
	}
	query += ` ORDER BY created_at DESC`
	if limit > 0 {
		query += ` LIMIT $` + string(rune('0'+argIdx))
		args = append(args, limit)
		argIdx++
	}
	if offset > 0 {
		query += ` OFFSET $` + string(rune('0'+argIdx))
		args = append(args, offset)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		var lastLogin sql.NullTime
		if err := rows.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.Email, &u.DisplayName, &u.Status, &u.CreatedAt, &u.UpdatedAt, &lastLogin); err != nil {
			return nil, err
		}
		if lastLogin.Valid {
			u.LastLoginAt = &lastLogin.Time
		}
		u.PasswordHash = "" // Don't expose
		users = append(users, u)
	}
	return users, nil
}

// GetByID returns user by ID.
func (s *UserStore) GetByID(id string) (*User, error) {
	var u User
	var lastLogin sql.NullTime
	err := s.db.QueryRow(`SELECT id, username, password_hash, role, COALESCE(email,''), COALESCE(display_name,''), COALESCE(status,'active'), created_at, updated_at, last_login_at FROM users WHERE id=$1`, id).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.Email, &u.DisplayName, &u.Status, &u.CreatedAt, &u.UpdatedAt, &lastLogin)
	if err != nil {
		return nil, err
	}
	if lastLogin.Valid {
		u.LastLoginAt = &lastLogin.Time
	}
	u.PasswordHash = ""
	return &u, nil
}

// Update updates user info (not password).
func (s *UserStore) Update(id string, role, email, displayName, status string) error {
	now := time.Now().UTC()
	_, err := s.db.Exec(`UPDATE users SET role=$1, email=$2, display_name=$3, status=$4, updated_at=$5 WHERE id=$6`,
		role, email, displayName, status, now, id)
	return err
}

// Delete removes a user.
func (s *UserStore) Delete(id string) error {
	_, err := s.db.Exec(`DELETE FROM users WHERE id=$1`, id)
	return err
}

// UpdatePassword changes user password.
func (s *UserStore) UpdatePassword(id, newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	_, err = s.db.Exec(`UPDATE users SET password_hash=$1, updated_at=$2 WHERE id=$3`, string(hash), now, id)
	return err
}

// UpdateLastLogin records login timestamp.
func (s *UserStore) UpdateLastLogin(id string) error {
	now := time.Now().UTC()
	_, err := s.db.Exec(`UPDATE users SET last_login_at=$1 WHERE id=$2`, now, id)
	return err
}
