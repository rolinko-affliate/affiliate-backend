package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/repository"
)

// CampaignService defines the interface for campaign operations
type CampaignService interface {
	CreateCampaign(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error)
	GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error)
	UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error
	ListCampaignsByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]*domain.Campaign, error)
	ListCampaignsByAdvertiser(ctx context.Context, advertiserID int64, page, pageSize int) ([]*domain.Campaign, error)
	DeleteCampaign(ctx context.Context, id int64) error
	
	// Provider offer methods
	CreateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) (*domain.CampaignProviderOffer, error)
	GetCampaignProviderOfferByID(ctx context.Context, id int64) (*domain.CampaignProviderOffer, error)
	UpdateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error
	ListCampaignProviderOffersByCampaign(ctx context.Context, campaignID int64) ([]*domain.CampaignProviderOffer, error)
	DeleteCampaignProviderOffer(ctx context.Context, id int64) error
}

// campaignService implements CampaignService
type campaignService struct {
	campaignRepo   repository.CampaignRepository
	advertiserRepo repository.AdvertiserRepository
	orgRepo        repository.OrganizationRepository
}

// NewCampaignService creates a new campaign service
func NewCampaignService(
	campaignRepo repository.CampaignRepository,
	advertiserRepo repository.AdvertiserRepository,
	orgRepo repository.OrganizationRepository,
) CampaignService {
	return &campaignService{
		campaignRepo:   campaignRepo,
		advertiserRepo: advertiserRepo,
		orgRepo:        orgRepo,
	}
}

// CreateCampaign creates a new campaign
func (s *campaignService) CreateCampaign(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error) {
	// Validate organization exists
	_, err := s.orgRepo.GetOrganizationByID(ctx, campaign.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	// Validate advertiser exists
	advertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, campaign.AdvertiserID)
	if err != nil {
		return nil, fmt.Errorf("advertiser not found: %w", err)
	}

	// Validate advertiser belongs to the organization
	if advertiser.OrganizationID != campaign.OrganizationID {
		return nil, fmt.Errorf("advertiser does not belong to the specified organization")
	}

	// Validate required fields
	if campaign.Name == "" {
		return nil, fmt.Errorf("campaign name cannot be empty")
	}

	// Set default status if not provided
	if campaign.Status == "" {
		campaign.Status = "draft"
	}

	// Validate status
	validStatuses := map[string]bool{
		"draft":    true,
		"active":   true,
		"paused":   true,
		"archived": true,
	}
	if !validStatuses[campaign.Status] {
		return nil, fmt.Errorf("invalid status: %s", campaign.Status)
	}

	// Validate dates
	if campaign.StartDate != nil && campaign.EndDate != nil {
		if campaign.StartDate.After(*campaign.EndDate) {
			return nil, fmt.Errorf("start date cannot be after end date")
		}
	}

	if err := s.campaignRepo.CreateCampaign(ctx, campaign); err != nil {
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	return campaign, nil
}

// GetCampaignByID retrieves a campaign by ID
func (s *campaignService) GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error) {
	return s.campaignRepo.GetCampaignByID(ctx, id)
}

// UpdateCampaign updates a campaign
func (s *campaignService) UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	// Validate required fields
	if campaign.Name == "" {
		return fmt.Errorf("campaign name cannot be empty")
	}

	// Validate status
	validStatuses := map[string]bool{
		"draft":    true,
		"active":   true,
		"paused":   true,
		"archived": true,
	}
	if !validStatuses[campaign.Status] {
		return fmt.Errorf("invalid status: %s", campaign.Status)
	}

	// Validate dates
	if campaign.StartDate != nil && campaign.EndDate != nil {
		if campaign.StartDate.After(*campaign.EndDate) {
			return fmt.Errorf("start date cannot be after end date")
		}
	}

	// Get existing campaign to verify organization and advertiser
	existingCampaign, err := s.campaignRepo.GetCampaignByID(ctx, campaign.CampaignID)
	if err != nil {
		return fmt.Errorf("failed to get existing campaign: %w", err)
	}

	// Ensure organization and advertiser IDs are not changed
	campaign.OrganizationID = existingCampaign.OrganizationID
	campaign.AdvertiserID = existingCampaign.AdvertiserID

	return s.campaignRepo.UpdateCampaign(ctx, campaign)
}

// ListCampaignsByOrganization retrieves a list of campaigns for an organization with pagination
func (s *campaignService) ListCampaignsByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]*domain.Campaign, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.campaignRepo.ListCampaignsByOrganization(ctx, orgID, pageSize, offset)
}

// ListCampaignsByAdvertiser retrieves a list of campaigns for an advertiser with pagination
func (s *campaignService) ListCampaignsByAdvertiser(ctx context.Context, advertiserID int64, page, pageSize int) ([]*domain.Campaign, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.campaignRepo.ListCampaignsByAdvertiser(ctx, advertiserID, pageSize, offset)
}

// DeleteCampaign deletes a campaign
func (s *campaignService) DeleteCampaign(ctx context.Context, id int64) error {
	return s.campaignRepo.DeleteCampaign(ctx, id)
}

// CreateCampaignProviderOffer creates a new campaign provider offer
func (s *campaignService) CreateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) (*domain.CampaignProviderOffer, error) {
	// Validate campaign exists
	_, err := s.campaignRepo.GetCampaignByID(ctx, offer.CampaignID)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	// Validate provider type
	if offer.ProviderType != "everflow" {
		return nil, fmt.Errorf("invalid provider type: %s", offer.ProviderType)
	}

	// Validate provider offer config JSON if provided
	if offer.ProviderOfferConfig != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*offer.ProviderOfferConfig), &jsonData); err != nil {
			return nil, fmt.Errorf("invalid provider offer config JSON: %w", err)
		}
	}

	// Set default values
	if offer.LastSyncedAt == nil {
		now := time.Now()
		offer.LastSyncedAt = &now
	}

	if err := s.campaignRepo.CreateCampaignProviderOffer(ctx, offer); err != nil {
		return nil, fmt.Errorf("failed to create campaign provider offer: %w", err)
	}

	return offer, nil
}

// GetCampaignProviderOfferByID retrieves a campaign provider offer by ID
func (s *campaignService) GetCampaignProviderOfferByID(ctx context.Context, id int64) (*domain.CampaignProviderOffer, error) {
	return s.campaignRepo.GetCampaignProviderOfferByID(ctx, id)
}

// UpdateCampaignProviderOffer updates a campaign provider offer
func (s *campaignService) UpdateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error {
	// Validate provider offer config JSON if provided
	if offer.ProviderOfferConfig != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*offer.ProviderOfferConfig), &jsonData); err != nil {
			return fmt.Errorf("invalid provider offer config JSON: %w", err)
		}
	}

	// Get existing offer to verify campaign ID
	existingOffer, err := s.campaignRepo.GetCampaignProviderOfferByID(ctx, offer.ProviderOfferID)
	if err != nil {
		return fmt.Errorf("failed to get existing campaign provider offer: %w", err)
	}

	// Ensure campaign ID is not changed
	offer.CampaignID = existingOffer.CampaignID
	offer.ProviderType = existingOffer.ProviderType

	// Update last synced time
	now := time.Now()
	offer.LastSyncedAt = &now

	return s.campaignRepo.UpdateCampaignProviderOffer(ctx, offer)
}

// ListCampaignProviderOffersByCampaign retrieves a list of campaign provider offers for a campaign
func (s *campaignService) ListCampaignProviderOffersByCampaign(ctx context.Context, campaignID int64) ([]*domain.CampaignProviderOffer, error) {
	return s.campaignRepo.ListCampaignProviderOffersByCampaign(ctx, campaignID)
}

// DeleteCampaignProviderOffer deletes a campaign provider offer
func (s *campaignService) DeleteCampaignProviderOffer(ctx context.Context, id int64) error {
	return s.campaignRepo.DeleteCampaignProviderOffer(ctx, id)
}