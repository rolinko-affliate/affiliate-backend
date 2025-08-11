package middleware

import (
	"net/http"

	"github.com/affiliate-backend/internal/platform/logger"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RBACMiddleware checks if the user's role is in the allowedRoles list
func RBACMiddleware(profileService service.ProfileService, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, exists := c.Get(UserIDKey)
		if !exists {
			logger.Error("User ID not found in context")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
			return
		}

		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			logger.Error("Error parsing User ID", "error", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid User ID format in context"})
			return
		}

		// Fetch profile (which includes role) from your database
		profile, err := profileService.GetProfileByID(c.Request.Context(), userID)
		if err != nil {
			logger.Error("Error fetching profile", "user_id", userID, "error", err)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User profile not found or access denied"})
			return
		}

		// Fetch role name based on profile.RoleID
		role, err := profileService.GetRoleByID(c.Request.Context(), profile.RoleID)
		if err != nil {
			logger.Error("Error fetching role", "user_id", userID, "role_id", profile.RoleID, "error", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not determine user role"})
			return
		}

		isAllowed := false
		for _, allowedRole := range allowedRoles {
			if role.Name == allowedRole {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource"})
			return
		}

		// Add role, organization_id, and profile to context for handlers
		c.Set(UserRoleKey, role.Name)
		c.Set("profile", profile) // Set the full profile object for handlers that need it
		if profile.OrganizationID != nil {
			c.Set("organizationID", *profile.OrganizationID)
		}

		c.Next()
	}
}
