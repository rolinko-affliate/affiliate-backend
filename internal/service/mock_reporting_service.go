package service

import (
	"context"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/logger"
)

// MockReportingService implements ReportingService using mock data
type MockReportingService struct {
	mockDataService MockDataService
	logger          *logger.Logger
}

// NewMockReportingService creates a new mock reporting service
func NewMockReportingService(mockDataService MockDataService, logger *logger.Logger) ReportingService {
	return &MockReportingService{
		mockDataService: mockDataService,
		logger:          logger,
	}
}

// GetPerformanceSummary retrieves aggregated performance metrics from mock data
func (s *MockReportingService) GetPerformanceSummary(ctx context.Context, filters domain.ReportingFilters, userProfile *domain.Profile) (*domain.PerformanceSummary, error) {
	orgID := int64(1) // Default organization ID
	if userProfile.OrganizationID != nil {
		orgID = *userProfile.OrganizationID
	}
	s.logger.Info("Getting performance summary from mock data", "org_id", orgID)
	return s.mockDataService.LoadPerformanceSummary(ctx, orgID, filters)
}

// GetPerformanceTimeSeries retrieves daily performance data for charts from mock data
func (s *MockReportingService) GetPerformanceTimeSeries(ctx context.Context, filters domain.ReportingFilters, userProfile *domain.Profile) ([]domain.PerformanceTimeSeriesPoint, error) {
	orgID := int64(1) // Default organization ID
	if userProfile.OrganizationID != nil {
		orgID = *userProfile.OrganizationID
	}
	s.logger.Info("Getting performance timeseries from mock data", "org_id", orgID)
	return s.mockDataService.LoadPerformanceTimeSeries(ctx, orgID, filters)
}

// GetDailyPerformanceReport retrieves daily performance breakdown from mock data
func (s *MockReportingService) GetDailyPerformanceReport(ctx context.Context, filters domain.ReportingFilters, pagination domain.PaginationParams, userProfile *domain.Profile) ([]domain.DailyPerformanceReport, *domain.PaginationResult, error) {
	orgID := int64(1) // Default organization ID
	if userProfile.OrganizationID != nil {
		orgID = *userProfile.OrganizationID
	}
	s.logger.Info("Getting daily performance report from mock data", "org_id", orgID)
	return s.mockDataService.LoadDailyPerformanceReport(ctx, orgID, filters, pagination)
}

// GetConversionsReport retrieves detailed conversion events from mock data
func (s *MockReportingService) GetConversionsReport(ctx context.Context, filters domain.ReportingFilters, pagination domain.PaginationParams, userProfile *domain.Profile) ([]domain.ConversionReport, *domain.PaginationResult, error) {
	orgID := int64(1) // Default organization ID
	if userProfile.OrganizationID != nil {
		orgID = *userProfile.OrganizationID
	}
	s.logger.Info("Getting conversions report from mock data", "org_id", orgID)
	return s.mockDataService.LoadConversionsReport(ctx, orgID, filters, pagination)
}

// GetClicksReport retrieves detailed click events from mock data
func (s *MockReportingService) GetClicksReport(ctx context.Context, filters domain.ReportingFilters, pagination domain.PaginationParams, userProfile *domain.Profile) ([]domain.ClickReport, *domain.PaginationResult, error) {
	orgID := int64(1) // Default organization ID
	if userProfile.OrganizationID != nil {
		orgID = *userProfile.OrganizationID
	}
	s.logger.Info("Getting clicks report from mock data", "org_id", orgID)
	return s.mockDataService.LoadClicksReport(ctx, orgID, filters, pagination)
}

// GetCampaignsList retrieves campaigns list for filter dropdown from mock data
func (s *MockReportingService) GetCampaignsList(ctx context.Context, affiliateID *string, status string, search *string, userProfile *domain.Profile) ([]domain.CampaignListItem, error) {
	orgID := int64(1) // Default organization ID
	if userProfile.OrganizationID != nil {
		orgID = *userProfile.OrganizationID
	}
	s.logger.Info("Getting campaigns list from mock data", "org_id", orgID)
	return s.mockDataService.LoadCampaignsList(ctx, orgID, affiliateID, status, search)
}
