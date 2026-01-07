package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"aiguardrails/internal/rbac"
)

// registerUserRoutes registers user management API routes.
func (s *Server) registerUserRoutes(r chi.Router) {
	// Platform admin: full user management
	r.With(rbac.RequirePerm(rbac.PermManageTenants)).Get("/users", s.listUsers)
	r.With(rbac.RequirePerm(rbac.PermManageTenants)).Post("/users", s.createUser)
	r.With(rbac.RequirePerm(rbac.PermManageTenants)).Get("/users/{userID}", s.getUser)
	r.With(rbac.RequirePerm(rbac.PermManageTenants)).Put("/users/{userID}", s.updateUser)
	r.With(rbac.RequirePerm(rbac.PermManageTenants)).Delete("/users/{userID}", s.deleteUser)
	r.With(rbac.RequirePerm(rbac.PermManageTenants)).Post("/users/{userID}/password", s.resetPassword)

	// Tenant user management
	r.With(rbac.RequirePerm(rbac.PermManageApps)).Get("/tenants/{tenantID}/users", s.listTenantUsers)
	r.With(rbac.RequirePerm(rbac.PermManageApps)).Post("/tenants/{tenantID}/users", s.addTenantUser)
	r.With(rbac.RequirePerm(rbac.PermManageApps)).Delete("/tenants/{tenantID}/users/{userID}", s.removeTenantUser)
	r.With(rbac.RequirePerm(rbac.PermManageApps)).Put("/tenants/{tenantID}/users/{userID}/role", s.updateTenantUserRole)
}

type createUserRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type updateUserRequest struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
}

type resetPasswordRequest struct {
	NewPassword string `json:"new_password"`
}

type addTenantUserRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type updateRoleRequest struct {
	Role string `json:"role"`
}

func (s *Server) listUsers(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")
	status := r.URL.Query().Get("status")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 50
	}

	users, err := s.userStore.List(role, status, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, users)
}

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.userStore.Create(req.Username, req.Password, req.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update additional fields
	if req.Email != "" || req.DisplayName != "" {
		_ = s.userStore.Update(user.ID, req.Role, req.Email, req.DisplayName, "active")
	}

	s.audit.RecordStore(s.auditStore, "user_created", map[string]string{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
	user.PasswordHash = ""
	s.writeJSON(w, http.StatusCreated, user)
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	user, err := s.userStore.GetByID(userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	s.writeJSON(w, http.StatusOK, user)
}

func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.userStore.Update(userID, req.Role, req.Email, req.DisplayName, req.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.audit.RecordStore(s.auditStore, "user_updated", map[string]string{"user_id": userID})
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if err := s.userStore.Delete(userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.audit.RecordStore(s.auditStore, "user_deleted", map[string]string{"user_id": userID})
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (s *Server) resetPassword(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	var req resetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.NewPassword == "" {
		http.Error(w, "new_password required", http.StatusBadRequest)
		return
	}

	if err := s.userStore.UpdatePassword(userID, req.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.audit.RecordStore(s.auditStore, "password_reset", map[string]string{"user_id": userID})
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "password_reset"})
}

func (s *Server) listTenantUsers(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantID")
	if !s.allowedTenant(r.Context(), tenantID) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	users, err := s.tenantUserStore.ListByTenant(tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, users)
}

func (s *Server) addTenantUser(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantID")
	if !s.allowedTenant(r.Context(), tenantID) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var req addTenantUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tu, err := s.tenantUserStore.Add(tenantID, req.UserID, req.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.audit.RecordStore(s.auditStore, "tenant_user_added", map[string]string{
		"tenant_id": tenantID,
		"user_id":   req.UserID,
		"role":      req.Role,
	})
	s.writeJSON(w, http.StatusCreated, tu)
}

func (s *Server) removeTenantUser(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantID")
	userID := chi.URLParam(r, "userID")
	if !s.allowedTenant(r.Context(), tenantID) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := s.tenantUserStore.Remove(tenantID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.audit.RecordStore(s.auditStore, "tenant_user_removed", map[string]string{
		"tenant_id": tenantID,
		"user_id":   userID,
	})
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "removed"})
}

func (s *Server) updateTenantUserRole(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantID")
	userID := chi.URLParam(r, "userID")
	if !s.allowedTenant(r.Context(), tenantID) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var req updateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.tenantUserStore.UpdateRole(tenantID, userID, req.Role); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.audit.RecordStore(s.auditStore, "tenant_user_role_updated", map[string]string{
		"tenant_id": tenantID,
		"user_id":   userID,
		"role":      req.Role,
	})
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "role_updated"})
}
