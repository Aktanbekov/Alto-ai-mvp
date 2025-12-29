package middleware


import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Type    string `json:"type"`
	jwt.RegisteredClaims
}


func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		
		// Extract token from "Bearer <token>"
		tok := strings.TrimPrefix(authHeader, "Bearer ")
		if tok == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		
		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.ParseWithClaims(tok, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(*MyClaims)
		
		// Verify this is an access token (not a refresh token)
		// Check the type claim from MapClaims if available
		if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
			if tokenType, ok := mapClaims["type"].(string); ok && tokenType != "access" {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		} else if claims.Type != "" && claims.Type != "access" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		
		c.Set("user", claims)
		c.Next()
	}
}