package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdvertiserRepository interface {
	CreateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error
	GetAdvertiserByID(ctx context.Context, id int64) (*domain.Advertiser, error)
	UpdateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error
	ListAdvertisersByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.Advertiser, error)
	DeleteAdvertiser(ctx context.Context, id int64) error
	
	// Extra info methods
	CreateAdvertiserExtraInfo(ctx context.Context, extraInfo *domain.AdvertiserExtraInfo) error
	GetAdvertiserExtraInfo(ctx context.Context, organizationID int64) (*domain.AdvertiserExtraInfo, error)
	UpdateAdvertiserExtraInfo(ctx context.Context, extraInfo *domain.AdvertiserExtraInfo) error
	UpsertAdvertiserExtraInfo(ctx context.Context, extraInfo *domain.AdvertiserExtraInfo) error
	DeleteAdvertiserExtraInfo(ctx context.Context, organizationID int64) error
	GetAdvertiserWithExtraInfo(ctx context.Context, advertiserID int64) (*domain.AdvertiserWithExtraInfo, error)
}

type pgxAdvertiserRepository struct {
	db *pgxpool.Pool
}

func NewPgxAdvertiserRepository(db *pgxpool.Pool) AdvertiserRepository {
	return &pgxAdvertiserRepository{db: db}
}

func (r *pgxAdvertiserRepository) CreateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error {
	query := `INSERT INTO public.advertisers (
		organization_id, name, contact_email, billing_details, status,
		internal_notes, default_currency_id, platform_name, platform_url, platform_username,
		accounting_contact_email, offer_id_macro, affiliate_id_macro, attribution_method,
		email_attribution_method, attribution_priority, reporting_timezone_id,
		created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
	) RETURNING advertiser_id, created_at, updated_at`

	billingDetailsJSON, err := marshalBillingDetails(advertiser.BillingDetails)
	if err != nil {
		return err
	}

	now := time.Now()
	err = r.db.QueryRow(ctx, query,
		advertiser.OrganizationID,
		advertiser.Name,
		advertiser.ContactEmail,
		billingDetailsJSON,
		advertiser.Status,
		advertiser.InternalNotes,
		advertiser.DefaultCurrencyID,
		advertiser.PlatformName,
		advertiser.PlatformURL,
		advertiser.PlatformUsername,
		advertiser.AccountingContactEmail,
		advertiser.OfferIDMacro,
		advertiser.AffiliateIDMacro,
		advertiser.AttributionMethod,
		advertiser.EmailAttributionMethod,
		advertiser.AttributionPriority,
		advertiser.ReportingTimezoneID,
		now,
		now,
	).Scan(
		&advertiser.AdvertiserID,
		&advertiser.CreatedAt,
		&advertiser.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating advertiser: %w", err)
	}
	return nil
}

func (r *pgxAdvertiserRepository) GetAdvertiserByID(ctx context.Context, id int64) (*domain.Advertiser, error) {
	query := `SELECT ` + advertiserSelectFields + ` FROM public.advertisers WHERE advertiser_id = $1`

	var advertiser domain.Advertiser
	var contactEmail, billingDetails sql.NullString
	var internalNotes, defaultCurrencyID, platformName, platformURL, platformUsername sql.NullString
	var accountingContactEmail, offerIDMacro, affiliateIDMacro sql.NullString
	var attributionMethod, emailAttributionMethod, attributionPriority sql.NullString
	var reportingTimezoneID sql.NullInt32

	err := r.db.QueryRow(ctx, query, id).Scan(
		&advertiser.AdvertiserID,
		&advertiser.OrganizationID,
		&advertiser.Name,
		&contactEmail,
		&billingDetails,
		&advertiser.Status,
		&internalNotes,
		&defaultCurrencyID,
		&platformName,
		&platformURL,
		&platformUsername,
		&accountingContactEmail,
		&offerIDMacro,
		&affiliateIDMacro,
		&attributionMethod,
		&emailAttributionMethod,
		&attributionPriority,
		&reportingTimezoneID,
		&advertiser.CreatedAt,
		&advertiser.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("advertiser not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("error getting advertiser by ID: %w", err)
	}

	err = scanNullableFields(&advertiser, contactEmail, billingDetails, internalNotes, defaultCurrencyID,
		platformName, platformURL, platformUsername, accountingContactEmail, offerIDMacro, affiliateIDMacro,
		attributionMethod, emailAttributionMethod, attributionPriority, reportingTimezoneID)
	if err != nil {
		return nil, err
	}

	return &advertiser, nil
}

func (r *pgxAdvertiserRepository) UpdateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error {
	query := `UPDATE public.advertisers SET 
		name = $1, contact_email = $2, billing_details = $3, status = $4,
		internal_notes = $5, default_currency_id = $6, platform_name = $7, platform_url = $8, platform_username = $9,
		accounting_contact_email = $10, offer_id_macro = $11, affiliate_id_macro = $12, attribution_method = $13,
		email_attribution_method = $14, attribution_priority = $15, reporting_timezone_id = $16,
		updated_at = $17
	WHERE advertiser_id = $18
	RETURNING updated_at`

	billingDetailsJSON, err := marshalBillingDetails(advertiser.BillingDetails)
	if err != nil {
		return err
	}

	now := time.Now()
	err = r.db.QueryRow(ctx, query,
		advertiser.Name,
		advertiser.ContactEmail,
		billingDetailsJSON,
		advertiser.Status,
		advertiser.InternalNotes,
		advertiser.DefaultCurrencyID,
		advertiser.PlatformName,
		advertiser.PlatformURL,
		advertiser.PlatformUsername,
		advertiser.AccountingContactEmail,
		advertiser.OfferIDMacro,
		advertiser.AffiliateIDMacro,
		advertiser.AttributionMethod,
		advertiser.EmailAttributionMethod,
		advertiser.AttributionPriority,
		advertiser.ReportingTimezoneID,
		now,
		advertiser.AdvertiserID,
	).Scan(&advertiser.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("advertiser not found: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("error updating advertiser: %w", err)
	}

	return nil
}

func (r *pgxAdvertiserRepository) ListAdvertisersByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.Advertiser, error) {
	query := `SELECT ` + advertiserSelectFields + `
	FROM public.advertisers
	WHERE organization_id = $1
	ORDER BY advertiser_id
	LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing advertisers: %w", err)
	}
	defer rows.Close()

	var advertisers []*domain.Advertiser
	for rows.Next() {
		var advertiser domain.Advertiser
		var contactEmail, billingDetails sql.NullString
		var internalNotes, defaultCurrencyID, platformName, platformURL, platformUsername sql.NullString
		var accountingContactEmail, offerIDMacro, affiliateIDMacro sql.NullString
		var attributionMethod, emailAttributionMethod, attributionPriority sql.NullString
		var reportingTimezoneID sql.NullInt32

		if err := rows.Scan(
			&advertiser.AdvertiserID,
			&advertiser.OrganizationID,
			&advertiser.Name,
			&contactEmail,
			&billingDetails,
			&advertiser.Status,
			&internalNotes,
			&defaultCurrencyID,
			&platformName,
			&platformURL,
			&platformUsername,
			&accountingContactEmail,
			&offerIDMacro,
			&affiliateIDMacro,
			&attributionMethod,
			&emailAttributionMethod,
			&attributionPriority,
			&reportingTimezoneID,
			&advertiser.CreatedAt,
			&advertiser.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning advertiser row: %w", err)
		}

		err = scanNullableFields(&advertiser, contactEmail, billingDetails, internalNotes, defaultCurrencyID,
			platformName, platformURL, platformUsername, accountingContactEmail, offerIDMacro, affiliateIDMacro,
			attributionMethod, emailAttributionMethod, attributionPriority, reportingTimezoneID)
		if err != nil {
			return nil, err
		}

		advertisers = append(advertisers, &advertiser)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating advertiser rows: %w", err)
	}

	return advertisers, nil
}

func (r *pgxAdvertiserRepository) DeleteAdvertiser(ctx context.Context, id int64) error {
	query := `DELETE FROM public.advertisers WHERE advertiser_id = $1`

	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting advertiser: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("advertiser not found: %w", domain.ErrNotFound)
	}

	return nil
}

// CreateAdvertiserExtraInfo creates extra info for an advertiser organization
func (r *pgxAdvertiserRepository) CreateAdvertiserExtraInfo(ctx context.Context, extraInfo *domain.AdvertiserExtraInfo) error {
	query := `INSERT INTO public.advertiser_extra_info (
		organization_id, website, website_type, company_size, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING extra_info_id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		extraInfo.OrganizationID,
		extraInfo.Website,
		extraInfo.WebsiteType,
		extraInfo.CompanySize,
		now,
		now,
	).Scan(&extraInfo.ExtraInfoID, &extraInfo.CreatedAt, &extraInfo.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating advertiser extra info: %w", err)
	}

	return nil
}

// GetAdvertiserExtraInfo retrieves extra info for an advertiser organization
func (r *pgxAdvertiserRepository) GetAdvertiserExtraInfo(ctx context.Context, organizationID int64) (*domain.AdvertiserExtraInfo, error) {
	query := `SELECT extra_info_id, organization_id, website, website_type, company_size, created_at, updated_at
		FROM public.advertiser_extra_info WHERE organization_id = $1`

	extraInfo := &domain.AdvertiserExtraInfo{}
	err := r.db.QueryRow(ctx, query, organizationID).Scan(
		&extraInfo.ExtraInfoID,
		&extraInfo.OrganizationID,
		&extraInfo.Website,
		&extraInfo.WebsiteType,
		&extraInfo.CompanySize,
		&extraInfo.CreatedAt,
		&extraInfo.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("error getting advertiser extra info: %w", err)
	}

	return extraInfo, nil
}

// UpdateAdvertiserExtraInfo updates extra info for an advertiser organization
func (r *pgxAdvertiserRepository) UpdateAdvertiserExtraInfo(ctx context.Context, extraInfo *domain.AdvertiserExtraInfo) error {
	query := `UPDATE public.advertiser_extra_info SET 
		website = $2, website_type = $3, company_size = $4, updated_at = $5
		WHERE organization_id = $1
		RETURNING updated_at`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		extraInfo.OrganizationID,
		extraInfo.Website,
		extraInfo.WebsiteType,
		extraInfo.CompanySize,
		now,
	).Scan(&extraInfo.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrNotFound
		}
		return fmt.Errorf("error updating advertiser extra info: %w", err)
	}

	return nil
}

// UpsertAdvertiserExtraInfo creates or updates extra info for an advertiser organization
func (r *pgxAdvertiserRepository) UpsertAdvertiserExtraInfo(ctx context.Context, extraInfo *domain.AdvertiserExtraInfo) error {
	query := `INSERT INTO public.advertiser_extra_info (
		organization_id, website, website_type, company_size, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (organization_id) DO UPDATE SET
		website = EXCLUDED.website,
		website_type = EXCLUDED.website_type,
		company_size = EXCLUDED.company_size,
		updated_at = EXCLUDED.updated_at
	RETURNING extra_info_id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		extraInfo.OrganizationID,
		extraInfo.Website,
		extraInfo.WebsiteType,
		extraInfo.CompanySize,
		now,
		now,
	).Scan(&extraInfo.ExtraInfoID, &extraInfo.CreatedAt, &extraInfo.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error upserting advertiser extra info: %w", err)
	}

	return nil
}

// DeleteAdvertiserExtraInfo deletes extra info for an advertiser organization
func (r *pgxAdvertiserRepository) DeleteAdvertiserExtraInfo(ctx context.Context, organizationID int64) error {
	query := `DELETE FROM public.advertiser_extra_info WHERE organization_id = $1`

	commandTag, err := r.db.Exec(ctx, query, organizationID)
	if err != nil {
		return fmt.Errorf("error deleting advertiser extra info: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// GetAdvertiserWithExtraInfo retrieves an advertiser with its extra info
func (r *pgxAdvertiserRepository) GetAdvertiserWithExtraInfo(ctx context.Context, advertiserID int64) (*domain.AdvertiserWithExtraInfo, error) {
	// Get the advertiser first
	advertiser, err := r.GetAdvertiserByID(ctx, advertiserID)
	if err != nil {
		return nil, err
	}

	// Get the extra info (may not exist) - use organization ID
	extraInfo, err := r.GetAdvertiserExtraInfo(ctx, advertiser.OrganizationID)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}

	result := &domain.AdvertiserWithExtraInfo{
		Advertiser: advertiser,
	}

	if err != domain.ErrNotFound {
		result.ExtraInfo = extraInfo
	}

	return result, nil
}
