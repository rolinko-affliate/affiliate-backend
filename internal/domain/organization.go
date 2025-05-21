package domain

import (
	"time"
)

// Organization represents an organization entity
type Organization struct {
	OrganizationID int64     `json:"organization_id" db:"organization_id"`
	Name           string    `json:"name" db:"name"`
	Status         string    `json:"status" db:"status"` // 'active', 'inactive', 'suspended'
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}