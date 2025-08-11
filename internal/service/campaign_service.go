package service

import (
	"context"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/logger"
	"github.com/affiliate-backend/internal/platform/provider"
	"github.com/affiliate-backend/internal/repository"
)

// CampaignService defines the interface for campaign business logic
type CampaignService interface {
	CreateCampaign(ctx context.Context, campaign *domain.Campaign) error
	GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error)
	UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error
	ListCampaignsByAdvertiser(ctx context.Context, advertiserID int64, limit, offset int) ([]*domain.Campaign, error)
	ListCampaignsByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.Campaign, error)
	DeleteCampaign(ctx context.Context, id int64) error
	GetProviderMapping(ctx context.Context, campaignID int64, providerType string) (*domain.CampaignProviderMapping, error)
}

// campaignService implements CampaignService
type campaignService struct {
	campaignRepo               repository.CampaignRepository
	campaignProviderMappingRepo repository.CampaignProviderMappingRepository
	integrationService         provider.IntegrationService
}

// NewCampaignService creates a new campaign service
func NewCampaignService(campaignRepo repository.CampaignRepository, campaignProviderMappingRepo repository.CampaignProviderMappingRepository, integrationService provider.IntegrationService) CampaignService {
	return &campaignService{
		campaignRepo:               campaignRepo,
		campaignProviderMappingRepo: campaignProviderMappingRepo,
		integrationService:         integrationService,
	}
}

// CreateCampaign creates a new campaign
func (s *campaignService) CreateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	logger.Info("Starting campaign creation", "campaign_id", campaign.CampaignID, "advertiser_id", campaign.AdvertiserID, "name", campaign.Name)
	
	// Validate campaign data
	logger.Debug("Validating campaign data", "campaign_id", campaign.CampaignID)
	if err := s.validateCampaign(campaign); err != nil {
		logger.Error("Campaign validation failed", "campaign_id", campaign.CampaignID, "error", err)
		return fmt.Errorf("campaign validation failed: %w", err)
	}
	logger.Debug("Campaign validation passed", "campaign_id", campaign.CampaignID)

	// Step 1: Create campaign in local repository
	logger.Debug("Creating campaign in local repository", "campaign_id", campaign.CampaignID)
	if err := s.campaignRepo.CreateCampaign(ctx, campaign); err != nil {
		logger.Error("Failed to create campaign in repository", "campaign_id", campaign.CampaignID, "error", err)
		return fmt.Errorf("failed to create campaign: %w", err)
	}
	logger.Info("Successfully created campaign in repository", "campaign_id", campaign.CampaignID)

	// Step 2: Call IntegrationService to create campaign in provider (Everflow)
	// The integration service handles provider mapping creation internally
	logger.Debug("Calling integration service to create campaign in provider", "campaign_id", campaign.CampaignID)
	_, err := s.integrationService.CreateCampaign(ctx, *campaign)
	if err != nil {
		// Log error but don't fail the operation since local creation succeeded
		logger.Warn("Failed to create campaign in provider", "campaign_id", campaign.CampaignID, "error", err)
		return nil
	}
	logger.Info("Successfully created campaign in provider", "campaign_id", campaign.CampaignID)

	logger.Info("Campaign creation completed successfully", "campaign_id", campaign.CampaignID)
	return nil
}

// GetCampaignByID retrieves a campaign by its ID
func (s *campaignService) GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error) {
	campaign, err := s.campaignRepo.GetCampaignByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign: %w", err)
	}

	return campaign, nil
}

// UpdateCampaign updates an existing campaign
func (s *campaignService) UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	logger.Info("Starting campaign update", "campaign_id", campaign.CampaignID, "name", campaign.Name)
	
	// Validate campaign data
	logger.Debug("Validating campaign data", "campaign_id", campaign.CampaignID)
	if err := s.validateCampaign(campaign); err != nil {
		logger.Error("Campaign validation failed", "campaign_id", campaign.CampaignID, "error", err)
		return fmt.Errorf("campaign validation failed: %w", err)
	}
	logger.Debug("Campaign validation passed", "campaign_id", campaign.CampaignID)

	// Step 1: Update campaign in local repository
	logger.Debug("Updating campaign in local repository", "campaign_id", campaign.CampaignID)
	if err := s.campaignRepo.UpdateCampaign(ctx, campaign); err != nil {
		logger.Error("Failed to update campaign in repository", "campaign_id", campaign.CampaignID, "error", err)
		return fmt.Errorf("failed to update campaign: %w", err)
	}
	logger.Info("Successfully updated campaign in repository", "campaign_id", campaign.CampaignID)

	// Step 2: Call IntegrationService to update campaign in provider (Everflow)
	logger.Debug("Calling integration service to update campaign in provider", "campaign_id", campaign.CampaignID)
	err := s.integrationService.UpdateCampaign(ctx, *campaign)
	if err != nil {
		// Log error but don't fail the operation since local update succeeded
		logger.Warn("Failed to update campaign in provider", "campaign_id", campaign.CampaignID, "error", err)
		return nil
	}
	logger.Info("Successfully updated campaign in provider", "campaign_id", campaign.CampaignID)

	logger.Info("Campaign update completed successfully", "campaign_id", campaign.CampaignID)
	return nil
}

// ListCampaignsByAdvertiser retrieves campaigns for a specific advertiser
func (s *campaignService) ListCampaignsByAdvertiser(ctx context.Context, advertiserID int64, limit, offset int) ([]*domain.Campaign, error) {
	campaigns, err := s.campaignRepo.ListCampaignsByAdvertiser(ctx, advertiserID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list campaigns by advertiser: %w", err)
	}

	return campaigns, nil
}

// ListCampaignsByOrganization retrieves campaigns for a specific organization
func (s *campaignService) ListCampaignsByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.Campaign, error) {
	campaigns, err := s.campaignRepo.ListCampaignsByOrganization(ctx, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list campaigns by organization: %w", err)
	}

	return campaigns, nil
}

// DeleteCampaign deletes a campaign by its ID
func (s *campaignService) DeleteCampaign(ctx context.Context, id int64) error {
	if err := s.campaignRepo.DeleteCampaign(ctx, id); err != nil {
		return fmt.Errorf("failed to delete campaign: %w", err)
	}

	return nil
}

// validateCampaign validates campaign business rules
func (s *campaignService) validateCampaign(campaign *domain.Campaign) error {
	if campaign.Name == "" {
		return fmt.Errorf("campaign name is required")
	}

	if campaign.OrganizationID <= 0 {
		return fmt.Errorf("valid organization ID is required")
	}

	if campaign.AdvertiserID <= 0 {
		return fmt.Errorf("valid advertiser ID is required")
	}

	// Validate status
	validStatuses := map[string]bool{
		"draft":    true,
		"active":   true,
		"paused":   true,
		"archived": true,
	}
	if !validStatuses[campaign.Status] {
		return fmt.Errorf("invalid campaign status: %s", campaign.Status)
	}

	// Validate date range if both dates are provided
	if campaign.StartDate != nil && campaign.EndDate != nil {
		if campaign.EndDate.Before(*campaign.StartDate) {
			return fmt.Errorf("end date cannot be before start date")
		}
	}

	return nil
}

// GetProviderMapping retrieves a campaign provider mapping
func (s *campaignService) GetProviderMapping(ctx context.Context, campaignID int64, providerType string) (*domain.CampaignProviderMapping, error) {
	return s.campaignProviderMappingRepo.GetCampaignProviderMapping(ctx, campaignID, providerType)
}
