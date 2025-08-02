package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/affiliate-backend/internal/api/middleware"
	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// OrganizationHandler handles organization-related requests
type OrganizationHandler struct {
	organizationService service.OrganizationService
	profileService      service.ProfileService
}

// NewOrganizationHandler creates a new organization handler
func NewOrganizationHandler(os service.OrganizationService, ps service.ProfileService) *OrganizationHandler {
	return &OrganizationHandler{
		organizationService: os,
		profileService:      ps,
	}
}

// checkOrganizationAccess verifies if the user has permission to access/modify the organization
// Returns true if the user has access, false otherwise
func (h *OrganizationHandler) checkOrganizationAccess(c *gin.Context, orgID int64) (bool, error) {
	// Get user role from context
	userRole, exists := c.Get(middleware.UserRoleKey)
	if !exists {
		return false, fmt.Errorf("user role not found in context")
	}

	// Admin can access all organizations
	if userRole.(string) == "Admin" {
		return true, nil
	}

	// Get user ID from context
	userIDStr, exists := c.Get(middleware.UserIDKey)
	if !exists {
		return false, fmt.Errorf("user ID not found in context")
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return false, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Get user's profile to check organization
	profile, err := h.profileService.GetProfileByID(c.Request.Context(), userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Check if user belongs to the organization
	if profile.OrganizationID == nil {
		return false, nil
	}

	return *profile.OrganizationID == orgID, nil
}

// AdvertiserExtraInfoRequest represents the extra info for advertiser organizations
type AdvertiserExtraInfoRequest struct {
	Website     *string `json:"website,omitempty"`
	WebsiteType *string `json:"website_type,omitempty" binding:"omitempty,oneof=shopify amazon shopline tiktok_shop other"`
	CompanySize *string `json:"company_size,omitempty" binding:"omitempty,oneof=startup small medium large enterprise"`
}

// AffiliateExtraInfoRequest represents the extra info for affiliate organizations
type AffiliateExtraInfoRequest struct {
	Website         *string `json:"website,omitempty"`
	AffiliateType   *string `json:"affiliate_type,omitempty" binding:"omitempty,oneof=cashback blog incentive content forum sub_affiliate_network other"`
	SelfDescription *string `json:"self_description,omitempty"`
	LogoURL         *string `json:"logo_url,omitempty"`
}

// CreateOrganizationRequest defines the request for creating an organization
type CreateOrganizationRequest struct {
	Name                string                       `json:"name" binding:"required"`
	Type                string                       `json:"type" binding:"required,oneof=advertiser affiliate platform_owner agency"`
	ContactEmail        string                       `json:"contact_email,omitempty"`
	Description         string                       `json:"description,omitempty"`
	AdvertiserExtraInfo *AdvertiserExtraInfoRequest `json:"advertiser_extra_info,omitempty"`
	AffiliateExtraInfo  *AffiliateExtraInfoRequest  `json:"affiliate_extra_info,omitempty"`
}

// CreateOrganization creates a new organization (authenticated endpoint)
// @Summary      Create a new organization (Admin only)
// @Description  Creates a new organization with the given name. Requires Admin role.
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        request  body      CreateOrganizationRequest  true  "Organization details"
// @Success      201      {object}  domain.Organization        "Created organization"
// @Failure      400      {object}  map[string]string          "Invalid request"
// @Failure      403      {object}  map[string]string          "Forbidden - Only admins can create organizations"
// @Failure      500      {object}  map[string]string          "Internal server error"
// @Security     BearerAuth
// @Router       /organizations [post]
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	// Check user role
	userRole, exists := c.Get(middleware.UserRoleKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User role not found in context"})
		return
	}

	if userRole.(string) != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only administrators can create organizations"})
		return
	}
	var req CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Convert string to OrganizationType
	orgType := domain.OrganizationType(req.Type)

	organization, err := h.organizationService.CreateOrganization(c.Request.Context(), req.Name, orgType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create organization: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, organization)
}

// CreateOrganizationPublic creates a new organization (public endpoint)
// @Summary      Create a new organization (Public)
// @Description  Creates a new organization with the given name and optional extra info. No authentication required.
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        request  body      CreateOrganizationRequest  true  "Organization details"
// @Success      201      {object}  domain.Organization        "Created organization"
// @Failure      400      {object}  map[string]string          "Invalid request"
// @Failure      500      {object}  map[string]string          "Internal server error"
// @Router       /public/organizations [post]
func (h *OrganizationHandler) CreateOrganizationPublic(c *gin.Context) {
	var req CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Convert string to OrganizationType
	orgType := domain.OrganizationType(req.Type)

	// Prepare service request
	serviceReq := &service.CreateOrganizationWithExtraInfoRequest{
		Name:         req.Name,
		Type:         orgType,
		ContactEmail: req.ContactEmail,
		Description:  req.Description,
	}

	// Convert extra info if provided
	if req.AdvertiserExtraInfo != nil {
		serviceReq.AdvertiserExtraInfo = &domain.AdvertiserExtraInfo{
			Website:     req.AdvertiserExtraInfo.Website,
			WebsiteType: req.AdvertiserExtraInfo.WebsiteType,
			CompanySize: req.AdvertiserExtraInfo.CompanySize,
		}
	}

	if req.AffiliateExtraInfo != nil {
		serviceReq.AffiliateExtraInfo = &domain.AffiliateExtraInfo{
			Website:         req.AffiliateExtraInfo.Website,
			AffiliateType:   req.AffiliateExtraInfo.AffiliateType,
			SelfDescription: req.AffiliateExtraInfo.SelfDescription,
			LogoURL:         req.AffiliateExtraInfo.LogoURL,
		}
	}

	organization, err := h.organizationService.CreateOrganizationWithExtraInfo(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create organization: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, organization)
}

// GetOrganization retrieves an organization by ID
// @Summary      Get organization by ID
// @Description  Retrieves an organization by its ID, optionally with extra info
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        id           path      int     true   "Organization ID"
// @Param        with_extra   query     bool    false  "Include extra info (advertiser_extra_info or affiliate_extra_info)"
// @Success      200          {object}  domain.Organization  "Organization details (basic)"
// @Success      200          {object}  domain.OrganizationWithExtraInfo  "Organization details (with extra info)"
// @Failure      400          {object}  map[string]string    "Invalid organization ID"
// @Failure      403          {object}  map[string]string    "Forbidden - User doesn't have permission"
// @Failure      404          {object}  map[string]string    "Organization not found"
// @Failure      500          {object}  map[string]string    "Internal server error"
// @Security     BearerAuth
// @Router       /organizations/{id} [get]
func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	// Check if extra info is requested
	withExtra := c.Query("with_extra") == "true"

	if withExtra {
		// Get organization with extra info
		organizationWithExtra, err := h.organizationService.GetOrganizationByIDWithExtraInfo(c.Request.Context(), id)
		if err != nil {
			if err.Error() == "organization not found: not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization: " + err.Error()})
			return
		}

		// Check if user has permission to view this organization
		hasAccess, err := h.checkOrganizationAccess(c, organizationWithExtra.OrganizationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions: " + err.Error()})
			return
		}
		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view this organization"})
			return
		}

		c.JSON(http.StatusOK, organizationWithExtra)
	} else {
		// Get basic organization info
		organization, err := h.organizationService.GetOrganizationByID(c.Request.Context(), id)
		if err != nil {
			if err.Error() == "organization not found: not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization: " + err.Error()})
			return
		}

		// Check if user has permission to view this organization
		hasAccess, err := h.checkOrganizationAccess(c, organization.OrganizationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions: " + err.Error()})
			return
		}
		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view this organization"})
			return
		}

		c.JSON(http.StatusOK, organization)
	}
}

// UpdateOrganizationRequest defines the request for updating an organization
type UpdateOrganizationRequest struct {
	Name                string                       `json:"name" binding:"required"`
	Type                string                       `json:"type" binding:"required,oneof=advertiser affiliate platform_owner agency"`
	ContactEmail        string                       `json:"contact_email,omitempty"`
	Description         string                       `json:"description,omitempty"`
	AdvertiserExtraInfo *AdvertiserExtraInfoRequest `json:"advertiser_extra_info,omitempty"`
	AffiliateExtraInfo  *AffiliateExtraInfoRequest  `json:"affiliate_extra_info,omitempty"`
}

// UpdateOrganization updates an organization
// @Summary      Update organization
// @Description  Updates an organization with the given details
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        id       path      int                       true  "Organization ID"
// @Param        request  body      UpdateOrganizationRequest  true  "Organization details"
// @Success      200      {object}  domain.Organization        "Updated organization"
// @Failure      400      {object}  map[string]string          "Invalid request"
// @Failure      403      {object}  map[string]string          "Forbidden - User doesn't have permission"
// @Failure      404      {object}  map[string]string          "Organization not found"
// @Failure      500      {object}  map[string]string          "Internal server error"
// @Security     BearerAuth
// @Router       /organizations/{id} [put]
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var req UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get existing organization
	organization, err := h.organizationService.GetOrganizationByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "organization not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization: " + err.Error()})
		return
	}

	// Check if user has permission to update this organization
	hasAccess, err := h.checkOrganizationAccess(c, organization.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions: " + err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this organization"})
		return
	}

	// Update organization
	organization.Name = req.Name
	organization.Type = domain.OrganizationType(req.Type)
	if err := h.organizationService.UpdateOrganization(c.Request.Context(), organization); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update organization: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, organization)
}

// ListOrganizations retrieves a list of organizations
// @Summary      List organizations
// @Description  Retrieves a list of organizations with pagination
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        page      query     int                      false  "Page number (default: 1)"
// @Param        pageSize  query     int                      false  "Page size (default: 10)"
// @Success      200       {array}   domain.Organization      "List of organizations"
// @Failure      500       {object}  map[string]string        "Internal server error"
// @Security     BearerAuth
// @Router       /organizations [get]
func (h *OrganizationHandler) ListOrganizations(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// Get user role from context
	userRole, exists := c.Get(middleware.UserRoleKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User role not found in context"})
		return
	}

	// Get all organizations
	organizations, err := h.organizationService.ListOrganizations(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list organizations: " + err.Error()})
		return
	}

	// If user is Admin, return all organizations
	if userRole.(string) == "Admin" {
		c.JSON(http.StatusOK, organizations)
		return
	}

	// For non-admin users, filter organizations to only include their own
	userOrgID, exists := c.Get("organizationID")
	if !exists {
		// If user doesn't have an organization, return empty list
		c.JSON(http.StatusOK, []*domain.Organization{})
		return
	}

	// Filter organizations to only include the user's organization
	var filteredOrgs []*domain.Organization
	for _, org := range organizations {
		if org.OrganizationID == userOrgID.(int64) {
			filteredOrgs = append(filteredOrgs, org)
		}
	}

	c.JSON(http.StatusOK, filteredOrgs)
}

// DeleteOrganization deletes an organization
// @Summary      Delete organization
// @Description  Deletes an organization by its ID
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        id   path      int                true  "Organization ID"
// @Success      204  {object}  nil                "No content"
// @Failure      400  {object}  map[string]string  "Invalid organization ID"
// @Failure      403  {object}  map[string]string  "Forbidden - User doesn't have permission"
// @Failure      404  {object}  map[string]string  "Organization not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /organizations/{id} [delete]
func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	// Get the organization first to check permissions
	organization, err := h.organizationService.GetOrganizationByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "organization not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization: " + err.Error()})
		return
	}

	// Check if user has permission to delete this organization
	hasAccess, err := h.checkOrganizationAccess(c, organization.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions: " + err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this organization"})
		return
	}

	if err := h.organizationService.DeleteOrganization(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
