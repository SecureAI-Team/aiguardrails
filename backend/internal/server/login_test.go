package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"

	"aiguardrails/internal/auth"
	"aiguardrails/internal/config"
	"aiguardrails/internal/rbac"
)

func TestLoginSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store := auth.NewUserStore(db)
	pwHash, _ := bcrypt.GenerateFromPassword([]byte("p"), bcrypt.DefaultCost)
	// Mock lookup for Verify
	mock.ExpectQuery("SELECT id, username, password_hash, role, created_at, updated_at FROM users").
		WithArgs("u").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "role", "created_at", "updated_at"}).
			AddRow("id1", "u", string(pwHash), rbac.RolePlatformAdmin, time.Now(), time.Now()))

	signer := &auth.JWTSigner{Secret: []byte("secret")}
	s := &Server{cfg: config.Default(), userStore: store, jwtSigner: signer}

	body, _ := json.Marshal(map[string]string{"username": "u", "password": "p"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	s.login(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var out map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if out["token"] == "" {
		t.Fatalf("expected token")
	}
}

func TestAdminAuthMiddleware(t *testing.T) {
	signer := &auth.JWTSigner{Secret: []byte("secret")}
	token, _ := signer.Sign("admin@example.com", rbac.RolePlatformAdmin, "", time.Hour)
	s := &Server{cfg: config.Default(), jwtSigner: signer}
	s.cfg.AdminToken = "adm"

	// via Admin token
	r1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r1.Header.Set("X-Admin-Token", "adm")
	w1 := httptest.NewRecorder()
	s.adminAuth()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(w1, r1)
	if w1.Result().StatusCode != http.StatusOK {
		t.Fatalf("admin token should pass")
	}

	// via JWT
	r2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	s.adminAuth()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(w2, r2)
	if w2.Result().StatusCode != http.StatusOK {
		t.Fatalf("admin jwt should pass")
	}

	// forbidden role
	userToken, _ := signer.Sign("u", "tenant_user", "", time.Hour)
	r3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r3.Header.Set("Authorization", "Bearer "+userToken)
	w3 := httptest.NewRecorder()
	s.adminAuth()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(w3, r3)
	if w3.Result().StatusCode != http.StatusForbidden {
		t.Fatalf("non-admin should be forbidden")
	}
}
