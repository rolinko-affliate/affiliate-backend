package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

// GetTrackingLinkQR generates a QR code for a tracking link
// @Summary Get QR code for tracking link
// @Description Generate a QR code for the specified tracking link
// @Tags tracking-links
// @Produce text/plain
// @Param organization_id path int true "Organization ID"
// @Param link_id path int true "Tracking Link ID"
// @Success 200 {string} string "Base64 encoded QR code"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /organizations/{organization_id}/tracking-links/{link_id}/qr [get]
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
	base64QR := base64.StdEncoding.EncodeToString(qrData)

	// Return base64 encoded string directly
	c.String(http.StatusOK, base64QR)
}

// CreateTrackingLinkClean creates a new tracking link with uniqueness guarantee
// @Summary Create a new tracking link
// @Description Create a new tracking link with uniqueness guarantee for campaign_id + affiliate_id combination
// @Tags tracking-links
// @Accept json
// @Produce json
// @Param request body models.TrackingLinkGenerationRequest true "Tracking link creation request"
// @Success 201 {object} models.TrackingLinkGenerationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /tracking-links [post]
func (h *TrackingLinkHandler) CreateTrackingLinkClean(c *gin.Context) {
	var req models.TrackingLinkGenerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Convert to upsert request to ensure uniqueness of campaign_id + affiliate_id combination
	upsertReq := &models.TrackingLinkUpsertRequest{
		CampaignID:          req.CampaignID,
		AffiliateID:         req.AffiliateID,
		Name:                req.Name,
		Description:         req.Description,
		SourceID:            req.SourceID,
		Sub1:                req.Sub1,
		Sub2:                req.Sub2,
		Sub3:                req.Sub3,
		Sub4:                req.Sub4,
		Sub5:                req.Sub5,
		IsEncryptParameters: req.IsEncryptParameters,
		IsRedirectLink:      req.IsRedirectLink,
		InternalNotes:       req.InternalNotes,
		Tags:                req.Tags,
	}

	// Convert to domain model
	domainUpsertReq := upsertReq.ToTrackingLinkUpsertDomain()

	response, err := h.trackingLinkService.UpsertTrackingLink(c.Request.Context(), domainUpsertReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create tracking link",
			Details: err.Error(),
		})
		return
	}

	// Convert to response model
	baseURL := fmt.Sprintf("%s://%s", getScheme(c), c.Request.Host)
	apiResponse := models.FromTrackingLinkUpsertDomain(response, baseURL)

	// Return 201 for created, 200 for updated
	statusCode := http.StatusCreated
	if !response.IsNew {
		statusCode = http.StatusOK
	}

	c.JSON(statusCode, apiResponse)
}

// GetTrackingLinkClean retrieves a tracking link by ID
// @Summary Get tracking link by ID
// @Description Retrieve a tracking link by its ID
// @Tags tracking-links
// @Produce json
// @Param id path int true "Tracking Link ID"
// @Success 200 {object} models.TrackingLinkResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /tracking-links/{id} [get]
func (h *TrackingLinkHandler) GetTrackingLinkClean(c *gin.Context) {
	trackingLinkID, err := strconv.ParseInt(c.Param("id"), 10, 64)
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

	// Convert to response model
	response := models.FromTrackingLinkDomain(trackingLink)
	c.JSON(http.StatusOK, response)
}

// UpdateTrackingLinkClean updates a tracking link and regenerates if key parameters change
// @Summary Update tracking link
// @Description Update a tracking link and regenerate if key parameters change
// @Tags tracking-links
// @Accept json
// @Produce json
// @Param id path int true "Tracking Link ID"
// @Param request body models.TrackingLinkUpdateRequest true "Tracking link update request"
// @Success 200 {object} models.TrackingLinkResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /tracking-links/{id} [put]
func (h *TrackingLinkHandler) UpdateTrackingLinkClean(c *gin.Context) {
	trackingLinkID, err := strconv.ParseInt(c.Param("id"), 10, 64)
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

	// Update the existing tracking link with request data
	req.UpdateTrackingLinkDomain(existingTrackingLink)

	// Update tracking link
	err = h.trackingLinkService.UpdateTrackingLink(c.Request.Context(), existingTrackingLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update tracking link",
			Details: err.Error(),
		})
		return
	}

	// Convert to response model
	response := models.FromTrackingLinkDomain(existingTrackingLink)
	c.JSON(http.StatusOK, response)
}

// DeleteTrackingLinkClean deletes a tracking link
// @Summary Delete tracking link
// @Description Delete a tracking link by its ID
// @Tags tracking-links
// @Param id path int true "Tracking Link ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /tracking-links/{id} [delete]
func (h *TrackingLinkHandler) DeleteTrackingLinkClean(c *gin.Context) {
	trackingLinkID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid tracking link ID",
			Details: "Tracking link ID must be a valid integer",
		})
		return
	}

	err = h.trackingLinkService.DeleteTrackingLink(c.Request.Context(), trackingLinkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to delete tracking link",
			Details: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListTrackingLinksClean lists tracking links with filtering
// @Summary List tracking links
// @Description List tracking links with optional filtering by affiliate IDs and campaign IDs
// @Tags tracking-links
// @Produce json
// @Param affiliate_ids query string false "Comma-separated list of affiliate IDs"
// @Param campaign_ids query string false "Comma-separated list of campaign IDs"
// @Param limit query int false "Number of items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} models.TrackingLinkListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /tracking-links [get]
func (h *TrackingLinkHandler) ListTrackingLinksClean(c *gin.Context) {
	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Parse affiliate IDs
	var affiliateIDs []int64
	if affiliateIDsStr := c.Query("affiliate_ids"); affiliateIDsStr != "" {
		var err error
		affiliateIDs, err = parseCommaSeparatedInt64s(affiliateIDsStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid affiliate_ids parameter",
				Details: err.Error(),
			})
			return
		}
	}

	// Parse campaign IDs
	var campaignIDs []int64
	if campaignIDsStr := c.Query("campaign_ids"); campaignIDsStr != "" {
		var err error
		campaignIDs, err = parseCommaSeparatedInt64s(campaignIDsStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid campaign_ids parameter",
				Details: err.Error(),
			})
			return
		}
	}

	// Get tracking links with filters
	trackingLinks, total, err := h.trackingLinkService.ListTrackingLinksWithFilters(
		c.Request.Context(),
		affiliateIDs,
		campaignIDs,
		limit,
		offset,
	)
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

	// Calculate pagination info
	page := (offset / limit) + 1
	totalPages := (total + limit - 1) / limit // Ceiling division

	// Create list response
	listResponse := &models.TrackingLinkListResponse{
		TrackingLinks: responses,
		Total:         total,
		Page:          page,
		PageSize:      limit,
		TotalPages:    totalPages,
	}

	c.JSON(http.StatusOK, listResponse)
}

// parseCommaSeparatedInt64s parses a comma-separated string of integers
func parseCommaSeparatedInt64s(s string) ([]int64, error) {
	if s == "" {
		return nil, nil
	}

	parts := strings.Split(s, ",")
	result := make([]int64, len(parts))

	for i, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}

		val, err := strconv.ParseInt(trimmed, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer: %s", trimmed)
		}
		result[i] = val
	}

	return result, nil
}

// getScheme returns the scheme (http or https) for the request
func getScheme(c *gin.Context) string {
	if c.Request.TLS != nil {
		return "https"
	}
	if scheme := c.GetHeader("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	return "http"
}
