package domain

import (
	"time"
)

// Campaign represents a campaign entity with full Everflow offer support
type Campaign struct {
	CampaignID     int64      `json:"campaign_id" db:"campaign_id"`
	OrganizationID int64      `json:"organization_id" db:"organization_id"`
	AdvertiserID   int64      `json:"advertiser_id" db:"advertiser_id"`
	Name           string     `json:"name" db:"name"`
	Description    *string    `json:"description,omitempty" db:"description"`
	Status         string     `json:"status" db:"status"` // 'draft', 'active', 'paused', 'archived'
	StartDate      *time.Time `json:"start_date,omitempty" db:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty" db:"end_date"`
	
	// Offer-specific fields for Everflow integration
	NetworkAdvertiserID *int32  `json:"network_advertiser_id,omitempty" db:"network_advertiser_id"`
	DestinationURL      *string `json:"destination_url,omitempty" db:"destination_url"`
	ThumbnailURL        *string `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	PreviewURL          *string `json:"preview_url,omitempty" db:"preview_url"`
	Visibility          *string `json:"visibility,omitempty" db:"visibility"`                   // 'public', 'require_approval', 'private'
	CurrencyID          *string `json:"currency_id,omitempty" db:"currency_id"`                 // 'USD', 'EUR', etc.
	ConversionMethod    *string `json:"conversion_method,omitempty" db:"conversion_method"`     // 'server_postback', 'pixel', etc.
	SessionDefinition   *string `json:"session_definition,omitempty" db:"session_definition"`   // 'cookie', 'ip', 'fingerprint'
	SessionDuration     *int32  `json:"session_duration,omitempty" db:"session_duration"`       // in hours
	CapsTimezoneID      *int32  `json:"caps_timezone_id,omitempty" db:"caps_timezone_id"`
	ProjectID           *string `json:"project_id,omitempty" db:"project_id"`
	DateLiveUntil       *time.Time `json:"date_live_until,omitempty" db:"date_live_until"`
	HTMLDescription     *string `json:"html_description,omitempty" db:"html_description"`
	InternalNotes       *string `json:"internal_notes,omitempty" db:"internal_notes"`
	TermsAndConditions  *string `json:"terms_and_conditions,omitempty" db:"terms_and_conditions"`
	IsUsingExplicitTermsAndConditions *bool `json:"is_using_explicit_terms_and_conditions,omitempty" db:"is_using_explicit_terms_and_conditions"`
	IsForceTermsAndConditions *bool `json:"is_force_terms_and_conditions,omitempty" db:"is_force_terms_and_conditions"`
	IsWhitelistCheckEnabled *bool `json:"is_whitelist_check_enabled,omitempty" db:"is_whitelist_check_enabled"`
	IsViewThroughEnabled *bool `json:"is_view_through_enabled,omitempty" db:"is_view_through_enabled"`
	ServerSideURL       *string `json:"server_side_url,omitempty" db:"server_side_url"`
	ViewThroughDestinationURL *string `json:"view_through_destination_url,omitempty" db:"view_through_destination_url"`
	IsDescriptionPlainText *bool `json:"is_description_plain_text,omitempty" db:"is_description_plain_text"`
	IsUseDirectLinking  *bool `json:"is_use_direct_linking,omitempty" db:"is_use_direct_linking"`
	AppIdentifier       *string `json:"app_identifier,omitempty" db:"app_identifier"`
	
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
	
	// Everflow tracking fields
	EncodedValue    *string `json:"encoded_value,omitempty" db:"encoded_value"`
	TodayClicks     *int    `json:"today_clicks,omitempty" db:"today_clicks"`
	TodayRevenue    *string `json:"today_revenue,omitempty" db:"today_revenue"`
	TimeCreated     *int    `json:"time_created,omitempty" db:"time_created"`
	TimeSaved       *int    `json:"time_saved,omitempty" db:"time_saved"`
	
	// Payout and revenue configuration
	PayoutType     *string  `json:"payout_type,omitempty" db:"payout_type"`         // 'cpa', 'cpc', 'cpm', etc.
	PayoutAmount   *float64 `json:"payout_amount,omitempty" db:"payout_amount"`
	RevenueType    *string  `json:"revenue_type,omitempty" db:"revenue_type"`       // 'rpa', 'rpc', 'rpm', etc.
	RevenueAmount  *float64 `json:"revenue_amount,omitempty" db:"revenue_amount"`
	
	// Additional configuration stored as JSON
	OfferConfig *string `json:"offer_config,omitempty" db:"offer_config"` // JSONB stored as string for additional Everflow-specific config
	
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CampaignProviderOffer represents an offer for a campaign in a provider
type CampaignProviderOffer struct {
	CampaignProviderOfferID int64      `json:"campaign_provider_offer_id" db:"campaign_provider_offer_id"`
	CampaignID              int64      `json:"campaign_id" db:"campaign_id"`
	ProviderType            string     `json:"provider_type" db:"provider_type"` // 'everflow' for MVP
	ProviderOfferRef        *string    `json:"provider_offer_ref,omitempty" db:"provider_offer_ref"` // Provider's Offer ID (e.g., Everflow's network_offer_id)
	ProviderOfferConfig     *string    `json:"provider_offer_config,omitempty" db:"provider_offer_config"` // JSONB stored as string
	IsActiveOnProvider      bool       `json:"is_active_on_provider" db:"is_active_on_provider"`
	LastSyncedAt            *time.Time `json:"last_synced_at,omitempty" db:"last_synced_at"`
	CreatedAt               time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at" db:"updated_at"`
}

// CampaignProviderMapping represents a mapping between a campaign and a provider
type CampaignProviderMapping struct {
	MappingID            int64     `json:"mapping_id" db:"mapping_id"`
	CampaignID           int64     `json:"campaign_id" db:"campaign_id"`
	ProviderType         string    `json:"provider_type" db:"provider_type"` // 'everflow' for MVP
	ProviderCampaignID   *string   `json:"provider_campaign_id,omitempty" db:"provider_campaign_id"` // Provider's Campaign/Offer ID
	APICredentials       *string   `json:"api_credentials,omitempty" db:"api_credentials"` // JSONB stored as string for API keys/tokens
	ProviderConfig       *string   `json:"provider_config,omitempty" db:"provider_config"` // JSONB stored as string
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}