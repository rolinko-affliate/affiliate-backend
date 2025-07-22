package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/affiliate-backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock service for testing
type MockFavoritePublisherListService struct {
	mock.Mock
}

func (m *MockFavoritePublisherListService) CreateList(ctx context.Context, organizationID int64, req *domain.CreateFavoritePublisherListRequest) (*domain.FavoritePublisherList, error) {
	args := m.Called(ctx, organizationID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FavoritePublisherList), args.Error(1)
}

func (m *MockFavoritePublisherListService) GetListByID(ctx context.Context, organizationID int64, listID int64) (*domain.FavoritePublisherList, error) {
	args := m.Called(ctx, organizationID, listID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FavoritePublisherList), args.Error(1)
}

func (m *MockFavoritePublisherListService) GetListsByOrganization(ctx context.Context, organizationID int64) ([]*domain.FavoritePublisherListWithStats, error) {
	args := m.Called(ctx, organizationID)
	return args.Get(0).([]*domain.FavoritePublisherListWithStats), args.Error(1)
}

func (m *MockFavoritePublisherListService) UpdateList(ctx context.Context, organizationID int64, listID int64, req *domain.UpdateFavoritePublisherListRequest) (*domain.FavoritePublisherList, error) {
	args := m.Called(ctx, organizationID, listID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FavoritePublisherList), args.Error(1)
}

func (m *MockFavoritePublisherListService) DeleteList(ctx context.Context, organizationID int64, listID int64) error {
	args := m.Called(ctx, organizationID, listID)
	return args.Error(0)
}

func (m *MockFavoritePublisherListService) AddPublisherToList(ctx context.Context, organizationID int64, listID int64, req *domain.AddPublisherToListRequest) (*domain.FavoritePublisherListItem, error) {
	args := m.Called(ctx, organizationID, listID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FavoritePublisherListItem), args.Error(1)
}

func (m *MockFavoritePublisherListService) RemovePublisherFromList(ctx context.Context, organizationID int64, listID int64, publisherDomain string) error {
	args := m.Called(ctx, organizationID, listID, publisherDomain)
	return args.Error(0)
}

func (m *MockFavoritePublisherListService) GetListItems(ctx context.Context, organizationID int64, listID int64, includeDetails bool) ([]*domain.FavoritePublisherListItem, error) {
	args := m.Called(ctx, organizationID, listID, includeDetails)
	return args.Get(0).([]*domain.FavoritePublisherListItem), args.Error(1)
}

func (m *MockFavoritePublisherListService) UpdatePublisherInList(ctx context.Context, organizationID int64, listID int64, publisherDomain string, req *domain.UpdatePublisherInListRequest) error {
	args := m.Called(ctx, organizationID, listID, publisherDomain, req)
	return args.Error(0)
}

func (m *MockFavoritePublisherListService) GetListsContainingPublisher(ctx context.Context, organizationID int64, publisherDomain string) ([]*domain.FavoritePublisherList, error) {
	args := m.Called(ctx, organizationID, publisherDomain)
	return args.Get(0).([]*domain.FavoritePublisherList), args.Error(1)
}

func setupFavoritePublisherListHandler() (*MockFavoritePublisherListService, *FavoritePublisherListHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockFavoritePublisherListService)
	handler := NewFavoritePublisherListHandler(mockService)
	router := gin.New()
	return mockService, handler, router
}

func TestFavoritePublisherListHandler_CreateList(t *testing.T) {
	mockService, handler, router := setupFavoritePublisherListHandler()

	router.POST("/favorite-publisher-lists", func(c *gin.Context) {
		c.Set("organizationID", int64(1))
		handler.CreateList(c)
	})

	t.Run("successful creation", func(t *testing.T) {
		req := domain.CreateFavoritePublisherListRequest{
			Name:        "Test List",
			Description: stringPtr("Test description"),
		}

		expectedList := &domain.FavoritePublisherList{
			ListID:         1,
			OrganizationID: 1,
			Name:           req.Name,
			Description:    req.Description,
		}

		mockService.On("CreateList", mock.Anything, int64(1), &req).Return(expectedList, nil).Once()

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/favorite-publisher-lists", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Favorite publisher list created successfully", response["message"])
		assert.NotNil(t, response["data"])

		mockService.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/favorite-publisher-lists", bytes.NewBuffer([]byte("invalid json")))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request body", response.Error)
	})

	t.Run("validation error", func(t *testing.T) {
		req := domain.CreateFavoritePublisherListRequest{
			Name: "", // Invalid empty name
		}

		mockService.On("CreateList", mock.Anything, int64(1), &req).Return(nil, domain.ErrInvalidInput).Once()

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/favorite-publisher-lists", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Invalid input", response.Error)

		mockService.AssertExpectations(t)
	})

	t.Run("missing organizationID", func(t *testing.T) {
		router.POST("/test", handler.CreateList) // No middleware to set organizationID

		req := domain.CreateFavoritePublisherListRequest{
			Name: "Test List",
		}

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Unauthorized", response.Error)
	})
}

func TestFavoritePublisherListHandler_GetLists(t *testing.T) {
	mockService, handler, router := setupFavoritePublisherListHandler()

	router.GET("/favorite-publisher-lists", func(c *gin.Context) {
		c.Set("organizationID", int64(1))
		handler.GetLists(c)
	})

	t.Run("successful retrieval", func(t *testing.T) {
		expectedLists := []*domain.FavoritePublisherListWithStats{
			{
				FavoritePublisherList: domain.FavoritePublisherList{
					ListID:         1,
					OrganizationID: 1,
					Name:           "List 1",
					Description:    stringPtr("Description 1"),
				},
				PublisherCount: 5,
			},
			{
				FavoritePublisherList: domain.FavoritePublisherList{
					ListID:         2,
					OrganizationID: 1,
					Name:           "List 2",
					Description:    stringPtr("Description 2"),
				},
				PublisherCount: 3,
			},
		}

		mockService.On("GetListsByOrganization", mock.Anything, int64(1)).Return(expectedLists, nil).Once()

		w := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/favorite-publisher-lists", nil)

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusOK, w.Code)

		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Favorite publisher lists retrieved successfully", response["message"])
		assert.NotNil(t, response["data"])

		mockService.AssertExpectations(t)
	})
}

func TestFavoritePublisherListHandler_GetListByID(t *testing.T) {
	mockService, handler, router := setupFavoritePublisherListHandler()

	router.GET("/favorite-publisher-lists/:list_id", func(c *gin.Context) {
		c.Set("organizationID", int64(1))
		handler.GetListByID(c)
	})

	t.Run("successful retrieval", func(t *testing.T) {
		expectedList := &domain.FavoritePublisherList{
			ListID:         1,
			OrganizationID: 1,
			Name:           "Test List",
			Description:    stringPtr("Test description"),
		}

		mockService.On("GetListByID", mock.Anything, int64(1), int64(1)).Return(expectedList, nil).Once()

		w := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/favorite-publisher-lists/1", nil)

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusOK, w.Code)

		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Favorite publisher list retrieved successfully", response["message"])
		assert.NotNil(t, response["data"])

		mockService.AssertExpectations(t)
	})

	t.Run("invalid list ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/favorite-publisher-lists/invalid", nil)

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Invalid list ID", response.Error)
	})

	t.Run("list not found", func(t *testing.T) {
		mockService.On("GetListByID", mock.Anything, int64(1), int64(999)).Return(nil, domain.ErrNotFound).Once()

		w := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/favorite-publisher-lists/999", nil)

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "List not found", response.Error)

		mockService.AssertExpectations(t)
	})
}

func TestFavoritePublisherListHandler_AddPublisherToList(t *testing.T) {
	mockService, handler, router := setupFavoritePublisherListHandler()

	router.POST("/favorite-publisher-lists/:list_id/publishers", func(c *gin.Context) {
		c.Set("organizationID", int64(1))
		handler.AddPublisherToList(c)
	})

	t.Run("successful addition", func(t *testing.T) {
		req := domain.AddPublisherToListRequest{
			PublisherDomain: "example.com",
			Notes:           stringPtr("Test notes"),
		}

		expectedItem := &domain.FavoritePublisherListItem{
			ItemID:          1,
			ListID:          1,
			PublisherDomain: req.PublisherDomain,
			Notes:           req.Notes,
		}

		mockService.On("AddPublisherToList", mock.Anything, int64(1), int64(1), &req).Return(expectedItem, nil).Once()

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/favorite-publisher-lists/1/publishers", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Publisher added to list successfully", response["message"])
		assert.NotNil(t, response["data"])

		mockService.AssertExpectations(t)
	})

	t.Run("publisher already in list", func(t *testing.T) {
		req := domain.AddPublisherToListRequest{
			PublisherDomain: "example.com",
			Notes:           stringPtr("Test notes"),
		}

		mockService.On("AddPublisherToList", mock.Anything, int64(1), int64(1), &req).Return(nil,
			errors.New("publisher example.com is already in the list")).Once()

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/favorite-publisher-lists/1/publishers", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusConflict, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Publisher already in list", response.Error)

		mockService.AssertExpectations(t)
	})
}

func TestFavoritePublisherListHandler_RemovePublisherFromList(t *testing.T) {
	mockService, handler, router := setupFavoritePublisherListHandler()

	router.DELETE("/favorite-publisher-lists/:list_id/publishers/:domain", func(c *gin.Context) {
		c.Set("organizationID", int64(1))
		handler.RemovePublisherFromList(c)
	})

	t.Run("successful removal", func(t *testing.T) {
		mockService.On("RemovePublisherFromList", mock.Anything, int64(1), int64(1), "example.com").Return(nil).Once()

		w := httptest.NewRecorder()
		request, _ := http.NewRequest("DELETE", "/favorite-publisher-lists/1/publishers/example.com", nil)

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusOK, w.Code)

		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Publisher removed from list successfully", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("publisher not found", func(t *testing.T) {
		mockService.On("RemovePublisherFromList", mock.Anything, int64(1), int64(1), "nonexistent.com").Return(domain.ErrNotFound).Once()

		w := httptest.NewRecorder()
		request, _ := http.NewRequest("DELETE", "/favorite-publisher-lists/1/publishers/nonexistent.com", nil)

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Publisher not found in list", response.Error)

		mockService.AssertExpectations(t)
	})
}

func TestFavoritePublisherListHandler_GetListItems(t *testing.T) {
	mockService, handler, router := setupFavoritePublisherListHandler()

	router.GET("/favorite-publisher-lists/:list_id/publishers", func(c *gin.Context) {
		c.Set("organizationID", int64(1))
		handler.GetListItems(c)
	})

	t.Run("get items without details", func(t *testing.T) {
		expectedItems := []*domain.FavoritePublisherListItem{
			{
				ItemID:          1,
				ListID:          1,
				PublisherDomain: "example1.com",
				Notes:           stringPtr("Notes 1"),
			},
			{
				ItemID:          2,
				ListID:          1,
				PublisherDomain: "example2.com",
				Notes:           stringPtr("Notes 2"),
			},
		}

		mockService.On("GetListItems", mock.Anything, int64(1), int64(1), false).Return(expectedItems, nil).Once()

		w := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/favorite-publisher-lists/1/publishers", nil)

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusOK, w.Code)

		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "List items retrieved successfully", response["message"])
		assert.NotNil(t, response["data"])

		mockService.AssertExpectations(t)
	})

	t.Run("get items with details", func(t *testing.T) {
		expectedItems := []*domain.FavoritePublisherListItem{
			{
				ItemID:          1,
				ListID:          1,
				PublisherDomain: "example1.com",
				Notes:           stringPtr("Notes 1"),
				Publisher: &domain.AnalyticsPublisher{
					ID:     1,
					Domain: "example1.com",
				},
			},
		}

		mockService.On("GetListItems", mock.Anything, int64(1), int64(1), true).Return(expectedItems, nil).Once()

		w := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/favorite-publisher-lists/1/publishers?include_details=true", nil)

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusOK, w.Code)

		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "List items retrieved successfully", response["message"])
		assert.NotNil(t, response["data"])

		mockService.AssertExpectations(t)
	})
}

func TestFavoritePublisherListHandler_GetListsContainingPublisher(t *testing.T) {
	mockService, handler, router := setupFavoritePublisherListHandler()

	router.GET("/favorite-publisher-lists/search", func(c *gin.Context) {
		c.Set("organizationID", int64(1))
		handler.GetListsContainingPublisher(c)
	})

	t.Run("successful search", func(t *testing.T) {
		expectedLists := []*domain.FavoritePublisherList{
			{
				ListID:         1,
				OrganizationID: 1,
				Name:           "List 1",
				Description:    stringPtr("Description 1"),
			},
		}

		mockService.On("GetListsContainingPublisher", mock.Anything, int64(1), "example.com").Return(expectedLists, nil).Once()

		w := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/favorite-publisher-lists/search?domain=example.com", nil)

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusOK, w.Code)

		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Lists containing publisher retrieved successfully", response["message"])
		assert.NotNil(t, response["data"])

		mockService.AssertExpectations(t)
	})

	t.Run("missing domain parameter", func(t *testing.T) {
		w := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/favorite-publisher-lists/search", nil)

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Invalid domain", response.Error)
	})
}


