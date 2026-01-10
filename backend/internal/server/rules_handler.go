package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"aiguardrails/internal/rules"
)

func (s *Server) listRules(w http.ResponseWriter, r *http.Request) {
	list, err := s.ruleStore.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, list)
}

func (s *Server) createRule(w http.ResponseWriter, r *http.Request) {
	var req rules.Rule
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.ID == "" {
		req.ID = uuid.NewString()
	}
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	req.IsSystem = false

	if err := s.ruleStore.Add(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If rule is OPA, we might need to trigger reload?
	// For now just save metadata.
	if req.Type == rules.RuleTypeOPA {
		s.syncOPARules()
	}

	s.writeJSON(w, http.StatusCreated, req)
}

func (s *Server) updateRule(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req rules.Rule
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure ID matches
	if req.ID != "" && req.ID != id {
		http.Error(w, "id mismatch", http.StatusBadRequest)
		return
	}
	req.ID = id
	req.UpdatedAt = time.Now()

	// Preserve CreatedAt if possible, or fetch old rule first.
	// For MemoryStore, Add() creates or overwrites.
	// But we should probably check if it exists or use Update?
	// The RuleStore interface has Add() and Get(). It might not have Update().
	// Inspecting RuleStore interface...
	// Assuming Add() overwrites.

	// Better: Get old rule to preserve CreatedAt
	old, err := s.ruleStore.Get(id)
	if err == nil {
		req.CreatedAt = old.CreatedAt
	} else {
		// Rule not found, treat as create? Or error?
		// Usually update should fail if not found.
		http.Error(w, "rule not found", http.StatusNotFound)
		return
	}

	if err := s.ruleStore.Add(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Type == rules.RuleTypeOPA {
		s.syncOPARules()
	}

	s.writeJSON(w, http.StatusOK, req)
}

func (s *Server) deleteRule(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.ruleStore.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.syncOPARules()
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
