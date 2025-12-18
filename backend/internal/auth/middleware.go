package auth

import (
	"context"
	"net/http"

	"aiguardrails/internal/tenant"
)

type contextKey string

const (
	appIDKey    contextKey = "appID"
	tenantIDKey contextKey = "tenantID"
	adminKey    contextKey = "admin"
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
			if err != nil || app.APISecret != secret || app.Revoked {
				http.Error(w, "invalid credentials", http.StatusUnauthorized)
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

