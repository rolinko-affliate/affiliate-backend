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
	"github.com/affiliate-backend/internal/platform/everflow/tracking"
	"github.com/affiliate-backend/internal/platform/logger"
	"github.com/affiliate-backend/internal/platform/provider"
	"github.com/google/uuid"
)

// IntegrationService implements the provider-agnostic IntegrationService interface for Everflow
type IntegrationService struct {
	advertiserClient *advertiser.APIClient
	affiliateClient  *affiliate.APIClient
	offerClient      *offer.APIClient
	trackingClient   *tracking.APIClient

	// Repository interfaces for provider mappings
	advertiserRepo AdvertiserRepository
	affiliateRepo  AffiliateRepository
	campaignRepo   CampaignRepository

	advertiserProviderMappingRepo AdvertiserProviderMappingRepository
	affiliateProviderMappingRepo  AffiliateProviderMappingRepository
	campaignProviderMappingRepo   CampaignProviderMappingRepository

	// Provider mappers
	affiliateProviderMapper  *AffiliateProviderMapper
	advertiserProviderMapper AdvertiserMapper
}

// Mapper interfaces
type AdvertiserMapper interface {
	MapAdvertiserToEverflowRequest(adv *domain.Advertiser, mapping *domain.AdvertiserProviderMapping) (*advertiser.CreateAdvertiserRequest, error)
	MapAdvertiserToEverflowRequestWithContext(adv *domain.Advertiser, mapping *domain.AdvertiserProviderMapping, ctx *provider.AdvertiserMappingContext) (*advertiser.CreateAdvertiserRequest, error)
	MapEverflowResponseToAdvertiser(resp *advertiser.Advertiser, adv *domain.Advertiser)
	MapEverflowResponseToProviderMapping(resp *advertiser.Advertiser, mapping *domain.AdvertiserProviderMapping) error
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
	trackingClient *tracking.APIClient,
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
		trackingClient:                trackingClient,
		advertiserRepo:                advertiserRepo,
		affiliateRepo:                 affiliateRepo,
		campaignRepo:                  campaignRepo,
		advertiserProviderMappingRepo: advertiserProviderMappingRepo,
		affiliateProviderMappingRepo:  affiliateProviderMappingRepo,
		campaignProviderMappingRepo:   campaignProviderMappingRepo,
		affiliateProviderMapper:       NewAffiliateProviderMapper(),
		advertiserProviderMapper:      NewSimpleAdvertiserProviderMapper(),
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

// CreateAdvertiserWithContext creates an advertiser in Everflow with additional context
func (s *IntegrationService) CreateAdvertiserWithContext(ctx context.Context, adv domain.Advertiser, mappingCtx *provider.AdvertiserMappingContext) (domain.Advertiser, error) {
	logger.Info("Starting advertiser creation in Everflow with context",
		"advertiser_id", adv.AdvertiserID,
		"name", adv.Name,
		"provider", "everflow")
	
	// Check if provider mapping already exists and is successful
	existingMapping, err := s.advertiserProviderMappingRepo.GetMappingByAdvertiserAndProvider(ctx, adv.AdvertiserID, "everflow")
	if err == nil && existingMapping != nil && existingMapping.SyncStatus != nil && *existingMapping.SyncStatus == "synced" {
		logger.Warn("Advertiser already has successful Everflow provider mapping",
			"advertiser_id", adv.AdvertiserID,
			"mapping_id", existingMapping.MappingID)
		return adv, fmt.Errorf("advertiser already has successful Everflow provider mapping")
	}
	
	// Map domain advertiser to Everflow request with context
	everflowReq, err := s.advertiserProviderMapper.MapAdvertiserToEverflowRequestWithContext(&adv, nil, mappingCtx)
	if err != nil {
		logger.Error("Failed to map advertiser to Everflow request",
			"advertiser_id", adv.AdvertiserID,
			"error", err)
		return adv, fmt.Errorf("failed to map advertiser to Everflow request: %w", err)
	}

	// Serialize the outbound request for payload storage
	requestPayload, err := json.Marshal(everflowReq)
	if err != nil {
		logger.Error("Failed to serialize request payload",
			"advertiser_id", adv.AdvertiserID,
			"error", err)
		return adv, fmt.Errorf("failed to serialize request payload: %w", err)
	}

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
		logger.ErrorSanitized("Failed to create advertiser in Everflow",
			"advertiser_id", adv.AdvertiserID,
			"error", err,
			"response_body", errorBody)
		return adv, fmt.Errorf("failed to create advertiser in Everflow: %w (response: %s)", err, errorBody)
	}
	defer httpResp.Body.Close()
	logger.Info("Successfully created advertiser in Everflow",
		"advertiser_id", adv.AdvertiserID)

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

// CreateAdvertiser creates an advertiser in Everflow (backward compatibility)
func (s *IntegrationService) CreateAdvertiser(ctx context.Context, adv domain.Advertiser) (domain.Advertiser, error) {
	// Call the context-aware version with nil context for backward compatibility
	return s.CreateAdvertiserWithContext(ctx, adv, nil)
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
	logger.Info("Starting affiliate creation in Everflow",
		"affiliate_id", aff.AffiliateID,
		"name", aff.Name,
		"provider", "everflow")
	
	// Check if provider mapping already exists and is successful
	existingMapping, err := s.affiliateProviderMappingRepo.GetAffiliateProviderMapping(ctx, aff.AffiliateID, "everflow")
	if err == nil && existingMapping != nil && existingMapping.SyncStatus != nil && *existingMapping.SyncStatus == "synced" {
		logger.Warn("Affiliate already has successful Everflow mapping",
			"affiliate_id", aff.AffiliateID,
			"mapping_id", existingMapping.MappingID)
		return aff, fmt.Errorf("affiliate already has successful Everflow provider mapping")
	}
	
	// Map domain affiliate to Everflow request (without existing mapping)
	everflowReq, err := s.affiliateProviderMapper.MapAffiliateToEverflowRequest(&aff, nil)
	if err != nil {
		logger.Error("Failed to map affiliate to Everflow request",
			"affiliate_id", aff.AffiliateID,
			"error", err)
		return aff, fmt.Errorf("failed to map affiliate to Everflow request: %w", err)
	}

	// Serialize the outbound request for payload storage
	requestPayload, err := json.Marshal(everflowReq)
	if err != nil {
		logger.Error("Failed to serialize request payload",
			"affiliate_id", aff.AffiliateID,
			"error", err)
		return aff, fmt.Errorf("failed to serialize request payload: %w", err)
	}

	// Create affiliate in Everflow
	resp, httpResp, err := s.affiliateClient.AffiliatesAPI.CreateAffiliate(ctx).CreateAffiliateRequest(*everflowReq).Execute()
	if err != nil {
		// Try to read the response body for more detailed error information
		var errorBody string
		var statusCode int
		if httpResp != nil {
			statusCode = httpResp.StatusCode
			if httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					errorBody = string(bodyBytes)
				}
			}
		}
		logger.ErrorSanitized("API call failed",
			"affiliate_id", aff.AffiliateID,
			"status_code", statusCode,
			"error", err,
			"response_body", errorBody)
		return aff, fmt.Errorf("failed to create affiliate in Everflow: %w (response: %s)", err, errorBody)
	}
	defer httpResp.Body.Close()
	
	logger.Info("Successfully created affiliate in Everflow",
		"affiliate_id", aff.AffiliateID)

	// Create or update provider mapping
	now := time.Now()
	syncStatus := "synced"

	var mapping *domain.AffiliateProviderMapping
	if existingMapping != nil {
		// Update existing failed mapping
		mapping = existingMapping
		mapping.SyncStatus = &syncStatus
		mapping.LastSyncAt = &now
		mapping.UpdatedAt = now
		mapping.SyncError = nil
	} else {
		// Create new mapping
		mapping = &domain.AffiliateProviderMapping{
			AffiliateID:  aff.AffiliateID,
			ProviderType: "everflow",
			SyncStatus:   &syncStatus,
			LastSyncAt:   &now,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
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

	// Create or update the provider mapping
	if existingMapping != nil {
		if err := s.affiliateProviderMappingRepo.UpdateAffiliateProviderMapping(ctx, mapping); err != nil {
			return aff, fmt.Errorf("failed to update affiliate provider mapping: %w", err)
		}
	} else {
		if err := s.affiliateProviderMappingRepo.CreateAffiliateProviderMapping(ctx, mapping); err != nil {
			return aff, fmt.Errorf("failed to create affiliate provider mapping: %w", err)
		}
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
	logger.Info("Starting campaign creation in Everflow", "campaign_id", camp.CampaignID, "advertiser_id", camp.AdvertiserID, "name", camp.Name, "provider", "everflow")
	
	// Check if provider mapping already exists and is successful
	logger.Debug("Checking for existing provider mapping", "campaign_id", camp.CampaignID)
	existingMapping, err := s.campaignProviderMappingRepo.GetCampaignProviderMapping(ctx, camp.CampaignID, "everflow")
	if err == nil && existingMapping != nil && existingMapping.IsActiveOnProvider != nil && *existingMapping.IsActiveOnProvider {
		logger.Warn("Campaign already has successful Everflow provider mapping", "campaign_id", camp.CampaignID)
		return camp, fmt.Errorf("campaign already has successful Everflow provider mapping")
	}
	if err != nil {
		logger.Debug("No existing mapping found (expected)", "campaign_id", camp.CampaignID, "error", err)
	} else if existingMapping != nil {
		logger.Debug("Found existing mapping but not active, will update it", "campaign_id", camp.CampaignID)
	}

	// Get advertiser provider mapping to get network_advertiser_id
	logger.Debug("Getting advertiser provider mapping", "campaign_id", camp.CampaignID, "advertiser_id", camp.AdvertiserID)
	advertiserMapping, err := s.advertiserProviderMappingRepo.GetMappingByAdvertiserAndProvider(ctx, camp.AdvertiserID, "everflow")
	if err != nil {
		logger.Error("Failed to get advertiser provider mapping", "campaign_id", camp.CampaignID, "advertiser_id", camp.AdvertiserID, "error", err)
		return camp, fmt.Errorf("failed to get advertiser provider mapping: %w", err)
	}
	if advertiserMapping.ProviderAdvertiserID == nil {
		logger.Error("Advertiser provider mapping exists but ProviderAdvertiserID is nil", "campaign_id", camp.CampaignID)
		return camp, fmt.Errorf("advertiser not found in Everflow")
	}

	networkAdvertiserID, err := strconv.ParseInt(*advertiserMapping.ProviderAdvertiserID, 10, 32)
	if err != nil {
		logger.Error("Invalid provider advertiser ID format", "campaign_id", camp.CampaignID, "error", err)
		return camp, fmt.Errorf("invalid provider advertiser ID: %w", err)
	}

	// Map domain campaign to Everflow offer request
	everflowReq, err := s.mapCampaignToEverflowRequest(&camp, int32(networkAdvertiserID))
	if err != nil {
		logger.Error("Failed to map campaign to Everflow request", "campaign_id", camp.CampaignID, "error", err)
		return camp, fmt.Errorf("failed to map campaign to Everflow request: %w", err)
	}

	// Serialize the outbound request for payload storage
	requestPayload, err := json.Marshal(everflowReq)
	if err != nil {
		logger.Error("Failed to serialize request payload", "campaign_id", camp.CampaignID, "error", err)
		return camp, fmt.Errorf("failed to serialize request payload: %w", err)
	}
	// Create offer in Everflow
	resp, httpResp, err := s.offerClient.OffersAPI.CreateOffer(ctx).CreateOfferRequest(*everflowReq).Execute()
	if err != nil {
		// Try to read the response body for more detailed error information
		var errorBody string
		var statusCode int
		if httpResp != nil {
			statusCode = httpResp.StatusCode
			if httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					errorBody = string(bodyBytes)
				}
			}
		}
		logger.ErrorSanitized("API call failed", "campaign_id", camp.CampaignID, "status_code", statusCode, "error", err, "response_body", errorBody)
		return camp, fmt.Errorf("failed to create offer in Everflow: %w (response: %s)", err, errorBody)
	}
	defer httpResp.Body.Close()
	logger.Info("Successfully created offer in Everflow", "campaign_id", camp.CampaignID)

	// Create or update provider mapping
	now := time.Now()

	var mapping *domain.CampaignProviderMapping
	if existingMapping != nil {
		// Update existing failed mapping
		mapping = existingMapping
		isActive := true
		mapping.IsActiveOnProvider = &isActive
		mapping.LastSyncedAt = &now
		mapping.UpdatedAt = now
	} else {
		// Create new mapping
		isActive := true
		mapping = &domain.CampaignProviderMapping{
			CampaignID:         camp.CampaignID,
			ProviderType:       "everflow",
			IsActiveOnProvider: &isActive,
			LastSyncedAt:       &now,
			CreatedAt:          now,
			UpdatedAt:          now,
		}
	}

	// Update mapping with Everflow response data
	if err := s.mapEverflowResponseToCampaignMapping(resp, mapping); err != nil {
		logger.Error("Failed to map Everflow response to provider mapping", "campaign_id", camp.CampaignID, "error", err)
		return camp, fmt.Errorf("failed to map Everflow response to provider mapping: %w", err)
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
		logger.Error("Failed to marshal payload", "campaign_id", camp.CampaignID, "error", err)
		return camp, fmt.Errorf("failed to marshal payload: %w", err)
	}

	payloadStr := string(payloadJSON)
	mapping.ProviderData = &payloadStr

	// Create or update the provider mapping
	if existingMapping != nil {
		if err := s.campaignProviderMappingRepo.UpdateCampaignProviderMapping(ctx, mapping); err != nil {
			logger.Error("Failed to update campaign provider mapping", "campaign_id", camp.CampaignID, "error", err)
			return camp, fmt.Errorf("failed to update campaign provider mapping: %w", err)
		}
		logger.Info("Successfully updated provider mapping", "campaign_id", camp.CampaignID)
	} else {
		if err := s.campaignProviderMappingRepo.CreateCampaignProviderMapping(ctx, mapping); err != nil {
			logger.Error("Failed to create campaign provider mapping", "campaign_id", camp.CampaignID, "error", err)
			return camp, fmt.Errorf("failed to create campaign provider mapping: %w", err)
		}
		logger.Info("Successfully created provider mapping", "campaign_id", camp.CampaignID)
	}

	// Update core campaign with non-provider-specific data from Everflow
	s.mapEverflowResponseToCampaign(resp, &camp)
	logger.Info("Successfully completed campaign creation", "campaign_id", camp.CampaignID)

	return camp, nil
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
func (s *IntegrationService) mapCampaignToEverflowRequest(camp *domain.Campaign, networkAdvertiserID int32) (*offer.CreateOfferRequest, error) {
	// Determine payout and revenue types based on campaign configuration
	var payoutType, revenueType string
	var payoutRevenue []offer.PayoutRevenue

	// Check if this is a click-based campaign
	if camp.FixedClickAmount != nil && *camp.FixedClickAmount > 0 {
		// Click-based campaign: RPC revenue model, CPC payout model
		revenueType = "rpc"  // Revenue Per Click
		payoutType = "cpc"   // Cost Per Click
		
		payoutRevenue = []offer.PayoutRevenue{
			*offer.NewPayoutRevenue(payoutType, revenueType, true, false),
		}
		
		// Set revenue to 0 (we don't earn revenue from clicks)
		payoutRevenue[0].SetRevenueAmount(0.0)
		
		// Set payout amount from user configuration
		payoutRevenue[0].SetPayoutAmount(*camp.FixedClickAmount)
		
	} else {
		// Conversion-based campaign: RPS revenue model (100%), payout depends on configuration
		revenueType = "rps"  // Revenue Per Sale
		
		// Determine payout type based on what's configured
		hasFixedAmount := camp.FixedConversionAmount != nil && *camp.FixedConversionAmount > 0
		hasPercentageAmount := camp.PercentageConversionAmount != nil && *camp.PercentageConversionAmount > 0
		
		if hasFixedAmount && hasPercentageAmount {
			payoutType = "cpa_cps"  // Cost Per Action + Cost Per Sale (mixed)
		} else if hasPercentageAmount {
			payoutType = "cps"      // Cost Per Sale (percentage only)
		} else {
			payoutType = "cpa"      // Cost Per Action (fixed amount, default)
		}
		
		payoutRevenue = []offer.PayoutRevenue{
			*offer.NewPayoutRevenue(payoutType, revenueType, true, false),
		}
		
		// Set revenue to 100% (we earn 100% of the sale)
		payoutRevenue[0].SetRevenuePercentage(100)
		
		// Set payout amounts based on configuration
		if hasFixedAmount {
			payoutRevenue[0].SetPayoutAmount(*camp.FixedConversionAmount)
		}
		
		if hasPercentageAmount {
			payoutRevenue[0].SetPayoutPercentage(int32(*camp.PercentageConversionAmount))
		}
		
		// Set default payout amount if nothing is configured
		if !hasFixedAmount && !hasPercentageAmount {
			payoutRevenue[0].SetPayoutAmount(0.0) // Default to 0 if not configured
		}
	}

	// Determine destination URL
	destinationURL := "https://example.com"
	if camp.DestinationURL != nil && *camp.DestinationURL != "" {
		destinationURL = *camp.DestinationURL
	}

	// Determine offer status based on campaign status
	offerStatus := "active"
	if camp.Status == "paused" {
		offerStatus = "paused"
	} else if camp.Status == "draft" {
		offerStatus = "inactive"
	}

	// Create the offer request
	req := offer.NewCreateOfferRequest(
		networkAdvertiserID,
		camp.Name,
		destinationURL,
		offerStatus,
		payoutRevenue,
	)

	// Set optional fields
	if camp.ThumbnailURL != nil {
		req.SetThumbnailUrl(*camp.ThumbnailURL)
	}

	if camp.PreviewURL != nil {
		req.SetPreviewUrl(*camp.PreviewURL)
	}

	if camp.InternalNotes != nil {
		req.SetInternalNotes(*camp.InternalNotes)
	}

	if camp.CurrencyID != nil {
		req.SetCurrencyId(*camp.CurrencyID)
	}

	if camp.ConversionMethod != nil {
		req.SetConversionMethod(*camp.ConversionMethod)
	}

	if camp.SessionDefinition != nil {
		req.SetSessionDefinition(*camp.SessionDefinition)
	}

	if camp.SessionDuration != nil {
		req.SetSessionDuration(*camp.SessionDuration)
	}

	if camp.Visibility != nil {
		req.SetVisibility(*camp.Visibility)
	}

	if camp.TermsAndConditions != nil {
		req.SetTermsAndConditions(*camp.TermsAndConditions)
		req.SetIsUsingExplicitTermsAndConditions(true)
	}

	// Set end date if available
	if camp.EndDate != nil {
		dateLiveUntil := camp.EndDate.Format("2006-01-02")
		req.SetDateLiveUntil(dateLiveUntil)
	}

	// Set description as HTML description
	if camp.Description != nil {
		req.SetHtmlDescription(*camp.Description)
	}

	// Set caps if enabled
	if camp.IsCapsEnabled != nil && *camp.IsCapsEnabled {
		req.SetIsCapsEnabled(true)
		if camp.DailyConversionCap != nil {
			req.SetDailyConversionCap(int32(*camp.DailyConversionCap))
		}
		if camp.WeeklyConversionCap != nil {
			req.SetWeeklyConversionCap(int32(*camp.WeeklyConversionCap))
		}
		if camp.MonthlyConversionCap != nil {
			req.SetMonthlyConversionCap(int32(*camp.MonthlyConversionCap))
		}
		if camp.GlobalConversionCap != nil {
			req.SetGlobalConversionCap(int32(*camp.GlobalConversionCap))
		}
		if camp.DailyClickCap != nil {
			req.SetDailyClickCap(int32(*camp.DailyClickCap))
		}
		if camp.WeeklyClickCap != nil {
			req.SetWeeklyClickCap(int32(*camp.WeeklyClickCap))
		}
		if camp.MonthlyClickCap != nil {
			req.SetMonthlyClickCap(int32(*camp.MonthlyClickCap))
		}
		if camp.GlobalClickCap != nil {
			req.SetGlobalClickCap(int32(*camp.GlobalClickCap))
		}
	}

	// Set constant values as mentioned by user
	networkTrackingDomainID := int32(12977)
	req.SetNetworkTrackingDomainId(networkTrackingDomainID)
	req.SetIsUseSecureLink(true)

	// Set default category ID
	req.SetNetworkCategoryId(1)

	// Create empty ruleset with default timezone
	ruleset := offer.NewRuleset()
	ruleset.SetDayPartingTimezoneId(58) // Default timezone from example
	req.SetRuleset(*ruleset)

	// Set attribution methods
	req.SetEmailAttributionMethod("first_affiliate_attribution")
	req.SetAttributionMethod("last_touch")

	return req, nil
}

// mapEverflowResponseToCampaign maps Everflow offer response to domain campaign
func (s *IntegrationService) mapEverflowResponseToCampaign(resp *offer.OfferResponse, camp *domain.Campaign) domain.Campaign {
	// Update campaign with response data if available
	if resp != nil {
		// Update basic fields that might have been modified by Everflow
		if resp.HasName() {
			camp.Name = resp.GetName()
		}
		
		if resp.HasDestinationUrl() {
			destinationURL := resp.GetDestinationUrl()
			camp.DestinationURL = &destinationURL
		}
		
		if resp.HasThumbnailUrl() {
			thumbnailURL := resp.GetThumbnailUrl()
			camp.ThumbnailURL = &thumbnailURL
		}
		
		if resp.HasPreviewUrl() {
			previewURL := resp.GetPreviewUrl()
			camp.PreviewURL = &previewURL
		}
		
		if resp.HasInternalNotes() {
			internalNotes := resp.GetInternalNotes()
			camp.InternalNotes = &internalNotes
		}
		
		if resp.HasCurrencyId() {
			currencyID := resp.GetCurrencyId()
			camp.CurrencyID = &currencyID
		}
		
		if resp.HasConversionMethod() {
			conversionMethod := resp.GetConversionMethod()
			camp.ConversionMethod = &conversionMethod
		}
		
		if resp.HasSessionDefinition() {
			sessionDefinition := resp.GetSessionDefinition()
			camp.SessionDefinition = &sessionDefinition
		}
		
		if resp.HasSessionDuration() {
			sessionDuration := resp.GetSessionDuration()
			camp.SessionDuration = &sessionDuration
		}
		
		if resp.HasVisibility() {
			visibility := resp.GetVisibility()
			camp.Visibility = &visibility
		}
		
		if resp.HasTermsAndConditions() {
			termsAndConditions := resp.GetTermsAndConditions()
			camp.TermsAndConditions = &termsAndConditions
		}
		
		// Map offer status back to campaign status
		if resp.HasOfferStatus() {
			offerStatus := resp.GetOfferStatus()
			switch offerStatus {
			case "active":
				camp.Status = "active"
			case "paused":
				camp.Status = "paused"
			case "inactive":
				camp.Status = "draft"
			default:
				camp.Status = "draft"
			}
		}
	}
	
	return *camp
}

// mapEverflowResponseToCampaignMapping maps Everflow offer response to campaign provider mapping
func (s *IntegrationService) mapEverflowResponseToCampaignMapping(resp *offer.OfferResponse, mapping *domain.CampaignProviderMapping) error {
	if resp == nil {
		return fmt.Errorf("invalid response data")
	}

	// Set provider offer ID (Everflow calls campaigns "offers")
	if resp.HasNetworkOfferId() {
		providerOfferID := strconv.FormatInt(int64(resp.GetNetworkOfferId()), 10)
		mapping.ProviderOfferID = &providerOfferID
	}

	// Create provider data structure
	providerData := &domain.EverflowCampaignProviderData{}
	
	if resp.HasNetworkOfferId() {
		networkCampaignID := resp.GetNetworkOfferId()
		providerData.NetworkCampaignID = &networkCampaignID
	}
	
	if resp.HasNetworkAdvertiserId() {
		networkAdvertiserID := resp.GetNetworkAdvertiserId()
		providerData.NetworkAdvertiserID = &networkAdvertiserID
	}
	
	if resp.HasCapsTimezoneId() {
		capsTimezoneID := resp.GetCapsTimezoneId()
		providerData.CapsTimezoneID = &capsTimezoneID
	}
	
	if resp.HasProjectId() {
		projectID := resp.GetProjectId()
		providerData.ProjectID = &projectID
	}
	
	if resp.HasHtmlDescription() {
		htmlDescription := resp.GetHtmlDescription()
		providerData.HTMLDescription = &htmlDescription
	}
	
	if resp.HasIsUsingExplicitTermsAndConditions() {
		isUsingExplicitTermsAndConditions := resp.GetIsUsingExplicitTermsAndConditions()
		providerData.IsUsingExplicitTermsAndConditions = &isUsingExplicitTermsAndConditions
	}
	
	if resp.HasIsForceTermsAndConditions() {
		isForceTermsAndConditions := resp.GetIsForceTermsAndConditions()
		providerData.IsForceTermsAndConditions = &isForceTermsAndConditions
	}
	
	if resp.HasIsViewThroughEnabled() {
		isViewThroughEnabled := resp.GetIsViewThroughEnabled()
		providerData.IsViewThroughEnabled = &isViewThroughEnabled
	}
	
	if resp.HasServerSideUrl() {
		serverSideURL := resp.GetServerSideUrl()
		providerData.ServerSideURL = &serverSideURL
	}
	
	if resp.HasViewThroughDestinationUrl() {
		viewThroughDestinationURL := resp.GetViewThroughDestinationUrl()
		providerData.ViewThroughDestinationURL = &viewThroughDestinationURL
	}
	
	if resp.HasIsDescriptionPlainText() {
		isDescriptionPlainText := resp.GetIsDescriptionPlainText()
		providerData.IsDescriptionPlainText = &isDescriptionPlainText
	}
	
	if resp.HasIsUseDirectLinking() {
		isUseDirectLinking := resp.GetIsUseDirectLinking()
		providerData.IsUseDirectLinking = &isUseDirectLinking
	}
	
	if resp.HasAppIdentifier() {
		appIdentifier := resp.GetAppIdentifier()
		providerData.AppIdentifier = &appIdentifier
	}
	
	if resp.HasTimeCreated() {
		timeCreated := int(resp.GetTimeCreated())
		providerData.TimeCreated = &timeCreated
	}
	
	if resp.HasTimeSaved() {
		timeSaved := int(resp.GetTimeSaved())
		providerData.TimeSaved = &timeSaved
	}

	// Convert provider data to JSON and store in mapping
	providerDataJSON, err := providerData.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize provider data: %w", err)
	}
	mapping.ProviderData = &providerDataJSON

	return nil
}

// mapCampaignToEverflowUpdateRequest maps domain campaign to Everflow offer update request
func (s *IntegrationService) mapCampaignToEverflowUpdateRequest(camp *domain.Campaign) (interface{}, error) {
	// TODO: Implement campaign update mapping when campaign functionality is needed
	return nil, fmt.Errorf("campaign update mapping not implemented")
}

// GenerateTrackingLink generates a tracking link via Everflow API
func (s *IntegrationService) GenerateTrackingLink(ctx context.Context, req *domain.TrackingLinkGenerationRequest, campaignMapping *domain.CampaignProviderMapping, affiliateMapping *domain.AffiliateProviderMapping) (*domain.TrackingLinkGenerationResponse, error) {
	logger.Info("Starting tracking link generation", "campaign_id", req.CampaignID, "affiliate_id", req.AffiliateID, "provider", "everflow")
	
	// Extract provider-specific IDs from mappings
	var networkOfferID, networkAffiliateID int32

	// Parse campaign provider data to get network_offer_id
	if campaignMapping != nil {
		logger.Debug("Campaign mapping found", "campaign_id", req.CampaignID, "mapping_id", campaignMapping.MappingID, "provider_type", campaignMapping.ProviderType)
		if campaignMapping.ProviderData != nil {
			logger.Debug("Parsing campaign provider data", "campaign_id", req.CampaignID)
			logger.Debug("Campaign provider data", "campaign_id", req.CampaignID, "provider_data", *campaignMapping.ProviderData)
			
			// Try to parse as complex payload structure first (request/response format)
			var payload map[string]interface{}
			if err := json.Unmarshal([]byte(*campaignMapping.ProviderData), &payload); err == nil {
				if response, ok := payload["response"].(map[string]interface{}); ok {
					if networkOfferIDFloat, ok := response["network_offer_id"].(float64); ok {
						networkOfferID = int32(networkOfferIDFloat)
						logger.Debug("Found network_offer_id from response payload", "campaign_id", req.CampaignID, "network_offer_id", networkOfferID)
					} else {
						logger.Warn("network_offer_id not found in response payload", "campaign_id", req.CampaignID)
					}
				} else {
					logger.Warn("response object not found in payload", "campaign_id", req.CampaignID)
				}
			} else {
				// Fallback: try to parse as simple EverflowCampaignProviderData structure
				logger.Debug("Trying fallback parsing as EverflowCampaignProviderData", "campaign_id", req.CampaignID)
				var campaignProviderData domain.EverflowCampaignProviderData
				if err := campaignProviderData.FromJSON(*campaignMapping.ProviderData); err == nil {
					if campaignProviderData.NetworkCampaignID != nil {
						networkOfferID = *campaignProviderData.NetworkCampaignID
						logger.Debug("Found network_offer_id from fallback", "campaign_id", req.CampaignID, "network_offer_id", networkOfferID)
					} else {
						logger.Warn("NetworkCampaignID is nil in campaign provider data", "campaign_id", req.CampaignID)
					}
				} else {
					logger.Warn("Failed to parse campaign provider data with both methods", "campaign_id", req.CampaignID, "error", err)
				}
			}
		} else {
			logger.Warn("Campaign mapping exists but ProviderData is nil", "campaign_id", req.CampaignID)
		}
	} else {
		logger.Warn("No campaign mapping available", "campaign_id", req.CampaignID)
	}

	// Parse affiliate provider data to get network_affiliate_id
	if affiliateMapping != nil {
		logger.Debug("Affiliate mapping found", "affiliate_id", req.AffiliateID, "mapping_id", affiliateMapping.MappingID, "provider_type", affiliateMapping.ProviderType)
		if affiliateMapping.ProviderData != nil {
			logger.Debug("Parsing affiliate provider data", "affiliate_id", req.AffiliateID)
			logger.Debug("Affiliate provider data", "affiliate_id", req.AffiliateID, "provider_data", *affiliateMapping.ProviderData)
			var affiliateProviderData domain.EverflowProviderData
			if err := json.Unmarshal([]byte(*affiliateMapping.ProviderData), &affiliateProviderData); err == nil {
				if affiliateProviderData.NetworkAffiliateID != nil {
					networkAffiliateID = *affiliateProviderData.NetworkAffiliateID
					logger.Debug("Found network_affiliate_id", "affiliate_id", req.AffiliateID, "network_affiliate_id", networkAffiliateID)
				} else {
					logger.Warn("NetworkAffiliateID is nil in affiliate provider data", "affiliate_id", req.AffiliateID)
				}
			} else {
				logger.Warn("Failed to parse affiliate provider data", "affiliate_id", req.AffiliateID, "error", err)
			}
		} else {
			logger.Warn("Affiliate mapping exists but ProviderData is nil", "affiliate_id", req.AffiliateID)
		}
	} else {
		logger.Warn("No affiliate mapping available", "affiliate_id", req.AffiliateID)
	}

	// If we don't have the required IDs, return an error
	if networkOfferID == 0 || networkAffiliateID == 0 {
		logger.Error("Missing required provider IDs", "campaign_id", req.CampaignID, "affiliate_id", req.AffiliateID, "network_offer_id", networkOfferID, "network_affiliate_id", networkAffiliateID)
		return nil, fmt.Errorf("missing required provider IDs: networkOfferID=%d, networkAffiliateID=%d", networkOfferID, networkAffiliateID)
	}

	// Create Everflow tracking link request
	logger.Debug("Creating Everflow request", "campaign_id", req.CampaignID, "affiliate_id", req.AffiliateID, "network_affiliate_id", networkAffiliateID, "network_offer_id", networkOfferID)
	everflowReq := tracking.NewCreateTrackingLinkRequest(networkAffiliateID, networkOfferID)

	// Add optional parameters
	if req.NetworkTrackingDomainID != nil {
		everflowReq.SetNetworkTrackingDomainId(*req.NetworkTrackingDomainID)
	}
	if req.NetworkOfferURLID != nil {
		everflowReq.SetNetworkOfferUrlId(*req.NetworkOfferURLID)
	}
	if req.CreativeID != nil {
		everflowReq.SetCreativeId(*req.CreativeID)
	}
	if req.NetworkTrafficSourceID != nil {
		everflowReq.SetNetworkTrafficSourceId(*req.NetworkTrafficSourceID)
	}
	if req.SourceID != nil {
		everflowReq.SetSourceId(*req.SourceID)
	}
	if req.Sub1 != nil {
		everflowReq.SetSub1(*req.Sub1)
	}
	if req.Sub2 != nil {
		everflowReq.SetSub2(*req.Sub2)
	}
	if req.Sub3 != nil {
		everflowReq.SetSub3(*req.Sub3)
	}
	if req.Sub4 != nil {
		everflowReq.SetSub4(*req.Sub4)
	}
	if req.Sub5 != nil {
		everflowReq.SetSub5(*req.Sub5)
	}
	if req.IsEncryptParameters != nil {
		everflowReq.SetIsEncryptParameters(*req.IsEncryptParameters)
	}
	if req.IsRedirectLink != nil {
		everflowReq.SetIsRedirectLink(*req.IsRedirectLink)
	}

	// Log the final request for debugging (sanitized)
	reqJSON, _ := json.MarshalIndent(everflowReq, "", "  ")
	logger.DebugSanitized("Sending request to Everflow", "campaign_id", req.CampaignID, "affiliate_id", req.AffiliateID, "request", string(reqJSON))

	// Make the actual API call to Everflow
	logger.Debug("Making API call to Everflow", "campaign_id", req.CampaignID, "affiliate_id", req.AffiliateID)
	resp, httpResp, err := s.trackingClient.TrackingAPI.CreateTrackingLink(ctx).CreateTrackingLinkRequest(*everflowReq).Execute()
	if err != nil {
		// Try to read the response body for more detailed error information
		var errorBody string
		var statusCode int
		if httpResp != nil {
			statusCode = httpResp.StatusCode
			if httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					errorBody = string(bodyBytes)
				}
			}
		}
		logger.ErrorSanitized("API call failed", "campaign_id", req.CampaignID, "affiliate_id", req.AffiliateID, "status_code", statusCode, "error", err, "response_body", errorBody)
		return nil, fmt.Errorf("failed to create tracking link in Everflow: %w (response: %s)", err, errorBody)
	}
	defer httpResp.Body.Close()
	
	logger.Info("Successfully generated tracking link", "campaign_id", req.CampaignID, "affiliate_id", req.AffiliateID)

	// Extract the generated URL from the response
	generatedURL := ""
	if resp.TrackingUrl != nil {
		generatedURL = *resp.TrackingUrl
	}

	// Create provider data from response
	providerData := &domain.EverflowTrackingLinkProviderData{
		NetworkOfferID:           resp.NetworkOfferId,
		NetworkAffiliateID:       resp.NetworkAffiliateId,
		NetworkTrackingDomainID:  resp.NetworkTrackingDomainId,
		NetworkOfferURLID:        resp.NetworkOfferUrlId,
		CreativeID:               resp.CreativeId,
		NetworkTrafficSourceID:   resp.NetworkTrafficSourceId,
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

// CreateTrackingLink creates a tracking link synchronization with Everflow
func (s *IntegrationService) CreateTrackingLink(ctx context.Context, trackingLink *domain.TrackingLink, campaignMapping *domain.CampaignProviderMapping, affiliateMapping *domain.AffiliateProviderMapping) (*domain.TrackingLinkProviderMapping, error) {
	logger.Info("Starting tracking link sync", "tracking_link_id", trackingLink.TrackingLinkID, "provider", "everflow")

	// Create a tracking link generation request to get the Everflow tracking URL
	req := &domain.TrackingLinkGenerationRequest{
		CampaignID:  trackingLink.CampaignID,
		AffiliateID: trackingLink.AffiliateID,
		SourceID:    trackingLink.SourceID,
		Sub1:        trackingLink.Sub1,
		Sub2:        trackingLink.Sub2,
	}

	// Generate the tracking link via Everflow API
	response, err := s.GenerateTrackingLink(ctx, req, campaignMapping, affiliateMapping)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tracking link via Everflow: %w", err)
	}

	// Create provider mapping
	syncStatus := "synced"
	now := time.Now()
	mapping := &domain.TrackingLinkProviderMapping{
		TrackingLinkID:         trackingLink.TrackingLinkID,
		ProviderType:           "everflow",
		ProviderTrackingLinkID: nil, // Everflow doesn't assign persistent IDs to tracking links
		ProviderData:           response.ProviderData,
		SyncStatus:             &syncStatus,
		LastSyncAt:             &now,
	}

	logger.Info("Successfully created tracking link sync", "tracking_link_id", trackingLink.TrackingLinkID)
	return mapping, nil
}

// GenerateTrackingLinkQR generates a QR code for a tracking link via Everflow API
func (s *IntegrationService) GenerateTrackingLinkQR(ctx context.Context, req *domain.TrackingLinkGenerationRequest, campaignMapping *domain.CampaignProviderMapping, affiliateMapping *domain.AffiliateProviderMapping) ([]byte, error) {
	// For now, simulate the Everflow QR API call since we don't have the tracking API client generated
	// In a real implementation, this would make an HTTP POST to /v1/networks/tracking/offers/clicks/qr

	// Return a mock QR code (in real implementation, this would be a PNG image from Everflow)
	return []byte("everflow-qr-code-png-data"), nil
}



