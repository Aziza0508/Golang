package middleware

import (
	"encoding/json"
	"net/http"
)

func APIKeyAuth(validAPIKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-KEY")
			if apiKey == "" || apiKey != validAPIKey {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized: invalid or missing API key"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
