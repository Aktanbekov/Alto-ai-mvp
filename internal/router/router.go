package router

import (
	"altoai_mvp/internal/auth"
	"altoai_mvp/internal/handlers"
	"altoai_mvp/internal/middleware"
	"altoai_mvp/internal/repository"
	"altoai_mvp/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), middleware.RequestLogger())

	// Initialize PostgreSQL repository
	userRepo, err := repository.NewPostgresRepo()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// wiring (DI)
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

	return r
}
