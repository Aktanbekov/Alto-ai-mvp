package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"altoai_mvp/internal/models"
	"altoai_mvp/internal/repository"
	"altoai_mvp/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func googleConf() *oauth2.Config {
	_ = godotenv.Load() // ok if .env not present; will use OS env

	// Get redirect URL from environment or use default
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if redirectURL == "" {
		// Default based on environment
		if os.Getenv("GIN_MODE") == "release" {
			redirectURL = "http://localhost:3000/auth/google/callback" // Docker default
		} else {
			redirectURL = "http://localhost:8080/auth/google/callback" // Local dev
		}
	}

	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

// JWT CLAIMS
type MyClaims struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	jwt.RegisteredClaims
}

// GET /auth/google
func HandleGoogleLogin(c *gin.Context) {
	conf := googleConf()
	if conf.ClientID == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "missing GOOGLE_CLIENT_ID"})
		return
	}
	url := conf.AuthCodeURL("state-123", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

// GET /auth/google/callback?code=...
func HandleGoogleCallback(c *gin.Context) {
	conf := googleConf()

	code := c.Query("code")
	if code == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing code"})
		return
	}

	tok, err := conf.Exchange(c, code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token exchange failed"})
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tok.AccessToken)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "failed to fetch userinfo"})
		return
	}
	defer resp.Body.Close()

	var gu googleUser
	if err := json.NewDecoder(resp.Body).Decode(&gu); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to decode user info"})
		return
	}

	// Initialize user service and create/update user
	userRepo, err := repository.NewPostgresRepo()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "database connection failed"})
		return
	}
	defer userRepo.Close()

	userService := services.NewUserService(userRepo)

	// Create user using service - ignore returned user since we don't need it
	_, err = userService.Create(c.Request.Context(), models.CreateUserDTO{
		Email: gu.Email,
		Name:  gu.Name,
	})
	if err != nil {
		// If error occurs (likely user already exists), we can ignore it
		// In a production environment, you might want to handle this differently
		// or update the existing user's information
	}

	// JWT CREATE
	secret := os.Getenv("JWT_SECRET")
	claims := MyClaims{
		Email:   gu.Email,
		Name:    gu.Name,
		Picture: gu.Picture,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "altoai_mvp",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not sign jwt"})
		return
	}

	// Get frontend URL from environment or use default
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		// Default to same origin in production, or dev server in development
		// Check if we're in production (no dev server running)
		if os.Getenv("GIN_MODE") == "release" {
			frontendURL = "http://localhost:3000" // Docker default
		} else {
			frontendURL = "http://localhost:5173" // Vite dev server
		}
	}

	// Issue HttpOnly session cookie
	// Use empty domain for same-origin, or specific domain if needed
	cookieDomain := os.Getenv("COOKIE_DOMAIN")
	if cookieDomain == "" {
		cookieDomain = "" // Empty means same origin
	}
	c.SetCookie("session", signed, 7*24*60*60, "/", cookieDomain, false, true)

	// Back to frontend
	c.Redirect(http.StatusFound, frontendURL+"/chat")
}
