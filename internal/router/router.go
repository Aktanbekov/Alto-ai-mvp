package router

import (
	"altoai_mvp/internal/auth"
	"altoai_mvp/internal/handlers"
	"altoai_mvp/internal/middleware"
	"altoai_mvp/internal/repository"
	"altoai_mvp/internal/services"
	"altoai_mvp/interview"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.CORS(), gin.Recovery(), middleware.RequestLogger())

	// Load interview questions
	questionsPath := "./interview/questions.json"
	if _, err := os.Stat(questionsPath); err == nil {
		if err := interview.LoadQuestions(questionsPath); err != nil {
			log.Printf("Warning: Failed to load interview questions: %v", err)
		} else {
			totalQuestions := 0
			for _, questions := range interview.QuestionsByCategory {
				totalQuestions += len(questions)
			}
			log.Printf("Loaded interview questions from %d categories (total: %d questions)", 
				len(interview.QuestionsByCategory), totalQuestions)
		}
	} else {
		log.Printf("Warning: Interview questions file not found at %s", questionsPath)
	}

	// wiring (DI)
	// Use PostgreSQL repository
	userRepo, err := repository.NewPostgresRepo()
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	// Note: We don't close the connection here as it's used throughout the server's lifetime
	// The connection will be closed when the server shuts down

	userSvc := services.NewUserService(userRepo)
	userH := handlers.NewUserHandler(userSvc)
	chatH := handlers.NewChatHandler()
	authSvc := services.NewAuthService(userRepo)
	authH := handlers.NewAuthHandler(authSvc)

	// Pass userRepo to Google auth handler
	auth.SetUserRepo(userRepo)

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
		// Auth endpoints
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authH.Login)
			auth.POST("/register", authH.Register)
			auth.POST("/verify-email", authH.VerifyEmail)
			auth.POST("/resend-verification", authH.ResendVerificationCode)
			auth.POST("/forgot-password", authH.ForgotPassword)
			auth.POST("/reset-password", authH.ResetPassword)
			auth.POST("/logout", authH.Logout)
			auth.POST("/refresh", authH.Refresh)
		}

		v1.GET("/users", userH.List)
		v1.POST("/users", userH.Create)
		v1.GET("/users/:id", userH.Get)
		v1.PUT("/users/:id", userH.Update)
		v1.DELETE("/users/:id", userH.Delete)
		v1.POST("/chat", chatH.Chat)

		// Interview endpoints
		interviewGroup := v1.Group("/interview")
		{
			interviewGroup.POST("/sessions", interview.CreateSessionHandler)
			interviewGroup.POST("/sessions/:id/answer", interview.SubmitAnswerHandler)
		}
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
