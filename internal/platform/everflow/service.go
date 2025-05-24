package everflow

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/crypto"
	"github.com/affiliate-backend/internal/repository"
)

// Service represents the Everflow integration service
type Service struct {
	client              *Client
	advertiserRepo      repository.AdvertiserRepository
	providerMappingRepo repository.AdvertiserProviderMappingRepository
	campaignRepo        repository.CampaignRepository
	cryptoService       crypto.Service
}

// NewService creates a new Everflow service
func NewService(
	apiKey string,
	advertiserRepo repository.AdvertiserRepository,
	providerMappingRepo repository.AdvertiserProviderMappingRepository,
	campaignRepo repository.CampaignRepository,
	cryptoService crypto.Service,
) *Service {
	return &Service{
		client:              NewClient(apiKey),
		advertiserRepo:      advertiserRepo,
		providerMappingRepo: providerMappingRepo,
		campaignRepo:        campaignRepo,
		cryptoService:       cryptoService,
	}
}

// CreateAdvertiserInEverflow creates an advertiser in Everflow and stores the mapping
func (s *Service) CreateAdvertiserInEverflow(ctx context.Context, advertiser *domain.Advertiser) error {
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

// ListAdvertisersFromEverflow retrieves advertisers from Everflow
func (s *Service) ListAdvertisersFromEverflow(ctx context.Context, page, pageSize *int) (*EverflowListAdvertisersResponse, error) {
	opts := &ListAdvertisersOptions{
		Page:     page,
		PageSize: pageSize,
	}

	return s.client.ListAdvertisers(ctx, opts)
}

// GetAdvertiserFromEverflow retrieves a single advertiser from Everflow by ID
func (s *Service) GetAdvertiserFromEverflow(ctx context.Context, networkAdvertiserID int64, relationships []string) (*Advertiser, error) {
	opts := &GetAdvertiserOptions{
		Relationships: relationships,
	}

	return s.client.GetAdvertiser(ctx, networkAdvertiserID, opts)
}

// GetAdvertiserFromEverflowByMapping retrieves an advertiser from Everflow using our internal advertiser ID
func (s *Service) GetAdvertiserFromEverflowByMapping(ctx context.Context, advertiserID int64, relationships []string) (*Advertiser, error) {
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

	return s.GetAdvertiserFromEverflow(ctx, networkAdvertiserID, relationships)
}

// UpdateAdvertiserInEverflow updates an advertiser in Everflow by ID
func (s *Service) UpdateAdvertiserInEverflow(ctx context.Context, networkAdvertiserID int64, req EverflowUpdateAdvertiserRequest) (*Advertiser, error) {
	return s.client.UpdateAdvertiser(ctx, networkAdvertiserID, req)
}

// UpdateAdvertiserInEverflowByMapping updates an advertiser in Everflow using our internal advertiser ID
func (s *Service) UpdateAdvertiserInEverflowByMapping(ctx context.Context, advertiserID int64, req EverflowUpdateAdvertiserRequest) (*Advertiser, error) {
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

	return s.UpdateAdvertiserInEverflow(ctx, networkAdvertiserID, req)
}

// CreateOfferInEverflow creates an offer in Everflow for a campaign and stores the mapping
func (s *Service) CreateOfferInEverflow(ctx context.Context, campaign *domain.Campaign) error {
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

// mapAdvertiserToEverflowRequest maps our advertiser to an Everflow advertiser request
func (s *Service) mapAdvertiserToEverflowRequest(advertiser *domain.Advertiser) (*EverflowCreateAdvertiserRequest, error) {
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

// mapCampaignToEverflowRequest maps our campaign to an Everflow offer request
func (s *Service) mapCampaignToEverflowRequest(campaign *domain.Campaign, networkAdvertiserID int64) (*OfferInput, error) {
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

	// Create a default destination URL (this would typically come from campaign details)
	destinationURL := fmt.Sprintf("https://example.com/campaigns/%d?click_id={transaction_id}", campaign.CampaignID)

	// Default values
	currencyID := "USD"
	visibility := "public"
	conversionMethod := "server_postback"
	sessionDefinition := "cookie"
	sessionDuration := 720 // 30 days in hours

	// Create default payout/revenue structure
	payoutAmount := 1.00
	revenueAmount := 2.00
	payoutRevenue := &PayoutRevenue{
		Entries: []PayoutRevenueEntry{
			{
				EntryName:     nil, // Empty entry name for default
				PayoutType:    "cpa",
				PayoutAmount:  &payoutAmount,
				RevenueType:   "rpa",
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

	return req, nil
}

// GetOfferFromEverflow retrieves a single offer from Everflow by ID
func (s *Service) GetOfferFromEverflow(ctx context.Context, networkOfferID int64, relationships []string) (*Offer, error) {
	opts := &GetOfferOptions{
		Relationships: relationships,
	}

	return s.client.GetOffer(ctx, networkOfferID, opts)
}

// GetOfferFromEverflowByMapping retrieves an offer from Everflow using our internal campaign ID
func (s *Service) GetOfferFromEverflowByMapping(ctx context.Context, campaignID int64, relationships []string) (*Offer, error) {
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

	return s.GetOfferFromEverflow(ctx, networkOfferID, relationships)
}

// UpdateOfferInEverflow updates an offer in Everflow by ID
func (s *Service) UpdateOfferInEverflow(ctx context.Context, networkOfferID int64, req OfferInput) (*Offer, error) {
	return s.client.UpdateOffer(ctx, networkOfferID, req)
}

// UpdateOfferInEverflowByMapping updates an offer in Everflow using our internal campaign ID
func (s *Service) UpdateOfferInEverflowByMapping(ctx context.Context, campaignID int64, req OfferInput) (*Offer, error) {
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

	// Update the offer
	updatedOffer, err := s.UpdateOfferInEverflow(ctx, networkOfferID, req)
	if err != nil {
		return nil, err
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

	return updatedOffer, nil
}

// ListOffersFromEverflow retrieves offers from Everflow with optional filtering
func (s *Service) ListOffersFromEverflow(ctx context.Context, req OffersTableRequest, opts *OffersTableOptions) (*OffersTableResponse, error) {
	return s.client.OffersTable(ctx, req, opts)
}

// SyncCampaignWithEverflowOffer synchronizes a campaign's data with its Everflow offer
func (s *Service) SyncCampaignWithEverflowOffer(ctx context.Context, campaignID int64) error {
	// Get the campaign
	campaign, err := s.campaignRepo.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return fmt.Errorf("failed to get campaign: %w", err)
	}

	// Get the Everflow offer
	offer, err := s.GetOfferFromEverflowByMapping(ctx, campaignID, []string{"reporting"})
	if err != nil {
		return fmt.Errorf("failed to get Everflow offer: %w", err)
	}

	// Update campaign status based on offer status
	var newStatus string
	switch offer.OfferStatus {
	case "active":
		newStatus = "active"
	case "paused":
		newStatus = "paused"
	case "pending":
		newStatus = "draft"
	default:
		newStatus = campaign.Status // Keep current status if unknown
	}

	if campaign.Status != newStatus {
		campaign.Status = newStatus
		if err := s.campaignRepo.UpdateCampaign(ctx, campaign); err != nil {
			return fmt.Errorf("failed to update campaign status: %w", err)
		}
	}

	// Update the provider offer mapping
	offers, err := s.campaignRepo.ListCampaignProviderOffersByCampaign(ctx, campaignID)
	if err != nil {
		return fmt.Errorf("failed to get campaign provider offers: %w", err)
	}

	for _, providerOffer := range offers {
		if providerOffer.ProviderType == "everflow" {
			now := time.Now()
			providerOffer.IsActiveOnProvider = offer.OfferStatus == "active"
			providerOffer.LastSyncedAt = &now

			// Update provider config
			providerConfig := map[string]interface{}{
				"network_offer_id":      offer.NetworkOfferID,
				"network_id":            offer.NetworkID,
				"network_advertiser_id": offer.NetworkAdvertiserID,
				"offer_status":          offer.OfferStatus,
				"currency_id":           offer.CurrencyID,
				"encoded_value":         offer.EncodedValue,
				"time_created":          offer.TimeCreated,
				"time_saved":            offer.TimeSaved,
				"today_clicks":          offer.TodayClicks,
				"today_revenue":         offer.TodayRevenue,
			}

			providerConfigJSON, err := json.Marshal(providerConfig)
			if err != nil {
				return fmt.Errorf("failed to marshal provider config: %w", err)
			}

			providerConfigStr := string(providerConfigJSON)
			providerOffer.ProviderOfferConfig = &providerConfigStr

			if err := s.campaignRepo.UpdateCampaignProviderOffer(ctx, providerOffer); err != nil {
				return fmt.Errorf("failed to update campaign provider offer: %w", err)
			}
			break
		}
	}

	return nil
}