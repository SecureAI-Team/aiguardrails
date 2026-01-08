package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"aiguardrails/internal/agent"
	"aiguardrails/internal/alert"
	"aiguardrails/internal/audit"
	"aiguardrails/internal/auth"
	"aiguardrails/internal/config"
	"aiguardrails/internal/mcp"
	"aiguardrails/internal/opa"
	"aiguardrails/internal/org"
	"aiguardrails/internal/policy"
	"aiguardrails/internal/promptfw"
	"aiguardrails/internal/rag"
	"aiguardrails/internal/rbac"
	"aiguardrails/internal/tenant"
	"aiguardrails/internal/tracing"
	"aiguardrails/internal/types"
	"aiguardrails/internal/usage"
)

// Server wires HTTP routes to services.
type Server struct {
	cfg             config.Config
	router          *chi.Mux
	tenant          tenant.Service
	policy          policy.Engine
	firewall        *promptfw.Firewall
	agent           *agent.Gateway
	rag             *rag.Security
	usage           *usage.Meter
	rate            *usage.RateLimiter
	audit           *audit.Logger
	auditStore      *audit.Store
	mcp             *mcp.Broker
	capStore        *mcp.Store
	rulesRepo       *policy.RulesRepository
	ruleStore       *policy.RuleStore
	tenantRuleStore *policy.TenantRuleStore
	userStore       *auth.UserStore
	tenantUserStore *auth.TenantUserStore
	socialAuth      *auth.SocialAuthStore
	smsStore        *auth.SMSStore
	jwtSigner       *auth.JWTSigner
	opaEval         *opa.Evaluator
	alertStore      *alert.RuleStore
	usageStore      *usage.UsageStore
	tracingStore    *tracing.Store
	orgStore        *org.Store
}

type ctxKey string

const authRoleCtxKey ctxKey = "role"

// New builds a Server with dependencies.
func New(cfg config.Config, tenantSvc tenant.Service, policyEng policy.Engine, firewall *promptfw.Firewall, agentGw *agent.Gateway, ragSec *rag.Security, usageMeter *usage.Meter, rateLimiter *usage.RateLimiter, auditLog *audit.Logger, auditStore *audit.Store, mcpBroker *mcp.Broker, capStore *mcp.Store, rulesRepo *policy.RulesRepository, ruleStore *policy.RuleStore, tenantRuleStore *policy.TenantRuleStore, userStore *auth.UserStore, tenantUserStore *auth.TenantUserStore, jwtSigner *auth.JWTSigner, opaEval *opa.Evaluator, alertStore *alert.RuleStore, usageStore *usage.UsageStore, tracingStore *tracing.Store, orgStore *org.Store) *Server {
	s := &Server{
		cfg:             cfg,
		router:          chi.NewRouter(),
		tenant:          tenantSvc,
		policy:          policyEng,
		firewall:        firewall,
		agent:           agentGw,
		rag:             ragSec,
		usage:           usageMeter,
		rate:            rateLimiter,
		audit:           auditLog,
		auditStore:      auditStore,
		mcp:             mcpBroker,
		capStore:        capStore,
		rulesRepo:       rulesRepo,
		ruleStore:       ruleStore,
		tenantRuleStore: tenantRuleStore,
		userStore:       userStore,
		tenantUserStore: tenantUserStore,
		jwtSigner:       jwtSigner,
		opaEval:         opaEval,
		alertStore:      alertStore,
		usageStore:      usageStore,
		tracingStore:    tracingStore,
		orgStore:        orgStore,
	}
	s.routes()
	return s
}

// Handler returns the configured router.
func (s *Server) Handler() http.Handler {
	return s.router
}

func cors(origins []string) func(http.Handler) http.Handler {
	allowed := map[string]struct{}{}
	for _, o := range origins {
		allowed[o] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				if len(allowed) == 0 {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					if _, ok := allowed["*"]; ok {
						w.Header().Set("Access-Control-Allow-Origin", "*")
					} else if _, ok := allowed[origin]; ok {
						w.Header().Set("Access-Control-Allow-Origin", origin)
					}
				}
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-App-Id, X-App-Secret")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (s *Server) routes() {
	r := s.router
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors(s.cfg.AllowedOrigins))
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	r.Post("/v1/auth/login", s.login)
	// Social auth routes (optional)
	if s.socialAuth != nil {
		r.Route("/v1", s.registerSocialAuthRoutes)
	}

	r.Route("/v1", func(r chi.Router) {
		// Admin-scoped endpoints (platform admin via token or local login)
		r.Group(func(r chi.Router) {
			r.Use(s.adminAuth())
			r.Post("/tenants", s.createTenant)
			r.Get("/tenants", s.listTenants)

			r.Post("/tenants/{tenantID}/apps", s.createApp)
			r.Get("/tenants/{tenantID}/apps", s.listApps)
			r.Post("/apps/{appID}/rotate", s.rotateApp)
			r.Post("/apps/{appID}/revoke", s.revokeApp)

			r.Post("/tenants/{tenantID}/policies", s.createPolicy)
			r.Put("/tenants/{tenantID}/policies/{policyID}", s.updatePolicy)
			r.Get("/tenants/{tenantID}/policies", s.listPolicies)

			r.Post("/capabilities", s.createCapability)
			r.Get("/audit", s.listAudit)
			r.Get("/tenants/{tenantID}/policies/history", s.listPolicyHistory)

			// User management routes
			if s.userStore != nil {
				s.registerUserRoutes(r)
			}

			// Alert routes
			if s.alertStore != nil {
				s.registerAlertRoutes(r)
			}

			// Usage/Stats routes
			if s.usageStore != nil {
				s.registerUsageRoutes(r)
			}

			// Tracing routes
			if s.tracingStore != nil {
				s.registerTracingRoutes(r)
			}

			// Organization routes
			if s.orgStore != nil {
				s.registerOrgRoutes(r)
			}
		})

		// Admin via OIDC/JWT (optional, only if OIDC is configured)
		if s.cfg.OidcJWKSURL != "" {
			r.Group(func(r chi.Router) {
				r.Use(auth.OIDCMiddleware(s.cfg.OidcJWKSURL, s.cfg.OIDCIssuer, s.cfg.OIDCAudience, "tenant_id", s.cfg.OIDCAdminRole, s.cfg.OIDCUserRole, s.cfg.OIDCTimeoutSec, s.cfg.OIDCCacheMin))
				// role comes from token; enforce permission per route
				r.With(rbac.RequirePerm(rbac.PermManageApps)).Post("/tenants/{tenantID}/apps", s.createApp)
				r.With(rbac.RequirePerm(rbac.PermManageApps)).Get("/tenants/{tenantID}/apps", s.listApps)
				r.With(rbac.RequirePerm(rbac.PermManageApps)).Post("/apps/{appID}/rotate", s.rotateApp)
				r.With(rbac.RequirePerm(rbac.PermManageApps)).Post("/apps/{appID}/revoke", s.revokeApp)
				r.With(rbac.RequirePerm(rbac.PermManagePolicy)).Post("/tenants/{tenantID}/policies", s.createPolicy)
				r.With(rbac.RequirePerm(rbac.PermManagePolicy)).Put("/tenants/{tenantID}/policies/{policyID}", s.updatePolicy)
				r.With(rbac.RequirePerm(rbac.PermManagePolicy)).Get("/tenants/{tenantID}/policies", s.listPolicies)
				r.With(rbac.RequirePerm(rbac.PermManagePolicy)).Get("/tenants/{tenantID}/policies/history", s.listPolicyHistory)
				r.With(rbac.RequirePerm(rbac.PermManageApps)).Post("/capabilities", s.createCapability)
				r.With(rbac.RequirePerm(rbac.PermViewLogs)).Get("/audit", s.listAudit)
				if s.rulesRepo != nil && s.ruleStore != nil {
					s.registerRulesRoutes(r)
				}
				if s.tenantRuleStore != nil {
					s.registerTenantRulesRoutes(r)
				}
				// User management routes
				if s.userStore != nil {
					s.registerUserRoutes(r)
				}
			})
		}

		// App-scoped endpoints require app credentials.
		r.Group(func(r chi.Router) {
			r.Use(auth.APIKeyMiddleware(s.tenant))
			r.Use(rbac.WithRole(rbac.RoleTenantUser))
			r.Post("/guardrails/prompt-check", s.checkPrompt)
			r.Post("/guardrails/output-filter", s.checkOutput)
			r.Post("/agent/plan", s.planAndAct)
			r.Get("/mcp/capabilities", s.listCapabilities)
		})
	})
}

type tenantRequest struct {
	Name string `json:"name"`
}

func (s *Server) createTenant(w http.ResponseWriter, r *http.Request) {
	var req tenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tenant, err := s.tenant.CreateTenant(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.writeJSON(w, http.StatusCreated, tenant)
}

func (s *Server) listTenants(w http.ResponseWriter, r *http.Request) {
	tenants, err := s.tenant.ListTenants()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, tenants)
}

type appRequest struct {
	Name       string `json:"name"`
	QuotaPerHr int64  `json:"quota_per_hr"`
}

func (s *Server) createApp(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantID")
	if !s.allowedTenant(r.Context(), tenantID) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var req appRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	app, err := s.tenant.CreateApp(tenantID, req.Name, req.QuotaPerHr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.audit.RecordStore(s.auditStore, "app_created", map[string]string{"app_id": app.ID, "tenant_id": tenantID})
	s.writeJSON(w, http.StatusCreated, app)
}

func (s *Server) listApps(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantID")
	if !s.allowedTenant(r.Context(), tenantID) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	apps, err := s.tenant.ListAppsByTenant(tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, apps)
}

func (s *Server) rotateApp(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")
	app, err := s.tenant.RotateSecret(appID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.audit.RecordStore(s.auditStore, "app_rotated", map[string]string{"app_id": appID})
	s.writeJSON(w, http.StatusOK, app)
}

func (s *Server) revokeApp(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")
	if err := s.tenant.RevokeApp(appID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.audit.RecordStore(s.auditStore, "app_revoked", map[string]string{"app_id": appID})
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "revoked"})
}

type policyRequest struct {
	Name           string   `json:"name"`
	PromptRules    []string `json:"prompt_rules"`
	ToolAllowList  []string `json:"tool_allowlist"`
	RAGNamespaces  []string `json:"rag_namespaces"`
	OutputFilters  []string `json:"output_filters"`
	SensitiveTerms []string `json:"sensitive_terms"`
}

func (s *Server) createPolicy(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantID")
	if !s.allowedTenant(r.Context(), tenantID) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var req policyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p, err := s.policy.CreatePolicy(types.Policy{
		TenantID:       tenantID,
		Name:           req.Name,
		PromptRules:    req.PromptRules,
		ToolAllowList:  req.ToolAllowList,
		RAGNamespaces:  req.RAGNamespaces,
		OutputFilters:  req.OutputFilters,
		SensitiveTerms: req.SensitiveTerms,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.audit.RecordStore(s.auditStore, "policy_created", map[string]string{"tenant_id": tenantID, "policy_id": p.ID})
	s.writeJSON(w, http.StatusCreated, p)
}

func (s *Server) updatePolicy(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantID")
	policyID := chi.URLParam(r, "policyID")
	if !s.allowedTenant(r.Context(), tenantID) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	oldPolicy, _ := s.policy.GetPolicy(tenantID, policyID)
	var req policyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p, err := s.policy.UpdatePolicy(types.Policy{
		ID:             policyID,
		TenantID:       tenantID,
		Name:           req.Name,
		PromptRules:    req.PromptRules,
		ToolAllowList:  req.ToolAllowList,
		RAGNamespaces:  req.RAGNamespaces,
		OutputFilters:  req.OutputFilters,
		SensitiveTerms: req.SensitiveTerms,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	diff := summarizePolicyDiff(oldPolicy, &p)
	s.audit.RecordStore(s.auditStore, "policy_updated", map[string]string{"tenant_id": tenantID, "policy_id": p.ID, "diff": diff})
	s.writeJSON(w, http.StatusOK, p)
}

func (s *Server) listPolicies(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantID")
	if !s.allowedTenant(r.Context(), tenantID) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	policies, err := s.policy.ListPolicies(tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, policies)
}

func (s *Server) listPolicyHistory(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantID")
	if !s.allowedTenant(r.Context(), tenantID) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	history, err := s.policy.ListHistory(tenantID, 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, history)
}

func (s *Server) listAudit(w http.ResponseWriter, r *http.Request) {
	if s.auditStore == nil {
		http.Error(w, "audit store unavailable", http.StatusInternalServerError)
		return
	}
	limit := 100
	if q := r.URL.Query().Get("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n > 0 && n <= 500 {
			limit = n
		}
	}
	eventLike := r.URL.Query().Get("event")
	tenantID := r.URL.Query().Get("tenant_id")
	events, err := s.auditStore.List(limit, eventLike, tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, events)
}

type promptCheckRequest struct {
	TenantID string `json:"tenant_id"`
	Prompt   string `json:"prompt"`
}

func (s *Server) checkPrompt(w http.ResponseWriter, r *http.Request) {
	var req promptCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tenantID := req.TenantID
	if tenantID == "" {
		tenantID = auth.TenantIDFromContext(r.Context())
	}
	// OPA check if enabled
	if s.opaEval != nil {
		allow, data, err := s.opaEval.Decide(r.Context(), opa.Input{
			TenantID: tenantID,
			AppID:    auth.AppIDFromContext(r.Context()),
			Mode:     "prompt_check",
			Prompt:   req.Prompt,
		})
		if err == nil && !allow {
			s.writeJSON(w, http.StatusOK, types.GuardrailResult{Allowed: false, Reason: "opa_block", Signals: []string{fmt.Sprint(data)}})
			return
		}
	}
	result := s.firewall.CheckPrompt(tenantID, req.Prompt)
	s.writeJSON(w, http.StatusOK, result)
}

type outputCheckRequest struct {
	Output   string `json:"output"`
	TenantID string `json:"tenant_id"`
}

func (s *Server) checkOutput(w http.ResponseWriter, r *http.Request) {
	var req outputCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tenantID := req.TenantID
	if tenantID == "" {
		tenantID = auth.TenantIDFromContext(r.Context())
	}
	if s.opaEval != nil {
		allow, data, err := s.opaEval.Decide(r.Context(), opa.Input{
			TenantID: tenantID,
			AppID:    auth.AppIDFromContext(r.Context()),
			Mode:     "output_filter",
			Output:   req.Output,
		})
		if err == nil && !allow {
			s.writeJSON(w, http.StatusOK, types.GuardrailResult{Allowed: false, Reason: "opa_block", Signals: []string{fmt.Sprint(data)}})
			return
		}
	}
	result := s.firewall.FilterOutput(tenantID, req.Output)
	s.writeJSON(w, http.StatusOK, result)
}

type planRequest struct {
	TenantID string   `json:"tenant_id"`
	Prompt   string   `json:"prompt"`
	Tools    []string `json:"tools"`
}

func (s *Server) planAndAct(w http.ResponseWriter, r *http.Request) {
	var req planRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tenantID := req.TenantID
	if tenantID == "" {
		tenantID = auth.TenantIDFromContext(r.Context())
	}
	appID := auth.AppIDFromContext(r.Context())
	if appID != "" {
		app, err := s.tenant.GetApp(appID)
		if err == nil && s.rate != nil {
			ok, count, err := s.rate.Allow(r.Context(), appID, app.QuotaPerHr)
			if err != nil {
				http.Error(w, "rate check error", http.StatusInternalServerError)
				return
			}
			if !ok {
				http.Error(w, "quota exceeded", http.StatusTooManyRequests)
				return
			}
			s.audit.RecordStore(s.auditStore, "usage_record", map[string]string{"app_id": appID, "count": fmt.Sprintf("%d", count)})
		}
	}
	result, err := s.agent.LegacyPlanAndAct(tenantID, req.Prompt, req.Tools)
	if err != nil {
		http.Error(w, result.Reason, http.StatusForbidden)
		return
	}
	if appID != "" {
		s.usage.Record(appID, 1)
	}
	s.writeJSON(w, http.StatusOK, result)
}

// listCapabilities returns registry entries.
func (s *Server) listCapabilities(w http.ResponseWriter, r *http.Request) {
	if s.capStore == nil {
		http.Error(w, "capabilities unavailable", http.StatusInternalServerError)
		return
	}
	tag := r.URL.Query().Get("tag")
	all, err := s.capStore.List(tag)
	if err != nil {
		http.Error(w, "capabilities unavailable", http.StatusInternalServerError)
		return
	}
	tenantID := auth.TenantIDFromContext(r.Context())
	if tenantID != "" {
		all = s.capStore.FilterAllowed(all, s.policyAllowList(tenantID))
	}
	s.writeJSON(w, http.StatusOK, all)
}

func (s *Server) policyAllowList(tenantID string) []string {
	policies, err := s.policy.ListPolicies(tenantID)
	if err != nil {
		return nil
	}
	set := map[string]struct{}{}
	for _, p := range policies {
		for _, tool := range p.ToolAllowList {
			set[tool] = struct{}{}
		}
	}
	out := make([]string, 0, len(set))
	for t := range set {
		out = append(out, t)
	}
	return out
}

func (s *Server) allowedTenant(ctx context.Context, tenantID string) bool {
	if tenantID == "" {
		return true
	}
	role := rbac.RoleFromContext(ctx)
	ctxTenant := auth.TenantIDFromContext(ctx)
	if role == rbac.RolePlatformAdmin {
		return true
	}
	if role == rbac.RoleTenantAdmin && ctxTenant == tenantID {
		return true
	}
	return false
}

// adminAuth allows either admin token header or JWT with admin role.
func (s *Server) adminAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Admin token path
			if r.Header.Get("X-Admin-Token") == s.cfg.AdminToken {
				next.ServeHTTP(w, r)
				return
			}
			// JWT path
			authz := r.Header.Get("Authorization")
			if authz == "" || !strings.HasPrefix(authz, "Bearer ") || s.jwtSigner == nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authz, "Bearer ")
			claims, err := s.jwtSigner.Parse(tokenStr)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if claims.Role != rbac.RolePlatformAdmin && claims.Role != rbac.RoleTenantAdmin {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			ctx := context.WithValue(r.Context(), authRoleCtxKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func summarizePolicyDiff(oldP, newP *types.Policy) string {
	if oldP == nil || newP == nil {
		return "policy_created_or_replaced"
	}
	changes := []string{}
	if oldP.Name != newP.Name {
		changes = append(changes, "name")
	}
	if strings.Join(oldP.PromptRules, ",") != strings.Join(newP.PromptRules, ",") {
		changes = append(changes, "prompt_rules")
	}
	if strings.Join(oldP.ToolAllowList, ",") != strings.Join(newP.ToolAllowList, ",") {
		changes = append(changes, "tool_allowlist")
	}
	if strings.Join(oldP.RAGNamespaces, ",") != strings.Join(newP.RAGNamespaces, ",") {
		changes = append(changes, "rag_namespaces")
	}
	if strings.Join(oldP.OutputFilters, ",") != strings.Join(newP.OutputFilters, ",") {
		changes = append(changes, "output_filters")
	}
	if strings.Join(oldP.SensitiveTerms, ",") != strings.Join(newP.SensitiveTerms, ",") {
		changes = append(changes, "sensitive_terms")
	}
	if len(changes) == 0 {
		return "no_changes"
	}
	return strings.Join(changes, "|")
}

type capabilityRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func (s *Server) createCapability(w http.ResponseWriter, r *http.Request) {
	if s.capStore == nil {
		http.Error(w, "cap store unavailable", http.StatusInternalServerError)
		return
	}
	var req capabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c, err := s.capStore.Add(mcp.Capability{
		Name:        req.Name,
		Description: req.Description,
		Tags:        req.Tags,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.audit.RecordStore(s.auditStore, "capability_created", map[string]string{"name": c.Name})
	s.writeJSON(w, http.StatusCreated, c)
}

func (s *Server) writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// Addr returns server address.
func (s *Server) Addr() string {
	return ":" + s.cfg.HTTPPort
}
