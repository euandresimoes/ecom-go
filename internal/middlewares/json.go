package middlewares

import "net/http"

func JSON(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		n.ServeHTTP(w, r)
	})
}