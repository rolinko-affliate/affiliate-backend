# Provider Integration Service

This package contains the provider-agnostic integration service interface and its mock implementation for testing.

## Overview

The `IntegrationService` interface defines a common contract for interacting with external providers (like Everflow) for managing advertisers, affiliates, and campaigns. This abstraction allows the application to work with different providers without changing the core business logic.

## Interface

```go
type IntegrationService interface {
    // Advertisers
    CreateAdvertiser(ctx context.Context, adv domain.Advertiser) (domain.Advertiser, error)
    UpdateAdvertiser(ctx context.Context, adv domain.Advertiser) error
    GetAdvertiser(ctx context.Context, id uuid.UUID) (domain.Advertiser, error)

    // Affiliates
    CreateAffiliate(ctx context.Context, aff domain.Affiliate) (domain.Affiliate, error)
    UpdateAffiliate(ctx context.Context, aff domain.Affiliate) error
    GetAffiliate(ctx context.Context, id uuid.UUID) (domain.Affiliate, error)

    // Campaigns
    CreateCampaign(ctx context.Context, camp domain.Campaign) (domain.Campaign, error)
    UpdateCampaign(ctx context.Context, camp domain.Campaign) error
    GetCampaign(ctx context.Context, id uuid.UUID) (domain.Campaign, error)
}
```

## Mock Implementation

### MockIntegrationService

A basic mock implementation using `testify/mock` that allows you to set up specific expectations for your tests.

```go
func TestMyService(t *testing.T) {
    mockService := provider.NewMockIntegrationService()
    ctx := context.Background()
    
    // Setup expectations
    affiliate := domain.Affiliate{Name: "Test Affiliate"}
    expectedAffiliate := domain.Affiliate{AffiliateID: 1, Name: "Test Affiliate"}
    mockService.On("CreateAffiliate", ctx, affiliate).Return(expectedAffiliate, nil)
    
    // Use in your service
    myService := NewMyService(mockService)
    result, err := myService.CreateAffiliate(ctx, affiliate)
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, expectedAffiliate, result)
    mockService.AssertExpectations(t)
}
```

### MockIntegrationServiceWithDefaults

A simple mock implementation with sensible default behaviors pre-configured. This implementation doesn't use testify/mock but provides a straightforward implementation that returns realistic test data. This is useful for tests where you don't need to verify specific integration service calls but need the service to work.

```go
func TestMyServiceWithDefaults(t *testing.T) {
    mockService := provider.NewMockIntegrationServiceWithDefaults()
    ctx := context.Background()
    
    // No need to setup expectations - defaults are built-in
    myService := NewMyService(mockService)
    
    affiliate := domain.Affiliate{Name: "Test Affiliate"}
    result, err := myService.CreateAffiliate(ctx, affiliate)
    
    // The mock will return the input with provider-assigned values
    assert.NoError(t, err)
    assert.Equal(t, int64(1), result.AffiliateID)
    assert.Equal(t, "Test Affiliate", result.Name)
    assert.NotZero(t, result.CreatedAt)
    assert.NotZero(t, result.UpdatedAt)
}
```

## Default Behaviors

The `MockIntegrationServiceWithDefaults` provides the following default behaviors:

### Advertisers
- **CreateAdvertiser**: Returns the input advertiser with ID=1 and current timestamps
- **UpdateAdvertiser**: Returns nil (success)
- **GetAdvertiser**: Returns a test advertiser with ID=1, Name="Test Advertiser", Status="active"

### Affiliates
- **CreateAffiliate**: Returns the input affiliate with ID=1 and current timestamps
- **UpdateAffiliate**: Returns nil (success)
- **GetAffiliate**: Returns a test affiliate with ID=1, Name="Test Affiliate", Status="active"

### Campaigns
- **CreateCampaign**: Returns the input campaign with ID=1, NetworkAdvertiserID=456, and current timestamps
- **UpdateCampaign**: Returns nil (success)
- **GetCampaign**: Returns a test campaign with ID=1, Name="Test Campaign", Status="active"

Note: The `MockIntegrationServiceWithDefaults` is a simple implementation that doesn't use testify/mock, so it doesn't support expectation verification. Use `MockIntegrationService` if you need to verify specific method calls.

## Usage in Service Tests

### Basic Mock Usage

```go
func TestCampaignService_CreateCampaign(t *testing.T) {
    mockIntegration := provider.NewMockIntegrationService()
    mockRepo := &MockCampaignRepository{}
    
    service := NewCampaignService(mockRepo, mockIntegration)
    
    // Setup expectations
    campaign := domain.Campaign{Name: "Test Campaign"}
    expectedCampaign := domain.Campaign{CampaignID: 1, Name: "Test Campaign"}
    mockIntegration.On("CreateCampaign", mock.Anything, campaign).Return(expectedCampaign, nil)
    
    result, err := service.CreateCampaign(context.Background(), &campaign)
    
    assert.NoError(t, err)
    assert.Equal(t, expectedCampaign, *result)
    mockIntegration.AssertExpectations(t)
}
```

### Defaults Mock Usage

```go
func TestCampaignService_CreateCampaign_WithDefaults(t *testing.T) {
    mockIntegration := provider.NewMockIntegrationServiceWithDefaults()
    mockRepo := &MockCampaignRepository{}
    
    service := NewCampaignService(mockRepo, mockIntegration)
    
    campaign := domain.Campaign{Name: "Test Campaign"}
    result, err := service.CreateCampaign(context.Background(), &campaign)
    
    // No need to setup expectations - defaults handle it
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, int64(1), result.CampaignID)
}
```

## Error Testing

You can also test error scenarios:

```go
func TestCampaignService_CreateCampaign_Error(t *testing.T) {
    mockIntegration := provider.NewMockIntegrationService()
    
    campaign := domain.Campaign{Name: "Test Campaign"}
    expectedError := errors.New("provider error")
    mockIntegration.On("CreateCampaign", mock.Anything, campaign).Return(domain.Campaign{}, expectedError)
    
    service := NewCampaignService(mockRepo, mockIntegration)
    result, err := service.CreateCampaign(context.Background(), &campaign)
    
    assert.Error(t, err)
    assert.Equal(t, expectedError, err)
    assert.Nil(t, result)
}
```

## Best Practices

1. **Use MockIntegrationService** when you need to verify specific calls to the integration service
2. **Use MockIntegrationServiceWithDefaults** when you just need the integration service to work without caring about the specific calls
3. **Always call AssertExpectations()** when using the basic mock to ensure all expected calls were made
4. **Use mock.Anything** for parameters you don't care about in your test
5. **Use specific matchers** when you need to verify exact parameter values

## Implementation Details

The mock implementations use the `testify/mock` package and implement the `IntegrationService` interface. They provide:

- Full interface compliance verification at compile time
- Flexible expectation setup
- Automatic assertion verification
- Support for both specific and default behaviors
- Comprehensive test coverage examples