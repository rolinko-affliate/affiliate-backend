package everflow

import (
	"github.com/affiliate-backend/internal/platform/everflow/advertiser"
	"github.com/affiliate-backend/internal/platform/everflow/affiliate"
	"github.com/affiliate-backend/internal/platform/everflow/offer"
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
	advertiserConfig := advertiser.NewConfiguration()
	advertiserConfig.Servers = []advertiser.ServerConfiguration{
		{
			URL: config.BaseURL,
		},
	}
	// Add Everflow API key header
	advertiserConfig.AddDefaultHeader("X-Eflow-API-Key", config.APIKey)
	advertiserClient := advertiser.NewAPIClient(advertiserConfig)

	// Configure affiliate client
	affiliateConfig := affiliate.NewConfiguration()
	affiliateConfig.Servers = []affiliate.ServerConfiguration{
		{
			URL: config.BaseURL + "/v1", // Affiliate API expects /v1 in server URL
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

	return NewIntegrationService(
		advertiserClient,
		affiliateClient,
		offerClient,
		advertiserRepo,
		affiliateRepo,
		campaignRepo,
		advertiserProviderMappingRepo,
		affiliateProviderMappingRepo,
		campaignProviderMappingRepo,
	)
}
