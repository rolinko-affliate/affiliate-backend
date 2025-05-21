package everflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	Name                 string                 `json:"name"`
	AccountStatus        string                 `json:"account_status"` // "active", "inactive", "suspended"
	DefaultCurrencyID    string                 `json:"default_currency_id"` // e.g., "USD"
	NetworkEmployeeID    *int                   `json:"network_employee_id,omitempty"`
	SalesManagerID       *int                   `json:"sales_manager_id,omitempty"`
	InternalNotes        *string                `json:"internal_notes,omitempty"`
	ReportingTimezoneID  *int                   `json:"reporting_timezone_id,omitempty"`
	AttributionMethod    *string                `json:"attribution_method,omitempty"` // "last_touch" or "first_touch"
	VerificationToken    *string                `json:"verification_token,omitempty"`
	IsContactAddressEnabled *bool               `json:"is_contact_address_enabled,omitempty"`
	ContactAddress       *AdvertiserAddress     `json:"contact_address,omitempty"`
	Labels               []string               `json:"labels,omitempty"`
	Users                []AdvertiserUser       `json:"users,omitempty"`
	Billing              *AdvertiserBilling     `json:"billing,omitempty"`
	Settings             *AdvertiserSettings    `json:"settings,omitempty"`
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
	TimezoneID      *int    `json:"timezone_id,omitempty"`
	CurrencyID      *string `json:"currency_id,omitempty"`
}

// AdvertiserBilling represents an advertiser's billing information
type AdvertiserBilling struct {
	BillingFrequency    string                 `json:"billing_frequency"` // "weekly", "monthly"
	TaxID               *string                `json:"tax_id,omitempty"`
	IsInvoiceCreationAuto *bool                `json:"is_invoice_creation_auto,omitempty"`
	Details             map[string]interface{} `json:"details,omitempty"`
	DefaultPaymentTerms *int                   `json:"default_payment_terms,omitempty"`
	DayOfMonth         *int                    `json:"day_of_month,omitempty"`
}

// AdvertiserSettings represents advertiser settings
type AdvertiserSettings struct {
	ExposedVariables map[string]bool `json:"exposed_variables,omitempty"`
}

// EverflowCreateAdvertiserResponse represents the response from creating an advertiser in Everflow
type EverflowCreateAdvertiserResponse struct {
	NetworkAdvertiserID int64                  `json:"network_advertiser_id"`
	Name                string                 `json:"name"`
	AccountStatus       string                 `json:"account_status"`
	DefaultCurrencyID   string                 `json:"default_currency_id"`
	TimeCreated         int64                  `json:"time_created"`
	TimeSaved           int64                  `json:"time_saved"`
	// Other fields omitted for brevity
}

// EverflowCreateOfferRequest represents the request to create an offer in Everflow
type EverflowCreateOfferRequest struct {
	Name                string                 `json:"name"`
	NetworkAdvertiserID int64                  `json:"network_advertiser_id"` // From advertiser_provider_mappings.provider_advertiser_id
	DestinationURL      string                 `json:"destination_url"`
	OfferStatus         string                 `json:"offer_status"` // e.g., "active", "pending", "paused"
	CurrencyID          string                 `json:"currency_id"`  // e.g., "USD"
	Visibility          string                 `json:"visibility"`   // e.g., "public", "private"
	ConversionMethod    string                 `json:"conversion_method"` // e.g., "server_postback", "http_image_pixel"
	NetworkCategoryID   *int                   `json:"network_category_id,omitempty"`
	PreviewURL          *string                `json:"preview_url,omitempty"`
	SessionDefinition   *string                `json:"session_definition,omitempty"` // e.g., "cookie", "ip"
	SessionDuration     *int                   `json:"session_duration,omitempty"`   // in hours
	InternalNotes       *string                `json:"internal_notes,omitempty"`
	Description         *string                `json:"description,omitempty"`
	IsCapsEnabled       *bool                  `json:"is_caps_enabled,omitempty"`
	PayoutRevenue       []PayoutRevenueItem    `json:"payout_revenue"`
	Tags                []string               `json:"tags,omitempty"`
}

// PayoutRevenueItem represents a payout/revenue structure for an offer
type PayoutRevenueItem struct {
	IsDefault           bool    `json:"is_default"`
	PayoutType          string  `json:"payout_type"` // "cpa", "cpc", "cpm", etc.
	PayoutAmount        float64 `json:"payout_amount"`
	RevenueType         string  `json:"revenue_type"` // "cpa", "cpc", "cpm", etc.
	RevenueAmount       float64 `json:"revenue_amount"`
	GoalName            *string `json:"goal_name,omitempty"`
	GoalTrackingMethod  *string `json:"goal_tracking_method,omitempty"`
}

// EverflowCreateOfferResponse represents the response from creating an offer in Everflow
type EverflowCreateOfferResponse struct {
	NetworkOfferID      int64   `json:"network_offer_id"`
	NetworkID           int64   `json:"network_id"`
	NetworkAdvertiserID int64   `json:"network_advertiser_id"`
	Name                string  `json:"name"`
	DestinationURL      string  `json:"destination_url"`
	OfferStatus         string  `json:"offer_status"`
	CurrencyID          string  `json:"currency_id"`
	OfferURL            string  `json:"offer_url,omitempty"`
	// Other fields omitted for brevity
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