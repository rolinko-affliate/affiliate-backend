package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/crypto"
	"github.com/affiliate-backend/internal/platform/provider"
	"github.com/affiliate-backend/internal/repository"
	"github.com/google/uuid"
)

// CampaignService defines the interface for campaign operations
type CampaignService interface {
	CreateCampaign(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error)
	GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error)
	UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error
	ListCampaignsByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]*domain.Campaign, error)
	ListCampaignsByAdvertiser(ctx context.Context, advertiserID int64, page, pageSize int) ([]*domain.Campaign, error)
	DeleteCampaign(ctx context.Context, id int64) error
	
	// Provider sync methods
	SyncCampaignToProvider(ctx context.Context, campaignID int64) error
	SyncCampaignFromProvider(ctx context.Context, campaignID int64) error
	
	// Provider offer methods
	CreateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) (*domain.CampaignProviderOffer, error)
	GetCampaignProviderOfferByID(ctx context.Context, id int64) (*domain.CampaignProviderOffer, error)
	UpdateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error
	ListCampaignProviderOffersByCampaign(ctx context.Context, campaignID int64) ([]*domain.CampaignProviderOffer, error)
	DeleteCampaignProviderOffer(ctx context.Context, id int64) error
}

// campaignService implements CampaignService
type campaignService struct {
	campaignRepo            repository.CampaignRepository
	campaignProviderMappingRepo repository.CampaignProviderMappingRepository
	advertiserRepo          repository.AdvertiserRepository
	orgRepo                 repository.OrganizationRepository
	cryptoService           crypto.Service
	integrationService      provider.IntegrationService
}

// NewCampaignService creates a new campaign service
func NewCampaignService(
	campaignRepo repository.CampaignRepository,
	campaignProviderMappingRepo repository.CampaignProviderMappingRepository,
	advertiserRepo repository.AdvertiserRepository,
	orgRepo repository.OrganizationRepository,
	cryptoService crypto.Service,
	integrationService provider.IntegrationService,
) CampaignService {
	return &campaignService{
		campaignRepo:                campaignRepo,
		campaignProviderMappingRepo: campaignProviderMappingRepo,
		advertiserRepo:              advertiserRepo,
		orgRepo:                     orgRepo,
		cryptoService:               cryptoService,
		integrationService:          integrationService,
	}
}

// CreateCampaign creates a new campaign
func (s *campaignService) CreateCampaign(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error) {
	// Validate organization exists
	_, err := s.orgRepo.GetOrganizationByID(ctx, campaign.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	// Validate advertiser exists
	advertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, campaign.AdvertiserID)
	if err != nil {
		return nil, fmt.Errorf("advertiser not found: %w", err)
	}

	// Validate advertiser belongs to the organization
	if advertiser.OrganizationID != campaign.OrganizationID {
		return nil, fmt.Errorf("advertiser does not belong to the specified organization")
	}

	// Validate required fields
	if campaign.Name == "" {
		return nil, fmt.Errorf("campaign name cannot be empty")
	}

	// Set default status if not provided
	if campaign.Status == "" {
		campaign.Status = "draft"
	}

	// Validate status
	validStatuses := map[string]bool{
		"draft":    true,
		"active":   true,
		"paused":   true,
		"archived": true,
	}
	if !validStatuses[campaign.Status] {
		return nil, fmt.Errorf("invalid status: %s", campaign.Status)
	}

	// Validate dates
	if campaign.StartDate != nil && campaign.EndDate != nil {
		if campaign.StartDate.After(*campaign.EndDate) {
			return nil, fmt.Errorf("start date cannot be after end date")
		}
	}

	// Set default values for offer fields if not provided
	s.setDefaultOfferValues(campaign)

	// Validate offer-specific fields
	if err := s.validateOfferFields(campaign); err != nil {
		return nil, fmt.Errorf("invalid offer configuration: %w", err)
	}

	// Step 1: Insert local record
	if err := s.campaignRepo.CreateCampaign(ctx, campaign); err != nil {
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	// Step 2: Call IntegrationService to create in provider
	providerCampaign, err := s.integrationService.CreateCampaign(ctx, *campaign)
	if err != nil {
		// Log error but don't fail the operation since local creation succeeded
		fmt.Printf("Warning: failed to create campaign in provider: %v\n", err)
		return campaign, nil
	}

	// Step 3: Create provider mapping with provider ID and payload
	var providerID *string
	if providerCampaign.NetworkAdvertiserID != nil {
		idStr := fmt.Sprintf("%d", *providerCampaign.NetworkAdvertiserID)
		providerID = &idStr
	}
	mapping := &domain.CampaignProviderMapping{
		CampaignID:         campaign.CampaignID,
		ProviderType:       "everflow",
		ProviderCampaignID: providerID,
		APICredentials:     nil, // Set by IntegrationService
		ProviderConfig:     nil, // Set by IntegrationService with full payload
	}

	if err := s.campaignProviderMappingRepo.CreateCampaignProviderMapping(ctx, mapping); err != nil {
		// Log error but don't fail the operation since campaign was created in provider
		fmt.Printf("Warning: failed to create provider mapping for campaign %d: %v\n", campaign.CampaignID, err)
	}

	return campaign, nil
}

// GetCampaignByID retrieves a campaign by ID
func (s *campaignService) GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error) {
	return s.campaignRepo.GetCampaignByID(ctx, id)
}

// UpdateCampaign updates a campaign
func (s *campaignService) UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	// Validate required fields
	if campaign.Name == "" {
		return fmt.Errorf("campaign name cannot be empty")
	}

	// Validate status
	validStatuses := map[string]bool{
		"draft":    true,
		"active":   true,
		"paused":   true,
		"archived": true,
	}
	if !validStatuses[campaign.Status] {
		return fmt.Errorf("invalid status: %s", campaign.Status)
	}

	// Validate dates
	if campaign.StartDate != nil && campaign.EndDate != nil {
		if campaign.StartDate.After(*campaign.EndDate) {
			return fmt.Errorf("start date cannot be after end date")
		}
	}

	// Get existing campaign to verify organization and advertiser
	existingCampaign, err := s.campaignRepo.GetCampaignByID(ctx, campaign.CampaignID)
	if err != nil {
		return fmt.Errorf("failed to get existing campaign: %w", err)
	}

	// Ensure organization and advertiser IDs are not changed
	campaign.OrganizationID = existingCampaign.OrganizationID
	campaign.AdvertiserID = existingCampaign.AdvertiserID

	// Step 1: Update local record first
	if err := s.campaignRepo.UpdateCampaign(ctx, campaign); err != nil {
		return fmt.Errorf("failed to update campaign: %w", err)
	}

	// Step 2: Check if provider mapping exists
	_, err = s.campaignProviderMappingRepo.GetCampaignProviderMapping(ctx, campaign.CampaignID, "everflow")
	if err != nil {
		// No provider mapping exists, skip provider sync
		return nil
	}

	// Step 3: Update in provider if mapping exists
	if err = s.integrationService.UpdateCampaign(ctx, *campaign); err != nil {
		// Log error but don't fail the operation since local update succeeded
		fmt.Printf("Warning: failed to update campaign in provider: %v\n", err)
	}

	return nil
}

// ListCampaignsByOrganization retrieves a list of campaigns for an organization with pagination
func (s *campaignService) ListCampaignsByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]*domain.Campaign, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.campaignRepo.ListCampaignsByOrganization(ctx, orgID, pageSize, offset)
}

// ListCampaignsByAdvertiser retrieves a list of campaigns for an advertiser with pagination
func (s *campaignService) ListCampaignsByAdvertiser(ctx context.Context, advertiserID int64, page, pageSize int) ([]*domain.Campaign, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.campaignRepo.ListCampaignsByAdvertiser(ctx, advertiserID, pageSize, offset)
}

// DeleteCampaign deletes a campaign
func (s *campaignService) DeleteCampaign(ctx context.Context, id int64) error {
	return s.campaignRepo.DeleteCampaign(ctx, id)
}

// CreateCampaignProviderOffer creates a new campaign provider offer
func (s *campaignService) CreateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) (*domain.CampaignProviderOffer, error) {
	// Validate campaign exists
	_, err := s.campaignRepo.GetCampaignByID(ctx, offer.CampaignID)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	// Validate provider type
	if offer.ProviderType != "everflow" {
		return nil, fmt.Errorf("invalid provider type: %s", offer.ProviderType)
	}

	// Validate provider offer config JSON if provided
	if offer.ProviderOfferConfig != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*offer.ProviderOfferConfig), &jsonData); err != nil {
			return nil, fmt.Errorf("invalid provider offer config JSON: %w", err)
		}
	}

	// Set default values
	if offer.LastSyncedAt == nil {
		now := time.Now()
		offer.LastSyncedAt = &now
	}

	if err := s.campaignRepo.CreateCampaignProviderOffer(ctx, offer); err != nil {
		return nil, fmt.Errorf("failed to create campaign provider offer: %w", err)
	}

	return offer, nil
}

// GetCampaignProviderOfferByID retrieves a campaign provider offer by ID
func (s *campaignService) GetCampaignProviderOfferByID(ctx context.Context, id int64) (*domain.CampaignProviderOffer, error) {
	return s.campaignRepo.GetCampaignProviderOfferByID(ctx, id)
}

// UpdateCampaignProviderOffer updates a campaign provider offer
func (s *campaignService) UpdateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error {
	// Validate provider offer config JSON if provided
	if offer.ProviderOfferConfig != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*offer.ProviderOfferConfig), &jsonData); err != nil {
			return fmt.Errorf("invalid provider offer config JSON: %w", err)
		}
	}

	// Get existing offer to verify campaign ID
	existingOffer, err := s.campaignRepo.GetCampaignProviderOfferByID(ctx, offer.CampaignProviderOfferID)
	if err != nil {
		return fmt.Errorf("failed to get existing campaign provider offer: %w", err)
	}

	// Ensure campaign ID is not changed
	offer.CampaignID = existingOffer.CampaignID
	offer.ProviderType = existingOffer.ProviderType

	// Update last synced time
	now := time.Now()
	offer.LastSyncedAt = &now

	return s.campaignRepo.UpdateCampaignProviderOffer(ctx, offer)
}

// ListCampaignProviderOffersByCampaign retrieves a list of campaign provider offers for a campaign
func (s *campaignService) ListCampaignProviderOffersByCampaign(ctx context.Context, campaignID int64) ([]*domain.CampaignProviderOffer, error) {
	return s.campaignRepo.ListCampaignProviderOffersByCampaign(ctx, campaignID)
}

// DeleteCampaignProviderOffer deletes a campaign provider offer
func (s *campaignService) DeleteCampaignProviderOffer(ctx context.Context, id int64) error {
	return s.campaignRepo.DeleteCampaignProviderOffer(ctx, id)
}

// setDefaultOfferValues sets default values for offer fields if not provided
func (s *campaignService) setDefaultOfferValues(campaign *domain.Campaign) {
	if campaign.Visibility == nil {
		visibility := "public"
		campaign.Visibility = &visibility
	}
	
	if campaign.CurrencyID == nil {
		currencyID := "USD"
		campaign.CurrencyID = &currencyID
	}
	
	if campaign.ConversionMethod == nil {
		conversionMethod := "server_postback"
		campaign.ConversionMethod = &conversionMethod
	}
	
	if campaign.SessionDefinition == nil {
		sessionDefinition := "cookie"
		campaign.SessionDefinition = &sessionDefinition
	}
	
	if campaign.SessionDuration == nil {
		sessionDuration := int32(24)
		campaign.SessionDuration = &sessionDuration
	}
	
	if campaign.PayoutType == nil {
		payoutType := "cpa"
		campaign.PayoutType = &payoutType
	}
	
	if campaign.PayoutAmount == nil {
		payoutAmount := 1.00
		campaign.PayoutAmount = &payoutAmount
	}
	
	if campaign.RevenueType == nil {
		revenueType := "rpa"
		campaign.RevenueType = &revenueType
	}
	
	if campaign.RevenueAmount == nil {
		revenueAmount := 2.00
		campaign.RevenueAmount = &revenueAmount
	}
	
	if campaign.IsForceTermsAndConditions == nil {
		isForce := false
		campaign.IsForceTermsAndConditions = &isForce
	}
	
	if campaign.IsCapsEnabled == nil {
		isCapsEnabled := false
		campaign.IsCapsEnabled = &isCapsEnabled
	}
}

// validateOfferFields validates offer-specific fields
func (s *campaignService) validateOfferFields(campaign *domain.Campaign) error {
	// Validate visibility
	if campaign.Visibility != nil {
		validVisibilities := map[string]bool{
			"public":           true,
			"require_approval": true,
			"private":          true,
		}
		if !validVisibilities[*campaign.Visibility] {
			return fmt.Errorf("invalid visibility: %s", *campaign.Visibility)
		}
	}
	
	// Validate conversion method
	if campaign.ConversionMethod != nil {
		validMethods := map[string]bool{
			"server_postback": true,
			"pixel":           true,
			"postback_url":    true,
			"api":             true,
		}
		if !validMethods[*campaign.ConversionMethod] {
			return fmt.Errorf("invalid conversion method: %s", *campaign.ConversionMethod)
		}
	}
	
	// Validate session definition
	if campaign.SessionDefinition != nil {
		validDefinitions := map[string]bool{
			"cookie":      true,
			"ip":          true,
			"fingerprint": true,
		}
		if !validDefinitions[*campaign.SessionDefinition] {
			return fmt.Errorf("invalid session definition: %s", *campaign.SessionDefinition)
		}
	}
	
	// Validate payout type
	if campaign.PayoutType != nil {
		validPayoutTypes := map[string]bool{
			"cpa":     true,
			"cpc":     true,
			"cpm":     true,
			"cps":     true,
			"cpa_cps": true,
			"prv":     true,
		}
		if !validPayoutTypes[*campaign.PayoutType] {
			return fmt.Errorf("invalid payout type: %s", *campaign.PayoutType)
		}
	}
	
	// Validate revenue type
	if campaign.RevenueType != nil {
		validRevenueTypes := map[string]bool{
			"rpa":     true,
			"rpc":     true,
			"rpm":     true,
			"rps":     true,
			"rpa_rps": true,
			"prv":     true,
		}
		if !validRevenueTypes[*campaign.RevenueType] {
			return fmt.Errorf("invalid revenue type: %s", *campaign.RevenueType)
		}
	}
	
	// Validate amounts are positive
	if campaign.PayoutAmount != nil && *campaign.PayoutAmount < 0 {
		return fmt.Errorf("payout amount must be non-negative")
	}
	
	if campaign.RevenueAmount != nil && *campaign.RevenueAmount < 0 {
		return fmt.Errorf("revenue amount must be non-negative")
	}
	
	// Validate session duration is positive
	if campaign.SessionDuration != nil && *campaign.SessionDuration <= 0 {
		return fmt.Errorf("session duration must be positive")
	}
	
	// Validate caps are non-negative
	caps := []*int{
		campaign.DailyConversionCap, campaign.WeeklyConversionCap,
		campaign.MonthlyConversionCap, campaign.GlobalConversionCap,
		campaign.DailyClickCap, campaign.WeeklyClickCap,
		campaign.MonthlyClickCap, campaign.GlobalClickCap,
	}
	
	for _, cap := range caps {
		if cap != nil && *cap < 0 {
			return fmt.Errorf("caps must be non-negative")
		}
	}
	
	return nil
}

// SyncCampaignToProvider syncs a campaign to the provider
func (s *campaignService) SyncCampaignToProvider(ctx context.Context, campaignID int64) error {
	// Get local campaign
	campaign, err := s.campaignRepo.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return fmt.Errorf("failed to get campaign: %w", err)
	}

	// Check if provider mapping exists
	_, err = s.campaignProviderMappingRepo.GetCampaignProviderMapping(ctx, campaignID, "everflow")
	if err != nil {
		// No mapping exists, create in provider
		return s.createCampaignInProvider(ctx, campaign)
	}

	// Mapping exists, update in provider
	if err = s.integrationService.UpdateCampaign(ctx, *campaign); err != nil {
		return fmt.Errorf("failed to sync campaign to provider: %w", err)
	}

	return nil
}

// SyncCampaignFromProvider syncs a campaign from the provider
func (s *campaignService) SyncCampaignFromProvider(ctx context.Context, campaignID int64) error {
	// Get provider mapping
	_, err := s.campaignProviderMappingRepo.GetCampaignProviderMapping(ctx, campaignID, "everflow")
	if err != nil {
		return fmt.Errorf("no provider mapping found for campaign %d: %w", campaignID, err)
	}

	// Convert campaign ID to UUID for IntegrationService
	campaignUUID := s.int64ToUUID(campaignID)
	
	// Get campaign from provider
	providerCampaign, err := s.integrationService.GetCampaign(ctx, campaignUUID)
	if err != nil {
		return fmt.Errorf("failed to get campaign from provider: %w", err)
	}

	// Update local campaign with provider data
	localCampaign, err := s.campaignRepo.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return fmt.Errorf("failed to get local campaign: %w", err)
	}

	// Merge provider data into local campaign
	s.mergeProviderDataIntoCampaign(localCampaign, &providerCampaign)

	// Update local record
	return s.campaignRepo.UpdateCampaign(ctx, localCampaign)
}

// Helper methods

// createCampaignInProvider creates a campaign in the provider when no mapping exists
func (s *campaignService) createCampaignInProvider(ctx context.Context, campaign *domain.Campaign) error {
	// Create in provider
	providerCampaign, err := s.integrationService.CreateCampaign(ctx, *campaign)
	if err != nil {
		return fmt.Errorf("failed to create campaign in provider: %w", err)
	}

	// Create provider mapping
	var providerID *string
	if providerCampaign.NetworkAdvertiserID != nil {
		idStr := fmt.Sprintf("%d", *providerCampaign.NetworkAdvertiserID)
		providerID = &idStr
	}
	mapping := &domain.CampaignProviderMapping{
		CampaignID:         campaign.CampaignID,
		ProviderType:       "everflow",
		ProviderCampaignID: providerID,
		APICredentials:     nil, // Set by IntegrationService
		ProviderConfig:     nil, // Set by IntegrationService with full payload
	}

	if err := s.campaignProviderMappingRepo.CreateCampaignProviderMapping(ctx, mapping); err != nil {
		fmt.Printf("Warning: failed to create provider mapping for campaign %d: %v\n", campaign.CampaignID, err)
	}

	return nil
}

// mergeProviderDataIntoCampaign merges provider data into local campaign
func (s *campaignService) mergeProviderDataIntoCampaign(local *domain.Campaign, provider *domain.Campaign) {
	// Merge relevant fields from provider into local
	if provider.NetworkAdvertiserID != nil {
		local.NetworkAdvertiserID = provider.NetworkAdvertiserID
	}
	// Add other fields as needed based on what the provider returns
}

// int64ToUUID converts int64 to UUID (copied from other services)
func (s *campaignService) int64ToUUID(id int64) uuid.UUID {
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