package domain

import (
	"time"
)

// ContactAddress represents contact address information for an affiliate
type ContactAddress struct {
	Address1       *string `json:"address1,omitempty"`
	Address2       *string `json:"address2,omitempty"`
	City           *string `json:"city,omitempty"`
	RegionCode     *string `json:"region_code,omitempty"`
	CountryCode    *string `json:"country_code,omitempty"`
	ZipPostalCode  *string `json:"zip_postal_code,omitempty"`
}

// HasData returns true if any contact address field has data
func (ca *ContactAddress) HasData() bool {
	if ca == nil {
		return false
	}
	return ca.Address1 != nil || ca.Address2 != nil || ca.City != nil ||
		ca.RegionCode != nil || ca.CountryCode != nil || ca.ZipPostalCode != nil
}

// Affiliate represents an affiliate entity
type Affiliate struct {
	AffiliateID    int64     `json:"affiliate_id" db:"affiliate_id"`
	OrganizationID int64     `json:"organization_id" db:"organization_id"`
	Name           string    `json:"name" db:"name"`
	ContactEmail   *string   `json:"contact_email,omitempty" db:"contact_email"`
	PaymentDetails *string   `json:"payment_details,omitempty" db:"payment_details"` // JSONB stored as string
	Status         string    `json:"status" db:"status"` // 'active', 'pending', 'rejected', 'inactive'
	
	// Enhanced fields for Everflow integration
	NetworkAffiliateID              *int32  `json:"network_affiliate_id,omitempty" db:"network_affiliate_id"`
	InternalNotes                   *string `json:"internal_notes,omitempty" db:"internal_notes"`
	DefaultCurrencyID               *string `json:"default_currency_id,omitempty" db:"default_currency_id"`
	EnableMediaCostTrackingLinks    *bool   `json:"enable_media_cost_tracking_links,omitempty" db:"enable_media_cost_tracking_links"`
	ReferrerID                      *int32  `json:"referrer_id,omitempty" db:"referrer_id"`
	IsContactAddressEnabled         *bool   `json:"is_contact_address_enabled,omitempty" db:"is_contact_address_enabled"`
	NetworkAffiliateTierID          *int32  `json:"network_affiliate_tier_id,omitempty" db:"network_affiliate_tier_id"`
	NetworkEmployeeID               *int32  `json:"network_employee_id,omitempty" db:"network_employee_id"`
	
	// Structured data stored as JSONB in database
	ContactAddress *ContactAddress `json:"contact_address,omitempty" db:"contact_address"`
	BillingInfo    *BillingDetails `json:"billing_info,omitempty" db:"billing_info"`
	Labels         *string         `json:"labels,omitempty" db:"labels"` // JSONB stored as string for labels array
	
	// Additional billing fields that don't fit in BillingDetails
	InvoiceAmountThreshold *float64 `json:"invoice_amount_threshold,omitempty" db:"invoice_amount_threshold"`
	DefaultPaymentTerms    *int32   `json:"default_payment_terms,omitempty" db:"default_payment_terms"`
	

	
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// AffiliateProviderMapping represents a mapping between an affiliate and a provider
type AffiliateProviderMapping struct {
	MappingID           int64     `json:"mapping_id" db:"mapping_id"`
	AffiliateID         int64     `json:"affiliate_id" db:"affiliate_id"`
	ProviderType        string    `json:"provider_type" db:"provider_type"` // 'everflow' for MVP
	ProviderAffiliateID *string   `json:"provider_affiliate_id,omitempty" db:"provider_affiliate_id"` // e.g., Everflow's network_affiliate_id
	APICredentials      *string   `json:"api_credentials,omitempty" db:"api_credentials"` // JSONB stored as string for API keys/tokens
	ProviderConfig      *string   `json:"provider_config,omitempty" db:"provider_config"` // JSONB stored as string
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}