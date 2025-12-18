package tenant

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"aiguardrails/internal/types"
)

// PGService implements Service backed by Postgres.
type PGService struct {
	db *sql.DB
}

// NewPGService constructs PGService.
func NewPGService(db *sql.DB) *PGService {
	return &PGService{db: db}
}

func (s *PGService) CreateTenant(name string) (*types.Tenant, error) {
	if name == "" {
		return nil, errors.New("tenant name required")
	}
	id := uuid.New()
	now := time.Now().UTC()
	_, err := s.db.Exec(`INSERT INTO tenants (id, name, created_at) VALUES ($1,$2,$3)`, id, name, now)
	if err != nil {
		return nil, err
	}
	return &types.Tenant{ID: id.String(), Name: name, CreatedAt: now}, nil
}

func (s *PGService) ListTenants() ([]types.Tenant, error) {
	rows, err := s.db.Query(`SELECT id, name, created_at FROM tenants ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []types.Tenant
	for rows.Next() {
		var t types.Tenant
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func (s *PGService) CreateApp(tenantID, name string, quotaPerHr int64) (*types.App, error) {
	if tenantID == "" || name == "" {
		return nil, errors.New("tenantID and name required")
	}
	id := uuid.New()
	apiKey := uuid.New()
	apiSecret := uuid.New()
	now := time.Now().UTC()
	_, err := s.db.Exec(`INSERT INTO apps (id, tenant_id, name, api_key, api_secret, quota_per_hr, created_at, revoked)
		VALUES ($1,$2,$3,$4,$5,$6,$7,false)`, id, tenantID, name, apiKey, apiSecret, quotaPerHr, now)
	if err != nil {
		return nil, err
	}
	return &types.App{
		ID:         id.String(),
		TenantID:   tenantID,
		Name:       name,
		APIKey:     apiKey.String(),
		APISecret:  apiSecret.String(),
		QuotaPerHr: quotaPerHr,
		CreatedAt:  now,
		Revoked:    false,
	}, nil
}

func (s *PGService) ListAppsByTenant(tenantID string) ([]types.App, error) {
	rows, err := s.db.Query(`SELECT id, tenant_id, name, api_key, api_secret, quota_per_hr, created_at, revoked FROM apps WHERE tenant_id=$1 ORDER BY created_at DESC`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []types.App
	for rows.Next() {
		var a types.App
		if err := rows.Scan(&a.ID, &a.TenantID, &a.Name, &a.APIKey, &a.APISecret, &a.QuotaPerHr, &a.CreatedAt, &a.Revoked); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

func (s *PGService) GetApp(appID string) (*types.App, error) {
	var a types.App
	err := s.db.QueryRow(`SELECT id, tenant_id, name, api_key, api_secret, quota_per_hr, created_at, revoked FROM apps WHERE id=$1`, appID).
		Scan(&a.ID, &a.TenantID, &a.Name, &a.APIKey, &a.APISecret, &a.QuotaPerHr, &a.CreatedAt, &a.Revoked)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("app not found")
		}
		return nil, err
	}
	return &a, nil
}

func (s *PGService) RotateSecret(appID string) (*types.App, error) {
	newSecret := uuid.New()
	_, err := s.db.Exec(`UPDATE apps SET api_secret=$1 WHERE id=$2`, newSecret, appID)
	if err != nil {
		return nil, err
	}
	return s.GetApp(appID)
}

func (s *PGService) RevokeApp(appID string) error {
	res, err := s.db.Exec(`UPDATE apps SET revoked=true WHERE id=$1`, appID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("app not found")
	}
	return nil
}

// Ensure interface compliance.
var _ Service = (*PGService)(nil)

