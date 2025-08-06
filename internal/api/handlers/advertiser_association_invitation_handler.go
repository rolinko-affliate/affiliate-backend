package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AdvertiserAssociationInvitationHandler handles HTTP requests for advertiser association invitations
type AdvertiserAssociationInvitationHandler struct {
	invitationService service.AdvertiserAssociationInvitationService
}

// NewAdvertiserAssociationInvitationHandler creates a new invitation handler
func NewAdvertiserAssociationInvitationHandler(invitationService service.AdvertiserAssociationInvitationService) *AdvertiserAssociationInvitationHandler {
	return &AdvertiserAssociationInvitationHandler{
		invitationService: invitationService,
	}
}

// CreateInvitation creates a new invitation
// @Summary Create invitation
// @Description Create a new advertiser association invitation link
// @Tags advertiser-association-invitations
// @Accept json
// @Produce json
// @Param request body domain.CreateInvitationRequest true "Create invitation request"
// @Success 201 {object} domain.AdvertiserAssociationInvitation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /advertiser-association-invitations [post]
func (h *AdvertiserAssociationInvitationHandler) CreateInvitation(c *gin.Context) {
	var req domain.CreateInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	invitation, err := h.invitationService.CreateInvitation(c.Request.Context(), &req, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to create invitation",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, invitation)
}

// GetInvitation retrieves an invitation by ID
// @Summary Get invitation
// @Description Get an advertiser association invitation by ID
// @Tags advertiser-association-invitations
// @Accept json
// @Produce json
// @Param id path int true "Invitation ID"
// @Success 200 {object} domain.AdvertiserAssociationInvitationWithDetails
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /advertiser-association-invitations/{id} [get]
func (h *AdvertiserAssociationInvitationHandler) GetInvitation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid invitation ID",
			Details: err.Error(),
		})
		return
	}

	invitation, err := h.invitationService.GetInvitationByIDWithDetails(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Invitation not found",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, invitation)
}

// GetInvitationByToken retrieves an invitation by token (public endpoint)
// @Summary Get invitation by token
// @Description Get an advertiser association invitation by token (public endpoint for invitation links)
// @Tags advertiser-association-invitations
// @Accept json
// @Produce json
// @Param token path string true "Invitation Token"
// @Success 200 {object} domain.AdvertiserAssociationInvitationWithDetails
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /public/invitations/{token} [get]
func (h *AdvertiserAssociationInvitationHandler) GetInvitationByToken(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invitation token is required",
		})
		return
	}

	invitation, err := h.invitationService.GetInvitationByTokenWithDetails(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Invitation not found",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, invitation)
}

// UpdateInvitation updates an existing invitation
// @Summary Update invitation
// @Description Update an advertiser association invitation
// @Tags advertiser-association-invitations
// @Accept json
// @Produce json
// @Param id path int true "Invitation ID"
// @Param request body domain.UpdateInvitationRequest true "Update invitation request"
// @Success 200 {object} domain.AdvertiserAssociationInvitation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /advertiser-association-invitations/{id} [put]
func (h *AdvertiserAssociationInvitationHandler) UpdateInvitation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid invitation ID",
			Details: err.Error(),
		})
		return
	}

	var req domain.UpdateInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	invitation, err := h.invitationService.UpdateInvitation(c.Request.Context(), id, &req, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to update invitation",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, invitation)
}

// DeleteInvitation deletes an invitation
// @Summary Delete invitation
// @Description Delete an advertiser association invitation
// @Tags advertiser-association-invitations
// @Accept json
// @Produce json
// @Param id path int true "Invitation ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /advertiser-association-invitations/{id} [delete]
func (h *AdvertiserAssociationInvitationHandler) DeleteInvitation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid invitation ID",
			Details: err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	err = h.invitationService.DeleteInvitation(c.Request.Context(), id, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to delete invitation",
			Details: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListInvitations lists invitations with optional filters
// @Summary List invitations
// @Description List advertiser association invitations with optional filters
// @Tags advertiser-association-invitations
// @Accept json
// @Produce json
// @Param advertiser_org_id query int false "Advertiser Organization ID"
// @Param status query string false "Invitation Status" Enums(active, disabled, expired)
// @Param created_by_user_id query string false "Created By User ID"
// @Param include_expired query bool false "Include Expired Invitations"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} domain.AdvertiserAssociationInvitationWithDetails
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /advertiser-association-invitations [get]
func (h *AdvertiserAssociationInvitationHandler) ListInvitations(c *gin.Context) {
	filter := &domain.InvitationListFilter{}

	// Parse query parameters
	if advertiserOrgIDStr := c.Query("advertiser_org_id"); advertiserOrgIDStr != "" {
		if advertiserOrgID, err := strconv.ParseInt(advertiserOrgIDStr, 10, 64); err == nil {
			filter.AdvertiserOrgID = &advertiserOrgID
		}
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := domain.InvitationStatus(statusStr)
		if status.IsValid() {
			filter.Status = &status
		}
	}

	if createdByUserID := c.Query("created_by_user_id"); createdByUserID != "" {
		filter.CreatedByUserID = &createdByUserID
	}

	if includeExpiredStr := c.Query("include_expired"); includeExpiredStr == "true" {
		filter.IncludeExpired = true
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	invitations, err := h.invitationService.ListInvitationsWithDetails(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to list invitations",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, invitations)
}

// UseInvitation uses an invitation to create an organization association
// @Summary Use invitation
// @Description Use an invitation link to create an organization association request
// @Tags advertiser-association-invitations
// @Accept json
// @Produce json
// @Param request body domain.UseInvitationRequest true "Use invitation request"
// @Success 200 {object} domain.UseInvitationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /advertiser-association-invitations/use [post]
func (h *AdvertiserAssociationInvitationHandler) UseInvitation(c *gin.Context) {
	var req domain.UseInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	// Get client IP and user agent for logging
	if req.IPAddress == nil {
		clientIP := c.ClientIP()
		req.IPAddress = &clientIP
	}
	if req.UserAgent == nil {
		userAgent := c.GetHeader("User-Agent")
		req.UserAgent = &userAgent
	}

	response, err := h.invitationService.UseInvitation(c.Request.Context(), &req, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to use invitation",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetInvitationUsageHistory retrieves usage history for an invitation
// @Summary Get invitation usage history
// @Description Get usage history for an advertiser association invitation
// @Tags advertiser-association-invitations
// @Accept json
// @Produce json
// @Param id path int true "Invitation ID"
// @Param limit query int false "Limit (default 50)"
// @Success 200 {array} domain.InvitationUsageLog
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /advertiser-association-invitations/{id}/usage-history [get]
func (h *AdvertiserAssociationInvitationHandler) GetInvitationUsageHistory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid invitation ID",
			Details: err.Error(),
		})
		return
	}

	limit := 50 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	usageHistory, err := h.invitationService.GetInvitationUsageHistory(c.Request.Context(), id, limit)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Failed to get usage history",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, usageHistory)
}

// GenerateInvitationLink generates a full invitation link
// @Summary Generate invitation link
// @Description Generate a full invitation link for sharing
// @Tags advertiser-association-invitations
// @Accept json
// @Produce json
// @Param id path int true "Invitation ID"
// @Param base_url query string false "Base URL (defaults to request host)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /advertiser-association-invitations/{id}/link [get]
func (h *AdvertiserAssociationInvitationHandler) GenerateInvitationLink(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid invitation ID",
			Details: err.Error(),
		})
		return
	}

	// Get base URL from query parameter or use request host
	baseURL := c.Query("base_url")
	if baseURL == "" {
		scheme := "https"
		if c.Request.TLS == nil {
			scheme = "http"
		}
		baseURL = fmt.Sprintf("%s://%s", scheme, c.Request.Host)
	}

	link, err := h.invitationService.GenerateInvitationLink(c.Request.Context(), id, baseURL)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Failed to generate invitation link",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"invitation_link": link,
		"invitation_id":   id,
	})
}