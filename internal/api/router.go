package api

import (
	"github.com/affiliate-backend/internal/api/handlers"
	"github.com/affiliate-backend/internal/api/middleware"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RouterOptions contains dependencies for the router
type RouterOptions struct {
	ProfileHandler                     *handlers.ProfileHandler
	ProfileService                     service.ProfileService
	OrganizationHandler                *handlers.OrganizationHandler
	OrganizationAssociationHandler     *handlers.OrganizationAssociationHandler
	AdvertiserHandler                  *handlers.AdvertiserHandler
	AffiliateHandler                   *handlers.AffiliateHandler
	CampaignHandler                    *handlers.CampaignHandler
	TrackingLinkHandler                *handlers.TrackingLinkHandler
	AnalyticsHandler                   *handlers.AnalyticsHandler
	FavoritePublisherListHandler       *handlers.FavoritePublisherListHandler
	PublisherMessagingHandler          *handlers.PublisherMessagingHandler
	BillingHandler                     *handlers.BillingHandler
	WebhookHandler                     *handlers.WebhookHandler
}

// SetupRouter sets up the API router
func SetupRouter(opts RouterOptions) *gin.Engine {
	r := gin.Default() // Starts with Logger and Recovery middleware

	// Apply CORS middleware (will only allow CORS in development)
	r.Use(middleware.CORSMiddleware())

	// Health Check
	r.GET("/health", handlers.HealthCheck)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes (e.g., Supabase webhook for profile creation, Stripe webhooks)
	public := r.Group("/api/v1/public")
	{
		public.POST("/webhooks/supabase/new-user", opts.ProfileHandler.HandleSupabaseNewUserWebhook)
		public.POST("/webhooks/stripe", opts.WebhookHandler.HandleStripeWebhook)
		// Organization creation endpoint (no authentication required)
		public.POST("/organizations", opts.OrganizationHandler.CreateOrganizationPublic)
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

	// Profile management routes - accessible to all authenticated users (JWT required)
	// TODO: Add granular RBAC after implementing more detailed role permissions
	profiles := v1.Group("/profiles")
	{
		profiles.POST("", opts.ProfileHandler.CreateProfile)
		profiles.POST("/upsert", opts.ProfileHandler.UpsertProfile)
		profiles.PUT("/:id", opts.ProfileHandler.UpdateProfile) // TODO: Add user-specific access control
		profiles.DELETE("/:id", opts.ProfileHandler.DeleteProfile) // TODO: Add appropriate RBAC restrictions
	}

	// --- Organization Routes ---
	organizations := v1.Group("/organizations")
	{
		// Basic organization operations - accessible to all authenticated users (JWT required, no RBAC)
		// These don't need ProfileMiddleware since users might not have profiles yet
		organizations.POST("", opts.OrganizationHandler.CreateOrganizationPublic) // Merged from public route
		organizations.GET("", opts.OrganizationHandler.ListOrganizations)
		organizations.GET("/:id", opts.OrganizationHandler.GetOrganization)

		// Admin-only routes for organization management - need ProfileMiddleware for RBAC
		organizations.PUT("/:id", profileMW(), rbacMW("Admin"), opts.OrganizationHandler.UpdateOrganization)
		organizations.DELETE("/:id", profileMW(), rbacMW("Admin"), opts.OrganizationHandler.DeleteOrganization)

		// Organization's resources - accessible by managers and admins - need ProfileMiddleware for RBAC
		organizations.GET("/:id/advertisers", profileMW(), rbacMW("Admin", "AdvertiserManager"), opts.AdvertiserHandler.ListAdvertisersByOrganization)
		organizations.GET("/:id/affiliates", profileMW(), rbacMW("Admin", "AffiliateManager"), opts.AffiliateHandler.ListAffiliatesByOrganization)
		organizations.GET("/:id/campaigns", profileMW(), rbacMW("Admin", "AdvertiserManager", "AffiliateManager"), opts.CampaignHandler.ListCampaignsByOrganization)
		
		// Organization associations - need ProfileMiddleware for RBAC
		organizations.GET("/:id/associations", profileMW(), rbacMW("Admin", "AdvertiserManager", "AffiliateManager"), opts.OrganizationAssociationHandler.GetAssociationsForOrganization)
		
		// Visibility endpoints - get visible resources based on associations - need ProfileMiddleware for RBAC
		organizations.GET("/:id/visible-affiliates", profileMW(), rbacMW("Admin", "AdvertiserManager"), opts.OrganizationAssociationHandler.GetVisibleAffiliatesForAdvertiser)
		organizations.GET("/:id/visible-campaigns", profileMW(), rbacMW("Admin", "AffiliateManager"), opts.OrganizationAssociationHandler.GetVisibleCampaignsForAffiliate)
	}

	// --- Advertiser Routes ---
	advertisers := v1.Group("/advertisers")
	advertisers.Use(profileMW()) // Load profile first to get user role
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
	advProviderMappings.Use(profileMW()) // Load profile first to get user role
	advProviderMappings.Use(rbacMW("AdvertiserManager", "Admin"))
	{
		advProviderMappings.POST("", opts.AdvertiserHandler.CreateProviderMapping)
		advProviderMappings.PUT("/:mappingId", opts.AdvertiserHandler.UpdateProviderMapping)
		advProviderMappings.DELETE("/:mappingId", opts.AdvertiserHandler.DeleteProviderMapping)
	}

	// --- Affiliate Routes ---
	affiliates := v1.Group("/affiliates")
	affiliates.Use(profileMW()) // Load profile first to get user role
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
	v1.POST("/affiliates/search", profileMW(), rbacMW("AdvertiserManager", "AffiliateManager", "Admin"), opts.AffiliateHandler.AffiliatesSearch)

	// Affiliate provider mappings
	affProviderMappings := v1.Group("/affiliate-provider-mappings")
	affProviderMappings.Use(profileMW()) // Load profile first to get user role
	affProviderMappings.Use(rbacMW("AffiliateManager", "Admin"))
	{
		affProviderMappings.POST("", opts.AffiliateHandler.CreateAffiliateProviderMapping)
		affProviderMappings.PUT("/:mappingId", opts.AffiliateHandler.UpdateAffiliateProviderMapping)
		affProviderMappings.DELETE("/:mappingId", opts.AffiliateHandler.DeleteAffiliateProviderMapping)
	}

	// --- Campaign Routes ---
	campaigns := v1.Group("/campaigns")
	campaigns.Use(profileMW()) // Load profile first to get user role
	campaigns.Use(rbacMW("AdvertiserManager", "Admin"))
	{
		campaigns.POST("", opts.CampaignHandler.CreateCampaign)
		campaigns.GET("/:id", opts.CampaignHandler.GetCampaign)
		campaigns.PUT("/:id", opts.CampaignHandler.UpdateCampaign)
		campaigns.DELETE("/:id", opts.CampaignHandler.DeleteCampaign)
	}

	// --- Tracking Link Routes ---
	// Organization-level tracking link routes - need ProfileMiddleware for RBAC
	organizations.GET("/:id/tracking-links", profileMW(), rbacMW("Admin", "AdvertiserManager", "AffiliateManager"), opts.TrackingLinkHandler.ListTrackingLinksByOrganization)
	organizations.POST("/:id/tracking-links", profileMW(), rbacMW("Admin", "AdvertiserManager"), opts.TrackingLinkHandler.CreateTrackingLink)
	organizations.POST("/:id/tracking-links/generate", profileMW(), rbacMW("Admin", "AdvertiserManager"), opts.TrackingLinkHandler.GenerateTrackingLink)
	organizations.GET("/:id/tracking-links/:link_id", profileMW(), rbacMW("Admin", "AdvertiserManager", "AffiliateManager"), opts.TrackingLinkHandler.GetTrackingLink)
	organizations.PUT("/:id/tracking-links/:link_id", profileMW(), rbacMW("Admin", "AdvertiserManager"), opts.TrackingLinkHandler.UpdateTrackingLink)
	organizations.DELETE("/:id/tracking-links/:link_id", profileMW(), rbacMW("Admin", "AdvertiserManager"), opts.TrackingLinkHandler.DeleteTrackingLink)
	organizations.POST("/:id/tracking-links/:link_id/regenerate", profileMW(), rbacMW("Admin", "AdvertiserManager"), opts.TrackingLinkHandler.RegenerateTrackingLink)
	organizations.GET("/:id/tracking-links/:link_id/qr", profileMW(), rbacMW("Admin", "AdvertiserManager", "AffiliateManager"), opts.TrackingLinkHandler.GetTrackingLinkQR)

	// Campaign-specific tracking link routes
	campaigns.GET("/:id/tracking-links", rbacMW("Admin", "AdvertiserManager", "AffiliateManager"), opts.TrackingLinkHandler.ListTrackingLinksByCampaign)

	// Affiliate-specific tracking link routes
	affiliates.GET("/:id/tracking-links", rbacMW("Admin", "AffiliateManager"), opts.TrackingLinkHandler.ListTrackingLinksByAffiliate)

	// --- Analytics Routes ---
	analytics := v1.Group("/analytics")
	analytics.Use(profileMW()) // Load profile first to get user role
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
	favoritePublisherLists.Use(profileMW()) // Load profile first to get user role
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
	publisherMessaging.Use(profileMW()) // Load profile first to get user role
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

	// --- Organization Association Routes ---
	orgAssociations := v1.Group("/organization-associations")
	orgAssociations.Use(profileMW()) // Load profile first to get user role
	orgAssociations.Use(rbacMW("Admin", "AdvertiserManager", "AffiliateManager"))
	{
		// Create invitations and requests
		orgAssociations.POST("/invitations", opts.OrganizationAssociationHandler.CreateInvitation)
		orgAssociations.POST("/requests", opts.OrganizationAssociationHandler.CreateRequest)
		
		// List and get associations
		orgAssociations.GET("", opts.OrganizationAssociationHandler.ListAssociations)
		orgAssociations.GET("/:id", opts.OrganizationAssociationHandler.GetAssociation)
		
		// Manage association status
		orgAssociations.POST("/:id/approve", opts.OrganizationAssociationHandler.ApproveAssociation)
		orgAssociations.POST("/:id/reject", opts.OrganizationAssociationHandler.RejectAssociation)
		orgAssociations.POST("/:id/suspend", opts.OrganizationAssociationHandler.SuspendAssociation)
		orgAssociations.POST("/:id/reactivate", opts.OrganizationAssociationHandler.ReactivateAssociation)
		
		// Update visibility settings
		orgAssociations.PUT("/:id/visibility", opts.OrganizationAssociationHandler.UpdateVisibility)
	}

	return r
}
