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
	CreateMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error
	GetMappingByID(ctx context.Context, id int64) (*domain.AdvertiserProviderMapping, error)
	GetMappingsByAdvertiserID(ctx context.Context, advertiserID int64) ([]*domain.AdvertiserProviderMapping, error)
	GetMappingByAdvertiserAndProvider(ctx context.Context, advertiserID int64, providerType string) (*domain.AdvertiserProviderMapping, error)
	UpdateMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error
	DeleteMapping(ctx context.Context, id int64) error
	UpdateSyncStatus(ctx context.Context, mappingID int64, status string, syncError *string) error
}

type pgxAdvertiserProviderMappingRepository struct {
	db *pgxpool.Pool
}

func NewAdvertiserProviderMappingRepository(db *pgxpool.Pool) AdvertiserProviderMappingRepository {
	return &pgxAdvertiserProviderMappingRepository{db: db}
}

func (r *pgxAdvertiserProviderMappingRepository) CreateMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error {
	query := `INSERT INTO public.advertiser_provider_mappings (
		advertiser_id, provider_type, provider_advertiser_id, api_credentials, provider_config,
		provider_data, sync_status, last_sync_at, sync_error, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
	) RETURNING mapping_id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		mapping.AdvertiserID,
		mapping.ProviderType,
		mapping.ProviderAdvertiserID,
		mapping.APICredentials,
		mapping.ProviderConfig,
		mapping.ProviderData,
		mapping.SyncStatus,
		mapping.LastSyncAt,
		mapping.SyncError,
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

func (r *pgxAdvertiserProviderMappingRepository) GetMappingByID(ctx context.Context, id int64) (*domain.AdvertiserProviderMapping, error) {
	query := `SELECT mapping_id, advertiser_id, provider_type, provider_advertiser_id, api_credentials, 
		provider_config, provider_data, sync_status, last_sync_at, sync_error, created_at, updated_at
		FROM public.advertiser_provider_mappings WHERE mapping_id = $1`

	var mapping domain.AdvertiserProviderMapping
	var providerAdvertiserID, apiCredentials, providerConfig, providerData sql.NullString
	var syncStatus, syncError sql.NullString
	var lastSyncAt sql.NullTime

	err := r.db.QueryRow(ctx, query, id).Scan(
		&mapping.MappingID,
		&mapping.AdvertiserID,
		&mapping.ProviderType,
		&providerAdvertiserID,
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
			return nil, fmt.Errorf("advertiser provider mapping not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("error getting advertiser provider mapping by ID: %w", err)
	}

	// Handle nullable fields
	if providerAdvertiserID.Valid {
		mapping.ProviderAdvertiserID = &providerAdvertiserID.String
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

func (r *pgxAdvertiserProviderMappingRepository) GetMappingsByAdvertiserID(ctx context.Context, advertiserID int64) ([]*domain.AdvertiserProviderMapping, error) {
	query := `SELECT mapping_id, advertiser_id, provider_type, provider_advertiser_id, api_credentials, 
		provider_config, provider_data, sync_status, last_sync_at, sync_error, created_at, updated_at
		FROM public.advertiser_provider_mappings WHERE advertiser_id = $1 ORDER BY mapping_id`

	rows, err := r.db.Query(ctx, query, advertiserID)
	if err != nil {
		return nil, fmt.Errorf("error listing advertiser provider mappings: %w", err)
	}
	defer rows.Close()

	var mappings []*domain.AdvertiserProviderMapping
	for rows.Next() {
		var mapping domain.AdvertiserProviderMapping
		var providerAdvertiserID, apiCredentials, providerConfig, providerData sql.NullString
		var syncStatus, syncError sql.NullString
		var lastSyncAt sql.NullTime

		if err := rows.Scan(
			&mapping.MappingID,
			&mapping.AdvertiserID,
			&mapping.ProviderType,
			&providerAdvertiserID,
			&apiCredentials,
			&providerConfig,
			&providerData,
			&syncStatus,
			&lastSyncAt,
			&syncError,
			&mapping.CreatedAt,
			&mapping.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning advertiser provider mapping row: %w", err)
		}

		// Handle nullable fields
		if providerAdvertiserID.Valid {
			mapping.ProviderAdvertiserID = &providerAdvertiserID.String
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

		mappings = append(mappings, &mapping)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating advertiser provider mapping rows: %w", err)
	}

	return mappings, nil
}

func (r *pgxAdvertiserProviderMappingRepository) GetMappingByAdvertiserAndProvider(ctx context.Context, advertiserID int64, providerType string) (*domain.AdvertiserProviderMapping, error) {
	query := `SELECT mapping_id, advertiser_id, provider_type, provider_advertiser_id, api_credentials, 
		provider_config, provider_data, sync_status, last_sync_at, sync_error, created_at, updated_at
		FROM public.advertiser_provider_mappings WHERE advertiser_id = $1 AND provider_type = $2`

	var mapping domain.AdvertiserProviderMapping
	var providerAdvertiserID, apiCredentials, providerConfig, providerData sql.NullString
	var syncStatus, syncError sql.NullString
	var lastSyncAt sql.NullTime

	err := r.db.QueryRow(ctx, query, advertiserID, providerType).Scan(
		&mapping.MappingID,
		&mapping.AdvertiserID,
		&mapping.ProviderType,
		&providerAdvertiserID,
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
			return nil, fmt.Errorf("advertiser provider mapping not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("error getting advertiser provider mapping: %w", err)
	}

	// Handle nullable fields
	if providerAdvertiserID.Valid {
		mapping.ProviderAdvertiserID = &providerAdvertiserID.String
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

func (r *pgxAdvertiserProviderMappingRepository) UpdateMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error {
	query := `UPDATE public.advertiser_provider_mappings SET 
		provider_advertiser_id = $1, api_credentials = $2, provider_config = $3, provider_data = $4,
		sync_status = $5, last_sync_at = $6, sync_error = $7, updated_at = $8
		WHERE mapping_id = $9
		RETURNING updated_at`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		mapping.ProviderAdvertiserID,
		mapping.APICredentials,
		mapping.ProviderConfig,
		mapping.ProviderData,
		mapping.SyncStatus,
		mapping.LastSyncAt,
		mapping.SyncError,
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

func (r *pgxAdvertiserProviderMappingRepository) DeleteMapping(ctx context.Context, id int64) error {
	query := `DELETE FROM public.advertiser_provider_mappings WHERE mapping_id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting advertiser provider mapping: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("advertiser provider mapping not found: %w", domain.ErrNotFound)
	}

	return nil
}

func (r *pgxAdvertiserProviderMappingRepository) UpdateSyncStatus(ctx context.Context, mappingID int64, status string, syncError *string) error {
	query := `UPDATE public.advertiser_provider_mappings SET 
		sync_status = $1, last_sync_at = $2, sync_error = $3, updated_at = $4
		WHERE mapping_id = $5`

	now := time.Now()
	_, err := r.db.Exec(ctx, query, status, now, syncError, now, mappingID)
	if err != nil {
		return fmt.Errorf("error updating sync status: %w", err)
	}

	return nil
}
