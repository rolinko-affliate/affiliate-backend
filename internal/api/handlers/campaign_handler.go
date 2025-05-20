package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// CampaignHandler handles campaign-related requests
type CampaignHandler struct {
	campaignService service.CampaignService
}

// NewCampaignHandler creates a new campaign handler
func NewCampaignHandler(cs service.CampaignService) *CampaignHandler {
	return &CampaignHandler{campaignService: cs}
}

// CreateCampaignRequest defines the request for creating a campaign
type CreateCampaignRequest struct {
	OrganizationID int64            `json:"organization_id" binding:"required"`
	AdvertiserID   int64            `json:"advertiser_id" binding:"required"`
	Name           string           `json:"name" binding:"required"`
	Description    *string          `json:"description,omitempty"`
	Status         string           `json:"status,omitempty"`
	StartDate      *string          `json:"start_date,omitempty"` // Format: YYYY-MM-DD
	EndDate        *string          `json:"end_date,omitempty"`   // Format: YYYY-MM-DD
}

// CreateCampaign creates a new campaign
// @Summary      Create a new campaign
// @Description  Creates a new campaign with the given details
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        request  body      CreateCampaignRequest  true  "Campaign details"
// @Success      201      {object}  domain.Campaign        "Created campaign"
// @Failure      400      {object}  map[string]string      "Invalid request"
// @Failure      500      {object}  map[string]string      "Internal server error"
// @Security     BearerAuth
// @Router       /campaigns [post]
func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var req CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	campaign := &domain.Campaign{
		OrganizationID: req.OrganizationID,
		AdvertiserID:   req.AdvertiserID,
		Name:           req.Name,
		Description:    req.Description,
		Status:         req.Status,
	}

	// Parse dates if provided
	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
			return
		}
		campaign.StartDate = &startDate
	}

	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
			return
		}
		campaign.EndDate = &endDate
	}

	createdCampaign, err := h.campaignService.CreateCampaign(c.Request.Context(), campaign)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdCampaign)
}

// GetCampaign retrieves a campaign by ID
// @Summary      Get campaign by ID
// @Description  Retrieves a campaign by its ID
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id   path      int               true  "Campaign ID"
// @Success      200  {object}  domain.Campaign  "Campaign details"
// @Failure      400  {object}  map[string]string "Invalid campaign ID"
// @Failure      404  {object}  map[string]string "Campaign not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Security     BearerAuth
// @Router       /campaigns/{id} [get]
func (h *CampaignHandler) GetCampaign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	campaign, err := h.campaignService.GetCampaignByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "campaign not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get campaign: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaign)
}

// UpdateCampaignRequest defines the request for updating a campaign
type UpdateCampaignRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	Status      string  `json:"status" binding:"required"`
	StartDate   *string `json:"start_date,omitempty"` // Format: YYYY-MM-DD
	EndDate     *string `json:"end_date,omitempty"`   // Format: YYYY-MM-DD
}

// UpdateCampaign updates a campaign
// @Summary      Update campaign
// @Description  Updates a campaign with the given details
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id       path      int                  true  "Campaign ID"
// @Param        request  body      UpdateCampaignRequest true  "Campaign details"
// @Success      200      {object}  domain.Campaign       "Updated campaign"
// @Failure      400      {object}  map[string]string     "Invalid request"
// @Failure      404      {object}  map[string]string     "Campaign not found"
// @Failure      500      {object}  map[string]string     "Internal server error"
// @Security     BearerAuth
// @Router       /campaigns/{id} [put]
func (h *CampaignHandler) UpdateCampaign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	var req UpdateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get existing campaign
	campaign, err := h.campaignService.GetCampaignByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "campaign not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get campaign: " + err.Error()})
		return
	}

	// Update campaign
	campaign.Name = req.Name
	campaign.Description = req.Description
	campaign.Status = req.Status

	// Parse dates if provided
	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
			return
		}
		campaign.StartDate = &startDate
	} else {
		campaign.StartDate = nil
	}

	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
			return
		}
		campaign.EndDate = &endDate
	} else {
		campaign.EndDate = nil
	}

	if err := h.campaignService.UpdateCampaign(c.Request.Context(), campaign); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campaign: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaign)
}

// ListCampaignsByOrganization retrieves a list of campaigns for an organization
// @Summary      List campaigns by organization
// @Description  Retrieves a list of campaigns for an organization with pagination
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id  path      int                 true   "Organization ID"
// @Param        page           query     int                 false  "Page number (default: 1)"
// @Param        pageSize       query     int                 false  "Page size (default: 10)"
// @Success      200            {array}   domain.Campaign     "List of campaigns"
// @Failure      400            {object}  map[string]string   "Invalid organization ID"
// @Failure      500            {object}  map[string]string   "Internal server error"
// @Security     BearerAuth
// @Router       /organizations/{id}/campaigns [get]
func (h *CampaignHandler) ListCampaignsByOrganization(c *gin.Context) {
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

	campaigns, err := h.campaignService.ListCampaignsByOrganization(c.Request.Context(), orgID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list campaigns: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}

// ListCampaignsByAdvertiser retrieves a list of campaigns for an advertiser
// @Summary      List campaigns by advertiser
// @Description  Retrieves a list of campaigns for an advertiser with pagination
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id  path      int                 true   "Advertiser ID"
// @Param        page         query     int                 false  "Page number (default: 1)"
// @Param        pageSize     query     int                 false  "Page size (default: 10)"
// @Success      200          {array}   domain.Campaign     "List of campaigns"
// @Failure      400          {object}  map[string]string   "Invalid advertiser ID"
// @Failure      500          {object}  map[string]string   "Internal server error"
// @Security     BearerAuth
// @Router       /advertisers/{id}/campaigns [get]
func (h *CampaignHandler) ListCampaignsByAdvertiser(c *gin.Context) {
	advertiserIDStr := c.Param("id")
	advertiserID, err := strconv.ParseInt(advertiserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
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

	campaigns, err := h.campaignService.ListCampaignsByAdvertiser(c.Request.Context(), advertiserID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list campaigns: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}

// DeleteCampaign deletes a campaign
// @Summary      Delete campaign
// @Description  Deletes a campaign by its ID
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id   path      int                true  "Campaign ID"
// @Success      204  {object}  nil                "No content"
// @Failure      400  {object}  map[string]string  "Invalid campaign ID"
// @Failure      404  {object}  map[string]string  "Campaign not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /campaigns/{id} [delete]
func (h *CampaignHandler) DeleteCampaign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	if err := h.campaignService.DeleteCampaign(c.Request.Context(), id); err != nil {
		if err.Error() == "campaign not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete campaign: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateCampaignProviderOfferRequest defines the request for creating a campaign provider offer
type CreateCampaignProviderOfferRequest struct {
	CampaignID         int64            `json:"campaign_id" binding:"required"`
	ProviderType       string           `json:"provider_type" binding:"required"`
	ProviderOfferRef   *string          `json:"provider_offer_ref,omitempty"`
	ProviderOfferConfig *json.RawMessage `json:"provider_offer_config,omitempty"`
	IsActiveOnProvider bool             `json:"is_active_on_provider"`
}

// CreateCampaignProviderOffer creates a new campaign provider offer
// @Summary      Create a new campaign provider offer
// @Description  Creates a new offer for a campaign on a provider
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        request  body      CreateCampaignProviderOfferRequest  true  "Offer details"
// @Success      201      {object}  domain.CampaignProviderOffer        "Created offer"
// @Failure      400      {object}  map[string]string                   "Invalid request"
// @Failure      500      {object}  map[string]string                   "Internal server error"
// @Security     BearerAuth
// @Router       /campaign-provider-offers [post]
func (h *CampaignHandler) CreateCampaignProviderOffer(c *gin.Context) {
	var req CreateCampaignProviderOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	offer := &domain.CampaignProviderOffer{
		CampaignID:        req.CampaignID,
		ProviderType:      req.ProviderType,
		ProviderOfferRef:  req.ProviderOfferRef,
		IsActiveOnProvider: req.IsActiveOnProvider,
	}

	if req.ProviderOfferConfig != nil {
		providerOfferConfigStr := string(*req.ProviderOfferConfig)
		offer.ProviderOfferConfig = &providerOfferConfigStr
	}

	createdOffer, err := h.campaignService.CreateCampaignProviderOffer(c.Request.Context(), offer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign provider offer: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdOffer)
}

// GetCampaignProviderOffer retrieves a campaign provider offer by ID
// @Summary      Get campaign provider offer by ID
// @Description  Retrieves a campaign provider offer by its ID
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id   path      int                          true  "Offer ID"
// @Success      200  {object}  domain.CampaignProviderOffer "Offer details"
// @Failure      400  {object}  map[string]string            "Invalid offer ID"
// @Failure      404  {object}  map[string]string            "Offer not found"
// @Failure      500  {object}  map[string]string            "Internal server error"
// @Security     BearerAuth
// @Router       /campaign-provider-offers/{id} [get]
func (h *CampaignHandler) GetCampaignProviderOffer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offer ID"})
		return
	}

	offer, err := h.campaignService.GetCampaignProviderOfferByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "campaign provider offer not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign provider offer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get campaign provider offer: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, offer)
}

// UpdateCampaignProviderOfferRequest defines the request for updating a campaign provider offer
type UpdateCampaignProviderOfferRequest struct {
	ProviderOfferRef   *string          `json:"provider_offer_ref,omitempty"`
	ProviderOfferConfig *json.RawMessage `json:"provider_offer_config,omitempty"`
	IsActiveOnProvider bool             `json:"is_active_on_provider"`
}

// UpdateCampaignProviderOffer updates a campaign provider offer
// @Summary      Update campaign provider offer
// @Description  Updates a campaign provider offer with the given details
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id       path      int                                true  "Offer ID"
// @Param        request  body      UpdateCampaignProviderOfferRequest true  "Offer details"
// @Success      200      {object}  domain.CampaignProviderOffer       "Updated offer"
// @Failure      400      {object}  map[string]string                  "Invalid request"
// @Failure      404      {object}  map[string]string                  "Offer not found"
// @Failure      500      {object}  map[string]string                  "Internal server error"
// @Security     BearerAuth
// @Router       /campaign-provider-offers/{id} [put]
func (h *CampaignHandler) UpdateCampaignProviderOffer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offer ID"})
		return
	}

	var req UpdateCampaignProviderOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get existing offer
	offer, err := h.campaignService.GetCampaignProviderOfferByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "campaign provider offer not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign provider offer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get campaign provider offer: " + err.Error()})
		return
	}

	// Update offer
	offer.ProviderOfferRef = req.ProviderOfferRef
	offer.IsActiveOnProvider = req.IsActiveOnProvider

	if req.ProviderOfferConfig != nil {
		providerOfferConfigStr := string(*req.ProviderOfferConfig)
		offer.ProviderOfferConfig = &providerOfferConfigStr
	}

	if err := h.campaignService.UpdateCampaignProviderOffer(c.Request.Context(), offer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campaign provider offer: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, offer)
}

// ListCampaignProviderOffersByCampaign retrieves a list of campaign provider offers for a campaign
// @Summary      List campaign provider offers by campaign
// @Description  Retrieves a list of provider offers for a campaign
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id  path      int                            true  "Campaign ID"
// @Success      200         {array}   domain.CampaignProviderOffer  "List of offers"
// @Failure      400         {object}  map[string]string             "Invalid campaign ID"
// @Failure      500         {object}  map[string]string             "Internal server error"
// @Security     BearerAuth
// @Router       /campaigns/{id}/provider-offers [get]
func (h *CampaignHandler) ListCampaignProviderOffersByCampaign(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := strconv.ParseInt(campaignIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	offers, err := h.campaignService.ListCampaignProviderOffersByCampaign(c.Request.Context(), campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list campaign provider offers: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, offers)
}

// DeleteCampaignProviderOffer deletes a campaign provider offer
// @Summary      Delete campaign provider offer
// @Description  Deletes a campaign provider offer by its ID
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id   path      int                true  "Offer ID"
// @Success      204  {object}  nil                "No content"
// @Failure      400  {object}  map[string]string  "Invalid offer ID"
// @Failure      404  {object}  map[string]string  "Offer not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /campaign-provider-offers/{id} [delete]
func (h *CampaignHandler) DeleteCampaignProviderOffer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offer ID"})
		return
	}

	if err := h.campaignService.DeleteCampaignProviderOffer(c.Request.Context(), id); err != nil {
		if err.Error() == "campaign provider offer not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign provider offer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete campaign provider offer: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}