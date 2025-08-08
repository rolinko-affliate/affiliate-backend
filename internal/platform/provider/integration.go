package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// AdvertiserMappingContext holds additional context information for advertiser mapping
type AdvertiserMappingContext struct {
	Organization *domain.Organization
	UserID       *string
}

// IntegrationService defines the provider-agnostic interface for advertiser, affiliate, campaign, and tracking link operations
type IntegrationService interface {
	// Advertisers
	CreateAdvertiser(ctx context.Context, adv domain.Advertiser) (domain.Advertiser, error)
	CreateAdvertiserWithContext(ctx context.Context, adv domain.Advertiser, mappingCtx *AdvertiserMappingContext) (domain.Advertiser, error)
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

	// Tracking Links
	CreateTrackingLink(ctx context.Context, trackingLink *domain.TrackingLink, campaignMapping *domain.CampaignProviderMapping, affiliateMapping *domain.AffiliateProviderMapping) (*domain.TrackingLinkProviderMapping, error)
	GenerateTrackingLink(ctx context.Context, req *domain.TrackingLinkGenerationRequest, campaignMapping *domain.CampaignProviderMapping, affiliateMapping *domain.AffiliateProviderMapping) (*domain.TrackingLinkGenerationResponse, error)
	GenerateTrackingLinkQR(ctx context.Context, req *domain.TrackingLinkGenerationRequest, campaignMapping *domain.CampaignProviderMapping, affiliateMapping *domain.AffiliateProviderMapping) ([]byte, error)
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

// CreateAdvertiserWithContext mocks advertiser creation with context
func (m *MockIntegrationService) CreateAdvertiserWithContext(ctx context.Context, adv domain.Advertiser, mappingCtx *AdvertiserMappingContext) (domain.Advertiser, error) {
	args := m.Called(ctx, adv, mappingCtx)
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

// CreateTrackingLink mocks tracking link creation
func (m *MockIntegrationService) CreateTrackingLink(ctx context.Context, trackingLink *domain.TrackingLink, campaignMapping *domain.CampaignProviderMapping, affiliateMapping *domain.AffiliateProviderMapping) (*domain.TrackingLinkProviderMapping, error) {
	args := m.Called(ctx, trackingLink, campaignMapping, affiliateMapping)
	return args.Get(0).(*domain.TrackingLinkProviderMapping), args.Error(1)
}

// GenerateTrackingLink mocks tracking link generation
func (m *MockIntegrationService) GenerateTrackingLink(ctx context.Context, req *domain.TrackingLinkGenerationRequest, campaignMapping *domain.CampaignProviderMapping, affiliateMapping *domain.AffiliateProviderMapping) (*domain.TrackingLinkGenerationResponse, error) {
	args := m.Called(ctx, req, campaignMapping, affiliateMapping)
	return args.Get(0).(*domain.TrackingLinkGenerationResponse), args.Error(1)
}

// GenerateTrackingLinkQR mocks tracking link QR code generation
func (m *MockIntegrationService) GenerateTrackingLinkQR(ctx context.Context, req *domain.TrackingLinkGenerationRequest, campaignMapping *domain.CampaignProviderMapping, affiliateMapping *domain.AffiliateProviderMapping) ([]byte, error) {
	args := m.Called(ctx, req, campaignMapping, affiliateMapping)
	return args.Get(0).([]byte), args.Error(1)
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

// CreateAdvertiserWithContext returns the input advertiser with default provider-assigned values and context
func (m *MockIntegrationServiceWithDefaults) CreateAdvertiserWithContext(ctx context.Context, adv domain.Advertiser, mappingCtx *AdvertiserMappingContext) (domain.Advertiser, error) {
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

// CreateTrackingLink returns a simulated tracking link provider mapping
func (m *MockIntegrationServiceWithDefaults) CreateTrackingLink(ctx context.Context, trackingLink *domain.TrackingLink, campaignMapping *domain.CampaignProviderMapping, affiliateMapping *domain.AffiliateProviderMapping) (*domain.TrackingLinkProviderMapping, error) {
	// Create a simulated tracking URL
	generatedURL := "http://mock-tracking-domain.test/ABC123/DEF456/?sub1=test&sub2=mock"

	// Create provider data
	providerData := &domain.EverflowTrackingLinkProviderData{
		NetworkOfferID:     int32Ptr(12345),
		NetworkAffiliateID: int32Ptr(67890),
		GeneratedURL:       &generatedURL,
	}

	providerDataJSON, _ := providerData.ToJSON()

	syncStatus := "synced"
	now := time.Now()
	return &domain.TrackingLinkProviderMapping{
		TrackingLinkID:         trackingLink.TrackingLinkID,
		ProviderType:           "everflow",
		ProviderTrackingLinkID: nil, // Everflow doesn't assign persistent IDs to tracking links
		ProviderData:           &providerDataJSON,
		SyncStatus:             &syncStatus,
		LastSyncAt:             &now,
	}, nil
}

// GenerateTrackingLink returns a simulated tracking link generation response
func (m *MockIntegrationServiceWithDefaults) GenerateTrackingLink(ctx context.Context, req *domain.TrackingLinkGenerationRequest, campaignMapping *domain.CampaignProviderMapping, affiliateMapping *domain.AffiliateProviderMapping) (*domain.TrackingLinkGenerationResponse, error) {
	// Create a simulated tracking URL
	generatedURL := "http://mock-tracking-domain.test/ABC123/DEF456/?sub1=test&sub2=mock"

	// Create provider data
	providerData := &domain.EverflowTrackingLinkProviderData{
		NetworkOfferID:     int32Ptr(12345),
		NetworkAffiliateID: int32Ptr(67890),
		GeneratedURL:       &generatedURL,
	}

	providerDataJSON, _ := providerData.ToJSON()

	return &domain.TrackingLinkGenerationResponse{
		GeneratedURL: generatedURL,
		ProviderData: &providerDataJSON,
	}, nil
}

// GenerateTrackingLinkQR returns a simulated QR code (empty byte slice for mock)
func (m *MockIntegrationServiceWithDefaults) GenerateTrackingLinkQR(ctx context.Context, req *domain.TrackingLinkGenerationRequest, campaignMapping *domain.CampaignProviderMapping, affiliateMapping *domain.AffiliateProviderMapping) ([]byte, error) {
	// Return a mock QR code (in real implementation, this would be a PNG image)
	return []byte("mock-qr-code-data"), nil
}

// LoggingMockIntegrationService is a mock integration service that logs all requests and returns simulated responses
// This is useful for development and testing when you want to see what requests would be made to the real provider
type LoggingMockIntegrationService struct{}

// Ensure LoggingMockIntegrationService implements IntegrationService
var _ IntegrationService = (*LoggingMockIntegrationService)(nil)

// NewLoggingMockIntegrationService creates a new logging mock integration service
func NewLoggingMockIntegrationService() *LoggingMockIntegrationService {
	logger.Info("Mock Integration Service initialized - all provider requests will be logged and simulated")
	return &LoggingMockIntegrationService{}
}

// logRequest logs the request details in a structured format
func (l *LoggingMockIntegrationService) logRequest(operation string, entityType string, data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		logger.Error("Mock request logging failed", "operation", operation, "entity_type", entityType, "error", err)
		return
	}

	logger.Debug("Mock request", "operation", operation, "entity_type", entityType, "data", string(jsonData))
}

// logResponse logs the response details
func (l *LoggingMockIntegrationService) logResponse(operation string, entityType string, data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		logger.Error("Mock response logging failed", "operation", operation, "entity_type", entityType, "error", err)
		return
	}

	logger.Debug("Mock response", "operation", operation, "entity_type", entityType, "data", string(jsonData))
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

// CreateAdvertiserWithContext logs the request with context and returns a simulated advertiser
func (l *LoggingMockIntegrationService) CreateAdvertiserWithContext(ctx context.Context, adv domain.Advertiser, mappingCtx *AdvertiserMappingContext) (domain.Advertiser, error) {
	requestData := map[string]interface{}{
		"advertiser": adv,
		"context":    mappingCtx,
	}
	l.logRequest("CREATE_WITH_CONTEXT", "ADVERTISER", requestData)

	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	// Create a simulated response
	response := adv
	response.AdvertiserID = 12345 // Simulate provider-assigned ID
	response.CreatedAt = time.Now()
	response.UpdatedAt = time.Now()

	l.logResponse("CREATE_WITH_CONTEXT", "ADVERTISER", response)
	return response, nil
}

// UpdateAdvertiser logs the request and simulates an update
func (l *LoggingMockIntegrationService) UpdateAdvertiser(ctx context.Context, adv domain.Advertiser) error {
	l.logRequest("UPDATE", "ADVERTISER", adv)

	// Simulate processing time
	time.Sleep(80 * time.Millisecond)

	logger.Info("Mock update advertiser completed", "advertiser_id", adv.AdvertiserID)
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

	logger.Info("Mock update affiliate completed", "affiliate_id", aff.AffiliateID)
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

	logger.Info("Mock update campaign completed", "campaign_id", camp.CampaignID)
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

// CreateTrackingLink logs the request and returns a simulated tracking link provider mapping
func (l *LoggingMockIntegrationService) CreateTrackingLink(ctx context.Context, trackingLink *domain.TrackingLink, campaignMapping *domain.CampaignProviderMapping, affiliateMapping *domain.AffiliateProviderMapping) (*domain.TrackingLinkProviderMapping, error) {
	l.logRequest("CREATE", "TRACKING_LINK", map[string]interface{}{
		"tracking_link":     trackingLink,
		"campaign_mapping":  campaignMapping,
		"affiliate_mapping": affiliateMapping,
	})

	// Simulate processing time
	time.Sleep(180 * time.Millisecond)

	// Create a simulated tracking URL
	sub1 := ""
	if trackingLink.Sub1 != nil {
		sub1 = *trackingLink.Sub1
	}
	sub2 := ""
	if trackingLink.Sub2 != nil {
		sub2 = *trackingLink.Sub2
	}
	generatedURL := fmt.Sprintf("http://mock-tracking-domain.test/TL%d/C%dA%d/?sub1=%s&sub2=%s", 
		trackingLink.TrackingLinkID, trackingLink.CampaignID, trackingLink.AffiliateID, sub1, sub2)

	// Create provider data
	providerData := &domain.EverflowTrackingLinkProviderData{
		NetworkOfferID:     int32Ptr(int32(trackingLink.CampaignID * 1000)), // Simulate network offer ID
		NetworkAffiliateID: int32Ptr(int32(trackingLink.AffiliateID * 1000)), // Simulate network affiliate ID
		GeneratedURL:       &generatedURL,
	}

	providerDataJSON, _ := providerData.ToJSON()

	// Create response
	syncStatus := "synced"
	now := time.Now()
	response := &domain.TrackingLinkProviderMapping{
		TrackingLinkID:         trackingLink.TrackingLinkID,
		ProviderType:           "everflow",
		ProviderTrackingLinkID: nil, // Everflow doesn't assign persistent IDs to tracking links
		ProviderData:           &providerDataJSON,
		SyncStatus:             &syncStatus,
		LastSyncAt:             &now,
	}

	l.logResponse("CREATE", "TRACKING_LINK", response)
	return response, nil
}

// GenerateTrackingLink logs the request and returns a simulated tracking link
func (l *LoggingMockIntegrationService) GenerateTrackingLink(ctx context.Context, req *domain.TrackingLinkGenerationRequest, campaignMapping *domain.CampaignProviderMapping, affiliateMapping *domain.AffiliateProviderMapping) (*domain.TrackingLinkGenerationResponse, error) {
	l.logRequest("GENERATE", "TRACKING_LINK", map[string]interface{}{
		"request":           req,
		"campaign_mapping":  campaignMapping,
		"affiliate_mapping": affiliateMapping,
	})

	// Simulate processing time
	time.Sleep(200 * time.Millisecond)

	// Create a simulated tracking URL with the request parameters
	baseURL := "http://mock-tracking-domain.test/ABC123/DEF456/"
	params := []string{}

	if req.SourceID != nil {
		params = append(params, fmt.Sprintf("source_id=%s", *req.SourceID))
	}
	if req.Sub1 != nil {
		params = append(params, fmt.Sprintf("sub1=%s", *req.Sub1))
	}
	if req.Sub2 != nil {
		params = append(params, fmt.Sprintf("sub2=%s", *req.Sub2))
	}
	if req.Sub3 != nil {
		params = append(params, fmt.Sprintf("sub3=%s", *req.Sub3))
	}
	if req.Sub4 != nil {
		params = append(params, fmt.Sprintf("sub4=%s", *req.Sub4))
	}
	if req.Sub5 != nil {
		params = append(params, fmt.Sprintf("sub5=%s", *req.Sub5))
	}

	generatedURL := baseURL
	if len(params) > 0 {
		generatedURL += "?" + strings.Join(params, "&")
	}

	// Create provider data
	providerData := &domain.EverflowTrackingLinkProviderData{
		NetworkOfferID:           int32Ptr(12345),
		NetworkCampaignID:        int32Ptr(54321),
		NetworkAffiliateID:       int32Ptr(67890),
		NetworkTrackingDomainID:  req.NetworkTrackingDomainID,
		NetworkOfferURLID:        req.NetworkOfferURLID,
		CreativeID:               req.CreativeID,
		NetworkTrafficSourceID:   req.NetworkTrafficSourceID,
		GeneratedURL:             &generatedURL,
		CanAffiliateRunAllOffers: boolPtr(true),
	}

	providerDataJSON, _ := providerData.ToJSON()

	response := &domain.TrackingLinkGenerationResponse{
		GeneratedURL: generatedURL,
		ProviderData: &providerDataJSON,
	}

	l.logResponse("GENERATE", "TRACKING_LINK", response)
	return response, nil
}

// GenerateTrackingLinkQR logs the request and returns a simulated QR code
func (l *LoggingMockIntegrationService) GenerateTrackingLinkQR(ctx context.Context, req *domain.TrackingLinkGenerationRequest, campaignMapping *domain.CampaignProviderMapping, affiliateMapping *domain.AffiliateProviderMapping) ([]byte, error) {
	l.logRequest("GENERATE", "TRACKING_LINK_QR", map[string]interface{}{
		"request":           req,
		"campaign_mapping":  campaignMapping,
		"affiliate_mapping": affiliateMapping,
	})

	// Simulate processing time
	time.Sleep(250 * time.Millisecond)

	// Return a mock QR code (in real implementation, this would be a PNG image)
	qrData := []byte("mock-qr-code-png-data-for-tracking-link")

	logger.Info("Mock QR code generation completed", "qr_data_size", len(qrData))
	return qrData, nil
}

// Helper functions for pointer creation
func int32Ptr(i int32) *int32 {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
