package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCampaignService is a mock implementation of CampaignService
type MockCampaignService struct {
	mock.Mock
}

func (m *MockCampaignService) CreateCampaign(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error) {
	args := m.Called(ctx, campaign)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Campaign), args.Error(1)
}

func (m *MockCampaignService) GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Campaign), args.Error(1)
}

func (m *MockCampaignService) UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	args := m.Called(ctx, campaign)
	return args.Error(0)
}

func (m *MockCampaignService) DeleteCampaign(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCampaignService) ListCampaignsByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]*domain.Campaign, error) {
	args := m.Called(ctx, orgID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Campaign), args.Error(1)
}

func (m *MockCampaignService) ListCampaignsByAdvertiser(ctx context.Context, advertiserID int64, page, pageSize int) ([]*domain.Campaign, error) {
	args := m.Called(ctx, advertiserID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Campaign), args.Error(1)
}

func (m *MockCampaignService) CreateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) (*domain.CampaignProviderOffer, error) {
	args := m.Called(ctx, offer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CampaignProviderOffer), args.Error(1)
}

func (m *MockCampaignService) GetCampaignProviderOfferByID(ctx context.Context, id int64) (*domain.CampaignProviderOffer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CampaignProviderOffer), args.Error(1)
}

func (m *MockCampaignService) UpdateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error {
	args := m.Called(ctx, offer)
	return args.Error(0)
}

func (m *MockCampaignService) ListCampaignProviderOffersByCampaign(ctx context.Context, campaignID int64) ([]*domain.CampaignProviderOffer, error) {
	args := m.Called(ctx, campaignID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.CampaignProviderOffer), args.Error(1)
}

func (m *MockCampaignService) DeleteCampaignProviderOffer(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCampaignService) SyncCampaignFromProvider(ctx context.Context, campaignID int64) error {
	args := m.Called(ctx, campaignID)
	return args.Error(0)
}

func (m *MockCampaignService) SyncCampaignToProvider(ctx context.Context, campaignID int64) error {
	args := m.Called(ctx, campaignID)
	return args.Error(0)
}

func TestCampaignHandler_CreateCampaign_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockCampaignService)
	handler := NewCampaignHandler(mockService)

	// Test data
	requestBody := CreateCampaignRequest{
		OrganizationID: 1,
		AdvertiserID:   1,
		Name:           "Test Campaign",
		Description:    stringPtr("Test Description"),
		Status:         "draft",
		
		// Offer fields
		DestinationURL:   stringPtr("https://example.com"),
		Visibility:       stringPtr("public"),
		CurrencyID:       stringPtr("USD"),
		PayoutType:       stringPtr("cpa"),
		PayoutAmount:     float64Ptr(1.0),
		RevenueType:      stringPtr("rpa"),
		RevenueAmount:    float64Ptr(1.5),
	}

	expectedCampaign := &domain.Campaign{
		CampaignID:     123,
		OrganizationID: 1,
		AdvertiserID:   1,
		Name:           "Test Campaign",
		Description:    stringPtr("Test Description"),
		Status:         "draft",
		DestinationURL: stringPtr("https://example.com"),
		Visibility:     stringPtr("public"),
		CurrencyID:     stringPtr("USD"),
		PayoutType:     stringPtr("cpa"),
		PayoutAmount:   float64Ptr(1.0),
		RevenueType:    stringPtr("rpa"),
		RevenueAmount:  float64Ptr(1.5),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Mock expectations
	mockService.On("CreateCampaign", mock.Anything, mock.AnythingOfType("*domain.Campaign")).Return(expectedCampaign, nil)

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/campaigns", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.CreateCampaign(c)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response domain.Campaign
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedCampaign.CampaignID, response.CampaignID)
	assert.Equal(t, expectedCampaign.Name, response.Name)

	mockService.AssertExpectations(t)
}

func TestCampaignHandler_CreateCampaign_ValidationError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockCampaignService)
	handler := NewCampaignHandler(mockService)

	// Test data with validation errors - valid JSON structure but invalid business rules
	requestBody := CreateCampaignRequest{
		OrganizationID: 1,   // Valid
		AdvertiserID:   1,   // Valid
		Name:           "Test Campaign", // Valid
		Status:         "invalid_status", // Invalid enum value
		
		// Invalid offer fields
		PayoutType:    stringPtr("invalid_payout_type"),
		PayoutAmount:  float64Ptr(-1.0), // Invalid - negative
		RevenueType:   stringPtr("invalid_revenue_type"),
		RevenueAmount: float64Ptr(-1.0), // Invalid - negative
	}

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/campaigns", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.CreateCampaign(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Validation failed", response["error"])
	assert.NotEmpty(t, response["validation_errors"])

	// Should not call service if validation fails
	mockService.AssertNotCalled(t, "CreateCampaign")
}

func TestCampaignHandler_GetCampaign_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockCampaignService)
	handler := NewCampaignHandler(mockService)

	campaignID := int64(123)
	expectedCampaign := &domain.Campaign{
		CampaignID:     campaignID,
		OrganizationID: 1,
		AdvertiserID:   1,
		Name:           "Test Campaign",
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Mock expectations
	mockService.On("GetCampaignByID", mock.Anything, campaignID).Return(expectedCampaign, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/campaigns/123", nil)

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "123"}}

	// Execute
	handler.GetCampaign(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response domain.Campaign
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedCampaign.CampaignID, response.CampaignID)
	assert.Equal(t, expectedCampaign.Name, response.Name)

	mockService.AssertExpectations(t)
}

func TestCampaignHandler_GetCampaign_NotFound(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockCampaignService)
	handler := NewCampaignHandler(mockService)

	campaignID := int64(999)

	// Mock expectations - return error that matches handler's expected error message
	notFoundErr := errors.New("campaign not found: not found")
	mockService.On("GetCampaignByID", mock.Anything, campaignID).Return(nil, notFoundErr)

	// Create request
	req, _ := http.NewRequest("GET", "/campaigns/999", nil)

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	// Execute
	handler.GetCampaign(c)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

func TestCampaignHandler_UpdateCampaign_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockCampaignService)
	handler := NewCampaignHandler(mockService)

	campaignID := int64(123)
	requestBody := UpdateCampaignRequest{
		Name:        "Updated Campaign",
		Description: stringPtr("Updated Description"),
		Status:      "active",
	}

	existingCampaign := &domain.Campaign{
		CampaignID:  campaignID,
		Name:        "Old Campaign",
		Description: stringPtr("Old Description"),
		Status:      "draft",
		UpdatedAt:   time.Now(),
	}

	// Mock expectations - first get the existing campaign, then update it
	mockService.On("GetCampaignByID", mock.Anything, campaignID).Return(existingCampaign, nil)
	mockService.On("UpdateCampaign", mock.Anything, mock.AnythingOfType("*domain.Campaign")).Return(nil)

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/campaigns/123", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "123"}}

	// Execute
	handler.UpdateCampaign(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response domain.Campaign
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, campaignID, response.CampaignID)
	assert.Equal(t, "Updated Campaign", response.Name)

	mockService.AssertExpectations(t)
}

func TestCampaignHandler_UpdateCampaign_ValidationError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockCampaignService)
	handler := NewCampaignHandler(mockService)

	requestBody := UpdateCampaignRequest{
		Name:          "Valid Campaign Name", // Valid
		Status:        "invalid_status", // Invalid enum value
		PayoutAmount:  float64Ptr(-1.0), // Invalid - negative
		RevenueAmount: float64Ptr(-1.0), // Invalid - negative
	}

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/campaigns/123", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "123"}}

	// Execute
	handler.UpdateCampaign(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Validation failed", response["error"])
	assert.NotEmpty(t, response["validation_errors"])

	// Should not call service if validation fails
	mockService.AssertNotCalled(t, "UpdateCampaign")
}

func TestCampaignHandler_DeleteCampaign_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockCampaignService)
	handler := NewCampaignHandler(mockService)

	campaignID := int64(123)

	// Mock expectations
	mockService.On("DeleteCampaign", mock.Anything, campaignID).Return(nil)

	// Create request
	req, _ := http.NewRequest("DELETE", "/campaigns/123", nil)

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "123"}}

	// Execute
	handler.DeleteCampaign(c)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)

	mockService.AssertExpectations(t)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}