package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UsageRecordRepository defines the interface for usage record operations
type UsageRecordRepository interface {
	Create(ctx context.Context, record *domain.UsageRecord) error
	GetByID(ctx context.Context, usageRecordID int64) (*domain.UsageRecord, error)
	GetByOrganizationAndDate(ctx context.Context, organizationID int64, date time.Time) (*domain.UsageRecord, error)
	GetByOrganizationID(ctx context.Context, organizationID int64, limit, offset int) ([]domain.UsageRecord, error)
	GetByDateRange(ctx context.Context, organizationID int64, startDate, endDate time.Time) ([]domain.UsageRecord, error)
	GetPendingRecords(ctx context.Context, limit int) ([]domain.UsageRecord, error)
	Update(ctx context.Context, record *domain.UsageRecord) error
	List(ctx context.Context, limit, offset int) ([]domain.UsageRecord, error)
	GetMonthlyUsage(ctx context.Context, organizationID int64, year int, month int) ([]domain.UsageRecord, error)
}

// PgxUsageRecordRepository implements UsageRecordRepository using pgx
type PgxUsageRecordRepository struct {
	db *pgxpool.Pool
}

// NewPgxUsageRecordRepository creates a new PgxUsageRecordRepository
func NewPgxUsageRecordRepository(db *pgxpool.Pool) UsageRecordRepository {
	return &PgxUsageRecordRepository{db: db}
}

// Create creates a new usage record
func (r *PgxUsageRecordRepository) Create(ctx context.Context, record *domain.UsageRecord) error {
	query := `
		INSERT INTO usage_records (
			organization_id, billing_account_id, usage_date, clicks, conversions, impressions,
			advertiser_spend, affiliate_payout, platform_revenue, currency, status,
			allocated_at, billed_at, campaign_breakdown, affiliate_breakdown, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		) RETURNING usage_record_id, created_at, updated_at`

	campaignBreakdownJSON, err := json.Marshal(record.CampaignBreakdown)
	if err != nil {
		return fmt.Errorf("failed to marshal campaign breakdown: %w", err)
	}

	affiliateBreakdownJSON, err := json.Marshal(record.AffiliateBreakdown)
	if err != nil {
		return fmt.Errorf("failed to marshal affiliate breakdown: %w", err)
	}

	metadataJSON, err := json.Marshal(record.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	err = r.db.QueryRow(ctx, query,
		record.OrganizationID,
		record.BillingAccountID,
		record.UsageDate,
		record.Clicks,
		record.Conversions,
		record.Impressions,
		record.AdvertiserSpend,
		record.AffiliatePayout,
		record.PlatformRevenue,
		record.Currency,
		record.Status,
		record.AllocatedAt,
		record.BilledAt,
		campaignBreakdownJSON,
		affiliateBreakdownJSON,
		metadataJSON,
	).Scan(&record.UsageRecordID, &record.CreatedAt, &record.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create usage record: %w", err)
	}

	return nil
}

// GetByID retrieves a usage record by ID
func (r *PgxUsageRecordRepository) GetByID(ctx context.Context, usageRecordID int64) (*domain.UsageRecord, error) {
	query := `
		SELECT usage_record_id, organization_id, billing_account_id, usage_date, clicks,
			   conversions, impressions, advertiser_spend, affiliate_payout, platform_revenue,
			   currency, status, allocated_at, billed_at, campaign_breakdown,
			   affiliate_breakdown, metadata, created_at, updated_at
		FROM usage_records
		WHERE usage_record_id = $1`

	record := &domain.UsageRecord{}
	var campaignBreakdownJSON, affiliateBreakdownJSON, metadataJSON []byte

	err := r.db.QueryRow(ctx, query, usageRecordID).Scan(
		&record.UsageRecordID,
		&record.OrganizationID,
		&record.BillingAccountID,
		&record.UsageDate,
		&record.Clicks,
		&record.Conversions,
		&record.Impressions,
		&record.AdvertiserSpend,
		&record.AffiliatePayout,
		&record.PlatformRevenue,
		&record.Currency,
		&record.Status,
		&record.AllocatedAt,
		&record.BilledAt,
		&campaignBreakdownJSON,
		&affiliateBreakdownJSON,
		&metadataJSON,
		&record.CreatedAt,
		&record.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("usage record not found")
		}
		return nil, fmt.Errorf("failed to get usage record: %w", err)
	}

	// Unmarshal JSON fields
	if len(campaignBreakdownJSON) > 0 {
		if err := json.Unmarshal(campaignBreakdownJSON, &record.CampaignBreakdown); err != nil {
			return nil, fmt.Errorf("failed to unmarshal campaign breakdown: %w", err)
		}
	}

	if len(affiliateBreakdownJSON) > 0 {
		if err := json.Unmarshal(affiliateBreakdownJSON, &record.AffiliateBreakdown); err != nil {
			return nil, fmt.Errorf("failed to unmarshal affiliate breakdown: %w", err)
		}
	}

	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &record.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return record, nil
}

// GetByOrganizationAndDate retrieves a usage record by organization and date
func (r *PgxUsageRecordRepository) GetByOrganizationAndDate(ctx context.Context, organizationID int64, date time.Time) (*domain.UsageRecord, error) {
	query := `
		SELECT usage_record_id, organization_id, billing_account_id, usage_date, clicks,
			   conversions, impressions, advertiser_spend, affiliate_payout, platform_revenue,
			   currency, status, allocated_at, billed_at, campaign_breakdown,
			   affiliate_breakdown, metadata, created_at, updated_at
		FROM usage_records
		WHERE organization_id = $1 AND usage_date = $2`

	record := &domain.UsageRecord{}
	var campaignBreakdownJSON, affiliateBreakdownJSON, metadataJSON []byte

	err := r.db.QueryRow(ctx, query, organizationID, date).Scan(
		&record.UsageRecordID,
		&record.OrganizationID,
		&record.BillingAccountID,
		&record.UsageDate,
		&record.Clicks,
		&record.Conversions,
		&record.Impressions,
		&record.AdvertiserSpend,
		&record.AffiliatePayout,
		&record.PlatformRevenue,
		&record.Currency,
		&record.Status,
		&record.AllocatedAt,
		&record.BilledAt,
		&campaignBreakdownJSON,
		&affiliateBreakdownJSON,
		&metadataJSON,
		&record.CreatedAt,
		&record.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("usage record not found")
		}
		return nil, fmt.Errorf("failed to get usage record: %w", err)
	}

	// Unmarshal JSON fields
	if len(campaignBreakdownJSON) > 0 {
		if err := json.Unmarshal(campaignBreakdownJSON, &record.CampaignBreakdown); err != nil {
			return nil, fmt.Errorf("failed to unmarshal campaign breakdown: %w", err)
		}
	}

	if len(affiliateBreakdownJSON) > 0 {
		if err := json.Unmarshal(affiliateBreakdownJSON, &record.AffiliateBreakdown); err != nil {
			return nil, fmt.Errorf("failed to unmarshal affiliate breakdown: %w", err)
		}
	}

	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &record.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return record, nil
}

// GetByOrganizationID retrieves usage records for an organization with pagination
func (r *PgxUsageRecordRepository) GetByOrganizationID(ctx context.Context, organizationID int64, limit, offset int) ([]domain.UsageRecord, error) {
	query := `
		SELECT usage_record_id, organization_id, billing_account_id, usage_date, clicks,
			   conversions, impressions, advertiser_spend, affiliate_payout, platform_revenue,
			   currency, status, allocated_at, billed_at, campaign_breakdown,
			   affiliate_breakdown, metadata, created_at, updated_at
		FROM usage_records
		WHERE organization_id = $1
		ORDER BY usage_date DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage records by organization: %w", err)
	}
	defer rows.Close()

	return r.scanUsageRecords(rows)
}

// GetByDateRange retrieves usage records for an organization within a date range
func (r *PgxUsageRecordRepository) GetByDateRange(ctx context.Context, organizationID int64, startDate, endDate time.Time) ([]domain.UsageRecord, error) {
	query := `
		SELECT usage_record_id, organization_id, billing_account_id, usage_date, clicks,
			   conversions, impressions, advertiser_spend, affiliate_payout, platform_revenue,
			   currency, status, allocated_at, billed_at, campaign_breakdown,
			   affiliate_breakdown, metadata, created_at, updated_at
		FROM usage_records
		WHERE organization_id = $1 AND usage_date >= $2 AND usage_date <= $3
		ORDER BY usage_date DESC`

	rows, err := r.db.Query(ctx, query, organizationID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage records by date range: %w", err)
	}
	defer rows.Close()

	return r.scanUsageRecords(rows)
}

// GetPendingRecords retrieves pending usage records for processing
func (r *PgxUsageRecordRepository) GetPendingRecords(ctx context.Context, limit int) ([]domain.UsageRecord, error) {
	query := `
		SELECT usage_record_id, organization_id, billing_account_id, usage_date, clicks,
			   conversions, impressions, advertiser_spend, affiliate_payout, platform_revenue,
			   currency, status, allocated_at, billed_at, campaign_breakdown,
			   affiliate_breakdown, metadata, created_at, updated_at
		FROM usage_records
		WHERE status = 'pending'
		ORDER BY usage_date ASC
		LIMIT $1`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending usage records: %w", err)
	}
	defer rows.Close()

	return r.scanUsageRecords(rows)
}

// Update updates a usage record
func (r *PgxUsageRecordRepository) Update(ctx context.Context, record *domain.UsageRecord) error {
	query := `
		UPDATE usage_records SET
			clicks = $2,
			conversions = $3,
			impressions = $4,
			advertiser_spend = $5,
			affiliate_payout = $6,
			platform_revenue = $7,
			currency = $8,
			status = $9,
			allocated_at = $10,
			billed_at = $11,
			campaign_breakdown = $12,
			affiliate_breakdown = $13,
			metadata = $14,
			updated_at = NOW()
		WHERE usage_record_id = $1
		RETURNING updated_at`

	campaignBreakdownJSON, err := json.Marshal(record.CampaignBreakdown)
	if err != nil {
		return fmt.Errorf("failed to marshal campaign breakdown: %w", err)
	}

	affiliateBreakdownJSON, err := json.Marshal(record.AffiliateBreakdown)
	if err != nil {
		return fmt.Errorf("failed to marshal affiliate breakdown: %w", err)
	}

	metadataJSON, err := json.Marshal(record.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	err = r.db.QueryRow(ctx, query,
		record.UsageRecordID,
		record.Clicks,
		record.Conversions,
		record.Impressions,
		record.AdvertiserSpend,
		record.AffiliatePayout,
		record.PlatformRevenue,
		record.Currency,
		record.Status,
		record.AllocatedAt,
		record.BilledAt,
		campaignBreakdownJSON,
		affiliateBreakdownJSON,
		metadataJSON,
	).Scan(&record.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update usage record: %w", err)
	}

	return nil
}

// List retrieves a list of usage records with pagination
func (r *PgxUsageRecordRepository) List(ctx context.Context, limit, offset int) ([]domain.UsageRecord, error) {
	query := `
		SELECT usage_record_id, organization_id, billing_account_id, usage_date, clicks,
			   conversions, impressions, advertiser_spend, affiliate_payout, platform_revenue,
			   currency, status, allocated_at, billed_at, campaign_breakdown,
			   affiliate_breakdown, metadata, created_at, updated_at
		FROM usage_records
		ORDER BY usage_date DESC, created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list usage records: %w", err)
	}
	defer rows.Close()

	return r.scanUsageRecords(rows)
}

// GetMonthlyUsage retrieves usage records for a specific month
func (r *PgxUsageRecordRepository) GetMonthlyUsage(ctx context.Context, organizationID int64, year int, month int) ([]domain.UsageRecord, error) {
	query := `
		SELECT usage_record_id, organization_id, billing_account_id, usage_date, clicks,
			   conversions, impressions, advertiser_spend, affiliate_payout, platform_revenue,
			   currency, status, allocated_at, billed_at, campaign_breakdown,
			   affiliate_breakdown, metadata, created_at, updated_at
		FROM usage_records
		WHERE organization_id = $1 
		AND EXTRACT(YEAR FROM usage_date) = $2
		AND EXTRACT(MONTH FROM usage_date) = $3
		ORDER BY usage_date ASC`

	rows, err := r.db.Query(ctx, query, organizationID, year, month)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly usage records: %w", err)
	}
	defer rows.Close()

	return r.scanUsageRecords(rows)
}

// scanUsageRecords is a helper method to scan multiple usage records from rows
func (r *PgxUsageRecordRepository) scanUsageRecords(rows pgx.Rows) ([]domain.UsageRecord, error) {
	var records []domain.UsageRecord
	for rows.Next() {
		record := domain.UsageRecord{}
		var campaignBreakdownJSON, affiliateBreakdownJSON, metadataJSON []byte

		err := rows.Scan(
			&record.UsageRecordID,
			&record.OrganizationID,
			&record.BillingAccountID,
			&record.UsageDate,
			&record.Clicks,
			&record.Conversions,
			&record.Impressions,
			&record.AdvertiserSpend,
			&record.AffiliatePayout,
			&record.PlatformRevenue,
			&record.Currency,
			&record.Status,
			&record.AllocatedAt,
			&record.BilledAt,
			&campaignBreakdownJSON,
			&affiliateBreakdownJSON,
			&metadataJSON,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan usage record: %w", err)
		}

		// Unmarshal JSON fields
		if len(campaignBreakdownJSON) > 0 {
			if err := json.Unmarshal(campaignBreakdownJSON, &record.CampaignBreakdown); err != nil {
				return nil, fmt.Errorf("failed to unmarshal campaign breakdown: %w", err)
			}
		}

		if len(affiliateBreakdownJSON) > 0 {
			if err := json.Unmarshal(affiliateBreakdownJSON, &record.AffiliateBreakdown); err != nil {
				return nil, fmt.Errorf("failed to unmarshal affiliate breakdown: %w", err)
			}
		}

		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &record.Metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating usage records: %w", err)
	}

	return records, nil
}