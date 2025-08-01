package handlers

import (
	"net/http"
	"strconv"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// OrganizationAssociationHandler handles HTTP requests for organization associations
type OrganizationAssociationHandler struct {
	associationService service.OrganizationAssociationService
}

// NewOrganizationAssociationHandler creates a new organization association handler
func NewOrganizationAssociationHandler(associationService service.OrganizationAssociationService) *OrganizationAssociationHandler {
	return &OrganizationAssociationHandler{
		associationService: associationService,
	}
}

// CreateInvitation creates a new invitation from advertiser to affiliate
// @Summary Create invitation
// @Description Create a new invitation from advertiser organization to affiliate organization
// @Tags organization-associations
// @Accept json
// @Produce json
// @Param request body domain.CreateAssociationRequest true "Create invitation request"
// @Success 201 {object} domain.OrganizationAssociation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/organization-associations/invitations [post]
func (h *OrganizationAssociationHandler) CreateInvitation(c *gin.Context) {
	var req domain.CreateAssociationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Set association type to invitation
	req.AssociationType = domain.AssociationTypeInvitation

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	association, err := h.associationService.CreateInvitation(c.Request.Context(), &req, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to create invitation",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, association)
}

// CreateRequest creates a new request from affiliate to advertiser
// @Summary Create request
// @Description Create a new request from affiliate organization to advertiser organization
// @Tags organization-associations
// @Accept json
// @Produce json
// @Param request body domain.CreateAssociationRequest true "Create request"
// @Success 201 {object} domain.OrganizationAssociation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/organization-associations/requests [post]
func (h *OrganizationAssociationHandler) CreateRequest(c *gin.Context) {
	var req domain.CreateAssociationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Set association type to request
	req.AssociationType = domain.AssociationTypeRequest

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	association, err := h.associationService.CreateRequest(c.Request.Context(), &req, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to create request",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, association)
}

// ApproveAssociation approves a pending association
// @Summary Approve association
// @Description Approve a pending organization association
// @Tags organization-associations
// @Accept json
// @Produce json
// @Param id path int true "Association ID"
// @Success 200 {object} domain.OrganizationAssociation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/organization-associations/{id}/approve [post]
func (h *OrganizationAssociationHandler) ApproveAssociation(c *gin.Context) {
	associationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid association ID",
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

	association, err := h.associationService.ApproveAssociation(c.Request.Context(), associationID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to approve association",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, association)
}

// RejectAssociation rejects a pending association
// @Summary Reject association
// @Description Reject a pending organization association
// @Tags organization-associations
// @Accept json
// @Produce json
// @Param id path int true "Association ID"
// @Success 200 {object} domain.OrganizationAssociation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/organization-associations/{id}/reject [post]
func (h *OrganizationAssociationHandler) RejectAssociation(c *gin.Context) {
	associationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid association ID",
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

	association, err := h.associationService.RejectAssociation(c.Request.Context(), associationID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to reject association",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, association)
}

// SuspendAssociation suspends an active association
// @Summary Suspend association
// @Description Suspend an active organization association
// @Tags organization-associations
// @Accept json
// @Produce json
// @Param id path int true "Association ID"
// @Success 200 {object} domain.OrganizationAssociation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/organization-associations/{id}/suspend [post]
func (h *OrganizationAssociationHandler) SuspendAssociation(c *gin.Context) {
	associationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid association ID",
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

	association, err := h.associationService.SuspendAssociation(c.Request.Context(), associationID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to suspend association",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, association)
}

// ReactivateAssociation reactivates a suspended association
// @Summary Reactivate association
// @Description Reactivate a suspended organization association
// @Tags organization-associations
// @Accept json
// @Produce json
// @Param id path int true "Association ID"
// @Success 200 {object} domain.OrganizationAssociation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/organization-associations/{id}/reactivate [post]
func (h *OrganizationAssociationHandler) ReactivateAssociation(c *gin.Context) {
	associationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid association ID",
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

	association, err := h.associationService.ReactivateAssociation(c.Request.Context(), associationID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to reactivate association",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, association)
}

// UpdateVisibility updates the visibility settings of an association
// @Summary Update visibility
// @Description Update the visibility settings of an organization association
// @Tags organization-associations
// @Accept json
// @Produce json
// @Param id path int true "Association ID"
// @Param request body domain.UpdateAssociationRequest true "Update visibility request"
// @Success 200 {object} domain.OrganizationAssociation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/organization-associations/{id}/visibility [put]
func (h *OrganizationAssociationHandler) UpdateVisibility(c *gin.Context) {
	associationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid association ID",
		})
		return
	}

	var req domain.UpdateAssociationRequest
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

	association, err := h.associationService.UpdateVisibility(c.Request.Context(), associationID, &req, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to update visibility",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, association)
}

// GetAssociation retrieves an association by ID
// @Summary Get association
// @Description Get an organization association by ID
// @Tags organization-associations
// @Accept json
// @Produce json
// @Param id path int true "Association ID"
// @Success 200 {object} domain.OrganizationAssociation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/organization-associations/{id} [get]
func (h *OrganizationAssociationHandler) GetAssociation(c *gin.Context) {
	associationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid association ID",
		})
		return
	}

	association, err := h.associationService.GetAssociationByID(c.Request.Context(), associationID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Association not found",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, association)
}

// ListAssociations lists organization associations with optional filtering
// @Summary List associations
// @Description List organization associations with optional filtering
// @Tags organization-associations
// @Accept json
// @Produce json
// @Param advertiser_org_id query int false "Advertiser organization ID"
// @Param affiliate_org_id query int false "Affiliate organization ID"
// @Param status query string false "Association status" Enums(pending,active,suspended,rejected)
// @Param association_type query string false "Association type" Enums(invitation,request)
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Param with_details query bool false "Include organization and user details" default(false)
// @Success 200 {array} domain.OrganizationAssociation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/organization-associations [get]
func (h *OrganizationAssociationHandler) ListAssociations(c *gin.Context) {
	filter := &domain.AssociationListFilter{}

	// Parse query parameters
	if advertiserOrgIDStr := c.Query("advertiser_org_id"); advertiserOrgIDStr != "" {
		if advertiserOrgID, err := strconv.ParseInt(advertiserOrgIDStr, 10, 64); err == nil {
			filter.AdvertiserOrgID = &advertiserOrgID
		}
	}

	if affiliateOrgIDStr := c.Query("affiliate_org_id"); affiliateOrgIDStr != "" {
		if affiliateOrgID, err := strconv.ParseInt(affiliateOrgIDStr, 10, 64); err == nil {
			filter.AffiliateOrgID = &affiliateOrgID
		}
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := domain.AssociationStatus(statusStr)
		if status.IsValid() {
			filter.Status = &status
		}
	}

	if associationTypeStr := c.Query("association_type"); associationTypeStr != "" {
		associationType := domain.AssociationType(associationTypeStr)
		if associationType.IsValid() {
			filter.AssociationType = &associationType
		}
	}

	// Parse pagination parameters
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}
	if filter.Limit == 0 {
		filter.Limit = 50 // Default limit
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	// Check if details are requested
	withDetails := c.Query("with_details") == "true"

	if withDetails {
		associations, err := h.associationService.ListAssociationsWithDetails(c.Request.Context(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to list associations with details",
				Details: err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, associations)
	} else {
		associations, err := h.associationService.ListAssociations(c.Request.Context(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to list associations",
				Details: err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, associations)
	}
}

// GetAssociationsForOrganization gets associations for a specific organization
// @Summary Get associations for organization
// @Description Get all associations for a specific organization (advertiser or affiliate)
// @Tags organization-associations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param status query string false "Association status" Enums(pending,active,suspended,rejected)
// @Param with_details query bool false "Include organization and user details" default(false)
// @Success 200 {array} domain.OrganizationAssociation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/organizations/{id}/associations [get]
func (h *OrganizationAssociationHandler) GetAssociationsForOrganization(c *gin.Context) {
	orgID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid organization ID",
		})
		return
	}

	var status *domain.AssociationStatus
	if statusStr := c.Query("status"); statusStr != "" {
		s := domain.AssociationStatus(statusStr)
		if s.IsValid() {
			status = &s
		}
	}

	withDetails := c.Query("with_details") == "true"

	// Create filter for both advertiser and affiliate organizations
	filter := &domain.AssociationListFilter{
		Status: status,
		Limit:  100, // Default limit
	}

	// We need to check both advertiser and affiliate associations
	// First try as advertiser
	filter.AdvertiserOrgID = &orgID
	var associations interface{}
	var listErr error

	if withDetails {
		advAssociations, err := h.associationService.ListAssociationsWithDetails(c.Request.Context(), filter)
		if err != nil {
			listErr = err
		} else {
			// Try as affiliate
			filter.AdvertiserOrgID = nil
			filter.AffiliateOrgID = &orgID
			affAssociations, err := h.associationService.ListAssociationsWithDetails(c.Request.Context(), filter)
			if err != nil {
				listErr = err
			} else {
				// Combine results
				allAssociations := append(advAssociations, affAssociations...)
				associations = allAssociations
			}
		}
	} else {
		advAssociations, err := h.associationService.ListAssociations(c.Request.Context(), filter)
		if err != nil {
			listErr = err
		} else {
			// Try as affiliate
			filter.AdvertiserOrgID = nil
			filter.AffiliateOrgID = &orgID
			affAssociations, err := h.associationService.ListAssociations(c.Request.Context(), filter)
			if err != nil {
				listErr = err
			} else {
				// Combine results
				allAssociations := append(advAssociations, affAssociations...)
				associations = allAssociations
			}
		}
	}

	if listErr != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get associations for organization",
			Details: listErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, associations)
}

// GetVisibleAffiliatesForAdvertiser gets all affiliates visible to an advertiser organization
// @Summary Get visible affiliates for advertiser
// @Description Get all affiliates from affiliate organizations that are visible to the specified advertiser organization based on active associations
// @Tags organization-associations
// @Accept json
// @Produce json
// @Param advertiser_org_id path int true "Advertiser Organization ID"
// @Param affiliate_org_id query int false "Filter by specific affiliate organization ID"
// @Success 200 {array} domain.Affiliate
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/organizations/{advertiser_org_id}/visible-affiliates [get]
func (h *OrganizationAssociationHandler) GetVisibleAffiliatesForAdvertiser(c *gin.Context) {
	advertiserOrgID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid advertiser organization ID",
		})
		return
	}

	var affiliateOrgID *int64
	if affiliateOrgIDStr := c.Query("affiliate_org_id"); affiliateOrgIDStr != "" {
		id, err := strconv.ParseInt(affiliateOrgIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid affiliate organization ID",
			})
			return
		}
		affiliateOrgID = &id
	}

	affiliates, err := h.associationService.GetVisibleAffiliatesForAdvertiser(c.Request.Context(), advertiserOrgID, affiliateOrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get visible affiliates",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, affiliates)
}

// GetVisibleCampaignsForAffiliate gets all campaigns visible to an affiliate organization
// @Summary Get visible campaigns for affiliate
// @Description Get all campaigns from advertiser organizations that are visible to the specified affiliate organization based on active associations
// @Tags organization-associations
// @Accept json
// @Produce json
// @Param affiliate_org_id path int true "Affiliate Organization ID"
// @Param advertiser_org_id query int false "Filter by specific advertiser organization ID"
// @Success 200 {array} domain.Campaign
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/organizations/{affiliate_org_id}/visible-campaigns [get]
func (h *OrganizationAssociationHandler) GetVisibleCampaignsForAffiliate(c *gin.Context) {
	affiliateOrgID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid affiliate organization ID",
		})
		return
	}

	var advertiserOrgID *int64
	if advertiserOrgIDStr := c.Query("advertiser_org_id"); advertiserOrgIDStr != "" {
		id, err := strconv.ParseInt(advertiserOrgIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid advertiser organization ID",
			})
			return
		}
		advertiserOrgID = &id
	}

	campaigns, err := h.associationService.GetVisibleCampaignsForAffiliate(c.Request.Context(), affiliateOrgID, advertiserOrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get visible campaigns",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}