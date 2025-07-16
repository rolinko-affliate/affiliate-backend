package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// FavoritePublisherListHandler handles favorite publisher list-related HTTP requests
type FavoritePublisherListHandler struct {
	favoriteListService service.FavoritePublisherListService
}

// NewFavoritePublisherListHandler creates a new favorite publisher list handler
func NewFavoritePublisherListHandler(favoriteListService service.FavoritePublisherListService) *FavoritePublisherListHandler {
	return &FavoritePublisherListHandler{
		favoriteListService: favoriteListService,
	}
}

// CreateList creates a new favorite publisher list
// @Summary Create a new favorite publisher list
// @Description Creates a new favorite publisher list for the user's organization
// @Tags favorite-publisher-lists
// @Accept json
// @Produce json
// @Param request body domain.CreateFavoritePublisherListRequest true "Create list request"
// @Success 201 {object} gin.H{message=string,data=domain.FavoritePublisherList}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /favorite-publisher-lists [post]
func (h *FavoritePublisherListHandler) CreateList(c *gin.Context) {
	// Get organization ID from context (set by RBAC middleware)
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Details: "Organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Details: "Invalid organization ID type",
		})
		return
	}

	var req domain.CreateFavoritePublisherListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	list, err := h.favoriteListService.CreateList(c.Request.Context(), orgID, &req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid input",
				Details: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create favorite publisher list",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Favorite publisher list created successfully",
		"data":    list,
	})
}

// GetLists retrieves all favorite publisher lists for the organization
// @Summary Get all favorite publisher lists
// @Description Retrieves all favorite publisher lists for the user's organization with statistics
// @Tags favorite-publisher-lists
// @Produce json
// @Success 200 {object} gin.H{message=string,data=[]domain.FavoritePublisherListWithStats}
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /favorite-publisher-lists [get]
func (h *FavoritePublisherListHandler) GetLists(c *gin.Context) {
	// Get organization ID from context (set by RBAC middleware)
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Details: "Organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Details: "Invalid organization ID type",
		})
		return
	}

	lists, err := h.favoriteListService.GetListsByOrganization(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve favorite publisher lists",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Favorite publisher lists retrieved successfully",
		"data":    lists,
	})
}

// GetListByID retrieves a specific favorite publisher list by ID
// @Summary Get favorite publisher list by ID
// @Description Retrieves a specific favorite publisher list by ID
// @Tags favorite-publisher-lists
// @Produce json
// @Param list_id path int true "List ID"
// @Success 200 {object} gin.H{message=string,data=domain.FavoritePublisherList}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /favorite-publisher-lists/{list_id} [get]
func (h *FavoritePublisherListHandler) GetListByID(c *gin.Context) {
	// Get organization ID from context (set by RBAC middleware)
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Details: "Organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Details: "Invalid organization ID type",
		})
		return
	}

	listIDStr := c.Param("list_id")
	listID, err := strconv.ParseInt(listIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid list ID",
			Details: "List ID must be a valid integer",
		})
		return
	}

	list, err := h.favoriteListService.GetListByID(c.Request.Context(), orgID, listID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "List not found",
				Details: "No favorite publisher list found with the specified ID",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve favorite publisher list",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Favorite publisher list retrieved successfully",
		"data":    list,
	})
}

// UpdateList updates a favorite publisher list
// @Summary Update favorite publisher list
// @Description Updates a favorite publisher list's name and/or description
// @Tags favorite-publisher-lists
// @Accept json
// @Produce json
// @Param list_id path int true "List ID"
// @Param request body domain.UpdateFavoritePublisherListRequest true "Update list request"
// @Success 200 {object} gin.H{message=string,data=domain.FavoritePublisherList}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /favorite-publisher-lists/{list_id} [put]
func (h *FavoritePublisherListHandler) UpdateList(c *gin.Context) {
	// Get organization ID from context (set by RBAC middleware)
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Details: "Organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Details: "Invalid organization ID type",
		})
		return
	}

	listIDStr := c.Param("list_id")
	listID, err := strconv.ParseInt(listIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid list ID",
			Details: "List ID must be a valid integer",
		})
		return
	}

	var req domain.UpdateFavoritePublisherListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	list, err := h.favoriteListService.UpdateList(c.Request.Context(), orgID, listID, &req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "List not found",
				Details: "No favorite publisher list found with the specified ID",
			})
			return
		}
		if errors.Is(err, domain.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid input",
				Details: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update favorite publisher list",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Favorite publisher list updated successfully",
		"data":    list,
	})
}

// DeleteList deletes a favorite publisher list
// @Summary Delete favorite publisher list
// @Description Deletes a favorite publisher list and all its items
// @Tags favorite-publisher-lists
// @Produce json
// @Param list_id path int true "List ID"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /favorite-publisher-lists/{list_id} [delete]
func (h *FavoritePublisherListHandler) DeleteList(c *gin.Context) {
	// Get organization ID from context (set by RBAC middleware)
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Details: "Organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Details: "Invalid organization ID type",
		})
		return
	}

	listIDStr := c.Param("list_id")
	listID, err := strconv.ParseInt(listIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid list ID",
			Details: "List ID must be a valid integer",
		})
		return
	}

	err = h.favoriteListService.DeleteList(c.Request.Context(), orgID, listID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "List not found",
				Details: "No favorite publisher list found with the specified ID",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to delete favorite publisher list",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Favorite publisher list deleted successfully",
	})
}

// AddPublisherToList adds a publisher to a favorite list
// @Summary Add publisher to favorite list
// @Description Adds a publisher domain to a favorite publisher list
// @Tags favorite-publisher-lists
// @Accept json
// @Produce json
// @Param list_id path int true "List ID"
// @Param request body domain.AddPublisherToListRequest true "Add publisher request"
// @Success 201 {object} gin.H{message=string,data=domain.FavoritePublisherListItem}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /favorite-publisher-lists/{list_id}/publishers [post]
func (h *FavoritePublisherListHandler) AddPublisherToList(c *gin.Context) {
	// Get organization ID from context (set by RBAC middleware)
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Details: "Organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Details: "Invalid organization ID type",
		})
		return
	}

	listIDStr := c.Param("list_id")
	listID, err := strconv.ParseInt(listIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid list ID",
			Details: "List ID must be a valid integer",
		})
		return
	}

	var req domain.AddPublisherToListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	item, err := h.favoriteListService.AddPublisherToList(c.Request.Context(), orgID, listID, &req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "List not found",
				Details: "No favorite publisher list found with the specified ID",
			})
			return
		}
		if errors.Is(err, domain.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid input",
				Details: err.Error(),
			})
			return
		}
		// Check if it's a duplicate publisher error
		if err.Error() == "publisher "+req.PublisherDomain+" is already in the list" {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "Publisher already in list",
				Details: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to add publisher to list",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Publisher added to list successfully",
		"data":    item,
	})
}

// RemovePublisherFromList removes a publisher from a favorite list
// @Summary Remove publisher from favorite list
// @Description Removes a publisher domain from a favorite publisher list
// @Tags favorite-publisher-lists
// @Produce json
// @Param list_id path int true "List ID"
// @Param domain path string true "Publisher domain"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /favorite-publisher-lists/{list_id}/publishers/{domain} [delete]
func (h *FavoritePublisherListHandler) RemovePublisherFromList(c *gin.Context) {
	// Get organization ID from context (set by RBAC middleware)
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Details: "Organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Details: "Invalid organization ID type",
		})
		return
	}

	listIDStr := c.Param("list_id")
	listID, err := strconv.ParseInt(listIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid list ID",
			Details: "List ID must be a valid integer",
		})
		return
	}

	domainParam := c.Param("domain")
	if domainParam == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid domain",
			Details: "Domain parameter is required",
		})
		return
	}

	err = h.favoriteListService.RemovePublisherFromList(c.Request.Context(), orgID, listID, domainParam)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Publisher not found in list",
				Details: "No publisher found with the specified domain in this list",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to remove publisher from list",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Publisher removed from list successfully",
	})
}

// GetListItems retrieves all items in a favorite list
// @Summary Get favorite list items
// @Description Retrieves all publisher items in a favorite publisher list
// @Tags favorite-publisher-lists
// @Produce json
// @Param list_id path int true "List ID"
// @Param include_details query bool false "Include publisher details from analytics"
// @Success 200 {object} gin.H{message=string,data=[]domain.FavoritePublisherListItem}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /favorite-publisher-lists/{list_id}/publishers [get]
func (h *FavoritePublisherListHandler) GetListItems(c *gin.Context) {
	// Get organization ID from context (set by RBAC middleware)
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Details: "Organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Details: "Invalid organization ID type",
		})
		return
	}

	listIDStr := c.Param("list_id")
	listID, err := strconv.ParseInt(listIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid list ID",
			Details: "List ID must be a valid integer",
		})
		return
	}

	includeDetails := c.Query("include_details") == "true"

	items, err := h.favoriteListService.GetListItems(c.Request.Context(), orgID, listID, includeDetails)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "List not found",
				Details: "No favorite publisher list found with the specified ID",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve list items",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "List items retrieved successfully",
		"data":    items,
	})
}

// UpdatePublisherInList updates a publisher's notes in a favorite list
// @Summary Update publisher notes in favorite list
// @Description Updates the notes for a publisher in a favorite publisher list
// @Tags favorite-publisher-lists
// @Accept json
// @Produce json
// @Param list_id path int true "List ID"
// @Param domain path string true "Publisher domain"
// @Param request body domain.UpdatePublisherInListRequest true "Update publisher request"
// @Success 200 {object} gin.H{message=string}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /favorite-publisher-lists/{list_id}/publishers/{domain} [put]
func (h *FavoritePublisherListHandler) UpdatePublisherInList(c *gin.Context) {
	// Get organization ID from context (set by RBAC middleware)
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Details: "Organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Details: "Invalid organization ID type",
		})
		return
	}

	listIDStr := c.Param("list_id")
	listID, err := strconv.ParseInt(listIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid list ID",
			Details: "List ID must be a valid integer",
		})
		return
	}

	domainParam := c.Param("domain")
	if domainParam == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid domain",
			Details: "Domain parameter is required",
		})
		return
	}

	var req domain.UpdatePublisherInListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	err = h.favoriteListService.UpdatePublisherInList(c.Request.Context(), orgID, listID, domainParam, &req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Publisher not found in list",
				Details: "No publisher found with the specified domain in this list",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update publisher in list",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Publisher notes updated successfully",
	})
}

// GetListsContainingPublisher retrieves all lists that contain a specific publisher
// @Summary Get lists containing publisher
// @Description Retrieves all favorite publisher lists that contain a specific publisher domain
// @Tags favorite-publisher-lists
// @Produce json
// @Param domain query string true "Publisher domain"
// @Success 200 {object} gin.H{message=string,data=[]domain.FavoritePublisherList}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /favorite-publisher-lists/search [get]
func (h *FavoritePublisherListHandler) GetListsContainingPublisher(c *gin.Context) {
	// Get organization ID from context (set by RBAC middleware)
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Details: "Organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Details: "Invalid organization ID type",
		})
		return
	}

	domainParam := c.Query("domain")
	if domainParam == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid domain",
			Details: "Domain query parameter is required",
		})
		return
	}

	lists, err := h.favoriteListService.GetListsContainingPublisher(c.Request.Context(), orgID, domainParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve lists containing publisher",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lists containing publisher retrieved successfully",
		"data":    lists,
	})
}