package domain

import (
	"time"
)

// Publisher status constants
const (
	PublisherStatusAdded     = "added"
	PublisherStatusContacted = "contacted"
	PublisherStatusAccepted  = "accepted"
)

// Valid status transitions
var ValidStatusTransitions = map[string][]string{
	PublisherStatusAdded:     {PublisherStatusContacted},
	PublisherStatusContacted: {PublisherStatusAccepted, PublisherStatusAdded}, // Allow going back to added
	PublisherStatusAccepted:  {PublisherStatusContacted},                     // Allow going back to contacted
}

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
	Status          string    `json:"status" db:"status"`
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
	Status          *string `json:"status,omitempty" binding:"omitempty,oneof=added contacted accepted"`
}

// UpdatePublisherInListRequest represents the request to update a publisher's notes in a list
type UpdatePublisherInListRequest struct {
	Notes *string `json:"notes,omitempty" binding:"omitempty,max=1000"`
}

// UpdatePublisherStatusRequest represents the request to update a publisher's status in a list
type UpdatePublisherStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=added contacted accepted"`
}

// FavoritePublisherListWithStats represents a list with additional statistics
type FavoritePublisherListWithStats struct {
	FavoritePublisherList
	PublisherCount int `json:"publisher_count"`
}

// Validation helpers
func validateStringLength(value string, minLen, maxLen int) error {
	if len(value) < minLen || len(value) > maxLen {
		return ErrInvalidInput
	}
	return nil
}

func validateOptionalStringLength(value *string, maxLen int) error {
	if value != nil && len(*value) > maxLen {
		return ErrInvalidInput
	}
	return nil
}

func validateStatus(status string) error {
	switch status {
	case PublisherStatusAdded, PublisherStatusContacted, PublisherStatusAccepted:
		return nil
	default:
		return ErrInvalidInput
	}
}

func ValidateStatusTransition(currentStatus, newStatus string) error {
	if currentStatus == newStatus {
		return nil // No change is valid
	}
	
	validTransitions, exists := ValidStatusTransitions[currentStatus]
	if !exists {
		return ErrInvalidInput
	}
	
	for _, validStatus := range validTransitions {
		if validStatus == newStatus {
			return nil
		}
	}
	
	return ErrInvalidInput
}

// Validation methods

// Validate validates the CreateFavoritePublisherListRequest
func (r *CreateFavoritePublisherListRequest) Validate() error {
	if err := validateStringLength(r.Name, 1, 255); err != nil {
		return err
	}
	return validateOptionalStringLength(r.Description, 1000)
}

// Validate validates the UpdateFavoritePublisherListRequest
func (r *UpdateFavoritePublisherListRequest) Validate() error {
	if r.Name != nil {
		if err := validateStringLength(*r.Name, 1, 255); err != nil {
			return err
		}
	}
	return validateOptionalStringLength(r.Description, 1000)
}

// Validate validates the AddPublisherToListRequest
func (r *AddPublisherToListRequest) Validate() error {
	if err := validateStringLength(r.PublisherDomain, 1, 255); err != nil {
		return err
	}
	if err := validateOptionalStringLength(r.Notes, 1000); err != nil {
		return err
	}
	if r.Status != nil {
		return validateStatus(*r.Status)
	}
	return nil
}

// Validate validates the UpdatePublisherInListRequest
func (r *UpdatePublisherInListRequest) Validate() error {
	return validateOptionalStringLength(r.Notes, 1000)
}

// Validate validates the UpdatePublisherStatusRequest
func (r *UpdatePublisherStatusRequest) Validate() error {
	return validateStatus(r.Status)
}
