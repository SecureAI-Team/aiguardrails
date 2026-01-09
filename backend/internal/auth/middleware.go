package auth

import (
	"context"
	"net/http"
	"strings"

	"aiguardrails/internal/tenant"
)

type contextKey string

const (
	appIDKey    contextKey = "appID"
	tenantIDKey contextKey = "tenantID"
	adminKey    contextKey = "admin"
	userKey     contextKey = "user"
)

// APIKeyMiddleware enforces app key/secret authentication.
func APIKeyMiddleware(svc tenant.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			appID := r.Header.Get("X-App-Id")
			secret := r.Header.Get("X-App-Secret")
			if appID == "" || secret == "" {
				http.Error(w, "missing credentials", http.StatusUnauthorized)
				return
			}
			app, err := svc.GetApp(appID)
			if err != nil {
				// Log detailed error
				println("APIKeyAuth Failed: App not found or DB error:", err.Error(), "AppID:", appID)
				http.Error(w, "invalid credentials (app lookup)", http.StatusUnauthorized)
				return
			}
			if app.Revoked {
				println("APIKeyAuth Failed: App revoked:", appID)
				http.Error(w, "invalid credentials (revoked)", http.StatusUnauthorized)
				return
			}
			if !strings.EqualFold(app.APISecret, secret) {
				http.Error(w, "invalid credentials (secret mismatch)", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), appIDKey, app.ID)
			ctx = context.WithValue(ctx, tenantIDKey, app.TenantID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AppIDFromContext extracts authenticated app ID.
func AppIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(appIDKey).(string); ok {
		return v
	}
	return ""
}

// TenantIDFromContext extracts authenticated tenant ID.
func TenantIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(tenantIDKey).(string); ok {
		return v
	}
	if v, ok := ctx.Value("tenantID").(string); ok {
		return v
	}
	return ""
}

// AdminTokenMiddleware checks for platform admin token.
func AdminTokenMiddleware(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if token == "" {
				http.Error(w, "admin token not configured", http.StatusUnauthorized)
				return
			}
			if r.Header.Get("X-Admin-Token") != token {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), adminKey, true)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// IsAdmin returns true if admin token passed.
func IsAdmin(ctx context.Context) bool {
	if v, ok := ctx.Value(adminKey).(bool); ok {
		return v
	}
	return false
}

// WithUser sets authenticated user context (local login).
func WithUser(next http.Handler, user *User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserFromContext returns authenticated user.
func UserFromContext(ctx context.Context) *User {
	if v, ok := ctx.Value(userKey).(*User); ok {
		return v
	}
	return nil
}
