package service

import (
	"context"
	"fmt"
	"time"

	"github.com/affiliate-backend/internal/domain"
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

	// Tracking link generation
	GenerateTrackingLink(ctx context.Context, req *domain.TrackingLinkGenerationRequest) (*domain.TrackingLinkGenerationResponse, error)
	RegenerateTrackingLink(ctx context.Context, trackingLinkID int64) (*domain.TrackingLinkGenerationResponse, error)

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
) TrackingLinkService {
	return &trackingLinkService{
		trackingLinkRepo:         trackingLinkRepo,
		trackingLinkProviderRepo: trackingLinkProviderRepo,
		campaignRepo:             campaignRepo,
		affiliateRepo:            affiliateRepo,
		campaignProviderRepo:     campaignProviderRepo,
		affiliateProviderRepo:    affiliateProviderRepo,
		integrationService:       integrationService,
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

	// Update tracking link in repository
	if err := s.trackingLinkRepo.UpdateTrackingLink(ctx, trackingLink); err != nil {
		return fmt.Errorf("failed to update tracking link: %w", err)
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
	_, err = s.affiliateRepo.GetAffiliateByID(ctx, req.AffiliateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get affiliate: %w", err)
	}

	// Check if tracking link already exists with same parameters
	existingLink, err := s.trackingLinkRepo.GetTrackingLinkByCampaignAndAffiliate(
		ctx, req.CampaignID, req.AffiliateID, req.SourceID, req.Sub1, req.Sub2, req.Sub3, req.Sub4, req.Sub5)
	if err == nil {
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

	// Create provider mapping for the tracking link
	providerMapping := &domain.TrackingLinkProviderMapping{
		TrackingLinkID:         trackingLink.TrackingLinkID,
		ProviderType:           "everflow",
		ProviderTrackingLinkID: nil, // Everflow doesn't return a specific tracking link ID
		ProviderData:           response.ProviderData,
		SyncStatus:             toStringPtr("synced"),
		LastSyncAt:             &time.Time{},
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}

	now := time.Now()
	providerMapping.LastSyncAt = &now

	// Save provider mapping
	if err := s.trackingLinkProviderRepo.CreateTrackingLinkProviderMapping(ctx, providerMapping); err != nil {
		return "", nil, fmt.Errorf("failed to create tracking link provider mapping: %w", err)
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

// Helper function to create string pointer
func toStringPtr(s string) *string {
	return &s
}
