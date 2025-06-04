package provider

import (
	"context"
	"testing"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMockIntegrationService_Interface(t *testing.T) {
	// Test that MockIntegrationService implements IntegrationService interface
	var _ IntegrationService = (*MockIntegrationService)(nil)
	
	mockService := NewMockIntegrationService()
	assert.NotNil(t, mockService)
}

func TestMockIntegrationService_CreateAffiliate(t *testing.T) {
	mockService := NewMockIntegrationService()
	ctx := context.Background()
	
	inputAffiliate := domain.Affiliate{
		Name:   "Test Affiliate",
		Status: "active",
	}
	
	expectedAffiliate := domain.Affiliate{
		AffiliateID:    1,
		Name:           "Test Affiliate",
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	
	// Setup mock expectation
	mockService.On("CreateAffiliate", ctx, inputAffiliate).Return(expectedAffiliate, nil)
	
	// Execute
	result, err := mockService.CreateAffiliate(ctx, inputAffiliate)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedAffiliate.AffiliateID, result.AffiliateID)
	assert.Equal(t, expectedAffiliate.Name, result.Name)
	assert.Equal(t, expectedAffiliate.Status, result.Status)
	
	mockService.AssertExpectations(t)
}

func TestMockIntegrationService_CreateAdvertiser(t *testing.T) {
	mockService := NewMockIntegrationService()
	ctx := context.Background()
	
	inputAdvertiser := domain.Advertiser{
		Name:   "Test Advertiser",
		Status: "active",
	}
	
	expectedAdvertiser := domain.Advertiser{
		AdvertiserID:   1,
		Name:           "Test Advertiser",
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	
	// Setup mock expectation
	mockService.On("CreateAdvertiser", ctx, inputAdvertiser).Return(expectedAdvertiser, nil)
	
	// Execute
	result, err := mockService.CreateAdvertiser(ctx, inputAdvertiser)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedAdvertiser.AdvertiserID, result.AdvertiserID)
	assert.Equal(t, expectedAdvertiser.Name, result.Name)
	assert.Equal(t, expectedAdvertiser.Status, result.Status)
	
	mockService.AssertExpectations(t)
}

func TestMockIntegrationService_CreateCampaign(t *testing.T) {
	mockService := NewMockIntegrationService()
	ctx := context.Background()
	
	inputCampaign := domain.Campaign{
		Name:   "Test Campaign",
		Status: "active",
	}
	
	expectedCampaign := domain.Campaign{
		CampaignID:     1,
		Name:           "Test Campaign",
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	
	// Setup mock expectation
	mockService.On("CreateCampaign", ctx, inputCampaign).Return(expectedCampaign, nil)
	
	// Execute
	result, err := mockService.CreateCampaign(ctx, inputCampaign)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedCampaign.CampaignID, result.CampaignID)
	assert.Equal(t, expectedCampaign.Name, result.Name)
	assert.Equal(t, expectedCampaign.Status, result.Status)
	
	mockService.AssertExpectations(t)
}

func TestMockIntegrationServiceWithDefaults_CreateAffiliate(t *testing.T) {
	mockService := NewMockIntegrationServiceWithDefaults()
	ctx := context.Background()
	
	inputAffiliate := domain.Affiliate{
		Name:   "Test Affiliate",
		Status: "active",
	}
	
	// Execute - no need to setup expectations, defaults are configured
	result, err := mockService.CreateAffiliate(ctx, inputAffiliate)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.AffiliateID)
	assert.Equal(t, "Test Affiliate", result.Name)
	assert.Equal(t, "active", result.Status)
	assert.NotZero(t, result.CreatedAt)
	assert.NotZero(t, result.UpdatedAt)
}

func TestMockIntegrationServiceWithDefaults_GetAffiliate(t *testing.T) {
	mockService := NewMockIntegrationServiceWithDefaults()
	ctx := context.Background()
	
	affiliateID := uuid.New()
	
	// Execute - no need to setup expectations, defaults are configured
	result, err := mockService.GetAffiliate(ctx, affiliateID)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.AffiliateID)
	assert.Equal(t, "Test Affiliate", result.Name)
	assert.Equal(t, "active", result.Status)
	assert.NotZero(t, result.CreatedAt)
	assert.NotZero(t, result.UpdatedAt)
}

func TestMockIntegrationService_UpdateOperations(t *testing.T) {
	mockService := NewMockIntegrationService()
	ctx := context.Background()
	
	// Test affiliate update
	affiliate := domain.Affiliate{AffiliateID: 1, Name: "Updated Affiliate"}
	mockService.On("UpdateAffiliate", ctx, affiliate).Return(nil)
	
	err := mockService.UpdateAffiliate(ctx, affiliate)
	assert.NoError(t, err)
	
	// Test advertiser update
	advertiser := domain.Advertiser{AdvertiserID: 1, Name: "Updated Advertiser"}
	mockService.On("UpdateAdvertiser", ctx, advertiser).Return(nil)
	
	err = mockService.UpdateAdvertiser(ctx, advertiser)
	assert.NoError(t, err)
	
	// Test campaign update
	campaign := domain.Campaign{CampaignID: 1, Name: "Updated Campaign"}
	mockService.On("UpdateCampaign", ctx, campaign).Return(nil)
	
	err = mockService.UpdateCampaign(ctx, campaign)
	assert.NoError(t, err)
	
	mockService.AssertExpectations(t)
}

func TestMockIntegrationService_GetOperations(t *testing.T) {
	mockService := NewMockIntegrationService()
	ctx := context.Background()
	
	// Test get affiliate
	affiliateID := uuid.New()
	expectedAffiliate := domain.Affiliate{AffiliateID: 1, Name: "Test Affiliate"}
	mockService.On("GetAffiliate", ctx, affiliateID).Return(expectedAffiliate, nil)
	
	result, err := mockService.GetAffiliate(ctx, affiliateID)
	assert.NoError(t, err)
	assert.Equal(t, expectedAffiliate, result)
	
	// Test get advertiser
	advertiserID := uuid.New()
	expectedAdvertiser := domain.Advertiser{AdvertiserID: 1, Name: "Test Advertiser"}
	mockService.On("GetAdvertiser", ctx, advertiserID).Return(expectedAdvertiser, nil)
	
	advertiserResult, err := mockService.GetAdvertiser(ctx, advertiserID)
	assert.NoError(t, err)
	assert.Equal(t, expectedAdvertiser, advertiserResult)
	
	// Test get campaign
	campaignID := uuid.New()
	expectedCampaign := domain.Campaign{CampaignID: 1, Name: "Test Campaign"}
	mockService.On("GetCampaign", ctx, campaignID).Return(expectedCampaign, nil)
	
	campaignResult, err := mockService.GetCampaign(ctx, campaignID)
	assert.NoError(t, err)
	assert.Equal(t, expectedCampaign, campaignResult)
	
	mockService.AssertExpectations(t)
}

func TestMockIntegrationService_ErrorScenarios(t *testing.T) {
	mockService := NewMockIntegrationService()
	ctx := context.Background()
	
	// Test error scenarios
	affiliate := domain.Affiliate{Name: "Test Affiliate"}
	expectedError := domain.ErrNotFound
	
	mockService.On("CreateAffiliate", ctx, affiliate).Return(domain.Affiliate{}, expectedError)
	
	result, err := mockService.CreateAffiliate(ctx, affiliate)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, domain.Affiliate{}, result)
	
	mockService.AssertExpectations(t)
}

// Example of how to use the mock in a service test
func ExampleMockIntegrationService_Usage() {
	// Create a new mock
	mockService := NewMockIntegrationService()
	
	// Or use the version with defaults for simpler testing
	mockServiceWithDefaults := NewMockIntegrationServiceWithDefaults()
	
	ctx := context.Background()
	
	// Setup specific expectations for the basic mock
	affiliate := domain.Affiliate{Name: "Test Affiliate"}
	expectedAffiliate := domain.Affiliate{AffiliateID: 1, Name: "Test Affiliate"}
	mockService.On("CreateAffiliate", ctx, affiliate).Return(expectedAffiliate, nil)
	
	// Use in your service
	result, err := mockService.CreateAffiliate(ctx, affiliate)
	if err == nil {
		// Handle success
		_ = result
	}
	
	// For the defaults version, no setup needed
	result2, err2 := mockServiceWithDefaults.CreateAffiliate(ctx, affiliate)
	if err2 == nil {
		// Handle success with default behavior
		_ = result2
	}
	
	// Verify expectations were met
	mockService.AssertExpectations(nil) // Pass your *testing.T in real tests
}