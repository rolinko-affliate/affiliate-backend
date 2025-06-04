package provider

import (
	"context"
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
	// Simulate provider-specific fields
	networkAdvID := int32(456)
	result.NetworkAdvertiserID = &networkAdvID
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