package domain

import (
	"time"
)

// Affiliate represents an affiliate entity
type Affiliate struct {
	AffiliateID    int64     `json:"affiliate_id" db:"affiliate_id"`
	OrganizationID int64     `json:"organization_id" db:"organization_id"`
	Name           string    `json:"name" db:"name"`
	ContactEmail   *string   `json:"contact_email,omitempty" db:"contact_email"`
	PaymentDetails *string   `json:"payment_details,omitempty" db:"payment_details"` // JSONB stored as string
	Status         string    `json:"status" db:"status"` // 'active', 'pending', 'rejected', 'inactive'
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// AffiliateProviderMapping represents a mapping between an affiliate and a provider
type AffiliateProviderMapping struct {
	MappingID           int64     `json:"mapping_id" db:"mapping_id"`
	AffiliateID         int64     `json:"affiliate_id" db:"affiliate_id"`
	ProviderType        string    `json:"provider_type" db:"provider_type"` // 'everflow' for MVP
	ProviderAffiliateID *string   `json:"provider_affiliate_id,omitempty" db:"provider_affiliate_id"` // e.g., Everflow's network_affiliate_id
	ProviderConfig      *string   `json:"provider_config,omitempty" db:"provider_config"` // JSONB stored as string
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}