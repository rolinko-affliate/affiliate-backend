package everflow

import (
	"context"
	"encoding/json"
	"fmt"
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
	affiliateProviderMapper *AffiliateProviderMapper
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
	// TODO: Implement advertiser creation when advertiser functionality is needed
	return adv, fmt.Errorf("advertiser creation not implemented")
}

// UpdateAdvertiser updates an advertiser in Everflow
func (s *IntegrationService) UpdateAdvertiser(ctx context.Context, adv domain.Advertiser) error {
	// TODO: Implement advertiser update when advertiser functionality is needed
	return fmt.Errorf("advertiser update not implemented")
}

// GetAdvertiser retrieves an advertiser from Everflow
func (s *IntegrationService) GetAdvertiser(ctx context.Context, id uuid.UUID) (domain.Advertiser, error) {
	// TODO: Implement advertiser retrieval when advertiser functionality is needed
	return domain.Advertiser{}, fmt.Errorf("advertiser retrieval not implemented")
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
		"request":               json.RawMessage(requestPayload),
		"response":              resp,
		"last_operation":        "create",
		"last_operation_time":   now,
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
		"request":               json.RawMessage(requestPayload),
		"response":              resp,
		"last_operation":        "update",
		"last_operation_time":   now,
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