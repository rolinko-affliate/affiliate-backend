package everflow

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/everflow/tracking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)



func TestGenerateTrackingLink(t *testing.T) {
	// Create a mock HTTP server to simulate Everflow API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/v1/networks/tracking/offers/clicks", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Mock response
		response := `{
			"tracking_url": "https://tracking.example.com/abc123/def456/?sub1=facebook&sub3=media_buying",
			"network_affiliate_id": 8,
			"network_offer_id": 20,
			"sub1": "facebook",
			"sub3": "media_buying",
			"is_encrypt_parameters": true
		}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer mockServer.Close()

	// Create tracking client with mock server
	config := tracking.NewConfiguration()
	config.Scheme = "http"
	config.Host = strings.TrimPrefix(mockServer.URL, "http://")
	trackingClient := tracking.NewAPIClient(config)

	// Create mock repositories
	mockAffiliateRepo := &MockAffiliateRepository{}
	mockCampaignRepo := &MockCampaignRepository{}
	mockAdvertiserRepo := &MockAdvertiserRepository{}
	mockAffiliateProviderMappingRepo := &MockAffiliateProviderMappingRepository{}
	mockAdvertiserProviderMappingRepo := &MockAdvertiserProviderMappingRepository{}
	mockCampaignProviderMappingRepo := &MockCampaignProviderMappingRepository{}

	// Create test mappings with proper provider data
	affiliateProviderData := domain.EverflowProviderData{
		NetworkAffiliateID: int32Ptr(8), // Use affiliate ID 8 as specified
	}
	affiliateProviderDataJSON, _ := json.Marshal(affiliateProviderData)
	affiliateProviderDataStr := string(affiliateProviderDataJSON)

	affiliateMapping := &domain.AffiliateProviderMapping{
		AffiliateID:         123,
		ProviderType:        "everflow",
		ProviderAffiliateID: stringPtr("8"),
		ProviderData:        &affiliateProviderDataStr,
	}

	campaignProviderData := domain.EverflowCampaignProviderData{
		NetworkCampaignID: int32Ptr(20), // Use offer ID 20 as specified
	}
	campaignProviderDataJSON, _ := campaignProviderData.ToJSON()

	campaignMapping := &domain.CampaignProviderMapping{
		CampaignID:         456,
		ProviderType:       "everflow",
		ProviderCampaignID: stringPtr("20"),
		ProviderData:       &campaignProviderDataJSON,
	}

	// Create integration service
	service := NewIntegrationService(
		nil, // advertiserClient not needed for this test
		nil, // affiliateClient not needed for this test
		nil, // offerClient not needed for this test
		trackingClient,
		mockAdvertiserRepo,
		mockAffiliateRepo,
		mockCampaignRepo,
		mockAdvertiserProviderMappingRepo,
		mockAffiliateProviderMappingRepo,
		mockCampaignProviderMappingRepo,
	)

	// Create test request
	sub1 := "facebook"
	sub3 := "media_buying"
	isEncrypt := true
	
	req := &domain.TrackingLinkGenerationRequest{
		AffiliateID:         123,
		CampaignID:          456,
		Sub1:                &sub1,
		Sub3:                &sub3,
		IsEncryptParameters: &isEncrypt,
	}

	// Execute the function
	ctx := context.Background()
	resp, err := service.GenerateTrackingLink(ctx, req, campaignMapping, affiliateMapping)

	// Verify the results
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "https://tracking.example.com/abc123/def456/?sub1=facebook&sub3=media_buying", resp.GeneratedURL)
	assert.NotNil(t, resp.ProviderData)

	// Verify provider data
	var providerData domain.EverflowTrackingLinkProviderData
	err = providerData.FromJSON(*resp.ProviderData)
	require.NoError(t, err)
	assert.Equal(t, int32(8), *providerData.NetworkAffiliateID)
	assert.Equal(t, int32(20), *providerData.NetworkOfferID)
	assert.Equal(t, "https://tracking.example.com/abc123/def456/?sub1=facebook&sub3=media_buying", *providerData.GeneratedURL)

}

// Mock repository implementations for testing
type MockAffiliateRepository struct{}
func (m *MockAffiliateRepository) GetAffiliateByID(ctx context.Context, id int64) (*domain.Affiliate, error) {
	return nil, nil
}

type MockCampaignRepository struct{}
func (m *MockCampaignRepository) GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error) {
	return nil, nil
}

type MockAdvertiserRepository struct{}
func (m *MockAdvertiserRepository) GetAdvertiserByID(ctx context.Context, id int64) (*domain.Advertiser, error) {
	return nil, nil
}

type MockAdvertiserProviderMappingRepository struct{}
func (m *MockAdvertiserProviderMappingRepository) CreateMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error {
	return nil
}
func (m *MockAdvertiserProviderMappingRepository) GetMappingByID(ctx context.Context, id int64) (*domain.AdvertiserProviderMapping, error) {
	return nil, nil
}
func (m *MockAdvertiserProviderMappingRepository) GetMappingsByAdvertiserID(ctx context.Context, advertiserID int64) ([]*domain.AdvertiserProviderMapping, error) {
	return nil, nil
}
func (m *MockAdvertiserProviderMappingRepository) GetMappingByAdvertiserAndProvider(ctx context.Context, advertiserID int64, providerType string) (*domain.AdvertiserProviderMapping, error) {
	return nil, nil
}
func (m *MockAdvertiserProviderMappingRepository) UpdateMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error {
	return nil
}
func (m *MockAdvertiserProviderMappingRepository) DeleteMapping(ctx context.Context, id int64) error {
	return nil
}
func (m *MockAdvertiserProviderMappingRepository) UpdateSyncStatus(ctx context.Context, mappingID int64, status string, syncError *string) error {
	return nil
}


type MockCampaignProviderMappingRepository struct {
	mock.Mock
}
func (m *MockCampaignProviderMappingRepository) GetCampaignProviderMapping(ctx context.Context, campaignID int64, providerType string) (*domain.CampaignProviderMapping, error) {
	args := m.Called(ctx, campaignID, providerType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CampaignProviderMapping), args.Error(1)
}
func (m *MockCampaignProviderMappingRepository) CreateCampaignProviderMapping(ctx context.Context, mapping *domain.CampaignProviderMapping) error {
	return nil
}
func (m *MockCampaignProviderMappingRepository) UpdateCampaignProviderMapping(ctx context.Context, mapping *domain.CampaignProviderMapping) error {
	return nil
}