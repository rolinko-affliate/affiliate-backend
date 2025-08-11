package integration_tests

import (
	"fmt"
	"testing"
)

// CleanupTracker tracks entities created during tests for cleanup
type CleanupTracker struct {
	config              *TestConfig
	createdProfiles     []string // IDs of profiles created in our platform
	createdAdvertisers  []int64 // IDs of advertisers created in our platform
	createdAffiliates   []string // IDs of affiliates created in our platform
	createdCampaigns    []string // IDs of campaigns created in our platform
	createdTrackingLinks []string // IDs of tracking links created in our platform
	createdOrganizations []int64 // IDs of organizations created in our platform
	
	// Everflow entity IDs for cleanup
	everflowAdvertiserIDs []int // Everflow advertiser IDs
	everflowPartnerIDs    []int // Everflow partner IDs
	everflowOfferIDs      []int // Everflow offer IDs
	everflowTrackingLinkIDs []int // Everflow tracking link IDs
}

// NewCleanupTracker creates a new cleanup tracker
func NewCleanupTracker(config *TestConfig) *CleanupTracker {
	return &CleanupTracker{
		config:                  config,
		createdProfiles:         make([]string, 0),
		createdAdvertisers:      make([]int64, 0),
		createdAffiliates:       make([]string, 0),
		createdCampaigns:        make([]string, 0),
		createdTrackingLinks:    make([]string, 0),
		createdOrganizations:    make([]int64, 0),
		everflowAdvertiserIDs:   make([]int, 0),
		everflowPartnerIDs:      make([]int, 0),
		everflowOfferIDs:        make([]int, 0),
		everflowTrackingLinkIDs: make([]int, 0),
	}
}

// TrackProfile adds a profile ID to the cleanup list
func (ct *CleanupTracker) TrackProfile(profileID string) {
	ct.createdProfiles = append(ct.createdProfiles, profileID)
}

// TrackAdvertiser adds an advertiser ID to the cleanup list
func (ct *CleanupTracker) TrackAdvertiser(advertiserID int64) {
	ct.createdAdvertisers = append(ct.createdAdvertisers, advertiserID)
}

// TrackAffiliate adds an affiliate ID to the cleanup list
func (ct *CleanupTracker) TrackAffiliate(affiliateID string) {
	ct.createdAffiliates = append(ct.createdAffiliates, affiliateID)
}

// TrackCampaign adds a campaign ID to the cleanup list
func (ct *CleanupTracker) TrackCampaign(campaignID string) {
	ct.createdCampaigns = append(ct.createdCampaigns, campaignID)
}

// TrackTrackingLink adds a tracking link ID to the cleanup list
func (ct *CleanupTracker) TrackTrackingLink(trackingLinkID string) {
	ct.createdTrackingLinks = append(ct.createdTrackingLinks, trackingLinkID)
}

// TrackOrganization adds an organization ID to the cleanup list
func (ct *CleanupTracker) TrackOrganization(organizationID int64) {
	ct.createdOrganizations = append(ct.createdOrganizations, organizationID)
}

// TrackEverflowAdvertiser adds an Everflow advertiser ID to the cleanup list
func (ct *CleanupTracker) TrackEverflowAdvertiser(advertiserID int) {
	ct.everflowAdvertiserIDs = append(ct.everflowAdvertiserIDs, advertiserID)
}

// TrackEverflowPartner adds an Everflow partner ID to the cleanup list
func (ct *CleanupTracker) TrackEverflowPartner(partnerID int) {
	ct.everflowPartnerIDs = append(ct.everflowPartnerIDs, partnerID)
}

// TrackEverflowOffer adds an Everflow offer ID to the cleanup list
func (ct *CleanupTracker) TrackEverflowOffer(offerID int) {
	ct.everflowOfferIDs = append(ct.everflowOfferIDs, offerID)
}

// TrackEverflowTrackingLink adds an Everflow tracking link ID to the cleanup list
func (ct *CleanupTracker) TrackEverflowTrackingLink(trackingLinkID int) {
	ct.everflowTrackingLinkIDs = append(ct.everflowTrackingLinkIDs, trackingLinkID)
}

// Cleanup performs cleanup of all tracked entities
func (ct *CleanupTracker) Cleanup(t *testing.T) {
	t.Log("Starting cleanup of test entities...")

	// Clean up tracking links first (they depend on campaigns and affiliates)
	ct.cleanupTrackingLinks(t)
	
	// Clean up campaigns (they depend on advertisers)
	ct.cleanupCampaigns(t)
	
	// Clean up affiliates and advertisers
	ct.cleanupAffiliates(t)
	ct.cleanupAdvertisers(t)
	
	// Clean up organizations
	ct.cleanupOrganizations(t)
	
	// Clean up profiles last (they might be referenced by other entities)
	ct.cleanupProfiles(t)
	
	// Clean up Everflow entities
	ct.cleanupEverflowEntities(t)
	
	t.Log("Cleanup completed")
}

// cleanupTrackingLinks cleans up tracking links from our platform
func (ct *CleanupTracker) cleanupTrackingLinks(t *testing.T) {
	for _, trackingLinkID := range ct.createdTrackingLinks {
		// Note: We need organization ID to delete tracking links
		// This is a simplified approach - in real tests you'd track org IDs too
		t.Logf("Would cleanup tracking link: %s", trackingLinkID)
		// resp := ct.config.PlatformAPIRequest(t, "DELETE", fmt.Sprintf("/api/v1/organizations/{org_id}/tracking-links/%s", trackingLinkID), nil)
		// if resp.StatusCode != 200 && resp.StatusCode != 404 {
		//     t.Logf("Warning: Failed to cleanup tracking link %s: %d", trackingLinkID, resp.StatusCode)
		// }
	}
}

// cleanupCampaigns cleans up campaigns from our platform
func (ct *CleanupTracker) cleanupCampaigns(t *testing.T) {
	for _, campaignID := range ct.createdCampaigns {
		resp := ct.config.PlatformAPIRequest(t, "DELETE", fmt.Sprintf("/api/v1/campaigns/%s", campaignID), nil)
		if resp.StatusCode != 200 && resp.StatusCode != 404 {
			t.Logf("Warning: Failed to cleanup campaign %s: %d", campaignID, resp.StatusCode)
		}
	}
}

// cleanupAffiliates cleans up affiliates from our platform
func (ct *CleanupTracker) cleanupAffiliates(t *testing.T) {
	for _, affiliateID := range ct.createdAffiliates {
		resp := ct.config.PlatformAPIRequest(t, "DELETE", fmt.Sprintf("/api/v1/affiliates/%s", affiliateID), nil)
		if resp.StatusCode != 200 && resp.StatusCode != 404 {
			t.Logf("Warning: Failed to cleanup affiliate %s: %d", affiliateID, resp.StatusCode)
		}
	}
}

// cleanupAdvertisers cleans up advertisers from our platform
func (ct *CleanupTracker) cleanupAdvertisers(t *testing.T) {
	for _, advertiserID := range ct.createdAdvertisers {
		resp := ct.config.PlatformAPIRequest(t, "DELETE", fmt.Sprintf("/api/v1/advertisers/%d", advertiserID), nil)
		if resp.StatusCode != 200 && resp.StatusCode != 404 {
			t.Logf("Warning: Failed to cleanup advertiser %d: %d", advertiserID, resp.StatusCode)
		}
	}
}

// cleanupOrganizations cleans up organizations from our platform
func (ct *CleanupTracker) cleanupOrganizations(t *testing.T) {
	for _, organizationID := range ct.createdOrganizations {
		resp := ct.config.PlatformAPIRequest(t, "DELETE", fmt.Sprintf("/api/v1/organizations/%d", organizationID), nil)
		if resp.StatusCode != 200 && resp.StatusCode != 404 {
			t.Logf("Warning: Failed to cleanup organization %d: %d", organizationID, resp.StatusCode)
		}
	}
}

// cleanupProfiles cleans up profiles from our platform
func (ct *CleanupTracker) cleanupProfiles(t *testing.T) {
	for _, profileID := range ct.createdProfiles {
		resp := ct.config.PlatformAPIRequest(t, "DELETE", fmt.Sprintf("/api/v1/profiles/%s", profileID), nil)
		if resp.StatusCode != 200 && resp.StatusCode != 404 {
			t.Logf("Warning: Failed to cleanup profile %s: %d", profileID, resp.StatusCode)
		}
	}
}

// cleanupEverflowEntities cleans up entities from Everflow
func (ct *CleanupTracker) cleanupEverflowEntities(t *testing.T) {
	// Clean up Everflow tracking links
	for _, trackingLinkID := range ct.everflowTrackingLinkIDs {
		resp := ct.config.EverflowAPIRequest(t, "DELETE", fmt.Sprintf("/networks/tracking_links/%d", trackingLinkID), nil)
		if resp.StatusCode != 200 && resp.StatusCode != 404 {
			t.Logf("Warning: Failed to cleanup Everflow tracking link %d: %d", trackingLinkID, resp.StatusCode)
		}
	}
	
	// Clean up Everflow offers
	for _, offerID := range ct.everflowOfferIDs {
		resp := ct.config.EverflowAPIRequest(t, "DELETE", fmt.Sprintf("/networks/offers/%d", offerID), nil)
		if resp.StatusCode != 200 && resp.StatusCode != 404 {
			t.Logf("Warning: Failed to cleanup Everflow offer %d: %d", offerID, resp.StatusCode)
		}
	}
	
	// Clean up Everflow partners
	for _, partnerID := range ct.everflowPartnerIDs {
		resp := ct.config.EverflowAPIRequest(t, "DELETE", fmt.Sprintf("/networks/affiliates/%d", partnerID), nil)
		if resp.StatusCode != 200 && resp.StatusCode != 404 {
			t.Logf("Warning: Failed to cleanup Everflow partner %d: %d", partnerID, resp.StatusCode)
		}
	}
	
	// Clean up Everflow advertisers
	for _, advertiserID := range ct.everflowAdvertiserIDs {
		resp := ct.config.EverflowAPIRequest(t, "DELETE", fmt.Sprintf("/networks/advertisers/%d", advertiserID), nil)
		if resp.StatusCode != 200 && resp.StatusCode != 404 {
			t.Logf("Warning: Failed to cleanup Everflow advertiser %d: %d", advertiserID, resp.StatusCode)
		}
	}
}

