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

// AffiliateRepository defines the interface for affiliate operations
type AffiliateRepository interface {
	CreateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error
	GetAffiliateByID(ctx context.Context, id int64) (*domain.Affiliate, error)
	UpdateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error
	ListAffiliatesByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.Affiliate, error)
	DeleteAffiliate(ctx context.Context, id int64) error
	
	// Provider mapping methods
	CreateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error
	GetAffiliateProviderMapping(ctx context.Context, affiliateID int64, providerType string) (*domain.AffiliateProviderMapping, error)
	UpdateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error
	DeleteAffiliateProviderMapping(ctx context.Context, mappingID int64) error
}

// pgxAffiliateRepository implements AffiliateRepository using pgx
type pgxAffiliateRepository struct {
	db *pgxpool.Pool
}

// NewPgxAffiliateRepository creates a new affiliate repository
func NewPgxAffiliateRepository(db *pgxpool.Pool) AffiliateRepository {
	return &pgxAffiliateRepository{db: db}
}

// CreateAffiliate creates a new affiliate in the database
func (r *pgxAffiliateRepository) CreateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error {
	query := `INSERT INTO public.affiliates (organization_id, name, contact_email, payment_details, status, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)
              RETURNING affiliate_id, created_at, updated_at`
	
	var contactEmail, paymentDetails sql.NullString
	
	if affiliate.ContactEmail != nil {
		contactEmail = sql.NullString{String: *affiliate.ContactEmail, Valid: true}
	}
	
	if affiliate.PaymentDetails != nil {
		paymentDetails = sql.NullString{String: *affiliate.PaymentDetails, Valid: true}
	}
	
	now := time.Now()
	err := r.db.QueryRow(ctx, query, 
		affiliate.OrganizationID, 
		affiliate.Name, 
		contactEmail, 
		paymentDetails, 
		affiliate.Status, 
		now, 
		now,
	).Scan(
		&affiliate.AffiliateID, 
		&affiliate.CreatedAt, 
		&affiliate.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("error creating affiliate: %w", err)
	}
	return nil
}

// GetAffiliateByID retrieves an affiliate by ID
func (r *pgxAffiliateRepository) GetAffiliateByID(ctx context.Context, id int64) (*domain.Affiliate, error) {
	query := `SELECT affiliate_id, organization_id, name, contact_email, payment_details, status, created_at, updated_at
              FROM public.affiliates WHERE affiliate_id = $1`
	
	var affiliate domain.Affiliate
	var contactEmail, paymentDetails sql.NullString
	
	err := r.db.QueryRow(ctx, query, id).Scan(
		&affiliate.AffiliateID,
		&affiliate.OrganizationID,
		&affiliate.Name,
		&contactEmail,
		&paymentDetails,
		&affiliate.Status,
		&affiliate.CreatedAt,
		&affiliate.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("affiliate not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("error getting affiliate by ID: %w", err)
	}
	
	if contactEmail.Valid {
		email := contactEmail.String
		affiliate.ContactEmail = &email
	}
	
	if paymentDetails.Valid {
		details := paymentDetails.String
		affiliate.PaymentDetails = &details
	}
	
	return &affiliate, nil
}

// UpdateAffiliate updates an affiliate in the database
func (r *pgxAffiliateRepository) UpdateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error {
	query := `UPDATE public.affiliates
              SET name = $1, contact_email = $2, payment_details = $3, status = $4
              WHERE affiliate_id = $5
              RETURNING updated_at`
	
	var contactEmail, paymentDetails sql.NullString
	
	if affiliate.ContactEmail != nil {
		contactEmail = sql.NullString{String: *affiliate.ContactEmail, Valid: true}
	}
	
	if affiliate.PaymentDetails != nil {
		paymentDetails = sql.NullString{String: *affiliate.PaymentDetails, Valid: true}
	}
	
	err := r.db.QueryRow(ctx, query, 
		affiliate.Name, 
		contactEmail, 
		paymentDetails, 
		affiliate.Status, 
		affiliate.AffiliateID,
	).Scan(&affiliate.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("affiliate not found: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("error updating affiliate: %w", err)
	}
	
	return nil
}

// ListAffiliatesByOrganization retrieves a list of affiliates for an organization with pagination
func (r *pgxAffiliateRepository) ListAffiliatesByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.Affiliate, error) {
	query := `SELECT affiliate_id, organization_id, name, contact_email, payment_details, status, created_at, updated_at
              FROM public.affiliates
              WHERE organization_id = $1
              ORDER BY affiliate_id
              LIMIT $2 OFFSET $3`
	
	rows, err := r.db.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing affiliates: %w", err)
	}
	defer rows.Close()
	
	var affiliates []*domain.Affiliate
	for rows.Next() {
		var affiliate domain.Affiliate
		var contactEmail, paymentDetails sql.NullString
		
		if err := rows.Scan(
			&affiliate.AffiliateID,
			&affiliate.OrganizationID,
			&affiliate.Name,
			&contactEmail,
			&paymentDetails,
			&affiliate.Status,
			&affiliate.CreatedAt,
			&affiliate.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning affiliate row: %w", err)
		}
		
		if contactEmail.Valid {
			email := contactEmail.String
			affiliate.ContactEmail = &email
		}
		
		if paymentDetails.Valid {
			details := paymentDetails.String
			affiliate.PaymentDetails = &details
		}
		
		affiliates = append(affiliates, &affiliate)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating affiliate rows: %w", err)
	}
	
	return affiliates, nil
}

// DeleteAffiliate deletes an affiliate from the database
func (r *pgxAffiliateRepository) DeleteAffiliate(ctx context.Context, id int64) error {
	query := `DELETE FROM public.affiliates WHERE affiliate_id = $1`
	
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting affiliate: %w", err)
	}
	
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("affiliate not found: %w", domain.ErrNotFound)
	}
	
	return nil
}

// CreateAffiliateProviderMapping creates a new affiliate provider mapping in the database
func (r *pgxAffiliateRepository) CreateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error {
	query := `INSERT INTO public.affiliate_provider_mappings 
              (affiliate_id, provider_type, provider_affiliate_id, provider_config, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6)
              RETURNING mapping_id, created_at, updated_at`
	
	var providerAffiliateID, providerConfig sql.NullString
	
	if mapping.ProviderAffiliateID != nil {
		providerAffiliateID = sql.NullString{String: *mapping.ProviderAffiliateID, Valid: true}
	}
	
	if mapping.ProviderConfig != nil {
		providerConfig = sql.NullString{String: *mapping.ProviderConfig, Valid: true}
	}
	
	now := time.Now()
	err := r.db.QueryRow(ctx, query, 
		mapping.AffiliateID, 
		mapping.ProviderType, 
		providerAffiliateID, 
		providerConfig, 
		now, 
		now,
	).Scan(
		&mapping.MappingID, 
		&mapping.CreatedAt, 
		&mapping.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("error creating affiliate provider mapping: %w", err)
	}
	
	return nil
}

// GetAffiliateProviderMapping retrieves an affiliate provider mapping
func (r *pgxAffiliateRepository) GetAffiliateProviderMapping(ctx context.Context, affiliateID int64, providerType string) (*domain.AffiliateProviderMapping, error) {
	query := `SELECT mapping_id, affiliate_id, provider_type, provider_affiliate_id, provider_config, created_at, updated_at
              FROM public.affiliate_provider_mappings 
              WHERE affiliate_id = $1 AND provider_type = $2`
	
	var mapping domain.AffiliateProviderMapping
	var providerAffiliateID, providerConfig sql.NullString
	
	err := r.db.QueryRow(ctx, query, affiliateID, providerType).Scan(
		&mapping.MappingID,
		&mapping.AffiliateID,
		&mapping.ProviderType,
		&providerAffiliateID,
		&providerConfig,
		&mapping.CreatedAt,
		&mapping.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("affiliate provider mapping not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("error getting affiliate provider mapping: %w", err)
	}
	
	if providerAffiliateID.Valid {
		id := providerAffiliateID.String
		mapping.ProviderAffiliateID = &id
	}
	
	if providerConfig.Valid {
		config := providerConfig.String
		mapping.ProviderConfig = &config
	}
	
	return &mapping, nil
}

// UpdateAffiliateProviderMapping updates an affiliate provider mapping in the database
func (r *pgxAffiliateRepository) UpdateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error {
	query := `UPDATE public.affiliate_provider_mappings
              SET provider_affiliate_id = $1, provider_config = $2
              WHERE mapping_id = $3
              RETURNING updated_at`
	
	var providerAffiliateID, providerConfig sql.NullString
	
	if mapping.ProviderAffiliateID != nil {
		providerAffiliateID = sql.NullString{String: *mapping.ProviderAffiliateID, Valid: true}
	}
	
	if mapping.ProviderConfig != nil {
		providerConfig = sql.NullString{String: *mapping.ProviderConfig, Valid: true}
	}
	
	err := r.db.QueryRow(ctx, query, 
		providerAffiliateID, 
		providerConfig, 
		mapping.MappingID,
	).Scan(&mapping.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("affiliate provider mapping not found: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("error updating affiliate provider mapping: %w", err)
	}
	
	return nil
}

// DeleteAffiliateProviderMapping deletes an affiliate provider mapping from the database
func (r *pgxAffiliateRepository) DeleteAffiliateProviderMapping(ctx context.Context, mappingID int64) error {
	query := `DELETE FROM public.affiliate_provider_mappings WHERE mapping_id = $1`
	
	commandTag, err := r.db.Exec(ctx, query, mappingID)
	if err != nil {
		return fmt.Errorf("error deleting affiliate provider mapping: %w", err)
	}
	
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("affiliate provider mapping not found: %w", domain.ErrNotFound)
	}
	
	return nil
}