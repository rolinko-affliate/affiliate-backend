package domain

import (
	"time"

	"github.com/google/uuid"
)

// Profile represents a user profile in the system
type Profile struct {
	ID             uuid.UUID  `json:"id" db:"id"` // Stores the auth.uid() from Supabase Auth
	OrganizationID *int64     `json:"organization_id,omitempty" db:"organization_id"` // Pointer for NULLable
	RoleID         int        `json:"role_id" db:"role_id"`
	RoleName       string     `json:"role_name" db:"role_name"` // Mandatory role name
	Email          string     `json:"email" db:"email"`
	FirstName      *string    `json:"first_name,omitempty" db:"first_name"` // Pointer for NULLable
	LastName       *string    `json:"last_name,omitempty" db:"last_name"`   // Pointer for NULLable
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// Role represents a user role in the system
type Role struct {
	RoleID      int     `json:"role_id" db:"role_id"`
	Name        string  `json:"name" db:"name"` // e.g., 'Admin', 'AdvertiserManager', 'AffiliateManager', 'Affiliate'
	Description *string `json:"description,omitempty" db:"description"` // Pointer for NULLable
}

// Note: Organization struct is defined in organization.go