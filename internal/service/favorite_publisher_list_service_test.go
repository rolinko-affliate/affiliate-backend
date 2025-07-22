package service

import (
	"context"
	"errors"
	"testing"

	"github.com/affiliate-backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock repository for testing
type MockFavoritePublisherListRepository struct {
	mock.Mock
}

func (m *MockFavoritePublisherListRepository) CreateList(ctx context.Context, list *domain.FavoritePublisherList) error {
	args := m.Called(ctx, list)
	return args.Error(0)
}

func (m *MockFavoritePublisherListRepository) GetListByID(ctx context.Context, listID int64) (*domain.FavoritePublisherList, error) {
	args := m.Called(ctx, listID)
	return args.Get(0).(*domain.FavoritePublisherList), args.Error(1)
}

func (m *MockFavoritePublisherListRepository) GetListsByOrganization(ctx context.Context, organizationID int64) ([]*domain.FavoritePublisherListWithStats, error) {
	args := m.Called(ctx, organizationID)
	return args.Get(0).([]*domain.FavoritePublisherListWithStats), args.Error(1)
}

func (m *MockFavoritePublisherListRepository) UpdateList(ctx context.Context, list *domain.FavoritePublisherList) error {
	args := m.Called(ctx, list)
	return args.Error(0)
}

func (m *MockFavoritePublisherListRepository) DeleteList(ctx context.Context, listID int64) error {
	args := m.Called(ctx, listID)
	return args.Error(0)
}

func (m *MockFavoritePublisherListRepository) AddPublisherToList(ctx context.Context, item *domain.FavoritePublisherListItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockFavoritePublisherListRepository) RemovePublisherFromList(ctx context.Context, listID int64, publisherDomain string) error {
	args := m.Called(ctx, listID, publisherDomain)
	return args.Error(0)
}

func (m *MockFavoritePublisherListRepository) GetListItems(ctx context.Context, listID int64) ([]*domain.FavoritePublisherListItem, error) {
	args := m.Called(ctx, listID)
	return args.Get(0).([]*domain.FavoritePublisherListItem), args.Error(1)
}

func (m *MockFavoritePublisherListRepository) GetListItemsWithPublisherDetails(ctx context.Context, listID int64) ([]*domain.FavoritePublisherListItem, error) {
	args := m.Called(ctx, listID)
	return args.Get(0).([]*domain.FavoritePublisherListItem), args.Error(1)
}

func (m *MockFavoritePublisherListRepository) UpdatePublisherInList(ctx context.Context, listID int64, publisherDomain string, notes *string) error {
	args := m.Called(ctx, listID, publisherDomain, notes)
	return args.Error(0)
}

func (m *MockFavoritePublisherListRepository) IsPublisherInList(ctx context.Context, listID int64, publisherDomain string) (bool, error) {
	args := m.Called(ctx, listID, publisherDomain)
	return args.Bool(0), args.Error(1)
}

func (m *MockFavoritePublisherListRepository) GetListsContainingPublisher(ctx context.Context, organizationID int64, publisherDomain string) ([]*domain.FavoritePublisherList, error) {
	args := m.Called(ctx, organizationID, publisherDomain)
	return args.Get(0).([]*domain.FavoritePublisherList), args.Error(1)
}

// Mock analytics repository for testing
type MockAnalyticsRepository struct {
	mock.Mock
}

func (m *MockAnalyticsRepository) GetPublisherByDomain(ctx context.Context, domainName string) (*domain.AnalyticsPublisher, error) {
	args := m.Called(ctx, domainName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AnalyticsPublisher), args.Error(1)
}

// Implement other required methods (not used in these tests)
func (m *MockAnalyticsRepository) CreateAdvertiser(ctx context.Context, advertiser *domain.AnalyticsAdvertiser) error {
	return nil
}
func (m *MockAnalyticsRepository) GetAdvertiserByID(ctx context.Context, id int64) (*domain.AnalyticsAdvertiser, error) {
	return nil, nil
}
func (m *MockAnalyticsRepository) GetAdvertiserByDomain(ctx context.Context, domainName string) (*domain.AnalyticsAdvertiser, error) {
	return nil, nil
}
func (m *MockAnalyticsRepository) UpdateAdvertiser(ctx context.Context, advertiser *domain.AnalyticsAdvertiser) error {
	return nil
}
func (m *MockAnalyticsRepository) DeleteAdvertiser(ctx context.Context, id int64) error {
	return nil
}
func (m *MockAnalyticsRepository) CreatePublisher(ctx context.Context, publisher *domain.AnalyticsPublisher) error {
	return nil
}
func (m *MockAnalyticsRepository) GetPublisherByID(ctx context.Context, id int64) (*domain.AnalyticsPublisher, error) {
	return nil, nil
}
func (m *MockAnalyticsRepository) UpdatePublisher(ctx context.Context, publisher *domain.AnalyticsPublisher) error {
	return nil
}
func (m *MockAnalyticsRepository) DeletePublisher(ctx context.Context, id int64) error {
	return nil
}
func (m *MockAnalyticsRepository) SearchAdvertisers(ctx context.Context, query string, limit int) ([]domain.AutocompleteResult, error) {
	return nil, nil
}
func (m *MockAnalyticsRepository) SearchPublishers(ctx context.Context, query string, limit int) ([]domain.AutocompleteResult, error) {
	return nil, nil
}
func (m *MockAnalyticsRepository) SearchBoth(ctx context.Context, query string, limit int) ([]domain.AutocompleteResult, error) {
	return nil, nil
}
func (m *MockAnalyticsRepository) AffiliatesSearch(ctx context.Context, domainFilter, country string, partnerDomains []string, verticals []string, limit int, offset int) (*AffiliatesSearchResult, error) {
	return nil, nil
}

func TestFavoritePublisherListService_CreateList(t *testing.T) {
	mockRepo := new(MockFavoritePublisherListRepository)
	mockAnalyticsRepo := new(MockAnalyticsRepository)
	service := NewFavoritePublisherListService(mockRepo, mockAnalyticsRepo)

	ctx := context.Background()
	organizationID := int64(1)

	t.Run("successful creation", func(t *testing.T) {
		req := &domain.CreateFavoritePublisherListRequest{
			Name:        "Test List",
			Description: stringPtr("Test description"),
		}

		mockRepo.On("CreateList", ctx, mock.MatchedBy(func(list *domain.FavoritePublisherList) bool {
			return list.OrganizationID == organizationID &&
				list.Name == req.Name &&
				list.Description == req.Description
		})).Return(nil).Once()

		result, err := service.CreateList(ctx, organizationID, req)
		require.NoError(t, err)
		assert.Equal(t, organizationID, result.OrganizationID)
		assert.Equal(t, req.Name, result.Name)
		assert.Equal(t, req.Description, result.Description)

		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		req := &domain.CreateFavoritePublisherListRequest{
			Name: "", // Invalid empty name
		}

		_, err := service.CreateList(ctx, organizationID, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")
	})
}

func TestFavoritePublisherListService_GetListByID(t *testing.T) {
	mockRepo := new(MockFavoritePublisherListRepository)
	mockAnalyticsRepo := new(MockAnalyticsRepository)
	service := NewFavoritePublisherListService(mockRepo, mockAnalyticsRepo)

	ctx := context.Background()
	organizationID := int64(1)
	listID := int64(100)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedList := &domain.FavoritePublisherList{
			ListID:         listID,
			OrganizationID: organizationID,
			Name:           "Test List",
			Description:    stringPtr("Test description"),
		}

		mockRepo.On("GetListByID", ctx, listID).Return(expectedList, nil).Once()

		result, err := service.GetListByID(ctx, organizationID, listID)
		require.NoError(t, err)
		assert.Equal(t, expectedList, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("list not found", func(t *testing.T) {
		mockRepo.On("GetListByID", ctx, listID).Return((*domain.FavoritePublisherList)(nil), domain.ErrNotFound).Once()

		_, err := service.GetListByID(ctx, organizationID, listID)
		assert.Equal(t, domain.ErrNotFound, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("list belongs to different organization", func(t *testing.T) {
		differentOrgList := &domain.FavoritePublisherList{
			ListID:         listID,
			OrganizationID: 2, // Different organization
			Name:           "Test List",
			Description:    stringPtr("Test description"),
		}

		mockRepo.On("GetListByID", ctx, listID).Return(differentOrgList, nil).Once()

		_, err := service.GetListByID(ctx, organizationID, listID)
		assert.Equal(t, domain.ErrNotFound, err)

		mockRepo.AssertExpectations(t)
	})
}

func TestFavoritePublisherListService_AddPublisherToList(t *testing.T) {
	mockRepo := new(MockFavoritePublisherListRepository)
	mockAnalyticsRepo := new(MockAnalyticsRepository)
	service := NewFavoritePublisherListService(mockRepo, mockAnalyticsRepo)

	ctx := context.Background()
	organizationID := int64(1)
	listID := int64(100)

	t.Run("successful addition", func(t *testing.T) {
		req := &domain.AddPublisherToListRequest{
			PublisherDomain: "example.com",
			Notes:           stringPtr("Test notes"),
		}

		list := &domain.FavoritePublisherList{
			ListID:         listID,
			OrganizationID: organizationID,
			Name:           "Test List",
		}

		publisher := &domain.AnalyticsPublisher{
			ID:     1,
			Domain: "example.com",
		}

		mockRepo.On("GetListByID", ctx, listID).Return(list, nil).Once()
		mockRepo.On("IsPublisherInList", ctx, listID, req.PublisherDomain).Return(false, nil).Once()
		mockAnalyticsRepo.On("GetPublisherByDomain", ctx, req.PublisherDomain).Return(publisher, nil).Once()
		mockRepo.On("AddPublisherToList", ctx, mock.MatchedBy(func(item *domain.FavoritePublisherListItem) bool {
			return item.ListID == listID &&
				item.PublisherDomain == req.PublisherDomain &&
				item.Notes == req.Notes
		})).Return(nil).Once()

		result, err := service.AddPublisherToList(ctx, organizationID, listID, req)
		require.NoError(t, err)
		assert.Equal(t, listID, result.ListID)
		assert.Equal(t, req.PublisherDomain, result.PublisherDomain)
		assert.Equal(t, req.Notes, result.Notes)

		mockRepo.AssertExpectations(t)
		mockAnalyticsRepo.AssertExpectations(t)
	})

	t.Run("publisher already in list", func(t *testing.T) {
		req := &domain.AddPublisherToListRequest{
			PublisherDomain: "example.com",
			Notes:           stringPtr("Test notes"),
		}

		list := &domain.FavoritePublisherList{
			ListID:         listID,
			OrganizationID: organizationID,
			Name:           "Test List",
		}

		mockRepo.On("GetListByID", ctx, listID).Return(list, nil).Once()
		mockRepo.On("IsPublisherInList", ctx, listID, req.PublisherDomain).Return(true, nil).Once()

		_, err := service.AddPublisherToList(ctx, organizationID, listID, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already in the list")

		mockRepo.AssertExpectations(t)
	})

	t.Run("list not found", func(t *testing.T) {
		req := &domain.AddPublisherToListRequest{
			PublisherDomain: "example.com",
			Notes:           stringPtr("Test notes"),
		}

		mockRepo.On("GetListByID", ctx, listID).Return((*domain.FavoritePublisherList)(nil), domain.ErrNotFound).Once()

		_, err := service.AddPublisherToList(ctx, organizationID, listID, req)
		assert.Equal(t, domain.ErrNotFound, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		req := &domain.AddPublisherToListRequest{
			PublisherDomain: "", // Invalid empty domain
		}

		_, err := service.AddPublisherToList(ctx, organizationID, listID, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "publisher_domain is required")
	})
}

func TestFavoritePublisherListService_RemovePublisherFromList(t *testing.T) {
	mockRepo := new(MockFavoritePublisherListRepository)
	mockAnalyticsRepo := new(MockAnalyticsRepository)
	service := NewFavoritePublisherListService(mockRepo, mockAnalyticsRepo)

	ctx := context.Background()
	organizationID := int64(1)
	listID := int64(100)
	publisherDomain := "example.com"

	t.Run("successful removal", func(t *testing.T) {
		list := &domain.FavoritePublisherList{
			ListID:         listID,
			OrganizationID: organizationID,
			Name:           "Test List",
		}

		mockRepo.On("GetListByID", ctx, listID).Return(list, nil).Once()
		mockRepo.On("RemovePublisherFromList", ctx, listID, publisherDomain).Return(nil).Once()

		err := service.RemovePublisherFromList(ctx, organizationID, listID, publisherDomain)
		require.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("list not found", func(t *testing.T) {
		mockRepo.On("GetListByID", ctx, listID).Return((*domain.FavoritePublisherList)(nil), domain.ErrNotFound).Once()

		err := service.RemovePublisherFromList(ctx, organizationID, listID, publisherDomain)
		assert.Equal(t, domain.ErrNotFound, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("publisher not in list", func(t *testing.T) {
		list := &domain.FavoritePublisherList{
			ListID:         listID,
			OrganizationID: organizationID,
			Name:           "Test List",
		}

		mockRepo.On("GetListByID", ctx, listID).Return(list, nil).Once()
		mockRepo.On("RemovePublisherFromList", ctx, listID, publisherDomain).Return(domain.ErrNotFound).Once()

		err := service.RemovePublisherFromList(ctx, organizationID, listID, publisherDomain)
		assert.Equal(t, domain.ErrNotFound, err)

		mockRepo.AssertExpectations(t)
	})
}

func TestFavoritePublisherListService_GetListItems(t *testing.T) {
	mockRepo := new(MockFavoritePublisherListRepository)
	mockAnalyticsRepo := new(MockAnalyticsRepository)
	service := NewFavoritePublisherListService(mockRepo, mockAnalyticsRepo)

	ctx := context.Background()
	organizationID := int64(1)
	listID := int64(100)

	t.Run("get items without details", func(t *testing.T) {
		list := &domain.FavoritePublisherList{
			ListID:         listID,
			OrganizationID: organizationID,
			Name:           "Test List",
		}

		expectedItems := []*domain.FavoritePublisherListItem{
			{
				ItemID:          1,
				ListID:          listID,
				PublisherDomain: "example1.com",
				Notes:           stringPtr("Notes 1"),
			},
			{
				ItemID:          2,
				ListID:          listID,
				PublisherDomain: "example2.com",
				Notes:           stringPtr("Notes 2"),
			},
		}

		mockRepo.On("GetListByID", ctx, listID).Return(list, nil).Once()
		mockRepo.On("GetListItems", ctx, listID).Return(expectedItems, nil).Once()

		result, err := service.GetListItems(ctx, organizationID, listID, false)
		require.NoError(t, err)
		assert.Equal(t, expectedItems, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("get items with details", func(t *testing.T) {
		list := &domain.FavoritePublisherList{
			ListID:         listID,
			OrganizationID: organizationID,
			Name:           "Test List",
		}

		expectedItems := []*domain.FavoritePublisherListItem{
			{
				ItemID:          1,
				ListID:          listID,
				PublisherDomain: "example1.com",
				Notes:           stringPtr("Notes 1"),
				Publisher: &domain.AnalyticsPublisher{
					ID:     1,
					Domain: "example1.com",
				},
			},
		}

		mockRepo.On("GetListByID", ctx, listID).Return(list, nil).Once()
		mockRepo.On("GetListItemsWithPublisherDetails", ctx, listID).Return(expectedItems, nil).Once()

		result, err := service.GetListItems(ctx, organizationID, listID, true)
		require.NoError(t, err)
		assert.Equal(t, expectedItems, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("list not found", func(t *testing.T) {
		mockRepo.On("GetListByID", ctx, listID).Return((*domain.FavoritePublisherList)(nil), domain.ErrNotFound).Once()

		_, err := service.GetListItems(ctx, organizationID, listID, false)
		assert.Equal(t, domain.ErrNotFound, err)

		mockRepo.AssertExpectations(t)
	})
}

func TestFavoritePublisherListService_UpdateList(t *testing.T) {
	mockRepo := new(MockFavoritePublisherListRepository)
	mockAnalyticsRepo := new(MockAnalyticsRepository)
	service := NewFavoritePublisherListService(mockRepo, mockAnalyticsRepo)

	ctx := context.Background()
	organizationID := int64(1)
	listID := int64(100)

	t.Run("successful update", func(t *testing.T) {
		existingList := &domain.FavoritePublisherList{
			ListID:         listID,
			OrganizationID: organizationID,
			Name:           "Original Name",
			Description:    stringPtr("Original description"),
		}

		req := &domain.UpdateFavoritePublisherListRequest{
			Name:        stringPtr("Updated Name"),
			Description: stringPtr("Updated description"),
		}

		mockRepo.On("GetListByID", ctx, listID).Return(existingList, nil).Once()
		mockRepo.On("UpdateList", ctx, mock.MatchedBy(func(list *domain.FavoritePublisherList) bool {
			return list.ListID == listID &&
				list.Name == "Updated Name" &&
				*list.Description == "Updated description"
		})).Return(nil).Once()

		result, err := service.UpdateList(ctx, organizationID, listID, req)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", result.Name)
		assert.Equal(t, "Updated description", *result.Description)

		mockRepo.AssertExpectations(t)
	})

	t.Run("partial update", func(t *testing.T) {
		existingList := &domain.FavoritePublisherList{
			ListID:         listID,
			OrganizationID: organizationID,
			Name:           "Original Name",
			Description:    stringPtr("Original description"),
		}

		req := &domain.UpdateFavoritePublisherListRequest{
			Name: stringPtr("Updated Name"),
			// Description not provided - should remain unchanged
		}

		mockRepo.On("GetListByID", ctx, listID).Return(existingList, nil).Once()
		mockRepo.On("UpdateList", ctx, mock.MatchedBy(func(list *domain.FavoritePublisherList) bool {
			return list.ListID == listID &&
				list.Name == "Updated Name" &&
				*list.Description == "Original description" // Should remain unchanged
		})).Return(nil).Once()

		result, err := service.UpdateList(ctx, organizationID, listID, req)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", result.Name)
		assert.Equal(t, "Original description", *result.Description)

		mockRepo.AssertExpectations(t)
	})
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}
