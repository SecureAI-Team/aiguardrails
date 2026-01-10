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

	s.writeJSON(w, http.StatusCreated, req)
}

func (s *Server) deleteRule(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.ruleStore.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
