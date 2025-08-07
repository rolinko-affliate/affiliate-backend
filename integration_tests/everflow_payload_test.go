package integration_tests

import (
	"testing"
)

// TestEverflowDirectAdvertiserCreation tests creating an advertiser directly with the working example payload
func TestEverflowDirectAdvertiserCreation(t *testing.T) {
	config := NewTestConfig()

	t.Log("=== Testing Direct Everflow Advertiser Creation ===")

	// Use the exact working payload you provided
	workingPayload := map[string]interface{}{
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
		"name":                       GenerateTestName("test_advertiser_direct"),
		"network_employee_id":        1,
		"offer_id_macro":             "",
		"platform_name":              "",
		"platform_url":               "",
		"platform_username":          "",
		"reporting_timezone_id":      80,
		"sales_manager_id":           1,
		"settings": map[string]interface{}{
			"exposed_variables": map[string]interface{}{
				"affiliate":    false,
				"affiliate_id": true,
				"offer_url":    false,
				"source_id":    false,
				"sub1":         true,
				"sub2":         true,
				"sub3":         false,
				"sub4":         false,
				"sub5":         false,
			},
		},
		"users": []map[string]interface{}{
			{
				"account_status":     "active",
				"currency_id":        "USD",
				"email":              GenerateTestEmail("advertiser_user"),
				"first_name":         "john",
				"initial_password":   "",
				"language_id":        1,
				"last_name":          "smith",
				"timezone_id":        80,
			},
		},
		"verification_token": "c7HIWpFUGnyQfN5wwBollBBGtUkeOm",
	}

	t.Log("Testing direct Everflow API call with working payload...")
	resp := callEverflowAPI(t, config, "POST", "/v1/networks/advertisers", workingPayload)
	
	t.Logf("Direct API Response: Status=%d", resp.StatusCode)
	t.Logf("Response Body: %s", resp.Body)

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		t.Log("Direct advertiser creation successful!")
		
		// TODO: Parse response and cleanup the created advertiser
		// For now, just log success
		
	} else {
		t.Logf("Direct advertiser creation failed with status %d", resp.StatusCode)
		t.Logf("Response: %s", resp.Body)
	}
}

// TestEverflowDirectAffiliateCreation tests creating an affiliate directly with the working example payload
func TestEverflowDirectAffiliateCreation(t *testing.T) {
	config := NewTestConfig()

	t.Log("=== Testing Direct Everflow Affiliate Creation ===")

	// Use the exact working payload you provided for affiliates
	workingPayload := map[string]interface{}{
		"account_status": "active",
		"billing": map[string]interface{}{
			"billing_frequency": "monthly",
			"details": map[string]interface{}{
				"day_of_month": 1,
			},
			"payment_type": "none",
			"tax_id":       "XXXXX",
		},
		"default_currency_id":                "USD",
		"enable_media_cost_tracking_links":   false,
		"internal_notes":                     "This is a test affiliate created using the API",
		"is_contact_address_enabled":         true,
		"name":                               GenerateTestName("test_affiliate_direct"),
		"network_employee_id":                1,
		"referrer_id":                        0,
	}

	t.Log("Testing direct Everflow API call with working affiliate payload...")
	resp := callEverflowAPI(t, config, "POST", "/v1/networks/affiliates", workingPayload)
	
	t.Logf("Direct Affiliate API Response: Status=%d", resp.StatusCode)
	t.Logf("Response Body: %s", resp.Body)

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		t.Log("Direct affiliate creation successful!")
	} else {
		t.Logf("Direct affiliate creation failed with status %d", resp.StatusCode)
		t.Logf("Response: %s", resp.Body)
	}
}