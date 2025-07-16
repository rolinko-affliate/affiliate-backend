package service

import (
	"context"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/repository"
)

// FavoritePublisherListService defines the interface for favorite publisher list business logic
type FavoritePublisherListService interface {
	// List management
	CreateList(ctx context.Context, organizationID int64, req *domain.CreateFavoritePublisherListRequest) (*domain.FavoritePublisherList, error)
	GetListByID(ctx context.Context, organizationID int64, listID int64) (*domain.FavoritePublisherList, error)
	GetListsByOrganization(ctx context.Context, organizationID int64) ([]*domain.FavoritePublisherListWithStats, error)
	UpdateList(ctx context.Context, organizationID int64, listID int64, req *domain.UpdateFavoritePublisherListRequest) (*domain.FavoritePublisherList, error)
	DeleteList(ctx context.Context, organizationID int64, listID int64) error
	
	// List item management
	AddPublisherToList(ctx context.Context, organizationID int64, listID int64, req *domain.AddPublisherToListRequest) (*domain.FavoritePublisherListItem, error)
	RemovePublisherFromList(ctx context.Context, organizationID int64, listID int64, publisherDomain string) error
	GetListItems(ctx context.Context, organizationID int64, listID int64, includeDetails bool) ([]*domain.FavoritePublisherListItem, error)
	UpdatePublisherInList(ctx context.Context, organizationID int64, listID int64, publisherDomain string, req *domain.UpdatePublisherInListRequest) error
	
	// Utility methods
	GetListsContainingPublisher(ctx context.Context, organizationID int64, publisherDomain string) ([]*domain.FavoritePublisherList, error)
}

// favoritePublisherListService implements FavoritePublisherListService
type favoritePublisherListService struct {
	favoriteListRepo repository.FavoritePublisherListRepository
	analyticsRepo    repository.AnalyticsRepository
}

// NewFavoritePublisherListService creates a new favorite publisher list service
func NewFavoritePublisherListService(favoriteListRepo repository.FavoritePublisherListRepository, analyticsRepo repository.AnalyticsRepository) FavoritePublisherListService {
	return &favoritePublisherListService{
		favoriteListRepo: favoriteListRepo,
		analyticsRepo:    analyticsRepo,
	}
}

// CreateList creates a new favorite publisher list
func (s *favoritePublisherListService) CreateList(ctx context.Context, organizationID int64, req *domain.CreateFavoritePublisherListRequest) (*domain.FavoritePublisherList, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	
	list := &domain.FavoritePublisherList{
		OrganizationID: organizationID,
		Name:           req.Name,
		Description:    req.Description,
	}
	
	err := s.favoriteListRepo.CreateList(ctx, list)
	if err != nil {
		return nil, fmt.Errorf("failed to create favorite publisher list: %w", err)
	}
	
	return list, nil
}

// GetListByID retrieves a favorite publisher list by ID, ensuring it belongs to the organization
func (s *favoritePublisherListService) GetListByID(ctx context.Context, organizationID int64, listID int64) (*domain.FavoritePublisherList, error) {
	list, err := s.favoriteListRepo.GetListByID(ctx, listID)
	if err != nil {
		return nil, err
	}
	
	// Ensure the list belongs to the organization
	if list.OrganizationID != organizationID {
		return nil, domain.ErrNotFound
	}
	
	return list, nil
}

// GetListsByOrganization retrieves all favorite publisher lists for an organization
func (s *favoritePublisherListService) GetListsByOrganization(ctx context.Context, organizationID int64) ([]*domain.FavoritePublisherListWithStats, error) {
	return s.favoriteListRepo.GetListsByOrganization(ctx, organizationID)
}

// UpdateList updates a favorite publisher list
func (s *favoritePublisherListService) UpdateList(ctx context.Context, organizationID int64, listID int64, req *domain.UpdateFavoritePublisherListRequest) (*domain.FavoritePublisherList, error) {
	// First, get the existing list to ensure it belongs to the organization
	list, err := s.GetListByID(ctx, organizationID, listID)
	if err != nil {
		return nil, err
	}
	
	// Update fields if provided
	if req.Name != nil {
		list.Name = *req.Name
	}
	if req.Description != nil {
		list.Description = req.Description
	}
	
	err = s.favoriteListRepo.UpdateList(ctx, list)
	if err != nil {
		return nil, fmt.Errorf("failed to update favorite publisher list: %w", err)
	}
	
	return list, nil
}

// DeleteList deletes a favorite publisher list
func (s *favoritePublisherListService) DeleteList(ctx context.Context, organizationID int64, listID int64) error {
	// First, ensure the list belongs to the organization
	_, err := s.GetListByID(ctx, organizationID, listID)
	if err != nil {
		return err
	}
	
	return s.favoriteListRepo.DeleteList(ctx, listID)
}

// AddPublisherToList adds a publisher to a favorite list
func (s *favoritePublisherListService) AddPublisherToList(ctx context.Context, organizationID int64, listID int64, req *domain.AddPublisherToListRequest) (*domain.FavoritePublisherListItem, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	
	// First, ensure the list belongs to the organization
	_, err := s.GetListByID(ctx, organizationID, listID)
	if err != nil {
		return nil, err
	}
	
	// Check if publisher is already in the list
	exists, err := s.favoriteListRepo.IsPublisherInList(ctx, listID, req.PublisherDomain)
	if err != nil {
		return nil, fmt.Errorf("failed to check if publisher exists in list: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("publisher %s is already in the list", req.PublisherDomain)
	}
	
	// Optionally validate that the publisher domain exists in analytics_publishers
	// This is a soft validation - we'll still allow adding domains that don't exist yet
	_, err = s.analyticsRepo.GetPublisherByDomain(ctx, req.PublisherDomain)
	if err != nil && err != domain.ErrNotFound {
		return nil, fmt.Errorf("failed to validate publisher domain: %w", err)
	}
	
	item := &domain.FavoritePublisherListItem{
		ListID:          listID,
		PublisherDomain: req.PublisherDomain,
		Notes:           req.Notes,
	}
	
	err = s.favoriteListRepo.AddPublisherToList(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("failed to add publisher to list: %w", err)
	}
	
	return item, nil
}

// RemovePublisherFromList removes a publisher from a favorite list
func (s *favoritePublisherListService) RemovePublisherFromList(ctx context.Context, organizationID int64, listID int64, publisherDomain string) error {
	// First, ensure the list belongs to the organization
	_, err := s.GetListByID(ctx, organizationID, listID)
	if err != nil {
		return err
	}
	
	return s.favoriteListRepo.RemovePublisherFromList(ctx, listID, publisherDomain)
}

// GetListItems retrieves all items in a favorite list
func (s *favoritePublisherListService) GetListItems(ctx context.Context, organizationID int64, listID int64, includeDetails bool) ([]*domain.FavoritePublisherListItem, error) {
	// First, ensure the list belongs to the organization
	_, err := s.GetListByID(ctx, organizationID, listID)
	if err != nil {
		return nil, err
	}
	
	if includeDetails {
		return s.favoriteListRepo.GetListItemsWithPublisherDetails(ctx, listID)
	}
	
	return s.favoriteListRepo.GetListItems(ctx, listID)
}

// UpdatePublisherInList updates the notes for a publisher in a list
func (s *favoritePublisherListService) UpdatePublisherInList(ctx context.Context, organizationID int64, listID int64, publisherDomain string, req *domain.UpdatePublisherInListRequest) error {
	// First, ensure the list belongs to the organization
	_, err := s.GetListByID(ctx, organizationID, listID)
	if err != nil {
		return err
	}
	
	return s.favoriteListRepo.UpdatePublisherInList(ctx, listID, publisherDomain, req.Notes)
}

// GetListsContainingPublisher retrieves all lists that contain a specific publisher for an organization
func (s *favoritePublisherListService) GetListsContainingPublisher(ctx context.Context, organizationID int64, publisherDomain string) ([]*domain.FavoritePublisherList, error) {
	return s.favoriteListRepo.GetListsContainingPublisher(ctx, organizationID, publisherDomain)
}