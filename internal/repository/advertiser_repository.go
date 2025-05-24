package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AdvertiserRepository defines the interface for advertiser operations
type AdvertiserRepository interface {
	CreateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error
	GetAdvertiserByID(ctx context.Context, id int64) (*domain.Advertiser, error)
	UpdateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error
	ListAdvertisersByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.Advertiser, error)
	DeleteAdvertiser(ctx context.Context, id int64) error
	
	// Provider mapping methods
	CreateAdvertiserProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error
	GetAdvertiserProviderMapping(ctx context.Context, advertiserID int64, providerType string) (*domain.AdvertiserProviderMapping, error)
	UpdateAdvertiserProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error
	DeleteAdvertiserProviderMapping(ctx context.Context, mappingID int64) error
}

// pgxAdvertiserRepository implements AdvertiserRepository using pgx
type pgxAdvertiserRepository struct {
	db *pgxpool.Pool
}

// NewPgxAdvertiserRepository creates a new advertiser repository
func NewPgxAdvertiserRepository(db *pgxpool.Pool) AdvertiserRepository {
	return &pgxAdvertiserRepository{db: db}
}

// CreateAdvertiser creates a new advertiser in the database
func (r *pgxAdvertiserRepository) CreateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error {
	query := `INSERT INTO public.advertisers (
		organization_id, name, contact_email, billing_details, status,
		internal_notes, default_currency_id, platform_name, platform_url, platform_username,
		accounting_contact_email, offer_id_macro, affiliate_id_macro, attribution_method,
		email_attribution_method, attribution_priority, reporting_timezone_id, is_expose_publisher_reporting,
		everflow_sync_status, last_everflow_sync_at, everflow_sync_error,
		created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23
	) RETURNING advertiser_id, created_at, updated_at`
	
	var billingDetailsJSON sql.NullString
	if advertiser.BillingDetails != nil {
		billingBytes, err := json.Marshal(advertiser.BillingDetails)
		if err != nil {
			return fmt.Errorf("failed to marshal billing details: %w", err)
		}
		billingDetailsJSON = sql.NullString{String: string(billingBytes), Valid: true}
	}
	
	now := time.Now()
	err := r.db.QueryRow(ctx, query, 
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
		advertiser.IsExposePublisherReporting,
		advertiser.EverflowSyncStatus,
		advertiser.LastEverflowSyncAt,
		advertiser.EverflowSyncError,
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

// GetAdvertiserByID retrieves an advertiser by ID
func (r *pgxAdvertiserRepository) GetAdvertiserByID(ctx context.Context, id int64) (*domain.Advertiser, error) {
	query := `SELECT 
		advertiser_id, organization_id, name, contact_email, billing_details, status,
		internal_notes, default_currency_id, platform_name, platform_url, platform_username,
		accounting_contact_email, offer_id_macro, affiliate_id_macro, attribution_method,
		email_attribution_method, attribution_priority, reporting_timezone_id, is_expose_publisher_reporting,
		everflow_sync_status, last_everflow_sync_at, everflow_sync_error,
		created_at, updated_at
	FROM public.advertisers WHERE advertiser_id = $1`
	
	var advertiser domain.Advertiser
	var contactEmail, billingDetails sql.NullString
	var internalNotes, defaultCurrencyID, platformName, platformURL, platformUsername sql.NullString
	var accountingContactEmail, offerIDMacro, affiliateIDMacro sql.NullString
	var attributionMethod, emailAttributionMethod, attributionPriority sql.NullString
	var reportingTimezoneID sql.NullInt32
	var isExposePublisherReporting sql.NullBool
	var everflowSyncStatus, everflowSyncError sql.NullString
	var lastEverflowSyncAt sql.NullTime
	
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
		&isExposePublisherReporting,
		&everflowSyncStatus,
		&lastEverflowSyncAt,
		&everflowSyncError,
		&advertiser.CreatedAt,
		&advertiser.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("advertiser not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("error getting advertiser by ID: %w", err)
	}
	
	// Handle nullable fields
	if contactEmail.Valid {
		advertiser.ContactEmail = &contactEmail.String
	}
	if billingDetails.Valid {
		var billing domain.BillingDetails
		if err := json.Unmarshal([]byte(billingDetails.String), &billing); err != nil {
			return nil, fmt.Errorf("failed to unmarshal billing details: %w", err)
		}
		advertiser.BillingDetails = &billing
	}
	if internalNotes.Valid {
		advertiser.InternalNotes = &internalNotes.String
	}
	if defaultCurrencyID.Valid {
		advertiser.DefaultCurrencyID = &defaultCurrencyID.String
	}
	if platformName.Valid {
		advertiser.PlatformName = &platformName.String
	}
	if platformURL.Valid {
		advertiser.PlatformURL = &platformURL.String
	}
	if platformUsername.Valid {
		advertiser.PlatformUsername = &platformUsername.String
	}
	if accountingContactEmail.Valid {
		advertiser.AccountingContactEmail = &accountingContactEmail.String
	}
	if offerIDMacro.Valid {
		advertiser.OfferIDMacro = &offerIDMacro.String
	}
	if affiliateIDMacro.Valid {
		advertiser.AffiliateIDMacro = &affiliateIDMacro.String
	}
	if attributionMethod.Valid {
		advertiser.AttributionMethod = &attributionMethod.String
	}
	if emailAttributionMethod.Valid {
		advertiser.EmailAttributionMethod = &emailAttributionMethod.String
	}
	if attributionPriority.Valid {
		advertiser.AttributionPriority = &attributionPriority.String
	}
	if reportingTimezoneID.Valid {
		timezoneID := int(reportingTimezoneID.Int32)
		advertiser.ReportingTimezoneID = &timezoneID
	}
	if isExposePublisherReporting.Valid {
		advertiser.IsExposePublisherReporting = &isExposePublisherReporting.Bool
	}
	if everflowSyncStatus.Valid {
		advertiser.EverflowSyncStatus = &everflowSyncStatus.String
	}
	if lastEverflowSyncAt.Valid {
		advertiser.LastEverflowSyncAt = &lastEverflowSyncAt.Time
	}
	if everflowSyncError.Valid {
		advertiser.EverflowSyncError = &everflowSyncError.String
	}
	
	return &advertiser, nil
}

// UpdateAdvertiser updates an advertiser in the database
func (r *pgxAdvertiserRepository) UpdateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error {
	query := `UPDATE public.advertisers SET 
		name = $1, contact_email = $2, billing_details = $3, status = $4,
		internal_notes = $5, default_currency_id = $6, platform_name = $7, platform_url = $8, platform_username = $9,
		accounting_contact_email = $10, offer_id_macro = $11, affiliate_id_macro = $12, attribution_method = $13,
		email_attribution_method = $14, attribution_priority = $15, reporting_timezone_id = $16, is_expose_publisher_reporting = $17,
		everflow_sync_status = $18, last_everflow_sync_at = $19, everflow_sync_error = $20,
		updated_at = $21
	WHERE advertiser_id = $22
	RETURNING updated_at`
	
	var billingDetailsJSON sql.NullString
	if advertiser.BillingDetails != nil {
		billingBytes, err := json.Marshal(advertiser.BillingDetails)
		if err != nil {
			return fmt.Errorf("failed to marshal billing details: %w", err)
		}
		billingDetailsJSON = sql.NullString{String: string(billingBytes), Valid: true}
	}
	
	now := time.Now()
	err := r.db.QueryRow(ctx, query, 
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
		advertiser.IsExposePublisherReporting,
		advertiser.EverflowSyncStatus,
		advertiser.LastEverflowSyncAt,
		advertiser.EverflowSyncError,
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

// ListAdvertisersByOrganization retrieves a list of advertisers for an organization with pagination
func (r *pgxAdvertiserRepository) ListAdvertisersByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.Advertiser, error) {
	query := `SELECT 
		advertiser_id, organization_id, name, contact_email, billing_details, status,
		internal_notes, default_currency_id, platform_name, platform_url, platform_username,
		accounting_contact_email, offer_id_macro, affiliate_id_macro, attribution_method,
		email_attribution_method, attribution_priority, reporting_timezone_id, is_expose_publisher_reporting,
		everflow_sync_status, last_everflow_sync_at, everflow_sync_error,
		created_at, updated_at
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
		var isExposePublisherReporting sql.NullBool
		var everflowSyncStatus, everflowSyncError sql.NullString
		var lastEverflowSyncAt sql.NullTime
		
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
			&isExposePublisherReporting,
			&everflowSyncStatus,
			&lastEverflowSyncAt,
			&everflowSyncError,
			&advertiser.CreatedAt,
			&advertiser.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning advertiser row: %w", err)
		}
		
		// Handle nullable fields
		if contactEmail.Valid {
			advertiser.ContactEmail = &contactEmail.String
		}
		if billingDetails.Valid {
			var billing domain.BillingDetails
			if err := json.Unmarshal([]byte(billingDetails.String), &billing); err != nil {
				return nil, fmt.Errorf("failed to unmarshal billing details: %w", err)
			}
			advertiser.BillingDetails = &billing
		}
		if internalNotes.Valid {
			advertiser.InternalNotes = &internalNotes.String
		}
		if defaultCurrencyID.Valid {
			advertiser.DefaultCurrencyID = &defaultCurrencyID.String
		}
		if platformName.Valid {
			advertiser.PlatformName = &platformName.String
		}
		if platformURL.Valid {
			advertiser.PlatformURL = &platformURL.String
		}
		if platformUsername.Valid {
			advertiser.PlatformUsername = &platformUsername.String
		}
		if accountingContactEmail.Valid {
			advertiser.AccountingContactEmail = &accountingContactEmail.String
		}
		if offerIDMacro.Valid {
			advertiser.OfferIDMacro = &offerIDMacro.String
		}
		if affiliateIDMacro.Valid {
			advertiser.AffiliateIDMacro = &affiliateIDMacro.String
		}
		if attributionMethod.Valid {
			advertiser.AttributionMethod = &attributionMethod.String
		}
		if emailAttributionMethod.Valid {
			advertiser.EmailAttributionMethod = &emailAttributionMethod.String
		}
		if attributionPriority.Valid {
			advertiser.AttributionPriority = &attributionPriority.String
		}
		if reportingTimezoneID.Valid {
			timezoneID := int(reportingTimezoneID.Int32)
			advertiser.ReportingTimezoneID = &timezoneID
		}
		if isExposePublisherReporting.Valid {
			advertiser.IsExposePublisherReporting = &isExposePublisherReporting.Bool
		}
		if everflowSyncStatus.Valid {
			advertiser.EverflowSyncStatus = &everflowSyncStatus.String
		}
		if lastEverflowSyncAt.Valid {
			advertiser.LastEverflowSyncAt = &lastEverflowSyncAt.Time
		}
		if everflowSyncError.Valid {
			advertiser.EverflowSyncError = &everflowSyncError.String
		}
		
		advertisers = append(advertisers, &advertiser)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating advertiser rows: %w", err)
	}
	
	return advertisers, nil
}

// DeleteAdvertiser deletes an advertiser from the database
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

// CreateAdvertiserProviderMapping creates a new advertiser provider mapping in the database
func (r *pgxAdvertiserRepository) CreateAdvertiserProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error {
	query := `INSERT INTO public.advertiser_provider_mappings 
              (advertiser_id, provider_type, provider_advertiser_id, api_credentials, provider_config, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)
              RETURNING mapping_id, created_at, updated_at`
	
	var providerAdvertiserID, apiCredentials, providerConfig sql.NullString
	
	if mapping.ProviderAdvertiserID != nil {
		providerAdvertiserID = sql.NullString{String: *mapping.ProviderAdvertiserID, Valid: true}
	}
	
	if mapping.APICredentials != nil {
		apiCredentials = sql.NullString{String: *mapping.APICredentials, Valid: true}
	}
	
	if mapping.ProviderConfig != nil {
		providerConfig = sql.NullString{String: *mapping.ProviderConfig, Valid: true}
	}
	
	now := time.Now()
	err := r.db.QueryRow(ctx, query, 
		mapping.AdvertiserID, 
		mapping.ProviderType, 
		providerAdvertiserID, 
		apiCredentials, 
		providerConfig, 
		now, 
		now,
	).Scan(
		&mapping.MappingID, 
		&mapping.CreatedAt, 
		&mapping.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("error creating advertiser provider mapping: %w", err)
	}
	
	return nil
}

// GetAdvertiserProviderMapping retrieves an advertiser provider mapping
func (r *pgxAdvertiserRepository) GetAdvertiserProviderMapping(ctx context.Context, advertiserID int64, providerType string) (*domain.AdvertiserProviderMapping, error) {
	query := `SELECT mapping_id, advertiser_id, provider_type, provider_advertiser_id, api_credentials, provider_config, created_at, updated_at
              FROM public.advertiser_provider_mappings 
              WHERE advertiser_id = $1 AND provider_type = $2`
	
	var mapping domain.AdvertiserProviderMapping
	var providerAdvertiserID, apiCredentials, providerConfig sql.NullString
	
	err := r.db.QueryRow(ctx, query, advertiserID, providerType).Scan(
		&mapping.MappingID,
		&mapping.AdvertiserID,
		&mapping.ProviderType,
		&providerAdvertiserID,
		&apiCredentials,
		&providerConfig,
		&mapping.CreatedAt,
		&mapping.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("advertiser provider mapping not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("error getting advertiser provider mapping: %w", err)
	}
	
	if providerAdvertiserID.Valid {
		id := providerAdvertiserID.String
		mapping.ProviderAdvertiserID = &id
	}
	
	if apiCredentials.Valid {
		creds := apiCredentials.String
		mapping.APICredentials = &creds
	}
	
	if providerConfig.Valid {
		config := providerConfig.String
		mapping.ProviderConfig = &config
	}
	
	return &mapping, nil
}

// UpdateAdvertiserProviderMapping updates an advertiser provider mapping in the database
func (r *pgxAdvertiserRepository) UpdateAdvertiserProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error {
	query := `UPDATE public.advertiser_provider_mappings
              SET provider_advertiser_id = $1, api_credentials = $2, provider_config = $3
              WHERE mapping_id = $4
              RETURNING updated_at`
	
	var providerAdvertiserID, apiCredentials, providerConfig sql.NullString
	
	if mapping.ProviderAdvertiserID != nil {
		providerAdvertiserID = sql.NullString{String: *mapping.ProviderAdvertiserID, Valid: true}
	}
	
	if mapping.APICredentials != nil {
		apiCredentials = sql.NullString{String: *mapping.APICredentials, Valid: true}
	}
	
	if mapping.ProviderConfig != nil {
		providerConfig = sql.NullString{String: *mapping.ProviderConfig, Valid: true}
	}
	
	err := r.db.QueryRow(ctx, query, 
		providerAdvertiserID, 
		apiCredentials, 
		providerConfig, 
		mapping.MappingID,
	).Scan(&mapping.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("advertiser provider mapping not found: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("error updating advertiser provider mapping: %w", err)
	}
	
	return nil
}

// DeleteAdvertiserProviderMapping deletes an advertiser provider mapping from the database
func (r *pgxAdvertiserRepository) DeleteAdvertiserProviderMapping(ctx context.Context, mappingID int64) error {
	query := `DELETE FROM public.advertiser_provider_mappings WHERE mapping_id = $1`
	
	commandTag, err := r.db.Exec(ctx, query, mappingID)
	if err != nil {
		return fmt.Errorf("error deleting advertiser provider mapping: %w", err)
	}
	
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("advertiser provider mapping not found: %w", domain.ErrNotFound)
	}
	
	return nil
}