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

type AdvertiserProviderMappingRepository interface {
	CreateAdvertiserProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error
	GetAdvertiserProviderMapping(ctx context.Context, advertiserID int64, providerType string) (*domain.AdvertiserProviderMapping, error)
	UpdateAdvertiserProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error
	DeleteAdvertiserProviderMapping(ctx context.Context, mappingID int64) error
}

type pgxAdvertiserProviderMappingRepository struct {
	db *pgxpool.Pool
}

func NewPgxAdvertiserProviderMappingRepository(db *pgxpool.Pool) AdvertiserProviderMappingRepository {
	return &pgxAdvertiserProviderMappingRepository{db: db}
}

func (r *pgxAdvertiserProviderMappingRepository) CreateAdvertiserProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error {
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

func (r *pgxAdvertiserProviderMappingRepository) GetAdvertiserProviderMapping(ctx context.Context, advertiserID int64, providerType string) (*domain.AdvertiserProviderMapping, error) {
	query := `SELECT mapping_id, advertiser_id, provider_type, provider_advertiser_id, api_credentials, provider_config,
	          created_at, updated_at
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
		mapping.ProviderAdvertiserID = &providerAdvertiserID.String
	}
	if apiCredentials.Valid {
		mapping.APICredentials = &apiCredentials.String
	}
	if providerConfig.Valid {
		mapping.ProviderConfig = &providerConfig.String
	}
	
	return &mapping, nil
}

func (r *pgxAdvertiserProviderMappingRepository) UpdateAdvertiserProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error {
	query := `UPDATE public.advertiser_provider_mappings SET 
	          provider_advertiser_id = $1, api_credentials = $2, provider_config = $3, updated_at = $4
	          WHERE mapping_id = $5
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
	
	now := time.Now()
	err := r.db.QueryRow(ctx, query, 
		providerAdvertiserID, 
		apiCredentials, 
		providerConfig, 
		now,
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

func (r *pgxAdvertiserProviderMappingRepository) DeleteAdvertiserProviderMapping(ctx context.Context, mappingID int64) error {
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