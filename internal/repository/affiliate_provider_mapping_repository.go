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
	CreateMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error
	GetMappingByID(ctx context.Context, id int64) (*domain.AffiliateProviderMapping, error)
	GetMappingsByAffiliateID(ctx context.Context, affiliateID int64) ([]*domain.AffiliateProviderMapping, error)
	GetMappingByAffiliateAndProvider(ctx context.Context, affiliateID int64, providerType string) (*domain.AffiliateProviderMapping, error)
	UpdateMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error
	DeleteMapping(ctx context.Context, id int64) error
	UpdateSyncStatus(ctx context.Context, mappingID int64, status string, syncError *string) error
	
	// Legacy methods for backward compatibility
	CreateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error
	GetAffiliateProviderMapping(ctx context.Context, affiliateID int64, providerType string) (*domain.AffiliateProviderMapping, error)
	UpdateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error
	DeleteAffiliateProviderMapping(ctx context.Context, mappingID int64) error
}

type pgxAffiliateProviderMappingRepository struct {
	db *pgxpool.Pool
}

func NewAffiliateProviderMappingRepository(db *pgxpool.Pool) AffiliateProviderMappingRepository {
	return &pgxAffiliateProviderMappingRepository{db: db}
}

func NewPgxAffiliateProviderMappingRepository(db *pgxpool.Pool) AffiliateProviderMappingRepository {
	return &pgxAffiliateProviderMappingRepository{db: db}
}

func (r *pgxAffiliateProviderMappingRepository) CreateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error {
	query := `INSERT INTO public.affiliate_provider_mappings 
              (affiliate_id, provider_type, provider_affiliate_id, api_credentials, provider_config, provider_data, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
              RETURNING mapping_id, created_at, updated_at`

	var providerAffiliateID, apiCredentials, providerConfig, providerData sql.NullString

	if mapping.ProviderAffiliateID != nil {
		providerAffiliateID = sql.NullString{String: *mapping.ProviderAffiliateID, Valid: true}
	}
	if mapping.APICredentials != nil {
		apiCredentials = sql.NullString{String: *mapping.APICredentials, Valid: true}
	}
	if mapping.ProviderConfig != nil {
		providerConfig = sql.NullString{String: *mapping.ProviderConfig, Valid: true}
	}
	if mapping.ProviderData != nil {
		providerData = sql.NullString{String: *mapping.ProviderData, Valid: true}
	}

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		mapping.AffiliateID,
		mapping.ProviderType,
		providerAffiliateID,
		apiCredentials,
		providerConfig,
		providerData,
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
	query := `SELECT mapping_id, affiliate_id, provider_type, provider_affiliate_id, api_credentials, provider_config, provider_data,
	          sync_status, last_sync_at, sync_error, created_at, updated_at
	          FROM public.affiliate_provider_mappings 
	          WHERE affiliate_id = $1 AND provider_type = $2`

	var mapping domain.AffiliateProviderMapping
	var providerAffiliateID, apiCredentials, providerConfig, providerData, syncStatus, syncError sql.NullString
	var lastSyncAt sql.NullTime

	err := r.db.QueryRow(ctx, query, affiliateID, providerType).Scan(
		&mapping.MappingID,
		&mapping.AffiliateID,
		&mapping.ProviderType,
		&providerAffiliateID,
		&apiCredentials,
		&providerConfig,
		&providerData,
		&syncStatus,
		&lastSyncAt,
		&syncError,
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
	if providerData.Valid {
		mapping.ProviderData = &providerData.String
	}
	if syncStatus.Valid {
		mapping.SyncStatus = &syncStatus.String
	}
	if lastSyncAt.Valid {
		mapping.LastSyncAt = &lastSyncAt.Time
	}
	if syncError.Valid {
		mapping.SyncError = &syncError.String
	}

	return &mapping, nil
}

func (r *pgxAffiliateProviderMappingRepository) UpdateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error {
	query := `UPDATE public.affiliate_provider_mappings SET 
	          provider_affiliate_id = $1, api_credentials = $2, provider_config = $3, provider_data = $4,
	          sync_status = $5, last_sync_at = $6, sync_error = $7, updated_at = $8
	          WHERE mapping_id = $9
	          RETURNING updated_at`

	var providerAffiliateID, apiCredentials, providerConfig, providerData, syncStatus, syncError sql.NullString
	var lastSyncAt sql.NullTime

	if mapping.ProviderAffiliateID != nil {
		providerAffiliateID = sql.NullString{String: *mapping.ProviderAffiliateID, Valid: true}
	}
	if mapping.APICredentials != nil {
		apiCredentials = sql.NullString{String: *mapping.APICredentials, Valid: true}
	}
	if mapping.ProviderConfig != nil {
		providerConfig = sql.NullString{String: *mapping.ProviderConfig, Valid: true}
	}
	if mapping.ProviderData != nil {
		providerData = sql.NullString{String: *mapping.ProviderData, Valid: true}
	}
	if mapping.SyncStatus != nil {
		syncStatus = sql.NullString{String: *mapping.SyncStatus, Valid: true}
	}
	if mapping.LastSyncAt != nil {
		lastSyncAt = sql.NullTime{Time: *mapping.LastSyncAt, Valid: true}
	}
	if mapping.SyncError != nil {
		syncError = sql.NullString{String: *mapping.SyncError, Valid: true}
	}

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		providerAffiliateID,
		apiCredentials,
		providerConfig,
		providerData,
		syncStatus,
		lastSyncAt,
		syncError,
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

// New interface methods for sync script compatibility

func (r *pgxAffiliateProviderMappingRepository) CreateMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error {
	return r.CreateAffiliateProviderMapping(ctx, mapping)
}

func (r *pgxAffiliateProviderMappingRepository) GetMappingByID(ctx context.Context, id int64) (*domain.AffiliateProviderMapping, error) {
	query := `SELECT mapping_id, affiliate_id, provider_type, provider_affiliate_id, api_credentials, 
	                 provider_config, provider_data, sync_status, last_sync_at, sync_error, created_at, updated_at
	          FROM public.affiliate_provider_mappings 
	          WHERE mapping_id = $1`

	mapping := &domain.AffiliateProviderMapping{}
	var providerAffiliateID, apiCredentials, providerConfig, providerData, syncStatus, syncError sql.NullString
	var lastSyncAt sql.NullTime

	err := r.db.QueryRow(ctx, query, id).Scan(
		&mapping.MappingID,
		&mapping.AffiliateID,
		&mapping.ProviderType,
		&providerAffiliateID,
		&apiCredentials,
		&providerConfig,
		&providerData,
		&syncStatus,
		&lastSyncAt,
		&syncError,
		&mapping.CreatedAt,
		&mapping.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("affiliate provider mapping not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("error getting affiliate provider mapping: %w", err)
	}

	// Handle nullable fields
	if providerAffiliateID.Valid {
		mapping.ProviderAffiliateID = &providerAffiliateID.String
	}
	if apiCredentials.Valid {
		mapping.APICredentials = &apiCredentials.String
	}
	if providerConfig.Valid {
		mapping.ProviderConfig = &providerConfig.String
	}
	if providerData.Valid {
		mapping.ProviderData = &providerData.String
	}
	if syncStatus.Valid {
		mapping.SyncStatus = &syncStatus.String
	}
	if lastSyncAt.Valid {
		mapping.LastSyncAt = &lastSyncAt.Time
	}
	if syncError.Valid {
		mapping.SyncError = &syncError.String
	}

	return mapping, nil
}

func (r *pgxAffiliateProviderMappingRepository) GetMappingsByAffiliateID(ctx context.Context, affiliateID int64) ([]*domain.AffiliateProviderMapping, error) {
	query := `SELECT mapping_id, affiliate_id, provider_type, provider_affiliate_id, api_credentials, 
	                 provider_config, provider_data, sync_status, last_sync_at, sync_error, created_at, updated_at
	          FROM public.affiliate_provider_mappings 
	          WHERE affiliate_id = $1`

	rows, err := r.db.Query(ctx, query, affiliateID)
	if err != nil {
		return nil, fmt.Errorf("error querying affiliate provider mappings: %w", err)
	}
	defer rows.Close()

	var mappings []*domain.AffiliateProviderMapping
	for rows.Next() {
		mapping := &domain.AffiliateProviderMapping{}
		var providerAffiliateID, apiCredentials, providerConfig, providerData, syncStatus, syncError sql.NullString
		var lastSyncAt sql.NullTime

		err := rows.Scan(
			&mapping.MappingID,
			&mapping.AffiliateID,
			&mapping.ProviderType,
			&providerAffiliateID,
			&apiCredentials,
			&providerConfig,
			&providerData,
			&syncStatus,
			&lastSyncAt,
			&syncError,
			&mapping.CreatedAt,
			&mapping.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning affiliate provider mapping: %w", err)
		}

		// Handle nullable fields
		if providerAffiliateID.Valid {
			mapping.ProviderAffiliateID = &providerAffiliateID.String
		}
		if apiCredentials.Valid {
			mapping.APICredentials = &apiCredentials.String
		}
		if providerConfig.Valid {
			mapping.ProviderConfig = &providerConfig.String
		}
		if providerData.Valid {
			mapping.ProviderData = &providerData.String
		}
		if syncStatus.Valid {
			mapping.SyncStatus = &syncStatus.String
		}
		if lastSyncAt.Valid {
			mapping.LastSyncAt = &lastSyncAt.Time
		}
		if syncError.Valid {
			mapping.SyncError = &syncError.String
		}

		mappings = append(mappings, mapping)
	}

	return mappings, nil
}

func (r *pgxAffiliateProviderMappingRepository) GetMappingByAffiliateAndProvider(ctx context.Context, affiliateID int64, providerType string) (*domain.AffiliateProviderMapping, error) {
	return r.GetAffiliateProviderMapping(ctx, affiliateID, providerType)
}

func (r *pgxAffiliateProviderMappingRepository) UpdateMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error {
	return r.UpdateAffiliateProviderMapping(ctx, mapping)
}

func (r *pgxAffiliateProviderMappingRepository) DeleteMapping(ctx context.Context, id int64) error {
	return r.DeleteAffiliateProviderMapping(ctx, id)
}

func (r *pgxAffiliateProviderMappingRepository) UpdateSyncStatus(ctx context.Context, mappingID int64, status string, syncError *string) error {
	query := `UPDATE public.affiliate_provider_mappings 
	          SET sync_status = $1, sync_error = $2, last_sync_at = $3, updated_at = $4
	          WHERE mapping_id = $5`

	var syncErrorValue sql.NullString
	if syncError != nil {
		syncErrorValue = sql.NullString{String: *syncError, Valid: true}
	}

	now := time.Now()
	var lastSyncAt *time.Time
	if status == "synced" {
		lastSyncAt = &now
	}

	commandTag, err := r.db.Exec(ctx, query, status, syncErrorValue, lastSyncAt, now, mappingID)
	if err != nil {
		return fmt.Errorf("error updating affiliate provider mapping sync status: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("affiliate provider mapping not found: %w", domain.ErrNotFound)
	}

	return nil
}
