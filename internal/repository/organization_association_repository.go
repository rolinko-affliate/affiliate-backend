package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OrganizationAssociationRepository defines the interface for organization association operations
type OrganizationAssociationRepository interface {
	CreateAssociation(ctx context.Context, association *domain.OrganizationAssociation) error
	GetAssociationByID(ctx context.Context, id int64) (*domain.OrganizationAssociation, error)
	GetAssociationByIDWithDetails(ctx context.Context, id int64) (*domain.OrganizationAssociationWithDetails, error)
	GetAssociationByOrganizations(ctx context.Context, advertiserOrgID, affiliateOrgID int64) (*domain.OrganizationAssociation, error)
	UpdateAssociation(ctx context.Context, association *domain.OrganizationAssociation) error
	ListAssociations(ctx context.Context, filter *domain.AssociationListFilter) ([]*domain.OrganizationAssociation, error)
	ListAssociationsWithDetails(ctx context.Context, filter *domain.AssociationListFilter) ([]*domain.OrganizationAssociationWithDetails, error)
	DeleteAssociation(ctx context.Context, id int64) error
	GetAssociationsByAdvertiserOrg(ctx context.Context, advertiserOrgID int64, status *domain.AssociationStatus) ([]*domain.OrganizationAssociation, error)
	GetAssociationsByAffiliateOrg(ctx context.Context, affiliateOrgID int64, status *domain.AssociationStatus) ([]*domain.OrganizationAssociation, error)
	CountAssociations(ctx context.Context, filter *domain.AssociationListFilter) (int64, error)
}

// pgxOrganizationAssociationRepository implements OrganizationAssociationRepository using pgx
type pgxOrganizationAssociationRepository struct {
	db *pgxpool.Pool
}

// NewPgxOrganizationAssociationRepository creates a new organization association repository
func NewPgxOrganizationAssociationRepository(db *pgxpool.Pool) OrganizationAssociationRepository {
	return &pgxOrganizationAssociationRepository{db: db}
}

// CreateAssociation creates a new organization association in the database
func (r *pgxOrganizationAssociationRepository) CreateAssociation(ctx context.Context, association *domain.OrganizationAssociation) error {
	query := `INSERT INTO public.organization_associations (
		advertiser_org_id, affiliate_org_id, status, association_type,
		visible_affiliate_ids, visible_campaign_ids, all_affiliates_visible, all_campaigns_visible,
		requested_by_user_id, message, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING association_id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		association.AdvertiserOrgID,
		association.AffiliateOrgID,
		association.Status,
		association.AssociationType,
		association.VisibleAffiliateIDs,
		association.VisibleCampaignIDs,
		association.AllAffiliatesVisible,
		association.AllCampaignsVisible,
		association.RequestedByUserID,
		association.Message,
		now,
		now,
	).Scan(&association.AssociationID, &association.CreatedAt, &association.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating organization association: %w", err)
	}
	return nil
}

// GetAssociationByID retrieves an organization association by ID
func (r *pgxOrganizationAssociationRepository) GetAssociationByID(ctx context.Context, id int64) (*domain.OrganizationAssociation, error) {
	query := `SELECT association_id, advertiser_org_id, affiliate_org_id, status, association_type,
		visible_affiliate_ids, visible_campaign_ids, all_affiliates_visible, all_campaigns_visible,
		requested_by_user_id, approved_by_user_id, message, created_at, updated_at, approved_at
		FROM public.organization_associations WHERE association_id = $1`

	association := &domain.OrganizationAssociation{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&association.AssociationID,
		&association.AdvertiserOrgID,
		&association.AffiliateOrgID,
		&association.Status,
		&association.AssociationType,
		&association.VisibleAffiliateIDs,
		&association.VisibleCampaignIDs,
		&association.AllAffiliatesVisible,
		&association.AllCampaignsVisible,
		&association.RequestedByUserID,
		&association.ApprovedByUserID,
		&association.Message,
		&association.CreatedAt,
		&association.UpdatedAt,
		&association.ApprovedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("organization association not found")
		}
		return nil, fmt.Errorf("error getting organization association: %w", err)
	}

	return association, nil
}

// GetAssociationByIDWithDetails retrieves an association by ID with organization and user details
func (r *pgxOrganizationAssociationRepository) GetAssociationByIDWithDetails(ctx context.Context, id int64) (*domain.OrganizationAssociationWithDetails, error) {
	query := `SELECT 
		oa.association_id, oa.advertiser_org_id, oa.affiliate_org_id, oa.status, oa.association_type,
		oa.visible_affiliate_ids, oa.visible_campaign_ids, oa.all_affiliates_visible, oa.all_campaigns_visible,
		oa.requested_by_user_id, oa.approved_by_user_id, oa.message, oa.created_at, oa.updated_at, oa.approved_at,
		adv_org.name as advertiser_name, adv_org.type as advertiser_type,
		aff_org.name as affiliate_name, aff_org.type as affiliate_type,
		req_user.first_name as req_first_name, req_user.last_name as req_last_name, req_user.email as req_email,
		app_user.first_name as app_first_name, app_user.last_name as app_last_name, app_user.email as app_email
		FROM public.organization_associations oa
		LEFT JOIN public.organizations adv_org ON oa.advertiser_org_id = adv_org.organization_id
		LEFT JOIN public.organizations aff_org ON oa.affiliate_org_id = aff_org.organization_id
		LEFT JOIN public.profiles req_user ON oa.requested_by_user_id = req_user.id
		LEFT JOIN public.profiles app_user ON oa.approved_by_user_id = app_user.id
		WHERE oa.association_id = $1`

	association := &domain.OrganizationAssociation{}
	advOrg := &domain.Organization{}
	affOrg := &domain.Organization{}
	var reqFirstName, reqLastName, reqEmail *string
	var appFirstName, appLastName, appEmail *string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&association.AssociationID,
		&association.AdvertiserOrgID,
		&association.AffiliateOrgID,
		&association.Status,
		&association.AssociationType,
		&association.VisibleAffiliateIDs,
		&association.VisibleCampaignIDs,
		&association.AllAffiliatesVisible,
		&association.AllCampaignsVisible,
		&association.RequestedByUserID,
		&association.ApprovedByUserID,
		&association.Message,
		&association.CreatedAt,
		&association.UpdatedAt,
		&association.ApprovedAt,
		&advOrg.Name,
		&advOrg.Type,
		&affOrg.Name,
		&affOrg.Type,
		&reqFirstName,
		&reqLastName,
		&reqEmail,
		&appFirstName,
		&appLastName,
		&appEmail,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("organization association not found")
		}
		return nil, fmt.Errorf("error getting organization association with details: %w", err)
	}

	// Set organization IDs
	advOrg.OrganizationID = association.AdvertiserOrgID
	affOrg.OrganizationID = association.AffiliateOrgID

	// Create the detailed association
	detailedAssociation := &domain.OrganizationAssociationWithDetails{
		OrganizationAssociation: association,
		AdvertiserOrganization:  advOrg,
		AffiliateOrganization:   affOrg,
	}

	// Add user details if available
	if association.RequestedByUserID != nil && (reqFirstName != nil || reqLastName != nil || reqEmail != nil) {
		reqUserID, _ := uuid.Parse(*association.RequestedByUserID)
		detailedAssociation.RequestedByUser = &domain.Profile{
			ID:        reqUserID,
			FirstName: reqFirstName,
			LastName:  reqLastName,
			Email:     *reqEmail,
		}
	}

	if association.ApprovedByUserID != nil && (appFirstName != nil || appLastName != nil || appEmail != nil) {
		appUserID, _ := uuid.Parse(*association.ApprovedByUserID)
		detailedAssociation.ApprovedByUser = &domain.Profile{
			ID:        appUserID,
			FirstName: appFirstName,
			LastName:  appLastName,
			Email:     *appEmail,
		}
	}

	return detailedAssociation, nil
}

// GetAssociationByOrganizations retrieves an association by advertiser and affiliate organization IDs
func (r *pgxOrganizationAssociationRepository) GetAssociationByOrganizations(ctx context.Context, advertiserOrgID, affiliateOrgID int64) (*domain.OrganizationAssociation, error) {
	query := `SELECT association_id, advertiser_org_id, affiliate_org_id, status, association_type,
		visible_affiliate_ids, visible_campaign_ids, all_affiliates_visible, all_campaigns_visible,
		requested_by_user_id, approved_by_user_id, message, created_at, updated_at, approved_at
		FROM public.organization_associations 
		WHERE advertiser_org_id = $1 AND affiliate_org_id = $2`

	association := &domain.OrganizationAssociation{}
	err := r.db.QueryRow(ctx, query, advertiserOrgID, affiliateOrgID).Scan(
		&association.AssociationID,
		&association.AdvertiserOrgID,
		&association.AffiliateOrgID,
		&association.Status,
		&association.AssociationType,
		&association.VisibleAffiliateIDs,
		&association.VisibleCampaignIDs,
		&association.AllAffiliatesVisible,
		&association.AllCampaignsVisible,
		&association.RequestedByUserID,
		&association.ApprovedByUserID,
		&association.Message,
		&association.CreatedAt,
		&association.UpdatedAt,
		&association.ApprovedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("organization association not found")
		}
		return nil, fmt.Errorf("error getting organization association: %w", err)
	}

	return association, nil
}

// UpdateAssociation updates an existing organization association
func (r *pgxOrganizationAssociationRepository) UpdateAssociation(ctx context.Context, association *domain.OrganizationAssociation) error {
	query := `UPDATE public.organization_associations SET
		status = $2, visible_affiliate_ids = $3, visible_campaign_ids = $4,
		all_affiliates_visible = $5, all_campaigns_visible = $6,
		approved_by_user_id = $7, approved_at = $8, updated_at = $9
		WHERE association_id = $1`

	now := time.Now()
	association.UpdatedAt = now

	_, err := r.db.Exec(ctx, query,
		association.AssociationID,
		association.Status,
		association.VisibleAffiliateIDs,
		association.VisibleCampaignIDs,
		association.AllAffiliatesVisible,
		association.AllCampaignsVisible,
		association.ApprovedByUserID,
		association.ApprovedAt,
		now,
	)

	if err != nil {
		return fmt.Errorf("error updating organization association: %w", err)
	}
	return nil
}

// ListAssociations retrieves a list of organization associations based on filter
func (r *pgxOrganizationAssociationRepository) ListAssociations(ctx context.Context, filter *domain.AssociationListFilter) ([]*domain.OrganizationAssociation, error) {
	query := `SELECT association_id, advertiser_org_id, affiliate_org_id, status, association_type,
		visible_affiliate_ids, visible_campaign_ids, all_affiliates_visible, all_campaigns_visible,
		requested_by_user_id, approved_by_user_id, message, created_at, updated_at, approved_at
		FROM public.organization_associations`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter != nil {
		if filter.AdvertiserOrgID != nil {
			conditions = append(conditions, fmt.Sprintf("advertiser_org_id = $%d", argIndex))
			args = append(args, *filter.AdvertiserOrgID)
			argIndex++
		}
		if filter.AffiliateOrgID != nil {
			conditions = append(conditions, fmt.Sprintf("affiliate_org_id = $%d", argIndex))
			args = append(args, *filter.AffiliateOrgID)
			argIndex++
		}
		if filter.Status != nil {
			conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
			args = append(args, *filter.Status)
			argIndex++
		}
		if filter.AssociationType != nil {
			conditions = append(conditions, fmt.Sprintf("association_type = $%d", argIndex))
			args = append(args, *filter.AssociationType)
			argIndex++
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	if filter != nil {
		if filter.Limit > 0 {
			query += fmt.Sprintf(" LIMIT $%d", argIndex)
			args = append(args, filter.Limit)
			argIndex++
		}
		if filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argIndex)
			args = append(args, filter.Offset)
		}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listing organization associations: %w", err)
	}
	defer rows.Close()

	var associations []*domain.OrganizationAssociation
	for rows.Next() {
		association := &domain.OrganizationAssociation{}
		err := rows.Scan(
			&association.AssociationID,
			&association.AdvertiserOrgID,
			&association.AffiliateOrgID,
			&association.Status,
			&association.AssociationType,
			&association.VisibleAffiliateIDs,
			&association.VisibleCampaignIDs,
			&association.AllAffiliatesVisible,
			&association.AllCampaignsVisible,
			&association.RequestedByUserID,
			&association.ApprovedByUserID,
			&association.Message,
			&association.CreatedAt,
			&association.UpdatedAt,
			&association.ApprovedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning organization association: %w", err)
		}
		associations = append(associations, association)
	}

	return associations, nil
}

// ListAssociationsWithDetails retrieves associations with organization and user details
func (r *pgxOrganizationAssociationRepository) ListAssociationsWithDetails(ctx context.Context, filter *domain.AssociationListFilter) ([]*domain.OrganizationAssociationWithDetails, error) {
	query := `SELECT 
		oa.association_id, oa.advertiser_org_id, oa.affiliate_org_id, oa.status, oa.association_type,
		oa.visible_affiliate_ids, oa.visible_campaign_ids, oa.all_affiliates_visible, oa.all_campaigns_visible,
		oa.requested_by_user_id, oa.approved_by_user_id, oa.message, oa.created_at, oa.updated_at, oa.approved_at,
		adv_org.name as advertiser_name, adv_org.type as advertiser_type,
		aff_org.name as affiliate_name, aff_org.type as affiliate_type,
		req_user.first_name as req_first_name, req_user.last_name as req_last_name, req_user.email as req_email,
		app_user.first_name as app_first_name, app_user.last_name as app_last_name, app_user.email as app_email
		FROM public.organization_associations oa
		LEFT JOIN public.organizations adv_org ON oa.advertiser_org_id = adv_org.organization_id
		LEFT JOIN public.organizations aff_org ON oa.affiliate_org_id = aff_org.organization_id
		LEFT JOIN public.profiles req_user ON oa.requested_by_user_id = req_user.id
		LEFT JOIN public.profiles app_user ON oa.approved_by_user_id = app_user.id`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter != nil {
		if filter.AdvertiserOrgID != nil {
			conditions = append(conditions, fmt.Sprintf("oa.advertiser_org_id = $%d", argIndex))
			args = append(args, *filter.AdvertiserOrgID)
			argIndex++
		}
		if filter.AffiliateOrgID != nil {
			conditions = append(conditions, fmt.Sprintf("oa.affiliate_org_id = $%d", argIndex))
			args = append(args, *filter.AffiliateOrgID)
			argIndex++
		}
		if filter.Status != nil {
			conditions = append(conditions, fmt.Sprintf("oa.status = $%d", argIndex))
			args = append(args, *filter.Status)
			argIndex++
		}
		if filter.AssociationType != nil {
			conditions = append(conditions, fmt.Sprintf("oa.association_type = $%d", argIndex))
			args = append(args, *filter.AssociationType)
			argIndex++
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY oa.created_at DESC"

	if filter != nil {
		if filter.Limit > 0 {
			query += fmt.Sprintf(" LIMIT $%d", argIndex)
			args = append(args, filter.Limit)
			argIndex++
		}
		if filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argIndex)
			args = append(args, filter.Offset)
		}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listing organization associations with details: %w", err)
	}
	defer rows.Close()

	var associations []*domain.OrganizationAssociationWithDetails
	for rows.Next() {
		association := &domain.OrganizationAssociation{}
		advOrg := &domain.Organization{}
		affOrg := &domain.Organization{}
		var reqFirstName, reqLastName, reqEmail *string
		var appFirstName, appLastName, appEmail *string

		err := rows.Scan(
			&association.AssociationID,
			&association.AdvertiserOrgID,
			&association.AffiliateOrgID,
			&association.Status,
			&association.AssociationType,
			&association.VisibleAffiliateIDs,
			&association.VisibleCampaignIDs,
			&association.AllAffiliatesVisible,
			&association.AllCampaignsVisible,
			&association.RequestedByUserID,
			&association.ApprovedByUserID,
			&association.Message,
			&association.CreatedAt,
			&association.UpdatedAt,
			&association.ApprovedAt,
			&advOrg.Name,
			&advOrg.Type,
			&affOrg.Name,
			&affOrg.Type,
			&reqFirstName,
			&reqLastName,
			&reqEmail,
			&appFirstName,
			&appLastName,
			&appEmail,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning organization association with details: %w", err)
		}

		// Set organization IDs
		advOrg.OrganizationID = association.AdvertiserOrgID
		affOrg.OrganizationID = association.AffiliateOrgID

		// Create the detailed association
		detailedAssociation := &domain.OrganizationAssociationWithDetails{
			OrganizationAssociation: association,
			AdvertiserOrganization:  advOrg,
			AffiliateOrganization:   affOrg,
		}

		// Add user details if available
		if reqFirstName != nil || reqLastName != nil || reqEmail != nil {
			reqUserID, _ := uuid.Parse(*association.RequestedByUserID)
			detailedAssociation.RequestedByUser = &domain.Profile{
				ID:        reqUserID,
				FirstName: reqFirstName,
				LastName:  reqLastName,
				Email:     *reqEmail,
			}
		}

		if appFirstName != nil || appLastName != nil || appEmail != nil {
			appUserID, _ := uuid.Parse(*association.ApprovedByUserID)
			detailedAssociation.ApprovedByUser = &domain.Profile{
				ID:        appUserID,
				FirstName: appFirstName,
				LastName:  appLastName,
				Email:     *appEmail,
			}
		}

		associations = append(associations, detailedAssociation)
	}

	return associations, nil
}

// DeleteAssociation deletes an organization association
func (r *pgxOrganizationAssociationRepository) DeleteAssociation(ctx context.Context, id int64) error {
	query := `DELETE FROM public.organization_associations WHERE association_id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting organization association: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("organization association not found")
	}

	return nil
}

// GetAssociationsByAdvertiserOrg retrieves associations for a specific advertiser organization
func (r *pgxOrganizationAssociationRepository) GetAssociationsByAdvertiserOrg(ctx context.Context, advertiserOrgID int64, status *domain.AssociationStatus) ([]*domain.OrganizationAssociation, error) {
	filter := &domain.AssociationListFilter{
		AdvertiserOrgID: &advertiserOrgID,
		Status:          status,
	}
	return r.ListAssociations(ctx, filter)
}

// GetAssociationsByAffiliateOrg retrieves associations for a specific affiliate organization
func (r *pgxOrganizationAssociationRepository) GetAssociationsByAffiliateOrg(ctx context.Context, affiliateOrgID int64, status *domain.AssociationStatus) ([]*domain.OrganizationAssociation, error) {
	filter := &domain.AssociationListFilter{
		AffiliateOrgID: &affiliateOrgID,
		Status:         status,
	}
	return r.ListAssociations(ctx, filter)
}

// CountAssociations counts the total number of associations matching the filter
func (r *pgxOrganizationAssociationRepository) CountAssociations(ctx context.Context, filter *domain.AssociationListFilter) (int64, error) {
	query := `SELECT COUNT(*) FROM public.organization_associations`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter != nil {
		if filter.AdvertiserOrgID != nil {
			conditions = append(conditions, fmt.Sprintf("advertiser_org_id = $%d", argIndex))
			args = append(args, *filter.AdvertiserOrgID)
			argIndex++
		}
		if filter.AffiliateOrgID != nil {
			conditions = append(conditions, fmt.Sprintf("affiliate_org_id = $%d", argIndex))
			args = append(args, *filter.AffiliateOrgID)
			argIndex++
		}
		if filter.Status != nil {
			conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
			args = append(args, *filter.Status)
			argIndex++
		}
		if filter.AssociationType != nil {
			conditions = append(conditions, fmt.Sprintf("association_type = $%d", argIndex))
			args = append(args, *filter.AssociationType)
			argIndex++
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting organization associations: %w", err)
	}

	return count, nil
}

// Helper function to convert slice of int64 to JSON string
func int64SliceToJSON(slice []int64) (*string, error) {
	if len(slice) == 0 {
		return nil, nil
	}
	
	jsonBytes, err := json.Marshal(slice)
	if err != nil {
		return nil, err
	}
	
	jsonStr := string(jsonBytes)
	return &jsonStr, nil
}

// Helper function to convert JSON string to slice of int64
func jsonToInt64Slice(jsonStr *string) ([]int64, error) {
	if jsonStr == nil || *jsonStr == "" {
		return nil, nil
	}
	
	var slice []int64
	err := json.Unmarshal([]byte(*jsonStr), &slice)
	if err != nil {
		return nil, err
	}
	
	return slice, nil
}