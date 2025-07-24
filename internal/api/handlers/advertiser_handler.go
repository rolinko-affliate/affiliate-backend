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

// CreateAdvertiser creates a new advertiser
// @Summary      Create a new advertiser
// @Description  Creates a new advertiser with the given details
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        request  body      models.CreateAdvertiserRequest  true  "Advertiser details"
// @Success      201      {object}  models.AdvertiserResponse       "Created advertiser"
// @Failure      400      {object}  map[string]string               "Invalid request"
// @Failure      401      {object}  map[string]string               "Unauthorized"
// @Failure      403      {object}  map[string]string               "Forbidden"
// @Failure      500      {object}  map[string]string               "Internal server error"
// @Security     BearerAuth
// @Router       /advertisers [post]
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

// GetAdvertiser retrieves an advertiser by ID
// @Summary      Get advertiser by ID
// @Description  Retrieves an advertiser by its ID
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        id   path      int                         true  "Advertiser ID"
// @Success      200  {object}  models.AdvertiserResponse   "Advertiser details"
// @Failure      400  {object}  map[string]string           "Invalid advertiser ID"
// @Failure      403  {object}  map[string]string           "Forbidden - User doesn't have permission"
// @Failure      404  {object}  map[string]string           "Advertiser not found"
// @Failure      500  {object}  map[string]string           "Internal server error"
// @Security     BearerAuth
// @Router       /advertisers/{id} [get]
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

// UpdateAdvertiser updates an advertiser
// @Summary      Update advertiser
// @Description  Updates an advertiser with the given details
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        id       path      int                           true  "Advertiser ID"
// @Param        request  body      models.UpdateAdvertiserRequest true  "Advertiser details"
// @Success      200      {object}  models.AdvertiserResponse     "Updated advertiser"
// @Failure      400      {object}  map[string]string             "Invalid request"
// @Failure      403      {object}  map[string]string             "Forbidden - User doesn't have permission"
// @Failure      404      {object}  map[string]string             "Advertiser not found"
// @Failure      500      {object}  map[string]string             "Internal server error"
// @Security     BearerAuth
// @Router       /advertisers/{id} [put]
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
// @Summary      List advertisers by organization
// @Description  Retrieves a list of advertisers for an organization with pagination
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        id  path      int                   true   "Organization ID"
// @Param        page           query     int                   false  "Page number (default: 1)"
// @Param        pageSize       query     int                   false  "Page size (default: 10)"
// @Success      200            {array}   domain.Advertiser     "List of advertisers"
// @Failure      400            {object}  map[string]string     "Invalid organization ID"
// @Failure      403            {object}  map[string]string     "Forbidden - User doesn't have permission"
// @Failure      500            {object}  map[string]string     "Internal server error"
// @Security     BearerAuth
// @Router       /organizations/{id}/advertisers [get]
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

// CreateProviderMapping creates a new advertiser provider mapping
// @Summary      Create advertiser provider mapping
// @Description  Creates a new advertiser provider mapping
// @Tags         advertiser-provider-mappings
// @Accept       json
// @Produce      json
// @Param        request  body      models.CreateAdvertiserProviderMappingRequest  true  "Provider mapping details"
// @Success      201      {object}  models.CreateAdvertiserProviderMappingResponse "Created provider mapping"
// @Failure      400      {object}  map[string]string                              "Invalid request"
// @Failure      500      {object}  map[string]string                              "Internal server error"
// @Security     BearerAuth
// @Router       /advertiser-provider-mappings [post]
func (h *AdvertiserHandler) CreateProviderMapping(c *gin.Context) {
	var req models.CreateAdvertiserProviderMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	mapping, err := h.advertiserService.CreateProviderMapping(c.Request.Context(), &req.ProviderMapping)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create provider mapping: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.CreateAdvertiserProviderMappingResponse{
		ProviderMapping: *mapping,
	})
}

// GetProviderMapping retrieves an advertiser provider mapping
// @Summary      Get advertiser provider mapping
// @Description  Retrieves an advertiser provider mapping by advertiser ID and provider type
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        id           path      int                                         true  "Advertiser ID"
// @Param        providerType path      string                                      true  "Provider type"
// @Success      200          {object}  models.GetAdvertiserProviderMappingResponse "Provider mapping details"
// @Failure      400          {object}  map[string]string                           "Invalid request"
// @Failure      404          {object}  map[string]string                           "Provider mapping not found"
// @Failure      500          {object}  map[string]string                           "Internal server error"
// @Security     BearerAuth
// @Router       /advertisers/{id}/provider-mappings/{providerType} [get]
func (h *AdvertiserHandler) GetProviderMapping(c *gin.Context) {
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

	mapping, err := h.advertiserService.GetProviderMapping(c.Request.Context(), advertiserID, providerType)
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

// UpdateProviderMapping updates an advertiser provider mapping
// @Summary      Update advertiser provider mapping
// @Description  Updates an advertiser provider mapping by mapping ID
// @Tags         advertiser-provider-mappings
// @Accept       json
// @Produce      json
// @Param        mappingId  path      int                                           true  "Mapping ID"
// @Param        request    body      models.UpdateAdvertiserProviderMappingRequest true  "Provider mapping details"
// @Success      200        {object}  map[string]string                             "Update successful"
// @Failure      400        {object}  map[string]string                             "Invalid request"
// @Failure      500        {object}  map[string]string                             "Internal server error"
// @Security     BearerAuth
// @Router       /advertiser-provider-mappings/{mappingId} [put]
func (h *AdvertiserHandler) UpdateProviderMapping(c *gin.Context) {
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
	err = h.advertiserService.UpdateProviderMapping(c.Request.Context(), &req.ProviderMapping)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update provider mapping: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Provider mapping updated successfully"})
}

// DeleteProviderMapping deletes an advertiser provider mapping
// @Summary      Delete advertiser provider mapping
// @Description  Deletes an advertiser provider mapping by mapping ID
// @Tags         advertiser-provider-mappings
// @Accept       json
// @Produce      json
// @Param        mappingId  path      int                true  "Mapping ID"
// @Success      204        {object}  nil                "No content"
// @Failure      400        {object}  map[string]string  "Invalid mapping ID"
// @Failure      500        {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /advertiser-provider-mappings/{mappingId} [delete]
func (h *AdvertiserHandler) DeleteProviderMapping(c *gin.Context) {
	mappingIDStr := c.Param("mappingId")
	mappingID, err := strconv.ParseInt(mappingIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mapping ID"})
		return
	}

	err = h.advertiserService.DeleteProviderMapping(c.Request.Context(), mappingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete provider mapping: " + err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// SyncAdvertiserToEverflow syncs an advertiser to Everflow
// @Summary      Sync advertiser to Everflow
// @Description  Syncs an advertiser to the Everflow platform
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        id   path      int                true  "Advertiser ID"
// @Success      200  {object}  map[string]string  "Sync initiated successfully"
// @Failure      400  {object}  map[string]string  "Invalid advertiser ID"
// @Failure      403  {object}  map[string]string  "Forbidden - User doesn't have permission"
// @Failure      404  {object}  map[string]string  "Advertiser not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /advertisers/{id}/sync-to-everflow [post]
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
// @Summary      Sync advertiser from Everflow
// @Description  Syncs an advertiser from the Everflow platform
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        id   path      int                true  "Advertiser ID"
// @Success      200  {object}  map[string]string  "Sync completed successfully"
// @Failure      400  {object}  map[string]string  "Invalid advertiser ID"
// @Failure      403  {object}  map[string]string  "Forbidden - User doesn't have permission"
// @Failure      404  {object}  map[string]string  "Advertiser not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /advertisers/{id}/sync-from-everflow [post]
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
// @Summary      Compare advertiser with Everflow
// @Description  Compares an advertiser with its Everflow counterpart and returns discrepancies
// @Tags         advertisers
// @Accept       json
// @Produce      json
// @Param        id   path      int                true  "Advertiser ID"
// @Success      200  {object}  map[string]interface{}  "Comparison results with discrepancies"
// @Failure      400  {object}  map[string]string       "Invalid advertiser ID"
// @Failure      403  {object}  map[string]string       "Forbidden - User doesn't have permission"
// @Failure      404  {object}  map[string]string       "Advertiser not found"
// @Failure      500  {object}  map[string]string       "Internal server error"
// @Security     BearerAuth
// @Router       /advertisers/{id}/compare-with-everflow [get]
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

// SyncAllAdvertisersToEverflow syncs all advertisers without Everflow mappings to Everflow
// @Summary      Sync all advertisers to Everflow
// @Description  Creates Everflow advertisers for all local advertisers that don't have provider mappings
// @Tags         advertisers
// @Produce      json
// @Success      200      {object}  domain.BulkSyncResult           "Sync results"
// @Failure      401      {object}  map[string]string               "Unauthorized"
// @Failure      403      {object}  map[string]string               "Forbidden"
// @Failure      500      {object}  map[string]string               "Internal server error"
// @Security     BearerAuth
// @Router       /advertisers/sync-all-to-everflow [post]
func (h *AdvertiserHandler) SyncAllAdvertisersToEverflow(c *gin.Context) {
	// Check if user has admin permissions
	userRole, exists := c.Get("userRole")
	if !exists || userRole != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	result, err := h.advertiserService.SyncAllAdvertisersToProvider(c.Request.Context(), "everflow")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync advertisers to Everflow: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
