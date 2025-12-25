package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS returns a Gin middleware handler for CORS
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Allow same-origin requests (for production when frontend is served from same server)
		if origin == "" {
			// Same-origin request, allow it
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			// Cross-origin request - allow common dev origins
			allowedOrigin := "http://localhost:5173"
			if origin == "http://localhost:3000" || origin == "http://localhost:8080" || origin == "http://127.0.0.1:5173" || origin == "http://127.0.0.1:3000" {
				allowedOrigin = origin
			}
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
		}
		
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// Legacy CORS for http.Handler (used in main.go)
func CORSLegacy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		// Allow same-origin requests (for production when frontend is served from same server)
		if origin == "" {
			// Same-origin request, allow it
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			// Cross-origin request - allow common dev origins
			allowedOrigin := "http://localhost:5173"
			if origin == "http://localhost:3000" || origin == "http://localhost:8080" || origin == "http://127.0.0.1:5173" || origin == "http://127.0.0.1:3000" {
				allowedOrigin = origin
			}
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		}
		
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
