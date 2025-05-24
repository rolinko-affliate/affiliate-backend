package everflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var everflowAPIBaseURL = "https://api.eflow.team/v1" // Everflow API base URL

// Client represents an Everflow API client
type Client struct {
	httpClient *http.Client
	apiKey     string
}

// NewClient creates a new Everflow API client
func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second}, // Increased timeout for potentially slow API responses
		apiKey:     apiKey,
	}
}

// EverflowCreateAdvertiserRequest represents the request to create an advertiser in Everflow
type EverflowCreateAdvertiserRequest struct {
	Name                           string              `json:"name"`
	AccountStatus                  string              `json:"account_status"`      // "active", "inactive", "suspended"
	DefaultCurrencyID              string              `json:"default_currency_id"` // e.g., "USD"
	NetworkEmployeeID              *int                `json:"network_employee_id,omitempty"`
	SalesManagerID                 *int                `json:"sales_manager_id,omitempty"`
	InternalNotes                  *string             `json:"internal_notes,omitempty"`
	ReportingTimezoneID            *int                `json:"reporting_timezone_id,omitempty"`
	AttributionMethod              *string             `json:"attribution_method,omitempty"`       // "last_touch" or "first_touch"
	EmailAttributionMethod         *string             `json:"email_attribution_method,omitempty"` // "last_affiliate_attribution", "first_affiliate_attribution"
	AttributionPriority            *string             `json:"attribution_priority,omitempty"`     // "click", "coupon_code"
	VerificationToken              *string             `json:"verification_token,omitempty"`
	OfferIDMacro                   *string             `json:"offer_id_macro,omitempty"`
	AffiliateIDMacro               *string             `json:"affiliate_id_macro,omitempty"`
	IsContactAddressEnabled        *bool               `json:"is_contact_address_enabled,omitempty"`
	IsExposePublisherReportingData *bool               `json:"is_expose_publisher_reporting_data,omitempty"`
	PlatformName                   *string             `json:"platform_name,omitempty"`
	PlatformURL                    *string             `json:"platform_url,omitempty"`
	PlatformUsername               *string             `json:"platform_username,omitempty"`
	AccountingContactEmail         *string             `json:"accounting_contact_email,omitempty"`
	ContactAddress                 *AdvertiserAddress  `json:"contact_address,omitempty"`
	Labels                         []string            `json:"labels,omitempty"`
	Users                          []AdvertiserUser    `json:"users,omitempty"`
	Billing                        *AdvertiserBilling  `json:"billing,omitempty"`
	Settings                       *AdvertiserSettings `json:"settings,omitempty"`
}

// AdvertiserAddress represents an advertiser's contact address
type AdvertiserAddress struct {
	Address1      string  `json:"address_1"`
	Address2      *string `json:"address_2,omitempty"`
	City          string  `json:"city"`
	ZipPostalCode string  `json:"zip_postal_code"`
	CountryCode   string  `json:"country_code"` // From metadata API
	RegionCode    string  `json:"region_code"`  // From metadata API
}

// AdvertiserUser represents a user associated with an advertiser
type AdvertiserUser struct {
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Email           string  `json:"email"`
	AccountStatus   string  `json:"account_status"` // "active", "inactive", "suspended"
	InitialPassword *string `json:"initial_password,omitempty"`
	LanguageID      *int    `json:"language_id,omitempty"`
	TimezoneID      *int    `json:"timezone_id,omitempty"`
	CurrencyID      *string `json:"currency_id,omitempty"`
}

// AdvertiserBilling represents an advertiser's billing information
type AdvertiserBilling struct {
	BillingFrequency           string                 `json:"billing_frequency"` // "weekly", "bimonthly", "monthly", "two_months", "quarterly", "manual", "other"
	TaxID                      *string                `json:"tax_id,omitempty"`
	IsInvoiceCreationAuto      *bool                  `json:"is_invoice_creation_auto,omitempty"`
	InvoiceAmountThreshold     *float64               `json:"invoice_amount_threshold,omitempty"`
	AutoInvoiceStartDate       *string                `json:"auto_invoice_start_date,omitempty"`
	DefaultInvoiceIsHidden     *bool                  `json:"default_invoice_is_hidden,omitempty"`
	InvoiceGenerationDaysDelay *int                   `json:"invoice_generation_days_delay,omitempty"`
	DefaultPaymentTerms        *int                   `json:"default_payment_terms,omitempty"`
	Details                    map[string]interface{} `json:"details,omitempty"`
}

// AdvertiserSettings represents advertiser settings
type AdvertiserSettings struct {
	ExposedVariables map[string]bool `json:"exposed_variables,omitempty"`
}

// Advertiser represents a complete advertiser object from Everflow
type Advertiser struct {
	NetworkAdvertiserID            int64                      `json:"network_advertiser_id"`
	NetworkID                      int64                      `json:"network_id"`
	Name                           string                     `json:"name"`
	AccountStatus                  string                     `json:"account_status"`
	NetworkEmployeeID              int                        `json:"network_employee_id"`
	InternalNotes                  string                     `json:"internal_notes"`
	AddressID                      int64                      `json:"address_id"`
	IsContactAddressEnabled        bool                       `json:"is_contact_address_enabled"`
	SalesManagerID                 int                        `json:"sales_manager_id"`
	IsExposePublisherReportingData *bool                      `json:"is_expose_publisher_reporting_data"`
	DefaultCurrencyID              string                     `json:"default_currency_id"`
	PlatformName                   string                     `json:"platform_name"`
	PlatformURL                    string                     `json:"platform_url"`
	PlatformUsername               string                     `json:"platform_username"`
	ReportingTimezoneID            int                        `json:"reporting_timezone_id"`
	AccountingContactEmail         string                     `json:"accounting_contact_email"`
	VerificationToken              string                     `json:"verification_token"`
	OfferIDMacro                   string                     `json:"offer_id_macro"`
	AffiliateIDMacro               string                     `json:"affiliate_id_macro"`
	TimeCreated                    int64                      `json:"time_created"`
	TimeSaved                      int64                      `json:"time_saved"`
	AttributionMethod              string                     `json:"attribution_method"`
	EmailAttributionMethod         string                     `json:"email_attribution_method"`
	AttributionPriority            string                     `json:"attribution_priority"`
	ContactAddress                 *AdvertiserAddress         `json:"contact_address,omitempty"`
	Labels                         []string                   `json:"labels,omitempty"`
	Users                          []AdvertiserUser           `json:"users,omitempty"`
	Billing                        *AdvertiserBillingResponse `json:"billing,omitempty"`
	Settings                       *AdvertiserSettings        `json:"settings,omitempty"`
	Relationship                   *AdvertiserRelationship    `json:"relationship,omitempty"`
}

// AdvertiserBillingResponse represents billing information in responses
type AdvertiserBillingResponse struct {
	NetworkID                  int64   `json:"network_id"`
	NetworkAdvertiserID        int64   `json:"network_advertiser_id"`
	BillingFrequency           string  `json:"billing_frequency"`
	InvoiceAmountThreshold     float64 `json:"invoice_amount_threshold"`
	TaxID                      string  `json:"tax_id"`
	IsInvoiceCreationAuto      bool    `json:"is_invoice_creation_auto"`
	AutoInvoiceStartDate       string  `json:"auto_invoice_start_date"`
	DefaultInvoiceIsHidden     bool    `json:"default_invoice_is_hidden"`
	InvoiceGenerationDaysDelay int     `json:"invoice_generation_days_delay"`
	DefaultPaymentTerms        int     `json:"default_payment_terms"`
}

// AdvertiserRelationship represents relationship data for an advertiser
type AdvertiserRelationship struct {
	Labels          *AdvertiserLabels          `json:"labels,omitempty"`
	AccountManager  *AccountManager            `json:"account_manager,omitempty"`
	Integrations    *AdvertiserIntegrations    `json:"integrations,omitempty"`
	Reporting       *AdvertiserReporting       `json:"reporting,omitempty"`
	APIKeys         *AdvertiserAPIKeys         `json:"api_keys,omitempty"`
	APIWhitelistIPs *AdvertiserAPIWhitelistIPs `json:"api_whitelist_ips,omitempty"`
	Billing         *AdvertiserBillingResponse `json:"billing,omitempty"`
	Settings        *AdvertiserSettings        `json:"settings,omitempty"`
	SalesManager    *SalesManager              `json:"sale_manager,omitempty"`
}

// AdvertiserLabels represents labels associated with an advertiser
type AdvertiserLabels struct {
	Total   int      `json:"total"`
	Entries []string `json:"entries"`
}

// AccountManager represents account manager information
type AccountManager struct {
	FirstName                  string `json:"first_name"`
	LastName                   string `json:"last_name"`
	Email                      string `json:"email"`
	WorkPhone                  string `json:"work_phone"`
	CellPhone                  string `json:"cell_phone"`
	InstantMessagingID         int    `json:"instant_messaging_id"`
	InstantMessagingIdentifier string `json:"instant_messaging_identifier"`
}

// AdvertiserIntegrations represents integration information
type AdvertiserIntegrations struct {
	AdvertiserDemandPartner *interface{} `json:"advertiser_demand_partner"`
}

// Paging represents pagination information
type Paging struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalCount int `json:"total_count"`
}

// AdvertiserReporting represents reporting statistics for an advertiser
type AdvertiserReporting struct {
	Imp            int     `json:"imp"`
	TotalClick     int     `json:"total_click"`
	UniqueClick    int     `json:"unique_click"`
	InvalidClick   int     `json:"invalid_click"`
	DuplicateClick int     `json:"duplicate_click"`
	GrossClick     int     `json:"gross_click"`
	CTR            float64 `json:"ctr"`
	CV             int     `json:"cv"`
	InvalidCVScrub int     `json:"invalid_cv_scrub"`
	ViewThroughCV  int     `json:"view_through_cv"`
	TotalCV        int     `json:"total_cv"`
	Event          int     `json:"event"`
	CVR            float64 `json:"cvr"`
	EVR            float64 `json:"evr"`
	CPC            float64 `json:"cpc"`
	CPM            float64 `json:"cpm"`
	CPA            float64 `json:"cpa"`
	EPC            float64 `json:"epc"`
	RPC            float64 `json:"rpc"`
	RPA            float64 `json:"rpa"`
	RPM            float64 `json:"rpm"`
	Payout         float64 `json:"payout"`
	Revenue        float64 `json:"revenue"`
}

// AdvertiserAPIKeys represents API keys information
type AdvertiserAPIKeys struct {
	Total   int           `json:"total"`
	Entries []interface{} `json:"entries"`
}

// AdvertiserAPIWhitelistIPs represents API whitelist IPs information
type AdvertiserAPIWhitelistIPs struct {
	Total   int           `json:"total"`
	Entries []interface{} `json:"entries"`
}

// SalesManager represents sales manager information
type SalesManager struct {
	FirstName                  string `json:"first_name"`
	LastName                   string `json:"last_name"`
	Email                      string `json:"email"`
	WorkPhone                  string `json:"work_phone"`
	CellPhone                  string `json:"cell_phone"`
	InstantMessagingID         int    `json:"instant_messaging_id"`
	InstantMessagingIdentifier string `json:"instant_messaging_identifier"`
}

// EverflowListAdvertisersResponse represents the response from listing advertisers
type EverflowListAdvertisersResponse struct {
	Advertisers []Advertiser `json:"advertisers"`
	Paging      Paging       `json:"paging"`
}

// EverflowCreateAdvertiserResponse represents the response from creating an advertiser in Everflow
type EverflowCreateAdvertiserResponse Advertiser

// PayoutRevenueEntry represents a single payout/revenue entry
type PayoutRevenueEntry struct {
	EntryName                   *string `json:"entry_name,omitempty"`
	PayoutType                  string  `json:"payout_type"`                            // "cpc", "cpa", "cpm", "cps", "cpa_cps", "prv"
	PayoutAmount                *float64 `json:"payout_amount,omitempty"`
	PayoutPercentage            *int    `json:"payout_percentage,omitempty"`
	RevenueType                 string  `json:"revenue_type"`                           // "rpc", "rpa", "rpm", "rps", "rpa_rps"
	RevenueAmount               *float64 `json:"revenue_amount,omitempty"`
	RevenuePercentage           *int    `json:"revenue_percentage,omitempty"`
	IsDefault                   bool    `json:"is_default"`
	IsPrivate                   bool    `json:"is_private"`
	IsPostbackDisabled          *bool   `json:"is_postback_disabled,omitempty"`
	GlobalAdvertiserEventID     *int    `json:"global_advertiser_event_id,omitempty"`
	IsMustApproveConversion     *bool   `json:"is_must_approve_conversion,omitempty"`
	IsAllowDuplicateConversion  *bool   `json:"is_allow_duplicate_conversion,omitempty"`
}

// PayoutRevenue represents the payout revenue structure
type PayoutRevenue struct {
	Entries []PayoutRevenueEntry `json:"entries"`
}

// OfferRelationshipRemainingCaps represents remaining caps in relationship data
type OfferRelationshipRemainingCaps struct {
	RemainingDailyConversionCap   *int `json:"remaining_daily_conversion_cap,omitempty"`
	RemainingMonthlyConversionCap *int `json:"remaining_monthly_conversion_cap,omitempty"`
	RemainingGlobalConversionCap  *int `json:"remaining_global_conversion_cap,omitempty"`
	RemainingDailyClickCap        *int `json:"remaining_daily_click_cap,omitempty"`
	RemainingWeeklyClickCap       *int `json:"remaining_weekly_click_cap,omitempty"`
	RemainingMonthlyClickCap      *int `json:"remaining_monthly_click_cap,omitempty"`
	RemainingGlobalClickCap       *int `json:"remaining_global_click_cap,omitempty"`
}

// OfferRelationshipRuleset represents ruleset in relationship data
type OfferRelationshipRuleset struct {
	NetworkRulesetID *int                   `json:"network_ruleset_id,omitempty"`
	Platforms        []map[string]interface{} `json:"platforms,omitempty"`
	DeviceTypes      []map[string]interface{} `json:"device_types,omitempty"`
	Countries        []map[string]interface{} `json:"countries,omitempty"`
}

// OfferRelationshipCategory represents category in relationship data
type OfferRelationshipCategory struct {
	NetworkCategoryID *int    `json:"network_category_id,omitempty"`
	Name              *string `json:"name,omitempty"`
}

// OfferRelationship represents relationship data for an offer
type OfferRelationship struct {
	RemainingCaps    *OfferRelationshipRemainingCaps `json:"remaining_caps,omitempty"`
	Ruleset          *OfferRelationshipRuleset       `json:"ruleset,omitempty"`
	Category         *OfferRelationshipCategory      `json:"category,omitempty"`
	PayoutRevenue    *PayoutRevenue                  `json:"payout_revenue,omitempty"`
	AdditionalData   map[string]interface{}          `json:"-"` // For any additional fields
}

// Offer represents a complete offer object from Everflow
type Offer struct {
	NetworkOfferID                              int64              `json:"network_offer_id"`
	NetworkID                                   int64              `json:"network_id"`
	NetworkAdvertiserID                         int64              `json:"network_advertiser_id"`
	NetworkOfferGroupID                         *int64             `json:"network_offer_group_id,omitempty"`
	NetworkCategoryID                           *int64             `json:"network_category_id,omitempty"`
	Name                                        string             `json:"name"`
	ThumbnailURL                                *string            `json:"thumbnail_url,omitempty"`
	InternalNotes                               *string            `json:"internal_notes,omitempty"`
	DestinationURL                              string             `json:"destination_url"`
	ServerSideURL                               *string            `json:"server_side_url,omitempty"`
	IsViewThroughEnabled                        *bool              `json:"is_view_through_enabled,omitempty"`
	ViewThroughDestinationURL                   *string            `json:"view_through_destination_url,omitempty"`
	PreviewURL                                  *string            `json:"preview_url,omitempty"`
	OfferStatus                                 string             `json:"offer_status"`
	Visibility                                  *string            `json:"visibility,omitempty"`
	CurrencyID                                  *string            `json:"currency_id,omitempty"`
	CapsTimezoneID                              *int               `json:"caps_timezone_id,omitempty"`
	ProjectID                                   *string            `json:"project_id,omitempty"`
	DateLiveUntil                               *string            `json:"date_live_until,omitempty"`
	HTMLDescription                             *string            `json:"html_description,omitempty"`
	IsUsingExplicitTermsAndConditions           *bool              `json:"is_using_explicit_terms_and_conditions,omitempty"`
	TermsAndConditions                          *string            `json:"terms_and_conditions,omitempty"`
	IsForceTermsAndConditions                   *bool              `json:"is_force_terms_and_conditions,omitempty"`
	IsCapsEnabled                               *bool              `json:"is_caps_enabled,omitempty"`
	DailyConversionCap                          *int               `json:"daily_conversion_cap,omitempty"`
	WeeklyConversionCap                         *int               `json:"weekly_conversion_cap,omitempty"`
	MonthlyConversionCap                        *int               `json:"monthly_conversion_cap,omitempty"`
	GlobalConversionCap                         *int               `json:"global_conversion_cap,omitempty"`
	DailyPayoutCap                              *int               `json:"daily_payout_cap,omitempty"`
	WeeklyPayoutCap                             *int               `json:"weekly_payout_cap,omitempty"`
	MonthlyPayoutCap                            *int               `json:"monthly_payout_cap,omitempty"`
	GlobalPayoutCap                             *int               `json:"global_payout_cap,omitempty"`
	DailyRevenueCap                             *int               `json:"daily_revenue_cap,omitempty"`
	WeeklyRevenueCap                            *int               `json:"weekly_revenue_cap,omitempty"`
	MonthlyRevenueCap                           *int               `json:"monthly_revenue_cap,omitempty"`
	GlobalRevenueCap                            *int               `json:"global_revenue_cap,omitempty"`
	DailyClickCap                               *int               `json:"daily_click_cap,omitempty"`
	WeeklyClickCap                              *int               `json:"weekly_click_cap,omitempty"`
	MonthlyClickCap                             *int               `json:"monthly_click_cap,omitempty"`
	GlobalClickCap                              *int               `json:"global_click_cap,omitempty"`
	RedirectMode                                *string            `json:"redirect_mode,omitempty"`
	IsUsingSuppressionList                      *bool              `json:"is_using_suppression_list,omitempty"`
	SuppressionListID                           *int               `json:"suppression_list_id,omitempty"`
	IsDuplicateFilterEnabled                    *bool              `json:"is_duplicate_filter_enabled,omitempty"`
	DuplicateFilterTargetingAction              *string            `json:"duplicate_filter_targeting_action,omitempty"`
	NetworkTrackingDomainID                     *int               `json:"network_tracking_domain_id,omitempty"`
	IsUseSecureLink                             *bool              `json:"is_use_secure_link,omitempty"`
	IsAllowDeepLink                             *bool              `json:"is_allow_deep_link,omitempty"`
	IsSessionTrackingEnabled                    *bool              `json:"is_session_tracking_enabled,omitempty"`
	SessionTrackingLifespanHour                 *int               `json:"session_tracking_lifespan_hour,omitempty"`
	SessionTrackingMinimumLifespanSecond        *int               `json:"session_tracking_minimum_lifespan_second,omitempty"`
	IsViewThroughSessionTrackingEnabled         *bool              `json:"is_view_through_session_tracking_enabled,omitempty"`
	ViewThroughSessionTrackingLifespanMinute    *int               `json:"view_through_session_tracking_lifespan_minute,omitempty"`
	ViewThroughSessionTrackingMinimalLifespanSecond *int           `json:"view_through_session_tracking_minimal_lifespan_second,omitempty"`
	IsBlockAlreadyConverted                     *bool              `json:"is_block_already_converted,omitempty"`
	AlreadyConvertedAction                      *string            `json:"already_converted_action,omitempty"`
	ConversionMethod                            *string            `json:"conversion_method,omitempty"`
	IsWhitelistCheckEnabled                     *bool              `json:"is_whitelist_check_enabled,omitempty"`
	IsUseScrubRate                              *bool              `json:"is_use_scrub_rate,omitempty"`
	ScrubRateStatus                             *string            `json:"scrub_rate_status,omitempty"`
	ScrubRatePercentage                         *int               `json:"scrub_rate_percentage,omitempty"`
	SessionDefinition                           *string            `json:"session_definition,omitempty"`
	SessionDuration                             *int               `json:"session_duration,omitempty"`
	AppIdentifier                               *string            `json:"app_identifier,omitempty"`
	IsDescriptionPlainText                      *bool              `json:"is_description_plain_text,omitempty"`
	IsUseDirectLinking                          *bool              `json:"is_use_direct_linking,omitempty"`
	EncodedValue                                *string            `json:"encoded_value,omitempty"`
	TodayClicks                                 *int               `json:"today_clicks,omitempty"`
	TodayRevenue                                *string            `json:"today_revenue,omitempty"`
	TimeCreated                                 *int64             `json:"time_created,omitempty"`
	TimeSaved                                   *int64             `json:"time_saved,omitempty"`
	NetworkAdvertiserName                       *string            `json:"network_advertiser_name,omitempty"`
	Category                                    *string            `json:"category,omitempty"`
	Payout                                      *string            `json:"payout,omitempty"`
	Revenue                                     *string            `json:"revenue,omitempty"`
	Labels                                      *string            `json:"labels,omitempty"`
	Channels                                    *string            `json:"channels,omitempty"`
	PayoutRevenue                               *PayoutRevenue     `json:"payout_revenue,omitempty"`
	Relationship                                *OfferRelationship `json:"relationship,omitempty"`
}

// OfferInput represents the input for creating or updating an offer
type OfferInput struct {
	NetworkAdvertiserID                         int64              `json:"network_advertiser_id"`
	NetworkOfferGroupID                         *int64             `json:"network_offer_group_id,omitempty"`
	Name                                        string             `json:"name"`
	ThumbnailURL                                *string            `json:"thumbnail_url,omitempty"`
	NetworkCategoryID                           *int64             `json:"network_category_id,omitempty"`
	InternalNotes                               *string            `json:"internal_notes,omitempty"`
	DestinationURL                              string             `json:"destination_url"`
	ServerSideURL                               *string            `json:"server_side_url,omitempty"`
	IsViewThroughEnabled                        *bool              `json:"is_view_through_enabled,omitempty"`
	ViewThroughDestinationURL                   *string            `json:"view_through_destination_url,omitempty"`
	PreviewURL                                  *string            `json:"preview_url,omitempty"`
	OfferStatus                                 string             `json:"offer_status"`
	Visibility                                  *string            `json:"visibility,omitempty"`
	CurrencyID                                  *string            `json:"currency_id,omitempty"`
	CapsTimezoneID                              *int               `json:"caps_timezone_id,omitempty"`
	ProjectID                                   *string            `json:"project_id,omitempty"`
	DateLiveUntil                               *string            `json:"date_live_until,omitempty"`
	HTMLDescription                             *string            `json:"html_description,omitempty"`
	IsUsingExplicitTermsAndConditions           *bool              `json:"is_using_explicit_terms_and_conditions,omitempty"`
	TermsAndConditions                          *string            `json:"terms_and_conditions,omitempty"`
	IsForceTermsAndConditions                   *bool              `json:"is_force_terms_and_conditions,omitempty"`
	IsCapsEnabled                               *bool              `json:"is_caps_enabled,omitempty"`
	DailyConversionCap                          *int               `json:"daily_conversion_cap,omitempty"`
	WeeklyConversionCap                         *int               `json:"weekly_conversion_cap,omitempty"`
	MonthlyConversionCap                        *int               `json:"monthly_conversion_cap,omitempty"`
	GlobalConversionCap                         *int               `json:"global_conversion_cap,omitempty"`
	DailyPayoutCap                              *int               `json:"daily_payout_cap,omitempty"`
	WeeklyPayoutCap                             *int               `json:"weekly_payout_cap,omitempty"`
	MonthlyPayoutCap                            *int               `json:"monthly_payout_cap,omitempty"`
	GlobalPayoutCap                             *int               `json:"global_payout_cap,omitempty"`
	DailyRevenueCap                             *int               `json:"daily_revenue_cap,omitempty"`
	WeeklyRevenueCap                            *int               `json:"weekly_revenue_cap,omitempty"`
	MonthlyRevenueCap                           *int               `json:"monthly_revenue_cap,omitempty"`
	GlobalRevenueCap                            *int               `json:"global_revenue_cap,omitempty"`
	DailyClickCap                               *int               `json:"daily_click_cap,omitempty"`
	WeeklyClickCap                              *int               `json:"weekly_click_cap,omitempty"`
	MonthlyClickCap                             *int               `json:"monthly_click_cap,omitempty"`
	GlobalClickCap                              *int               `json:"global_click_cap,omitempty"`
	RedirectMode                                *string            `json:"redirect_mode,omitempty"`
	IsUsingSuppressionList                      *bool              `json:"is_using_suppression_list,omitempty"`
	SuppressionListID                           *int               `json:"suppression_list_id,omitempty"`
	IsDuplicateFilterEnabled                    *bool              `json:"is_duplicate_filter_enabled,omitempty"`
	DuplicateFilterTargetingAction              *string            `json:"duplicate_filter_targeting_action,omitempty"`
	NetworkTrackingDomainID                     *int               `json:"network_tracking_domain_id,omitempty"`
	IsUseSecureLink                             *bool              `json:"is_use_secure_link,omitempty"`
	IsAllowDeepLink                             *bool              `json:"is_allow_deep_link,omitempty"`
	IsSessionTrackingEnabled                    *bool              `json:"is_session_tracking_enabled,omitempty"`
	SessionTrackingLifespanHour                 *int               `json:"session_tracking_lifespan_hour,omitempty"`
	SessionTrackingMinimumLifespanSecond        *int               `json:"session_tracking_minimum_lifespan_second,omitempty"`
	IsViewThroughSessionTrackingEnabled         *bool              `json:"is_view_through_session_tracking_enabled,omitempty"`
	ViewThroughSessionTrackingLifespanMinute    *int               `json:"view_through_session_tracking_lifespan_minute,omitempty"`
	ViewThroughSessionTrackingMinimalLifespanSecond *int           `json:"view_through_session_tracking_minimal_lifespan_second,omitempty"`
	IsBlockAlreadyConverted                     *bool              `json:"is_block_already_converted,omitempty"`
	AlreadyConvertedAction                      *string            `json:"already_converted_action,omitempty"`
	ConversionMethod                            *string            `json:"conversion_method,omitempty"`
	IsWhitelistCheckEnabled                     *bool              `json:"is_whitelist_check_enabled,omitempty"`
	IsUseScrubRate                              *bool              `json:"is_use_scrub_rate,omitempty"`
	ScrubRateStatus                             *string            `json:"scrub_rate_status,omitempty"`
	ScrubRatePercentage                         *int               `json:"scrub_rate_percentage,omitempty"`
	SessionDefinition                           *string            `json:"session_definition,omitempty"`
	SessionDuration                             *int               `json:"session_duration,omitempty"`
	AppIdentifier                               *string            `json:"app_identifier,omitempty"`
	IsDescriptionPlainText                      *bool              `json:"is_description_plain_text,omitempty"`
	IsUseDirectLinking                          *bool              `json:"is_use_direct_linking,omitempty"`
	EncodedValue                                *string            `json:"encoded_value,omitempty"`
	Meta                                        *OfferMeta         `json:"meta,omitempty"`
	PayoutRevenue                               *PayoutRevenue     `json:"payout_revenue,omitempty"`
}

// OfferMeta represents metadata for an offer
type OfferMeta struct {
	AdvertiserCampaignName *string `json:"advertiser_campaign_name,omitempty"`
}

// OffersTableSearchTerm represents a search term for the offers table
type OffersTableSearchTerm struct {
	SearchType string `json:"search_type"` // "name", "encoded_value", "advertiser_name"
	Value      string `json:"value"`
}

// OffersTableRequest represents the request for the offers table endpoint
type OffersTableRequest struct {
	SearchTerms []OffersTableSearchTerm    `json:"search_terms,omitempty"`
	Filters     map[string]interface{}     `json:"filters,omitempty"`
}

// OffersTableOptions represents options for the offers table request
type OffersTableOptions struct {
	Page          *int     `json:"page,omitempty"`
	PageSize      *int     `json:"page_size,omitempty"`
	Relationships []string `json:"relationships,omitempty"`
}

// OffersTableResponse represents the response from the offers table endpoint
type OffersTableResponse struct {
	Offers []Offer `json:"offers"`
}

// ListAdvertisersOptions represents options for listing advertisers
type ListAdvertisersOptions struct {
	Page     *int `json:"page,omitempty"`
	PageSize *int `json:"page_size,omitempty"`
}

// ListAdvertisers retrieves all advertisers with optional filters
func (c *Client) ListAdvertisers(ctx context.Context, opts *ListAdvertisersOptions) (*EverflowListAdvertisersResponse, error) {
	reqURL := everflowAPIBaseURL + "/networks/advertisers"

	// Add query parameters if options are provided
	if opts != nil {
		params := url.Values{}
		if opts.Page != nil {
			params.Add("page", strconv.Itoa(*opts.Page))
		}
		if opts.PageSize != nil {
			params.Add("page_size", strconv.Itoa(*opts.PageSize))
		}
		if len(params) > 0 {
			reqURL += "?" + params.Encode()
		}
	}

	httpReq, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("X-Eflow-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to Everflow: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read error body for more details
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return nil, fmt.Errorf("Everflow API request failed with status %d: %v", resp.StatusCode, errorResp)
		}
		return nil, fmt.Errorf("Everflow API request failed with status %d", resp.StatusCode)
	}

	var listResp EverflowListAdvertisersResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode Everflow list advertisers response: %w", err)
	}

	return &listResp, nil
}

// CreateAdvertiser creates a new advertiser in Everflow
func (c *Client) CreateAdvertiser(ctx context.Context, req EverflowCreateAdvertiserRequest) (*EverflowCreateAdvertiserResponse, error) {
	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create advertiser request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", everflowAPIBaseURL+"/networks/advertisers", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Eflow-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to Everflow: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		// Read error body for more details
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return nil, fmt.Errorf("Everflow API request failed with status %d: %v", resp.StatusCode, errorResp)
		}
		return nil, fmt.Errorf("Everflow API request failed with status %d", resp.StatusCode)
	}

	var createResp EverflowCreateAdvertiserResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return nil, fmt.Errorf("failed to decode Everflow create advertiser response: %w", err)
	}

	return &createResp, nil
}

// GetAdvertiserOptions represents options for getting a single advertiser
type GetAdvertiserOptions struct {
	Relationships []string `json:"relationships,omitempty"`
}

// Valid relationship values for GetAdvertiser
const (
	RelationshipReporting     = "reporting"
	RelationshipLabels        = "labels"
	RelationshipDemandPartner = "demand_partner"
	RelationshipBilling       = "billing"
	RelationshipIntegrations  = "integrations"
	RelationshipAPI           = "api"
)

// GetAdvertiser retrieves a single advertiser by ID from Everflow
func (c *Client) GetAdvertiser(ctx context.Context, networkAdvertiserID int64, opts *GetAdvertiserOptions) (*Advertiser, error) {
	reqURL := fmt.Sprintf("%s/networks/advertisers/%d", everflowAPIBaseURL, networkAdvertiserID)

	// Add query parameters if options are provided
	if opts != nil && len(opts.Relationships) > 0 {
		params := url.Values{}
		for _, relationship := range opts.Relationships {
			params.Add("relationship", relationship)
		}
		reqURL += "?" + params.Encode()
	}

	httpReq, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("X-Eflow-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to Everflow: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("advertiser with ID %d not found", networkAdvertiserID)
	}

	if resp.StatusCode != http.StatusOK {
		// Read error body for more details
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return nil, fmt.Errorf("Everflow API request failed with status %d: %v", resp.StatusCode, errorResp)
		}
		return nil, fmt.Errorf("Everflow API request failed with status %d", resp.StatusCode)
	}

	var advertiser Advertiser
	if err := json.NewDecoder(resp.Body).Decode(&advertiser); err != nil {
		return nil, fmt.Errorf("failed to decode Everflow get advertiser response: %w", err)
	}

	return &advertiser, nil
}

// EverflowUpdateAdvertiserRequest represents the request to update an advertiser in Everflow
type EverflowUpdateAdvertiserRequest struct {
	Name                           string              `json:"name"`
	AccountStatus                  string              `json:"account_status"`
	NetworkEmployeeID              int                 `json:"network_employee_id"`
	InternalNotes                  *string             `json:"internal_notes,omitempty"`
	AddressID                      *int                `json:"address_id,omitempty"`
	IsContactAddressEnabled        *bool               `json:"is_contact_address_enabled,omitempty"`
	SalesManagerID                 *int                `json:"sales_manager_id,omitempty"`
	IsExposePublisherReportingData *bool               `json:"is_expose_publisher_reporting_data,omitempty"`
	DefaultCurrencyID              string              `json:"default_currency_id"`
	PlatformName                   *string             `json:"platform_name,omitempty"`
	PlatformURL                    *string             `json:"platform_url,omitempty"`
	PlatformUsername               *string             `json:"platform_username,omitempty"`
	ReportingTimezoneID            int                 `json:"reporting_timezone_id"`
	AttributionMethod              *string             `json:"attribution_method,omitempty"`
	EmailAttributionMethod         *string             `json:"email_attribution_method,omitempty"`
	AttributionPriority            *string             `json:"attribution_priority,omitempty"`
	AccountingContactEmail         *string             `json:"accounting_contact_email,omitempty"`
	VerificationToken              *string             `json:"verification_token,omitempty"`
	OfferIDMacro                   *string             `json:"offer_id_macro,omitempty"`
	AffiliateIDMacro               *string             `json:"affiliate_id_macro,omitempty"`
	Labels                         []string            `json:"labels,omitempty"`
	ContactAddress                 *AdvertiserAddress  `json:"contact_address,omitempty"`
	Billing                        *AdvertiserBilling  `json:"billing,omitempty"`
	Settings                       *AdvertiserSettings `json:"settings,omitempty"`
}

// UpdateAdvertiser updates an existing advertiser in Everflow
func (c *Client) UpdateAdvertiser(ctx context.Context, networkAdvertiserID int64, req EverflowUpdateAdvertiserRequest) (*Advertiser, error) {
	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update advertiser request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/networks/advertisers/%d", everflowAPIBaseURL, networkAdvertiserID), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Eflow-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to Everflow: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("advertiser with ID %d not found", networkAdvertiserID)
	}

	if resp.StatusCode != http.StatusOK {
		// Read error body for more details
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return nil, fmt.Errorf("Everflow API request failed with status %d: %v", resp.StatusCode, errorResp)
		}
		return nil, fmt.Errorf("Everflow API request failed with status %d", resp.StatusCode)
	}

	var advertiser Advertiser
	if err := json.NewDecoder(resp.Body).Decode(&advertiser); err != nil {
		return nil, fmt.Errorf("failed to decode Everflow update advertiser response: %w", err)
	}

	return &advertiser, nil
}

// GetOfferOptions represents options for getting a single offer
type GetOfferOptions struct {
	Relationships []string `json:"relationships,omitempty"`
}

// Valid relationship values for GetOffer
const (
	OfferRelationshipCampaigns                = "campaigns"
	OfferRelationshipAdvertiserGlobalEvents   = "advertiser_global_events"
	OfferRelationshipOfferEmail               = "offer_email"
	OfferRelationshipOfferEmailOptout         = "offer_email_optout"
	OfferRelationshipReporting                = "reporting"
)

// GetOffer retrieves a single offer by ID from Everflow
func (c *Client) GetOffer(ctx context.Context, networkOfferID int64, opts *GetOfferOptions) (*Offer, error) {
	reqURL := fmt.Sprintf("%s/networks/offers/%d", everflowAPIBaseURL, networkOfferID)

	// Add query parameters if options are provided
	if opts != nil && len(opts.Relationships) > 0 {
		params := url.Values{}
		for _, relationship := range opts.Relationships {
			params.Add("relationship", relationship)
		}
		reqURL += "?" + params.Encode()
	}

	httpReq, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("X-Eflow-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to Everflow: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("offer with ID %d not found", networkOfferID)
	}

	if resp.StatusCode != http.StatusOK {
		// Read error body for more details
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return nil, fmt.Errorf("Everflow API request failed with status %d: %v", resp.StatusCode, errorResp)
		}
		return nil, fmt.Errorf("Everflow API request failed with status %d", resp.StatusCode)
	}

	var offer Offer
	if err := json.NewDecoder(resp.Body).Decode(&offer); err != nil {
		return nil, fmt.Errorf("failed to decode Everflow get offer response: %w", err)
	}

	return &offer, nil
}

// CreateOffer creates a new offer in Everflow
func (c *Client) CreateOffer(ctx context.Context, req OfferInput) (*Offer, error) {
	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create offer request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", everflowAPIBaseURL+"/networks/offers", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Eflow-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to Everflow: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		// Read error body for more details
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return nil, fmt.Errorf("Everflow API request failed with status %d: %v", resp.StatusCode, errorResp)
		}
		return nil, fmt.Errorf("Everflow API request failed with status %d", resp.StatusCode)
	}

	var offer Offer
	if err := json.NewDecoder(resp.Body).Decode(&offer); err != nil {
		return nil, fmt.Errorf("failed to decode Everflow create offer response: %w", err)
	}

	return &offer, nil
}

// UpdateOffer updates an existing offer in Everflow
func (c *Client) UpdateOffer(ctx context.Context, networkOfferID int64, req OfferInput) (*Offer, error) {
	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update offer request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/networks/offers/%d", everflowAPIBaseURL, networkOfferID), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Eflow-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to Everflow: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("offer with ID %d not found", networkOfferID)
	}

	if resp.StatusCode != http.StatusOK {
		// Read error body for more details
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return nil, fmt.Errorf("Everflow API request failed with status %d: %v", resp.StatusCode, errorResp)
		}
		return nil, fmt.Errorf("Everflow API request failed with status %d", resp.StatusCode)
	}

	var offer Offer
	if err := json.NewDecoder(resp.Body).Decode(&offer); err != nil {
		return nil, fmt.Errorf("failed to decode Everflow update offer response: %w", err)
	}

	return &offer, nil
}

// OffersTable retrieves a list of offers with optional filtering and search
func (c *Client) OffersTable(ctx context.Context, req OffersTableRequest, opts *OffersTableOptions) (*OffersTableResponse, error) {
	reqURL := everflowAPIBaseURL + "/networks/offerstable"

	// Add query parameters if options are provided
	if opts != nil {
		params := url.Values{}
		if opts.Page != nil {
			params.Add("page", strconv.Itoa(*opts.Page))
		}
		if opts.PageSize != nil {
			params.Add("page_size", strconv.Itoa(*opts.PageSize))
		}
		if len(opts.Relationships) > 0 {
			for _, relationship := range opts.Relationships {
				params.Add("relationship", relationship)
			}
		}
		if len(params) > 0 {
			reqURL += "?" + params.Encode()
		}
	}

	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal offers table request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Eflow-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to Everflow: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read error body for more details
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return nil, fmt.Errorf("Everflow API request failed with status %d: %v", resp.StatusCode, errorResp)
		}
		return nil, fmt.Errorf("Everflow API request failed with status %d", resp.StatusCode)
	}

	var offersResp OffersTableResponse
	if err := json.NewDecoder(resp.Body).Decode(&offersResp); err != nil {
		return nil, fmt.Errorf("failed to decode Everflow offers table response: %w", err)
	}

	return &offersResp, nil
}

// AddTagsToAdvertiser adds tags to an advertiser in Everflow
func (c *Client) AddTagsToAdvertiser(ctx context.Context, networkAdvertiserID int64, tags []string) error {
	if len(tags) == 0 {
		return nil // No tags to add
	}

	type tagRequest struct {
		Tags []string `json:"tags"`
	}

	req := tagRequest{
		Tags: tags,
	}

	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal add tags request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/networks/advertisers/%d/tags", everflowAPIBaseURL, networkAdvertiserID),
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Eflow-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to execute request to Everflow: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		// Read error body for more details
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return fmt.Errorf("Everflow API request failed with status %d: %v", resp.StatusCode, errorResp)
		}
		return fmt.Errorf("Everflow API request failed with status %d", resp.StatusCode)
	}

	return nil
}

// AddTagsToOffer adds tags to an offer in Everflow
func (c *Client) AddTagsToOffer(ctx context.Context, networkOfferID int64, tags []string) error {
	if len(tags) == 0 {
		return nil // No tags to add
	}

	type tagRequest struct {
		Tags []string `json:"tags"`
	}

	req := tagRequest{
		Tags: tags,
	}

	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal add tags request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/networks/offers/%d/tags", everflowAPIBaseURL, networkOfferID),
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Eflow-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to execute request to Everflow: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		// Read error body for more details
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return fmt.Errorf("Everflow API request failed with status %d: %v", resp.StatusCode, errorResp)
		}
		return fmt.Errorf("Everflow API request failed with status %d", resp.StatusCode)
	}

	return nil
}

// GenerateTrackingLink generates a tracking link for an offer
func (c *Client) GenerateTrackingLink(ctx context.Context, networkOfferID int64, networkAffiliateID int64, subIDs map[string]string) (string, error) {
	type trackingLinkRequest struct {
		NetworkOfferID     int64             `json:"network_offer_id"`
		NetworkAffiliateID int64             `json:"network_affiliate_id"`
		SubIDs             map[string]string `json:"sub_ids,omitempty"`
	}

	req := trackingLinkRequest{
		NetworkOfferID:     networkOfferID,
		NetworkAffiliateID: networkAffiliateID,
		SubIDs:             subIDs,
	}

	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal tracking link request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", everflowAPIBaseURL+"/networks/tracking_links", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Eflow-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to execute request to Everflow: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		// Read error body for more details
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return "", fmt.Errorf("Everflow API request failed with status %d: %v", resp.StatusCode, errorResp)
		}
		return "", fmt.Errorf("Everflow API request failed with status %d", resp.StatusCode)
	}

	var response struct {
		TrackingURL string `json:"tracking_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode Everflow tracking link response: %w", err)
	}

	return response.TrackingURL, nil
}
