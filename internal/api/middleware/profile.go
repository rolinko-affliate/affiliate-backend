package middleware

import (
	"log"
	"net/http"

	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProfileMiddleware loads the user profile and sets it in context
// This middleware should be used after AuthMiddleware to load the full profile
func ProfileMiddleware(profileService service.ProfileService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, exists := c.Get(UserIDKey)
		if !exists {
			log.Println("User ID not found in context")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
			return
		}

		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			log.Printf("Error parsing User ID: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid User ID format in context"})
			return
		}

		// Fetch profile from database
		profile, err := profileService.GetProfileByID(c.Request.Context(), userID)
		if err != nil {
			log.Printf("Error fetching profile: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User profile not found"})
			return
		}

		// Set profile in context for handlers
		c.Set("profile", profile)
		c.Set(UserRoleKey, profile.RoleName) // Set user role for RBAC
		if profile.OrganizationID != nil {
			c.Set("organizationID", *profile.OrganizationID)
		}

		c.Next()
	}
}
