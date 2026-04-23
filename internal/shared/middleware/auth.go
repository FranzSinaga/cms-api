package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/FranzSinaga/blogcms/internal/shared"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaim struct {
	UserID string
	Email  string
	Role   string
	Name   string
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read cookie
		cookie, err := r.Cookie(shared.AuthTokenCookieName)
		if err != nil {
			shared.WriteError(w, "Unauthorized: missing authentication token", http.StatusUnauthorized)
			return
		}

		// Validate JWT
		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized: invalid token claims", http.StatusUnauthorized)
			return
		}

		// Safely extract user information from claims
		userID, ok := claims["user_id"].(string)
		if !ok {
			http.Error(w, "Unauthorized: invalid token claims", http.StatusUnauthorized)
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			http.Error(w, "Unauthorized: invalid token claims", http.StatusUnauthorized)
			return
		}

		name, ok := claims["name"].(string)
		if !ok {
			http.Error(w, "Unauthorized: invalid token claims", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "Unauthorized: invalid token claims", http.StatusUnauthorized)
			return
		}

		// Store user claims in context
		userClaims := &UserClaim{
			UserID: userID,
			Email:  email,
			Role:   role,
			Name:   name,
		}

		ctx := context.WithValue(r.Context(), shared.UserContextKey, userClaims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
