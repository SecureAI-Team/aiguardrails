package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"aiguardrails/internal/alert"
)

// registerAlertRoutes 注册告警API路由
func (s *Server) registerAlertRoutes(r chi.Router) {
	// 告警规则
	r.Get("/alerts/rules", s.listAlertRules)
	r.Post("/alerts/rules", s.createAlertRule)
	r.Put("/alerts/rules/{id}", s.updateAlertRule)
	r.Delete("/alerts/rules/{id}", s.deleteAlertRule)

	// 告警历史
	r.Get("/alerts/history", s.listAlertHistory)
	r.Post("/alerts/history/{id}/ack", s.acknowledgeAlert)

	// 通知渠道
	r.Get("/alerts/channels", s.listNotificationChannels)
	r.Post("/alerts/channels", s.createNotificationChannel)
	r.Post("/alerts/channels/{id}/test", s.testNotificationChannel)
}

func (s *Server) listAlertRules(w http.ResponseWriter, r *http.Request) {
	rules, err := s.alertStore.ListRules(nil, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, rules)
}

func (s *Server) createAlertRule(w http.ResponseWriter, r *http.Request) {
	var rule alert.AlertRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.alertStore.CreateRule(&rule); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusCreated, rule)
}

func (s *Server) updateAlertRule(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var rule alert.AlertRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.alertStore.UpdateRule(id, &rule); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, rule)
}

func (s *Server) deleteAlertRule(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.alertStore.DeleteRule(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listAlertHistory(w http.ResponseWriter, r *http.Request) {
	severity := r.URL.Query().Get("severity")
	limit := 100

	history, err := s.alertStore.ListHistory(nil, severity, nil, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, history)
}

func (s *Server) acknowledgeAlert(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := "admin" // TODO: 从context获取当前用户

	if err := s.alertStore.AcknowledgeAlert(id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "acknowledged"})
}

func (s *Server) listNotificationChannels(w http.ResponseWriter, r *http.Request) {
	channels, err := s.alertStore.ListChannels(nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, channels)
}

func (s *Server) createNotificationChannel(w http.ResponseWriter, r *http.Request) {
	var channel alert.NotificationChannel
	if err := json.NewDecoder(r.Body).Decode(&channel); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.alertStore.SaveChannel(&channel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusCreated, channel)
}

func (s *Server) testNotificationChannel(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现渠道测试
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "test_sent"})
}
