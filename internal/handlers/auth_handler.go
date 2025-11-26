package handlers

import (
	"altoai_mvp/internal/models"
	"altoai_mvp/internal/services"
	errs "altoai_mvp/pkg/errors"
	"altoai_mvp/pkg/response"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc services.AuthService
}

func NewAuthHandler(authSvc services.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var dto models.LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, errs.FromBinding(err))
		return
	}

	token, user, err := h.authSvc.Login(c.Request.Context(), dto)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Set cookie
	cookieDomain := os.Getenv("COOKIE_DOMAIN")
	if cookieDomain == "" {
		cookieDomain = ""
	}
	c.SetCookie("session", token, 7*24*60*60, "/", cookieDomain, false, true)

	response.OK(c, gin.H{
		"user": gin.H{
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var dto models.CreateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, errs.FromBinding(err))
		return
	}

	err := h.authSvc.Register(c.Request.Context(), dto)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(c, gin.H{
		"message": "Registration successful. Please check your email for verification code.",
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var dto models.VerifyEmailDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, errs.FromBinding(err))
		return
	}

	token, user, err := h.authSvc.VerifyEmail(c.Request.Context(), dto)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Set cookie
	cookieDomain := os.Getenv("COOKIE_DOMAIN")
	if cookieDomain == "" {
		cookieDomain = ""
	}
	c.SetCookie("session", token, 7*24*60*60, "/", cookieDomain, false, true)

	response.OK(c, gin.H{
		"user": gin.H{
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func (h *AuthHandler) ResendVerificationCode(c *gin.Context) {
	var dto models.ResendVerificationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, errs.FromBinding(err))
		return
	}

	err := h.authSvc.ResendVerificationCode(c.Request.Context(), dto)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Don't reveal if user exists or not
	response.OK(c, gin.H{
		"message": "If an account exists with this email and is not verified, a new verification code has been sent.",
	})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var dto models.ForgotPasswordDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, errs.FromBinding(err))
		return
	}

	err := h.authSvc.ForgotPassword(c.Request.Context(), dto)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Don't reveal if user exists or not
	response.OK(c, gin.H{
		"message": "If an account exists with this email, a password reset code has been sent.",
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var dto models.ResetPasswordDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, errs.FromBinding(err))
		return
	}

	err := h.authSvc.ResetPassword(c.Request.Context(), dto)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(c, gin.H{
		"message": "Password reset successful. You can now login with your new password.",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear the session cookie
	cookieDomain := os.Getenv("COOKIE_DOMAIN")
	if cookieDomain == "" {
		cookieDomain = ""
	}
	c.SetCookie("session", "", -1, "/", cookieDomain, false, true)
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	// Get token from cookie
	token, err := c.Cookie("session")
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "no session found")
		return
	}

	newToken, err := h.authSvc.RefreshToken(c.Request.Context(), token)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Set new cookie
	cookieDomain := os.Getenv("COOKIE_DOMAIN")
	if cookieDomain == "" {
		cookieDomain = ""
	}
	c.SetCookie("session", newToken, 7*24*60*60, "/", cookieDomain, false, true)

	response.OK(c, gin.H{"message": "token refreshed"})
}
