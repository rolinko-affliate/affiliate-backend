package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/logger"
)

// MockDataService provides methods to load mock data from CSV files
type MockDataService interface {
	// Dashboard data methods - New enhanced methods
	LoadAdvertiserMetrics(ctx context.Context, orgID int64) (*domain.AdvertiserMetrics, error)
	LoadOffersPaginated(ctx context.Context, orgID int64, page, perPage int) (*domain.OffersPaginated, error)
	LoadAgencyPerformanceOverview(ctx context.Context, orgID int64) (*domain.AgencyPerformanceOverview, error)
	LoadAdvertiserOrganizations(ctx context.Context, orgID int64) ([]domain.AdvertiserOrganization, error)
	LoadCampaignsOverview(ctx context.Context, orgID int64) (*domain.CampaignsOverview, error)
	LoadAgencyPerformanceChart(ctx context.Context, orgID int64, period string) (*domain.AgencyPerformanceChart, error)
	LoadPlatformOverview(ctx context.Context) (*domain.PlatformOverview, error)
	LoadUserActivityMetrics(ctx context.Context) (*domain.UserActivityMetrics, error)
	LoadSystemHealthMetrics(ctx context.Context) (*domain.SystemHealthMetrics, error)
	LoadRevenueBySource(ctx context.Context, orgID int64) ([]domain.RevenueBySource, error)
	LoadGeographicDistribution(ctx context.Context, orgID int64) ([]domain.GeographicData, error)
	LoadEnhancedRevenueChart(ctx context.Context, orgID int64, period string) (*domain.RevenueChart, error)
	
	// Dashboard data methods - Legacy support
	LoadAdvertiserSummary(ctx context.Context, orgID int64) (*domain.AdvertiserSummary, error)
	LoadCampaignPerformance(ctx context.Context, orgID int64) (*domain.CampaignPerformance, error)
	LoadRevenueChart(ctx context.Context, orgID int64, period string) (*domain.RevenueChart, error)
	LoadBillingInfo(ctx context.Context, orgID int64) (*domain.BillingInfo, error)
	LoadAgencySummary(ctx context.Context, orgID int64) (*domain.AgencySummary, error)
	LoadClientPerformance(ctx context.Context, orgID int64) (*domain.ClientPerformance, error)
	LoadAgencyRevenueChart(ctx context.Context, orgID int64, period string) (*domain.AgencyRevenueChart, error)
	LoadPlatformSummary(ctx context.Context) (*domain.PlatformSummary, error)
	LoadUserMetrics(ctx context.Context) (*domain.UserMetrics, error)
	LoadRevenueMetrics(ctx context.Context) (*domain.RevenueMetrics, error)
	LoadSystemHealth(ctx context.Context) (*domain.SystemHealth, error)
	LoadActivities(ctx context.Context, orgID int64, limit, offset int, activityTypes []string, since *time.Time) (*domain.ActivityResponse, error)
	LoadCampaignDetail(ctx context.Context, campaignID int64, orgID int64) (*domain.CampaignDetail, error)

	// Reporting data methods
	LoadPerformanceSummary(ctx context.Context, orgID int64, filters domain.ReportingFilters) (*domain.PerformanceSummary, error)
	LoadPerformanceTimeSeries(ctx context.Context, orgID int64, filters domain.ReportingFilters) ([]domain.PerformanceTimeSeriesPoint, error)
	LoadDailyPerformanceReport(ctx context.Context, orgID int64, filters domain.ReportingFilters, pagination domain.PaginationParams) ([]domain.DailyPerformanceReport, *domain.PaginationResult, error)
	LoadConversionsReport(ctx context.Context, orgID int64, filters domain.ReportingFilters, pagination domain.PaginationParams) ([]domain.ConversionReport, *domain.PaginationResult, error)
	LoadClicksReport(ctx context.Context, orgID int64, filters domain.ReportingFilters, pagination domain.PaginationParams) ([]domain.ClickReport, *domain.PaginationResult, error)
	LoadCampaignsList(ctx context.Context, orgID int64, affiliateID *string, status string, search *string) ([]domain.CampaignListItem, error)
}

type mockDataService struct {
	logger   *logger.Logger
	dataPath string
}

// NewMockDataService creates a new mock data service
func NewMockDataService(logger *logger.Logger) MockDataService {
	return &mockDataService{
		logger:   logger,
		dataPath: "mock-data",
	}
}

// LoadAdvertiserSummary loads advertiser summary data from CSV
func (s *mockDataService) LoadAdvertiserSummary(ctx context.Context, orgID int64) (*domain.AdvertiserSummary, error) {
	records, err := s.loadCSV("advertiser_summary.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load advertiser summary: %w", err)
	}

	for _, record := range records[1:] { // Skip header
		if len(record) < 5 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if recordOrgID == orgID {
			totalClicks, _ := strconv.ParseInt(record[1], 10, 64)
			conversions, _ := strconv.ParseInt(record[2], 10, 64)
			revenue, _ := strconv.ParseFloat(record[3], 64)
			conversionRate, _ := strconv.ParseFloat(record[4], 64)

			return &domain.AdvertiserSummary{
				TotalClicks:    totalClicks,
				Conversions:    conversions,
				Revenue:        revenue,
				ConversionRate: conversionRate,
			}, nil
		}
	}

	// Return default data if not found
	return &domain.AdvertiserSummary{
		TotalClicks:    10000,
		Conversions:    800,
		Revenue:        25000.00,
		ConversionRate: 8.0,
	}, nil
}

// LoadCampaignPerformance loads campaign performance data from CSV
func (s *mockDataService) LoadCampaignPerformance(ctx context.Context, orgID int64) (*domain.CampaignPerformance, error) {
	records, err := s.loadCSV("campaign_performance.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load campaign performance: %w", err)
	}

	var campaigns []domain.CampaignSummary
	totalCampaigns := 0
	activeCampaigns := 0

	for _, record := range records[1:] { // Skip header
		if len(record) < 10 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			continue
		}

		if recordOrgID == orgID {
			campaignID, _ := strconv.ParseInt(record[0], 10, 64)
			clicks, _ := strconv.ParseInt(record[3], 10, 64)
			conversions, _ := strconv.ParseInt(record[4], 10, 64)
			revenue, _ := strconv.ParseFloat(record[5], 64)
			conversionRate, _ := strconv.ParseFloat(record[6], 64)
			totalCampaigns, _ = strconv.Atoi(record[8])
			activeCampaigns, _ = strconv.Atoi(record[9])

			campaign := domain.CampaignSummary{
				ID:             campaignID,
				Name:           record[2],
				Clicks:         clicks,
				Conversions:    conversions,
				Revenue:        revenue,
				ConversionRate: conversionRate,
				Status:         record[7],
			}
			campaigns = append(campaigns, campaign)
		}
	}

	if len(campaigns) == 0 {
		// Return default data if not found
		campaigns = []domain.CampaignSummary{
			{
				ID:             1,
				Name:           "Default Campaign",
				Clicks:         5000,
				Conversions:    400,
				Revenue:        12000.00,
				ConversionRate: 8.0,
				Status:         "active",
			},
		}
		totalCampaigns = 1
		activeCampaigns = 1
	}

	return &domain.CampaignPerformance{
		Campaigns:       campaigns,
		TotalCampaigns:  totalCampaigns,
		ActiveCampaigns: activeCampaigns,
	}, nil
}

// LoadRevenueChart loads revenue chart data from CSV
func (s *mockDataService) LoadRevenueChart(ctx context.Context, orgID int64, period string) (*domain.RevenueChart, error) {
	records, err := s.loadCSV("revenue_chart_data.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load revenue chart: %w", err)
	}

	var dataPoints []domain.RevenueDataPoint

	for _, record := range records[1:] { // Skip header
		if len(record) < 6 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if recordOrgID == orgID && record[5] == period {
			revenue, _ := strconv.ParseFloat(record[2], 64)
			clicks, _ := strconv.ParseInt(record[3], 10, 64)
			conversions, _ := strconv.ParseInt(record[4], 10, 64)

			dataPoint := domain.RevenueDataPoint{
				Date:        record[1],
				Revenue:     revenue,
				Clicks:      clicks,
				Conversions: conversions,
			}
			dataPoints = append(dataPoints, dataPoint)
		}
	}

	if len(dataPoints) == 0 {
		// Return default data if not found
		now := time.Now()
		for i := 6; i >= 0; i-- {
			date := now.AddDate(0, 0, -i)
			dataPoints = append(dataPoints, domain.RevenueDataPoint{
				Date:        date.Format("2006-01-02"),
				Revenue:     1000.0 + float64(i*100),
				Clicks:      500 + int64(i*50),
				Conversions: 40 + int64(i*5),
			})
		}
	}

	return &domain.RevenueChart{
		Data:   dataPoints,
		Period: period,
	}, nil
}

// LoadBillingInfo loads billing information from CSV
func (s *mockDataService) LoadBillingInfo(ctx context.Context, orgID int64) (*domain.BillingInfo, error) {
	records, err := s.loadCSV("billing_info.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load billing info: %w", err)
	}

	for _, record := range records[1:] { // Skip header
		if len(record) < 4 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if recordOrgID == orgID {
			currentBalance, _ := strconv.ParseFloat(record[1], 64)
			monthlySpend, _ := strconv.ParseFloat(record[2], 64)

			return &domain.BillingInfo{
				CurrentBalance: currentBalance,
				MonthlySpend:   monthlySpend,
				Currency:       record[3],
			}, nil
		}
	}

	// Return default data if not found
	return &domain.BillingInfo{
		CurrentBalance: 1500.00,
		MonthlySpend:   5000.00,
		Currency:       "USD",
	}, nil
}

// LoadAgencySummary loads agency summary data from CSV
func (s *mockDataService) LoadAgencySummary(ctx context.Context, orgID int64) (*domain.AgencySummary, error) {
	records, err := s.loadCSV("agency_summary.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load agency summary: %w", err)
	}

	for _, record := range records[1:] { // Skip header
		if len(record) < 5 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if recordOrgID == orgID {
			totalClients, _ := strconv.Atoi(record[1])
			totalRevenue, _ := strconv.ParseFloat(record[2], 64)
			totalConversions, _ := strconv.ParseInt(record[3], 10, 64)
			avgConversionRate, _ := strconv.ParseFloat(record[4], 64)

			return &domain.AgencySummary{
				TotalClients:          totalClients,
				TotalRevenue:          totalRevenue,
				TotalConversions:      totalConversions,
				AverageConversionRate: avgConversionRate,
			}, nil
		}
	}

	// Return default data if not found
	return &domain.AgencySummary{
		TotalClients:          5,
		TotalRevenue:          75000.00,
		TotalConversions:      2500,
		AverageConversionRate: 7.5,
	}, nil
}

// LoadClientPerformance loads client performance data from CSV
func (s *mockDataService) LoadClientPerformance(ctx context.Context, orgID int64) (*domain.ClientPerformance, error) {
	clientRecords, err := s.loadCSV("client_performance.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load client performance: %w", err)
	}

	topPerformerRecords, err := s.loadCSV("top_performers.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load top performers: %w", err)
	}

	var clients []domain.ClientSummary
	var topPerformers []domain.TopPerformer

	// Load clients
	for _, record := range clientRecords[1:] { // Skip header
		if len(record) < 7 {
			continue
		}

		agencyID, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			continue
		}

		if agencyID == orgID {
			clientID, _ := strconv.ParseInt(record[0], 10, 64)
			revenue, _ := strconv.ParseFloat(record[3], 64)
			conversions, _ := strconv.ParseInt(record[4], 10, 64)
			conversionRate, _ := strconv.ParseFloat(record[5], 64)
			activeCampaigns, _ := strconv.Atoi(record[6])

			client := domain.ClientSummary{
				ID:              clientID,
				Name:            record[2],
				Revenue:         revenue,
				Conversions:     conversions,
				ConversionRate:  conversionRate,
				ActiveCampaigns: activeCampaigns,
			}
			clients = append(clients, client)
		}
	}

	// Load top performers
	for _, record := range topPerformerRecords[1:] { // Skip header
		if len(record) < 5 {
			continue
		}

		agencyID, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			continue
		}

		if agencyID == orgID {
			clientID, _ := strconv.ParseInt(record[0], 10, 64)
			revenue, _ := strconv.ParseFloat(record[3], 64)
			growth, _ := strconv.ParseFloat(record[4], 64)

			topPerformer := domain.TopPerformer{
				ClientID:   clientID,
				ClientName: record[2],
				Revenue:    revenue,
				Growth:     growth,
			}
			topPerformers = append(topPerformers, topPerformer)
		}
	}

	if len(clients) == 0 {
		// Return default data if not found
		clients = []domain.ClientSummary{
			{
				ID:              1,
				Name:            "Default Client",
				Revenue:         15000.00,
				Conversions:     300,
				ConversionRate:  7.5,
				ActiveCampaigns: 2,
			},
		}
	}

	if len(topPerformers) == 0 {
		topPerformers = []domain.TopPerformer{
			{
				ClientID:   1,
				ClientName: "Default Client",
				Revenue:    15000.00,
				Growth:     10.5,
			},
		}
	}

	return &domain.ClientPerformance{
		Clients:       clients,
		TopPerformers: topPerformers,
	}, nil
}

// LoadAgencyRevenueChart loads agency revenue chart data from CSV
func (s *mockDataService) LoadAgencyRevenueChart(ctx context.Context, orgID int64, period string) (*domain.AgencyRevenueChart, error) {
	records, err := s.loadCSV("agency_revenue_chart_data.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load agency revenue chart: %w", err)
	}

	dataPointsMap := make(map[string]*domain.AgencyRevenueDataPoint)

	for _, record := range records[1:] { // Skip header
		if len(record) < 7 {
			continue
		}

		agencyID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if agencyID == orgID && record[6] == period {
			date := record[1]
			totalRevenue, _ := strconv.ParseFloat(record[2], 64)
			clientID, _ := strconv.ParseInt(record[3], 10, 64)
			clientRevenue, _ := strconv.ParseFloat(record[5], 64)

			if dataPoint, exists := dataPointsMap[date]; exists {
				// Add client breakdown to existing data point
				breakdown := domain.ClientRevenueBreakdown{
					ClientID:   clientID,
					ClientName: record[4],
					Revenue:    clientRevenue,
				}
				dataPoint.ClientBreakdown = append(dataPoint.ClientBreakdown, breakdown)
			} else {
				// Create new data point
				breakdown := domain.ClientRevenueBreakdown{
					ClientID:   clientID,
					ClientName: record[4],
					Revenue:    clientRevenue,
				}
				dataPointsMap[date] = &domain.AgencyRevenueDataPoint{
					Date:            date,
					TotalRevenue:    totalRevenue,
					ClientBreakdown: []domain.ClientRevenueBreakdown{breakdown},
				}
			}
		}
	}

	var dataPoints []domain.AgencyRevenueDataPoint
	for _, dataPoint := range dataPointsMap {
		dataPoints = append(dataPoints, *dataPoint)
	}

	if len(dataPoints) == 0 {
		// Return default data if not found
		now := time.Now()
		for i := 6; i >= 0; i-- {
			date := now.AddDate(0, 0, -i)
			dataPoints = append(dataPoints, domain.AgencyRevenueDataPoint{
				Date:         date.Format("2006-01-02"),
				TotalRevenue: 2000.0 + float64(i*200),
				ClientBreakdown: []domain.ClientRevenueBreakdown{
					{
						ClientID:   1,
						ClientName: "Default Client",
						Revenue:    2000.0 + float64(i*200),
					},
				},
			})
		}
	}

	return &domain.AgencyRevenueChart{
		Data:   dataPoints,
		Period: period,
	}, nil
}

// LoadPlatformSummary loads platform summary data from CSV
func (s *mockDataService) LoadPlatformSummary(ctx context.Context) (*domain.PlatformSummary, error) {
	records, err := s.loadCSV("platform_summary.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load platform summary: %w", err)
	}

	if len(records) > 1 && len(records[1]) >= 4 {
		record := records[1]
		totalUsers, _ := strconv.Atoi(record[0])
		totalRevenue, _ := strconv.ParseFloat(record[1], 64)
		totalTransactions, _ := strconv.ParseInt(record[2], 10, 64)
		platformGrowth, _ := strconv.ParseFloat(record[3], 64)

		return &domain.PlatformSummary{
			TotalUsers:        totalUsers,
			TotalRevenue:      totalRevenue,
			TotalTransactions: totalTransactions,
			PlatformGrowth:    platformGrowth,
		}, nil
	}

	// Return default data if not found
	return &domain.PlatformSummary{
		TotalUsers:        1000,
		TotalRevenue:      500000.00,
		TotalTransactions: 25000,
		PlatformGrowth:    15.5,
	}, nil
}

// LoadUserMetrics loads user metrics data from CSV
func (s *mockDataService) LoadUserMetrics(ctx context.Context) (*domain.UserMetrics, error) {
	records, err := s.loadCSV("user_metrics.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load user metrics: %w", err)
	}

	usersByType := make(map[string]int)
	var activeUsers, newUsers int
	var userGrowthRate float64

	for _, record := range records[1:] { // Skip header
		if len(record) < 6 {
			continue
		}

		activeUsers, _ = strconv.Atoi(record[0])
		newUsers, _ = strconv.Atoi(record[1])
		userGrowthRate, _ = strconv.ParseFloat(record[2], 64)
		userType := record[3]
		userCount, _ := strconv.Atoi(record[4])

		usersByType[userType] = userCount
	}

	if len(usersByType) == 0 {
		// Return default data if not found
		usersByType = map[string]int{
			"advertiser": 300,
			"agency":     200,
			"affiliate":  500,
		}
		activeUsers = 800
		newUsers = 50
		userGrowthRate = 8.5
	}

	return &domain.UserMetrics{
		ActiveUsers:    activeUsers,
		NewUsers:       newUsers,
		UserGrowthRate: userGrowthRate,
		UsersByType:    usersByType,
	}, nil
}

// LoadRevenueMetrics loads revenue metrics data from CSV
func (s *mockDataService) LoadRevenueMetrics(ctx context.Context) (*domain.RevenueMetrics, error) {
	records, err := s.loadCSV("revenue_metrics.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load revenue metrics: %w", err)
	}

	var totalRevenue, revenueGrowth, avgRevenuePerUser float64
	var revenueBySource []domain.RevenueBySource

	for _, record := range records[1:] { // Skip header
		if len(record) < 6 {
			continue
		}

		totalRevenue, _ = strconv.ParseFloat(record[0], 64)
		revenueGrowth, _ = strconv.ParseFloat(record[1], 64)
		avgRevenuePerUser, _ = strconv.ParseFloat(record[2], 64)

		sourceRevenue, _ := strconv.ParseFloat(record[4], 64)
		sourcePercentage, _ := strconv.ParseFloat(record[5], 64)

		revenueBySource = append(revenueBySource, domain.RevenueBySource{
			Source:     record[3],
			Revenue:    sourceRevenue,
			Percentage: sourcePercentage,
		})
	}

	if len(revenueBySource) == 0 {
		// Return default data if not found
		totalRevenue = 500000.00
		revenueGrowth = 15.3
		avgRevenuePerUser = 455.67
		revenueBySource = []domain.RevenueBySource{
			{Source: "Advertising Fees", Revenue: 300000.00, Percentage: 60.0},
			{Source: "Transaction Fees", Revenue: 100000.00, Percentage: 20.0},
			{Source: "Subscription Fees", Revenue: 75000.00, Percentage: 15.0},
			{Source: "Other", Revenue: 25000.00, Percentage: 5.0},
		}
	}

	return &domain.RevenueMetrics{
		TotalRevenue:          totalRevenue,
		RevenueGrowth:         revenueGrowth,
		AverageRevenuePerUser: avgRevenuePerUser,
		RevenueBySource:       revenueBySource,
	}, nil
}

// LoadSystemHealth loads system health data from CSV
func (s *mockDataService) LoadSystemHealth(ctx context.Context) (*domain.SystemHealth, error) {
	records, err := s.loadCSV("system_health.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load system health: %w", err)
	}

	if len(records) > 1 && len(records[1]) >= 4 {
		record := records[1]
		uptime, _ := strconv.ParseFloat(record[0], 64)
		responseTime, _ := strconv.ParseFloat(record[1], 64)
		errorRate, _ := strconv.ParseFloat(record[2], 64)
		activeConnections, _ := strconv.Atoi(record[3])

		return &domain.SystemHealth{
			Uptime:            uptime,
			ResponseTime:      responseTime,
			ErrorRate:         errorRate,
			ActiveConnections: activeConnections,
		}, nil
	}

	// Return default data if not found
	return &domain.SystemHealth{
		Uptime:            99.9,
		ResponseTime:      120.0,
		ErrorRate:         0.1,
		ActiveConnections: 500,
	}, nil
}

// LoadActivities loads activity data from CSV
func (s *mockDataService) LoadActivities(ctx context.Context, orgID int64, limit, offset int, activityTypes []string, since *time.Time) (*domain.ActivityResponse, error) {
	records, err := s.loadCSV("activities.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load activities: %w", err)
	}

	var activities []domain.Activity
	total := 0

	for _, record := range records[1:] { // Skip header
		if len(record) < 11 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			continue
		}

		if recordOrgID == orgID {
			total++

			// Apply pagination
			if total <= offset {
				continue
			}
			if len(activities) >= limit {
				break
			}

			timestamp, _ := time.Parse(time.RFC3339, record[4])

			var campaignID *int64
			var clientID *int64
			var userID *uuid.UUID

			if record[5] != "" {
				cid, _ := strconv.ParseInt(record[5], 10, 64)
				campaignID = &cid
			}
			if record[7] != "" {
				cid, _ := strconv.ParseInt(record[7], 10, 64)
				clientID = &cid
			}
			if record[10] != "" {
				uid, _ := uuid.Parse(record[10])
				userID = &uid
			}

			var campaignName, clientName, severity *string
			if record[6] != "" {
				campaignName = &record[6]
			}
			if record[8] != "" {
				clientName = &record[8]
			}
			if record[9] != "" {
				severity = &record[9]
			}

			activity := domain.Activity{
				ID:             record[0],
				Type:           domain.ActivityType(record[2]),
				Description:    record[3],
				Timestamp:      timestamp,
				CampaignID:     campaignID,
				CampaignName:   campaignName,
				ClientID:       clientID,
				ClientName:     clientName,
				Severity:       severity,
				OrganizationID: recordOrgID,
				UserID:         userID,
			}
			activities = append(activities, activity)
		}
	}

	hasMore := total > offset+len(activities)

	return &domain.ActivityResponse{
		Activities: activities,
		Total:      total,
		HasMore:    hasMore,
	}, nil
}

// LoadCampaignDetail loads campaign detail data from CSV
func (s *mockDataService) LoadCampaignDetail(ctx context.Context, campaignID int64, orgID int64) (*domain.CampaignDetail, error) {
	campaignRecords, err := s.loadCSV("campaign_details.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load campaign details: %w", err)
	}

	dailyStatsRecords, err := s.loadCSV("daily_stats.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load daily stats: %w", err)
	}

	// Find campaign info
	var campaignInfo *domain.CampaignInfo
	var campaignMetrics *domain.CampaignMetrics

	for _, record := range campaignRecords[1:] { // Skip header
		if len(record) < 20 {
			continue
		}

		recordCampaignID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			continue
		}

		if recordCampaignID == campaignID && recordOrgID == orgID {
			startDate, _ := time.Parse(time.RFC3339, record[4])
			var endDate *time.Time
			if record[5] != "" {
				ed, _ := time.Parse(time.RFC3339, record[5])
				endDate = &ed
			}
			budget, _ := strconv.ParseFloat(record[6], 64)
			spent, _ := strconv.ParseFloat(record[7], 64)
			clicks, _ := strconv.ParseInt(record[8], 10, 64)
			impressions, _ := strconv.ParseInt(record[9], 10, 64)
			conversions, _ := strconv.ParseInt(record[10], 10, 64)
			revenue, _ := strconv.ParseFloat(record[11], 64)
			ctr, _ := strconv.ParseFloat(record[12], 64)
			cpc, _ := strconv.ParseFloat(record[13], 64)
			cpa, _ := strconv.ParseFloat(record[14], 64)

			campaignInfo = &domain.CampaignInfo{
				ID:        campaignID,
				Name:      record[2],
				Status:    record[3],
				StartDate: startDate,
				EndDate:   endDate,
				Budget:    budget,
				Spent:     spent,
			}

			campaignMetrics = &domain.CampaignMetrics{
				Clicks:      clicks,
				Impressions: impressions,
				Conversions: conversions,
				Revenue:     revenue,
				CTR:         ctr,
				CPC:         cpc,
				CPA:         cpa,
			}
			break
		}
	}

	if campaignInfo == nil {
		return nil, fmt.Errorf("campaign not found")
	}

	// Load daily stats
	var dailyStats []domain.DailyStat
	for _, record := range dailyStatsRecords[1:] { // Skip header
		if len(record) < 6 {
			continue
		}

		recordCampaignID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if recordCampaignID == campaignID {
			clicks, _ := strconv.ParseInt(record[2], 10, 64)
			impressions, _ := strconv.ParseInt(record[3], 10, 64)
			conversions, _ := strconv.ParseInt(record[4], 10, 64)
			revenue, _ := strconv.ParseFloat(record[5], 64)

			dailyStat := domain.DailyStat{
				Date:        record[1],
				Clicks:      clicks,
				Impressions: impressions,
				Conversions: conversions,
				Revenue:     revenue,
			}
			dailyStats = append(dailyStats, dailyStat)
		}
	}

	return &domain.CampaignDetail{
		Campaign:    *campaignInfo,
		Performance: *campaignMetrics,
		DailyStats:  dailyStats,
	}, nil
}

// loadCSV loads a CSV file and returns the records
func (s *mockDataService) loadCSV(filename string) ([][]string, error) {
	filePath := filepath.Join(s.dataPath, filename)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file %s: %w", filePath, err)
	}

	return records, nil
}

// LoadPerformanceSummary loads performance summary data from CSV
func (s *mockDataService) LoadPerformanceSummary(ctx context.Context, orgID int64, filters domain.ReportingFilters) (*domain.PerformanceSummary, error) {
	records, err := s.loadCSV("performance_summary.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load performance summary: %w", err)
	}

	for _, record := range records[1:] { // Skip header
		if len(record) < 8 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if recordOrgID == orgID {
			totalClicks, _ := strconv.ParseInt(record[1], 10, 64)
			totalConversions, _ := strconv.ParseInt(record[2], 10, 64)
			totalRevenue, _ := strconv.ParseFloat(record[3], 64)
			conversionRate, _ := strconv.ParseFloat(record[4], 64)
			averageRevenue, _ := strconv.ParseFloat(record[5], 64)
			clickThroughRate, _ := strconv.ParseFloat(record[6], 64)
			totalImpressions, _ := strconv.ParseInt(record[7], 10, 64)

			return &domain.PerformanceSummary{
				TotalClicks:      totalClicks,
				TotalConversions: totalConversions,
				TotalRevenue:     totalRevenue,
				ConversionRate:   conversionRate,
				AverageRevenue:   averageRevenue,
				ClickThroughRate: clickThroughRate,
				TotalImpressions: totalImpressions,
			}, nil
		}
	}

	// Return default data if not found
	return &domain.PerformanceSummary{
		TotalClicks:      1000,
		TotalConversions: 80,
		TotalRevenue:     2500.00,
		ConversionRate:   8.0,
		AverageRevenue:   31.25,
		ClickThroughRate: 2.5,
		TotalImpressions: 40000,
	}, nil
}

// LoadPerformanceTimeSeries loads performance time series data from CSV
func (s *mockDataService) LoadPerformanceTimeSeries(ctx context.Context, orgID int64, filters domain.ReportingFilters) ([]domain.PerformanceTimeSeriesPoint, error) {
	records, err := s.loadCSV("performance_timeseries.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load performance timeseries: %w", err)
	}

	var timeSeries []domain.PerformanceTimeSeriesPoint

	for _, record := range records[1:] { // Skip header
		if len(record) < 8 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if recordOrgID == orgID {
			// Apply date filtering if provided
			if filters.StartDate != "" && record[1] < filters.StartDate {
				continue
			}
			if filters.EndDate != "" && record[1] > filters.EndDate {
				continue
			}

			clicks, _ := strconv.ParseInt(record[2], 10, 64)
			impressions, _ := strconv.ParseInt(record[3], 10, 64)
			conversions, _ := strconv.ParseInt(record[4], 10, 64)
			revenue, _ := strconv.ParseFloat(record[5], 64)
			conversionRate, _ := strconv.ParseFloat(record[6], 64)
			clickThroughRate, _ := strconv.ParseFloat(record[7], 64)

			point := domain.PerformanceTimeSeriesPoint{
				Date:             record[1],
				Clicks:           clicks,
				Impressions:      impressions,
				Conversions:      conversions,
				Revenue:          revenue,
				ConversionRate:   conversionRate,
				ClickThroughRate: clickThroughRate,
			}
			timeSeries = append(timeSeries, point)
		}
	}

	return timeSeries, nil
}

// LoadDailyPerformanceReport loads daily performance report data from CSV
func (s *mockDataService) LoadDailyPerformanceReport(ctx context.Context, orgID int64, filters domain.ReportingFilters, pagination domain.PaginationParams) ([]domain.DailyPerformanceReport, *domain.PaginationResult, error) {
	records, err := s.loadCSV("daily_performance_report.csv")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load daily performance report: %w", err)
	}

	var reports []domain.DailyPerformanceReport

	for _, record := range records[1:] { // Skip header
		if len(record) < 11 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if recordOrgID == orgID {
			// Apply date filtering if provided
			if filters.StartDate != "" && record[1] < filters.StartDate {
				continue
			}
			if filters.EndDate != "" && record[1] > filters.EndDate {
				continue
			}

			// Apply campaign filtering if provided
			if len(filters.CampaignIDs) > 0 {
				campaignMatch := false
				for _, campaignID := range filters.CampaignIDs {
					if record[2] == campaignID {
						campaignMatch = true
						break
					}
				}
				if !campaignMatch {
					continue
				}
			}

			clicks, _ := strconv.ParseInt(record[4], 10, 64)
			impressions, _ := strconv.ParseInt(record[5], 10, 64)
			conversions, _ := strconv.ParseInt(record[6], 10, 64)
			revenue, _ := strconv.ParseFloat(record[7], 64)
			conversionRate, _ := strconv.ParseFloat(record[8], 64)
			clickThroughRate, _ := strconv.ParseFloat(record[9], 64)
			payouts, _ := strconv.ParseFloat(record[10], 64)

			report := domain.DailyPerformanceReport{
				Date:             record[1],
				CampaignID:       record[2],
				CampaignName:     record[3],
				Clicks:           clicks,
				Impressions:      impressions,
				Conversions:      conversions,
				Revenue:          revenue,
				ConversionRate:   conversionRate,
				ClickThroughRate: clickThroughRate,
				Payouts:          payouts,
			}
			reports = append(reports, report)
		}
	}

	// Apply pagination
	totalItems := len(reports)
	startIndex := (pagination.Page - 1) * pagination.Limit
	endIndex := startIndex + pagination.Limit

	if startIndex >= totalItems {
		reports = []domain.DailyPerformanceReport{}
	} else if endIndex > totalItems {
		reports = reports[startIndex:]
	} else {
		reports = reports[startIndex:endIndex]
	}

	totalPages := (totalItems + pagination.Limit - 1) / pagination.Limit
	paginationResult := &domain.PaginationResult{
		CurrentPage:     pagination.Page,
		TotalPages:      totalPages,
		TotalItems:      totalItems,
		ItemsPerPage:    pagination.Limit,
		HasNextPage:     pagination.Page < totalPages,
		HasPreviousPage: pagination.Page > 1,
	}

	return reports, paginationResult, nil
}

// LoadConversionsReport loads conversions report data from CSV
func (s *mockDataService) LoadConversionsReport(ctx context.Context, orgID int64, filters domain.ReportingFilters, pagination domain.PaginationParams) ([]domain.ConversionReport, *domain.PaginationResult, error) {
	records, err := s.loadCSV("conversions_report.csv")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load conversions report: %w", err)
	}

	var conversions []domain.ConversionReport

	for _, record := range records[1:] { // Skip header
		if len(record) < 18 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if recordOrgID == orgID {
			// Apply campaign filtering if provided
			if len(filters.CampaignIDs) > 0 {
				campaignMatch := false
				for _, campaignID := range filters.CampaignIDs {
					if record[4] == campaignID {
						campaignMatch = true
						break
					}
				}
				if !campaignMatch {
					continue
				}
			}

			// Apply affiliate filtering if provided
			if filters.AffiliateID != nil && *filters.AffiliateID != "" && record[11] != *filters.AffiliateID {
				continue
			}

			// Apply status filtering if provided
			if filters.Status != nil && *filters.Status != "" && record[8] != *filters.Status {
				continue
			}

			timestamp, _ := time.Parse(time.RFC3339, record[2])
			payout, _ := strconv.ParseFloat(record[9], 64)
			var conversionValue *float64
			if record[14] != "" {
				cv, _ := strconv.ParseFloat(record[14], 64)
				conversionValue = &cv
			}

			var clickID, sub1, sub2, sub3 *string
			if record[13] != "" {
				clickID = &record[13]
			}
			if record[15] != "" {
				sub1 = &record[15]
			}
			if record[16] != "" {
				sub2 = &record[16]
			}
			if record[17] != "" {
				sub3 = &record[17]
			}

			conversion := domain.ConversionReport{
				ID:              record[1],
				Timestamp:       timestamp,
				TransactionID:   record[3],
				CampaignID:      record[4],
				CampaignName:    record[5],
				OfferID:         record[6],
				OfferName:       record[7],
				Status:          record[8],
				Payout:          payout,
				Currency:        record[10],
				AffiliateID:     record[11],
				AffiliateName:   record[12],
				ClickID:         clickID,
				ConversionValue: conversionValue,
				Sub1:            sub1,
				Sub2:            sub2,
				Sub3:            sub3,
			}
			conversions = append(conversions, conversion)
		}
	}

	// Sort by timestamp DESC (most recent first)
	sort.Slice(conversions, func(i, j int) bool {
		return conversions[i].Timestamp.After(conversions[j].Timestamp)
	})

	// Apply pagination
	totalItems := len(conversions)
	startIndex := (pagination.Page - 1) * pagination.Limit
	endIndex := startIndex + pagination.Limit

	if startIndex >= totalItems {
		conversions = []domain.ConversionReport{}
	} else if endIndex > totalItems {
		conversions = conversions[startIndex:]
	} else {
		conversions = conversions[startIndex:endIndex]
	}

	totalPages := (totalItems + pagination.Limit - 1) / pagination.Limit
	paginationResult := &domain.PaginationResult{
		CurrentPage:     pagination.Page,
		TotalPages:      totalPages,
		TotalItems:      totalItems,
		ItemsPerPage:    pagination.Limit,
		HasNextPage:     pagination.Page < totalPages,
		HasPreviousPage: pagination.Page > 1,
	}

	return conversions, paginationResult, nil
}

// LoadClicksReport loads clicks report data from CSV
func (s *mockDataService) LoadClicksReport(ctx context.Context, orgID int64, filters domain.ReportingFilters, pagination domain.PaginationParams) ([]domain.ClickReport, *domain.PaginationResult, error) {
	records, err := s.loadCSV("clicks_report.csv")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load clicks report: %w", err)
	}

	var clicks []domain.ClickReport

	for _, record := range records[1:] { // Skip header
		if len(record) < 21 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if recordOrgID == orgID {
			// Apply campaign filtering if provided
			if len(filters.CampaignIDs) > 0 {
				campaignMatch := false
				for _, campaignID := range filters.CampaignIDs {
					if record[3] == campaignID {
						campaignMatch = true
						break
					}
				}
				if !campaignMatch {
					continue
				}
			}

			// Apply affiliate filtering if provided
			if filters.AffiliateID != nil && *filters.AffiliateID != "" && record[7] != *filters.AffiliateID {
				continue
			}

			timestamp, _ := time.Parse(time.RFC3339, record[2])
			converted, _ := strconv.ParseBool(record[19])

			var region, city, referrerURL, sub1, sub2, sub3, conversionID *string
			if record[12] != "" {
				region = &record[12]
			}
			if record[13] != "" {
				city = &record[13]
			}
			if record[14] != "" {
				referrerURL = &record[14]
			}
			if record[16] != "" {
				sub1 = &record[16]
			}
			if record[17] != "" {
				sub2 = &record[17]
			}
			if record[18] != "" {
				sub3 = &record[18]
			}
			if record[20] != "" {
				conversionID = &record[20]
			}

			click := domain.ClickReport{
				ID:             record[1],
				Timestamp:      timestamp,
				CampaignID:     record[3],
				CampaignName:   record[4],
				OfferID:        record[5],
				OfferName:      record[6],
				AffiliateID:    record[7],
				AffiliateName:  record[8],
				IPAddress:      record[9],
				UserAgent:      record[10],
				Country:        record[11],
				Region:         region,
				City:           city,
				ReferrerURL:    referrerURL,
				LandingPageURL: record[15],
				Sub1:           sub1,
				Sub2:           sub2,
				Sub3:           sub3,
				Converted:      converted,
				ConversionID:   conversionID,
			}
			clicks = append(clicks, click)
		}
	}

	// Sort by timestamp DESC (most recent first)
	sort.Slice(clicks, func(i, j int) bool {
		return clicks[i].Timestamp.After(clicks[j].Timestamp)
	})

	// Apply pagination
	totalItems := len(clicks)
	startIndex := (pagination.Page - 1) * pagination.Limit
	endIndex := startIndex + pagination.Limit

	if startIndex >= totalItems {
		clicks = []domain.ClickReport{}
	} else if endIndex > totalItems {
		clicks = clicks[startIndex:]
	} else {
		clicks = clicks[startIndex:endIndex]
	}

	totalPages := (totalItems + pagination.Limit - 1) / pagination.Limit
	paginationResult := &domain.PaginationResult{
		CurrentPage:     pagination.Page,
		TotalPages:      totalPages,
		TotalItems:      totalItems,
		ItemsPerPage:    pagination.Limit,
		HasNextPage:     pagination.Page < totalPages,
		HasPreviousPage: pagination.Page > 1,
	}

	return clicks, paginationResult, nil
}

// LoadCampaignsList loads campaigns list data from CSV
func (s *mockDataService) LoadCampaignsList(ctx context.Context, orgID int64, affiliateID *string, status string, search *string) ([]domain.CampaignListItem, error) {
	records, err := s.loadCSV("campaigns_list.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load campaigns list: %w", err)
	}

	var campaigns []domain.CampaignListItem

	for _, record := range records[1:] { // Skip header
		if len(record) < 4 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if recordOrgID == orgID {
			// Apply status filtering if provided
			if status != "" && record[3] != status {
				continue
			}

			// Apply search filtering if provided
			if search != nil && *search != "" {
				searchTerm := strings.ToLower(*search)
				if !strings.Contains(strings.ToLower(record[2]), searchTerm) {
					continue
				}
			}

			campaign := domain.CampaignListItem{
				ID:     record[1],
				Name:   record[2],
				Status: record[3],
			}
			campaigns = append(campaigns, campaign)
		}
	}

	return campaigns, nil
}

// LoadAdvertiserMetrics loads advertiser metrics with historical data
func (s *mockDataService) LoadAdvertiserMetrics(ctx context.Context, orgID int64) (*domain.AdvertiserMetrics, error) {
	records, err := s.loadCSV("historical_metrics.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load historical metrics: %w", err)
	}

	for _, record := range records[1:] { // Skip header
		if len(record) < 7 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if recordOrgID == orgID {
			// Parse metrics for this organization
			metrics := &domain.AdvertiserMetrics{}
			
			// Find all metrics for this organization
			for _, r := range records[1:] {
				if len(r) < 7 {
					continue
				}
				
				rOrgID, _ := strconv.ParseInt(r[0], 10, 64)
				if rOrgID != orgID {
					continue
				}
				
				metricType := r[1]
				today, _ := strconv.ParseFloat(r[2], 64)
				yesterday, _ := strconv.ParseFloat(r[3], 64)
				currentMonth, _ := strconv.ParseFloat(r[4], 64)
				lastMonth, _ := strconv.ParseFloat(r[5], 64)
				changePercentage, _ := strconv.ParseFloat(r[6], 64)
				
				metric := domain.MetricWithHistory{
					Today:           today,
					Yesterday:       yesterday,
					CurrentMonth:    currentMonth,
					LastMonth:       lastMonth,
					ChangePercentage: changePercentage,
				}
				
				switch metricType {
				case "total_clicks":
					metrics.TotalClicks = metric
				case "conversions":
					metrics.Conversions = metric
				case "revenue":
					metrics.Revenue = metric
				case "conversion_rate":
					metrics.ConversionRate = metric
				case "events":
					metrics.Events = metric
				case "event_rate":
					metrics.EventRate = metric
				}
			}
			
			return metrics, nil
		}
	}

	// Return default data if not found
	return &domain.AdvertiserMetrics{
		TotalClicks: domain.MetricWithHistory{
			Today: 1250, Yesterday: 1180, CurrentMonth: 35000, LastMonth: 32000, ChangePercentage: 9.4,
		},
		Conversions: domain.MetricWithHistory{
			Today: 89, Yesterday: 85, CurrentMonth: 2450, LastMonth: 2200, ChangePercentage: 11.4,
		},
		Revenue: domain.MetricWithHistory{
			Today: 4450.50, Yesterday: 4200.00, CurrentMonth: 122500.00, LastMonth: 110000.00, ChangePercentage: 11.4,
		},
		ConversionRate: domain.MetricWithHistory{
			Today: 0.071, Yesterday: 0.072, CurrentMonth: 0.070, LastMonth: 0.069, ChangePercentage: 1.4,
		},
		Events: domain.MetricWithHistory{
			Today: 156, Yesterday: 148, CurrentMonth: 4200, LastMonth: 3800, ChangePercentage: 10.5,
		},
		EventRate: domain.MetricWithHistory{
			Today: 0.125, Yesterday: 0.125, CurrentMonth: 0.120, LastMonth: 0.119, ChangePercentage: 0.8,
		},
	}, nil
}

// LoadOffersPaginated loads paginated offers data
func (s *mockDataService) LoadOffersPaginated(ctx context.Context, orgID int64, page, perPage int) (*domain.OffersPaginated, error) {
	records, err := s.loadCSV("offers.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load offers: %w", err)
	}

	var allOffers []domain.Offer
	for _, record := range records[1:] { // Skip header
		if len(record) < 11 {
			continue
		}

		advertiserOrgID, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			continue
		}

		if advertiserOrgID == orgID {
			id, _ := strconv.ParseInt(record[0], 10, 64)
			payout, _ := strconv.ParseFloat(record[4], 64)
			
			// Parse countries JSON-like string
			countriesStr := strings.Trim(record[7], "[]\"")
			countries := strings.Split(countriesStr, "\",\"")
			for i, country := range countries {
				countries[i] = strings.Trim(country, "\"")
			}
			
			var description *string
			if record[3] != "" {
				description = &record[3]
			}

			offer := domain.Offer{
				ID:             id,
				Name:           record[2],
				Description:    description,
				Payout:         payout,
				Currency:       record[5],
				Status:         record[6],
				Category:       record[8],
				Countries:      countries,
				ConversionFlow: record[9],
				CreatedAt:      record[10],
			}
			allOffers = append(allOffers, offer)
		}
	}

	// Calculate pagination
	totalCount := len(allOffers)
	startIndex := (page - 1) * perPage
	endIndex := startIndex + perPage
	
	if startIndex >= totalCount {
		return &domain.OffersPaginated{
			Items:      []domain.Offer{},
			TotalCount: totalCount,
			Page:       page,
			PerPage:    perPage,
			HasNext:    false,
		}, nil
	}
	
	if endIndex > totalCount {
		endIndex = totalCount
	}
	
	items := allOffers[startIndex:endIndex]
	hasNext := endIndex < totalCount

	return &domain.OffersPaginated{
		Items:      items,
		TotalCount: totalCount,
		Page:       page,
		PerPage:    perPage,
		HasNext:    hasNext,
	}, nil
}

// LoadAgencyPerformanceOverview loads agency performance overview
func (s *mockDataService) LoadAgencyPerformanceOverview(ctx context.Context, orgID int64) (*domain.AgencyPerformanceOverview, error) {
	records, err := s.loadCSV("historical_metrics.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load historical metrics: %w", err)
	}

	overview := &domain.AgencyPerformanceOverview{}
	
	for _, record := range records[1:] { // Skip header
		if len(record) < 7 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if recordOrgID == orgID {
			metricType := record[1]
			today, _ := strconv.ParseFloat(record[2], 64)
			yesterday, _ := strconv.ParseFloat(record[3], 64)
			currentMonth, _ := strconv.ParseFloat(record[4], 64)
			lastMonth, _ := strconv.ParseFloat(record[5], 64)
			changePercentage, _ := strconv.ParseFloat(record[6], 64)
			
			metric := domain.MetricWithHistory{
				Today:           today,
				Yesterday:       yesterday,
				CurrentMonth:    currentMonth,
				LastMonth:       lastMonth,
				ChangePercentage: changePercentage,
			}
			
			switch metricType {
			case "total_clicks":
				overview.TotalClicks = metric
				overview.ClicksPerDay = domain.MetricWithHistory{
					Today: today, Yesterday: yesterday, 
					CurrentMonth: currentMonth/30, LastMonth: lastMonth/30, 
					ChangePercentage: changePercentage,
				}
			case "conversions":
				overview.TotalConversions = metric
				overview.ConversionsPerDay = domain.MetricWithHistory{
					Today: today, Yesterday: yesterday,
					CurrentMonth: currentMonth/30, LastMonth: lastMonth/30,
					ChangePercentage: changePercentage,
				}
			case "revenue":
				overview.TotalEarnings = metric
			}
		}
	}

	// Return default if not found
	if overview.TotalClicks.Today == 0 {
		return &domain.AgencyPerformanceOverview{
			TotalConversions: domain.MetricWithHistory{
				Today: 450, Yesterday: 420, CurrentMonth: 12500, LastMonth: 11200, ChangePercentage: 11.6,
			},
			TotalClicks: domain.MetricWithHistory{
				Today: 6200, Yesterday: 5800, CurrentMonth: 175000, LastMonth: 162000, ChangePercentage: 8.0,
			},
			ConversionsPerDay: domain.MetricWithHistory{
				Today: 450, Yesterday: 420, CurrentMonth: 417, LastMonth: 373, ChangePercentage: 11.6,
			},
			ClicksPerDay: domain.MetricWithHistory{
				Today: 6200, Yesterday: 5800, CurrentMonth: 5833, LastMonth: 5400, ChangePercentage: 8.0,
			},
			TotalEarnings: domain.MetricWithHistory{
				Today: 22500.00, Yesterday: 21000.00, CurrentMonth: 625000.00, LastMonth: 560000.00, ChangePercentage: 11.6,
			},
		}, nil
	}

	return overview, nil
}

// LoadAdvertiserOrganizations loads advertiser organizations for agencies
func (s *mockDataService) LoadAdvertiserOrganizations(ctx context.Context, orgID int64) ([]domain.AdvertiserOrganization, error) {
	records, err := s.loadCSV("advertiser_organizations.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load advertiser organizations: %w", err)
	}

	var organizations []domain.AdvertiserOrganization
	for _, record := range records[1:] { // Skip header
		if len(record) < 7 {
			continue
		}

		agencyOrgID, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			continue
		}

		if agencyOrgID == orgID {
			id, _ := strconv.ParseInt(record[0], 10, 64)
			campaignsCount, _ := strconv.Atoi(record[4])
			revenue, _ := strconv.ParseFloat(record[5], 64)
			conversionRate, _ := strconv.ParseFloat(record[6], 64)

			org := domain.AdvertiserOrganization{
				ID:             id,
				Name:           record[2],
				Status:         record[3],
				CampaignsCount: campaignsCount,
				Revenue:        revenue,
				ConversionRate: conversionRate,
			}
			organizations = append(organizations, org)
		}
	}

	return organizations, nil
}

// LoadCampaignsOverview loads campaigns overview for agencies
func (s *mockDataService) LoadCampaignsOverview(ctx context.Context, orgID int64) (*domain.CampaignsOverview, error) {
	records, err := s.loadCSV("enhanced_campaign_performance.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load enhanced campaign performance: %w", err)
	}

	var campaigns []domain.AgencyCampaignSummary
	activeCount := 0
	
	for _, record := range records[1:] { // Skip header
		if len(record) < 15 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			continue
		}

		if recordOrgID == orgID {
			id, _ := strconv.ParseInt(record[0], 10, 64)
			advertiserOrgID, _ := strconv.ParseInt(record[2], 10, 64)
			clicks, _ := strconv.ParseInt(record[5], 10, 64)
			conversions, _ := strconv.ParseInt(record[6], 10, 64)
			totalCost, _ := strconv.ParseFloat(record[14], 64)
			conversionRate, _ := strconv.ParseFloat(record[8], 64)

			campaign := domain.AgencyCampaignSummary{
				ID:               id,
				Name:             record[4],
				AdvertiserOrgID:  advertiserOrgID,
				AdvertiserName:   record[3],
				Status:           record[11],
				Clicks:           clicks,
				Conversions:      conversions,
				TotalCost:        totalCost,
				ConversionRate:   conversionRate,
			}
			
			if campaign.Status == "active" {
				activeCount++
			}
			
			campaigns = append(campaigns, campaign)
		}
	}

	return &domain.CampaignsOverview{
		Campaigns:   campaigns,
		TotalCount:  len(campaigns),
		ActiveCount: activeCount,
	}, nil
}

// LoadAgencyPerformanceChart loads agency performance chart data
func (s *mockDataService) LoadAgencyPerformanceChart(ctx context.Context, orgID int64, period string) (*domain.AgencyPerformanceChart, error) {
	records, err := s.loadCSV("agency_performance_chart.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load agency performance chart: %w", err)
	}

	var dataPoints []domain.AgencyPerformanceDataPoint
	currentDate := ""
	var currentPoint *domain.AgencyPerformanceDataPoint
	
	for _, record := range records[1:] { // Skip header
		if len(record) < 8 {
			continue
		}

		agencyOrgID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if agencyOrgID == orgID {
			date := record[1]
			
			// If this is a new date, save the previous point and start a new one
			if date != currentDate {
				if currentPoint != nil {
					dataPoints = append(dataPoints, *currentPoint)
				}
				
				conversions, _ := strconv.ParseInt(record[2], 10, 64)
				clicks, _ := strconv.ParseInt(record[3], 10, 64)
				revenue, _ := strconv.ParseFloat(record[4], 64)
				
				currentPoint = &domain.AgencyPerformanceDataPoint{
					Date:                date,
					Conversions:         conversions,
					Clicks:              clicks,
					Revenue:             revenue,
					AdvertiserBreakdown: []domain.AdvertiserBreakdownPoint{},
				}
				currentDate = date
			}
			
			// Add advertiser breakdown
			if len(record) >= 8 {
				advertiserID, _ := strconv.ParseInt(record[5], 10, 64)
				advertiserConversions, _ := strconv.ParseInt(record[7], 10, 64)
				advertiserRevenue, _ := strconv.ParseFloat(record[8], 64)
				
				breakdown := domain.AdvertiserBreakdownPoint{
					AdvertiserID:   advertiserID,
					AdvertiserName: record[6],
					Conversions:    advertiserConversions,
					Revenue:        advertiserRevenue,
				}
				currentPoint.AdvertiserBreakdown = append(currentPoint.AdvertiserBreakdown, breakdown)
			}
		}
	}
	
	// Add the last point
	if currentPoint != nil {
		dataPoints = append(dataPoints, *currentPoint)
	}

	return &domain.AgencyPerformanceChart{
		Data:   dataPoints,
		Period: period,
	}, nil
}

// LoadPlatformOverview loads platform overview with historical data
func (s *mockDataService) LoadPlatformOverview(ctx context.Context) (*domain.PlatformOverview, error) {
	records, err := s.loadCSV("historical_metrics.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load historical metrics: %w", err)
	}

	overview := &domain.PlatformOverview{}
	
	for _, record := range records[1:] { // Skip header
		if len(record) < 7 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		// Platform metrics have orgID 100
		if recordOrgID == 100 {
			metricType := record[1]
			today, _ := strconv.ParseFloat(record[2], 64)
			yesterday, _ := strconv.ParseFloat(record[3], 64)
			currentMonth, _ := strconv.ParseFloat(record[4], 64)
			lastMonth, _ := strconv.ParseFloat(record[5], 64)
			changePercentage, _ := strconv.ParseFloat(record[6], 64)
			
			metric := domain.MetricWithHistory{
				Today:           today,
				Yesterday:       yesterday,
				CurrentMonth:    currentMonth,
				LastMonth:       lastMonth,
				ChangePercentage: changePercentage,
			}
			
			switch metricType {
			case "total_organizations":
				overview.TotalOrganizations = metric
			case "total_users":
				overview.TotalUsers = metric
			case "total_clicks":
				overview.TotalRevenue = domain.MetricWithHistory{
					Today: today * 0.1, Yesterday: yesterday * 0.1,
					CurrentMonth: currentMonth * 0.1, LastMonth: lastMonth * 0.1,
					ChangePercentage: changePercentage,
				}
			case "monthly_growth":
				overview.MonthlyGrowth = metric
			case "new_registrations":
				overview.NewRegistrations = metric
			}
		}
	}

	// Return default if not found
	if overview.TotalUsers.Today == 0 {
		return &domain.PlatformOverview{
			TotalOrganizations: domain.MetricWithHistory{
				Today: 1250, Yesterday: 1248, CurrentMonth: 1250, LastMonth: 1180, ChangePercentage: 5.9,
			},
			TotalUsers: domain.MetricWithHistory{
				Today: 5600, Yesterday: 5580, CurrentMonth: 5600, LastMonth: 5200, ChangePercentage: 7.7,
			},
			TotalRevenue: domain.MetricWithHistory{
				Today: 22500.00, Yesterday: 21000.00, CurrentMonth: 625000.00, LastMonth: 560000.00, ChangePercentage: 11.6,
			},
			MonthlyGrowth: domain.MetricWithHistory{
				Today: 0.094, Yesterday: 0.092, CurrentMonth: 0.094, LastMonth: 0.087, ChangePercentage: 8.0,
			},
			NewRegistrations: domain.MetricWithHistory{
				Today: 25, Yesterday: 22, CurrentMonth: 680, LastMonth: 620, ChangePercentage: 9.7,
			},
		}, nil
	}

	return overview, nil
}

// LoadUserActivityMetrics loads user activity metrics
func (s *mockDataService) LoadUserActivityMetrics(ctx context.Context) (*domain.UserActivityMetrics, error) {
	records, err := s.loadCSV("enhanced_platform_metrics.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load platform metrics: %w", err)
	}

	if len(records) > 1 && len(records[1]) >= 6 {
		record := records[1]
		dailyActiveUsers, _ := strconv.Atoi(record[1])
		weeklyActiveUsers, _ := strconv.Atoi(record[2])
		monthlyActiveUsers, _ := strconv.Atoi(record[3])
		activeAdvertisers, _ := strconv.Atoi(record[4])
		activeAffiliates, _ := strconv.Atoi(record[5])

		return &domain.UserActivityMetrics{
			DailyActiveUsers:   dailyActiveUsers,
			WeeklyActiveUsers:  weeklyActiveUsers,
			MonthlyActiveUsers: monthlyActiveUsers,
			ActiveAdvertisers:  activeAdvertisers,
			ActiveAffiliates:   activeAffiliates,
		}, nil
	}

	// Return default data
	return &domain.UserActivityMetrics{
		DailyActiveUsers:   2800,
		WeeklyActiveUsers:  4200,
		MonthlyActiveUsers: 5600,
		ActiveAdvertisers:  450,
		ActiveAffiliates:   800,
	}, nil
}

// LoadSystemHealthMetrics loads system health metrics
func (s *mockDataService) LoadSystemHealthMetrics(ctx context.Context) (*domain.SystemHealthMetrics, error) {
	records, err := s.loadCSV("enhanced_platform_metrics.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load platform metrics: %w", err)
	}

	if len(records) > 1 && len(records[1]) >= 12 {
		record := records[1]
		totalCampaigns, _ := strconv.Atoi(record[6])
		requestsPerMinute, _ := strconv.Atoi(record[7])
		successRate, _ := strconv.ParseFloat(record[8], 64)
		rateLimitHits, _ := strconv.Atoi(record[9])
		averageQueryTime, _ := strconv.ParseFloat(record[10], 64)
		connectionPoolUsage, _ := strconv.ParseFloat(record[11], 64)

		return &domain.SystemHealthMetrics{
			TotalCampaigns:       totalCampaigns,
			RequestsPerMinute:    requestsPerMinute,
			SuccessRate:          successRate,
			RateLimitHits:        rateLimitHits,
			AverageQueryTime:     averageQueryTime,
			ConnectionPoolUsage:  connectionPoolUsage,
		}, nil
	}

	// Return default data
	return &domain.SystemHealthMetrics{
		TotalCampaigns:       2500,
		RequestsPerMinute:    2500,
		SuccessRate:          99.95,
		RateLimitHits:        12,
		AverageQueryTime:     25,
		ConnectionPoolUsage:  65,
	}, nil
}

// LoadRevenueBySource loads revenue by source data
func (s *mockDataService) LoadRevenueBySource(ctx context.Context, orgID int64) ([]domain.RevenueBySource, error) {
	records, err := s.loadCSV("revenue_by_source.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load revenue by source: %w", err)
	}

	var sources []domain.RevenueBySource
	for _, record := range records[1:] { // Skip header
		if len(record) < 5 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if recordOrgID == orgID {
			revenue, _ := strconv.ParseFloat(record[2], 64)
			percentage, _ := strconv.ParseFloat(record[3], 64)
			growth, _ := strconv.ParseFloat(record[4], 64)

			source := domain.RevenueBySource{
				Source:     record[1],
				Revenue:    revenue,
				Percentage: percentage,
				Growth:     growth,
			}
			sources = append(sources, source)
		}
	}

	return sources, nil
}

// LoadGeographicDistribution loads geographic distribution data
func (s *mockDataService) LoadGeographicDistribution(ctx context.Context, orgID int64) ([]domain.GeographicData, error) {
	records, err := s.loadCSV("geographic_distribution.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load geographic distribution: %w", err)
	}

	var geoData []domain.GeographicData
	for _, record := range records[1:] { // Skip header
		if len(record) < 5 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if recordOrgID == orgID {
			revenue, _ := strconv.ParseFloat(record[2], 64)
			users, _ := strconv.Atoi(record[3])
			conversionRate, _ := strconv.ParseFloat(record[4], 64)

			data := domain.GeographicData{
				Country:        record[1],
				Revenue:        revenue,
				Users:          users,
				ConversionRate: conversionRate,
			}
			geoData = append(geoData, data)
		}
	}

	return geoData, nil
}

// LoadEnhancedRevenueChart loads enhanced revenue chart with events data
func (s *mockDataService) LoadEnhancedRevenueChart(ctx context.Context, orgID int64, period string) (*domain.RevenueChart, error) {
	records, err := s.loadCSV("enhanced_revenue_chart.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load enhanced revenue chart: %w", err)
	}

	var dataPoints []domain.RevenueDataPoint
	for _, record := range records[1:] { // Skip header
		if len(record) < 6 {
			continue
		}

		recordOrgID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if recordOrgID == orgID {
			revenue, _ := strconv.ParseFloat(record[2], 64)
			clicks, _ := strconv.ParseInt(record[3], 10, 64)
			conversions, _ := strconv.ParseInt(record[4], 10, 64)
			events, _ := strconv.ParseInt(record[5], 10, 64)

			dataPoint := domain.RevenueDataPoint{
				Date:        record[1],
				Revenue:     revenue,
				Clicks:      clicks,
				Conversions: conversions,
				Events:      events,
			}
			dataPoints = append(dataPoints, dataPoint)
		}
	}

	return &domain.RevenueChart{
		Data:   dataPoints,
		Period: period,
	}, nil
}
