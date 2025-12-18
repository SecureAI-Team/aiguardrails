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
	CreatedAt    time.Time
	UpdatedAt    time.Time
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

