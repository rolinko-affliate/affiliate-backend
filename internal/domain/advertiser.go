package domain

import (
	"time"
)

// Advertiser represents an advertiser entity (clean domain model)
type Advertiser struct {
	AdvertiserID    int64     `json:"advertiser_id" db:"advertiser_id"`
	OrganizationID  int64     `json:"organization_id" db:"organization_id"`
	Name            string    `json:"name" db:"name"`
	ContactEmail    *string   `json:"contact_email,omitempty" db:"contact_email"`
	BillingDetails  *string   `json:"billing_details,omitempty" db:"billing_details"` // JSONB stored as string
	Status          string    `json:"status" db:"status"` // 'active', 'pending', 'inactive', 'rejected'
	
	// General purpose fields (provider-agnostic)
	InternalNotes              *string `json:"internal_notes,omitempty" db:"internal_notes"`
	DefaultCurrencyID          *string `json:"default_currency_id,omitempty" db:"default_currency_id"`
	PlatformName               *string `json:"platform_name,omitempty" db:"platform_name"`
	PlatformURL                *string `json:"platform_url,omitempty" db:"platform_url"`
	PlatformUsername           *string `json:"platform_username,omitempty" db:"platform_username"`
	AccountingContactEmail     *string `json:"accounting_contact_email,omitempty" db:"accounting_contact_email"`
	OfferIDMacro               *string `json:"offer_id_macro,omitempty" db:"offer_id_macro"`
	AffiliateIDMacro           *string `json:"affiliate_id_macro,omitempty" db:"affiliate_id_macro"`
	AttributionMethod          *string `json:"attribution_method,omitempty" db:"attribution_method"`
	EmailAttributionMethod     *string `json:"email_attribution_method,omitempty" db:"email_attribution_method"`
	AttributionPriority        *string `json:"attribution_priority,omitempty" db:"attribution_priority"`
	ReportingTimezoneID        *int32  `json:"reporting_timezone_id,omitempty" db:"reporting_timezone_id"`
	IsExposePublisherReporting *bool   `json:"is_expose_publisher_reporting,omitempty" db:"is_expose_publisher_reporting"`
	
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
	
	// Provider-specific data (stored as JSONB) - contains all Everflow-specific fields
	ProviderData         *string   `json:"provider_data,omitempty" db:"provider_data"` // JSONB for provider-specific fields
	
	// Synchronization metadata
	SyncStatus           *string   `json:"sync_status,omitempty" db:"sync_status"` // 'synced', 'pending', 'error'
	LastSyncAt           *time.Time `json:"last_sync_at,omitempty" db:"last_sync_at"`
	SyncError            *string   `json:"sync_error,omitempty" db:"sync_error"`
	
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// AdvertiserDiscrepancy represents a discrepancy between local and provider data
type AdvertiserDiscrepancy struct {
	Field         string      `json:"field"`
	LocalValue    interface{} `json:"local_value"`
	ProviderValue interface{} `json:"provider_value"`
	Severity      string      `json:"severity"` // 'low', 'medium', 'high', 'critical'
}

// AdvertiserWithProviderData represents an advertiser with provider comparison data
type AdvertiserWithProviderData struct {
	*Advertiser
	ProviderData   interface{}                `json:"provider_data,omitempty"`
	Discrepancies  []AdvertiserDiscrepancy    `json:"discrepancies,omitempty"`
	SyncStatus     string                     `json:"sync_status"` // 'synced', 'out_of_sync', 'not_synced', 'error'
}

// AdvertiserWithEverflowData is an alias for backward compatibility
type AdvertiserWithEverflowData = AdvertiserWithProviderData

// EverflowAdvertiserProviderData represents Everflow-specific data stored in ProviderData field
type EverflowAdvertiserProviderData struct {
	// Everflow-specific fields only (general purpose fields moved to main Advertiser model)
	NetworkAdvertiserID         *int32          `json:"network_advertiser_id,omitempty"`
	NetworkEmployeeID           *int32          `json:"network_employee_id,omitempty"`
	IsExposePublisherReporting  *bool           `json:"is_expose_publisher_reporting,omitempty"`
	
	// Everflow-specific structured data
	Settings                    *interface{}    `json:"settings,omitempty"`
	ReportingData               *interface{}    `json:"reporting_data,omitempty"`
	
	// Additional fields for extensibility
	AdditionalFields            map[string]interface{} `json:"additional_fields,omitempty"`
}