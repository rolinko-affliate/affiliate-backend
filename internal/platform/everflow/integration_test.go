package everflow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/everflow/affiliate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockAffiliateProviderMappingRepository is a mock implementation of the repository
type MockAffiliateProviderMappingRepository struct {
	mock.Mock
}

func (m *MockAffiliateProviderMappingRepository) GetAffiliateProviderMapping(ctx context.Context, affiliateID int64, providerType string) (*domain.AffiliateProviderMapping, error) {
	args := m.Called(ctx, affiliateID, providerType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AffiliateProviderMapping), args.Error(1)
}

func (m *MockAffiliateProviderMappingRepository) CreateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error {
	args := m.Called(ctx, mapping)
	return args.Error(0)
}

func (m *MockAffiliateProviderMappingRepository) UpdateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error {
	args := m.Called(ctx, mapping)
	return args.Error(0)
}

func (m *MockAffiliateProviderMappingRepository) ListAffiliateProviderMappings(ctx context.Context, affiliateID int64) ([]*domain.AffiliateProviderMapping, error) {
	args := m.Called(ctx, affiliateID)
	return args.Get(0).([]*domain.AffiliateProviderMapping), args.Error(1)
}

func (m *MockAffiliateProviderMappingRepository) DeleteAffiliateProviderMapping(ctx context.Context, affiliateID int64, providerType string) error {
	args := m.Called(ctx, affiliateID, providerType)
	return args.Error(0)
}

// TestEverflowIntegrationService_CreateAffiliate tests the complete flow of creating an affiliate
func TestEverflowIntegrationService_CreateAffiliate(t *testing.T) {
	tests := []struct {
		name                    string
		inputAffiliate          domain.Affiliate
		existingMapping         *domain.AffiliateProviderMapping
		everflowResponse        *affiliate.Affiliate
		everflowStatusCode      int
		expectError             bool
		expectedErrorContains   string
		expectMappingCreate     bool
		expectMappingUpdate     bool
		validateMapping         func(t *testing.T, mapping *domain.AffiliateProviderMapping)
		validateEverflowRequest func(t *testing.T, req *affiliate.CreateAffiliateRequest)
	}{
		{
			name: "successful_affiliate_creation_new_mapping",
			inputAffiliate: domain.Affiliate{
				AffiliateID:            1,
				OrganizationID:         1,
				Name:                   "Example Affiliate",
				ContactEmail:           stringPtr("affiliate@example.com"),
				PaymentDetails:         stringPtr(`{"bank_account":"123456789","routing_number":"987654321","payment_method":"bank_transfer","account_name":"Example Affiliate LLC"}`),
				Status:                 "active",
				InternalNotes:          stringPtr("High-performing affiliate with good conversion rates"),
				DefaultCurrencyID:      stringPtr("USD"),
				ContactAddress:         stringPtr(`{"address1":"123 Main St","city":"New York","region_code":"NY","country_code":"US","zip_postal_code":"10001"}`),
				BillingInfo:            stringPtr(`{"billing_frequency":"monthly","payment_type":"none","tax_id":"12-3456789","details":{"day_of_month":1}}`),
				Labels:                 stringPtr(`["premium","high-volume","trusted"]`),
				InvoiceAmountThreshold: float64Ptr(1000.00),
				DefaultPaymentTerms:    int32Ptr(30),
			},
			existingMapping: nil,
			everflowResponse: &affiliate.Affiliate{
				NetworkAffiliateId:            int32Ptr(12345),
				Name:                          stringPtr("Example Affiliate"),
				AccountStatus:                 stringPtr("active"),
				DefaultCurrencyId:             stringPtr("USD"),
				InternalNotes:                 stringPtr("High-performing affiliate with good conversion rates"),
				NetworkEmployeeId:             int32Ptr(1),
				ReferrerId:                    int32Ptr(0),
				EnableMediaCostTrackingLinks:  boolPtr(false),
				IsContactAddressEnabled:       boolPtr(true),
				NetworkId:                     int32Ptr(1),
				AccountManagerId:              int32Ptr(100),
				AccountManagerName:            stringPtr("John Doe"),
				TimeCreated:                   int64Ptr(1704067200), // 2024-01-01T00:00:00Z
				TimeSaved:                     int64Ptr(1704067200), // 2024-01-01T00:00:00Z
			},
			everflowStatusCode:  200,
			expectError:         false,
			expectMappingCreate: true,
			expectMappingUpdate: false,
			validateMapping: func(t *testing.T, mapping *domain.AffiliateProviderMapping) {
				assert.Equal(t, int64(1), mapping.AffiliateID)
				assert.Equal(t, "everflow", mapping.ProviderType)
				assert.NotNil(t, mapping.SyncStatus)
				assert.Equal(t, "synced", *mapping.SyncStatus)
				assert.NotNil(t, mapping.LastSyncAt)
				assert.NotNil(t, mapping.ProviderConfig)

				// Validate provider data
				assert.NotNil(t, mapping.ProviderData)
				var everflowData domain.EverflowProviderData
				err := json.Unmarshal([]byte(*mapping.ProviderData), &everflowData)
				assert.NoError(t, err)
				assert.NotNil(t, everflowData.NetworkAffiliateID)
				assert.Equal(t, int32(12345), *everflowData.NetworkAffiliateID)
				assert.NotNil(t, everflowData.NetworkEmployeeID)
				assert.Equal(t, int32(1), *everflowData.NetworkEmployeeID)

				// Validate provider config contains request and response
				var config map[string]interface{}
				err = json.Unmarshal([]byte(*mapping.ProviderConfig), &config)
				assert.NoError(t, err)
				assert.Contains(t, config, "request")
				assert.Contains(t, config, "response")
				assert.Contains(t, config, "last_operation")
				assert.Equal(t, "create", config["last_operation"])
			},
			validateEverflowRequest: func(t *testing.T, req *affiliate.CreateAffiliateRequest) {
				assert.Equal(t, "Example Affiliate", req.GetName())
				assert.Equal(t, "active", req.GetAccountStatus())
				assert.Equal(t, "USD", req.GetDefaultCurrencyId())
				assert.Equal(t, "High-performing affiliate with good conversion rates", req.GetInternalNotes())
				assert.Equal(t, int32(1), req.GetNetworkEmployeeId())
				assert.Equal(t, int32(0), req.GetReferrerId())
				assert.False(t, req.GetEnableMediaCostTrackingLinks())

				// Validate contact address
				if req.HasContactAddress() {
					contactAddr := req.GetContactAddress()
					assert.Equal(t, "123 Main St", contactAddr.GetAddress1())
					assert.Equal(t, "New York", contactAddr.GetCity())
					assert.Equal(t, "NY", contactAddr.GetRegionCode())
					assert.Equal(t, "US", contactAddr.GetCountryCode())
					assert.Equal(t, "10001", contactAddr.GetZipPostalCode())
				}

				// Validate billing info
				if req.HasBilling() {
					billing := req.GetBilling()
					assert.Equal(t, "monthly", billing.GetBillingFrequency())
					assert.Equal(t, "none", billing.GetPaymentType())
					assert.Equal(t, "12-3456789", billing.GetTaxId())
					assert.Equal(t, int32(30), billing.GetDefaultPaymentTerms())
					assert.Equal(t, float64(1000.00), billing.GetInvoiceAmountThreshold())

					if billing.HasPayment() {
						payment := billing.GetPayment()
						assert.Equal(t, "Example Affiliate LLC", payment.GetAccountName())
						assert.Equal(t, "123456789", payment.GetAccountNumber())
						assert.Equal(t, "987654321", payment.GetRoutingNumber())
					}
				}

				// Validate labels
				labels := req.GetLabels()
				assert.Len(t, labels, 3)
				assert.Contains(t, labels, "premium")
				assert.Contains(t, labels, "high-volume")
				assert.Contains(t, labels, "trusted")
			},
		},
		{
			name: "successful_affiliate_creation_update_failed_mapping",
			inputAffiliate: domain.Affiliate{
				AffiliateID:    1,
				OrganizationID: 1,
				Name:           "Example Affiliate",
				Status:         "active",
			},
			existingMapping: &domain.AffiliateProviderMapping{
				AffiliateID:  1,
				ProviderType: "everflow",
				SyncStatus:   stringPtr("failed"),
				SyncError:    stringPtr("Previous sync failed"),
				CreatedAt:    time.Now().Add(-1 * time.Hour),
				UpdatedAt:    time.Now().Add(-1 * time.Hour),
			},
			everflowResponse: &affiliate.Affiliate{
				NetworkAffiliateId: int32Ptr(12345),
				Name:               stringPtr("Example Affiliate"),
				AccountStatus:      stringPtr("active"),
			},
			everflowStatusCode:  200,
			expectError:         false,
			expectMappingCreate: false,
			expectMappingUpdate: true,
			validateMapping: func(t *testing.T, mapping *domain.AffiliateProviderMapping) {
				assert.Equal(t, int64(1), mapping.AffiliateID)
				assert.Equal(t, "everflow", mapping.ProviderType)
				assert.NotNil(t, mapping.SyncStatus)
				assert.Equal(t, "synced", *mapping.SyncStatus)
				assert.Nil(t, mapping.SyncError) // Should be cleared after successful sync
			},
		},
		{
			name: "error_affiliate_already_synced",
			inputAffiliate: domain.Affiliate{
				AffiliateID: 1,
				Name:        "Example Affiliate",
			},
			existingMapping: &domain.AffiliateProviderMapping{
				AffiliateID:  1,
				ProviderType: "everflow",
				SyncStatus:   stringPtr("synced"),
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			expectError:           true,
			expectedErrorContains: "already has successful Everflow provider mapping",
			expectMappingCreate:   false,
			expectMappingUpdate:   false,
		},
		{
			name: "error_everflow_api_failure",
			inputAffiliate: domain.Affiliate{
				AffiliateID: 1,
				Name:        "Example Affiliate",
			},
			existingMapping:       nil,
			everflowStatusCode:    500,
			expectError:           true,
			expectedErrorContains: "failed to create affiliate in Everflow",
			expectMappingCreate:   false,
			expectMappingUpdate:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock repository
			mockRepo := &MockAffiliateProviderMappingRepository{}

			// Setup repository expectations
			mockRepo.On("GetAffiliateProviderMapping", mock.Anything, tt.inputAffiliate.AffiliateID, "everflow").
				Return(tt.existingMapping, func() error {
					if tt.existingMapping == nil {
						return fmt.Errorf("not found")
					}
					return nil
				}())

			if tt.expectMappingCreate {
				mockRepo.On("CreateAffiliateProviderMapping", mock.Anything, mock.MatchedBy(func(mapping *domain.AffiliateProviderMapping) bool {
					if tt.validateMapping != nil {
						tt.validateMapping(t, mapping)
					}
					return true
				})).Return(nil)
			}

			if tt.expectMappingUpdate {
				mockRepo.On("UpdateAffiliateProviderMapping", mock.Anything, mock.MatchedBy(func(mapping *domain.AffiliateProviderMapping) bool {
					if tt.validateMapping != nil {
						tt.validateMapping(t, mapping)
					}
					return true
				})).Return(nil)
			}

			// Setup mock Everflow server
			var capturedRequest *affiliate.CreateAffiliateRequest
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.everflowStatusCode != 200 {
					w.WriteHeader(tt.everflowStatusCode)
					w.Write([]byte(`{"error": "Internal server error"}`))
					return
				}

				// Capture the request for validation
				var req affiliate.CreateAffiliateRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				require.NoError(t, err)
				capturedRequest = &req

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				
				responseJSON, err := json.Marshal(tt.everflowResponse)
				require.NoError(t, err)
				w.Write(responseJSON)
			}))
			defer server.Close()

			// Create Everflow client configuration
			config := affiliate.NewConfiguration()
			config.Servers = []affiliate.ServerConfiguration{
				{
					URL: server.URL,
				},
			}
			affiliateClient := affiliate.NewAPIClient(config)

			// Create integration service
			mapper := NewAffiliateProviderMapper()
			service := &IntegrationService{
				affiliateClient:                affiliateClient,
				affiliateProviderMapper:        mapper,
				affiliateProviderMappingRepo:   mockRepo,
			}

			// Execute the test
			ctx := context.Background()
			result, err := service.CreateAffiliate(ctx, tt.inputAffiliate)

			// Validate results
			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErrorContains != "" {
					assert.Contains(t, err.Error(), tt.expectedErrorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.inputAffiliate.AffiliateID, result.AffiliateID)
				assert.Equal(t, tt.inputAffiliate.Name, result.Name)

				// Validate the request sent to Everflow
				if tt.validateEverflowRequest != nil && capturedRequest != nil {
					tt.validateEverflowRequest(t, capturedRequest)
				}
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestAffiliateProviderMapper_FullRoundTrip tests the complete mapping round trip
func TestAffiliateProviderMapper_FullRoundTrip(t *testing.T) {
	mapper := NewAffiliateProviderMapper()

	// Create a comprehensive test affiliate
	originalAffiliate := &domain.Affiliate{
		AffiliateID:            1,
		OrganizationID:         1,
		Name:                   "Test Affiliate",
		ContactEmail:           stringPtr("test@example.com"),
		PaymentDetails:         stringPtr(`{"bank_account":"123456789","routing_number":"987654321","payment_method":"bank_transfer","account_name":"Test LLC"}`),
		Status:                 "active",
		InternalNotes:          stringPtr("Test notes"),
		DefaultCurrencyID:      stringPtr("USD"),
		ContactAddress:         stringPtr(`{"address1":"123 Test St","city":"Test City","region_code":"TC","country_code":"US","zip_postal_code":"12345"}`),
		BillingInfo:            stringPtr(`{"billing_frequency":"monthly","payment_type":"wire","tax_id":"12-3456789","details":{"day_of_month":15}}`),
		Labels:                 stringPtr(`["test","premium"]`),
		InvoiceAmountThreshold: float64Ptr(500.00),
		DefaultPaymentTerms:    int32Ptr(15),
	}

	// Test mapping to Everflow request
	everflowReq, err := mapper.MapAffiliateToEverflowRequest(originalAffiliate, nil)
	require.NoError(t, err)
	require.NotNil(t, everflowReq)

	// Validate key fields in the request
	assert.Equal(t, "Test Affiliate", everflowReq.GetName())
	assert.Equal(t, "active", everflowReq.GetAccountStatus())
	assert.Equal(t, "USD", everflowReq.GetDefaultCurrencyId())
	assert.Equal(t, "Test notes", everflowReq.GetInternalNotes())

	// Create a mock Everflow response
	everflowResp := &affiliate.Affiliate{
		NetworkAffiliateId:           int32Ptr(12345),
		Name:                         stringPtr("Test Affiliate"),
		AccountStatus:                stringPtr("active"),
		DefaultCurrencyId:            stringPtr("USD"),
		InternalNotes:                stringPtr("Test notes"),
		NetworkEmployeeId:            int32Ptr(1),
		ReferrerId:                   int32Ptr(0),
		EnableMediaCostTrackingLinks: boolPtr(false),
		IsContactAddressEnabled:      boolPtr(true),
	}

	// Test mapping response to provider data
	providerData, err := mapper.MapEverflowResponseToProviderData(everflowResp)
	require.NoError(t, err)
	require.NotNil(t, providerData)

	assert.NotNil(t, providerData.NetworkAffiliateID)
	assert.Equal(t, int32(12345), *providerData.NetworkAffiliateID)
	assert.NotNil(t, providerData.NetworkEmployeeID)
	assert.Equal(t, int32(1), *providerData.NetworkEmployeeID)

	// Test mapping response to affiliate (should update the original)
	updatedAffiliate := *originalAffiliate
	err = mapper.MapEverflowResponseToAffiliate(everflowResp, &updatedAffiliate)
	require.NoError(t, err)

	// The affiliate should retain its original data since MapEverflowResponseToAffiliate
	// only updates fields that come from Everflow but aren't provider-specific
	assert.Equal(t, originalAffiliate.Name, updatedAffiliate.Name)
	assert.Equal(t, originalAffiliate.Status, updatedAffiliate.Status)

	// Test mapping response to provider mapping
	mapping := &domain.AffiliateProviderMapping{
		AffiliateID:  1,
		ProviderType: "everflow",
	}

	err = mapper.MapEverflowResponseToProviderMapping(everflowResp, mapping)
	require.NoError(t, err)

	// Validate that provider data was set
	require.NotNil(t, mapping.ProviderData)
	var mappedProviderData domain.EverflowProviderData
	err = json.Unmarshal([]byte(*mapping.ProviderData), &mappedProviderData)
	require.NoError(t, err)

	assert.NotNil(t, mappedProviderData.NetworkAffiliateID)
	assert.Equal(t, int32(12345), *mappedProviderData.NetworkAffiliateID)
}

