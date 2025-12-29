package router

import (
	"altoai_mvp/internal/auth"
	"altoai_mvp/internal/handlers"
	"altoai_mvp/internal/middleware"
	"altoai_mvp/internal/repository"
	"altoai_mvp/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), middleware.RequestLogger())

	// wiring (DI)
	userRepo := repository.NewUserMemoryRepo()
	userSvc := services.NewUserService(userRepo)
	authSvc := services.NewAuthService(userRepo)
	userH := handlers.NewUserHandler(userSvc)
	authH := handlers.NewAuthHandler(authSvc)
	chatH := handlers.NewChatHandler()

	// Initialize Google auth with the user repository
	auth.SetUserRepo(userRepo)

	// health
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// AUTH - Google
	r.GET("/auth/google", auth.HandleGoogleLogin)
	r.GET("/auth/google/callback", auth.HandleGoogleCallback)
	
	// User info endpoint (requires auth)
	r.GET("/me", middleware.JWTAuth(), func(c *gin.Context) {
		user := c.MustGet("user").(*middleware.MyClaims)
		c.JSON(http.StatusOK, gin.H{
			"email": user.Email,
			"name": user.Name,
			"picture": user.Picture,
		})
	})

	// versioned API
	v1 := r.Group("/api/v1")
	{
		// Auth routes
		v1.POST("/auth/login", authH.Login)
		v1.POST("/auth/register", authH.Register)
		v1.POST("/auth/verify-email", authH.VerifyEmail)
		v1.POST("/auth/refresh", authH.Refresh) // No auth middleware needed
		v1.POST("/auth/logout", authH.Logout)
		v1.POST("/auth/forgot-password", authH.ForgotPassword)
		v1.POST("/auth/reset-password", authH.ResetPassword)
		v1.POST("/auth/resend-verification", authH.ResendVerificationCode)
		
		// User routes
		v1.GET("/users", userH.List)
		v1.POST("/users", userH.Create)
		v1.GET("/users/:id", userH.Get)
		v1.PUT("/users/:id", userH.Update)
		v1.DELETE("/users/:id", userH.Delete)
		
		// Chat route (requires auth)
		v1.POST("/chat", middleware.JWTAuth(), chatH.Chat)
	}

	return r
}
