package everflow

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/everflow/affiliate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAffiliateMapping(t *testing.T) {
	service := &IntegrationService{}

	t.Run("mapAffiliateToEverflowRequest", func(t *testing.T) {
		t.Run("with all fields", func(t *testing.T) {
			// Create domain affiliate with all fields
			aff := &domain.Affiliate{
				AffiliateID:                     1,
				OrganizationID:                  100,
				Name:                           "Test Affiliate",
				ContactEmail:                   stringPtr("test@example.com"),
				Status:                         "active",
				NetworkAffiliateID:             int32Ptr(1234),
				InternalNotes:                  stringPtr("This is a test affiliate created using the API"),
				DefaultCurrencyID:              stringPtr("USD"),
				EnableMediaCostTrackingLinks:   boolPtr(false),
				ReferrerID:                     int32Ptr(0),
				IsContactAddressEnabled:        boolPtr(true),
				NetworkAffiliateTierID:         int32Ptr(1),
				NetworkEmployeeID:              int32Ptr(1),
				CreatedAt:                      time.Now(),
				UpdatedAt:                      time.Now(),
			}

			req, err := service.mapAffiliateToEverflowRequest(aff)
			require.NoError(t, err)
			require.NotNil(t, req)

			// Verify required fields
			assert.Equal(t, "Test Affiliate", req.GetName())
			assert.Equal(t, "active", req.GetAccountStatus())
			assert.Equal(t, int32(1), req.GetNetworkEmployeeId())

			// Verify optional fields
			assert.Equal(t, "This is a test affiliate created using the API", req.GetInternalNotes())
			assert.Equal(t, "USD", req.GetDefaultCurrencyId())
			assert.Equal(t, false, req.GetEnableMediaCostTrackingLinks())
			assert.Equal(t, int32(0), req.GetReferrerId())
			assert.Equal(t, true, req.GetIsContactAddressEnabled())
			assert.Equal(t, int32(1), req.GetNetworkAffiliateTierId())
		})

		t.Run("with minimal fields", func(t *testing.T) {
			// Create domain affiliate with only required fields
			aff := &domain.Affiliate{
				AffiliateID:    1,
				OrganizationID: 100,
				Name:           "Minimal Affiliate",
				Status:         "pending",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			req, err := service.mapAffiliateToEverflowRequest(aff)
			require.NoError(t, err)
			require.NotNil(t, req)

			// Verify required fields
			assert.Equal(t, "Minimal Affiliate", req.GetName())
			assert.Equal(t, "pending", req.GetAccountStatus())
			assert.Equal(t, int32(1), req.GetNetworkEmployeeId()) // Default value

			// Verify optional fields are not set (except those with defaults)
			assert.False(t, req.HasInternalNotes())
			assert.False(t, req.HasDefaultCurrencyId())
			assert.False(t, req.HasNetworkAffiliateTierId())
			
			// These fields have default values set by the constructor
			assert.True(t, req.HasEnableMediaCostTrackingLinks())
			assert.Equal(t, false, req.GetEnableMediaCostTrackingLinks()) // Default value
			assert.True(t, req.HasReferrerId())
			assert.Equal(t, int32(0), req.GetReferrerId()) // Default value
			assert.True(t, req.HasIsContactAddressEnabled())
			assert.Equal(t, false, req.GetIsContactAddressEnabled()) // Default value
		})

		t.Run("with contact address and billing info", func(t *testing.T) {
			// Create domain affiliate with new structured data
			labelsJSON := `["premium", "trusted"]`
			aff := &domain.Affiliate{
				AffiliateID:    1,
				OrganizationID: 100,
				Name:           "Full Featured Affiliate",
				Status:         "active",
				NetworkEmployeeID: int32Ptr(1),
				
				// Contact Address
				ContactAddress: &domain.ContactAddress{
					Address1:       stringPtr("123 Main Street"),
					City:           stringPtr("New York"),
					RegionCode:     stringPtr("NY"),
					CountryCode:    stringPtr("US"),
					ZipPostalCode:  stringPtr("10001"),
				},
				
				// Billing Information
				BillingInfo: &domain.BillingDetails{
					Frequency:             (*domain.BillingFrequency)(stringPtr("monthly")),
					PaymentType:           (*domain.PaymentType)(stringPtr("wire")),
					TaxID:                 stringPtr("12-3456789"),
					IsInvoiceCreationAuto: boolPtr(true),
					Schedule: &domain.BillingSchedule{
						DayOfMonth:    int32Ptr(15),
						StartingMonth: int32Ptr(1),
					},
					PaymentDetails: &domain.PaymentDetails{
						Type:          (*domain.PaymentDetailsType)(stringPtr("wire")),
						BankName:      stringPtr("Test Bank"),
						AccountNumber: stringPtr("123456789"),
						RoutingNumber: stringPtr("987654321"),
						SwiftCode:     stringPtr("TESTSWIFT"),
					},
				},
				
				// Labels
				Labels: &labelsJSON,
				
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			req, err := service.mapAffiliateToEverflowRequest(aff)
			require.NoError(t, err)
			require.NotNil(t, req)

			// Verify basic fields
			assert.Equal(t, "Full Featured Affiliate", req.GetName())
			assert.Equal(t, "active", req.GetAccountStatus())

			// Verify contact address is mapped
			assert.True(t, req.HasContactAddress())
			contactAddr := req.GetContactAddress()
			assert.Equal(t, "123 Main Street", contactAddr.GetAddress1())
			assert.Equal(t, "New York", contactAddr.GetCity())
			assert.Equal(t, "NY", contactAddr.GetRegionCode())
			assert.Equal(t, "US", contactAddr.GetCountryCode())
			assert.Equal(t, "10001", contactAddr.GetZipPostalCode())

			// Verify billing information is mapped
			assert.True(t, req.HasBilling())
			billing := req.GetBilling()
			assert.Equal(t, "monthly", billing.GetBillingFrequency())
			assert.Equal(t, "wire", billing.GetPaymentType())
			assert.Equal(t, "12-3456789", billing.GetTaxId())
			assert.Equal(t, true, billing.GetIsInvoiceCreationAuto())

			// Verify billing details (schedule fields are in the main billing object)
			assert.True(t, billing.HasDetails())
			details := billing.GetDetails()
			assert.Equal(t, int32(15), details.GetDayOfMonth())
			assert.Equal(t, int32(1), details.GetStartingMonth())

			// Verify labels
			assert.True(t, req.HasLabels())
			labels := req.GetLabels()
			assert.Len(t, labels, 2)
			assert.Contains(t, labels, "premium")
			assert.Contains(t, labels, "trusted")
		})
	})

	t.Run("mapAffiliateToEverflowUpdateRequest", func(t *testing.T) {
		t.Run("with all fields", func(t *testing.T) {
			aff := &domain.Affiliate{
				AffiliateID:                     1,
				OrganizationID:                  100,
				Name:                           "Updated Affiliate",
				Status:                         "inactive",
				NetworkEmployeeID:              int32Ptr(2),
				InternalNotes:                  stringPtr("Updated notes"),
				DefaultCurrencyID:              stringPtr("EUR"),
				EnableMediaCostTrackingLinks:   boolPtr(true),
				ReferrerID:                     int32Ptr(5),
				IsContactAddressEnabled:        boolPtr(false),
				NetworkAffiliateTierID:         int32Ptr(2),
			}

			req, err := service.mapAffiliateToEverflowUpdateRequest(aff)
			require.NoError(t, err)
			require.NotNil(t, req)

			// Verify required fields
			assert.Equal(t, "Updated Affiliate", req.GetName())
			assert.Equal(t, "inactive", req.GetAccountStatus())
			assert.Equal(t, int32(2), req.GetNetworkEmployeeId())

			// Verify optional fields
			assert.Equal(t, "Updated notes", req.GetInternalNotes())
			assert.Equal(t, "EUR", req.GetDefaultCurrencyId())
			assert.Equal(t, true, req.GetEnableMediaCostTrackingLinks())
			assert.Equal(t, int32(5), req.GetReferrerId())
			assert.Equal(t, false, req.GetIsContactAddressEnabled())
			assert.Equal(t, int32(2), req.GetNetworkAffiliateTierId())
		})
	})

	t.Run("mapEverflowCreateResponseToAffiliate", func(t *testing.T) {
		// Load test data from JSON file
		var responseData map[string]interface{}
		err := json.Unmarshal([]byte(createAffiliateResponseJSON), &responseData)
		require.NoError(t, err)

		// Create Everflow response object
		resp := &affiliate.Affiliate{}
		resp.SetNetworkAffiliateId(int32(responseData["network_affiliate_id"].(float64)))
		resp.SetName(responseData["name"].(string))
		resp.SetAccountStatus(responseData["account_status"].(string))

		// Original domain affiliate
		originalAff := &domain.Affiliate{
			AffiliateID:    1,
			OrganizationID: 100,
			Name:           "Original Name",
			Status:         "pending",
		}

		result := service.mapEverflowCreateResponseToAffiliate(resp, originalAff)

		// Verify the mapping
		assert.Equal(t, int64(1), result.AffiliateID) // Original ID preserved
		assert.Equal(t, int64(100), result.OrganizationID) // Original org ID preserved
		assert.Equal(t, "Test Affiliate", result.Name) // Updated from response
		assert.Equal(t, int32(1234), *result.NetworkAffiliateID) // Set from response
	})

	t.Run("mapEverflowResponseToAffiliate", func(t *testing.T) {
		// Create a comprehensive Everflow response
		resp := &affiliate.AffiliateWithRelationships{}
		resp.SetName("Test Affiliate")
		resp.SetAccountStatus("active")
		resp.SetInternalNotes("This is a test affiliate created using the API")
		resp.SetDefaultCurrencyId("USD")
		resp.SetEnableMediaCostTrackingLinks(false)
		resp.SetReferrerId(0)
		resp.SetIsContactAddressEnabled(true)
		resp.SetNetworkEmployeeId(1)

		// Original domain affiliate
		originalAff := &domain.Affiliate{
			AffiliateID:    1,
			OrganizationID: 100,
			Name:           "Original Name",
			Status:         "pending",
		}

		result := service.mapEverflowResponseToAffiliate(resp, originalAff)

		// Verify the mapping
		assert.Equal(t, int64(1), result.AffiliateID) // Original ID preserved
		assert.Equal(t, int64(100), result.OrganizationID) // Original org ID preserved
		assert.Equal(t, "Test Affiliate", result.Name) // Updated from response
		assert.Equal(t, "active", result.Status) // Updated from response
		assert.Equal(t, "This is a test affiliate created using the API", *result.InternalNotes)
		assert.Equal(t, "USD", *result.DefaultCurrencyID)
		assert.Equal(t, false, *result.EnableMediaCostTrackingLinks)
		assert.Equal(t, int32(0), *result.ReferrerID)
		assert.Equal(t, true, *result.IsContactAddressEnabled)
		assert.Equal(t, int32(1), *result.NetworkEmployeeID)
	})
}

func TestStatusMapping(t *testing.T) {
	service := &IntegrationService{}

	t.Run("mapDomainStatusToEverflowStatus", func(t *testing.T) {
		testCases := []struct {
			domainStatus   string
			expectedStatus string
		}{
			{"active", "active"},
			{"pending", "pending"},
			{"rejected", "rejected"},
			{"inactive", "inactive"},
			{"unknown", "pending"}, // Default case
		}

		for _, tc := range testCases {
			t.Run(tc.domainStatus, func(t *testing.T) {
				result := service.mapDomainStatusToEverflowStatus(tc.domainStatus)
				assert.Equal(t, tc.expectedStatus, result)
			})
		}
	})

	t.Run("mapEverflowStatusToDomainStatus", func(t *testing.T) {
		testCases := []struct {
			everflowStatus string
			expectedStatus string
		}{
			{"active", "active"},
			{"pending", "pending"},
			{"rejected", "rejected"},
			{"inactive", "inactive"},
			{"unknown", "pending"}, // Default case
		}

		for _, tc := range testCases {
			t.Run(tc.everflowStatus, func(t *testing.T) {
				result := service.mapEverflowStatusToDomainStatus(tc.everflowStatus)
				assert.Equal(t, tc.expectedStatus, result)
			})
		}
	})
}

func TestHelperFunctions(t *testing.T) {
	service := &IntegrationService{}

	t.Run("getDefaultNetworkEmployeeID", func(t *testing.T) {
		t.Run("with provided ID", func(t *testing.T) {
			result := service.getDefaultNetworkEmployeeID(int32Ptr(5))
			assert.Equal(t, int32(5), result)
		})

		t.Run("with nil ID", func(t *testing.T) {
			result := service.getDefaultNetworkEmployeeID(nil)
			assert.Equal(t, int32(1), result) // Default value
		})
	})
}

// Test data based on the example files
const createAffiliateResponseJSON = `{
  "network_affiliate_id": 1234,
  "network_id": 346,
  "name": "Test Affiliate",
  "account_status": "active",
  "network_employee_id": 1,
  "account_manager_id": 1,
  "account_manager_name": "John Doe",
  "account_executive_id": 0,
  "account_executive_name": "",
  "internal_notes": "This is a test affiliate created using the API",
  "has_notifications": true,
  "network_traffic_source_id": 0,
  "adress_id": 5678,
  "default_currency_id": "USD",
  "is_contact_address_enabled": true,
  "enable_media_cost_tracking_links": false,
  "network_affiliate_tier_id": 1,
  "today_revenue": "$0.00",
  "time_created": 1735650123,
  "time_saved": 1735650123,
  "labels": [
    "test",
    "type 1"
  ],
  "balance": 0,
  "last_login": 0,
  "global_tracking_domain_url": "",
  "network_country_code": "US",
  "is_payable": true,
  "payment_type": "wire",
  "referrer_id": 0
}`

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}