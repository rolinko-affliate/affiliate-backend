package provider

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// IntegrationService defines the provider-agnostic interface for advertiser, affiliate, and campaign operations
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

// ProviderAdvertiserService defines the interface for advertiser operations
type ProviderAdvertiserService interface {
	CreateAdvertiserInProvider(ctx context.Context, adv domain.Advertiser) (domain.Advertiser, error)
	UpdateAdvertiserInProvider(ctx context.Context, adv domain.Advertiser) error
	GetAdvertiserFromProvider(ctx context.Context, id uuid.UUID) (domain.Advertiser, error)
}

// ProviderCampaignService defines the interface for campaign operations
type ProviderCampaignService interface {
	CreateOfferInProvider(ctx context.Context, camp domain.Campaign) (domain.Campaign, error)
	UpdateOfferInProvider(ctx context.Context, camp domain.Campaign) error
	GetOfferFromProvider(ctx context.Context, id uuid.UUID) (domain.Campaign, error)
}

// MockIntegrationService is a comprehensive mock implementation of IntegrationService for testing
type MockIntegrationService struct {
	mock.Mock
}

// Ensure MockIntegrationService implements IntegrationService
var _ IntegrationService = (*MockIntegrationService)(nil)

// CreateAdvertiser mocks advertiser creation
func (m *MockIntegrationService) CreateAdvertiser(ctx context.Context, adv domain.Advertiser) (domain.Advertiser, error) {
	args := m.Called(ctx, adv)
	return args.Get(0).(domain.Advertiser), args.Error(1)
}

// UpdateAdvertiser mocks advertiser update
func (m *MockIntegrationService) UpdateAdvertiser(ctx context.Context, adv domain.Advertiser) error {
	args := m.Called(ctx, adv)
	return args.Error(0)
}

// GetAdvertiser mocks advertiser retrieval
func (m *MockIntegrationService) GetAdvertiser(ctx context.Context, id uuid.UUID) (domain.Advertiser, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Advertiser), args.Error(1)
}

// CreateAffiliate mocks affiliate creation
func (m *MockIntegrationService) CreateAffiliate(ctx context.Context, aff domain.Affiliate) (domain.Affiliate, error) {
	args := m.Called(ctx, aff)
	return args.Get(0).(domain.Affiliate), args.Error(1)
}

// UpdateAffiliate mocks affiliate update
func (m *MockIntegrationService) UpdateAffiliate(ctx context.Context, aff domain.Affiliate) error {
	args := m.Called(ctx, aff)
	return args.Error(0)
}

// GetAffiliate mocks affiliate retrieval
func (m *MockIntegrationService) GetAffiliate(ctx context.Context, id uuid.UUID) (domain.Affiliate, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Affiliate), args.Error(1)
}

// CreateCampaign mocks campaign creation
func (m *MockIntegrationService) CreateCampaign(ctx context.Context, camp domain.Campaign) (domain.Campaign, error) {
	args := m.Called(ctx, camp)
	return args.Get(0).(domain.Campaign), args.Error(1)
}

// UpdateCampaign mocks campaign update
func (m *MockIntegrationService) UpdateCampaign(ctx context.Context, camp domain.Campaign) error {
	args := m.Called(ctx, camp)
	return args.Error(0)
}

// GetCampaign mocks campaign retrieval
func (m *MockIntegrationService) GetCampaign(ctx context.Context, id uuid.UUID) (domain.Campaign, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Campaign), args.Error(1)
}

// NewMockIntegrationService creates a new mock integration service
func NewMockIntegrationService() *MockIntegrationService {
	return &MockIntegrationService{}
}

// MockIntegrationServiceWithDefaults provides a mock with sensible default behaviors for testing
// This implementation doesn't use testify mock but provides a simple implementation
type MockIntegrationServiceWithDefaults struct{}

// Ensure MockIntegrationServiceWithDefaults implements IntegrationService
var _ IntegrationService = (*MockIntegrationServiceWithDefaults)(nil)

// NewMockIntegrationServiceWithDefaults creates a mock with default behaviors
func NewMockIntegrationServiceWithDefaults() *MockIntegrationServiceWithDefaults {
	return &MockIntegrationServiceWithDefaults{}
}

// CreateAdvertiser returns the input advertiser with default provider-assigned values
func (m *MockIntegrationServiceWithDefaults) CreateAdvertiser(ctx context.Context, adv domain.Advertiser) (domain.Advertiser, error) {
	result := adv
	result.AdvertiserID = 1
	result.CreatedAt = time.Now()
	result.UpdatedAt = time.Now()
	return result, nil
}

// UpdateAdvertiser returns nil (success)
func (m *MockIntegrationServiceWithDefaults) UpdateAdvertiser(ctx context.Context, adv domain.Advertiser) error {
	return nil
}

// GetAdvertiser returns a default test advertiser
func (m *MockIntegrationServiceWithDefaults) GetAdvertiser(ctx context.Context, id uuid.UUID) (domain.Advertiser, error) {
	return domain.Advertiser{
		AdvertiserID:   1,
		OrganizationID: 1,
		Name:           "Test Advertiser",
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

// CreateAffiliate returns the input affiliate with default provider-assigned values
func (m *MockIntegrationServiceWithDefaults) CreateAffiliate(ctx context.Context, aff domain.Affiliate) (domain.Affiliate, error) {
	result := aff
	result.AffiliateID = 1
	result.CreatedAt = time.Now()
	result.UpdatedAt = time.Now()
	return result, nil
}

// UpdateAffiliate returns nil (success)
func (m *MockIntegrationServiceWithDefaults) UpdateAffiliate(ctx context.Context, aff domain.Affiliate) error {
	return nil
}

// GetAffiliate returns a default test affiliate
func (m *MockIntegrationServiceWithDefaults) GetAffiliate(ctx context.Context, id uuid.UUID) (domain.Affiliate, error) {
	return domain.Affiliate{
		AffiliateID:    1,
		OrganizationID: 1,
		Name:           "Test Affiliate",
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

// CreateCampaign returns the input campaign with default provider-assigned values
func (m *MockIntegrationServiceWithDefaults) CreateCampaign(ctx context.Context, camp domain.Campaign) (domain.Campaign, error) {
	result := camp
	result.CampaignID = 1
	result.CreatedAt = time.Now()
	result.UpdatedAt = time.Now()

	return result, nil
}

// UpdateCampaign returns nil (success)
func (m *MockIntegrationServiceWithDefaults) UpdateCampaign(ctx context.Context, camp domain.Campaign) error {
	return nil
}

// GetCampaign returns a default test campaign
func (m *MockIntegrationServiceWithDefaults) GetCampaign(ctx context.Context, id uuid.UUID) (domain.Campaign, error) {
	return domain.Campaign{
		CampaignID:     1,
		OrganizationID: 1,
		AdvertiserID:   1,
		Name:           "Test Campaign",
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

// LoggingMockIntegrationService is a mock integration service that logs all requests and returns simulated responses
// This is useful for development and testing when you want to see what requests would be made to the real provider
type LoggingMockIntegrationService struct{}

// Ensure LoggingMockIntegrationService implements IntegrationService
var _ IntegrationService = (*LoggingMockIntegrationService)(nil)

// NewLoggingMockIntegrationService creates a new logging mock integration service
func NewLoggingMockIntegrationService() *LoggingMockIntegrationService {
	log.Println("🔧 Mock Integration Service initialized - all provider requests will be logged and simulated")
	return &LoggingMockIntegrationService{}
}

// logRequest logs the request details in a structured format
func (l *LoggingMockIntegrationService) logRequest(operation string, entityType string, data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("🚀 MOCK %s %s: (failed to marshal data: %v)", operation, entityType, err)
		return
	}
	
	log.Printf("🚀 MOCK %s %s:\n%s", operation, entityType, string(jsonData))
}

// logResponse logs the response details
func (l *LoggingMockIntegrationService) logResponse(operation string, entityType string, data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("✅ MOCK %s %s Response: (failed to marshal data: %v)", operation, entityType, err)
		return
	}
	
	log.Printf("✅ MOCK %s %s Response:\n%s", operation, entityType, string(jsonData))
}

// CreateAdvertiser logs the request and returns a simulated advertiser
func (l *LoggingMockIntegrationService) CreateAdvertiser(ctx context.Context, adv domain.Advertiser) (domain.Advertiser, error) {
	l.logRequest("CREATE", "ADVERTISER", adv)
	
	// Simulate processing time
	time.Sleep(100 * time.Millisecond)
	
	// Create a simulated response
	response := adv
	response.AdvertiserID = 12345 // Simulate provider-assigned ID
	response.CreatedAt = time.Now()
	response.UpdatedAt = time.Now()
	
	// Note: Advertiser domain doesn't have ProviderData field
	// Provider-specific data would be stored in AdvertiserProviderMapping
	
	l.logResponse("CREATE", "ADVERTISER", response)
	return response, nil
}

// UpdateAdvertiser logs the request and simulates an update
func (l *LoggingMockIntegrationService) UpdateAdvertiser(ctx context.Context, adv domain.Advertiser) error {
	l.logRequest("UPDATE", "ADVERTISER", adv)
	
	// Simulate processing time
	time.Sleep(80 * time.Millisecond)
	
	log.Printf("✅ MOCK UPDATE ADVERTISER: Successfully updated advertiser ID %d", adv.AdvertiserID)
	return nil
}

// GetAdvertiser logs the request and returns a simulated advertiser
func (l *LoggingMockIntegrationService) GetAdvertiser(ctx context.Context, id uuid.UUID) (domain.Advertiser, error) {
	l.logRequest("GET", "ADVERTISER", map[string]interface{}{"id": id})
	
	// Simulate processing time
	time.Sleep(50 * time.Millisecond)
	
	// Create a simulated response
	response := domain.Advertiser{
		AdvertiserID:   12345,
		OrganizationID: 1,
		Name:           "Mock Advertiser",
		Status:         "active",
		CreatedAt:      time.Now().Add(-24 * time.Hour),
		UpdatedAt:      time.Now(),
	}
	
	// Note: Advertiser domain doesn't have ProviderData field
	// Provider-specific data would be stored in AdvertiserProviderMapping
	
	l.logResponse("GET", "ADVERTISER", response)
	return response, nil
}

// CreateAffiliate logs the request and returns a simulated affiliate
func (l *LoggingMockIntegrationService) CreateAffiliate(ctx context.Context, aff domain.Affiliate) (domain.Affiliate, error) {
	l.logRequest("CREATE", "AFFILIATE", aff)
	
	// Simulate processing time
	time.Sleep(120 * time.Millisecond)
	
	// Create a simulated response
	response := aff
	response.AffiliateID = 67890 // Simulate provider-assigned ID
	response.CreatedAt = time.Now()
	response.UpdatedAt = time.Now()
	
	// Note: Affiliate domain doesn't have ProviderData field
	// Provider-specific data would be stored in AffiliateProviderMapping
	
	l.logResponse("CREATE", "AFFILIATE", response)
	return response, nil
}

// UpdateAffiliate logs the request and simulates an update
func (l *LoggingMockIntegrationService) UpdateAffiliate(ctx context.Context, aff domain.Affiliate) error {
	l.logRequest("UPDATE", "AFFILIATE", aff)
	
	// Simulate processing time
	time.Sleep(90 * time.Millisecond)
	
	log.Printf("✅ MOCK UPDATE AFFILIATE: Successfully updated affiliate ID %d", aff.AffiliateID)
	return nil
}

// GetAffiliate logs the request and returns a simulated affiliate
func (l *LoggingMockIntegrationService) GetAffiliate(ctx context.Context, id uuid.UUID) (domain.Affiliate, error) {
	l.logRequest("GET", "AFFILIATE", map[string]interface{}{"id": id})
	
	// Simulate processing time
	time.Sleep(60 * time.Millisecond)
	
	// Create a simulated response
	response := domain.Affiliate{
		AffiliateID:    67890,
		OrganizationID: 1,
		Name:           "Mock Affiliate",
		Status:         "active",
		CreatedAt:      time.Now().Add(-48 * time.Hour),
		UpdatedAt:      time.Now(),
	}
	
	// Note: Affiliate domain doesn't have ProviderData field
	// Provider-specific data would be stored in AffiliateProviderMapping
	
	l.logResponse("GET", "AFFILIATE", response)
	return response, nil
}

// CreateCampaign logs the request and returns a simulated campaign
func (l *LoggingMockIntegrationService) CreateCampaign(ctx context.Context, camp domain.Campaign) (domain.Campaign, error) {
	l.logRequest("CREATE", "CAMPAIGN", camp)
	
	// Simulate processing time
	time.Sleep(150 * time.Millisecond)
	
	// Create a simulated response
	response := camp
	response.CampaignID = 54321 // Simulate provider-assigned ID
	response.CreatedAt = time.Now()
	response.UpdatedAt = time.Now()
	
	// Note: Campaign domain doesn't have ProviderData field
	// Provider-specific data would be stored in CampaignProviderMapping
	
	l.logResponse("CREATE", "CAMPAIGN", response)
	return response, nil
}

// UpdateCampaign logs the request and simulates an update
func (l *LoggingMockIntegrationService) UpdateCampaign(ctx context.Context, camp domain.Campaign) error {
	l.logRequest("UPDATE", "CAMPAIGN", camp)
	
	// Simulate processing time
	time.Sleep(100 * time.Millisecond)
	
	log.Printf("✅ MOCK UPDATE CAMPAIGN: Successfully updated campaign ID %d", camp.CampaignID)
	return nil
}

// GetCampaign logs the request and returns a simulated campaign
func (l *LoggingMockIntegrationService) GetCampaign(ctx context.Context, id uuid.UUID) (domain.Campaign, error) {
	l.logRequest("GET", "CAMPAIGN", map[string]interface{}{"id": id})
	
	// Simulate processing time
	time.Sleep(70 * time.Millisecond)
	
	// Create a simulated response
	response := domain.Campaign{
		CampaignID:     54321,
		OrganizationID: 1,
		AdvertiserID:   12345,
		Name:           "Mock Campaign",
		Status:         "active",
		CreatedAt:      time.Now().Add(-72 * time.Hour),
		UpdatedAt:      time.Now(),
	}
	
	// Note: Campaign domain doesn't have ProviderData field
	// Provider-specific data would be stored in CampaignProviderMapping
	
	l.logResponse("GET", "CAMPAIGN", response)
	return response, nil
}