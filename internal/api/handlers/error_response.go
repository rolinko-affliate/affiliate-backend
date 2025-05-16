package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// RespondWithError sends a standardized error response
func RespondWithError(c *gin.Context, code int, message string, details ...string) {
	resp := ErrorResponse{Error: message}
	if len(details) > 0 {
		resp.Details = details[0]
	}
	c.AbortWithStatusJSON(code, resp)
}

// HealthCheck handler for health checks
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}