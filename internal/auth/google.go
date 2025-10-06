package auth

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
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
	_ = json.NewDecoder(resp.Body).Decode(&gu)

	// Issue HttpOnly session cookie (7 days). In production, sign/encode it or use JWT/sessions.
	c.SetCookie("session", "email="+gu.Email, 7*24*60*60, "/", "", false, true)

	// Back to frontend
	c.Redirect(http.StatusFound, "http://localhost:5173/dashboard")
}
