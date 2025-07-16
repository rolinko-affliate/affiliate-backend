package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/affiliate-backend/internal/api/models"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// TrackingLinkHandler handles HTTP requests for tracking links
type TrackingLinkHandler struct {
	trackingLinkService service.TrackingLinkService
}

// NewTrackingLinkHandler creates a new tracking link handler
func NewTrackingLinkHandler(trackingLinkService service.TrackingLinkService) *TrackingLinkHandler {
	return &TrackingLinkHandler{
		trackingLinkService: trackingLinkService,
	}
}

// CreateTrackingLink creates a new tracking link
// @Summary Create a new tracking link
// @Description Create a new tracking link for a campaign and affiliate
// @Tags tracking-links
// @Accept json
// @Produce json
// @Param organization_id path int true "Organization ID"
// @Param request body models.TrackingLinkRequest true "Tracking link creation request"
// @Success 201 {object} models.TrackingLinkResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{organization_id}/tracking-links [post]
func (h *TrackingLinkHandler) CreateTrackingLink(c *gin.Context) {
	organizationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid organization ID",
			Details: "Organization ID must be a valid integer",
		})
		return
	}

	var req models.TrackingLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Convert to domain model
	trackingLink := req.ToTrackingLinkDomain(organizationID)

	// Create tracking link
	if err := h.trackingLinkService.CreateTrackingLink(c.Request.Context(), trackingLink); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create tracking link",
			Details: err.Error(),
		})
		return
	}

	// Convert to response model
	response := models.FromTrackingLinkDomain(trackingLink)
	c.JSON(http.StatusCreated, response)
}

// GenerateTrackingLink generates a new tracking link with provider integration
// @Summary Generate a new tracking link
// @Description Generate a new tracking link with provider integration for a campaign and affiliate
// @Tags tracking-links
// @Accept json
// @Produce json
// @Param organization_id path int true "Organization ID"
// @Param request body models.TrackingLinkGenerationRequest true "Tracking link generation request"
// @Success 201 {object} models.TrackingLinkGenerationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{organization_id}/tracking-links/generate [post]
func (h *TrackingLinkHandler) GenerateTrackingLink(c *gin.Context) {
	_, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid organization ID",
			Details: "Organization ID must be a valid integer",
		})
		return
	}

	var req models.TrackingLinkGenerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Convert to domain model
	domainReq := req.ToTrackingLinkGenerationDomain()

	// Generate tracking link
	response, err := h.trackingLinkService.GenerateTrackingLink(c.Request.Context(), domainReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to generate tracking link",
			Details: err.Error(),
		})
		return
	}

	// Convert to API response
	baseURL := getBaseURL(c)
	apiResponse := models.FromTrackingLinkGenerationDomain(response, baseURL)
	c.JSON(http.StatusCreated, apiResponse)
}

// GetTrackingLink retrieves a tracking link by ID
// @Summary Get a tracking link
// @Description Retrieve a tracking link by its ID
// @Tags tracking-links
// @Produce json
// @Param organization_id path int true "Organization ID"
// @Param tracking_link_id path int true "Tracking Link ID"
// @Success 200 {object} models.TrackingLinkResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{organization_id}/tracking-links/{tracking_link_id} [get]
func (h *TrackingLinkHandler) GetTrackingLink(c *gin.Context) {
	trackingLinkID, err := strconv.ParseInt(c.Param("link_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid tracking link ID",
			Details: "Tracking link ID must be a valid integer",
		})
		return
	}

	trackingLink, err := h.trackingLinkService.GetTrackingLinkByID(c.Request.Context(), trackingLinkID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Tracking link not found",
			Details: err.Error(),
		})
		return
	}

	response := models.FromTrackingLinkDomain(trackingLink)
	c.JSON(http.StatusOK, response)
}

// UpdateTrackingLink updates an existing tracking link
// @Summary Update a tracking link
// @Description Update an existing tracking link
// @Tags tracking-links
// @Accept json
// @Produce json
// @Param organization_id path int true "Organization ID"
// @Param tracking_link_id path int true "Tracking Link ID"
// @Param request body models.TrackingLinkUpdateRequest true "Tracking link update request"
// @Success 200 {object} models.TrackingLinkResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{organization_id}/tracking-links/{tracking_link_id} [put]
func (h *TrackingLinkHandler) UpdateTrackingLink(c *gin.Context) {
	organizationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid organization ID",
			Details: "Organization ID must be a valid integer",
		})
		return
	}

	trackingLinkID, err := strconv.ParseInt(c.Param("link_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid tracking link ID",
			Details: "Tracking link ID must be a valid integer",
		})
		return
	}

	var req models.TrackingLinkUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Get existing tracking link
	existingTrackingLink, err := h.trackingLinkService.GetTrackingLinkByID(c.Request.Context(), trackingLinkID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Tracking link not found",
			Details: err.Error(),
		})
		return
	}

	// Verify tracking link belongs to the organization
	if existingTrackingLink.OrganizationID != organizationID {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Tracking link not found",
			Details: "Tracking link does not belong to the specified organization",
		})
		return
	}

	// Update fields from request
	req.UpdateTrackingLinkDomain(existingTrackingLink)

	// Update tracking link
	if err := h.trackingLinkService.UpdateTrackingLink(c.Request.Context(), existingTrackingLink); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update tracking link",
			Details: err.Error(),
		})
		return
	}

	// Get updated tracking link
	updatedTrackingLink, err := h.trackingLinkService.GetTrackingLinkByID(c.Request.Context(), trackingLinkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve updated tracking link",
			Details: err.Error(),
		})
		return
	}

	response := models.FromTrackingLinkDomain(updatedTrackingLink)
	c.JSON(http.StatusOK, response)
}

// DeleteTrackingLink deletes a tracking link
// @Summary Delete a tracking link
// @Description Delete a tracking link by its ID
// @Tags tracking-links
// @Param organization_id path int true "Organization ID"
// @Param tracking_link_id path int true "Tracking Link ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{organization_id}/tracking-links/{tracking_link_id} [delete]
func (h *TrackingLinkHandler) DeleteTrackingLink(c *gin.Context) {
	trackingLinkID, err := strconv.ParseInt(c.Param("link_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid tracking link ID",
			Details: "Tracking link ID must be a valid integer",
		})
		return
	}

	if err := h.trackingLinkService.DeleteTrackingLink(c.Request.Context(), trackingLinkID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to delete tracking link",
			Details: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListTrackingLinksByCampaign lists tracking links for a specific campaign
// @Summary List tracking links by campaign
// @Description Retrieve a list of tracking links for a specific campaign
// @Tags tracking-links
// @Produce json
// @Param id path int true "Campaign ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} models.TrackingLinkListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /campaigns/{id}/tracking-links [get]
func (h *TrackingLinkHandler) ListTrackingLinksByCampaign(c *gin.Context) {
	campaignID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid campaign ID",
			Details: "Campaign ID must be a valid integer",
		})
		return
	}

	page, pageSize := getPaginationParams(c)
	offset := (page - 1) * pageSize

	trackingLinks, err := h.trackingLinkService.ListTrackingLinksByCampaign(c.Request.Context(), campaignID, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to list tracking links",
			Details: err.Error(),
		})
		return
	}

	// Convert to response models
	var responses []*models.TrackingLinkResponse
	for _, trackingLink := range trackingLinks {
		responses = append(responses, models.FromTrackingLinkDomain(trackingLink))
	}

	// Calculate total pages (simplified - in production you'd get total count from service)
	totalPages := (len(responses) + pageSize - 1) / pageSize

	response := models.TrackingLinkListResponse{
		TrackingLinks: responses,
		Total:         len(responses),
		Page:          page,
		PageSize:      pageSize,
		TotalPages:    totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// ListTrackingLinksByAffiliate lists tracking links for a specific affiliate
// @Summary List tracking links by affiliate
// @Description Retrieve a list of tracking links for a specific affiliate
// @Tags tracking-links
// @Produce json
// @Param id path int true "Affiliate ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} models.TrackingLinkListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /affiliates/{id}/tracking-links [get]
func (h *TrackingLinkHandler) ListTrackingLinksByAffiliate(c *gin.Context) {
	affiliateID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid affiliate ID",
			Details: "Affiliate ID must be a valid integer",
		})
		return
	}

	page, pageSize := getPaginationParams(c)
	offset := (page - 1) * pageSize

	trackingLinks, err := h.trackingLinkService.ListTrackingLinksByAffiliate(c.Request.Context(), affiliateID, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to list tracking links",
			Details: err.Error(),
		})
		return
	}

	// Convert to response models
	var responses []*models.TrackingLinkResponse
	for _, trackingLink := range trackingLinks {
		responses = append(responses, models.FromTrackingLinkDomain(trackingLink))
	}

	// Calculate total pages (simplified - in production you'd get total count from service)
	totalPages := (len(responses) + pageSize - 1) / pageSize

	response := models.TrackingLinkListResponse{
		TrackingLinks: responses,
		Total:         len(responses),
		Page:          page,
		PageSize:      pageSize,
		TotalPages:    totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// ListTrackingLinksByOrganization lists tracking links for a specific organization
// @Summary List tracking links by organization
// @Description Retrieve a list of tracking links for a specific organization
// @Tags tracking-links
// @Produce json
// @Param organization_id path int true "Organization ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} models.TrackingLinkListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{organization_id}/tracking-links [get]
func (h *TrackingLinkHandler) ListTrackingLinksByOrganization(c *gin.Context) {
	organizationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid organization ID",
			Details: "Organization ID must be a valid integer",
		})
		return
	}

	page, pageSize := getPaginationParams(c)
	offset := (page - 1) * pageSize

	trackingLinks, err := h.trackingLinkService.ListTrackingLinksByOrganization(c.Request.Context(), organizationID, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to list tracking links",
			Details: err.Error(),
		})
		return
	}

	// Convert to response models
	var responses []*models.TrackingLinkResponse
	for _, trackingLink := range trackingLinks {
		responses = append(responses, models.FromTrackingLinkDomain(trackingLink))
	}

	// Calculate total pages (simplified - in production you'd get total count from service)
	totalPages := (len(responses) + pageSize - 1) / pageSize

	response := models.TrackingLinkListResponse{
		TrackingLinks: responses,
		Total:         len(responses),
		Page:          page,
		PageSize:      pageSize,
		TotalPages:    totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// RegenerateTrackingLink regenerates an existing tracking link
// @Summary Regenerate a tracking link
// @Description Regenerate an existing tracking link with provider integration
// @Tags tracking-links
// @Produce json
// @Param organization_id path int true "Organization ID"
// @Param tracking_link_id path int true "Tracking Link ID"
// @Success 200 {object} models.TrackingLinkGenerationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{organization_id}/tracking-links/{tracking_link_id}/regenerate [post]
func (h *TrackingLinkHandler) RegenerateTrackingLink(c *gin.Context) {
	trackingLinkID, err := strconv.ParseInt(c.Param("link_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid tracking link ID",
			Details: "Tracking link ID must be a valid integer",
		})
		return
	}

	response, err := h.trackingLinkService.RegenerateTrackingLink(c.Request.Context(), trackingLinkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to regenerate tracking link",
			Details: err.Error(),
		})
		return
	}

	// Convert to API response
	baseURL := getBaseURL(c)
	apiResponse := models.FromTrackingLinkGenerationDomain(response, baseURL)
	c.JSON(http.StatusOK, apiResponse)
}

// GetTrackingLinkQR generates and returns a QR code for a tracking link
// @Summary Get tracking link QR code
// @Description Generate and return a QR code image for a tracking link
// @Tags tracking-links
// @Produce image/png
// @Param organization_id path int true "Organization ID"
// @Param tracking_link_id path int true "Tracking Link ID"
// @Success 200 {string} binary "QR code image" encoded with base64
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{organization_id}/tracking-links/{tracking_link_id}/qr [get]
func (h *TrackingLinkHandler) GetTrackingLinkQR(c *gin.Context) {
	trackingLinkID, err := strconv.ParseInt(c.Param("link_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid tracking link ID",
			Details: "Tracking link ID must be a valid integer",
		})
		return
	}

	// Get tracking link to verify it exists
	trackingLink, err := h.trackingLinkService.GetTrackingLinkByID(c.Request.Context(), trackingLinkID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Tracking link not found",
			Details: err.Error(),
		})
		return
	}

	// For now, return a simple mock QR code
	// In a real implementation, you would generate a QR code from the tracking URL
	qrData := []byte(trackingLink.Name)

	// Convert qrData from bytes slice to base64 encoded string
	// Convert qrData from bytes slice to base64 encoded string
	base64QR := base64.StdEncoding.EncodeToString(qrData)

	// Return base64 encoded string directly
	c.String(http.StatusOK, base64QR)
}

// Helper functions

// getBaseURL extracts the base URL from the request
func getBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}
