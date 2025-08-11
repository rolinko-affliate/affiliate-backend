package integration_tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/affiliate-backend/internal/platform/everflow/offer"
	"github.com/stretchr/testify/assert"
)

// TestDebugOfferClient tests the generated offer client directly
func TestDebugOfferClient(t *testing.T) {
	t.Log("=== Testing Generated Offer Client ===")

	// Configure the offer client exactly like our application does
	offerConfig := offer.NewConfiguration()
	offerConfig.Servers = []offer.ServerConfiguration{
		{
			URL: "https://api.eflow.team/v1",
		},
	}
	offerConfig.AddDefaultHeader("X-Eflow-API-Key", "fqGImoDQSr6zDnT758O6JA")
	offerClient := offer.NewAPIClient(offerConfig)

	// First create an advertiser to get a valid network_advertiser_id
	config := NewTestConfig()
	advertiserPayload := map[string]interface{}{
		"account_status":             "active",
		"accounting_contact_email":   "",
		"address_id":                 1,
		"affiliate_id_macro":         "",
		"attribution_method":         "last_touch",
		"attribution_priority":       "click",
		"billing": map[string]interface{}{
			"billing_frequency":       "other",
			"default_payment_terms":   0,
			"details":                 map[string]interface{}{},
			"tax_id":                  "123456789",
		},
		"contact_address": map[string]interface{}{
			"address_1":        "4110 rue St-Laurent",
			"address_2":        "202",
			"city":             "Montreal",
			"country_code":     "CA",
			"country_id":       36,
			"region_code":      "QC",
			"zip_postal_code":  "H2R 0A1",
		},
		"default_currency_id":        "USD",
		"email_attribution_method":   "last_affiliate_attribution",
		"internal_notes":             "Some notes not visible to the advertiser",
		"is_contact_address_enabled": false,
		"labels":                     []string{"DTC Brand"},
		"name":                       "Test Advertiser for Generated Client",
		"network_employee_id":        1,
		"offer_id_macro":             "",
		"platform_name":              "",
		"platform_url":               "",
		"platform_username":          "",
		"reporting_timezone_id":      80,
		"sales_manager_id":           1,
		"settings": map[string]interface{}{
			"exposed_variables": map[string]interface{}{
				"affiliate":     false,
				"affiliate_id":  true,
				"offer_url":     false,
				"source_id":     false,
				"sub1":          true,
				"sub2":          true,
				"sub3":          false,
				"sub4":          false,
				"sub5":          false,
			},
		},
		"users": []map[string]interface{}{
			{
				"account_status": "active",
				"currency_id":    "USD",
				"email":          fmt.Sprintf("test_generated_client_%d@example.com", time.Now().Unix()),
				"first_name":     "User",
				"language_id":    1,
				"last_name":      "Account",
				"timezone_id":    80,
			},
		},
	}

	t.Log("Creating advertiser first...")
	advertiserResp := config.EverflowAPIRequest(t, "POST", "/networks/advertisers", advertiserPayload)
	LogResponse(t, "Advertiser Creation", advertiserResp)

	if advertiserResp.StatusCode != 200 {
		t.Skipf("Cannot create advertiser (status %d), skipping offer test", advertiserResp.StatusCode)
		return
	}

	var advertiserResult struct {
		NetworkAdvertiserID int `json:"network_advertiser_id"`
	}
	ParseJSONResponse(t, advertiserResp, &advertiserResult)

	// Now test the generated client
	t.Log("Creating offer using generated client...")

	// Helper function to create pointers
	stringPtr := func(s string) *string { return &s }
	int32Ptr := func(i int32) *int32 { return &i }
	boolPtr := func(b bool) *bool { return &b }
	float64Ptr := func(f float64) *float64 { return &f }

	// Create the offer request using the generated types (matching working direct API call)
	offerRequest := offer.CreateOfferRequest{
		Name:                "Test Offer Generated Client",
		NetworkAdvertiserId: int32(advertiserResult.NetworkAdvertiserID),
		OfferStatus:        "active",
		DestinationUrl:     "https://example.com",
		CurrencyId:         stringPtr("USD"),
		NetworkCategoryId:  int32Ptr(1),
		// Add missing required fields from working direct API call
		NetworkTrackingDomainId:    int32Ptr(12977),
		AttributionMethod:          stringPtr("last_touch"),
		ConversionMethod:           stringPtr("server_postback"),
		RedirectMode:               stringPtr("standard"),
		SessionDefinition:          stringPtr("cookie"),
		SessionDuration:            int32Ptr(24),
		EmailAttributionMethod:     stringPtr("first_affiliate_attribution"),
		Visibility:                 stringPtr("public"),
		PayoutRevenue: []offer.PayoutRevenue{
			{
				EntryName:                  stringPtr("Base"),
				PayoutType:                 "cpa_cps",
				RevenueType:                "rpa_rps",
				PayoutAmount:               float64Ptr(2.0),
				RevenueAmount:              float64Ptr(5.0),
				PayoutPercentage:           int32Ptr(5),
				RevenuePercentage:          int32Ptr(10),
				IsDefault:                  true,
				IsPrivate:                  false,
				IsAllowDuplicateConversion: boolPtr(true),
			},
		},
	}

	ctx := context.Background()
	resp, httpResp, err := offerClient.OffersAPI.CreateOffer(ctx).CreateOfferRequest(offerRequest).Execute()

	t.Logf("Generated Client Response - Status: %d", httpResp.StatusCode)
	if err != nil {
		t.Logf("Generated Client Error: %v", err)
	}
	if resp != nil {
		t.Logf("Generated Client Success - Offer ID: %d", resp.GetNetworkOfferId())
	}

	// Check if it worked
	if httpResp.StatusCode == 200 {
		t.Log("Generated client worked!")
		assert.NotNil(t, resp)
		assert.Greater(t, resp.GetNetworkOfferId(), int32(0))
	} else {
		t.Logf("Generated client failed with status %d", httpResp.StatusCode)
		t.Fail()
	}

	// Cleanup
	t.Log("Cleaning up advertiser...")
	cleanupResp := config.EverflowAPIRequest(t, "DELETE", fmt.Sprintf("/networks/advertisers/%d", advertiserResult.NetworkAdvertiserID), nil)
	if cleanupResp.StatusCode != 200 && cleanupResp.StatusCode != 204 {
		t.Logf("Warning: Failed to delete advertiser %d (status %d)", advertiserResult.NetworkAdvertiserID, cleanupResp.StatusCode)
	}
}