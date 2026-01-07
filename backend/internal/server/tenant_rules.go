package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"aiguardrails/internal/policy"
	"aiguardrails/internal/rbac"
)

// registerTenantRulesRoutes 注册租户规则API路由
func (s *Server) registerTenantRulesRoutes(r chi.Router) {
	r.Get("/tenants/{tenantID}/rules", s.listTenantRules)
	r.Post("/tenants/{tenantID}/rules", s.createTenantRule)
	r.Get("/tenants/{tenantID}/rules/{ruleID}", s.getTenantRule)
	r.Put("/tenants/{tenantID}/rules/{ruleID}", s.updateTenantRule)
	r.Delete("/tenants/{tenantID}/rules/{ruleID}", s.deleteTenantRule)
	r.Get("/rules/templates", s.listRuleTemplates)
	r.Get("/platform/rules", s.listPlatformRules)
}

type tenantRuleRequest struct {
	RuleType    string          `json:"rule_type"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Config      json.RawMessage `json:"config"`
	Enabled     bool            `json:"enabled"`
	Priority    int             `json:"priority,omitempty"`
}

func (s *Server) listTenantRules(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantID")
	if !s.allowedTenant(r.Context(), tenantID) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	ruleType := policy.TenantRuleType(r.URL.Query().Get("type"))
	rules, err := s.tenantRuleStore.List(tenantID, ruleType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, rules)
}

func (s *Server) createTenantRule(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantID")
	if !s.allowedTenant(r.Context(), tenantID) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// 检查权限
	role := rbac.RoleFromContext(r.Context())
	if role != rbac.RolePlatformAdmin && role != rbac.RoleTenantAdmin {
		http.Error(w, "tenant admin required", http.StatusForbidden)
		return
	}

	var req tenantRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rule := policy.TenantRule{
		TenantID:    tenantID,
		RuleType:    policy.TenantRuleType(req.RuleType),
		Name:        req.Name,
		Description: req.Description,
		Config:      req.Config,
		Enabled:     req.Enabled,
		Priority:    req.Priority,
	}

	created, err := s.tenantRuleStore.Create(rule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.audit.RecordStore(s.auditStore, "tenant_rule_created", map[string]string{
		"tenant_id": tenantID,
		"rule_id":   created.ID,
		"rule_type": string(created.RuleType),
		"name":      created.Name,
	})
	s.writeJSON(w, http.StatusCreated, created)
}

func (s *Server) getTenantRule(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantID")
	ruleID := chi.URLParam(r, "ruleID")
	if !s.allowedTenant(r.Context(), tenantID) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	rule, err := s.tenantRuleStore.Get(tenantID, ruleID)
	if err != nil {
		if err == policy.ErrRuleNotFound {
			http.Error(w, "rule not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, rule)
}

func (s *Server) updateTenantRule(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantID")
	ruleID := chi.URLParam(r, "ruleID")
	if !s.allowedTenant(r.Context(), tenantID) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	role := rbac.RoleFromContext(r.Context())
	if role != rbac.RolePlatformAdmin && role != rbac.RoleTenantAdmin {
		http.Error(w, "tenant admin required", http.StatusForbidden)
		return
	}

	var req tenantRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rule := policy.TenantRule{
		ID:          ruleID,
		TenantID:    tenantID,
		RuleType:    policy.TenantRuleType(req.RuleType),
		Name:        req.Name,
		Description: req.Description,
		Config:      req.Config,
		Enabled:     req.Enabled,
		Priority:    req.Priority,
	}

	updated, err := s.tenantRuleStore.Update(rule)
	if err != nil {
		if err == policy.ErrRuleNotFound {
			http.Error(w, "rule not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.audit.RecordStore(s.auditStore, "tenant_rule_updated", map[string]string{
		"tenant_id": tenantID,
		"rule_id":   ruleID,
		"name":      updated.Name,
	})
	s.writeJSON(w, http.StatusOK, updated)
}

func (s *Server) deleteTenantRule(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantID")
	ruleID := chi.URLParam(r, "ruleID")
	if !s.allowedTenant(r.Context(), tenantID) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	role := rbac.RoleFromContext(r.Context())
	if role != rbac.RolePlatformAdmin && role != rbac.RoleTenantAdmin {
		http.Error(w, "tenant admin required", http.StatusForbidden)
		return
	}

	if err := s.tenantRuleStore.Delete(tenantID, ruleID); err != nil {
		if err == policy.ErrRuleNotFound {
			http.Error(w, "rule not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.audit.RecordStore(s.auditStore, "tenant_rule_deleted", map[string]string{
		"tenant_id": tenantID,
		"rule_id":   ruleID,
	})
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (s *Server) listRuleTemplates(w http.ResponseWriter, r *http.Request) {
	ruleType := policy.TenantRuleType(r.URL.Query().Get("type"))
	templates, err := s.tenantRuleStore.ListTemplates(ruleType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, templates)
}

func (s *Server) listPlatformRules(w http.ResponseWriter, r *http.Request) {
	// 返回平台级法规规则（只读）
	if s.rulesRepo == nil {
		s.writeJSON(w, http.StatusOK, []policy.Rule{})
		return
	}
	rules := s.rulesRepo.List(map[string]string{})
	s.writeJSON(w, http.StatusOK, rules)
}
