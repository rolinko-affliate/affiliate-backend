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
	GetRoleByID(ctx context.Context, roleID int) (*domain.Role, error)
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

// GetRoleByID retrieves a role by ID
func (s *profileService) GetRoleByID(ctx context.Context, roleID int) (*domain.Role, error) {
	return s.profileRepo.GetRoleByID(ctx, roleID)
}