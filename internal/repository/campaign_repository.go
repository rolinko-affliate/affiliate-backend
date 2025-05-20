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

// CampaignRepository defines the interface for campaign operations
type CampaignRepository interface {
	CreateCampaign(ctx context.Context, campaign *domain.Campaign) error
	GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error)
	UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error
	ListCampaignsByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.Campaign, error)
	ListCampaignsByAdvertiser(ctx context.Context, advertiserID int64, limit, offset int) ([]*domain.Campaign, error)
	DeleteCampaign(ctx context.Context, id int64) error
	
	// Provider offer methods
	CreateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error
	GetCampaignProviderOfferByID(ctx context.Context, id int64) (*domain.CampaignProviderOffer, error)
	UpdateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error
	ListCampaignProviderOffersByCampaign(ctx context.Context, campaignID int64) ([]*domain.CampaignProviderOffer, error)
	DeleteCampaignProviderOffer(ctx context.Context, id int64) error
}

// pgxCampaignRepository implements CampaignRepository using pgx
type pgxCampaignRepository struct {
	db *pgxpool.Pool
}

// NewPgxCampaignRepository creates a new campaign repository
func NewPgxCampaignRepository(db *pgxpool.Pool) CampaignRepository {
	return &pgxCampaignRepository{db: db}
}

// CreateCampaign creates a new campaign in the database
func (r *pgxCampaignRepository) CreateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	query := `INSERT INTO public.campaigns 
              (organization_id, advertiser_id, name, description, status, start_date, end_date, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
              RETURNING campaign_id, created_at, updated_at`
	
	var description sql.NullString
	var startDate, endDate sql.NullTime
	
	if campaign.Description != nil {
		description = sql.NullString{String: *campaign.Description, Valid: true}
	}
	
	if campaign.StartDate != nil {
		startDate = sql.NullTime{Time: *campaign.StartDate, Valid: true}
	}
	
	if campaign.EndDate != nil {
		endDate = sql.NullTime{Time: *campaign.EndDate, Valid: true}
	}
	
	now := time.Now()
	err := r.db.QueryRow(ctx, query, 
		campaign.OrganizationID, 
		campaign.AdvertiserID, 
		campaign.Name, 
		description, 
		campaign.Status, 
		startDate, 
		endDate, 
		now, 
		now,
	).Scan(
		&campaign.CampaignID, 
		&campaign.CreatedAt, 
		&campaign.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("error creating campaign: %w", err)
	}
	return nil
}

// GetCampaignByID retrieves a campaign by ID
func (r *pgxCampaignRepository) GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error) {
	query := `SELECT campaign_id, organization_id, advertiser_id, name, description, status, start_date, end_date, created_at, updated_at
              FROM public.campaigns WHERE campaign_id = $1`
	
	var campaign domain.Campaign
	var description sql.NullString
	var startDate, endDate sql.NullTime
	
	err := r.db.QueryRow(ctx, query, id).Scan(
		&campaign.CampaignID,
		&campaign.OrganizationID,
		&campaign.AdvertiserID,
		&campaign.Name,
		&description,
		&campaign.Status,
		&startDate,
		&endDate,
		&campaign.CreatedAt,
		&campaign.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("campaign not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("error getting campaign by ID: %w", err)
	}
	
	if description.Valid {
		desc := description.String
		campaign.Description = &desc
	}
	
	if startDate.Valid {
		date := startDate.Time
		campaign.StartDate = &date
	}
	
	if endDate.Valid {
		date := endDate.Time
		campaign.EndDate = &date
	}
	
	return &campaign, nil
}

// UpdateCampaign updates a campaign in the database
func (r *pgxCampaignRepository) UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	query := `UPDATE public.campaigns
              SET name = $1, description = $2, status = $3, start_date = $4, end_date = $5
              WHERE campaign_id = $6
              RETURNING updated_at`
	
	var description sql.NullString
	var startDate, endDate sql.NullTime
	
	if campaign.Description != nil {
		description = sql.NullString{String: *campaign.Description, Valid: true}
	}
	
	if campaign.StartDate != nil {
		startDate = sql.NullTime{Time: *campaign.StartDate, Valid: true}
	}
	
	if campaign.EndDate != nil {
		endDate = sql.NullTime{Time: *campaign.EndDate, Valid: true}
	}
	
	err := r.db.QueryRow(ctx, query, 
		campaign.Name, 
		description, 
		campaign.Status, 
		startDate, 
		endDate, 
		campaign.CampaignID,
	).Scan(&campaign.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("campaign not found: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("error updating campaign: %w", err)
	}
	
	return nil
}

// ListCampaignsByOrganization retrieves a list of campaigns for an organization with pagination
func (r *pgxCampaignRepository) ListCampaignsByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.Campaign, error) {
	query := `SELECT campaign_id, organization_id, advertiser_id, name, description, status, start_date, end_date, created_at, updated_at
              FROM public.campaigns
              WHERE organization_id = $1
              ORDER BY campaign_id
              LIMIT $2 OFFSET $3`
	
	return r.listCampaigns(ctx, query, orgID, limit, offset)
}

// ListCampaignsByAdvertiser retrieves a list of campaigns for an advertiser with pagination
func (r *pgxCampaignRepository) ListCampaignsByAdvertiser(ctx context.Context, advertiserID int64, limit, offset int) ([]*domain.Campaign, error) {
	query := `SELECT campaign_id, organization_id, advertiser_id, name, description, status, start_date, end_date, created_at, updated_at
              FROM public.campaigns
              WHERE advertiser_id = $1
              ORDER BY campaign_id
              LIMIT $2 OFFSET $3`
	
	return r.listCampaigns(ctx, query, advertiserID, limit, offset)
}

// listCampaigns is a helper function to list campaigns with a given query
func (r *pgxCampaignRepository) listCampaigns(ctx context.Context, query string, param int64, limit, offset int) ([]*domain.Campaign, error) {
	rows, err := r.db.Query(ctx, query, param, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing campaigns: %w", err)
	}
	defer rows.Close()
	
	var campaigns []*domain.Campaign
	for rows.Next() {
		var campaign domain.Campaign
		var description sql.NullString
		var startDate, endDate sql.NullTime
		
		if err := rows.Scan(
			&campaign.CampaignID,
			&campaign.OrganizationID,
			&campaign.AdvertiserID,
			&campaign.Name,
			&description,
			&campaign.Status,
			&startDate,
			&endDate,
			&campaign.CreatedAt,
			&campaign.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning campaign row: %w", err)
		}
		
		if description.Valid {
			desc := description.String
			campaign.Description = &desc
		}
		
		if startDate.Valid {
			date := startDate.Time
			campaign.StartDate = &date
		}
		
		if endDate.Valid {
			date := endDate.Time
			campaign.EndDate = &date
		}
		
		campaigns = append(campaigns, &campaign)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating campaign rows: %w", err)
	}
	
	return campaigns, nil
}

// DeleteCampaign deletes a campaign from the database
func (r *pgxCampaignRepository) DeleteCampaign(ctx context.Context, id int64) error {
	query := `DELETE FROM public.campaigns WHERE campaign_id = $1`
	
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting campaign: %w", err)
	}
	
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("campaign not found: %w", domain.ErrNotFound)
	}
	
	return nil
}

// CreateCampaignProviderOffer creates a new campaign provider offer in the database
func (r *pgxCampaignRepository) CreateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error {
	query := `INSERT INTO public.campaign_provider_offers 
              (campaign_id, provider_type, provider_offer_ref, provider_offer_config, is_active_on_provider, last_synced_at, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
              RETURNING provider_offer_id, created_at, updated_at`
	
	var providerOfferRef, providerOfferConfig sql.NullString
	var lastSyncedAt sql.NullTime
	
	if offer.ProviderOfferRef != nil {
		providerOfferRef = sql.NullString{String: *offer.ProviderOfferRef, Valid: true}
	}
	
	if offer.ProviderOfferConfig != nil {
		providerOfferConfig = sql.NullString{String: *offer.ProviderOfferConfig, Valid: true}
	}
	
	if offer.LastSyncedAt != nil {
		lastSyncedAt = sql.NullTime{Time: *offer.LastSyncedAt, Valid: true}
	}
	
	now := time.Now()
	err := r.db.QueryRow(ctx, query, 
		offer.CampaignID, 
		offer.ProviderType, 
		providerOfferRef, 
		providerOfferConfig, 
		offer.IsActiveOnProvider, 
		lastSyncedAt, 
		now, 
		now,
	).Scan(
		&offer.ProviderOfferID, 
		&offer.CreatedAt, 
		&offer.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("error creating campaign provider offer: %w", err)
	}
	
	return nil
}

// GetCampaignProviderOfferByID retrieves a campaign provider offer by ID
func (r *pgxCampaignRepository) GetCampaignProviderOfferByID(ctx context.Context, id int64) (*domain.CampaignProviderOffer, error) {
	query := `SELECT provider_offer_id, campaign_id, provider_type, provider_offer_ref, provider_offer_config, 
              is_active_on_provider, last_synced_at, created_at, updated_at
              FROM public.campaign_provider_offers WHERE provider_offer_id = $1`
	
	var offer domain.CampaignProviderOffer
	var providerOfferRef, providerOfferConfig sql.NullString
	var lastSyncedAt sql.NullTime
	
	err := r.db.QueryRow(ctx, query, id).Scan(
		&offer.ProviderOfferID,
		&offer.CampaignID,
		&offer.ProviderType,
		&providerOfferRef,
		&providerOfferConfig,
		&offer.IsActiveOnProvider,
		&lastSyncedAt,
		&offer.CreatedAt,
		&offer.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("campaign provider offer not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("error getting campaign provider offer by ID: %w", err)
	}
	
	if providerOfferRef.Valid {
		ref := providerOfferRef.String
		offer.ProviderOfferRef = &ref
	}
	
	if providerOfferConfig.Valid {
		config := providerOfferConfig.String
		offer.ProviderOfferConfig = &config
	}
	
	if lastSyncedAt.Valid {
		synced := lastSyncedAt.Time
		offer.LastSyncedAt = &synced
	}
	
	return &offer, nil
}

// UpdateCampaignProviderOffer updates a campaign provider offer in the database
func (r *pgxCampaignRepository) UpdateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error {
	query := `UPDATE public.campaign_provider_offers
              SET provider_offer_ref = $1, provider_offer_config = $2, is_active_on_provider = $3, last_synced_at = $4
              WHERE provider_offer_id = $5
              RETURNING updated_at`
	
	var providerOfferRef, providerOfferConfig sql.NullString
	var lastSyncedAt sql.NullTime
	
	if offer.ProviderOfferRef != nil {
		providerOfferRef = sql.NullString{String: *offer.ProviderOfferRef, Valid: true}
	}
	
	if offer.ProviderOfferConfig != nil {
		providerOfferConfig = sql.NullString{String: *offer.ProviderOfferConfig, Valid: true}
	}
	
	if offer.LastSyncedAt != nil {
		lastSyncedAt = sql.NullTime{Time: *offer.LastSyncedAt, Valid: true}
	}
	
	err := r.db.QueryRow(ctx, query, 
		providerOfferRef, 
		providerOfferConfig, 
		offer.IsActiveOnProvider, 
		lastSyncedAt, 
		offer.ProviderOfferID,
	).Scan(&offer.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("campaign provider offer not found: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("error updating campaign provider offer: %w", err)
	}
	
	return nil
}

// ListCampaignProviderOffersByCampaign retrieves a list of campaign provider offers for a campaign
func (r *pgxCampaignRepository) ListCampaignProviderOffersByCampaign(ctx context.Context, campaignID int64) ([]*domain.CampaignProviderOffer, error) {
	query := `SELECT provider_offer_id, campaign_id, provider_type, provider_offer_ref, provider_offer_config, 
              is_active_on_provider, last_synced_at, created_at, updated_at
              FROM public.campaign_provider_offers
              WHERE campaign_id = $1
              ORDER BY provider_offer_id`
	
	rows, err := r.db.Query(ctx, query, campaignID)
	if err != nil {
		return nil, fmt.Errorf("error listing campaign provider offers: %w", err)
	}
	defer rows.Close()
	
	var offers []*domain.CampaignProviderOffer
	for rows.Next() {
		var offer domain.CampaignProviderOffer
		var providerOfferRef, providerOfferConfig sql.NullString
		var lastSyncedAt sql.NullTime
		
		if err := rows.Scan(
			&offer.ProviderOfferID,
			&offer.CampaignID,
			&offer.ProviderType,
			&providerOfferRef,
			&providerOfferConfig,
			&offer.IsActiveOnProvider,
			&lastSyncedAt,
			&offer.CreatedAt,
			&offer.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning campaign provider offer row: %w", err)
		}
		
		if providerOfferRef.Valid {
			ref := providerOfferRef.String
			offer.ProviderOfferRef = &ref
		}
		
		if providerOfferConfig.Valid {
			config := providerOfferConfig.String
			offer.ProviderOfferConfig = &config
		}
		
		if lastSyncedAt.Valid {
			synced := lastSyncedAt.Time
			offer.LastSyncedAt = &synced
		}
		
		offers = append(offers, &offer)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating campaign provider offer rows: %w", err)
	}
	
	return offers, nil
}

// DeleteCampaignProviderOffer deletes a campaign provider offer from the database
func (r *pgxCampaignRepository) DeleteCampaignProviderOffer(ctx context.Context, id int64) error {
	query := `DELETE FROM public.campaign_provider_offers WHERE provider_offer_id = $1`
	
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting campaign provider offer: %w", err)
	}
	
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("campaign provider offer not found: %w", domain.ErrNotFound)
	}
	
	return nil
}