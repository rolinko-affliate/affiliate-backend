package api

import (
	"github.com/affiliate-backend/internal/api/handlers"
	"github.com/affiliate-backend/internal/api/middleware"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// RouterOptions contains dependencies for the router
type RouterOptions struct {
	ProfileHandler               *handlers.ProfileHandler
	ProfileService               service.ProfileService
	OrganizationHandler          *handlers.OrganizationHandler
	AdvertiserHandler            *handlers.AdvertiserHandler
	AffiliateHandler             *handlers.AffiliateHandler
	CampaignHandler              *handlers.CampaignHandler
	TrackingLinkHandler          *handlers.TrackingLinkHandler
	AnalyticsHandler             *handlers.AnalyticsHandler
	FavoritePublisherListHandler *handlers.FavoritePublisherListHandler
	PublisherMessagingHandler    *handlers.PublisherMessagingHandler
	BillingHandler               *handlers.BillingHandler
	WebhookHandler               *handlers.WebhookHandler
}

// SetupRouter sets up the API router
func SetupRouter(opts RouterOptions) *gin.Engine {
	r := gin.Default() // Starts with Logger and Recovery middleware

	// Apply CORS middleware (will only allow CORS in development)
	r.Use(middleware.CORSMiddleware())

	// Health Check
	r.GET("/health", handlers.HealthCheck)

	// Public routes (e.g., Supabase webhook for profile creation, Stripe webhooks)
	public := r.Group("/api/v1/public")
	{
		public.POST("/webhooks/supabase/new-user", opts.ProfileHandler.HandleSupabaseNewUserWebhook)
		public.POST("/webhooks/stripe", opts.WebhookHandler.HandleStripeWebhook)
	}

	// Authenticated routes
	v1 := r.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware()) // Apply JWT auth middleware

	// Create RBAC middleware factory
	rbacMW := func(allowedRoles ...string) gin.HandlerFunc {
		return middleware.RBACMiddleware(opts.ProfileService, allowedRoles...)
	}

	// Create Profile middleware factory (for endpoints that need profile but not role restrictions)
	profileMW := func() gin.HandlerFunc {
		return middleware.ProfileMiddleware(opts.ProfileService)
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
	{
		// Admin-only routes for organization management
		organizations.POST("", rbacMW("Admin"), opts.OrganizationHandler.CreateOrganization)
		organizations.PUT("/:id", rbacMW("Admin"), opts.OrganizationHandler.UpdateOrganization)
		organizations.DELETE("/:id", rbacMW("Admin"), opts.OrganizationHandler.DeleteOrganization)

		// Read-only routes accessible by users who belong to organizations
		organizations.GET("", rbacMW("Admin", "AdvertiserManager", "AffiliateManager", "User"), opts.OrganizationHandler.ListOrganizations)
		organizations.GET("/:id", rbacMW("Admin", "AdvertiserManager", "AffiliateManager", "User"), opts.OrganizationHandler.GetOrganization)

		// Organization's resources - accessible by managers and admins
		organizations.GET("/:id/advertisers", rbacMW("Admin", "AdvertiserManager"), opts.AdvertiserHandler.ListAdvertisersByOrganization)
		organizations.GET("/:id/affiliates", rbacMW("Admin", "AffiliateManager"), opts.AffiliateHandler.ListAffiliatesByOrganization)
		organizations.GET("/:id/campaigns", rbacMW("Admin", "AdvertiserManager", "AffiliateManager"), opts.CampaignHandler.ListCampaignsByOrganization)
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

	// Affiliate Search - accessible by both advertisers and affiliate managers
	v1.POST("/affiliates/search", rbacMW("AdvertiserManager", "AffiliateManager", "Admin"), opts.AffiliateHandler.AffiliatesSearch)

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

	// --- Tracking Link Routes ---
	// Organization-level tracking link routes
	organizations.GET("/:id/tracking-links", rbacMW("Admin", "AdvertiserManager", "AffiliateManager"), opts.TrackingLinkHandler.ListTrackingLinksByOrganization)
	organizations.POST("/:id/tracking-links", rbacMW("Admin", "AdvertiserManager"), opts.TrackingLinkHandler.CreateTrackingLink)
	organizations.POST("/:id/tracking-links/generate", rbacMW("Admin", "AdvertiserManager"), opts.TrackingLinkHandler.GenerateTrackingLink)
	organizations.GET("/:id/tracking-links/:link_id", rbacMW("Admin", "AdvertiserManager", "AffiliateManager"), opts.TrackingLinkHandler.GetTrackingLink)
	organizations.PUT("/:id/tracking-links/:link_id", rbacMW("Admin", "AdvertiserManager"), opts.TrackingLinkHandler.UpdateTrackingLink)
	organizations.DELETE("/:id/tracking-links/:link_id", rbacMW("Admin", "AdvertiserManager"), opts.TrackingLinkHandler.DeleteTrackingLink)
	organizations.POST("/:id/tracking-links/:link_id/regenerate", rbacMW("Admin", "AdvertiserManager"), opts.TrackingLinkHandler.RegenerateTrackingLink)
	organizations.GET("/:id/tracking-links/:link_id/qr", rbacMW("Admin", "AdvertiserManager", "AffiliateManager"), opts.TrackingLinkHandler.GetTrackingLinkQR)

	// Campaign-specific tracking link routes
	campaigns.GET("/:id/tracking-links", rbacMW("Admin", "AdvertiserManager", "AffiliateManager"), opts.TrackingLinkHandler.ListTrackingLinksByCampaign)

	// Affiliate-specific tracking link routes
	affiliates.GET("/:id/tracking-links", rbacMW("Admin", "AffiliateManager"), opts.TrackingLinkHandler.ListTrackingLinksByAffiliate)

	// --- Analytics Routes ---
	analytics := v1.Group("/analytics")
	analytics.Use(rbacMW("AdvertiserManager", "AffiliateManager", "Admin")) // Allow all managers and admins
	{
		// Autocompletion endpoint
		analytics.GET("/autocomplete", opts.AnalyticsHandler.AutocompleteOrganizations)

		// Advertiser analytics endpoints
		analytics.GET("/advertisers/:id", opts.AnalyticsHandler.GetAdvertiserByID)
		analytics.POST("/advertisers", opts.AnalyticsHandler.CreateAdvertiser) // For future data management

		// Publisher/Affiliate analytics endpoints
		analytics.GET("/affiliates/:id", opts.AnalyticsHandler.GetPublisherByID)
		analytics.GET("/affiliates/domain/:domain", opts.AnalyticsHandler.GetPublisherByDomain)
		analytics.POST("/affiliates", opts.AnalyticsHandler.CreatePublisher) // For future data management
	}

	// --- Favorite Publisher Lists Routes ---
	favoritePublisherLists := v1.Group("/favorite-publisher-lists")
	favoritePublisherLists.Use(rbacMW("AdvertiserManager", "AffiliateManager", "Admin")) // Allow all managers and admins
	{
		// List management
		favoritePublisherLists.POST("", opts.FavoritePublisherListHandler.CreateList)
		favoritePublisherLists.GET("", opts.FavoritePublisherListHandler.GetLists)
		favoritePublisherLists.GET("/:list_id", opts.FavoritePublisherListHandler.GetListByID)
		favoritePublisherLists.PUT("/:list_id", opts.FavoritePublisherListHandler.UpdateList)
		favoritePublisherLists.DELETE("/:list_id", opts.FavoritePublisherListHandler.DeleteList)

		// Publisher management within lists
		favoritePublisherLists.POST("/:list_id/publishers", opts.FavoritePublisherListHandler.AddPublisherToList)
		favoritePublisherLists.GET("/:list_id/publishers", opts.FavoritePublisherListHandler.GetListItems)
		favoritePublisherLists.PUT("/:list_id/publishers/:domain", opts.FavoritePublisherListHandler.UpdatePublisherInList)
		favoritePublisherLists.PATCH("/:list_id/publishers/:domain/status", opts.FavoritePublisherListHandler.UpdatePublisherStatus)
		favoritePublisherLists.DELETE("/:list_id/publishers/:domain", opts.FavoritePublisherListHandler.RemovePublisherFromList)

		// Utility endpoints
		favoritePublisherLists.GET("/search", opts.FavoritePublisherListHandler.GetListsContainingPublisher)
	}

	// --- Publisher Messaging Routes ---
	publisherMessaging := v1.Group("/publisher-messaging")
	publisherMessaging.Use(rbacMW("AdvertiserManager", "AffiliateManager", "Admin")) // Allow all managers and admins
	{
		// Conversation management
		publisherMessaging.POST("/conversations", opts.PublisherMessagingHandler.CreateConversation)
		publisherMessaging.GET("/conversations", opts.PublisherMessagingHandler.GetConversations)
		publisherMessaging.GET("/conversations/:conversation_id", opts.PublisherMessagingHandler.GetConversation)
		publisherMessaging.PUT("/conversations/:conversation_id/status", opts.PublisherMessagingHandler.UpdateConversationStatus)
		publisherMessaging.DELETE("/conversations/:conversation_id", opts.PublisherMessagingHandler.DeleteConversation)

		// Message management
		publisherMessaging.POST("/conversations/:conversation_id/messages", opts.PublisherMessagingHandler.AddMessage)

		// External service integration (no RBAC required for external services)
		publisherMessaging.POST("/conversations/:conversation_id/external-messages", opts.PublisherMessagingHandler.AddExternalMessage)
	}

	// --- Billing Routes ---
	billing := v1.Group("/billing")
	billing.Use(profileMW()) // Load profile for access control validation in handlers
	{
		// Billing dashboard and account management
		billing.GET("/dashboard", opts.BillingHandler.GetBillingDashboard)
		billing.PUT("/config", opts.BillingHandler.UpdateBillingConfig)

		// Payment methods
		billing.POST("/payment-methods", opts.BillingHandler.AddPaymentMethod)
		billing.DELETE("/payment-methods/:id", opts.BillingHandler.RemovePaymentMethod)

		// Transactions and recharging
		billing.POST("/recharge", opts.BillingHandler.Recharge)
		billing.GET("/transactions", opts.BillingHandler.GetTransactionHistory)
	}

	return r
}
