package rbac

import (
	"context"
	"net/http"
)

type ctxKey string

const roleKey ctxKey = "role"

// WithRole injects role into context (placeholder for real JWT/OIDC parsing).
func WithRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), roleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequirePerm enforces permission based on role in context.
func RequirePerm(perm string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, _ := r.Context().Value(roleKey).(string)
			if !HasPermission(role, perm) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RoleFromContext returns role string.
func RoleFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(roleKey).(string); ok {
		return v
	}
	if v, ok := ctx.Value("role").(string); ok {
		return v
	}
	return ""
}

