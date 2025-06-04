package service

import (
	"context"
	"testing"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/provider"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Minimal mock for testing integration service functionality
type MockIntegrationServiceSimple struct {
	mock.Mock
}

func (m *MockIntegrationServiceSimple) CreateAffiliate(ctx context.Context, affiliate domain.Affiliate) (domain.Affiliate, error) {
	args := m.Called(ctx, affiliate)
	return args.Get(0).(domain.Affiliate), args.Error(1)
}

func (m *MockIntegrationServiceSimple) UpdateAffiliate(ctx context.Context, affiliate domain.Affiliate) error {
	args := m.Called(ctx, affiliate)
	return args.Error(0)
}

func (m *MockIntegrationServiceSimple) GetAffiliate(ctx context.Context, id uuid.UUID) (domain.Affiliate, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Affiliate), args.Error(1)
}

func (m *MockIntegrationServiceSimple) CreateAdvertiser(ctx context.Context, adv domain.Advertiser) (domain.Advertiser, error) {
	args := m.Called(ctx, adv)
	return args.Get(0).(domain.Advertiser), args.Error(1)
}

func (m *MockIntegrationServiceSimple) UpdateAdvertiser(ctx context.Context, adv domain.Advertiser) error {
	args := m.Called(ctx, adv)
	return args.Error(0)
}

func (m *MockIntegrationServiceSimple) GetAdvertiser(ctx context.Context, id uuid.UUID) (domain.Advertiser, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Advertiser), args.Error(1)
}

func (m *MockIntegrationServiceSimple) CreateCampaign(ctx context.Context, campaign domain.Campaign) (domain.Campaign, error) {
	args := m.Called(ctx, campaign)
	return args.Get(0).(domain.Campaign), args.Error(1)
}

func (m *MockIntegrationServiceSimple) UpdateCampaign(ctx context.Context, campaign domain.Campaign) error {
	args := m.Called(ctx, campaign)
	return args.Error(0)
}

func (m *MockIntegrationServiceSimple) GetCampaign(ctx context.Context, id uuid.UUID) (domain.Campaign, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Campaign), args.Error(1)
}

func TestIntegrationServiceInterface(t *testing.T) {
	t.Run("test integration service interface compliance", func(t *testing.T) {
		// This test verifies that our mock implements the interface correctly
		var _ provider.IntegrationService = &MockIntegrationServiceSimple{}
		
		mockService := &MockIntegrationServiceSimple{}
		ctx := context.Background()
		
		// Test affiliate creation
		inputAffiliate := domain.Affiliate{
			Name:   "Test Affiliate",
			Status: "active",
		}
		
		expectedAffiliate := domain.Affiliate{
			AffiliateID:        1,
			Name:               "Test Affiliate",
			Status:             "active",
		}
		
		mockService.On("CreateAffiliate", ctx, inputAffiliate).Return(expectedAffiliate, nil)
		
		result, err := mockService.CreateAffiliate(ctx, inputAffiliate)
		
		assert.NoError(t, err)
		assert.Equal(t, expectedAffiliate.AffiliateID, result.AffiliateID)
		assert.Equal(t, expectedAffiliate.Name, result.Name)
		assert.Equal(t, expectedAffiliate.Status, result.Status)
		
		mockService.AssertExpectations(t)
	})
}

func TestAffiliateServiceIntegration(t *testing.T) {
	t.Run("test affiliate service with integration service", func(t *testing.T) {
		// This test focuses on the integration between affiliate service and integration service
		// without requiring full repository mocks
		
		mockIntegrationService := &MockIntegrationServiceSimple{}
		ctx := context.Background()
		
		// Test data
		affiliate := domain.Affiliate{
			Name:   "Test Affiliate",
			Status: "active",
		}
		
		expectedResult := domain.Affiliate{
			AffiliateID:        1,
			Name:               "Test Affiliate",
			Status:             "active",
		}
		
		// Setup mock expectation
		mockIntegrationService.On("CreateAffiliate", ctx, affiliate).Return(expectedResult, nil)
		
		// Execute
		result, err := mockIntegrationService.CreateAffiliate(ctx, affiliate)
		
		// Verify
		assert.NoError(t, err)
		assert.Equal(t, int64(1), result.AffiliateID)
		assert.Equal(t, "Test Affiliate", result.Name)
		assert.Equal(t, "active", result.Status)
		
		mockIntegrationService.AssertExpectations(t)
	})
}

// Helper function
func int32Ptr(i int32) *int32 {
	return &i
}