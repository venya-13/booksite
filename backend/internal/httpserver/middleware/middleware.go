package middleware

import (
	"context"
	"net/http"
	"strings"

	"google-auth-demo/backend/internal/jwt"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing or invalid token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// You can save claims in context for handlers
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", claims.GoogleID)
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "is_admin", claims.IsAdmin)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
