package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AnalyticsHandler handles analytics-related HTTP requests
type AnalyticsHandler struct {
	analyticsService service.AnalyticsService
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(analyticsService service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// AutocompleteRequest represents the request parameters for autocompletion
type AutocompleteRequest struct {
	Query string `form:"q" binding:"required,min=3" json:"q"`
	Type  string `form:"type" json:"type"` // "advertiser", "publisher", "both", or empty (defaults to "both")
	Limit int    `form:"limit" json:"limit"`
}

// AutocompleteOrganizations handles autocompletion search
// @Summary Search organizations for autocompletion
// @Description Search advertisers and/or publishers by domain name for autocompletion (minimum 3 characters)
// @Tags Analytics
// @Accept json
// @Produce json
// @Param q query string true "Search query (minimum 3 characters)" minlength(3)
// @Param type query string false "Organization type filter" Enums(advertiser,publisher,both) default(both)
// @Param limit query int false "Maximum number of results" default(10) minimum(1) maximum(50)
// @Success 200 {object} SuccessResponse{data=[]domain.AutocompleteResult} "Autocompletion results"
// @Failure 400 {object} ErrorResponse "Bad request - invalid parameters"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/analytics/autocomplete [get]
func (h *AnalyticsHandler) AutocompleteOrganizations(c *gin.Context) {
	var req AutocompleteRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request parameters",
			Details: err.Error(),
		})
		return
	}

	// Set default limit if not provided
	if req.Limit <= 0 {
		req.Limit = 10
	}

	results, err := h.analyticsService.SearchOrganizations(c.Request.Context(), req.Query, req.Type, req.Limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Search failed",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Autocompletion results retrieved successfully",
		"data":    results,
	})
}

// GetAdvertiserByID retrieves advertiser analytics data by ID
// @Summary Get advertiser analytics data
// @Description Retrieve detailed analytics data for a specific advertiser
// @Tags Analytics
// @Accept json
// @Produce json
// @Param id path int true "Advertiser ID"
// @Success 200 {object} SuccessResponse{data=domain.AnalyticsAdvertiserResponse} "Advertiser analytics data"
// @Failure 400 {object} ErrorResponse "Bad request - invalid ID"
// @Failure 404 {object} ErrorResponse "Advertiser not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/analytics/advertisers/{id} [get]
func (h *AnalyticsHandler) GetAdvertiserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid advertiser ID",
			Details: "ID must be a valid integer",
		})
		return
	}

	advertiser, err := h.analyticsService.GetAdvertiserByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Advertiser not found",
				Details: "No advertiser found with the specified ID",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve advertiser",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Advertiser retrieved successfully",
		"data":    advertiser,
	})
}

// GetPublisherByID retrieves publisher analytics data by ID
// @Summary Get publisher analytics data
// @Description Retrieve detailed analytics data for a specific publisher (affiliate)
// @Tags Analytics
// @Accept json
// @Produce json
// @Param id path int true "Publisher ID"
// @Success 200 {object} SuccessResponse{data=domain.AnalyticsPublisherResponse} "Publisher analytics data"
// @Failure 400 {object} ErrorResponse "Bad request - invalid ID"
// @Failure 404 {object} ErrorResponse "Publisher not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/analytics/affiliates/{id} [get]
func (h *AnalyticsHandler) GetPublisherByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid publisher ID",
			Details: "ID must be a valid integer",
		})
		return
	}

	publisher, err := h.analyticsService.GetPublisherByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Publisher not found",
				Details: "No publisher found with the specified ID",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve publisher",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Publisher retrieved successfully",
		"data":    publisher,
	})
}

// Additional CRUD endpoints for managing analytics data (optional, for data management)

// CreateAdvertiserRequest represents the request body for creating an advertiser
type CreateAdvertiserRequest struct {
	Domain string                 `json:"domain" binding:"required"`
	Data   map[string]interface{} `json:"data" binding:"required"`
}

// CreateAdvertiser creates a new advertiser analytics entry
// @Summary Create advertiser analytics data
// @Description Create a new advertiser analytics entry
// @Tags Analytics
// @Accept json
// @Produce json
// @Param advertiser body CreateAdvertiserRequest true "Advertiser data"
// @Success 201 {object} SuccessResponse{data=domain.AnalyticsAdvertiser} "Advertiser created successfully"
// @Failure 400 {object} ErrorResponse "Bad request - invalid data"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/analytics/advertisers [post]
func (h *AnalyticsHandler) CreateAdvertiser(c *gin.Context) {
	var req CreateAdvertiserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Convert data to JSON string for storage
	dataJSON, err := json.Marshal(req.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid data format",
			Details: err.Error(),
		})
		return
	}

	// Create advertiser object
	dataStr := string(dataJSON)
	advertiser := &domain.AnalyticsAdvertiser{
		Domain:         req.Domain,
		AdditionalData: &dataStr,
	}

	// Create via service
	if err := h.analyticsService.CreateAdvertiser(c.Request.Context(), advertiser); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create advertiser",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, advertiser)
}

// CreatePublisherRequest represents the request body for creating a publisher
type CreatePublisherRequest struct {
	Domain string                 `json:"domain" binding:"required"`
	Data   map[string]interface{} `json:"data" binding:"required"`
}

// CreatePublisher creates a new publisher analytics entry
// @Summary Create publisher analytics data
// @Description Create a new publisher analytics entry
// @Tags Analytics
// @Accept json
// @Produce json
// @Param publisher body CreatePublisherRequest true "Publisher data"
// @Success 201 {object} SuccessResponse{data=domain.AnalyticsPublisher} "Publisher created successfully"
// @Failure 400 {object} ErrorResponse "Bad request - invalid data"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/analytics/affiliates [post]
func (h *AnalyticsHandler) CreatePublisher(c *gin.Context) {
	var req CreatePublisherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Convert data to JSON string for storage
	dataJSON, err := json.Marshal(req.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid data format",
			Details: err.Error(),
		})
		return
	}

	// Create publisher object
	dataStr := string(dataJSON)
	publisher := &domain.AnalyticsPublisher{
		Domain:         req.Domain,
		AdditionalData: &dataStr,
	}

	// Create via service
	if err := h.analyticsService.CreatePublisher(c.Request.Context(), publisher); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create publisher",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, publisher)
}