package integration_tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestEverflowDirectOfferCreation tests creating an offer directly via Everflow API
func TestEverflowDirectOfferCreation(t *testing.T) {
	config := NewTestConfig()

	t.Log("=== Testing Direct Everflow Offer Creation ===")

	// Use the exact working advertiser payload from our previous tests
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
		"name":                       "Test Advertiser for Offer",
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
				"email":          fmt.Sprintf("test_offer_%d@example.com", time.Now().Unix()),
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
	}

	var advertiserResult struct {
		NetworkAdvertiserID int `json:"network_advertiser_id"`
	}
	ParseJSONResponse(t, advertiserResp, &advertiserResult)

	// Now create an offer using the advertiser ID - use the same payload structure as our platform
	offerPayload := map[string]interface{}{
		"name":                    "Test Offer Direct",
		"network_advertiser_id":   advertiserResult.NetworkAdvertiserID,
		"offer_status":           "active",
		"destination_url":        "https://example.com",
		"currency_id":            "USD",
		"network_category_id":    1,
		"network_tracking_domain_id": 12977,
		"attribution_method":     "last_touch",
		"conversion_method":      "server_postback",
		"redirect_mode":          "standard",
		"session_definition":     "cookie",
		"session_duration":       24,
		"email_attribution_method": "first_affiliate_attribution",
		"visibility":             "public",
		"payout_revenue": []map[string]interface{}{
			{
				"entry_name":                    "Base",
				"payout_type":                  "cpa_cps",
				"revenue_type":                 "rpa_rps",
				"payout_amount":                2.0,
				"revenue_amount":               5.0,
				"payout_percentage":            5.0,
				"revenue_percentage":           10.0,
				"is_default":                   true,
				"is_allow_duplicate_conversion": true,
			},
		},
	}

	t.Log("Creating offer via direct Everflow API...")
	offerResp := config.EverflowAPIRequest(t, "POST", "/networks/offers", offerPayload)
	LogResponse(t, "Offer Creation", offerResp)

	// Check if the offer was created successfully
	if offerResp.StatusCode == 200 {
		var offerResult struct {
			NetworkOfferID int    `json:"network_offer_id"`
			Name          string `json:"name"`
		}
		ParseJSONResponse(t, offerResp, &offerResult)

		assert.Greater(t, offerResult.NetworkOfferID, 0, "Offer should have a valid ID")
		assert.Equal(t, "Test Offer Direct", offerResult.Name, "Offer name should match")

		t.Logf("✅ Successfully created offer with ID: %d", offerResult.NetworkOfferID)

		// Clean up - delete the offer
		t.Log("Cleaning up offer...")
		deleteResp := config.EverflowAPIRequest(t, "DELETE", 
			fmt.Sprintf("/networks/offers/%d", offerResult.NetworkOfferID), nil)
		if deleteResp.StatusCode != 200 {
			t.Logf("Warning: Failed to delete offer %d (status %d)", offerResult.NetworkOfferID, deleteResp.StatusCode)
		}

		// Clean up - delete the advertiser
		t.Log("Cleaning up advertiser...")
		deleteAdvResp := config.EverflowAPIRequest(t, "DELETE", 
			fmt.Sprintf("/networks/advertisers/%d", advertiserResult.NetworkAdvertiserID), nil)
		if deleteAdvResp.StatusCode != 200 {
			t.Logf("Warning: Failed to delete advertiser %d (status %d)", advertiserResult.NetworkAdvertiserID, deleteAdvResp.StatusCode)
		}
	} else {
		t.Logf("❌ Offer creation failed with status %d", offerResp.StatusCode)
		
		// Still clean up the advertiser
		t.Log("Cleaning up advertiser...")
		deleteAdvResp := config.EverflowAPIRequest(t, "DELETE", 
			fmt.Sprintf("/networks/advertisers/%d", advertiserResult.NetworkAdvertiserID), nil)
		if deleteAdvResp.StatusCode != 200 {
			t.Logf("Warning: Failed to delete advertiser %d (status %d)", advertiserResult.NetworkAdvertiserID, deleteAdvResp.StatusCode)
		}
		
		t.Fail()
	}
}