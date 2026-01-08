package org

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Organization 组织
type Organization struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	LogoURL     string    `json:"logo_url,omitempty"`
	SSOEnabled  bool      `json:"sso_enabled"`
	SSOProvider string    `json:"sso_provider,omitempty"`
	MaxTeams    int       `json:"max_teams"`
	MaxMembers  int       `json:"max_members"`
	MaxTenants  int       `json:"max_tenants"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Team 团队
type Team struct {
	ID          string    `json:"id"`
	OrgID       string    `json:"org_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DefaultRole string    `json:"default_role"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
}

// Member 成员
type Member struct {
	ID        string     `json:"id"`
	OrgID     string     `json:"org_id"`
	UserID    string     `json:"user_id"`
	Role      string     `json:"role"`
	TeamIDs   []string   `json:"team_ids"`
	Status    string     `json:"status"`
	InvitedBy *string    `json:"invited_by,omitempty"`
	JoinedAt  *time.Time `json:"joined_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	// Joined fields
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

// IPWhitelist IP白名单
type IPWhitelist struct {
	ID          string    `json:"id"`
	ScopeType   string    `json:"scope_type"`
	ScopeID     string    `json:"scope_id"`
	IPAddress   string    `json:"ip_address,omitempty"`
	IPCIDR      string    `json:"ip_cidr,omitempty"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
}

// Store 组织存储
type Store struct {
	db *sql.DB
}

// NewStore 创建存储
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// CreateOrg 创建组织
func (s *Store) CreateOrg(org *Organization) error {
	org.ID = uuid.NewString()
	now := time.Now().UTC()
	org.CreatedAt = now
	org.UpdatedAt = now
	org.Status = "active"

	_, err := s.db.Exec(`INSERT INTO organizations (id, name, slug, description, logo_url, max_teams, max_members, max_tenants, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		org.ID, org.Name, org.Slug, org.Description, org.LogoURL, org.MaxTeams, org.MaxMembers, org.MaxTenants, org.Status, org.CreatedAt, org.UpdatedAt)
	return err
}

// GetOrg 获取组织
func (s *Store) GetOrg(id string) (*Organization, error) {
	var org Organization
	err := s.db.QueryRow(`SELECT id, name, slug, description, logo_url, sso_enabled, COALESCE(sso_provider,''), 
		max_teams, max_members, max_tenants, status, created_at, updated_at
		FROM organizations WHERE id = $1`, id).
		Scan(&org.ID, &org.Name, &org.Slug, &org.Description, &org.LogoURL, &org.SSOEnabled, &org.SSOProvider,
			&org.MaxTeams, &org.MaxMembers, &org.MaxTenants, &org.Status, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// ListOrgs 列出组织
func (s *Store) ListOrgs() ([]Organization, error) {
	rows, err := s.db.Query(`SELECT id, name, slug, description, logo_url, sso_enabled, status, created_at
		FROM organizations WHERE status = 'active' ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []Organization
	for rows.Next() {
		var o Organization
		err := rows.Scan(&o.ID, &o.Name, &o.Slug, &o.Description, &o.LogoURL, &o.SSOEnabled, &o.Status, &o.CreatedAt)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, o)
	}
	return orgs, nil
}

// CreateTeam 创建团队
func (s *Store) CreateTeam(team *Team) error {
	team.ID = uuid.NewString()
	team.CreatedAt = time.Now().UTC()
	_, err := s.db.Exec(`INSERT INTO teams (id, org_id, name, description, default_role, permissions, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		team.ID, team.OrgID, team.Name, team.Description, team.DefaultRole, pq.Array(team.Permissions), team.CreatedAt)
	return err
}

// ListTeams 列出团队
func (s *Store) ListTeams(orgID string) ([]Team, error) {
	rows, err := s.db.Query(`SELECT id, org_id, name, description, default_role, permissions, created_at
		FROM teams WHERE org_id = $1 ORDER BY name`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var t Team
		err := rows.Scan(&t.ID, &t.OrgID, &t.Name, &t.Description, &t.DefaultRole, pq.Array(&t.Permissions), &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, nil
}

// AddMember 添加成员
func (s *Store) AddMember(m *Member) error {
	m.ID = uuid.NewString()
	m.CreatedAt = time.Now().UTC()
	m.Status = "active"
	_, err := s.db.Exec(`INSERT INTO org_members (id, org_id, user_id, role, team_ids, status, invited_by, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		m.ID, m.OrgID, m.UserID, m.Role, pq.Array(m.TeamIDs), m.Status, m.InvitedBy, m.CreatedAt)
	return err
}

// ListMembers 列出成员
func (s *Store) ListMembers(orgID string) ([]Member, error) {
	rows, err := s.db.Query(`SELECT m.id, m.org_id, m.user_id, m.role, m.team_ids, m.status, m.joined_at, m.created_at,
		COALESCE(u.username,'') as username
		FROM org_members m LEFT JOIN users u ON m.user_id = u.id
		WHERE m.org_id = $1 ORDER BY m.created_at`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var m Member
		err := rows.Scan(&m.ID, &m.OrgID, &m.UserID, &m.Role, pq.Array(&m.TeamIDs), &m.Status, &m.JoinedAt, &m.CreatedAt, &m.Username)
		if err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

// RemoveMember 移除成员
func (s *Store) RemoveMember(orgID, userID string) error {
	_, err := s.db.Exec(`DELETE FROM org_members WHERE org_id = $1 AND user_id = $2`, orgID, userID)
	return err
}

// AddIPWhitelist 添加IP白名单
func (s *Store) AddIPWhitelist(w *IPWhitelist) error {
	w.ID = uuid.NewString()
	w.CreatedAt = time.Now().UTC()
	_, err := s.db.Exec(`INSERT INTO ip_whitelist (id, scope_type, scope_id, ip_address, ip_cidr, description, enabled, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		w.ID, w.ScopeType, w.ScopeID, w.IPAddress, w.IPCIDR, w.Description, w.Enabled, w.CreatedAt)
	return err
}

// ListIPWhitelist 列出IP白名单
func (s *Store) ListIPWhitelist(scopeType, scopeID string) ([]IPWhitelist, error) {
	rows, err := s.db.Query(`SELECT id, scope_type, scope_id, ip_address, ip_cidr, description, enabled, created_at
		FROM ip_whitelist WHERE scope_type = $1 AND scope_id = $2 ORDER BY created_at`, scopeType, scopeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []IPWhitelist
	for rows.Next() {
		var w IPWhitelist
		err := rows.Scan(&w.ID, &w.ScopeType, &w.ScopeID, &w.IPAddress, &w.IPCIDR, &w.Description, &w.Enabled, &w.CreatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, w)
	}
	return list, nil
}

// DeleteIPWhitelist 删除IP白名单
func (s *Store) DeleteIPWhitelist(id string) error {
	_, err := s.db.Exec(`DELETE FROM ip_whitelist WHERE id = $1`, id)
	return err
}

// CheckIPAllowed 检查IP是否允许
func (s *Store) CheckIPAllowed(scopeType, scopeID, ipAddress string) (bool, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM ip_whitelist WHERE scope_type = $1 AND scope_id = $2 AND enabled = true`, scopeType, scopeID).Scan(&count)
	if err != nil {
		return true, err // 错误时默认允许
	}
	if count == 0 {
		return true, nil // 没有白名单=允许所有
	}

	// 检查IP是否在白名单中
	var exists int
	err = s.db.QueryRow(`SELECT COUNT(*) FROM ip_whitelist WHERE scope_type = $1 AND scope_id = $2 AND enabled = true
		AND (ip_address = $3 OR $3 <<= ip_cidr::inet)`, scopeType, scopeID, ipAddress).Scan(&exists)
	return exists > 0, err
}
