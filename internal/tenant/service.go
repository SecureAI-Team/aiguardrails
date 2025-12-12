package tenant

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"

	"aiguardrails/internal/types"
)

// Service manages tenants and their applications.
type Service interface {
	CreateTenant(name string) (*types.Tenant, error)
	ListTenants() ([]types.Tenant, error)
	CreateApp(tenantID, name string, quotaPerHr int64) (*types.App, error)
	ListAppsByTenant(tenantID string) ([]types.App, error)
	GetApp(appID string) (*types.App, error)
	RotateSecret(appID string) (*types.App, error)
	RevokeApp(appID string) error
}

// MemoryService is an in-memory implementation suitable for prototyping.
type MemoryService struct {
	mu      sync.RWMutex
	tenants map[string]types.Tenant
	apps    map[string]types.App
}

// NewMemoryService constructs a MemoryService.
func NewMemoryService() *MemoryService {
	return &MemoryService{
		tenants: map[string]types.Tenant{},
		apps:    map[string]types.App{},
	}
}

// CreateTenant registers a new tenant.
func (s *MemoryService) CreateTenant(name string) (*types.Tenant, error) {
	if name == "" {
		return nil, errors.New("tenant name required")
	}
	id := uuid.NewString()
	tenant := types.Tenant{ID: id, Name: name, CreatedAt: time.Now().UTC()}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tenants[id] = tenant
	return &tenant, nil
}

// ListTenants returns all tenants.
func (s *MemoryService) ListTenants() ([]types.Tenant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]types.Tenant, 0, len(s.tenants))
	for _, t := range s.tenants {
		out = append(out, t)
	}
	return out, nil
}

// CreateApp issues a new application credential.
func (s *MemoryService) CreateApp(tenantID, name string, quotaPerHr int64) (*types.App, error) {
	if tenantID == "" {
		return nil, errors.New("tenantID required")
	}
	if name == "" {
		return nil, errors.New("app name required")
	}

	s.mu.RLock()
	_, ok := s.tenants[tenantID]
	s.mu.RUnlock()
	if !ok {
		return nil, errors.New("tenant not found")
	}

	app := types.App{
		ID:         uuid.NewString(),
		TenantID:   tenantID,
		Name:       name,
		APIKey:     uuid.NewString(),
		APISecret:  uuid.NewString(),
		QuotaPerHr: quotaPerHr,
		CreatedAt:  time.Now().UTC(),
		Revoked:    false,
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.apps[app.ID] = app
	return &app, nil
}

// ListAppsByTenant returns apps for a tenant.
func (s *MemoryService) ListAppsByTenant(tenantID string) ([]types.App, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []types.App
	for _, app := range s.apps {
		if app.TenantID == tenantID {
			out = append(out, app)
		}
	}
	return out, nil
}

// GetApp finds an app by ID.
func (s *MemoryService) GetApp(appID string) (*types.App, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	app, ok := s.apps[appID]
	if !ok {
		return nil, errors.New("app not found")
	}
	return &app, nil
}

// RotateSecret issues a new secret for the app.
func (s *MemoryService) RotateSecret(appID string) (*types.App, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	app, ok := s.apps[appID]
	if !ok {
		return nil, errors.New("app not found")
	}
	app.APISecret = uuid.NewString()
	s.apps[appID] = app
	return &app, nil
}

// RevokeApp marks an app as revoked.
func (s *MemoryService) RevokeApp(appID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	app, ok := s.apps[appID]
	if !ok {
		return errors.New("app not found")
	}
	app.Revoked = true
	s.apps[appID] = app
	return nil
}

