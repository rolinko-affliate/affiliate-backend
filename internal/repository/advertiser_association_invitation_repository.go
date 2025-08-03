package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AdvertiserAssociationInvitationRepository defines the interface for invitation operations
type AdvertiserAssociationInvitationRepository interface {
	// Invitation CRUD operations
	CreateInvitation(ctx context.Context, invitation *domain.AdvertiserAssociationInvitation) error
	GetInvitationByID(ctx context.Context, id int64) (*domain.AdvertiserAssociationInvitation, error)
	GetInvitationByToken(ctx context.Context, token string) (*domain.AdvertiserAssociationInvitation, error)
	GetInvitationByIDWithDetails(ctx context.Context, id int64) (*domain.AdvertiserAssociationInvitationWithDetails, error)
	GetInvitationByTokenWithDetails(ctx context.Context, token string) (*domain.AdvertiserAssociationInvitationWithDetails, error)
	UpdateInvitation(ctx context.Context, invitation *domain.AdvertiserAssociationInvitation) error
	DeleteInvitation(ctx context.Context, id int64) error
	ListInvitations(ctx context.Context, filter *domain.InvitationListFilter) ([]*domain.AdvertiserAssociationInvitation, error)
	ListInvitationsWithDetails(ctx context.Context, filter *domain.InvitationListFilter) ([]*domain.AdvertiserAssociationInvitationWithDetails, error)
	CountInvitations(ctx context.Context, filter *domain.InvitationListFilter) (int64, error)
	
	// Usage tracking
	LogInvitationUsage(ctx context.Context, usage *domain.InvitationUsageLog) error
	GetInvitationUsageHistory(ctx context.Context, invitationID int64, limit int) ([]*domain.InvitationUsageLog, error)
	IncrementInvitationUsage(ctx context.Context, invitationID int64) error
	
	// Utility methods
	GetInvitationsByAdvertiser(ctx context.Context, advertiserOrgID int64, status *domain.InvitationStatus) ([]*domain.AdvertiserAssociationInvitation, error)
	ExpireInvitations(ctx context.Context) (int64, error) // Returns number of expired invitations
}

// pgxAdvertiserAssociationInvitationRepository implements AdvertiserAssociationInvitationRepository using pgx
type pgxAdvertiserAssociationInvitationRepository struct {
	db *pgxpool.Pool
}

// NewPgxAdvertiserAssociationInvitationRepository creates a new invitation repository
func NewPgxAdvertiserAssociationInvitationRepository(db *pgxpool.Pool) AdvertiserAssociationInvitationRepository {
	return &pgxAdvertiserAssociationInvitationRepository{db: db}
}

// CreateInvitation creates a new invitation in the database
func (r *pgxAdvertiserAssociationInvitationRepository) CreateInvitation(ctx context.Context, invitation *domain.AdvertiserAssociationInvitation) error {
	query := `INSERT INTO public.advertiser_association_invitations (
		advertiser_org_id, invitation_token, name, description, allowed_affiliate_org_ids,
		max_uses, current_uses, expires_at, status, created_by_user_id, message,
		default_all_affiliates_visible, default_all_campaigns_visible,
		default_visible_affiliate_ids, default_visible_campaign_ids,
		created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	RETURNING invitation_id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		invitation.AdvertiserOrgID,
		invitation.InvitationToken,
		invitation.Name,
		invitation.Description,
		invitation.AllowedAffiliateOrgIDs,
		invitation.MaxUses,
		invitation.CurrentUses,
		invitation.ExpiresAt,
		invitation.Status,
		invitation.CreatedByUserID,
		invitation.Message,
		invitation.DefaultAllAffiliatesVisible,
		invitation.DefaultAllCampaignsVisible,
		invitation.DefaultVisibleAffiliateIDs,
		invitation.DefaultVisibleCampaignIDs,
		now,
		now,
	).Scan(&invitation.InvitationID, &invitation.CreatedAt, &invitation.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating invitation: %w", err)
	}
	return nil
}

// GetInvitationByID retrieves an invitation by ID
func (r *pgxAdvertiserAssociationInvitationRepository) GetInvitationByID(ctx context.Context, id int64) (*domain.AdvertiserAssociationInvitation, error) {
	query := `SELECT invitation_id, advertiser_org_id, invitation_token, name, description,
		allowed_affiliate_org_ids, max_uses, current_uses, expires_at, status,
		created_by_user_id, message, default_all_affiliates_visible, default_all_campaigns_visible,
		default_visible_affiliate_ids, default_visible_campaign_ids, created_at, updated_at
		FROM public.advertiser_association_invitations WHERE invitation_id = $1`

	invitation := &domain.AdvertiserAssociationInvitation{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&invitation.InvitationID,
		&invitation.AdvertiserOrgID,
		&invitation.InvitationToken,
		&invitation.Name,
		&invitation.Description,
		&invitation.AllowedAffiliateOrgIDs,
		&invitation.MaxUses,
		&invitation.CurrentUses,
		&invitation.ExpiresAt,
		&invitation.Status,
		&invitation.CreatedByUserID,
		&invitation.Message,
		&invitation.DefaultAllAffiliatesVisible,
		&invitation.DefaultAllCampaignsVisible,
		&invitation.DefaultVisibleAffiliateIDs,
		&invitation.DefaultVisibleCampaignIDs,
		&invitation.CreatedAt,
		&invitation.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, fmt.Errorf("error getting invitation: %w", err)
	}
	return invitation, nil
}

// GetInvitationByToken retrieves an invitation by token
func (r *pgxAdvertiserAssociationInvitationRepository) GetInvitationByToken(ctx context.Context, token string) (*domain.AdvertiserAssociationInvitation, error) {
	query := `SELECT invitation_id, advertiser_org_id, invitation_token, name, description,
		allowed_affiliate_org_ids, max_uses, current_uses, expires_at, status,
		created_by_user_id, message, default_all_affiliates_visible, default_all_campaigns_visible,
		default_visible_affiliate_ids, default_visible_campaign_ids, created_at, updated_at
		FROM public.advertiser_association_invitations WHERE invitation_token = $1`

	invitation := &domain.AdvertiserAssociationInvitation{}
	err := r.db.QueryRow(ctx, query, token).Scan(
		&invitation.InvitationID,
		&invitation.AdvertiserOrgID,
		&invitation.InvitationToken,
		&invitation.Name,
		&invitation.Description,
		&invitation.AllowedAffiliateOrgIDs,
		&invitation.MaxUses,
		&invitation.CurrentUses,
		&invitation.ExpiresAt,
		&invitation.Status,
		&invitation.CreatedByUserID,
		&invitation.Message,
		&invitation.DefaultAllAffiliatesVisible,
		&invitation.DefaultAllCampaignsVisible,
		&invitation.DefaultVisibleAffiliateIDs,
		&invitation.DefaultVisibleCampaignIDs,
		&invitation.CreatedAt,
		&invitation.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, fmt.Errorf("error getting invitation: %w", err)
	}
	return invitation, nil
}

// GetInvitationByIDWithDetails retrieves an invitation by ID with additional details
func (r *pgxAdvertiserAssociationInvitationRepository) GetInvitationByIDWithDetails(ctx context.Context, id int64) (*domain.AdvertiserAssociationInvitationWithDetails, error) {
	query := `SELECT 
		i.invitation_id, i.advertiser_org_id, i.invitation_token, i.name, i.description,
		i.allowed_affiliate_org_ids, i.max_uses, i.current_uses, i.expires_at, i.status,
		i.created_by_user_id, i.message, i.default_all_affiliates_visible, i.default_all_campaigns_visible,
		i.default_visible_affiliate_ids, i.default_visible_campaign_ids, i.created_at, i.updated_at,
		o.organization_id, o.name, o.type, o.created_at, o.updated_at,
		p.id, p.email, p.first_name, p.last_name, p.role_id, p.created_at, p.updated_at
		FROM public.advertiser_association_invitations i
		LEFT JOIN public.organizations o ON i.advertiser_org_id = o.organization_id
		LEFT JOIN public.profiles p ON i.created_by_user_id = p.id
		WHERE i.invitation_id = $1`

	invitation := &domain.AdvertiserAssociationInvitation{}
	org := &domain.Organization{}
	profile := &domain.Profile{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&invitation.InvitationID,
		&invitation.AdvertiserOrgID,
		&invitation.InvitationToken,
		&invitation.Name,
		&invitation.Description,
		&invitation.AllowedAffiliateOrgIDs,
		&invitation.MaxUses,
		&invitation.CurrentUses,
		&invitation.ExpiresAt,
		&invitation.Status,
		&invitation.CreatedByUserID,
		&invitation.Message,
		&invitation.DefaultAllAffiliatesVisible,
		&invitation.DefaultAllCampaignsVisible,
		&invitation.DefaultVisibleAffiliateIDs,
		&invitation.DefaultVisibleCampaignIDs,
		&invitation.CreatedAt,
		&invitation.UpdatedAt,
		// Organization fields
		&org.OrganizationID,
		&org.Name,
		&org.Type,
		&org.CreatedAt,
		&org.UpdatedAt,
		// Profile fields
		&profile.ID,
		&profile.Email,
		&profile.FirstName,
		&profile.LastName,
		&profile.RoleID,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, fmt.Errorf("error getting invitation with details: %w", err)
	}

	details := &domain.AdvertiserAssociationInvitationWithDetails{
		AdvertiserAssociationInvitation: invitation,
		AdvertiserOrganization:          org,
		CreatedByUser:                   profile,
		UsageCount:                      invitation.CurrentUses,
	}

	return details, nil
}

// GetInvitationByTokenWithDetails retrieves an invitation by token with additional details
func (r *pgxAdvertiserAssociationInvitationRepository) GetInvitationByTokenWithDetails(ctx context.Context, token string) (*domain.AdvertiserAssociationInvitationWithDetails, error) {
	query := `SELECT 
		i.invitation_id, i.advertiser_org_id, i.invitation_token, i.name, i.description,
		i.allowed_affiliate_org_ids, i.max_uses, i.current_uses, i.expires_at, i.status,
		i.created_by_user_id, i.message, i.default_all_affiliates_visible, i.default_all_campaigns_visible,
		i.default_visible_affiliate_ids, i.default_visible_campaign_ids, i.created_at, i.updated_at,
		o.organization_id, o.name, o.type, o.created_at, o.updated_at,
		p.id, p.email, p.first_name, p.last_name, p.role_id, p.created_at, p.updated_at
		FROM public.advertiser_association_invitations i
		LEFT JOIN public.organizations o ON i.advertiser_org_id = o.organization_id
		LEFT JOIN public.profiles p ON i.created_by_user_id = p.id
		WHERE i.invitation_token = $1`

	invitation := &domain.AdvertiserAssociationInvitation{}
	org := &domain.Organization{}
	profile := &domain.Profile{}

	err := r.db.QueryRow(ctx, query, token).Scan(
		&invitation.InvitationID,
		&invitation.AdvertiserOrgID,
		&invitation.InvitationToken,
		&invitation.Name,
		&invitation.Description,
		&invitation.AllowedAffiliateOrgIDs,
		&invitation.MaxUses,
		&invitation.CurrentUses,
		&invitation.ExpiresAt,
		&invitation.Status,
		&invitation.CreatedByUserID,
		&invitation.Message,
		&invitation.DefaultAllAffiliatesVisible,
		&invitation.DefaultAllCampaignsVisible,
		&invitation.DefaultVisibleAffiliateIDs,
		&invitation.DefaultVisibleCampaignIDs,
		&invitation.CreatedAt,
		&invitation.UpdatedAt,
		// Organization fields
		&org.OrganizationID,
		&org.Name,
		&org.Type,
		&org.CreatedAt,
		&org.UpdatedAt,
		// Profile fields
		&profile.ID,
		&profile.Email,
		&profile.FirstName,
		&profile.LastName,
		&profile.RoleID,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, fmt.Errorf("error getting invitation with details: %w", err)
	}

	details := &domain.AdvertiserAssociationInvitationWithDetails{
		AdvertiserAssociationInvitation: invitation,
		AdvertiserOrganization:          org,
		CreatedByUser:                   profile,
		UsageCount:                      invitation.CurrentUses,
	}

	return details, nil
}

// UpdateInvitation updates an existing invitation
func (r *pgxAdvertiserAssociationInvitationRepository) UpdateInvitation(ctx context.Context, invitation *domain.AdvertiserAssociationInvitation) error {
	query := `UPDATE public.advertiser_association_invitations SET
		name = $2, description = $3, allowed_affiliate_org_ids = $4, max_uses = $5,
		current_uses = $6, expires_at = $7, status = $8, message = $9,
		default_all_affiliates_visible = $10, default_all_campaigns_visible = $11,
		default_visible_affiliate_ids = $12, default_visible_campaign_ids = $13,
		updated_at = $14
		WHERE invitation_id = $1`

	now := time.Now()
	result, err := r.db.Exec(ctx, query,
		invitation.InvitationID,
		invitation.Name,
		invitation.Description,
		invitation.AllowedAffiliateOrgIDs,
		invitation.MaxUses,
		invitation.CurrentUses,
		invitation.ExpiresAt,
		invitation.Status,
		invitation.Message,
		invitation.DefaultAllAffiliatesVisible,
		invitation.DefaultAllCampaignsVisible,
		invitation.DefaultVisibleAffiliateIDs,
		invitation.DefaultVisibleCampaignIDs,
		now,
	)

	if err != nil {
		return fmt.Errorf("error updating invitation: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("invitation not found")
	}

	invitation.UpdatedAt = now
	return nil
}

// DeleteInvitation deletes an invitation
func (r *pgxAdvertiserAssociationInvitationRepository) DeleteInvitation(ctx context.Context, id int64) error {
	query := `DELETE FROM public.advertiser_association_invitations WHERE invitation_id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting invitation: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("invitation not found")
	}

	return nil
}

// ListInvitations lists invitations based on filter
func (r *pgxAdvertiserAssociationInvitationRepository) ListInvitations(ctx context.Context, filter *domain.InvitationListFilter) ([]*domain.AdvertiserAssociationInvitation, error) {
	query := `SELECT invitation_id, advertiser_org_id, invitation_token, name, description,
		allowed_affiliate_org_ids, max_uses, current_uses, expires_at, status,
		created_by_user_id, message, default_all_affiliates_visible, default_all_campaigns_visible,
		default_visible_affiliate_ids, default_visible_campaign_ids, created_at, updated_at
		FROM public.advertiser_association_invitations`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.AdvertiserOrgID != nil {
		conditions = append(conditions, fmt.Sprintf("advertiser_org_id = $%d", argIndex))
		args = append(args, *filter.AdvertiserOrgID)
		argIndex++
	}

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}

	if filter.CreatedByUserID != nil {
		conditions = append(conditions, fmt.Sprintf("created_by_user_id = $%d", argIndex))
		args = append(args, *filter.CreatedByUserID)
		argIndex++
	}

	if !filter.IncludeExpired {
		conditions = append(conditions, "(expires_at IS NULL OR expires_at > NOW())")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listing invitations: %w", err)
	}
	defer rows.Close()

	var invitations []*domain.AdvertiserAssociationInvitation
	for rows.Next() {
		invitation := &domain.AdvertiserAssociationInvitation{}
		err := rows.Scan(
			&invitation.InvitationID,
			&invitation.AdvertiserOrgID,
			&invitation.InvitationToken,
			&invitation.Name,
			&invitation.Description,
			&invitation.AllowedAffiliateOrgIDs,
			&invitation.MaxUses,
			&invitation.CurrentUses,
			&invitation.ExpiresAt,
			&invitation.Status,
			&invitation.CreatedByUserID,
			&invitation.Message,
			&invitation.DefaultAllAffiliatesVisible,
			&invitation.DefaultAllCampaignsVisible,
			&invitation.DefaultVisibleAffiliateIDs,
			&invitation.DefaultVisibleCampaignIDs,
			&invitation.CreatedAt,
			&invitation.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning invitation: %w", err)
		}
		invitations = append(invitations, invitation)
	}

	return invitations, nil
}

// ListInvitationsWithDetails lists invitations with additional details
func (r *pgxAdvertiserAssociationInvitationRepository) ListInvitationsWithDetails(ctx context.Context, filter *domain.InvitationListFilter) ([]*domain.AdvertiserAssociationInvitationWithDetails, error) {
	// Implementation would be similar to ListInvitations but with JOINs
	// For brevity, returning a simplified version
	invitations, err := r.ListInvitations(ctx, filter)
	if err != nil {
		return nil, err
	}

	var details []*domain.AdvertiserAssociationInvitationWithDetails
	for _, invitation := range invitations {
		detail := &domain.AdvertiserAssociationInvitationWithDetails{
			AdvertiserAssociationInvitation: invitation,
			UsageCount:                      invitation.CurrentUses,
		}
		details = append(details, detail)
	}

	return details, nil
}

// CountInvitations counts invitations based on filter
func (r *pgxAdvertiserAssociationInvitationRepository) CountInvitations(ctx context.Context, filter *domain.InvitationListFilter) (int64, error) {
	query := `SELECT COUNT(*) FROM public.advertiser_association_invitations`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.AdvertiserOrgID != nil {
		conditions = append(conditions, fmt.Sprintf("advertiser_org_id = $%d", argIndex))
		args = append(args, *filter.AdvertiserOrgID)
		argIndex++
	}

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}

	if filter.CreatedByUserID != nil {
		conditions = append(conditions, fmt.Sprintf("created_by_user_id = $%d", argIndex))
		args = append(args, *filter.CreatedByUserID)
		argIndex++
	}

	if !filter.IncludeExpired {
		conditions = append(conditions, "(expires_at IS NULL OR expires_at > NOW())")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting invitations: %w", err)
	}

	return count, nil
}

// LogInvitationUsage logs the usage of an invitation
func (r *pgxAdvertiserAssociationInvitationRepository) LogInvitationUsage(ctx context.Context, usage *domain.InvitationUsageLog) error {
	query := `INSERT INTO public.invitation_usage_log (
		invitation_id, affiliate_org_id, used_by_user_id, association_id,
		ip_address, user_agent, success, error_message, used_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING usage_id, used_at`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		usage.InvitationID,
		usage.AffiliateOrgID,
		usage.UsedByUserID,
		usage.AssociationID,
		usage.IPAddress,
		usage.UserAgent,
		usage.Success,
		usage.ErrorMessage,
		now,
	).Scan(&usage.UsageID, &usage.UsedAt)

	if err != nil {
		return fmt.Errorf("error logging invitation usage: %w", err)
	}
	return nil
}

// GetInvitationUsageHistory retrieves usage history for an invitation
func (r *pgxAdvertiserAssociationInvitationRepository) GetInvitationUsageHistory(ctx context.Context, invitationID int64, limit int) ([]*domain.InvitationUsageLog, error) {
	query := `SELECT usage_id, invitation_id, affiliate_org_id, used_by_user_id, association_id,
		ip_address, user_agent, success, error_message, used_at
		FROM public.invitation_usage_log 
		WHERE invitation_id = $1 
		ORDER BY used_at DESC`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := r.db.Query(ctx, query, invitationID)
	if err != nil {
		return nil, fmt.Errorf("error getting invitation usage history: %w", err)
	}
	defer rows.Close()

	var usages []*domain.InvitationUsageLog
	for rows.Next() {
		usage := &domain.InvitationUsageLog{}
		var ipAddress interface{}
		err := rows.Scan(
			&usage.UsageID,
			&usage.InvitationID,
			&usage.AffiliateOrgID,
			&usage.UsedByUserID,
			&usage.AssociationID,
			&ipAddress,
			&usage.UserAgent,
			&usage.Success,
			&usage.ErrorMessage,
			&usage.UsedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning usage log: %w", err)
		}
		
		// Handle inet type conversion
		if ipAddress != nil {
			ipStr := fmt.Sprintf("%v", ipAddress)
			usage.IPAddress = &ipStr
		}
		
		usages = append(usages, usage)
	}

	return usages, nil
}

// IncrementInvitationUsage increments the usage count for an invitation
func (r *pgxAdvertiserAssociationInvitationRepository) IncrementInvitationUsage(ctx context.Context, invitationID int64) error {
	query := `UPDATE public.advertiser_association_invitations 
		SET current_uses = current_uses + 1, updated_at = NOW() 
		WHERE invitation_id = $1`

	result, err := r.db.Exec(ctx, query, invitationID)
	if err != nil {
		return fmt.Errorf("error incrementing invitation usage: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("invitation not found")
	}

	return nil
}

// GetInvitationsByAdvertiser retrieves invitations for a specific advertiser
func (r *pgxAdvertiserAssociationInvitationRepository) GetInvitationsByAdvertiser(ctx context.Context, advertiserOrgID int64, status *domain.InvitationStatus) ([]*domain.AdvertiserAssociationInvitation, error) {
	filter := &domain.InvitationListFilter{
		AdvertiserOrgID: &advertiserOrgID,
		Status:          status,
	}
	return r.ListInvitations(ctx, filter)
}

// ExpireInvitations marks expired invitations as expired
func (r *pgxAdvertiserAssociationInvitationRepository) ExpireInvitations(ctx context.Context) (int64, error) {
	query := `UPDATE public.advertiser_association_invitations 
		SET status = 'expired', updated_at = NOW() 
		WHERE status = 'active' AND expires_at IS NOT NULL AND expires_at <= NOW()`

	result, err := r.db.Exec(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("error expiring invitations: %w", err)
	}

	return result.RowsAffected(), nil
}