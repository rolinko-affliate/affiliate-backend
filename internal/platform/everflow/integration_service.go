package everflow

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/everflow/advertiser"
	"github.com/affiliate-backend/internal/platform/everflow/affiliate"
	"github.com/affiliate-backend/internal/platform/everflow/offer"
	"github.com/google/uuid"
)

// IntegrationService implements the provider-agnostic IntegrationService interface for Everflow
type IntegrationService struct {
	advertiserClient *advertiser.APIClient
	affiliateClient  *affiliate.APIClient
	offerClient      *offer.APIClient

	// Repository interfaces for provider mappings
	advertiserRepo AdvertiserRepository
	affiliateRepo  AffiliateRepository
	campaignRepo   CampaignRepository

	advertiserProviderMappingRepo AdvertiserProviderMappingRepository
	affiliateProviderMappingRepo  AffiliateProviderMappingRepository
	campaignProviderMappingRepo   CampaignProviderMappingRepository

	// Provider mappers
	affiliateProviderMapper  *AffiliateProviderMapper
	advertiserProviderMapper *AdvertiserProviderMapper
}

// Repository interfaces
type AdvertiserRepository interface {
	GetAdvertiserByID(ctx context.Context, id int64) (*domain.Advertiser, error)
}

type AffiliateRepository interface {
	GetAffiliateByID(ctx context.Context, id int64) (*domain.Affiliate, error)
}

type CampaignRepository interface {
	GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error)
}

type AdvertiserProviderMappingRepository interface {
	GetMappingByAdvertiserAndProvider(ctx context.Context, advertiserID int64, providerType string) (*domain.AdvertiserProviderMapping, error)
	CreateMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error
	UpdateMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error
}

type AffiliateProviderMappingRepository interface {
	GetAffiliateProviderMapping(ctx context.Context, affiliateID int64, providerType string) (*domain.AffiliateProviderMapping, error)
	CreateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error
	UpdateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error
}

type CampaignProviderMappingRepository interface {
	GetCampaignProviderMapping(ctx context.Context, campaignID int64, providerType string) (*domain.CampaignProviderMapping, error)
	CreateCampaignProviderMapping(ctx context.Context, mapping *domain.CampaignProviderMapping) error
	UpdateCampaignProviderMapping(ctx context.Context, mapping *domain.CampaignProviderMapping) error
}

// NewIntegrationService creates a new Everflow integration service
func NewIntegrationService(
	advertiserClient *advertiser.APIClient,
	affiliateClient *affiliate.APIClient,
	offerClient *offer.APIClient,
	advertiserRepo AdvertiserRepository,
	affiliateRepo AffiliateRepository,
	campaignRepo CampaignRepository,
	advertiserProviderMappingRepo AdvertiserProviderMappingRepository,
	affiliateProviderMappingRepo AffiliateProviderMappingRepository,
	campaignProviderMappingRepo CampaignProviderMappingRepository,
) *IntegrationService {
	return &IntegrationService{
		advertiserClient:              advertiserClient,
		affiliateClient:               affiliateClient,
		offerClient:                   offerClient,
		advertiserRepo:                advertiserRepo,
		affiliateRepo:                 affiliateRepo,
		campaignRepo:                  campaignRepo,
		advertiserProviderMappingRepo: advertiserProviderMappingRepo,
		affiliateProviderMappingRepo:  affiliateProviderMappingRepo,
		campaignProviderMappingRepo:   campaignProviderMappingRepo,
		affiliateProviderMapper:       NewAffiliateProviderMapper(),
		advertiserProviderMapper:      NewAdvertiserProviderMapper(),
	}
}

// UUID conversion helpers
func uuidToInt64(id uuid.UUID) (int64, error) {
	// Convert UUID to string and parse as int64
	// This is a simplified approach - in production you might want a more sophisticated mapping
	str := id.String()
	// Remove hyphens and take first 16 characters, then parse as hex
	cleaned := str[:8] + str[9:13] + str[14:18] + str[19:23]
	val, err := strconv.ParseInt(cleaned[:15], 16, 64) // Use 15 chars to avoid overflow
	if err != nil {
		return 0, fmt.Errorf("failed to convert UUID to int64: %w", err)
	}
	return val, nil
}

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

// CreateAdvertiser creates an advertiser in Everflow
func (s *IntegrationService) CreateAdvertiser(ctx context.Context, adv domain.Advertiser) (domain.Advertiser, error) {
	// Check if provider mapping already exists and is successful
	existingMapping, err := s.advertiserProviderMappingRepo.GetMappingByAdvertiserAndProvider(ctx, adv.AdvertiserID, "everflow")
	if err == nil && existingMapping != nil && existingMapping.SyncStatus != nil && *existingMapping.SyncStatus == "synced" {
		return adv, fmt.Errorf("advertiser already has successful Everflow provider mapping")
	}

	// Map domain advertiser to Everflow request (without existing mapping)
	everflowReq, err := s.advertiserProviderMapper.MapAdvertiserToEverflowRequest(&adv, nil)
	if err != nil {
		return adv, fmt.Errorf("failed to map advertiser to Everflow request: %w", err)
	}

	// Serialize the outbound request for payload storage
	requestPayload, err := json.Marshal(everflowReq)
	if err != nil {
		return adv, fmt.Errorf("failed to serialize request payload: %w", err)
	}

	// Log the request payload for debugging
	fmt.Printf("DEBUG: Sending Everflow advertiser request: %s\n", string(requestPayload))

	// Create advertiser in Everflow
	resp, httpResp, err := s.advertiserClient.DefaultAPI.CreateAdvertiser(ctx).CreateAdvertiserRequest(*everflowReq).Execute()
	if err != nil {
		// Try to read the response body for more detailed error information
		var errorBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				errorBody = string(bodyBytes)
			}
		}
		return adv, fmt.Errorf("failed to create advertiser in Everflow: %w (response: %s)", err, errorBody)
	}
	defer httpResp.Body.Close()

	// Create or update provider mapping
	now := time.Now()
	syncStatus := "synced"

	var mapping *domain.AdvertiserProviderMapping
	if existingMapping != nil {
		// Update existing failed mapping
		mapping = existingMapping
		mapping.SyncStatus = &syncStatus
		mapping.LastSyncAt = &now
		mapping.UpdatedAt = now
		mapping.SyncError = nil
	} else {
		// Create new mapping
		mapping = &domain.AdvertiserProviderMapping{
			AdvertiserID: adv.AdvertiserID,
			ProviderType: "everflow",
			SyncStatus:   &syncStatus,
			LastSyncAt:   &now,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
	}

	// Update mapping with Everflow response data
	if err := s.advertiserProviderMapper.MapEverflowResponseToProviderMapping(resp, mapping); err != nil {
		return adv, fmt.Errorf("failed to map Everflow response to provider mapping: %w", err)
	}

	// Store request/response payload in provider config
	payload := map[string]interface{}{
		"request":             json.RawMessage(requestPayload),
		"response":            resp,
		"last_operation":      "create",
		"last_operation_time": now,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return adv, fmt.Errorf("failed to marshal payload: %w", err)
	}

	payloadStr := string(payloadJSON)
	mapping.ProviderConfig = &payloadStr

	// Create or update the provider mapping
	if existingMapping != nil {
		if err := s.advertiserProviderMappingRepo.UpdateMapping(ctx, mapping); err != nil {
			return adv, fmt.Errorf("failed to update advertiser provider mapping: %w", err)
		}
	} else {
		if err := s.advertiserProviderMappingRepo.CreateMapping(ctx, mapping); err != nil {
			return adv, fmt.Errorf("failed to create advertiser provider mapping: %w", err)
		}
	}

	// Update core advertiser with non-provider-specific data from Everflow
	s.advertiserProviderMapper.MapEverflowResponseToAdvertiser(resp, &adv)

	return adv, nil
}

// UpdateAdvertiser updates an advertiser in Everflow
func (s *IntegrationService) UpdateAdvertiser(ctx context.Context, adv domain.Advertiser) error {
	// Get provider mapping
	mapping, err := s.advertiserProviderMappingRepo.GetMappingByAdvertiserAndProvider(ctx, adv.AdvertiserID, "everflow")
	if err != nil {
		return fmt.Errorf("failed to get advertiser provider mapping: %w", err)
	}

	if mapping.ProviderAdvertiserID == nil {
		return fmt.Errorf("advertiser not found in Everflow")
	}

	providerID, err := strconv.ParseInt(*mapping.ProviderAdvertiserID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid provider advertiser ID: %w", err)
	}

	// Map domain advertiser to Everflow update request using existing mapping
	everflowReq, err := s.advertiserProviderMapper.MapAdvertiserToEverflowRequest(&adv, mapping)
	if err != nil {
		return fmt.Errorf("failed to map advertiser to Everflow update request: %w", err)
	}

	// Convert to update request (different constructor requirements)
	updateReq := advertiser.NewUpdateAdvertiserRequest(
		everflowReq.GetName(),
		everflowReq.GetAccountStatus(),
		everflowReq.GetNetworkEmployeeId(),
		everflowReq.GetDefaultCurrencyId(),
		everflowReq.GetReportingTimezoneId(),
		everflowReq.GetAttributionMethod(),
		everflowReq.GetEmailAttributionMethod(),
		everflowReq.GetAttributionPriority(),
	)

	if everflowReq.HasInternalNotes() {
		updateReq.SetInternalNotes(everflowReq.GetInternalNotes())
	}
	if everflowReq.HasAddressId() {
		updateReq.SetAddressId(everflowReq.GetAddressId())
	}
	if everflowReq.HasIsContactAddressEnabled() {
		updateReq.SetIsContactAddressEnabled(everflowReq.GetIsContactAddressEnabled())
	}
	if everflowReq.HasSalesManagerId() {
		updateReq.SetSalesManagerId(everflowReq.GetSalesManagerId())
	}
	if everflowReq.HasPlatformName() {
		updateReq.SetPlatformName(everflowReq.GetPlatformName())
	}
	if everflowReq.HasPlatformUrl() {
		updateReq.SetPlatformUrl(everflowReq.GetPlatformUrl())
	}
	if everflowReq.HasPlatformUsername() {
		updateReq.SetPlatformUsername(everflowReq.GetPlatformUsername())
	}
	if everflowReq.HasAccountingContactEmail() {
		updateReq.SetAccountingContactEmail(everflowReq.GetAccountingContactEmail())
	}
	if everflowReq.HasVerificationToken() {
		updateReq.SetVerificationToken(everflowReq.GetVerificationToken())
	}
	if everflowReq.HasOfferIdMacro() {
		updateReq.SetOfferIdMacro(everflowReq.GetOfferIdMacro())
	}
	if everflowReq.HasAffiliateIdMacro() {
		updateReq.SetAffiliateIdMacro(everflowReq.GetAffiliateIdMacro())
	}
	if everflowReq.HasLabels() {
		updateReq.SetLabels(everflowReq.GetLabels())
	}
	if everflowReq.HasUsers() {
		updateReq.SetUsers(everflowReq.GetUsers())
	}
	if everflowReq.HasContactAddress() {
		updateReq.SetContactAddress(everflowReq.GetContactAddress())
	}
	if everflowReq.HasBilling() {
		updateReq.SetBilling(everflowReq.GetBilling())
	}
	if everflowReq.HasSettings() {
		updateReq.SetSettings(everflowReq.GetSettings())
	}

	// Update advertiser in Everflow
	resp, httpResp, err := s.advertiserClient.DefaultAPI.UpdateAdvertiser(ctx, int32(providerID)).UpdateAdvertiserRequest(*updateReq).Execute()
	if err != nil {
		return fmt.Errorf("failed to update advertiser in Everflow: %w", err)
	}
	defer httpResp.Body.Close()

	// Update provider mapping with response data
	if err := s.advertiserProviderMapper.MapEverflowResponseToProviderMapping(resp, mapping); err != nil {
		return fmt.Errorf("failed to update provider mapping with response: %w", err)
	}

	// Update sync metadata
	now := time.Now()
	syncStatus := "synced"
	mapping.SyncStatus = &syncStatus
	mapping.LastSyncAt = &now
	mapping.UpdatedAt = now

	// Update provider config with request/response payload
	requestPayload, _ := json.Marshal(updateReq)
	payload := map[string]interface{}{
		"request":             json.RawMessage(requestPayload),
		"response":            resp,
		"last_operation":      "update",
		"last_operation_time": now,
	}

	payloadJSON, _ := json.Marshal(payload)
	payloadStr := string(payloadJSON)
	mapping.ProviderConfig = &payloadStr

	return s.advertiserProviderMappingRepo.UpdateMapping(ctx, mapping)
}

// GetAdvertiser retrieves an advertiser from Everflow
func (s *IntegrationService) GetAdvertiser(ctx context.Context, id uuid.UUID) (domain.Advertiser, error) {
	// Convert UUID to int64
	advertiserID, err := uuidToInt64(id)
	if err != nil {
		return domain.Advertiser{}, fmt.Errorf("failed to convert UUID to int64: %w", err)
	}

	// Get local advertiser
	adv, err := s.advertiserRepo.GetAdvertiserByID(ctx, advertiserID)
	if err != nil {
		return domain.Advertiser{}, fmt.Errorf("failed to get local advertiser: %w", err)
	}

	// Get provider mapping
	mapping, err := s.advertiserProviderMappingRepo.GetMappingByAdvertiserAndProvider(ctx, advertiserID, "everflow")
	if err != nil {
		return *adv, fmt.Errorf("failed to get advertiser provider mapping: %w", err)
	}

	if mapping.ProviderAdvertiserID == nil {
		return *adv, fmt.Errorf("advertiser not found in Everflow")
	}

	providerID, err := strconv.ParseInt(*mapping.ProviderAdvertiserID, 10, 32)
	if err != nil {
		return *adv, fmt.Errorf("invalid provider advertiser ID: %w", err)
	}

	// Get advertiser from Everflow
	resp, httpResp, err := s.advertiserClient.DefaultAPI.GetAdvertiserById(ctx, int32(providerID)).Execute()
	if err != nil {
		return *adv, fmt.Errorf("failed to get advertiser from Everflow: %w", err)
	}
	defer httpResp.Body.Close()

	// Map Everflow response to domain model
	s.advertiserProviderMapper.MapEverflowResponseToAdvertiser(resp, adv)
	return *adv, nil
}

// CreateAffiliate creates an affiliate in Everflow
func (s *IntegrationService) CreateAffiliate(ctx context.Context, aff domain.Affiliate) (domain.Affiliate, error) {
	// Check if provider mapping already exists
	existingMapping, err := s.affiliateProviderMappingRepo.GetAffiliateProviderMapping(ctx, aff.AffiliateID, "everflow")
	if err == nil && existingMapping != nil {
		return aff, fmt.Errorf("affiliate already has Everflow provider mapping")
	}

	// Map domain affiliate to Everflow request (without existing mapping)
	everflowReq, err := s.affiliateProviderMapper.MapAffiliateToEverflowRequest(&aff, nil)
	if err != nil {
		return aff, fmt.Errorf("failed to map affiliate to Everflow request: %w", err)
	}

	// Serialize the outbound request for payload storage
	requestPayload, err := json.Marshal(everflowReq)
	if err != nil {
		return aff, fmt.Errorf("failed to serialize request payload: %w", err)
	}

	// Create affiliate in Everflow
	resp, httpResp, err := s.affiliateClient.AffiliatesAPI.CreateAffiliate(ctx).CreateAffiliateRequest(*everflowReq).Execute()
	if err != nil {
		return aff, fmt.Errorf("failed to create affiliate in Everflow: %w", err)
	}
	defer httpResp.Body.Close()

	// Create provider mapping
	now := time.Now()
	syncStatus := "synced"

	mapping := &domain.AffiliateProviderMapping{
		AffiliateID:  aff.AffiliateID,
		ProviderType: "everflow",
		SyncStatus:   &syncStatus,
		LastSyncAt:   &now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Update mapping with Everflow response data
	if err := s.affiliateProviderMapper.MapEverflowResponseToProviderMapping(resp, mapping); err != nil {
		return aff, fmt.Errorf("failed to map Everflow response to provider mapping: %w", err)
	}

	// Store request/response payload in provider config
	payload := map[string]interface{}{
		"request":             json.RawMessage(requestPayload),
		"response":            resp,
		"last_operation":      "create",
		"last_operation_time": now,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return aff, fmt.Errorf("failed to marshal payload: %w", err)
	}

	payloadStr := string(payloadJSON)
	mapping.ProviderConfig = &payloadStr

	// Create the provider mapping
	if err := s.affiliateProviderMappingRepo.CreateAffiliateProviderMapping(ctx, mapping); err != nil {
		return aff, fmt.Errorf("failed to create affiliate provider mapping: %w", err)
	}

	// Update core affiliate with non-provider-specific data from Everflow
	s.affiliateProviderMapper.MapEverflowResponseToAffiliate(resp, &aff)

	return aff, nil
}

// UpdateAffiliate updates an affiliate in Everflow
func (s *IntegrationService) UpdateAffiliate(ctx context.Context, aff domain.Affiliate) error {
	// Get provider mapping
	mapping, err := s.affiliateProviderMappingRepo.GetAffiliateProviderMapping(ctx, aff.AffiliateID, "everflow")
	if err != nil {
		return fmt.Errorf("failed to get affiliate provider mapping: %w", err)
	}

	if mapping.ProviderAffiliateID == nil {
		return fmt.Errorf("affiliate not found in Everflow")
	}

	providerID, err := strconv.ParseInt(*mapping.ProviderAffiliateID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid provider affiliate ID: %w", err)
	}

	// Map domain affiliate to Everflow update request using existing mapping
	everflowReq, err := s.affiliateProviderMapper.MapAffiliateToEverflowRequest(&aff, mapping)
	if err != nil {
		return fmt.Errorf("failed to map affiliate to Everflow update request: %w", err)
	}

	// Convert to update request (different constructor requirements)
	updateReq := affiliate.NewUpdateAffiliateRequest(
		everflowReq.GetName(),
		everflowReq.GetAccountStatus(),
		everflowReq.GetNetworkEmployeeId(),
	)

	if everflowReq.HasInternalNotes() {
		updateReq.SetInternalNotes(everflowReq.GetInternalNotes())
	}
	if everflowReq.HasDefaultCurrencyId() {
		updateReq.SetDefaultCurrencyId(everflowReq.GetDefaultCurrencyId())
	}
	if everflowReq.HasEnableMediaCostTrackingLinks() {
		updateReq.SetEnableMediaCostTrackingLinks(everflowReq.GetEnableMediaCostTrackingLinks())
	}
	if everflowReq.HasReferrerId() {
		updateReq.SetReferrerId(everflowReq.GetReferrerId())
	}
	if everflowReq.HasIsContactAddressEnabled() {
		updateReq.SetIsContactAddressEnabled(everflowReq.GetIsContactAddressEnabled())
	}
	if everflowReq.HasNetworkAffiliateTierId() {
		updateReq.SetNetworkAffiliateTierId(everflowReq.GetNetworkAffiliateTierId())
	}
	if everflowReq.HasContactAddress() {
		updateReq.SetContactAddress(everflowReq.GetContactAddress())
	}
	if everflowReq.HasBilling() {
		// Map to billing info for update request
		billingInfo := everflowReq.GetBilling()
		updateReq.SetBilling(billingInfo)
	}
	if everflowReq.HasLabels() {
		updateReq.SetLabels(everflowReq.GetLabels())
	}

	// Update affiliate in Everflow
	resp, httpResp, err := s.affiliateClient.AffiliatesAPI.UpdateAffiliate(ctx, int32(providerID)).UpdateAffiliateRequest(*updateReq).Execute()
	if err != nil {
		return fmt.Errorf("failed to update affiliate in Everflow: %w", err)
	}
	defer httpResp.Body.Close()

	// Update provider mapping with response data
	if err := s.affiliateProviderMapper.MapEverflowResponseToProviderMapping(resp, mapping); err != nil {
		return fmt.Errorf("failed to update provider mapping with response: %w", err)
	}

	// Update sync metadata
	now := time.Now()
	syncStatus := "synced"
	mapping.SyncStatus = &syncStatus
	mapping.LastSyncAt = &now
	mapping.UpdatedAt = now

	// Update provider config with request/response payload
	requestPayload, _ := json.Marshal(updateReq)
	payload := map[string]interface{}{
		"request":             json.RawMessage(requestPayload),
		"response":            resp,
		"last_operation":      "update",
		"last_operation_time": now,
	}

	payloadJSON, _ := json.Marshal(payload)
	payloadStr := string(payloadJSON)
	mapping.ProviderConfig = &payloadStr

	return s.affiliateProviderMappingRepo.UpdateAffiliateProviderMapping(ctx, mapping)
}

// GetAffiliate retrieves an affiliate from Everflow
func (s *IntegrationService) GetAffiliate(ctx context.Context, id uuid.UUID) (domain.Affiliate, error) {
	// Convert UUID to int64
	affiliateID, err := uuidToInt64(id)
	if err != nil {
		return domain.Affiliate{}, fmt.Errorf("failed to convert UUID to int64: %w", err)
	}

	// Get local affiliate
	aff, err := s.affiliateRepo.GetAffiliateByID(ctx, affiliateID)
	if err != nil {
		return domain.Affiliate{}, fmt.Errorf("failed to get local affiliate: %w", err)
	}

	// Get provider mapping
	mapping, err := s.affiliateProviderMappingRepo.GetAffiliateProviderMapping(ctx, affiliateID, "everflow")
	if err != nil {
		return *aff, fmt.Errorf("failed to get affiliate provider mapping: %w", err)
	}

	if mapping.ProviderAffiliateID == nil {
		return *aff, fmt.Errorf("affiliate not found in Everflow")
	}

	providerID, err := strconv.ParseInt(*mapping.ProviderAffiliateID, 10, 32)
	if err != nil {
		return *aff, fmt.Errorf("invalid provider affiliate ID: %w", err)
	}

	// Get affiliate from Everflow
	resp, httpResp, err := s.affiliateClient.AffiliatesAPI.GetAffiliateById(ctx, int32(providerID)).Execute()
	if err != nil {
		return *aff, fmt.Errorf("failed to get affiliate from Everflow: %w", err)
	}
	defer httpResp.Body.Close()

	// Map Everflow response to domain model
	updatedAffiliate := s.mapEverflowResponseToAffiliate(resp, aff)
	return updatedAffiliate, nil
}

// CreateCampaign creates a campaign (offer) in Everflow
func (s *IntegrationService) CreateCampaign(ctx context.Context, camp domain.Campaign) (domain.Campaign, error) {
	// TODO: Implement campaign creation when campaign functionality is needed
	return camp, fmt.Errorf("campaign creation not implemented")
}

// UpdateCampaign updates a campaign (offer) in Everflow
func (s *IntegrationService) UpdateCampaign(ctx context.Context, camp domain.Campaign) error {
	// TODO: Implement campaign update when campaign functionality is needed
	return fmt.Errorf("campaign update not implemented")
}

// GetCampaign retrieves a campaign from Everflow
func (s *IntegrationService) GetCampaign(ctx context.Context, id uuid.UUID) (domain.Campaign, error) {
	// TODO: Implement campaign retrieval when campaign functionality is needed
	return domain.Campaign{}, fmt.Errorf("campaign retrieval not implemented")
}

// #############################################################################
// ## Missing Mapper Methods - Placeholder Implementations
// #############################################################################

// mapAdvertiserToEverflowRequest maps domain advertiser to Everflow request
func (s *IntegrationService) mapAdvertiserToEverflowRequest(adv *domain.Advertiser) (interface{}, error) {
	// TODO: Implement advertiser mapping when advertiser functionality is needed
	return nil, fmt.Errorf("advertiser mapping not implemented")
}

// mapEverflowResponseToAdvertiser maps Everflow response to domain advertiser
func (s *IntegrationService) mapEverflowResponseToAdvertiser(resp interface{}, adv *domain.Advertiser) domain.Advertiser {
	// TODO: Implement response mapping when advertiser functionality is needed
	return *adv
}

// mapAdvertiserToEverflowUpdateRequest maps domain advertiser to Everflow update request
func (s *IntegrationService) mapAdvertiserToEverflowUpdateRequest(adv *domain.Advertiser) (interface{}, error) {
	// TODO: Implement advertiser update mapping when advertiser functionality is needed
	return nil, fmt.Errorf("advertiser update mapping not implemented")
}

// mapEverflowResponseToAffiliate maps Everflow response to domain affiliate
func (s *IntegrationService) mapEverflowResponseToAffiliate(resp interface{}, aff *domain.Affiliate) domain.Affiliate {
	// TODO: Implement response mapping when affiliate sync functionality is needed
	return *aff
}

// mapCampaignToEverflowRequest maps domain campaign to Everflow offer request
func (s *IntegrationService) mapCampaignToEverflowRequest(camp *domain.Campaign, networkAdvertiserID int32) (interface{}, error) {
	// TODO: Implement campaign mapping when campaign functionality is needed
	return nil, fmt.Errorf("campaign mapping not implemented")
}

// mapEverflowResponseToCampaign maps Everflow offer response to domain campaign
func (s *IntegrationService) mapEverflowResponseToCampaign(resp interface{}, camp *domain.Campaign) domain.Campaign {
	// TODO: Implement response mapping when campaign functionality is needed
	return *camp
}

// mapCampaignToEverflowUpdateRequest maps domain campaign to Everflow offer update request
func (s *IntegrationService) mapCampaignToEverflowUpdateRequest(camp *domain.Campaign) (interface{}, error) {
	// TODO: Implement campaign update mapping when campaign functionality is needed
	return nil, fmt.Errorf("campaign update mapping not implemented")
}

// GenerateTrackingLink generates a tracking link via Everflow API
func (s *IntegrationService) GenerateTrackingLink(ctx context.Context, req *domain.TrackingLinkGenerationRequest, campaignMapping *domain.CampaignProviderMapping, affiliateMapping *domain.AffiliateProviderMapping) (*domain.TrackingLinkGenerationResponse, error) {
	// Extract provider-specific IDs from mappings
	var networkOfferID, networkAffiliateID int32

	// Parse campaign provider data to get network_offer_id
	if campaignMapping != nil && campaignMapping.ProviderData != nil {
		var campaignProviderData domain.EverflowCampaignProviderData
		if err := campaignProviderData.FromJSON(*campaignMapping.ProviderData); err == nil {
			if campaignProviderData.NetworkCampaignID != nil {
				networkOfferID = *campaignProviderData.NetworkCampaignID
			}
		}
	}

	// Parse affiliate provider data to get network_affiliate_id
	if affiliateMapping != nil && affiliateMapping.ProviderData != nil {
		var affiliateProviderData domain.EverflowProviderData
		if err := json.Unmarshal([]byte(*affiliateMapping.ProviderData), &affiliateProviderData); err == nil {
			if affiliateProviderData.NetworkAffiliateID != nil {
				networkAffiliateID = *affiliateProviderData.NetworkAffiliateID
			}
		}
	}

	// If we don't have the required IDs, return an error
	if networkOfferID == 0 || networkAffiliateID == 0 {
		return nil, fmt.Errorf("missing required provider IDs: networkOfferID=%d, networkAffiliateID=%d", networkOfferID, networkAffiliateID)
	}

	// Create Everflow tracking link request
	everflowReq := map[string]interface{}{
		"network_offer_id":     networkOfferID,
		"network_affiliate_id": networkAffiliateID,
	}

	// Add optional parameters
	if req.NetworkTrackingDomainID != nil {
		everflowReq["network_tracking_domain_id"] = *req.NetworkTrackingDomainID
	}
	if req.NetworkOfferURLID != nil {
		everflowReq["network_offer_url_id"] = *req.NetworkOfferURLID
	}
	if req.CreativeID != nil {
		everflowReq["creative_id"] = *req.CreativeID
	}
	if req.NetworkTrafficSourceID != nil {
		everflowReq["network_traffic_source_id"] = *req.NetworkTrafficSourceID
	}
	if req.SourceID != nil {
		everflowReq["source_id"] = *req.SourceID
	}
	if req.Sub1 != nil {
		everflowReq["sub1"] = *req.Sub1
	}
	if req.Sub2 != nil {
		everflowReq["sub2"] = *req.Sub2
	}
	if req.Sub3 != nil {
		everflowReq["sub3"] = *req.Sub3
	}
	if req.Sub4 != nil {
		everflowReq["sub4"] = *req.Sub4
	}
	if req.Sub5 != nil {
		everflowReq["sub5"] = *req.Sub5
	}
	if req.IsEncryptParameters != nil {
		everflowReq["is_encrypt_parameters"] = *req.IsEncryptParameters
	}
	if req.IsRedirectLink != nil {
		everflowReq["is_redirect_link"] = *req.IsRedirectLink
	}

	// For now, simulate the Everflow API call since we don't have the tracking API client generated
	// In a real implementation, this would make an HTTP POST to /v1/networks/tracking/offers/clicks

	// Simulate the response
	baseURL := "http://tracking-domain.everflow.test"
	trackingPath := fmt.Sprintf("/%s/%s/", generateTrackingCode(), generateTrackingCode())

	// Build query parameters
	params := []string{}
	if req.SourceID != nil {
		params = append(params, fmt.Sprintf("source_id=%s", *req.SourceID))
	}
	if req.Sub1 != nil {
		params = append(params, fmt.Sprintf("sub1=%s", *req.Sub1))
	}
	if req.Sub2 != nil {
		params = append(params, fmt.Sprintf("sub2=%s", *req.Sub2))
	}
	if req.Sub3 != nil {
		params = append(params, fmt.Sprintf("sub3=%s", *req.Sub3))
	}
	if req.Sub4 != nil {
		params = append(params, fmt.Sprintf("sub4=%s", *req.Sub4))
	}
	if req.Sub5 != nil {
		params = append(params, fmt.Sprintf("sub5=%s", *req.Sub5))
	}

	generatedURL := baseURL + trackingPath
	if len(params) > 0 {
		generatedURL += "?" + joinParams(params)
	}

	// Create provider data
	providerData := &domain.EverflowTrackingLinkProviderData{
		NetworkOfferID:           &networkOfferID,
		NetworkAffiliateID:       &networkAffiliateID,
		NetworkTrackingDomainID:  req.NetworkTrackingDomainID,
		NetworkOfferURLID:        req.NetworkOfferURLID,
		CreativeID:               req.CreativeID,
		NetworkTrafficSourceID:   req.NetworkTrafficSourceID,
		GeneratedURL:             &generatedURL,
		CanAffiliateRunAllOffers: boolPtr(true),
	}

	providerDataJSON, err := providerData.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize provider data: %w", err)
	}

	return &domain.TrackingLinkGenerationResponse{
		GeneratedURL: generatedURL,
		ProviderData: &providerDataJSON,
	}, nil
}

// GenerateTrackingLinkQR generates a QR code for a tracking link via Everflow API
func (s *IntegrationService) GenerateTrackingLinkQR(ctx context.Context, req *domain.TrackingLinkGenerationRequest, campaignMapping *domain.CampaignProviderMapping, affiliateMapping *domain.AffiliateProviderMapping) ([]byte, error) {
	// For now, simulate the Everflow QR API call since we don't have the tracking API client generated
	// In a real implementation, this would make an HTTP POST to /v1/networks/tracking/offers/clicks/qr

	// Return a mock QR code (in real implementation, this would be a PNG image from Everflow)
	return []byte("everflow-qr-code-png-data"), nil
}

// Helper functions
func generateTrackingCode() string {
	// Generate a random tracking code (simplified)
	codes := []string{"ABC123", "DEF456", "GHI789", "JKL012", "MNO345"}
	return codes[time.Now().Nanosecond()%len(codes)]
}

func joinParams(params []string) string {
	result := ""
	for i, param := range params {
		if i > 0 {
			result += "&"
		}
		result += param
	}
	return result
}

func boolPtr(b bool) *bool {
	return &b
}
