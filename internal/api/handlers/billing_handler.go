package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// BillingHandler handles billing-related HTTP requests
type BillingHandler struct {
	billingService *service.BillingService
	profileService service.ProfileService
}

// NewBillingHandler creates a new billing handler
func NewBillingHandler(billingService *service.BillingService, profileService service.ProfileService) *BillingHandler {
	return &BillingHandler{
		billingService: billingService,
		profileService: profileService,
	}
}

// validateBillingAccess checks if the user has access to billing resources and returns the target organization ID
func (h *BillingHandler) validateBillingAccess(userProfile *domain.Profile) (int64, error) {
	// Platform owner (Admin) can access any organization's billing
	if userProfile.RoleID == 1 { // Admin (platform owner)
		if userProfile.OrganizationID == nil {
			return 0, fmt.Errorf("platform owner must be associated with an organization")
		}
		return *userProfile.OrganizationID, nil
	}

	// Regular users can only access their own organization's billing
	if userProfile.OrganizationID == nil {
		return 0, fmt.Errorf("user is not associated with an organization")
	}

	return *userProfile.OrganizationID, nil
}

// GetBillingDashboard godoc
// @Summary Get billing dashboard
// @Description Get billing dashboard data for the authenticated user's organization
// @Tags billing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.BillingDashboardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /billing/dashboard [get]
func (h *BillingHandler) GetBillingDashboard(c *gin.Context) {
	// Get user profile from context (set by auth middleware)
	profile, exists := c.Get("profile")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Details: "User profile not found in context",
		})
		return
	}

	userProfile := profile.(*domain.Profile)

	// Validate billing access and get target organization ID
	targetOrganizationID, err := h.validateBillingAccess(userProfile)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "Forbidden",
			Details: err.Error(),
		})
		return
	}

	// Get billing dashboard
	dashboard, err := h.billingService.GetBillingDashboard(c.Request.Context(), targetOrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to get billing dashboard: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

// AddPaymentMethod godoc
// @Summary Add payment method
// @Description Add a new payment method for the organization
// @Tags billing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreatePaymentMethodRequest true "Payment method details"
// @Success 201 {object} domain.StripePaymentMethod
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /billing/payment-methods [post]
func (h *BillingHandler) AddPaymentMethod(c *gin.Context) {
	// Get user profile from context
	profile, exists := c.Get("profile")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Details: "User profile not found in context",
		})
		return
	}

	userProfile := profile.(*domain.Profile)

	// Validate billing access and get target organization ID
	targetOrganizationID, err := h.validateBillingAccess(userProfile)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "Forbidden",
			Details: err.Error(),
		})
		return
	}

	// Parse request
	var req domain.CreatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Details: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Add payment method
	paymentMethod, err := h.billingService.AddPaymentMethod(c.Request.Context(), targetOrganizationID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to add payment method: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, paymentMethod)
}

// RemovePaymentMethod godoc
// @Summary Remove payment method
// @Description Remove a payment method from the organization
// @Tags billing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Payment Method ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /billing/payment-methods/{id} [delete]
func (h *BillingHandler) RemovePaymentMethod(c *gin.Context) {
	// Get user profile from context
	profile, exists := c.Get("profile")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Details: "User profile not found in context",
		})
		return
	}

	userProfile := profile.(*domain.Profile)

	// Validate billing access and get target organization ID
	targetOrganizationID, err := h.validateBillingAccess(userProfile)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "Forbidden",
			Details: err.Error(),
		})
		return
	}

	// Parse payment method ID
	paymentMethodIDStr := c.Param("id")
	paymentMethodID, err := strconv.ParseInt(paymentMethodIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Details: "Invalid payment method ID",
		})
		return
	}

	// Remove payment method
	err = h.billingService.RemovePaymentMethod(c.Request.Context(), targetOrganizationID, paymentMethodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to remove payment method: " + err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// Recharge godoc
// @Summary Recharge account
// @Description Add funds to the organization's account
// @Tags billing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.RechargeRequest true "Recharge details"
// @Success 201 {object} domain.Transaction
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /billing/recharge [post]
func (h *BillingHandler) Recharge(c *gin.Context) {
	// Get user profile from context
	profile, exists := c.Get("profile")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Details: "User profile not found in context",
		})
		return
	}

	userProfile := profile.(*domain.Profile)

	// Validate billing access and get target organization ID
	targetOrganizationID, err := h.validateBillingAccess(userProfile)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "Forbidden",
			Details: err.Error(),
		})
		return
	}

	// Parse request
	var req domain.RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Details: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Recharge account
	transaction, err := h.billingService.Recharge(c.Request.Context(), targetOrganizationID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to recharge account: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// UpdateBillingConfig godoc
// @Summary Update billing configuration
// @Description Update billing configuration for the organization
// @Tags billing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.UpdateBillingConfigRequest true "Billing configuration"
// @Success 200 {object} domain.BillingAccount
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /billing/config [put]
func (h *BillingHandler) UpdateBillingConfig(c *gin.Context) {
	// Get user profile from context
	profile, exists := c.Get("profile")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Details: "User profile not found in context",
		})
		return
	}

	userProfile := profile.(*domain.Profile)

	// Validate billing access and get target organization ID
	targetOrganizationID, err := h.validateBillingAccess(userProfile)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "Forbidden",
			Details: err.Error(),
		})
		return
	}

	// Parse request
	var req domain.UpdateBillingConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Details: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Update billing config
	account, err := h.billingService.UpdateBillingConfig(c.Request.Context(), targetOrganizationID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to update billing configuration: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, account)
}

// GetTransactionHistory godoc
// @Summary Get transaction history
// @Description Get transaction history for the organization
// @Tags billing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} domain.Transaction
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /billing/transactions [get]
func (h *BillingHandler) GetTransactionHistory(c *gin.Context) {
	// Get user profile from context
	profile, exists := c.Get("profile")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Details: "User profile not found in context",
		})
		return
	}

	userProfile := profile.(*domain.Profile)

	// Validate billing access and get target organization ID
	targetOrganizationID, err := h.validateBillingAccess(userProfile)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "Forbidden",
			Details: err.Error(),
		})
		return
	}

	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get transaction history
	transactions, err := h.billingService.GetTransactionHistory(c.Request.Context(), targetOrganizationID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to get transaction history: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
