package handlers

import (
	"net/http"
	"strconv"

	"github.com/affiliate-backend/internal/api/models"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// CampaignHandler handles campaign-related requests following clean architecture
type CampaignHandler struct {
	campaignService service.CampaignService
}

// NewCampaignHandler creates a new campaign handler
func NewCampaignHandler(campaignService service.CampaignService) *CampaignHandler {
	return &CampaignHandler{
		campaignService: campaignService,
	}
}

// CreateCampaign creates a new campaign
// @Summary Create a new campaign
// @Description Create a new campaign with the provided details
// @Tags campaigns
// @Accept json
// @Produce json
// @Param campaign body models.CreateCampaignRequest true "Campaign creation request"
// @Success 201 {object} models.CampaignResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /campaigns [post]
func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var req models.CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Convert to domain model
	campaign := req.ToCampaignDomain()

	// Create campaign
	if err := h.campaignService.CreateCampaign(c.Request.Context(), campaign); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create campaign",
			Details: err.Error(),
		})
		return
	}

	// Return response
	response := models.FromCampaignDomain(campaign)
	c.JSON(http.StatusCreated, response)
}

// GetCampaign retrieves a campaign by ID
// @Summary Get a campaign by ID
// @Description Retrieve a campaign by its ID
// @Tags campaigns
// @Produce json
// @Param id path int true "Campaign ID"
// @Success 200 {object} models.CampaignResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /campaigns/{id} [get]
func (h *CampaignHandler) GetCampaign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid campaign ID",
			Details: "Campaign ID must be a valid integer",
		})
		return
	}

	campaign, err := h.campaignService.GetCampaignByID(c.Request.Context(), id)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Campaign not found",
				Details: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get campaign",
			Details: err.Error(),
		})
		return
	}

	response := models.FromCampaignDomain(campaign)
	c.JSON(http.StatusOK, response)
}

// UpdateCampaign updates an existing campaign
// @Summary Update a campaign
// @Description Update an existing campaign with the provided details
// @Tags campaigns
// @Accept json
// @Produce json
// @Param id path int true "Campaign ID"
// @Param campaign body models.UpdateCampaignRequest true "Campaign update request"
// @Success 200 {object} models.CampaignResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /campaigns/{id} [put]
func (h *CampaignHandler) UpdateCampaign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid campaign ID",
			Details: "Campaign ID must be a valid integer",
		})
		return
	}

	var req models.UpdateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Get existing campaign
	campaign, err := h.campaignService.GetCampaignByID(c.Request.Context(), id)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Campaign not found",
				Details: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get campaign",
			Details: err.Error(),
		})
		return
	}

	// Update campaign with request data
	req.UpdateCampaignDomain(campaign)

	// Update campaign
	if err := h.campaignService.UpdateCampaign(c.Request.Context(), campaign); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update campaign",
			Details: err.Error(),
		})
		return
	}

	// Return response
	response := models.FromCampaignDomain(campaign)
	c.JSON(http.StatusOK, response)
}

// ListCampaignsByAdvertiser lists campaigns for a specific advertiser
// @Summary List campaigns by advertiser
// @Description Retrieve campaigns for a specific advertiser with pagination
// @Tags campaigns
// @Produce json
// @Param advertiser_id path int true "Advertiser ID"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Success 200 {object} models.CampaignListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /advertisers/{advertiser_id}/campaigns [get]
func (h *CampaignHandler) ListCampaignsByAdvertiser(c *gin.Context) {
	advertiserIDStr := c.Param("advertiser_id")
	advertiserID, err := strconv.ParseInt(advertiserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid advertiser ID",
			Details: "Advertiser ID must be a valid integer",
		})
		return
	}

	page, pageSize := getPaginationParams(c)
	offset := (page - 1) * pageSize

	campaigns, err := h.campaignService.ListCampaignsByAdvertiser(c.Request.Context(), advertiserID, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to list campaigns",
			Details: err.Error(),
		})
		return
	}

	// For simplicity, we're not implementing total count here
	// In a real application, you'd want to get the total count for proper pagination
	response := models.FromCampaignDomainList(campaigns, len(campaigns), page, pageSize)
	c.JSON(http.StatusOK, response)
}

// ListCampaignsByOrganization lists campaigns for a specific organization
// @Summary List campaigns by organization
// @Description Retrieve campaigns for a specific organization with pagination
// @Tags campaigns
// @Produce json
// @Param organization_id path int true "Organization ID"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Success 200 {object} models.CampaignListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /organizations/{organization_id}/campaigns [get]
func (h *CampaignHandler) ListCampaignsByOrganization(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid organization ID",
			Details: "Organization ID must be a valid integer",
		})
		return
	}

	page, pageSize := getPaginationParams(c)
	offset := (page - 1) * pageSize

	campaigns, err := h.campaignService.ListCampaignsByOrganization(c.Request.Context(), orgID, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to list campaigns",
			Details: err.Error(),
		})
		return
	}

	// For simplicity, we're not implementing total count here
	// In a real application, you'd want to get the total count for proper pagination
	response := models.FromCampaignDomainList(campaigns, len(campaigns), page, pageSize)
	c.JSON(http.StatusOK, response)
}

// DeleteCampaign deletes a campaign by ID
// @Summary Delete a campaign
// @Description Delete a campaign by its ID
// @Tags campaigns
// @Param id path int true "Campaign ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /campaigns/{id} [delete]
func (h *CampaignHandler) DeleteCampaign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid campaign ID",
			Details: "Campaign ID must be a valid integer",
		})
		return
	}

	if err := h.campaignService.DeleteCampaign(c.Request.Context(), id); err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Campaign not found",
				Details: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to delete campaign",
			Details: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
