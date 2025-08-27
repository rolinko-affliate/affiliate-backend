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
	GetDashboardData(ctx context.Context, userID uuid.UUID, period string, startDate, endDate *time.Time, timezone string, advertiserIDs, campaignIDs, affiliateIDs []int64) (*domain.DashboardData, error)
	GetCampaignDetail(ctx context.Context, userID uuid.UUID, campaignID int64) (*domain.CampaignDetail, error)
	GetRecentActivity(ctx context.Context, userID uuid.UUID, limit, offset int, activityTypes []string, since *time.Time) (*domain.ActivityResponse, error)

	// Enhanced dashboard methods
	GetOffersPaginated(ctx context.Context, userID uuid.UUID, page, perPage int) (*domain.OffersPaginated, error)

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
	cacheRepo           repository.DashboardCacheRepository
	everflowRepo        repository.EverflowRepository
	reportingService    ReportingService
	profileService      ProfileService
	organizationService OrganizationService
	mockDataService     MockDataService
	logger              *logger.Logger
	mockMode            bool
}

// NewDashboardService creates a new dashboard service
func NewDashboardService(
	cacheRepo repository.DashboardCacheRepository,
	everflowRepo repository.EverflowRepository,
	reportingService ReportingService,
	profileService ProfileService,
	organizationService OrganizationService,
	mockDataService MockDataService,
	logger *logger.Logger,
	mockMode bool,
) DashboardService {
	return &dashboardService{
		cacheRepo:           cacheRepo,
		everflowRepo:        everflowRepo,
		reportingService:    reportingService,
		profileService:      profileService,
		organizationService: organizationService,
		mockDataService:     mockDataService,
		logger:              logger,
		mockMode:            mockMode,
	}
}

// GetDashboardData retrieves dashboard data based on user's organization type
func (s *dashboardService) GetDashboardData(ctx context.Context, userID uuid.UUID, period string, startDate, endDate *time.Time, timezone string, advertiserIDs, campaignIDs, affiliateIDs []int64) (*domain.DashboardData, error) {
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

	// Get date range
	from, to := query.GetDateRange()

	// For now, return empty activities since we're focusing on Everflow data
	activities := []domain.Activity{}

	// In mock mode, use default organization type (advertiser)
	if s.mockMode && s.mockDataService != nil {
		summary, err := s.getAdvertiserSummary(ctx, orgID, from, to)
		if err != nil {
			return nil, fmt.Errorf("failed to get dashboard summary: %w", err)
		}

		dashboardData := &domain.DashboardData{
			OrganizationType: domain.OrganizationTypeAdvertiser,
			Summary:          summary,
			RecentActivity:   activities,
			LastUpdated:      time.Now(),
		}

		return dashboardData, nil
	}

	// Get organization details
	org, err := s.organizationService.GetOrganizationByID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	// Check cache first
	cachedData, err := s.cacheRepo.GetCachedDashboardData(ctx, orgID, org.Type)
	if err != nil {
		s.logger.Warn("Failed to get cached dashboard data", "error", err)
	} else if cachedData != nil {
		s.logger.Debug("Returning cached dashboard data", "org_id", orgID)
		return cachedData, nil
	}

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
	// Use mock data if mock mode is enabled
	if s.mockMode && s.mockDataService != nil {
		return s.getAdvertiserSummaryFromMock(ctx, orgID)
	}

	// Get dashboard summary from Everflow
	// For timezone, we'll use UTC (timezone ID 67) as default
	dashboardSummary, err := s.everflowRepo.GetDashboardSummary(ctx, 67)
	if err != nil {
		return nil, fmt.Errorf("failed to get Everflow dashboard summary: %w", err)
	}

	// Convert Everflow data to our domain model with historical metrics
	advertiserMetrics := domain.AdvertiserMetrics{
		TotalClicks: domain.MetricWithHistory{
			Today:           float64(dashboardSummary.Click.Today),
			Yesterday:       float64(dashboardSummary.Click.Yesterday),
			CurrentMonth:    float64(dashboardSummary.Click.CurrentMonth),
			LastMonth:       float64(dashboardSummary.Click.LastMonth),
			ChangePercentage: calculateChangePercentage(float64(dashboardSummary.Click.Today), float64(dashboardSummary.Click.Yesterday)),
		},
		Conversions: domain.MetricWithHistory{
			Today:           float64(dashboardSummary.Conversion.Today),
			Yesterday:       float64(dashboardSummary.Conversion.Yesterday),
			CurrentMonth:    float64(dashboardSummary.Conversion.CurrentMonth),
			LastMonth:       float64(dashboardSummary.Conversion.LastMonth),
			ChangePercentage: calculateChangePercentage(float64(dashboardSummary.Conversion.Today), float64(dashboardSummary.Conversion.Yesterday)),
		},
		Revenue: domain.MetricWithHistory{
			Today:           float64(dashboardSummary.Cost.Today),
			Yesterday:       float64(dashboardSummary.Cost.Yesterday),
			CurrentMonth:    float64(dashboardSummary.Cost.CurrentMonth),
			LastMonth:       float64(dashboardSummary.Cost.LastMonth),
			ChangePercentage: calculateChangePercentage(float64(dashboardSummary.Cost.Today), float64(dashboardSummary.Cost.Yesterday)),
		},
		ConversionRate: domain.MetricWithHistory{
			Today:           dashboardSummary.CVR.Today,
			Yesterday:       dashboardSummary.CVR.Yesterday,
			CurrentMonth:    dashboardSummary.CVR.CurrentMonth,
			LastMonth:       dashboardSummary.CVR.LastMonth,
			ChangePercentage: calculateChangePercentage(dashboardSummary.CVR.Today, dashboardSummary.CVR.Yesterday),
		},
		Events: domain.MetricWithHistory{
			Today:           float64(dashboardSummary.Events.Today),
			Yesterday:       float64(dashboardSummary.Events.Yesterday),
			CurrentMonth:    float64(dashboardSummary.Events.CurrentMonth),
			LastMonth:       float64(dashboardSummary.Events.LastMonth),
			ChangePercentage: calculateChangePercentage(float64(dashboardSummary.Events.Today), float64(dashboardSummary.Events.Yesterday)),
		},
		EventRate: domain.MetricWithHistory{
			Today:           dashboardSummary.EVR.Today,
			Yesterday:       dashboardSummary.EVR.Yesterday,
			CurrentMonth:    dashboardSummary.EVR.CurrentMonth,
			LastMonth:       dashboardSummary.EVR.LastMonth,
			ChangePercentage: calculateChangePercentage(dashboardSummary.EVR.Today, dashboardSummary.EVR.Yesterday),
		},
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
		Metrics:             advertiserMetrics,
		CampaignPerformance: campaignPerformance,
		RevenueChart:        revenueChart,
		Billing:             billing,
	}, nil
}

// getAdvertiserSummaryFromMock builds advertiser dashboard using mock data
func (s *dashboardService) getAdvertiserSummaryFromMock(ctx context.Context, orgID int64) (*domain.AdvertiserDashboard, error) {
	// Load advertiser metrics with historical data
	metrics, err := s.mockDataService.LoadAdvertiserMetrics(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to load mock advertiser metrics: %w", err)
	}

	// Load campaign performance
	campaignPerformance, err := s.mockDataService.LoadCampaignPerformance(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to load mock campaign performance: %w", err)
	}

	// Load enhanced revenue chart (default to 30d period)
	revenueChart, err := s.mockDataService.LoadEnhancedRevenueChart(ctx, orgID, "30d")
	if err != nil {
		return nil, fmt.Errorf("failed to load mock revenue chart: %w", err)
	}

	// Load offers with pagination (first page, 10 items)
	offers, err := s.mockDataService.LoadOffersPaginated(ctx, orgID, 1, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to load mock offers: %w", err)
	}

	// Load billing info
	billing, err := s.mockDataService.LoadBillingInfo(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to load mock billing info: %w", err)
	}

	return &domain.AdvertiserDashboard{
		Metrics:             *metrics,
		CampaignPerformance: *campaignPerformance,
		RevenueChart:        *revenueChart,
		Offers:              offers,
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

// calculateChangePercentage calculates the percentage change between current and previous values
func calculateChangePercentage(current, previous float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100 // 100% increase from 0
	}
	return ((current - previous) / previous) * 100
}

// getAgencySummary builds agency-specific dashboard summary using Everflow data
func (s *dashboardService) getAgencySummary(ctx context.Context, orgID int64, from, to time.Time) (*domain.AgencyDashboard, error) {
	// Use mock data if mock mode is enabled
	if s.mockMode && s.mockDataService != nil {
		return s.getAgencySummaryFromMock(ctx, orgID)
	}

	// Get dashboard summary from Everflow
	dashboardSummary, err := s.everflowRepo.GetDashboardSummary(ctx, 67)
	if err != nil {
		return nil, fmt.Errorf("failed to get Everflow dashboard summary: %w", err)
	}

	// Create agency performance overview from Everflow data
	agencyPerformanceOverview := domain.AgencyPerformanceOverview{
		TotalConversions: domain.MetricWithHistory{
			Today:            float64(dashboardSummary.Conversion.Today),
			Yesterday:        float64(dashboardSummary.Conversion.Yesterday),
			CurrentMonth:     float64(dashboardSummary.Conversion.CurrentMonth),
			LastMonth:        float64(dashboardSummary.Conversion.LastMonth),
			ChangePercentage: dashboardSummary.Conversion.TrendingPercentage,
		},
		TotalClicks: domain.MetricWithHistory{
			Today:            float64(dashboardSummary.Click.Today),
			Yesterday:        float64(dashboardSummary.Click.Yesterday),
			CurrentMonth:     float64(dashboardSummary.Click.CurrentMonth),
			LastMonth:        float64(dashboardSummary.Click.LastMonth),
			ChangePercentage: dashboardSummary.Click.TrendingPercentage,
		},
		ConversionsPerDay: domain.MetricWithHistory{
			Today:            float64(dashboardSummary.Conversion.Today),
			Yesterday:        float64(dashboardSummary.Conversion.Yesterday),
			CurrentMonth:     float64(dashboardSummary.Conversion.CurrentMonth),
			LastMonth:        float64(dashboardSummary.Conversion.LastMonth),
			ChangePercentage: dashboardSummary.Conversion.TrendingPercentage,
		},
		ClicksPerDay: domain.MetricWithHistory{
			Today:            float64(dashboardSummary.Click.Today),
			Yesterday:        float64(dashboardSummary.Click.Yesterday),
			CurrentMonth:     float64(dashboardSummary.Click.CurrentMonth),
			LastMonth:        float64(dashboardSummary.Click.LastMonth),
			ChangePercentage: dashboardSummary.Click.TrendingPercentage,
		},
		TotalEarnings: domain.MetricWithHistory{
			Today:            float64(dashboardSummary.Cost.Today),
			Yesterday:        float64(dashboardSummary.Cost.Yesterday),
			CurrentMonth:     float64(dashboardSummary.Cost.CurrentMonth),
			LastMonth:        float64(dashboardSummary.Cost.LastMonth),
			ChangePercentage: dashboardSummary.Cost.TrendingPercentage,
		},
	}

	// Create agency performance chart from Everflow data
	agencyPerformanceChart := domain.AgencyPerformanceChart{
		Data: []domain.AgencyPerformanceDataPoint{
			{
				Date:        time.Now().AddDate(0, 0, -6).Format("2006-01-02"),
				Revenue:     float64(dashboardSummary.Cost.Today) * 0.8,
				Conversions: int64(dashboardSummary.Conversion.Today) * 80 / 100,
				Clicks:      int64(dashboardSummary.Click.Today) * 80 / 100,
			},
			{
				Date:        time.Now().AddDate(0, 0, -5).Format("2006-01-02"),
				Revenue:     float64(dashboardSummary.Cost.Today) * 0.9,
				Conversions: int64(dashboardSummary.Conversion.Today) * 90 / 100,
				Clicks:      int64(dashboardSummary.Click.Today) * 90 / 100,
			},
			{
				Date:        time.Now().Format("2006-01-02"),
				Revenue:     float64(dashboardSummary.Cost.Today),
				Conversions: int64(dashboardSummary.Conversion.Today),
				Clicks:      int64(dashboardSummary.Click.Today),
			},
		},
		Period: "7d",
	}

	advertiserOrgs := []domain.AdvertiserOrganization{
		{
			ID:             1,
			Name:           "Client A",
			Status:         "active",
			CampaignsCount: 3,
			Revenue:        float64(dashboardSummary.Cost.Today) * 0.4,
			ConversionRate: dashboardSummary.CVR.Today,
		},
		{
			ID:             2,
			Name:           "Client B", 
			Status:         "active",
			CampaignsCount: 2,
			Revenue:        float64(dashboardSummary.Cost.Today) * 0.3,
			ConversionRate: dashboardSummary.CVR.Today,
		},
	}

	campaignsOverview := domain.CampaignsOverview{
		Campaigns: []domain.AgencyCampaignSummary{
			{
				ID:               100,
				Name:             "Campaign A",
				AdvertiserOrgID:  1,
				AdvertiserName:   "Client A",
				Status:           "active",
				Clicks:           int64(dashboardSummary.Click.Today) * 40 / 100,
				Conversions:      int64(dashboardSummary.Conversion.Today) * 40 / 100,
				TotalCost:        float64(dashboardSummary.Cost.Today) * 0.4,
				ConversionRate:   dashboardSummary.CVR.Today,
			},
			{
				ID:               101,
				Name:             "Campaign B",
				AdvertiserOrgID:  2,
				AdvertiserName:   "Client B",
				Status:           "active",
				Clicks:           int64(dashboardSummary.Click.Today) * 30 / 100,
				Conversions:      int64(dashboardSummary.Conversion.Today) * 30 / 100,
				TotalCost:        float64(dashboardSummary.Cost.Today) * 0.3,
				ConversionRate:   dashboardSummary.CVR.Today,
			},
		},
		TotalCount:  2,
		ActiveCount: 2,
	}



	return &domain.AgencyDashboard{
		PerformanceOverview:     agencyPerformanceOverview,
		AdvertiserOrganizations: advertiserOrgs,
		CampaignsOverview:       campaignsOverview,
		PerformanceChart:        agencyPerformanceChart,
	}, nil
}

// getAgencySummaryFromMock builds agency dashboard using mock data
func (s *dashboardService) getAgencySummaryFromMock(ctx context.Context, orgID int64) (*domain.AgencyDashboard, error) {
	// Load agency performance overview
	performanceOverview, err := s.mockDataService.LoadAgencyPerformanceOverview(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to load mock agency performance overview: %w", err)
	}

	// Load advertiser organizations
	advertiserOrgs, err := s.mockDataService.LoadAdvertiserOrganizations(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to load mock advertiser organizations: %w", err)
	}

	// Load campaigns overview
	campaignsOverview, err := s.mockDataService.LoadCampaignsOverview(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to load mock campaigns overview: %w", err)
	}

	// Load agency performance chart (default to 30d period)
	performanceChart, err := s.mockDataService.LoadAgencyPerformanceChart(ctx, orgID, "30d")
	if err != nil {
		return nil, fmt.Errorf("failed to load mock agency performance chart: %w", err)
	}

	return &domain.AgencyDashboard{
		PerformanceOverview:      *performanceOverview,
		AdvertiserOrganizations:  advertiserOrgs,
		CampaignsOverview:        *campaignsOverview,
		PerformanceChart:         *performanceChart,
	}, nil
}

// getPlatformOwnerSummary builds platform owner dashboard summary using Everflow data
func (s *dashboardService) getPlatformOwnerSummary(ctx context.Context, from, to time.Time) (*domain.PlatformOwnerDashboard, error) {
	// Use mock data if mock mode is enabled
	if s.mockMode && s.mockDataService != nil {
		return s.getPlatformOwnerSummaryFromMock(ctx)
	}

	// Get dashboard summary from Everflow
	dashboardSummary, err := s.everflowRepo.GetDashboardSummary(ctx, 67)
	if err != nil {
		return nil, fmt.Errorf("failed to get Everflow dashboard summary: %w", err)
	}

	// Create platform overview from Everflow data
	platformOverview := domain.PlatformOverview{
		TotalOrganizations: domain.MetricWithHistory{
			Today:            10, // Mock data
			Yesterday:        9,
			CurrentMonth:     10,
			LastMonth:        9,
			ChangePercentage: 11.1,
		},
		TotalUsers: domain.MetricWithHistory{
			Today:            1000, // Mock data
			Yesterday:        995,
			CurrentMonth:     1000,
			LastMonth:        950,
			ChangePercentage: 5.3,
		},
		TotalRevenue: domain.MetricWithHistory{
			Today:            float64(dashboardSummary.Cost.Today),
			Yesterday:        float64(dashboardSummary.Cost.Yesterday),
			CurrentMonth:     float64(dashboardSummary.Cost.CurrentMonth),
			LastMonth:        float64(dashboardSummary.Cost.LastMonth),
			ChangePercentage: dashboardSummary.Cost.TrendingPercentage,
		},
		MonthlyGrowth: domain.MetricWithHistory{
			Today:            5.2, // Mock data
			Yesterday:        5.1,
			CurrentMonth:     5.2,
			LastMonth:        4.8,
			ChangePercentage: 8.3,
		},
		NewRegistrations: domain.MetricWithHistory{
			Today:            50, // Mock data
			Yesterday:        45,
			CurrentMonth:     1500,
			LastMonth:        1400,
			ChangePercentage: 7.1,
		},
	}

	userActivityMetrics := domain.UserActivityMetrics{
		DailyActiveUsers:   800,
		WeeklyActiveUsers:  5600,
		MonthlyActiveUsers: 24000,
		ActiveAdvertisers:  300,
		ActiveAffiliates:   500,
	}

	systemHealthMetrics := domain.SystemHealthMetrics{
		TotalCampaigns:         50,
		RequestsPerMinute:      1200,
		SuccessRate:            99.9,
		RateLimitHits:          5,
		AverageQueryTime:       120.5,
		ConnectionPoolUsage:    75.0,
	}

	return &domain.PlatformOwnerDashboard{
		PlatformOverview: platformOverview,
		UserActivity:     userActivityMetrics,
		SystemHealth:     systemHealthMetrics,
		RevenueBySource: []domain.RevenueBySource{
			{Source: "direct", Revenue: float64(dashboardSummary.Cost.Today) * 0.6, Percentage: 60.0, Growth: 5.2},
			{Source: "referral", Revenue: float64(dashboardSummary.Cost.Today) * 0.4, Percentage: 40.0, Growth: 3.1},
		},
		GeographicDistribution: []domain.GeographicData{
			{Country: "US", Revenue: float64(dashboardSummary.Cost.Today) * 0.5, Users: 400, ConversionRate: dashboardSummary.CVR.Today},
			{Country: "CA", Revenue: float64(dashboardSummary.Cost.Today) * 0.3, Users: 240, ConversionRate: dashboardSummary.CVR.Today * 0.9},
			{Country: "UK", Revenue: float64(dashboardSummary.Cost.Today) * 0.2, Users: 160, ConversionRate: dashboardSummary.CVR.Today * 0.8},
		},
	}, nil
}

// getPlatformOwnerSummaryFromMock builds platform owner dashboard using mock data
func (s *dashboardService) getPlatformOwnerSummaryFromMock(ctx context.Context) (*domain.PlatformOwnerDashboard, error) {
	// Load platform overview with historical data
	platformOverview, err := s.mockDataService.LoadPlatformOverview(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load mock platform overview: %w", err)
	}

	// Load user activity metrics
	userActivity, err := s.mockDataService.LoadUserActivityMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load mock user activity metrics: %w", err)
	}

	// Load system health metrics
	systemHealth, err := s.mockDataService.LoadSystemHealthMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load mock system health metrics: %w", err)
	}

	// Load revenue by source (use platform orgID 100)
	revenueBySource, err := s.mockDataService.LoadRevenueBySource(ctx, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to load mock revenue by source: %w", err)
	}

	// Load geographic distribution (use platform orgID 100)
	geographicDistribution, err := s.mockDataService.LoadGeographicDistribution(ctx, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to load mock geographic distribution: %w", err)
	}

	return &domain.PlatformOwnerDashboard{
		PlatformOverview:       *platformOverview,
		UserActivity:           *userActivity,
		SystemHealth:           *systemHealth,
		RevenueBySource:        revenueBySource,
		GeographicDistribution: geographicDistribution,
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

	// Use mock data if mock mode is enabled
	if s.mockMode && s.mockDataService != nil {
		return s.mockDataService.LoadCampaignDetail(ctx, campaignID, *profile.OrganizationID)
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

	// Use mock data if mock mode is enabled
	if s.mockMode && s.mockDataService != nil {
		return s.mockDataService.LoadActivities(ctx, *profile.OrganizationID, limit, offset, activityTypes, since)
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

	// Use mock data if mock mode is enabled
	if s.mockMode && s.mockDataService != nil {
		return s.mockDataService.LoadSystemHealth(ctx)
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

// GetOffersPaginated retrieves paginated offers for the user's organization
func (s *dashboardService) GetOffersPaginated(ctx context.Context, userID uuid.UUID, page, perPage int) (*domain.OffersPaginated, error) {
	// Get user profile to determine organization
	profile, err := s.profileService.GetProfileByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	if profile.OrganizationID == nil {
		return nil, errors.New("user is not associated with any organization")
	}

	orgID := *profile.OrganizationID

	// Use mock data if mock mode is enabled
	if s.mockMode && s.mockDataService != nil {
		return s.mockDataService.LoadOffersPaginated(ctx, orgID, page, perPage)
	}

	// For non-mock mode, return empty offers for now
	return &domain.OffersPaginated{
		Items:      []domain.Offer{},
		TotalCount: 0,
		Page:       page,
		PerPage:    perPage,
		HasNext:    false,
	}, nil
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
