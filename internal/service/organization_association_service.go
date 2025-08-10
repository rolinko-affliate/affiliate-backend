package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/repository"
	"github.com/google/uuid"
)

// OrganizationAssociationService defines the interface for organization association business logic
type OrganizationAssociationService interface {
	// Association management
	CreateInvitation(ctx context.Context, req *domain.CreateAssociationRequest, requestedByUserID string) (*domain.OrganizationAssociation, error)
	CreateRequest(ctx context.Context, req *domain.CreateAssociationRequest, requestedByUserID string) (*domain.OrganizationAssociation, error)
	ApproveAssociation(ctx context.Context, associationID int64, approvedByUserID string) (*domain.OrganizationAssociation, error)
	RejectAssociation(ctx context.Context, associationID int64, approvedByUserID string) (*domain.OrganizationAssociation, error)
	SuspendAssociation(ctx context.Context, associationID int64, suspendedByUserID string) (*domain.OrganizationAssociation, error)
	ReactivateAssociation(ctx context.Context, associationID int64, reactivatedByUserID string) (*domain.OrganizationAssociation, error)
	
	// Visibility management
	UpdateVisibility(ctx context.Context, associationID int64, req *domain.UpdateAssociationRequest, updatedByUserID string) (*domain.OrganizationAssociation, error)
	
	// Retrieval
	GetAssociationByID(ctx context.Context, id int64) (*domain.OrganizationAssociation, error)
	GetAssociationByIDWithDetails(ctx context.Context, id int64) (*domain.OrganizationAssociationWithDetails, error)
	GetAssociationByOrganizations(ctx context.Context, advertiserOrgID, affiliateOrgID int64) (*domain.OrganizationAssociation, error)
	ListAssociations(ctx context.Context, filter *domain.AssociationListFilter) ([]*domain.OrganizationAssociation, error)
	ListAssociationsWithDetails(ctx context.Context, filter *domain.AssociationListFilter) ([]*domain.OrganizationAssociationWithDetails, error)
	GetAssociationsForAdvertiser(ctx context.Context, advertiserOrgID int64, status *domain.AssociationStatus) ([]*domain.OrganizationAssociation, error)
	GetAssociationsForAffiliate(ctx context.Context, affiliateOrgID int64, status *domain.AssociationStatus) ([]*domain.OrganizationAssociation, error)
	
	// Validation and utilities
	ValidateAssociationAccess(ctx context.Context, associationID int64, userOrgID int64) error
	CanUserManageAssociation(ctx context.Context, associationID int64, userID string, userOrgID int64) (bool, error)
	GetVisibleAffiliates(ctx context.Context, associationID int64) ([]int64, error)
	GetVisibleCampaigns(ctx context.Context, associationID int64) ([]int64, error)
	
	// Visibility queries
	GetVisibleAffiliatesForAdvertiser(ctx context.Context, advertiserOrgID int64, affiliateOrgID *int64) ([]*domain.Affiliate, error)
	GetVisibleCampaignsForAffiliate(ctx context.Context, affiliateOrgID int64, advertiserOrgID *int64) ([]*domain.Campaign, error)
}

// organizationAssociationService implements OrganizationAssociationService
type organizationAssociationService struct {
	associationRepo repository.OrganizationAssociationRepository
	orgRepo         repository.OrganizationRepository
	profileRepo     repository.ProfileRepository
	affiliateRepo   repository.AffiliateRepository
	campaignRepo    repository.CampaignRepository
}

// NewOrganizationAssociationService creates a new organization association service
func NewOrganizationAssociationService(
	associationRepo repository.OrganizationAssociationRepository,
	orgRepo repository.OrganizationRepository,
	profileRepo repository.ProfileRepository,
	affiliateRepo repository.AffiliateRepository,
	campaignRepo repository.CampaignRepository,
) OrganizationAssociationService {
	return &organizationAssociationService{
		associationRepo: associationRepo,
		orgRepo:         orgRepo,
		profileRepo:     profileRepo,
		affiliateRepo:   affiliateRepo,
		campaignRepo:    campaignRepo,
	}
}

// CreateInvitation creates a new invitation from advertiser to affiliate
func (s *organizationAssociationService) CreateInvitation(ctx context.Context, req *domain.CreateAssociationRequest, requestedByUserID string) (*domain.OrganizationAssociation, error) {
	// Validate that the advertiser organization exists and is of correct type
	advOrg, err := s.orgRepo.GetOrganizationByID(ctx, req.AdvertiserOrgID)
	if err != nil {
		return nil, fmt.Errorf("advertiser organization not found: %w", err)
	}
	if advOrg.Type != domain.OrganizationTypeAdvertiser {
		return nil, fmt.Errorf("organization %d is not an advertiser organization", req.AdvertiserOrgID)
	}

	// Validate that the affiliate organization exists and is of correct type
	affOrg, err := s.orgRepo.GetOrganizationByID(ctx, req.AffiliateOrgID)
	if err != nil {
		return nil, fmt.Errorf("affiliate organization not found: %w", err)
	}
	if affOrg.Type != domain.OrganizationTypeAffiliate {
		return nil, fmt.Errorf("organization %d is not an affiliate organization", req.AffiliateOrgID)
	}

	// Check if association already exists
	existing, err := s.associationRepo.GetAssociationByOrganizations(ctx, req.AdvertiserOrgID, req.AffiliateOrgID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("association already exists between organizations %d and %d", req.AdvertiserOrgID, req.AffiliateOrgID)
	}

	// Convert visible affiliate IDs to JSON
	var visibleAffiliateIDs *string
	if len(req.VisibleAffiliateIDs) > 0 {
		jsonBytes, err := json.Marshal(req.VisibleAffiliateIDs)
		if err != nil {
			return nil, fmt.Errorf("error marshaling visible affiliate IDs: %w", err)
		}
		jsonStr := string(jsonBytes)
		visibleAffiliateIDs = &jsonStr
	}

	// Set default visibility for campaigns (all visible for invitations)
	allCampaignsVisible := true
	if req.AllCampaignsVisible != nil {
		allCampaignsVisible = *req.AllCampaignsVisible
	}

	// Set default visibility for affiliates
	allAffiliatesVisible := true
	if req.AllAffiliatesVisible != nil {
		allAffiliatesVisible = *req.AllAffiliatesVisible
	}

	association := &domain.OrganizationAssociation{
		AdvertiserOrgID:      req.AdvertiserOrgID,
		AffiliateOrgID:       req.AffiliateOrgID,
		Status:               domain.AssociationStatusPending,
		AssociationType:      domain.AssociationTypeInvitation,
		VisibleAffiliateIDs:  visibleAffiliateIDs,
		AllAffiliatesVisible: allAffiliatesVisible,
		AllCampaignsVisible:  allCampaignsVisible,
		RequestedByUserID:    &requestedByUserID,
		Message:              req.Message,
	}

	if err := association.Validate(); err != nil {
		return nil, fmt.Errorf("invalid association data: %w", err)
	}

	if err := s.associationRepo.CreateAssociation(ctx, association); err != nil {
		return nil, fmt.Errorf("error creating invitation: %w", err)
	}

	return association, nil
}

// CreateRequest creates a new request from affiliate to advertiser
func (s *organizationAssociationService) CreateRequest(ctx context.Context, req *domain.CreateAssociationRequest, requestedByUserID string) (*domain.OrganizationAssociation, error) {
	// Validate that the advertiser organization exists and is of correct type
	advOrg, err := s.orgRepo.GetOrganizationByID(ctx, req.AdvertiserOrgID)
	if err != nil {
		return nil, fmt.Errorf("advertiser organization not found: %w", err)
	}
	if advOrg.Type != domain.OrganizationTypeAdvertiser {
		return nil, fmt.Errorf("organization %d is not an advertiser organization", req.AdvertiserOrgID)
	}

	// Validate that the affiliate organization exists and is of correct type
	affOrg, err := s.orgRepo.GetOrganizationByID(ctx, req.AffiliateOrgID)
	if err != nil {
		return nil, fmt.Errorf("affiliate organization not found: %w", err)
	}
	if affOrg.Type != domain.OrganizationTypeAffiliate {
		return nil, fmt.Errorf("organization %d is not an affiliate organization", req.AffiliateOrgID)
	}

	// Check if association already exists
	existing, err := s.associationRepo.GetAssociationByOrganizations(ctx, req.AdvertiserOrgID, req.AffiliateOrgID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("association already exists between organizations %d and %d", req.AdvertiserOrgID, req.AffiliateOrgID)
	}

	// Convert visible campaign IDs to JSON
	var visibleCampaignIDs *string
	if len(req.VisibleCampaignIDs) > 0 {
		jsonBytes, err := json.Marshal(req.VisibleCampaignIDs)
		if err != nil {
			return nil, fmt.Errorf("error marshaling visible campaign IDs: %w", err)
		}
		jsonStr := string(jsonBytes)
		visibleCampaignIDs = &jsonStr
	}

	// Set default visibility for affiliates (all visible for requests)
	allAffiliatesVisible := true
	if req.AllAffiliatesVisible != nil {
		allAffiliatesVisible = *req.AllAffiliatesVisible
	}

	// Set default visibility for campaigns
	allCampaignsVisible := true
	if req.AllCampaignsVisible != nil {
		allCampaignsVisible = *req.AllCampaignsVisible
	}

	association := &domain.OrganizationAssociation{
		AdvertiserOrgID:      req.AdvertiserOrgID,
		AffiliateOrgID:       req.AffiliateOrgID,
		Status:               domain.AssociationStatusPending,
		AssociationType:      domain.AssociationTypeRequest,
		VisibleCampaignIDs:   visibleCampaignIDs,
		AllAffiliatesVisible: allAffiliatesVisible,
		AllCampaignsVisible:  allCampaignsVisible,
		RequestedByUserID:    &requestedByUserID,
		Message:              req.Message,
	}

	if err := association.Validate(); err != nil {
		return nil, fmt.Errorf("invalid association data: %w", err)
	}

	if err := s.associationRepo.CreateAssociation(ctx, association); err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	return association, nil
}

// ApproveAssociation approves a pending association
func (s *organizationAssociationService) ApproveAssociation(ctx context.Context, associationID int64, approvedByUserID string) (*domain.OrganizationAssociation, error) {
	association, err := s.associationRepo.GetAssociationByID(ctx, associationID)
	if err != nil {
		return nil, fmt.Errorf("association not found: %w", err)
	}

	if !association.CanBeActivated() {
		return nil, fmt.Errorf("association cannot be approved in current status: %s", association.Status)
	}

	association.Status = domain.AssociationStatusActive
	association.ApprovedByUserID = &approvedByUserID
	now := time.Now()
	association.ApprovedAt = &now

	if err := s.associationRepo.UpdateAssociation(ctx, association); err != nil {
		return nil, fmt.Errorf("error approving association: %w", err)
	}

	return association, nil
}

// RejectAssociation rejects a pending association
func (s *organizationAssociationService) RejectAssociation(ctx context.Context, associationID int64, approvedByUserID string) (*domain.OrganizationAssociation, error) {
	association, err := s.associationRepo.GetAssociationByID(ctx, associationID)
	if err != nil {
		return nil, fmt.Errorf("association not found: %w", err)
	}

	if !association.IsPending() {
		return nil, fmt.Errorf("association cannot be rejected in current status: %s", association.Status)
	}

	association.Status = domain.AssociationStatusRejected
	association.ApprovedByUserID = &approvedByUserID
	now := time.Now()
	association.ApprovedAt = &now

	if err := s.associationRepo.UpdateAssociation(ctx, association); err != nil {
		return nil, fmt.Errorf("error rejecting association: %w", err)
	}

	return association, nil
}

// SuspendAssociation suspends an active association
func (s *organizationAssociationService) SuspendAssociation(ctx context.Context, associationID int64, suspendedByUserID string) (*domain.OrganizationAssociation, error) {
	association, err := s.associationRepo.GetAssociationByID(ctx, associationID)
	if err != nil {
		return nil, fmt.Errorf("association not found: %w", err)
	}

	if !association.CanBeSuspended() {
		return nil, fmt.Errorf("association cannot be suspended in current status: %s", association.Status)
	}

	association.Status = domain.AssociationStatusSuspended
	association.ApprovedByUserID = &suspendedByUserID

	if err := s.associationRepo.UpdateAssociation(ctx, association); err != nil {
		return nil, fmt.Errorf("error suspending association: %w", err)
	}

	return association, nil
}

// ReactivateAssociation reactivates a suspended association
func (s *organizationAssociationService) ReactivateAssociation(ctx context.Context, associationID int64, reactivatedByUserID string) (*domain.OrganizationAssociation, error) {
	association, err := s.associationRepo.GetAssociationByID(ctx, associationID)
	if err != nil {
		return nil, fmt.Errorf("association not found: %w", err)
	}

	if !association.CanBeReactivated() {
		return nil, fmt.Errorf("association cannot be reactivated in current status: %s", association.Status)
	}

	association.Status = domain.AssociationStatusActive
	association.ApprovedByUserID = &reactivatedByUserID

	if err := s.associationRepo.UpdateAssociation(ctx, association); err != nil {
		return nil, fmt.Errorf("error reactivating association: %w", err)
	}

	return association, nil
}

// UpdateVisibility updates the visibility settings of an association
func (s *organizationAssociationService) UpdateVisibility(ctx context.Context, associationID int64, req *domain.UpdateAssociationRequest, updatedByUserID string) (*domain.OrganizationAssociation, error) {
	association, err := s.associationRepo.GetAssociationByID(ctx, associationID)
	if err != nil {
		return nil, fmt.Errorf("association not found: %w", err)
	}

	// Update visible affiliate IDs if provided (including empty list)
	if req.VisibleAffiliateIDs != nil {
		jsonBytes, err := json.Marshal(req.VisibleAffiliateIDs)
		if err != nil {
			return nil, fmt.Errorf("error marshaling visible affiliate IDs: %w", err)
		}
		jsonStr := string(jsonBytes)
		association.VisibleAffiliateIDs = &jsonStr
	}

	// Update visible campaign IDs if provided (including empty list)
	if req.VisibleCampaignIDs != nil {
		jsonBytes, err := json.Marshal(req.VisibleCampaignIDs)
		if err != nil {
			return nil, fmt.Errorf("error marshaling visible campaign IDs: %w", err)
		}
		jsonStr := string(jsonBytes)
		association.VisibleCampaignIDs = &jsonStr
	}

	// Update visibility flags if provided
	if req.AllAffiliatesVisible != nil {
		association.AllAffiliatesVisible = *req.AllAffiliatesVisible
	}
	if req.AllCampaignsVisible != nil {
		association.AllCampaignsVisible = *req.AllCampaignsVisible
	}

	// Update status if provided
	if req.Status != nil {
		association.Status = *req.Status
	}

	if err := s.associationRepo.UpdateAssociation(ctx, association); err != nil {
		return nil, fmt.Errorf("error updating association visibility: %w", err)
	}

	return association, nil
}

// GetAssociationByID retrieves an association by ID
func (s *organizationAssociationService) GetAssociationByID(ctx context.Context, id int64) (*domain.OrganizationAssociation, error) {
	return s.associationRepo.GetAssociationByID(ctx, id)
}

// GetAssociationByIDWithDetails retrieves an association by ID with organization and user details
func (s *organizationAssociationService) GetAssociationByIDWithDetails(ctx context.Context, id int64) (*domain.OrganizationAssociationWithDetails, error) {
	return s.associationRepo.GetAssociationByIDWithDetails(ctx, id)
}

// GetAssociationByOrganizations retrieves an association by organization IDs
func (s *organizationAssociationService) GetAssociationByOrganizations(ctx context.Context, advertiserOrgID, affiliateOrgID int64) (*domain.OrganizationAssociation, error) {
	return s.associationRepo.GetAssociationByOrganizations(ctx, advertiserOrgID, affiliateOrgID)
}

// ListAssociations lists associations based on filter
func (s *organizationAssociationService) ListAssociations(ctx context.Context, filter *domain.AssociationListFilter) ([]*domain.OrganizationAssociation, error) {
	return s.associationRepo.ListAssociations(ctx, filter)
}

// ListAssociationsWithDetails lists associations with organization and user details
func (s *organizationAssociationService) ListAssociationsWithDetails(ctx context.Context, filter *domain.AssociationListFilter) ([]*domain.OrganizationAssociationWithDetails, error) {
	return s.associationRepo.ListAssociationsWithDetails(ctx, filter)
}

// GetAssociationsForAdvertiser retrieves associations for a specific advertiser organization
func (s *organizationAssociationService) GetAssociationsForAdvertiser(ctx context.Context, advertiserOrgID int64, status *domain.AssociationStatus) ([]*domain.OrganizationAssociation, error) {
	return s.associationRepo.GetAssociationsByAdvertiserOrg(ctx, advertiserOrgID, status)
}

// GetAssociationsForAffiliate retrieves associations for a specific affiliate organization
func (s *organizationAssociationService) GetAssociationsForAffiliate(ctx context.Context, affiliateOrgID int64, status *domain.AssociationStatus) ([]*domain.OrganizationAssociation, error) {
	return s.associationRepo.GetAssociationsByAffiliateOrg(ctx, affiliateOrgID, status)
}

// ValidateAssociationAccess validates that a user's organization has access to an association
func (s *organizationAssociationService) ValidateAssociationAccess(ctx context.Context, associationID int64, userOrgID int64) error {
	association, err := s.associationRepo.GetAssociationByID(ctx, associationID)
	if err != nil {
		return fmt.Errorf("association not found: %w", err)
	}

	if association.AdvertiserOrgID != userOrgID && association.AffiliateOrgID != userOrgID {
		return fmt.Errorf("user organization does not have access to this association")
	}

	return nil
}

// CanUserManageAssociation checks if a user can manage a specific association
func (s *organizationAssociationService) CanUserManageAssociation(ctx context.Context, associationID int64, userID string, userOrgID int64) (bool, error) {
	// First validate access
	if err := s.ValidateAssociationAccess(ctx, associationID, userOrgID); err != nil {
		return false, err
	}

	// Get user profile to check role
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return false, fmt.Errorf("invalid user ID format: %w", err)
	}
	
	profile, err := s.profileRepo.GetProfileByID(ctx, userUUID)
	if err != nil {
		return false, fmt.Errorf("user profile not found: %w", err)
	}

	// Check if user belongs to the organization
	if profile.OrganizationID == nil || *profile.OrganizationID != userOrgID {
		return false, nil
	}

	// For now, allow all users in the organization to manage associations
	// This can be extended with role-based checks if needed
	return true, nil
}

// GetVisibleAffiliates returns the list of affiliate IDs visible in an association
func (s *organizationAssociationService) GetVisibleAffiliates(ctx context.Context, associationID int64) ([]int64, error) {
	association, err := s.associationRepo.GetAssociationByID(ctx, associationID)
	if err != nil {
		return nil, fmt.Errorf("association not found: %w", err)
	}

	// If all affiliates are visible OR if no specific affiliate IDs are set, return all affiliate IDs
	if association.AllAffiliatesVisible || association.VisibleAffiliateIDs == nil {
		// Get all affiliates for the affiliate organization
		affiliates, err := s.affiliateRepo.GetAffiliatesByOrganization(ctx, association.AffiliateOrgID)
		if err != nil {
			return nil, fmt.Errorf("error getting affiliates for organization %d: %w", association.AffiliateOrgID, err)
		}
		
		// Extract affiliate IDs
		var affiliateIDs []int64
		for _, affiliate := range affiliates {
			affiliateIDs = append(affiliateIDs, affiliate.AffiliateID)
		}
		
		return affiliateIDs, nil
	}

	// Parse the specific affiliate IDs from JSON
	var affiliateIDs []int64
	if err := json.Unmarshal([]byte(*association.VisibleAffiliateIDs), &affiliateIDs); err != nil {
		return nil, fmt.Errorf("error unmarshaling visible affiliate IDs: %w", err)
	}

	// If the parsed list is empty, return all affiliate IDs
	if len(affiliateIDs) == 0 {
		affiliates, err := s.affiliateRepo.GetAffiliatesByOrganization(ctx, association.AffiliateOrgID)
		if err != nil {
			return nil, fmt.Errorf("error getting affiliates for organization %d: %w", association.AffiliateOrgID, err)
		}
		
		for _, affiliate := range affiliates {
			affiliateIDs = append(affiliateIDs, affiliate.AffiliateID)
		}
	}

	return affiliateIDs, nil
}

// GetVisibleCampaigns returns the list of campaign IDs visible in an association
func (s *organizationAssociationService) GetVisibleCampaigns(ctx context.Context, associationID int64) ([]int64, error) {
	association, err := s.associationRepo.GetAssociationByID(ctx, associationID)
	if err != nil {
		return nil, fmt.Errorf("association not found: %w", err)
	}

	// If all campaigns are visible OR if no specific campaign IDs are set, return all campaign IDs
	if association.AllCampaignsVisible || association.VisibleCampaignIDs == nil {
		// Get all campaigns for the advertiser organization
		campaigns, err := s.campaignRepo.ListCampaignsByOrganization(ctx, association.AdvertiserOrgID, 1000, 0) // Large limit to get all
		if err != nil {
			return nil, fmt.Errorf("error getting campaigns for organization %d: %w", association.AdvertiserOrgID, err)
		}
		
		// Extract campaign IDs
		var campaignIDs []int64
		for _, campaign := range campaigns {
			campaignIDs = append(campaignIDs, campaign.CampaignID)
		}
		
		return campaignIDs, nil
	}

	// Parse the specific campaign IDs from JSON
	var campaignIDs []int64
	if err := json.Unmarshal([]byte(*association.VisibleCampaignIDs), &campaignIDs); err != nil {
		return nil, fmt.Errorf("error unmarshaling visible campaign IDs: %w", err)
	}

	// If the parsed list is empty, return all campaign IDs
	if len(campaignIDs) == 0 {
		campaigns, err := s.campaignRepo.ListCampaignsByOrganization(ctx, association.AdvertiserOrgID, 1000, 0) // Large limit to get all
		if err != nil {
			return nil, fmt.Errorf("error getting campaigns for organization %d: %w", association.AdvertiserOrgID, err)
		}
		
		for _, campaign := range campaigns {
			campaignIDs = append(campaignIDs, campaign.CampaignID)
		}
	}

	return campaignIDs, nil
}

// GetVisibleAffiliatesForAdvertiser gets all affiliates visible to an advertiser organization
func (s *organizationAssociationService) GetVisibleAffiliatesForAdvertiser(ctx context.Context, advertiserOrgID int64, affiliateOrgID *int64) ([]*domain.Affiliate, error) {
	// Get all active associations for the advertiser
	activeStatus := domain.AssociationStatusActive
	filter := &domain.AssociationListFilter{
		AdvertiserOrgID: &advertiserOrgID,
		Status:          &activeStatus,
	}
	
	if affiliateOrgID != nil {
		filter.AffiliateOrgID = affiliateOrgID
	}

	associations, err := s.associationRepo.ListAssociations(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error getting associations: %w", err)
	}

	var allAffiliates []*domain.Affiliate
	
	for _, association := range associations {
		// Get all affiliates from this affiliate organization
		orgAffiliates, err := s.affiliateRepo.GetAffiliatesByOrganization(ctx, association.AffiliateOrgID)
		if err != nil {
			return nil, fmt.Errorf("error getting affiliates for organization %d: %w", association.AffiliateOrgID, err)
		}

		// Filter based on visibility settings
		if association.AllAffiliatesVisible {
			// All affiliates are visible
			allAffiliates = append(allAffiliates, orgAffiliates...)
		} else if association.VisibleAffiliateIDs != nil {
			// Only specific affiliates are visible
			var visibleIDs []int64
			if err := json.Unmarshal([]byte(*association.VisibleAffiliateIDs), &visibleIDs); err != nil {
				return nil, fmt.Errorf("error unmarshaling visible affiliate IDs: %w", err)
			}
			
			// Create a map for quick lookup
			visibleMap := make(map[int64]bool)
			for _, id := range visibleIDs {
				visibleMap[id] = true
			}
			
			// Filter affiliates
			for _, affiliate := range orgAffiliates {
				if visibleMap[affiliate.AffiliateID] {
					allAffiliates = append(allAffiliates, affiliate)
				}
			}
		}
		// If neither all_affiliates_visible nor visible_affiliate_ids is set, no affiliates are visible
	}

	return allAffiliates, nil
}

// GetVisibleCampaignsForAffiliate gets all campaigns visible to an affiliate organization
func (s *organizationAssociationService) GetVisibleCampaignsForAffiliate(ctx context.Context, affiliateOrgID int64, advertiserOrgID *int64) ([]*domain.Campaign, error) {
	// Get all active associations for the affiliate
	activeStatus := domain.AssociationStatusActive
	filter := &domain.AssociationListFilter{
		AffiliateOrgID: &affiliateOrgID,
		Status:         &activeStatus,
	}
	
	if advertiserOrgID != nil {
		filter.AdvertiserOrgID = advertiserOrgID
	}

	associations, err := s.associationRepo.ListAssociations(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error getting associations: %w", err)
	}

	var allCampaigns []*domain.Campaign
	
	for _, association := range associations {
		// Get all campaigns from this advertiser organization
		orgCampaigns, err := s.campaignRepo.ListCampaignsByOrganization(ctx, association.AdvertiserOrgID, 1000, 0) // Large limit to get all
		if err != nil {
			return nil, fmt.Errorf("error getting campaigns for organization %d: %w", association.AdvertiserOrgID, err)
		}

		// Filter based on visibility settings
		if association.AllCampaignsVisible {
			// All campaigns are visible
			allCampaigns = append(allCampaigns, orgCampaigns...)
		} else if association.VisibleCampaignIDs != nil {
			// Only specific campaigns are visible
			var visibleIDs []int64
			if err := json.Unmarshal([]byte(*association.VisibleCampaignIDs), &visibleIDs); err != nil {
				return nil, fmt.Errorf("error unmarshaling visible campaign IDs: %w", err)
			}
			
			// Create a map for quick lookup
			visibleMap := make(map[int64]bool)
			for _, id := range visibleIDs {
				visibleMap[id] = true
			}
			
			// Filter campaigns
			for _, campaign := range orgCampaigns {
				if visibleMap[campaign.CampaignID] {
					allCampaigns = append(allCampaigns, campaign)
				}
			}
		}
		// If neither all_campaigns_visible nor visible_campaign_ids is set, no campaigns are visible
	}

	return allCampaigns, nil
}