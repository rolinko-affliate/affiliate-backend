package provider

import (
	"context"

	"github.com/affiliate-backend/internal/domain"
	"github.com/google/uuid"
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