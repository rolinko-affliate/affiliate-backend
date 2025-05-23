package everflow

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_ListAdvertisers(t *testing.T) {
	// Mock response data based on the OpenAPI spec example
	mockResponse := EverflowListAdvertisersResponse{
		Advertisers: []Advertiser{
			{
				NetworkAdvertiserID:            13,
				NetworkID:                      1,
				Name:                           "Gabielle Deleon Inc.",
				AccountStatus:                  "active",
				NetworkEmployeeID:              1,
				InternalNotes:                  "",
				AddressID:                      0,
				IsContactAddressEnabled:        false,
				SalesManagerID:                 0,
				IsExposePublisherReportingData: nil,
				DefaultCurrencyID:              "USD",
				PlatformName:                   "",
				PlatformURL:                    "",
				PlatformUsername:               "",
				ReportingTimezoneID:            67,
				AccountingContactEmail:         "",
				VerificationToken:              "",
				OfferIDMacro:                   "oid",
				AffiliateIDMacro:               "affid",
				TimeCreated:                    1582295424,
				TimeSaved:                      1582295424,
				AttributionMethod:              "last_touch",
				EmailAttributionMethod:         "last_affiliate_attribution",
				AttributionPriority:            "click",
				Relationship: &AdvertiserRelationship{
					Labels: &AdvertiserLabels{
						Total:   0,
						Entries: []string{},
					},
					AccountManager: &AccountManager{
						FirstName:                  "Bob",
						LastName:                   "Smith",
						Email:                      "my.everflow@gmail.com",
						WorkPhone:                  "8734215936",
						CellPhone:                  "",
						InstantMessagingID:         0,
						InstantMessagingIdentifier: "",
					},
					Integrations: &AdvertiserIntegrations{
						AdvertiserDemandPartner: nil,
					},
				},
			},
		},
		Paging: Paging{
			Page:       1,
			PageSize:   50,
			TotalCount: 1,
		},
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/networks/advertisers", r.URL.Path)
		assert.Equal(t, "test-api-key", r.Header.Get("X-Eflow-API-Key"))

		// Check query parameters
		if r.URL.Query().Get("page") != "" {
			assert.Equal(t, "1", r.URL.Query().Get("page"))
		}
		if r.URL.Query().Get("page_size") != "" {
			assert.Equal(t, "50", r.URL.Query().Get("page_size"))
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client with test server URL
	client := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		apiKey:     "test-api-key",
	}

	// Override the base URL for testing
	originalURL := everflowAPIBaseURL
	everflowAPIBaseURL = server.URL
	defer func() { everflowAPIBaseURL = originalURL }()

	// Test without options
	t.Run("without options", func(t *testing.T) {
		ctx := context.Background()
		resp, err := client.ListAdvertisers(ctx, nil)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Advertisers, 1)
		assert.Equal(t, "Gabielle Deleon Inc.", resp.Advertisers[0].Name)
		assert.Equal(t, int64(13), resp.Advertisers[0].NetworkAdvertiserID)
		assert.Equal(t, 1, resp.Paging.Page)
		assert.Equal(t, 50, resp.Paging.PageSize)
		assert.Equal(t, 1, resp.Paging.TotalCount)
	})

	// Test with options
	t.Run("with options", func(t *testing.T) {
		ctx := context.Background()
		page := 1
		pageSize := 50
		opts := &ListAdvertisersOptions{
			Page:     &page,
			PageSize: &pageSize,
		}

		resp, err := client.ListAdvertisers(ctx, opts)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Advertisers, 1)
	})
}

func TestClient_CreateAdvertiser(t *testing.T) {
	// Mock response data based on the OpenAPI spec example
	mockResponse := Advertiser{
		NetworkAdvertiserID:            289,
		NetworkID:                      5,
		Name:                           "Some Brand Inc.",
		AccountStatus:                  "active",
		NetworkEmployeeID:              264,
		InternalNotes:                  "Some notes not visible to the advertiser",
		SalesManagerID:                 227,
		IsExposePublisherReportingData: nil,
		DefaultCurrencyID:              "USD",
		PlatformName:                   "",
		PlatformURL:                    "",
		PlatformUsername:               "",
		ReportingTimezoneID:            80,
		AccountingContactEmail:         "",
		VerificationToken:              "c7HIWpFUGnyQfN5wwBollBBGtUkeOm",
		OfferIDMacro:                   "oid",
		AffiliateIDMacro:               "affid",
		TimeCreated:                    1727214292,
		TimeSaved:                      1727214610,
		AttributionMethod:              "last_touch",
		EmailAttributionMethod:         "last_affiliate_attribution",
		AttributionPriority:            "coupon_code",
		IsContactAddressEnabled:        true,
		AddressID:                      84636,
		ContactAddress: &AdvertiserAddress{
			Address1:      "4110 rue St-Laurent",
			Address2:      stringPtr("202"),
			City:          "Montreal",
			ZipPostalCode: "H2R 0A1",
			CountryCode:   "CA",
			RegionCode:    "QC",
		},
		Labels: []string{"DTC Brand"},
		Users:  []AdvertiserUser{},
		Billing: &AdvertiserBillingResponse{
			NetworkID:                  5,
			NetworkAdvertiserID:        289,
			BillingFrequency:           "other",
			InvoiceAmountThreshold:     0,
			TaxID:                      "123456789",
			IsInvoiceCreationAuto:      false,
			AutoInvoiceStartDate:       "2019-06-01 00:00:00",
			DefaultInvoiceIsHidden:     false,
			InvoiceGenerationDaysDelay: 0,
			DefaultPaymentTerms:        0,
		},
		Settings: &AdvertiserSettings{
			ExposedVariables: map[string]bool{
				"affiliate_id": true,
				"affiliate":    false,
				"sub1":         true,
				"sub2":         true,
				"sub3":         false,
				"sub4":         false,
				"sub5":         false,
				"source_id":    false,
				"offer_url":    false,
			},
		},
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/v1/networks/advertisers", r.URL.Path)
		assert.Equal(t, "test-api-key", r.Header.Get("X-Eflow-API-Key"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Verify request body
		var req EverflowCreateAdvertiserRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "Some Brand Inc.", req.Name)
		assert.Equal(t, "active", req.AccountStatus)
		assert.Equal(t, "USD", req.DefaultCurrencyID)

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client with test server URL
	client := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		apiKey:     "test-api-key",
	}

	// Override the base URL for testing
	originalURL := everflowAPIBaseURL
	everflowAPIBaseURL = server.URL
	defer func() { everflowAPIBaseURL = originalURL }()

	// Test create advertiser
	ctx := context.Background()
	networkEmployeeID := 264
	salesManagerID := 227
	reportingTimezoneID := 80
	attributionMethod := "last_touch"
	emailAttributionMethod := "last_affiliate_attribution"
	attributionPriority := "coupon_code"
	verificationToken := "c7HIWpFUGnyQfN5wwBollBBGtUkeOm"
	internalNotes := "Some notes not visible to the advertiser"
	isContactAddressEnabled := true

	req := EverflowCreateAdvertiserRequest{
		Name:                    "Some Brand Inc.",
		AccountStatus:           "active",
		DefaultCurrencyID:       "USD",
		NetworkEmployeeID:       &networkEmployeeID,
		SalesManagerID:          &salesManagerID,
		ReportingTimezoneID:     &reportingTimezoneID,
		AttributionMethod:       &attributionMethod,
		EmailAttributionMethod:  &emailAttributionMethod,
		AttributionPriority:     &attributionPriority,
		VerificationToken:       &verificationToken,
		InternalNotes:           &internalNotes,
		IsContactAddressEnabled: &isContactAddressEnabled,
		ContactAddress: &AdvertiserAddress{
			Address1:      "4110 rue St-Laurent",
			Address2:      stringPtr("202"),
			City:          "Montreal",
			ZipPostalCode: "H2R 0A1",
			CountryCode:   "CA",
			RegionCode:    "QC",
		},
		Labels: []string{"DTC Brand"},
		Users: []AdvertiserUser{
			{
				AccountStatus:   "active",
				LanguageID:      intPtr(1),
				TimezoneID:      &reportingTimezoneID,
				CurrencyID:      stringPtr("USD"),
				FirstName:       "John",
				LastName:        "Doe",
				Email:           "john.doe@example.com",
				InitialPassword: stringPtr(""),
			},
		},
		Billing: &AdvertiserBilling{
			BillingFrequency:    "other",
			TaxID:               stringPtr("123456789"),
			DefaultPaymentTerms: intPtr(0),
			Details:             map[string]interface{}{},
		},
		Settings: &AdvertiserSettings{
			ExposedVariables: map[string]bool{
				"affiliate_id": true,
				"affiliate":    false,
				"sub1":         true,
				"sub2":         true,
				"sub3":         false,
				"sub4":         false,
				"sub5":         false,
				"source_id":    false,
				"offer_url":    false,
			},
		},
	}

	resp, err := client.CreateAdvertiser(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Some Brand Inc.", resp.Name)
	assert.Equal(t, int64(289), resp.NetworkAdvertiserID)
	assert.Equal(t, "active", resp.AccountStatus)
	assert.Equal(t, "USD", resp.DefaultCurrencyID)
}

func TestClient_GetAdvertiser(t *testing.T) {
	// Mock response data based on the OpenAPI spec example
	mockResponse := Advertiser{
		NetworkAdvertiserID:            1,
		NetworkID:                      1,
		Name:                           "Google",
		AccountStatus:                  "active",
		NetworkEmployeeID:              11,
		InternalNotes:                  "",
		AddressID:                      0,
		IsContactAddressEnabled:        false,
		SalesManagerID:                 17,
		IsExposePublisherReportingData: nil,
		DefaultCurrencyID:              "USD",
		PlatformName:                   "",
		PlatformURL:                    "",
		PlatformUsername:               "",
		ReportingTimezoneID:            67,
		AccountingContactEmail:         "",
		VerificationToken:              "",
		OfferIDMacro:                   "oid",
		AffiliateIDMacro:               "affid",
		AttributionMethod:              "last_touch",
		EmailAttributionMethod:         "last_affiliate_attribution",
		AttributionPriority:            "click",
		TimeCreated:                    1559919745,
		TimeSaved:                      1559919745,
		Relationship: &AdvertiserRelationship{
			Labels: &AdvertiserLabels{
				Total:   0,
				Entries: []string{},
			},
			AccountManager: &AccountManager{
				FirstName:                  "Bob",
				LastName:                   "Smith",
				Email:                      "my.everflow@gmail.com",
				WorkPhone:                  "",
				CellPhone:                  "",
				InstantMessagingID:         0,
				InstantMessagingIdentifier: "",
			},
			Reporting: &AdvertiserReporting{
				Imp:            0,
				TotalClick:     0,
				UniqueClick:    0,
				InvalidClick:   0,
				DuplicateClick: 0,
				GrossClick:     0,
				CTR:            0,
				CV:             0,
				InvalidCVScrub: 0,
				ViewThroughCV:  0,
				TotalCV:        0,
				Event:          0,
				CVR:            0,
				EVR:            0,
				CPC:            0,
				CPM:            0,
				CPA:            0,
				EPC:            0,
				RPC:            0,
				RPA:            0,
				RPM:            0,
				Payout:         0,
				Revenue:        0,
			},
			APIKeys: &AdvertiserAPIKeys{
				Total:   0,
				Entries: []interface{}{},
			},
			APIWhitelistIPs: &AdvertiserAPIWhitelistIPs{
				Total:   0,
				Entries: []interface{}{},
			},
			Billing: &AdvertiserBillingResponse{
				NetworkID:                  63,
				NetworkAdvertiserID:        13,
				BillingFrequency:           "other",
				InvoiceAmountThreshold:     0,
				TaxID:                      "",
				IsInvoiceCreationAuto:      false,
				AutoInvoiceStartDate:       "2019-06-01 00:00:00",
				DefaultInvoiceIsHidden:     false,
				InvoiceGenerationDaysDelay: 0,
				DefaultPaymentTerms:        0,
			},
			Settings: &AdvertiserSettings{
				ExposedVariables: map[string]bool{
					"affiliate_id": false,
					"affiliate":    false,
					"sub1":         false,
					"sub2":         false,
					"sub3":         false,
					"sub4":         false,
					"sub5":         false,
					"source_id":    false,
				},
			},
			SalesManager: &SalesManager{
				FirstName:                  "Bob",
				LastName:                   "Smith",
				Email:                      "my.everflow@gmail.com",
				WorkPhone:                  "4878866676",
				CellPhone:                  "",
				InstantMessagingID:         0,
				InstantMessagingIdentifier: "",
			},
		},
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/networks/advertisers/1", r.URL.Path)
		assert.Equal(t, "test-api-key", r.Header.Get("X-Eflow-API-Key"))

		// Check query parameters for relationships
		relationships := r.URL.Query()["relationship"]
		if len(relationships) > 0 {
			expectedRelationships := []string{"reporting", "labels", "billing"}
			assert.ElementsMatch(t, expectedRelationships, relationships)
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client with test server URL
	client := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		apiKey:     "test-api-key",
	}

	// Override the base URL for testing
	originalURL := everflowAPIBaseURL
	everflowAPIBaseURL = server.URL
	defer func() { everflowAPIBaseURL = originalURL }()

	// Test without relationships
	t.Run("without relationships", func(t *testing.T) {
		ctx := context.Background()
		resp, err := client.GetAdvertiser(ctx, 1, nil)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "Google", resp.Name)
		assert.Equal(t, int64(1), resp.NetworkAdvertiserID)
		assert.Equal(t, "active", resp.AccountStatus)
	})

	// Test with relationships
	t.Run("with relationships", func(t *testing.T) {
		ctx := context.Background()
		opts := &GetAdvertiserOptions{
			Relationships: []string{RelationshipReporting, RelationshipLabels, RelationshipBilling},
		}

		resp, err := client.GetAdvertiser(ctx, 1, opts)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "Google", resp.Name)
		assert.NotNil(t, resp.Relationship)
		assert.NotNil(t, resp.Relationship.Reporting)
		assert.NotNil(t, resp.Relationship.Labels)
		assert.NotNil(t, resp.Relationship.Billing)
	})

	// Test not found
	t.Run("not found", func(t *testing.T) {
		// Create a server that returns 404
		notFoundServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer notFoundServer.Close()

		// Override the base URL for testing
		everflowAPIBaseURL = notFoundServer.URL
		defer func() { everflowAPIBaseURL = originalURL }()

		ctx := context.Background()
		resp, err := client.GetAdvertiser(ctx, 999, nil)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "advertiser with ID 999 not found")
	})
}

func TestClient_UpdateAdvertiser(t *testing.T) {
	// Mock response data
	mockResponse := EverflowUpdateAdvertiserResponse{
		Result: true,
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/v1/networks/advertisers/123", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "test-api-key", r.Header.Get("X-Eflow-API-Key"))

		// Verify request body
		var req EverflowUpdateAdvertiserRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Brand Name", req.Name)
		assert.Equal(t, "active", req.AccountStatus)
		assert.Equal(t, 123, req.NetworkEmployeeID)
		assert.Equal(t, "USD", req.DefaultCurrencyID)
		assert.Equal(t, 80, req.ReportingTimezoneID)
		assert.NotNil(t, req.InternalNotes)
		assert.Equal(t, "Updated notes", *req.InternalNotes)
		assert.NotNil(t, req.AttributionMethod)
		assert.Equal(t, "last_touch", *req.AttributionMethod)
		assert.NotNil(t, req.EmailAttributionMethod)
		assert.Equal(t, "last_affiliate_attribution", *req.EmailAttributionMethod)
		assert.NotNil(t, req.AttributionPriority)
		assert.Equal(t, "click", *req.AttributionPriority)
		assert.Equal(t, []string{"Enterprise", "Updated"}, req.Labels)

		// Verify nested objects
		assert.NotNil(t, req.ContactAddress)
		assert.Equal(t, "123 New Ave", req.ContactAddress.Address1)
		assert.Equal(t, "Suite 100", *req.ContactAddress.Address2)
		assert.Equal(t, "New York", req.ContactAddress.City)
		assert.Equal(t, "10001", req.ContactAddress.ZipPostalCode)
		assert.Equal(t, "US", req.ContactAddress.CountryCode)
		assert.Equal(t, "NY", req.ContactAddress.RegionCode)

		assert.NotNil(t, req.Billing)
		assert.Equal(t, "monthly", req.Billing.BillingFrequency)
		assert.NotNil(t, req.Billing.TaxID)
		assert.Equal(t, "ABCD1234", *req.Billing.TaxID)
		assert.NotNil(t, req.Billing.IsInvoiceCreationAuto)
		assert.True(t, *req.Billing.IsInvoiceCreationAuto)
		assert.NotNil(t, req.Billing.DefaultPaymentTerms)
		assert.Equal(t, 30, *req.Billing.DefaultPaymentTerms)

		assert.NotNil(t, req.Settings)
		assert.NotNil(t, req.Settings.ExposedVariables)
		assert.True(t, req.Settings.ExposedVariables["affiliate_id"])
		assert.False(t, req.Settings.ExposedVariables["affiliate"])
		assert.True(t, req.Settings.ExposedVariables["sub1"])
		assert.True(t, req.Settings.ExposedVariables["sub2"])
		assert.False(t, req.Settings.ExposedVariables["sub3"])
		assert.False(t, req.Settings.ExposedVariables["sub4"])
		assert.False(t, req.Settings.ExposedVariables["sub5"])
		assert.False(t, req.Settings.ExposedVariables["source_id"])
		assert.True(t, req.Settings.ExposedVariables["offer_url"])

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client with test server URL
	client := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		apiKey:     "test-api-key",
	}

	// Override the base URL for testing
	originalURL := everflowAPIBaseURL
	everflowAPIBaseURL = server.URL
	defer func() { everflowAPIBaseURL = originalURL }()

	// Test successful update
	t.Run("successful update", func(t *testing.T) {
		ctx := context.Background()
		req := EverflowUpdateAdvertiserRequest{
			Name:                    "Updated Brand Name",
			AccountStatus:           "active",
			NetworkEmployeeID:       123,
			InternalNotes:           stringPtr("Updated notes"),
			AddressID:               intPtr(1234),
			IsContactAddressEnabled: boolPtr(true),
			SalesManagerID:          intPtr(321),
			DefaultCurrencyID:       "USD",
			PlatformName:            stringPtr("MyPlatform"),
			PlatformURL:             stringPtr("https://myplatform.example.com"),
			PlatformUsername:        stringPtr("brand_user"),
			ReportingTimezoneID:     80,
			AttributionMethod:       stringPtr("last_touch"),
			EmailAttributionMethod:  stringPtr("last_affiliate_attribution"),
			AttributionPriority:     stringPtr("click"),
			AccountingContactEmail:  stringPtr("accounting@brand.com"),
			VerificationToken:       stringPtr("NewToken123"),
			OfferIDMacro:            stringPtr("oid"),
			AffiliateIDMacro:        stringPtr("affid"),
			Labels:                  []string{"Enterprise", "Updated"},
			ContactAddress: &AdvertiserAddress{
				Address1:      "123 New Ave",
				Address2:      stringPtr("Suite 100"),
				City:          "New York",
				ZipPostalCode: "10001",
				CountryCode:   "US",
				RegionCode:    "NY",
			},
			Billing: &AdvertiserBilling{
				BillingFrequency:      "monthly",
				TaxID:                 stringPtr("ABCD1234"),
				IsInvoiceCreationAuto: boolPtr(true),
				DefaultPaymentTerms:   intPtr(30),
				Details:               map[string]interface{}{},
			},
			Settings: &AdvertiserSettings{
				ExposedVariables: map[string]bool{
					"affiliate_id": true,
					"affiliate":    false,
					"sub1":         true,
					"sub2":         true,
					"sub3":         false,
					"sub4":         false,
					"sub5":         false,
					"source_id":    false,
					"offer_url":    true,
				},
			},
		}

		resp, err := client.UpdateAdvertiser(ctx, 123, req)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.Result)
	})

	// Test not found
	t.Run("not found", func(t *testing.T) {
		// Create a server that returns 404
		notFoundServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer notFoundServer.Close()

		// Override the base URL for testing
		everflowAPIBaseURL = notFoundServer.URL
		defer func() { everflowAPIBaseURL = originalURL }()

		ctx := context.Background()
		req := EverflowUpdateAdvertiserRequest{
			Name:                "Test",
			AccountStatus:       "active",
			NetworkEmployeeID:   123,
			DefaultCurrencyID:   "USD",
			ReportingTimezoneID: 80,
		}

		resp, err := client.UpdateAdvertiser(ctx, 999, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "advertiser with ID 999 not found")
	})

	// Test minimal required fields
	t.Run("minimal required fields", func(t *testing.T) {
		// Create a simple server for this test
		minimalServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request body has required fields
			var req EverflowUpdateAdvertiserRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, "Minimal Brand", req.Name)
			assert.Equal(t, "active", req.AccountStatus)
			assert.Equal(t, 456, req.NetworkEmployeeID)
			assert.Equal(t, "EUR", req.DefaultCurrencyID)
			assert.Equal(t, 67, req.ReportingTimezoneID)

			// Return success response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(EverflowUpdateAdvertiserResponse{Result: true})
		}))
		defer minimalServer.Close()

		// Override the base URL for testing
		everflowAPIBaseURL = minimalServer.URL
		defer func() { everflowAPIBaseURL = originalURL }()

		ctx := context.Background()
		req := EverflowUpdateAdvertiserRequest{
			Name:                "Minimal Brand",
			AccountStatus:       "active",
			NetworkEmployeeID:   456,
			DefaultCurrencyID:   "EUR",
			ReportingTimezoneID: 67,
		}

		resp, err := client.UpdateAdvertiser(ctx, 123, req)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.Result)
	})
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
