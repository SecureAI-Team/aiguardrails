package server

import (
	"encoding/json"
	"net/http"
	"sync"
)

type SettingsStore struct {
	mu       sync.RWMutex
	settings map[string]string
}

func NewSettingsStore() *SettingsStore {
	return &SettingsStore{
		settings: make(map[string]string),
	}
}

func (s *SettingsStore) Get(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.settings[key]
}

func (s *SettingsStore) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.settings[key] = value
}

func (s *SettingsStore) GetAll() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := make(map[string]string)
	for k, v := range s.settings {
		copy[k] = v
	}
	return copy
}

type settingsRequest struct {
	Settings map[string]string `json:"settings"`
}

func (s *Server) updateSettings(w http.ResponseWriter, r *http.Request) {
	var req settingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for k, v := range req.Settings {
		s.settings.Set(k, v)
		if k == "qwen_api_key" && s.llmGuard != nil {
			s.llmGuard.SetAPIKey(v)
		}
	}

	s.writeJSON(w, http.StatusOK, s.settings.GetAll())
}

func (s *Server) getSettings(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, s.settings.GetAll())
}
