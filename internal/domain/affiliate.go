package domain

import (
	"fmt"
	"time"
)

// ContactAddress represents contact address information for an affiliate
type ContactAddress struct {
	Address1      *string `json:"address1,omitempty"`
	Address2      *string `json:"address2,omitempty"`
	City          *string `json:"city,omitempty"`
	RegionCode    *string `json:"region_code,omitempty"`
	CountryCode   *string `json:"country_code,omitempty"`
	ZipPostalCode *string `json:"zip_postal_code,omitempty"`
}

// HasData returns true if any contact address field has data
func (ca *ContactAddress) HasData() bool {
	if ca == nil {
		return false
	}
	return ca.Address1 != nil || ca.Address2 != nil || ca.City != nil ||
		ca.RegionCode != nil || ca.CountryCode != nil || ca.ZipPostalCode != nil
}

// Affiliate represents an affiliate entity (clean domain model)
type Affiliate struct {
	AffiliateID    int64   `json:"affiliate_id" db:"affiliate_id"`
	OrganizationID int64   `json:"organization_id" db:"organization_id"`
	Name           string  `json:"name" db:"name"`
	ContactEmail   *string `json:"contact_email,omitempty" db:"contact_email"`
	PaymentDetails *string `json:"payment_details,omitempty" db:"payment_details"` // JSONB stored as string
	Status         string  `json:"status" db:"status"`                             // 'active', 'pending', 'rejected', 'inactive'

	// General purpose fields moved from EverflowProviderData
	InternalNotes          *string  `json:"internal_notes,omitempty" db:"internal_notes"`
	DefaultCurrencyID      *string  `json:"default_currency_id,omitempty" db:"default_currency_id"`
	ContactAddress         *string  `json:"contact_address,omitempty" db:"contact_address"` // JSONB stored as string
	BillingInfo            *string  `json:"billing_info,omitempty" db:"billing_info"`       // JSONB stored as string
	Labels                 *string  `json:"labels,omitempty" db:"labels"`                   // JSONB stored as string (array of strings)
	InvoiceAmountThreshold *float64 `json:"invoice_amount_threshold,omitempty" db:"invoice_amount_threshold"`
	DefaultPaymentTerms    *int32   `json:"default_payment_terms,omitempty" db:"default_payment_terms"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// AffiliateProviderMapping represents a mapping between an affiliate and a provider
type AffiliateProviderMapping struct {
	MappingID           int64   `json:"mapping_id" db:"mapping_id"`
	AffiliateID         int64   `json:"affiliate_id" db:"affiliate_id"`
	ProviderType        string  `json:"provider_type" db:"provider_type"`                           // 'everflow' for MVP
	ProviderAffiliateID *string `json:"provider_affiliate_id,omitempty" db:"provider_affiliate_id"` // e.g., Everflow's network_affiliate_id
	APICredentials      *string `json:"api_credentials,omitempty" db:"api_credentials"`             // JSONB stored as string for API keys/tokens
	ProviderConfig      *string `json:"provider_config,omitempty" db:"provider_config"`             // JSONB stored as string

	// Provider-specific data (stored as JSONB) - contains all Everflow-specific fields
	ProviderData *string `json:"provider_data,omitempty" db:"provider_data"` // JSONB for provider-specific fields

	// Synchronization metadata
	SyncStatus *string    `json:"sync_status,omitempty" db:"sync_status"` // 'synced', 'pending', 'error'
	LastSyncAt *time.Time `json:"last_sync_at,omitempty" db:"last_sync_at"`
	SyncError  *string    `json:"sync_error,omitempty" db:"sync_error"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// AffiliateExtraInfo represents additional information for an affiliate organization
type AffiliateExtraInfo struct {
	ExtraInfoID     int64     `json:"extra_info_id" db:"extra_info_id"`
	OrganizationID  int64     `json:"organization_id" db:"organization_id"`
	Website         *string   `json:"website,omitempty" db:"website"`
	AffiliateType   *string   `json:"affiliate_type,omitempty" db:"affiliate_type"` // 'cashback', 'blog', 'incentive', 'content', 'forum', 'sub_affiliate_network'
	SelfDescription *string   `json:"self_description,omitempty" db:"self_description"`
	LogoURL         *string   `json:"logo_url,omitempty" db:"logo_url"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Validate validates the affiliate extra info data
func (aei *AffiliateExtraInfo) Validate() error {
	if aei.OrganizationID <= 0 {
		return fmt.Errorf("valid organization ID is required")
	}
	
	if aei.AffiliateType != nil {
		validTypes := map[string]bool{
			"cashback": true, "blog": true, "incentive": true, 
			"content": true, "forum": true, "sub_affiliate_network": true,
		}
		if !validTypes[*aei.AffiliateType] {
			return fmt.Errorf("invalid affiliate type: %s", *aei.AffiliateType)
		}
	}
	
	return nil
}

// AffiliateWithExtraInfo represents an affiliate with its extra information
type AffiliateWithExtraInfo struct {
	*Affiliate
	ExtraInfo *AffiliateExtraInfo `json:"extra_info,omitempty"`
}

// EverflowProviderData represents Everflow-specific data stored in ProviderData field
type EverflowProviderData struct {
	// Everflow-specific fields only (general purpose fields moved to main Affiliate model)
	NetworkAffiliateID           *int32 `json:"network_affiliate_id,omitempty"`
	EnableMediaCostTrackingLinks *bool  `json:"enable_media_cost_tracking_links,omitempty"`
	ReferrerID                   *int32 `json:"referrer_id,omitempty"`
	IsContactAddressEnabled      *bool  `json:"is_contact_address_enabled,omitempty"`
	NetworkAffiliateTierID       *int32 `json:"network_affiliate_tier_id,omitempty"`
	NetworkEmployeeID            *int32 `json:"network_employee_id,omitempty"`

	// Everflow-specific structured data
	Users *[]interface{} `json:"users,omitempty"`

	// Additional fields for extensibility
	AdditionalFields map[string]interface{} `json:"additional_fields,omitempty"`
}
