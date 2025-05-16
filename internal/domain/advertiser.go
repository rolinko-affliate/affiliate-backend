package domain

import (
	"time"
)

// Advertiser represents an advertiser entity
type Advertiser struct {
	AdvertiserID    int64     `json:"advertiser_id" db:"advertiser_id"`
	OrganizationID  int64     `json:"organization_id" db:"organization_id"`
	Name            string    `json:"name" db:"name"`
	ContactEmail    *string   `json:"contact_email,omitempty" db:"contact_email"`
	BillingDetails  *string   `json:"billing_details,omitempty" db:"billing_details"` // JSONB stored as string
	Status          string    `json:"status" db:"status"` // 'active', 'pending', 'inactive', 'rejected'
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// AdvertiserProviderMapping represents a mapping between an advertiser and a provider
type AdvertiserProviderMapping struct {
	MappingID            int64     `json:"mapping_id" db:"mapping_id"`
	AdvertiserID         int64     `json:"advertiser_id" db:"advertiser_id"`
	ProviderType         string    `json:"provider_type" db:"provider_type"` // 'everflow' for MVP
	ProviderAdvertiserID *string   `json:"provider_advertiser_id,omitempty" db:"provider_advertiser_id"` // e.g., Everflow's network_advertiser_id
	APICredentials       *string   `json:"-" db:"api_credentials"` // Exclude from JSON, JSONB stored as string, ENCRYPTED
	ProviderConfig       *string   `json:"provider_config,omitempty" db:"provider_config"` // JSONB stored as string
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}