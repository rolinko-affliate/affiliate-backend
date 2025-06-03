package provider

import (
	"context"

	"github.com/affiliate-backend/internal/domain"
)

// ProviderAdvertiserService defines the interface for advertiser operations with external providers
type ProviderAdvertiserService interface {
	// CreateAdvertiserInProvider creates an advertiser in the external provider and stores the mapping
	CreateAdvertiserInProvider(ctx context.Context, advertiser *domain.Advertiser) error
	
	// GetAdvertiserFromProvider retrieves an advertiser from the external provider using our internal advertiser ID
	// Returns a domain.Advertiser object with provider data mapped to our domain model
	GetAdvertiserFromProvider(ctx context.Context, advertiserID int64, relationships []string) (*domain.Advertiser, error)
	
	// UpdateAdvertiserInProvider updates an advertiser in the external provider using our internal advertiser ID
	// Returns a domain.Advertiser object with updated provider data mapped to our domain model
	UpdateAdvertiserInProvider(ctx context.Context, advertiserID int64, advertiser *domain.Advertiser) (*domain.Advertiser, error)
}

// ProviderCampaignService defines the interface for campaign operations with external providers
type ProviderCampaignService interface {
	// CreateOfferInProvider creates an offer in the external provider for a campaign and stores the mapping
	CreateOfferInProvider(ctx context.Context, campaign *domain.Campaign) error
	
	// GetOfferFromProvider retrieves an offer from the external provider using our internal campaign ID
	// Returns a domain.Campaign object with provider offer data mapped to our domain model
	GetOfferFromProvider(ctx context.Context, campaignID int64, relationships []string) (*domain.Campaign, error)
	
	// UpdateOfferInProvider updates an offer in the external provider using our internal campaign ID
	// Returns a domain.Campaign object with updated provider offer data mapped to our domain model
	UpdateOfferInProvider(ctx context.Context, campaignID int64, campaign *domain.Campaign) (*domain.Campaign, error)
}
