package integration_tests

import (
	"fmt"
	"testing"
	"time"
)

// TestEverflowDirectAPI tests calling Everflow API directly to understand the issue
func TestEverflowDirectAPI(t *testing.T) {
	config := NewTestConfig()
	
	if config.EverflowAPIKey == "" {
		t.Skip("EVERFLOW_API_KEY not set")
	}

	t.Log("=== Testing Direct Everflow API Call ===")

	// Test 1: Try just the absolute minimum required fields
	timestamp := time.Now().Unix()
	minimalPayload := map[string]interface{}{
		"name": fmt.Sprintf("test_advertiser_direct_%d", timestamp),
	}
	
	t.Log("Testing absolute minimal advertiser creation...")
	resp := callEverflowAPI(t, config, "POST", "/v1/networks/advertisers", minimalPayload)
	t.Logf("Absolute minimal request response: Status=%d, Body=%s", resp.StatusCode, resp.Body)
	
	// Test 2: Add required fields one by one
	basicPayload := map[string]interface{}{
		"name":                 fmt.Sprintf("test_advertiser_basic_%d", timestamp),
		"account_status":       "active",
		"network_employee_id":  1,
		"default_currency_id":  "USD",
		"reporting_timezone_id": 80,
	}

	t.Log("Testing basic advertiser creation...")
	resp2 := callEverflowAPI(t, config, "POST", "/v1/networks/advertisers", basicPayload)
	t.Logf("Basic request response: Status=%d, Body=%s", resp2.StatusCode, resp2.Body)

	// Test 3: Try creating an affiliate to check API key permissions
	if resp.StatusCode != 201 && resp2.StatusCode != 201 {
		t.Log("Advertiser creation failed, trying affiliate creation to check API key permissions...")
		
		affiliatePayload := map[string]interface{}{
			"name":                 fmt.Sprintf("test_affiliate_%d", timestamp),
			"account_status":       "active",
			"network_employee_id":  1,
			"default_currency_id":  "USD",
			"reporting_timezone_id": 80,
		}
		
		resp3 := callEverflowAPI(t, config, "POST", "/v1/networks/affiliates", affiliatePayload)
		t.Logf("Affiliate creation response: Status=%d, Body=%s", resp3.StatusCode, resp3.Body)

		t.Log("Getting advertisers list...")
		advertisersResp := callEverflowAPI(t, config, "GET", "/v1/networks/advertisers", nil)
		t.Logf("Advertisers response: Status=%d, Body=%s", advertisersResp.StatusCode, advertisersResp.Body)
	}
}

