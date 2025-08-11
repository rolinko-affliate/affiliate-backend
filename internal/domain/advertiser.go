package domain

import (
	"fmt"
	"time"
)

// Advertiser represents an advertiser entity (clean domain model)
type Advertiser struct {
	AdvertiserID   int64           `json:"advertiser_id" db:"advertiser_id"`
	OrganizationID int64           `json:"organization_id" db:"organization_id"`
	Name           string          `json:"name" db:"name"`
	ContactEmail   *string         `json:"contact_email,omitempty" db:"contact_email"`
	BillingDetails *BillingDetails `json:"billing_details,omitempty" db:"billing_details"` // JSONB
	Status         string          `json:"status" db:"status"`                             // 'active', 'pending', 'inactive', 'rejected'

	// General purpose fields (provider-agnostic)
	InternalNotes          *string `json:"internal_notes,omitempty" db:"internal_notes"`
	DefaultCurrencyID      *string `json:"default_currency_id,omitempty" db:"default_currency_id"`
	PlatformName           *string `json:"platform_name,omitempty" db:"platform_name"`
	PlatformURL            *string `json:"platform_url,omitempty" db:"platform_url"`
	PlatformUsername       *string `json:"platform_username,omitempty" db:"platform_username"`
	AccountingContactEmail *string `json:"accounting_contact_email,omitempty" db:"accounting_contact_email"`
	OfferIDMacro           *string `json:"offer_id_macro,omitempty" db:"offer_id_macro"`
	AffiliateIDMacro       *string `json:"affiliate_id_macro,omitempty" db:"affiliate_id_macro"`
	AttributionMethod      *string `json:"attribution_method,omitempty" db:"attribution_method"`
	EmailAttributionMethod *string `json:"email_attribution_method,omitempty" db:"email_attribution_method"`
	AttributionPriority    *string `json:"attribution_priority,omitempty" db:"attribution_priority"`
	ReportingTimezoneID    *int32  `json:"reporting_timezone_id,omitempty" db:"reporting_timezone_id"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Validate validates the advertiser data
func (a *Advertiser) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("advertiser name is required")
	}
	if a.OrganizationID <= 0 {
		return fmt.Errorf("valid organization ID is required")
	}
	if a.Status == "" {
		return fmt.Errorf("advertiser status is required")
	}
	validStatuses := map[string]bool{
		"active": true, "pending": true, "inactive": true, "rejected": true,
	}
	if !validStatuses[a.Status] {
		return fmt.Errorf("invalid status: %s", a.Status)
	}
	return nil
}

// HasBillingDetails returns true if billing details are present
func (a *Advertiser) HasBillingDetails() bool {
	return a.BillingDetails != nil && a.BillingDetails.HasData()
}

// AdvertiserProviderMapping represents a mapping between an advertiser and a provider
type AdvertiserProviderMapping struct {
	MappingID            int64   `json:"mapping_id" db:"mapping_id"`
	AdvertiserID         int64   `json:"advertiser_id" db:"advertiser_id"`
	ProviderType         string  `json:"provider_type" db:"provider_type"`                             // 'everflow' for MVP
	ProviderAdvertiserID *string `json:"provider_advertiser_id,omitempty" db:"provider_advertiser_id"` // e.g., Everflow's network_advertiser_id
	APICredentials       *string `json:"-" db:"api_credentials"`                                       // Exclude from JSON, JSONB stored as string, ENCRYPTED
	ProviderConfig       *string `json:"provider_config,omitempty" db:"provider_config"`               // JSONB stored as string

	// Provider-specific data (stored as JSONB) - contains all Everflow-specific fields
	ProviderData *string `json:"provider_data,omitempty" db:"provider_data"` // JSONB for provider-specific fields

	// Synchronization metadata
	SyncStatus *string    `json:"sync_status,omitempty" db:"sync_status"` // 'synced', 'pending', 'error'
	LastSyncAt *time.Time `json:"last_sync_at,omitempty" db:"last_sync_at"`
	SyncError  *string    `json:"sync_error,omitempty" db:"sync_error"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// IsSynced returns true if the mapping is successfully synced
func (apm *AdvertiserProviderMapping) IsSynced() bool {
	return apm.SyncStatus != nil && *apm.SyncStatus == "synced"
}

// HasSyncError returns true if there's a sync error
func (apm *AdvertiserProviderMapping) HasSyncError() bool {
	return apm.SyncStatus != nil && *apm.SyncStatus == "error"
}

// GetSyncStatusString returns the sync status as a string, defaulting to "not_synced"
func (apm *AdvertiserProviderMapping) GetSyncStatusString() string {
	if apm.SyncStatus == nil {
		return "not_synced"
	}
	return *apm.SyncStatus
}

// EverflowAdvertiserProviderData represents Everflow-specific data stored in ProviderData field
type EverflowAdvertiserProviderData struct {
	// Everflow-specific fields only (general purpose fields moved to main Advertiser model)
	NetworkAdvertiserID         *int32 `json:"network_advertiser_id,omitempty"`
	NetworkEmployeeID           *int32 `json:"network_employee_id,omitempty"`
	SalesManagerID              *int32 `json:"sales_manager_id,omitempty"`
	AddressID                   *int32 `json:"address_id,omitempty"`
	IsContactAddressEnabled     *bool  `json:"is_contact_address_enabled,omitempty"`
	VerificationToken           *string `json:"verification_token,omitempty"`
	
	// Everflow-specific structured data
	Users           *[]interface{} `json:"users,omitempty"`
	ContactAddress  *interface{}   `json:"contact_address,omitempty"`
	Settings        *interface{}   `json:"settings,omitempty"`
	Labels          *interface{}   `json:"labels,omitempty"`
	Billing         *interface{}   `json:"billing,omitempty"`
	
	// Additional fields for extensibility
	AdditionalFields map[string]interface{} `json:"additional_fields,omitempty"`
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
	ProviderData  interface{}             `json:"provider_data,omitempty"`
	Discrepancies []AdvertiserDiscrepancy `json:"discrepancies,omitempty"`
	SyncStatus    string                  `json:"sync_status"` // 'synced', 'out_of_sync', 'not_synced', 'error'
}

// AdvertiserExtraInfo represents additional information for an advertiser organization
type AdvertiserExtraInfo struct {
	ExtraInfoID    int64     `json:"extra_info_id" db:"extra_info_id"`
	OrganizationID int64     `json:"organization_id" db:"organization_id"`
	Website        *string   `json:"website,omitempty" db:"website"`
	WebsiteType    *string   `json:"website_type,omitempty" db:"website_type"` // 'shopify', 'amazon', 'shopline', 'tiktok_shop'
	CompanySize    *string   `json:"company_size,omitempty" db:"company_size"`  // 'startup', 'small', 'medium', 'large', 'enterprise'
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// Validate validates the advertiser extra info data
func (aei *AdvertiserExtraInfo) Validate() error {
	if aei.OrganizationID <= 0 {
		return fmt.Errorf("valid organization ID is required")
	}
	
	if aei.WebsiteType != nil {
		validTypes := map[string]bool{
			"shopify": true, "amazon": true, "shopline": true, "tiktok_shop": true,
		}
		if !validTypes[*aei.WebsiteType] {
			return fmt.Errorf("invalid website type: %s", *aei.WebsiteType)
		}
	}
	
	if aei.CompanySize != nil {
		validSizes := map[string]bool{
			"startup": true, "small": true, "medium": true, "large": true, "enterprise": true,
		}
		if !validSizes[*aei.CompanySize] {
			return fmt.Errorf("invalid company size: %s", *aei.CompanySize)
		}
	}
	
	return nil
}

// AdvertiserWithExtraInfo represents an advertiser with its extra information
type AdvertiserWithExtraInfo struct {
	*Advertiser
	ExtraInfo *AdvertiserExtraInfo `json:"extra_info,omitempty"`
}

// AdvertiserWithEverflowData is an alias for backward compatibility
type AdvertiserWithEverflowData = AdvertiserWithProviderData

// BulkSyncResult represents the result of a bulk sync operation
type BulkSyncResult struct {
	TotalProcessed int                    `json:"total_processed"`
	SuccessCount   int                    `json:"success_count"`
	FailureCount   int                    `json:"failure_count"`
	Successes      []BulkSyncItemResult   `json:"successes,omitempty"`
	Failures       []BulkSyncItemResult   `json:"failures,omitempty"`
	StartedAt      time.Time              `json:"started_at"`
	CompletedAt    time.Time              `json:"completed_at"`
}

// BulkSyncItemResult represents the result of syncing a single item
type BulkSyncItemResult struct {
	AdvertiserID         int64  `json:"advertiser_id"`
	AdvertiserName       string `json:"advertiser_name"`
	ProviderAdvertiserID string `json:"provider_advertiser_id,omitempty"`
	Error                string `json:"error,omitempty"`
}
