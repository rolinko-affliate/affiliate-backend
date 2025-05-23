package middleware

import (
	"github.com/affiliate-backend/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware returns a middleware that handles CORS based on the environment
func CORSMiddleware() gin.HandlerFunc {
	// In development, allow all origins
	if config.AppConfig.Environment == "development" {
		return cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		})
	}

	// In production, don't allow CORS (default behavior)
	return func(c *gin.Context) {
		c.Next()
	}
}