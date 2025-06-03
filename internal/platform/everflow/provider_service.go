package everflow

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/affiliate-backend/internal/config"
	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/crypto"
	"github.com/affiliate-backend/internal/platform/provider"
	"github.com/affiliate-backend/internal/repository"
)

// ProviderService implements both ProviderAdvertiserService and ProviderCampaignService for Everflow
type ProviderService struct {
	client              *Client
	advertiserRepo      repository.AdvertiserRepository
	providerMappingRepo repository.AdvertiserProviderMappingRepository
	campaignRepo        repository.CampaignRepository
	cryptoService       crypto.Service
}

// NewProviderService creates a new Everflow provider service
func NewProviderService(
	apiKey string,
	cfg *config.Config,
	advertiserRepo repository.AdvertiserRepository,
	providerMappingRepo repository.AdvertiserProviderMappingRepository,
	campaignRepo repository.CampaignRepository,
	cryptoService crypto.Service,
) *ProviderService {
	return &ProviderService{
		client:              NewClient(apiKey, cfg),
		advertiserRepo:      advertiserRepo,
		providerMappingRepo: providerMappingRepo,
		campaignRepo:        campaignRepo,
		cryptoService:       cryptoService,
	}
}

// Ensure ProviderService implements the interfaces
var _ provider.ProviderAdvertiserService = (*ProviderService)(nil)
var _ provider.ProviderCampaignService = (*ProviderService)(nil)

// ProviderAdvertiserService implementation

// CreateAdvertiserInProvider creates an advertiser in Everflow and stores the mapping
func (s *ProviderService) CreateAdvertiserInProvider(ctx context.Context, advertiser *domain.Advertiser) error {
	// Map our advertiser to Everflow advertiser
	everflowReq, err := s.mapAdvertiserToEverflowRequest(advertiser)
	if err != nil {
		return fmt.Errorf("failed to map advertiser to Everflow request: %w", err)
	}

	// Create advertiser in Everflow
	everflowResp, err := s.client.CreateAdvertiser(ctx, *everflowReq)
	if err != nil {
		return fmt.Errorf("failed to create advertiser in Everflow: %w", err)
	}

	// Add tags to the advertiser in Everflow
	tags := []string{
		fmt.Sprintf("advertiser_id:%d", advertiser.AdvertiserID),
		fmt.Sprintf("organization_id:%d", advertiser.OrganizationID),
	}

	if err := s.client.AddTagsToAdvertiser(ctx, everflowResp.NetworkAdvertiserID, tags); err != nil {
		// Log the error but continue, as this is not critical
		fmt.Printf("Warning: failed to add tags to advertiser in Everflow: %v\n", err)
	}

	// Create mapping in our database
	providerAdvertiserID := strconv.FormatInt(everflowResp.NetworkAdvertiserID, 10)

	// Create a provider config to store additional Everflow-specific data
	providerConfig := map[string]interface{}{
		"network_advertiser_id": everflowResp.NetworkAdvertiserID,
		"account_status":        everflowResp.AccountStatus,
		"default_currency_id":   everflowResp.DefaultCurrencyID,
		"time_created":          everflowResp.TimeCreated,
		"time_saved":            everflowResp.TimeSaved,
	}

	providerConfigJSON, err := json.Marshal(providerConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal provider config: %w", err)
	}

	providerConfigStr := string(providerConfigJSON)

	mapping := &domain.AdvertiserProviderMapping{
		AdvertiserID:         advertiser.AdvertiserID,
		ProviderType:         "everflow",
		ProviderAdvertiserID: &providerAdvertiserID,
		ProviderConfig:       &providerConfigStr,
	}

	if err := s.providerMappingRepo.CreateAdvertiserProviderMapping(ctx, mapping); err != nil {
		return fmt.Errorf("failed to create advertiser provider mapping: %w", err)
	}

	return nil
}

// GetAdvertiserFromProvider retrieves an advertiser from Everflow using our internal advertiser ID
func (s *ProviderService) GetAdvertiserFromProvider(ctx context.Context, advertiserID int64, relationships []string) (*domain.Advertiser, error) {
	// Get the advertiser's Everflow mapping
	mapping, err := s.providerMappingRepo.GetAdvertiserProviderMapping(ctx, advertiserID, "everflow")
	if err != nil {
		return nil, fmt.Errorf("failed to get advertiser provider mapping: %w", err)
	}

	if mapping.ProviderAdvertiserID == nil {
		return nil, fmt.Errorf("advertiser does not have an Everflow ID")
	}

	// Convert the provider advertiser ID to int64
	networkAdvertiserID, err := strconv.ParseInt(*mapping.ProviderAdvertiserID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid provider advertiser ID: %w", err)
	}

	// Get advertiser from Everflow
	opts := &GetAdvertiserOptions{
		Relationships: relationships,
	}

	everflowAdvertiser, err := s.client.GetAdvertiser(ctx, networkAdvertiserID, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get advertiser from Everflow: %w", err)
	}

	// Get the local advertiser
	localAdvertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, advertiserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get local advertiser: %w", err)
	}

	// Map Everflow data to our domain model
	result := s.mapEverflowAdvertiserToDomain(everflowAdvertiser, localAdvertiser)
	return result, nil
}

// UpdateAdvertiserInProvider updates an advertiser in Everflow using our internal advertiser ID
func (s *ProviderService) UpdateAdvertiserInProvider(ctx context.Context, advertiserID int64, advertiser *domain.Advertiser) (*domain.Advertiser, error) {
	// Get the advertiser's Everflow mapping
	mapping, err := s.providerMappingRepo.GetAdvertiserProviderMapping(ctx, advertiserID, "everflow")
	if err != nil {
		return nil, fmt.Errorf("failed to get advertiser provider mapping: %w", err)
	}

	if mapping.ProviderAdvertiserID == nil {
		return nil, fmt.Errorf("advertiser does not have an Everflow ID")
	}

	// Convert the provider advertiser ID to int64
	networkAdvertiserID, err := strconv.ParseInt(*mapping.ProviderAdvertiserID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid provider advertiser ID: %w", err)
	}

	// Map our advertiser to Everflow update request
	updateReq, err := s.mapAdvertiserToEverflowUpdateRequest(advertiser)
	if err != nil {
		return nil, fmt.Errorf("failed to map advertiser to Everflow update request: %w", err)
	}

	// Update advertiser in Everflow
	updatedEverflowAdvertiser, err := s.client.UpdateAdvertiser(ctx, networkAdvertiserID, *updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update advertiser in Everflow: %w", err)
	}

	// Map updated Everflow data back to our domain model
	result := s.mapEverflowAdvertiserToDomain(updatedEverflowAdvertiser, advertiser)
	return result, nil
}

// ProviderCampaignService implementation

// CreateOfferInProvider creates an offer in Everflow for a campaign and stores the mapping
func (s *ProviderService) CreateOfferInProvider(ctx context.Context, campaign *domain.Campaign) error {
	// Get the advertiser
	advertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, campaign.AdvertiserID)
	if err != nil {
		return fmt.Errorf("failed to get advertiser: %w", err)
	}

	// Get the advertiser's Everflow mapping
	mapping, err := s.providerMappingRepo.GetAdvertiserProviderMapping(ctx, advertiser.AdvertiserID, "everflow")
	if err != nil {
		return fmt.Errorf("failed to get advertiser provider mapping: %w", err)
	}

	if mapping.ProviderAdvertiserID == nil {
		return fmt.Errorf("advertiser does not have an Everflow ID")
	}

	// Convert the provider advertiser ID to int64
	networkAdvertiserID, err := strconv.ParseInt(*mapping.ProviderAdvertiserID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid provider advertiser ID: %w", err)
	}

	// Map our campaign to Everflow offer
	everflowReq, err := s.mapCampaignToEverflowRequest(campaign, networkAdvertiserID)
	if err != nil {
		return fmt.Errorf("failed to map campaign to Everflow request: %w", err)
	}

	// Create offer in Everflow
	everflowResp, err := s.client.CreateOffer(ctx, *everflowReq)
	if err != nil {
		return fmt.Errorf("failed to create offer in Everflow: %w", err)
	}

	// Add tags to the offer in Everflow
	tags := []string{
		fmt.Sprintf("campaign_id:%d", campaign.CampaignID),
		fmt.Sprintf("advertiser_id:%d", campaign.AdvertiserID),
		fmt.Sprintf("organization_id:%d", campaign.OrganizationID),
	}

	if err := s.client.AddTagsToOffer(ctx, everflowResp.NetworkOfferID, tags); err != nil {
		// Log the error but continue, as this is not critical
		fmt.Printf("Warning: failed to add tags to offer in Everflow: %v\n", err)
	}

	// Create mapping in our database
	providerOfferRef := strconv.FormatInt(everflowResp.NetworkOfferID, 10)

	// Create a provider config to store additional Everflow-specific data
	providerConfig := map[string]interface{}{
		"network_offer_id":      everflowResp.NetworkOfferID,
		"network_id":            everflowResp.NetworkID,
		"network_advertiser_id": everflowResp.NetworkAdvertiserID,
		"offer_status":          everflowResp.OfferStatus,
		"currency_id":           everflowResp.CurrencyID,
		"encoded_value":         everflowResp.EncodedValue,
		"time_created":          everflowResp.TimeCreated,
		"time_saved":            everflowResp.TimeSaved,
	}

	providerConfigJSON, err := json.Marshal(providerConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal provider config: %w", err)
	}

	providerConfigStr := string(providerConfigJSON)
	now := time.Now()

	offer := &domain.CampaignProviderOffer{
		CampaignID:          campaign.CampaignID,
		ProviderType:        "everflow",
		ProviderOfferRef:    &providerOfferRef,
		ProviderOfferConfig: &providerConfigStr,
		IsActiveOnProvider:  everflowResp.OfferStatus == "active",
		LastSyncedAt:        &now,
	}

	if err := s.campaignRepo.CreateCampaignProviderOffer(ctx, offer); err != nil {
		return fmt.Errorf("failed to create campaign provider offer: %w", err)
	}

	return nil
}

// GetOfferFromProvider retrieves an offer from Everflow using our internal campaign ID
func (s *ProviderService) GetOfferFromProvider(ctx context.Context, campaignID int64, relationships []string) (*domain.Campaign, error) {
	// Get the campaign's Everflow offer mapping
	offers, err := s.campaignRepo.ListCampaignProviderOffersByCampaign(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign provider offers: %w", err)
	}

	// Find the Everflow offer
	var everflowOffer *domain.CampaignProviderOffer
	for _, offer := range offers {
		if offer.ProviderType == "everflow" && offer.ProviderOfferRef != nil {
			everflowOffer = offer
			break
		}
	}

	if everflowOffer == nil {
		return nil, fmt.Errorf("campaign does not have an Everflow offer")
	}

	// Convert the provider offer ref to int64
	networkOfferID, err := strconv.ParseInt(*everflowOffer.ProviderOfferRef, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid provider offer ref: %w", err)
	}

	// Get offer from Everflow
	opts := &GetOfferOptions{
		Relationships: relationships,
	}

	everflowOfferData, err := s.client.GetOffer(ctx, networkOfferID, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get offer from Everflow: %w", err)
	}

	// Get the local campaign
	localCampaign, err := s.campaignRepo.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get local campaign: %w", err)
	}

	// Map Everflow offer data to our domain model
	result := s.mapEverflowOfferToDomain(everflowOfferData, localCampaign)
	return result, nil
}

// UpdateOfferInProvider updates an offer in Everflow using our internal campaign ID
func (s *ProviderService) UpdateOfferInProvider(ctx context.Context, campaignID int64, campaign *domain.Campaign) (*domain.Campaign, error) {
	// Get the campaign's Everflow offer mapping
	offers, err := s.campaignRepo.ListCampaignProviderOffersByCampaign(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign provider offers: %w", err)
	}

	// Find the Everflow offer
	var everflowOffer *domain.CampaignProviderOffer
	for _, offer := range offers {
		if offer.ProviderType == "everflow" && offer.ProviderOfferRef != nil {
			everflowOffer = offer
			break
		}
	}

	if everflowOffer == nil {
		return nil, fmt.Errorf("campaign does not have an Everflow offer")
	}

	// Convert the provider offer ref to int64
	networkOfferID, err := strconv.ParseInt(*everflowOffer.ProviderOfferRef, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid provider offer ref: %w", err)
	}

	// Get the advertiser's Everflow mapping for network_advertiser_id
	advertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, campaign.AdvertiserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get advertiser: %w", err)
	}

	mapping, err := s.providerMappingRepo.GetAdvertiserProviderMapping(ctx, advertiser.AdvertiserID, "everflow")
	if err != nil {
		return nil, fmt.Errorf("failed to get advertiser provider mapping: %w", err)
	}

	if mapping.ProviderAdvertiserID == nil {
		return nil, fmt.Errorf("advertiser does not have an Everflow ID")
	}

	networkAdvertiserID, err := strconv.ParseInt(*mapping.ProviderAdvertiserID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid provider advertiser ID: %w", err)
	}

	// Map our campaign to Everflow offer update request
	updateReq, err := s.mapCampaignToEverflowRequest(campaign, networkAdvertiserID)
	if err != nil {
		return nil, fmt.Errorf("failed to map campaign to Everflow request: %w", err)
	}

	// Update the offer in Everflow
	updatedOffer, err := s.client.UpdateOffer(ctx, networkOfferID, *updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update offer in Everflow: %w", err)
	}

	// Update the local mapping with new status and sync time
	now := time.Now()
	everflowOffer.IsActiveOnProvider = updatedOffer.OfferStatus == "active"
	everflowOffer.LastSyncedAt = &now

	// Update provider config with new data
	providerConfig := map[string]interface{}{
		"network_offer_id":      updatedOffer.NetworkOfferID,
		"network_id":            updatedOffer.NetworkID,
		"network_advertiser_id": updatedOffer.NetworkAdvertiserID,
		"offer_status":          updatedOffer.OfferStatus,
		"currency_id":           updatedOffer.CurrencyID,
		"encoded_value":         updatedOffer.EncodedValue,
		"time_created":          updatedOffer.TimeCreated,
		"time_saved":            updatedOffer.TimeSaved,
	}

	providerConfigJSON, err := json.Marshal(providerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal provider config: %w", err)
	}

	providerConfigStr := string(providerConfigJSON)
	everflowOffer.ProviderOfferConfig = &providerConfigStr

	if err := s.campaignRepo.UpdateCampaignProviderOffer(ctx, everflowOffer); err != nil {
		// Log the error but don't fail the operation
		fmt.Printf("Warning: failed to update campaign provider offer mapping: %v\n", err)
	}

	// Map updated Everflow offer data back to our domain model
	result := s.mapEverflowOfferToDomain(updatedOffer, campaign)
	return result, nil
}

// Helper methods for mapping between domain models and Everflow models

// mapAdvertiserToEverflowRequest maps our advertiser to an Everflow advertiser request
func (s *ProviderService) mapAdvertiserToEverflowRequest(advertiser *domain.Advertiser) (*EverflowCreateAdvertiserRequest, error) {
	// Map status
	accountStatus := "active"
	switch advertiser.Status {
	case "active":
		accountStatus = "active"
	case "inactive", "rejected":
		accountStatus = "inactive"
	case "pending":
		accountStatus = "pending"
	}

	// Determine default currency
	defaultCurrency := "USD"
	if advertiser.DefaultCurrencyID != nil {
		defaultCurrency = *advertiser.DefaultCurrencyID
	}

	// Create basic request
	req := &EverflowCreateAdvertiserRequest{
		Name:              advertiser.Name,
		AccountStatus:     accountStatus,
		DefaultCurrencyID: defaultCurrency,
	}

	// Map Everflow-specific fields
	if advertiser.InternalNotes != nil {
		req.InternalNotes = advertiser.InternalNotes
	} else if advertiser.ContactEmail != nil {
		// Fallback: Add contact email as internal notes if no internal notes but contact email exists
		notes := fmt.Sprintf("Contact Email: %s", *advertiser.ContactEmail)
		req.InternalNotes = &notes
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

	if advertiser.ReportingTimezoneID != nil {
		req.ReportingTimezoneID = advertiser.ReportingTimezoneID
	}

	if advertiser.IsExposePublisherReporting != nil {
		req.IsExposePublisherReportingData = advertiser.IsExposePublisherReporting
	}

	// Map billing details if available
	if advertiser.BillingDetails != nil {
		// Map billing information to Everflow billing structure
		req.Billing = &AdvertiserBilling{
			BillingFrequency:           advertiser.BillingDetails.BillingFrequency,
			TaxID:                      advertiser.BillingDetails.TaxID,
			IsInvoiceCreationAuto:      advertiser.BillingDetails.IsInvoiceCreationAuto,
			InvoiceAmountThreshold:     advertiser.BillingDetails.InvoiceAmountThreshold,
			AutoInvoiceStartDate:       advertiser.BillingDetails.AutoInvoiceStartDate,
			DefaultInvoiceIsHidden:     advertiser.BillingDetails.DefaultInvoiceIsHidden,
			InvoiceGenerationDaysDelay: advertiser.BillingDetails.InvoiceGenerationDaysDelay,
			DefaultPaymentTerms:        advertiser.BillingDetails.DefaultPaymentTerms,
			Details:                    advertiser.BillingDetails.AdditionalDetails,
		}

		// If we have an address in billing details, add it to the request
		if advertiser.BillingDetails.Address != nil {
			isContactAddressEnabled := true
			req.IsContactAddressEnabled = &isContactAddressEnabled

			address := advertiser.BillingDetails.Address
			address2 := ""
			if address.Line2 != nil {
				address2 = *address.Line2
			}
			region := "CA" // Default region
			if address.State != nil {
				region = *address.State
			}

			req.ContactAddress = &AdvertiserAddress{
				Address1:      address.Line1,
				Address2:      &address2,
				City:          address.City,
				ZipPostalCode: address.PostalCode,
				CountryCode:   address.Country,
				RegionCode:    region,
			}
		}
	}

	return req, nil
}

// mapAdvertiserToEverflowUpdateRequest maps our advertiser to an Everflow update request
func (s *ProviderService) mapAdvertiserToEverflowUpdateRequest(advertiser *domain.Advertiser) (*EverflowUpdateAdvertiserRequest, error) {
	// Create update request based on the advertiser data
	req := &EverflowUpdateAdvertiserRequest{
		Name: advertiser.Name,
	}

	// Map status
	if advertiser.Status != "" {
		var accountStatus string
		switch advertiser.Status {
		case "active":
			accountStatus = "active"
		case "inactive", "rejected":
			accountStatus = "inactive"
		case "pending":
			accountStatus = "pending"
		}
		req.AccountStatus = accountStatus
	}

	// Map other fields if they exist
	if advertiser.DefaultCurrencyID != nil {
		req.DefaultCurrencyID = *advertiser.DefaultCurrencyID
	}

	if advertiser.InternalNotes != nil {
		req.InternalNotes = advertiser.InternalNotes
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

	if advertiser.ReportingTimezoneID != nil {
		req.ReportingTimezoneID = *advertiser.ReportingTimezoneID
	}

	if advertiser.IsExposePublisherReporting != nil {
		req.IsExposePublisherReportingData = advertiser.IsExposePublisherReporting
	}

	return req, nil
}

// mapCampaignToEverflowRequest maps our campaign to an Everflow offer request
func (s *ProviderService) mapCampaignToEverflowRequest(campaign *domain.Campaign, networkAdvertiserID int64) (*OfferInput, error) {
	// Map status
	offerStatus := "pending"
	switch campaign.Status {
	case "active":
		offerStatus = "active"
	case "paused":
		offerStatus = "paused"
	case "draft", "archived":
		offerStatus = "pending"
	}

	// Create a default destination URL if not provided
	destinationURL := fmt.Sprintf("https://example.com/campaigns/%d?click_id={transaction_id}", campaign.CampaignID)
	if campaign.DestinationURL != nil && *campaign.DestinationURL != "" {
		destinationURL = *campaign.DestinationURL
	}

	// Default values
	currencyID := "USD"
	visibility := "public"
	conversionMethod := "server_postback"
	sessionDefinition := "cookie"
	sessionDuration := 720 // 30 days in hours

	// Use campaign values if provided
	if campaign.CurrencyID != nil {
		currencyID = *campaign.CurrencyID
	}
	if campaign.Visibility != nil {
		visibility = *campaign.Visibility
	}
	if campaign.ConversionMethod != nil {
		conversionMethod = *campaign.ConversionMethod
	}
	if campaign.SessionDefinition != nil {
		sessionDefinition = *campaign.SessionDefinition
	}
	if campaign.SessionDuration != nil {
		sessionDuration = *campaign.SessionDuration
	}

	// Create payout/revenue structure
	payoutAmount := 1.00
	revenueAmount := 2.00
	payoutType := "cpa"
	revenueType := "rpa"

	if campaign.PayoutAmount != nil {
		payoutAmount = *campaign.PayoutAmount
	}
	if campaign.RevenueAmount != nil {
		revenueAmount = *campaign.RevenueAmount
	}
	if campaign.PayoutType != nil {
		payoutType = *campaign.PayoutType
	}
	if campaign.RevenueType != nil {
		revenueType = *campaign.RevenueType
	}

	payoutRevenue := &PayoutRevenue{
		Entries: []PayoutRevenueEntry{
			{
				EntryName:     nil, // Empty entry name for default
				PayoutType:    payoutType,
				PayoutAmount:  &payoutAmount,
				RevenueType:   revenueType,
				RevenueAmount: &revenueAmount,
				IsDefault:     true,
				IsPrivate:     false,
			},
		},
	}

	// Create basic request
	req := &OfferInput{
		Name:                campaign.Name,
		NetworkAdvertiserID: networkAdvertiserID,
		DestinationURL:      destinationURL,
		OfferStatus:         offerStatus,
		CurrencyID:          &currencyID,
		Visibility:          &visibility,
		ConversionMethod:    &conversionMethod,
		SessionDefinition:   &sessionDefinition,
		SessionDuration:     &sessionDuration,
		PayoutRevenue:       payoutRevenue,
	}

	// Add description if available
	if campaign.Description != nil {
		req.HTMLDescription = campaign.Description
	}

	// Add other optional fields
	if campaign.ThumbnailURL != nil {
		req.ThumbnailURL = campaign.ThumbnailURL
	}
	if campaign.PreviewURL != nil {
		req.PreviewURL = campaign.PreviewURL
	}
	if campaign.InternalNotes != nil {
		req.InternalNotes = campaign.InternalNotes
	}
	if campaign.TermsAndConditions != nil {
		req.TermsAndConditions = campaign.TermsAndConditions
	}
	if campaign.IsForceTermsAndConditions != nil {
		req.IsForceTermsAndConditions = campaign.IsForceTermsAndConditions
	}

	return req, nil
}

// mapEverflowAdvertiserToDomain maps Everflow advertiser data to our domain model
func (s *ProviderService) mapEverflowAdvertiserToDomain(everflowAdvertiser *Advertiser, localAdvertiser *domain.Advertiser) *domain.Advertiser {
	// Start with the local advertiser as base
	result := *localAdvertiser

	// Update with Everflow data
	result.Name = everflowAdvertiser.Name

	// Map status
	switch everflowAdvertiser.AccountStatus {
	case "active":
		result.Status = "active"
	case "inactive":
		result.Status = "inactive"
	case "pending":
		result.Status = "pending"
	default:
		// Keep existing status if unknown
	}

	// Map Everflow-specific fields
	if everflowAdvertiser.DefaultCurrencyID != "" {
		result.DefaultCurrencyID = &everflowAdvertiser.DefaultCurrencyID
	}

	if everflowAdvertiser.InternalNotes != "" {
		result.InternalNotes = &everflowAdvertiser.InternalNotes
	}

	if everflowAdvertiser.PlatformName != "" {
		result.PlatformName = &everflowAdvertiser.PlatformName
	}

	if everflowAdvertiser.PlatformURL != "" {
		result.PlatformURL = &everflowAdvertiser.PlatformURL
	}

	if everflowAdvertiser.PlatformUsername != "" {
		result.PlatformUsername = &everflowAdvertiser.PlatformUsername
	}

	if everflowAdvertiser.AccountingContactEmail != "" {
		result.AccountingContactEmail = &everflowAdvertiser.AccountingContactEmail
	}

	if everflowAdvertiser.OfferIDMacro != "" {
		result.OfferIDMacro = &everflowAdvertiser.OfferIDMacro
	}

	if everflowAdvertiser.AffiliateIDMacro != "" {
		result.AffiliateIDMacro = &everflowAdvertiser.AffiliateIDMacro
	}

	if everflowAdvertiser.AttributionMethod != "" {
		result.AttributionMethod = &everflowAdvertiser.AttributionMethod
	}

	if everflowAdvertiser.EmailAttributionMethod != "" {
		result.EmailAttributionMethod = &everflowAdvertiser.EmailAttributionMethod
	}

	if everflowAdvertiser.AttributionPriority != "" {
		result.AttributionPriority = &everflowAdvertiser.AttributionPriority
	}

	if everflowAdvertiser.ReportingTimezoneID != 0 {
		result.ReportingTimezoneID = &everflowAdvertiser.ReportingTimezoneID
	}

	if everflowAdvertiser.IsExposePublisherReportingData != nil {
		result.IsExposePublisherReporting = everflowAdvertiser.IsExposePublisherReportingData
	}

	return &result
}

// mapEverflowOfferToDomain maps Everflow offer data to our domain model
func (s *ProviderService) mapEverflowOfferToDomain(everflowOffer *Offer, localCampaign *domain.Campaign) *domain.Campaign {
	// Start with the local campaign as base
	result := *localCampaign

	// Update with Everflow data
	result.Name = everflowOffer.Name

	// Map status
	switch everflowOffer.OfferStatus {
	case "active":
		result.Status = "active"
	case "paused":
		result.Status = "paused"
	case "pending":
		result.Status = "draft"
	default:
		// Keep existing status if unknown
	}

	// Map offer fields
	if everflowOffer.DestinationURL != "" {
		result.DestinationURL = &everflowOffer.DestinationURL
	}

	if everflowOffer.ThumbnailURL != nil {
		result.ThumbnailURL = everflowOffer.ThumbnailURL
	}

	if everflowOffer.PreviewURL != nil {
		result.PreviewURL = everflowOffer.PreviewURL
	}

	if everflowOffer.Visibility != nil {
		result.Visibility = everflowOffer.Visibility
	}

	if everflowOffer.CurrencyID != nil && *everflowOffer.CurrencyID != "" {
		result.CurrencyID = everflowOffer.CurrencyID
	}

	if everflowOffer.ConversionMethod != nil {
		result.ConversionMethod = everflowOffer.ConversionMethod
	}

	if everflowOffer.SessionDefinition != nil {
		result.SessionDefinition = everflowOffer.SessionDefinition
	}

	if everflowOffer.SessionDuration != nil {
		result.SessionDuration = everflowOffer.SessionDuration
	}

	if everflowOffer.InternalNotes != nil {
		result.InternalNotes = everflowOffer.InternalNotes
	}

	if everflowOffer.TermsAndConditions != nil {
		result.TermsAndConditions = everflowOffer.TermsAndConditions
	}

	if everflowOffer.IsForceTermsAndConditions != nil {
		result.IsForceTermsAndConditions = everflowOffer.IsForceTermsAndConditions
	}

	// Map tracking fields
	if everflowOffer.EncodedValue != nil {
		result.EncodedValue = everflowOffer.EncodedValue
	}

	if everflowOffer.TodayClicks != nil {
		result.TodayClicks = everflowOffer.TodayClicks
	}

	if everflowOffer.TodayRevenue != nil {
		result.TodayRevenue = everflowOffer.TodayRevenue
	}

	if everflowOffer.TimeCreated != nil {
		timeCreated := int(*everflowOffer.TimeCreated)
		result.TimeCreated = &timeCreated
	}

	if everflowOffer.TimeSaved != nil {
		timeSaved := int(*everflowOffer.TimeSaved)
		result.TimeSaved = &timeSaved
	}

	// Map payout/revenue from the first entry if available
	if everflowOffer.PayoutRevenue != nil && len(everflowOffer.PayoutRevenue.Entries) > 0 {
		entry := everflowOffer.PayoutRevenue.Entries[0]
		if entry.PayoutType != "" {
			result.PayoutType = &entry.PayoutType
		}
		if entry.PayoutAmount != nil {
			result.PayoutAmount = entry.PayoutAmount
		}
		if entry.RevenueType != "" {
			result.RevenueType = &entry.RevenueType
		}
		if entry.RevenueAmount != nil {
			result.RevenueAmount = entry.RevenueAmount
		}
	}

	return &result
}