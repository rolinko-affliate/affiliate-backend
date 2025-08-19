package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/logger"
)

// MockDataService provides methods to load mock data from CSV files
type MockDataService interface {
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
