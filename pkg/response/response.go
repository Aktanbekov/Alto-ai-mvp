package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{"data": data})
}

func Error(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

func ValidationError(c *gin.Context, details map[string]string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   "validation_error",
		"details": details,
	})
}
