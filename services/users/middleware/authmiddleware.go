package middleware

import (
	"context"
	"net/http"
	"strings"

	"spotiftn/users/jwt"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "role"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// IZVLAČENJE IZ MapClaims

		userID, ok := claims["userId"].(string)
		if !ok {
			http.Error(w, "invalid userId in token", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "invalid role in token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, RoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
