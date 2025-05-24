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
		"offer_url":             everflowResp.OfferURL,
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
func (s *Service) mapCampaignToEverflowRequest(campaign *domain.Campaign, networkAdvertiserID int64) (*EverflowCreateOfferRequest, error) {
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

	// Create basic request
	req := &EverflowCreateOfferRequest{
		Name:                campaign.Name,
		NetworkAdvertiserID: networkAdvertiserID,
		DestinationURL:      destinationURL,
		OfferStatus:         offerStatus,
		CurrencyID:          "USD",             // Default to USD, could be configurable
		Visibility:          "public",          // Default to public, could be configurable
		ConversionMethod:    "server_postback", // Default to server postback, could be configurable

		// Add default payout/revenue structure
		PayoutRevenue: []PayoutRevenueItem{
			{
				IsDefault:     true,
				PayoutType:    "cpa",
				PayoutAmount:  1.00, // Default payout amount, should be configurable
				RevenueType:   "cpa",
				RevenueAmount: 2.00, // Default revenue amount, should be configurable
			},
		},
	}

	// Add description if available
	if campaign.Description != nil {
		req.Description = campaign.Description
	}

	// Add session definition and duration
	sessionDefinition := "cookie"
	sessionDuration := 720 // 30 days in hours
	req.SessionDefinition = &sessionDefinition
	req.SessionDuration = &sessionDuration

	// Add tags
	req.Tags = []string{
		fmt.Sprintf("campaign_id:%d", campaign.CampaignID),
		fmt.Sprintf("advertiser_id:%d", campaign.AdvertiserID),
		fmt.Sprintf("organization_id:%d", campaign.OrganizationID),
	}

	return req, nil
}