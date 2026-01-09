package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// registerTracingRoutes 注册追踪和P1/P2路由
func (s *Server) registerTracingRoutes(r chi.Router) {
	// 请求追踪
	r.Get("/traces", s.listTraces)
	r.Get("/traces/{id}", s.getTrace)

	// 模型目录
	r.Get("/models", s.listModels)

	// 成本
	r.Get("/cost", s.getCostSummary)

	// 导出
	r.Get("/exports", s.listExports)
	r.Post("/exports", s.createExport)
}

func (s *Server) listTraces(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}
	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	var blocked *bool
	if b := r.URL.Query().Get("blocked"); b != "" {
		val := b == "true"
		blocked = &val
	}

	traces, err := s.tracingStore.ListTraces(tenantID, blocked, nil, limit)
	if err != nil {
		log.Printf("listTraces error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, traces)
}

func (s *Server) getTrace(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	trace, err := s.tracingStore.GetTrace(id)
	if err != nil {
		http.Error(w, "trace not found", http.StatusNotFound)
		return
	}
	s.writeJSON(w, http.StatusOK, trace)
}

func (s *Server) listModels(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	models, err := s.tracingStore.ListModels(provider)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, models)
}

func (s *Server) getCostSummary(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	records, err := s.tracingStore.GetCostSummary(tenantID, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, records)
}

func (s *Server) listExports(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}
	jobs, err := s.tracingStore.ListExportJobs(tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, jobs)
}

type createExportRequest struct {
	TenantID string          `json:"tenant_id"`
	Type     string          `json:"type"`
	Filters  json.RawMessage `json:"filters"`
}

func (s *Server) createExport(w http.ResponseWriter, r *http.Request) {
	var req createExportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.TenantID == "" {
		req.TenantID = "default"
	}

	job, err := s.tracingStore.CreateExportJob(req.TenantID, req.Type, req.Filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// TODO: 启动异步导出任务
	s.writeJSON(w, http.StatusCreated, job)
}
