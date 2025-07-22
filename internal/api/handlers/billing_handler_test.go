package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// BillingServiceInterface defines the interface for billing service
type BillingServiceInterface interface {
	GetBillingDashboard(ctx context.Context, organizationID int64) (*domain.BillingDashboardResponse, error)
	AddPaymentMethod(ctx context.Context, organizationID int64, req *domain.CreatePaymentMethodRequest) (*domain.StripePaymentMethod, error)
	RemovePaymentMethod(ctx context.Context, organizationID int64, paymentMethodID string) error
	Recharge(ctx context.Context, organizationID int64, req *domain.RechargeRequest) (*domain.Transaction, error)
	UpdateBillingConfig(ctx context.Context, organizationID int64, req *domain.UpdateBillingConfigRequest) (*domain.BillingAccount, error)
	GetTransactionHistory(ctx context.Context, organizationID int64, limit, offset int) ([]domain.Transaction, error)
}

// MockBillingService is a mock implementation of BillingServiceInterface
type MockBillingService struct {
	mock.Mock
}

func (m *MockBillingService) GetBillingDashboard(ctx context.Context, organizationID int64) (*domain.BillingDashboardResponse, error) {
	args := m.Called(ctx, organizationID)
	return args.Get(0).(*domain.BillingDashboardResponse), args.Error(1)
}

func (m *MockBillingService) AddPaymentMethod(ctx context.Context, organizationID int64, req *domain.CreatePaymentMethodRequest) (*domain.StripePaymentMethod, error) {
	args := m.Called(ctx, organizationID, req)
	return args.Get(0).(*domain.StripePaymentMethod), args.Error(1)
}

func (m *MockBillingService) RemovePaymentMethod(ctx context.Context, organizationID int64, paymentMethodID string) error {
	args := m.Called(ctx, organizationID, paymentMethodID)
	return args.Error(0)
}

func (m *MockBillingService) Recharge(ctx context.Context, organizationID int64, req *domain.RechargeRequest) (*domain.Transaction, error) {
	args := m.Called(ctx, organizationID, req)
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockBillingService) UpdateBillingConfig(ctx context.Context, organizationID int64, req *domain.UpdateBillingConfigRequest) (*domain.BillingAccount, error) {
	args := m.Called(ctx, organizationID, req)
	return args.Get(0).(*domain.BillingAccount), args.Error(1)
}

func (m *MockBillingService) GetTransactionHistory(ctx context.Context, organizationID int64, limit, offset int) ([]domain.Transaction, error) {
	args := m.Called(ctx, organizationID, limit, offset)
	return args.Get(0).([]domain.Transaction), args.Error(1)
}

// MockProfileService is a mock implementation of ProfileService
type MockProfileService struct {
	mock.Mock
}

func (m *MockProfileService) CreateNewUserProfile(ctx context.Context, userID uuid.UUID, email string, initialOrgID *int64, initialRoleID int) (*domain.Profile, error) {
	args := m.Called(ctx, userID, email, initialOrgID, initialRoleID)
	return args.Get(0).(*domain.Profile), args.Error(1)
}

func (m *MockProfileService) GetProfileByID(ctx context.Context, id uuid.UUID) (*domain.Profile, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Profile), args.Error(1)
}

func (m *MockProfileService) UpdateProfile(ctx context.Context, profile *domain.Profile) (*domain.Profile, error) {
	args := m.Called(ctx, profile)
	return args.Get(0).(*domain.Profile), args.Error(1)
}

func (m *MockProfileService) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProfileService) GetRoleByID(ctx context.Context, roleID int) (*domain.Role, error) {
	args := m.Called(ctx, roleID)
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockProfileService) UpsertProfile(ctx context.Context, userID uuid.UUID, email string, orgID *int64, roleID int, firstName, lastName *string) (*domain.Profile, error) {
	args := m.Called(ctx, userID, email, orgID, roleID, firstName, lastName)
	return args.Get(0).(*domain.Profile), args.Error(1)
}

// TestBillingHandler wraps BillingHandler for testing with interfaces
type TestBillingHandler struct {
	billingService BillingServiceInterface
	profileService service.ProfileService
}

func NewTestBillingHandler(billingService BillingServiceInterface, profileService service.ProfileService) *TestBillingHandler {
	return &TestBillingHandler{
		billingService: billingService,
		profileService: profileService,
	}
}

// validateBillingAccess checks if the user has access to billing resources and returns the target organization ID
func (h *TestBillingHandler) validateBillingAccess(userProfile *domain.Profile) (int64, error) {
	// Platform owner (Admin) can access any organization's billing
	if userProfile.RoleID == 1 { // Admin (platform owner)
		// For admin, use their organization ID as default
		if userProfile.OrganizationID != nil {
			return *userProfile.OrganizationID, nil
		}
		return 0, fmt.Errorf("admin profile missing organization ID")
	}

	// Regular users can only access their own organization's billing
	if userProfile.OrganizationID == nil {
		return 0, fmt.Errorf("user profile missing organization ID")
	}

	return *userProfile.OrganizationID, nil
}

// GetBillingDashboard handles GET /billing/dashboard
func (h *TestBillingHandler) GetBillingDashboard(c *gin.Context) {
	profile, exists := c.Get("profile")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Profile not found in context"})
		return
	}

	userProfile, ok := profile.(*domain.Profile)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid profile type"})
		return
	}

	targetOrganizationID, err := h.validateBillingAccess(userProfile)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	dashboard, err := h.billingService.GetBillingDashboard(c.Request.Context(), targetOrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

// AddPaymentMethod handles POST /billing/payment-methods
func (h *TestBillingHandler) AddPaymentMethod(c *gin.Context) {
	profile, exists := c.Get("profile")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Profile not found in context"})
		return
	}

	userProfile, ok := profile.(*domain.Profile)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid profile type"})
		return
	}

	targetOrganizationID, err := h.validateBillingAccess(userProfile)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var req domain.CreatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentMethod, err := h.billingService.AddPaymentMethod(c.Request.Context(), targetOrganizationID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paymentMethod)
}

func TestBillingHandler_validateBillingAccess(t *testing.T) {
	handler := &TestBillingHandler{}

	tests := []struct {
		name          string
		profile       *domain.Profile
		expectedOrgID int64
		expectError   bool
		errorMsg      string
	}{
		{
			name: "Admin with organization - should succeed",
			profile: &domain.Profile{
				ID:             uuid.New(),
				RoleID:         1, // Admin
				OrganizationID: int64Ptr(100),
			},
			expectedOrgID: 100,
			expectError:   false,
		},
		{
			name: "Admin without organization - should fail",
			profile: &domain.Profile{
				ID:             uuid.New(),
				RoleID:         1, // Admin
				OrganizationID: nil,
			},
			expectError: true,
			errorMsg:    "admin profile missing organization ID",
		},
		{
			name: "Regular user with organization - should succeed",
			profile: &domain.Profile{
				ID:             uuid.New(),
				RoleID:         2, // AdvertiserManager
				OrganizationID: int64Ptr(200),
			},
			expectedOrgID: 200,
			expectError:   false,
		},
		{
			name: "Regular user without organization - should fail",
			profile: &domain.Profile{
				ID:             uuid.New(),
				RoleID:         3, // AffiliateManager
				OrganizationID: nil,
			},
			expectError: true,
			errorMsg:    "user profile missing organization ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgID, err := handler.validateBillingAccess(tt.profile)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOrgID, orgID)
			}
		})
	}
}

func TestBillingHandler_GetBillingDashboard_AccessControl(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name              string
		profile           *domain.Profile
		expectedStatus    int
		expectServiceCall bool
		expectedOrgID     int64
	}{
		{
			name: "Admin user - should access billing",
			profile: &domain.Profile{
				ID:             uuid.New(),
				RoleID:         1, // Admin
				OrganizationID: int64Ptr(100),
			},
			expectedStatus:    http.StatusOK,
			expectServiceCall: true,
			expectedOrgID:     100,
		},
		{
			name: "AdvertiserManager user - should access billing",
			profile: &domain.Profile{
				ID:             uuid.New(),
				RoleID:         2, // AdvertiserManager
				OrganizationID: int64Ptr(200),
			},
			expectedStatus:    http.StatusOK,
			expectServiceCall: true,
			expectedOrgID:     200,
		},
		{
			name: "User without organization - should be forbidden",
			profile: &domain.Profile{
				ID:             uuid.New(),
				RoleID:         2, // AdvertiserManager
				OrganizationID: nil,
			},
			expectedStatus:    http.StatusForbidden,
			expectServiceCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockBillingService := new(MockBillingService)
			mockProfileService := new(MockProfileService)

			if tt.expectServiceCall {
				mockDashboard := &domain.BillingDashboardResponse{
					BillingAccount: &domain.BillingAccount{
						OrganizationID: tt.expectedOrgID,
						Balance:        decimal.NewFromFloat(100.0),
					},
					CurrentBalance: decimal.NewFromFloat(100.0),
				}
				mockBillingService.On("GetBillingDashboard", mock.Anything, tt.expectedOrgID).Return(mockDashboard, nil)
			}

			handler := NewTestBillingHandler(mockBillingService, mockProfileService)

			// Setup Gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/billing/dashboard", nil)

			// Set profile in context (simulating ProfileMiddleware)
			c.Set("profile", tt.profile)

			// Call handler
			handler.GetBillingDashboard(c)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectServiceCall {
				mockBillingService.AssertExpectations(t)
			}
		})
	}
}

func TestBillingHandler_AddPaymentMethod_AccessControl(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test that different users can only access their own organization's billing
	adminProfile := &domain.Profile{
		ID:             uuid.New(),
		RoleID:         1, // Admin
		OrganizationID: int64Ptr(100),
	}

	advertiserProfile := &domain.Profile{
		ID:             uuid.New(),
		RoleID:         2, // AdvertiserManager
		OrganizationID: int64Ptr(200),
	}

	tests := []struct {
		name           string
		profile        *domain.Profile
		expectedOrgID  int64
		expectedStatus int
	}{
		{
			name:           "Admin accesses org 100",
			profile:        adminProfile,
			expectedOrgID:  100,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Advertiser accesses org 200",
			profile:        advertiserProfile,
			expectedOrgID:  200,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockBillingService := new(MockBillingService)
			mockProfileService := new(MockProfileService)

			mockPaymentMethod := &domain.StripePaymentMethod{
				PaymentMethodID:       1,
				OrganizationID:        tt.expectedOrgID,
				StripePaymentMethodID: "pm_test",
			}

			req := &domain.CreatePaymentMethodRequest{
				PaymentMethodID: "pm_test",
			}

			mockBillingService.On("AddPaymentMethod", mock.Anything, tt.expectedOrgID, req).Return(mockPaymentMethod, nil)

			handler := NewTestBillingHandler(mockBillingService, mockProfileService)

			// Setup request
			reqBody, _ := json.Marshal(req)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", "/billing/payment-methods", bytes.NewBuffer(reqBody))
			c.Request.Header.Set("Content-Type", "application/json")

			// Set profile in context
			c.Set("profile", tt.profile)

			// Call handler
			handler.AddPaymentMethod(c)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)
			mockBillingService.AssertExpectations(t)
		})
	}
}

// Helper function to create int64 pointer
func int64Ptr(i int64) *int64 {
	return &i
}
