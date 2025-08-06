package everflow

import (
	"encoding/json"
	"testing"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/everflow/affiliate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAffiliateProviderMapper(t *testing.T) {
	mapper := NewAffiliateProviderMapper()

	t.Run("MapAffiliateToEverflowRequest", func(t *testing.T) {
		t.Run("handles nil affiliate", func(t *testing.T) {
			_, err := mapper.MapAffiliateToEverflowRequest(nil, nil)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "affiliate cannot be nil")
		})

		t.Run("maps basic affiliate fields", func(t *testing.T) {
			affiliate := &domain.Affiliate{
				AffiliateID:       1,
				OrganizationID:    1,
				Name:              "Test Affiliate",
				Status:            "active",
				DefaultCurrencyID: stringPtr("USD"),
				InternalNotes:     stringPtr("Test notes"),
			}

			req, err := mapper.MapAffiliateToEverflowRequest(affiliate, nil)
			assert.NoError(t, err)
			assert.NotNil(t, req)

			assert.Equal(t, "Test Affiliate", req.GetName())
			assert.Equal(t, "active", req.GetAccountStatus())
			assert.Equal(t, "USD", req.GetDefaultCurrencyId())
			assert.Equal(t, "Test notes", req.GetInternalNotes())
			assert.Equal(t, int32(1), req.GetNetworkEmployeeId())
			assert.Equal(t, int32(0), req.GetReferrerId())
			assert.False(t, req.GetEnableMediaCostTrackingLinks())
		})

		t.Run("maps contact address", func(t *testing.T) {
			contactAddr := `{"address1":"123 Main St","city":"New York","region_code":"NY","country_code":"US","zip_postal_code":"10001"}`
			affiliate := &domain.Affiliate{
				AffiliateID:    1,
				Name:           "Test Affiliate",
				ContactAddress: &contactAddr,
			}

			req, err := mapper.MapAffiliateToEverflowRequest(affiliate, nil)
			assert.NoError(t, err)
			assert.True(t, req.HasContactAddress())

			addr := req.GetContactAddress()
			assert.Equal(t, "123 Main St", addr.GetAddress1())
			assert.Equal(t, "New York", addr.GetCity())
			assert.Equal(t, "NY", addr.GetRegionCode())
			assert.Equal(t, "US", addr.GetCountryCode())
			assert.Equal(t, "10001", addr.GetZipPostalCode())
		})

		t.Run("maps billing info", func(t *testing.T) {
			billingInfo := `{"billing_frequency":"monthly","payment_type":"wire","tax_id":"12-3456789","details":{"day_of_month":15}}`
			affiliate := &domain.Affiliate{
				AffiliateID:            1,
				Name:                   "Test Affiliate",
				BillingInfo:            &billingInfo,
				DefaultPaymentTerms:    int32Ptr(30),
				InvoiceAmountThreshold: float64Ptr(1000.00),
			}

			req, err := mapper.MapAffiliateToEverflowRequest(affiliate, nil)
			assert.NoError(t, err)
			assert.True(t, req.HasBilling())

			billing := req.GetBilling()
			assert.Equal(t, "monthly", billing.GetBillingFrequency())
			assert.Equal(t, "wire", billing.GetPaymentType())
			assert.Equal(t, "12-3456789", billing.GetTaxId())
			assert.Equal(t, int32(30), billing.GetDefaultPaymentTerms())
			assert.Equal(t, float64(1000.00), billing.GetInvoiceAmountThreshold())
		})

		t.Run("maps labels", func(t *testing.T) {
			labels := `["premium","high-volume","trusted"]`
			affiliate := &domain.Affiliate{
				AffiliateID: 1,
				Name:        "Test Affiliate",
				Labels:      &labels,
			}

			req, err := mapper.MapAffiliateToEverflowRequest(affiliate, nil)
			assert.NoError(t, err)

			reqLabels := req.GetLabels()
			assert.Len(t, reqLabels, 3)
			assert.Contains(t, reqLabels, "premium")
			assert.Contains(t, reqLabels, "high-volume")
			assert.Contains(t, reqLabels, "trusted")
		})
	})

	t.Run("MapEverflowResponseToProviderMapping", func(t *testing.T) {
		t.Run("handles nil inputs", func(t *testing.T) {
			err := mapper.MapEverflowResponseToProviderMapping(nil, nil)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid response or mapping")
		})

		t.Run("maps response to provider mapping", func(t *testing.T) {
			resp := &affiliate.Affiliate{
				NetworkAffiliateId:           int32Ptr(12345),
				Name:                         stringPtr("Test Affiliate"),
				AccountStatus:                stringPtr("active"),
				NetworkEmployeeId:            int32Ptr(1),
				EnableMediaCostTrackingLinks: boolPtr(false),
			}

			mapping := &domain.AffiliateProviderMapping{
				AffiliateID:  1,
				ProviderType: "everflow",
			}

			err := mapper.MapEverflowResponseToProviderMapping(resp, mapping)
			assert.NoError(t, err)
			assert.NotNil(t, mapping.ProviderData)

			// Validate provider data
			var providerData domain.EverflowProviderData
			err = json.Unmarshal([]byte(*mapping.ProviderData), &providerData)
			assert.NoError(t, err)
			assert.NotNil(t, providerData.NetworkAffiliateID)
			assert.Equal(t, int32(12345), *providerData.NetworkAffiliateID)
		})
	})

	t.Run("MapEverflowResponseToAffiliate", func(t *testing.T) {
		t.Run("handles nil inputs", func(t *testing.T) {
			err := mapper.MapEverflowResponseToAffiliate(nil, nil)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "response and affiliate cannot be nil")
		})

		t.Run("maps response to affiliate", func(t *testing.T) {
			resp := &affiliate.Affiliate{
				Name:              stringPtr("Updated Name"),
				AccountStatus:     stringPtr("active"),
				DefaultCurrencyId: stringPtr("EUR"),
				InternalNotes:     stringPtr("Updated notes"),
			}

			aff := &domain.Affiliate{
				AffiliateID: 1,
				Name:        "Original Name",
			}

			err := mapper.MapEverflowResponseToAffiliate(resp, aff)
			assert.NoError(t, err)

			// The affiliate should be updated with response data
			assert.Equal(t, "Updated Name", aff.Name)
			assert.Equal(t, "active", aff.Status)
			assert.NotNil(t, aff.DefaultCurrencyID)
			assert.Equal(t, "EUR", *aff.DefaultCurrencyID)
			assert.NotNil(t, aff.InternalNotes)
			assert.Equal(t, "Updated notes", *aff.InternalNotes)
		})
	})

	t.Run("MapEverflowResponseToProviderData", func(t *testing.T) {
		t.Run("handles nil response", func(t *testing.T) {
			_, err := mapper.MapEverflowResponseToProviderData(nil)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "response cannot be nil")
		})

		t.Run("maps response to provider data", func(t *testing.T) {
			resp := &affiliate.Affiliate{
				NetworkAffiliateId:           int32Ptr(12345),
				NetworkEmployeeId:            int32Ptr(1),
				EnableMediaCostTrackingLinks: boolPtr(true),
				ReferrerId:                   int32Ptr(999),
				IsContactAddressEnabled:      boolPtr(false),
				NetworkId:                    int32Ptr(100),
				AccountManagerId:             int32Ptr(200),
				AccountManagerName:           stringPtr("John Doe"),
				TimeCreated:                  int64Ptr(1704067200), // 2024-01-01T00:00:00Z
				TimeSaved:                    int64Ptr(1704067200), // 2024-01-01T00:00:00Z
			}

			providerData, err := mapper.MapEverflowResponseToProviderData(resp)
			require.NoError(t, err)
			require.NotNil(t, providerData)

			assert.NotNil(t, providerData.NetworkAffiliateID)
			assert.Equal(t, int32(12345), *providerData.NetworkAffiliateID)
			assert.NotNil(t, providerData.NetworkEmployeeID)
			assert.Equal(t, int32(1), *providerData.NetworkEmployeeID)
			assert.NotNil(t, providerData.EnableMediaCostTrackingLinks)
			assert.True(t, *providerData.EnableMediaCostTrackingLinks)
			assert.NotNil(t, providerData.ReferrerID)
			assert.Equal(t, int32(999), *providerData.ReferrerID)
			assert.NotNil(t, providerData.IsContactAddressEnabled)
			assert.False(t, *providerData.IsContactAddressEnabled)

			// Check additional fields
			assert.NotNil(t, providerData.AdditionalFields)
			assert.Equal(t, int32(100), providerData.AdditionalFields["network_id"])
			assert.Equal(t, int32(200), providerData.AdditionalFields["account_manager_id"])
			assert.Equal(t, "John Doe", providerData.AdditionalFields["account_manager_name"])
			assert.Equal(t, int64(1704067200), providerData.AdditionalFields["time_created"])
			assert.Equal(t, int64(1704067200), providerData.AdditionalFields["time_saved"])
		})
	})
}


