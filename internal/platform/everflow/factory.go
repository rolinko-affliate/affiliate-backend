package everflow

import (
	"strings"

	"github.com/affiliate-backend/internal/platform/everflow/advertiser"
	"github.com/affiliate-backend/internal/platform/everflow/affiliate"
	"github.com/affiliate-backend/internal/platform/everflow/offer"
	"github.com/affiliate-backend/internal/platform/everflow/reporting"
	"github.com/affiliate-backend/internal/platform/everflow/tracking"
)

// Config holds the configuration for Everflow clients
type Config struct {
	BaseURL string
	APIKey  string
}

// NewIntegrationServiceWithClients creates a new IntegrationService with configured Everflow clients
func NewIntegrationServiceWithClients(
	config Config,
	advertiserRepo AdvertiserRepository,
	affiliateRepo AffiliateRepository,
	campaignRepo CampaignRepository,
	advertiserProviderMappingRepo AdvertiserProviderMappingRepository,
	affiliateProviderMappingRepo AffiliateProviderMappingRepository,
	campaignProviderMappingRepo CampaignProviderMappingRepository,
) *IntegrationService {
	// Configure advertiser client
	// Note: Advertiser client uses /v1/networks/advertisers path, so we need base URL without /v1
	advertiserBaseURL := strings.TrimSuffix(config.BaseURL, "/v1")
	advertiserConfig := advertiser.NewConfiguration()
	advertiserConfig.Servers = []advertiser.ServerConfiguration{
		{
			URL: advertiserBaseURL,
		},
	}
	// Add Everflow API key header
	advertiserConfig.AddDefaultHeader("X-Eflow-API-Key", config.APIKey)
	advertiserClient := advertiser.NewAPIClient(advertiserConfig)

	// Configure affiliate client
	affiliateConfig := affiliate.NewConfiguration()
	affiliateConfig.Servers = []affiliate.ServerConfiguration{
		{
			URL: config.BaseURL, // BaseURL already includes /v1
		},
	}
	// Add Everflow API key header
	affiliateConfig.AddDefaultHeader("X-Eflow-API-Key", config.APIKey)
	affiliateClient := affiliate.NewAPIClient(affiliateConfig)

	// Configure offer client
	offerConfig := offer.NewConfiguration()
	offerConfig.Servers = []offer.ServerConfiguration{
		{
			URL: config.BaseURL,
		},
	}
	// Add Everflow API key header
	offerConfig.AddDefaultHeader("X-Eflow-API-Key", config.APIKey)
	offerClient := offer.NewAPIClient(offerConfig)

	// Configure tracking client
	trackingConfig := tracking.NewConfiguration()
	trackingConfig.Servers = []tracking.ServerConfiguration{
		{
			URL: config.BaseURL,
		},
	}
	// Add Everflow API key header
	trackingConfig.AddDefaultHeader("X-Eflow-API-Key", config.APIKey)
	trackingClient := tracking.NewAPIClient(trackingConfig)

	return NewIntegrationService(
		advertiserClient,
		affiliateClient,
		offerClient,
		trackingClient,
		advertiserRepo,
		affiliateRepo,
		campaignRepo,
		advertiserProviderMappingRepo,
		affiliateProviderMappingRepo,
		campaignProviderMappingRepo,
	)
}

// NewReportingClient creates a new Everflow reporting client
func NewReportingClient(config Config) *reporting.Client {
	// Use the base URL without /v1 suffix for reporting API
	baseURL := strings.TrimSuffix(config.BaseURL, "/v1")
	return reporting.NewClient(baseURL, config.APIKey)
}
