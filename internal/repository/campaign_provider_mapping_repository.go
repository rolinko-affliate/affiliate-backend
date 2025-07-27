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

// CampaignProviderMappingRepository defines the interface for campaign provider mapping operations
type CampaignProviderMappingRepository interface {
	CreateCampaignProviderMapping(ctx context.Context, mapping *domain.CampaignProviderMapping) error
	GetCampaignProviderMapping(ctx context.Context, campaignID int64, providerType string) (*domain.CampaignProviderMapping, error)
	GetCampaignProviderMappingByID(ctx context.Context, mappingID int64) (*domain.CampaignProviderMapping, error)
	UpdateCampaignProviderMapping(ctx context.Context, mapping *domain.CampaignProviderMapping) error
	ListCampaignProviderMappingsByCampaign(ctx context.Context, campaignID int64) ([]*domain.CampaignProviderMapping, error)
	DeleteCampaignProviderMapping(ctx context.Context, mappingID int64) error
}

// pgxCampaignProviderMappingRepository implements CampaignProviderMappingRepository using pgx
type pgxCampaignProviderMappingRepository struct {
	db *pgxpool.Pool
}

// NewPgxCampaignProviderMappingRepository creates a new campaign provider mapping repository
func NewPgxCampaignProviderMappingRepository(db *pgxpool.Pool) CampaignProviderMappingRepository {
	return &pgxCampaignProviderMappingRepository{db: db}
}

// CreateCampaignProviderMapping creates a new campaign provider mapping
func (r *pgxCampaignProviderMappingRepository) CreateCampaignProviderMapping(ctx context.Context, mapping *domain.CampaignProviderMapping) error {
	query := `INSERT INTO public.campaign_provider_mappings 
              (campaign_id, provider_type, provider_campaign_id, provider_offer_id, provider_config, is_active_on_provider, last_synced_at, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
              RETURNING mapping_id, created_at, updated_at`

	var providerCampaignID sql.NullString
	var providerOfferID sql.NullString
	var providerData sql.NullString
	var isActiveOnProvider sql.NullBool
	var lastSyncedAt sql.NullTime

	if mapping.ProviderCampaignID != nil {
		providerCampaignID = sql.NullString{String: *mapping.ProviderCampaignID, Valid: true}
	}
	if mapping.ProviderOfferID != nil {
		providerOfferID = sql.NullString{String: *mapping.ProviderOfferID, Valid: true}
	}
	if mapping.ProviderData != nil {
		providerData = sql.NullString{String: *mapping.ProviderData, Valid: true}
	}
	if mapping.IsActiveOnProvider != nil {
		isActiveOnProvider = sql.NullBool{Bool: *mapping.IsActiveOnProvider, Valid: true}
	}
	if mapping.LastSyncedAt != nil {
		lastSyncedAt = sql.NullTime{Time: *mapping.LastSyncedAt, Valid: true}
	}

	now := time.Now()

	err := r.db.QueryRow(ctx, query,
		mapping.CampaignID, mapping.ProviderType, providerCampaignID, providerOfferID, providerData,
		isActiveOnProvider, lastSyncedAt, now, now,
	).Scan(&mapping.MappingID, &mapping.CreatedAt, &mapping.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create campaign provider mapping: %w", err)
	}

	return nil
}

// GetCampaignProviderMapping retrieves a campaign provider mapping by campaign ID and provider type
func (r *pgxCampaignProviderMappingRepository) GetCampaignProviderMapping(ctx context.Context, campaignID int64, providerType string) (*domain.CampaignProviderMapping, error) {
	query := `SELECT mapping_id, campaign_id, provider_type, provider_campaign_id, provider_offer_id, provider_config, 
              is_active_on_provider, last_synced_at, created_at, updated_at
              FROM public.campaign_provider_mappings 
              WHERE campaign_id = $1 AND provider_type = $2`

	mapping := &domain.CampaignProviderMapping{}
	var providerCampaignID sql.NullString
	var providerOfferID sql.NullString
	var providerData sql.NullString
	var isActiveOnProvider sql.NullBool
	var lastSyncedAt sql.NullTime

	err := r.db.QueryRow(ctx, query, campaignID, providerType).Scan(
		&mapping.MappingID, &mapping.CampaignID, &mapping.ProviderType,
		&providerCampaignID, &providerOfferID, &providerData, &isActiveOnProvider, &lastSyncedAt,
		&mapping.CreatedAt, &mapping.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("campaign provider mapping not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get campaign provider mapping: %w", err)
	}

	// Handle nullable fields
	if providerCampaignID.Valid {
		mapping.ProviderCampaignID = &providerCampaignID.String
	}
	if providerOfferID.Valid {
		mapping.ProviderOfferID = &providerOfferID.String
	}
	if providerData.Valid {
		mapping.ProviderData = &providerData.String
	}
	if isActiveOnProvider.Valid {
		mapping.IsActiveOnProvider = &isActiveOnProvider.Bool
	}
	if lastSyncedAt.Valid {
		mapping.LastSyncedAt = &lastSyncedAt.Time
	}

	return mapping, nil
}

// GetCampaignProviderMappingByID retrieves a campaign provider mapping by its ID
func (r *pgxCampaignProviderMappingRepository) GetCampaignProviderMappingByID(ctx context.Context, mappingID int64) (*domain.CampaignProviderMapping, error) {
	query := `SELECT mapping_id, campaign_id, provider_type, provider_campaign_id, provider_offer_id, provider_config, 
              is_active_on_provider, last_synced_at, created_at, updated_at
              FROM public.campaign_provider_mappings 
              WHERE mapping_id = $1`

	mapping := &domain.CampaignProviderMapping{}
	var providerCampaignID sql.NullString
	var providerOfferID sql.NullString
	var providerData sql.NullString
	var isActiveOnProvider sql.NullBool
	var lastSyncedAt sql.NullTime

	err := r.db.QueryRow(ctx, query, mappingID).Scan(
		&mapping.MappingID, &mapping.CampaignID, &mapping.ProviderType,
		&providerCampaignID, &providerOfferID, &providerData, &isActiveOnProvider, &lastSyncedAt,
		&mapping.CreatedAt, &mapping.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("campaign provider mapping not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get campaign provider mapping: %w", err)
	}

	// Handle nullable fields
	if providerCampaignID.Valid {
		mapping.ProviderCampaignID = &providerCampaignID.String
	}
	if providerOfferID.Valid {
		mapping.ProviderOfferID = &providerOfferID.String
	}
	if providerData.Valid {
		mapping.ProviderData = &providerData.String
	}
	if isActiveOnProvider.Valid {
		mapping.IsActiveOnProvider = &isActiveOnProvider.Bool
	}
	if lastSyncedAt.Valid {
		mapping.LastSyncedAt = &lastSyncedAt.Time
	}

	return mapping, nil
}

// UpdateCampaignProviderMapping updates an existing campaign provider mapping
func (r *pgxCampaignProviderMappingRepository) UpdateCampaignProviderMapping(ctx context.Context, mapping *domain.CampaignProviderMapping) error {
	query := `UPDATE public.campaign_provider_mappings SET 
              campaign_id = $2, provider_type = $3, provider_campaign_id = $4, provider_offer_id = $5, provider_config = $6,
              is_active_on_provider = $7, last_synced_at = $8, updated_at = $9
              WHERE mapping_id = $1`

	var providerCampaignID sql.NullString
	var providerOfferID sql.NullString
	var providerData sql.NullString
	var isActiveOnProvider sql.NullBool
	var lastSyncedAt sql.NullTime

	if mapping.ProviderCampaignID != nil {
		providerCampaignID = sql.NullString{String: *mapping.ProviderCampaignID, Valid: true}
	}
	if mapping.ProviderOfferID != nil {
		providerOfferID = sql.NullString{String: *mapping.ProviderOfferID, Valid: true}
	}
	if mapping.ProviderData != nil {
		providerData = sql.NullString{String: *mapping.ProviderData, Valid: true}
	}
	if mapping.IsActiveOnProvider != nil {
		isActiveOnProvider = sql.NullBool{Bool: *mapping.IsActiveOnProvider, Valid: true}
	}
	if mapping.LastSyncedAt != nil {
		lastSyncedAt = sql.NullTime{Time: *mapping.LastSyncedAt, Valid: true}
	}

	now := time.Now()

	result, err := r.db.Exec(ctx, query,
		mapping.MappingID, mapping.CampaignID, mapping.ProviderType, providerCampaignID, providerOfferID, providerData,
		isActiveOnProvider, lastSyncedAt, now,
	)

	if err != nil {
		return fmt.Errorf("failed to update campaign provider mapping: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("campaign provider mapping not found: not found")
	}

	mapping.UpdatedAt = now
	return nil
}

// ListCampaignProviderMappingsByCampaign retrieves all provider mappings for a specific campaign
func (r *pgxCampaignProviderMappingRepository) ListCampaignProviderMappingsByCampaign(ctx context.Context, campaignID int64) ([]*domain.CampaignProviderMapping, error) {
	query := `SELECT mapping_id, campaign_id, provider_type, provider_campaign_id, provider_offer_id, provider_config, 
              is_active_on_provider, last_synced_at, created_at, updated_at
              FROM public.campaign_provider_mappings 
              WHERE campaign_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to list campaign provider mappings: %w", err)
	}
	defer rows.Close()

	var mappings []*domain.CampaignProviderMapping
	for rows.Next() {
		mapping := &domain.CampaignProviderMapping{}
		var providerCampaignID sql.NullString
		var providerOfferID sql.NullString
		var providerData sql.NullString
		var isActiveOnProvider sql.NullBool
		var lastSyncedAt sql.NullTime

		err := rows.Scan(
			&mapping.MappingID, &mapping.CampaignID, &mapping.ProviderType,
			&providerCampaignID, &providerOfferID, &providerData, &isActiveOnProvider, &lastSyncedAt,
			&mapping.CreatedAt, &mapping.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan campaign provider mapping: %w", err)
		}

		// Handle nullable fields
		if providerCampaignID.Valid {
			mapping.ProviderCampaignID = &providerCampaignID.String
		}
		if providerOfferID.Valid {
			mapping.ProviderOfferID = &providerOfferID.String
		}
		if providerData.Valid {
			mapping.ProviderData = &providerData.String
		}
		if isActiveOnProvider.Valid {
			mapping.IsActiveOnProvider = &isActiveOnProvider.Bool
		}
		if lastSyncedAt.Valid {
			mapping.LastSyncedAt = &lastSyncedAt.Time
		}

		mappings = append(mappings, mapping)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate campaign provider mappings: %w", err)
	}

	return mappings, nil
}

// DeleteCampaignProviderMapping deletes a campaign provider mapping by its ID
func (r *pgxCampaignProviderMappingRepository) DeleteCampaignProviderMapping(ctx context.Context, mappingID int64) error {
	query := `DELETE FROM public.campaign_provider_mappings WHERE mapping_id = $1`

	result, err := r.db.Exec(ctx, query, mappingID)
	if err != nil {
		return fmt.Errorf("failed to delete campaign provider mapping: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("campaign provider mapping not found: not found")
	}

	return nil
}
