package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	jwt "github.com/golang-jwt/jwt/v4"
)

// Claims represents expected JWT claims.
// OIDCMiddleware validates JWT using JWKS, issuer, and audience.
func OIDCMiddleware(jwksURL, issuer, audience, tenantClaim, adminRole, userRole string, timeoutSec, cacheMin int) func(http.Handler) http.Handler {
	kf, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshTimeout: time.Duration(timeoutSec) * time.Second,
		RefreshInterval: time.Duration(cacheMin) * time.Minute,
		RefreshErrorHandler: func(err error) {
			_ = err
		},
	})
	if err != nil || kf == nil {
		// fallback to empty JWKS to keep handler alive; will reject tokens
		kf = keyfunc.NewGiven(map[string]keyfunc.GivenKey{})
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authz := r.Header.Get("Authorization")
			if authz == "" || !strings.HasPrefix(authz, "Bearer ") {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authz, "Bearer ")
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, kf.Keyfunc)
			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			if issuer != "" && claims["iss"] != issuer {
				http.Error(w, "issuer mismatch", http.StatusUnauthorized)
				return
			}
			if audience != "" {
				if aud, ok := claims["aud"]; ok {
					switch v := aud.(type) {
					case string:
						if v != audience {
							http.Error(w, "aud mismatch", http.StatusUnauthorized)
							return
						}
					case []interface{}:
						found := false
						for _, item := range v {
							if s, ok := item.(string); ok && s == audience {
								found = true
								break
							}
						}
						if !found {
							http.Error(w, "aud mismatch", http.StatusUnauthorized)
							return
						}
					}
				}
			}
			roleVal := ""
			if v, ok := claims["role"].(string); ok {
				roleVal = v
			}
			roleVal = MapRole(roleVal, adminRole, userRole)
			ctx := context.WithValue(r.Context(), "role", roleVal)
			if tenantClaim != "" {
				if v, ok := claims[tenantClaim].(string); ok && v != "" {
					ctx = context.WithValue(ctx, tenantIDKey, v)
				}
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

