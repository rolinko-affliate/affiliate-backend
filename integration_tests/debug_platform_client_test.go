package integration_tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/affiliate-backend/internal/platform/everflow/offer"
)

// TestDebugPlatformClient tests the offer client using the exact same configuration as the platform
func TestDebugPlatformClient(t *testing.T) {
	t.Log("=== Testing Platform Client Configuration ===")

	// Initialize test config
	config := NewTestConfig()

	// Create advertiser first (using exact working payload from debug_offer_client_test.go)
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
		"name":                       "Test Advertiser for Platform Client",
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
				"email":          fmt.Sprintf("test_platform_client_%d@example.com", time.Now().Unix()),
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
		t.Logf("Cannot create advertiser (status %d), skipping offer test", advertiserResp.StatusCode)
		t.Skip()
		return
	}

	var advertiserResult struct {
		NetworkAdvertiserID int    `json:"network_advertiser_id"`
		Name               string `json:"name"`
	}
	ParseJSONResponse(t, advertiserResp, &advertiserResult)

	// Now create offer client using EXACT same configuration as platform
	t.Log("Creating offer client with platform configuration...")
	
	// Use the same configuration as the platform
	offerConfig := offer.NewConfiguration()
	offerConfig.Servers = []offer.ServerConfiguration{
		{
			URL: "https://api.eflow.team/v1", // Same as platform base URL
		},
	}
	offerConfig.AddDefaultHeader("X-Eflow-API-Key", "fqGImoDQSr6zDnT758O6JA") // Same API key
	offerClient := offer.NewAPIClient(offerConfig)

	// Create offer request using the same method as platform (constructor + setters)
	t.Log("Creating offer using platform-style request construction...")
	
	// Create basic payout revenue structure (same as platform)
	payoutRevenue := []offer.PayoutRevenue{
		*offer.NewPayoutRevenue("cpa_cps", "rpa_rps", true, false),
	}
	payoutRevenue[0].SetPayoutAmount(2.0)
	payoutRevenue[0].SetPayoutPercentage(5)
	payoutRevenue[0].SetRevenueAmount(5.0)
	payoutRevenue[0].SetRevenuePercentage(10)

	// Create the offer request using constructor (same as platform)
	req := offer.NewCreateOfferRequest(
		int32(advertiserResult.NetworkAdvertiserID),
		"Test Offer Platform Client",
		"https://example.com",
		"active",
		payoutRevenue,
	)

	// Set the same optional fields as platform
	req.SetNetworkTrackingDomainId(12977)
	req.SetIsUseSecureLink(true)
	req.SetNetworkCategoryId(1)
	
	// Create empty ruleset with default timezone (same as platform)
	ruleset := offer.NewRuleset()
	ruleset.SetDayPartingTimezoneId(58)
	req.SetRuleset(*ruleset)
	
	// Set attribution methods (same as platform)
	req.SetEmailAttributionMethod("first_affiliate_attribution")
	req.SetAttributionMethod("last_touch")
	
	// Set additional required fields that platform sets
	req.SetConversionMethod("server_postback")
	req.SetRedirectMode("standard")
	req.SetSessionDefinition("cookie")
	req.SetSessionDuration(24)
	req.SetVisibility("public")
	req.SetHtmlDescription("Test campaign for platform client")

	ctx := context.Background()
	resp, httpResp, err := offerClient.OffersAPI.CreateOffer(ctx).CreateOfferRequest(*req).Execute()

	t.Logf("Platform Client Response - Status: %d", httpResp.StatusCode)
	if err != nil {
		t.Logf("Platform Client Error: %v", err)
	}

	if httpResp.StatusCode == 200 && resp != nil {
		t.Logf("Platform Client Success - Offer ID: %d", resp.GetNetworkOfferId())
		t.Log("Platform client configuration works!")
		
		// Clean up offer
		t.Log("Cleaning up offer...")
		deleteResp := config.EverflowAPIRequest(t, "DELETE", 
			fmt.Sprintf("/networks/offers/%d", resp.GetNetworkOfferId()), nil)
		if deleteResp.StatusCode != 200 {
			t.Logf("Warning: Failed to delete offer %d (status %d)", resp.GetNetworkOfferId(), deleteResp.StatusCode)
		}
	} else {
		t.Logf("Platform client failed with status %d", httpResp.StatusCode)
		if httpResp != nil && httpResp.Body != nil {
			// Try to read response body for debugging
			t.Log("Response body available for debugging")
		}
		t.Fail()
	}

	// Cleanup advertiser
	t.Log("Cleaning up advertiser...")
	cleanupResp := config.EverflowAPIRequest(t, "DELETE", fmt.Sprintf("/networks/advertisers/%d", advertiserResult.NetworkAdvertiserID), nil)
	if cleanupResp.StatusCode != 200 && cleanupResp.StatusCode != 204 {
		t.Logf("Warning: Failed to delete advertiser %d (status %d)", advertiserResult.NetworkAdvertiserID, cleanupResp.StatusCode)
	}
}