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
	query := `INSERT INTO public.advertisers (organization_id, name, contact_email, billing_details, status, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)
              RETURNING advertiser_id, created_at, updated_at`
	
	var billingDetailsJSON sql.NullString
	if advertiser.BillingDetails != nil {
		billingDetailsJSON = sql.NullString{String: *advertiser.BillingDetails, Valid: true}
	}
	
	now := time.Now()
	err := r.db.QueryRow(ctx, query, 
		advertiser.OrganizationID, 
		advertiser.Name, 
		advertiser.ContactEmail, 
		billingDetailsJSON, 
		advertiser.Status, 
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
	query := `SELECT advertiser_id, organization_id, name, contact_email, billing_details, status, created_at, updated_at
              FROM public.advertisers WHERE advertiser_id = $1`
	
	var advertiser domain.Advertiser
	var contactEmail, billingDetails sql.NullString
	
	err := r.db.QueryRow(ctx, query, id).Scan(
		&advertiser.AdvertiserID,
		&advertiser.OrganizationID,
		&advertiser.Name,
		&contactEmail,
		&billingDetails,
		&advertiser.Status,
		&advertiser.CreatedAt,
		&advertiser.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("advertiser not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("error getting advertiser by ID: %w", err)
	}
	
	if contactEmail.Valid {
		email := contactEmail.String
		advertiser.ContactEmail = &email
	}
	
	if billingDetails.Valid {
		details := billingDetails.String
		advertiser.BillingDetails = &details
	}
	
	return &advertiser, nil
}

// UpdateAdvertiser updates an advertiser in the database
func (r *pgxAdvertiserRepository) UpdateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error {
	query := `UPDATE public.advertisers
              SET name = $1, contact_email = $2, billing_details = $3, status = $4
              WHERE advertiser_id = $5
              RETURNING updated_at`
	
	var billingDetailsJSON sql.NullString
	if advertiser.BillingDetails != nil {
		billingDetailsJSON = sql.NullString{String: *advertiser.BillingDetails, Valid: true}
	}
	
	err := r.db.QueryRow(ctx, query, 
		advertiser.Name, 
		advertiser.ContactEmail, 
		billingDetailsJSON, 
		advertiser.Status, 
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
	query := `SELECT advertiser_id, organization_id, name, contact_email, billing_details, status, created_at, updated_at
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
		
		if err := rows.Scan(
			&advertiser.AdvertiserID,
			&advertiser.OrganizationID,
			&advertiser.Name,
			&contactEmail,
			&billingDetails,
			&advertiser.Status,
			&advertiser.CreatedAt,
			&advertiser.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning advertiser row: %w", err)
		}
		
		if contactEmail.Valid {
			email := contactEmail.String
			advertiser.ContactEmail = &email
		}
		
		if billingDetails.Valid {
			details := billingDetails.String
			advertiser.BillingDetails = &details
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