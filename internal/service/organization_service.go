package service

import (
	"context"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/repository"
)

// OrganizationService defines the interface for organization operations
type OrganizationService interface {
	CreateOrganization(ctx context.Context, name string) (*domain.Organization, error)
	GetOrganizationByID(ctx context.Context, id int64) (*domain.Organization, error)
	UpdateOrganization(ctx context.Context, org *domain.Organization) error
	ListOrganizations(ctx context.Context, page, pageSize int) ([]*domain.Organization, error)
	DeleteOrganization(ctx context.Context, id int64) error
}

// organizationService implements OrganizationService
type organizationService struct {
	orgRepo repository.OrganizationRepository
}

// NewOrganizationService creates a new organization service
func NewOrganizationService(orgRepo repository.OrganizationRepository) OrganizationService {
	return &organizationService{orgRepo: orgRepo}
}

// CreateOrganization creates a new organization
func (s *organizationService) CreateOrganization(ctx context.Context, name string) (*domain.Organization, error) {
	if name == "" {
		return nil, fmt.Errorf("organization name cannot be empty")
	}

	org := &domain.Organization{
		Name: name,
	}

	if err := s.orgRepo.CreateOrganization(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	return org, nil
}

// GetOrganizationByID retrieves an organization by ID
func (s *organizationService) GetOrganizationByID(ctx context.Context, id int64) (*domain.Organization, error) {
	return s.orgRepo.GetOrganizationByID(ctx, id)
}

// UpdateOrganization updates an organization
func (s *organizationService) UpdateOrganization(ctx context.Context, org *domain.Organization) error {
	if org.Name == "" {
		return fmt.Errorf("organization name cannot be empty")
	}

	return s.orgRepo.UpdateOrganization(ctx, org)
}

// ListOrganizations retrieves a list of organizations with pagination
func (s *organizationService) ListOrganizations(ctx context.Context, page, pageSize int) ([]*domain.Organization, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.orgRepo.ListOrganizations(ctx, pageSize, offset)
}

// DeleteOrganization deletes an organization
func (s *organizationService) DeleteOrganization(ctx context.Context, id int64) error {
	return s.orgRepo.DeleteOrganization(ctx, id)
}