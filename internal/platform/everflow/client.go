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
	EmailAttributionMethod         *string             `json:"email_attribution_method,omitempty"` // "last_affiliate_attribution"
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
	BillingFrequency           string                 `json:"billing_frequency"` // "weekly", "monthly", "other"
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

// EverflowCreateOfferRequest represents the request to create an offer in Everflow
type EverflowCreateOfferRequest struct {
	Name                string              `json:"name"`
	NetworkAdvertiserID int64               `json:"network_advertiser_id"` // From advertiser_provider_mappings.provider_advertiser_id
	DestinationURL      string              `json:"destination_url"`
	OfferStatus         string              `json:"offer_status"`      // e.g., "active", "pending", "paused"
	CurrencyID          string              `json:"currency_id"`       // e.g., "USD"
	Visibility          string              `json:"visibility"`        // e.g., "public", "private"
	ConversionMethod    string              `json:"conversion_method"` // e.g., "server_postback", "http_image_pixel"
	NetworkCategoryID   *int                `json:"network_category_id,omitempty"`
	PreviewURL          *string             `json:"preview_url,omitempty"`
	SessionDefinition   *string             `json:"session_definition,omitempty"` // e.g., "cookie", "ip"
	SessionDuration     *int                `json:"session_duration,omitempty"`   // in hours
	InternalNotes       *string             `json:"internal_notes,omitempty"`
	Description         *string             `json:"description,omitempty"`
	IsCapsEnabled       *bool               `json:"is_caps_enabled,omitempty"`
	PayoutRevenue       []PayoutRevenueItem `json:"payout_revenue"`
	Tags                []string            `json:"tags,omitempty"`
}

// PayoutRevenueItem represents a payout/revenue structure for an offer
type PayoutRevenueItem struct {
	IsDefault          bool    `json:"is_default"`
	PayoutType         string  `json:"payout_type"` // "cpa", "cpc", "cpm", etc.
	PayoutAmount       float64 `json:"payout_amount"`
	RevenueType        string  `json:"revenue_type"` // "cpa", "cpc", "cpm", etc.
	RevenueAmount      float64 `json:"revenue_amount"`
	GoalName           *string `json:"goal_name,omitempty"`
	GoalTrackingMethod *string `json:"goal_tracking_method,omitempty"`
}

// EverflowCreateOfferResponse represents the response from creating an offer in Everflow
type EverflowCreateOfferResponse struct {
	NetworkOfferID      int64  `json:"network_offer_id"`
	NetworkID           int64  `json:"network_id"`
	NetworkAdvertiserID int64  `json:"network_advertiser_id"`
	Name                string `json:"name"`
	DestinationURL      string `json:"destination_url"`
	OfferStatus         string `json:"offer_status"`
	CurrencyID          string `json:"currency_id"`
	OfferURL            string `json:"offer_url,omitempty"`
	// Other fields omitted for brevity
}

// ListAdvertisersOptions represents options for listing advertisers
type ListAdvertisersOptions struct {
	Page     *int `json:"page,omitempty"`
	PageSize *int `json:"page_size,omitempty"`
}

// ListAdvertisers retrieves all advertisers with optional filters
func (c *Client) ListAdvertisers(ctx context.Context, opts *ListAdvertisersOptions) (*EverflowListAdvertisersResponse, error) {
	reqURL := everflowAPIBaseURL + "/v1/networks/advertisers"

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

	httpReq, err := http.NewRequestWithContext(ctx, "POST", everflowAPIBaseURL+"/v1/networks/advertisers", bytes.NewBuffer(payloadBytes))
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
	reqURL := fmt.Sprintf("%s/v1/networks/advertisers/%d", everflowAPIBaseURL, networkAdvertiserID)

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
	Name                    string              `json:"name"`
	AccountStatus           string              `json:"account_status"`
	NetworkEmployeeID       int                 `json:"network_employee_id"`
	InternalNotes           *string             `json:"internal_notes,omitempty"`
	AddressID               *int                `json:"address_id,omitempty"`
	IsContactAddressEnabled *bool               `json:"is_contact_address_enabled,omitempty"`
	SalesManagerID          *int                `json:"sales_manager_id,omitempty"`
	DefaultCurrencyID       string              `json:"default_currency_id"`
	PlatformName            *string             `json:"platform_name,omitempty"`
	PlatformURL             *string             `json:"platform_url,omitempty"`
	PlatformUsername        *string             `json:"platform_username,omitempty"`
	ReportingTimezoneID     int                 `json:"reporting_timezone_id"`
	AttributionMethod       *string             `json:"attribution_method,omitempty"`
	EmailAttributionMethod  *string             `json:"email_attribution_method,omitempty"`
	AttributionPriority     *string             `json:"attribution_priority,omitempty"`
	AccountingContactEmail  *string             `json:"accounting_contact_email,omitempty"`
	VerificationToken       *string             `json:"verification_token,omitempty"`
	OfferIDMacro            *string             `json:"offer_id_macro,omitempty"`
	AffiliateIDMacro        *string             `json:"affiliate_id_macro,omitempty"`
	Labels                  []string            `json:"labels,omitempty"`
	ContactAddress          *AdvertiserAddress  `json:"contact_address,omitempty"`
	Billing                 *AdvertiserBilling  `json:"billing,omitempty"`
	Settings                *AdvertiserSettings `json:"settings,omitempty"`
}

// EverflowUpdateAdvertiserResponse represents the response from updating an advertiser
type EverflowUpdateAdvertiserResponse struct {
	Result bool `json:"result"`
}

// UpdateAdvertiser updates an existing advertiser in Everflow
func (c *Client) UpdateAdvertiser(ctx context.Context, networkAdvertiserID int64, req EverflowUpdateAdvertiserRequest) (*EverflowUpdateAdvertiserResponse, error) {
	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update advertiser request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/v1/networks/advertisers/%d", everflowAPIBaseURL, networkAdvertiserID), bytes.NewBuffer(payloadBytes))
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

	var updateResp EverflowUpdateAdvertiserResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateResp); err != nil {
		return nil, fmt.Errorf("failed to decode Everflow update advertiser response: %w", err)
	}

	return &updateResp, nil
}

// CreateOffer creates a new offer in Everflow
func (c *Client) CreateOffer(ctx context.Context, req EverflowCreateOfferRequest) (*EverflowCreateOfferResponse, error) {
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

	var createResp EverflowCreateOfferResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return nil, fmt.Errorf("failed to decode Everflow create offer response: %w", err)
	}

	return &createResp, nil
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
		fmt.Sprintf("%s/v1/networks/advertisers/%d/tags", everflowAPIBaseURL, networkAdvertiserID),
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
