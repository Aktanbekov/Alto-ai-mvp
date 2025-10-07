package middleware


import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	jwt.RegisteredClaims
}


func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tok, err := c.Cookie("session")
		if err != nil {
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
		c.Set("user", claims)
		c.Next()
	}
}