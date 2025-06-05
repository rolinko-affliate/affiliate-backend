package api

import (
	"github.com/affiliate-backend/internal/api/handlers"
	"github.com/affiliate-backend/internal/api/middleware"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// RouterOptions contains dependencies for the router
type RouterOptions struct {
	ProfileHandler       *handlers.ProfileHandler
	ProfileService       service.ProfileService
	OrganizationHandler  *handlers.OrganizationHandler
	AdvertiserHandler    *handlers.AdvertiserHandler
	AffiliateHandler     *handlers.AffiliateHandler
	CampaignHandler      *handlers.CampaignHandler
}

// SetupRouter sets up the API router
func SetupRouter(opts RouterOptions) *gin.Engine {
	r := gin.Default() // Starts with Logger and Recovery middleware
	
	// Apply CORS middleware (will only allow CORS in development)
	r.Use(middleware.CORSMiddleware())

	// Health Check
	r.GET("/health", handlers.HealthCheck)

	// Public routes (e.g., Supabase webhook for profile creation)
	public := r.Group("/api/v1/public")
	{
		public.POST("/webhooks/supabase/new-user", opts.ProfileHandler.HandleSupabaseNewUserWebhook)
	}

	// Authenticated routes
	v1 := r.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware()) // Apply JWT auth middleware

	// Create RBAC middleware factory
	rbacMW := func(allowedRoles ...string) gin.HandlerFunc {
		return middleware.RBACMiddleware(opts.ProfileService, allowedRoles...)
	}

	// --- Profile Routes ---
	v1.GET("/users/me", opts.ProfileHandler.GetMyProfile)
		
	// Profile management routes
	profiles := v1.Group("/profiles")
	{
		profiles.POST("", opts.ProfileHandler.CreateProfile)
		profiles.POST("/upsert", opts.ProfileHandler.UpsertProfile)
		profiles.PUT("/:id", opts.ProfileHandler.UpdateProfile)
		profiles.DELETE("/:id", opts.ProfileHandler.DeleteProfile)
	}

	// --- Organization Routes ---
	organizations := v1.Group("/organizations")
	organizations.Use(rbacMW("Admin")) // Only admins can manage organizations
	{
		organizations.POST("", opts.OrganizationHandler.CreateOrganization)
		organizations.GET("", opts.OrganizationHandler.ListOrganizations)
		organizations.GET("/:id", opts.OrganizationHandler.GetOrganization)
		organizations.PUT("/:id", opts.OrganizationHandler.UpdateOrganization)
		organizations.DELETE("/:id", opts.OrganizationHandler.DeleteOrganization)
		
		// Organization's advertisers
		organizations.GET("/:id/advertisers", opts.AdvertiserHandler.ListAdvertisersByOrganization)
		
		// Organization's affiliates
		organizations.GET("/:id/affiliates", opts.AffiliateHandler.ListAffiliatesByOrganization)
		
		// Organization's campaigns
		organizations.GET("/:id/campaigns", opts.CampaignHandler.ListCampaignsByOrganization)
	}

	// --- Advertiser Routes ---
	advertisers := v1.Group("/advertisers")
	advertisers.Use(rbacMW("AdvertiserManager", "Admin"))
	{
		advertisers.POST("", opts.AdvertiserHandler.CreateAdvertiser)
		advertisers.GET("/:id", opts.AdvertiserHandler.GetAdvertiser)
		advertisers.PUT("/:id", opts.AdvertiserHandler.UpdateAdvertiser)
		advertisers.DELETE("/:id", opts.AdvertiserHandler.DeleteAdvertiser)
		
		// Everflow sync endpoints
		advertisers.POST("/:id/sync-to-everflow", opts.AdvertiserHandler.SyncAdvertiserToEverflow)
		advertisers.POST("/:id/sync-from-everflow", opts.AdvertiserHandler.SyncAdvertiserFromEverflow)
		advertisers.GET("/:id/compare-with-everflow", opts.AdvertiserHandler.CompareAdvertiserWithEverflow)
		
		// Advertiser's campaigns
		advertisers.GET("/:id/campaigns", opts.CampaignHandler.ListCampaignsByAdvertiser)
		
		// Advertiser's provider mappings
		advertisers.GET("/:id/provider-mappings/:providerType", opts.AdvertiserHandler.GetProviderMapping)
	}
	
	// Advertiser provider mappings
	advProviderMappings := v1.Group("/advertiser-provider-mappings")
	advProviderMappings.Use(rbacMW("AdvertiserManager", "Admin"))
	{
		advProviderMappings.POST("", opts.AdvertiserHandler.CreateProviderMapping)
		advProviderMappings.PUT("/:mappingId", opts.AdvertiserHandler.UpdateProviderMapping)
		advProviderMappings.DELETE("/:mappingId", opts.AdvertiserHandler.DeleteProviderMapping)
	}

	// --- Affiliate Routes ---
	affiliates := v1.Group("/affiliates")
	affiliates.Use(rbacMW("AffiliateManager", "Admin"))
	{
		affiliates.POST("", opts.AffiliateHandler.CreateAffiliate)
		affiliates.GET("/:id", opts.AffiliateHandler.GetAffiliate)
		affiliates.PUT("/:id", opts.AffiliateHandler.UpdateAffiliate)
		affiliates.DELETE("/:id", opts.AffiliateHandler.DeleteAffiliate)
		
		// Affiliate's provider mappings
		affiliates.GET("/:id/provider-mappings/:providerType", opts.AffiliateHandler.GetAffiliateProviderMapping)
	}
	
	// Affiliate provider mappings
	affProviderMappings := v1.Group("/affiliate-provider-mappings")
	affProviderMappings.Use(rbacMW("AffiliateManager", "Admin"))
	{
		affProviderMappings.POST("", opts.AffiliateHandler.CreateAffiliateProviderMapping)
		affProviderMappings.PUT("/:mappingId", opts.AffiliateHandler.UpdateAffiliateProviderMapping)
		affProviderMappings.DELETE("/:mappingId", opts.AffiliateHandler.DeleteAffiliateProviderMapping)
	}

	// --- Campaign Routes ---
	campaigns := v1.Group("/campaigns")
	campaigns.Use(rbacMW("AdvertiserManager", "Admin"))
	{
		campaigns.POST("", opts.CampaignHandler.CreateCampaign)
		campaigns.GET("/:id", opts.CampaignHandler.GetCampaign)
		campaigns.PUT("/:id", opts.CampaignHandler.UpdateCampaign)
		campaigns.DELETE("/:id", opts.CampaignHandler.DeleteCampaign)
	}

	return r
}