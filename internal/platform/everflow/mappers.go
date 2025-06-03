package everflow

import (
	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/everflow/advertiser"
	"github.com/affiliate-backend/internal/platform/everflow/affiliate"
	"github.com/affiliate-backend/internal/platform/everflow/offer"
)

// Advertiser mapping functions

func (s *IntegrationService) mapAdvertiserToEverflowRequest(adv *domain.Advertiser) (*advertiser.CreateAdvertiserRequest, error) {
	req := advertiser.NewCreateAdvertiserRequest(
		adv.Name,
		s.mapDomainStatusToEverflowStatus(adv.Status),
		s.getDefaultNetworkEmployeeID(adv.NetworkEmployeeID),
		s.getDefaultCurrencyID(adv.DefaultCurrencyID),
		s.getDefaultReportingTimezoneID(adv.ReportingTimezoneID),
		s.getDefaultAttributionMethod(adv.AttributionMethod),
		s.getDefaultEmailAttributionMethod(adv.EmailAttributionMethod),
		s.getDefaultAttributionPriority(adv.AttributionPriority),
	)

	// Set optional fields
	if adv.InternalNotes != nil {
		req.SetInternalNotes(*adv.InternalNotes)
	}
	if adv.PlatformName != nil {
		req.SetPlatformName(*adv.PlatformName)
	}
	if adv.PlatformURL != nil {
		req.SetPlatformUrl(*adv.PlatformURL)
	}
	if adv.PlatformUsername != nil {
		req.SetPlatformUsername(*adv.PlatformUsername)
	}
	if adv.AccountingContactEmail != nil {
		req.SetAccountingContactEmail(*adv.AccountingContactEmail)
	}
	if adv.OfferIDMacro != nil {
		req.SetOfferIdMacro(*adv.OfferIDMacro)
	}
	if adv.AffiliateIDMacro != nil {
		req.SetAffiliateIdMacro(*adv.AffiliateIDMacro)
	}

	return req, nil
}

func (s *IntegrationService) mapAdvertiserToEverflowUpdateRequest(adv *domain.Advertiser) (*advertiser.UpdateAdvertiserRequest, error) {
	req := advertiser.NewUpdateAdvertiserRequestWithDefaults()
	
	req.SetName(adv.Name)
	req.SetAccountStatus(s.mapDomainStatusToEverflowStatus(adv.Status))
	
	if adv.NetworkEmployeeID != nil {
		req.SetNetworkEmployeeId(*adv.NetworkEmployeeID)
	}
	if adv.DefaultCurrencyID != nil {
		req.SetDefaultCurrencyId(*adv.DefaultCurrencyID)
	}
	if adv.ReportingTimezoneID != nil {
		req.SetReportingTimezoneId(*adv.ReportingTimezoneID)
	}
	if adv.AttributionMethod != nil {
		req.SetAttributionMethod(*adv.AttributionMethod)
	}
	if adv.EmailAttributionMethod != nil {
		req.SetEmailAttributionMethod(*adv.EmailAttributionMethod)
	}
	if adv.AttributionPriority != nil {
		req.SetAttributionPriority(*adv.AttributionPriority)
	}

	// Set optional fields
	if adv.InternalNotes != nil {
		req.SetInternalNotes(*adv.InternalNotes)
	}
	if adv.PlatformName != nil {
		req.SetPlatformName(*adv.PlatformName)
	}
	if adv.PlatformURL != nil {
		req.SetPlatformUrl(*adv.PlatformURL)
	}
	if adv.PlatformUsername != nil {
		req.SetPlatformUsername(*adv.PlatformUsername)
	}
	if adv.AccountingContactEmail != nil {
		req.SetAccountingContactEmail(*adv.AccountingContactEmail)
	}
	if adv.OfferIDMacro != nil {
		req.SetOfferIdMacro(*adv.OfferIDMacro)
	}
	if adv.AffiliateIDMacro != nil {
		req.SetAffiliateIdMacro(*adv.AffiliateIDMacro)
	}

	return req, nil
}

func (s *IntegrationService) mapEverflowResponseToAdvertiser(resp *advertiser.Advertiser, adv *domain.Advertiser) domain.Advertiser {
	result := *adv // Copy the original

	// Update fields from Everflow response
	if resp.Name != nil {
		result.Name = *resp.Name
	}
	if resp.AccountStatus != nil {
		result.Status = s.mapEverflowStatusToDomainStatus(*resp.AccountStatus)
	}
	if resp.InternalNotes != nil {
		result.InternalNotes = resp.InternalNotes
	}
	if resp.DefaultCurrencyId != nil {
		result.DefaultCurrencyID = resp.DefaultCurrencyId
	}
	if resp.PlatformName != nil {
		result.PlatformName = resp.PlatformName
	}
	if resp.PlatformUrl != nil {
		result.PlatformURL = resp.PlatformUrl
	}
	if resp.PlatformUsername != nil {
		result.PlatformUsername = resp.PlatformUsername
	}
	if resp.ReportingTimezoneId != nil {
		result.ReportingTimezoneID = resp.ReportingTimezoneId
	}
	if resp.AttributionMethod != nil {
		result.AttributionMethod = resp.AttributionMethod
	}
	if resp.EmailAttributionMethod != nil {
		result.EmailAttributionMethod = resp.EmailAttributionMethod
	}
	if resp.AttributionPriority != nil {
		result.AttributionPriority = resp.AttributionPriority
	}
	if resp.AccountingContactEmail != nil {
		result.AccountingContactEmail = resp.AccountingContactEmail
	}
	if resp.OfferIdMacro != nil {
		result.OfferIDMacro = resp.OfferIdMacro
	}
	if resp.AffiliateIdMacro != nil {
		result.AffiliateIDMacro = resp.AffiliateIdMacro
	}
	if resp.NetworkEmployeeId != nil {
		result.NetworkEmployeeID = resp.NetworkEmployeeId
	}
	if resp.IsExposePublisherReportingData.IsSet() {
		val := resp.IsExposePublisherReportingData.Get()
		result.IsExposePublisherReporting = val
	}

	return result
}

// Affiliate mapping functions

func (s *IntegrationService) mapAffiliateToEverflowRequest(aff *domain.Affiliate) (*affiliate.CreateAffiliateRequest, error) {
	req := affiliate.NewCreateAffiliateRequest(
		aff.Name,
		s.mapDomainStatusToEverflowStatus(aff.Status),
		s.getDefaultNetworkEmployeeID(aff.NetworkEmployeeID),
	)

	// Set optional fields
	if aff.InternalNotes != nil {
		req.SetInternalNotes(*aff.InternalNotes)
	}
	if aff.DefaultCurrencyID != nil {
		req.SetDefaultCurrencyId(*aff.DefaultCurrencyID)
	}
	if aff.EnableMediaCostTrackingLinks != nil {
		req.SetEnableMediaCostTrackingLinks(*aff.EnableMediaCostTrackingLinks)
	}
	if aff.ReferrerID != nil {
		req.SetReferrerId(*aff.ReferrerID)
	}
	if aff.IsContactAddressEnabled != nil {
		req.SetIsContactAddressEnabled(*aff.IsContactAddressEnabled)
	}
	if aff.NetworkAffiliateTierID != nil {
		req.SetNetworkAffiliateTierId(*aff.NetworkAffiliateTierID)
	}

	return req, nil
}

func (s *IntegrationService) mapAffiliateToEverflowUpdateRequest(aff *domain.Affiliate) (*affiliate.UpdateAffiliateRequest, error) {
	req := affiliate.NewUpdateAffiliateRequestWithDefaults()
	
	req.SetName(aff.Name)
	req.SetAccountStatus(s.mapDomainStatusToEverflowStatus(aff.Status))
	
	if aff.NetworkEmployeeID != nil {
		req.SetNetworkEmployeeId(*aff.NetworkEmployeeID)
	}

	// Set optional fields
	if aff.InternalNotes != nil {
		req.SetInternalNotes(*aff.InternalNotes)
	}
	if aff.DefaultCurrencyID != nil {
		req.SetDefaultCurrencyId(*aff.DefaultCurrencyID)
	}
	if aff.EnableMediaCostTrackingLinks != nil {
		req.SetEnableMediaCostTrackingLinks(*aff.EnableMediaCostTrackingLinks)
	}
	if aff.ReferrerID != nil {
		req.SetReferrerId(*aff.ReferrerID)
	}
	if aff.IsContactAddressEnabled != nil {
		req.SetIsContactAddressEnabled(*aff.IsContactAddressEnabled)
	}
	if aff.NetworkAffiliateTierID != nil {
		req.SetNetworkAffiliateTierId(*aff.NetworkAffiliateTierID)
	}

	return req, nil
}

func (s *IntegrationService) mapEverflowCreateResponseToAffiliate(resp *affiliate.Affiliate, aff *domain.Affiliate) domain.Affiliate {
	result := *aff // Copy the original

	// Update fields from Everflow response
	if resp.Name != nil {
		result.Name = *resp.Name
	}
	if resp.NetworkAffiliateId != nil {
		result.NetworkAffiliateID = resp.NetworkAffiliateId
	}

	return result
}

func (s *IntegrationService) mapEverflowResponseToAffiliate(resp *affiliate.AffiliateWithRelationships, aff *domain.Affiliate) domain.Affiliate {
	result := *aff // Copy the original

	// Update fields from Everflow response
	if resp.Name != nil {
		result.Name = *resp.Name
	}
	if resp.AccountStatus != nil {
		result.Status = s.mapEverflowStatusToDomainStatus(*resp.AccountStatus)
	}
	if resp.InternalNotes != nil {
		result.InternalNotes = resp.InternalNotes
	}
	if resp.DefaultCurrencyId != nil {
		result.DefaultCurrencyID = resp.DefaultCurrencyId
	}
	if resp.EnableMediaCostTrackingLinks != nil {
		result.EnableMediaCostTrackingLinks = resp.EnableMediaCostTrackingLinks
	}
	if resp.ReferrerId != nil {
		result.ReferrerID = resp.ReferrerId
	}
	if resp.IsContactAddressEnabled != nil {
		result.IsContactAddressEnabled = resp.IsContactAddressEnabled
	}
	// NetworkAffiliateTierID is not available in AffiliateWithRelationships response
	if resp.NetworkEmployeeId != nil {
		result.NetworkEmployeeID = resp.NetworkEmployeeId
	}

	return result
}

// Campaign (Offer) mapping functions

func (s *IntegrationService) mapCampaignToEverflowRequest(camp *domain.Campaign, networkAdvertiserID int32) (*offer.CreateOfferRequest, error) {
	// Create default payout revenue
	defaultPayout := offer.NewPayoutRevenueWithDefaults()
	defaultPayout.SetPayoutType("cpa")
	defaultPayout.SetPayoutAmount(0.0)
	
	req := offer.NewCreateOfferRequest(
		networkAdvertiserID,
		camp.Name,
		s.getDefaultDestinationURL(camp.DestinationURL),
		s.getDefaultOfferStatus(camp.Status),
		[]offer.PayoutRevenue{*defaultPayout},
	)

	// Set optional fields
	if camp.ThumbnailURL != nil {
		req.SetThumbnailUrl(*camp.ThumbnailURL)
	}
	if camp.InternalNotes != nil {
		req.SetInternalNotes(*camp.InternalNotes)
	}
	if camp.ServerSideURL != nil {
		req.SetServerSideUrl(*camp.ServerSideURL)
	}
	if camp.IsViewThroughEnabled != nil {
		req.SetIsViewThroughEnabled(*camp.IsViewThroughEnabled)
	}
	if camp.ViewThroughDestinationURL != nil {
		req.SetViewThroughDestinationUrl(*camp.ViewThroughDestinationURL)
	}
	if camp.PreviewURL != nil {
		req.SetPreviewUrl(*camp.PreviewURL)
	}
	if camp.CurrencyID != nil {
		req.SetCurrencyId(*camp.CurrencyID)
	}
	if camp.CapsTimezoneID != nil {
		req.SetCapsTimezoneId(*camp.CapsTimezoneID)
	}
	if camp.SessionDuration != nil {
		req.SetSessionDuration(*camp.SessionDuration)
	}

	return req, nil
}

func (s *IntegrationService) mapCampaignToEverflowUpdateRequest(camp *domain.Campaign) (*offer.UpdateOfferRequest, error) {
	// Create default payout revenue
	defaultPayout := offer.NewPayoutRevenueWithDefaults()
	defaultPayout.SetPayoutType("cpa")
	defaultPayout.SetPayoutAmount(0.0)
	
	req := offer.NewUpdateOfferRequestWithDefaults()
	
	// Set required fields
	req.SetName(camp.Name)
	req.SetOfferStatus(s.getDefaultOfferStatus(camp.Status))
	req.SetDestinationUrl(s.getDefaultDestinationURL(camp.DestinationURL))
	if camp.NetworkAdvertiserID != nil {
		req.SetNetworkAdvertiserId(*camp.NetworkAdvertiserID)
	}
	req.SetPayoutRevenue([]offer.PayoutRevenue{*defaultPayout})
	
	if camp.DestinationURL != nil {
		req.SetDestinationUrl(*camp.DestinationURL)
	}

	// Set optional fields
	if camp.ThumbnailURL != nil {
		req.SetThumbnailUrl(*camp.ThumbnailURL)
	}
	if camp.InternalNotes != nil {
		req.SetInternalNotes(*camp.InternalNotes)
	}
	if camp.ServerSideURL != nil {
		req.SetServerSideUrl(*camp.ServerSideURL)
	}
	if camp.IsViewThroughEnabled != nil {
		req.SetIsViewThroughEnabled(*camp.IsViewThroughEnabled)
	}
	if camp.ViewThroughDestinationURL != nil {
		req.SetViewThroughDestinationUrl(*camp.ViewThroughDestinationURL)
	}
	if camp.PreviewURL != nil {
		req.SetPreviewUrl(*camp.PreviewURL)
	}
	if camp.CurrencyID != nil {
		req.SetCurrencyId(*camp.CurrencyID)
	}
	if camp.CapsTimezoneID != nil {
		req.SetCapsTimezoneId(*camp.CapsTimezoneID)
	}
	if camp.SessionDuration != nil {
		req.SetSessionDuration(*camp.SessionDuration)
	}

	return req, nil
}

func (s *IntegrationService) mapEverflowResponseToCampaign(resp *offer.OfferResponse, camp *domain.Campaign) domain.Campaign {
	result := *camp // Copy the original

	// Update fields from Everflow response
	if resp.Name != nil {
		result.Name = *resp.Name
	}
	if resp.OfferStatus != nil {
		result.Status = s.mapEverflowOfferStatusToDomainStatus(*resp.OfferStatus)
	}
	if resp.DestinationUrl != nil {
		result.DestinationURL = resp.DestinationUrl
	}
	if resp.ThumbnailUrl != nil {
		result.ThumbnailURL = resp.ThumbnailUrl
	}
	if resp.InternalNotes != nil {
		result.InternalNotes = resp.InternalNotes
	}
	if resp.ServerSideUrl != nil {
		result.ServerSideURL = resp.ServerSideUrl
	}
	if resp.IsViewThroughEnabled != nil {
		result.IsViewThroughEnabled = resp.IsViewThroughEnabled
	}
	if resp.ViewThroughDestinationUrl != nil {
		result.ViewThroughDestinationURL = resp.ViewThroughDestinationUrl
	}
	if resp.PreviewUrl != nil {
		result.PreviewURL = resp.PreviewUrl
	}
	if resp.CurrencyId != nil {
		result.CurrencyID = resp.CurrencyId
	}
	if resp.CapsTimezoneId != nil {
		result.CapsTimezoneID = resp.CapsTimezoneId
	}
	if resp.SessionDuration != nil {
		result.SessionDuration = resp.SessionDuration
	}
	if resp.NetworkAdvertiserId != nil {
		result.NetworkAdvertiserID = resp.NetworkAdvertiserId
	}

	return result
}

// Helper functions for default values and status mapping

func (s *IntegrationService) mapDomainStatusToEverflowStatus(status string) string {
	switch status {
	case "active":
		return "active"
	case "pending":
		return "pending"
	case "inactive":
		return "inactive"
	case "rejected":
		return "rejected"
	default:
		return "pending"
	}
}

func (s *IntegrationService) mapEverflowStatusToDomainStatus(status string) string {
	switch status {
	case "active":
		return "active"
	case "pending":
		return "pending"
	case "inactive":
		return "inactive"
	case "rejected":
		return "rejected"
	default:
		return "pending"
	}
}

func (s *IntegrationService) mapEverflowOfferStatusToDomainStatus(status string) string {
	switch status {
	case "active":
		return "active"
	case "paused":
		return "paused"
	case "pending":
		return "draft"
	case "archived":
		return "archived"
	default:
		return "draft"
	}
}

func (s *IntegrationService) getDefaultNetworkEmployeeID(id *int32) int32 {
	if id != nil {
		return *id
	}
	return 1 // Default employee ID
}

func (s *IntegrationService) getDefaultCurrencyID(id *string) string {
	if id != nil {
		return *id
	}
	return "USD" // Default currency
}

func (s *IntegrationService) getDefaultReportingTimezoneID(id *int32) int32 {
	if id != nil {
		return *id
	}
	return 1 // Default timezone ID (UTC)
}

func (s *IntegrationService) getDefaultAttributionMethod(method *string) string {
	if method != nil {
		return *method
	}
	return "first_click" // Default attribution method
}

func (s *IntegrationService) getDefaultEmailAttributionMethod(method *string) string {
	if method != nil {
		return *method
	}
	return "first_click" // Default email attribution method
}

func (s *IntegrationService) getDefaultAttributionPriority(priority *string) string {
	if priority != nil {
		return *priority
	}
	return "click" // Default attribution priority
}

func (s *IntegrationService) getDefaultDestinationURL(url *string) string {
	if url != nil {
		return *url
	}
	return "https://example.com" // Default destination URL
}

func (s *IntegrationService) getDefaultOfferStatus(status string) string {
	switch status {
	case "active":
		return "active"
	case "paused":
		return "paused"
	case "draft":
		return "pending"
	case "archived":
		return "archived"
	default:
		return "pending"
	}
}