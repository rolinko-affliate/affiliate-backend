package service

import (
	"context"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/repository"
)

// CreateOrganizationWithExtraInfoRequest represents the request for creating an organization with extra info
type CreateOrganizationWithExtraInfoRequest struct {
	Name                string
	Type                domain.OrganizationType
	ContactEmail        string
	Description         string
	AdvertiserExtraInfo *domain.AdvertiserExtraInfo
	AffiliateExtraInfo  *domain.AffiliateExtraInfo
}

// OrganizationService defines the interface for organization operations
type OrganizationService interface {
	CreateOrganization(ctx context.Context, name string, orgType domain.OrganizationType) (*domain.Organization, error)
	CreateOrganizationWithExtraInfo(ctx context.Context, req *CreateOrganizationWithExtraInfoRequest) (*domain.Organization, error)
	GetOrganizationByID(ctx context.Context, id int64) (*domain.Organization, error)
	UpdateOrganization(ctx context.Context, org *domain.Organization) error
	ListOrganizations(ctx context.Context, page, pageSize int) ([]*domain.Organization, error)
	DeleteOrganization(ctx context.Context, id int64) error
}

// organizationService implements OrganizationService
type organizationService struct {
	orgRepo        repository.OrganizationRepository
	advertiserRepo repository.AdvertiserRepository
	affiliateRepo  repository.AffiliateRepository
}

// NewOrganizationService creates a new organization service
func NewOrganizationService(orgRepo repository.OrganizationRepository, advertiserRepo repository.AdvertiserRepository, affiliateRepo repository.AffiliateRepository) OrganizationService {
	return &organizationService{
		orgRepo:        orgRepo,
		advertiserRepo: advertiserRepo,
		affiliateRepo:  affiliateRepo,
	}
}

// CreateOrganization creates a new organization
func (s *organizationService) CreateOrganization(ctx context.Context, name string, orgType domain.OrganizationType) (*domain.Organization, error) {
	if name == "" {
		return nil, fmt.Errorf("organization name cannot be empty")
	}
	if orgType == "" {
		return nil, fmt.Errorf("organization type cannot be empty")
	}

	// Validate organization type using the enum's IsValid method
	if !orgType.IsValid() {
		return nil, fmt.Errorf("invalid organization type: %s. Valid types are: %s, %s, %s",
			orgType,
			domain.OrganizationTypeAdvertiser,
			domain.OrganizationTypeAffiliate,
			domain.OrganizationTypePlatformOwner)
	}

	org := &domain.Organization{
		Name: name,
		Type: orgType,
	}

	if err := s.orgRepo.CreateOrganization(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	return org, nil
}

// CreateOrganizationWithExtraInfo creates a new organization with extra info
func (s *organizationService) CreateOrganizationWithExtraInfo(ctx context.Context, req *CreateOrganizationWithExtraInfoRequest) (*domain.Organization, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("organization name cannot be empty")
	}
	if req.Type == "" {
		return nil, fmt.Errorf("organization type cannot be empty")
	}

	// Validate organization type
	if !req.Type.IsValid() {
		return nil, fmt.Errorf("invalid organization type: %s", req.Type)
	}

	// Create the organization first
	org := &domain.Organization{
		Name: req.Name,
		Type: req.Type,
	}

	if err := s.orgRepo.CreateOrganization(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// Handle extra info based on organization type - link directly to organization
	switch req.Type {
	case domain.OrganizationTypeAdvertiser:
		if req.AdvertiserExtraInfo != nil {
			// Create advertiser extra info directly linked to organization
			req.AdvertiserExtraInfo.OrganizationID = org.OrganizationID
			if err := s.advertiserRepo.CreateAdvertiserExtraInfo(ctx, req.AdvertiserExtraInfo); err != nil {
				return nil, fmt.Errorf("failed to create advertiser extra info: %w", err)
			}
		}

	case domain.OrganizationTypeAffiliate:
		if req.AffiliateExtraInfo != nil {
			// Create affiliate extra info directly linked to organization
			req.AffiliateExtraInfo.OrganizationID = org.OrganizationID
			if err := s.affiliateRepo.CreateAffiliateExtraInfo(ctx, req.AffiliateExtraInfo); err != nil {
				return nil, fmt.Errorf("failed to create affiliate extra info: %w", err)
			}
		}
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
	if org.Type == "" {
		return fmt.Errorf("organization type cannot be empty")
	}

	// Validate organization type using the enum's IsValid method
	if !org.Type.IsValid() {
		return fmt.Errorf("invalid organization type: %s. Valid types are: %s, %s, %s",
			org.Type,
			domain.OrganizationTypeAdvertiser,
			domain.OrganizationTypeAffiliate,
			domain.OrganizationTypePlatformOwner)
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
