package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/logger"
	"github.com/affiliate-backend/internal/platform/provider"
	"github.com/affiliate-backend/internal/repository"
)

// TrackingLinkService defines the interface for tracking link business logic
type TrackingLinkService interface {
	// Core CRUD operations
	CreateTrackingLink(ctx context.Context, trackingLink *domain.TrackingLink) error
	GetTrackingLinkByID(ctx context.Context, id int64) (*domain.TrackingLink, error)
	UpdateTrackingLink(ctx context.Context, trackingLink *domain.TrackingLink) error
	DeleteTrackingLink(ctx context.Context, id int64) error

	// List operations
	ListTrackingLinksByCampaign(ctx context.Context, campaignID int64, limit, offset int) ([]*domain.TrackingLink, error)
	ListTrackingLinksByAffiliate(ctx context.Context, affiliateID int64, limit, offset int) ([]*domain.TrackingLink, error)
	ListTrackingLinksByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.TrackingLink, error)
	ListTrackingLinksByCampaignAndAffiliate(ctx context.Context, campaignID, affiliateID int64, limit, offset int) ([]*domain.TrackingLink, error)

	// Tracking link generation
	GenerateTrackingLink(ctx context.Context, req *domain.TrackingLinkGenerationRequest) (*domain.TrackingLinkGenerationResponse, error)
	RegenerateTrackingLink(ctx context.Context, trackingLinkID int64) (*domain.TrackingLinkGenerationResponse, error)
	
	// Tracking link upsert
	UpsertTrackingLink(ctx context.Context, req *domain.TrackingLinkUpsertRequest) (*domain.TrackingLinkUpsertResponse, error)

	// Provider sync operations
	SyncTrackingLinkToProvider(ctx context.Context, trackingLinkID int64) error
	SyncTrackingLinkFromProvider(ctx context.Context, trackingLinkID int64) error

	// Provider mapping operations
	CreateTrackingLinkProviderMapping(ctx context.Context, mapping *domain.TrackingLinkProviderMapping) (*domain.TrackingLinkProviderMapping, error)
	GetTrackingLinkProviderMapping(ctx context.Context, trackingLinkID int64, providerType string) (*domain.TrackingLinkProviderMapping, error)
	UpdateTrackingLinkProviderMapping(ctx context.Context, mapping *domain.TrackingLinkProviderMapping) error
	DeleteTrackingLinkProviderMapping(ctx context.Context, mappingID int64) error
}

// trackingLinkService implements TrackingLinkService
type trackingLinkService struct {
	trackingLinkRepo         repository.TrackingLinkRepository
	trackingLinkProviderRepo repository.TrackingLinkProviderMappingRepository
	campaignRepo             repository.CampaignRepository
	affiliateRepo            repository.AffiliateRepository
	campaignProviderRepo     repository.CampaignProviderMappingRepository
	affiliateProviderRepo    repository.AffiliateProviderMappingRepository
	integrationService       provider.IntegrationService
	orgAssociationService    OrganizationAssociationService
}

// NewTrackingLinkService creates a new tracking link service
func NewTrackingLinkService(
	trackingLinkRepo repository.TrackingLinkRepository,
	trackingLinkProviderRepo repository.TrackingLinkProviderMappingRepository,
	campaignRepo repository.CampaignRepository,
	affiliateRepo repository.AffiliateRepository,
	campaignProviderRepo repository.CampaignProviderMappingRepository,
	affiliateProviderRepo repository.AffiliateProviderMappingRepository,
	integrationService provider.IntegrationService,
	orgAssociationService OrganizationAssociationService,
) TrackingLinkService {
	return &trackingLinkService{
		trackingLinkRepo:         trackingLinkRepo,
		trackingLinkProviderRepo: trackingLinkProviderRepo,
		campaignRepo:             campaignRepo,
		affiliateRepo:            affiliateRepo,
		campaignProviderRepo:     campaignProviderRepo,
		affiliateProviderRepo:    affiliateProviderRepo,
		integrationService:       integrationService,
		orgAssociationService:    orgAssociationService,
	}
}

// CreateTrackingLink creates a new tracking link
func (s *trackingLinkService) CreateTrackingLink(ctx context.Context, trackingLink *domain.TrackingLink) error {
	// Validate tracking link data
	if err := s.validateTrackingLink(trackingLink); err != nil {
		return fmt.Errorf("tracking link validation failed: %w", err)
	}

	// Set default values
	if trackingLink.Status == "" {
		trackingLink.Status = "active"
	}

	// Set timestamps
	now := time.Now()
	trackingLink.CreatedAt = now
	trackingLink.UpdatedAt = now

	// Create tracking link in repository
	if err := s.trackingLinkRepo.CreateTrackingLink(ctx, trackingLink); err != nil {
		return fmt.Errorf("failed to create tracking link: %w", err)
	}

	// Synchronize with provider (Everflow) - run in background to avoid blocking
	go func() {
		if err := s.synchronizeTrackingLinkWithProvider(context.Background(), trackingLink); err != nil {
			logger.Warn("Failed to synchronize tracking link with provider", 
				"tracking_link_id", trackingLink.TrackingLinkID, 
				"error", err)
		}
	}()

	return nil
}

// GetTrackingLinkByID retrieves a tracking link by its ID
func (s *trackingLinkService) GetTrackingLinkByID(ctx context.Context, id int64) (*domain.TrackingLink, error) {
	trackingLink, err := s.trackingLinkRepo.GetTrackingLinkByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracking link: %w", err)
	}

	return trackingLink, nil
}

// UpdateTrackingLink updates an existing tracking link
func (s *trackingLinkService) UpdateTrackingLink(ctx context.Context, trackingLink *domain.TrackingLink) error {
	// Validate tracking link data
	if err := s.validateTrackingLink(trackingLink); err != nil {
		return fmt.Errorf("tracking link validation failed: %w", err)
	}

	// Get the existing tracking link to compare parameters
	existingLink, err := s.trackingLinkRepo.GetTrackingLinkByID(ctx, trackingLink.TrackingLinkID)
	if err != nil {
		return fmt.Errorf("failed to get existing tracking link: %w", err)
	}

	// Check if tracking parameters have changed
	parametersChanged := s.hasTrackingParametersChanged(existingLink, trackingLink)
	
	// Update tracking link in repository
	if err := s.trackingLinkRepo.UpdateTrackingLink(ctx, trackingLink); err != nil {
		return fmt.Errorf("failed to update tracking link: %w", err)
	}

	// If tracking parameters changed, regenerate the tracking link
	if parametersChanged {
		logger.Info("Tracking parameters changed, regenerating tracking link", 
			"tracking_link_id", trackingLink.TrackingLinkID,
			"campaign_id", trackingLink.CampaignID,
			"affiliate_id", trackingLink.AffiliateID)
		
		response, err := s.RegenerateTrackingLink(ctx, trackingLink.TrackingLinkID)
		if err != nil {
			logger.Error("Failed to regenerate tracking link after parameter change", 
				"tracking_link_id", trackingLink.TrackingLinkID,
				"error", err)
			// Don't fail the update operation, just log the error
			// The tracking link update was successful, regeneration is a bonus
		} else {
			logger.Info("Successfully regenerated tracking link after parameter change", 
				"tracking_link_id", trackingLink.TrackingLinkID,
				"new_url", response.GeneratedURL)
		}
	}

	return nil
}

// DeleteTrackingLink deletes a tracking link by its ID
func (s *trackingLinkService) DeleteTrackingLink(ctx context.Context, id int64) error {
	if err := s.trackingLinkRepo.DeleteTrackingLink(ctx, id); err != nil {
		return fmt.Errorf("failed to delete tracking link: %w", err)
	}

	return nil
}

// ListTrackingLinksByCampaign retrieves tracking links for a specific campaign
func (s *trackingLinkService) ListTrackingLinksByCampaign(ctx context.Context, campaignID int64, limit, offset int) ([]*domain.TrackingLink, error) {
	trackingLinks, err := s.trackingLinkRepo.ListTrackingLinksByCampaign(ctx, campaignID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tracking links by campaign: %w", err)
	}

	return trackingLinks, nil
}

// ListTrackingLinksByAffiliate retrieves tracking links for a specific affiliate
func (s *trackingLinkService) ListTrackingLinksByAffiliate(ctx context.Context, affiliateID int64, limit, offset int) ([]*domain.TrackingLink, error) {
	trackingLinks, err := s.trackingLinkRepo.ListTrackingLinksByAffiliate(ctx, affiliateID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tracking links by affiliate: %w", err)
	}

	return trackingLinks, nil
}

// ListTrackingLinksByOrganization retrieves tracking links for a specific organization
func (s *trackingLinkService) ListTrackingLinksByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.TrackingLink, error) {
	trackingLinks, err := s.trackingLinkRepo.ListTrackingLinksByOrganization(ctx, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tracking links by organization: %w", err)
	}

	return trackingLinks, nil
}

// ListTrackingLinksByCampaignAndAffiliate lists tracking links by campaign and affiliate
func (s *trackingLinkService) ListTrackingLinksByCampaignAndAffiliate(ctx context.Context, campaignID, affiliateID int64, limit, offset int) ([]*domain.TrackingLink, error) {
	// Validate campaign exists
	campaign, err := s.campaignRepo.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign: %w", err)
	}

	// Validate affiliate exists
	affiliate, err := s.affiliateRepo.GetAffiliateByID(ctx, affiliateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get affiliate: %w", err)
	}

	// Verify that there's an active association between the advertiser and affiliate organizations
	// and that the campaign and affiliate are visible to each other
	err = s.verifyAssociationAndVisibility(ctx, campaign.OrganizationID, affiliate.OrganizationID, campaignID, affiliateID)
	if err != nil {
		return nil, fmt.Errorf("association verification failed: %w", err)
	}

	trackingLinks, err := s.trackingLinkRepo.ListTrackingLinksByCampaignAndAffiliate(ctx, campaignID, affiliateID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tracking links by campaign and affiliate: %w", err)
	}

	return trackingLinks, nil
}

// GenerateTrackingLink generates a new tracking link with provider integration
func (s *trackingLinkService) GenerateTrackingLink(ctx context.Context, req *domain.TrackingLinkGenerationRequest) (*domain.TrackingLinkGenerationResponse, error) {
	// Validate the request
	if err := s.validateTrackingLinkGenerationRequest(req); err != nil {
		return nil, fmt.Errorf("tracking link generation request validation failed: %w", err)
	}

	// Get campaign and affiliate information
	campaign, err := s.campaignRepo.GetCampaignByID(ctx, req.CampaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign: %w", err)
	}

	// Validate affiliate exists
	affiliate, err := s.affiliateRepo.GetAffiliateByID(ctx, req.AffiliateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get affiliate: %w", err)
	}

	// Verify that there's an active association between the advertiser and affiliate organizations
	// and that the campaign and affiliate are visible to each other
	err = s.verifyAssociationAndVisibility(ctx, campaign.OrganizationID, affiliate.OrganizationID, req.CampaignID, req.AffiliateID)
	if err != nil {
		return nil, fmt.Errorf("association verification failed: %w", err)
	}

	// Check if tracking link already exists with same parameters
	existingLink, err := s.trackingLinkRepo.GetTrackingLinkByCampaignAndAffiliate(
		ctx, req.CampaignID, req.AffiliateID, req.SourceID, req.Sub1, req.Sub2, req.Sub3, req.Sub4, req.Sub5)
	if err == nil && existingLink.TrackingURL != nil {
		// Tracking link already exists, return it
		return &domain.TrackingLinkGenerationResponse{
			TrackingLink: existingLink,
			GeneratedURL: *existingLink.TrackingURL,
		}, nil
	}

	// Create new tracking link entity
	trackingLink := &domain.TrackingLink{
		OrganizationID:      campaign.OrganizationID,
		CampaignID:          req.CampaignID,
		AffiliateID:         req.AffiliateID,
		Name:                req.Name,
		Description:         req.Description,
		Status:              "active",
		SourceID:            req.SourceID,
		Sub1:                req.Sub1,
		Sub2:                req.Sub2,
		Sub3:                req.Sub3,
		Sub4:                req.Sub4,
		Sub5:                req.Sub5,
		IsEncryptParameters: req.IsEncryptParameters,
		IsRedirectLink:      req.IsRedirectLink,
	}

	// Set timestamps
	now := time.Now()
	trackingLink.CreatedAt = now
	trackingLink.UpdatedAt = now

	// Create tracking link in database
	if err := s.trackingLinkRepo.CreateTrackingLink(ctx, trackingLink); err != nil {
		return nil, fmt.Errorf("failed to create tracking link: %w", err)
	}

	// Generate tracking link via provider integration
	generatedURL, providerData, err := s.generateTrackingLinkViaProvider(ctx, trackingLink, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tracking link via provider: %w", err)
	}

	// Update tracking link with generated URL
	trackingLink.TrackingURL = &generatedURL
	if err := s.trackingLinkRepo.UpdateTrackingLink(ctx, trackingLink); err != nil {
		return nil, fmt.Errorf("failed to update tracking link with generated URL: %w", err)
	}

	return &domain.TrackingLinkGenerationResponse{
		TrackingLink: trackingLink,
		GeneratedURL: generatedURL,
		ProviderData: providerData,
	}, nil
}

// RegenerateTrackingLink regenerates an existing tracking link
func (s *trackingLinkService) RegenerateTrackingLink(ctx context.Context, trackingLinkID int64) (*domain.TrackingLinkGenerationResponse, error) {
	// Get existing tracking link
	trackingLink, err := s.trackingLinkRepo.GetTrackingLinkByID(ctx, trackingLinkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracking link: %w", err)
	}

	// Create generation request from existing tracking link
	req := &domain.TrackingLinkGenerationRequest{
		CampaignID:          trackingLink.CampaignID,
		AffiliateID:         trackingLink.AffiliateID,
		Name:                trackingLink.Name,
		Description:         trackingLink.Description,
		SourceID:            trackingLink.SourceID,
		Sub1:                trackingLink.Sub1,
		Sub2:                trackingLink.Sub2,
		Sub3:                trackingLink.Sub3,
		Sub4:                trackingLink.Sub4,
		Sub5:                trackingLink.Sub5,
		IsEncryptParameters: trackingLink.IsEncryptParameters,
		IsRedirectLink:      trackingLink.IsRedirectLink,
	}

	// Generate new tracking link via provider integration
	generatedURL, providerData, err := s.generateTrackingLinkViaProvider(ctx, trackingLink, req)
	if err != nil {
		return nil, fmt.Errorf("failed to regenerate tracking link via provider: %w", err)
	}

	// Update tracking link with new generated URL
	trackingLink.TrackingURL = &generatedURL
	if err := s.trackingLinkRepo.UpdateTrackingLink(ctx, trackingLink); err != nil {
		return nil, fmt.Errorf("failed to update tracking link with regenerated URL: %w", err)
	}

	return &domain.TrackingLinkGenerationResponse{
		TrackingLink: trackingLink,
		GeneratedURL: generatedURL,
		ProviderData: providerData,
	}, nil
}

// UpsertTrackingLink creates or updates a tracking link based on campaign_id and affiliate_id
func (s *trackingLinkService) UpsertTrackingLink(ctx context.Context, req *domain.TrackingLinkUpsertRequest) (*domain.TrackingLinkUpsertResponse, error) {
	// Validate the request
	if err := s.validateTrackingLinkUpsertRequest(req); err != nil {
		return nil, fmt.Errorf("tracking link upsert request validation failed: %w", err)
	}

	// Get campaign and affiliate information
	campaign, err := s.campaignRepo.GetCampaignByID(ctx, req.CampaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign: %w", err)
	}

	// Validate affiliate exists
	affiliate, err := s.affiliateRepo.GetAffiliateByID(ctx, req.AffiliateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get affiliate: %w", err)
	}

	// Verify that there's an active association between the advertiser and affiliate organizations
	// and that the campaign and affiliate are visible to each other
	err = s.verifyAssociationAndVisibility(ctx, campaign.OrganizationID, affiliate.OrganizationID, req.CampaignID, req.AffiliateID)
	if err != nil {
		return nil, fmt.Errorf("association verification failed: %w", err)
	}

	// Check if tracking link already exists with same campaign and affiliate
	existingLink, err := s.trackingLinkRepo.GetTrackingLinkByCampaignAndAffiliate(
		ctx, req.CampaignID, req.AffiliateID, req.SourceID, req.Sub1, req.Sub2, req.Sub3, req.Sub4, req.Sub5)
	
	var trackingLink *domain.TrackingLink
	var isNew bool
	
	if err != nil {
		// No existing tracking link found, create a new one
		isNew = true
		
		// Create new tracking link entity
		trackingLink = &domain.TrackingLink{
			OrganizationID:      campaign.OrganizationID,
			CampaignID:          req.CampaignID,
			AffiliateID:         req.AffiliateID,
			Name:                req.Name,
			Description:         req.Description,
			Status:              "active",
			SourceID:            req.SourceID,
			Sub1:                req.Sub1,
			Sub2:                req.Sub2,
			Sub3:                req.Sub3,
			Sub4:                req.Sub4,
			Sub5:                req.Sub5,
			IsEncryptParameters: req.IsEncryptParameters,
			IsRedirectLink:      req.IsRedirectLink,
			InternalNotes:       req.InternalNotes,
			Tags:                req.Tags,
		}

		// Set timestamps
		now := time.Now()
		trackingLink.CreatedAt = now
		trackingLink.UpdatedAt = now

		// Create tracking link in database
		if err := s.trackingLinkRepo.CreateTrackingLink(ctx, trackingLink); err != nil {
			return nil, fmt.Errorf("failed to create tracking link: %w", err)
		}
	} else {
		// Existing tracking link found, update it
		isNew = false
		trackingLink = existingLink
		
		// Check if tracking parameters have changed
		parametersChanged := s.hasTrackingParametersChangedFromRequest(existingLink, req)
		
		// Update fields from request
		trackingLink.Name = req.Name
		trackingLink.Description = req.Description
		if req.SourceID != nil {
			trackingLink.SourceID = req.SourceID
		}
		if req.Sub1 != nil {
			trackingLink.Sub1 = req.Sub1
		}
		if req.Sub2 != nil {
			trackingLink.Sub2 = req.Sub2
		}
		if req.Sub3 != nil {
			trackingLink.Sub3 = req.Sub3
		}
		if req.Sub4 != nil {
			trackingLink.Sub4 = req.Sub4
		}
		if req.Sub5 != nil {
			trackingLink.Sub5 = req.Sub5
		}
		if req.IsEncryptParameters != nil {
			trackingLink.IsEncryptParameters = req.IsEncryptParameters
		}
		if req.IsRedirectLink != nil {
			trackingLink.IsRedirectLink = req.IsRedirectLink
		}
		if req.InternalNotes != nil {
			trackingLink.InternalNotes = req.InternalNotes
		}
		if req.Tags != nil {
			trackingLink.Tags = req.Tags
		}
		
		// Update timestamp
		trackingLink.UpdatedAt = time.Now()
		
		// Update tracking link in database
		if err := s.trackingLinkRepo.UpdateTrackingLink(ctx, trackingLink); err != nil {
			return nil, fmt.Errorf("failed to update tracking link: %w", err)
		}
		
		// If tracking parameters changed, we need to regenerate the tracking link
		if parametersChanged {
			logger.Info("Tracking parameters changed during upsert, regenerating tracking link", 
				"tracking_link_id", trackingLink.TrackingLinkID,
				"campaign_id", trackingLink.CampaignID,
				"affiliate_id", trackingLink.AffiliateID)
		}
	}

	var generatedURL string
	var providerData *string

	// Generate or regenerate tracking link via provider integration only if needed
	if isNew || (existingLink != nil && s.hasTrackingParametersChangedFromRequest(existingLink, req)) {
		generationReq := &domain.TrackingLinkGenerationRequest{
			CampaignID:              req.CampaignID,
			AffiliateID:             req.AffiliateID,
			Name:                    req.Name,
			Description:             req.Description,
			SourceID:                req.SourceID,
			Sub1:                    req.Sub1,
			Sub2:                    req.Sub2,
			Sub3:                    req.Sub3,
			Sub4:                    req.Sub4,
			Sub5:                    req.Sub5,
			IsEncryptParameters:     req.IsEncryptParameters,
			IsRedirectLink:          req.IsRedirectLink,
			NetworkTrackingDomainID: req.NetworkTrackingDomainID,
			NetworkOfferURLID:       req.NetworkOfferURLID,
			CreativeID:              req.CreativeID,
			NetworkTrafficSourceID:  req.NetworkTrafficSourceID,
		}

		var err error
		generatedURL, providerData, err = s.generateTrackingLinkViaProvider(ctx, trackingLink, generationReq)
		if err != nil {
			return nil, fmt.Errorf("failed to generate tracking link via provider: %w", err)
		}

		// Update tracking link with generated URL
		trackingLink.TrackingURL = &generatedURL
		if err := s.trackingLinkRepo.UpdateTrackingLink(ctx, trackingLink); err != nil {
			return nil, fmt.Errorf("failed to update tracking link with generated URL: %w", err)
		}
	} else {
		// Use existing tracking URL and provider data
		if trackingLink.TrackingURL != nil {
			generatedURL = *trackingLink.TrackingURL
		}
		
		// Get existing provider data if available
		if existingMapping, err := s.trackingLinkProviderRepo.GetTrackingLinkProviderMapping(ctx, trackingLink.TrackingLinkID, "everflow"); err == nil {
			providerData = existingMapping.ProviderData
		}
	}

	return &domain.TrackingLinkUpsertResponse{
		TrackingLink: trackingLink,
		GeneratedURL: generatedURL,
		ProviderData: providerData,
		IsNew:        isNew,
	}, nil
}

// generateTrackingLinkViaProvider generates tracking link via provider integration
func (s *trackingLinkService) generateTrackingLinkViaProvider(ctx context.Context, trackingLink *domain.TrackingLink, req *domain.TrackingLinkGenerationRequest) (string, *string, error) {
	// Get campaign provider mapping to get provider-specific campaign ID
	campaignMapping, err := s.campaignProviderRepo.GetCampaignProviderMapping(ctx, req.CampaignID, "everflow")
	if err != nil {
		return "", nil, fmt.Errorf("failed to get campaign provider mapping: %w", err)
	}

	// Get affiliate provider mapping to get provider-specific affiliate ID
	affiliateMapping, err := s.affiliateProviderRepo.GetAffiliateProviderMapping(ctx, req.AffiliateID, "everflow")
	if err != nil {
		return "", nil, fmt.Errorf("failed to get affiliate provider mapping: %w", err)
	}

	// Create provider-specific tracking link generation request
	providerReq := &domain.TrackingLinkGenerationRequest{
		CampaignID:              req.CampaignID,
		AffiliateID:             req.AffiliateID,
		Name:                    req.Name,
		Description:             req.Description,
		SourceID:                req.SourceID,
		Sub1:                    req.Sub1,
		Sub2:                    req.Sub2,
		Sub3:                    req.Sub3,
		Sub4:                    req.Sub4,
		Sub5:                    req.Sub5,
		IsEncryptParameters:     req.IsEncryptParameters,
		IsRedirectLink:          req.IsRedirectLink,
		NetworkTrackingDomainID: req.NetworkTrackingDomainID,
		NetworkOfferURLID:       req.NetworkOfferURLID,
		CreativeID:              req.CreativeID,
		NetworkTrafficSourceID:  req.NetworkTrafficSourceID,
	}

	// Generate tracking link via integration service
	response, err := s.integrationService.GenerateTrackingLink(ctx, providerReq, campaignMapping, affiliateMapping)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate tracking link via integration service: %w", err)
	}

	// Check if provider mapping already exists
	existingMapping, err := s.trackingLinkProviderRepo.GetTrackingLinkProviderMapping(ctx, trackingLink.TrackingLinkID, "everflow")
	now := time.Now()
	
	if err != nil {
		// No existing mapping found, create a new one
		providerMapping := &domain.TrackingLinkProviderMapping{
			TrackingLinkID:         trackingLink.TrackingLinkID,
			ProviderType:           "everflow",
			ProviderTrackingLinkID: nil, // Everflow doesn't return a specific tracking link ID
			ProviderData:           response.ProviderData,
			SyncStatus:             toStringPtr("synced"),
			LastSyncAt:             &now,
			CreatedAt:              now,
			UpdatedAt:              now,
		}

		// Create provider mapping
		if err := s.trackingLinkProviderRepo.CreateTrackingLinkProviderMapping(ctx, providerMapping); err != nil {
			return "", nil, fmt.Errorf("failed to create tracking link provider mapping: %w", err)
		}
	} else {
		// Update existing mapping
		existingMapping.ProviderData = response.ProviderData
		existingMapping.SyncStatus = toStringPtr("synced")
		existingMapping.LastSyncAt = &now
		existingMapping.UpdatedAt = now

		// Update provider mapping
		if err := s.trackingLinkProviderRepo.UpdateTrackingLinkProviderMapping(ctx, existingMapping); err != nil {
			return "", nil, fmt.Errorf("failed to update tracking link provider mapping: %w", err)
		}
	}

	return response.GeneratedURL, response.ProviderData, nil
}

// SyncTrackingLinkToProvider syncs a tracking link to the provider
func (s *trackingLinkService) SyncTrackingLinkToProvider(ctx context.Context, trackingLinkID int64) error {
	// Implementation would sync tracking link data to provider
	// For now, this is a placeholder
	return nil
}

// SyncTrackingLinkFromProvider syncs a tracking link from the provider
func (s *trackingLinkService) SyncTrackingLinkFromProvider(ctx context.Context, trackingLinkID int64) error {
	// Implementation would sync tracking link data from provider
	// For now, this is a placeholder
	return nil
}

// CreateTrackingLinkProviderMapping creates a new tracking link provider mapping
func (s *trackingLinkService) CreateTrackingLinkProviderMapping(ctx context.Context, mapping *domain.TrackingLinkProviderMapping) (*domain.TrackingLinkProviderMapping, error) {
	// Set timestamps
	now := time.Now()
	mapping.CreatedAt = now
	mapping.UpdatedAt = now

	if err := s.trackingLinkProviderRepo.CreateTrackingLinkProviderMapping(ctx, mapping); err != nil {
		return nil, fmt.Errorf("failed to create tracking link provider mapping: %w", err)
	}

	return mapping, nil
}

// GetTrackingLinkProviderMapping retrieves a tracking link provider mapping
func (s *trackingLinkService) GetTrackingLinkProviderMapping(ctx context.Context, trackingLinkID int64, providerType string) (*domain.TrackingLinkProviderMapping, error) {
	mapping, err := s.trackingLinkProviderRepo.GetTrackingLinkProviderMapping(ctx, trackingLinkID, providerType)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracking link provider mapping: %w", err)
	}

	return mapping, nil
}

// UpdateTrackingLinkProviderMapping updates an existing tracking link provider mapping
func (s *trackingLinkService) UpdateTrackingLinkProviderMapping(ctx context.Context, mapping *domain.TrackingLinkProviderMapping) error {
	if err := s.trackingLinkProviderRepo.UpdateTrackingLinkProviderMapping(ctx, mapping); err != nil {
		return fmt.Errorf("failed to update tracking link provider mapping: %w", err)
	}

	return nil
}

// DeleteTrackingLinkProviderMapping deletes a tracking link provider mapping
func (s *trackingLinkService) DeleteTrackingLinkProviderMapping(ctx context.Context, mappingID int64) error {
	if err := s.trackingLinkProviderRepo.DeleteTrackingLinkProviderMapping(ctx, mappingID); err != nil {
		return fmt.Errorf("failed to delete tracking link provider mapping: %w", err)
	}

	return nil
}

// validateTrackingLink validates tracking link business rules
func (s *trackingLinkService) validateTrackingLink(trackingLink *domain.TrackingLink) error {
	if trackingLink.Name == "" {
		return fmt.Errorf("tracking link name is required")
	}

	if trackingLink.OrganizationID <= 0 {
		return fmt.Errorf("valid organization ID is required")
	}

	if trackingLink.CampaignID <= 0 {
		return fmt.Errorf("valid campaign ID is required")
	}

	if trackingLink.AffiliateID <= 0 {
		return fmt.Errorf("valid affiliate ID is required")
	}

	// Validate status
	validStatuses := map[string]bool{
		"active":   true,
		"paused":   true,
		"archived": true,
	}
	if !validStatuses[trackingLink.Status] {
		return fmt.Errorf("invalid tracking link status: %s", trackingLink.Status)
	}

	return nil
}

// validateTrackingLinkGenerationRequest validates tracking link generation request
func (s *trackingLinkService) validateTrackingLinkGenerationRequest(req *domain.TrackingLinkGenerationRequest) error {
	if req.Name == "" {
		return fmt.Errorf("tracking link name is required")
	}

	if req.CampaignID <= 0 {
		return fmt.Errorf("valid campaign ID is required")
	}

	if req.AffiliateID <= 0 {
		return fmt.Errorf("valid affiliate ID is required")
	}

	return nil
}

// verifyAssociationAndVisibility verifies that there's an active association between organizations
// and that the campaign and affiliate are visible to each other
func (s *trackingLinkService) verifyAssociationAndVisibility(ctx context.Context, advertiserOrgID, affiliateOrgID, campaignID, affiliateID int64) error {
	// Get the association between the organizations
	association, err := s.orgAssociationService.GetAssociationByOrganizations(ctx, advertiserOrgID, affiliateOrgID)
	if err != nil {
		return fmt.Errorf("no association found between organizations %d and %d: %w", advertiserOrgID, affiliateOrgID, err)
	}

	// Verify the association is active
	if association.Status != "active" {
		return fmt.Errorf("association between organizations %d and %d is not active (status: %s)", advertiserOrgID, affiliateOrgID, association.Status)
	}

	// Check if the campaign is visible to the affiliate organization
	if !association.AllCampaignsVisible {
		// Parse the visible campaign IDs from JSON
		if association.VisibleCampaignIDs == nil {
			return fmt.Errorf("campaign %d is not visible to affiliate organization %d (no visible campaigns configured)", campaignID, affiliateOrgID)
		}
		
		var visibleCampaignIDs []int64
		if err := json.Unmarshal([]byte(*association.VisibleCampaignIDs), &visibleCampaignIDs); err != nil {
			return fmt.Errorf("failed to parse visible campaign IDs: %w", err)
		}
		
		// Check if the specific campaign is in the visible list
		campaignVisible := false
		for _, visibleCampaignID := range visibleCampaignIDs {
			if visibleCampaignID == campaignID {
				campaignVisible = true
				break
			}
		}
		if !campaignVisible {
			return fmt.Errorf("campaign %d is not visible to affiliate organization %d", campaignID, affiliateOrgID)
		}
	}

	// Check if the affiliate is visible to the advertiser organization
	if !association.AllAffiliatesVisible {
		// Parse the visible affiliate IDs from JSON
		if association.VisibleAffiliateIDs == nil {
			return fmt.Errorf("affiliate %d is not visible to advertiser organization %d (no visible affiliates configured)", affiliateID, advertiserOrgID)
		}
		
		var visibleAffiliateIDs []int64
		if err := json.Unmarshal([]byte(*association.VisibleAffiliateIDs), &visibleAffiliateIDs); err != nil {
			return fmt.Errorf("failed to parse visible affiliate IDs: %w", err)
		}
		
		// Check if the specific affiliate is in the visible list
		affiliateVisible := false
		for _, visibleAffiliateID := range visibleAffiliateIDs {
			if visibleAffiliateID == affiliateID {
				affiliateVisible = true
				break
			}
		}
		if !affiliateVisible {
			return fmt.Errorf("affiliate %d is not visible to advertiser organization %d", affiliateID, advertiserOrgID)
		}
	}

	return nil
}

// synchronizeTrackingLinkWithProvider synchronizes a tracking link with the provider (Everflow)
func (s *trackingLinkService) synchronizeTrackingLinkWithProvider(ctx context.Context, trackingLink *domain.TrackingLink) error {
	logger.Info("Starting tracking link synchronization", "tracking_link_id", trackingLink.TrackingLinkID)

	// Get campaign provider mapping
	campaignMapping, err := s.campaignProviderRepo.GetCampaignProviderMapping(ctx, trackingLink.CampaignID, "everflow")
	if err != nil {
		return fmt.Errorf("failed to get campaign provider mapping: %w", err)
	}

	// Get affiliate provider mapping
	affiliateMapping, err := s.affiliateProviderRepo.GetAffiliateProviderMapping(ctx, trackingLink.AffiliateID, "everflow")
	if err != nil {
		return fmt.Errorf("failed to get affiliate provider mapping: %w", err)
	}

	// Create tracking link in provider
	providerMapping, err := s.integrationService.CreateTrackingLink(ctx, trackingLink, campaignMapping, affiliateMapping)
	if err != nil {
		return fmt.Errorf("failed to create tracking link in provider: %w", err)
	}

	// Store provider mapping
	if err := s.trackingLinkProviderRepo.CreateTrackingLinkProviderMapping(ctx, providerMapping); err != nil {
		logger.Warn("Failed to store tracking link provider mapping", 
			"tracking_link_id", trackingLink.TrackingLinkID, 
			"error", err)
		// Don't return error here as the tracking link was successfully created in the provider
	}

	logger.Info("Successfully synchronized tracking link with provider", "tracking_link_id", trackingLink.TrackingLinkID)
	return nil
}

// hasTrackingParametersChanged checks if any tracking parameters have changed
func (s *trackingLinkService) hasTrackingParametersChanged(existing, updated *domain.TrackingLink) bool {
	// Compare source_id
	if !stringPtrEqual(existing.SourceID, updated.SourceID) {
		logger.Debug("SourceID changed", 
			"existing", stringPtrValue(existing.SourceID), 
			"updated", stringPtrValue(updated.SourceID))
		return true
	}

	// Compare sub1
	if !stringPtrEqual(existing.Sub1, updated.Sub1) {
		logger.Debug("Sub1 changed", 
			"existing", stringPtrValue(existing.Sub1), 
			"updated", stringPtrValue(updated.Sub1))
		return true
	}

	// Compare sub2
	if !stringPtrEqual(existing.Sub2, updated.Sub2) {
		logger.Debug("Sub2 changed", 
			"existing", stringPtrValue(existing.Sub2), 
			"updated", stringPtrValue(updated.Sub2))
		return true
	}

	// Compare sub3
	if !stringPtrEqual(existing.Sub3, updated.Sub3) {
		logger.Debug("Sub3 changed", 
			"existing", stringPtrValue(existing.Sub3), 
			"updated", stringPtrValue(updated.Sub3))
		return true
	}

	// Compare sub4
	if !stringPtrEqual(existing.Sub4, updated.Sub4) {
		logger.Debug("Sub4 changed", 
			"existing", stringPtrValue(existing.Sub4), 
			"updated", stringPtrValue(updated.Sub4))
		return true
	}

	// Compare sub5
	if !stringPtrEqual(existing.Sub5, updated.Sub5) {
		logger.Debug("Sub5 changed", 
			"existing", stringPtrValue(existing.Sub5), 
			"updated", stringPtrValue(updated.Sub5))
		return true
	}

	// Compare tags
	if !stringPtrEqual(existing.Tags, updated.Tags) {
		logger.Debug("Tags changed", 
			"existing", stringPtrValue(existing.Tags), 
			"updated", stringPtrValue(updated.Tags))
		return true
	}

	return false
}

// validateTrackingLinkUpsertRequest validates tracking link upsert request
func (s *trackingLinkService) validateTrackingLinkUpsertRequest(req *domain.TrackingLinkUpsertRequest) error {
	if req.Name == "" {
		return fmt.Errorf("tracking link name is required")
	}

	if req.CampaignID <= 0 {
		return fmt.Errorf("valid campaign ID is required")
	}

	if req.AffiliateID <= 0 {
		return fmt.Errorf("valid affiliate ID is required")
	}

	return nil
}

// hasTrackingParametersChangedFromRequest checks if any tracking parameters have changed from request
func (s *trackingLinkService) hasTrackingParametersChangedFromRequest(existing *domain.TrackingLink, req *domain.TrackingLinkUpsertRequest) bool {
	// Compare source_id
	if !stringPtrEqual(existing.SourceID, req.SourceID) {
		logger.Debug("SourceID changed from request", 
			"existing", stringPtrValue(existing.SourceID), 
			"request", stringPtrValue(req.SourceID))
		return true
	}

	// Compare sub1
	if !stringPtrEqual(existing.Sub1, req.Sub1) {
		logger.Debug("Sub1 changed from request", 
			"existing", stringPtrValue(existing.Sub1), 
			"request", stringPtrValue(req.Sub1))
		return true
	}

	// Compare sub2
	if !stringPtrEqual(existing.Sub2, req.Sub2) {
		logger.Debug("Sub2 changed from request", 
			"existing", stringPtrValue(existing.Sub2), 
			"request", stringPtrValue(req.Sub2))
		return true
	}

	// Compare sub3
	if !stringPtrEqual(existing.Sub3, req.Sub3) {
		logger.Debug("Sub3 changed from request", 
			"existing", stringPtrValue(existing.Sub3), 
			"request", stringPtrValue(req.Sub3))
		return true
	}

	// Compare sub4
	if !stringPtrEqual(existing.Sub4, req.Sub4) {
		logger.Debug("Sub4 changed from request", 
			"existing", stringPtrValue(existing.Sub4), 
			"request", stringPtrValue(req.Sub4))
		return true
	}

	// Compare sub5
	if !stringPtrEqual(existing.Sub5, req.Sub5) {
		logger.Debug("Sub5 changed from request", 
			"existing", stringPtrValue(existing.Sub5), 
			"request", stringPtrValue(req.Sub5))
		return true
	}

	// Compare tags
	if !stringPtrEqual(existing.Tags, req.Tags) {
		logger.Debug("Tags changed from request", 
			"existing", stringPtrValue(existing.Tags), 
			"request", stringPtrValue(req.Tags))
		return true
	}

	return false
}

// stringPtrEqual compares two string pointers for equality
func stringPtrEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// stringPtrValue returns the value of a string pointer or empty string if nil
func stringPtrValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Helper function to create string pointer
func toStringPtr(s string) *string {
	return &s
}
