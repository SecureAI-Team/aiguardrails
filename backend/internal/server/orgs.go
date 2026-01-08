package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"aiguardrails/internal/org"
)

// registerOrgRoutes 注册组织管理路由
func (s *Server) registerOrgRoutes(r chi.Router) {
	// 组织
	r.Get("/orgs", s.listOrgs)
	r.Post("/orgs", s.createOrg)
	r.Get("/orgs/{id}", s.getOrg)

	// 团队
	r.Get("/orgs/{orgId}/teams", s.listTeams)
	r.Post("/orgs/{orgId}/teams", s.createTeam)

	// 成员
	r.Get("/orgs/{orgId}/members", s.listMembers)
	r.Post("/orgs/{orgId}/members", s.addMember)
	r.Delete("/orgs/{orgId}/members/{userId}", s.removeMember)

	// IP白名单
	r.Get("/orgs/{orgId}/whitelist", s.listIPWhitelist)
	r.Post("/orgs/{orgId}/whitelist", s.addIPWhitelist)
	r.Delete("/whitelist/{id}", s.deleteIPWhitelist)
}

func (s *Server) listOrgs(w http.ResponseWriter, r *http.Request) {
	orgs, err := s.orgStore.ListOrgs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, orgs)
}

func (s *Server) createOrg(w http.ResponseWriter, r *http.Request) {
	var o org.Organization
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.orgStore.CreateOrg(&o); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusCreated, o)
}

func (s *Server) getOrg(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	o, err := s.orgStore.GetOrg(id)
	if err != nil {
		http.Error(w, "org not found", http.StatusNotFound)
		return
	}
	s.writeJSON(w, http.StatusOK, o)
}

func (s *Server) listTeams(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	teams, err := s.orgStore.ListTeams(orgID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, teams)
}

func (s *Server) createTeam(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	var t org.Team
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.OrgID = orgID
	if err := s.orgStore.CreateTeam(&t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusCreated, t)
}

func (s *Server) listMembers(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	members, err := s.orgStore.ListMembers(orgID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, members)
}

func (s *Server) addMember(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	var m org.Member
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m.OrgID = orgID
	if err := s.orgStore.AddMember(&m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusCreated, m)
}

func (s *Server) removeMember(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	userID := chi.URLParam(r, "userId")
	if err := s.orgStore.RemoveMember(orgID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listIPWhitelist(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	list, err := s.orgStore.ListIPWhitelist("org", orgID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, list)
}

func (s *Server) addIPWhitelist(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	var wl org.IPWhitelist
	if err := json.NewDecoder(r.Body).Decode(&wl); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	wl.ScopeType = "org"
	wl.ScopeID = orgID
	if err := s.orgStore.AddIPWhitelist(&wl); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusCreated, wl)
}

func (s *Server) deleteIPWhitelist(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.orgStore.DeleteIPWhitelist(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
