package blitzhttp

import (
	"log"
	"net/http"
)

// Logger logs requests
func Logger() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

// Auth rejects unauthorized requests
func Auth() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return // Short-circuits
			}
			next.ServeHTTP(w, r)
		})
	}
}
