package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/affiliate-backend/internal/api/middleware"
	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/logger"
	"github.com/affiliate-backend/internal/service"
)

// DashboardHandler handles dashboard-related HTTP requests
type DashboardHandler struct {
	service service.DashboardService
	logger  *logger.Logger
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(service service.DashboardService, log *logger.Logger) *DashboardHandler {
	return &DashboardHandler{
		service: service,
		logger:  log,
	}
}

// GetDashboard handles GET /dashboard
// @Summary      Get dashboard data
// @Description  Returns dashboard data based on user's organization type
// @Tags         dashboard
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        period      query     string  false  "Time period (today, 7d, 30d, 90d, custom)"  default(30d)
// @Param        start_date  query     string  false  "Start date for custom period (YYYY-MM-DD)"
// @Param        end_date    query     string  false  "End date for custom period (YYYY-MM-DD)"
// @Param        timezone    query     string  false  "Timezone identifier"  default(UTC)
// @Success      200         {object}  domain.DashboardData
// @Failure      400         {object}  ErrorResponse
// @Failure      401         {object}  ErrorResponse
// @Failure      403         {object}  ErrorResponse
// @Failure      500         {object}  ErrorResponse
// @Router       /dashboard [get]
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get(middleware.UserIDKey)
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	// Parse query parameters
	period := c.DefaultQuery("period", "30d")
	timezone := c.DefaultQuery("timezone", "UTC")

	var startDate, endDate *time.Time

	// Parse custom date range if provided
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err != nil {
			RespondWithError(c, http.StatusBadRequest, "Invalid start_date format, expected YYYY-MM-DD")
			return
		} else {
			startDate = &parsed
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err != nil {
			RespondWithError(c, http.StatusBadRequest, "Invalid end_date format, expected YYYY-MM-DD")
			return
		} else {
			endDate = &parsed
		}
	}

	// Get dashboard data
	dashboardData, err := h.service.GetDashboardData(c.Request.Context(), userID, period, startDate, endDate, timezone)
	if err != nil {
		h.logger.Error("Failed to get dashboard data",
			"user_id", userID,
			"error", err,
		)

		if errors.Is(err, domain.ErrNotFound) {
			RespondWithError(c, http.StatusNotFound, err.Error())
			return
		}

		if errors.Is(err, domain.ErrUnauthorized) || errors.Is(err, domain.ErrForbidden) {
			RespondWithError(c, http.StatusForbidden, err.Error())
			return
		}

		if errors.Is(err, domain.ErrInvalidInput) {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}

		RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve dashboard data")
		return
	}

	c.JSON(http.StatusOK, dashboardData)
}

// GetCampaignDetail handles GET /dashboard/campaigns/:campaignId
// @Summary      Get campaign detail
// @Description  Returns detailed performance data for a specific campaign
// @Tags         dashboard
// @Accept       json
// @Produce      json
// @Param        campaignId  path      int  true  "Campaign ID"
// @Success      200         {object}  domain.CampaignDetail
// @Failure      400         {object}  ErrorResponse
// @Failure      401         {object}  ErrorResponse
// @Failure      403         {object}  ErrorResponse
// @Failure      404         {object}  ErrorResponse
// @Failure      500         {object}  ErrorResponse
// @Router       /dashboard/campaigns/{campaignId} [get]
func (h *DashboardHandler) GetCampaignDetail(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get(middleware.UserIDKey)
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	// Parse campaign ID
	campaignIDStr := c.Param("id")
	campaignID, err := strconv.ParseInt(campaignIDStr, 10, 64)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid campaign ID format")
		return
	}

	// Get campaign detail
	detail, err := h.service.GetCampaignDetail(c.Request.Context(), userID, campaignID)
	if err != nil {
		h.logger.Error("Failed to get campaign detail",
			"user_id", userID,
			"campaign_id", campaignID,
			"error", err,
		)

		if errors.Is(err, domain.ErrInvalidInput) {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, domain.ErrUnauthorized) || errors.Is(err, domain.ErrForbidden) {
			RespondWithError(c, http.StatusForbidden, err.Error())
			return
		}

		if errors.Is(err, domain.ErrNotFound) {
			RespondWithError(c, http.StatusNotFound, err.Error())
			return
		}

		RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve campaign detail")
		return
	}

	c.JSON(http.StatusOK, detail)
}

// GetRecentActivity handles GET /dashboard/activity
// @Summary      Get recent activity
// @Description  Returns paginated recent activity feed
// @Tags         dashboard
// @Accept       json
// @Produce      json
// @Param        limit   query     int       false  "Number of items per page"  default(10)  maximum(50)
// @Param        offset  query     int       false  "Number of items to skip"   default(0)
// @Param        type    query     []string  false  "Filter by activity types"
// @Param        since   query     string    false  "Filter activities since this timestamp (RFC3339)"
// @Success      200     {object}  domain.ActivityResponse
// @Failure      400     {object}  ErrorResponse
// @Failure      401     {object}  ErrorResponse
// @Failure      403     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /dashboard/activity [get]
func (h *DashboardHandler) GetRecentActivity(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get(middleware.UserIDKey)
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	// Parse pagination parameters
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	// Parse activity types filter
	var activityTypes []string
	if typeStr := c.Query("type"); typeStr != "" {
		activityTypes = strings.Split(typeStr, ",")
	}

	// Parse since parameter
	var since *time.Time
	if sinceStr := c.Query("since"); sinceStr != "" {
		if parsed, err := time.Parse(time.RFC3339, sinceStr); err != nil {
			RespondWithError(c, http.StatusBadRequest, "Invalid since format, expected RFC3339")
			return
		} else {
			since = &parsed
		}
	}

	// Get recent activity
	response, err := h.service.GetRecentActivity(c.Request.Context(), userID, limit, offset, activityTypes, since)
	if err != nil {
		h.logger.Error("Failed to get recent activity",
			"user_id", userID,
			"error", err,
		)

		if errors.Is(err, domain.ErrInvalidInput) {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, domain.ErrUnauthorized) || errors.Is(err, domain.ErrForbidden) {
			RespondWithError(c, http.StatusForbidden, err.Error())
			return
		}

		RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve recent activity")
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetSystemHealth handles GET /dashboard/system/health
// @Summary      Get system health metrics
// @Description  Returns system health metrics (Platform Owner only)
// @Tags         dashboard
// @Accept       json
// @Produce      json
// @Success      200  {object}  domain.SystemHealth
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /dashboard/system/health [get]
func (h *DashboardHandler) GetSystemHealth(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get(middleware.UserIDKey)
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	// Get system health
	health, err := h.service.GetSystemHealth(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get system health",
			"user_id", userID,
			"error", err,
		)

		if errors.Is(err, domain.ErrUnauthorized) || errors.Is(err, domain.ErrForbidden) {
			RespondWithError(c, http.StatusForbidden, err.Error())
			return
		}

		RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve system health")
		return
	}

	c.JSON(http.StatusOK, health)
}

// DashboardHealthCheck handles GET /dashboard/health
// @Summary      Dashboard health check
// @Description  Returns the health status of the dashboard service
// @Tags         dashboard
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Health status"
// @Failure      503  {object}  map[string]interface{}  "Service unavailable"
// @Router       /dashboard/health [get]
func (h *DashboardHandler) DashboardHealthCheck(c *gin.Context) {
	ctx := c.Request.Context()

	// Check database connectivity
	if err := h.service.CheckDatabaseHealth(ctx); err != nil {
		h.logger.Error("Dashboard database health check failed", "error", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "DOWN",
			"checks": gin.H{
				"database": "FAIL",
				"error":    err.Error(),
			},
		})
		return
	}

	// Check cache connectivity
	if err := h.service.CheckCacheHealth(ctx); err != nil {
		h.logger.Error("Dashboard cache health check failed", "error", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "DOWN",
			"checks": gin.H{
				"database": "OK",
				"cache":    "FAIL",
				"error":    err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
		"checks": gin.H{
			"database": "OK",
			"cache":    "OK",
		},
	})
}

// InvalidateCache handles POST /dashboard/cache/invalidate
// @Summary      Invalidate dashboard cache
// @Description  Invalidates dashboard cache for the user's organization
// @Tags         dashboard
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string  "Cache invalidated"
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /dashboard/cache/invalidate [post]
func (h *DashboardHandler) InvalidateCache(c *gin.Context) {
	// Get user profile from context (set by profile middleware)
	profile, exists := c.Get("profile")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, "User profile not found in context")
		return
	}

	userProfile := profile.(*domain.Profile)
	if userProfile.OrganizationID == nil {
		RespondWithError(c, http.StatusBadRequest, "User is not associated with any organization")
		return
	}

	orgID := *userProfile.OrganizationID

	// Invalidate cache
	if err := h.service.InvalidateCache(c.Request.Context(), orgID); err != nil {
		h.logger.Error("Failed to invalidate dashboard cache",
			"org_id", orgID,
			"user_id", userProfile.ID,
			"error", err,
		)
		RespondWithError(c, http.StatusInternalServerError, "Failed to invalidate cache")
		return
	}

	h.logger.Info("Dashboard cache invalidated",
		"org_id", orgID,
		"user_id", userProfile.ID,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Dashboard cache invalidated successfully",
	})
}

// TrackActivity handles POST /dashboard/activity
// @Summary      Track dashboard activity
// @Description  Creates a new activity record for the user's organization
// @Tags         dashboard
// @Accept       json
// @Produce      json
// @Param        request  body      TrackActivityRequest  true  "Activity data"
// @Success      201      {object}  map[string]string     "Activity tracked"
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /dashboard/activity [post]
func (h *DashboardHandler) TrackActivity(c *gin.Context) {
	// Get user profile from context
	profile, exists := c.Get("profile")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, "User profile not found in context")
		return
	}

	userProfile := profile.(*domain.Profile)
	if userProfile.OrganizationID == nil {
		RespondWithError(c, http.StatusBadRequest, "User is not associated with any organization")
		return
	}

	orgID := *userProfile.OrganizationID

	// Parse request body
	var req TrackActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validate activity type
	activityType := domain.ActivityType(req.Type)
	if !activityType.IsValid() {
		RespondWithError(c, http.StatusBadRequest, "Invalid activity type")
		return
	}

	// Track activity
	if err := h.service.TrackActivityWithUser(c.Request.Context(), userProfile.ID, orgID, activityType, req.Description, req.Metadata); err != nil {
		h.logger.Error("Failed to track activity",
			"org_id", orgID,
			"user_id", userProfile.ID,
			"activity_type", req.Type,
			"error", err,
		)
		RespondWithError(c, http.StatusInternalServerError, "Failed to track activity")
		return
	}

	h.logger.Info("Activity tracked",
		"org_id", orgID,
		"user_id", userProfile.ID,
		"activity_type", req.Type,
	)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Activity tracked successfully",
	})
}

// TrackActivityRequest represents the request body for tracking activity
type TrackActivityRequest struct {
	Type        string                 `json:"type" binding:"required"`
	Description string                 `json:"description" binding:"required"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// DashboardErrorResponse represents dashboard-specific error responses
type DashboardErrorResponse struct {
	Error     string    `json:"error"`
	Details   string    `json:"details,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"request_id,omitempty"`
}

// Common dashboard error codes
const (
	ErrCodeDashboardUnauthorized    = "DASHBOARD_UNAUTHORIZED"
	ErrCodeDashboardNotFound        = "DASHBOARD_NOT_FOUND"
	ErrCodeDashboardInvalidPeriod   = "DASHBOARD_INVALID_PERIOD"
	ErrCodeDashboardCacheError      = "DASHBOARD_CACHE_ERROR"
	ErrCodeDashboardDataUnavailable = "DASHBOARD_DATA_UNAVAILABLE"
)

// RespondWithDashboardError sends a dashboard-specific error response
func RespondWithDashboardError(c *gin.Context, code int, errorCode, message, details string) {
	requestID, _ := c.Get("request_id")
	requestIDStr, _ := requestID.(string)

	resp := DashboardErrorResponse{
		Error:     message,
		Details:   details,
		Timestamp: time.Now(),
		RequestID: requestIDStr,
	}

	c.Header("X-Error-Code", errorCode)
	c.AbortWithStatusJSON(code, resp)
}