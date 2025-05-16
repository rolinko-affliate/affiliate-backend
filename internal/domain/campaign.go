package domain

import (
	"time"
)

// Campaign represents a campaign entity
type Campaign struct {
	CampaignID     int64      `json:"campaign_id" db:"campaign_id"`
	OrganizationID int64      `json:"organization_id" db:"organization_id"`
	AdvertiserID   int64      `json:"advertiser_id" db:"advertiser_id"`
	Name           string     `json:"name" db:"name"`
	Description    *string    `json:"description,omitempty" db:"description"`
	Status         string     `json:"status" db:"status"` // 'draft', 'active', 'paused', 'archived'
	StartDate      *time.Time `json:"start_date,omitempty" db:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty" db:"end_date"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// CampaignProviderOffer represents an offer for a campaign in a provider
type CampaignProviderOffer struct {
	ProviderOfferID    int64      `json:"provider_offer_id" db:"provider_offer_id"`
	CampaignID         int64      `json:"campaign_id" db:"campaign_id"`
	ProviderType       string     `json:"provider_type" db:"provider_type"` // 'everflow' for MVP
	ProviderOfferRef   *string    `json:"provider_offer_ref,omitempty" db:"provider_offer_ref"` // Provider's Offer ID (e.g., Everflow's network_offer_id)
	ProviderOfferConfig *string    `json:"provider_offer_config,omitempty" db:"provider_offer_config"` // JSONB stored as string
	IsActiveOnProvider  bool       `json:"is_active_on_provider" db:"is_active_on_provider"`
	LastSyncedAt        *time.Time `json:"last_synced_at,omitempty" db:"last_synced_at"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
}