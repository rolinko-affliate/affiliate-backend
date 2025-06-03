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
	CreateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error
	UpdateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error
	ListCampaignProviderOffersByCampaign(ctx context.Context, campaignID int64) ([]*domain.CampaignProviderOffer, error)
}

type AdvertiserProviderMappingRepository interface {
	GetAdvertiserProviderMapping(ctx context.Context, advertiserID int64, providerType string) (*domain.AdvertiserProviderMapping, error)
	CreateAdvertiserProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error
	UpdateAdvertiserProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error
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
	// Map domain advertiser to Everflow request
	everflowReq, err := s.mapAdvertiserToEverflowRequest(&adv)
	if err != nil {
		return adv, fmt.Errorf("failed to map advertiser to Everflow request: %w", err)
	}

	// Serialize the outbound request for payload storage
	requestPayload, err := json.Marshal(everflowReq)
	if err != nil {
		return adv, fmt.Errorf("failed to serialize request payload: %w", err)
	}

	// Create advertiser in Everflow
	resp, httpResp, err := s.advertiserClient.DefaultAPI.CreateAdvertiser(ctx).CreateAdvertiserRequest(*everflowReq).Execute()
	if err != nil {
		return adv, fmt.Errorf("failed to create advertiser in Everflow: %w", err)
	}
	defer httpResp.Body.Close()

	// Create provider mapping with full payload
	providerAdvertiserID := strconv.FormatInt(int64(resp.GetNetworkAdvertiserId()), 10)
	
	// Create payload with request and response
	payload := map[string]interface{}{
		"request":               json.RawMessage(requestPayload),
		"response":              resp,
		"provider_id":           resp.GetNetworkAdvertiserId(),
		"last_operation":        "create",
		"last_operation_time":   time.Now(),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return adv, fmt.Errorf("failed to marshal payload: %w", err)
	}

	payloadStr := string(payloadJSON)
	now := time.Now()

	mapping := &domain.AdvertiserProviderMapping{
		AdvertiserID:         adv.AdvertiserID,
		ProviderType:         "everflow",
		ProviderAdvertiserID: &providerAdvertiserID,
		ProviderConfig:       &payloadStr,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if err := s.advertiserProviderMappingRepo.CreateAdvertiserProviderMapping(ctx, mapping); err != nil {
		return adv, fmt.Errorf("failed to create advertiser provider mapping: %w", err)
	}

	// Map Everflow response back to domain model
	updatedAdvertiser := s.mapEverflowResponseToAdvertiser(resp, &adv)
	return updatedAdvertiser, nil
}

// UpdateAdvertiser updates an advertiser in Everflow
func (s *IntegrationService) UpdateAdvertiser(ctx context.Context, adv domain.Advertiser) error {
	// Get provider mapping
	mapping, err := s.advertiserProviderMappingRepo.GetAdvertiserProviderMapping(ctx, adv.AdvertiserID, "everflow")
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

	// Map domain advertiser to Everflow update request
	everflowReq, err := s.mapAdvertiserToEverflowUpdateRequest(&adv)
	if err != nil {
		return fmt.Errorf("failed to map advertiser to Everflow update request: %w", err)
	}

	// Update advertiser in Everflow
	resp, httpResp, err := s.advertiserClient.DefaultAPI.UpdateAdvertiser(ctx, int32(providerID)).UpdateAdvertiserRequest(*everflowReq).Execute()
	if err != nil {
		return fmt.Errorf("failed to update advertiser in Everflow: %w", err)
	}
	defer httpResp.Body.Close()

	// Update provider mapping payload
	requestPayload, _ := json.Marshal(everflowReq)
	payload := map[string]interface{}{
		"request":               json.RawMessage(requestPayload),
		"response":              resp,
		"provider_id":           resp.GetNetworkAdvertiserId(),
		"last_operation":        "update",
		"last_operation_time":   time.Now(),
	}

	payloadJSON, _ := json.Marshal(payload)
	payloadStr := string(payloadJSON)
	mapping.ProviderConfig = &payloadStr
	mapping.UpdatedAt = time.Now()

	return s.advertiserProviderMappingRepo.UpdateAdvertiserProviderMapping(ctx, mapping)
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
	mapping, err := s.advertiserProviderMappingRepo.GetAdvertiserProviderMapping(ctx, advertiserID, "everflow")
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
	updatedAdvertiser := s.mapEverflowResponseToAdvertiser(resp, adv)
	return updatedAdvertiser, nil
}

// CreateAffiliate creates an affiliate in Everflow
func (s *IntegrationService) CreateAffiliate(ctx context.Context, aff domain.Affiliate) (domain.Affiliate, error) {
	// Map domain affiliate to Everflow request
	everflowReq, err := s.mapAffiliateToEverflowRequest(&aff)
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

	// Create provider mapping with full payload
	providerAffiliateID := strconv.FormatInt(int64(resp.GetNetworkAffiliateId()), 10)
	
	// Create payload with request and response
	payload := map[string]interface{}{
		"request":               json.RawMessage(requestPayload),
		"response":              resp,
		"provider_id":           resp.GetNetworkAffiliateId(),
		"last_operation":        "create",
		"last_operation_time":   time.Now(),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return aff, fmt.Errorf("failed to marshal payload: %w", err)
	}

	payloadStr := string(payloadJSON)
	now := time.Now()

	mapping := &domain.AffiliateProviderMapping{
		AffiliateID:         aff.AffiliateID,
		ProviderType:        "everflow",
		ProviderAffiliateID: &providerAffiliateID,
		ProviderConfig:      &payloadStr,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := s.affiliateProviderMappingRepo.CreateAffiliateProviderMapping(ctx, mapping); err != nil {
		return aff, fmt.Errorf("failed to create affiliate provider mapping: %w", err)
	}

	// Map Everflow response back to domain model
	updatedAffiliate := s.mapEverflowCreateResponseToAffiliate(resp, &aff)
	return updatedAffiliate, nil
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

	// Map domain affiliate to Everflow update request
	everflowReq, err := s.mapAffiliateToEverflowUpdateRequest(&aff)
	if err != nil {
		return fmt.Errorf("failed to map affiliate to Everflow update request: %w", err)
	}

	// Update affiliate in Everflow
	resp, httpResp, err := s.affiliateClient.AffiliatesAPI.UpdateAffiliate(ctx, int32(providerID)).UpdateAffiliateRequest(*everflowReq).Execute()
	if err != nil {
		return fmt.Errorf("failed to update affiliate in Everflow: %w", err)
	}
	defer httpResp.Body.Close()

	// Update provider mapping payload
	requestPayload, _ := json.Marshal(everflowReq)
	payload := map[string]interface{}{
		"request":               json.RawMessage(requestPayload),
		"response":              resp,
		"provider_id":           resp.GetNetworkAffiliateId(),
		"last_operation":        "update",
		"last_operation_time":   time.Now(),
	}

	payloadJSON, _ := json.Marshal(payload)
	payloadStr := string(payloadJSON)
	mapping.ProviderConfig = &payloadStr
	mapping.UpdatedAt = time.Now()

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
	// Get the network advertiser ID from the advertiser provider mapping
	advertiserMapping, err := s.advertiserProviderMappingRepo.GetAdvertiserProviderMapping(ctx, camp.AdvertiserID, "everflow")
	if err != nil {
		return camp, fmt.Errorf("failed to get advertiser provider mapping: %w", err)
	}

	if advertiserMapping.ProviderAdvertiserID == nil {
		return camp, fmt.Errorf("advertiser not found in Everflow")
	}

	networkAdvertiserID, err := strconv.ParseInt(*advertiserMapping.ProviderAdvertiserID, 10, 32)
	if err != nil {
		return camp, fmt.Errorf("invalid network advertiser ID: %w", err)
	}

	// Map domain campaign to Everflow offer request
	everflowReq, err := s.mapCampaignToEverflowRequest(&camp, int32(networkAdvertiserID))
	if err != nil {
		return camp, fmt.Errorf("failed to map campaign to Everflow request: %w", err)
	}

	// Serialize the outbound request for payload storage
	requestPayload, err := json.Marshal(everflowReq)
	if err != nil {
		return camp, fmt.Errorf("failed to serialize request payload: %w", err)
	}

	// Create offer in Everflow
	resp, httpResp, err := s.offerClient.OffersAPI.CreateOffer(ctx).CreateOfferRequest(*everflowReq).Execute()
	if err != nil {
		return camp, fmt.Errorf("failed to create offer in Everflow: %w", err)
	}
	defer httpResp.Body.Close()

	// Create provider mapping with full payload
	providerOfferID := strconv.FormatInt(int64(resp.GetNetworkOfferId()), 10)
	
	// Create payload with request and response
	payload := map[string]interface{}{
		"request":               json.RawMessage(requestPayload),
		"response":              resp,
		"provider_id":           resp.GetNetworkOfferId(),
		"last_operation":        "create",
		"last_operation_time":   time.Now(),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return camp, fmt.Errorf("failed to marshal payload: %w", err)
	}

	payloadStr := string(payloadJSON)
	now := time.Now()

	// Create both CampaignProviderOffer (existing) and CampaignProviderMapping (new)
	providerOffer := &domain.CampaignProviderOffer{
		CampaignID:          camp.CampaignID,
		ProviderType:        "everflow",
		ProviderOfferRef:    &providerOfferID,
		ProviderOfferConfig: &payloadStr,
		IsActiveOnProvider:  true,
		LastSyncedAt:        &now,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := s.campaignRepo.CreateCampaignProviderOffer(ctx, providerOffer); err != nil {
		return camp, fmt.Errorf("failed to create campaign provider offer: %w", err)
	}

	mapping := &domain.CampaignProviderMapping{
		CampaignID:         camp.CampaignID,
		ProviderType:       "everflow",
		ProviderCampaignID: &providerOfferID,
		ProviderConfig:     &payloadStr,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.campaignProviderMappingRepo.CreateCampaignProviderMapping(ctx, mapping); err != nil {
		return camp, fmt.Errorf("failed to create campaign provider mapping: %w", err)
	}

	// Map Everflow response back to domain model
	updatedCampaign := s.mapEverflowResponseToCampaign(resp, &camp)
	return updatedCampaign, nil
}

// UpdateCampaign updates a campaign (offer) in Everflow
func (s *IntegrationService) UpdateCampaign(ctx context.Context, camp domain.Campaign) error {
	// Get provider mapping
	mapping, err := s.campaignProviderMappingRepo.GetCampaignProviderMapping(ctx, camp.CampaignID, "everflow")
	if err != nil {
		return fmt.Errorf("failed to get campaign provider mapping: %w", err)
	}

	if mapping.ProviderCampaignID == nil {
		return fmt.Errorf("campaign not found in Everflow")
	}

	providerID, err := strconv.ParseInt(*mapping.ProviderCampaignID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid provider campaign ID: %w", err)
	}

	// Map domain campaign to Everflow update request
	everflowReq, err := s.mapCampaignToEverflowUpdateRequest(&camp)
	if err != nil {
		return fmt.Errorf("failed to map campaign to Everflow update request: %w", err)
	}

	// Update offer in Everflow
	resp, httpResp, err := s.offerClient.OffersAPI.UpdateOffer(ctx, int32(providerID)).UpdateOfferRequest(*everflowReq).Execute()
	if err != nil {
		return fmt.Errorf("failed to update offer in Everflow: %w", err)
	}
	defer httpResp.Body.Close()

	// Update provider mapping payload
	requestPayload, _ := json.Marshal(everflowReq)
	payload := map[string]interface{}{
		"request":               json.RawMessage(requestPayload),
		"response":              resp,
		"provider_id":           resp.GetNetworkOfferId(),
		"last_operation":        "update",
		"last_operation_time":   time.Now(),
	}

	payloadJSON, _ := json.Marshal(payload)
	payloadStr := string(payloadJSON)
	mapping.ProviderConfig = &payloadStr
	mapping.UpdatedAt = time.Now()

	// Update both mappings
	if err := s.campaignProviderMappingRepo.UpdateCampaignProviderMapping(ctx, mapping); err != nil {
		return fmt.Errorf("failed to update campaign provider mapping: %w", err)
	}

	// Also update the CampaignProviderOffer
	offers, err := s.campaignRepo.ListCampaignProviderOffersByCampaign(ctx, camp.CampaignID)
	if err != nil {
		return fmt.Errorf("failed to get campaign provider offers: %w", err)
	}

	for _, offer := range offers {
		if offer.ProviderType == "everflow" {
			offer.ProviderOfferConfig = &payloadStr
			offer.UpdatedAt = time.Now()
			if err := s.campaignRepo.UpdateCampaignProviderOffer(ctx, offer); err != nil {
				return fmt.Errorf("failed to update campaign provider offer: %w", err)
			}
			break
		}
	}

	return nil
}

// GetCampaign retrieves a campaign from Everflow
func (s *IntegrationService) GetCampaign(ctx context.Context, id uuid.UUID) (domain.Campaign, error) {
	// Convert UUID to int64
	campaignID, err := uuidToInt64(id)
	if err != nil {
		return domain.Campaign{}, fmt.Errorf("failed to convert UUID to int64: %w", err)
	}

	// Get local campaign
	camp, err := s.campaignRepo.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return domain.Campaign{}, fmt.Errorf("failed to get local campaign: %w", err)
	}

	// Get provider mapping
	mapping, err := s.campaignProviderMappingRepo.GetCampaignProviderMapping(ctx, campaignID, "everflow")
	if err != nil {
		return *camp, fmt.Errorf("failed to get campaign provider mapping: %w", err)
	}

	if mapping.ProviderCampaignID == nil {
		return *camp, fmt.Errorf("campaign not found in Everflow")
	}

	providerID, err := strconv.ParseInt(*mapping.ProviderCampaignID, 10, 32)
	if err != nil {
		return *camp, fmt.Errorf("invalid provider campaign ID: %w", err)
	}

	// Get offer from Everflow
	resp, httpResp, err := s.offerClient.OffersAPI.GetOfferById(ctx, int32(providerID)).Execute()
	if err != nil {
		return *camp, fmt.Errorf("failed to get offer from Everflow: %w", err)
	}
	defer httpResp.Body.Close()

	// Map Everflow response to domain model
	updatedCampaign := s.mapEverflowResponseToCampaign(resp, camp)
	return updatedCampaign, nil
}