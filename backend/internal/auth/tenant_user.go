package auth

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// TenantUser represents user membership in a tenant.
type TenantUser struct {
	ID        string
	TenantID  string
	UserID    string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
	// Joined fields
	Username    string
	Email       string
	DisplayName string
}

// TenantUserStore manages tenant-user associations.
type TenantUserStore struct {
	db *sql.DB
}

// NewTenantUserStore creates store instance.
func NewTenantUserStore(db *sql.DB) *TenantUserStore {
	return &TenantUserStore{db: db}
}

// Add associates a user with a tenant.
func (s *TenantUserStore) Add(tenantID, userID, role string) (*TenantUser, error) {
	id := uuid.NewString()
	now := time.Now().UTC()
	if role == "" {
		role = "tenant_user"
	}
	_, err := s.db.Exec(`INSERT INTO tenant_users (id, tenant_id, user_id, role, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		id, tenantID, userID, role, now, now)
	if err != nil {
		return nil, err
	}
	return &TenantUser{ID: id, TenantID: tenantID, UserID: userID, Role: role, CreatedAt: now, UpdatedAt: now}, nil
}

// Remove removes user from tenant.
func (s *TenantUserStore) Remove(tenantID, userID string) error {
	_, err := s.db.Exec(`DELETE FROM tenant_users WHERE tenant_id=$1 AND user_id=$2`, tenantID, userID)
	return err
}

// UpdateRole changes user role in tenant.
func (s *TenantUserStore) UpdateRole(tenantID, userID, role string) error {
	now := time.Now().UTC()
	_, err := s.db.Exec(`UPDATE tenant_users SET role=$1, updated_at=$2 WHERE tenant_id=$3 AND user_id=$4`,
		role, now, tenantID, userID)
	return err
}

// ListByTenant returns all users in a tenant.
func (s *TenantUserStore) ListByTenant(tenantID string) ([]TenantUser, error) {
	rows, err := s.db.Query(`
		SELECT tu.id, tu.tenant_id, tu.user_id, tu.role, tu.created_at, tu.updated_at,
		       u.username, COALESCE(u.email,''), COALESCE(u.display_name,'')
		FROM tenant_users tu
		JOIN users u ON tu.user_id = u.id
		WHERE tu.tenant_id = $1
		ORDER BY tu.created_at`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []TenantUser
	for rows.Next() {
		var tu TenantUser
		if err := rows.Scan(&tu.ID, &tu.TenantID, &tu.UserID, &tu.Role, &tu.CreatedAt, &tu.UpdatedAt,
			&tu.Username, &tu.Email, &tu.DisplayName); err != nil {
			return nil, err
		}
		users = append(users, tu)
	}
	return users, nil
}

// ListByUser returns all tenants a user belongs to.
func (s *TenantUserStore) ListByUser(userID string) ([]TenantUser, error) {
	rows, err := s.db.Query(`SELECT id, tenant_id, user_id, role, created_at, updated_at FROM tenant_users WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberships []TenantUser
	for rows.Next() {
		var tu TenantUser
		if err := rows.Scan(&tu.ID, &tu.TenantID, &tu.UserID, &tu.Role, &tu.CreatedAt, &tu.UpdatedAt); err != nil {
			return nil, err
		}
		memberships = append(memberships, tu)
	}
	return memberships, nil
}

// Get returns specific tenant-user association.
func (s *TenantUserStore) Get(tenantID, userID string) (*TenantUser, error) {
	var tu TenantUser
	err := s.db.QueryRow(`SELECT id, tenant_id, user_id, role, created_at, updated_at FROM tenant_users WHERE tenant_id=$1 AND user_id=$2`,
		tenantID, userID).Scan(&tu.ID, &tu.TenantID, &tu.UserID, &tu.Role, &tu.CreatedAt, &tu.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &tu, nil
}
