package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TrackingLinkProviderMappingRepository defines the interface for tracking link provider mapping data access
type TrackingLinkProviderMappingRepository interface {
	CreateTrackingLinkProviderMapping(ctx context.Context, mapping *domain.TrackingLinkProviderMapping) error
	GetTrackingLinkProviderMapping(ctx context.Context, trackingLinkID int64, providerType string) (*domain.TrackingLinkProviderMapping, error)
	GetTrackingLinkProviderMappingByID(ctx context.Context, mappingID int64) (*domain.TrackingLinkProviderMapping, error)
	UpdateTrackingLinkProviderMapping(ctx context.Context, mapping *domain.TrackingLinkProviderMapping) error
	DeleteTrackingLinkProviderMapping(ctx context.Context, mappingID int64) error
	ListTrackingLinkProviderMappingsByTrackingLink(ctx context.Context, trackingLinkID int64) ([]*domain.TrackingLinkProviderMapping, error)
	ListTrackingLinkProviderMappingsByProvider(ctx context.Context, providerType string, limit, offset int) ([]*domain.TrackingLinkProviderMapping, error)
}

// trackingLinkProviderMappingRepository implements TrackingLinkProviderMappingRepository
type trackingLinkProviderMappingRepository struct {
	db *pgxpool.Pool
}

// NewTrackingLinkProviderMappingRepository creates a new tracking link provider mapping repository
func NewTrackingLinkProviderMappingRepository(db *pgxpool.Pool) TrackingLinkProviderMappingRepository {
	return &trackingLinkProviderMappingRepository{db: db}
}

// CreateTrackingLinkProviderMapping creates a new tracking link provider mapping
func (r *trackingLinkProviderMappingRepository) CreateTrackingLinkProviderMapping(ctx context.Context, mapping *domain.TrackingLinkProviderMapping) error {
	query := `
		INSERT INTO public.tracking_link_provider_mappings (
			tracking_link_id, provider_type, provider_tracking_link_id,
			provider_data, sync_status, last_sync_at, sync_error,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING mapping_id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		mapping.TrackingLinkID,
		mapping.ProviderType,
		mapping.ProviderTrackingLinkID,
		mapping.ProviderData,
		mapping.SyncStatus,
		mapping.LastSyncAt,
		mapping.SyncError,
		mapping.CreatedAt,
		mapping.UpdatedAt,
	).Scan(&mapping.MappingID, &mapping.CreatedAt, &mapping.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create tracking link provider mapping: %w", err)
	}

	return nil
}

// GetTrackingLinkProviderMapping retrieves a tracking link provider mapping by tracking link ID and provider type
func (r *trackingLinkProviderMappingRepository) GetTrackingLinkProviderMapping(ctx context.Context, trackingLinkID int64, providerType string) (*domain.TrackingLinkProviderMapping, error) {
	query := `
		SELECT mapping_id, tracking_link_id, provider_type, provider_tracking_link_id,
			   provider_data, sync_status, last_sync_at, sync_error,
			   created_at, updated_at
		FROM public.tracking_link_provider_mappings
		WHERE tracking_link_id = $1 AND provider_type = $2`

	mapping := &domain.TrackingLinkProviderMapping{}
	err := r.db.QueryRow(ctx, query, trackingLinkID, providerType).Scan(
		&mapping.MappingID,
		&mapping.TrackingLinkID,
		&mapping.ProviderType,
		&mapping.ProviderTrackingLinkID,
		&mapping.ProviderData,
		&mapping.SyncStatus,
		&mapping.LastSyncAt,
		&mapping.SyncError,
		&mapping.CreatedAt,
		&mapping.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tracking link provider mapping not found")
		}
		return nil, fmt.Errorf("failed to get tracking link provider mapping: %w", err)
	}

	return mapping, nil
}

// GetTrackingLinkProviderMappingByID retrieves a tracking link provider mapping by its ID
func (r *trackingLinkProviderMappingRepository) GetTrackingLinkProviderMappingByID(ctx context.Context, mappingID int64) (*domain.TrackingLinkProviderMapping, error) {
	query := `
		SELECT mapping_id, tracking_link_id, provider_type, provider_tracking_link_id,
			   provider_data, sync_status, last_sync_at, sync_error,
			   created_at, updated_at
		FROM public.tracking_link_provider_mappings
		WHERE mapping_id = $1`

	mapping := &domain.TrackingLinkProviderMapping{}
	err := r.db.QueryRow(ctx, query, mappingID).Scan(
		&mapping.MappingID,
		&mapping.TrackingLinkID,
		&mapping.ProviderType,
		&mapping.ProviderTrackingLinkID,
		&mapping.ProviderData,
		&mapping.SyncStatus,
		&mapping.LastSyncAt,
		&mapping.SyncError,
		&mapping.CreatedAt,
		&mapping.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tracking link provider mapping not found")
		}
		return nil, fmt.Errorf("failed to get tracking link provider mapping: %w", err)
	}

	return mapping, nil
}

// UpdateTrackingLinkProviderMapping updates an existing tracking link provider mapping
func (r *trackingLinkProviderMappingRepository) UpdateTrackingLinkProviderMapping(ctx context.Context, mapping *domain.TrackingLinkProviderMapping) error {
	query := `
		UPDATE public.tracking_link_provider_mappings SET
			provider_tracking_link_id = $2,
			provider_data = $3,
			sync_status = $4,
			last_sync_at = $5,
			sync_error = $6,
			updated_at = CURRENT_TIMESTAMP
		WHERE mapping_id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query,
		mapping.MappingID,
		mapping.ProviderTrackingLinkID,
		mapping.ProviderData,
		mapping.SyncStatus,
		mapping.LastSyncAt,
		mapping.SyncError,
	).Scan(&mapping.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("tracking link provider mapping not found")
		}
		return fmt.Errorf("failed to update tracking link provider mapping: %w", err)
	}

	return nil
}

// DeleteTrackingLinkProviderMapping deletes a tracking link provider mapping by its ID
func (r *trackingLinkProviderMappingRepository) DeleteTrackingLinkProviderMapping(ctx context.Context, mappingID int64) error {
	query := `DELETE FROM public.tracking_link_provider_mappings WHERE mapping_id = $1`

	result, err := r.db.Exec(ctx, query, mappingID)
	if err != nil {
		return fmt.Errorf("failed to delete tracking link provider mapping: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("tracking link provider mapping not found")
	}

	return nil
}

// ListTrackingLinkProviderMappingsByTrackingLink retrieves all provider mappings for a specific tracking link
func (r *trackingLinkProviderMappingRepository) ListTrackingLinkProviderMappingsByTrackingLink(ctx context.Context, trackingLinkID int64) ([]*domain.TrackingLinkProviderMapping, error) {
	query := `
		SELECT mapping_id, tracking_link_id, provider_type, provider_tracking_link_id,
			   provider_data, sync_status, last_sync_at, sync_error,
			   created_at, updated_at
		FROM public.tracking_link_provider_mappings
		WHERE tracking_link_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, trackingLinkID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tracking link provider mappings by tracking link: %w", err)
	}
	defer rows.Close()

	mappings := make([]*domain.TrackingLinkProviderMapping, 0)
	for rows.Next() {
		mapping := &domain.TrackingLinkProviderMapping{}
		err := rows.Scan(
			&mapping.MappingID,
			&mapping.TrackingLinkID,
			&mapping.ProviderType,
			&mapping.ProviderTrackingLinkID,
			&mapping.ProviderData,
			&mapping.SyncStatus,
			&mapping.LastSyncAt,
			&mapping.SyncError,
			&mapping.CreatedAt,
			&mapping.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tracking link provider mapping: %w", err)
		}
		mappings = append(mappings, mapping)
	}

	return mappings, nil
}

// ListTrackingLinkProviderMappingsByProvider retrieves tracking link provider mappings for a specific provider
func (r *trackingLinkProviderMappingRepository) ListTrackingLinkProviderMappingsByProvider(ctx context.Context, providerType string, limit, offset int) ([]*domain.TrackingLinkProviderMapping, error) {
	query := `
		SELECT mapping_id, tracking_link_id, provider_type, provider_tracking_link_id,
			   provider_data, sync_status, last_sync_at, sync_error,
			   created_at, updated_at
		FROM public.tracking_link_provider_mappings
		WHERE provider_type = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, providerType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tracking link provider mappings by provider: %w", err)
	}
	defer rows.Close()

	mappings := make([]*domain.TrackingLinkProviderMapping, 0)
	for rows.Next() {
		mapping := &domain.TrackingLinkProviderMapping{}
		err := rows.Scan(
			&mapping.MappingID,
			&mapping.TrackingLinkID,
			&mapping.ProviderType,
			&mapping.ProviderTrackingLinkID,
			&mapping.ProviderData,
			&mapping.SyncStatus,
			&mapping.LastSyncAt,
			&mapping.SyncError,
			&mapping.CreatedAt,
			&mapping.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tracking link provider mapping: %w", err)
		}
		mappings = append(mappings, mapping)
	}

	return mappings, nil
}
