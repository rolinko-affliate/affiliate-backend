package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/crypto"
	"github.com/affiliate-backend/internal/platform/provider"
	"github.com/affiliate-backend/internal/repository"
	"github.com/google/uuid"
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
	integrationService  provider.IntegrationService
}

func NewAdvertiserService(
	advertiserRepo repository.AdvertiserRepository,
	providerMappingRepo repository.AdvertiserProviderMappingRepository,
	orgRepo repository.OrganizationRepository,
	cryptoService crypto.Service,
	integrationService provider.IntegrationService,
) AdvertiserService {
	return &advertiserService{
		advertiserRepo:      advertiserRepo,
		providerMappingRepo: providerMappingRepo,
		orgRepo:             orgRepo,
		cryptoService:       cryptoService,
		integrationService:  integrationService,
	}
}

func (s *advertiserService) CreateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) (*domain.Advertiser, error) {
	// Validate organization exists
	_, err := s.orgRepo.GetOrganizationByID(ctx, advertiser.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	setDefaultStatus(advertiser)
	
	if err := validateAdvertiser(advertiser); err != nil {
		return nil, err
	}

	// Set status to "creating" during provider sync
	syncStatus := "creating"
	advertiser.EverflowSyncStatus = &syncStatus

	// Step 1: Insert local record with status = "creating"
	if err := s.advertiserRepo.CreateAdvertiser(ctx, advertiser); err != nil {
		return nil, fmt.Errorf("failed to create advertiser: %w", err)
	}

	// Step 2: Call IntegrationService to create in provider
	providerAdvertiser, err := s.integrationService.CreateAdvertiser(ctx, *advertiser)
	if err != nil {
		// Rollback: update status to "failed"
		failedStatus := "failed"
		advertiser.EverflowSyncStatus = &failedStatus
		s.advertiserRepo.UpdateAdvertiser(ctx, advertiser)
		return nil, fmt.Errorf("failed to create advertiser in provider: %w", err)
	}

	// Step 3: Create provider mapping with provider ID and payload
	var providerID *string
	if providerAdvertiser.NetworkEmployeeID != nil {
		idStr := fmt.Sprintf("%d", *providerAdvertiser.NetworkEmployeeID)
		providerID = &idStr
	}
	mapping := &domain.AdvertiserProviderMapping{
		AdvertiserID:         advertiser.AdvertiserID,
		ProviderType:         "everflow",
		ProviderAdvertiserID: providerID,
		APICredentials:       nil, // Set by IntegrationService
		ProviderConfig:       nil, // Set by IntegrationService with full payload
	}

	if err := s.providerMappingRepo.CreateAdvertiserProviderMapping(ctx, mapping); err != nil {
		// Log error but don't fail the operation since advertiser was created in provider
		fmt.Printf("Warning: failed to create provider mapping for advertiser %d: %v\n", advertiser.AdvertiserID, err)
	}

	// Step 4: Update status to "active"
	activeStatus := "active"
	advertiser.EverflowSyncStatus = &activeStatus
	if err := s.advertiserRepo.UpdateAdvertiser(ctx, advertiser); err != nil {
		// Log error but don't fail since advertiser was created successfully
		fmt.Printf("Warning: failed to update advertiser status to active: %v\n", err)
	}

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

	// Step 1: Update local record first
	if err := s.advertiserRepo.UpdateAdvertiser(ctx, advertiser); err != nil {
		return fmt.Errorf("failed to update advertiser: %w", err)
	}

	// Step 2: Check if provider mapping exists
	_, err := s.providerMappingRepo.GetAdvertiserProviderMapping(ctx, advertiser.AdvertiserID, "everflow")
	if err != nil {
		// No provider mapping exists, skip provider sync
		return nil
	}

	// Step 3: Update in provider if mapping exists
	if err := s.integrationService.UpdateAdvertiser(ctx, *advertiser); err != nil {
		// Log error but don't fail the operation since local update succeeded
		fmt.Printf("Warning: failed to update advertiser in provider: %v\n", err)
		
		// Update sync status to indicate sync failure
		syncStatus := "sync_failed"
		advertiser.EverflowSyncStatus = &syncStatus
		s.advertiserRepo.UpdateAdvertiser(ctx, advertiser)
	} else {
		// Update sync status to indicate successful sync
		syncStatus := "active"
		advertiser.EverflowSyncStatus = &syncStatus
		s.advertiserRepo.UpdateAdvertiser(ctx, advertiser)
	}

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
	// Get local advertiser
	advertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, advertiserID)
	if err != nil {
		return fmt.Errorf("failed to get advertiser: %w", err)
	}

	// Check if provider mapping exists
	_, err = s.providerMappingRepo.GetAdvertiserProviderMapping(ctx, advertiserID, "everflow")
	if err != nil {
		// No mapping exists, create in provider
		return s.createAdvertiserInProvider(ctx, advertiser)
	}

	// Mapping exists, update in provider
	if err := s.integrationService.UpdateAdvertiser(ctx, *advertiser); err != nil {
		syncStatus := "sync_failed"
		advertiser.EverflowSyncStatus = &syncStatus
		s.advertiserRepo.UpdateAdvertiser(ctx, advertiser)
		return fmt.Errorf("failed to sync advertiser to provider: %w", err)
	}

	// Update sync status
	syncStatus := "active"
	advertiser.EverflowSyncStatus = &syncStatus
	return s.advertiserRepo.UpdateAdvertiser(ctx, advertiser)
}

func (s *advertiserService) SyncAdvertiserFromProvider(ctx context.Context, advertiserID int64) error {
	// Get provider mapping
	_, err := s.providerMappingRepo.GetAdvertiserProviderMapping(ctx, advertiserID, "everflow")
	if err != nil {
		return fmt.Errorf("no provider mapping found for advertiser %d: %w", advertiserID, err)
	}

	// Convert advertiser ID to UUID for IntegrationService
	advertiserUUID := int64ToUUID(advertiserID)
	
	// Get advertiser from provider
	providerAdvertiser, err := s.integrationService.GetAdvertiser(ctx, advertiserUUID)
	if err != nil {
		return fmt.Errorf("failed to get advertiser from provider: %w", err)
	}

	// Update local advertiser with provider data
	localAdvertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, advertiserID)
	if err != nil {
		return fmt.Errorf("failed to get local advertiser: %w", err)
	}

	// Merge provider data into local advertiser
	s.mergeProviderDataIntoAdvertiser(localAdvertiser, &providerAdvertiser)

	// Update local record
	syncStatus := "active"
	localAdvertiser.EverflowSyncStatus = &syncStatus
	return s.advertiserRepo.UpdateAdvertiser(ctx, localAdvertiser)
}

func (s *advertiserService) CompareAdvertiserWithProvider(ctx context.Context, advertiserID int64) ([]domain.AdvertiserDiscrepancy, error) {
	// Get local advertiser
	localAdvertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, advertiserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get local advertiser: %w", err)
	}

	// Get provider mapping
	_, err = s.providerMappingRepo.GetAdvertiserProviderMapping(ctx, advertiserID, "everflow")
	if err != nil {
		return []domain.AdvertiserDiscrepancy{{
			Field:         "provider_mapping",
			LocalValue:    "exists",
			ProviderValue: "missing",
			Severity:      "critical",
		}}, nil
	}

	// Convert advertiser ID to UUID for IntegrationService
	advertiserUUID := int64ToUUID(advertiserID)
	
	// Get advertiser from provider
	providerAdvertiser, err := s.integrationService.GetAdvertiser(ctx, advertiserUUID)
	if err != nil {
		return []domain.AdvertiserDiscrepancy{{
			Field:         "provider_record",
			LocalValue:    "exists",
			ProviderValue: "missing",
			Severity:      "critical",
		}}, nil
	}

	// Compare fields and return discrepancies
	return s.compareAdvertiserFields(localAdvertiser, &providerAdvertiser), nil
}

// Helper methods

// createAdvertiserInProvider creates an advertiser in the provider when no mapping exists
func (s *advertiserService) createAdvertiserInProvider(ctx context.Context, advertiser *domain.Advertiser) error {
	// Create in provider
	providerAdvertiser, err := s.integrationService.CreateAdvertiser(ctx, *advertiser)
	if err != nil {
		syncStatus := "sync_failed"
		advertiser.EverflowSyncStatus = &syncStatus
		s.advertiserRepo.UpdateAdvertiser(ctx, advertiser)
		return fmt.Errorf("failed to create advertiser in provider: %w", err)
	}

	// Create provider mapping
	var providerID *string
	if providerAdvertiser.NetworkEmployeeID != nil {
		idStr := fmt.Sprintf("%d", *providerAdvertiser.NetworkEmployeeID)
		providerID = &idStr
	}
	mapping := &domain.AdvertiserProviderMapping{
		AdvertiserID:         advertiser.AdvertiserID,
		ProviderType:         "everflow",
		ProviderAdvertiserID: providerID,
		APICredentials:       nil, // Set by IntegrationService
		ProviderConfig:       nil, // Set by IntegrationService with full payload
	}

	if err := s.providerMappingRepo.CreateAdvertiserProviderMapping(ctx, mapping); err != nil {
		fmt.Printf("Warning: failed to create provider mapping for advertiser %d: %v\n", advertiser.AdvertiserID, err)
	}

	// Update sync status
	syncStatus := "active"
	advertiser.EverflowSyncStatus = &syncStatus
	return s.advertiserRepo.UpdateAdvertiser(ctx, advertiser)
}

// mergeProviderDataIntoAdvertiser merges provider data into local advertiser
func (s *advertiserService) mergeProviderDataIntoAdvertiser(local *domain.Advertiser, provider *domain.Advertiser) {
	// Merge relevant fields from provider into local
	if provider.NetworkEmployeeID != nil {
		local.NetworkEmployeeID = provider.NetworkEmployeeID
	}
	// Add other fields as needed based on what the provider returns
}

// compareAdvertiserFields compares local and provider advertiser fields
func (s *advertiserService) compareAdvertiserFields(local *domain.Advertiser, provider *domain.Advertiser) []domain.AdvertiserDiscrepancy {
	var discrepancies []domain.AdvertiserDiscrepancy

	// Compare name
	if local.Name != provider.Name {
		discrepancies = append(discrepancies, domain.AdvertiserDiscrepancy{
			Field:         "name",
			LocalValue:    local.Name,
			ProviderValue: provider.Name,
			Severity:      "medium",
		})
	}

	// Compare contact email
	if (local.ContactEmail == nil) != (provider.ContactEmail == nil) ||
		(local.ContactEmail != nil && provider.ContactEmail != nil && *local.ContactEmail != *provider.ContactEmail) {
		discrepancies = append(discrepancies, domain.AdvertiserDiscrepancy{
			Field:         "contact_email",
			LocalValue:    local.ContactEmail,
			ProviderValue: provider.ContactEmail,
			Severity:      "medium",
		})
	}

	// Compare status
	if local.Status != provider.Status {
		discrepancies = append(discrepancies, domain.AdvertiserDiscrepancy{
			Field:         "status",
			LocalValue:    local.Status,
			ProviderValue: provider.Status,
			Severity:      "high",
		})
	}

	// Compare NetworkEmployeeID if both exist
	if local.NetworkEmployeeID != nil && provider.NetworkEmployeeID != nil {
		if *local.NetworkEmployeeID != *provider.NetworkEmployeeID {
			discrepancies = append(discrepancies, domain.AdvertiserDiscrepancy{
				Field:         "network_employee_id",
				LocalValue:    *local.NetworkEmployeeID,
				ProviderValue: *provider.NetworkEmployeeID,
				Severity:      "critical",
			})
		}
	}

	return discrepancies
}

// UUID conversion helpers (copied from IntegrationService)
func int64ToUUID(id int64) uuid.UUID {
	// Convert int64 back to UUID format
	// This is a simplified approach - in production you might want a more sophisticated mapping
	hex := fmt.Sprintf("%015x", id)
	// Pad to 32 characters
	for len(hex) < 32 {
		hex = "0" + hex
	}
	// Format as UUID
	uuidStr := fmt.Sprintf("%s-%s-%s-%s-%s", hex[:8], hex[8:12], hex[12:16], hex[16:20], hex[20:32])
	parsed, _ := uuid.Parse(uuidStr)
	return parsed
}

// setDefaultStatus sets default status for advertiser if not provided
func setDefaultStatus(advertiser *domain.Advertiser) {
	if advertiser.Status == "" {
		advertiser.Status = "pending"
	}
}

// validateAdvertiser validates advertiser fields
func validateAdvertiser(advertiser *domain.Advertiser) error {
	if advertiser.Name == "" {
		return fmt.Errorf("advertiser name cannot be empty")
	}
	
	if advertiser.ContactEmail == nil || *advertiser.ContactEmail == "" {
		return fmt.Errorf("advertiser contact email cannot be empty")
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
	
	return nil
}
