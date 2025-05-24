package everflow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type MockAdvertiserRepository struct {
	mock.Mock
}

type MockOrganizationRepository struct {
	mock.Mock
}

func (m *MockOrganizationRepository) GetOrganizationByID(ctx context.Context, id int64) (*domain.Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) CreateOrganization(ctx context.Context, org *domain.Organization) error {
	args := m.Called(ctx, org)
	return args.Error(0)
}

func (m *MockOrganizationRepository) UpdateOrganization(ctx context.Context, org *domain.Organization) error {
	args := m.Called(ctx, org)
	return args.Error(0)
}

func (m *MockOrganizationRepository) ListOrganizations(ctx context.Context, limit, offset int) ([]*domain.Organization, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) DeleteOrganization(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAdvertiserRepository) CreateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error {
	args := m.Called(ctx, advertiser)
	return args.Error(0)
}

func (m *MockAdvertiserRepository) GetAdvertiserByID(ctx context.Context, id int64) (*domain.Advertiser, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Advertiser), args.Error(1)
}

func (m *MockAdvertiserRepository) UpdateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error {
	args := m.Called(ctx, advertiser)
	return args.Error(0)
}

func (m *MockAdvertiserRepository) ListAdvertisersByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.Advertiser, error) {
	args := m.Called(ctx, orgID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Advertiser), args.Error(1)
}

func (m *MockAdvertiserRepository) DeleteAdvertiser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Mock provider mapping repository
type MockAdvertiserProviderMappingRepository struct {
	mock.Mock
}

func (m *MockAdvertiserProviderMappingRepository) CreateAdvertiserProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error {
	args := m.Called(ctx, mapping)
	return args.Error(0)
}

func (m *MockAdvertiserProviderMappingRepository) GetAdvertiserProviderMapping(ctx context.Context, advertiserID int64, providerType string) (*domain.AdvertiserProviderMapping, error) {
	args := m.Called(ctx, advertiserID, providerType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AdvertiserProviderMapping), args.Error(1)
}

func (m *MockAdvertiserProviderMappingRepository) UpdateAdvertiserProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error {
	args := m.Called(ctx, mapping)
	return args.Error(0)
}

func (m *MockAdvertiserProviderMappingRepository) DeleteAdvertiserProviderMapping(ctx context.Context, mappingID int64) error {
	args := m.Called(ctx, mappingID)
	return args.Error(0)
}

type MockCampaignRepository struct {
	mock.Mock
}

func (m *MockCampaignRepository) CreateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	args := m.Called(ctx, campaign)
	return args.Error(0)
}

func (m *MockCampaignRepository) GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Campaign), args.Error(1)
}

func (m *MockCampaignRepository) UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	args := m.Called(ctx, campaign)
	return args.Error(0)
}

func (m *MockCampaignRepository) ListCampaignsByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.Campaign, error) {
	args := m.Called(ctx, orgID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Campaign), args.Error(1)
}

func (m *MockCampaignRepository) ListCampaignsByAdvertiser(ctx context.Context, advertiserID int64, limit, offset int) ([]*domain.Campaign, error) {
	args := m.Called(ctx, advertiserID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Campaign), args.Error(1)
}

func (m *MockCampaignRepository) DeleteCampaign(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCampaignRepository) CreateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error {
	args := m.Called(ctx, offer)
	return args.Error(0)
}

func (m *MockCampaignRepository) GetCampaignProviderOfferByID(ctx context.Context, id int64) (*domain.CampaignProviderOffer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CampaignProviderOffer), args.Error(1)
}

func (m *MockCampaignRepository) UpdateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error {
	args := m.Called(ctx, offer)
	return args.Error(0)
}

func (m *MockCampaignRepository) ListCampaignProviderOffersByCampaign(ctx context.Context, campaignID int64) ([]*domain.CampaignProviderOffer, error) {
	args := m.Called(ctx, campaignID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.CampaignProviderOffer), args.Error(1)
}

func (m *MockCampaignRepository) DeleteCampaignProviderOffer(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockCryptoService implements crypto.Service for testing
type MockCryptoService struct {
	mock.Mock
}

// Ensure MockCryptoService implements crypto.Service
var _ crypto.Service = (*MockCryptoService)(nil)

func (m *MockCryptoService) Encrypt(plaintext string) (string, error) {
	args := m.Called(plaintext)
	return args.String(0), args.Error(1)
}

func (m *MockCryptoService) Decrypt(ciphertext string) (string, error) {
	args := m.Called(ciphertext)
	return args.String(0), args.Error(1)
}

// Test CreateAdvertiserInEverflow
func TestCreateAdvertiserInEverflow(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		assert.Equal(t, "POST", r.Method)

		if r.URL.Path == "/v1/networks/advertisers" {
			// Return a successful response for advertiser creation
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{
				"network_advertiser_id": 12345,
				"name": "Test Advertiser",
				"account_status": "active",
				"default_currency_id": "USD",
				"time_created": 1621234567,
				"time_saved": 1621234567
			}`)
		} else if r.URL.Path == "/v1/networks/advertisers/12345/tags" {
			// Return a successful response for adding tags
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success": true}`)
		} else {
			// Unexpected path
			t.Fatalf("Unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	// Create test advertiser
	advertiser := &domain.Advertiser{
		AdvertiserID:   123,
		OrganizationID: 456,
		Name:           "Test Advertiser",
		Status:         "active",
	}

	// Setup mocks
	mockAdvertiserRepo := new(MockAdvertiserRepository)
	mockProviderMappingRepo := new(MockAdvertiserProviderMappingRepository)
	mockCampaignRepo := new(MockCampaignRepository)
	mockCryptoService := new(MockCryptoService)

	// Setup expectations
	mockProviderMappingRepo.On("CreateAdvertiserProviderMapping", mock.Anything, mock.MatchedBy(func(mapping *domain.AdvertiserProviderMapping) bool {
		// Verify mapping properties
		assert.Equal(t, advertiser.AdvertiserID, mapping.AdvertiserID)
		assert.Equal(t, "everflow", mapping.ProviderType)
		assert.NotNil(t, mapping.ProviderAdvertiserID)
		assert.Equal(t, "12345", *mapping.ProviderAdvertiserID)
		assert.NotNil(t, mapping.ProviderConfig)

		// Verify provider config contains expected data
		var providerConfig map[string]interface{}
		err := json.Unmarshal([]byte(*mapping.ProviderConfig), &providerConfig)
		assert.NoError(t, err)
		assert.Equal(t, float64(12345), providerConfig["network_advertiser_id"])

		return true
	})).Return(nil)

	// Create service with mock client
	client := NewClient("test-api-key")
	client.httpClient = server.Client()

	service := &Service{
		client:         client,
		advertiserRepo:        mockAdvertiserRepo,
		providerMappingRepo:   mockProviderMappingRepo,
		campaignRepo:   mockCampaignRepo,
		cryptoService:  mockCryptoService,
	}

	// Override the base URL to point to our test server
	origBaseURL := everflowAPIBaseURL
	everflowAPIBaseURL = server.URL + "/v1"
	defer func() { everflowAPIBaseURL = origBaseURL }()

	// Test the service
	err := service.CreateAdvertiserInEverflow(context.Background(), advertiser)
	assert.NoError(t, err)

	// Verify expectations
	mockAdvertiserRepo.AssertExpectations(t)
}

// Test CreateOfferInEverflow
func TestCreateOfferInEverflow(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		assert.Equal(t, "POST", r.Method)

		if r.URL.Path == "/v1/networks/offers" {
			// Return a successful response for offer creation
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{
				"network_offer_id": 54321,
				"network_id": 1,
				"network_advertiser_id": 12345,
				"name": "Test Campaign",
				"destination_url": "https://example.com/campaigns/789?click_id={transaction_id}",
				"offer_status": "active",
				"currency_id": "USD",
				"encoded_value": "ABC123",
				"time_created": 1621234567,
				"time_saved": 1621234567,
				"visibility": "public",
				"conversion_method": "server_postback"
			}`)
		} else if r.URL.Path == "/v1/networks/offers/54321/tags" {
			// Return a successful response for adding tags
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success": true}`)
		} else {
			// Unexpected path
			t.Fatalf("Unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	// Create test campaign
	campaign := &domain.Campaign{
		CampaignID:     789,
		OrganizationID: 456,
		AdvertiserID:   123,
		Name:           "Test Campaign",
		Status:         "active",
	}

	// Create test advertiser
	advertiser := &domain.Advertiser{
		AdvertiserID:   123,
		OrganizationID: 456,
		Name:           "Test Advertiser",
		Status:         "active",
	}

	// Create test mapping
	providerAdvertiserID := "12345"
	mapping := &domain.AdvertiserProviderMapping{
		MappingID:            1,
		AdvertiserID:         123,
		ProviderType:         "everflow",
		ProviderAdvertiserID: &providerAdvertiserID,
	}

	// Setup mocks
	mockAdvertiserRepo := new(MockAdvertiserRepository)
	mockProviderMappingRepo := new(MockAdvertiserProviderMappingRepository)
	mockCampaignRepo := new(MockCampaignRepository)
	mockCryptoService := new(MockCryptoService)

	// Setup expectations
	mockAdvertiserRepo.On("GetAdvertiserByID", mock.Anything, int64(123)).Return(advertiser, nil)
	mockProviderMappingRepo.On("GetAdvertiserProviderMapping", mock.Anything, int64(123), "everflow").Return(mapping, nil)

	mockCampaignRepo.On("CreateCampaignProviderOffer", mock.Anything, mock.MatchedBy(func(offer *domain.CampaignProviderOffer) bool {
		// Verify offer properties
		assert.Equal(t, campaign.CampaignID, offer.CampaignID)
		assert.Equal(t, "everflow", offer.ProviderType)
		assert.NotNil(t, offer.ProviderOfferRef)
		assert.Equal(t, "54321", *offer.ProviderOfferRef)
		assert.NotNil(t, offer.ProviderOfferConfig)
		assert.True(t, offer.IsActiveOnProvider)
		assert.NotNil(t, offer.LastSyncedAt)

		// Verify provider config contains expected data
		var providerConfig map[string]interface{}
		err := json.Unmarshal([]byte(*offer.ProviderOfferConfig), &providerConfig)
		assert.NoError(t, err)
		assert.Equal(t, float64(54321), providerConfig["network_offer_id"])
		assert.Equal(t, float64(12345), providerConfig["network_advertiser_id"])
		assert.Equal(t, "ABC123", providerConfig["encoded_value"])
		assert.Equal(t, float64(1621234567), providerConfig["time_created"])
		assert.Equal(t, float64(1621234567), providerConfig["time_saved"])

		return true
	})).Return(nil)

	// Create service with mock client
	client := NewClient("test-api-key")
	client.httpClient = server.Client()

	service := &Service{
		client:         client,
		advertiserRepo:        mockAdvertiserRepo,
		providerMappingRepo:   mockProviderMappingRepo,
		campaignRepo:   mockCampaignRepo,
		cryptoService:  mockCryptoService,
	}

	// Override the base URL to point to our test server
	origBaseURL := everflowAPIBaseURL
	everflowAPIBaseURL = server.URL + "/v1"
	defer func() { everflowAPIBaseURL = origBaseURL }()

	// Test the service
	err := service.CreateOfferInEverflow(context.Background(), campaign)
	assert.NoError(t, err)

	// Verify expectations
	mockAdvertiserRepo.AssertExpectations(t)
	mockCampaignRepo.AssertExpectations(t)
}

// Test mapAdvertiserToEverflowRequest
func TestMapAdvertiserToEverflowRequest(t *testing.T) {
	// Setup
	mockAdvertiserRepo := new(MockAdvertiserRepository)
	mockProviderMappingRepo := new(MockAdvertiserProviderMappingRepository)
	mockCampaignRepo := new(MockCampaignRepository)
	mockCryptoService := new(MockCryptoService)

	service := &Service{
		client:         NewClient("test-api-key"),
		advertiserRepo:        mockAdvertiserRepo,
		providerMappingRepo:   mockProviderMappingRepo,
		campaignRepo:   mockCampaignRepo,
		cryptoService:  mockCryptoService,
	}

	// Test case 1: Basic advertiser
	t.Run("Basic advertiser", func(t *testing.T) {
		advertiser := &domain.Advertiser{
			AdvertiserID:   123,
			OrganizationID: 456,
			Name:           "Test Advertiser",
			Status:         "active",
		}

		req, err := service.mapAdvertiserToEverflowRequest(advertiser)
		assert.NoError(t, err)
		assert.Equal(t, "Test Advertiser", req.Name)
		assert.Equal(t, "active", req.AccountStatus)
		assert.Equal(t, "USD", req.DefaultCurrencyID)
		assert.Nil(t, req.ContactAddress)
		assert.Nil(t, req.Billing)
	})

	// Test case 2: Advertiser with contact email
	t.Run("Advertiser with contact email", func(t *testing.T) {
		email := "test@example.com"
		advertiser := &domain.Advertiser{
			AdvertiserID:   123,
			OrganizationID: 456,
			Name:           "Test Advertiser",
			Status:         "active",
			ContactEmail:   &email,
		}

		req, err := service.mapAdvertiserToEverflowRequest(advertiser)
		assert.NoError(t, err)
		assert.NotNil(t, req.InternalNotes)
		assert.Contains(t, *req.InternalNotes, email)
	})

	// Test case 3: Advertiser with billing details
	t.Run("Advertiser with billing details", func(t *testing.T) {
		line2 := "Suite 100"
		state := "CA"
		taxID := "123456789"
		billingDetails := &domain.BillingDetails{
			BillingFrequency: "monthly",
			TaxID:           &taxID,
			Address: &domain.BillingAddress{
				Line1:      "123 Main St",
				Line2:      &line2,
				City:       "San Francisco",
				State:      &state,
				PostalCode: "94105",
				Country:    "US",
			},
		}

		advertiser := &domain.Advertiser{
			AdvertiserID:   123,
			OrganizationID: 456,
			Name:           "Test Advertiser",
			Status:         "active",
			BillingDetails: billingDetails,
		}

		req, err := service.mapAdvertiserToEverflowRequest(advertiser)
		assert.NoError(t, err)
		assert.NotNil(t, req.IsContactAddressEnabled)
		assert.True(t, *req.IsContactAddressEnabled)
		assert.NotNil(t, req.ContactAddress)
		assert.Equal(t, "123 Main St", req.ContactAddress.Address1)
		assert.Equal(t, "Suite 100", *req.ContactAddress.Address2)
		assert.Equal(t, "San Francisco", req.ContactAddress.City)
		assert.Equal(t, "94105", req.ContactAddress.ZipPostalCode)
		assert.Equal(t, "US", req.ContactAddress.CountryCode)
		assert.Equal(t, "CA", req.ContactAddress.RegionCode)

		assert.NotNil(t, req.Billing)
		assert.Equal(t, "monthly", req.Billing.BillingFrequency)
		assert.Equal(t, "123456789", *req.Billing.TaxID)
	})
}

// Test mapCampaignToEverflowRequest
func TestMapCampaignToEverflowRequest(t *testing.T) {
	// Setup
	mockAdvertiserRepo := new(MockAdvertiserRepository)
	mockProviderMappingRepo := new(MockAdvertiserProviderMappingRepository)
	mockCampaignRepo := new(MockCampaignRepository)
	mockCryptoService := new(MockCryptoService)

	service := &Service{
		client:         NewClient("test-api-key"),
		advertiserRepo:        mockAdvertiserRepo,
		providerMappingRepo:   mockProviderMappingRepo,
		campaignRepo:   mockCampaignRepo,
		cryptoService:  mockCryptoService,
	}

	// Test case 1: Basic campaign
	t.Run("Basic campaign", func(t *testing.T) {
		campaign := &domain.Campaign{
			CampaignID:     789,
			OrganizationID: 456,
			AdvertiserID:   123,
			Name:           "Test Campaign",
			Status:         "active",
		}

		req, err := service.mapCampaignToEverflowRequest(campaign, 12345)
		assert.NoError(t, err)
		assert.Equal(t, "Test Campaign", req.Name)
		assert.Equal(t, int64(12345), req.NetworkAdvertiserID)
		assert.Equal(t, "active", req.OfferStatus)
		assert.Equal(t, "USD", *req.CurrencyID)
		assert.Equal(t, "public", *req.Visibility)
		assert.Equal(t, "server_postback", *req.ConversionMethod)
		assert.Equal(t, fmt.Sprintf("https://example.com/campaigns/%d?click_id={transaction_id}", campaign.CampaignID), req.DestinationURL)

		// Verify payout/revenue
		assert.NotNil(t, req.PayoutRevenue)
		assert.Len(t, req.PayoutRevenue.Entries, 1)
		assert.True(t, req.PayoutRevenue.Entries[0].IsDefault)
		assert.Equal(t, "cpa", req.PayoutRevenue.Entries[0].PayoutType)
		assert.Equal(t, 1.00, *req.PayoutRevenue.Entries[0].PayoutAmount)
		assert.Equal(t, "rpa", req.PayoutRevenue.Entries[0].RevenueType)
		assert.Equal(t, 2.00, *req.PayoutRevenue.Entries[0].RevenueAmount)

		// Verify session settings
		assert.NotNil(t, req.SessionDefinition)
		assert.Equal(t, "cookie", *req.SessionDefinition)
		assert.NotNil(t, req.SessionDuration)
		assert.Equal(t, 720, *req.SessionDuration)

		// Note: Tags are not part of the OfferInput structure
	})

	// Test case 2: Campaign with description
	t.Run("Campaign with description", func(t *testing.T) {
		description := "This is a test campaign description"
		campaign := &domain.Campaign{
			CampaignID:     789,
			OrganizationID: 456,
			AdvertiserID:   123,
			Name:           "Test Campaign",
			Status:         "active",
			Description:    &description,
		}

		req, err := service.mapCampaignToEverflowRequest(campaign, 12345)
		assert.NoError(t, err)
		assert.NotNil(t, req.HTMLDescription)
		assert.Equal(t, description, *req.HTMLDescription)
	})

	// Test case 3: Campaign with different status
	t.Run("Campaign with different status", func(t *testing.T) {
		campaign := &domain.Campaign{
			CampaignID:     789,
			OrganizationID: 456,
			AdvertiserID:   123,
			Name:           "Test Campaign",
			Status:         "paused",
		}

		req, err := service.mapCampaignToEverflowRequest(campaign, 12345)
		assert.NoError(t, err)
		assert.Equal(t, "paused", req.OfferStatus)
	})
}

// Test GetOfferFromEverflow
func TestGetOfferFromEverflow(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/networks/offers/54321", r.URL.Path)

		// Return a successful response for offer retrieval
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"network_offer_id": 54321,
			"network_id": 1,
			"network_advertiser_id": 12345,
			"name": "Test Offer",
			"destination_url": "https://example.com/offer",
			"offer_status": "active",
			"currency_id": "USD",
			"encoded_value": "ABC123",
			"time_created": 1621234567,
			"time_saved": 1621234567,
			"visibility": "public",
			"conversion_method": "server_postback",
			"today_clicks": 100,
			"today_revenue": "$50.00",
			"payout": "CPA: $1.00",
			"revenue": "RPA: $2.00"
		}`)
	}))
	defer server.Close()

	// Setup mocks
	mockAdvertiserRepo := new(MockAdvertiserRepository)
	mockProviderMappingRepo := new(MockAdvertiserProviderMappingRepository)
	mockCampaignRepo := new(MockCampaignRepository)
	mockCryptoService := new(MockCryptoService)

	// Create service with mock client
	client := NewClient("test-api-key")
	client.httpClient = server.Client()

	service := &Service{
		client:              client,
		advertiserRepo:      mockAdvertiserRepo,
		providerMappingRepo: mockProviderMappingRepo,
		campaignRepo:        mockCampaignRepo,
		cryptoService:       mockCryptoService,
	}

	// Override the base URL to point to our test server
	origBaseURL := everflowAPIBaseURL
	everflowAPIBaseURL = server.URL + "/v1"
	defer func() { everflowAPIBaseURL = origBaseURL }()

	// Test the service
	offer, err := service.client.GetOffer(context.Background(), 54321, nil)
	assert.NoError(t, err)
	assert.NotNil(t, offer)
	assert.Equal(t, int64(54321), offer.NetworkOfferID)
	assert.Equal(t, "Test Offer", offer.Name)
	assert.Equal(t, "active", offer.OfferStatus)
	assert.Equal(t, "ABC123", *offer.EncodedValue)
	assert.Equal(t, 100, *offer.TodayClicks)
	assert.Equal(t, "$50.00", *offer.TodayRevenue)
}

// Test UpdateOfferInEverflow
func TestUpdateOfferInEverflow(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/v1/networks/offers/54321", r.URL.Path)

		// Return a successful response for offer update
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"network_offer_id": 54321,
			"network_id": 1,
			"network_advertiser_id": 12345,
			"name": "Updated Test Offer",
			"destination_url": "https://example.com/updated-offer",
			"offer_status": "paused",
			"currency_id": "USD",
			"encoded_value": "ABC123",
			"time_created": 1621234567,
			"time_saved": 1621234568,
			"visibility": "public",
			"conversion_method": "server_postback"
		}`)
	}))
	defer server.Close()

	// Setup mocks
	mockAdvertiserRepo := new(MockAdvertiserRepository)
	mockProviderMappingRepo := new(MockAdvertiserProviderMappingRepository)
	mockCampaignRepo := new(MockCampaignRepository)
	mockCryptoService := new(MockCryptoService)

	// Create service with mock client
	client := NewClient("test-api-key")
	client.httpClient = server.Client()

	service := &Service{
		client:              client,
		advertiserRepo:      mockAdvertiserRepo,
		providerMappingRepo: mockProviderMappingRepo,
		campaignRepo:        mockCampaignRepo,
		cryptoService:       mockCryptoService,
	}

	// Override the base URL to point to our test server
	origBaseURL := everflowAPIBaseURL
	everflowAPIBaseURL = server.URL + "/v1"
	defer func() { everflowAPIBaseURL = origBaseURL }()

	// Create update request
	updateReq := OfferInput{
		Name:           "Updated Test Offer",
		DestinationURL: "https://example.com/updated-offer",
		OfferStatus:    "paused",
	}

	// Test the service
	offer, err := service.client.UpdateOffer(context.Background(), 54321, updateReq)
	assert.NoError(t, err)
	assert.NotNil(t, offer)
	assert.Equal(t, int64(54321), offer.NetworkOfferID)
	assert.Equal(t, "Updated Test Offer", offer.Name)
	assert.Equal(t, "paused", offer.OfferStatus)
	assert.Equal(t, "https://example.com/updated-offer", offer.DestinationURL)
}

// Test OffersTable
func TestOffersTable(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/v1/networks/offerstable", r.URL.Path)

		// Return a successful response for offers table
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"offers": [
				{
					"network_offer_id": 54321,
					"network_id": 1,
					"network_advertiser_id": 12345,
					"name": "Test Offer 1",
					"offer_status": "active",
					"currency_id": "USD",
					"encoded_value": "ABC123",
					"visibility": "public",
					"today_clicks": 50,
					"today_revenue": "$25.00"
				},
				{
					"network_offer_id": 54322,
					"network_id": 1,
					"network_advertiser_id": 12345,
					"name": "Test Offer 2",
					"offer_status": "paused",
					"currency_id": "USD",
					"encoded_value": "DEF456",
					"visibility": "public",
					"today_clicks": 25,
					"today_revenue": "$12.50"
				}
			]
		}`)
	}))
	defer server.Close()

	// Setup mocks
	mockAdvertiserRepo := new(MockAdvertiserRepository)
	mockProviderMappingRepo := new(MockAdvertiserProviderMappingRepository)
	mockCampaignRepo := new(MockCampaignRepository)
	mockCryptoService := new(MockCryptoService)

	// Create service with mock client
	client := NewClient("test-api-key")
	client.httpClient = server.Client()

	service := &Service{
		client:              client,
		advertiserRepo:      mockAdvertiserRepo,
		providerMappingRepo: mockProviderMappingRepo,
		campaignRepo:        mockCampaignRepo,
		cryptoService:       mockCryptoService,
	}

	// Override the base URL to point to our test server
	origBaseURL := everflowAPIBaseURL
	everflowAPIBaseURL = server.URL + "/v1"
	defer func() { everflowAPIBaseURL = origBaseURL }()

	// Create request with filters
	req := OffersTableRequest{
		Filters: map[string]interface{}{
			"offer_status": "active",
		},
		SearchTerms: []OffersTableSearchTerm{
			{
				SearchType: "name",
				Value:      "Test",
			},
		},
	}

	// Test the service
	response, err := service.client.OffersTable(context.Background(), req, nil)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Offers, 2)
	
	// Verify first offer
	offer1 := response.Offers[0]
	assert.Equal(t, int64(54321), offer1.NetworkOfferID)
	assert.Equal(t, "Test Offer 1", offer1.Name)
	assert.Equal(t, "active", offer1.OfferStatus)
	assert.Equal(t, 50, *offer1.TodayClicks)
	assert.Equal(t, "$25.00", *offer1.TodayRevenue)

	// Verify second offer
	offer2 := response.Offers[1]
	assert.Equal(t, int64(54322), offer2.NetworkOfferID)
	assert.Equal(t, "Test Offer 2", offer2.Name)
	assert.Equal(t, "paused", offer2.OfferStatus)
	assert.Equal(t, 25, *offer2.TodayClicks)
	assert.Equal(t, "$12.50", *offer2.TodayRevenue)
}

// Test service methods that use the new client functionality
func TestServiceWithNewClientMethods(t *testing.T) {
	// Setup mocks
	mockAdvertiserRepo := new(MockAdvertiserRepository)
	mockProviderMappingRepo := new(MockAdvertiserProviderMappingRepository)
	mockCampaignRepo := new(MockCampaignRepository)
	mockCryptoService := new(MockCryptoService)

	service := &Service{
		client:              NewClient("test-api-key"),
		advertiserRepo:      mockAdvertiserRepo,
		providerMappingRepo: mockProviderMappingRepo,
		campaignRepo:        mockCampaignRepo,
		cryptoService:       mockCryptoService,
	}

	t.Run("mapCampaignToEverflowRequest returns OfferInput", func(t *testing.T) {
		campaign := &domain.Campaign{
			CampaignID:     789,
			OrganizationID: 456,
			AdvertiserID:   123,
			Name:           "Test Campaign",
			Status:         "active",
		}

		req, err := service.mapCampaignToEverflowRequest(campaign, 12345)
		assert.NoError(t, err)
		assert.IsType(t, &OfferInput{}, req)
		assert.Equal(t, "Test Campaign", req.Name)
		assert.Equal(t, int64(12345), req.NetworkAdvertiserID)
		assert.Equal(t, "active", req.OfferStatus)
		assert.NotNil(t, req.CurrencyID)
		assert.Equal(t, "USD", *req.CurrencyID)
		assert.NotNil(t, req.Visibility)
		assert.Equal(t, "public", *req.Visibility)
		assert.NotNil(t, req.ConversionMethod)
		assert.Equal(t, "server_postback", *req.ConversionMethod)
		assert.NotNil(t, req.PayoutRevenue)
		assert.Len(t, req.PayoutRevenue.Entries, 1)
		
		// Verify payout revenue entry
		entry := req.PayoutRevenue.Entries[0]
		assert.Equal(t, "cpa", entry.PayoutType)
		assert.Equal(t, "rpa", entry.RevenueType)
		assert.Equal(t, 1.00, *entry.PayoutAmount)
		assert.Equal(t, 2.00, *entry.RevenueAmount)
		assert.True(t, entry.IsDefault)
		assert.False(t, entry.IsPrivate)
	})
}

// Test new service methods
func TestServiceNewMethods(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		switch {
		case r.Method == "GET" && r.URL.Path == "/v1/networks/offers/54321":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"network_offer_id": 54321,
				"network_id": 1,
				"network_advertiser_id": 12345,
				"name": "Test Offer",
				"offer_status": "active",
				"currency_id": "USD",
				"encoded_value": "ABC123",
				"time_created": 1621234567,
				"time_saved": 1621234567
			}`)
		case r.Method == "PUT" && r.URL.Path == "/v1/networks/offers/54321":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"network_offer_id": 54321,
				"network_id": 1,
				"network_advertiser_id": 12345,
				"name": "Updated Test Offer",
				"offer_status": "paused",
				"currency_id": "USD",
				"encoded_value": "ABC123",
				"time_created": 1621234567,
				"time_saved": 1621234568
			}`)
		case r.Method == "POST" && r.URL.Path == "/v1/networks/offerstable":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"offers": [
					{
						"network_offer_id": 54321,
						"name": "Test Offer 1",
						"offer_status": "active"
					}
				]
			}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Setup mocks
	mockAdvertiserRepo := new(MockAdvertiserRepository)
	mockProviderMappingRepo := new(MockAdvertiserProviderMappingRepository)
	mockCampaignRepo := new(MockCampaignRepository)
	mockCryptoService := new(MockCryptoService)

	// Create service with mock client
	client := NewClient("test-api-key")
	client.httpClient = server.Client()

	service := &Service{
		client:              client,
		advertiserRepo:      mockAdvertiserRepo,
		providerMappingRepo: mockProviderMappingRepo,
		campaignRepo:        mockCampaignRepo,
		cryptoService:       mockCryptoService,
	}

	// Override the base URL to point to our test server
	origBaseURL := everflowAPIBaseURL
	everflowAPIBaseURL = server.URL + "/v1"
	defer func() { everflowAPIBaseURL = origBaseURL }()

	t.Run("GetOfferFromEverflow", func(t *testing.T) {
		offer, err := service.GetOfferFromEverflow(context.Background(), 54321, nil)
		assert.NoError(t, err)
		assert.NotNil(t, offer)
		assert.Equal(t, int64(54321), offer.NetworkOfferID)
		assert.Equal(t, "Test Offer", offer.Name)
		assert.Equal(t, "active", offer.OfferStatus)
	})

	t.Run("UpdateOfferInEverflow", func(t *testing.T) {
		updateReq := OfferInput{
			Name:        "Updated Test Offer",
			OfferStatus: "paused",
		}

		offer, err := service.UpdateOfferInEverflow(context.Background(), 54321, updateReq)
		assert.NoError(t, err)
		assert.NotNil(t, offer)
		assert.Equal(t, int64(54321), offer.NetworkOfferID)
		assert.Equal(t, "Updated Test Offer", offer.Name)
		assert.Equal(t, "paused", offer.OfferStatus)
	})

	t.Run("ListOffersFromEverflow", func(t *testing.T) {
		req := OffersTableRequest{
			Filters: map[string]interface{}{
				"offer_status": "active",
			},
		}

		response, err := service.ListOffersFromEverflow(context.Background(), req, nil)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Offers, 1)
		assert.Equal(t, int64(54321), response.Offers[0].NetworkOfferID)
	})

	t.Run("GetOfferFromEverflowByMapping", func(t *testing.T) {
		// Setup campaign provider offer
		providerOfferRef := "54321"
		campaignOffer := &domain.CampaignProviderOffer{
			CampaignID:       789,
			ProviderType:     "everflow",
			ProviderOfferRef: &providerOfferRef,
		}

		mockCampaignRepo.On("ListCampaignProviderOffersByCampaign", mock.Anything, int64(789)).Return([]*domain.CampaignProviderOffer{campaignOffer}, nil)

		offer, err := service.GetOfferFromEverflowByMapping(context.Background(), 789, nil)
		assert.NoError(t, err)
		assert.NotNil(t, offer)
		assert.Equal(t, int64(54321), offer.NetworkOfferID)

		mockCampaignRepo.AssertExpectations(t)
		
		// Clear mock expectations for next test
		mockCampaignRepo.ExpectedCalls = nil
		mockCampaignRepo.Calls = nil
	})

	t.Run("UpdateOfferInEverflowByMapping", func(t *testing.T) {
		// Setup campaign provider offer
		providerOfferRef := "54321"
		campaignOffer := &domain.CampaignProviderOffer{
			CampaignProviderOfferID: 1,
			CampaignID:              789,
			ProviderType:            "everflow",
			ProviderOfferRef:        &providerOfferRef,
		}

		mockCampaignRepo.On("ListCampaignProviderOffersByCampaign", mock.Anything, int64(789)).Return([]*domain.CampaignProviderOffer{campaignOffer}, nil)
		mockCampaignRepo.On("UpdateCampaignProviderOffer", mock.Anything, mock.MatchedBy(func(offer *domain.CampaignProviderOffer) bool {
			assert.Equal(t, int64(1), offer.CampaignProviderOfferID)
			assert.False(t, offer.IsActiveOnProvider) // Should be false because offer status is "paused"
			assert.NotNil(t, offer.LastSyncedAt)
			assert.NotNil(t, offer.ProviderOfferConfig)
			return true
		})).Return(nil)

		updateReq := OfferInput{
			Name:        "Updated Test Offer",
			OfferStatus: "paused",
		}

		offer, err := service.UpdateOfferInEverflowByMapping(context.Background(), 789, updateReq)
		assert.NoError(t, err)
		assert.NotNil(t, offer)
		assert.Equal(t, "Updated Test Offer", offer.Name)
		assert.Equal(t, "paused", offer.OfferStatus)

		mockCampaignRepo.AssertExpectations(t)
		
		// Clear mock expectations for next test
		mockCampaignRepo.ExpectedCalls = nil
		mockCampaignRepo.Calls = nil
	})

	t.Run("SyncCampaignWithEverflowOffer", func(t *testing.T) {
		// Setup campaign
		campaign := &domain.Campaign{
			CampaignID: 789,
			Name:       "Test Campaign",
			Status:     "draft",
		}

		// Setup campaign provider offer
		providerOfferRef := "54321"
		campaignOffer := &domain.CampaignProviderOffer{
			CampaignProviderOfferID: 1,
			CampaignID:              789,
			ProviderType:            "everflow",
			ProviderOfferRef:        &providerOfferRef,
		}

		mockCampaignRepo.On("GetCampaignByID", mock.Anything, int64(789)).Return(campaign, nil)
		mockCampaignRepo.On("ListCampaignProviderOffersByCampaign", mock.Anything, int64(789)).Return([]*domain.CampaignProviderOffer{campaignOffer}, nil).Twice()
		mockCampaignRepo.On("UpdateCampaign", mock.Anything, mock.MatchedBy(func(c *domain.Campaign) bool {
			assert.Equal(t, "active", c.Status) // Should be updated to active
			return true
		})).Return(nil)
		mockCampaignRepo.On("UpdateCampaignProviderOffer", mock.Anything, mock.MatchedBy(func(offer *domain.CampaignProviderOffer) bool {
			assert.True(t, offer.IsActiveOnProvider) // Should be true because offer status is "active"
			assert.NotNil(t, offer.LastSyncedAt)
			return true
		})).Return(nil)

		err := service.SyncCampaignWithEverflowOffer(context.Background(), 789)
		assert.NoError(t, err)

		mockCampaignRepo.AssertExpectations(t)
	})
}
