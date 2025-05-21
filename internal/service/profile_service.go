package service

import (
	"context"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/repository"
	"github.com/google/uuid"
)

// ProfileService defines the interface for profile operations
type ProfileService interface {
	CreateNewUserProfile(ctx context.Context, userID uuid.UUID, email string, initialOrgID *int64, initialRoleID int) (*domain.Profile, error)
	GetProfileByID(ctx context.Context, id uuid.UUID) (*domain.Profile, error)
	UpdateProfile(ctx context.Context, profile *domain.Profile) (*domain.Profile, error)
	DeleteProfile(ctx context.Context, id uuid.UUID) error
	GetRoleByID(ctx context.Context, roleID int) (*domain.Role, error)
	UpsertProfile(ctx context.Context, userID uuid.UUID, email string, orgID *int64, roleID int, firstName, lastName *string) (*domain.Profile, error)
}

// profileService implements ProfileService
type profileService struct {
	profileRepo repository.ProfileRepository
}

// NewProfileService creates a new profile service
func NewProfileService(profileRepo repository.ProfileRepository) ProfileService {
	return &profileService{profileRepo: profileRepo}
}

// CreateNewUserProfile creates a new user profile
func (s *profileService) CreateNewUserProfile(ctx context.Context, userID uuid.UUID, email string, initialOrgID *int64, initialRoleID int) (*domain.Profile, error) {
	// Business logic: e.g., assign default role/org if not provided, validate
	// For MVP, we might assume initialRoleID is provided.
	// If initialOrgID is nil, a new organization might need to be created, or it's an error.
	// This depends on your signup flow logic.

	now := time.Now()
	profile := &domain.Profile{
		ID:             userID,
		Email:          email,
		OrganizationID: initialOrgID,
		RoleID:         initialRoleID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err := s.profileRepo.CreateProfile(ctx, profile)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// GetProfileByID retrieves a profile by ID
func (s *profileService) GetProfileByID(ctx context.Context, id uuid.UUID) (*domain.Profile, error) {
	return s.profileRepo.GetProfileByID(ctx, id)
}

// UpdateProfile updates an existing profile
func (s *profileService) UpdateProfile(ctx context.Context, profile *domain.Profile) (*domain.Profile, error) {
	// Update the updated_at timestamp
	profile.UpdatedAt = time.Now()
	
	// Perform the update
	err := s.profileRepo.UpdateProfile(ctx, profile)
	if err != nil {
		return nil, err
	}
	
	// Return the updated profile
	return profile, nil
}

// DeleteProfile deletes a profile by ID
func (s *profileService) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	return s.profileRepo.DeleteProfile(ctx, id)
}

// GetRoleByID retrieves a role by ID
func (s *profileService) GetRoleByID(ctx context.Context, roleID int) (*domain.Role, error) {
	return s.profileRepo.GetRoleByID(ctx, roleID)
}

// UpsertProfile creates or updates a profile
func (s *profileService) UpsertProfile(ctx context.Context, userID uuid.UUID, email string, orgID *int64, roleID int, firstName, lastName *string) (*domain.Profile, error) {
	// Check if profile exists
	existingProfile, err := s.profileRepo.GetProfileByID(ctx, userID)
	
	now := time.Now()
	var profile *domain.Profile
	
	if err == nil {
		// Profile exists, update it
		profile = existingProfile
		profile.OrganizationID = orgID
		profile.RoleID = roleID
		profile.Email = email
		profile.FirstName = firstName
		profile.LastName = lastName
		profile.UpdatedAt = now
	} else {
		// Profile doesn't exist or error occurred
		profile = &domain.Profile{
			ID:             userID,
			Email:          email,
			OrganizationID: orgID,
			RoleID:         roleID,
			FirstName:      firstName,
			LastName:       lastName,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
	}
	
	// Use the repository's upsert method
	err = s.profileRepo.UpsertProfile(ctx, profile)
	if err != nil {
		return nil, err
	}
	
	return profile, nil
}