package service

import (
	"context"
	"testing"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Basic mock implementations for testing
type BasicMockCampaignRepository struct {
	mock.Mock
}

func (m *BasicMockCampaignRepository) CreateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	args := m.Called(ctx, campaign)
	// Simulate setting ID and timestamps
	campaign.CampaignID = 123
	campaign.CreatedAt = time.Now()
	campaign.UpdatedAt = time.Now()
	return args.Error(0)
}

func (m *BasicMockCampaignRepository) GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Campaign), args.Error(1)
}

func (m *BasicMockCampaignRepository) UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	args := m.Called(ctx, campaign)
	return args.Error(0)
}

func (m *BasicMockCampaignRepository) DeleteCampaign(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *BasicMockCampaignRepository) ListCampaignsByOrganization(ctx context.Context, organizationID int64, limit, offset int) ([]*domain.Campaign, error) {
	args := m.Called(ctx, organizationID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Campaign), args.Error(1)
}

func (m *BasicMockCampaignRepository) ListCampaignsByAdvertiser(ctx context.Context, advertiserID int64, limit, offset int) ([]*domain.Campaign, error) {
	args := m.Called(ctx, advertiserID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Campaign), args.Error(1)
}

func (m *BasicMockCampaignRepository) CreateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error {
	args := m.Called(ctx, offer)
	return args.Error(0)
}

func (m *BasicMockCampaignRepository) GetCampaignProviderOfferByID(ctx context.Context, id int64) (*domain.CampaignProviderOffer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CampaignProviderOffer), args.Error(1)
}

func (m *BasicMockCampaignRepository) UpdateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error {
	args := m.Called(ctx, offer)
	return args.Error(0)
}

func (m *BasicMockCampaignRepository) DeleteCampaignProviderOffer(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *BasicMockCampaignRepository) ListCampaignProviderOffersByCampaign(ctx context.Context, campaignID int64) ([]*domain.CampaignProviderOffer, error) {
	args := m.Called(ctx, campaignID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.CampaignProviderOffer), args.Error(1)
}

type BasicMockAdvertiserRepository struct {
	mock.Mock
}

func (m *BasicMockAdvertiserRepository) CreateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error {
	return nil
}

func (m *BasicMockAdvertiserRepository) GetAdvertiserByID(ctx context.Context, id int64) (*domain.Advertiser, error) {
	return &domain.Advertiser{
		AdvertiserID:   id,
		OrganizationID: 1, // Match the organization ID in the test
		Name:           "Test Advertiser",
	}, nil
}

func (m *BasicMockAdvertiserRepository) UpdateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error {
	return nil
}

func (m *BasicMockAdvertiserRepository) DeleteAdvertiser(ctx context.Context, id int64) error {
	return nil
}

func (m *BasicMockAdvertiserRepository) ListAdvertisersByOrganization(ctx context.Context, organizationID int64, limit, offset int) ([]*domain.Advertiser, error) {
	return []*domain.Advertiser{}, nil
}

type BasicMockOrganizationRepository struct {
	mock.Mock
}

func (m *BasicMockOrganizationRepository) CreateOrganization(ctx context.Context, org *domain.Organization) error {
	return nil
}

func (m *BasicMockOrganizationRepository) GetOrganizationByID(ctx context.Context, id int64) (*domain.Organization, error) {
	return &domain.Organization{OrganizationID: id}, nil
}

func (m *BasicMockOrganizationRepository) UpdateOrganization(ctx context.Context, org *domain.Organization) error {
	return nil
}

func (m *BasicMockOrganizationRepository) DeleteOrganization(ctx context.Context, id int64) error {
	return nil
}

func (m *BasicMockOrganizationRepository) ListOrganizations(ctx context.Context, limit, offset int) ([]*domain.Organization, error) {
	return []*domain.Organization{}, nil
}

type BasicMockCryptoService struct{}

func (m *BasicMockCryptoService) Encrypt(plaintext string) (string, error) {
	return plaintext, nil
}

func (m *BasicMockCryptoService) Decrypt(ciphertext string) (string, error) {
	return ciphertext, nil
}

// Mock provider service for testing
type BasicMockProviderCampaignService struct {
	mock.Mock
}

func (m *BasicMockProviderCampaignService) CreateOfferInProvider(ctx context.Context, campaign *domain.Campaign) error {
	args := m.Called(ctx, campaign)
	return args.Error(0)
}

func (m *BasicMockProviderCampaignService) GetOfferFromProvider(ctx context.Context, campaignID int64, relationships []string) (*domain.Campaign, error) {
	args := m.Called(ctx, campaignID, relationships)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Campaign), args.Error(1)
}

func (m *BasicMockProviderCampaignService) UpdateOfferInProvider(ctx context.Context, campaignID int64, campaign *domain.Campaign) (*domain.Campaign, error) {
	args := m.Called(ctx, campaignID, campaign)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Campaign), args.Error(1)
}

func TestCampaignService_CreateCampaign_BasicSuccess(t *testing.T) {
	// Setup
	mockRepo := new(BasicMockCampaignRepository)
	mockAdvertiserRepo := new(BasicMockAdvertiserRepository)
	mockOrgRepo := new(BasicMockOrganizationRepository)
	mockCrypto := &BasicMockCryptoService{}
	
	// Create a mock provider service
	var mockProviderSvc provider.ProviderCampaignService = &BasicMockProviderCampaignService{}
	
	service := NewCampaignService(mockRepo, mockAdvertiserRepo, mockOrgRepo, mockProviderSvc, mockCrypto)

	ctx := context.Background()
	
	// Test data
	campaign := &domain.Campaign{
		OrganizationID: 1,
		AdvertiserID:   1,
		Name:           "Test Campaign",
		Status:         "draft",
	}

	// Mock expectations
	mockRepo.On("CreateCampaign", ctx, mock.AnythingOfType("*domain.Campaign")).Return(nil)

	// Execute
	result, err := service.CreateCampaign(ctx, campaign)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(123), result.CampaignID)
	assert.Equal(t, "Test Campaign", result.Name)
	assert.Equal(t, "draft", result.Status)

	mockRepo.AssertExpectations(t)
}

func TestCampaignService_GetCampaignByID_BasicSuccess(t *testing.T) {
	// Setup
	mockRepo := new(BasicMockCampaignRepository)
	mockAdvertiserRepo := new(BasicMockAdvertiserRepository)
	mockOrgRepo := new(BasicMockOrganizationRepository)
	mockCrypto := &BasicMockCryptoService{}
	var mockProviderSvc provider.ProviderCampaignService = &BasicMockProviderCampaignService{}
	
	service := NewCampaignService(mockRepo, mockAdvertiserRepo, mockOrgRepo, mockProviderSvc, mockCrypto)

	ctx := context.Background()
	campaignID := int64(123)

	expectedCampaign := &domain.Campaign{
		CampaignID:     123,
		OrganizationID: 1,
		AdvertiserID:   1,
		Name:           "Test Campaign",
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Mock expectations
	mockRepo.On("GetCampaignByID", ctx, campaignID).Return(expectedCampaign, nil)

	// Execute
	result, err := service.GetCampaignByID(ctx, campaignID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCampaign.CampaignID, result.CampaignID)
	assert.Equal(t, expectedCampaign.Name, result.Name)

	mockRepo.AssertExpectations(t)
}

func TestCampaignService_GetCampaignByID_NotFound(t *testing.T) {
	// Setup
	mockRepo := new(BasicMockCampaignRepository)
	mockAdvertiserRepo := new(BasicMockAdvertiserRepository)
	mockOrgRepo := new(BasicMockOrganizationRepository)
	mockCrypto := &BasicMockCryptoService{}
	var mockProviderSvc provider.ProviderCampaignService = &BasicMockProviderCampaignService{}
	
	service := NewCampaignService(mockRepo, mockAdvertiserRepo, mockOrgRepo, mockProviderSvc, mockCrypto)

	ctx := context.Background()
	campaignID := int64(999)

	// Mock expectations
	mockRepo.On("GetCampaignByID", ctx, campaignID).Return(nil, domain.ErrNotFound)

	// Execute
	result, err := service.GetCampaignByID(ctx, campaignID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrNotFound, err)

	mockRepo.AssertExpectations(t)
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