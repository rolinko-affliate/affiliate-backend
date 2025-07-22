package handlers

import (
	"net/http"
	"strconv"
	"strings"

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

// isNotFoundError checks if an error indicates a resource was not found
func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "not found")
}

// getPaginationParams extracts pagination parameters from query string
func getPaginationParams(c *gin.Context) (page, pageSize int) {
	page = 1
	pageSize = 20

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	return page, pageSize
}

// HealthCheck handler for health checks
// @Summary      Health check endpoint
// @Description  Returns the health status of the API
// @Tags         system
// @Produce      json
// @Success      200  {object}  map[string]string  "Status UP"
// @Router       /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}
