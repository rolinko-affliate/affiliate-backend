package domain

import (
	"fmt"
	"time"
)

// InvitationStatus represents the status of an advertiser association invitation
type InvitationStatus string

const (
	InvitationStatusActive   InvitationStatus = "active"
	InvitationStatusDisabled InvitationStatus = "disabled"
	InvitationStatusExpired  InvitationStatus = "expired"
)

// IsValid checks if the invitation status is valid
func (is InvitationStatus) IsValid() bool {
	switch is {
	case InvitationStatusActive, InvitationStatusDisabled, InvitationStatusExpired:
		return true
	default:
		return false
	}
}

// String returns the string representation of the invitation status
func (is InvitationStatus) String() string {
	return string(is)
}

// AdvertiserAssociationInvitation represents an invitation link created by an advertiser
type AdvertiserAssociationInvitation struct {
	InvitationID      int64            `json:"invitation_id" db:"invitation_id"`
	AdvertiserOrgID   int64            `json:"advertiser_org_id" db:"advertiser_org_id"`
	InvitationToken   string           `json:"invitation_token" db:"invitation_token"`
	Name              string           `json:"name" db:"name"`
	Description       *string          `json:"description,omitempty" db:"description"`
	
	// Access control
	AllowedAffiliateOrgIDs *string `json:"allowed_affiliate_org_ids,omitempty" db:"allowed_affiliate_org_ids"` // JSONB array
	
	// Usage limits
	MaxUses     *int `json:"max_uses,omitempty" db:"max_uses"`
	CurrentUses int  `json:"current_uses" db:"current_uses"`
	
	// Expiration
	ExpiresAt *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	
	// Status and metadata
	Status            InvitationStatus `json:"status" db:"status"`
	CreatedByUserID   string           `json:"created_by_user_id" db:"created_by_user_id"`
	Message           *string          `json:"message,omitempty" db:"message"`
	
	// Default visibility settings for associations created through this invitation
	DefaultAllAffiliatesVisible bool    `json:"default_all_affiliates_visible" db:"default_all_affiliates_visible"`
	DefaultAllCampaignsVisible  bool    `json:"default_all_campaigns_visible" db:"default_all_campaigns_visible"`
	DefaultVisibleAffiliateIDs  *string `json:"default_visible_affiliate_ids,omitempty" db:"default_visible_affiliate_ids"` // JSONB array
	DefaultVisibleCampaignIDs   *string `json:"default_visible_campaign_ids,omitempty" db:"default_visible_campaign_ids"`   // JSONB array
	
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Validate validates the advertiser association invitation data
func (aai *AdvertiserAssociationInvitation) Validate() error {
	if aai.AdvertiserOrgID <= 0 {
		return fmt.Errorf("valid advertiser organization ID is required")
	}
	if aai.InvitationToken == "" {
		return fmt.Errorf("invitation token is required")
	}
	if aai.Name == "" {
		return fmt.Errorf("invitation name is required")
	}
	if !aai.Status.IsValid() {
		return fmt.Errorf("invalid invitation status: %s", aai.Status)
	}
	if aai.CreatedByUserID == "" {
		return fmt.Errorf("created by user ID is required")
	}
	if aai.MaxUses != nil && *aai.MaxUses <= 0 {
		return fmt.Errorf("max uses must be greater than 0 if specified")
	}
	if aai.CurrentUses < 0 {
		return fmt.Errorf("current uses cannot be negative")
	}
	if aai.MaxUses != nil && aai.CurrentUses > *aai.MaxUses {
		return fmt.Errorf("current uses cannot exceed max uses")
	}
	return nil
}

// IsActive returns true if the invitation is active and can be used
func (aai *AdvertiserAssociationInvitation) IsActive() bool {
	return aai.Status == InvitationStatusActive
}

// IsExpired returns true if the invitation has expired
func (aai *AdvertiserAssociationInvitation) IsExpired() bool {
	if aai.Status == InvitationStatusExpired {
		return true
	}
	if aai.ExpiresAt != nil && time.Now().After(*aai.ExpiresAt) {
		return true
	}
	return false
}

// IsUsageLimitReached returns true if the invitation has reached its usage limit
func (aai *AdvertiserAssociationInvitation) IsUsageLimitReached() bool {
	return aai.MaxUses != nil && aai.CurrentUses >= *aai.MaxUses
}

// CanBeUsed returns true if the invitation can be used (active, not expired, not at usage limit)
func (aai *AdvertiserAssociationInvitation) CanBeUsed() bool {
	return aai.IsActive() && !aai.IsExpired() && !aai.IsUsageLimitReached()
}

// CanBeUsedByAffiliate returns true if the invitation can be used by the specified affiliate organization
func (aai *AdvertiserAssociationInvitation) CanBeUsedByAffiliate(affiliateOrgID int64) bool {
	if !aai.CanBeUsed() {
		return false
	}
	
	// If no restrictions, any affiliate can use it
	if aai.AllowedAffiliateOrgIDs == nil {
		return true
	}
	
	// Check if the affiliate is in the allowed list
	// This would need to be implemented with JSON parsing in the service layer
	return true // Placeholder - actual implementation in service layer
}

// InvitationUsageLog represents a log entry for invitation usage
type InvitationUsageLog struct {
	UsageID       int64      `json:"usage_id" db:"usage_id"`
	InvitationID  int64      `json:"invitation_id" db:"invitation_id"`
	AffiliateOrgID int64     `json:"affiliate_org_id" db:"affiliate_org_id"`
	UsedByUserID  *string    `json:"used_by_user_id,omitempty" db:"used_by_user_id"`
	AssociationID *int64     `json:"association_id,omitempty" db:"association_id"`
	
	// Usage metadata
	IPAddress    *string `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent    *string `json:"user_agent,omitempty" db:"user_agent"`
	Success      bool    `json:"success" db:"success"`
	ErrorMessage *string `json:"error_message,omitempty" db:"error_message"`
	
	UsedAt time.Time `json:"used_at" db:"used_at"`
}

// AdvertiserAssociationInvitationWithDetails represents an invitation with additional details
type AdvertiserAssociationInvitationWithDetails struct {
	*AdvertiserAssociationInvitation
	AdvertiserOrganization *Organization `json:"advertiser_organization,omitempty"`
	CreatedByUser          *Profile      `json:"created_by_user,omitempty"`
	UsageCount             int           `json:"usage_count"`
	RecentUsages           []*InvitationUsageLog `json:"recent_usages,omitempty"`
}

// CreateInvitationRequest represents a request to create a new invitation
type CreateInvitationRequest struct {
	AdvertiserOrgID             int64     `json:"advertiser_org_id" binding:"required"`
	Name                        string    `json:"name" binding:"required"`
	Description                 *string   `json:"description,omitempty"`
	AllowedAffiliateOrgIDs      []int64   `json:"allowed_affiliate_org_ids,omitempty"`
	MaxUses                     *int      `json:"max_uses,omitempty"`
	ExpiresAt                   *time.Time `json:"expires_at,omitempty"`
	Message                     *string   `json:"message,omitempty"`
	DefaultAllAffiliatesVisible *bool     `json:"default_all_affiliates_visible,omitempty"`
	DefaultAllCampaignsVisible  *bool     `json:"default_all_campaigns_visible,omitempty"`
	DefaultVisibleAffiliateIDs  []int64   `json:"default_visible_affiliate_ids,omitempty"`
	DefaultVisibleCampaignIDs   []int64   `json:"default_visible_campaign_ids,omitempty"`
}

// UpdateInvitationRequest represents a request to update an existing invitation
type UpdateInvitationRequest struct {
	Name                        *string    `json:"name,omitempty"`
	Description                 *string    `json:"description,omitempty"`
	AllowedAffiliateOrgIDs      []int64    `json:"allowed_affiliate_org_ids,omitempty"`
	MaxUses                     *int       `json:"max_uses,omitempty"`
	ExpiresAt                   *time.Time `json:"expires_at,omitempty"`
	Status                      *InvitationStatus `json:"status,omitempty"`
	Message                     *string    `json:"message,omitempty"`
	DefaultAllAffiliatesVisible *bool      `json:"default_all_affiliates_visible,omitempty"`
	DefaultAllCampaignsVisible  *bool      `json:"default_all_campaigns_visible,omitempty"`
	DefaultVisibleAffiliateIDs  []int64    `json:"default_visible_affiliate_ids,omitempty"`
	DefaultVisibleCampaignIDs   []int64    `json:"default_visible_campaign_ids,omitempty"`
}

// InvitationListFilter represents filters for listing invitations
type InvitationListFilter struct {
	AdvertiserOrgID *int64            `json:"advertiser_org_id,omitempty"`
	Status          *InvitationStatus `json:"status,omitempty"`
	CreatedByUserID *string           `json:"created_by_user_id,omitempty"`
	IncludeExpired  bool              `json:"include_expired,omitempty"`
	Limit           int               `json:"limit,omitempty"`
	Offset          int               `json:"offset,omitempty"`
}

// UseInvitationRequest represents a request to use an invitation
type UseInvitationRequest struct {
	InvitationToken string  `json:"invitation_token" binding:"required"`
	AffiliateOrgID  int64   `json:"affiliate_org_id" binding:"required"`
	Message         *string `json:"message,omitempty"`
	IPAddress       *string `json:"ip_address,omitempty"`
	UserAgent       *string `json:"user_agent,omitempty"`
}

// UseInvitationResponse represents the response when using an invitation
type UseInvitationResponse struct {
	Success       bool                     `json:"success"`
	Association   *OrganizationAssociation `json:"association,omitempty"`
	Invitation    *AdvertiserAssociationInvitation `json:"invitation,omitempty"`
	ErrorMessage  *string                  `json:"error_message,omitempty"`
}