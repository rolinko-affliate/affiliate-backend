package integration_tests

import (
	"encoding/json"
	"testing"
)

// TestDebugPlatformAdvertiserPayload captures what our platform is actually sending to Everflow
func TestDebugPlatformAdvertiserPayload(t *testing.T) {
	config := NewTestConfig()
	cleanup := NewCleanupTracker(config)
	defer cleanup.Cleanup(t)

	t.Log("=== Debugging Platform Advertiser Payload ===")

	// Step 1: Create a user profile first
	profilePayload := map[string]interface{}{
		"id":         config.UserID,
		"email":      GenerateTestEmail("debug_user"),
		"first_name": "Debug",
		"last_name":  "User",
		"role_id":    1,
	}

	t.Log("Creating test user profile...")
	profileResp := config.PlatformAPIRequest(t, "POST", "/api/v1/profiles", profilePayload)
	LogResponse(t, "Profile Creation", profileResp)
	AssertSuccessResponse(t, profileResp, 201)
	// Profile cleanup is handled automatically

	// Step 2: Create organization
	orgPayload := map[string]interface{}{
		"name": GenerateTestName("debug_org"),
		"type": "advertiser",
	}

	t.Log("Creating test organization...")
	orgResp := config.PlatformAPIRequest(t, "POST", "/api/v1/organizations", orgPayload)
	LogResponse(t, "Organization Creation", orgResp)
	AssertSuccessResponse(t, orgResp, 201)

	var orgResult struct {
		OrganizationID int64 `json:"organization_id"`
	}
	ParseJSONResponse(t, orgResp, &orgResult)
	cleanup.TrackOrganization(orgResult.OrganizationID)

	// Step 3: Create advertiser with minimal payload to see what gets sent to Everflow
	advertiserPayload := map[string]interface{}{
		"organization_id": orgResult.OrganizationID,
		"name":           "Debug Test Advertiser",
		"description":    "Debug test advertiser",
		"website_url":    "https://debug-test.com",
		"contact_email":  "debug@test.com",
		"status":         "active",
	}

	t.Log("Creating advertiser via our API to debug payload...")
	
	// We expect this to fail, but we want to see what payload was sent
	advertiserResp := config.PlatformAPIRequest(t, "POST", "/api/v1/advertisers", advertiserPayload)
	LogResponse(t, "Advertiser Creation (Debug)", advertiserResp)
	
	// Parse the error message to see if we can extract the payload that was sent
	if advertiserResp.StatusCode == 500 {
		var errorResp struct {
			Error string `json:"error"`
		}
		
		err := json.Unmarshal(advertiserResp.Body, &errorResp)
		if err == nil {
			t.Logf("Error message: %s", errorResp.Error)
			
			// The error message should contain details about what was sent to Everflow
			// This will help us understand the format difference
		}
	}
	
	t.Log("This test is for debugging - failure is expected until we fix the payload format")
}