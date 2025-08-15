package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

// ReportingRepository handles reporting data persistence and caching
type ReportingRepository interface {
	// Cache management
	GetCachedSummary(ctx context.Context, key string) (*domain.PerformanceSummary, error)
	SetCachedSummary(ctx context.Context, key string, summary *domain.PerformanceSummary, ttl time.Duration) error
	GetCachedTimeSeries(ctx context.Context, key string) ([]domain.PerformanceTimeSeriesPoint, error)
	SetCachedTimeSeries(ctx context.Context, key string, data []domain.PerformanceTimeSeriesPoint, ttl time.Duration) error

	// Campaign data
	GetCampaignsByOrganization(ctx context.Context, organizationID int64, status string) ([]domain.CampaignListItem, error)
	GetCampaignsByAffiliate(ctx context.Context, affiliateID int64, status string) ([]domain.CampaignListItem, error)
}

type reportingRepository struct {
	db *sqlx.DB
	// Add Redis client here if using Redis for caching
}

// NewReportingRepository creates a new reporting repository
func NewReportingRepository(db *sqlx.DB) ReportingRepository {
	return &reportingRepository{
		db: db,
	}
}

// GetCachedSummary retrieves cached performance summary
func (r *reportingRepository) GetCachedSummary(ctx context.Context, key string) (*domain.PerformanceSummary, error) {
	// Implementation depends on caching strategy (Redis, in-memory, etc.)
	// For now, return nil to indicate no cache hit
	return nil, nil
}

// SetCachedSummary stores performance summary in cache
func (r *reportingRepository) SetCachedSummary(ctx context.Context, key string, summary *domain.PerformanceSummary, ttl time.Duration) error {
	// Implementation depends on caching strategy
	return nil
}

// GetCachedTimeSeries retrieves cached time series data
func (r *reportingRepository) GetCachedTimeSeries(ctx context.Context, key string) ([]domain.PerformanceTimeSeriesPoint, error) {
	// Implementation depends on caching strategy
	return nil, nil
}

// SetCachedTimeSeries stores time series data in cache
func (r *reportingRepository) SetCachedTimeSeries(ctx context.Context, key string, data []domain.PerformanceTimeSeriesPoint, ttl time.Duration) error {
	// Implementation depends on caching strategy
	return nil
}

// GetCampaignsByOrganization retrieves campaigns for an organization
func (r *reportingRepository) GetCampaignsByOrganization(ctx context.Context, organizationID int64, status string) ([]domain.CampaignListItem, error) {
	query := `
		SELECT campaign_id, name, status 
		FROM campaigns 
		WHERE organization_id = $1
	`
	args := []interface{}{organizationID}

	if status != "all" {
		query += " AND status = $2"
		args = append(args, status)
	}

	query += " ORDER BY name"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query campaigns by organization: %w", err)
	}
	defer rows.Close()

	var campaigns []domain.CampaignListItem
	for rows.Next() {
		var campaign domain.CampaignListItem
		var campaignID int64
		
		err := rows.Scan(&campaignID, &campaign.Name, &campaign.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan campaign row: %w", err)
		}
		
		campaign.ID = strconv.FormatInt(campaignID, 10)
		campaigns = append(campaigns, campaign)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating campaign rows: %w", err)
	}

	return campaigns, nil
}

// GetCampaignsByAffiliate retrieves campaigns visible to an affiliate
func (r *reportingRepository) GetCampaignsByAffiliate(ctx context.Context, affiliateID int64, status string) ([]domain.CampaignListItem, error) {
	// This would need to join with organization associations to get visible campaigns
	query := `
		SELECT DISTINCT c.campaign_id, c.name, c.status 
		FROM campaigns c
		INNER JOIN organization_associations oa ON c.organization_id = oa.advertiser_organization_id
		WHERE oa.affiliate_organization_id = (
			SELECT organization_id FROM affiliates WHERE affiliate_id = $1
		)
		AND oa.status = 'active'
	`
	args := []interface{}{affiliateID}

	if status != "all" {
		query += " AND c.status = $2"
		args = append(args, status)
	}

	query += " ORDER BY c.name"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query campaigns by affiliate: %w", err)
	}
	defer rows.Close()

	var campaigns []domain.CampaignListItem
	for rows.Next() {
		var campaign domain.CampaignListItem
		var campaignID int64
		
		err := rows.Scan(&campaignID, &campaign.Name, &campaign.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan campaign row: %w", err)
		}
		
		campaign.ID = strconv.FormatInt(campaignID, 10)
		campaigns = append(campaigns, campaign)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating campaign rows: %w", err)
	}

	return campaigns, nil
}