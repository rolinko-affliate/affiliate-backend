package domain

import (
	"fmt"
	"time"
)

// DelegationPermission represents a specific permission that can be delegated
type DelegationPermission string

const (
	// Campaign permissions
	PermissionCampaignCreate     DelegationPermission = "campaign_create"
	PermissionCampaignManage     DelegationPermission = "campaign_manage"
	PermissionCampaignView       DelegationPermission = "campaign_view"
	PermissionCampaignDelete     DelegationPermission = "campaign_delete"
	
	// Association permissions
	PermissionAssociationInvite  DelegationPermission = "association_invite"
	PermissionAssociationManage  DelegationPermission = "association_manage"
	PermissionAssociationView    DelegationPermission = "association_view"
	
	// Invitation permissions
	PermissionInvitationCreate   DelegationPermission = "invitation_create"
	PermissionInvitationManage   DelegationPermission = "invitation_manage"
	PermissionInvitationView     DelegationPermission = "invitation_view"
	
	// Analytics permissions
	PermissionAnalyticsView      DelegationPermission = "analytics_view"
	PermissionAnalyticsExport    DelegationPermission = "analytics_export"
	
	// Billing permissions
	PermissionBillingView        DelegationPermission = "billing_view"
	PermissionBillingManage      DelegationPermission = "billing_manage"
	
	// Organization permissions
	PermissionOrganizationView   DelegationPermission = "organization_view"
	PermissionOrganizationManage DelegationPermission = "organization_manage"
)

// IsValid checks if the delegation permission is valid
func (dp DelegationPermission) IsValid() bool {
	switch dp {
	case PermissionCampaignCreate, PermissionCampaignManage, PermissionCampaignView, PermissionCampaignDelete,
		 PermissionAssociationInvite, PermissionAssociationManage, PermissionAssociationView,
		 PermissionInvitationCreate, PermissionInvitationManage, PermissionInvitationView,
		 PermissionAnalyticsView, PermissionAnalyticsExport,
		 PermissionBillingView, PermissionBillingManage,
		 PermissionOrganizationView, PermissionOrganizationManage:
		return true
	default:
		return false
	}
}

// String returns the string representation of the delegation permission
func (dp DelegationPermission) String() string {
	return string(dp)
}

// GetAllDelegationPermissions returns all valid delegation permissions
func GetAllDelegationPermissions() []DelegationPermission {
	return []DelegationPermission{
		PermissionCampaignCreate, PermissionCampaignManage, PermissionCampaignView, PermissionCampaignDelete,
		PermissionAssociationInvite, PermissionAssociationManage, PermissionAssociationView,
		PermissionInvitationCreate, PermissionInvitationManage, PermissionInvitationView,
		PermissionAnalyticsView, PermissionAnalyticsExport,
		PermissionBillingView, PermissionBillingManage,
		PermissionOrganizationView, PermissionOrganizationManage,
	}
}

// DelegationStatus represents the status of an agency delegation
type DelegationStatus string

const (
	DelegationStatusPending   DelegationStatus = "pending"
	DelegationStatusActive    DelegationStatus = "active"
	DelegationStatusSuspended DelegationStatus = "suspended"
	DelegationStatusRevoked   DelegationStatus = "revoked"
)

// IsValid checks if the delegation status is valid
func (ds DelegationStatus) IsValid() bool {
	switch ds {
	case DelegationStatusPending, DelegationStatusActive, DelegationStatusSuspended, DelegationStatusRevoked:
		return true
	default:
		return false
	}
}

// String returns the string representation of the delegation status
func (ds DelegationStatus) String() string {
	return string(ds)
}

// AgencyDelegation represents the delegation relationship between an agency and an advertiser organization
type AgencyDelegation struct {
	DelegationID      int64            `json:"delegation_id" db:"delegation_id"`
	AgencyOrgID       int64            `json:"agency_org_id" db:"agency_org_id"`
	AdvertiserOrgID   int64            `json:"advertiser_org_id" db:"advertiser_org_id"`
	Status            DelegationStatus `json:"status" db:"status"`
	
	// Permissions granted to the agency (JSONB array of permission strings)
	Permissions       *string          `json:"permissions,omitempty" db:"permissions"`
	
	// Delegation metadata
	DelegatedByUserID *string          `json:"delegated_by_user_id,omitempty" db:"delegated_by_user_id"` // UUID of advertiser user who created delegation
	AcceptedByUserID  *string          `json:"accepted_by_user_id,omitempty" db:"accepted_by_user_id"`   // UUID of agency user who accepted delegation
	Message           *string          `json:"message,omitempty" db:"message"`                           // Optional message with delegation
	
	// Expiration settings
	ExpiresAt         *time.Time       `json:"expires_at,omitempty" db:"expires_at"`                     // Optional expiration date
	
	CreatedAt         time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at" db:"updated_at"`
	AcceptedAt        *time.Time       `json:"accepted_at,omitempty" db:"accepted_at"`
}

// Validate validates the agency delegation data
func (ad *AgencyDelegation) Validate() error {
	if ad.AgencyOrgID <= 0 {
		return fmt.Errorf("valid agency organization ID is required")
	}
	if ad.AdvertiserOrgID <= 0 {
		return fmt.Errorf("valid advertiser organization ID is required")
	}
	if ad.AgencyOrgID == ad.AdvertiserOrgID {
		return fmt.Errorf("agency and advertiser organization IDs cannot be the same")
	}
	if !ad.Status.IsValid() {
		return fmt.Errorf("invalid delegation status: %s", ad.Status)
	}
	
	// Validate expiration date if set
	if ad.ExpiresAt != nil && ad.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("expiration date cannot be in the past")
	}
	
	return nil
}

// IsActive returns true if the delegation is active and not expired
func (ad *AgencyDelegation) IsActive() bool {
	if ad.Status != DelegationStatusActive {
		return false
	}
	if ad.ExpiresAt != nil && ad.ExpiresAt.Before(time.Now()) {
		return false
	}
	return true
}

// IsPending returns true if the delegation is pending acceptance
func (ad *AgencyDelegation) IsPending() bool {
	return ad.Status == DelegationStatusPending
}

// IsSuspended returns true if the delegation is suspended
func (ad *AgencyDelegation) IsSuspended() bool {
	return ad.Status == DelegationStatusSuspended
}

// IsRevoked returns true if the delegation is revoked
func (ad *AgencyDelegation) IsRevoked() bool {
	return ad.Status == DelegationStatusRevoked
}

// IsExpired returns true if the delegation has expired
func (ad *AgencyDelegation) IsExpired() bool {
	return ad.ExpiresAt != nil && ad.ExpiresAt.Before(time.Now())
}

// CanBeAccepted returns true if the delegation can be accepted (from pending status)
func (ad *AgencyDelegation) CanBeAccepted() bool {
	return ad.Status == DelegationStatusPending && !ad.IsExpired()
}

// CanBeSuspended returns true if the delegation can be suspended (from active status)
func (ad *AgencyDelegation) CanBeSuspended() bool {
	return ad.Status == DelegationStatusActive
}

// CanBeReactivated returns true if the delegation can be reactivated (from suspended status)
func (ad *AgencyDelegation) CanBeReactivated() bool {
	return ad.Status == DelegationStatusSuspended && !ad.IsExpired()
}

// CanBeRevoked returns true if the delegation can be revoked
func (ad *AgencyDelegation) CanBeRevoked() bool {
	return ad.Status == DelegationStatusPending || ad.Status == DelegationStatusActive || ad.Status == DelegationStatusSuspended
}

// AgencyDelegationWithDetails represents a delegation with additional organization details
type AgencyDelegationWithDetails struct {
	*AgencyDelegation
	AgencyOrganization     *Organization `json:"agency_organization,omitempty"`
	AdvertiserOrganization *Organization `json:"advertiser_organization,omitempty"`
	DelegatedByUser        *Profile      `json:"delegated_by_user,omitempty"`
	AcceptedByUser         *Profile      `json:"accepted_by_user,omitempty"`
}

// CreateDelegationRequest represents a request to create a new agency delegation
type CreateDelegationRequest struct {
	AgencyOrgID       int64                  `json:"agency_org_id" binding:"required"`
	AdvertiserOrgID   int64                  `json:"advertiser_org_id" binding:"required"`
	Permissions       []DelegationPermission `json:"permissions" binding:"required"`
	Message           *string                `json:"message,omitempty"`
	ExpiresAt         *time.Time             `json:"expires_at,omitempty"`
}

// UpdateDelegationRequest represents a request to update an existing delegation
type UpdateDelegationRequest struct {
	Status      *DelegationStatus      `json:"status,omitempty"`
	Permissions []DelegationPermission `json:"permissions,omitempty"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
}

// DelegationListFilter represents filters for listing delegations
type DelegationListFilter struct {
	AgencyOrgID     *int64            `json:"agency_org_id,omitempty"`
	AdvertiserOrgID *int64            `json:"advertiser_org_id,omitempty"`
	Status          *DelegationStatus `json:"status,omitempty"`
	IncludeExpired  bool              `json:"include_expired,omitempty"`
	Limit           int               `json:"limit,omitempty"`
	Offset          int               `json:"offset,omitempty"`
}

// PermissionCheckRequest represents a request to check if an agency has specific permissions
type PermissionCheckRequest struct {
	AgencyOrgID     int64                  `json:"agency_org_id" binding:"required"`
	AdvertiserOrgID int64                  `json:"advertiser_org_id" binding:"required"`
	Permissions     []DelegationPermission `json:"permissions" binding:"required"`
}

// PermissionCheckResponse represents the response to a permission check
type PermissionCheckResponse struct {
	HasPermissions    bool                            `json:"has_permissions"`
	PermissionResults map[DelegationPermission]bool   `json:"permission_results"`
	DelegationStatus  DelegationStatus                `json:"delegation_status"`
	IsExpired         bool                            `json:"is_expired"`
}