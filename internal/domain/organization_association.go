package domain

import (
	"fmt"
	"time"
)

// AssociationStatus represents the status of an organization association
type AssociationStatus string

const (
	AssociationStatusPending   AssociationStatus = "pending"
	AssociationStatusActive    AssociationStatus = "active"
	AssociationStatusSuspended AssociationStatus = "suspended"
	AssociationStatusRejected  AssociationStatus = "rejected"
)

// IsValid checks if the association status is valid
func (as AssociationStatus) IsValid() bool {
	switch as {
	case AssociationStatusPending, AssociationStatusActive, AssociationStatusSuspended, AssociationStatusRejected:
		return true
	default:
		return false
	}
}

// String returns the string representation of the association status
func (as AssociationStatus) String() string {
	return string(as)
}

// AssociationType represents the type of association request
type AssociationType string

const (
	AssociationTypeInvitation AssociationType = "invitation" // Advertiser invites affiliate
	AssociationTypeRequest    AssociationType = "request"    // Affiliate requests to join advertiser
)

// IsValid checks if the association type is valid
func (at AssociationType) IsValid() bool {
	switch at {
	case AssociationTypeInvitation, AssociationTypeRequest:
		return true
	default:
		return false
	}
}

// String returns the string representation of the association type
func (at AssociationType) String() string {
	return string(at)
}

// OrganizationAssociation represents the association between an advertiser and affiliate organization
type OrganizationAssociation struct {
	AssociationID          int64             `json:"association_id" db:"association_id"`
	AdvertiserOrgID        int64             `json:"advertiser_org_id" db:"advertiser_org_id"`
	AffiliateOrgID         int64             `json:"affiliate_org_id" db:"affiliate_org_id"`
	Status                 AssociationStatus `json:"status" db:"status"`
	AssociationType        AssociationType   `json:"association_type" db:"association_type"`
	
	// Visibility settings - JSON arrays of IDs
	VisibleAffiliateIDs    *string `json:"visible_affiliate_ids,omitempty" db:"visible_affiliate_ids"`    // JSONB array of affiliate IDs visible to advertiser
	VisibleCampaignIDs     *string `json:"visible_campaign_ids,omitempty" db:"visible_campaign_ids"`      // JSONB array of campaign IDs visible to affiliate
	
	// Default visibility flags (when true, all affiliates/campaigns are visible)
	AllAffiliatesVisible   bool    `json:"all_affiliates_visible" db:"all_affiliates_visible"`
	AllCampaignsVisible    bool    `json:"all_campaigns_visible" db:"all_campaigns_visible"`
	
	// Request/invitation metadata
	RequestedByUserID      *string `json:"requested_by_user_id,omitempty" db:"requested_by_user_id"`      // UUID of user who initiated
	ApprovedByUserID       *string `json:"approved_by_user_id,omitempty" db:"approved_by_user_id"`        // UUID of user who approved
	Message                *string `json:"message,omitempty" db:"message"`                                // Optional message with request/invitation
	
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
	ApprovedAt             *time.Time `json:"approved_at,omitempty" db:"approved_at"`
}

// Validate validates the organization association data
func (oa *OrganizationAssociation) Validate() error {
	if oa.AdvertiserOrgID <= 0 {
		return fmt.Errorf("valid advertiser organization ID is required")
	}
	if oa.AffiliateOrgID <= 0 {
		return fmt.Errorf("valid affiliate organization ID is required")
	}
	if oa.AdvertiserOrgID == oa.AffiliateOrgID {
		return fmt.Errorf("advertiser and affiliate organization IDs cannot be the same")
	}
	if !oa.Status.IsValid() {
		return fmt.Errorf("invalid association status: %s", oa.Status)
	}
	if !oa.AssociationType.IsValid() {
		return fmt.Errorf("invalid association type: %s", oa.AssociationType)
	}
	return nil
}

// IsActive returns true if the association is active
func (oa *OrganizationAssociation) IsActive() bool {
	return oa.Status == AssociationStatusActive
}

// IsPending returns true if the association is pending approval
func (oa *OrganizationAssociation) IsPending() bool {
	return oa.Status == AssociationStatusPending
}

// IsSuspended returns true if the association is suspended
func (oa *OrganizationAssociation) IsSuspended() bool {
	return oa.Status == AssociationStatusSuspended
}

// IsRejected returns true if the association is rejected
func (oa *OrganizationAssociation) IsRejected() bool {
	return oa.Status == AssociationStatusRejected
}

// CanBeActivated returns true if the association can be activated (from pending status)
func (oa *OrganizationAssociation) CanBeActivated() bool {
	return oa.Status == AssociationStatusPending
}

// CanBeSuspended returns true if the association can be suspended (from active status)
func (oa *OrganizationAssociation) CanBeSuspended() bool {
	return oa.Status == AssociationStatusActive
}

// CanBeReactivated returns true if the association can be reactivated (from suspended status)
func (oa *OrganizationAssociation) CanBeReactivated() bool {
	return oa.Status == AssociationStatusSuspended
}

// OrganizationAssociationWithDetails represents an association with additional organization details
type OrganizationAssociationWithDetails struct {
	*OrganizationAssociation
	AdvertiserOrganization *Organization `json:"advertiser_organization,omitempty"`
	AffiliateOrganization  *Organization `json:"affiliate_organization,omitempty"`
	RequestedByUser        *Profile      `json:"requested_by_user,omitempty"`
	ApprovedByUser         *Profile      `json:"approved_by_user,omitempty"`
}

// CreateAssociationRequest represents a request to create a new association
type CreateAssociationRequest struct {
	AdvertiserOrgID        int64             `json:"advertiser_org_id" binding:"required"`
	AffiliateOrgID         int64             `json:"affiliate_org_id" binding:"required"`
	AssociationType        AssociationType   `json:"association_type" binding:"required"`
	Message                *string           `json:"message,omitempty"`
	VisibleAffiliateIDs    []int64           `json:"visible_affiliate_ids,omitempty"`    // For invitations from advertiser
	VisibleCampaignIDs     []int64           `json:"visible_campaign_ids,omitempty"`     // For requests from affiliate
	AllAffiliatesVisible   *bool             `json:"all_affiliates_visible,omitempty"`
	AllCampaignsVisible    *bool             `json:"all_campaigns_visible,omitempty"`
}

// UpdateAssociationRequest represents a request to update an existing association
type UpdateAssociationRequest struct {
	Status                 *AssociationStatus `json:"status,omitempty"`
	VisibleAffiliateIDs    []int64            `json:"visible_affiliate_ids,omitempty"`
	VisibleCampaignIDs     []int64            `json:"visible_campaign_ids,omitempty"`
	AllAffiliatesVisible   *bool              `json:"all_affiliates_visible,omitempty"`
	AllCampaignsVisible    *bool              `json:"all_campaigns_visible,omitempty"`
}

// AssociationListFilter represents filters for listing associations
type AssociationListFilter struct {
	AdvertiserOrgID *int64             `json:"advertiser_org_id,omitempty"`
	AffiliateOrgID  *int64             `json:"affiliate_org_id,omitempty"`
	Status          *AssociationStatus `json:"status,omitempty"`
	AssociationType *AssociationType   `json:"association_type,omitempty"`
	Limit           int                `json:"limit,omitempty"`
	Offset          int                `json:"offset,omitempty"`
}