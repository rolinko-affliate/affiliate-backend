package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/logger"
	"github.com/affiliate-backend/internal/repository"
)

const (
	DefaultCacheTTL = 15 * time.Minute
)

// DashboardService defines the interface for dashboard business logic
type DashboardService interface {
	// Main dashboard methods
	GetDashboardData(ctx context.Context, userID uuid.UUID, period string, startDate, endDate *time.Time, timezone string) (*domain.DashboardData, error)
	GetCampaignDetail(ctx context.Context, userID uuid.UUID, campaignID int64) (*domain.CampaignDetail, error)
	GetRecentActivity(ctx context.Context, userID uuid.UUID, limit, offset int, activityTypes []string, since *time.Time) (*domain.ActivityResponse, error)

	// System health (platform owner only)
	GetSystemHealth(ctx context.Context, userID uuid.UUID) (*domain.SystemHealth, error)

	// Activity tracking
	TrackActivity(ctx context.Context, orgID int64, activityType domain.ActivityType, description string, metadata map[string]interface{}) error
	TrackActivityWithUser(ctx context.Context, userID uuid.UUID, orgID int64, activityType domain.ActivityType, description string, metadata map[string]interface{}) error

	// Cache management
	InvalidateCache(ctx context.Context, orgID int64) error

	// Health checks
	CheckDatabaseHealth(ctx context.Context) error
	CheckCacheHealth(ctx context.Context) error
}

// dashboardService implements DashboardService
type dashboardService struct {
	cacheRepo          repository.DashboardCacheRepository
	everflowRepo       repository.EverflowRepository
	reportingService   ReportingService
	profileService     ProfileService
	organizationService OrganizationService
	logger             *logger.Logger
}

// NewDashboardService creates a new dashboard service
func NewDashboardService(
	cacheRepo repository.DashboardCacheRepository,
	everflowRepo repository.EverflowRepository,
	reportingService ReportingService,
	profileService ProfileService,
	organizationService OrganizationService,
	logger *logger.Logger,
) DashboardService {
	return &dashboardService{
		cacheRepo:          cacheRepo,
		everflowRepo:       everflowRepo,
		reportingService:   reportingService,
		profileService:     profileService,
		organizationService: organizationService,
		logger:             logger,
	}
}

// GetDashboardData retrieves dashboard data based on user's organization type
func (s *dashboardService) GetDashboardData(ctx context.Context, userID uuid.UUID, period string, startDate, endDate *time.Time, timezone string) (*domain.DashboardData, error) {
	start := time.Now()

	s.logger.Info("Dashboard data request started",
		"user_id", userID,
		"period", period,
		"start_date", startDate,
		"end_date", endDate,
		"timezone", timezone,
	)

	defer func() {
		s.logger.Info("Dashboard data request completed",
			"user_id", userID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	}()

	// Get user profile to determine organization and permissions
	profile, err := s.profileService.GetProfileByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	if profile.OrganizationID == nil {
		return nil, errors.New("user is not associated with any organization")
	}

	orgID := *profile.OrganizationID

	// Get organization details
	org, err := s.organizationService.GetOrganizationByID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	// Parse query parameters
	query := &domain.DashboardQuery{
		Period:    period,
		StartDate: startDate,
		EndDate:   endDate,
		Timezone:  timezone,
	}

	if err := query.Validate(); err != nil {
		return nil, err
	}

	// Check cache first
	cachedData, err := s.cacheRepo.GetCachedDashboardData(ctx, orgID, org.Type)
	if err != nil {
		s.logger.Warn("Failed to get cached dashboard data", "error", err)
	} else if cachedData != nil {
		s.logger.Debug("Returning cached dashboard data", "org_id", orgID)
		return cachedData, nil
	}

	// Get date range
	from, to := query.GetDateRange()

	// For now, return empty activities since we're focusing on Everflow data
	activities := []domain.Activity{}

	// Build dashboard data based on organization type
	var summary interface{}
	switch org.Type {
	case domain.OrganizationTypeAdvertiser:
		summary, err = s.getAdvertiserSummary(ctx, orgID, from, to)
	case domain.OrganizationTypeAgency:
		summary, err = s.getAgencySummary(ctx, orgID, from, to)
	case domain.OrganizationTypePlatformOwner:
		summary, err = s.getPlatformOwnerSummary(ctx, from, to)
	default:
		return nil, errors.New("unsupported organization type")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard summary: %w", err)
	}

	dashboardData := &domain.DashboardData{
		OrganizationType: org.Type,
		Summary:          summary,
		RecentActivity:   activities,
		LastUpdated:      time.Now(),
	}

	// Cache the result
	if err := s.cacheRepo.SetCachedDashboardData(ctx, orgID, org.Type, dashboardData, DefaultCacheTTL); err != nil {
		s.logger.Warn("Failed to cache dashboard data", "error", err)
	}

	return dashboardData, nil
}

// getAdvertiserSummary builds advertiser-specific dashboard summary using Everflow data
func (s *dashboardService) getAdvertiserSummary(ctx context.Context, orgID int64, from, to time.Time) (*domain.AdvertiserDashboard, error) {
	// Get dashboard summary from Everflow
	// For timezone, we'll use UTC (timezone ID 67) as default
	dashboardSummary, err := s.everflowRepo.GetDashboardSummary(ctx, 67)
	if err != nil {
		return nil, fmt.Errorf("failed to get Everflow dashboard summary: %w", err)
	}

	// Convert Everflow data to our domain model
	advertiserSummary := domain.AdvertiserSummary{
		TotalClicks:    int64(dashboardSummary.Click.Today),
		Conversions:    int64(dashboardSummary.Conversion.Today),
		Revenue:        float64(dashboardSummary.Cost.Today),
		ConversionRate: dashboardSummary.CVR.Today,
	}

	// Get entity report for campaign-level data
	entityReq := &domain.EverflowEntityRequest{
		From:       from.Format("2006-01-02"),
		To:         to.Format("2006-01-02"),
		TimezoneID: 67, // UTC
		CurrencyID: "USD",
		Query: domain.EverflowQuery{
			Filters:       []domain.EverflowFilter{},
			MetricFilters: []domain.EverflowMetricFilter{},
			Exclusions:    []domain.EverflowFilter{},
			SearchTerms:   []string{},
		},
		Columns: []domain.EverflowColumn{
			{Column: "offer"},
		},
	}

	entityResponse, err := s.everflowRepo.GetEntityReport(ctx, entityReq)
	if err != nil {
		s.logger.Warn("Failed to get entity report", "error", err)
		// Continue with empty campaign data
	}

	// Convert entity data to campaign summaries
	var campaignSummaries []domain.CampaignSummary
	activeCampaigns := 0
	
	if entityResponse != nil {
		for _, row := range entityResponse.Table.Rows {
			if offerName, ok := row["offer"].(string); ok {
				clicks := int64(0)
				conversions := int64(0)
				revenue := float64(0)
				
				if clicksVal, ok := row["clicks"]; ok {
					if clicksFloat, ok := clicksVal.(float64); ok {
						clicks = int64(clicksFloat)
					}
				}
				if conversionsVal, ok := row["conversions"]; ok {
					if conversionsFloat, ok := conversionsVal.(float64); ok {
						conversions = int64(conversionsFloat)
					}
				}
				if revenueVal, ok := row["revenue"]; ok {
					if revenueFloat, ok := revenueVal.(float64); ok {
						revenue = revenueFloat
					}
				}
				
				conversionRate := float64(0)
				if clicks > 0 {
					conversionRate = float64(conversions) / float64(clicks) * 100
				}
				
				campaignSummaries = append(campaignSummaries, domain.CampaignSummary{
					ID:             int64(len(campaignSummaries) + 1), // Generate ID
					Name:           offerName,
					Status:         "active",
					Clicks:         clicks,
					Conversions:    conversions,
					Revenue:        revenue,
					ConversionRate: conversionRate,
				})
				activeCampaigns++
			}
		}
	}

	campaignPerformance := domain.CampaignPerformance{
		Campaigns:       campaignSummaries,
		TotalCampaigns:  len(campaignSummaries),
		ActiveCampaigns: activeCampaigns,
	}

	// Build revenue chart from Everflow performance data
	revenueChart := s.buildRevenueChartFromEverflow(entityResponse)

	// Create basic billing info from Everflow cost data
	billing := &domain.BillingInfo{
		CurrentBalance: float64(dashboardSummary.Cost.CurrentMonth),
		MonthlySpend:   float64(dashboardSummary.Cost.Today),
		Currency:       "USD",
	}

	return &domain.AdvertiserDashboard{
		Summary:             advertiserSummary,
		CampaignPerformance: campaignPerformance,
		RevenueChart:        revenueChart,
		Billing:             billing,
	}, nil
}

// buildRevenueChartFromEverflow builds revenue chart from Everflow performance data
func (s *dashboardService) buildRevenueChartFromEverflow(entityResponse *domain.EverflowEntityResponse) domain.RevenueChart {
	var dataPoints []domain.RevenueDataPoint
	
	if entityResponse != nil && entityResponse.Performance.Data != nil {
		for _, point := range entityResponse.Performance.Data {
			revenue := float64(0)
			clicks := int64(0)
			conversions := int64(0)
			
			if revenueVal, ok := point.Metrics["revenue"]; ok {
				if revenueFloat, ok := revenueVal.(float64); ok {
					revenue = revenueFloat
				}
			}
			if clicksVal, ok := point.Metrics["clicks"]; ok {
				if clicksFloat, ok := clicksVal.(float64); ok {
					clicks = int64(clicksFloat)
				}
			}
			if conversionsVal, ok := point.Metrics["conversions"]; ok {
				if conversionsFloat, ok := conversionsVal.(float64); ok {
					conversions = int64(conversionsFloat)
				}
			}
			
			dataPoints = append(dataPoints, domain.RevenueDataPoint{
				Date:        point.Timestamp,
				Revenue:     revenue,
				Clicks:      clicks,
				Conversions: conversions,
			})
		}
	}
	
	// If no performance data, create a basic chart with current data
	if len(dataPoints) == 0 {
		dataPoints = []domain.RevenueDataPoint{
			{
				Date:        time.Now().Format("2006-01-02"),
				Revenue:     0,
				Clicks:      0,
				Conversions: 0,
			},
		}
	}
	
	return domain.RevenueChart{
		Period: "daily",
		Data:   dataPoints,
	}
}

// getAgencySummary builds agency-specific dashboard summary using Everflow data
func (s *dashboardService) getAgencySummary(ctx context.Context, orgID int64, from, to time.Time) (*domain.AgencyDashboard, error) {
	// Get dashboard summary from Everflow
	dashboardSummary, err := s.everflowRepo.GetDashboardSummary(ctx, 67)
	if err != nil {
		return nil, fmt.Errorf("failed to get Everflow dashboard summary: %w", err)
	}

	agencySummary := domain.AgencySummary{
		TotalClients:          5, // Mock data for now
		TotalRevenue:          float64(dashboardSummary.Cost.CurrentMonth),
		TotalConversions:      int64(dashboardSummary.Conversion.CurrentMonth),
		AverageConversionRate: dashboardSummary.CVR.CurrentMonth,
	}

	// Create mock client performance data
	topPerformers := []domain.TopPerformer{
		{
			ClientID:   1,
			ClientName: "Client A",
			Revenue:    float64(dashboardSummary.Cost.Today) * 0.4,
			Growth:     5.2,
		},
		{
			ClientID:   2,
			ClientName: "Client B", 
			Revenue:    float64(dashboardSummary.Cost.Today) * 0.3,
			Growth:     3.1,
		},
	}

	clientSummaries := []domain.ClientSummary{
		{
			ID:              1,
			Name:            "Client A",
			Revenue:         float64(dashboardSummary.Cost.Today) * 0.4,
			Conversions:     int64(dashboardSummary.Conversion.Today) * 40 / 100,
			ConversionRate:  dashboardSummary.CVR.Today,
			ActiveCampaigns: 3,
		},
		{
			ID:              2,
			Name:            "Client B",
			Revenue:         float64(dashboardSummary.Cost.Today) * 0.3,
			Conversions:     int64(dashboardSummary.Conversion.Today) * 30 / 100,
			ConversionRate:  dashboardSummary.CVR.Today,
			ActiveCampaigns: 2,
		},
	}

	clientPerformance := domain.ClientPerformance{
		Clients:       clientSummaries,
		TopPerformers: topPerformers,
	}

	// Build revenue chart from Everflow data
	revenueChart := s.buildRevenueChartFromEverflow(nil)
	
	// Convert to agency revenue chart with client breakdown
	agencyRevenueChart := domain.AgencyRevenueChart{
		Period: revenueChart.Period,
		Data: []domain.AgencyRevenueDataPoint{
			{
				Date:         time.Now().Format("2006-01-02"),
				TotalRevenue: float64(dashboardSummary.Cost.Today),
				ClientBreakdown: []domain.ClientRevenueBreakdown{
					{ClientID: 1, ClientName: "Client A", Revenue: float64(dashboardSummary.Cost.Today) * 0.4},
					{ClientID: 2, ClientName: "Client B", Revenue: float64(dashboardSummary.Cost.Today) * 0.3},
					{ClientID: 3, ClientName: "Client C", Revenue: float64(dashboardSummary.Cost.Today) * 0.3},
				},
			},
		},
	}

	return &domain.AgencyDashboard{
		Summary:           agencySummary,
		ClientPerformance: clientPerformance,
		RevenueChart:      agencyRevenueChart,
	}, nil
}

// getPlatformOwnerSummary builds platform owner dashboard summary using Everflow data
func (s *dashboardService) getPlatformOwnerSummary(ctx context.Context, from, to time.Time) (*domain.PlatformOwnerDashboard, error) {
	// Get dashboard summary from Everflow
	dashboardSummary, err := s.everflowRepo.GetDashboardSummary(ctx, 67)
	if err != nil {
		return nil, fmt.Errorf("failed to get Everflow dashboard summary: %w", err)
	}

	// Create mock platform metrics based on Everflow data
	platformSummary := domain.PlatformSummary{
		TotalUsers:        1000, // Mock data
		TotalRevenue:      float64(dashboardSummary.Cost.CurrentMonth),
		TotalTransactions: int64(dashboardSummary.Conversion.CurrentMonth),
		PlatformGrowth:    dashboardSummary.Cost.TrendingPercentage,
	}

	userMetrics := domain.UserMetrics{
		ActiveUsers:     800,
		NewUsers:        50,
		UserGrowthRate:  5.2,
		UsersByType:     map[string]int{"advertiser": 300, "agency": 200, "affiliate": 500},
	}

	revenueMetrics := domain.RevenueMetrics{
		TotalRevenue:          float64(dashboardSummary.Cost.CurrentMonth),
		RevenueGrowth:         dashboardSummary.Cost.TrendingPercentage,
		AverageRevenuePerUser: float64(dashboardSummary.Cost.Today) / 800, // Average revenue per user
		RevenueBySource: []domain.RevenueBySource{
			{Source: "direct", Revenue: float64(dashboardSummary.Cost.Today) * 0.6, Percentage: 60.0},
			{Source: "referral", Revenue: float64(dashboardSummary.Cost.Today) * 0.4, Percentage: 40.0},
		},
	}

	systemHealth := domain.SystemHealth{
		Uptime:            99.9,
		ResponseTime:      120,
		ErrorRate:         0.1,
		ActiveConnections: 500,
	}

	return &domain.PlatformOwnerDashboard{
		Summary:        platformSummary,
		UserMetrics:    userMetrics,
		RevenueMetrics: revenueMetrics,
		SystemHealth:   systemHealth,
	}, nil
}



// GetCampaignDetail retrieves detailed campaign information
func (s *dashboardService) GetCampaignDetail(ctx context.Context, userID uuid.UUID, campaignID int64) (*domain.CampaignDetail, error) {
	// Get user profile to verify access
	profile, err := s.profileService.GetProfileByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	if profile.OrganizationID == nil {
		return nil, errors.New("user is not associated with any organization")
	}

	// For now, return mock campaign detail since we're focusing on Everflow data
	// In a real implementation, this would fetch from Everflow API
	endDate := time.Now().AddDate(0, 1, 0)
	detail := &domain.CampaignDetail{
		Campaign: domain.CampaignInfo{
			ID:        campaignID,
			Name:      fmt.Sprintf("Campaign %d", campaignID),
			Status:    "active",
			StartDate: time.Now().AddDate(0, -1, 0),
			EndDate:   &endDate,
			Budget:    10000.0,
			Spent:     5000.0,
		},
		Performance: domain.CampaignMetrics{
			Clicks:      5000,
			Impressions: 100000,
			Conversions: 250,
			Revenue:     12500.0,
			CTR:         5.0,
			CPC:         1.0,
			CPA:         20.0,
		},
		DailyStats: []domain.DailyStat{
			{
				Date:        time.Now().Format("2006-01-02"),
				Clicks:      1000,
				Impressions: 20000,
				Conversions: 50,
				Revenue:     2500.0,
			},
		},
	}

	return detail, nil
}

// GetRecentActivity retrieves recent activity for the user's organization
func (s *dashboardService) GetRecentActivity(ctx context.Context, userID uuid.UUID, limit, offset int, activityTypes []string, since *time.Time) (*domain.ActivityResponse, error) {
	// Get user profile
	profile, err := s.profileService.GetProfileByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	if profile.OrganizationID == nil {
		return nil, errors.New("user is not associated with any organization")
	}

	// Convert string activity types to domain types and back to valid strings
	var validActivityTypes []string
	for _, typeStr := range activityTypes {
		activityType := domain.ActivityType(typeStr)
		if activityType.IsValid() {
			validActivityTypes = append(validActivityTypes, string(activityType))
		}
	}

	// For now, return empty activities since we're focusing on Everflow data
	activities := []domain.Activity{}
	total := 0
	hasMore := false

	return &domain.ActivityResponse{
		Activities: activities,
		Total:      total,
		HasMore:    hasMore,
	}, nil
}

// GetSystemHealth retrieves system health metrics (platform owner only)
func (s *dashboardService) GetSystemHealth(ctx context.Context, userID uuid.UUID) (*domain.SystemHealth, error) {
	// Verify user has platform owner permissions
	profile, err := s.profileService.GetProfileByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	if profile.RoleName != "Admin" && profile.RoleName != "PlatformOwner" {
		return nil, errors.New("insufficient privileges for system health access")
	}

	// Return mock system health data
	health := &domain.SystemHealth{
		Uptime:            99.9,
		ResponseTime:      120,
		ErrorRate:         0.1,
		ActiveConnections: 500,
	}

	return health, nil
}

// TrackActivity creates a new activity record
func (s *dashboardService) TrackActivity(ctx context.Context, orgID int64, activityType domain.ActivityType, description string, metadata map[string]interface{}) error {
	// For now, just log the activity since we're not storing in database
	s.logger.Info("Activity tracked",
		"org_id", orgID,
		"type", activityType,
		"description", description,
	)

	// Invalidate cache for the organization
	if err := s.cacheRepo.InvalidateDashboardCache(ctx, orgID); err != nil {
		s.logger.Warn("Failed to invalidate dashboard cache", "org_id", orgID, "error", err)
	}

	return nil
}

// TrackActivityWithUser creates a new activity record with user context
func (s *dashboardService) TrackActivityWithUser(ctx context.Context, userID uuid.UUID, orgID int64, activityType domain.ActivityType, description string, metadata map[string]interface{}) error {
	// For now, just log the activity since we're not storing in database
	s.logger.Info("Activity tracked with user",
		"user_id", userID,
		"org_id", orgID,
		"type", activityType,
		"description", description,
	)

	// Invalidate cache for the organization
	if err := s.cacheRepo.InvalidateDashboardCache(ctx, orgID); err != nil {
		s.logger.Warn("Failed to invalidate dashboard cache", "org_id", orgID, "error", err)
	}

	return nil
}

// InvalidateCache invalidates dashboard cache for an organization
func (s *dashboardService) InvalidateCache(ctx context.Context, orgID int64) error {
	return s.cacheRepo.InvalidateDashboardCache(ctx, orgID)
}

// CheckDatabaseHealth checks database connectivity (not used with Everflow architecture)
func (s *dashboardService) CheckDatabaseHealth(ctx context.Context) error {
	// Since we're using Everflow API, we don't have database dependency
	return nil
}

// CheckCacheHealth checks cache connectivity
func (s *dashboardService) CheckCacheHealth(ctx context.Context) error {
	// Try to set and get a test value
	testData := &domain.DashboardData{
		OrganizationType: domain.OrganizationTypeAdvertiser,
		LastUpdated:      time.Now(),
	}

	if err := s.cacheRepo.SetCachedDashboardData(ctx, -1, domain.OrganizationTypeAdvertiser, testData, time.Minute); err != nil {
		return fmt.Errorf("cache health check failed (set): %w", err)
	}

	_, err := s.cacheRepo.GetCachedDashboardData(ctx, -1, domain.OrganizationTypeAdvertiser)
	if err != nil {
		return fmt.Errorf("cache health check failed (get): %w", err)
	}

	// Clean up test data
	_ = s.cacheRepo.InvalidateDashboardCache(ctx, -1)

	return nil
}