package domain

import (
	"encoding/json"
	"time"
)

// Campaign represents a clean campaign entity following clean architecture principles
type Campaign struct {
	CampaignID     int64      `json:"campaign_id" db:"campaign_id"`
	OrganizationID int64      `json:"organization_id" db:"organization_id"`
	AdvertiserID   int64      `json:"advertiser_id" db:"advertiser_id"`
	Name           string     `json:"name" db:"name"`
	Description    *string    `json:"description,omitempty" db:"description"`
	Status         string     `json:"status" db:"status"` // 'draft', 'active', 'paused', 'archived'
	StartDate      *time.Time `json:"start_date,omitempty" db:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty" db:"end_date"`
	InternalNotes  *string    `json:"internal_notes,omitempty" db:"internal_notes"`
	
	// Core campaign fields (provider-agnostic)
	DestinationURL      *string `json:"destination_url,omitempty" db:"destination_url"`
	ThumbnailURL        *string `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	PreviewURL          *string `json:"preview_url,omitempty" db:"preview_url"`
	Visibility          *string `json:"visibility,omitempty" db:"visibility"`                   // 'public', 'require_approval', 'private'
	CurrencyID          *string `json:"currency_id,omitempty" db:"currency_id"`                 // 'USD', 'EUR', etc.
	ConversionMethod    *string `json:"conversion_method,omitempty" db:"conversion_method"`     // 'server_postback', 'pixel', etc.
	SessionDefinition   *string `json:"session_definition,omitempty" db:"session_definition"`   // 'cookie', 'ip', 'fingerprint'
	SessionDuration     *int32  `json:"session_duration,omitempty" db:"session_duration"`       // in hours
	TermsAndConditions  *string `json:"terms_and_conditions,omitempty" db:"terms_and_conditions"`
	
	// Caps and limits
	IsCapsEnabled         *bool `json:"is_caps_enabled,omitempty" db:"is_caps_enabled"`
	DailyConversionCap    *int  `json:"daily_conversion_cap,omitempty" db:"daily_conversion_cap"`
	WeeklyConversionCap   *int  `json:"weekly_conversion_cap,omitempty" db:"weekly_conversion_cap"`
	MonthlyConversionCap  *int  `json:"monthly_conversion_cap,omitempty" db:"monthly_conversion_cap"`
	GlobalConversionCap   *int  `json:"global_conversion_cap,omitempty" db:"global_conversion_cap"`
	DailyClickCap         *int  `json:"daily_click_cap,omitempty" db:"daily_click_cap"`
	WeeklyClickCap        *int  `json:"weekly_click_cap,omitempty" db:"weekly_click_cap"`
	MonthlyClickCap       *int  `json:"monthly_click_cap,omitempty" db:"monthly_click_cap"`
	GlobalClickCap        *int  `json:"global_click_cap,omitempty" db:"global_click_cap"`
	
	// Payout and revenue configuration
	BillingModel      *string  `json:"billing_model,omitempty" db:"billing_model"`           // 'click' or 'conversion'
	PayoutStructure   *string  `json:"payout_structure,omitempty" db:"payout_structure"`     // 'fixed' or 'percentage' (only for conversion)
	PayoutAmount      *float64 `json:"payout_amount,omitempty" db:"payout_amount"`           // Fixed amount or percentage value
	RevenueStructure  *string  `json:"revenue_structure,omitempty" db:"revenue_structure"`   // 'fixed' or 'percentage' (only for conversion)
	RevenueAmount     *float64 `json:"revenue_amount,omitempty" db:"revenue_amount"`         // Fixed amount or percentage value
	
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CampaignProviderMapping represents a mapping between a campaign and a provider following clean architecture
type CampaignProviderMapping struct {
	MappingID            int64      `json:"mapping_id" db:"mapping_id"`
	CampaignID           int64      `json:"campaign_id" db:"campaign_id"`
	ProviderType         string     `json:"provider_type" db:"provider_type"` // 'everflow' for MVP
	ProviderCampaignID   *string    `json:"provider_campaign_id,omitempty" db:"provider_campaign_id"` // Provider's Campaign ID
	
	// Provider-specific data stored as JSONB
	ProviderData         *string    `json:"provider_data,omitempty" db:"provider_data"`
	
	// Synchronization metadata
	SyncStatus           *string    `json:"sync_status,omitempty" db:"sync_status"`
	LastSyncAt           *time.Time `json:"last_sync_at,omitempty" db:"last_sync_at"`
	SyncError            *string    `json:"sync_error,omitempty" db:"sync_error"`
	
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}

// EverflowCampaignProviderData represents Everflow-specific campaign data
type EverflowCampaignProviderData struct {
	NetworkCampaignID                 *int32                 `json:"network_campaign_id,omitempty"`
	NetworkAdvertiserID               *int32                 `json:"network_advertiser_id,omitempty"`
	CapsTimezoneID                    *int32                 `json:"caps_timezone_id,omitempty"`
	ProjectID                         *string                `json:"project_id,omitempty"`
	DateLiveUntil                     *time.Time             `json:"date_live_until,omitempty"`
	HTMLDescription                   *string                `json:"html_description,omitempty"`
	IsUsingExplicitTermsAndConditions *bool                  `json:"is_using_explicit_terms_and_conditions,omitempty"`
	IsForceTermsAndConditions         *bool                  `json:"is_force_terms_and_conditions,omitempty"`
	IsWhitelistCheckEnabled           *bool                  `json:"is_whitelist_check_enabled,omitempty"`
	IsViewThroughEnabled              *bool                  `json:"is_view_through_enabled,omitempty"`
	ServerSideURL                     *string                `json:"server_side_url,omitempty"`
	ViewThroughDestinationURL         *string                `json:"view_through_destination_url,omitempty"`
	IsDescriptionPlainText            *bool                  `json:"is_description_plain_text,omitempty"`
	IsUseDirectLinking                *bool                  `json:"is_use_direct_linking,omitempty"`
	AppIdentifier                     *string                `json:"app_identifier,omitempty"`
	EncodedValue                      *string                `json:"encoded_value,omitempty"`
	TodayClicks                       *int                   `json:"today_clicks,omitempty"`
	TodayRevenue                      *string                `json:"today_revenue,omitempty"`
	TimeCreated                       *int                   `json:"time_created,omitempty"`
	TimeSaved                         *int                   `json:"time_saved,omitempty"`
	AdditionalFields                  map[string]interface{} `json:"additional_fields,omitempty"`
}

// ToJSON converts EverflowCampaignProviderData to JSON string
func (e *EverflowCampaignProviderData) ToJSON() (string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON populates EverflowCampaignProviderData from JSON string
func (e *EverflowCampaignProviderData) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), e)
}