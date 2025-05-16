package everflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	everflowAPIBaseURL = "https://api.eflow.team/v1" // Everflow API base URL
)

// Client represents an Everflow API client
type Client struct {
	httpClient *http.Client
	apiKey     string
}

// NewClient creates a new Everflow API client
func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiKey:     apiKey,
	}
}

// EverflowCreateOfferRequest represents the request to create an offer in Everflow
type EverflowCreateOfferRequest struct {
	Name                string  `json:"name"`
	NetworkAdvertiserID int64   `json:"network_advertiser_id"` // From advertiser_provider_mappings.provider_advertiser_id
	DestinationURL      string  `json:"destination_url"`
	OfferStatus         string  `json:"offer_status"` // e.g., "active", "pending", "paused"
	CurrencyID          string  `json:"currency_id"`  // e.g., "USD"
	Visibility          string  `json:"visibility"`   // e.g., "public", "private"
	// Add other fields based on Everflow's API documentation
}

// EverflowCreateOfferResponse represents the response from creating an offer in Everflow
type EverflowCreateOfferResponse struct {
	NetworkOfferID int64  `json:"network_offer_id"`
	OfferURL       string `json:"offer_url,omitempty"`
	// Add other fields based on Everflow's API documentation
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

// GenerateTrackingLink generates a tracking link for an offer
func (c *Client) GenerateTrackingLink(ctx context.Context, networkOfferID int64, networkAffiliateID int64, subIDs map[string]string) (string, error) {
	// Implementation depends on Everflow's API for generating tracking links
	// This is a placeholder - you'll need to implement based on Everflow's documentation
	
	type trackingLinkRequest struct {
		NetworkOfferID    int64             `json:"network_offer_id"`
		NetworkAffiliateID int64            `json:"network_affiliate_id"`
		SubIDs            map[string]string `json:"sub_ids,omitempty"`
	}
	
	req := trackingLinkRequest{
		NetworkOfferID:    networkOfferID,
		NetworkAffiliateID: networkAffiliateID,
		SubIDs:            subIDs,
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