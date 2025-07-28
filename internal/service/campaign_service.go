package service

import (
	"context"
	"fmt"
	"log"

	"github.com/affiliate-backend/internal/domain"
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
}

// campaignService implements CampaignService
type campaignService struct {
	campaignRepo       repository.CampaignRepository
	integrationService provider.IntegrationService
}

// NewCampaignService creates a new campaign service
func NewCampaignService(campaignRepo repository.CampaignRepository, integrationService provider.IntegrationService) CampaignService {
	return &campaignService{
		campaignRepo:       campaignRepo,
		integrationService: integrationService,
	}
}

// CreateCampaign creates a new campaign
func (s *campaignService) CreateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	log.Printf("🚀 CAMPAIGN SERVICE: Starting campaign creation for campaign_id=%d, advertiser_id=%d, name='%s'", campaign.CampaignID, campaign.AdvertiserID, campaign.Name)
	
	// Validate campaign data
	log.Printf("🔍 CAMPAIGN SERVICE: Validating campaign data...")
	if err := s.validateCampaign(campaign); err != nil {
		log.Printf("❌ CAMPAIGN SERVICE: Campaign validation failed: %v", err)
		return fmt.Errorf("campaign validation failed: %w", err)
	}
	log.Printf("✅ CAMPAIGN SERVICE: Campaign validation passed")

	// Step 1: Create campaign in local repository
	log.Printf("🔄 CAMPAIGN SERVICE: Creating campaign in local repository...")
	if err := s.campaignRepo.CreateCampaign(ctx, campaign); err != nil {
		log.Printf("❌ CAMPAIGN SERVICE: Failed to create campaign in repository: %v", err)
		return fmt.Errorf("failed to create campaign: %w", err)
	}
	log.Printf("✅ CAMPAIGN SERVICE: Successfully created campaign in repository with ID=%d", campaign.CampaignID)

	// Step 2: Call IntegrationService to create campaign in provider (Everflow)
	// The integration service handles provider mapping creation internally
	log.Printf("🔄 CAMPAIGN SERVICE: Calling integration service to create campaign in provider...")
	_, err := s.integrationService.CreateCampaign(ctx, *campaign)
	if err != nil {
		// Log error but don't fail the operation since local creation succeeded
		log.Printf("⚠️  CAMPAIGN SERVICE: Failed to create campaign in provider: %v", err)
		fmt.Printf("Warning: failed to create campaign in provider: %v\n", err)
		return nil
	}
	log.Printf("✅ CAMPAIGN SERVICE: Successfully created campaign in provider")

	log.Printf("✅ CAMPAIGN SERVICE: Campaign creation completed successfully")
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
	// Validate campaign data
	if err := s.validateCampaign(campaign); err != nil {
		return fmt.Errorf("campaign validation failed: %w", err)
	}

	// Update campaign in repository
	if err := s.campaignRepo.UpdateCampaign(ctx, campaign); err != nil {
		return fmt.Errorf("failed to update campaign: %w", err)
	}

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
