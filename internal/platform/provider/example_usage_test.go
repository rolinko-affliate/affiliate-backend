package provider

import (
	"context"
	"testing"

	"github.com/affiliate-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ExampleService demonstrates how a service might use the IntegrationService
type ExampleService struct {
	integration IntegrationService
}

func NewExampleService(integration IntegrationService) *ExampleService {
	return &ExampleService{
		integration: integration,
	}
}

func (s *ExampleService) ProcessAffiliate(ctx context.Context, affiliate domain.Affiliate) (*domain.Affiliate, error) {
	// Some business logic here...
	
	// Create affiliate in provider
	result, err := s.integration.CreateAffiliate(ctx, affiliate)
	if err != nil {
		return nil, err
	}
	
	// More business logic...
	
	return &result, nil
}

func (s *ExampleService) SyncAffiliate(ctx context.Context, affiliateID uuid.UUID) (*domain.Affiliate, error) {
	// Get affiliate from provider
	result, err := s.integration.GetAffiliate(ctx, affiliateID)
	if err != nil {
		return nil, err
	}
	
	// Update local data...
	
	return &result, nil
}

// Example test using the basic mock
func TestExampleService_ProcessAffiliate_WithBasicMock(t *testing.T) {
	// Setup
	mockIntegration := NewMockIntegrationService()
	service := NewExampleService(mockIntegration)
	ctx := context.Background()
	
	// Test data
	inputAffiliate := domain.Affiliate{
		Name:   "Test Affiliate",
		Status: "pending",
	}
	
	expectedAffiliate := domain.Affiliate{
		AffiliateID:    123,
		Name:           "Test Affiliate",
		Status:         "active", // Provider might change status
		OrganizationID: 1,
	}
	
	// Setup mock expectations
	mockIntegration.On("CreateAffiliate", ctx, inputAffiliate).Return(expectedAffiliate, nil)
	
	// Execute
	result, err := service.ProcessAffiliate(ctx, inputAffiliate)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(123), result.AffiliateID)
	assert.Equal(t, "Test Affiliate", result.Name)
	assert.Equal(t, "active", result.Status)
	
	// Verify all expectations were met
	mockIntegration.AssertExpectations(t)
}

// Example test using the defaults mock
func TestExampleService_ProcessAffiliate_WithDefaultsMock(t *testing.T) {
	// Setup - no need to configure expectations
	mockIntegration := NewMockIntegrationServiceWithDefaults()
	service := NewExampleService(mockIntegration)
	ctx := context.Background()
	
	// Test data
	inputAffiliate := domain.Affiliate{
		Name:   "Test Affiliate",
		Status: "pending",
	}
	
	// Execute
	result, err := service.ProcessAffiliate(ctx, inputAffiliate)
	
	// Assert - defaults will provide sensible values
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.AffiliateID) // Default ID
	assert.Equal(t, "Test Affiliate", result.Name)
	assert.NotZero(t, result.CreatedAt) // Default sets timestamps
	assert.NotZero(t, result.UpdatedAt)
}

// Example test for error scenarios
func TestExampleService_ProcessAffiliate_Error(t *testing.T) {
	// Setup
	mockIntegration := NewMockIntegrationService()
	service := NewExampleService(mockIntegration)
	ctx := context.Background()
	
	// Test data
	inputAffiliate := domain.Affiliate{
		Name:   "Test Affiliate",
		Status: "pending",
	}
	
	// Setup mock to return error
	expectedError := domain.ErrNotFound
	mockIntegration.On("CreateAffiliate", ctx, inputAffiliate).Return(domain.Affiliate{}, expectedError)
	
	// Execute
	result, err := service.ProcessAffiliate(ctx, inputAffiliate)
	
	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
	
	mockIntegration.AssertExpectations(t)
}

// Example test for sync operation
func TestExampleService_SyncAffiliate(t *testing.T) {
	// Setup
	mockIntegration := NewMockIntegrationService()
	service := NewExampleService(mockIntegration)
	ctx := context.Background()
	
	// Test data
	affiliateID := uuid.New()
	expectedAffiliate := domain.Affiliate{
		AffiliateID:    123,
		Name:           "Synced Affiliate",
		Status:         "active",
		OrganizationID: 1,
	}
	
	// Setup mock expectations
	mockIntegration.On("GetAffiliate", ctx, affiliateID).Return(expectedAffiliate, nil)
	
	// Execute
	result, err := service.SyncAffiliate(ctx, affiliateID)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedAffiliate.AffiliateID, result.AffiliateID)
	assert.Equal(t, expectedAffiliate.Name, result.Name)
	
	mockIntegration.AssertExpectations(t)
}

// Example test using mock.MatchedBy for complex parameter matching
func TestExampleService_ProcessAffiliate_WithParameterMatching(t *testing.T) {
	// Setup
	mockIntegration := NewMockIntegrationService()
	service := NewExampleService(mockIntegration)
	ctx := context.Background()
	
	// Test data
	inputAffiliate := domain.Affiliate{
		Name:   "Test Affiliate",
		Status: "pending",
	}
	
	expectedAffiliate := domain.Affiliate{
		AffiliateID: 123,
		Name:        "Test Affiliate",
		Status:      "active",
	}
	
	// Setup mock with parameter matching
	mockIntegration.On("CreateAffiliate", 
		ctx, 
		mock.MatchedBy(func(aff domain.Affiliate) bool {
			return aff.Name == "Test Affiliate" && aff.Status == "pending"
		})).Return(expectedAffiliate, nil)
	
	// Execute
	result, err := service.ProcessAffiliate(ctx, inputAffiliate)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedAffiliate.AffiliateID, result.AffiliateID)
	
	mockIntegration.AssertExpectations(t)
}

// Example test showing how to verify call count
func TestExampleService_ProcessAffiliate_CallCount(t *testing.T) {
	// Setup
	mockIntegration := NewMockIntegrationService()
	service := NewExampleService(mockIntegration)
	ctx := context.Background()
	
	// Test data
	inputAffiliate := domain.Affiliate{Name: "Test Affiliate"}
	expectedAffiliate := domain.Affiliate{AffiliateID: 123, Name: "Test Affiliate"}
	
	// Setup mock to expect exactly one call
	mockIntegration.On("CreateAffiliate", ctx, inputAffiliate).Return(expectedAffiliate, nil).Once()
	
	// Execute
	result, err := service.ProcessAffiliate(ctx, inputAffiliate)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// This will verify that CreateAffiliate was called exactly once
	mockIntegration.AssertExpectations(t)
	
	// You can also check specific method calls
	mockIntegration.AssertCalled(t, "CreateAffiliate", ctx, inputAffiliate)
	mockIntegration.AssertNumberOfCalls(t, "CreateAffiliate", 1)
}