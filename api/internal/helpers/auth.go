package helpers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// AuthMiddleware validates JWT token and checks allowed roles
func AuthMiddleware(next http.Handler, allowedRoles ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			JSON(w, map[string]string{"error": "Unauthorized: Bearer token required"}, http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			JSON(w, map[string]string{"error": "Invalid token"}, http.StatusUnauthorized)
			return
		}

		claims, err := ValidateToken(token)
		if err != nil {
			JSON(w, map[string]string{"error": "Invalid token"}, http.StatusUnauthorized)
			return
		}

		// Role check
		if len(allowedRoles) > 0 {
			roleAllowed := false
			for _, role := range allowedRoles {
				if claims.Role == role {
					roleAllowed = true
					break
				}
			}
			if !roleAllowed {
				JSON(w, map[string]string{"error": "Forbidden"}, http.StatusForbidden)
				return
			}
		}

		// Pass claims to the request context
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper to respond with JSON
func JSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
