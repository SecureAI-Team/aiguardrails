package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if s.userStore == nil {
		http.Error(w, "user store unavailable", http.StatusInternalServerError)
		return
	}
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u, err := s.userStore.Verify(req.Username, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	token, err := s.jwtSigner.Sign(u.Username, u.Role, "", 24*time.Hour)
	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, loginResponse{Token: token})
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	if s.userStore == nil {
		http.Error(w, "user store unavailable", http.StatusInternalServerError)
		return
	}
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password required", http.StatusBadRequest)
		return
	}
	// Check if user exists
	if _, err := s.userStore.GetByUsername(req.Username); err == nil {
		http.Error(w, "username already taken", http.StatusConflict)
		return
	}
	// Create user
	u, err := s.userStore.Create(req.Username, req.Password, "tenant_admin") // Default role
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Auto login
	token, err := s.jwtSigner.Sign(u.Username, u.Role, "", 24*time.Hour)
	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusCreated, loginResponse{Token: token})
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authz := r.Header.Get("Authorization")
		if authz == "" || !strings.HasPrefix(authz, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authz, "Bearer ")
		claims, err := s.jwtSigner.Parse(tokenStr)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), authRoleCtxKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
