package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// Helper functions for type conversion
func intToInt32Ptr(i *int) *int32 {
	if i == nil {
		return nil
	}
	val := int32(*i)
	return &val
}

func int32ToIntPtr(i *int32) *int {
	if i == nil {
		return nil
	}
	val := int(*i)
	return &val
}

// CampaignHandler handles campaign-related requests
type CampaignHandler struct {
	campaignService service.CampaignService
}

// NewCampaignHandler creates a new campaign handler
func NewCampaignHandler(cs service.CampaignService) *CampaignHandler {
	return &CampaignHandler{campaignService: cs}
}

// CreateCampaignRequest defines the request for creating a campaign
type CreateCampaignRequest struct {
	OrganizationID int64            `json:"organization_id" binding:"required"`
	AdvertiserID   int64            `json:"advertiser_id" binding:"required"`
	Name           string           `json:"name" binding:"required"`
	Description    *string          `json:"description,omitempty"`
	Status         string           `json:"status,omitempty"`
	StartDate      *string          `json:"start_date,omitempty"` // Format: YYYY-MM-DD
	EndDate        *string          `json:"end_date,omitempty"`   // Format: YYYY-MM-DD
	
	// Offer-specific fields for Everflow integration
	DestinationURL      *string  `json:"destination_url,omitempty"`
	ThumbnailURL        *string  `json:"thumbnail_url,omitempty"`
	PreviewURL          *string  `json:"preview_url,omitempty"`
	Visibility          *string  `json:"visibility,omitempty"`                   // 'public', 'require_approval', 'private'
	CurrencyID          *string  `json:"currency_id,omitempty"`                 // 'USD', 'EUR', etc.
	ConversionMethod    *string  `json:"conversion_method,omitempty"`           // 'server_postback', 'pixel', etc.
	SessionDefinition   *string  `json:"session_definition,omitempty"`          // 'cookie', 'ip', 'fingerprint'
	SessionDuration     *int     `json:"session_duration,omitempty"`            // in hours
	InternalNotes       *string  `json:"internal_notes,omitempty"`
	TermsAndConditions  *string  `json:"terms_and_conditions,omitempty"`
	IsForceTermsAndConditions *bool `json:"is_force_terms_and_conditions,omitempty"`
	
	// Caps and limits
	IsCapsEnabled         *bool `json:"is_caps_enabled,omitempty"`
	DailyConversionCap    *int  `json:"daily_conversion_cap,omitempty"`
	WeeklyConversionCap   *int  `json:"weekly_conversion_cap,omitempty"`
	MonthlyConversionCap  *int  `json:"monthly_conversion_cap,omitempty"`
	GlobalConversionCap   *int  `json:"global_conversion_cap,omitempty"`
	DailyClickCap         *int  `json:"daily_click_cap,omitempty"`
	WeeklyClickCap        *int  `json:"weekly_click_cap,omitempty"`
	MonthlyClickCap       *int  `json:"monthly_click_cap,omitempty"`
	GlobalClickCap        *int  `json:"global_click_cap,omitempty"`
	
	// Payout and revenue configuration
	PayoutType     *string  `json:"payout_type,omitempty"`         // 'cpa', 'cpc', 'cpm', etc.
	PayoutAmount   *float64 `json:"payout_amount,omitempty"`
	RevenueType    *string  `json:"revenue_type,omitempty"`        // 'rpa', 'rpc', 'rpm', etc.
	RevenueAmount  *float64 `json:"revenue_amount,omitempty"`
	
	// Additional configuration
	OfferConfig *json.RawMessage `json:"offer_config,omitempty"` // Additional Everflow-specific config
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// validateCreateCampaignRequest validates the create campaign request
func validateCreateCampaignRequest(req *CreateCampaignRequest) []ValidationError {
	var errors []ValidationError

	// Validate required fields
	if req.OrganizationID <= 0 {
		errors = append(errors, ValidationError{
			Field:   "organization_id",
			Message: "organization_id must be a positive integer",
		})
	}

	if req.AdvertiserID <= 0 {
		errors = append(errors, ValidationError{
			Field:   "advertiser_id",
			Message: "advertiser_id must be a positive integer",
		})
	}

	if strings.TrimSpace(req.Name) == "" {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "name is required and cannot be empty",
		})
	}

	if len(req.Name) > 255 {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "name cannot exceed 255 characters",
		})
	}

	// Validate status if provided
	if req.Status != "" {
		validStatuses := []string{"draft", "active", "paused", "archived"}
		if !contains(validStatuses, req.Status) {
			errors = append(errors, ValidationError{
				Field:   "status",
				Message: fmt.Sprintf("status must be one of: %s", strings.Join(validStatuses, ", ")),
			})
		}
	}

	// Validate visibility if provided
	if req.Visibility != nil && *req.Visibility != "" {
		validVisibilities := []string{"public", "require_approval", "private"}
		if !contains(validVisibilities, *req.Visibility) {
			errors = append(errors, ValidationError{
				Field:   "visibility",
				Message: fmt.Sprintf("visibility must be one of: %s", strings.Join(validVisibilities, ", ")),
			})
		}
	}

	// Validate conversion method if provided
	if req.ConversionMethod != nil && *req.ConversionMethod != "" {
		validMethods := []string{"server_postback", "pixel", "hybrid"}
		if !contains(validMethods, *req.ConversionMethod) {
			errors = append(errors, ValidationError{
				Field:   "conversion_method",
				Message: fmt.Sprintf("conversion_method must be one of: %s", strings.Join(validMethods, ", ")),
			})
		}
	}

	// Validate session definition if provided
	if req.SessionDefinition != nil && *req.SessionDefinition != "" {
		validDefinitions := []string{"cookie", "ip", "fingerprint"}
		if !contains(validDefinitions, *req.SessionDefinition) {
			errors = append(errors, ValidationError{
				Field:   "session_definition",
				Message: fmt.Sprintf("session_definition must be one of: %s", strings.Join(validDefinitions, ", ")),
			})
		}
	}

	// Validate session duration if provided
	if req.SessionDuration != nil && *req.SessionDuration <= 0 {
		errors = append(errors, ValidationError{
			Field:   "session_duration",
			Message: "session_duration must be a positive integer (hours)",
		})
	}

	// Validate payout type if provided
	if req.PayoutType != nil && *req.PayoutType != "" {
		validPayoutTypes := []string{"cpa", "cpc", "cpm", "cps", "cpa_cps", "prv"}
		if !contains(validPayoutTypes, *req.PayoutType) {
			errors = append(errors, ValidationError{
				Field:   "payout_type",
				Message: fmt.Sprintf("payout_type must be one of: %s", strings.Join(validPayoutTypes, ", ")),
			})
		}
	}

	// Validate revenue type if provided
	if req.RevenueType != nil && *req.RevenueType != "" {
		validRevenueTypes := []string{"rpa", "rpc", "rpm", "rps", "rpa_rps", "prv"}
		if !contains(validRevenueTypes, *req.RevenueType) {
			errors = append(errors, ValidationError{
				Field:   "revenue_type",
				Message: fmt.Sprintf("revenue_type must be one of: %s", strings.Join(validRevenueTypes, ", ")),
			})
		}
	}

	// Validate payout amount if provided
	if req.PayoutAmount != nil && *req.PayoutAmount < 0 {
		errors = append(errors, ValidationError{
			Field:   "payout_amount",
			Message: "payout_amount cannot be negative",
		})
	}

	// Validate revenue amount if provided
	if req.RevenueAmount != nil && *req.RevenueAmount < 0 {
		errors = append(errors, ValidationError{
			Field:   "revenue_amount",
			Message: "revenue_amount cannot be negative",
		})
	}

	// Validate caps are non-negative if provided
	capFields := map[string]*int{
		"daily_conversion_cap":   req.DailyConversionCap,
		"weekly_conversion_cap":  req.WeeklyConversionCap,
		"monthly_conversion_cap": req.MonthlyConversionCap,
		"global_conversion_cap":  req.GlobalConversionCap,
		"daily_click_cap":        req.DailyClickCap,
		"weekly_click_cap":       req.WeeklyClickCap,
		"monthly_click_cap":      req.MonthlyClickCap,
		"global_click_cap":       req.GlobalClickCap,
	}

	for fieldName, value := range capFields {
		if value != nil && *value < 0 {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("%s cannot be negative", fieldName),
			})
		}
	}

	// Validate date formats if provided
	if req.StartDate != nil && *req.StartDate != "" {
		if _, err := time.Parse("2006-01-02", *req.StartDate); err != nil {
			errors = append(errors, ValidationError{
				Field:   "start_date",
				Message: "start_date must be in YYYY-MM-DD format",
			})
		}
	}

	if req.EndDate != nil && *req.EndDate != "" {
		if _, err := time.Parse("2006-01-02", *req.EndDate); err != nil {
			errors = append(errors, ValidationError{
				Field:   "end_date",
				Message: "end_date must be in YYYY-MM-DD format",
			})
		}
	}

	// Validate date logic if both dates are provided
	if req.StartDate != nil && req.EndDate != nil && *req.StartDate != "" && *req.EndDate != "" {
		startDate, startErr := time.Parse("2006-01-02", *req.StartDate)
		endDate, endErr := time.Parse("2006-01-02", *req.EndDate)
		
		if startErr == nil && endErr == nil && endDate.Before(startDate) {
			errors = append(errors, ValidationError{
				Field:   "end_date",
				Message: "end_date cannot be before start_date",
			})
		}
	}

	return errors
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// validateUpdateCampaignRequest validates the update campaign request
func validateUpdateCampaignRequest(req *UpdateCampaignRequest) []ValidationError {
	var errors []ValidationError

	// Validate required fields
	if strings.TrimSpace(req.Name) == "" {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "name is required and cannot be empty",
		})
	}

	if len(req.Name) > 255 {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "name cannot exceed 255 characters",
		})
	}

	// Validate status (required for update)
	validStatuses := []string{"draft", "active", "paused", "archived"}
	if !contains(validStatuses, req.Status) {
		errors = append(errors, ValidationError{
			Field:   "status",
			Message: fmt.Sprintf("status must be one of: %s", strings.Join(validStatuses, ", ")),
		})
	}

	// Validate visibility if provided
	if req.Visibility != nil && *req.Visibility != "" {
		validVisibilities := []string{"public", "require_approval", "private"}
		if !contains(validVisibilities, *req.Visibility) {
			errors = append(errors, ValidationError{
				Field:   "visibility",
				Message: fmt.Sprintf("visibility must be one of: %s", strings.Join(validVisibilities, ", ")),
			})
		}
	}

	// Validate conversion method if provided
	if req.ConversionMethod != nil && *req.ConversionMethod != "" {
		validMethods := []string{"server_postback", "pixel", "hybrid"}
		if !contains(validMethods, *req.ConversionMethod) {
			errors = append(errors, ValidationError{
				Field:   "conversion_method",
				Message: fmt.Sprintf("conversion_method must be one of: %s", strings.Join(validMethods, ", ")),
			})
		}
	}

	// Validate session definition if provided
	if req.SessionDefinition != nil && *req.SessionDefinition != "" {
		validDefinitions := []string{"cookie", "ip", "fingerprint"}
		if !contains(validDefinitions, *req.SessionDefinition) {
			errors = append(errors, ValidationError{
				Field:   "session_definition",
				Message: fmt.Sprintf("session_definition must be one of: %s", strings.Join(validDefinitions, ", ")),
			})
		}
	}

	// Validate session duration if provided
	if req.SessionDuration != nil && *req.SessionDuration <= 0 {
		errors = append(errors, ValidationError{
			Field:   "session_duration",
			Message: "session_duration must be a positive integer (hours)",
		})
	}

	// Validate payout type if provided
	if req.PayoutType != nil && *req.PayoutType != "" {
		validPayoutTypes := []string{"cpa", "cpc", "cpm", "cps", "cpa_cps", "prv"}
		if !contains(validPayoutTypes, *req.PayoutType) {
			errors = append(errors, ValidationError{
				Field:   "payout_type",
				Message: fmt.Sprintf("payout_type must be one of: %s", strings.Join(validPayoutTypes, ", ")),
			})
		}
	}

	// Validate revenue type if provided
	if req.RevenueType != nil && *req.RevenueType != "" {
		validRevenueTypes := []string{"rpa", "rpc", "rpm", "rps", "rpa_rps", "prv"}
		if !contains(validRevenueTypes, *req.RevenueType) {
			errors = append(errors, ValidationError{
				Field:   "revenue_type",
				Message: fmt.Sprintf("revenue_type must be one of: %s", strings.Join(validRevenueTypes, ", ")),
			})
		}
	}

	// Validate payout amount if provided
	if req.PayoutAmount != nil && *req.PayoutAmount < 0 {
		errors = append(errors, ValidationError{
			Field:   "payout_amount",
			Message: "payout_amount cannot be negative",
		})
	}

	// Validate revenue amount if provided
	if req.RevenueAmount != nil && *req.RevenueAmount < 0 {
		errors = append(errors, ValidationError{
			Field:   "revenue_amount",
			Message: "revenue_amount cannot be negative",
		})
	}

	// Validate caps are non-negative if provided
	capFields := map[string]*int{
		"daily_conversion_cap":   req.DailyConversionCap,
		"weekly_conversion_cap":  req.WeeklyConversionCap,
		"monthly_conversion_cap": req.MonthlyConversionCap,
		"global_conversion_cap":  req.GlobalConversionCap,
		"daily_click_cap":        req.DailyClickCap,
		"weekly_click_cap":       req.WeeklyClickCap,
		"monthly_click_cap":      req.MonthlyClickCap,
		"global_click_cap":       req.GlobalClickCap,
	}

	for fieldName, value := range capFields {
		if value != nil && *value < 0 {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("%s cannot be negative", fieldName),
			})
		}
	}

	// Validate date formats if provided
	if req.StartDate != nil && *req.StartDate != "" {
		if _, err := time.Parse("2006-01-02", *req.StartDate); err != nil {
			errors = append(errors, ValidationError{
				Field:   "start_date",
				Message: "start_date must be in YYYY-MM-DD format",
			})
		}
	}

	if req.EndDate != nil && *req.EndDate != "" {
		if _, err := time.Parse("2006-01-02", *req.EndDate); err != nil {
			errors = append(errors, ValidationError{
				Field:   "end_date",
				Message: "end_date must be in YYYY-MM-DD format",
			})
		}
	}

	// Validate date logic if both dates are provided
	if req.StartDate != nil && req.EndDate != nil && *req.StartDate != "" && *req.EndDate != "" {
		startDate, startErr := time.Parse("2006-01-02", *req.StartDate)
		endDate, endErr := time.Parse("2006-01-02", *req.EndDate)
		
		if startErr == nil && endErr == nil && endDate.Before(startDate) {
			errors = append(errors, ValidationError{
				Field:   "end_date",
				Message: "end_date cannot be before start_date",
			})
		}
	}

	return errors
}

// CreateCampaign creates a new campaign
// @Summary      Create a new campaign
// @Description  Creates a new campaign with the given details
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        request  body      CreateCampaignRequest  true  "Campaign details"
// @Success      201      {object}  domain.Campaign        "Created campaign"
// @Failure      400      {object}  map[string]string      "Invalid request"
// @Failure      500      {object}  map[string]string      "Internal server error"
// @Security     BearerAuth
// @Router       /campaigns [post]
func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var req CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Validate the request
	if validationErrors := validateCreateCampaignRequest(&req); len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"validation_errors": validationErrors,
		})
		return
	}

	campaign := &domain.Campaign{
		OrganizationID: req.OrganizationID,
		AdvertiserID:   req.AdvertiserID,
		Name:           req.Name,
		Description:    req.Description,
		Status:         req.Status,
		
		// Offer-specific fields
		DestinationURL:            req.DestinationURL,
		ThumbnailURL:              req.ThumbnailURL,
		PreviewURL:                req.PreviewURL,
		Visibility:                req.Visibility,
		CurrencyID:                req.CurrencyID,
		ConversionMethod:          req.ConversionMethod,
		SessionDefinition:         req.SessionDefinition,
		SessionDuration:           intToInt32Ptr(req.SessionDuration),
		InternalNotes:             req.InternalNotes,
		TermsAndConditions:        req.TermsAndConditions,
		IsForceTermsAndConditions: req.IsForceTermsAndConditions,
		
		// Caps and limits
		IsCapsEnabled:         req.IsCapsEnabled,
		DailyConversionCap:    req.DailyConversionCap,
		WeeklyConversionCap:   req.WeeklyConversionCap,
		MonthlyConversionCap:  req.MonthlyConversionCap,
		GlobalConversionCap:   req.GlobalConversionCap,
		DailyClickCap:         req.DailyClickCap,
		WeeklyClickCap:        req.WeeklyClickCap,
		MonthlyClickCap:       req.MonthlyClickCap,
		GlobalClickCap:        req.GlobalClickCap,
		
		// Payout and revenue
		PayoutType:    req.PayoutType,
		PayoutAmount:  req.PayoutAmount,
		RevenueType:   req.RevenueType,
		RevenueAmount: req.RevenueAmount,
	}

	// Handle offer config JSON
	if req.OfferConfig != nil {
		offerConfigStr := string(*req.OfferConfig)
		campaign.OfferConfig = &offerConfigStr
	}

	// Parse dates if provided
	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
			return
		}
		campaign.StartDate = &startDate
	}

	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
			return
		}
		campaign.EndDate = &endDate
	}

	createdCampaign, err := h.campaignService.CreateCampaign(c.Request.Context(), campaign)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdCampaign)
}

// GetCampaign retrieves a campaign by ID
// @Summary      Get campaign by ID
// @Description  Retrieves a campaign by its ID
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id   path      int               true  "Campaign ID"
// @Success      200  {object}  domain.Campaign  "Campaign details"
// @Failure      400  {object}  map[string]string "Invalid campaign ID"
// @Failure      404  {object}  map[string]string "Campaign not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Security     BearerAuth
// @Router       /campaigns/{id} [get]
func (h *CampaignHandler) GetCampaign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	campaign, err := h.campaignService.GetCampaignByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "campaign not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get campaign: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaign)
}

// UpdateCampaignRequest defines the request for updating a campaign
type UpdateCampaignRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	Status      string  `json:"status" binding:"required"`
	StartDate   *string `json:"start_date,omitempty"` // Format: YYYY-MM-DD
	EndDate     *string `json:"end_date,omitempty"`   // Format: YYYY-MM-DD
	
	// Offer-specific fields for Everflow integration
	DestinationURL      *string  `json:"destination_url,omitempty"`
	ThumbnailURL        *string  `json:"thumbnail_url,omitempty"`
	PreviewURL          *string  `json:"preview_url,omitempty"`
	Visibility          *string  `json:"visibility,omitempty"`                   // 'public', 'require_approval', 'private'
	CurrencyID          *string  `json:"currency_id,omitempty"`                 // 'USD', 'EUR', etc.
	ConversionMethod    *string  `json:"conversion_method,omitempty"`           // 'server_postback', 'pixel', etc.
	SessionDefinition   *string  `json:"session_definition,omitempty"`          // 'cookie', 'ip', 'fingerprint'
	SessionDuration     *int     `json:"session_duration,omitempty"`            // in hours
	InternalNotes       *string  `json:"internal_notes,omitempty"`
	TermsAndConditions  *string  `json:"terms_and_conditions,omitempty"`
	IsForceTermsAndConditions *bool `json:"is_force_terms_and_conditions,omitempty"`
	
	// Caps and limits
	IsCapsEnabled         *bool `json:"is_caps_enabled,omitempty"`
	DailyConversionCap    *int  `json:"daily_conversion_cap,omitempty"`
	WeeklyConversionCap   *int  `json:"weekly_conversion_cap,omitempty"`
	MonthlyConversionCap  *int  `json:"monthly_conversion_cap,omitempty"`
	GlobalConversionCap   *int  `json:"global_conversion_cap,omitempty"`
	DailyClickCap         *int  `json:"daily_click_cap,omitempty"`
	WeeklyClickCap        *int  `json:"weekly_click_cap,omitempty"`
	MonthlyClickCap       *int  `json:"monthly_click_cap,omitempty"`
	GlobalClickCap        *int  `json:"global_click_cap,omitempty"`
	
	// Payout and revenue configuration
	PayoutType     *string  `json:"payout_type,omitempty"`         // 'cpa', 'cpc', 'cpm', etc.
	PayoutAmount   *float64 `json:"payout_amount,omitempty"`
	RevenueType    *string  `json:"revenue_type,omitempty"`        // 'rpa', 'rpc', 'rpm', etc.
	RevenueAmount  *float64 `json:"revenue_amount,omitempty"`
	
	// Additional configuration
	OfferConfig *json.RawMessage `json:"offer_config,omitempty"` // Additional Everflow-specific config
}

// UpdateCampaign updates a campaign
// @Summary      Update campaign
// @Description  Updates a campaign with the given details
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id       path      int                  true  "Campaign ID"
// @Param        request  body      UpdateCampaignRequest true  "Campaign details"
// @Success      200      {object}  domain.Campaign       "Updated campaign"
// @Failure      400      {object}  map[string]string     "Invalid request"
// @Failure      404      {object}  map[string]string     "Campaign not found"
// @Failure      500      {object}  map[string]string     "Internal server error"
// @Security     BearerAuth
// @Router       /campaigns/{id} [put]
func (h *CampaignHandler) UpdateCampaign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	var req UpdateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Validate the request
	if validationErrors := validateUpdateCampaignRequest(&req); len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"validation_errors": validationErrors,
		})
		return
	}

	// Get existing campaign
	campaign, err := h.campaignService.GetCampaignByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "campaign not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get campaign: " + err.Error()})
		return
	}

	// Update campaign basic fields
	campaign.Name = req.Name
	campaign.Description = req.Description
	campaign.Status = req.Status
	
	// Update offer-specific fields
	campaign.DestinationURL = req.DestinationURL
	campaign.ThumbnailURL = req.ThumbnailURL
	campaign.PreviewURL = req.PreviewURL
	campaign.Visibility = req.Visibility
	campaign.CurrencyID = req.CurrencyID
	campaign.ConversionMethod = req.ConversionMethod
	campaign.SessionDefinition = req.SessionDefinition
	campaign.SessionDuration = intToInt32Ptr(req.SessionDuration)
	campaign.InternalNotes = req.InternalNotes
	campaign.TermsAndConditions = req.TermsAndConditions
	campaign.IsForceTermsAndConditions = req.IsForceTermsAndConditions
	
	// Update caps and limits
	campaign.IsCapsEnabled = req.IsCapsEnabled
	campaign.DailyConversionCap = req.DailyConversionCap
	campaign.WeeklyConversionCap = req.WeeklyConversionCap
	campaign.MonthlyConversionCap = req.MonthlyConversionCap
	campaign.GlobalConversionCap = req.GlobalConversionCap
	campaign.DailyClickCap = req.DailyClickCap
	campaign.WeeklyClickCap = req.WeeklyClickCap
	campaign.MonthlyClickCap = req.MonthlyClickCap
	campaign.GlobalClickCap = req.GlobalClickCap
	
	// Update payout and revenue
	campaign.PayoutType = req.PayoutType
	campaign.PayoutAmount = req.PayoutAmount
	campaign.RevenueType = req.RevenueType
	campaign.RevenueAmount = req.RevenueAmount
	
	// Handle offer config JSON
	if req.OfferConfig != nil {
		offerConfigStr := string(*req.OfferConfig)
		campaign.OfferConfig = &offerConfigStr
	}

	// Parse dates if provided
	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
			return
		}
		campaign.StartDate = &startDate
	} else {
		campaign.StartDate = nil
	}

	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
			return
		}
		campaign.EndDate = &endDate
	} else {
		campaign.EndDate = nil
	}

	if err := h.campaignService.UpdateCampaign(c.Request.Context(), campaign); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campaign: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaign)
}

// ListCampaignsByOrganization retrieves a list of campaigns for an organization
// @Summary      List campaigns by organization
// @Description  Retrieves a list of campaigns for an organization with pagination
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id  path      int                 true   "Organization ID"
// @Param        page           query     int                 false  "Page number (default: 1)"
// @Param        pageSize       query     int                 false  "Page size (default: 10)"
// @Success      200            {array}   domain.Campaign     "List of campaigns"
// @Failure      400            {object}  map[string]string   "Invalid organization ID"
// @Failure      500            {object}  map[string]string   "Internal server error"
// @Security     BearerAuth
// @Router       /organizations/{id}/campaigns [get]
func (h *CampaignHandler) ListCampaignsByOrganization(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
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

	campaigns, err := h.campaignService.ListCampaignsByOrganization(c.Request.Context(), orgID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list campaigns: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}

// ListCampaignsByAdvertiser retrieves a list of campaigns for an advertiser
// @Summary      List campaigns by advertiser
// @Description  Retrieves a list of campaigns for an advertiser with pagination
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id  path      int                 true   "Advertiser ID"
// @Param        page         query     int                 false  "Page number (default: 1)"
// @Param        pageSize     query     int                 false  "Page size (default: 10)"
// @Success      200          {array}   domain.Campaign     "List of campaigns"
// @Failure      400          {object}  map[string]string   "Invalid advertiser ID"
// @Failure      500          {object}  map[string]string   "Internal server error"
// @Security     BearerAuth
// @Router       /advertisers/{id}/campaigns [get]
func (h *CampaignHandler) ListCampaignsByAdvertiser(c *gin.Context) {
	advertiserIDStr := c.Param("id")
	advertiserID, err := strconv.ParseInt(advertiserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
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

	campaigns, err := h.campaignService.ListCampaignsByAdvertiser(c.Request.Context(), advertiserID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list campaigns: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}

// DeleteCampaign deletes a campaign
// @Summary      Delete campaign
// @Description  Deletes a campaign by its ID
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id   path      int                true  "Campaign ID"
// @Success      204  {object}  nil                "No content"
// @Failure      400  {object}  map[string]string  "Invalid campaign ID"
// @Failure      404  {object}  map[string]string  "Campaign not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /campaigns/{id} [delete]
func (h *CampaignHandler) DeleteCampaign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	if err := h.campaignService.DeleteCampaign(c.Request.Context(), id); err != nil {
		if err.Error() == "campaign not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete campaign: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
	c.Writer.Write([]byte{})
}

// CreateCampaignProviderOfferRequest defines the request for creating a campaign provider offer
// swagger:model
type CreateCampaignProviderOfferRequest struct {
	// Campaign ID
	CampaignID         int64  `json:"campaign_id" binding:"required" example:"1"`
	// Provider type (e.g., 'everflow')
	ProviderType       string `json:"provider_type" binding:"required" example:"everflow"`
	// Provider's offer reference
	ProviderOfferRef   *string `json:"provider_offer_ref,omitempty" example:"offer-12345"`
	// Provider offer configuration in JSON format
	// swagger:strfmt json
	ProviderOfferConfig *json.RawMessage `json:"provider_offer_config,omitempty" swaggertype:"object"`
	// Whether the offer is active on the provider
	IsActiveOnProvider bool `json:"is_active_on_provider" example:"true"`
}

// CreateCampaignProviderOffer creates a new campaign provider offer
// @Summary      Create a new campaign provider offer
// @Description  Creates a new offer for a campaign on a provider
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        request  body      CreateCampaignProviderOfferRequest  true  "Offer details"
// @Success      201      {object}  domain.CampaignProviderOffer        "Created offer"
// @Failure      400      {object}  map[string]string                   "Invalid request"
// @Failure      500      {object}  map[string]string                   "Internal server error"
// @Security     BearerAuth
// @Router       /campaign-provider-offers [post]
func (h *CampaignHandler) CreateCampaignProviderOffer(c *gin.Context) {
	var req CreateCampaignProviderOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	offer := &domain.CampaignProviderOffer{
		CampaignID:        req.CampaignID,
		ProviderType:      req.ProviderType,
		ProviderOfferRef:  req.ProviderOfferRef,
		IsActiveOnProvider: req.IsActiveOnProvider,
	}

	if req.ProviderOfferConfig != nil {
		providerOfferConfigStr := string(*req.ProviderOfferConfig)
		offer.ProviderOfferConfig = &providerOfferConfigStr
	}

	createdOffer, err := h.campaignService.CreateCampaignProviderOffer(c.Request.Context(), offer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign provider offer: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdOffer)
}

// GetCampaignProviderOffer retrieves a campaign provider offer by ID
// @Summary      Get campaign provider offer by ID
// @Description  Retrieves a campaign provider offer by its ID
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id   path      int                          true  "Offer ID"
// @Success      200  {object}  domain.CampaignProviderOffer "Offer details"
// @Failure      400  {object}  map[string]string            "Invalid offer ID"
// @Failure      404  {object}  map[string]string            "Offer not found"
// @Failure      500  {object}  map[string]string            "Internal server error"
// @Security     BearerAuth
// @Router       /campaign-provider-offers/{id} [get]
func (h *CampaignHandler) GetCampaignProviderOffer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offer ID"})
		return
	}

	offer, err := h.campaignService.GetCampaignProviderOfferByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "campaign provider offer not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign provider offer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get campaign provider offer: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, offer)
}

// UpdateCampaignProviderOfferRequest defines the request for updating a campaign provider offer
// swagger:model
type UpdateCampaignProviderOfferRequest struct {
	// Provider's offer reference
	ProviderOfferRef   *string `json:"provider_offer_ref,omitempty" example:"offer-12345"`
	// Provider offer configuration in JSON format
	// swagger:strfmt json
	ProviderOfferConfig *json.RawMessage `json:"provider_offer_config,omitempty" swaggertype:"object"`
	// Whether the offer is active on the provider
	IsActiveOnProvider bool `json:"is_active_on_provider" example:"true"`
}

// UpdateCampaignProviderOffer updates a campaign provider offer
// @Summary      Update campaign provider offer
// @Description  Updates a campaign provider offer with the given details
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id       path      int                                true  "Offer ID"
// @Param        request  body      UpdateCampaignProviderOfferRequest true  "Offer details"
// @Success      200      {object}  domain.CampaignProviderOffer       "Updated offer"
// @Failure      400      {object}  map[string]string                  "Invalid request"
// @Failure      404      {object}  map[string]string                  "Offer not found"
// @Failure      500      {object}  map[string]string                  "Internal server error"
// @Security     BearerAuth
// @Router       /campaign-provider-offers/{id} [put]
func (h *CampaignHandler) UpdateCampaignProviderOffer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offer ID"})
		return
	}

	var req UpdateCampaignProviderOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get existing offer
	offer, err := h.campaignService.GetCampaignProviderOfferByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "campaign provider offer not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign provider offer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get campaign provider offer: " + err.Error()})
		return
	}

	// Update offer
	offer.ProviderOfferRef = req.ProviderOfferRef
	offer.IsActiveOnProvider = req.IsActiveOnProvider

	if req.ProviderOfferConfig != nil {
		providerOfferConfigStr := string(*req.ProviderOfferConfig)
		offer.ProviderOfferConfig = &providerOfferConfigStr
	}

	if err := h.campaignService.UpdateCampaignProviderOffer(c.Request.Context(), offer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campaign provider offer: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, offer)
}

// ListCampaignProviderOffersByCampaign retrieves a list of campaign provider offers for a campaign
// @Summary      List campaign provider offers by campaign
// @Description  Retrieves a list of provider offers for a campaign
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id  path      int                            true  "Campaign ID"
// @Success      200         {array}   domain.CampaignProviderOffer  "List of offers"
// @Failure      400         {object}  map[string]string             "Invalid campaign ID"
// @Failure      500         {object}  map[string]string             "Internal server error"
// @Security     BearerAuth
// @Router       /campaigns/{id}/provider-offers [get]
func (h *CampaignHandler) ListCampaignProviderOffersByCampaign(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := strconv.ParseInt(campaignIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	offers, err := h.campaignService.ListCampaignProviderOffersByCampaign(c.Request.Context(), campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list campaign provider offers: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, offers)
}

// DeleteCampaignProviderOffer deletes a campaign provider offer
// @Summary      Delete campaign provider offer
// @Description  Deletes a campaign provider offer by its ID
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        id   path      int                true  "Offer ID"
// @Success      204  {object}  nil                "No content"
// @Failure      400  {object}  map[string]string  "Invalid offer ID"
// @Failure      404  {object}  map[string]string  "Offer not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /campaign-provider-offers/{id} [delete]
func (h *CampaignHandler) DeleteCampaignProviderOffer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offer ID"})
		return
	}

	if err := h.campaignService.DeleteCampaignProviderOffer(c.Request.Context(), id); err != nil {
		if err.Error() == "campaign provider offer not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign provider offer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete campaign provider offer: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
	c.Writer.Write([]byte{})
}