package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/crypto"
	"github.com/affiliate-backend/internal/platform/everflow"
	"github.com/affiliate-backend/internal/repository"
)

// AdvertiserService defines the interface for advertiser operations
type AdvertiserService interface {
	CreateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) (*domain.Advertiser, error)
	GetAdvertiserByID(ctx context.Context, id int64) (*domain.Advertiser, error)
	UpdateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error
	ListAdvertisersByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]*domain.Advertiser, error)
	DeleteAdvertiser(ctx context.Context, id int64) error
	
	// Provider mapping methods
	CreateAdvertiserProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) (*domain.AdvertiserProviderMapping, error)
	GetAdvertiserProviderMapping(ctx context.Context, advertiserID int64, providerType string) (*domain.AdvertiserProviderMapping, error)
	UpdateAdvertiserProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error
	DeleteAdvertiserProviderMapping(ctx context.Context, mappingID int64) error
}

// advertiserService implements AdvertiserService
type advertiserService struct {
	advertiserRepo repository.AdvertiserRepository
	orgRepo        repository.OrganizationRepository
	everflowService *everflow.Service
	cryptoService  crypto.Service
}

// NewAdvertiserService creates a new advertiser service
func NewAdvertiserService(
	advertiserRepo repository.AdvertiserRepository, 
	orgRepo repository.OrganizationRepository,
	everflowService *everflow.Service,
	cryptoService crypto.Service,
) AdvertiserService {
	return &advertiserService{
		advertiserRepo:  advertiserRepo,
		orgRepo:         orgRepo,
		everflowService: everflowService,
		cryptoService:   cryptoService,
	}
}

// CreateAdvertiser creates a new advertiser
func (s *advertiserService) CreateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) (*domain.Advertiser, error) {
	// Validate organization exists
	_, err := s.orgRepo.GetOrganizationByID(ctx, advertiser.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	// Validate required fields
	if advertiser.Name == "" {
		return nil, fmt.Errorf("advertiser name cannot be empty")
	}

	// Set default status if not provided
	if advertiser.Status == "" {
		advertiser.Status = "pending"
	}

	// Validate status
	validStatuses := map[string]bool{
		"active":   true,
		"pending":  true,
		"inactive": true,
		"rejected": true,
	}
	if !validStatuses[advertiser.Status] {
		return nil, fmt.Errorf("invalid status: %s", advertiser.Status)
	}

	// Validate billing details JSON if provided
	if advertiser.BillingDetails != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*advertiser.BillingDetails), &jsonData); err != nil {
			return nil, fmt.Errorf("invalid billing details JSON: %w", err)
		}
	}

	if err := s.advertiserRepo.CreateAdvertiser(ctx, advertiser); err != nil {
		return nil, fmt.Errorf("failed to create advertiser: %w", err)
	}

	// Create advertiser in Everflow if the service is available
	if s.everflowService != nil {
		go func() {
			// Use a background context since this is a fire-and-forget operation
			bgCtx := context.Background()
			if err := s.everflowService.CreateAdvertiserInEverflow(bgCtx, advertiser); err != nil {
				// Log the error but don't fail the advertiser creation
				log.Printf("Error creating advertiser in Everflow: %v", err)
			}
		}()
	}

	return advertiser, nil
}

// GetAdvertiserByID retrieves an advertiser by ID
func (s *advertiserService) GetAdvertiserByID(ctx context.Context, id int64) (*domain.Advertiser, error) {
	return s.advertiserRepo.GetAdvertiserByID(ctx, id)
}

// UpdateAdvertiser updates an advertiser
func (s *advertiserService) UpdateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error {
	// Validate required fields
	if advertiser.Name == "" {
		return fmt.Errorf("advertiser name cannot be empty")
	}

	// Validate status
	validStatuses := map[string]bool{
		"active":   true,
		"pending":  true,
		"inactive": true,
		"rejected": true,
	}
	if !validStatuses[advertiser.Status] {
		return fmt.Errorf("invalid status: %s", advertiser.Status)
	}

	// Validate billing details JSON if provided
	if advertiser.BillingDetails != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*advertiser.BillingDetails), &jsonData); err != nil {
			return fmt.Errorf("invalid billing details JSON: %w", err)
		}
	}

	return s.advertiserRepo.UpdateAdvertiser(ctx, advertiser)
}

// ListAdvertisersByOrganization retrieves a list of advertisers for an organization with pagination
func (s *advertiserService) ListAdvertisersByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]*domain.Advertiser, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.advertiserRepo.ListAdvertisersByOrganization(ctx, orgID, pageSize, offset)
}

// DeleteAdvertiser deletes an advertiser
func (s *advertiserService) DeleteAdvertiser(ctx context.Context, id int64) error {
	return s.advertiserRepo.DeleteAdvertiser(ctx, id)
}

// CreateAdvertiserProviderMapping creates a new advertiser provider mapping
func (s *advertiserService) CreateAdvertiserProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) (*domain.AdvertiserProviderMapping, error) {
	// Validate advertiser exists
	_, err := s.advertiserRepo.GetAdvertiserByID(ctx, mapping.AdvertiserID)
	if err != nil {
		return nil, fmt.Errorf("advertiser not found: %w", err)
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

	// Validate API credentials JSON if provided
	if mapping.APICredentials != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*mapping.APICredentials), &jsonData); err != nil {
			return nil, fmt.Errorf("invalid API credentials JSON: %w", err)
		}
	}

	if err := s.advertiserRepo.CreateAdvertiserProviderMapping(ctx, mapping); err != nil {
		return nil, fmt.Errorf("failed to create advertiser provider mapping: %w", err)
	}

	return mapping, nil
}

// GetAdvertiserProviderMapping retrieves an advertiser provider mapping
func (s *advertiserService) GetAdvertiserProviderMapping(ctx context.Context, advertiserID int64, providerType string) (*domain.AdvertiserProviderMapping, error) {
	return s.advertiserRepo.GetAdvertiserProviderMapping(ctx, advertiserID, providerType)
}

// UpdateAdvertiserProviderMapping updates an advertiser provider mapping
func (s *advertiserService) UpdateAdvertiserProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error {
	// Validate provider config JSON if provided
	if mapping.ProviderConfig != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*mapping.ProviderConfig), &jsonData); err != nil {
			return fmt.Errorf("invalid provider config JSON: %w", err)
		}
	}

	// Validate API credentials JSON if provided
	if mapping.APICredentials != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*mapping.APICredentials), &jsonData); err != nil {
			return fmt.Errorf("invalid API credentials JSON: %w", err)
		}
	}

	return s.advertiserRepo.UpdateAdvertiserProviderMapping(ctx, mapping)
}

// DeleteAdvertiserProviderMapping deletes an advertiser provider mapping
func (s *advertiserService) DeleteAdvertiserProviderMapping(ctx context.Context, mappingID int64) error {
	return s.advertiserRepo.DeleteAdvertiserProviderMapping(ctx, mappingID)
}