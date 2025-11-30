package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/euandresimoes/ecom-go/backend/internal/infra/security"
	"github.com/euandresimoes/ecom-go/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(jwtManager *security.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]any{
					"status": http.StatusUnauthorized,
					"error":  "missing authorization header",
				})
				return
			}

			tokenString := strings.Replace(header, "Bearer ", "", 1)
			if tokenString == "" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]any{
					"status": http.StatusUnauthorized,
					"error":  "invalid authorization header",
				})
				return
			}

			token, err := jwtManager.Verify(tokenString)
			if err != nil || !token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]any{
					"status": http.StatusUnauthorized,
					"error":  "invalid token",
				})
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]any{
					"status": http.StatusUnauthorized,
					"error":  "invalid claims",
				})
				return
			}

			id := claims["id"].(float64)
			role := claims["role"].(string)

			ctx := r.Context()
			ctx = context.WithValue(ctx, models.UserIDKey, id)
			ctx = context.WithValue(ctx, models.UserRoleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
