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

type CampaignProviderMappingRepository interface {
	CreateCampaignProviderMapping(ctx context.Context, mapping *domain.CampaignProviderMapping) error
	GetCampaignProviderMapping(ctx context.Context, campaignID int64, providerType string) (*domain.CampaignProviderMapping, error)
	UpdateCampaignProviderMapping(ctx context.Context, mapping *domain.CampaignProviderMapping) error
	DeleteCampaignProviderMapping(ctx context.Context, mappingID int64) error
}

type pgxCampaignProviderMappingRepository struct {
	db *pgxpool.Pool
}

func NewPgxCampaignProviderMappingRepository(db *pgxpool.Pool) CampaignProviderMappingRepository {
	return &pgxCampaignProviderMappingRepository{db: db}
}

func (r *pgxCampaignProviderMappingRepository) CreateCampaignProviderMapping(ctx context.Context, mapping *domain.CampaignProviderMapping) error {
	query := `INSERT INTO public.campaign_provider_mappings 
              (campaign_id, provider_type, provider_campaign_id, api_credentials, provider_config, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)
              RETURNING mapping_id, created_at, updated_at`
	
	var providerCampaignID, apiCredentials, providerConfig sql.NullString
	
	if mapping.ProviderCampaignID != nil {
		providerCampaignID = sql.NullString{String: *mapping.ProviderCampaignID, Valid: true}
	}
	if mapping.APICredentials != nil {
		apiCredentials = sql.NullString{String: *mapping.APICredentials, Valid: true}
	}
	if mapping.ProviderConfig != nil {
		providerConfig = sql.NullString{String: *mapping.ProviderConfig, Valid: true}
	}
	
	now := time.Now()
	err := r.db.QueryRow(ctx, query, 
		mapping.CampaignID, 
		mapping.ProviderType, 
		providerCampaignID, 
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
		return fmt.Errorf("error creating campaign provider mapping: %w", err)
	}
	
	return nil
}

func (r *pgxCampaignProviderMappingRepository) GetCampaignProviderMapping(ctx context.Context, campaignID int64, providerType string) (*domain.CampaignProviderMapping, error) {
	query := `SELECT mapping_id, campaign_id, provider_type, provider_campaign_id, api_credentials, provider_config,
	          created_at, updated_at
	          FROM public.campaign_provider_mappings 
	          WHERE campaign_id = $1 AND provider_type = $2`
	
	var mapping domain.CampaignProviderMapping
	var providerCampaignID, apiCredentials, providerConfig sql.NullString
	
	err := r.db.QueryRow(ctx, query, campaignID, providerType).Scan(
		&mapping.MappingID,
		&mapping.CampaignID,
		&mapping.ProviderType,
		&providerCampaignID,
		&apiCredentials,
		&providerConfig,
		&mapping.CreatedAt,
		&mapping.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("campaign provider mapping not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("error getting campaign provider mapping: %w", err)
	}
	
	if providerCampaignID.Valid {
		mapping.ProviderCampaignID = &providerCampaignID.String
	}
	if apiCredentials.Valid {
		mapping.APICredentials = &apiCredentials.String
	}
	if providerConfig.Valid {
		mapping.ProviderConfig = &providerConfig.String
	}
	
	return &mapping, nil
}

func (r *pgxCampaignProviderMappingRepository) UpdateCampaignProviderMapping(ctx context.Context, mapping *domain.CampaignProviderMapping) error {
	query := `UPDATE public.campaign_provider_mappings SET 
	          provider_campaign_id = $1, api_credentials = $2, provider_config = $3, updated_at = $4
	          WHERE mapping_id = $5
	          RETURNING updated_at`
	
	var providerCampaignID, apiCredentials, providerConfig sql.NullString
	
	if mapping.ProviderCampaignID != nil {
		providerCampaignID = sql.NullString{String: *mapping.ProviderCampaignID, Valid: true}
	}
	if mapping.APICredentials != nil {
		apiCredentials = sql.NullString{String: *mapping.APICredentials, Valid: true}
	}
	if mapping.ProviderConfig != nil {
		providerConfig = sql.NullString{String: *mapping.ProviderConfig, Valid: true}
	}
	
	now := time.Now()
	err := r.db.QueryRow(ctx, query, 
		providerCampaignID, 
		apiCredentials, 
		providerConfig, 
		now,
		mapping.MappingID,
	).Scan(&mapping.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("campaign provider mapping not found: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("error updating campaign provider mapping: %w", err)
	}
	
	return nil
}

func (r *pgxCampaignProviderMappingRepository) DeleteCampaignProviderMapping(ctx context.Context, mappingID int64) error {
	query := `DELETE FROM public.campaign_provider_mappings WHERE mapping_id = $1`
	
	commandTag, err := r.db.Exec(ctx, query, mappingID)
	if err != nil {
		return fmt.Errorf("error deleting campaign provider mapping: %w", err)
	}
	
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("campaign provider mapping not found: %w", domain.ErrNotFound)
	}
	
	return nil
}