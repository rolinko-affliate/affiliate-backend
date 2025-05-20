package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/repository"
)

// AffiliateService defines the interface for affiliate operations
type AffiliateService interface {
	CreateAffiliate(ctx context.Context, affiliate *domain.Affiliate) (*domain.Affiliate, error)
	GetAffiliateByID(ctx context.Context, id int64) (*domain.Affiliate, error)
	UpdateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error
	ListAffiliatesByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]*domain.Affiliate, error)
	DeleteAffiliate(ctx context.Context, id int64) error
	
	// Provider mapping methods
	CreateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) (*domain.AffiliateProviderMapping, error)
	GetAffiliateProviderMapping(ctx context.Context, affiliateID int64, providerType string) (*domain.AffiliateProviderMapping, error)
	UpdateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error
	DeleteAffiliateProviderMapping(ctx context.Context, mappingID int64) error
}

// affiliateService implements AffiliateService
type affiliateService struct {
	affiliateRepo repository.AffiliateRepository
	orgRepo       repository.OrganizationRepository
}

// NewAffiliateService creates a new affiliate service
func NewAffiliateService(affiliateRepo repository.AffiliateRepository, orgRepo repository.OrganizationRepository) AffiliateService {
	return &affiliateService{
		affiliateRepo: affiliateRepo,
		orgRepo:       orgRepo,
	}
}

// CreateAffiliate creates a new affiliate
func (s *affiliateService) CreateAffiliate(ctx context.Context, affiliate *domain.Affiliate) (*domain.Affiliate, error) {
	// Validate organization exists
	_, err := s.orgRepo.GetOrganizationByID(ctx, affiliate.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	// Validate required fields
	if affiliate.Name == "" {
		return nil, fmt.Errorf("affiliate name cannot be empty")
	}

	// Set default status if not provided
	if affiliate.Status == "" {
		affiliate.Status = "pending"
	}

	// Validate status
	validStatuses := map[string]bool{
		"active":   true,
		"pending":  true,
		"inactive": true,
		"rejected": true,
	}
	if !validStatuses[affiliate.Status] {
		return nil, fmt.Errorf("invalid status: %s", affiliate.Status)
	}

	// Validate payment details JSON if provided
	if affiliate.PaymentDetails != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*affiliate.PaymentDetails), &jsonData); err != nil {
			return nil, fmt.Errorf("invalid payment details JSON: %w", err)
		}
	}

	if err := s.affiliateRepo.CreateAffiliate(ctx, affiliate); err != nil {
		return nil, fmt.Errorf("failed to create affiliate: %w", err)
	}

	return affiliate, nil
}

// GetAffiliateByID retrieves an affiliate by ID
func (s *affiliateService) GetAffiliateByID(ctx context.Context, id int64) (*domain.Affiliate, error) {
	return s.affiliateRepo.GetAffiliateByID(ctx, id)
}

// UpdateAffiliate updates an affiliate
func (s *affiliateService) UpdateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error {
	// Validate required fields
	if affiliate.Name == "" {
		return fmt.Errorf("affiliate name cannot be empty")
	}

	// Validate status
	validStatuses := map[string]bool{
		"active":   true,
		"pending":  true,
		"inactive": true,
		"rejected": true,
	}
	if !validStatuses[affiliate.Status] {
		return fmt.Errorf("invalid status: %s", affiliate.Status)
	}

	// Validate payment details JSON if provided
	if affiliate.PaymentDetails != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*affiliate.PaymentDetails), &jsonData); err != nil {
			return fmt.Errorf("invalid payment details JSON: %w", err)
		}
	}

	return s.affiliateRepo.UpdateAffiliate(ctx, affiliate)
}

// ListAffiliatesByOrganization retrieves a list of affiliates for an organization with pagination
func (s *affiliateService) ListAffiliatesByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]*domain.Affiliate, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.affiliateRepo.ListAffiliatesByOrganization(ctx, orgID, pageSize, offset)
}

// DeleteAffiliate deletes an affiliate
func (s *affiliateService) DeleteAffiliate(ctx context.Context, id int64) error {
	return s.affiliateRepo.DeleteAffiliate(ctx, id)
}

// CreateAffiliateProviderMapping creates a new affiliate provider mapping
func (s *affiliateService) CreateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) (*domain.AffiliateProviderMapping, error) {
	// Validate affiliate exists
	_, err := s.affiliateRepo.GetAffiliateByID(ctx, mapping.AffiliateID)
	if err != nil {
		return nil, fmt.Errorf("affiliate not found: %w", err)
	}

	// Validate provider type
	if mapping.ProviderType != "everflow" {
		return nil, fmt.Errorf("invalid provider type: %s", mapping.ProviderType)
	}

	// Validate provider config JSON if provided
	if mapping.ProviderConfig != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*mapping.ProviderConfig), &jsonData); err != nil {
			return nil, fmt.Errorf("invalid provider config JSON: %w", err)
		}
	}

	if err := s.affiliateRepo.CreateAffiliateProviderMapping(ctx, mapping); err != nil {
		return nil, fmt.Errorf("failed to create affiliate provider mapping: %w", err)
	}

	return mapping, nil
}

// GetAffiliateProviderMapping retrieves an affiliate provider mapping
func (s *affiliateService) GetAffiliateProviderMapping(ctx context.Context, affiliateID int64, providerType string) (*domain.AffiliateProviderMapping, error) {
	return s.affiliateRepo.GetAffiliateProviderMapping(ctx, affiliateID, providerType)
}

// UpdateAffiliateProviderMapping updates an affiliate provider mapping
func (s *affiliateService) UpdateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error {
	// Validate provider config JSON if provided
	if mapping.ProviderConfig != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*mapping.ProviderConfig), &jsonData); err != nil {
			return fmt.Errorf("invalid provider config JSON: %w", err)
		}
	}

	return s.affiliateRepo.UpdateAffiliateProviderMapping(ctx, mapping)
}

// DeleteAffiliateProviderMapping deletes an affiliate provider mapping
func (s *affiliateService) DeleteAffiliateProviderMapping(ctx context.Context, mappingID int64) error {
	return s.affiliateRepo.DeleteAffiliateProviderMapping(ctx, mappingID)
}