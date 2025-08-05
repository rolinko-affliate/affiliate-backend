package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AgencyDelegationHandler handles HTTP requests for agency delegations
type AgencyDelegationHandler struct {
	delegationService service.AgencyDelegationService
}

// NewAgencyDelegationHandler creates a new agency delegation handler
func NewAgencyDelegationHandler(delegationService service.AgencyDelegationService) *AgencyDelegationHandler {
	return &AgencyDelegationHandler{
		delegationService: delegationService,
	}
}

// CreateDelegation creates a new agency delegation
// @Summary Create agency delegation
// @Description Create a new delegation from advertiser organization to agency organization
// @Tags agency-delegations
// @Accept json
// @Produce json
// @Param request body domain.CreateDelegationRequest true "Create delegation request"
// @Success 201 {object} domain.AgencyDelegation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /agency-delegations [post]
func (h *AgencyDelegationHandler) CreateDelegation(c *gin.Context) {
	var req domain.CreateDelegationRequest
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

	delegation, err := h.delegationService.CreateDelegation(c.Request.Context(), &req, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to create delegation",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, delegation)
}

// AcceptDelegation accepts a pending delegation
// @Summary Accept delegation
// @Description Accept a pending agency delegation
// @Tags agency-delegations
// @Accept json
// @Produce json
// @Param id path int true "Delegation ID"
// @Success 200 {object} domain.AgencyDelegation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /agency-delegations/{id}/accept [post]
func (h *AgencyDelegationHandler) AcceptDelegation(c *gin.Context) {
	delegationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid delegation ID",
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

	delegation, err := h.delegationService.AcceptDelegation(c.Request.Context(), delegationID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to accept delegation",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, delegation)
}

// RejectDelegation rejects a pending delegation
// @Summary Reject delegation
// @Description Reject a pending agency delegation
// @Tags agency-delegations
// @Accept json
// @Produce json
// @Param id path int true "Delegation ID"
// @Success 200 {object} domain.AgencyDelegation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /agency-delegations/{id}/reject [post]
func (h *AgencyDelegationHandler) RejectDelegation(c *gin.Context) {
	delegationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid delegation ID",
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

	delegation, err := h.delegationService.RejectDelegation(c.Request.Context(), delegationID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to reject delegation",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, delegation)
}

// SuspendDelegation suspends an active delegation
// @Summary Suspend delegation
// @Description Suspend an active agency delegation
// @Tags agency-delegations
// @Accept json
// @Produce json
// @Param id path int true "Delegation ID"
// @Success 200 {object} domain.AgencyDelegation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /agency-delegations/{id}/suspend [post]
func (h *AgencyDelegationHandler) SuspendDelegation(c *gin.Context) {
	delegationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid delegation ID",
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

	delegation, err := h.delegationService.SuspendDelegation(c.Request.Context(), delegationID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to suspend delegation",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, delegation)
}

// ReactivateDelegation reactivates a suspended delegation
// @Summary Reactivate delegation
// @Description Reactivate a suspended agency delegation
// @Tags agency-delegations
// @Accept json
// @Produce json
// @Param id path int true "Delegation ID"
// @Success 200 {object} domain.AgencyDelegation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /agency-delegations/{id}/reactivate [post]
func (h *AgencyDelegationHandler) ReactivateDelegation(c *gin.Context) {
	delegationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid delegation ID",
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

	delegation, err := h.delegationService.ReactivateDelegation(c.Request.Context(), delegationID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to reactivate delegation",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, delegation)
}

// RevokeDelegation revokes a delegation
// @Summary Revoke delegation
// @Description Revoke an agency delegation
// @Tags agency-delegations
// @Accept json
// @Produce json
// @Param id path int true "Delegation ID"
// @Success 200 {object} domain.AgencyDelegation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /agency-delegations/{id}/revoke [post]
func (h *AgencyDelegationHandler) RevokeDelegation(c *gin.Context) {
	delegationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid delegation ID",
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

	delegation, err := h.delegationService.RevokeDelegation(c.Request.Context(), delegationID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to revoke delegation",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, delegation)
}

// UpdatePermissions updates the permissions of a delegation
// @Summary Update delegation permissions
// @Description Update the permissions of an agency delegation
// @Tags agency-delegations
// @Accept json
// @Produce json
// @Param id path int true "Delegation ID"
// @Param request body UpdatePermissionsRequest true "Update permissions request"
// @Success 200 {object} domain.AgencyDelegation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /agency-delegations/{id}/permissions [put]
func (h *AgencyDelegationHandler) UpdatePermissions(c *gin.Context) {
	delegationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid delegation ID",
		})
		return
	}

	var req UpdatePermissionsRequest
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

	delegation, err := h.delegationService.UpdatePermissions(c.Request.Context(), delegationID, req.Permissions, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to update permissions",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, delegation)
}

// UpdateExpiration updates the expiration date of a delegation
// @Summary Update delegation expiration
// @Description Update the expiration date of an agency delegation
// @Tags agency-delegations
// @Accept json
// @Produce json
// @Param id path int true "Delegation ID"
// @Param request body UpdateExpirationRequest true "Update expiration request"
// @Success 200 {object} domain.AgencyDelegation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /agency-delegations/{id}/expiration [put]
func (h *AgencyDelegationHandler) UpdateExpiration(c *gin.Context) {
	delegationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid delegation ID",
		})
		return
	}

	var req UpdateExpirationRequest
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

	delegation, err := h.delegationService.UpdateExpiration(c.Request.Context(), delegationID, req.ExpiresAt, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to update expiration",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, delegation)
}

// GetDelegation retrieves a delegation by ID
// @Summary Get delegation
// @Description Get an agency delegation by ID
// @Tags agency-delegations
// @Accept json
// @Produce json
// @Param id path int true "Delegation ID"
// @Param with_details query bool false "Include organization and user details" default(false)
// @Success 200 {object} domain.AgencyDelegation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /agency-delegations/{id} [get]
func (h *AgencyDelegationHandler) GetDelegation(c *gin.Context) {
	delegationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid delegation ID",
		})
		return
	}

	// Check if details are requested
	withDetails := c.Query("with_details") == "true"

	if withDetails {
		delegation, err := h.delegationService.GetDelegationByIDWithDetails(c.Request.Context(), delegationID)
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Delegation not found",
				Details: err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, delegation)
	} else {
		delegation, err := h.delegationService.GetDelegationByID(c.Request.Context(), delegationID)
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Delegation not found",
				Details: err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, delegation)
	}
}

// ListDelegations lists agency delegations with optional filtering
// @Summary List delegations
// @Description List agency delegations with optional filtering
// @Tags agency-delegations
// @Accept json
// @Produce json
// @Param agency_org_id query int false "Agency organization ID"
// @Param advertiser_org_id query int false "Advertiser organization ID"
// @Param status query string false "Delegation status" Enums(pending,active,suspended,revoked)
// @Param include_expired query bool false "Include expired delegations" default(false)
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Param with_details query bool false "Include organization and user details" default(false)
// @Success 200 {array} domain.AgencyDelegation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /agency-delegations [get]
func (h *AgencyDelegationHandler) ListDelegations(c *gin.Context) {
	filter := &domain.DelegationListFilter{}

	// Parse query parameters
	if agencyOrgIDStr := c.Query("agency_org_id"); agencyOrgIDStr != "" {
		if agencyOrgID, err := strconv.ParseInt(agencyOrgIDStr, 10, 64); err == nil {
			filter.AgencyOrgID = &agencyOrgID
		}
	}

	if advertiserOrgIDStr := c.Query("advertiser_org_id"); advertiserOrgIDStr != "" {
		if advertiserOrgID, err := strconv.ParseInt(advertiserOrgIDStr, 10, 64); err == nil {
			filter.AdvertiserOrgID = &advertiserOrgID
		}
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := domain.DelegationStatus(statusStr)
		if status.IsValid() {
			filter.Status = &status
		}
	}

	filter.IncludeExpired = c.Query("include_expired") == "true"

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
		delegations, err := h.delegationService.ListDelegationsWithDetails(c.Request.Context(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to list delegations with details",
				Details: err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, delegations)
	} else {
		delegations, err := h.delegationService.ListDelegations(c.Request.Context(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to list delegations",
				Details: err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, delegations)
	}
}

// CheckPermissions checks if an agency has specific permissions for an advertiser
// @Summary Check permissions
// @Description Check if an agency has specific permissions for an advertiser organization
// @Tags agency-delegations
// @Accept json
// @Produce json
// @Param request body domain.PermissionCheckRequest true "Permission check request"
// @Success 200 {object} domain.PermissionCheckResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /agency-delegations/check-permissions [post]
func (h *AgencyDelegationHandler) CheckPermissions(c *gin.Context) {
	var req domain.PermissionCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	response, err := h.delegationService.CheckPermissions(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to check permissions",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAgencyDelegations retrieves all active delegations for an agency
// @Summary Get agency delegations
// @Description Get all active delegations for an agency organization
// @Tags agency-delegations
// @Accept json
// @Produce json
// @Param agency_org_id path int true "Agency organization ID"
// @Success 200 {array} domain.AgencyDelegation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /agency-delegations/agency/{agency_org_id} [get]
func (h *AgencyDelegationHandler) GetAgencyDelegations(c *gin.Context) {
	agencyOrgID, err := strconv.ParseInt(c.Param("agency_org_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid agency organization ID",
		})
		return
	}

	delegations, err := h.delegationService.GetAgencyDelegations(c.Request.Context(), agencyOrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get agency delegations",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, delegations)
}

// GetAdvertiserDelegations retrieves all active delegations for an advertiser
// @Summary Get advertiser delegations
// @Description Get all active delegations for an advertiser organization
// @Tags agency-delegations
// @Accept json
// @Produce json
// @Param advertiser_org_id path int true "Advertiser organization ID"
// @Success 200 {array} domain.AgencyDelegation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /agency-delegations/advertiser/{advertiser_org_id} [get]
func (h *AgencyDelegationHandler) GetAdvertiserDelegations(c *gin.Context) {
	advertiserOrgID, err := strconv.ParseInt(c.Param("advertiser_org_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid advertiser organization ID",
		})
		return
	}

	delegations, err := h.delegationService.GetAdvertiserDelegations(c.Request.Context(), advertiserOrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get advertiser delegations",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, delegations)
}

// GetAvailablePermissions returns all available delegation permissions
// @Summary Get available permissions
// @Description Get all available delegation permissions
// @Tags agency-delegations
// @Accept json
// @Produce json
// @Success 200 {array} string
// @Security BearerAuth
// @Router /agency-delegations/permissions [get]
func (h *AgencyDelegationHandler) GetAvailablePermissions(c *gin.Context) {
	permissions := domain.GetAllDelegationPermissions()
	
	// Convert to strings for JSON response
	permissionStrings := make([]string, len(permissions))
	for i, permission := range permissions {
		permissionStrings[i] = permission.String()
	}

	c.JSON(http.StatusOK, permissionStrings)
}

// Request/Response types for specific endpoints

// UpdatePermissionsRequest represents a request to update delegation permissions
type UpdatePermissionsRequest struct {
	Permissions []domain.DelegationPermission `json:"permissions" binding:"required"`
}

// UpdateExpirationRequest represents a request to update delegation expiration
type UpdateExpirationRequest struct {
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}