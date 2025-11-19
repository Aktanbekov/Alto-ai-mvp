package router

import (
	"altoai_mvp/internal/auth"
	"altoai_mvp/internal/handlers"
	"altoai_mvp/internal/middleware"
	"altoai_mvp/internal/repository"
	"altoai_mvp/internal/services"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), middleware.RequestLogger())

	// wiring (DI)
	userRepo := repository.NewUserMemoryRepo()
	userSvc := services.NewUserService(userRepo)
	userH := handlers.NewUserHandler(userSvc)

	// health
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// AUTH - Google
	r.GET("/auth/google", auth.HandleGoogleLogin)
	r.GET("/auth/google/callback", auth.HandleGoogleCallback)
	r.GET("/me", middleware.JWTAuth(), func(c *gin.Context) {
		user := c.MustGet("user").(*middleware.MyClaims)
		c.JSON(http.StatusOK, gin.H{
			"email":   user.Email,
			"name":    user.Name,
			"picture": user.Picture,
		})
	})

	// versioned API
	v1 := r.Group("/api/v1")
	{
		v1.GET("/users", userH.List)
		v1.POST("/users", userH.Create)
		v1.GET("/users/:id", userH.Get)
		v1.PUT("/users/:id", userH.Update)
		v1.DELETE("/users/:id", userH.Delete)
	}

	// Serve static files from frontend/dist (for production)
	// Check if frontend/dist exists, if not, skip static file serving (for development)
	frontendDist := "./frontend/dist"
	if _, err := os.Stat(frontendDist); err == nil {
		// Static assets
		if _, err := os.Stat(filepath.Join(frontendDist, "assets")); err == nil {
			r.Static("/assets", filepath.Join(frontendDist, "assets"))
		}

		// Static files
		if _, err := os.Stat(filepath.Join(frontendDist, "vite.svg")); err == nil {
			r.StaticFile("/vite.svg", filepath.Join(frontendDist, "vite.svg"))
		}
		if _, err := os.Stat(filepath.Join(frontendDist, "favicon.ico")); err == nil {
			r.StaticFile("/favicon.ico", filepath.Join(frontendDist, "favicon.ico"))
		}

		// Serve index.html for all non-API routes (SPA routing)
		r.NoRoute(func(c *gin.Context) {
			indexPath := filepath.Join(frontendDist, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				c.File(indexPath)
			} else {
				c.JSON(404, gin.H{"error": "Not found"})
			}
		})
	}

	return r
}
