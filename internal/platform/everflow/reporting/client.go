package reporting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents the Everflow reporting API client
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

// NewClient creates a new Everflow reporting client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

// EntityReportRequest represents the request for entity reporting
type EntityReportRequest struct {
	From       string                 `json:"from"`
	To         string                 `json:"to"`
	TimezoneID int                    `json:"timezone_id"`
	CurrencyID string                 `json:"currency_id"`
	Query      EntityReportQuery      `json:"query"`
	Columns    []EntityReportColumn   `json:"columns"`
}

type EntityReportQuery struct {
	Filters       []EntityReportFilter       `json:"filters,omitempty"`
	MetricFilters []EntityReportMetricFilter `json:"metric_filters,omitempty"`
	Exclusions    []EntityReportFilter       `json:"exclusions,omitempty"`
}

type EntityReportFilter struct {
	ResourceType  string `json:"resource_type"`
	FilterIDValue string `json:"filter_id_value"`
}

type EntityReportMetricFilter struct {
	MetricType  string      `json:"metric_type"`
	Operator    string      `json:"operator"`
	MetricValue interface{} `json:"metric_value"`
}

type EntityReportColumn struct {
	Column string `json:"column"`
}

// EntityReportResponse represents the response from entity reporting
type EntityReportResponse struct {
	Table             []map[string]interface{} `json:"table"`
	Performance       []map[string]interface{} `json:"performance"`
	Summary           map[string]interface{}   `json:"summary"`
	IncompleteResults bool                     `json:"incomplete_results,omitempty"`
}

// ConversionsRequest represents the request for conversions reporting
type ConversionsRequest struct {
	ShowConversions bool             `json:"show_conversions"`
	ShowEvents      bool             `json:"show_events"`
	ShowOnlyVT      bool             `json:"show_only_vt"`
	ShowOnlyCT      bool             `json:"show_only_ct"`
	From            string           `json:"from"`
	To              string           `json:"to"`
	TimezoneID      int              `json:"timezone_id"`
	CurrencyID      string           `json:"currency_id"`
	Query           ConversionsQuery `json:"query"`
}

type ConversionsQuery struct {
	Filters     []EntityReportFilter `json:"filters,omitempty"`
	SearchTerms []string             `json:"search_terms,omitempty"`
}

// ConversionsResponse represents the response from conversions reporting
type ConversionsResponse struct {
	Conversions []ConversionData `json:"conversions"`
	Paging      PagingInfo       `json:"paging"`
}

type ConversionData struct {
	ConversionID            string                 `json:"conversion_id"`
	ConversionUnixTimestamp int64                  `json:"conversion_unix_timestamp"`
	CostType                string                 `json:"cost_type"`
	Cost                    float64                `json:"cost"`
	SessionUserIP           string                 `json:"session_user_ip"`
	ConversionUserIP        string                 `json:"conversion_user_ip"`
	Country                 string                 `json:"country"`
	Region                  string                 `json:"region"`
	City                    string                 `json:"city"`
	DMA                     int                    `json:"dma"`
	Carrier                 string                 `json:"carrier"`
	Platform                string                 `json:"platform"`
	OSVersion               string                 `json:"os_version"`
	DeviceType              string                 `json:"device_type"`
	Brand                   string                 `json:"brand"`
	Browser                 string                 `json:"browser"`
	Language                string                 `json:"language"`
	HTTPUserAgent           string                 `json:"http_user_agent"`
	IsEvent                 bool                   `json:"is_event"`
	TransactionID           string                 `json:"transaction_id"`
	ClickUnixTimestamp      int64                  `json:"click_unix_timestamp"`
	Event                   string                 `json:"event"`
	CurrencyID              string                 `json:"currency_id"`
	ISP                     string                 `json:"isp"`
	Adv1                    string                 `json:"adv1"`
	Adv2                    string                 `json:"adv2"`
	Adv3                    string                 `json:"adv3"`
	Adv4                    string                 `json:"adv4"`
	Adv5                    string                 `json:"adv5"`
	OrderID                 string                 `json:"order_id"`
	SaleAmount              float64                `json:"sale_amount"`
	Relationship            ConversionRelationship `json:"relationship"`
	CouponCode              string                 `json:"coupon_code"`
}

type ConversionRelationship struct {
	Offer       ConversionOffer     `json:"offer"`
	EventsCount int                 `json:"events_count"`
	AffiliateID int                 `json:"affiliate_id"`
	Affiliate   ConversionAffiliate `json:"affiliate"`
	Sub1        string              `json:"sub1"`
	Sub2        string              `json:"sub2"`
	Sub3        string              `json:"sub3"`
	Sub4        string              `json:"sub4"`
	Sub5        string              `json:"sub5"`
	SourceID    string              `json:"source_id"`
	OfferURL    *string             `json:"offer_url"`
}

type ConversionOffer struct {
	NetworkOfferID int    `json:"network_offer_id"`
	NetworkID      int    `json:"network_id"`
	Name           string `json:"name"`
	OfferStatus    string `json:"offer_status"`
}

type ConversionAffiliate struct {
	NetworkAffiliateID int    `json:"network_affiliate_id"`
	NetworkID          int    `json:"network_id"`
	Name               string `json:"name"`
	AccountStatus      string `json:"account_status"`
}

type PagingInfo struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalCount int `json:"total_count"`
}

// DashboardSummaryRequest represents the request for dashboard summary
type DashboardSummaryRequest struct {
	TimezoneID int `json:"timezone_id"`
}

// DashboardSummaryResponse represents the response from dashboard summary
type DashboardSummaryResponse struct {
	Click      DashboardMetric `json:"click"`
	Conversion DashboardMetric `json:"conversion"`
	Cost       DashboardMetric `json:"cost"`
	CVR        DashboardMetric `json:"cvr"`
	Events     DashboardMetric `json:"events"`
	EVR        DashboardMetric `json:"evr"`
	Imp        DashboardMetric `json:"imp"`
}

type DashboardMetric struct {
	CurrentMonth       float64 `json:"current_month"`
	LastMonth          float64 `json:"last_month"`
	Today              float64 `json:"today"`
	TrendingPercentage float64 `json:"trending_percentage"`
	Yesterday          float64 `json:"yesterday"`
}

// GetEntityReport retrieves entity reporting data from Everflow
func (c *Client) GetEntityReport(ctx context.Context, req EntityReportRequest) (*EntityReportResponse, error) {
	result, err := c.postRequest(ctx, "/v1/advertisers/reporting/entity", req, &EntityReportResponse{})
	if err != nil {
		return nil, err
	}
	return result.(*EntityReportResponse), nil
}

// GetConversions retrieves conversions data from Everflow
func (c *Client) GetConversions(ctx context.Context, req ConversionsRequest) (*ConversionsResponse, error) {
	result, err := c.postRequest(ctx, "/v1/advertisers/reporting/conversions", req, &ConversionsResponse{})
	if err != nil {
		return nil, err
	}
	return result.(*ConversionsResponse), nil
}

// GetConversionByID retrieves a single conversion by ID from Everflow
func (c *Client) GetConversionByID(ctx context.Context, conversionID string) (*ConversionData, error) {
	result, err := c.getRequest(ctx, fmt.Sprintf("/v1/advertisers/reporting/conversions/%s", conversionID), &ConversionData{})
	if err != nil {
		return nil, err
	}
	return result.(*ConversionData), nil
}

// GetDashboardSummary retrieves dashboard summary data from Everflow
func (c *Client) GetDashboardSummary(ctx context.Context, req DashboardSummaryRequest) (*DashboardSummaryResponse, error) {
	result, err := c.postRequest(ctx, "/v1/advertisers/dashboard/summary", req, &DashboardSummaryResponse{})
	if err != nil {
		return nil, err
	}
	return result.(*DashboardSummaryResponse), nil
}

// postRequest makes a POST request to the Everflow API
func (c *Client) postRequest(ctx context.Context, endpoint string, payload interface{}, result interface{}) (interface{}, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Eflow-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

// getRequest makes a GET request to the Everflow API
func (c *Client) getRequest(ctx context.Context, endpoint string, result interface{}) (interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Eflow-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}