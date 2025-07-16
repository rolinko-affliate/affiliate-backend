package domain

import (
	"time"
)

// FavoritePublisherList represents a favorite publisher list belonging to an organization
type FavoritePublisherList struct {
	ListID         int64     `json:"list_id" db:"list_id"`
	OrganizationID int64     `json:"organization_id" db:"organization_id"`
	Name           string    `json:"name" db:"name"`
	Description    *string   `json:"description,omitempty" db:"description"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	
	// Optional: Include items when fetching with details
	Items []FavoritePublisherListItem `json:"items,omitempty" db:"-"`
}

// FavoritePublisherListItem represents a publisher domain within a favorite list
type FavoritePublisherListItem struct {
	ItemID          int64     `json:"item_id" db:"item_id"`
	ListID          int64     `json:"list_id" db:"list_id"`
	PublisherDomain string    `json:"publisher_domain" db:"publisher_domain"`
	Notes           *string   `json:"notes,omitempty" db:"notes"`
	AddedAt         time.Time `json:"added_at" db:"added_at"`
	
	// Optional: Include publisher details when fetching with details
	Publisher *AnalyticsPublisher `json:"publisher,omitempty" db:"-"`
}

// CreateFavoritePublisherListRequest represents the request to create a new favorite list
type CreateFavoritePublisherListRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=1000"`
}

// UpdateFavoritePublisherListRequest represents the request to update a favorite list
type UpdateFavoritePublisherListRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=1000"`
}

// AddPublisherToListRequest represents the request to add a publisher to a favorite list
type AddPublisherToListRequest struct {
	PublisherDomain string  `json:"publisher_domain" binding:"required,min=1,max=255"`
	Notes           *string `json:"notes,omitempty" binding:"omitempty,max=1000"`
}

// UpdatePublisherInListRequest represents the request to update a publisher's notes in a list
type UpdatePublisherInListRequest struct {
	Notes *string `json:"notes,omitempty" binding:"omitempty,max=1000"`
}

// FavoritePublisherListWithStats represents a list with additional statistics
type FavoritePublisherListWithStats struct {
	FavoritePublisherList
	PublisherCount int `json:"publisher_count"`
}

// Validation methods

// Validate validates the CreateFavoritePublisherListRequest
func (r *CreateFavoritePublisherListRequest) Validate() error {
	if r.Name == "" {
		return ErrInvalidInput
	}
	if len(r.Name) > 255 {
		return ErrInvalidInput
	}
	if r.Description != nil && len(*r.Description) > 1000 {
		return ErrInvalidInput
	}
	return nil
}

// Validate validates the AddPublisherToListRequest
func (r *AddPublisherToListRequest) Validate() error {
	if r.PublisherDomain == "" {
		return ErrInvalidInput
	}
	if len(r.PublisherDomain) > 255 {
		return ErrInvalidInput
	}
	if r.Notes != nil && len(*r.Notes) > 1000 {
		return ErrInvalidInput
	}
	return nil
}



// Validate validates the update publisher in list request
func (r *UpdatePublisherInListRequest) Validate() error {
	if r.Notes != nil && len(*r.Notes) > 1000 {
		return ErrInvalidInput
	}
	return nil
}