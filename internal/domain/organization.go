package domain

import (
	"time"
)

// OrganizationType represents the type of organization
type OrganizationType string

const (
	OrganizationTypeAdvertiser    OrganizationType = "advertiser"
	OrganizationTypeAffiliate     OrganizationType = "affiliate"
	OrganizationTypePlatformOwner OrganizationType = "platform_owner"
	OrganizationTypeAgency        OrganizationType = "agency"
)

// IsValid checks if the organization type is valid
func (ot OrganizationType) IsValid() bool {
	switch ot {
	case OrganizationTypeAdvertiser, OrganizationTypeAffiliate, OrganizationTypePlatformOwner:
		return true
	default:
		return false
	}
}

// String returns the string representation of the organization type
func (ot OrganizationType) String() string {
	return string(ot)
}

// GetValidOrganizationTypes returns all valid organization types
func GetValidOrganizationTypes() []OrganizationType {
	return []OrganizationType{
		OrganizationTypeAdvertiser,
		OrganizationTypeAffiliate,
		OrganizationTypePlatformOwner,
	}
}

// Organization represents an organization entity
type Organization struct {
	OrganizationID int64            `json:"organization_id" db:"organization_id"`
	Name           string           `json:"name" db:"name"`
	Type           OrganizationType `json:"type" db:"type"`
	CreatedAt      time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at" db:"updated_at"`
}
