package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/crypto"
	"github.com/affiliate-backend/internal/repository"
)

type AdvertiserService interface {
	CreateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) (*domain.Advertiser, error)
	GetAdvertiserByID(ctx context.Context, id int64) (*domain.Advertiser, error)
	GetAdvertiserWithProviderData(ctx context.Context, id int64) (*domain.AdvertiserWithProviderData, error)
	UpdateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error
	ListAdvertisersByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]*domain.Advertiser, error)
	DeleteAdvertiser(ctx context.Context, id int64) error
	
	SyncAdvertiserToProvider(ctx context.Context, advertiserID int64) error
	SyncAdvertiserFromProvider(ctx context.Context, advertiserID int64) error
	CompareAdvertiserWithProvider(ctx context.Context, advertiserID int64) ([]domain.AdvertiserDiscrepancy, error)
	
	CreateAdvertiserProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) (*domain.AdvertiserProviderMapping, error)
	GetAdvertiserProviderMapping(ctx context.Context, advertiserID int64, providerType string) (*domain.AdvertiserProviderMapping, error)
	UpdateAdvertiserProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error
	DeleteAdvertiserProviderMapping(ctx context.Context, mappingID int64) error
}

type advertiserService struct {
	advertiserRepo      repository.AdvertiserRepository
	providerMappingRepo repository.AdvertiserProviderMappingRepository
	orgRepo             repository.OrganizationRepository
	cryptoService       crypto.Service
}

func NewAdvertiserService(
	advertiserRepo repository.AdvertiserRepository,
	providerMappingRepo repository.AdvertiserProviderMappingRepository,
	orgRepo repository.OrganizationRepository,
	cryptoService crypto.Service,
) AdvertiserService {
	return &advertiserService{
		advertiserRepo:      advertiserRepo,
		providerMappingRepo: providerMappingRepo,
		orgRepo:             orgRepo,
		cryptoService:       cryptoService,
	}
}

func (s *advertiserService) CreateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) (*domain.Advertiser, error) {
	_, err := s.orgRepo.GetOrganizationByID(ctx, advertiser.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	setDefaultStatus(advertiser)
	
	if err := validateAdvertiser(advertiser); err != nil {
		return nil, err
	}

	syncStatus := "not_synced"
	advertiser.EverflowSyncStatus = &syncStatus

	if err := s.advertiserRepo.CreateAdvertiser(ctx, advertiser); err != nil {
		return nil, fmt.Errorf("failed to create advertiser: %w", err)
	}

	// TODO: Add provider sync using IntegrationService if needed

	return advertiser, nil
}

// GetAdvertiserByID retrieves an advertiser by ID
func (s *advertiserService) GetAdvertiserByID(ctx context.Context, id int64) (*domain.Advertiser, error) {
	return s.advertiserRepo.GetAdvertiserByID(ctx, id)
}

// UpdateAdvertiser updates an advertiser with Everflow synchronization
func (s *advertiserService) UpdateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error {
	if err := validateAdvertiser(advertiser); err != nil {
		return err
	}

	if err := s.advertiserRepo.UpdateAdvertiser(ctx, advertiser); err != nil {
		return fmt.Errorf("failed to update advertiser: %w", err)
	}

	// TODO: Add provider sync using IntegrationService if needed

	return nil
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

	if err := s.providerMappingRepo.CreateAdvertiserProviderMapping(ctx, mapping); err != nil {
		return nil, fmt.Errorf("failed to create advertiser provider mapping: %w", err)
	}

	return mapping, nil
}

// GetAdvertiserProviderMapping retrieves an advertiser provider mapping
func (s *advertiserService) GetAdvertiserProviderMapping(ctx context.Context, advertiserID int64, providerType string) (*domain.AdvertiserProviderMapping, error) {
	return s.providerMappingRepo.GetAdvertiserProviderMapping(ctx, advertiserID, providerType)
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

	return s.providerMappingRepo.UpdateAdvertiserProviderMapping(ctx, mapping)
}

// DeleteAdvertiserProviderMapping deletes an advertiser provider mapping
func (s *advertiserService) DeleteAdvertiserProviderMapping(ctx context.Context, mappingID int64) error {
	return s.providerMappingRepo.DeleteAdvertiserProviderMapping(ctx, mappingID)
}

func (s *advertiserService) GetAdvertiserWithProviderData(ctx context.Context, id int64) (*domain.AdvertiserWithProviderData, error) {
	advertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get advertiser: %w", err)
	}
	return &domain.AdvertiserWithProviderData{
		Advertiser: advertiser,
		SyncStatus: "not_implemented",
	}, nil
}

func (s *advertiserService) SyncAdvertiserToProvider(ctx context.Context, advertiserID int64) error {
	// TODO: Implement using IntegrationService
	return fmt.Errorf("sync to provider not implemented")
}

func (s *advertiserService) SyncAdvertiserFromProvider(ctx context.Context, advertiserID int64) error {
	// TODO: Implement using IntegrationService
	return fmt.Errorf("sync from provider not implemented")
}

func (s *advertiserService) CompareAdvertiserWithProvider(ctx context.Context, advertiserID int64) ([]domain.AdvertiserDiscrepancy, error) {
	// TODO: Implement using IntegrationService
	return nil, fmt.Errorf("compare with provider not implemented")
}
