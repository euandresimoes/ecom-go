package middlewares

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/euandresimoes/ecom-go/internal/domain/auth"
)

func Auth(jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
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

			next.ServeHTTP(w, r)
		})
	}
}
