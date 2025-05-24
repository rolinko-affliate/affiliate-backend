package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/affiliate-backend/internal/api/middleware"
	"github.com/affiliate-backend/internal/api/models"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AdvertiserHandler struct {
	advertiserService service.AdvertiserService
	profileService    service.ProfileService
}

func NewAdvertiserHandler(as service.AdvertiserService, ps service.ProfileService) *AdvertiserHandler {
	return &AdvertiserHandler{
		advertiserService: as,
		profileService:    ps,
	}
}

func (h *AdvertiserHandler) CreateAdvertiser(c *gin.Context) {
	var req models.CreateAdvertiserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userOrgID, exists := c.Get("organizationID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return
	}

	userRole, exists := c.Get(middleware.UserRoleKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in context"})
		return
	}

	if userRole.(string) != "Admin" && userOrgID.(int64) != req.OrganizationID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot create advertiser for different organization"})
		return
	}

	advertiser := req.ToDomain()
	createdAdvertiser, err := h.advertiserService.CreateAdvertiser(c.Request.Context(), advertiser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.ToAdvertiserResponse(createdAdvertiser))
}

func (h *AdvertiserHandler) GetAdvertiser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	advertiser, err := h.advertiserService.GetAdvertiserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
		return
	}

	hasAccess, err := h.checkAdvertiserAccess(c, advertiser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, models.ToAdvertiserResponse(advertiser))
}

func (h *AdvertiserHandler) GetAdvertiserWithEverflowData(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	advertiser, err := h.advertiserService.GetAdvertiserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
		return
	}

	hasAccess, err := h.checkAdvertiserAccess(c, advertiser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	data, err := h.advertiserService.GetAdvertiserWithProviderData(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.ToAdvertiserWithEverflowResponse(data))
}

func (h *AdvertiserHandler) UpdateAdvertiser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	var req models.UpdateAdvertiserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingAdvertiser, err := h.advertiserService.GetAdvertiserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
		return
	}

	hasAccess, err := h.checkAdvertiserAccess(c, existingAdvertiser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	advertiser := req.ToDomain(id, existingAdvertiser.OrganizationID)
	if err := h.advertiserService.UpdateAdvertiser(c.Request.Context(), advertiser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedAdvertiser, err := h.advertiserService.GetAdvertiserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated advertiser"})
		return
	}

	c.JSON(http.StatusOK, models.ToAdvertiserResponse(updatedAdvertiser))
}

func (h *AdvertiserHandler) ListAdvertisers(c *gin.Context) {
	userOrgID, exists := c.Get("organizationID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})
		return
	}

	userRole, exists := c.Get(middleware.UserRoleKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in context"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var orgID int64
	if userRole.(string) == "Admin" {
		if orgIDParam := c.Query("organization_id"); orgIDParam != "" {
			if parsedOrgID, err := strconv.ParseInt(orgIDParam, 10, 64); err == nil {
				orgID = parsedOrgID
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization_id parameter"})
				return
			}
		} else {
			orgID = userOrgID.(int64)
		}
	} else {
		orgID = userOrgID.(int64)
	}

	advertisers, err := h.advertiserService.ListAdvertisersByOrganization(c.Request.Context(), orgID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*models.AdvertiserResponse
	for _, advertiser := range advertisers {
		responses = append(responses, models.ToAdvertiserResponse(advertiser))
	}

	c.JSON(http.StatusOK, models.ListAdvertisersResponse{
		Advertisers: responses,
		Page:        page,
		PageSize:    pageSize,
		Total:       len(responses),
	})
}

func (h *AdvertiserHandler) DeleteAdvertiser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	advertiser, err := h.advertiserService.GetAdvertiserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
		return
	}

	hasAccess, err := h.checkAdvertiserAccess(c, advertiser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.advertiserService.DeleteAdvertiser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *AdvertiserHandler) SyncToEverflow(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	advertiser, err := h.advertiserService.GetAdvertiserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
		return
	}

	hasAccess, err := h.checkAdvertiserAccess(c, advertiser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.advertiserService.SyncAdvertiserToProvider(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sync to Everflow initiated successfully"})
}

func (h *AdvertiserHandler) SyncFromEverflow(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	advertiser, err := h.advertiserService.GetAdvertiserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
		return
	}

	hasAccess, err := h.checkAdvertiserAccess(c, advertiser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.advertiserService.SyncAdvertiserFromProvider(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sync from Everflow completed successfully"})
}

func (h *AdvertiserHandler) CompareWithEverflow(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	advertiser, err := h.advertiserService.GetAdvertiserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
		return
	}

	hasAccess, err := h.checkAdvertiserAccess(c, advertiser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	discrepancies, err := h.advertiserService.CompareAdvertiserWithProvider(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"discrepancies": discrepancies})
}

func (h *AdvertiserHandler) checkAdvertiserAccess(c *gin.Context, advertiserOrgID int64) (bool, error) {
	userRole, exists := c.Get(middleware.UserRoleKey)
	if !exists {
		return false, fmt.Errorf("user role not found in context")
	}

	if userRole.(string) == "Admin" {
		return true, nil
	}

	userOrgID, exists := c.Get("organizationID")
	if !exists {
		return false, fmt.Errorf("user organization ID not found in context")
	}

	return userOrgID.(int64) == advertiserOrgID, nil
}

// ListAdvertisersByOrganization retrieves a list of advertisers for an organization
func (h *AdvertiserHandler) ListAdvertisersByOrganization(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

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

// CreateAdvertiserProviderMapping creates a new advertiser provider mapping
func (h *AdvertiserHandler) CreateAdvertiserProviderMapping(c *gin.Context) {
	var req models.CreateAdvertiserProviderMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	mapping, err := h.advertiserService.CreateAdvertiserProviderMapping(c.Request.Context(), &req.ProviderMapping)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create provider mapping: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.CreateAdvertiserProviderMappingResponse{
		ProviderMapping: *mapping,
	})
}

// GetAdvertiserProviderMapping retrieves an advertiser provider mapping
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get provider mapping: " + err.Error()})
		return
	}

	if mapping == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider mapping not found"})
		return
	}

	c.JSON(http.StatusOK, models.GetAdvertiserProviderMappingResponse{
		ProviderMapping: *mapping,
	})
}

// UpdateAdvertiserProviderMapping updates an advertiser provider mapping
func (h *AdvertiserHandler) UpdateAdvertiserProviderMapping(c *gin.Context) {
	mappingIDStr := c.Param("mappingId")
	mappingID, err := strconv.ParseInt(mappingIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mapping ID"})
		return
	}

	var req models.UpdateAdvertiserProviderMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	req.ProviderMapping.MappingID = mappingID
	err = h.advertiserService.UpdateAdvertiserProviderMapping(c.Request.Context(), &req.ProviderMapping)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update provider mapping: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Provider mapping updated successfully"})
}

// DeleteAdvertiserProviderMapping deletes an advertiser provider mapping
func (h *AdvertiserHandler) DeleteAdvertiserProviderMapping(c *gin.Context) {
	mappingIDStr := c.Param("mappingId")
	mappingID, err := strconv.ParseInt(mappingIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mapping ID"})
		return
	}

	err = h.advertiserService.DeleteAdvertiserProviderMapping(c.Request.Context(), mappingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete provider mapping: " + err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// SyncAdvertiserToEverflow syncs an advertiser to Everflow
func (h *AdvertiserHandler) SyncAdvertiserToEverflow(c *gin.Context) {
	advertiserIDStr := c.Param("id")
	advertiserID, err := strconv.ParseInt(advertiserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	advertiser, err := h.advertiserService.GetAdvertiserByID(c.Request.Context(), advertiserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get advertiser: " + err.Error()})
		return
	}

	if advertiser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
		return
	}

	hasAccess, err := h.checkAdvertiserAccess(c, advertiser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions: " + err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to sync this advertiser"})
		return
	}

	err = h.advertiserService.SyncAdvertiserToProvider(c.Request.Context(), advertiserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync advertiser to Everflow: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Advertiser synced to Everflow successfully"})
}

// SyncAdvertiserFromEverflow syncs an advertiser from Everflow
func (h *AdvertiserHandler) SyncAdvertiserFromEverflow(c *gin.Context) {
	advertiserIDStr := c.Param("id")
	advertiserID, err := strconv.ParseInt(advertiserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	advertiser, err := h.advertiserService.GetAdvertiserByID(c.Request.Context(), advertiserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get advertiser: " + err.Error()})
		return
	}

	if advertiser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
		return
	}

	hasAccess, err := h.checkAdvertiserAccess(c, advertiser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions: " + err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to sync this advertiser"})
		return
	}

	err = h.advertiserService.SyncAdvertiserFromProvider(c.Request.Context(), advertiserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync advertiser from Everflow: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Advertiser synced from Everflow successfully"})
}

// CompareAdvertiserWithEverflow compares an advertiser with Everflow data
func (h *AdvertiserHandler) CompareAdvertiserWithEverflow(c *gin.Context) {
	advertiserIDStr := c.Param("id")
	advertiserID, err := strconv.ParseInt(advertiserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	advertiser, err := h.advertiserService.GetAdvertiserByID(c.Request.Context(), advertiserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get advertiser: " + err.Error()})
		return
	}

	if advertiser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
		return
	}

	hasAccess, err := h.checkAdvertiserAccess(c, advertiser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions: " + err.Error()})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to compare this advertiser"})
		return
	}

	comparison, err := h.advertiserService.CompareAdvertiserWithProvider(c.Request.Context(), advertiserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to compare advertiser with Everflow: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, comparison)
}