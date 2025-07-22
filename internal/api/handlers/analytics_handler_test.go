package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/affiliate-backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAnalyticsService is a mock implementation of AnalyticsService
type MockAnalyticsService struct {
	mock.Mock
}

func (m *MockAnalyticsService) GetPublisherByID(ctx context.Context, publisherID int64) (*domain.AnalyticsPublisherResponse, error) {
	args := m.Called(ctx, publisherID)
	return args.Get(0).(*domain.AnalyticsPublisherResponse), args.Error(1)
}

func (m *MockAnalyticsService) GetPublisherByDomain(ctx context.Context, domainName string) (*domain.AnalyticsPublisherResponse, error) {
	args := m.Called(ctx, domainName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AnalyticsPublisherResponse), args.Error(1)
}

func (m *MockAnalyticsService) CreatePublisher(ctx context.Context, req *domain.CreateAnalyticsPublisherRequest) (*domain.AnalyticsPublisherResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*domain.AnalyticsPublisherResponse), args.Error(1)
}

func (m *MockAnalyticsService) GetAdvertiserByID(ctx context.Context, advertiserID int64) (*domain.AnalyticsAdvertiserResponse, error) {
	args := m.Called(ctx, advertiserID)
	return args.Get(0).(*domain.AnalyticsAdvertiserResponse), args.Error(1)
}

func (m *MockAnalyticsService) CreateAdvertiser(ctx context.Context, req *domain.CreateAnalyticsAdvertiserRequest) (*domain.AnalyticsAdvertiserResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*domain.AnalyticsAdvertiserResponse), args.Error(1)
}

func TestAnalyticsHandler_GetPublisherByDomain(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		domain         string
		mockSetup      func(*MockAnalyticsService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful retrieval",
			domain: "example.com",
			mockSetup: func(m *MockAnalyticsService) {
				response := &domain.AnalyticsPublisherResponse{
					Publisher: domain.AnalyticsPublisher{
						Domain:       "example.com",
						Relevance:    85.5,
						TrafficScore: 1250.75,
						Known: domain.BoolValue{
							Value: true,
						},
						MetaData: domain.PublisherMetaData{
							Description:        "Test publisher",
							FaviconImageURL:    "https://example.com/favicon.ico",
							ScreenshotImageURL: "https://example.com/screenshot.png",
						},
					},
				}
				m.On("GetPublisherByDomain", mock.Anything, "example.com").Return(response, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "publisher not found",
			domain: "nonexistent.com",
			mockSetup: func(m *MockAnalyticsService) {
				m.On("GetPublisherByDomain", mock.Anything, "nonexistent.com").Return(nil, domain.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Publisher not found",
		},
		{
			name:   "empty domain parameter",
			domain: "",
			mockSetup: func(m *MockAnalyticsService) {
				// No mock setup needed as validation should fail before service call
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Domain parameter is required",
		},
		{
			name:   "internal server error",
			domain: "error.com",
			mockSetup: func(m *MockAnalyticsService) {
				m.On("GetPublisherByDomain", mock.Anything, "error.com").Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to retrieve publisher",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockAnalyticsService)
			tt.mockSetup(mockService)

			handler := &AnalyticsHandler{
				analyticsService: mockService,
			}

			// Create request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{
				{Key: "domain", Value: tt.domain},
			}

			// Execute
			handler.GetPublisherByDomain(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, "Publisher retrieved successfully", response["message"])
				assert.NotNil(t, response["data"])

				data := response["data"].(map[string]interface{})
				publisher := data["publisher"].(map[string]interface{})
				assert.Equal(t, tt.domain, publisher["domain"])
			} else {
				assert.Contains(t, response["error"].(string), tt.expectedError)
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestAnalyticsHandler_GetPublisherByDomain_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// This test demonstrates the full request flow
	mockService := new(MockAnalyticsService)

	// Setup mock response
	expectedResponse := &domain.AnalyticsPublisherResponse{
		Publisher: domain.AnalyticsPublisher{
			Domain:       "integration-test.com",
			Relevance:    92.3,
			TrafficScore: 2500.50,
			Known: domain.BoolValue{
				Value: true,
			},
			MetaData: domain.PublisherMetaData{
				Description:        "Integration test publisher",
				FaviconImageURL:    "https://integration-test.com/favicon.ico",
				ScreenshotImageURL: "https://integration-test.com/screenshot.png",
			},
		},
	}

	mockService.On("GetPublisherByDomain", mock.Anything, "integration-test.com").Return(expectedResponse, nil)

	handler := &AnalyticsHandler{
		analyticsService: mockService,
	}

	// Create HTTP request
	req, _ := http.NewRequest("GET", "/api/v1/analytics/affiliates/domain/integration-test.com", nil)
	w := httptest.NewRecorder()

	// Setup Gin context
	router := gin.New()
	router.GET("/api/v1/analytics/affiliates/domain/:domain", handler.GetPublisherByDomain)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Publisher retrieved successfully", response["message"])

	data := response["data"].(map[string]interface{})
	publisher := data["publisher"].(map[string]interface{})

	assert.Equal(t, "integration-test.com", publisher["domain"])
	assert.Equal(t, 92.3, publisher["relevance"])
	assert.Equal(t, 2500.50, publisher["trafficScore"])

	known := publisher["known"].(map[string]interface{})
	assert.Equal(t, true, known["value"])

	metaData := publisher["metaData"].(map[string]interface{})
	assert.Equal(t, "Integration test publisher", metaData["description"])
	assert.Equal(t, "https://integration-test.com/favicon.ico", metaData["faviconImageUrl"])

	mockService.AssertExpectations(t)
}
