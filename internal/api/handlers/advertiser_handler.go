package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/affiliate-backend/internal/api/middleware"
	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AdvertiserHandler handles advertiser-related requests
type AdvertiserHandler struct {
	advertiserService service.AdvertiserService
	profileService    service.ProfileService
}

// NewAdvertiserHandler creates a new advertiser handler
func NewAdvertiserHandler(as service.AdvertiserService, ps service.ProfileService) *AdvertiserHandler {
	return &AdvertiserHandler{
		advertiserService: as,
		profileService:    ps,
	}
}

// checkAdvertiserAccess verifies if the user has permission to access/modify the advertiser
// Returns true if the user has access, false otherwise
func (h *AdvertiserHandler) checkAdvertiserAccess(c *gin.Context, advertiserOrgID int64) (bool, error) {
	// Get user role from context
	userRole, exists := c.Get(middleware.UserRoleKey)
	if !exists {
		return false, fmt.Errorf("user role not found in context")
	}
	
	// Admin can access all advertisers
	if userRole.(string) == "Admin" {
		return true, nil
	}
	
	// Get user's organization ID from context
	userOrgID, exists := c.Get("organizationID")
	if !exists {
		return false, fmt.Errorf("user organization ID not found in context")
	}
	
	// Check if user belongs to the same organization as the advertiser
	return userOrgID.(int64) == advertiserOrgID, nil
}

// CreateAdvertiserRequest defines the request for creating an advertiser
// swagger:model
type CreateAdvertiserRequest struct {
	// Organization ID
	OrganizationID int64  `json:"organization_id" binding:"required" example:"1"`
	// Advertiser name
	Name           string `json:"name" binding:"required" example:"Example Advertiser"`
	// Contact email address
	ContactEmail   *string `json:"contact_email,omitempty" example:"contact@example.com"`
	// Billing details in JSON format
	// swagger:strfmt json
	BillingDetails *json.RawMessage `json:"billing_details,omitempty" swaggertype:"object"`
	// Status of the advertiser (active, pending, inactive, rejected)
	Status         string `json:"status,omitempty" example:"active"`
	
	// Enhanced fields for Everflow integration
	InternalNotes              *string `json:"internal_notes,omitempty" example:"Internal notes about the advertiser"`
	DefaultCurrencyID          *string `json:"default_currency_id,omitempty" example:"USD"`
	PlatformName               *string `json:"platform_name,omitempty" example:"Example Platform"`
	PlatformURL                *string `json:"platform_url,omitempty" example:"https://platform.example.com"`
	PlatformUsername           *string `json:"platform_username,omitempty" example:"platform_user"`
	AccountingContactEmail     *string `json:"accounting_contact_email,omitempty" example:"accounting@example.com"`
	OfferIDMacro               *string `json:"offer_id_macro,omitempty" example:"oid"`
	AffiliateIDMacro           *string `json:"affiliate_id_macro,omitempty" example:"aid"`
	AttributionMethod          *string `json:"attribution_method,omitempty" example:"first_click"`
	EmailAttributionMethod     *string `json:"email_attribution_method,omitempty" example:"first_click"`
	AttributionPriority        *string `json:"attribution_priority,omitempty" example:"high"`
	ReportingTimezoneID        *int    `json:"reporting_timezone_id,omitempty" example:"67"`
	IsExposePublisherReporting *bool   `json:"is_expose_publisher_reporting,omitempty" example:"false"`
	
	// Synchronization options
	SyncToEverflow *bool `json:"sync_to_everflow,omitempty" example:"true"`
}

// CreateAdvertiser creates a new advertiser
// @Summary      Create a new advertiser
// @Description  Creates a new advertiser with the given details
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        request  body      CreateAdvertiserRequest  true  "Advertiser details"
// @Success      201      {object}  domain.Advertiser        "Created advertiser"
// @Failure      400      {object}  map[string]string        "Invalid request"
// @Failure      403      {object}  map[string]string        "Forbidden - User doesn't have permission"
// @Failure      500      {object}  map[string]string        "Internal server error"
// @Security     BearerAuth
// @Router       /advertisers [post]
func (h *AdvertiserHandler) CreateAdvertiser(c *gin.Context) {
	var req CreateAdvertiserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Check if user has permission to create an advertiser for this organization
	hasAccess, err := h.checkAdvertiserAccess(c, req.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions: " + err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to create an advertiser for this organization"})
		return
	}

	advertiser := &domain.Advertiser{
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		ContactEmail:   req.ContactEmail,
		Status:         req.Status,
		
		// Enhanced fields for Everflow integration
		InternalNotes:              req.InternalNotes,
		DefaultCurrencyID:          req.DefaultCurrencyID,
		PlatformName:               req.PlatformName,
		PlatformURL:                req.PlatformURL,
		PlatformUsername:           req.PlatformUsername,
		AccountingContactEmail:     req.AccountingContactEmail,
		OfferIDMacro:               req.OfferIDMacro,
		AffiliateIDMacro:           req.AffiliateIDMacro,
		AttributionMethod:          req.AttributionMethod,
		EmailAttributionMethod:     req.EmailAttributionMethod,
		AttributionPriority:        req.AttributionPriority,
		ReportingTimezoneID:        req.ReportingTimezoneID,
		IsExposePublisherReporting: req.IsExposePublisherReporting,
	}

	if req.BillingDetails != nil {
		var billingDetails domain.BillingDetails
		if err := json.Unmarshal(*req.BillingDetails, &billingDetails); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid billing details format: " + err.Error()})
			return
		}
		advertiser.BillingDetails = &billingDetails
	}

	createdAdvertiser, err := h.advertiserService.CreateAdvertiser(c.Request.Context(), advertiser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create advertiser: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdAdvertiser)
}

// GetAdvertiser retrieves an advertiser by ID with optional Everflow data comparison
// @Summary      Get advertiser by ID
// @Description  Retrieves an advertiser by its ID with optional Everflow data comparison
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        id              path      int                 true   "Advertiser ID"
// @Param        include_everflow query     bool                false  "Include Everflow data comparison"
// @Success      200             {object}  domain.Advertiser   "Advertiser details"
// @Success      200             {object}  domain.AdvertiserWithEverflowData  "Advertiser with Everflow comparison (when include_everflow=true)"
// @Failure      400             {object}  map[string]string   "Invalid advertiser ID"
// @Failure      403             {object}  map[string]string   "Forbidden - User doesn't have permission"
// @Failure      404             {object}  map[string]string   "Advertiser not found"
// @Failure      500             {object}  map[string]string   "Internal server error"
// @Security     BearerAuth
// @Router       /advertisers/{id} [get]
func (h *AdvertiserHandler) GetAdvertiser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	// Check if Everflow data should be included
	includeEverflow := c.Query("include_everflow") == "true"

	if includeEverflow {
		// Get advertiser with Everflow data comparison
		advertiserWithEverflow, err := h.advertiserService.GetAdvertiserWithEverflowData(c.Request.Context(), id)
		if err != nil {
			if err.Error() == "advertiser not found: not found" || err.Error() == "failed to get advertiser: not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get advertiser: " + err.Error()})
			return
		}

		// Check if user has permission to view this advertiser
		hasAccess, err := h.checkAdvertiserAccess(c, advertiserWithEverflow.Advertiser.OrganizationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions: " + err.Error()})
			return
		}
		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view this advertiser"})
			return
		}

		c.JSON(http.StatusOK, advertiserWithEverflow)
	} else {
		// Get basic advertiser data
		advertiser, err := h.advertiserService.GetAdvertiserByID(c.Request.Context(), id)
		if err != nil {
			if err.Error() == "advertiser not found: not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get advertiser: " + err.Error()})
			return
		}

		// Check if user has permission to view this advertiser
		hasAccess, err := h.checkAdvertiserAccess(c, advertiser.OrganizationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions: " + err.Error()})
			return
		}
		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view this advertiser"})
			return
		}

		c.JSON(http.StatusOK, advertiser)
	}
}

// UpdateAdvertiserRequest defines the request for updating an advertiser
// swagger:model
type UpdateAdvertiserRequest struct {
	// Advertiser name
	Name           string `json:"name" binding:"required" example:"Updated Advertiser"`
	// Contact email address
	ContactEmail   *string `json:"contact_email,omitempty" example:"updated@example.com"`
	// Billing details in JSON format
	// swagger:strfmt json
	BillingDetails *json.RawMessage `json:"billing_details,omitempty" swaggertype:"object"`
	// Status of the advertiser (active, pending, inactive, rejected)
	Status         string `json:"status" binding:"required" example:"active"`
	
	// Enhanced fields for Everflow integration
	InternalNotes              *string `json:"internal_notes,omitempty" example:"Updated internal notes"`
	DefaultCurrencyID          *string `json:"default_currency_id,omitempty" example:"USD"`
	PlatformName               *string `json:"platform_name,omitempty" example:"Updated Platform"`
	PlatformURL                *string `json:"platform_url,omitempty" example:"https://updated-platform.example.com"`
	PlatformUsername           *string `json:"platform_username,omitempty" example:"updated_user"`
	AccountingContactEmail     *string `json:"accounting_contact_email,omitempty" example:"updated-accounting@example.com"`
	OfferIDMacro               *string `json:"offer_id_macro,omitempty" example:"oid"`
	AffiliateIDMacro           *string `json:"affiliate_id_macro,omitempty" example:"aid"`
	AttributionMethod          *string `json:"attribution_method,omitempty" example:"last_click"`
	EmailAttributionMethod     *string `json:"email_attribution_method,omitempty" example:"last_click"`
	AttributionPriority        *string `json:"attribution_priority,omitempty" example:"medium"`
	ReportingTimezoneID        *int    `json:"reporting_timezone_id,omitempty" example:"67"`
	IsExposePublisherReporting *bool   `json:"is_expose_publisher_reporting,omitempty" example:"true"`
	
	// Synchronization options
	SyncToEverflow *bool `json:"sync_to_everflow,omitempty" example:"true"`
}

// UpdateAdvertiser updates an advertiser
// @Summary      Update advertiser
// @Description  Updates an advertiser with the given details
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        id       path      int                     true  "Advertiser ID"
// @Param        request  body      UpdateAdvertiserRequest  true  "Advertiser details"
// @Success      200      {object}  domain.Advertiser        "Updated advertiser"
// @Failure      400      {object}  map[string]string        "Invalid request"
// @Failure      403      {object}  map[string]string        "Forbidden - User doesn't have permission"
// @Failure      404      {object}  map[string]string        "Advertiser not found"
// @Failure      500      {object}  map[string]string        "Internal server error"
// @Security     BearerAuth
// @Router       /advertisers/{id} [put]
func (h *AdvertiserHandler) UpdateAdvertiser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	var req UpdateAdvertiserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get existing advertiser
	advertiser, err := h.advertiserService.GetAdvertiserByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "advertiser not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get advertiser: " + err.Error()})
		return
	}

	// Check if user has permission to update this advertiser
	hasAccess, err := h.checkAdvertiserAccess(c, advertiser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions: " + err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this advertiser"})
		return
	}

	// Update advertiser
	advertiser.Name = req.Name
	advertiser.ContactEmail = req.ContactEmail
	advertiser.Status = req.Status
	
	// Update enhanced fields for Everflow integration
	advertiser.InternalNotes = req.InternalNotes
	advertiser.DefaultCurrencyID = req.DefaultCurrencyID
	advertiser.PlatformName = req.PlatformName
	advertiser.PlatformURL = req.PlatformURL
	advertiser.PlatformUsername = req.PlatformUsername
	advertiser.AccountingContactEmail = req.AccountingContactEmail
	advertiser.OfferIDMacro = req.OfferIDMacro
	advertiser.AffiliateIDMacro = req.AffiliateIDMacro
	advertiser.AttributionMethod = req.AttributionMethod
	advertiser.EmailAttributionMethod = req.EmailAttributionMethod
	advertiser.AttributionPriority = req.AttributionPriority
	advertiser.ReportingTimezoneID = req.ReportingTimezoneID
	advertiser.IsExposePublisherReporting = req.IsExposePublisherReporting

	if req.BillingDetails != nil {
		var billingDetails domain.BillingDetails
		if err := json.Unmarshal(*req.BillingDetails, &billingDetails); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid billing details format: " + err.Error()})
			return
		}
		advertiser.BillingDetails = &billingDetails
	}

	if err := h.advertiserService.UpdateAdvertiser(c.Request.Context(), advertiser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update advertiser: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, advertiser)
}

// ListAdvertisersByOrganization retrieves a list of advertisers for an organization
// @Summary      List advertisers by organization
// @Description  Retrieves a list of advertisers for an organization with pagination
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        id  path      int                    true   "Organization ID"
// @Param        page           query     int                    false  "Page number (default: 1)"
// @Param        pageSize       query     int                    false  "Page size (default: 10)"
// @Success      200            {array}   domain.Advertiser      "List of advertisers"
// @Failure      400            {object}  map[string]string      "Invalid organization ID"
// @Failure      403            {object}  map[string]string      "Forbidden - User doesn't have permission"
// @Failure      500            {object}  map[string]string      "Internal server error"
// @Security     BearerAuth
// @Router       /organizations/{id}/advertisers [get]
func (h *AdvertiserHandler) ListAdvertisersByOrganization(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	// Check if user has permission to list advertisers for this organization
	hasAccess, err := h.checkAdvertiserAccess(c, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions: " + err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to list advertisers for this organization"})
		return
	}

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

	advertisers, err := h.advertiserService.ListAdvertisersByOrganization(c.Request.Context(), orgID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list advertisers: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, advertisers)
}

// DeleteAdvertiser deletes an advertiser
// @Summary      Delete advertiser
// @Description  Deletes an advertiser by its ID
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        id   path      int                true  "Advertiser ID"
// @Success      204  {object}  nil                "No content"
// @Failure      400  {object}  map[string]string  "Invalid advertiser ID"
// @Failure      403  {object}  map[string]string  "Forbidden - User doesn't have permission"
// @Failure      404  {object}  map[string]string  "Advertiser not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /advertisers/{id} [delete]
func (h *AdvertiserHandler) DeleteAdvertiser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	// Get the advertiser first to check permissions
	advertiser, err := h.advertiserService.GetAdvertiserByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "advertiser not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get advertiser: " + err.Error()})
		return
	}

	// Check if user has permission to delete this advertiser
	hasAccess, err := h.checkAdvertiserAccess(c, advertiser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions: " + err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this advertiser"})
		return
	}

	if err := h.advertiserService.DeleteAdvertiser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete advertiser: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateAdvertiserProviderMappingRequest defines the request for creating an advertiser provider mapping
// swagger:model
type CreateAdvertiserProviderMappingRequest struct {
	// Advertiser ID
	AdvertiserID         int64  `json:"advertiser_id" binding:"required" example:"1"`
	// Provider type (e.g., 'everflow')
	ProviderType         string `json:"provider_type" binding:"required" example:"everflow"`
	// Provider's advertiser ID
	ProviderAdvertiserID *string `json:"provider_advertiser_id,omitempty" example:"adv-12345"`
	// API credentials in JSON format
	// swagger:strfmt json
	APICredentials       *json.RawMessage `json:"api_credentials,omitempty" swaggertype:"object"`
	// Provider configuration in JSON format
	// swagger:strfmt json
	ProviderConfig       *json.RawMessage `json:"provider_config,omitempty" swaggertype:"object"`
}

// CreateAdvertiserProviderMapping creates a new advertiser provider mapping
// @Summary      Create a new advertiser provider mapping
// @Description  Creates a new mapping between an advertiser and a provider
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        request  body      CreateAdvertiserProviderMappingRequest  true  "Mapping details"
// @Success      201      {object}  domain.AdvertiserProviderMapping        "Created mapping"
// @Failure      400      {object}  map[string]string                       "Invalid request"
// @Failure      500      {object}  map[string]string                       "Internal server error"
// @Security     BearerAuth
// @Router       /advertiser-provider-mappings [post]
func (h *AdvertiserHandler) CreateAdvertiserProviderMapping(c *gin.Context) {
	var req CreateAdvertiserProviderMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	mapping := &domain.AdvertiserProviderMapping{
		AdvertiserID:         req.AdvertiserID,
		ProviderType:         req.ProviderType,
		ProviderAdvertiserID: req.ProviderAdvertiserID,
	}

	if req.APICredentials != nil {
		apiCredentialsStr := string(*req.APICredentials)
		mapping.APICredentials = &apiCredentialsStr
	}

	if req.ProviderConfig != nil {
		providerConfigStr := string(*req.ProviderConfig)
		mapping.ProviderConfig = &providerConfigStr
	}

	createdMapping, err := h.advertiserService.CreateAdvertiserProviderMapping(c.Request.Context(), mapping)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create advertiser provider mapping: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdMapping)
}

// GetAdvertiserProviderMapping retrieves an advertiser provider mapping
// @Summary      Get advertiser provider mapping
// @Description  Retrieves a mapping between an advertiser and a provider
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        id   path      int                                true  "Advertiser ID"
// @Param        providerType   path      string                             true  "Provider Type"
// @Success      200            {object}  domain.AdvertiserProviderMapping  "Mapping details"
// @Failure      400            {object}  map[string]string                 "Invalid request"
// @Failure      404            {object}  map[string]string                 "Mapping not found"
// @Failure      500            {object}  map[string]string                 "Internal server error"
// @Security     BearerAuth
// @Router       /advertisers/{id}/provider-mappings/{providerType} [get]
func (h *AdvertiserHandler) GetAdvertiserProviderMapping(c *gin.Context) {
	advertiserIDStr := c.Param("id")
	advertiserID, err := strconv.ParseInt(advertiserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	providerType := c.Param("providerType")
	if providerType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider type is required"})
		return
	}

	mapping, err := h.advertiserService.GetAdvertiserProviderMapping(c.Request.Context(), advertiserID, providerType)
	if err != nil {
		if err.Error() == "advertiser provider mapping not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser provider mapping not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get advertiser provider mapping: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapping)
}

// UpdateAdvertiserProviderMappingRequest defines the request for updating an advertiser provider mapping
// swagger:model
type UpdateAdvertiserProviderMappingRequest struct {
	// Provider's advertiser ID
	ProviderAdvertiserID *string `json:"provider_advertiser_id,omitempty" example:"adv-12345"`
	// API credentials in JSON format
	// swagger:strfmt json
	APICredentials       *json.RawMessage `json:"api_credentials,omitempty" swaggertype:"object"`
	// Provider configuration in JSON format
	// swagger:strfmt json
	ProviderConfig       *json.RawMessage `json:"provider_config,omitempty" swaggertype:"object"`
}

// UpdateAdvertiserProviderMapping updates an advertiser provider mapping
// @Summary      Update advertiser provider mapping
// @Description  Updates a mapping between an advertiser and a provider
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        mappingId  path      int                                      true  "Mapping ID"
// @Param        request    body      UpdateAdvertiserProviderMappingRequest  true  "Mapping details"
// @Success      200        {object}  domain.AdvertiserProviderMapping        "Updated mapping"
// @Failure      400        {object}  map[string]string                       "Invalid request"
// @Failure      404        {object}  map[string]string                       "Mapping not found"
// @Failure      500        {object}  map[string]string                       "Internal server error"
// @Security     BearerAuth
// @Router       /advertiser-provider-mappings/{mappingId} [put]
func (h *AdvertiserHandler) UpdateAdvertiserProviderMapping(c *gin.Context) {
	mappingIDStr := c.Param("mappingId")
	mappingID, err := strconv.ParseInt(mappingIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mapping ID"})
		return
	}

	var req UpdateAdvertiserProviderMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get existing mapping
	mapping, err := h.advertiserService.GetAdvertiserProviderMapping(c.Request.Context(), 0, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get advertiser provider mapping: " + err.Error()})
		return
	}

	// Update mapping
	mapping.MappingID = mappingID
	mapping.ProviderAdvertiserID = req.ProviderAdvertiserID

	if req.APICredentials != nil {
		apiCredentialsStr := string(*req.APICredentials)
		mapping.APICredentials = &apiCredentialsStr
	}

	if req.ProviderConfig != nil {
		providerConfigStr := string(*req.ProviderConfig)
		mapping.ProviderConfig = &providerConfigStr
	}

	if err := h.advertiserService.UpdateAdvertiserProviderMapping(c.Request.Context(), mapping); err != nil {
		if err.Error() == "advertiser provider mapping not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser provider mapping not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update advertiser provider mapping: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapping)
}

// DeleteAdvertiserProviderMapping deletes an advertiser provider mapping
// @Summary      Delete advertiser provider mapping
// @Description  Deletes a mapping between an advertiser and a provider
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        mappingId  path      int                true  "Mapping ID"
// @Success      204        {object}  nil                "No content"
// @Failure      400        {object}  map[string]string  "Invalid mapping ID"
// @Failure      404        {object}  map[string]string  "Mapping not found"
// @Failure      500        {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /advertiser-provider-mappings/{mappingId} [delete]
func (h *AdvertiserHandler) DeleteAdvertiserProviderMapping(c *gin.Context) {
	mappingIDStr := c.Param("mappingId")
	mappingID, err := strconv.ParseInt(mappingIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mapping ID"})
		return
	}

	if err := h.advertiserService.DeleteAdvertiserProviderMapping(c.Request.Context(), mappingID); err != nil {
		if err.Error() == "advertiser provider mapping not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser provider mapping not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete advertiser provider mapping: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SyncAdvertiserToEverflow synchronizes an advertiser to Everflow
// @Summary      Sync advertiser to Everflow
// @Description  Synchronizes an advertiser from local database to Everflow
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        id   path      int                true  "Advertiser ID"
// @Success      200  {object}  map[string]string  "Sync successful"
// @Failure      400  {object}  map[string]string  "Invalid advertiser ID"
// @Failure      403  {object}  map[string]string  "Forbidden - User doesn't have permission"
// @Failure      404  {object}  map[string]string  "Advertiser not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /advertisers/{id}/sync-to-everflow [post]
func (h *AdvertiserHandler) SyncAdvertiserToEverflow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	// Get advertiser to check permissions
	advertiser, err := h.advertiserService.GetAdvertiserByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "advertiser not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get advertiser: " + err.Error()})
		return
	}

	// Check if user has permission to sync this advertiser
	hasAccess, err := h.checkAdvertiserAccess(c, advertiser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions: " + err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to sync this advertiser"})
		return
	}

	// Sync to Everflow
	if err := h.advertiserService.SyncAdvertiserToEverflow(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync advertiser to Everflow: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Advertiser successfully synced to Everflow"})
}

// SyncAdvertiserFromEverflow synchronizes an advertiser from Everflow
// @Summary      Sync advertiser from Everflow
// @Description  Synchronizes an advertiser from Everflow to local database
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        id   path      int                true  "Advertiser ID"
// @Success      200  {object}  map[string]string  "Sync successful"
// @Failure      400  {object}  map[string]string  "Invalid advertiser ID"
// @Failure      403  {object}  map[string]string  "Forbidden - User doesn't have permission"
// @Failure      404  {object}  map[string]string  "Advertiser not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /advertisers/{id}/sync-from-everflow [post]
func (h *AdvertiserHandler) SyncAdvertiserFromEverflow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	// Get advertiser to check permissions
	advertiser, err := h.advertiserService.GetAdvertiserByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "advertiser not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get advertiser: " + err.Error()})
		return
	}

	// Check if user has permission to sync this advertiser
	hasAccess, err := h.checkAdvertiserAccess(c, advertiser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions: " + err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to sync this advertiser"})
		return
	}

	// Sync from Everflow
	if err := h.advertiserService.SyncAdvertiserFromEverflow(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync advertiser from Everflow: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Advertiser successfully synced from Everflow"})
}

// CompareAdvertiserWithEverflow compares an advertiser with Everflow data
// @Summary      Compare advertiser with Everflow
// @Description  Compares local advertiser data with Everflow data and returns discrepancies
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        id   path      int                                true  "Advertiser ID"
// @Success      200  {object}  []domain.AdvertiserDiscrepancy    "Discrepancies found"
// @Failure      400  {object}  map[string]string                 "Invalid advertiser ID"
// @Failure      403  {object}  map[string]string                 "Forbidden - User doesn't have permission"
// @Failure      404  {object}  map[string]string                 "Advertiser not found"
// @Failure      500  {object}  map[string]string                 "Internal server error"
// @Security     BearerAuth
// @Router       /advertisers/{id}/compare-with-everflow [get]
func (h *AdvertiserHandler) CompareAdvertiserWithEverflow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	// Get advertiser to check permissions
	advertiser, err := h.advertiserService.GetAdvertiserByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "advertiser not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get advertiser: " + err.Error()})
		return
	}

	// Check if user has permission to view this advertiser
	hasAccess, err := h.checkAdvertiserAccess(c, advertiser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions: " + err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view this advertiser"})
		return
	}

	// Compare with Everflow
	discrepancies, err := h.advertiserService.CompareAdvertiserWithEverflow(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to compare advertiser with Everflow: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"discrepancies": discrepancies,
		"total_count":   len(discrepancies),
	})
}