package handlers

import (
	"net/http"

	"github.com/affiliate-backend/internal/api/middleware"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProfileHandler handles profile-related requests
type ProfileHandler struct {
	profileService service.ProfileService
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(ps service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: ps}
}

// SupabaseUserWebhookPayload defines the expected structure from Supabase new user webhook
type SupabaseUserWebhookPayload struct {
	Type   string `json:"type"` // e.g., "INSERT"
	Table  string `json:"table"` // e.g., "users"
	Record struct {
		ID    uuid.UUID `json:"id"`
		Email string    `json:"email"`
		// Other fields from Supabase user record
	} `json:"record"`
}

// HandleSupabaseNewUserWebhook handles the webhook from Supabase for new user creation
// @Summary      Handle Supabase new user webhook
// @Description  Process webhook from Supabase when a new user is created
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        payload  body      SupabaseUserWebhookPayload  true  "Webhook payload"
// @Success      200      {object}  map[string]string           "User profile created successfully"
// @Failure      400      {object}  map[string]string           "Invalid webhook payload"
// @Failure      500      {object}  map[string]string           "Internal server error"
// @Router       /public/webhooks/supabase/new-user [post]
func (h *ProfileHandler) HandleSupabaseNewUserWebhook(c *gin.Context) {
	var payload SupabaseUserWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook payload: " + err.Error()})
		return
	}

	// Verify it's a new user creation event for auth.users
	if payload.Type != "INSERT" || payload.Table != "users" {
		c.JSON(http.StatusOK, gin.H{"message": "Event ignored, not a new user creation."})
		return
	}

	userID := payload.Record.ID
	email := payload.Record.Email

	if userID == uuid.Nil || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID or email in webhook payload"})
		return
	}

	// MVP: Assign a default role and potentially create/assign to a default organization
	// This logic needs to be defined based on your business requirements.
	var defaultOrgID int64 = 1 // Example: Assume org ID 1 exists
	var defaultRoleID int = 4  // Example: 'Affiliate' role ID

	_, err := h.profileService.CreateNewUserProfile(c.Request.Context(), userID, email, &defaultOrgID, defaultRoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user profile: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User profile created successfully from webhook"})
}

// GetMyProfile retrieves the current user's profile
// @Summary      Get current user profile
// @Description  Retrieves the profile of the currently authenticated user
// @Tags         profile
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "User profile"
// @Failure      400  {object}  map[string]string       "Invalid user ID format"
// @Failure      404  {object}  map[string]string       "Profile not found"
// @Security     BearerAuth
// @Router       /users/me [get]
func (h *ProfileHandler) GetMyProfile(c *gin.Context) {
	userIDStr, _ := c.Get(middleware.UserIDKey)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	profile, err := h.profileService.GetProfileByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}
	c.JSON(http.StatusOK, profile)
}