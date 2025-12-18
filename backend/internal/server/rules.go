package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"aiguardrails/internal/policy"
)

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

type rulesHandler struct {
	repo      *policy.RulesRepository
	ruleStore *policy.RuleStore
}

func (s *Server) registerRulesRoutes(r chi.Router) {
	h := rulesHandler{repo: s.rulesRepo, ruleStore: s.ruleStore}
	r.Get("/rules", h.listRules)
	r.Post("/policies/{policyID}/rules/{ruleID}", h.attachRule)
	r.Get("/policies/{policyID}/rules", h.listPolicyRules)
}

func (h rulesHandler) listRules(w http.ResponseWriter, r *http.Request) {
	filter := map[string]string{
		"jurisdiction": r.URL.Query().Get("jurisdiction"),
		"regulation":   r.URL.Query().Get("regulation"),
		"vendor":       r.URL.Query().Get("vendor"),
		"product":      r.URL.Query().Get("product"),
	}
	rules := h.repo.List(filter)
	writeJSON(w, http.StatusOK, rules)
}

func (h rulesHandler) attachRule(w http.ResponseWriter, r *http.Request) {
	policyID := chi.URLParam(r, "policyID")
	ruleID := chi.URLParam(r, "ruleID")
	if policyID == "" || ruleID == "" {
		http.Error(w, "missing ids", http.StatusBadRequest)
		return
	}
	if _, err := h.repo.Get(ruleID); err != nil {
		http.Error(w, "rule not found", http.StatusNotFound)
		return
	}
	if err := h.ruleStore.Attach(policyID, ruleID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "attached"})
}

func (h rulesHandler) listPolicyRules(w http.ResponseWriter, r *http.Request) {
	policyID := chi.URLParam(r, "policyID")
	if policyID == "" {
		http.Error(w, "missing policyID", http.StatusBadRequest)
		return
	}
	ids, err := h.ruleStore.ListByPolicy(policyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var rules []policy.Rule
	for _, id := range ids {
		if rule, err := h.repo.Get(id); err == nil {
			rules = append(rules, *rule)
		}
	}
	writeJSON(w, http.StatusOK, rules)
}

