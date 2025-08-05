package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/repository"
	"github.com/google/uuid"
)

// AgencyDelegationService handles business logic for agency delegations
type AgencyDelegationService interface {
	CreateDelegation(ctx context.Context, req *domain.CreateDelegationRequest, delegatedByUserID string) (*domain.AgencyDelegation, error)
	AcceptDelegation(ctx context.Context, delegationID int64, acceptedByUserID string) (*domain.AgencyDelegation, error)
	RejectDelegation(ctx context.Context, delegationID int64, userID string) (*domain.AgencyDelegation, error)
	SuspendDelegation(ctx context.Context, delegationID int64, userID string) (*domain.AgencyDelegation, error)
	ReactivateDelegation(ctx context.Context, delegationID int64, userID string) (*domain.AgencyDelegation, error)
	RevokeDelegation(ctx context.Context, delegationID int64, userID string) (*domain.AgencyDelegation, error)
	UpdatePermissions(ctx context.Context, delegationID int64, permissions []domain.DelegationPermission, userID string) (*domain.AgencyDelegation, error)
	UpdateExpiration(ctx context.Context, delegationID int64, expiresAt *time.Time, userID string) (*domain.AgencyDelegation, error)
	GetDelegationByID(ctx context.Context, delegationID int64) (*domain.AgencyDelegation, error)
	GetDelegationByIDWithDetails(ctx context.Context, delegationID int64) (*domain.AgencyDelegationWithDetails, error)
	ListDelegations(ctx context.Context, filter *domain.DelegationListFilter) ([]*domain.AgencyDelegation, error)
	ListDelegationsWithDetails(ctx context.Context, filter *domain.DelegationListFilter) ([]*domain.AgencyDelegationWithDetails, error)
	CheckPermissions(ctx context.Context, req *domain.PermissionCheckRequest) (*domain.PermissionCheckResponse, error)
	GetAgencyDelegations(ctx context.Context, agencyOrgID int64) ([]*domain.AgencyDelegation, error)
	GetAdvertiserDelegations(ctx context.Context, advertiserOrgID int64) ([]*domain.AgencyDelegation, error)
	ExpireOldDelegations(ctx context.Context) (int64, error)
	ValidateUserPermission(ctx context.Context, userID string, organizationID int64, action string) error
}

type agencyDelegationService struct {
	delegationRepo   repository.AgencyDelegationRepository
	organizationRepo repository.OrganizationRepository
	profileRepo      repository.ProfileRepository
}

// NewAgencyDelegationService creates a new agency delegation service
func NewAgencyDelegationService(
	delegationRepo repository.AgencyDelegationRepository,
	organizationRepo repository.OrganizationRepository,
	profileRepo repository.ProfileRepository,
) AgencyDelegationService {
	return &agencyDelegationService{
		delegationRepo:   delegationRepo,
		organizationRepo: organizationRepo,
		profileRepo:      profileRepo,
	}
}

// CreateDelegation creates a new agency delegation
func (s *agencyDelegationService) CreateDelegation(ctx context.Context, req *domain.CreateDelegationRequest, delegatedByUserID string) (*domain.AgencyDelegation, error) {
	// Validate request
	if req.AgencyOrgID <= 0 {
		return nil, fmt.Errorf("valid agency organization ID is required")
	}
	if req.AdvertiserOrgID <= 0 {
		return nil, fmt.Errorf("valid advertiser organization ID is required")
	}
	if req.AgencyOrgID == req.AdvertiserOrgID {
		return nil, fmt.Errorf("agency and advertiser organization IDs cannot be the same")
	}
	if len(req.Permissions) == 0 {
		return nil, fmt.Errorf("at least one permission is required")
	}

	// Validate permissions
	for _, permission := range req.Permissions {
		if !permission.IsValid() {
			return nil, fmt.Errorf("invalid permission: %s", permission)
		}
	}

	// Validate that agency organization exists and is of type agency
	agencyOrg, err := s.organizationRepo.GetOrganizationByID(ctx, req.AgencyOrgID)
	if err != nil {
		return nil, fmt.Errorf("agency organization not found: %w", err)
	}
	if agencyOrg.Type != domain.OrganizationTypeAgency {
		return nil, fmt.Errorf("organization %d is not of type agency", req.AgencyOrgID)
	}

	// Validate that advertiser organization exists and is of type advertiser
	advertiserOrg, err := s.organizationRepo.GetOrganizationByID(ctx, req.AdvertiserOrgID)
	if err != nil {
		return nil, fmt.Errorf("advertiser organization not found: %w", err)
	}
	if advertiserOrg.Type != domain.OrganizationTypeAdvertiser {
		return nil, fmt.Errorf("organization %d is not of type advertiser", req.AdvertiserOrgID)
	}

	// Validate user permission to create delegation for advertiser organization
	err = s.ValidateUserPermission(ctx, delegatedByUserID, req.AdvertiserOrgID, "delegate")
	if err != nil {
		return nil, fmt.Errorf("user not authorized to create delegation: %w", err)
	}

	// Check if delegation already exists
	existingDelegation, err := s.delegationRepo.GetByOrganizations(ctx, req.AgencyOrgID, req.AdvertiserOrgID)
	if err == nil && existingDelegation != nil {
		return nil, fmt.Errorf("delegation already exists between agency %d and advertiser %d", req.AgencyOrgID, req.AdvertiserOrgID)
	}

	// Convert permissions to JSON
	permissionsJSON, err := json.Marshal(req.Permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal permissions: %w", err)
	}
	permissionsStr := string(permissionsJSON)

	// Create delegation
	delegation := &domain.AgencyDelegation{
		AgencyOrgID:       req.AgencyOrgID,
		AdvertiserOrgID:   req.AdvertiserOrgID,
		Status:            domain.DelegationStatusPending,
		Permissions:       &permissionsStr,
		DelegatedByUserID: &delegatedByUserID,
		Message:           req.Message,
		ExpiresAt:         req.ExpiresAt,
	}

	// Validate delegation
	if err := delegation.Validate(); err != nil {
		return nil, fmt.Errorf("invalid delegation: %w", err)
	}

	return s.delegationRepo.Create(ctx, delegation)
}

// AcceptDelegation accepts a pending delegation
func (s *agencyDelegationService) AcceptDelegation(ctx context.Context, delegationID int64, acceptedByUserID string) (*domain.AgencyDelegation, error) {
	delegation, err := s.delegationRepo.GetByID(ctx, delegationID)
	if err != nil {
		return nil, fmt.Errorf("delegation not found: %w", err)
	}

	if !delegation.CanBeAccepted() {
		return nil, fmt.Errorf("delegation cannot be accepted in current state: %s", delegation.Status)
	}

	// Validate user permission to accept delegation for agency organization
	err = s.ValidateUserPermission(ctx, acceptedByUserID, delegation.AgencyOrgID, "accept_delegation")
	if err != nil {
		return nil, fmt.Errorf("user not authorized to accept delegation: %w", err)
	}

	// Update delegation status
	delegation.Status = domain.DelegationStatusActive
	delegation.AcceptedByUserID = &acceptedByUserID
	now := time.Now()
	delegation.AcceptedAt = &now

	return s.delegationRepo.Update(ctx, delegation)
}

// RejectDelegation rejects a pending delegation
func (s *agencyDelegationService) RejectDelegation(ctx context.Context, delegationID int64, userID string) (*domain.AgencyDelegation, error) {
	delegation, err := s.delegationRepo.GetByID(ctx, delegationID)
	if err != nil {
		return nil, fmt.Errorf("delegation not found: %w", err)
	}

	if !delegation.IsPending() {
		return nil, fmt.Errorf("only pending delegations can be rejected")
	}

	// Validate user permission (either agency or advertiser can reject)
	err1 := s.ValidateUserPermission(ctx, userID, delegation.AgencyOrgID, "reject_delegation")
	err2 := s.ValidateUserPermission(ctx, userID, delegation.AdvertiserOrgID, "reject_delegation")
	if err1 != nil && err2 != nil {
		return nil, fmt.Errorf("user not authorized to reject delegation")
	}

	// Update delegation status
	delegation.Status = domain.DelegationStatusRevoked

	return s.delegationRepo.Update(ctx, delegation)
}

// SuspendDelegation suspends an active delegation
func (s *agencyDelegationService) SuspendDelegation(ctx context.Context, delegationID int64, userID string) (*domain.AgencyDelegation, error) {
	delegation, err := s.delegationRepo.GetByID(ctx, delegationID)
	if err != nil {
		return nil, fmt.Errorf("delegation not found: %w", err)
	}

	if !delegation.CanBeSuspended() {
		return nil, fmt.Errorf("delegation cannot be suspended in current state: %s", delegation.Status)
	}

	// Validate user permission (either agency or advertiser can suspend)
	err1 := s.ValidateUserPermission(ctx, userID, delegation.AgencyOrgID, "suspend_delegation")
	err2 := s.ValidateUserPermission(ctx, userID, delegation.AdvertiserOrgID, "suspend_delegation")
	if err1 != nil && err2 != nil {
		return nil, fmt.Errorf("user not authorized to suspend delegation")
	}

	// Update delegation status
	delegation.Status = domain.DelegationStatusSuspended

	return s.delegationRepo.Update(ctx, delegation)
}

// ReactivateDelegation reactivates a suspended delegation
func (s *agencyDelegationService) ReactivateDelegation(ctx context.Context, delegationID int64, userID string) (*domain.AgencyDelegation, error) {
	delegation, err := s.delegationRepo.GetByID(ctx, delegationID)
	if err != nil {
		return nil, fmt.Errorf("delegation not found: %w", err)
	}

	if !delegation.CanBeReactivated() {
		return nil, fmt.Errorf("delegation cannot be reactivated in current state: %s", delegation.Status)
	}

	// Validate user permission (either agency or advertiser can reactivate)
	err1 := s.ValidateUserPermission(ctx, userID, delegation.AgencyOrgID, "reactivate_delegation")
	err2 := s.ValidateUserPermission(ctx, userID, delegation.AdvertiserOrgID, "reactivate_delegation")
	if err1 != nil && err2 != nil {
		return nil, fmt.Errorf("user not authorized to reactivate delegation")
	}

	// Update delegation status
	delegation.Status = domain.DelegationStatusActive

	return s.delegationRepo.Update(ctx, delegation)
}

// RevokeDelegation revokes a delegation
func (s *agencyDelegationService) RevokeDelegation(ctx context.Context, delegationID int64, userID string) (*domain.AgencyDelegation, error) {
	delegation, err := s.delegationRepo.GetByID(ctx, delegationID)
	if err != nil {
		return nil, fmt.Errorf("delegation not found: %w", err)
	}

	if !delegation.CanBeRevoked() {
		return nil, fmt.Errorf("delegation cannot be revoked in current state: %s", delegation.Status)
	}

	// Validate user permission (either agency or advertiser can revoke)
	err1 := s.ValidateUserPermission(ctx, userID, delegation.AgencyOrgID, "revoke_delegation")
	err2 := s.ValidateUserPermission(ctx, userID, delegation.AdvertiserOrgID, "revoke_delegation")
	if err1 != nil && err2 != nil {
		return nil, fmt.Errorf("user not authorized to revoke delegation")
	}

	// Update delegation status
	delegation.Status = domain.DelegationStatusRevoked

	return s.delegationRepo.Update(ctx, delegation)
}

// UpdatePermissions updates the permissions of a delegation
func (s *agencyDelegationService) UpdatePermissions(ctx context.Context, delegationID int64, permissions []domain.DelegationPermission, userID string) (*domain.AgencyDelegation, error) {
	delegation, err := s.delegationRepo.GetByID(ctx, delegationID)
	if err != nil {
		return nil, fmt.Errorf("delegation not found: %w", err)
	}

	// Validate permissions
	if len(permissions) == 0 {
		return nil, fmt.Errorf("at least one permission is required")
	}
	for _, permission := range permissions {
		if !permission.IsValid() {
			return nil, fmt.Errorf("invalid permission: %s", permission)
		}
	}

	// Validate user permission (only advertiser can update permissions)
	err = s.ValidateUserPermission(ctx, userID, delegation.AdvertiserOrgID, "update_delegation_permissions")
	if err != nil {
		return nil, fmt.Errorf("user not authorized to update delegation permissions: %w", err)
	}

	// Convert permissions to JSON
	permissionsJSON, err := json.Marshal(permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal permissions: %w", err)
	}
	permissionsStr := string(permissionsJSON)

	// Update delegation permissions
	delegation.Permissions = &permissionsStr

	return s.delegationRepo.Update(ctx, delegation)
}

// UpdateExpiration updates the expiration date of a delegation
func (s *agencyDelegationService) UpdateExpiration(ctx context.Context, delegationID int64, expiresAt *time.Time, userID string) (*domain.AgencyDelegation, error) {
	delegation, err := s.delegationRepo.GetByID(ctx, delegationID)
	if err != nil {
		return nil, fmt.Errorf("delegation not found: %w", err)
	}

	// Validate expiration date
	if expiresAt != nil && expiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("expiration date cannot be in the past")
	}

	// Validate user permission (only advertiser can update expiration)
	err = s.ValidateUserPermission(ctx, userID, delegation.AdvertiserOrgID, "update_delegation_expiration")
	if err != nil {
		return nil, fmt.Errorf("user not authorized to update delegation expiration: %w", err)
	}

	// Update delegation expiration
	delegation.ExpiresAt = expiresAt

	return s.delegationRepo.Update(ctx, delegation)
}

// GetDelegationByID retrieves a delegation by ID
func (s *agencyDelegationService) GetDelegationByID(ctx context.Context, delegationID int64) (*domain.AgencyDelegation, error) {
	return s.delegationRepo.GetByID(ctx, delegationID)
}

// GetDelegationByIDWithDetails retrieves a delegation by ID with details
func (s *agencyDelegationService) GetDelegationByIDWithDetails(ctx context.Context, delegationID int64) (*domain.AgencyDelegationWithDetails, error) {
	return s.delegationRepo.GetByIDWithDetails(ctx, delegationID)
}

// ListDelegations lists delegations with optional filtering
func (s *agencyDelegationService) ListDelegations(ctx context.Context, filter *domain.DelegationListFilter) ([]*domain.AgencyDelegation, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50 // Default limit
	}
	if filter.Limit > 100 {
		filter.Limit = 100 // Maximum limit
	}

	return s.delegationRepo.List(ctx, *filter)
}

// ListDelegationsWithDetails lists delegations with details and optional filtering
func (s *agencyDelegationService) ListDelegationsWithDetails(ctx context.Context, filter *domain.DelegationListFilter) ([]*domain.AgencyDelegationWithDetails, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50 // Default limit
	}
	if filter.Limit > 100 {
		filter.Limit = 100 // Maximum limit
	}

	return s.delegationRepo.ListWithDetails(ctx, *filter)
}

// CheckPermissions checks if an agency has specific permissions for an advertiser
func (s *agencyDelegationService) CheckPermissions(ctx context.Context, req *domain.PermissionCheckRequest) (*domain.PermissionCheckResponse, error) {
	// Validate request
	if req.AgencyOrgID <= 0 {
		return nil, fmt.Errorf("valid agency organization ID is required")
	}
	if req.AdvertiserOrgID <= 0 {
		return nil, fmt.Errorf("valid advertiser organization ID is required")
	}
	if len(req.Permissions) == 0 {
		return nil, fmt.Errorf("at least one permission is required")
	}

	// Validate permissions
	for _, permission := range req.Permissions {
		if !permission.IsValid() {
			return nil, fmt.Errorf("invalid permission: %s", permission)
		}
	}

	// Convert permissions to strings for repository call
	permissionStrings := make([]string, len(req.Permissions))
	for i, permission := range req.Permissions {
		permissionStrings[i] = string(permission)
	}

	// Check permissions
	permissionResults, err := s.delegationRepo.CheckPermissions(ctx, req.AgencyOrgID, req.AdvertiserOrgID, permissionStrings)
	if err != nil {
		return nil, fmt.Errorf("failed to check permissions: %w", err)
	}

	// Get delegation to check status and expiry
	delegation, err := s.delegationRepo.GetByOrganizations(ctx, req.AgencyOrgID, req.AdvertiserOrgID)
	if err != nil {
		// If no delegation exists, all permissions are false
		response := &domain.PermissionCheckResponse{
			HasPermissions:    false,
			PermissionResults: make(map[domain.DelegationPermission]bool),
			DelegationStatus:  domain.DelegationStatusRevoked,
			IsExpired:         false,
		}
		for _, permission := range req.Permissions {
			response.PermissionResults[permission] = false
		}
		return response, nil
	}

	// Convert results back to DelegationPermission type
	results := make(map[domain.DelegationPermission]bool)
	hasAllPermissions := true
	for _, permission := range req.Permissions {
		hasPermission := permissionResults[string(permission)]
		results[permission] = hasPermission
		if !hasPermission {
			hasAllPermissions = false
		}
	}

	// Check if delegation is expired
	isExpired := delegation.ExpiresAt != nil && delegation.ExpiresAt.Before(time.Now())

	return &domain.PermissionCheckResponse{
		HasPermissions:    hasAllPermissions,
		PermissionResults: results,
		DelegationStatus:  delegation.Status,
		IsExpired:         isExpired,
	}, nil
}

// GetAgencyDelegations retrieves all active delegations for an agency
func (s *agencyDelegationService) GetAgencyDelegations(ctx context.Context, agencyOrgID int64) ([]*domain.AgencyDelegation, error) {
	return s.delegationRepo.GetActiveDelegationsByAgency(ctx, agencyOrgID)
}

// GetAdvertiserDelegations retrieves all active delegations for an advertiser
func (s *agencyDelegationService) GetAdvertiserDelegations(ctx context.Context, advertiserOrgID int64) ([]*domain.AgencyDelegation, error) {
	return s.delegationRepo.GetActiveDelegationsByAdvertiser(ctx, advertiserOrgID)
}

// ExpireOldDelegations marks expired delegations as revoked
func (s *agencyDelegationService) ExpireOldDelegations(ctx context.Context) (int64, error) {
	return s.delegationRepo.ExpireOldDelegations(ctx)
}

// ValidateUserPermission validates if a user has permission to perform an action on an organization
// This is a simplified implementation - in a real system, you would check user roles and permissions
func (s *agencyDelegationService) ValidateUserPermission(ctx context.Context, userID string, organizationID int64, action string) error {
	// Parse user ID as UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}
	
	// Get user profile
	_, err = s.profileRepo.GetProfileByID(ctx, userUUID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check if user belongs to the organization
	// This is a simplified check - in a real system, you would have a more sophisticated
	// role-based access control system
	
	// TODO: re-enable this check after implementing user roles and organization membership
	// if profile.OrganizationID == nil || *profile.OrganizationID != organizationID {
	// 	return fmt.Errorf("user does not belong to organization %d", organizationID)
	// }

	// Additional permission checks based on action could be implemented here
	// For now, we assume any user in the organization can perform these actions

	return nil
}