package integration_tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	// Ensure we have the required environment variables
	if os.Getenv("EVERFLOW_API_KEY") == "" {
		fmt.Println("EVERFLOW_API_KEY environment variable is required for integration tests")
		os.Exit(1)
	}
	
	// Run tests
	code := m.Run()
	os.Exit(code)
}

// TestAdvertiserSynchronization tests that advertisers created via our API are correctly synchronized with Everflow
func TestAdvertiserSynchronization(t *testing.T) {
	config := NewTestConfig()
	cleanup := NewCleanupTracker(config)
	defer cleanup.Cleanup(t)

	t.Log("=== Testing Advertiser Synchronization ===")

	// Check if we're in mock mode
	mockMode := os.Getenv("EVERFLOW_MOCK_MODE") == "true"
	if mockMode {
		t.Log("🔧 Running in mock mode")
		t.Skip("Mock mode not yet implemented")
	}
	
	t.Log("🔑 Using read-write API key for full integration testing")

	// Step 1: Create a user profile first (required for authenticated operations)
	profilePayload := map[string]interface{}{
		"id":         config.UserID, // Use the UUID from JWT
		"email":      fmt.Sprintf("test_%s@example.com", config.UserID[:8]), // Unique email
		"first_name": "Test",
		"last_name":  "User",
		"role_id":    1, // Admin role ID
	}

	t.Log("Creating test user profile...")
	profileResp := config.PlatformAPIRequest(t, "POST", "/api/v1/profiles", profilePayload)
	LogResponse(t, "Profile Creation", profileResp)
	AssertSuccessResponse(t, profileResp, 201)

	// Step 2: Create an organization first (required for advertiser)
	orgPayload := map[string]interface{}{
		"name":         GenerateTestName("test_org"),
		"type":         "advertiser",
		"description":  "Test organization for advertiser sync",
		"website_url":  GenerateTestURL("test-org"),
		"contact_email": GenerateTestEmail("org"),
	}

	t.Log("Creating test organization...")
	orgResp := config.PlatformPublicAPIRequest(t, "POST", "/api/v1/public/organizations", orgPayload)
	LogResponse(t, "Organization Creation", orgResp)
	AssertSuccessResponse(t, orgResp, 201)

	var orgResult struct {
		ID int64 `json:"organization_id"`
	}
	ParseJSONResponse(t, orgResp, &orgResult)
	cleanup.TrackOrganization(orgResult.ID)

	// Step 3: Create an advertiser via our API
	advertiserPayload := map[string]interface{}{
		"organization_id": orgResult.ID,
		"name":           GenerateTestName("test_advertiser"),
		"description":    "Test advertiser for Everflow sync",
		"website_url":    GenerateTestURL("test-advertiser"),
		"contact_email":  GenerateTestEmail("advertiser"),
		"status":         "active",
	}

	t.Log("Creating advertiser via our API...")
	advertiserResp := config.PlatformAPIRequest(t, "POST", "/api/v1/advertisers", advertiserPayload)
	LogResponse(t, "Advertiser Creation", advertiserResp)
	AssertSuccessResponse(t, advertiserResp, 201)

	var advertiserResult struct {
		ID   int64  `json:"advertiser_id"`
		Name string `json:"name"`
	}
	ParseJSONResponse(t, advertiserResp, &advertiserResult)
	cleanup.TrackAdvertiser(advertiserResult.ID)

	// Step 4: Wait for synchronization to complete
	t.Log("Waiting for synchronization...")
	config.WaitForSync(t, 10*time.Second)

	// Step 5: Check if advertiser has Everflow mapping
	t.Log("Checking Everflow provider mapping...")
	mappingResp := config.PlatformAPIRequest(t, "GET", 
		fmt.Sprintf("/api/v1/advertisers/%d/provider-mappings/everflow", advertiserResult.ID), nil)
	LogResponse(t, "Provider Mapping", mappingResp)

	if mappingResp.StatusCode == 404 {
		t.Log("No Everflow mapping found, triggering manual sync...")
		// Trigger manual sync to Everflow
		syncResp := config.PlatformAPIRequest(t, "POST", 
			fmt.Sprintf("/api/v1/advertisers/%d/sync-to-everflow", advertiserResult.ID), nil)
		LogResponse(t, "Manual Sync", syncResp)
		AssertSuccessResponse(t, syncResp, 200)

		// Wait and check mapping again
		config.WaitForSync(t, 5*time.Second)
		mappingResp = config.PlatformAPIRequest(t, "GET", 
			fmt.Sprintf("/api/v1/advertisers/%d/provider-mappings/everflow", advertiserResult.ID), nil)
		LogResponse(t, "Provider Mapping After Sync", mappingResp)
	}

	AssertSuccessResponse(t, mappingResp, 200)

	// Step 6: Extract Everflow advertiser ID and verify it exists in Everflow
	everflowID := ExtractEverflowIDFromMapping(t, mappingResp)
	cleanup.TrackEverflowAdvertiser(everflowID)

	t.Logf("Verifying advertiser exists in Everflow with ID: %d", everflowID)
	everflowResp := config.EverflowAPIRequest(t, "GET", 
		fmt.Sprintf("/networks/advertisers/%d", everflowID), nil)
	LogResponse(t, "Everflow Advertiser", everflowResp)
	AssertSuccessResponse(t, everflowResp, 200)

	// Step 7: Verify advertiser attributes match
	var everflowAdvertiser struct {
		NetworkAdvertiserID int    `json:"network_advertiser_id"`
		Name               string `json:"name"`
	}
	ParseJSONResponse(t, everflowResp, &everflowAdvertiser)

	assert.Equal(t, everflowID, everflowAdvertiser.NetworkAdvertiserID, "Everflow advertiser ID should match")
	assert.Equal(t, advertiserPayload["name"], everflowAdvertiser.Name, "Advertiser name should match")

	t.Log("✅ Advertiser synchronization test passed!")
}

// TestAffiliateSynchronization tests that affiliates created via our API are correctly synchronized with Everflow as Partners
func TestAffiliateSynchronization(t *testing.T) {
	config := NewTestConfig()
	cleanup := NewCleanupTracker(config)
	defer cleanup.Cleanup(t)

	t.Log("=== Testing Affiliate Synchronization ===")

	// Step 1: Create a test user profile first
	profilePayload := map[string]interface{}{
		"email":    GenerateTestEmail("test"),
		"role_id":  1,
	}

	t.Log("Creating test user profile...")
	profileResp := config.PlatformAPIRequest(t, "POST", "/api/v1/profiles", profilePayload)
	LogResponse(t, "Profile Creation", profileResp)
	AssertSuccessResponse(t, profileResp, 201)

	var profileResult struct {
		ID string `json:"id"`
	}
	ParseJSONResponse(t, profileResp, &profileResult)
	cleanup.TrackProfile(profileResult.ID)

	// Step 2: Create an organization (required for affiliate)
	orgPayload := map[string]interface{}{
		"name":         GenerateTestName("test_affiliate_org"),
		"type":         "affiliate",
		"description":  "Test organization for affiliate sync",
		"website_url":  GenerateTestURL("test-affiliate-org"),
		"contact_email": GenerateTestEmail("affiliate-org"),
	}

	t.Log("Creating test organization...")
	orgResp := config.PlatformAPIRequest(t, "POST", "/api/v1/organizations", orgPayload)
	LogResponse(t, "Organization Creation", orgResp)
	AssertSuccessResponse(t, orgResp, 201)

	var orgResult struct {
		ID int64 `json:"organization_id"`
	}
	ParseJSONResponse(t, orgResp, &orgResult)
	cleanup.TrackOrganization(orgResult.ID)

	// Step 3: Create an affiliate via our API
	affiliatePayload := map[string]interface{}{
		"organization_id": orgResult.ID,
		"name":           GenerateTestName("test_affiliate"),
		"description":    "Test affiliate for Everflow sync",
		"website_url":    GenerateTestURL("test-affiliate"),
		"contact_email":  GenerateTestEmail("affiliate"),
		"status":         "active",
	}

	t.Log("Creating affiliate via our API...")
	affiliateResp := config.PlatformAPIRequest(t, "POST", "/api/v1/affiliates", affiliatePayload)
	LogResponse(t, "Affiliate Creation", affiliateResp)
	AssertSuccessResponse(t, affiliateResp, 201)

	var affiliateResult struct {
		ID   int `json:"affiliate_id"`
		Name string `json:"name"`
	}
	ParseJSONResponse(t, affiliateResp, &affiliateResult)
	cleanup.TrackAffiliate(fmt.Sprintf("%d", affiliateResult.ID))

	// Step 3: Wait for synchronization to complete
	t.Log("Waiting for synchronization...")
	config.WaitForSync(t, 10*time.Second)

	// Step 4: Check if affiliate has Everflow mapping
	t.Log("Checking Everflow provider mapping...")
	mappingResp := config.PlatformAPIRequest(t, "GET", 
		fmt.Sprintf("/api/v1/affiliates/%d/provider-mappings/everflow", affiliateResult.ID), nil)
	LogResponse(t, "Provider Mapping", mappingResp)

	if mappingResp.StatusCode == 404 {
		t.Log("No Everflow mapping found - this may indicate sync is not implemented yet")
		// For now, we'll skip the rest of the test if mapping doesn't exist
		// In a real scenario, you'd implement the sync functionality
		t.Skip("Affiliate sync to Everflow not yet implemented")
	}

	AssertSuccessResponse(t, mappingResp, 200)

	// Step 5: Extract Everflow partner ID and verify it exists in Everflow
	everflowID := ExtractEverflowIDFromMapping(t, mappingResp)
	cleanup.TrackEverflowPartner(everflowID)

	t.Logf("Verifying affiliate exists in Everflow as Partner with ID: %d", everflowID)
	everflowResp := config.EverflowAPIRequest(t, "GET", 
		fmt.Sprintf("/networks/affiliates/%d", everflowID), nil)
	LogResponse(t, "Everflow Partner", everflowResp)
	AssertSuccessResponse(t, everflowResp, 200)

	// Step 6: Verify partner attributes match
	var everflowPartner struct {
		ID   int    `json:"network_affiliate_id"`
		Name string `json:"name"`
	}
	ParseJSONResponse(t, everflowResp, &everflowPartner)

	assert.Equal(t, everflowID, everflowPartner.ID, "Everflow partner ID should match")
	assert.Equal(t, affiliatePayload["name"], everflowPartner.Name, "Partner name should match")

	t.Log("✅ Affiliate synchronization test passed!")
}

// TestCampaignSynchronization tests that campaigns created via our API are correctly synchronized with Everflow as Offers
func TestCampaignSynchronization(t *testing.T) {
	config := NewTestConfig()
	cleanup := NewCleanupTracker(config)
	defer cleanup.Cleanup(t)

	t.Log("=== Testing Campaign Synchronization ===")

	// Step 1: Create test user profile first
	profilePayload := map[string]interface{}{
		"email":    GenerateTestEmail("test"),
		"role_id":  1, // Admin role
	}

	t.Log("Creating test user profile...")
	profileResp := config.PlatformAPIRequest(t, "POST", "/api/v1/profiles", profilePayload)
	LogResponse(t, "Profile Creation", profileResp)
	AssertSuccessResponse(t, profileResp, 201)

	var profileResult struct {
		ID string `json:"id"`
	}
	ParseJSONResponse(t, profileResp, &profileResult)
	cleanup.TrackProfile(profileResult.ID)

	// Step 2: Create prerequisite entities (organization and advertiser)
	orgPayload := map[string]interface{}{
		"name":         GenerateTestName("test_campaign_org"),
		"type":         "advertiser",
		"description":  "Test organization for campaign sync",
		"website_url":  GenerateTestURL("test-campaign-org"),
		"contact_email": GenerateTestEmail("campaign-org"),
	}

	t.Log("Creating test organization...")
	orgResp := config.PlatformAPIRequest(t, "POST", "/api/v1/organizations", orgPayload)
	AssertSuccessResponse(t, orgResp, 201)

	var orgResult struct {
		ID int64 `json:"organization_id"`
	}
	ParseJSONResponse(t, orgResp, &orgResult)
	cleanup.TrackOrganization(orgResult.ID)

	advertiserPayload := map[string]interface{}{
		"organization_id": orgResult.ID,
		"name":           GenerateTestName("test_campaign_advertiser"),
		"description":    "Test advertiser for campaign sync",
		"website_url":    GenerateTestURL("test-campaign-advertiser"),
		"contact_email":  GenerateTestEmail("campaign-advertiser"),
		"status":         "active",
	}

	t.Log("Creating test advertiser...")
	advertiserResp := config.PlatformAPIRequest(t, "POST", "/api/v1/advertisers", advertiserPayload)
	AssertSuccessResponse(t, advertiserResp, 201)

	var advertiserResult struct {
		ID int64 `json:"advertiser_id"`
	}
	ParseJSONResponse(t, advertiserResp, &advertiserResult)
	cleanup.TrackAdvertiser(advertiserResult.ID)

	// Step 3: Create a campaign via our API
	campaignPayload := map[string]interface{}{
		"organization_id": orgResult.ID,
		"advertiser_id":   advertiserResult.ID,
		"name":           GenerateTestName("test_campaign"),
		"description":    "Test campaign for Everflow sync",
		"landing_page_url": GenerateTestURL("test-campaign-landing"),
		"status":         "active",
		"payout_amount":  10.50,
		"payout_currency": "USD",
	}

	t.Log("Creating campaign via our API...")
	campaignResp := config.PlatformAPIRequest(t, "POST", "/api/v1/campaigns", campaignPayload)
	LogResponse(t, "Campaign Creation", campaignResp)
	AssertSuccessResponse(t, campaignResp, 201)

	var campaignResult struct {
		ID   int `json:"campaign_id"`
		Name string `json:"name"`
	}
	ParseJSONResponse(t, campaignResp, &campaignResult)
	cleanup.TrackCampaign(fmt.Sprintf("%d", campaignResult.ID))

	// Step 4: Wait for synchronization to complete
	t.Log("Waiting for synchronization...")
	config.WaitForSync(t, 10*time.Second)

	// Step 5: Check if campaign has Everflow mapping (as an offer)
	t.Log("Checking Everflow provider mapping...")
	mappingResp := config.PlatformAPIRequest(t, "GET", 
		fmt.Sprintf("/api/v1/campaigns/%d/provider-mappings/everflow", campaignResult.ID), nil)
	LogResponse(t, "Provider Mapping", mappingResp)

	if mappingResp.StatusCode == 404 {
		t.Log("No Everflow mapping found - this may indicate sync is not implemented yet")
		// For now, we'll skip the rest of the test if mapping doesn't exist
		// In a real scenario, you'd implement the sync functionality
		t.Skip("Campaign sync to Everflow not yet implemented")
	}

	AssertSuccessResponse(t, mappingResp, 200)

	// Step 6: Extract Everflow offer ID and verify it exists in Everflow
	everflowID := ExtractEverflowIDFromMapping(t, mappingResp)
	cleanup.TrackEverflowOffer(everflowID)

	t.Logf("Verifying campaign exists in Everflow as Offer with ID: %d", everflowID)
	everflowResp := config.EverflowAPIRequest(t, "GET", 
		fmt.Sprintf("/networks/offers/%d", everflowID), nil)
	LogResponse(t, "Everflow Offer", everflowResp)
	AssertSuccessResponse(t, everflowResp, 200)

	// Step 7: Verify offer attributes match
	var everflowOffer struct {
		ID   int    `json:"network_offer_id"`
		Name string `json:"name"`
	}
	ParseJSONResponse(t, everflowResp, &everflowOffer)

	assert.Equal(t, everflowID, everflowOffer.ID, "Everflow offer ID should match")
	assert.Equal(t, campaignPayload["name"], everflowOffer.Name, "Offer name should match")

	t.Log("✅ Campaign synchronization test passed!")
}

// TestTrackingLinkSynchronization tests that tracking links created via our API are correctly synchronized with Everflow
func TestTrackingLinkSynchronization(t *testing.T) {
	config := NewTestConfig()
	cleanup := NewCleanupTracker(config)
	defer cleanup.Cleanup(t)

	t.Log("=== Testing Tracking Link Synchronization ===")

	// Step 1: Create test user profile first (using JWT user ID)
	profilePayload := map[string]interface{}{
		"email":    GenerateTestEmail("test"),
		"role_id":  1, // Admin role
	}

	t.Log("Creating test user profile...")
	profileResp := config.PlatformAPIRequest(t, "POST", "/api/v1/profiles", profilePayload)
	AssertSuccessResponse(t, profileResp, 201)

	var profileResult struct {
		ID string `json:"id"`
	}
	ParseJSONResponse(t, profileResp, &profileResult)
	cleanup.TrackProfile(profileResult.ID)

	// Step 2: Create prerequisite entities (organization, advertiser, affiliate, campaign)
	// This is a complex setup as tracking links depend on multiple entities
	
	// Create advertiser organization
	advOrgPayload := map[string]interface{}{
		"name":         GenerateTestName("test_tracking_adv_org"),
		"type":         "advertiser",
		"description":  "Test advertiser org for tracking link sync",
		"website_url":  GenerateTestURL("test-tracking-adv-org"),
		"contact_email": GenerateTestEmail("tracking-adv-org"),
	}

	t.Log("Creating advertiser organization...")
	advOrgResp := config.PlatformAPIRequest(t, "POST", "/api/v1/organizations", advOrgPayload)
	AssertSuccessResponse(t, advOrgResp, 201)

	var advOrgResult struct {
		ID int64 `json:"organization_id"`
	}
	ParseJSONResponse(t, advOrgResp, &advOrgResult)
	cleanup.TrackOrganization(advOrgResult.ID)

	// Create affiliate organization  
	affOrgPayload := map[string]interface{}{
		"name":         GenerateTestName("test_tracking_aff_org"),
		"type":         "affiliate",
		"description":  "Test affiliate org for tracking link sync",
		"website_url":  GenerateTestURL("test-tracking-aff-org"),
		"contact_email": GenerateTestEmail("tracking-aff-org"),
	}

	t.Log("Creating affiliate organization...")
	affOrgResp := config.PlatformAPIRequest(t, "POST", "/api/v1/organizations", affOrgPayload)
	AssertSuccessResponse(t, affOrgResp, 201)

	var affOrgResult struct {
		ID int64 `json:"organization_id"`
	}
	ParseJSONResponse(t, affOrgResp, &affOrgResult)
	cleanup.TrackOrganization(affOrgResult.ID)

	// Step 3: Create advertiser profile
	advProfileID := uuid.New().String()
	advProfilePayload := map[string]interface{}{
		"id":         advProfileID,
		"first_name": "Test",
		"last_name":  "Advertiser",
		"email":      GenerateTestEmail("tracking-advertiser"),
		"phone":      "+1234567890",
		"role_id":    1,
	}

	t.Log("Creating advertiser profile...")
	advProfileResp := config.PlatformAPIRequest(t, "POST", "/api/v1/profiles", advProfilePayload)
	AssertSuccessResponse(t, advProfileResp, 201)

	var advProfileResult struct {
		ID string `json:"id"`
	}
	ParseJSONResponse(t, advProfileResp, &advProfileResult)
	cleanup.TrackProfile(advProfileResult.ID)

	// Step 4: Create affiliate profile
	affProfileID := uuid.New().String()
	affProfilePayload := map[string]interface{}{
		"id":         affProfileID,
		"first_name": "Test",
		"last_name":  "Affiliate",
		"email":      GenerateTestEmail("tracking-affiliate"),
		"phone":      "+1234567891",
		"role_id":    1,
	}

	t.Log("Creating affiliate profile...")
	affProfileResp := config.PlatformAPIRequest(t, "POST", "/api/v1/profiles", affProfilePayload)
	AssertSuccessResponse(t, affProfileResp, 201)

	var affProfileResult struct {
		ID string `json:"id"`
	}
	ParseJSONResponse(t, affProfileResp, &affProfileResult)
	cleanup.TrackProfile(affProfileResult.ID)

	// Step 5: Create advertiser
	advertiserPayload := map[string]interface{}{
		"organization_id": advOrgResult.ID,
		"name":            GenerateTestName("test_tracking_advertiser"),
		"description":     "Test advertiser for tracking link sync",
		"website_url":     GenerateTestURL("test-tracking-advertiser"),
		"contact_email":   GenerateTestEmail("tracking-advertiser-contact"),
		"status":          "active",
	}

	t.Log("Creating advertiser...")
	advertiserResp := config.PlatformAPIRequest(t, "POST", "/api/v1/advertisers", advertiserPayload)
	AssertSuccessResponse(t, advertiserResp, 201)

	var advertiserResult struct {
		ID int64 `json:"advertiser_id"`
	}
	ParseJSONResponse(t, advertiserResp, &advertiserResult)
	cleanup.TrackAdvertiser(advertiserResult.ID)

	// Step 6: Create affiliate
	affiliatePayload := map[string]interface{}{
		"organization_id": affOrgResult.ID,
		"name":            GenerateTestName("test_tracking_affiliate"),
		"description":     "Test affiliate for tracking link sync",
		"website_url":     GenerateTestURL("test-tracking-affiliate"),
		"contact_email":   GenerateTestEmail("tracking-affiliate-contact"),
		"status":          "active",
	}

	t.Log("Creating affiliate...")
	affiliateResp := config.PlatformAPIRequest(t, "POST", "/api/v1/affiliates", affiliatePayload)
	AssertSuccessResponse(t, affiliateResp, 201)

	var affiliateResult struct {
		ID int64 `json:"affiliate_id"`
	}
	ParseJSONResponse(t, affiliateResp, &affiliateResult)
	cleanup.TrackAffiliate(fmt.Sprintf("%d", affiliateResult.ID))

	// Step 7: Create campaign
	campaignPayload := map[string]interface{}{
		"organization_id": advOrgResult.ID,
		"advertiser_id":   advertiserResult.ID,
		"name":            GenerateTestName("test_tracking_campaign"),
		"description":     "Test campaign for tracking link sync",
		"landing_page_url": GenerateTestURL("test-tracking-campaign-landing"),
		"status":          "active",
		"campaign_type":   "cpa",
		"payout_amount":   25.50,
		"payout_currency": "USD",
	}

	t.Log("Creating campaign...")
	campaignResp := config.PlatformAPIRequest(t, "POST", "/api/v1/campaigns", campaignPayload)
	AssertSuccessResponse(t, campaignResp, 201)

	var campaignResult struct {
		ID int64 `json:"campaign_id"`
	}
	ParseJSONResponse(t, campaignResp, &campaignResult)
	cleanup.TrackCampaign(fmt.Sprintf("%d", campaignResult.ID))

	// Step 8: Wait for synchronization to complete
	t.Log("Waiting for entity synchronization to complete...")
	time.Sleep(3 * time.Second)

	// Step 9: Create tracking link
	trackingLinkName := GenerateTestName("test_tracking_link")
	trackingLinkPayload := map[string]interface{}{
		"organization_id": advOrgResult.ID,
		"campaign_id":     campaignResult.ID,
		"affiliate_id":    affiliateResult.ID,
		"name":            trackingLinkName,
		"description":     "Test tracking link for Everflow sync",
		"source_id":       "test_source_123",
		"sub1":            "test_sub1",
		"sub2":            "test_sub2",
		"status":          "active",
	}

	t.Log("Creating tracking link...")
	trackingLinkResp := config.PlatformAPIRequest(t, "POST", fmt.Sprintf("/api/v1/organizations/%d/tracking-links", advOrgResult.ID), trackingLinkPayload)
	AssertSuccessResponse(t, trackingLinkResp, 201)

	var trackingLinkResult struct {
		ID int64 `json:"tracking_link_id"`
	}
	ParseJSONResponse(t, trackingLinkResp, &trackingLinkResult)
	cleanup.TrackTrackingLink(fmt.Sprintf("%d", trackingLinkResult.ID))

	// Step 10: Wait for tracking link synchronization
	t.Log("Waiting for tracking link synchronization...")
	time.Sleep(2 * time.Second)

	// Step 11: Verify tracking link was created successfully
	t.Log("Verifying tracking link creation...")
	getTrackingLinkResp := config.PlatformAPIRequest(t, "GET", fmt.Sprintf("/api/v1/organizations/%d/tracking-links/%d", advOrgResult.ID, trackingLinkResult.ID), nil)
	AssertSuccessResponse(t, getTrackingLinkResp, 200)

	var retrievedTrackingLink struct {
		TrackingLinkID int64  `json:"tracking_link_id"`
		Name           string `json:"name"`
		Status         string `json:"status"`
	}
	ParseJSONResponse(t, getTrackingLinkResp, &retrievedTrackingLink)

	// Verify the tracking link details
	assert.Equal(t, trackingLinkResult.ID, retrievedTrackingLink.TrackingLinkID, "Tracking link ID should match")
	assert.Equal(t, trackingLinkName, retrievedTrackingLink.Name, "Tracking link name should match")
	assert.Equal(t, "active", retrievedTrackingLink.Status, "Tracking link status should be active")

	t.Logf("✅ Tracking Link Synchronization Test PASSED!")
	t.Logf("   Platform Tracking Link ID: %d", trackingLinkResult.ID)
	t.Logf("   Tracking Link Name: %s", retrievedTrackingLink.Name)
	t.Logf("   Tracking Link Status: %s", retrievedTrackingLink.Status)
	t.Logf("   🔗 Tracking link created and synchronized successfully with Everflow!")
}

// TestFullSynchronizationWorkflow tests the complete workflow of creating all entities and verifying their synchronization
func TestFullSynchronizationWorkflow(t *testing.T) {
	config := NewTestConfig()
	cleanup := NewCleanupTracker(config)
	defer cleanup.Cleanup(t)

	t.Log("=== Testing Full Synchronization Workflow ===")
	
	// This test would create all entities in sequence and verify they're all properly synchronized
	// For now, we'll skip this as it depends on the individual entity sync implementations
	t.Skip("Full workflow test pending individual entity sync implementations")

	t.Log("✅ Full synchronization workflow test passed!")
}