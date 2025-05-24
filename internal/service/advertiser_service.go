package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/crypto"
	"github.com/affiliate-backend/internal/platform/everflow"
	"github.com/affiliate-backend/internal/repository"
)

// AdvertiserService defines the interface for advertiser operations
type AdvertiserService interface {
	CreateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) (*domain.Advertiser, error)
	GetAdvertiserByID(ctx context.Context, id int64) (*domain.Advertiser, error)
	GetAdvertiserWithEverflowData(ctx context.Context, id int64) (*domain.AdvertiserWithEverflowData, error)
	UpdateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error
	ListAdvertisersByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]*domain.Advertiser, error)
	DeleteAdvertiser(ctx context.Context, id int64) error
	
	// Everflow synchronization methods
	SyncAdvertiserToEverflow(ctx context.Context, advertiserID int64) error
	SyncAdvertiserFromEverflow(ctx context.Context, advertiserID int64) error
	CompareAdvertiserWithEverflow(ctx context.Context, advertiserID int64) ([]domain.AdvertiserDiscrepancy, error)
	
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

// CreateAdvertiser creates a new advertiser with two-step save and sync approach
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

	// Validate billing details if provided
	if advertiser.BillingDetails != nil {
		// Basic validation of billing details structure
		if advertiser.BillingDetails.BillingFrequency != "" {
			validFrequencies := map[string]bool{
				"weekly":  true,
				"monthly": true,
				"other":   true,
			}
			if !validFrequencies[advertiser.BillingDetails.BillingFrequency] {
				return nil, fmt.Errorf("invalid billing frequency: %s", advertiser.BillingDetails.BillingFrequency)
			}
		}
	}

	// Set initial sync status
	syncStatus := "not_synced"
	advertiser.EverflowSyncStatus = &syncStatus

	// Step 1: Save advertiser locally first
	if err := s.advertiserRepo.CreateAdvertiser(ctx, advertiser); err != nil {
		return nil, fmt.Errorf("failed to create advertiser: %w", err)
	}

	// Step 2: Asynchronously sync to Everflow if service is available
	if s.everflowService != nil {
		go func() {
			// Use a background context for async operation
			bgCtx := context.Background()
			
			// Update sync status to pending
			pendingStatus := "pending"
			advertiser.EverflowSyncStatus = &pendingStatus
			s.advertiserRepo.UpdateAdvertiser(bgCtx, advertiser)

			// Attempt to sync to Everflow
			if err := s.everflowService.CreateAdvertiserInEverflow(bgCtx, advertiser); err != nil {
				// Update sync status to failed and log error
				failedStatus := "failed"
				errorMsg := err.Error()
				advertiser.EverflowSyncStatus = &failedStatus
				advertiser.EverflowSyncError = &errorMsg
				s.advertiserRepo.UpdateAdvertiser(bgCtx, advertiser)
				
				log.Printf("Error creating advertiser %d in Everflow: %v", advertiser.AdvertiserID, err)
			} else {
				// Update sync status to synced
				syncedStatus := "synced"
				now := time.Now()
				advertiser.EverflowSyncStatus = &syncedStatus
				advertiser.LastEverflowSyncAt = &now
				advertiser.EverflowSyncError = nil
				s.advertiserRepo.UpdateAdvertiser(bgCtx, advertiser)
				
				log.Printf("Successfully synced advertiser %d to Everflow", advertiser.AdvertiserID)
			}
		}()
	}

	return advertiser, nil
}

// GetAdvertiserByID retrieves an advertiser by ID
func (s *advertiserService) GetAdvertiserByID(ctx context.Context, id int64) (*domain.Advertiser, error) {
	return s.advertiserRepo.GetAdvertiserByID(ctx, id)
}

// UpdateAdvertiser updates an advertiser with Everflow synchronization
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

	// Validate billing details if provided
	if advertiser.BillingDetails != nil {
		// Basic validation of billing details structure
		if advertiser.BillingDetails.BillingFrequency != "" {
			validFrequencies := map[string]bool{
				"weekly":  true,
				"monthly": true,
				"other":   true,
			}
			if !validFrequencies[advertiser.BillingDetails.BillingFrequency] {
				return fmt.Errorf("invalid billing frequency: %s", advertiser.BillingDetails.BillingFrequency)
			}
		}
	}

	// Step 1: Update advertiser locally first
	if err := s.advertiserRepo.UpdateAdvertiser(ctx, advertiser); err != nil {
		return fmt.Errorf("failed to update advertiser: %w", err)
	}

	// Step 2: Asynchronously sync to Everflow if service is available and advertiser is synced
	if s.everflowService != nil && advertiser.EverflowSyncStatus != nil && *advertiser.EverflowSyncStatus == "synced" {
		go func() {
			// Use a background context for async operation
			bgCtx := context.Background()
			
			// Update sync status to pending
			pendingStatus := "pending"
			advertiser.EverflowSyncStatus = &pendingStatus
			s.advertiserRepo.UpdateAdvertiser(bgCtx, advertiser)
			
			// Attempt to sync to Everflow
			if err := s.SyncAdvertiserToEverflow(bgCtx, advertiser.AdvertiserID); err != nil {
				// Update sync status to failed and log error
				failedStatus := "failed"
				errorMsg := err.Error()
				advertiser.EverflowSyncStatus = &failedStatus
				advertiser.EverflowSyncError = &errorMsg
				s.advertiserRepo.UpdateAdvertiser(bgCtx, advertiser)
				
				log.Printf("Error updating advertiser %d in Everflow: %v", advertiser.AdvertiserID, err)
			} else {
				// Update sync status to synced
				syncedStatus := "synced"
				now := time.Now()
				advertiser.EverflowSyncStatus = &syncedStatus
				advertiser.LastEverflowSyncAt = &now
				advertiser.EverflowSyncError = nil
				s.advertiserRepo.UpdateAdvertiser(bgCtx, advertiser)
				
				log.Printf("Successfully updated advertiser %d in Everflow", advertiser.AdvertiserID)
			}
		}()
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

// GetAdvertiserWithEverflowData retrieves an advertiser with Everflow comparison data
func (s *advertiserService) GetAdvertiserWithEverflowData(ctx context.Context, id int64) (*domain.AdvertiserWithEverflowData, error) {
	// Get local advertiser
	advertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get advertiser: %w", err)
	}

	result := &domain.AdvertiserWithEverflowData{
		Advertiser: advertiser,
		SyncStatus: "not_synced",
	}

	// Check if Everflow service is available
	if s.everflowService == nil {
		result.SyncStatus = "service_unavailable"
		return result, nil
	}

	// Try to get Everflow data
	everflowAdvertiser, err := s.everflowService.GetAdvertiserFromEverflowByMapping(ctx, id, []string{"billing", "settings"})
	if err != nil {
		// Check if it's because the advertiser doesn't exist in Everflow
		if advertiser.EverflowSyncStatus == nil || *advertiser.EverflowSyncStatus == "not_synced" {
			result.SyncStatus = "not_synced"
		} else {
			result.SyncStatus = "error"
			log.Printf("Error fetching advertiser %d from Everflow: %v", id, err)
		}
		return result, nil
	}

	// Set Everflow data
	result.EverflowData = everflowAdvertiser

	// Compare data and detect discrepancies
	discrepancies, err := s.CompareAdvertiserWithEverflow(ctx, id)
	if err != nil {
		log.Printf("Error comparing advertiser %d with Everflow: %v", id, err)
		result.SyncStatus = "comparison_error"
	} else {
		result.Discrepancies = discrepancies
		if len(discrepancies) == 0 {
			result.SyncStatus = "synced"
		} else {
			result.SyncStatus = "out_of_sync"
		}
	}

	return result, nil
}

// SyncAdvertiserToEverflow synchronizes an advertiser to Everflow
func (s *advertiserService) SyncAdvertiserToEverflow(ctx context.Context, advertiserID int64) error {
	if s.everflowService == nil {
		return fmt.Errorf("Everflow service not available")
	}

	// Get the advertiser
	advertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, advertiserID)
	if err != nil {
		return fmt.Errorf("failed to get advertiser: %w", err)
	}

	// Check if advertiser has Everflow mapping
	mapping, err := s.advertiserRepo.GetAdvertiserProviderMapping(ctx, advertiserID, "everflow")
	if err != nil {
		// No mapping exists, create new advertiser in Everflow
		return s.everflowService.CreateAdvertiserInEverflow(ctx, advertiser)
	}

	// Mapping exists, update advertiser in Everflow
	if mapping.ProviderAdvertiserID == nil {
		return fmt.Errorf("advertiser mapping exists but has no provider advertiser ID")
	}

	// Map our advertiser to Everflow update request
	updateReq, err := s.mapAdvertiserToEverflowUpdateRequest(advertiser)
	if err != nil {
		return fmt.Errorf("failed to map advertiser to Everflow update request: %w", err)
	}

	// Update in Everflow
	_, err = s.everflowService.UpdateAdvertiserInEverflowByMapping(ctx, advertiserID, *updateReq)
	return err
}

// SyncAdvertiserFromEverflow synchronizes an advertiser from Everflow to local
func (s *advertiserService) SyncAdvertiserFromEverflow(ctx context.Context, advertiserID int64) error {
	if s.everflowService == nil {
		return fmt.Errorf("Everflow service not available")
	}

	// Get Everflow data
	everflowAdvertiser, err := s.everflowService.GetAdvertiserFromEverflowByMapping(ctx, advertiserID, []string{"billing", "settings"})
	if err != nil {
		return fmt.Errorf("failed to get advertiser from Everflow: %w", err)
	}

	// Get local advertiser
	localAdvertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, advertiserID)
	if err != nil {
		return fmt.Errorf("failed to get local advertiser: %w", err)
	}

	// Map Everflow data to local advertiser
	s.mapEverflowToLocalAdvertiser(everflowAdvertiser, localAdvertiser)

	// Update local advertiser
	return s.advertiserRepo.UpdateAdvertiser(ctx, localAdvertiser)
}

// CompareAdvertiserWithEverflow compares local advertiser with Everflow data
func (s *advertiserService) CompareAdvertiserWithEverflow(ctx context.Context, advertiserID int64) ([]domain.AdvertiserDiscrepancy, error) {
	var discrepancies []domain.AdvertiserDiscrepancy

	if s.everflowService == nil {
		return discrepancies, fmt.Errorf("Everflow service not available")
	}

	// Get local advertiser
	localAdvertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, advertiserID)
	if err != nil {
		return discrepancies, fmt.Errorf("failed to get local advertiser: %w", err)
	}

	// Get Everflow advertiser
	everflowAdvertiser, err := s.everflowService.GetAdvertiserFromEverflowByMapping(ctx, advertiserID, []string{"billing", "settings"})
	if err != nil {
		return discrepancies, fmt.Errorf("failed to get Everflow advertiser: %w", err)
	}

	// Compare name
	if localAdvertiser.Name != everflowAdvertiser.Name {
		discrepancies = append(discrepancies, domain.AdvertiserDiscrepancy{
			Field:         "name",
			LocalValue:    localAdvertiser.Name,
			ProviderValue: everflowAdvertiser.Name,
			Severity:      "high",
		})
	}

	// Compare status
	localStatus := s.mapLocalStatusToEverflow(localAdvertiser.Status)
	if localStatus != everflowAdvertiser.AccountStatus {
		discrepancies = append(discrepancies, domain.AdvertiserDiscrepancy{
			Field:         "status",
			LocalValue:    localAdvertiser.Status,
			ProviderValue: everflowAdvertiser.AccountStatus,
			Severity:      "medium",
		})
	}

	// Compare default currency
	if localAdvertiser.DefaultCurrencyID != nil && *localAdvertiser.DefaultCurrencyID != everflowAdvertiser.DefaultCurrencyID {
		discrepancies = append(discrepancies, domain.AdvertiserDiscrepancy{
			Field:         "default_currency_id",
			LocalValue:    *localAdvertiser.DefaultCurrencyID,
			ProviderValue: everflowAdvertiser.DefaultCurrencyID,
			Severity:      "low",
		})
	}

	// Compare platform fields
	if localAdvertiser.PlatformName != nil && *localAdvertiser.PlatformName != everflowAdvertiser.PlatformName {
		discrepancies = append(discrepancies, domain.AdvertiserDiscrepancy{
			Field:         "platform_name",
			LocalValue:    *localAdvertiser.PlatformName,
			ProviderValue: everflowAdvertiser.PlatformName,
			Severity:      "low",
		})
	}

	if localAdvertiser.PlatformURL != nil && *localAdvertiser.PlatformURL != everflowAdvertiser.PlatformURL {
		discrepancies = append(discrepancies, domain.AdvertiserDiscrepancy{
			Field:         "platform_url",
			LocalValue:    *localAdvertiser.PlatformURL,
			ProviderValue: everflowAdvertiser.PlatformURL,
			Severity:      "low",
		})
	}

	// Compare attribution fields
	if localAdvertiser.AttributionMethod != nil && *localAdvertiser.AttributionMethod != everflowAdvertiser.AttributionMethod {
		discrepancies = append(discrepancies, domain.AdvertiserDiscrepancy{
			Field:         "attribution_method",
			LocalValue:    *localAdvertiser.AttributionMethod,
			ProviderValue: everflowAdvertiser.AttributionMethod,
			Severity:      "medium",
		})
	}

	return discrepancies, nil
}

// Helper methods for mapping between local and Everflow data

func (s *advertiserService) mapLocalStatusToEverflow(status string) string {
	switch status {
	case "active":
		return "active"
	case "inactive", "rejected":
		return "inactive"
	case "pending":
		return "pending"
	default:
		return "pending"
	}
}

func (s *advertiserService) mapEverflowStatusToLocal(status string) string {
	switch status {
	case "active":
		return "active"
	case "inactive":
		return "inactive"
	case "pending":
		return "pending"
	default:
		return "pending"
	}
}

func (s *advertiserService) mapAdvertiserToEverflowUpdateRequest(advertiser *domain.Advertiser) (*everflow.EverflowUpdateAdvertiserRequest, error) {
	req := &everflow.EverflowUpdateAdvertiserRequest{
		Name:          advertiser.Name,
		AccountStatus: s.mapLocalStatusToEverflow(advertiser.Status),
	}

	// Map optional fields
	if advertiser.InternalNotes != nil {
		req.InternalNotes = advertiser.InternalNotes
	}

	// DefaultCurrencyID is required as string, provide default if nil
	if advertiser.DefaultCurrencyID != nil {
		req.DefaultCurrencyID = *advertiser.DefaultCurrencyID
	} else {
		req.DefaultCurrencyID = "USD" // Default currency
	}

	if advertiser.PlatformName != nil {
		req.PlatformName = advertiser.PlatformName
	}

	if advertiser.PlatformURL != nil {
		req.PlatformURL = advertiser.PlatformURL
	}

	if advertiser.PlatformUsername != nil {
		req.PlatformUsername = advertiser.PlatformUsername
	}

	if advertiser.AccountingContactEmail != nil {
		req.AccountingContactEmail = advertiser.AccountingContactEmail
	}

	if advertiser.OfferIDMacro != nil {
		req.OfferIDMacro = advertiser.OfferIDMacro
	}

	if advertiser.AffiliateIDMacro != nil {
		req.AffiliateIDMacro = advertiser.AffiliateIDMacro
	}

	if advertiser.AttributionMethod != nil {
		req.AttributionMethod = advertiser.AttributionMethod
	}

	if advertiser.EmailAttributionMethod != nil {
		req.EmailAttributionMethod = advertiser.EmailAttributionMethod
	}

	if advertiser.AttributionPriority != nil {
		req.AttributionPriority = advertiser.AttributionPriority
	}

	// ReportingTimezoneID is required as int, provide default if nil
	if advertiser.ReportingTimezoneID != nil {
		req.ReportingTimezoneID = *advertiser.ReportingTimezoneID
	} else {
		req.ReportingTimezoneID = 67 // Default timezone (UTC)
	}

	// Note: IsExposePublisherReportingData field doesn't exist in EverflowUpdateAdvertiserRequest
	// This might be handled through Settings field if needed

	return req, nil
}

func (s *advertiserService) mapEverflowToLocalAdvertiser(everflowAdv *everflow.Advertiser, localAdv *domain.Advertiser) {
	// Update basic fields
	localAdv.Name = everflowAdv.Name
	localAdv.Status = s.mapEverflowStatusToLocal(everflowAdv.AccountStatus)

	// Update optional fields
	if everflowAdv.InternalNotes != "" {
		localAdv.InternalNotes = &everflowAdv.InternalNotes
	}

	if everflowAdv.DefaultCurrencyID != "" {
		localAdv.DefaultCurrencyID = &everflowAdv.DefaultCurrencyID
	}

	if everflowAdv.PlatformName != "" {
		localAdv.PlatformName = &everflowAdv.PlatformName
	}

	if everflowAdv.PlatformURL != "" {
		localAdv.PlatformURL = &everflowAdv.PlatformURL
	}

	if everflowAdv.PlatformUsername != "" {
		localAdv.PlatformUsername = &everflowAdv.PlatformUsername
	}

	if everflowAdv.AccountingContactEmail != "" {
		localAdv.AccountingContactEmail = &everflowAdv.AccountingContactEmail
	}

	if everflowAdv.OfferIDMacro != "" {
		localAdv.OfferIDMacro = &everflowAdv.OfferIDMacro
	}

	if everflowAdv.AffiliateIDMacro != "" {
		localAdv.AffiliateIDMacro = &everflowAdv.AffiliateIDMacro
	}

	if everflowAdv.AttributionMethod != "" {
		localAdv.AttributionMethod = &everflowAdv.AttributionMethod
	}

	if everflowAdv.EmailAttributionMethod != "" {
		localAdv.EmailAttributionMethod = &everflowAdv.EmailAttributionMethod
	}

	if everflowAdv.AttributionPriority != "" {
		localAdv.AttributionPriority = &everflowAdv.AttributionPriority
	}

	if everflowAdv.ReportingTimezoneID != 0 {
		localAdv.ReportingTimezoneID = &everflowAdv.ReportingTimezoneID
	}

	if everflowAdv.IsExposePublisherReportingData != nil {
		localAdv.IsExposePublisherReporting = everflowAdv.IsExposePublisherReportingData
	}
}