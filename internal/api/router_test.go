package api

import (
	"testing"

	"github.com/affiliate-backend/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

// TestRouterSetup tests that the router can be set up without panicking
func TestRouterSetup(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a minimal RouterOptions with nil handlers
	opts := RouterOptions{
		ProfileHandler:      &handlers.ProfileHandler{},
		OrganizationHandler: &handlers.OrganizationHandler{},
		AdvertiserHandler:   &handlers.AdvertiserHandler{},
		AffiliateHandler:    &handlers.AffiliateHandler{},
		CampaignHandler:     &handlers.CampaignHandler{},
	}

	// This should not panic if our route fixes are correct
	router := SetupRouter(opts)
	if router == nil {
		t.Error("Expected router to be non-nil")
	}
}
