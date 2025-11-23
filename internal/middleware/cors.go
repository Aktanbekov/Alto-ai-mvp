package middleware

import "net/http"

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		// Allow same-origin requests (for production when frontend is served from same server)
		if origin == "" {
			// Same-origin request, allow it
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Host"))
		} else {
			// Cross-origin request - allow localhost:5173 for dev or same origin
			allowedOrigin := "http://localhost:5173"
			if origin == "http://localhost:3000" || origin == "http://localhost:8080" {
				allowedOrigin = origin
			}
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		}
		
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
