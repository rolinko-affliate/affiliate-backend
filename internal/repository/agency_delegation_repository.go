package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/affiliate-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AgencyDelegationRepository handles database operations for agency delegations
type AgencyDelegationRepository interface {
	Create(ctx context.Context, delegation *domain.AgencyDelegation) (*domain.AgencyDelegation, error)
	GetByID(ctx context.Context, delegationID int64) (*domain.AgencyDelegation, error)
	GetByIDWithDetails(ctx context.Context, delegationID int64) (*domain.AgencyDelegationWithDetails, error)
	GetByOrganizations(ctx context.Context, agencyOrgID, advertiserOrgID int64) (*domain.AgencyDelegation, error)
	Update(ctx context.Context, delegation *domain.AgencyDelegation) (*domain.AgencyDelegation, error)
	Delete(ctx context.Context, delegationID int64) error
	List(ctx context.Context, filter domain.DelegationListFilter) ([]*domain.AgencyDelegation, error)
	ListWithDetails(ctx context.Context, filter domain.DelegationListFilter) ([]*domain.AgencyDelegationWithDetails, error)
	GetByAgencyOrg(ctx context.Context, agencyOrgID int64, filter domain.DelegationListFilter) ([]*domain.AgencyDelegation, error)
	GetByAdvertiserOrg(ctx context.Context, advertiserOrgID int64, filter domain.DelegationListFilter) ([]*domain.AgencyDelegation, error)
	HasPermission(ctx context.Context, agencyOrgID, advertiserOrgID int64, permission string) (bool, error)
	GetActivePermissions(ctx context.Context, agencyOrgID, advertiserOrgID int64) ([]string, error)
	CheckPermissions(ctx context.Context, agencyOrgID, advertiserOrgID int64, permissions []string) (map[string]bool, error)
	GetActiveDelegationsByAgency(ctx context.Context, agencyOrgID int64) ([]*domain.AgencyDelegation, error)
	GetActiveDelegationsByAdvertiser(ctx context.Context, advertiserOrgID int64) ([]*domain.AgencyDelegation, error)
	ExpireOldDelegations(ctx context.Context) (int64, error)
}

type agencyDelegationRepository struct {
	db *pgxpool.Pool
}

// NewPgxAgencyDelegationRepository creates a new agency delegation repository
func NewPgxAgencyDelegationRepository(db *pgxpool.Pool) AgencyDelegationRepository {
	return &agencyDelegationRepository{db: db}
}

// Create creates a new agency delegation
func (r *agencyDelegationRepository) Create(ctx context.Context, delegation *domain.AgencyDelegation) (*domain.AgencyDelegation, error) {
	query := `
		INSERT INTO agency_delegations (
			agency_org_id, advertiser_org_id, status, permissions,
			delegated_by_user_id, message, expires_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		) RETURNING delegation_id, created_at, updated_at`

	// Convert permissions to JSON string if provided
	var permissionsJSON string
	if delegation.Permissions != nil {
		permissionsJSON = *delegation.Permissions
	} else {
		permissionsJSON = "[]"
	}

	err := r.db.QueryRow(ctx, query,
		delegation.AgencyOrgID,
		delegation.AdvertiserOrgID,
		delegation.Status,
		permissionsJSON,
		delegation.DelegatedByUserID,
		delegation.Message,
		delegation.ExpiresAt,
	).Scan(&delegation.DelegationID, &delegation.CreatedAt, &delegation.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create agency delegation: %w", err)
	}

	return delegation, nil
}

// GetByID retrieves an agency delegation by ID
func (r *agencyDelegationRepository) GetByID(ctx context.Context, delegationID int64) (*domain.AgencyDelegation, error) {
	query := `
		SELECT delegation_id, agency_org_id, advertiser_org_id, status, permissions,
			   delegated_by_user_id, accepted_by_user_id, message, expires_at,
			   created_at, updated_at, accepted_at
		FROM agency_delegations
		WHERE delegation_id = $1`

	var delegation domain.AgencyDelegation

	err := r.db.QueryRow(ctx, query, delegationID).Scan(
		&delegation.DelegationID,
		&delegation.AgencyOrgID,
		&delegation.AdvertiserOrgID,
		&delegation.Status,
		&delegation.Permissions,
		&delegation.DelegatedByUserID,
		&delegation.AcceptedByUserID,
		&delegation.Message,
		&delegation.ExpiresAt,
		&delegation.CreatedAt,
		&delegation.UpdatedAt,
		&delegation.AcceptedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("agency delegation not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get agency delegation: %w", err)
	}

	return &delegation, nil
}

// GetByIDWithDetails retrieves an agency delegation with organization details
func (r *agencyDelegationRepository) GetByIDWithDetails(ctx context.Context, delegationID int64) (*domain.AgencyDelegationWithDetails, error) {
	// First get the basic delegation
	delegation, err := r.GetByID(ctx, delegationID)
	if err != nil {
		return nil, err
	}

	// Then get the organization and user details
	query := `
		SELECT ao.organization_id, ao.name, ao.type, ao.created_at, ao.updated_at,
			   ado.organization_id, ado.name, ado.type, ado.created_at, ado.updated_at,
			   dp.id, dp.first_name, dp.last_name, dp.email, dp.created_at, dp.updated_at,
			   ap.id, ap.first_name, ap.last_name, ap.email, ap.created_at, ap.updated_at
		FROM agency_delegations ad
		JOIN organizations ao ON ad.agency_org_id = ao.organization_id
		JOIN organizations ado ON ad.advertiser_org_id = ado.organization_id
		LEFT JOIN profiles dp ON ad.delegated_by_user_id = dp.id
		LEFT JOIN profiles ap ON ad.accepted_by_user_id = ap.id
		WHERE ad.delegation_id = $1`

	var agencyOrg domain.Organization
	var advertiserOrg domain.Organization
	var delegatedByUser, acceptedByUser sql.NullString
	var delegatedByProfile, acceptedByProfile domain.Profile

	err = r.db.QueryRow(ctx, query, delegationID).Scan(
		&agencyOrg.OrganizationID, &agencyOrg.Name, &agencyOrg.Type, &agencyOrg.CreatedAt, &agencyOrg.UpdatedAt,
		&advertiserOrg.OrganizationID, &advertiserOrg.Name, &advertiserOrg.Type, &advertiserOrg.CreatedAt, &advertiserOrg.UpdatedAt,
		&delegatedByUser, &delegatedByProfile.FirstName, &delegatedByProfile.LastName, &delegatedByProfile.Email, &delegatedByProfile.CreatedAt, &delegatedByProfile.UpdatedAt,
		&acceptedByUser, &acceptedByProfile.FirstName, &acceptedByProfile.LastName, &acceptedByProfile.Email, &acceptedByProfile.CreatedAt, &acceptedByProfile.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("agency delegation not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get agency delegation details: %w", err)
	}

	result := &domain.AgencyDelegationWithDetails{
		AgencyDelegation:       delegation,
		AgencyOrganization:     &agencyOrg,
		AdvertiserOrganization: &advertiserOrg,
	}

	if delegatedByUser.Valid {
		result.DelegatedByUser = &delegatedByProfile
	}
	if acceptedByUser.Valid {
		result.AcceptedByUser = &acceptedByProfile
	}

	return result, nil
}

// GetByOrganizations retrieves an agency delegation by organization IDs
func (r *agencyDelegationRepository) GetByOrganizations(ctx context.Context, agencyOrgID, advertiserOrgID int64) (*domain.AgencyDelegation, error) {
	query := `
		SELECT delegation_id, agency_org_id, advertiser_org_id, status, permissions,
			   delegated_by_user_id, accepted_by_user_id, message, expires_at,
			   created_at, updated_at, accepted_at
		FROM agency_delegations
		WHERE agency_org_id = $1 AND advertiser_org_id = $2`

	var delegation domain.AgencyDelegation

	err := r.db.QueryRow(ctx, query, agencyOrgID, advertiserOrgID).Scan(
		&delegation.DelegationID,
		&delegation.AgencyOrgID,
		&delegation.AdvertiserOrgID,
		&delegation.Status,
		&delegation.Permissions,
		&delegation.DelegatedByUserID,
		&delegation.AcceptedByUserID,
		&delegation.Message,
		&delegation.ExpiresAt,
		&delegation.CreatedAt,
		&delegation.UpdatedAt,
		&delegation.AcceptedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("agency delegation not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get agency delegation: %w", err)
	}

	return &delegation, nil
}

// Update updates an existing agency delegation
func (r *agencyDelegationRepository) Update(ctx context.Context, delegation *domain.AgencyDelegation) (*domain.AgencyDelegation, error) {
	query := `
		UPDATE agency_delegations
		SET status = $2, permissions = $3, accepted_by_user_id = $4, message = $5, 
			expires_at = $6, accepted_at = $7, updated_at = CURRENT_TIMESTAMP
		WHERE delegation_id = $1
		RETURNING updated_at`

	// Convert permissions to JSON string
	var permissionsJSON string
	if delegation.Permissions != nil {
		permissionsJSON = *delegation.Permissions
	} else {
		permissionsJSON = "[]"
	}

	err := r.db.QueryRow(ctx, query,
		delegation.DelegationID,
		delegation.Status,
		permissionsJSON,
		delegation.AcceptedByUserID,
		delegation.Message,
		delegation.ExpiresAt,
		delegation.AcceptedAt,
	).Scan(&delegation.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("agency delegation not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to update agency delegation: %w", err)
	}

	return delegation, nil
}

// Delete deletes an agency delegation
func (r *agencyDelegationRepository) Delete(ctx context.Context, delegationID int64) error {
	query := `DELETE FROM agency_delegations WHERE delegation_id = $1`

	result, err := r.db.Exec(ctx, query, delegationID)
	if err != nil {
		return fmt.Errorf("failed to delete agency delegation: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("agency delegation not found: %w", domain.ErrNotFound)
	}

	return nil
}

// List retrieves agency delegations with filtering
func (r *agencyDelegationRepository) List(ctx context.Context, filter domain.DelegationListFilter) ([]*domain.AgencyDelegation, error) {
	query := `
		SELECT delegation_id, agency_org_id, advertiser_org_id, status, permissions,
			   delegated_by_user_id, accepted_by_user_id, message, expires_at,
			   created_at, updated_at, accepted_at
		FROM agency_delegations`

	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if filter.AgencyOrgID != nil {
		conditions = append(conditions, fmt.Sprintf("agency_org_id = $%d", argIndex))
		args = append(args, *filter.AgencyOrgID)
		argIndex++
	}

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

	// Handle expired delegations filter
	if !filter.IncludeExpired {
		conditions = append(conditions, "(expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add ordering
	query += " ORDER BY created_at DESC"

	// Add pagination
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
		argIndex++
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list agency delegations: %w", err)
	}
	defer rows.Close()

	// Initialize with empty slice to ensure we never return nil
	delegations := make([]*domain.AgencyDelegation, 0)
	
	for rows.Next() {
		var delegation domain.AgencyDelegation

		err := rows.Scan(
			&delegation.DelegationID,
			&delegation.AgencyOrgID,
			&delegation.AdvertiserOrgID,
			&delegation.Status,
			&delegation.Permissions,
			&delegation.DelegatedByUserID,
			&delegation.AcceptedByUserID,
			&delegation.Message,
			&delegation.ExpiresAt,
			&delegation.CreatedAt,
			&delegation.UpdatedAt,
			&delegation.AcceptedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan agency delegation: %w", err)
		}

		delegations = append(delegations, &delegation)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over agency delegations: %w", err)
	}

	return delegations, nil
}

// ListWithDetails retrieves agency delegations with organization details
func (r *agencyDelegationRepository) ListWithDetails(ctx context.Context, filter domain.DelegationListFilter) ([]*domain.AgencyDelegationWithDetails, error) {
	// First get the basic delegations
	delegations, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Initialize with empty slice to ensure we never return nil
	result := make([]*domain.AgencyDelegationWithDetails, 0)
	
	// Then get details for each delegation
	for _, delegation := range delegations {
		details, err := r.GetByIDWithDetails(ctx, delegation.DelegationID)
		if err != nil {
			return nil, fmt.Errorf("failed to get details for delegation %d: %w", delegation.DelegationID, err)
		}
		result = append(result, details)
	}

	return result, nil
}

// GetByAgencyOrg retrieves delegations for a specific agency organization
func (r *agencyDelegationRepository) GetByAgencyOrg(ctx context.Context, agencyOrgID int64, filter domain.DelegationListFilter) ([]*domain.AgencyDelegation, error) {
	filter.AgencyOrgID = &agencyOrgID
	return r.List(ctx, filter)
}

// GetByAdvertiserOrg retrieves delegations for a specific advertiser organization
func (r *agencyDelegationRepository) GetByAdvertiserOrg(ctx context.Context, advertiserOrgID int64, filter domain.DelegationListFilter) ([]*domain.AgencyDelegation, error) {
	filter.AdvertiserOrgID = &advertiserOrgID
	return r.List(ctx, filter)
}

// HasPermission checks if an agency has a specific permission for an advertiser
func (r *agencyDelegationRepository) HasPermission(ctx context.Context, agencyOrgID, advertiserOrgID int64, permission string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM agency_delegations
			WHERE agency_org_id = $1 
			  AND advertiser_org_id = $2 
			  AND status = 'active'
			  AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
			  AND permissions @> $3
		)`

	permissionJSON, err := json.Marshal([]string{permission})
	if err != nil {
		return false, fmt.Errorf("failed to marshal permission: %w", err)
	}

	var exists bool
	err = r.db.QueryRow(ctx, query, agencyOrgID, advertiserOrgID, permissionJSON).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}

	return exists, nil
}

// GetActivePermissions retrieves all active permissions for an agency-advertiser relationship
func (r *agencyDelegationRepository) GetActivePermissions(ctx context.Context, agencyOrgID, advertiserOrgID int64) ([]string, error) {
	query := `
		SELECT permissions
		FROM agency_delegations
		WHERE agency_org_id = $1 
		  AND advertiser_org_id = $2 
		  AND status = 'active'
		  AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)`

	var permissionsJSON []byte
	err := r.db.QueryRow(ctx, query, agencyOrgID, advertiserOrgID).Scan(&permissionsJSON)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to get active permissions: %w", err)
	}

	var permissions []string
	if len(permissionsJSON) > 0 {
		if err := json.Unmarshal(permissionsJSON, &permissions); err != nil {
			return nil, fmt.Errorf("failed to unmarshal permissions: %w", err)
		}
	}

	return permissions, nil
}

// ExpireOldDelegations marks expired delegations as revoked
func (r *agencyDelegationRepository) ExpireOldDelegations(ctx context.Context) (int64, error) {
	query := `
		UPDATE agency_delegations
		SET status = 'revoked', updated_at = CURRENT_TIMESTAMP
		WHERE status IN ('pending', 'active')
		  AND expires_at IS NOT NULL
		  AND expires_at <= CURRENT_TIMESTAMP`

	result, err := r.db.Exec(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to expire old delegations: %w", err)
	}

	return result.RowsAffected(), nil
}

// CheckPermissions checks multiple permissions for an agency-advertiser relationship
func (r *agencyDelegationRepository) CheckPermissions(ctx context.Context, agencyOrgID, advertiserOrgID int64, permissions []string) (map[string]bool, error) {
	result := make(map[string]bool)
	
	for _, permission := range permissions {
		hasPermission, err := r.HasPermission(ctx, agencyOrgID, advertiserOrgID, permission)
		if err != nil {
			return nil, fmt.Errorf("failed to check permission %s: %w", permission, err)
		}
		result[permission] = hasPermission
	}
	
	return result, nil
}

// GetActiveDelegationsByAgency retrieves all active delegations for an agency
func (r *agencyDelegationRepository) GetActiveDelegationsByAgency(ctx context.Context, agencyOrgID int64) ([]*domain.AgencyDelegation, error) {
	status := domain.DelegationStatusActive
	filter := domain.DelegationListFilter{
		AgencyOrgID: &agencyOrgID,
		Status:      &status,
	}
	return r.List(ctx, filter)
}

// GetActiveDelegationsByAdvertiser retrieves all active delegations for an advertiser
func (r *agencyDelegationRepository) GetActiveDelegationsByAdvertiser(ctx context.Context, advertiserOrgID int64) ([]*domain.AgencyDelegation, error) {
	status := domain.DelegationStatusActive
	filter := domain.DelegationListFilter{
		AdvertiserOrgID: &advertiserOrgID,
		Status:          &status,
	}
	return r.List(ctx, filter)
}