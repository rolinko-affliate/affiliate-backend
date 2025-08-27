package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// DashboardData represents the main dashboard response
type DashboardData struct {
	OrganizationType OrganizationType `json:"organization_type"`
	Summary          interface{}      `json:"summary"`          // Type varies by org type
	RecentActivity   []Activity       `json:"recent_activity"`
	LastUpdated      time.Time        `json:"last_updated"`
}

// Activity represents a dashboard activity item
type Activity struct {
	ID          string                 `json:"id" db:"id"`
	Type        ActivityType           `json:"type" db:"activity_type"`
	Description string                 `json:"description" db:"description"`
	Timestamp   time.Time              `json:"timestamp" db:"created_at"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" db:"metadata"`

	// Context-specific fields
	CampaignID   *int64  `json:"campaign_id,omitempty" db:"campaign_id"`
	CampaignName *string `json:"campaign_name,omitempty" db:"campaign_name"`
	ClientID     *int64  `json:"client_id,omitempty" db:"client_id"`
	ClientName   *string `json:"client_name,omitempty" db:"client_name"`
	Severity     *string `json:"severity,omitempty" db:"severity"`

	// Internal fields
	OrganizationID int64      `json:"-" db:"organization_id"`
	UserID         *uuid.UUID `json:"-" db:"user_id"`
}

// ActivityType represents different types of activities
type ActivityType string

const (
	// Advertiser activities
	ActivityCampaignCreated ActivityType = "campaign_created"
	ActivityCampaignUpdated ActivityType = "campaign_updated"
	ActivityConversion      ActivityType = "conversion"
	ActivityClick           ActivityType = "click"

	// Agency activities
	ActivityClientAdded      ActivityType = "client_added"
	ActivityCampaignLaunched ActivityType = "campaign_launched"
	ActivityPayment          ActivityType = "payment"

	// Platform owner activities
	ActivityUserRegistered   ActivityType = "user_registered"
	ActivitySystemAlert      ActivityType = "system_alert"
	ActivityRevenueMilestone ActivityType = "revenue_milestone"
	ActivityError            ActivityType = "error"
)

// IsValid checks if the activity type is valid
func (at ActivityType) IsValid() bool {
	switch at {
	case ActivityCampaignCreated, ActivityCampaignUpdated, ActivityConversion, ActivityClick,
		ActivityClientAdded, ActivityCampaignLaunched, ActivityPayment,
		ActivityUserRegistered, ActivitySystemAlert, ActivityRevenueMilestone, ActivityError:
		return true
	default:
		return false
	}
}

// String returns the string representation of the activity type
func (at ActivityType) String() string {
	return string(at)
}

// ActivityResponse represents paginated activity response
type ActivityResponse struct {
	Activities []Activity `json:"activities"`
	Total      int        `json:"total"`
	HasMore    bool       `json:"has_more"`
}

// MetricWithHistory represents a metric with historical comparison data
type MetricWithHistory struct {
	Today           float64 `json:"today"`
	Yesterday       float64 `json:"yesterday"`
	CurrentMonth    float64 `json:"current_month"`
	LastMonth       float64 `json:"last_month"`
	ChangePercentage float64 `json:"change_percentage"`
}

// AdvertiserDashboard represents advertiser-specific dashboard data
type AdvertiserDashboard struct {
	Metrics             AdvertiserMetrics   `json:"metrics"`
	CampaignPerformance CampaignPerformance `json:"campaign_performance"`
	RevenueChart        RevenueChart        `json:"revenue_chart"`
	Offers              *OffersPaginated    `json:"offers,omitempty"`
	Billing             *BillingInfo        `json:"billing,omitempty"`
}

// AdvertiserMetrics contains key metrics for advertisers with historical data
type AdvertiserMetrics struct {
	TotalClicks    MetricWithHistory `json:"total_clicks"`
	Conversions    MetricWithHistory `json:"conversions"`
	Revenue        MetricWithHistory `json:"revenue"`
	ConversionRate MetricWithHistory `json:"conversion_rate"`
	Events         MetricWithHistory `json:"events"`
	EventRate      MetricWithHistory `json:"event_rate"`
}

// AdvertiserSummary contains key metrics for advertisers (legacy support)
type AdvertiserSummary struct {
	TotalClicks    int64   `json:"total_clicks"`
	Conversions    int64   `json:"conversions"`
	Revenue        float64 `json:"revenue"`
	ConversionRate float64 `json:"conversion_rate"` // Percentage (0-100)
}

// OffersPaginated represents paginated offers data
type OffersPaginated struct {
	Items      []Offer `json:"items"`
	TotalCount int     `json:"total_count"`
	Page       int     `json:"page"`
	PerPage    int     `json:"per_page"`
	HasNext    bool    `json:"has_next"`
}

// Offer represents an offer/campaign offer
type Offer struct {
	ID             int64    `json:"id"`
	Name           string   `json:"name"`
	Description    *string  `json:"description,omitempty"`
	Payout         float64  `json:"payout"`
	Currency       string   `json:"currency"`
	Status         string   `json:"status"`
	Category       string   `json:"category"`
	Countries      []string `json:"countries"`
	ConversionFlow string   `json:"conversion_flow"`
	CreatedAt      string   `json:"created_at"`
}

// CampaignPerformance contains campaign-related metrics
type CampaignPerformance struct {
	Campaigns       []CampaignSummary `json:"campaigns"`
	TotalCampaigns  int               `json:"total_campaigns"`
	ActiveCampaigns int               `json:"active_campaigns"`
}

// CampaignSummary represents individual campaign performance
type CampaignSummary struct {
	ID               int64   `json:"id" db:"campaign_id"`
	Name             string  `json:"name" db:"name"`
	Clicks           int64   `json:"clicks" db:"clicks"`
	Conversions      int64   `json:"conversions" db:"conversions"`
	Revenue          float64 `json:"revenue" db:"revenue"`
	ConversionRate   float64 `json:"conversion_rate" db:"conversion_rate"`
	Events           int64   `json:"events" db:"events"`
	EventRate        float64 `json:"event_rate" db:"event_rate"`
	Status           string  `json:"status" db:"status"`
	StartDate        string  `json:"start_date" db:"start_date"`
	EndDate          *string `json:"end_date,omitempty" db:"end_date"`
	AdvertiserOrgID  int64   `json:"advertiser_org_id" db:"advertiser_org_id"`
}

// RevenueChart contains time-series revenue data
type RevenueChart struct {
	Data   []RevenueDataPoint `json:"data"`
	Period string             `json:"period"`
}

// RevenueDataPoint represents a single data point in revenue chart
type RevenueDataPoint struct {
	Date        string  `json:"date"`
	Revenue     float64 `json:"revenue"`
	Clicks      int64   `json:"clicks"`
	Conversions int64   `json:"conversions"`
	Events      int64   `json:"events"`
}

// BillingInfo contains billing-related information
type BillingInfo struct {
	CurrentBalance float64 `json:"current_balance"`
	MonthlySpend   float64 `json:"monthly_spend"`
	Currency       string  `json:"currency"`
}

// AgencyDashboard represents agency-specific dashboard data
type AgencyDashboard struct {
	PerformanceOverview      AgencyPerformanceOverview `json:"performance_overview"`
	AdvertiserOrganizations  []AdvertiserOrganization  `json:"advertiser_organizations"`
	CampaignsOverview        CampaignsOverview         `json:"campaigns_overview"`
	PerformanceChart         AgencyPerformanceChart    `json:"performance_chart"`
}

// AgencyPerformanceOverview contains aggregated performance across advertisers
type AgencyPerformanceOverview struct {
	TotalConversions   MetricWithHistory `json:"total_conversions"`
	TotalClicks        MetricWithHistory `json:"total_clicks"`
	ConversionsPerDay  MetricWithHistory `json:"conversions_per_day"`
	ClicksPerDay       MetricWithHistory `json:"clicks_per_day"`
	TotalEarnings      MetricWithHistory `json:"total_earnings"`
}

// AdvertiserOrganization represents an advertiser organization for agencies
type AdvertiserOrganization struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Status         string  `json:"status"`
	CampaignsCount int     `json:"campaigns_count"`
	Revenue        float64 `json:"revenue"`
	ConversionRate float64 `json:"conversion_rate"`
}

// CampaignsOverview contains campaigns overview across all advertisers
type CampaignsOverview struct {
	Campaigns   []AgencyCampaignSummary `json:"campaigns"`
	TotalCount  int                     `json:"total_count"`
	ActiveCount int                     `json:"active_count"`
}

// AgencyCampaignSummary represents campaign summary for agencies
type AgencyCampaignSummary struct {
	ID               int64   `json:"id"`
	Name             string  `json:"name"`
	AdvertiserOrgID  int64   `json:"advertiser_org_id"`
	AdvertiserName   string  `json:"advertiser_name"`
	Status           string  `json:"status"`
	Clicks           int64   `json:"clicks"`
	Conversions      int64   `json:"conversions"`
	TotalCost        float64 `json:"total_cost"`
	ConversionRate   float64 `json:"conversion_rate"`
}

// AgencyPerformanceChart contains agency-specific performance data
type AgencyPerformanceChart struct {
	Data   []AgencyPerformanceDataPoint `json:"data"`
	Period string                       `json:"period"`
}

// AgencyPerformanceDataPoint represents performance data with advertiser breakdown
type AgencyPerformanceDataPoint struct {
	Date                string                      `json:"date"`
	Conversions         int64                       `json:"conversions"`
	Clicks              int64                       `json:"clicks"`
	Revenue             float64                     `json:"revenue"`
	AdvertiserBreakdown []AdvertiserBreakdownPoint  `json:"advertiser_breakdown,omitempty"`
}

// AdvertiserBreakdownPoint represents revenue contribution by advertiser
type AdvertiserBreakdownPoint struct {
	AdvertiserID   int64   `json:"advertiser_id"`
	AdvertiserName string  `json:"advertiser_name"`
	Conversions    int64   `json:"conversions"`
	Revenue        float64 `json:"revenue"`
}

// AgencySummary contains key metrics for agencies (legacy support)
type AgencySummary struct {
	TotalClients          int     `json:"total_clients"`
	TotalRevenue          float64 `json:"total_revenue"`
	TotalConversions      int64   `json:"total_conversions"`
	AverageConversionRate float64 `json:"average_conversion_rate"`
}

// ClientPerformance contains client-related metrics
type ClientPerformance struct {
	Clients       []ClientSummary `json:"clients"`
	TopPerformers []TopPerformer  `json:"top_performers"`
}

// ClientSummary represents individual client performance
type ClientSummary struct {
	ID              int64   `json:"id" db:"organization_id"`
	Name            string  `json:"name" db:"name"`
	Revenue         float64 `json:"revenue" db:"revenue"`
	Conversions     int64   `json:"conversions" db:"conversions"`
	ConversionRate  float64 `json:"conversion_rate" db:"conversion_rate"`
	ActiveCampaigns int     `json:"active_campaigns" db:"active_campaigns"`
}

// TopPerformer represents top performing clients
type TopPerformer struct {
	ClientID   int64   `json:"client_id"`
	ClientName string  `json:"client_name"`
	Revenue    float64 `json:"revenue"`
	Growth     float64 `json:"growth"` // Percentage change
}

// AgencyRevenueChart contains agency-specific revenue data
type AgencyRevenueChart struct {
	Data   []AgencyRevenueDataPoint `json:"data"`
	Period string                   `json:"period"`
}

// AgencyRevenueDataPoint represents revenue data with client breakdown
type AgencyRevenueDataPoint struct {
	Date            string                   `json:"date"`
	TotalRevenue    float64                  `json:"total_revenue"`
	ClientBreakdown []ClientRevenueBreakdown `json:"client_breakdown"`
}

// ClientRevenueBreakdown represents revenue contribution by client
type ClientRevenueBreakdown struct {
	ClientID   int64   `json:"client_id"`
	ClientName string  `json:"client_name"`
	Revenue    float64 `json:"revenue"`
}

// PlatformOwnerDashboard represents platform owner dashboard data
type PlatformOwnerDashboard struct {
	PlatformOverview       PlatformOverview       `json:"platform_overview"`
	UserActivity           UserActivityMetrics    `json:"user_activity"`
	SystemHealth           SystemHealthMetrics    `json:"system_health"`
	RevenueBySource        []RevenueBySource      `json:"revenue_by_source"`
	GeographicDistribution []GeographicData       `json:"geographic_distribution"`
}

// PlatformOverview contains platform-wide metrics with historical data
type PlatformOverview struct {
	TotalOrganizations MetricWithHistory `json:"total_organizations"`
	TotalUsers         MetricWithHistory `json:"total_users"`
	TotalRevenue       MetricWithHistory `json:"total_revenue"`
	MonthlyGrowth      MetricWithHistory `json:"monthly_growth"`
	NewRegistrations   MetricWithHistory `json:"new_registrations"`
}

// GeographicData represents geographic distribution of users/revenue
type GeographicData struct {
	Country        string  `json:"country"`
	Revenue        float64 `json:"revenue"`
	Users          int     `json:"users"`
	ConversionRate float64 `json:"conversion_rate"`
}

// PlatformSummary contains high-level platform metrics (legacy support)
type PlatformSummary struct {
	TotalUsers        int     `json:"total_users"`
	TotalRevenue      float64 `json:"total_revenue"`
	TotalTransactions int64   `json:"total_transactions"`
	PlatformGrowth    float64 `json:"platform_growth"` // Percentage
}

// UserActivityMetrics contains user activity metrics
type UserActivityMetrics struct {
	DailyActiveUsers   int `json:"daily_active_users"`
	WeeklyActiveUsers  int `json:"weekly_active_users"`
	MonthlyActiveUsers int `json:"monthly_active_users"`
	ActiveAdvertisers  int `json:"active_advertisers"`
	ActiveAffiliates   int `json:"active_affiliates"`
}

// SystemHealthMetrics contains system health metrics
type SystemHealthMetrics struct {
	TotalCampaigns       int     `json:"total_campaigns"`
	RequestsPerMinute    int     `json:"requests_per_minute"`
	SuccessRate          float64 `json:"success_rate"`
	RateLimitHits        int     `json:"rate_limit_hits"`
	AverageQueryTime     float64 `json:"average_query_time"`
	ConnectionPoolUsage  float64 `json:"connection_pool_usage"`
}

// UserMetrics contains user-related metrics (legacy support)
type UserMetrics struct {
	ActiveUsers    int            `json:"active_users"`
	NewUsers       int            `json:"new_users"`
	UserGrowthRate float64        `json:"user_growth_rate"`
	UsersByType    map[string]int `json:"users_by_type"`
}

// RevenueMetrics contains revenue-related metrics (legacy support)
type RevenueMetrics struct {
	TotalRevenue          float64           `json:"total_revenue"`
	RevenueGrowth         float64           `json:"revenue_growth"`
	AverageRevenuePerUser float64           `json:"average_revenue_per_user"`
	RevenueBySource       []RevenueBySource `json:"revenue_by_source"`
}

// RevenueBySource represents revenue breakdown by source
type RevenueBySource struct {
	Source     string  `json:"source"`
	Revenue    float64 `json:"revenue"`
	Percentage float64 `json:"percentage"`
	Growth     float64 `json:"growth"`
}

// SystemHealth contains system health metrics (legacy support)
type SystemHealth struct {
	Uptime            float64 `json:"uptime"`             // Percentage
	ResponseTime      float64 `json:"response_time"`      // Average in milliseconds
	ErrorRate         float64 `json:"error_rate"`         // Percentage
	ActiveConnections int     `json:"active_connections"`
}

// CampaignDetail represents detailed campaign information
type CampaignDetail struct {
	Campaign    CampaignInfo    `json:"campaign"`
	Performance CampaignMetrics `json:"performance"`
	DailyStats  []DailyStat     `json:"daily_stats"`
}

// CampaignInfo contains basic campaign information
type CampaignInfo struct {
	ID        int64      `json:"id" db:"campaign_id"`
	Name      string     `json:"name" db:"name"`
	Status    string     `json:"status" db:"status"`
	StartDate time.Time  `json:"start_date" db:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty" db:"end_date"`
	Budget    float64    `json:"budget" db:"budget"`
	Spent     float64    `json:"spent" db:"spent"`
}

// ClientInfo represents basic client information for agencies
type ClientInfo struct {
	ID          int64   `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Status      string  `json:"status" db:"status"`
	TotalSpend  float64 `json:"total_spend" db:"total_spend"`
	Campaigns   int     `json:"campaigns" db:"campaigns"`
}

// CampaignMetrics contains campaign performance metrics
type CampaignMetrics struct {
	Clicks      int64   `json:"clicks"`
	Impressions int64   `json:"impressions"`
	Conversions int64   `json:"conversions"`
	Revenue     float64 `json:"revenue"`
	CTR         float64 `json:"ctr"` // Click-through rate
	CPC         float64 `json:"cpc"` // Cost per click
	CPA         float64 `json:"cpa"` // Cost per acquisition
}

// DailyStat represents daily performance statistics
type DailyStat struct {
	Date        string  `json:"date" db:"stat_date"`
	Clicks      int64   `json:"clicks" db:"clicks"`
	Impressions int64   `json:"impressions" db:"impressions"`
	Conversions int64   `json:"conversions" db:"conversions"`
	Revenue     float64 `json:"revenue" db:"revenue"`
}

// DashboardMetric represents pre-calculated dashboard metrics
type DashboardMetric struct {
	ID             uuid.UUID              `json:"id" db:"id"`
	OrganizationID int64                  `json:"organization_id" db:"organization_id"`
	MetricType     string                 `json:"metric_type" db:"metric_type"`
	MetricDate     time.Time              `json:"metric_date" db:"metric_date"`
	MetricData     map[string]interface{} `json:"metric_data" db:"metric_data"`
	CreatedAt      time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at" db:"updated_at"`
}

// SystemHealthMetric represents system health data points
type SystemHealthMetric struct {
	ID                    uuid.UUID `json:"id" db:"id"`
	MetricTimestamp       time.Time `json:"metric_timestamp" db:"metric_timestamp"`
	UptimePercentage      *float64  `json:"uptime_percentage" db:"uptime_percentage"`
	AvgResponseTimeMs     *float64  `json:"avg_response_time_ms" db:"avg_response_time_ms"`
	ErrorRatePercentage   *float64  `json:"error_rate_percentage" db:"error_rate_percentage"`
	ActiveConnections     *int      `json:"active_connections" db:"active_connections"`
	TotalUsers            *int      `json:"total_users" db:"total_users"`
	TotalRevenue          *float64  `json:"total_revenue" db:"total_revenue"`
	TotalTransactions     *int64    `json:"total_transactions" db:"total_transactions"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
}

// DashboardQuery represents query parameters for dashboard requests
type DashboardQuery struct {
	Period    string     `json:"period"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Timezone  string     `json:"timezone"`
}

// Validate validates the dashboard query parameters
func (dq *DashboardQuery) Validate() error {
	// Validate period
	validPeriods := map[string]bool{
		"today": true,
		"7d":    true,
		"30d":   true,
		"90d":   true,
		"custom": true,
	}

	if dq.Period != "" && !validPeriods[dq.Period] {
		return errors.New("invalid period value")
	}

	// If custom period, start and end dates are required
	if dq.Period == "custom" {
		if dq.StartDate == nil || dq.EndDate == nil {
			return errors.New("start_date and end_date are required for custom period")
		}

		if dq.EndDate.Before(*dq.StartDate) {
			return errors.New("end_date must be after start_date")
		}

		// Limit to 1 year maximum
		if dq.EndDate.Sub(*dq.StartDate) > 365*24*time.Hour {
			return errors.New("date range cannot exceed 1 year")
		}
	}

	return nil
}

// GetDateRange returns the actual date range based on the query parameters
func (dq *DashboardQuery) GetDateRange() (time.Time, time.Time) {
	now := time.Now().UTC()

	switch dq.Period {
	case "today":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		return start, now
	case "7d":
		return now.AddDate(0, 0, -7), now
	case "30d":
		return now.AddDate(0, 0, -30), now
	case "90d":
		return now.AddDate(0, 0, -90), now
	case "custom":
		if dq.StartDate != nil && dq.EndDate != nil {
			return *dq.StartDate, *dq.EndDate
		}
	}

	// Default to last 30 days
	return now.AddDate(0, 0, -30), now
}

// Everflow API Models
type EverflowEntityRequest struct {
	From       string           `json:"from"`
	To         string           `json:"to"`
	TimezoneID int             `json:"timezone_id"`
	CurrencyID string          `json:"currency_id"`
	Query      EverflowQuery   `json:"query"`
	Columns    []EverflowColumn `json:"columns"`
}

type EverflowQuery struct {
	Filters       []EverflowFilter       `json:"filters"`
	MetricFilters []EverflowMetricFilter `json:"metric_filters"`
	Exclusions    []EverflowFilter       `json:"exclusions"`
	SearchTerms   []string               `json:"search_terms"`
}

type EverflowFilter struct {
	ResourceType    string `json:"resource_type"`
	FilterIDValue   string `json:"filter_id_value"`
}

type EverflowMetricFilter struct {
	MetricType  string      `json:"metric_type"`
	Operator    string      `json:"operator"`
	MetricValue interface{} `json:"metric_value"`
}

type EverflowColumn struct {
	Column string `json:"column"`
}

type EverflowEntityResponse struct {
	Table             EverflowTable       `json:"table"`
	Performance       EverflowPerformance `json:"performance"`
	Summary           EverflowSummary     `json:"summary"`
	IncompleteResults bool                `json:"incomplete_results"`
}

type EverflowTable struct {
	Rows []map[string]interface{} `json:"rows"`
}

type EverflowPerformance struct {
	Data []EverflowPerformancePoint `json:"data"`
}

type EverflowPerformancePoint struct {
	Timestamp string                 `json:"timestamp"`
	Metrics   map[string]interface{} `json:"metrics"`
}

type EverflowSummary struct {
	Metrics map[string]interface{} `json:"metrics"`
}

type EverflowConversionRequest struct {
	ShowConversions bool          `json:"show_conversions"`
	ShowEvents      bool          `json:"show_events"`
	ShowOnlyVT      bool          `json:"show_only_vt"`
	ShowOnlyCT      bool          `json:"show_only_ct"`
	From            string        `json:"from"`
	To              string        `json:"to"`
	TimezoneID      int          `json:"timezone_id"`
	CurrencyID      string       `json:"currency_id"`
	Query           EverflowQuery `json:"query"`
}

type EverflowConversionResponse struct {
	Conversions []EverflowConversion `json:"conversions"`
	Paging      EverflowPaging       `json:"paging"`
}

type EverflowConversion struct {
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
	Relationship            EverflowRelationship   `json:"relationship"`
	CouponCode              string                 `json:"coupon_code"`
}

type EverflowRelationship struct {
	Offer       EverflowOffer     `json:"offer"`
	EventsCount int               `json:"events_count"`
	AffiliateID int               `json:"affiliate_id"`
	Affiliate   EverflowAffiliate `json:"affiliate"`
	Sub1        string            `json:"sub1"`
	Sub2        string            `json:"sub2"`
	Sub3        string            `json:"sub3"`
	Sub4        string            `json:"sub4"`
	Sub5        string            `json:"sub5"`
	SourceID    string            `json:"source_id"`
	OfferURL    *string           `json:"offer_url"`
}

type EverflowOffer struct {
	NetworkOfferID int    `json:"network_offer_id"`
	NetworkID      int    `json:"network_id"`
	Name           string `json:"name"`
	OfferStatus    string `json:"offer_status"`
}

type EverflowAffiliate struct {
	NetworkAffiliateID int    `json:"network_affiliate_id"`
	NetworkID          int    `json:"network_id"`
	Name               string `json:"name"`
	AccountStatus      string `json:"account_status"`
}

type EverflowPaging struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalCount int `json:"total_count"`
}

type EverflowExportRequest struct {
	From            string        `json:"from"`
	To              string        `json:"to"`
	TimezoneID      int          `json:"timezone_id"`
	ShowConversions bool         `json:"show_conversions"`
	ShowEvents      bool         `json:"show_events"`
	Query           EverflowQuery `json:"query"`
	Format          string       `json:"format"` // "csv" or "json"
}

type EverflowDashboardSummary struct {
	Click      EverflowMetricSummary `json:"click"`
	Conversion EverflowMetricSummary `json:"conversion"`
	Cost       EverflowMetricSummary `json:"cost"`
	CVR        EverflowCVRSummary    `json:"cvr"`
	Events     EverflowMetricSummary `json:"events"`
	EVR        EverflowCVRSummary    `json:"evr"`
	Imp        EverflowMetricSummary `json:"imp"`
}

type EverflowMetricSummary struct {
	CurrentMonth       int     `json:"current_month"`
	LastMonth          int     `json:"last_month"`
	Today              int     `json:"today"`
	TrendingPercentage float64 `json:"trending_percentage"`
	Yesterday          int     `json:"yesterday"`
}

type EverflowCVRSummary struct {
	CurrentMonth float64 `json:"current_month"`
	LastMonth    float64 `json:"last_month"`
	Today        float64 `json:"today"`
	Yesterday    float64 `json:"yesterday"`
}