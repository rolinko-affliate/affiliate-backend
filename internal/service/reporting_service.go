package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/everflow/reporting"
	"github.com/affiliate-backend/internal/platform/logger"
	"github.com/affiliate-backend/internal/repository"
)

// ReportingService handles reporting business logic
type ReportingService interface {
	GetPerformanceSummary(ctx context.Context, filters domain.ReportingFilters, userProfile *domain.Profile) (*domain.PerformanceSummary, error)
	GetPerformanceTimeSeries(ctx context.Context, filters domain.ReportingFilters, userProfile *domain.Profile) ([]domain.PerformanceTimeSeriesPoint, error)
	GetDailyPerformanceReport(ctx context.Context, filters domain.ReportingFilters, pagination domain.PaginationParams, userProfile *domain.Profile) ([]domain.DailyPerformanceReport, *domain.PaginationResult, error)
	GetConversionsReport(ctx context.Context, filters domain.ReportingFilters, pagination domain.PaginationParams, userProfile *domain.Profile) ([]domain.ConversionReport, *domain.PaginationResult, error)
	GetClicksReport(ctx context.Context, filters domain.ReportingFilters, pagination domain.PaginationParams, userProfile *domain.Profile) ([]domain.ClickReport, *domain.PaginationResult, error)
	GetCampaignsList(ctx context.Context, affiliateID *string, status string, search *string, userProfile *domain.Profile) ([]domain.CampaignListItem, error)
}

type reportingService struct {
	everflowClient *reporting.Client
	repo           repository.ReportingRepository
	campaignRepo   repository.CampaignRepository
}

// NewReportingService creates a new reporting service
func NewReportingService(
	everflowClient *reporting.Client,
	repo repository.ReportingRepository,
	campaignRepo repository.CampaignRepository,
) ReportingService {
	return &reportingService{
		everflowClient: everflowClient,
		repo:           repo,
		campaignRepo:   campaignRepo,
	}
}

// GetPerformanceSummary retrieves aggregated performance metrics
func (s *reportingService) GetPerformanceSummary(ctx context.Context, filters domain.ReportingFilters, userProfile *domain.Profile) (*domain.PerformanceSummary, error) {
	// Generate cache key
	cacheKey := s.generateCacheKey("summary", filters, userProfile)

	// Try to get from cache first
	if cached, err := s.repo.GetCachedSummary(ctx, cacheKey); err == nil && cached != nil {
		logger.Info("Returning cached performance summary", "cache_key", cacheKey)
		return cached, nil
	}

	// Build Everflow entity report request
	req := s.buildEntityReportRequest(filters, []string{"offer"}) // Group by offer for summary

	// Apply user-specific filters
	if err := s.applyUserFilters(&req, userProfile); err != nil {
		return nil, fmt.Errorf("failed to apply user filters: %w", err)
	}

	// Get data from Everflow
	resp, err := s.everflowClient.GetEntityReport(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get entity report from Everflow: %w", err)
	}

	// Transform Everflow response to domain model
	summary := s.transformToPerformanceSummary(resp)

	// Cache the result
	if err := s.repo.SetCachedSummary(ctx, cacheKey, summary, 15*time.Minute); err != nil {
		logger.Warn("Failed to cache performance summary", "error", err)
	}

	return summary, nil
}

// GetPerformanceTimeSeries retrieves daily performance data for charts
func (s *reportingService) GetPerformanceTimeSeries(ctx context.Context, filters domain.ReportingFilters, userProfile *domain.Profile) ([]domain.PerformanceTimeSeriesPoint, error) {
	// Generate cache key
	cacheKey := s.generateCacheKey("timeseries", filters, userProfile)

	// Try to get from cache first
	if cached, err := s.repo.GetCachedTimeSeries(ctx, cacheKey); err == nil && cached != nil {
		logger.Info("Returning cached time series data", "cache_key", cacheKey)
		return cached, nil
	}

	// Determine grouping column based on granularity
	granularity := "daily"
	if filters.Granularity != nil {
		granularity = *filters.Granularity
	}

	var columns []string
	switch granularity {
	case "hourly":
		columns = []string{"hour"}
	case "weekly":
		columns = []string{"week"}
	case "monthly":
		columns = []string{"month"}
	default:
		columns = []string{"date"}
	}

	// Build Everflow entity report request
	req := s.buildEntityReportRequest(filters, columns)

	// Apply user-specific filters
	if err := s.applyUserFilters(&req, userProfile); err != nil {
		return nil, fmt.Errorf("failed to apply user filters: %w", err)
	}

	// Get data from Everflow
	resp, err := s.everflowClient.GetEntityReport(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get entity report from Everflow: %w", err)
	}

	// Transform Everflow response to domain model
	timeSeries := s.transformToTimeSeriesData(resp)

	// Cache the result
	if err := s.repo.SetCachedTimeSeries(ctx, cacheKey, timeSeries, 10*time.Minute); err != nil {
		logger.Warn("Failed to cache time series data", "error", err)
	}

	return timeSeries, nil
}

// GetDailyPerformanceReport retrieves detailed daily performance breakdown
func (s *reportingService) GetDailyPerformanceReport(ctx context.Context, filters domain.ReportingFilters, pagination domain.PaginationParams, userProfile *domain.Profile) ([]domain.DailyPerformanceReport, *domain.PaginationResult, error) {
	// Build Everflow entity report request with date and offer grouping
	req := s.buildEntityReportRequest(filters, []string{"date", "offer"})

	// Apply user-specific filters
	if err := s.applyUserFilters(&req, userProfile); err != nil {
		return nil, nil, fmt.Errorf("failed to apply user filters: %w", err)
	}

	// Get data from Everflow
	resp, err := s.everflowClient.GetEntityReport(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get entity report from Everflow: %w", err)
	}

	// Transform Everflow response to domain model
	dailyReports := s.transformToDailyPerformanceReport(resp)

	// Apply search filter if provided
	if filters.Search != nil && *filters.Search != "" {
		dailyReports = s.filterDailyReportsBySearch(dailyReports, *filters.Search)
	}

	// Apply sorting
	s.sortDailyReports(dailyReports, pagination.SortBy, pagination.SortOrder)

	// Apply pagination
	paginatedReports, paginationResult := s.paginateDailyReports(dailyReports, pagination)

	return paginatedReports, paginationResult, nil
}

// GetConversionsReport retrieves detailed conversion events
func (s *reportingService) GetConversionsReport(ctx context.Context, filters domain.ReportingFilters, pagination domain.PaginationParams, userProfile *domain.Profile) ([]domain.ConversionReport, *domain.PaginationResult, error) {
	// Build Everflow conversions request
	req := reporting.ConversionsRequest{
		ShowConversions: true,
		ShowEvents:      true,
		ShowOnlyVT:      false,
		ShowOnlyCT:      false,
		From:            filters.StartDate,
		To:              filters.EndDate,
		TimezoneID:      90, // UTC
		CurrencyID:      "USD",
		Query: reporting.ConversionsQuery{
			Filters:     s.buildEverflowFilters(filters, userProfile),
			SearchTerms: []string{},
		},
	}

	// Add search terms if provided
	if filters.Search != nil && *filters.Search != "" {
		req.Query.SearchTerms = []string{*filters.Search}
	}

	// Get data from Everflow
	resp, err := s.everflowClient.GetConversions(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get conversions from Everflow: %w", err)
	}

	// Transform Everflow response to domain model
	conversions := s.transformToConversionReport(resp.Conversions)

	// Apply status filter if provided
	if filters.Status != nil && *filters.Status != "all" {
		conversions = s.filterConversionsByStatus(conversions, *filters.Status)
	}

	// Apply sorting
	s.sortConversions(conversions, pagination.SortBy, pagination.SortOrder)

	// Apply pagination
	paginatedConversions, paginationResult := s.paginateConversions(conversions, pagination)

	return paginatedConversions, paginationResult, nil
}

// GetClicksReport retrieves detailed click events
func (s *reportingService) GetClicksReport(ctx context.Context, filters domain.ReportingFilters, pagination domain.PaginationParams, userProfile *domain.Profile) ([]domain.ClickReport, *domain.PaginationResult, error) {
	// Note: Everflow doesn't have a direct clicks endpoint in the advertiser API
	// This would need to be implemented using the raw clicks endpoint or entity reporting
	// For now, return empty data with a note

	logger.Warn("Clicks report not yet implemented - Everflow advertiser API doesn't provide raw clicks endpoint")

	return []domain.ClickReport{}, &domain.PaginationResult{
		CurrentPage:     pagination.Page,
		TotalPages:      0,
		TotalItems:      0,
		ItemsPerPage:    pagination.Limit,
		HasNextPage:     false,
		HasPreviousPage: false,
	}, nil
}

// GetCampaignsList retrieves campaigns for filter dropdown
func (s *reportingService) GetCampaignsList(ctx context.Context, affiliateID *string, status string, search *string, userProfile *domain.Profile) ([]domain.CampaignListItem, error) {
	var campaigns []domain.CampaignListItem
	var err error

	if affiliateID != nil {
		// Get campaigns visible to affiliate
		affID, parseErr := strconv.ParseInt(*affiliateID, 10, 64)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid affiliate ID: %w", parseErr)
		}
		campaigns, err = s.repo.GetCampaignsByAffiliate(ctx, affID, status)
	} else {
		// Get campaigns for user's organization
		if userProfile.OrganizationID == nil {
			return nil, fmt.Errorf("user is not associated with any organization")
		}
		campaigns, err = s.repo.GetCampaignsByOrganization(ctx, *userProfile.OrganizationID, status)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get campaigns list: %w", err)
	}

	// Apply search filter if provided
	if search != nil && *search != "" {
		campaigns = s.filterCampaignsBySearch(campaigns, *search)
	}

	return campaigns, nil
}

// Helper methods

func (s *reportingService) generateCacheKey(reportType string, filters domain.ReportingFilters, userProfile *domain.Profile) string {
	// Create a hash of the filters and user context for cache key
	data := fmt.Sprintf("%s:%s:%s:%v:%v:%d",
		reportType,
		filters.StartDate,
		filters.EndDate,
		filters.CampaignIDs,
		filters.AffiliateID,
		userProfile.OrganizationID,
	)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

func (s *reportingService) buildEntityReportRequest(filters domain.ReportingFilters, columns []string) reporting.EntityReportRequest {
	req := reporting.EntityReportRequest{
		From:       filters.StartDate,
		To:         filters.EndDate,
		TimezoneID: 90, // UTC
		CurrencyID: "USD",
		Query: reporting.EntityReportQuery{
			Filters:    []reporting.EntityReportFilter{},
			Exclusions: []reporting.EntityReportFilter{},
		},
		Columns: make([]reporting.EntityReportColumn, len(columns)),
	}

	// Set columns
	for i, col := range columns {
		req.Columns[i] = reporting.EntityReportColumn{Column: col}
	}

	return req
}

func (s *reportingService) applyUserFilters(req *reporting.EntityReportRequest, userProfile *domain.Profile) error {
	// Apply organization-specific filters based on user role
	// This would need to be implemented based on your business logic

	// For example, if user is an affiliate manager, only show their affiliates' data
	if userProfile.RoleName == "AffiliateManager" {
		// Add affiliate filter based on user's organization
		// This would require mapping organization to affiliate IDs
	}

	return nil
}

func (s *reportingService) buildEverflowFilters(filters domain.ReportingFilters, userProfile *domain.Profile) []reporting.EntityReportFilter {
	var everflowFilters []reporting.EntityReportFilter

	// Add campaign filters
	if len(filters.CampaignIDs) > 0 {
		for _, campaignID := range filters.CampaignIDs {
			everflowFilters = append(everflowFilters, reporting.EntityReportFilter{
				ResourceType:  "offer",
				FilterIDValue: campaignID,
			})
		}
	}

	// Add affiliate filter
	if filters.AffiliateID != nil {
		everflowFilters = append(everflowFilters, reporting.EntityReportFilter{
			ResourceType:  "affiliate",
			FilterIDValue: *filters.AffiliateID,
		})
	}

	return everflowFilters
}

func (s *reportingService) transformToPerformanceSummary(resp *reporting.EntityReportResponse) *domain.PerformanceSummary {
	summary := &domain.PerformanceSummary{}

	if resp.Summary != nil {
		// Extract metrics from Everflow summary
		if clicks, ok := resp.Summary["clicks"].(float64); ok {
			summary.TotalClicks = int64(clicks)
		}
		if conversions, ok := resp.Summary["conversions"].(float64); ok {
			summary.TotalConversions = int64(conversions)
		}
		if revenue, ok := resp.Summary["revenue"].(float64); ok {
			summary.TotalRevenue = revenue
		}
		if impressions, ok := resp.Summary["impressions"].(float64); ok {
			summary.TotalImpressions = int64(impressions)
		}

		// Calculate derived metrics
		if summary.TotalClicks > 0 {
			summary.ConversionRate = float64(summary.TotalConversions) / float64(summary.TotalClicks) * 100
		}
		if summary.TotalConversions > 0 {
			summary.AverageRevenue = summary.TotalRevenue / float64(summary.TotalConversions)
		}
		if summary.TotalImpressions > 0 {
			summary.ClickThroughRate = float64(summary.TotalClicks) / float64(summary.TotalImpressions) * 100
		}
	}

	return summary
}

func (s *reportingService) transformToTimeSeriesData(resp *reporting.EntityReportResponse) []domain.PerformanceTimeSeriesPoint {
	var timeSeries []domain.PerformanceTimeSeriesPoint

	for _, row := range resp.Table {
		point := domain.PerformanceTimeSeriesPoint{}

		// Extract date
		if date, ok := row["date"].(string); ok {
			point.Date = date
		}

		// Extract metrics
		if clicks, ok := row["clicks"].(float64); ok {
			point.Clicks = int64(clicks)
		}
		if conversions, ok := row["conversions"].(float64); ok {
			point.Conversions = int64(conversions)
		}
		if revenue, ok := row["revenue"].(float64); ok {
			point.Revenue = revenue
		}
		if impressions, ok := row["impressions"].(float64); ok {
			point.Impressions = int64(impressions)
		}

		// Calculate derived metrics
		if point.Clicks > 0 {
			point.ConversionRate = float64(point.Conversions) / float64(point.Clicks) * 100
		}
		if point.Impressions > 0 {
			point.ClickThroughRate = float64(point.Clicks) / float64(point.Impressions) * 100
		}

		timeSeries = append(timeSeries, point)
	}

	return timeSeries
}

func (s *reportingService) transformToDailyPerformanceReport(resp *reporting.EntityReportResponse) []domain.DailyPerformanceReport {
	var reports []domain.DailyPerformanceReport

	for _, row := range resp.Table {
		report := domain.DailyPerformanceReport{}

		// Extract basic fields
		if date, ok := row["date"].(string); ok {
			report.Date = date
		}
		if offerID, ok := row["offer_id"].(string); ok {
			report.CampaignID = offerID
		}
		if offerName, ok := row["offer_name"].(string); ok {
			report.CampaignName = offerName
		}

		// Extract metrics
		if clicks, ok := row["clicks"].(float64); ok {
			report.Clicks = int64(clicks)
		}
		if conversions, ok := row["conversions"].(float64); ok {
			report.Conversions = int64(conversions)
		}
		if revenue, ok := row["revenue"].(float64); ok {
			report.Revenue = revenue
		}
		if impressions, ok := row["impressions"].(float64); ok {
			report.Impressions = int64(impressions)
		}
		if payouts, ok := row["payouts"].(float64); ok {
			report.Payouts = payouts
		}

		// Calculate derived metrics
		if report.Clicks > 0 {
			report.ConversionRate = float64(report.Conversions) / float64(report.Clicks) * 100
		}
		if report.Impressions > 0 {
			report.ClickThroughRate = float64(report.Clicks) / float64(report.Impressions) * 100
		}

		reports = append(reports, report)
	}

	return reports
}

func (s *reportingService) transformToConversionReport(conversions []reporting.ConversionData) []domain.ConversionReport {
	var reports []domain.ConversionReport

	for _, conv := range conversions {
		report := domain.ConversionReport{
			ID:            conv.ConversionID,
			Timestamp:     time.Unix(conv.ConversionUnixTimestamp, 0),
			TransactionID: conv.TransactionID,
			CampaignID:    strconv.Itoa(conv.Relationship.Offer.NetworkOfferID),
			CampaignName:  conv.Relationship.Offer.Name,
			OfferID:       strconv.Itoa(conv.Relationship.Offer.NetworkOfferID),
			OfferName:     conv.Relationship.Offer.Name,
			Status:        s.mapConversionStatus(conv.Relationship.Offer.OfferStatus),
			Payout:        conv.Cost,
			Currency:      conv.CurrencyID,
			AffiliateID:   strconv.Itoa(conv.Relationship.AffiliateID),
			AffiliateName: conv.Relationship.Affiliate.Name,
		}

		// Set optional fields
		if conv.Relationship.Sub1 != "" {
			report.Sub1 = &conv.Relationship.Sub1
		}
		if conv.Relationship.Sub2 != "" {
			report.Sub2 = &conv.Relationship.Sub2
		}
		if conv.Relationship.Sub3 != "" {
			report.Sub3 = &conv.Relationship.Sub3
		}
		if conv.SaleAmount > 0 {
			report.ConversionValue = &conv.SaleAmount
		}

		reports = append(reports, report)
	}

	return reports
}

func (s *reportingService) mapConversionStatus(everflowStatus string) string {
	// Map Everflow status to our domain status
	switch everflowStatus {
	case "active":
		return "approved"
	case "paused":
		return "pending"
	default:
		return "pending"
	}
}

func (s *reportingService) filterDailyReportsBySearch(reports []domain.DailyPerformanceReport, search string) []domain.DailyPerformanceReport {
	search = strings.ToLower(search)
	var filtered []domain.DailyPerformanceReport

	for _, report := range reports {
		if strings.Contains(strings.ToLower(report.CampaignName), search) ||
			strings.Contains(strings.ToLower(report.CampaignID), search) {
			filtered = append(filtered, report)
		}
	}

	return filtered
}

func (s *reportingService) sortDailyReports(reports []domain.DailyPerformanceReport, sortBy, sortOrder string) {
	sort.Slice(reports, func(i, j int) bool {
		var less bool
		switch sortBy {
		case "clicks":
			less = reports[i].Clicks < reports[j].Clicks
		case "conversions":
			less = reports[i].Conversions < reports[j].Conversions
		case "revenue":
			less = reports[i].Revenue < reports[j].Revenue
		default: // date
			less = reports[i].Date < reports[j].Date
		}

		if sortOrder == "asc" {
			return less
		}
		return !less
	})
}

func (s *reportingService) paginateDailyReports(reports []domain.DailyPerformanceReport, pagination domain.PaginationParams) ([]domain.DailyPerformanceReport, *domain.PaginationResult) {
	totalItems := len(reports)
	totalPages := (totalItems + pagination.Limit - 1) / pagination.Limit

	start := (pagination.Page - 1) * pagination.Limit
	end := start + pagination.Limit

	if start >= totalItems {
		return []domain.DailyPerformanceReport{}, &domain.PaginationResult{
			CurrentPage:     pagination.Page,
			TotalPages:      totalPages,
			TotalItems:      totalItems,
			ItemsPerPage:    pagination.Limit,
			HasNextPage:     false,
			HasPreviousPage: pagination.Page > 1,
		}
	}

	if end > totalItems {
		end = totalItems
	}

	paginatedReports := reports[start:end]

	return paginatedReports, &domain.PaginationResult{
		CurrentPage:     pagination.Page,
		TotalPages:      totalPages,
		TotalItems:      totalItems,
		ItemsPerPage:    pagination.Limit,
		HasNextPage:     pagination.Page < totalPages,
		HasPreviousPage: pagination.Page > 1,
	}
}

func (s *reportingService) filterConversionsByStatus(conversions []domain.ConversionReport, status string) []domain.ConversionReport {
	var filtered []domain.ConversionReport

	for _, conv := range conversions {
		if conv.Status == status {
			filtered = append(filtered, conv)
		}
	}

	return filtered
}

func (s *reportingService) sortConversions(conversions []domain.ConversionReport, sortBy, sortOrder string) {
	sort.Slice(conversions, func(i, j int) bool {
		var less bool
		switch sortBy {
		case "payout":
			less = conversions[i].Payout < conversions[j].Payout
		case "campaign":
			less = conversions[i].CampaignName < conversions[j].CampaignName
		case "status":
			less = conversions[i].Status < conversions[j].Status
		default: // timestamp
			less = conversions[i].Timestamp.Before(conversions[j].Timestamp)
		}

		if sortOrder == "asc" {
			return less
		}
		return !less
	})
}

func (s *reportingService) paginateConversions(conversions []domain.ConversionReport, pagination domain.PaginationParams) ([]domain.ConversionReport, *domain.PaginationResult) {
	totalItems := len(conversions)
	totalPages := (totalItems + pagination.Limit - 1) / pagination.Limit

	start := (pagination.Page - 1) * pagination.Limit
	end := start + pagination.Limit

	if start >= totalItems {
		return []domain.ConversionReport{}, &domain.PaginationResult{
			CurrentPage:     pagination.Page,
			TotalPages:      totalPages,
			TotalItems:      totalItems,
			ItemsPerPage:    pagination.Limit,
			HasNextPage:     false,
			HasPreviousPage: pagination.Page > 1,
		}
	}

	if end > totalItems {
		end = totalItems
	}

	paginatedConversions := conversions[start:end]

	return paginatedConversions, &domain.PaginationResult{
		CurrentPage:     pagination.Page,
		TotalPages:      totalPages,
		TotalItems:      totalItems,
		ItemsPerPage:    pagination.Limit,
		HasNextPage:     pagination.Page < totalPages,
		HasPreviousPage: pagination.Page > 1,
	}
}

func (s *reportingService) filterCampaignsBySearch(campaigns []domain.CampaignListItem, search string) []domain.CampaignListItem {
	search = strings.ToLower(search)
	var filtered []domain.CampaignListItem

	for _, campaign := range campaigns {
		if strings.Contains(strings.ToLower(campaign.Name), search) ||
			strings.Contains(strings.ToLower(campaign.ID), search) {
			filtered = append(filtered, campaign)
		}
	}

	return filtered
}