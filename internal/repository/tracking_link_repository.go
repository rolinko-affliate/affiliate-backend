package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TrackingLinkRepository defines the interface for tracking link data access
type TrackingLinkRepository interface {
	CreateTrackingLink(ctx context.Context, trackingLink *domain.TrackingLink) error
	GetTrackingLinkByID(ctx context.Context, trackingLinkID int64) (*domain.TrackingLink, error)
	UpdateTrackingLink(ctx context.Context, trackingLink *domain.TrackingLink) error
	DeleteTrackingLink(ctx context.Context, trackingLinkID int64) error
	ListTrackingLinksByCampaign(ctx context.Context, campaignID int64, limit, offset int) ([]*domain.TrackingLink, error)
	ListTrackingLinksByAffiliate(ctx context.Context, affiliateID int64, limit, offset int) ([]*domain.TrackingLink, error)
	ListTrackingLinksByOrganization(ctx context.Context, organizationID int64, limit, offset int) ([]*domain.TrackingLink, error)
	GetTrackingLinkByCampaignAndAffiliate(ctx context.Context, campaignID, affiliateID int64, sourceID, sub1, sub2, sub3, sub4, sub5 *string) (*domain.TrackingLink, error)

}

// trackingLinkRepository implements TrackingLinkRepository
type trackingLinkRepository struct {
	db *pgxpool.Pool
}

// NewTrackingLinkRepository creates a new tracking link repository
func NewTrackingLinkRepository(db *pgxpool.Pool) TrackingLinkRepository {
	return &trackingLinkRepository{db: db}
}

// CreateTrackingLink creates a new tracking link
func (r *trackingLinkRepository) CreateTrackingLink(ctx context.Context, trackingLink *domain.TrackingLink) error {
	query := `
		INSERT INTO public.tracking_links (
			organization_id, campaign_id, affiliate_id, name, description, status,
			tracking_url, source_id, sub1, sub2, sub3, sub4, sub5,
			is_encrypt_parameters, is_redirect_link,
			internal_notes, tags,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		RETURNING tracking_link_id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		trackingLink.OrganizationID,
		trackingLink.CampaignID,
		trackingLink.AffiliateID,
		trackingLink.Name,
		trackingLink.Description,
		trackingLink.Status,
		trackingLink.TrackingURL,
		trackingLink.SourceID,
		trackingLink.Sub1,
		trackingLink.Sub2,
		trackingLink.Sub3,
		trackingLink.Sub4,
		trackingLink.Sub5,
		trackingLink.IsEncryptParameters,
		trackingLink.IsRedirectLink,
		trackingLink.InternalNotes,
		trackingLink.Tags,
		trackingLink.CreatedAt,
		trackingLink.UpdatedAt,
	).Scan(&trackingLink.TrackingLinkID, &trackingLink.CreatedAt, &trackingLink.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create tracking link: %w", err)
	}

	return nil
}

// GetTrackingLinkByID retrieves a tracking link by its ID
func (r *trackingLinkRepository) GetTrackingLinkByID(ctx context.Context, trackingLinkID int64) (*domain.TrackingLink, error) {
	query := `
		SELECT tracking_link_id, organization_id, campaign_id, affiliate_id, name, description, status,
			   tracking_url, source_id, sub1, sub2, sub3, sub4, sub5,
			   is_encrypt_parameters, is_redirect_link,
			   internal_notes, tags,
			   created_at, updated_at
		FROM public.tracking_links
		WHERE tracking_link_id = $1`

	trackingLink := &domain.TrackingLink{}
	err := r.db.QueryRow(ctx, query, trackingLinkID).Scan(
		&trackingLink.TrackingLinkID,
		&trackingLink.OrganizationID,
		&trackingLink.CampaignID,
		&trackingLink.AffiliateID,
		&trackingLink.Name,
		&trackingLink.Description,
		&trackingLink.Status,
		&trackingLink.TrackingURL,
		&trackingLink.SourceID,
		&trackingLink.Sub1,
		&trackingLink.Sub2,
		&trackingLink.Sub3,
		&trackingLink.Sub4,
		&trackingLink.Sub5,
		&trackingLink.IsEncryptParameters,
		&trackingLink.IsRedirectLink,
		&trackingLink.InternalNotes,
		&trackingLink.Tags,
		&trackingLink.CreatedAt,
		&trackingLink.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tracking link not found")
		}
		return nil, fmt.Errorf("failed to get tracking link: %w", err)
	}

	return trackingLink, nil
}

// UpdateTrackingLink updates an existing tracking link
func (r *trackingLinkRepository) UpdateTrackingLink(ctx context.Context, trackingLink *domain.TrackingLink) error {
	query := `
		UPDATE public.tracking_links SET
			name = $2, description = $3, status = $4,
			tracking_url = $5, source_id = $6, sub1 = $7, sub2 = $8, sub3 = $9, sub4 = $10, sub5 = $11,
			is_encrypt_parameters = $12, is_redirect_link = $13,
			internal_notes = $14, tags = $15,
			updated_at = CURRENT_TIMESTAMP
		WHERE tracking_link_id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query,
		trackingLink.TrackingLinkID,
		trackingLink.Name,
		trackingLink.Description,
		trackingLink.Status,
		trackingLink.TrackingURL,
		trackingLink.SourceID,
		trackingLink.Sub1,
		trackingLink.Sub2,
		trackingLink.Sub3,
		trackingLink.Sub4,
		trackingLink.Sub5,
		trackingLink.IsEncryptParameters,
		trackingLink.IsRedirectLink,
		trackingLink.InternalNotes,
		trackingLink.Tags,
	).Scan(&trackingLink.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("tracking link not found")
		}
		return fmt.Errorf("failed to update tracking link: %w", err)
	}

	return nil
}

// DeleteTrackingLink deletes a tracking link by its ID
func (r *trackingLinkRepository) DeleteTrackingLink(ctx context.Context, trackingLinkID int64) error {
	query := `DELETE FROM public.tracking_links WHERE tracking_link_id = $1`

	result, err := r.db.Exec(ctx, query, trackingLinkID)
	if err != nil {
		return fmt.Errorf("failed to delete tracking link: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("tracking link not found")
	}

	return nil
}

// ListTrackingLinksByCampaign retrieves tracking links for a specific campaign
func (r *trackingLinkRepository) ListTrackingLinksByCampaign(ctx context.Context, campaignID int64, limit, offset int) ([]*domain.TrackingLink, error) {
	query := `
		SELECT tracking_link_id, organization_id, campaign_id, affiliate_id, name, description, status,
			   tracking_url, source_id, sub1, sub2, sub3, sub4, sub5,
			   is_encrypt_parameters, is_redirect_link,
			   internal_notes, tags,
			   created_at, updated_at
		FROM public.tracking_links
		WHERE campaign_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, campaignID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tracking links by campaign: %w", err)
	}
	defer rows.Close()

	var trackingLinks []*domain.TrackingLink
	for rows.Next() {
		trackingLink := &domain.TrackingLink{}
		err := rows.Scan(
			&trackingLink.TrackingLinkID,
			&trackingLink.OrganizationID,
			&trackingLink.CampaignID,
			&trackingLink.AffiliateID,
			&trackingLink.Name,
			&trackingLink.Description,
			&trackingLink.Status,
			&trackingLink.TrackingURL,
			&trackingLink.SourceID,
			&trackingLink.Sub1,
			&trackingLink.Sub2,
			&trackingLink.Sub3,
			&trackingLink.Sub4,
			&trackingLink.Sub5,
			&trackingLink.IsEncryptParameters,
			&trackingLink.IsRedirectLink,
			&trackingLink.InternalNotes,
			&trackingLink.Tags,
			&trackingLink.CreatedAt,
			&trackingLink.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tracking link: %w", err)
		}
		trackingLinks = append(trackingLinks, trackingLink)
	}

	return trackingLinks, nil
}

// ListTrackingLinksByAffiliate retrieves tracking links for a specific affiliate
func (r *trackingLinkRepository) ListTrackingLinksByAffiliate(ctx context.Context, affiliateID int64, limit, offset int) ([]*domain.TrackingLink, error) {
	query := `
		SELECT tracking_link_id, organization_id, campaign_id, affiliate_id, name, description, status,
			   tracking_url, source_id, sub1, sub2, sub3, sub4, sub5,
			   is_encrypt_parameters, is_redirect_link,
			   
			   internal_notes, tags,
			   created_at, updated_at
		FROM public.tracking_links
		WHERE affiliate_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, affiliateID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tracking links by affiliate: %w", err)
	}
	defer rows.Close()

	var trackingLinks []*domain.TrackingLink
	for rows.Next() {
		trackingLink := &domain.TrackingLink{}
		err := rows.Scan(
			&trackingLink.TrackingLinkID,
			&trackingLink.OrganizationID,
			&trackingLink.CampaignID,
			&trackingLink.AffiliateID,
			&trackingLink.Name,
			&trackingLink.Description,
			&trackingLink.Status,
			&trackingLink.TrackingURL,
			&trackingLink.SourceID,
			&trackingLink.Sub1,
			&trackingLink.Sub2,
			&trackingLink.Sub3,
			&trackingLink.Sub4,
			&trackingLink.Sub5,
			&trackingLink.IsEncryptParameters,
			&trackingLink.IsRedirectLink,
			&trackingLink.InternalNotes,
			&trackingLink.Tags,
			&trackingLink.CreatedAt,
			&trackingLink.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tracking link: %w", err)
		}
		trackingLinks = append(trackingLinks, trackingLink)
	}

	return trackingLinks, nil
}

// ListTrackingLinksByOrganization retrieves tracking links for a specific organization
func (r *trackingLinkRepository) ListTrackingLinksByOrganization(ctx context.Context, organizationID int64, limit, offset int) ([]*domain.TrackingLink, error) {
	query := `
		SELECT tracking_link_id, organization_id, campaign_id, affiliate_id, name, description, status,
			   tracking_url, source_id, sub1, sub2, sub3, sub4, sub5,
			   is_encrypt_parameters, is_redirect_link,
			   
			   internal_notes, tags,
			   created_at, updated_at
		FROM public.tracking_links
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tracking links by organization: %w", err)
	}
	defer rows.Close()

	var trackingLinks []*domain.TrackingLink
	for rows.Next() {
		trackingLink := &domain.TrackingLink{}
		err := rows.Scan(
			&trackingLink.TrackingLinkID,
			&trackingLink.OrganizationID,
			&trackingLink.CampaignID,
			&trackingLink.AffiliateID,
			&trackingLink.Name,
			&trackingLink.Description,
			&trackingLink.Status,
			&trackingLink.TrackingURL,
			&trackingLink.SourceID,
			&trackingLink.Sub1,
			&trackingLink.Sub2,
			&trackingLink.Sub3,
			&trackingLink.Sub4,
			&trackingLink.Sub5,
			&trackingLink.IsEncryptParameters,
			&trackingLink.IsRedirectLink,
			&trackingLink.InternalNotes,
			&trackingLink.Tags,
			&trackingLink.CreatedAt,
			&trackingLink.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tracking link: %w", err)
		}
		trackingLinks = append(trackingLinks, trackingLink)
	}

	return trackingLinks, nil
}

// GetTrackingLinkByCampaignAndAffiliate retrieves a tracking link by campaign, affiliate, and tracking parameters
func (r *trackingLinkRepository) GetTrackingLinkByCampaignAndAffiliate(ctx context.Context, campaignID, affiliateID int64, sourceID, sub1, sub2, sub3, sub4, sub5 *string) (*domain.TrackingLink, error) {
	query := `
		SELECT tracking_link_id, organization_id, campaign_id, affiliate_id, name, description, status,
			   tracking_url, source_id, sub1, sub2, sub3, sub4, sub5,
			   is_encrypt_parameters, is_redirect_link,
			   
			   internal_notes, tags,
			   created_at, updated_at
		FROM public.tracking_links
		WHERE campaign_id = $1 AND affiliate_id = $2 
		  AND (source_id = $3 OR (source_id IS NULL AND $3 IS NULL))
		  AND (sub1 = $4 OR (sub1 IS NULL AND $4 IS NULL))
		  AND (sub2 = $5 OR (sub2 IS NULL AND $5 IS NULL))
		  AND (sub3 = $6 OR (sub3 IS NULL AND $6 IS NULL))
		  AND (sub4 = $7 OR (sub4 IS NULL AND $7 IS NULL))
		  AND (sub5 = $8 OR (sub5 IS NULL AND $8 IS NULL))`

	trackingLink := &domain.TrackingLink{}
	err := r.db.QueryRow(ctx, query, campaignID, affiliateID, sourceID, sub1, sub2, sub3, sub4, sub5).Scan(
		&trackingLink.TrackingLinkID,
		&trackingLink.OrganizationID,
		&trackingLink.CampaignID,
		&trackingLink.AffiliateID,
		&trackingLink.Name,
		&trackingLink.Description,
		&trackingLink.Status,
		&trackingLink.TrackingURL,
		&trackingLink.SourceID,
		&trackingLink.Sub1,
		&trackingLink.Sub2,
		&trackingLink.Sub3,
		&trackingLink.Sub4,
		&trackingLink.Sub5,
		&trackingLink.IsEncryptParameters,
		&trackingLink.IsRedirectLink,
		&trackingLink.InternalNotes,
		&trackingLink.Tags,
		&trackingLink.CreatedAt,
		&trackingLink.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tracking link not found")
		}
		return nil, fmt.Errorf("failed to get tracking link: %w", err)
	}

	return trackingLink, nil
}

