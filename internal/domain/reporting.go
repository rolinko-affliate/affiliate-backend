package domain

import (
	"time"
)

// PerformanceSummary represents aggregated performance metrics
type PerformanceSummary struct {
	TotalClicks       int64   `json:"total_clicks"`
	TotalConversions  int64   `json:"total_conversions"`
	TotalRevenue      float64 `json:"total_revenue"`
	ConversionRate    float64 `json:"conversion_rate"`
	AverageRevenue    float64 `json:"average_revenue"`
	ClickThroughRate  float64 `json:"click_through_rate"`
	TotalImpressions  int64   `json:"total_impressions"`
}

// PerformanceTimeSeriesPoint represents a single data point in time series
type PerformanceTimeSeriesPoint struct {
	Date             string  `json:"date"`
	Clicks           int64   `json:"clicks"`
	Impressions      int64   `json:"impressions"`
	Conversions      int64   `json:"conversions"`
	Revenue          float64 `json:"revenue"`
	ConversionRate   float64 `json:"conversion_rate"`
	ClickThroughRate float64 `json:"click_through_rate"`
}

// DailyPerformanceReport represents daily performance breakdown
type DailyPerformanceReport struct {
	Date             string  `json:"date"`
	CampaignID       string  `json:"campaign_id"`
	CampaignName     string  `json:"campaign_name"`
	Clicks           int64   `json:"clicks"`
	Impressions      int64   `json:"impressions"`
	Conversions      int64   `json:"conversions"`
	Revenue          float64 `json:"revenue"`
	ConversionRate   float64 `json:"conversion_rate"`
	ClickThroughRate float64 `json:"click_through_rate"`
	Payouts          float64 `json:"payouts"`
}

// ConversionReport represents detailed conversion events
type ConversionReport struct {
	ID              string    `json:"id"`
	Timestamp       time.Time `json:"timestamp"`
	TransactionID   string    `json:"transaction_id"`
	CampaignID      string    `json:"campaign_id"`
	CampaignName    string    `json:"campaign_name"`
	OfferID         string    `json:"offer_id"`
	OfferName       string    `json:"offer_name"`
	Status          string    `json:"status"`
	Payout          float64   `json:"payout"`
	Currency        string    `json:"currency"`
	AffiliateID     string    `json:"affiliate_id"`
	AffiliateName   string    `json:"affiliate_name"`
	ClickID         *string   `json:"click_id,omitempty"`
	ConversionValue *float64  `json:"conversion_value,omitempty"`
	Sub1            *string   `json:"sub1,omitempty"`
	Sub2            *string   `json:"sub2,omitempty"`
	Sub3            *string   `json:"sub3,omitempty"`
}

// ClickReport represents detailed click events
type ClickReport struct {
	ID             string    `json:"id"`
	Timestamp      time.Time `json:"timestamp"`
	CampaignID     string    `json:"campaign_id"`
	CampaignName   string    `json:"campaign_name"`
	OfferID        string    `json:"offer_id"`
	OfferName      string    `json:"offer_name"`
	AffiliateID    string    `json:"affiliate_id"`
	AffiliateName  string    `json:"affiliate_name"`
	IPAddress      string    `json:"ip_address"`
	UserAgent      string    `json:"user_agent"`
	Country        string    `json:"country"`
	Region         *string   `json:"region,omitempty"`
	City           *string   `json:"city,omitempty"`
	ReferrerURL    *string   `json:"referrer_url,omitempty"`
	LandingPageURL string    `json:"landing_page_url"`
	Sub1           *string   `json:"sub1,omitempty"`
	Sub2           *string   `json:"sub2,omitempty"`
	Sub3           *string   `json:"sub3,omitempty"`
	Converted      bool      `json:"converted"`
	ConversionID   *string   `json:"conversion_id,omitempty"`
}

// ReportingFilters represents common filters for reporting queries
type ReportingFilters struct {
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"`
	CampaignIDs []string `json:"campaign_ids,omitempty"`
	AffiliateID *string  `json:"affiliate_id,omitempty"`
	Search      *string  `json:"search,omitempty"`
	Status      *string  `json:"status,omitempty"`
	Granularity *string  `json:"granularity,omitempty"`
}

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

// PaginationResult represents pagination metadata
type PaginationResult struct {
	CurrentPage     int  `json:"current_page"`
	TotalPages      int  `json:"total_pages"`
	TotalItems      int  `json:"total_items"`
	ItemsPerPage    int  `json:"items_per_page"`
	HasNextPage     bool `json:"has_next_page"`
	HasPreviousPage bool `json:"has_previous_page"`
}

// CampaignListItem represents a campaign for filter dropdown
type CampaignListItem struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}