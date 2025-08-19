package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/logger"
)

// MockDataService provides methods to load mock data from CSV files
type MockDataService interface {
	// Dashboard data methods
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
