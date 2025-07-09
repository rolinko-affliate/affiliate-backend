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

// AffiliateHandler handles affiliate-related requests
type AffiliateHandler struct {
	affiliateService service.AffiliateService
	profileService   service.ProfileService
}

// NewAffiliateHandler creates a new affiliate handler
func NewAffiliateHandler(as service.AffiliateService, ps service.ProfileService) *AffiliateHandler {
	return &AffiliateHandler{
		affiliateService: as,
		profileService:   ps,
	}
}

// checkAffiliateAccess verifies if the user has permission to access/modify the affiliate
// Returns true if the user has access, false otherwise
func (h *AffiliateHandler) checkAffiliateAccess(c *gin.Context, affiliateOrgID int64) (bool, error) {
	// Get user role from context
	userRole, exists := c.Get(middleware.UserRoleKey)
	if !exists {
		return false, fmt.Errorf("user role not found in context")
	}
	
	// Admin can access all affiliates
	if userRole.(string) == "Admin" {
		return true, nil
	}
	
	// Get user's organization ID from context
	userOrgID, exists := c.Get("organizationID")
	if !exists {
		return false, fmt.Errorf("user organization ID not found in context")
	}
	
	// Check if user belongs to the same organization as the affiliate
	return userOrgID.(int64) == affiliateOrgID, nil
}

// CreateAffiliateRequest defines the request for creating an affiliate
// swagger:model
type CreateAffiliateRequest struct {
	// Organization ID
	OrganizationID int64  `json:"organization_id" binding:"required" example:"1"`
	// Affiliate name
	Name           string `json:"name" binding:"required" example:"Example Affiliate"`
	// Contact email address
	ContactEmail   *string `json:"contact_email,omitempty" example:"affiliate@example.com"`
	// Payment details in JSON format
	// swagger:strfmt json
	PaymentDetails *json.RawMessage `json:"payment_details,omitempty" swaggertype:"object"`
	// Status of the affiliate (active, pending, inactive, rejected)
	Status         string `json:"status,omitempty" example:"active"`
}

// CreateAffiliate creates a new affiliate
// @Summary      Create a new affiliate
// @Description  Creates a new affiliate with the given details
// @Tags         affiliates
// @Accept       json
// @Produce      json
// @Param        request  body      CreateAffiliateRequest  true  "Affiliate details"
// @Success      201      {object}  domain.Affiliate        "Created affiliate"
// @Failure      400      {object}  map[string]string       "Invalid request"
// @Failure      500      {object}  map[string]string       "Internal server error"
// @Security     BearerAuth
// @Router       /affiliates [post]
func (h *AffiliateHandler) CreateAffiliate(c *gin.Context) {
	var req CreateAffiliateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	affiliate := &domain.Affiliate{
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		ContactEmail:   req.ContactEmail,
		Status:         req.Status,
	}

	if req.PaymentDetails != nil {
		paymentDetailsStr := string(*req.PaymentDetails)
		affiliate.PaymentDetails = &paymentDetailsStr
	}

	createdAffiliate, err := h.affiliateService.CreateAffiliate(c.Request.Context(), affiliate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create affiliate: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdAffiliate)
}

// GetAffiliate retrieves an affiliate by ID
// @Summary      Get affiliate by ID
// @Description  Retrieves an affiliate by its ID
// @Tags         affiliates
// @Accept       json
// @Produce      json
// @Param        id   path      int                true  "Affiliate ID"
// @Success      200  {object}  domain.Affiliate  "Affiliate details"
// @Failure      400  {object}  map[string]string "Invalid affiliate ID"
// @Failure      404  {object}  map[string]string "Affiliate not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Security     BearerAuth
// @Router       /affiliates/{id} [get]
func (h *AffiliateHandler) GetAffiliate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid affiliate ID"})
		return
	}

	affiliate, err := h.affiliateService.GetAffiliateByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "affiliate not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Affiliate not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get affiliate: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, affiliate)
}

// UpdateAffiliateRequest defines the request for updating an affiliate
// swagger:model
type UpdateAffiliateRequest struct {
	// Affiliate name
	Name           string `json:"name" binding:"required" example:"Updated Affiliate"`
	// Contact email address
	ContactEmail   *string `json:"contact_email,omitempty" example:"updated@example.com"`
	// Payment details in JSON format
	// swagger:strfmt json
	PaymentDetails *json.RawMessage `json:"payment_details,omitempty" swaggertype:"object"`
	// Status of the affiliate (active, pending, inactive, rejected)
	Status         string `json:"status" binding:"required" example:"active"`
}

// UpdateAffiliate updates an affiliate
// @Summary      Update affiliate
// @Description  Updates an affiliate with the given details
// @Tags         affiliates
// @Accept       json
// @Produce      json
// @Param        id       path      int                    true  "Affiliate ID"
// @Param        request  body      UpdateAffiliateRequest  true  "Affiliate details"
// @Success      200      {object}  domain.Affiliate        "Updated affiliate"
// @Failure      400      {object}  map[string]string       "Invalid request"
// @Failure      404      {object}  map[string]string       "Affiliate not found"
// @Failure      500      {object}  map[string]string       "Internal server error"
// @Security     BearerAuth
// @Router       /affiliates/{id} [put]
func (h *AffiliateHandler) UpdateAffiliate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid affiliate ID"})
		return
	}

	var req UpdateAffiliateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get existing affiliate
	affiliate, err := h.affiliateService.GetAffiliateByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "affiliate not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Affiliate not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get affiliate: " + err.Error()})
		return
	}

	// Update affiliate
	affiliate.Name = req.Name
	affiliate.ContactEmail = req.ContactEmail
	affiliate.Status = req.Status

	if req.PaymentDetails != nil {
		paymentDetailsStr := string(*req.PaymentDetails)
		affiliate.PaymentDetails = &paymentDetailsStr
	}

	if err := h.affiliateService.UpdateAffiliate(c.Request.Context(), affiliate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update affiliate: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, affiliate)
}

// ListAffiliatesByOrganization retrieves a list of affiliates for an organization
// @Summary      List affiliates by organization
// @Description  Retrieves a list of affiliates for an organization with pagination
// @Tags         affiliates
// @Accept       json
// @Produce      json
// @Param        id  path      int                   true   "Organization ID"
// @Param        page           query     int                   false  "Page number (default: 1)"
// @Param        pageSize       query     int                   false  "Page size (default: 10)"
// @Success      200            {array}   domain.Affiliate      "List of affiliates"
// @Failure      400            {object}  map[string]string     "Invalid organization ID"
// @Failure      500            {object}  map[string]string     "Internal server error"
// @Security     BearerAuth
// @Router       /organizations/{id}/affiliates [get]
func (h *AffiliateHandler) ListAffiliatesByOrganization(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
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

	affiliates, err := h.affiliateService.ListAffiliatesByOrganization(c.Request.Context(), orgID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list affiliates: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, affiliates)
}

// DeleteAffiliate deletes an affiliate
// @Summary      Delete affiliate
// @Description  Deletes an affiliate by its ID
// @Tags         affiliates
// @Accept       json
// @Produce      json
// @Param        id   path      int                true  "Affiliate ID"
// @Success      204  {object}  nil                "No content"
// @Failure      400  {object}  map[string]string  "Invalid affiliate ID"
// @Failure      404  {object}  map[string]string  "Affiliate not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /affiliates/{id} [delete]
func (h *AffiliateHandler) DeleteAffiliate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid affiliate ID"})
		return
	}

	if err := h.affiliateService.DeleteAffiliate(c.Request.Context(), id); err != nil {
		if err.Error() == "affiliate not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Affiliate not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete affiliate: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateAffiliateProviderMappingRequest defines the request for creating an affiliate provider mapping
// swagger:model
type CreateAffiliateProviderMappingRequest struct {
	// Affiliate ID
	AffiliateID         int64  `json:"affiliate_id" binding:"required" example:"1"`
	// Provider type (e.g., 'everflow')
	ProviderType        string `json:"provider_type" binding:"required" example:"everflow"`
	// Provider's affiliate ID
	ProviderAffiliateID *string `json:"provider_affiliate_id,omitempty" example:"aff-12345"`
	// Provider configuration in JSON format
	// swagger:strfmt json
	ProviderConfig      *json.RawMessage `json:"provider_config,omitempty" swaggertype:"object"`
}

// CreateAffiliateProviderMapping creates a new affiliate provider mapping
// @Summary      Create a new affiliate provider mapping
// @Description  Creates a new mapping between an affiliate and a provider
// @Tags         affiliates
// @Accept       json
// @Produce      json
// @Param        request  body      CreateAffiliateProviderMappingRequest  true  "Mapping details"
// @Success      201      {object}  domain.AffiliateProviderMapping        "Created mapping"
// @Failure      400      {object}  map[string]string                      "Invalid request"
// @Failure      500      {object}  map[string]string                      "Internal server error"
// @Security     BearerAuth
// @Router       /affiliate-provider-mappings [post]
func (h *AffiliateHandler) CreateAffiliateProviderMapping(c *gin.Context) {
	var req CreateAffiliateProviderMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	mapping := &domain.AffiliateProviderMapping{
		AffiliateID:         req.AffiliateID,
		ProviderType:        req.ProviderType,
		ProviderAffiliateID: req.ProviderAffiliateID,
	}

	if req.ProviderConfig != nil {
		providerConfigStr := string(*req.ProviderConfig)
		mapping.ProviderConfig = &providerConfigStr
	}

	createdMapping, err := h.affiliateService.CreateAffiliateProviderMapping(c.Request.Context(), mapping)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create affiliate provider mapping: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdMapping)
}

// GetAffiliateProviderMapping retrieves an affiliate provider mapping
// @Summary      Get affiliate provider mapping
// @Description  Retrieves a mapping between an affiliate and a provider
// @Tags         affiliates
// @Accept       json
// @Produce      json
// @Param        id   path      int                               true  "Affiliate ID"
// @Param        providerType  path      string                            true  "Provider Type"
// @Success      200           {object}  domain.AffiliateProviderMapping  "Mapping details"
// @Failure      400           {object}  map[string]string                "Invalid request"
// @Failure      404           {object}  map[string]string                "Mapping not found"
// @Failure      500           {object}  map[string]string                "Internal server error"
// @Security     BearerAuth
// @Router       /affiliates/{id}/provider-mappings/{providerType} [get]
func (h *AffiliateHandler) GetAffiliateProviderMapping(c *gin.Context) {
	affiliateIDStr := c.Param("id")
	affiliateID, err := strconv.ParseInt(affiliateIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid affiliate ID"})
		return
	}

	providerType := c.Param("providerType")
	if providerType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider type is required"})
		return
	}

	mapping, err := h.affiliateService.GetAffiliateProviderMapping(c.Request.Context(), affiliateID, providerType)
	if err != nil {
		if err.Error() == "affiliate provider mapping not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Affiliate provider mapping not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get affiliate provider mapping: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapping)
}

// UpdateAffiliateProviderMappingRequest defines the request for updating an affiliate provider mapping
// swagger:model
type UpdateAffiliateProviderMappingRequest struct {
	// Provider's affiliate ID
	ProviderAffiliateID *string `json:"provider_affiliate_id,omitempty" example:"aff-12345"`
	// Provider configuration in JSON format
	// swagger:strfmt json
	ProviderConfig      *json.RawMessage `json:"provider_config,omitempty" swaggertype:"object"`
}

// UpdateAffiliateProviderMapping updates an affiliate provider mapping
// @Summary      Update affiliate provider mapping
// @Description  Updates a mapping between an affiliate and a provider
// @Tags         affiliates
// @Accept       json
// @Produce      json
// @Param        mappingId  path      int                                     true  "Mapping ID"
// @Param        request    body      UpdateAffiliateProviderMappingRequest  true  "Mapping details"
// @Success      200        {object}  domain.AffiliateProviderMapping        "Updated mapping"
// @Failure      400        {object}  map[string]string                      "Invalid request"
// @Failure      404        {object}  map[string]string                      "Mapping not found"
// @Failure      500        {object}  map[string]string                      "Internal server error"
// @Security     BearerAuth
// @Router       /affiliate-provider-mappings/{mappingId} [put]
func (h *AffiliateHandler) UpdateAffiliateProviderMapping(c *gin.Context) {
	mappingIDStr := c.Param("mappingId")
	mappingID, err := strconv.ParseInt(mappingIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mapping ID"})
		return
	}

	var req UpdateAffiliateProviderMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get existing mapping
	mapping, err := h.affiliateService.GetAffiliateProviderMapping(c.Request.Context(), 0, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get affiliate provider mapping: " + err.Error()})
		return
	}

	// Update mapping
	mapping.MappingID = mappingID
	mapping.ProviderAffiliateID = req.ProviderAffiliateID

	if req.ProviderConfig != nil {
		providerConfigStr := string(*req.ProviderConfig)
		mapping.ProviderConfig = &providerConfigStr
	}

	if err := h.affiliateService.UpdateAffiliateProviderMapping(c.Request.Context(), mapping); err != nil {
		if err.Error() == "affiliate provider mapping not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Affiliate provider mapping not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update affiliate provider mapping: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapping)
}

// DeleteAffiliateProviderMapping deletes an affiliate provider mapping
// @Summary      Delete affiliate provider mapping
// @Description  Deletes a mapping between an affiliate and a provider
// @Tags         affiliates
// @Accept       json
// @Produce      json
// @Param        mappingId  path      int                true  "Mapping ID"
// @Success      204        {object}  nil                "No content"
// @Failure      400        {object}  map[string]string  "Invalid mapping ID"
// @Failure      404        {object}  map[string]string  "Mapping not found"
// @Failure      500        {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /affiliate-provider-mappings/{mappingId} [delete]
func (h *AffiliateHandler) DeleteAffiliateProviderMapping(c *gin.Context) {
	mappingIDStr := c.Param("mappingId")
	mappingID, err := strconv.ParseInt(mappingIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mapping ID"})
		return
	}

	if err := h.affiliateService.DeleteAffiliateProviderMapping(c.Request.Context(), mappingID); err != nil {
		if err.Error() == "affiliate provider mapping not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Affiliate provider mapping not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete affiliate provider mapping: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}