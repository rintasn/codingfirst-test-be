// middleware/auth_middleware.go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"main/utils"
)

// AuthMiddleware adalah middleware untuk memeriksa JWT token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Dapatkan header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Periksa format header (harus "Bearer TOKEN")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format, expected 'Bearer TOKEN'", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Validasi token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Tambahkan user ID ke konteks request
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
