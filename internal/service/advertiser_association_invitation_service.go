package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/repository"
)

// AdvertiserAssociationInvitationService defines the interface for invitation business logic
type AdvertiserAssociationInvitationService interface {
	// Invitation management
	CreateInvitation(ctx context.Context, req *domain.CreateInvitationRequest, createdByUserID string) (*domain.AdvertiserAssociationInvitation, error)
	GetInvitationByID(ctx context.Context, id int64) (*domain.AdvertiserAssociationInvitation, error)
	GetInvitationByToken(ctx context.Context, token string) (*domain.AdvertiserAssociationInvitation, error)
	GetInvitationByIDWithDetails(ctx context.Context, id int64) (*domain.AdvertiserAssociationInvitationWithDetails, error)
	GetInvitationByTokenWithDetails(ctx context.Context, token string) (*domain.AdvertiserAssociationInvitationWithDetails, error)
	UpdateInvitation(ctx context.Context, id int64, req *domain.UpdateInvitationRequest, updatedByUserID string) (*domain.AdvertiserAssociationInvitation, error)
	DeleteInvitation(ctx context.Context, id int64, deletedByUserID string) error
	ListInvitations(ctx context.Context, filter *domain.InvitationListFilter) ([]*domain.AdvertiserAssociationInvitation, error)
	ListInvitationsWithDetails(ctx context.Context, filter *domain.InvitationListFilter) ([]*domain.AdvertiserAssociationInvitationWithDetails, error)
	
	// Invitation usage
	UseInvitation(ctx context.Context, req *domain.UseInvitationRequest, usedByUserID string) (*domain.UseInvitationResponse, error)
	GetInvitationUsageHistory(ctx context.Context, invitationID int64, limit int) ([]*domain.InvitationUsageLog, error)
	
	// Utility methods
	GenerateInvitationToken() (string, error)
	ValidateInvitationAccess(ctx context.Context, invitationID int64, userOrgID int64) error
	CanUserManageInvitation(ctx context.Context, invitationID int64, userID string, userOrgID int64) (bool, error)
	GetInvitationsForAdvertiser(ctx context.Context, advertiserOrgID int64, status *domain.InvitationStatus) ([]*domain.AdvertiserAssociationInvitation, error)
	ExpireInvitations(ctx context.Context) (int64, error)
	
	// Link generation
	GenerateInvitationLink(ctx context.Context, invitationID int64, baseURL string) (string, error)
	GenerateInvitationLinkByToken(token string, baseURL string) string
}

// advertiserAssociationInvitationService implements AdvertiserAssociationInvitationService
type advertiserAssociationInvitationService struct {
	invitationRepo  repository.AdvertiserAssociationInvitationRepository
	associationRepo repository.OrganizationAssociationRepository
	orgRepo         repository.OrganizationRepository
	profileRepo     repository.ProfileRepository
	associationService OrganizationAssociationService
}

// NewAdvertiserAssociationInvitationService creates a new invitation service
func NewAdvertiserAssociationInvitationService(
	invitationRepo repository.AdvertiserAssociationInvitationRepository,
	associationRepo repository.OrganizationAssociationRepository,
	orgRepo repository.OrganizationRepository,
	profileRepo repository.ProfileRepository,
	associationService OrganizationAssociationService,
) AdvertiserAssociationInvitationService {
	return &advertiserAssociationInvitationService{
		invitationRepo:     invitationRepo,
		associationRepo:    associationRepo,
		orgRepo:            orgRepo,
		profileRepo:        profileRepo,
		associationService: associationService,
	}
}

// CreateInvitation creates a new invitation
func (s *advertiserAssociationInvitationService) CreateInvitation(ctx context.Context, req *domain.CreateInvitationRequest, createdByUserID string) (*domain.AdvertiserAssociationInvitation, error) {
	// Validate that the advertiser organization exists and is of correct type
	advOrg, err := s.orgRepo.GetOrganizationByID(ctx, req.AdvertiserOrgID)
	if err != nil {
		return nil, fmt.Errorf("advertiser organization not found: %w", err)
	}
	if advOrg.Type != domain.OrganizationTypeAdvertiser {
		return nil, fmt.Errorf("organization %d is not an advertiser organization", req.AdvertiserOrgID)
	}

	// Generate unique invitation token
	token, err := s.GenerateInvitationToken()
	if err != nil {
		return nil, fmt.Errorf("error generating invitation token: %w", err)
	}

	// Convert allowed affiliate org IDs to JSON
	var allowedAffiliateOrgIDs *string
	if len(req.AllowedAffiliateOrgIDs) > 0 {
		jsonBytes, err := json.Marshal(req.AllowedAffiliateOrgIDs)
		if err != nil {
			return nil, fmt.Errorf("error marshaling allowed affiliate org IDs: %w", err)
		}
		jsonStr := string(jsonBytes)
		allowedAffiliateOrgIDs = &jsonStr
	}

	// Convert default visible affiliate IDs to JSON
	var defaultVisibleAffiliateIDs *string
	if len(req.DefaultVisibleAffiliateIDs) > 0 {
		jsonBytes, err := json.Marshal(req.DefaultVisibleAffiliateIDs)
		if err != nil {
			return nil, fmt.Errorf("error marshaling default visible affiliate IDs: %w", err)
		}
		jsonStr := string(jsonBytes)
		defaultVisibleAffiliateIDs = &jsonStr
	}

	// Convert default visible campaign IDs to JSON
	var defaultVisibleCampaignIDs *string
	if len(req.DefaultVisibleCampaignIDs) > 0 {
		jsonBytes, err := json.Marshal(req.DefaultVisibleCampaignIDs)
		if err != nil {
			return nil, fmt.Errorf("error marshaling default visible campaign IDs: %w", err)
		}
		jsonStr := string(jsonBytes)
		defaultVisibleCampaignIDs = &jsonStr
	}

	// Set default visibility settings
	defaultAllAffiliatesVisible := true
	if req.DefaultAllAffiliatesVisible != nil {
		defaultAllAffiliatesVisible = *req.DefaultAllAffiliatesVisible
	}

	defaultAllCampaignsVisible := true
	if req.DefaultAllCampaignsVisible != nil {
		defaultAllCampaignsVisible = *req.DefaultAllCampaignsVisible
	}

	invitation := &domain.AdvertiserAssociationInvitation{
		AdvertiserOrgID:             req.AdvertiserOrgID,
		InvitationToken:             token,
		Name:                        req.Name,
		Description:                 req.Description,
		AllowedAffiliateOrgIDs:      allowedAffiliateOrgIDs,
		MaxUses:                     req.MaxUses,
		CurrentUses:                 0,
		ExpiresAt:                   req.ExpiresAt,
		Status:                      domain.InvitationStatusActive,
		CreatedByUserID:             createdByUserID,
		Message:                     req.Message,
		DefaultAllAffiliatesVisible: defaultAllAffiliatesVisible,
		DefaultAllCampaignsVisible:  defaultAllCampaignsVisible,
		DefaultVisibleAffiliateIDs:  defaultVisibleAffiliateIDs,
		DefaultVisibleCampaignIDs:   defaultVisibleCampaignIDs,
	}

	if err := invitation.Validate(); err != nil {
		return nil, fmt.Errorf("invalid invitation data: %w", err)
	}

	if err := s.invitationRepo.CreateInvitation(ctx, invitation); err != nil {
		return nil, fmt.Errorf("error creating invitation: %w", err)
	}

	return invitation, nil
}

// GetInvitationByID retrieves an invitation by ID
func (s *advertiserAssociationInvitationService) GetInvitationByID(ctx context.Context, id int64) (*domain.AdvertiserAssociationInvitation, error) {
	return s.invitationRepo.GetInvitationByID(ctx, id)
}

// GetInvitationByToken retrieves an invitation by token
func (s *advertiserAssociationInvitationService) GetInvitationByToken(ctx context.Context, token string) (*domain.AdvertiserAssociationInvitation, error) {
	return s.invitationRepo.GetInvitationByToken(ctx, token)
}

// GetInvitationByIDWithDetails retrieves an invitation by ID with additional details
func (s *advertiserAssociationInvitationService) GetInvitationByIDWithDetails(ctx context.Context, id int64) (*domain.AdvertiserAssociationInvitationWithDetails, error) {
	return s.invitationRepo.GetInvitationByIDWithDetails(ctx, id)
}

// GetInvitationByTokenWithDetails retrieves an invitation by token with additional details
func (s *advertiserAssociationInvitationService) GetInvitationByTokenWithDetails(ctx context.Context, token string) (*domain.AdvertiserAssociationInvitationWithDetails, error) {
	return s.invitationRepo.GetInvitationByTokenWithDetails(ctx, token)
}

// UpdateInvitation updates an existing invitation
func (s *advertiserAssociationInvitationService) UpdateInvitation(ctx context.Context, id int64, req *domain.UpdateInvitationRequest, updatedByUserID string) (*domain.AdvertiserAssociationInvitation, error) {
	invitation, err := s.invitationRepo.GetInvitationByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("invitation not found: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		invitation.Name = *req.Name
	}
	if req.Description != nil {
		invitation.Description = req.Description
	}
	if req.MaxUses != nil {
		invitation.MaxUses = req.MaxUses
	}
	if req.ExpiresAt != nil {
		invitation.ExpiresAt = req.ExpiresAt
	}
	if req.Status != nil {
		invitation.Status = *req.Status
	}
	if req.Message != nil {
		invitation.Message = req.Message
	}
	if req.DefaultAllAffiliatesVisible != nil {
		invitation.DefaultAllAffiliatesVisible = *req.DefaultAllAffiliatesVisible
	}
	if req.DefaultAllCampaignsVisible != nil {
		invitation.DefaultAllCampaignsVisible = *req.DefaultAllCampaignsVisible
	}

	// Update allowed affiliate org IDs if provided
	if len(req.AllowedAffiliateOrgIDs) > 0 {
		jsonBytes, err := json.Marshal(req.AllowedAffiliateOrgIDs)
		if err != nil {
			return nil, fmt.Errorf("error marshaling allowed affiliate org IDs: %w", err)
		}
		jsonStr := string(jsonBytes)
		invitation.AllowedAffiliateOrgIDs = &jsonStr
	}

	// Update default visible affiliate IDs if provided
	if len(req.DefaultVisibleAffiliateIDs) > 0 {
		jsonBytes, err := json.Marshal(req.DefaultVisibleAffiliateIDs)
		if err != nil {
			return nil, fmt.Errorf("error marshaling default visible affiliate IDs: %w", err)
		}
		jsonStr := string(jsonBytes)
		invitation.DefaultVisibleAffiliateIDs = &jsonStr
	}

	// Update default visible campaign IDs if provided
	if len(req.DefaultVisibleCampaignIDs) > 0 {
		jsonBytes, err := json.Marshal(req.DefaultVisibleCampaignIDs)
		if err != nil {
			return nil, fmt.Errorf("error marshaling default visible campaign IDs: %w", err)
		}
		jsonStr := string(jsonBytes)
		invitation.DefaultVisibleCampaignIDs = &jsonStr
	}

	if err := invitation.Validate(); err != nil {
		return nil, fmt.Errorf("invalid invitation data: %w", err)
	}

	if err := s.invitationRepo.UpdateInvitation(ctx, invitation); err != nil {
		return nil, fmt.Errorf("error updating invitation: %w", err)
	}

	return invitation, nil
}

// DeleteInvitation deletes an invitation
func (s *advertiserAssociationInvitationService) DeleteInvitation(ctx context.Context, id int64, deletedByUserID string) error {
	// Verify invitation exists
	_, err := s.invitationRepo.GetInvitationByID(ctx, id)
	if err != nil {
		return fmt.Errorf("invitation not found: %w", err)
	}

	return s.invitationRepo.DeleteInvitation(ctx, id)
}

// ListInvitations lists invitations based on filter
func (s *advertiserAssociationInvitationService) ListInvitations(ctx context.Context, filter *domain.InvitationListFilter) ([]*domain.AdvertiserAssociationInvitation, error) {
	return s.invitationRepo.ListInvitations(ctx, filter)
}

// ListInvitationsWithDetails lists invitations with additional details
func (s *advertiserAssociationInvitationService) ListInvitationsWithDetails(ctx context.Context, filter *domain.InvitationListFilter) ([]*domain.AdvertiserAssociationInvitationWithDetails, error) {
	return s.invitationRepo.ListInvitationsWithDetails(ctx, filter)
}

// UseInvitation uses an invitation to create an organization association
func (s *advertiserAssociationInvitationService) UseInvitation(ctx context.Context, req *domain.UseInvitationRequest, usedByUserID string) (*domain.UseInvitationResponse, error) {
	// Get invitation by token
	invitation, err := s.invitationRepo.GetInvitationByToken(ctx, req.InvitationToken)
	if err != nil {
		errorMsg := "Invitation not found"
		return &domain.UseInvitationResponse{
			Success:      false,
			ErrorMessage: &errorMsg,
		}, nil
	}

	// Check if invitation can be used
	if !invitation.CanBeUsed() {
		var errorMsg string
		if invitation.IsExpired() {
			errorMsg = "Invitation has expired"
		} else if invitation.IsUsageLimitReached() {
			errorMsg = "Invitation usage limit has been reached"
		} else {
			errorMsg = "Invitation is not active"
		}
		
		// Log failed usage attempt
		s.logInvitationUsage(ctx, invitation.InvitationID, req.AffiliateOrgID, &usedByUserID, nil, req.IPAddress, req.UserAgent, false, &errorMsg)
		
		return &domain.UseInvitationResponse{
			Success:      false,
			ErrorMessage: &errorMsg,
			Invitation:   invitation,
		}, nil
	}

	// Validate that the affiliate organization exists and is of correct type
	affOrg, err := s.orgRepo.GetOrganizationByID(ctx, req.AffiliateOrgID)
	if err != nil {
		errorMsg := "Affiliate organization not found"
		s.logInvitationUsage(ctx, invitation.InvitationID, req.AffiliateOrgID, &usedByUserID, nil, req.IPAddress, req.UserAgent, false, &errorMsg)
		return &domain.UseInvitationResponse{
			Success:      false,
			ErrorMessage: &errorMsg,
			Invitation:   invitation,
		}, nil
	}
	if affOrg.Type != domain.OrganizationTypeAffiliate {
		errorMsg := "Organization is not an affiliate organization"
		s.logInvitationUsage(ctx, invitation.InvitationID, req.AffiliateOrgID, &usedByUserID, nil, req.IPAddress, req.UserAgent, false, &errorMsg)
		return &domain.UseInvitationResponse{
			Success:      false,
			ErrorMessage: &errorMsg,
			Invitation:   invitation,
		}, nil
	}

	// Check if affiliate is allowed to use this invitation
	if !s.isAffiliateAllowed(invitation, req.AffiliateOrgID) {
		errorMsg := "Your organization is not allowed to use this invitation"
		s.logInvitationUsage(ctx, invitation.InvitationID, req.AffiliateOrgID, &usedByUserID, nil, req.IPAddress, req.UserAgent, false, &errorMsg)
		return &domain.UseInvitationResponse{
			Success:      false,
			ErrorMessage: &errorMsg,
			Invitation:   invitation,
		}, nil
	}

	// Check if association already exists
	existing, err := s.associationRepo.GetAssociationByOrganizations(ctx, invitation.AdvertiserOrgID, req.AffiliateOrgID)
	if err == nil && existing != nil {
		errorMsg := "Association already exists between these organizations"
		s.logInvitationUsage(ctx, invitation.InvitationID, req.AffiliateOrgID, &usedByUserID, &existing.AssociationID, req.IPAddress, req.UserAgent, false, &errorMsg)
		return &domain.UseInvitationResponse{
			Success:      false,
			ErrorMessage: &errorMsg,
			Invitation:   invitation,
			Association:  existing,
		}, nil
	}

	// Create association request using the invitation's default settings
	createAssocReq := &domain.CreateAssociationRequest{
		AdvertiserOrgID:        invitation.AdvertiserOrgID,
		AffiliateOrgID:         req.AffiliateOrgID,
		AssociationType:        domain.AssociationTypeRequest, // Invitation usage creates a request
		Message:                req.Message,
		AllAffiliatesVisible:   &invitation.DefaultAllAffiliatesVisible,
		AllCampaignsVisible:    &invitation.DefaultAllCampaignsVisible,
	}

	// Parse default visible affiliate IDs
	if invitation.DefaultVisibleAffiliateIDs != nil {
		var affiliateIDs []int64
		if err := json.Unmarshal([]byte(*invitation.DefaultVisibleAffiliateIDs), &affiliateIDs); err == nil {
			createAssocReq.VisibleAffiliateIDs = affiliateIDs
		}
	}

	// Parse default visible campaign IDs
	if invitation.DefaultVisibleCampaignIDs != nil {
		var campaignIDs []int64
		if err := json.Unmarshal([]byte(*invitation.DefaultVisibleCampaignIDs), &campaignIDs); err == nil {
			createAssocReq.VisibleCampaignIDs = campaignIDs
		}
	}

	// Create the association
	association, err := s.associationService.CreateRequest(ctx, createAssocReq, usedByUserID)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to create association: %s", err.Error())
		s.logInvitationUsage(ctx, invitation.InvitationID, req.AffiliateOrgID, &usedByUserID, nil, req.IPAddress, req.UserAgent, false, &errorMsg)
		return &domain.UseInvitationResponse{
			Success:      false,
			ErrorMessage: &errorMsg,
			Invitation:   invitation,
		}, nil
	}

	// Increment invitation usage count
	if err := s.invitationRepo.IncrementInvitationUsage(ctx, invitation.InvitationID); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: Failed to increment invitation usage count: %v\n", err)
	}

	// Log successful usage
	s.logInvitationUsage(ctx, invitation.InvitationID, req.AffiliateOrgID, &usedByUserID, &association.AssociationID, req.IPAddress, req.UserAgent, true, nil)

	return &domain.UseInvitationResponse{
		Success:     true,
		Association: association,
		Invitation:  invitation,
	}, nil
}

// GetInvitationUsageHistory retrieves usage history for an invitation
func (s *advertiserAssociationInvitationService) GetInvitationUsageHistory(ctx context.Context, invitationID int64, limit int) ([]*domain.InvitationUsageLog, error) {
	return s.invitationRepo.GetInvitationUsageHistory(ctx, invitationID, limit)
}

// GenerateInvitationToken generates a unique invitation token
func (s *advertiserAssociationInvitationService) GenerateInvitationToken() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("error generating random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// ValidateInvitationAccess validates that a user has access to an invitation
func (s *advertiserAssociationInvitationService) ValidateInvitationAccess(ctx context.Context, invitationID int64, userOrgID int64) error {
	invitation, err := s.invitationRepo.GetInvitationByID(ctx, invitationID)
	if err != nil {
		return fmt.Errorf("invitation not found: %w", err)
	}

	if invitation.AdvertiserOrgID != userOrgID {
		return fmt.Errorf("access denied: invitation belongs to different organization")
	}

	return nil
}

// CanUserManageInvitation checks if a user can manage an invitation
func (s *advertiserAssociationInvitationService) CanUserManageInvitation(ctx context.Context, invitationID int64, userID string, userOrgID int64) (bool, error) {
	err := s.ValidateInvitationAccess(ctx, invitationID, userOrgID)
	return err == nil, err
}

// GetInvitationsForAdvertiser retrieves invitations for a specific advertiser
func (s *advertiserAssociationInvitationService) GetInvitationsForAdvertiser(ctx context.Context, advertiserOrgID int64, status *domain.InvitationStatus) ([]*domain.AdvertiserAssociationInvitation, error) {
	return s.invitationRepo.GetInvitationsByAdvertiser(ctx, advertiserOrgID, status)
}

// ExpireInvitations marks expired invitations as expired
func (s *advertiserAssociationInvitationService) ExpireInvitations(ctx context.Context) (int64, error) {
	return s.invitationRepo.ExpireInvitations(ctx)
}

// GenerateInvitationLink generates a full invitation link
func (s *advertiserAssociationInvitationService) GenerateInvitationLink(ctx context.Context, invitationID int64, baseURL string) (string, error) {
	invitation, err := s.invitationRepo.GetInvitationByID(ctx, invitationID)
	if err != nil {
		return "", fmt.Errorf("invitation not found: %w", err)
	}

	return s.GenerateInvitationLinkByToken(invitation.InvitationToken, baseURL), nil
}

// GenerateInvitationLinkByToken generates a full invitation link by token
func (s *advertiserAssociationInvitationService) GenerateInvitationLinkByToken(token string, baseURL string) string {
	return fmt.Sprintf("%s/invitations/%s", baseURL, token)
}

// Helper methods

// isAffiliateAllowed checks if an affiliate is allowed to use the invitation
func (s *advertiserAssociationInvitationService) isAffiliateAllowed(invitation *domain.AdvertiserAssociationInvitation, affiliateOrgID int64) bool {
	// If no restrictions, any affiliate can use it
	if invitation.AllowedAffiliateOrgIDs == nil {
		return true
	}

	// Parse allowed affiliate org IDs
	var allowedIDs []int64
	if err := json.Unmarshal([]byte(*invitation.AllowedAffiliateOrgIDs), &allowedIDs); err != nil {
		// If parsing fails, assume unrestricted
		return true
	}

	// Check if affiliate is in the allowed list
	for _, id := range allowedIDs {
		if id == affiliateOrgID {
			return true
		}
	}

	return false
}

// logInvitationUsage logs the usage of an invitation
func (s *advertiserAssociationInvitationService) logInvitationUsage(ctx context.Context, invitationID, affiliateOrgID int64, usedByUserID *string, associationID *int64, ipAddress, userAgent *string, success bool, errorMessage *string) {
	usage := &domain.InvitationUsageLog{
		InvitationID:   invitationID,
		AffiliateOrgID: affiliateOrgID,
		UsedByUserID:   usedByUserID,
		AssociationID:  associationID,
		IPAddress:      ipAddress,
		UserAgent:      userAgent,
		Success:        success,
		ErrorMessage:   errorMessage,
	}

	if err := s.invitationRepo.LogInvitationUsage(ctx, usage); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: Failed to log invitation usage: %v\n", err)
	}
}

