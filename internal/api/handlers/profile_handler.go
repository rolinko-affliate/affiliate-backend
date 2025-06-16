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

// ProfileRequest represents the request body for creating a profile
type ProfileRequest struct {
	OrganizationID *int64  `json:"organization_id,omitempty"`
	RoleID         int     `json:"role_id"`
	Email          string  `json:"email"`
	FirstName      *string `json:"first_name,omitempty"`
	LastName       *string `json:"last_name,omitempty"`
}

// UpdateProfileRequest represents the request body for updating a profile
// All fields are nullable to allow partial updates
type UpdateProfileRequest struct {
	OrganizationID *int64  `json:"organization_id,omitempty"`
	RoleID         *int    `json:"role_id,omitempty"`
	Email          *string `json:"email,omitempty"`
	FirstName      *string `json:"first_name,omitempty"`
	LastName       *string `json:"last_name,omitempty"`
	// Note: RoleName is not included as it's derived from RoleID
}

// UpsertProfileRequest represents the request body for upserting a profile
type UpsertProfileRequest struct {
	ID             string  `json:"id"`
	OrganizationID *int64  `json:"organization_id,omitempty"`
	RoleID         int     `json:"role_id"`
	Email          string  `json:"email"`
	FirstName      *string `json:"first_name,omitempty"`
	LastName       *string `json:"last_name,omitempty"`
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

// CreateProfile creates a new profile
// @Summary      Create a new profile
// @Description  Creates a new user profile
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param        profile  body      ProfileRequest  true  "Profile information"
// @Success      201      {object}  domain.Profile  "Created profile"
// @Failure      400      {object}  map[string]string  "Invalid request"
// @Failure      500      {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /profiles [post]
func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	var req ProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Validate required fields
	if req.Email == "" || req.RoleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and role_id are required"})
		return
	}

	// Get the user ID from the JWT token (set by auth middleware)
	userIDStr, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}
	
	profileID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Create the profile using the service
	createdProfile, err := h.profileService.CreateNewUserProfile(
		c.Request.Context(),
		profileID,
		req.Email,
		req.OrganizationID,
		req.RoleID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdProfile)
}

// UpdateProfile updates an existing profile
// @Summary      Update a profile
// @Description  Updates an existing user profile with only the provided fields
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param        id       path      string               true  "Profile ID"
// @Param        profile  body      UpdateProfileRequest true  "Updated profile information"
// @Success      200      {object}  domain.Profile       "Updated profile"
// @Failure      400      {object}  map[string]string    "Invalid request"
// @Failure      404      {object}  map[string]string    "Profile not found"
// @Failure      500      {object}  map[string]string    "Internal server error"
// @Security     BearerAuth
// @Router       /profiles/{id} [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	// Parse profile ID from URL
	profileIDStr := c.Param("id")
	profileID, err := uuid.Parse(profileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID format"})
		return
	}

	// Get existing profile
	existingProfile, err := h.profileService.GetProfileByID(c.Request.Context(), profileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	// Parse request body
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Only update fields that are provided in the request
	if req.OrganizationID != nil {
		existingProfile.OrganizationID = req.OrganizationID
	}
	if req.RoleID != nil {
		existingProfile.RoleID = *req.RoleID
		// Note: RoleName will be updated in the repository based on the new RoleID
	}
	if req.Email != nil {
		existingProfile.Email = *req.Email
	}
	if req.FirstName != nil {
		existingProfile.FirstName = req.FirstName
	}
	if req.LastName != nil {
		existingProfile.LastName = req.LastName
	}

	// Save updated profile
	updatedProfile, err := h.profileService.UpdateProfile(c.Request.Context(), existingProfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProfile)
}

// DeleteProfile deletes a profile
// @Summary      Delete a profile
// @Description  Deletes an existing user profile
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Profile ID"
// @Success      204  {object}  nil  "Profile deleted successfully"
// @Failure      400  {object}  map[string]string  "Invalid profile ID format"
// @Failure      404  {object}  map[string]string  "Profile not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /profiles/{id} [delete]
func (h *ProfileHandler) DeleteProfile(c *gin.Context) {
	// Parse profile ID from URL
	profileIDStr := c.Param("id")
	profileID, err := uuid.Parse(profileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID format"})
		return
	}

	// Delete the profile
	err = h.profileService.DeleteProfile(c.Request.Context(), profileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete profile: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpsertProfile creates or updates a profile
// @Summary      Upsert a profile
// @Description  Creates a new profile if it doesn't exist, or updates an existing one
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param        profile  body      UpsertProfileRequest  true  "Profile information"
// @Success      200      {object}  domain.Profile  "Upserted profile"
// @Failure      400      {object}  map[string]string  "Invalid request"
// @Failure      500      {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /profiles/upsert [post]
func (h *ProfileHandler) UpsertProfile(c *gin.Context) {
	var req UpsertProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Validate required fields
	if req.ID == "" || req.Email == "" || req.RoleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID, email, and role_id are required"})
		return
	}

	// Parse the UUID
	profileID, err := uuid.Parse(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID format"})
		return
	}

	// Upsert the profile using the service
	upsertedProfile, err := h.profileService.UpsertProfile(
		c.Request.Context(),
		profileID,
		req.Email,
		req.OrganizationID,
		req.RoleID,
		req.FirstName,
		req.LastName,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert profile: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, upsertedProfile)
}