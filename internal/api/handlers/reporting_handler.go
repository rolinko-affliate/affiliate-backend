package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// ReportingHandler handles reporting API endpoints
type ReportingHandler struct {
	reportingService service.ReportingService
}

// NewReportingHandler creates a new reporting handler
func NewReportingHandler(reportingService service.ReportingService) *ReportingHandler {
	return &ReportingHandler{
		reportingService: reportingService,
	}
}

// GetPerformanceSummary handles GET /api/v1/reports/performance/summary
func (h *ReportingHandler) GetPerformanceSummary(c *gin.Context) {
	// Parse query parameters
	filters, err := h.parseReportingFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid query parameters",
			"code":    "INVALID_PARAMETERS",
			"details": err.Error(),
		})
		return
	}

	// Validate date range
	if err := h.validateDateRange(filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid date range",
			"code":    "INVALID_DATE_RANGE",
			"details": err.Error(),
		})
		return
	}

	// Get user profile from context
	userProfile, exists := c.Get("profile")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User profile not found",
			"code":    "UNAUTHORIZED",
		})
		return
	}

	profile := userProfile.(*domain.Profile)

	// Get performance summary
	summary, err := h.reportingService.GetPerformanceSummary(c.Request.Context(), filters, profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get performance summary",
			"code":    "INTERNAL_ERROR",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": summary,
		"dateRange": gin.H{
			"startDate": filters.StartDate,
			"endDate":   filters.EndDate,
		},
		"status": "success",
	})
}

// GetPerformanceTimeSeries handles GET /api/v1/reports/performance/timeseries
func (h *ReportingHandler) GetPerformanceTimeSeries(c *gin.Context) {
	// Parse query parameters
	filters, err := h.parseReportingFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid query parameters",
			"code":    "INVALID_PARAMETERS",
			"details": err.Error(),
		})
		return
	}

	// Validate date range
	if err := h.validateDateRange(filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid date range",
			"code":    "INVALID_DATE_RANGE",
			"details": err.Error(),
		})
		return
	}

	// Get user profile from context
	userProfile, exists := c.Get("profile")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User profile not found",
			"code":    "UNAUTHORIZED",
		})
		return
	}

	profile := userProfile.(*domain.Profile)

	// Get time series data
	timeSeries, err := h.reportingService.GetPerformanceTimeSeries(c.Request.Context(), filters, profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get performance time series",
			"code":    "INTERNAL_ERROR",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   timeSeries,
		"status": "success",
	})
}

// GetDailyPerformanceReport handles GET /api/v1/reports/performance/daily
func (h *ReportingHandler) GetDailyPerformanceReport(c *gin.Context) {
	// Parse query parameters
	filters, err := h.parseReportingFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid query parameters",
			"code":    "INVALID_PARAMETERS",
			"details": err.Error(),
		})
		return
	}

	// Parse pagination parameters
	pagination, err := h.parsePaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid pagination parameters",
			"code":    "INVALID_PARAMETERS",
			"details": err.Error(),
		})
		return
	}

	// Validate date range
	if err := h.validateDateRange(filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid date range",
			"code":    "INVALID_DATE_RANGE",
			"details": err.Error(),
		})
		return
	}

	// Get user profile from context
	userProfile, exists := c.Get("profile")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User profile not found",
			"code":    "UNAUTHORIZED",
		})
		return
	}

	profile := userProfile.(*domain.Profile)

	// Get daily performance report
	reports, paginationResult, err := h.reportingService.GetDailyPerformanceReport(c.Request.Context(), filters, pagination, profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get daily performance report",
			"code":    "INTERNAL_ERROR",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       reports,
		"pagination": paginationResult,
		"status":     "success",
	})
}

// GetConversionsReport handles GET /api/v1/reports/conversions
func (h *ReportingHandler) GetConversionsReport(c *gin.Context) {
	// Parse query parameters
	filters, err := h.parseReportingFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid query parameters",
			"code":    "INVALID_PARAMETERS",
			"details": err.Error(),
		})
		return
	}

	// Parse pagination parameters
	pagination, err := h.parsePaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid pagination parameters",
			"code":    "INVALID_PARAMETERS",
			"details": err.Error(),
		})
		return
	}

	// Validate date range
	if err := h.validateDateRange(filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid date range",
			"code":    "INVALID_DATE_RANGE",
			"details": err.Error(),
		})
		return
	}

	// Get user profile from context
	userProfile, exists := c.Get("profile")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User profile not found",
			"code":    "UNAUTHORIZED",
		})
		return
	}

	profile := userProfile.(*domain.Profile)

	// Get conversions report
	conversions, paginationResult, err := h.reportingService.GetConversionsReport(c.Request.Context(), filters, pagination, profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get conversions report",
			"code":    "INTERNAL_ERROR",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       conversions,
		"pagination": paginationResult,
		"status":     "success",
	})
}

// GetClicksReport handles GET /api/v1/reports/clicks
func (h *ReportingHandler) GetClicksReport(c *gin.Context) {
	// Parse query parameters
	filters, err := h.parseReportingFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid query parameters",
			"code":    "INVALID_PARAMETERS",
			"details": err.Error(),
		})
		return
	}

	// Parse pagination parameters
	pagination, err := h.parsePaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid pagination parameters",
			"code":    "INVALID_PARAMETERS",
			"details": err.Error(),
		})
		return
	}

	// Validate date range
	if err := h.validateDateRange(filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid date range",
			"code":    "INVALID_DATE_RANGE",
			"details": err.Error(),
		})
		return
	}

	// Get user profile from context
	userProfile, exists := c.Get("profile")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User profile not found",
			"code":    "UNAUTHORIZED",
		})
		return
	}

	profile := userProfile.(*domain.Profile)

	// Get clicks report
	clicks, paginationResult, err := h.reportingService.GetClicksReport(c.Request.Context(), filters, pagination, profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get clicks report",
			"code":    "INTERNAL_ERROR",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       clicks,
		"pagination": paginationResult,
		"status":     "success",
	})
}

// GetCampaignsList handles GET /api/v1/campaigns
func (h *ReportingHandler) GetCampaignsList(c *gin.Context) {
	// Parse query parameters
	affiliateID := c.Query("affiliateId")
	status := c.DefaultQuery("status", "active")
	search := c.Query("search")

	var affiliateIDPtr *string
	if affiliateID != "" {
		affiliateIDPtr = &affiliateID
	}

	var searchPtr *string
	if search != "" {
		searchPtr = &search
	}

	// Get user profile from context
	userProfile, exists := c.Get("profile")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User profile not found",
			"code":    "UNAUTHORIZED",
		})
		return
	}

	profile := userProfile.(*domain.Profile)

	// Get campaigns list
	campaigns, err := h.reportingService.GetCampaignsList(c.Request.Context(), affiliateIDPtr, status, searchPtr, profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get campaigns list",
			"code":    "INTERNAL_ERROR",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   campaigns,
		"status": "success",
	})
}

// Helper methods

func (h *ReportingHandler) parseReportingFilters(c *gin.Context) (domain.ReportingFilters, error) {
	filters := domain.ReportingFilters{
		StartDate: c.Query("startDate"),
		EndDate:   c.Query("endDate"),
	}

	// Parse campaign IDs
	if campaignIDsStr := c.Query("campaignIds"); campaignIDsStr != "" {
		filters.CampaignIDs = strings.Split(campaignIDsStr, ",")
	}

	// Parse affiliate ID
	if affiliateID := c.Query("affiliateId"); affiliateID != "" {
		filters.AffiliateID = &affiliateID
	}

	// Parse search
	if search := c.Query("search"); search != "" {
		filters.Search = &search
	}

	// Parse status
	if status := c.Query("status"); status != "" {
		filters.Status = &status
	}

	// Parse granularity
	if granularity := c.Query("granularity"); granularity != "" {
		filters.Granularity = &granularity
	}

	return filters, nil
}

func (h *ReportingHandler) parsePaginationParams(c *gin.Context) (domain.PaginationParams, error) {
	pagination := domain.PaginationParams{
		Page:      1,
		Limit:     10,
		SortBy:    "date",
		SortOrder: "desc",
	}

	// Parse page
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			pagination.Page = page
		}
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			pagination.Limit = limit
		}
	}

	// Parse sort by
	if sortBy := c.Query("sortBy"); sortBy != "" {
		pagination.SortBy = sortBy
	}

	// Parse sort order
	if sortOrder := c.Query("sortOrder"); sortOrder != "" {
		pagination.SortOrder = sortOrder
	}

	return pagination, nil
}

func (h *ReportingHandler) validateDateRange(filters domain.ReportingFilters) error {
	// Parse start and end dates
	startDate, err := time.Parse("2006-01-02", filters.StartDate)
	if err != nil {
		return err
	}

	endDate, err := time.Parse("2006-01-02", filters.EndDate)
	if err != nil {
		return err
	}

	// Check if start date is before end date
	if startDate.After(endDate) {
		return fmt.Errorf("start date must be before end date")
	}

	// Check if date range is not too large (max 1 year)
	if endDate.Sub(startDate) > 365*24*time.Hour {
		return fmt.Errorf("date range cannot exceed 1 year")
	}

	return nil
}