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

type AffiliateProviderMappingRepository interface {
	CreateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error
	GetAffiliateProviderMapping(ctx context.Context, affiliateID int64, providerType string) (*domain.AffiliateProviderMapping, error)
	UpdateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error
	DeleteAffiliateProviderMapping(ctx context.Context, mappingID int64) error
}

type pgxAffiliateProviderMappingRepository struct {
	db *pgxpool.Pool
}

func NewPgxAffiliateProviderMappingRepository(db *pgxpool.Pool) AffiliateProviderMappingRepository {
	return &pgxAffiliateProviderMappingRepository{db: db}
}

func (r *pgxAffiliateProviderMappingRepository) CreateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error {
	query := `INSERT INTO public.affiliate_provider_mappings 
              (affiliate_id, provider_type, provider_affiliate_id, api_credentials, provider_config, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)
              RETURNING mapping_id, created_at, updated_at`
	
	var providerAffiliateID, apiCredentials, providerConfig sql.NullString
	
	if mapping.ProviderAffiliateID != nil {
		providerAffiliateID = sql.NullString{String: *mapping.ProviderAffiliateID, Valid: true}
	}
	if mapping.APICredentials != nil {
		apiCredentials = sql.NullString{String: *mapping.APICredentials, Valid: true}
	}
	if mapping.ProviderConfig != nil {
		providerConfig = sql.NullString{String: *mapping.ProviderConfig, Valid: true}
	}
	
	now := time.Now()
	err := r.db.QueryRow(ctx, query, 
		mapping.AffiliateID, 
		mapping.ProviderType, 
		providerAffiliateID, 
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
		return fmt.Errorf("error creating affiliate provider mapping: %w", err)
	}
	
	return nil
}

func (r *pgxAffiliateProviderMappingRepository) GetAffiliateProviderMapping(ctx context.Context, affiliateID int64, providerType string) (*domain.AffiliateProviderMapping, error) {
	query := `SELECT mapping_id, affiliate_id, provider_type, provider_affiliate_id, api_credentials, provider_config,
	          created_at, updated_at
	          FROM public.affiliate_provider_mappings 
	          WHERE affiliate_id = $1 AND provider_type = $2`
	
	var mapping domain.AffiliateProviderMapping
	var providerAffiliateID, apiCredentials, providerConfig sql.NullString
	
	err := r.db.QueryRow(ctx, query, affiliateID, providerType).Scan(
		&mapping.MappingID,
		&mapping.AffiliateID,
		&mapping.ProviderType,
		&providerAffiliateID,
		&apiCredentials,
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
		mapping.ProviderAffiliateID = &providerAffiliateID.String
	}
	if apiCredentials.Valid {
		mapping.APICredentials = &apiCredentials.String
	}
	if providerConfig.Valid {
		mapping.ProviderConfig = &providerConfig.String
	}
	
	return &mapping, nil
}

func (r *pgxAffiliateProviderMappingRepository) UpdateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error {
	query := `UPDATE public.affiliate_provider_mappings SET 
	          provider_affiliate_id = $1, api_credentials = $2, provider_config = $3, updated_at = $4
	          WHERE mapping_id = $5
	          RETURNING updated_at`
	
	var providerAffiliateID, apiCredentials, providerConfig sql.NullString
	
	if mapping.ProviderAffiliateID != nil {
		providerAffiliateID = sql.NullString{String: *mapping.ProviderAffiliateID, Valid: true}
	}
	if mapping.APICredentials != nil {
		apiCredentials = sql.NullString{String: *mapping.APICredentials, Valid: true}
	}
	if mapping.ProviderConfig != nil {
		providerConfig = sql.NullString{String: *mapping.ProviderConfig, Valid: true}
	}
	
	now := time.Now()
	err := r.db.QueryRow(ctx, query, 
		providerAffiliateID, 
		apiCredentials, 
		providerConfig, 
		now,
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

func (r *pgxAffiliateProviderMappingRepository) DeleteAffiliateProviderMapping(ctx context.Context, mappingID int64) error {
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