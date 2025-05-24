package integration

import (
	"context"
	"testing"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// SimpleCampaignIntegrationTestSuite tests the DB-first campaign creation flow
type SimpleCampaignIntegrationTestSuite struct {
	suite.Suite
	mockCampaignRepo   *SimpleMockCampaignRepository
	mockAdvertiserRepo *SimpleMockAdvertiserRepository
	mockOrgRepo        *SimpleMockOrganizationRepository
	mockCryptoService  *SimpleMockCryptoService
	service            service.CampaignService
}

// SimpleMockCampaignRepository - minimal implementation
type SimpleMockCampaignRepository struct {
	mock.Mock
}

func (m *SimpleMockCampaignRepository) CreateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	args := m.Called(ctx, campaign)
	// Simulate setting ID and timestamps
	campaign.CampaignID = 123
	campaign.CreatedAt = time.Now()
	campaign.UpdatedAt = time.Now()
	return args.Error(0)
}

func (m *SimpleMockCampaignRepository) GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Campaign), args.Error(1)
}

func (m *SimpleMockCampaignRepository) UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	args := m.Called(ctx, campaign)
	return args.Error(0)
}

func (m *SimpleMockCampaignRepository) DeleteCampaign(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *SimpleMockCampaignRepository) ListCampaignsByOrganization(ctx context.Context, organizationID int64, limit, offset int) ([]*domain.Campaign, error) {
	args := m.Called(ctx, organizationID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Campaign), args.Error(1)
}

func (m *SimpleMockCampaignRepository) ListCampaignsByAdvertiser(ctx context.Context, advertiserID int64, limit, offset int) ([]*domain.Campaign, error) {
	args := m.Called(ctx, advertiserID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Campaign), args.Error(1)
}

func (m *SimpleMockCampaignRepository) CreateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error {
	args := m.Called(ctx, offer)
	return args.Error(0)
}

func (m *SimpleMockCampaignRepository) GetCampaignProviderOfferByID(ctx context.Context, id int64) (*domain.CampaignProviderOffer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CampaignProviderOffer), args.Error(1)
}

func (m *SimpleMockCampaignRepository) UpdateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error {
	args := m.Called(ctx, offer)
	return args.Error(0)
}

func (m *SimpleMockCampaignRepository) DeleteCampaignProviderOffer(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *SimpleMockCampaignRepository) ListCampaignProviderOffersByCampaign(ctx context.Context, campaignID int64) ([]*domain.CampaignProviderOffer, error) {
	args := m.Called(ctx, campaignID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.CampaignProviderOffer), args.Error(1)
}

// SimpleMockAdvertiserRepository
type SimpleMockAdvertiserRepository struct {
	mock.Mock
}

func (m *SimpleMockAdvertiserRepository) CreateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error {
	return nil
}

func (m *SimpleMockAdvertiserRepository) GetAdvertiserByID(ctx context.Context, id int64) (*domain.Advertiser, error) {
	args := m.Called(ctx, id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return &domain.Advertiser{
		AdvertiserID:   id,
		OrganizationID: 1,
		Name:           "Test Advertiser",
		Status:         "active",
	}, nil
}

func (m *SimpleMockAdvertiserRepository) UpdateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error {
	return nil
}

func (m *SimpleMockAdvertiserRepository) DeleteAdvertiser(ctx context.Context, id int64) error {
	return nil
}

func (m *SimpleMockAdvertiserRepository) ListAdvertisersByOrganization(ctx context.Context, organizationID int64, limit, offset int) ([]*domain.Advertiser, error) {
	return nil, nil
}

// SimpleMockOrganizationRepository
type SimpleMockOrganizationRepository struct {
	mock.Mock
}

func (m *SimpleMockOrganizationRepository) CreateOrganization(ctx context.Context, org *domain.Organization) error {
	return nil
}

func (m *SimpleMockOrganizationRepository) GetOrganizationByID(ctx context.Context, id int64) (*domain.Organization, error) {
	args := m.Called(ctx, id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return &domain.Organization{
		OrganizationID: id,
		Name:           "Test Organization",
		Status:         "active",
	}, nil
}

func (m *SimpleMockOrganizationRepository) UpdateOrganization(ctx context.Context, org *domain.Organization) error {
	return nil
}

func (m *SimpleMockOrganizationRepository) DeleteOrganization(ctx context.Context, id int64) error {
	return nil
}

func (m *SimpleMockOrganizationRepository) ListOrganizations(ctx context.Context, limit, offset int) ([]*domain.Organization, error) {
	return nil, nil
}

// SimpleMockCryptoService
type SimpleMockCryptoService struct {
	mock.Mock
}

func (m *SimpleMockCryptoService) Encrypt(plaintext string) (string, error) {
	args := m.Called(plaintext)
	if args.Error(1) != nil {
		return "", args.Error(1)
	}
	return "encrypted_" + plaintext, nil
}

func (m *SimpleMockCryptoService) Decrypt(ciphertext string) (string, error) {
	args := m.Called(ciphertext)
	if args.Error(1) != nil {
		return "", args.Error(1)
	}
	return "decrypted_" + ciphertext, nil
}

func (suite *SimpleCampaignIntegrationTestSuite) SetupTest() {
	suite.mockCampaignRepo = &SimpleMockCampaignRepository{}
	suite.mockAdvertiserRepo = &SimpleMockAdvertiserRepository{}
	suite.mockOrgRepo = &SimpleMockOrganizationRepository{}
	suite.mockCryptoService = &SimpleMockCryptoService{}
	
	// Create service without Everflow integration (nil)
	suite.service = service.NewCampaignService(
		suite.mockCampaignRepo,
		suite.mockAdvertiserRepo,
		suite.mockOrgRepo,
		nil, // No Everflow service for DB-only tests
		suite.mockCryptoService,
	)
}

func (suite *SimpleCampaignIntegrationTestSuite) TestDBOnlyCampaignCreation() {
	ctx := context.Background()

	// Test data - DB-only campaign
	campaign := &domain.Campaign{
		OrganizationID: 1,
		AdvertiserID:   1,
		Name:           "DB Only Campaign",
		Description:    stringPtr("Test campaign for DB-only integration testing"),
		Status:         "draft", // Draft status should not trigger Everflow sync
		
		// Offer fields
		DestinationURL:   stringPtr("https://db-test.com"),
		Visibility:       stringPtr("public"),
		CurrencyID:       stringPtr("USD"),
		PayoutType:       stringPtr("cpa"),
		PayoutAmount:     float64Ptr(1.5),
		RevenueType:      stringPtr("rpa"),
		RevenueAmount:    float64Ptr(2.0),
	}

	// Mock organization lookup
	suite.mockOrgRepo.On("GetOrganizationByID", ctx, int64(1)).Return(nil, nil)
	
	// Mock advertiser lookup
	suite.mockAdvertiserRepo.On("GetAdvertiserByID", ctx, int64(1)).Return(nil, nil)

	// Mock database creation
	suite.mockCampaignRepo.On("CreateCampaign", ctx, mock.AnythingOfType("*domain.Campaign")).Return(nil)

	// Execute
	result, err := suite.service.CreateCampaign(ctx, campaign)

	// Assertions
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(int64(123), result.CampaignID) // Mock sets ID to 123
	suite.Equal("DB Only Campaign", result.Name)
	suite.Equal("draft", result.Status)

	// Verify all mocks were called as expected
	suite.mockCampaignRepo.AssertExpectations(suite.T())
	suite.mockAdvertiserRepo.AssertExpectations(suite.T())
	suite.mockOrgRepo.AssertExpectations(suite.T())
}

func (suite *SimpleCampaignIntegrationTestSuite) TestCampaignRetrievalAfterCreation() {
	ctx := context.Background()

	// First create a campaign
	campaign := &domain.Campaign{
		OrganizationID: 1,
		AdvertiserID:   1,
		Name:           "Retrieval Test Campaign",
		Status:         "active",
		DestinationURL: stringPtr("https://retrieval-test.com"),
	}

	// Mock dependencies for creation
	suite.mockOrgRepo.On("GetOrganizationByID", ctx, int64(1)).Return(nil, nil)
	suite.mockAdvertiserRepo.On("GetAdvertiserByID", ctx, int64(1)).Return(nil, nil)
	suite.mockCampaignRepo.On("CreateCampaign", ctx, mock.AnythingOfType("*domain.Campaign")).Return(nil)

	// Create campaign
	createdCampaign, err := suite.service.CreateCampaign(ctx, campaign)
	suite.NoError(err)
	suite.NotNil(createdCampaign)

	// Mock retrieval - return the created campaign
	suite.mockCampaignRepo.On("GetCampaignByID", ctx, createdCampaign.CampaignID).Return(createdCampaign, nil)

	// Retrieve campaign
	retrievedCampaign, err := suite.service.GetCampaignByID(ctx, createdCampaign.CampaignID)
	suite.NoError(err)
	suite.NotNil(retrievedCampaign)
	suite.Equal(createdCampaign.CampaignID, retrievedCampaign.CampaignID)
	suite.Equal("Retrieval Test Campaign", retrievedCampaign.Name)

	// Verify mocks
	suite.mockCampaignRepo.AssertExpectations(suite.T())
}

func (suite *SimpleCampaignIntegrationTestSuite) TestCampaignWithAllOfferFields() {
	ctx := context.Background()

	// Test data with comprehensive offer fields
	dateLiveUntil, _ := time.Parse("2006-01-02", "2024-12-31")
	campaign := &domain.Campaign{
		OrganizationID: 1,
		AdvertiserID:   1,
		Name:           "Comprehensive Offer Campaign",
		Description:    stringPtr("Campaign with all offer fields"),
		Status:         "active",
		
		// Core offer fields
		DestinationURL:                stringPtr("https://comprehensive-test.com"),
		Visibility:                    stringPtr("require_approval"),
		CurrencyID:                    stringPtr("EUR"),
		PayoutType:                    stringPtr("cps"),
		PayoutAmount:                  float64Ptr(5.0),
		RevenueType:                   stringPtr("rps"),
		RevenueAmount:                 float64Ptr(7.5),
		
		// Additional offer fields
		ThumbnailURL:                  stringPtr("https://example.com/thumb.jpg"),
		InternalNotes:                 stringPtr("Internal test notes"),
		ServerSideURL:                 stringPtr("https://server.example.com/postback"),
		IsViewThroughEnabled:          boolPtr(true),
		ViewThroughDestinationURL:     stringPtr("https://viewthrough.example.com"),
		PreviewURL:                    stringPtr("https://preview.example.com"),
		CapsTimezoneID:                intPtr(5),
		ProjectID:                     stringPtr("PROJECT123"),
		DateLiveUntil:                 &dateLiveUntil,
		HTMLDescription:               stringPtr("<p>HTML description</p>"),
		IsUsingExplicitTermsAndConditions: boolPtr(true),
		TermsAndConditions:            stringPtr("Terms and conditions text"),
		IsForceTermsAndConditions:     boolPtr(false),
		IsCapsEnabled:                 boolPtr(true),
		ConversionMethod:              stringPtr("server_postback"),
		IsWhitelistCheckEnabled:       boolPtr(false),
		SessionDefinition:             stringPtr("cookie"),
		SessionDuration:               intPtr(48),
		AppIdentifier:                 stringPtr("com.example.app"),
		IsDescriptionPlainText:        boolPtr(false),
		IsUseDirectLinking:            boolPtr(true),
		
		// Caps
		DailyConversionCap:            intPtr(100),
		WeeklyConversionCap:           intPtr(500),
		MonthlyConversionCap:          intPtr(2000),
		GlobalConversionCap:           intPtr(10000),
		DailyClickCap:                 intPtr(1000),
		WeeklyClickCap:               intPtr(5000),
		MonthlyClickCap:              intPtr(20000),
		GlobalClickCap:               intPtr(100000),
	}

	// Mock dependencies
	suite.mockOrgRepo.On("GetOrganizationByID", ctx, int64(1)).Return(nil, nil)
	suite.mockAdvertiserRepo.On("GetAdvertiserByID", ctx, int64(1)).Return(nil, nil)
	suite.mockCampaignRepo.On("CreateCampaign", ctx, mock.AnythingOfType("*domain.Campaign")).Return(nil)

	// Execute
	result, err := suite.service.CreateCampaign(ctx, campaign)

	// Assertions
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal("Comprehensive Offer Campaign", result.Name)
	suite.Equal("require_approval", *result.Visibility)
	suite.Equal("EUR", *result.CurrencyID)
	suite.Equal(5.0, *result.PayoutAmount)
	suite.Equal(7.5, *result.RevenueAmount)
	suite.True(*result.IsViewThroughEnabled)
	suite.Equal(100, *result.DailyConversionCap)
	suite.Equal(1000, *result.DailyClickCap)

	// Verify mocks
	suite.mockCampaignRepo.AssertExpectations(suite.T())
}

func TestSimpleCampaignIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(SimpleCampaignIntegrationTestSuite))
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