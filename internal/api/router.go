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
	ProfileHandler                         *handlers.ProfileHandler
	ProfileService                         service.ProfileService
	OrganizationHandler                    *handlers.OrganizationHandler
	OrganizationAssociationHandler         *handlers.OrganizationAssociationHandler
	AgencyDelegationHandler                *handlers.AgencyDelegationHandler
	AdvertiserAssociationInvitationHandler *handlers.AdvertiserAssociationInvitationHandler
	AdvertiserHandler                      *handlers.AdvertiserHandler
	AffiliateHandler                       *handlers.AffiliateHandler
	CampaignHandler                        *handlers.CampaignHandler
	TrackingLinkHandler                    *handlers.TrackingLinkHandler
	AnalyticsHandler                       *handlers.AnalyticsHandler
	FavoritePublisherListHandler           *handlers.FavoritePublisherListHandler
	PublisherMessagingHandler              *handlers.PublisherMessagingHandler
	BillingHandler                         *handlers.BillingHandler
	WebhookHandler                         *handlers.WebhookHandler
	ReportingHandler                       *handlers.ReportingHandler
	DashboardHandler                       *handlers.DashboardHandler
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
		// Public invitation endpoint (no authentication required for viewing invitations)
		public.GET("/invitations/:token", opts.AdvertiserAssociationInvitationHandler.GetInvitationByToken)
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
		// TODO: Remove open access to CreateProfile - should have proper access control in production
		profiles.POST("", opts.ProfileHandler.CreateProfile) // Temporarily without access control
		profiles.POST("/upsert", opts.ProfileHandler.UpsertProfile)
		profiles.PUT("/:id", opts.ProfileHandler.UpdateProfile)    // TODO: Add user-specific access control
		profiles.DELETE("/:id", opts.ProfileHandler.DeleteProfile) // TODO: Add appropriate RBAC restrictions
	}

	// --- Organization Routes ---
	organizations := v1.Group("/organizations")
	{
		// Basic organization operations - accessible to all authenticated users (JWT required, no RBAC)
		// POST doesn't need ProfileMiddleware since users might not have profiles yet
		organizations.POST("", opts.OrganizationHandler.CreateOrganizationPublic) // Merged from public route

		// GET operations need ProfileMiddleware for access control checks
		organizations.GET("", profileMW(), opts.OrganizationHandler.ListOrganizations)
		organizations.GET("/:id", profileMW(), opts.OrganizationHandler.GetOrganization)

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
		advertisers.POST("/sync-all-to-everflow", rbacMW("Admin"), opts.AdvertiserHandler.SyncAllAdvertisersToEverflow)

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

		// Campaign provider mappings
		campaigns.GET("/:id/provider-mappings/:providerType", opts.CampaignHandler.GetProviderMapping)
	}

	// --- Legacy Tracking Link Routes (QR code only) ---
	// Keep only the QR code endpoint for backward compatibility
	organizations.GET("/:id/tracking-links/:link_id/qr", profileMW(), rbacMW("Admin", "AdvertiserManager", "AffiliateManager"), opts.TrackingLinkHandler.GetTrackingLinkQR)

	// --- Analytics Routes ---
	analytics := v1.Group("/analytics")
	analytics.Use(profileMW())                                              // Load profile first to get user role
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
	favoritePublisherLists.Use(profileMW())                                              // Load profile first to get user role
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
	publisherMessaging.Use(profileMW())                                              // Load profile first to get user role
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

	// --- Reporting Routes ---
	reports := v1.Group("/reports")
	reports.Use(profileMW()) // Load profile first to get user role
	reports.Use(rbacMW("AdvertiserManager", "AffiliateManager", "Admin")) // Allow all managers and admins
	{
		// Performance reporting
		reports.GET("/performance/summary", opts.ReportingHandler.GetPerformanceSummary)
		reports.GET("/performance/timeseries", opts.ReportingHandler.GetPerformanceTimeSeries)
		reports.GET("/performance/daily", opts.ReportingHandler.GetDailyPerformanceReport)

		// Event reporting
		reports.GET("/conversions", opts.ReportingHandler.GetConversionsReport)
		reports.GET("/clicks", opts.ReportingHandler.GetClicksReport)
	}

	// Campaigns list endpoint (for filters) - moved from campaigns group to be accessible by reporting
	v1.GET("/campaigns", profileMW(), rbacMW("AdvertiserManager", "AffiliateManager", "Admin"), opts.ReportingHandler.GetCampaignsList)

	// --- Organization Association Routes ---
	orgAssociations := v1.Group("/organization-associations")
	orgAssociations.Use(profileMW()) // Load profile first to get user role
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

	// --- Advertiser Association Invitation Routes ---
	advInvitations := v1.Group("/advertiser-association-invitations")
	advInvitations.Use(profileMW()) // Load profile first to get user role
	{
		// Invitation management - primarily for advertisers
		advInvitations.POST("", rbacMW("AdvertiserManager", "Admin"), opts.AdvertiserAssociationInvitationHandler.CreateInvitation)
		advInvitations.GET("", rbacMW("AdvertiserManager", "AffiliateManager", "Admin"), opts.AdvertiserAssociationInvitationHandler.ListInvitations)
		advInvitations.GET("/:id", rbacMW("AdvertiserManager", "AffiliateManager", "Admin"), opts.AdvertiserAssociationInvitationHandler.GetInvitation)
		advInvitations.PUT("/:id", rbacMW("AdvertiserManager", "Admin"), opts.AdvertiserAssociationInvitationHandler.UpdateInvitation)
		advInvitations.DELETE("/:id", rbacMW("AdvertiserManager", "Admin"), opts.AdvertiserAssociationInvitationHandler.DeleteInvitation)

		// Invitation usage - for affiliates to use invitations
		advInvitations.POST("/use", rbacMW("AffiliateManager", "Admin"), opts.AdvertiserAssociationInvitationHandler.UseInvitation)

		// Invitation analytics and management
		advInvitations.GET("/:id/usage-history", rbacMW("AdvertiserManager", "Admin"), opts.AdvertiserAssociationInvitationHandler.GetInvitationUsageHistory)
		advInvitations.GET("/:id/link", rbacMW("AdvertiserManager", "Admin"), opts.AdvertiserAssociationInvitationHandler.GenerateInvitationLink)
	}

	// --- Agency Delegation Routes ---
	agencyDelegations := v1.Group("/agency-delegations")
	agencyDelegations.Use(profileMW()) // Load profile first to get user role
	{
		// Delegation management - Platform owners can create delegations between any organizations
		agencyDelegations.POST("", rbacMW("AdvertiserManager", "PlatformOwner", "Admin"), opts.AgencyDelegationHandler.CreateDelegation)
		agencyDelegations.GET("", rbacMW("AdvertiserManager", "AgencyManager", "PlatformOwner", "Admin"), opts.AgencyDelegationHandler.ListDelegations)
		agencyDelegations.GET("/:id", rbacMW("AdvertiserManager", "AgencyManager", "PlatformOwner", "Admin"), opts.AgencyDelegationHandler.GetDelegation)

		// Delegation status management
		agencyDelegations.POST("/:id/accept", rbacMW("AgencyManager", "PlatformOwner", "Admin"), opts.AgencyDelegationHandler.AcceptDelegation)
		agencyDelegations.POST("/:id/reject", rbacMW("AdvertiserManager", "AgencyManager", "PlatformOwner", "Admin"), opts.AgencyDelegationHandler.RejectDelegation)
		agencyDelegations.POST("/:id/suspend", rbacMW("AdvertiserManager", "AgencyManager", "PlatformOwner", "Admin"), opts.AgencyDelegationHandler.SuspendDelegation)
		agencyDelegations.POST("/:id/reactivate", rbacMW("AdvertiserManager", "AgencyManager", "PlatformOwner", "Admin"), opts.AgencyDelegationHandler.ReactivateDelegation)
		agencyDelegations.POST("/:id/revoke", rbacMW("AdvertiserManager", "AgencyManager", "PlatformOwner", "Admin"), opts.AgencyDelegationHandler.RevokeDelegation)

		// Delegation configuration - Platform owners can manage all delegations
		agencyDelegations.PUT("/:id/permissions", rbacMW("AdvertiserManager", "PlatformOwner", "Admin"), opts.AgencyDelegationHandler.UpdatePermissions)
		agencyDelegations.PUT("/:id/expiration", rbacMW("AdvertiserManager", "PlatformOwner", "Admin"), opts.AgencyDelegationHandler.UpdateExpiration)

		// Permission checking and utility endpoints
		agencyDelegations.POST("/check-permissions", rbacMW("AdvertiserManager", "AgencyManager", "PlatformOwner", "Admin"), opts.AgencyDelegationHandler.CheckPermissions)
		agencyDelegations.GET("/permissions", opts.AgencyDelegationHandler.GetAvailablePermissions)

		// Organization-specific delegation endpoints
		agencyDelegations.GET("/agency/:agency_org_id", rbacMW("AgencyManager", "PlatformOwner", "Admin"), opts.AgencyDelegationHandler.GetAgencyDelegations)
		agencyDelegations.GET("/advertiser/:advertiser_org_id", rbacMW("AdvertiserManager", "PlatformOwner", "Admin"), opts.AgencyDelegationHandler.GetAdvertiserDelegations)
	}

	// --- Clean Tracking Link Routes ---
	trackingLinks := v1.Group("/tracking-links")
	trackingLinks.Use(profileMW()) // Load profile first to get user role
	trackingLinks.Use(rbacMW("Admin", "AdvertiserManager", "AffiliateManager"))
	{
		trackingLinks.POST("", opts.TrackingLinkHandler.CreateTrackingLinkClean)
		trackingLinks.GET("", opts.TrackingLinkHandler.ListTrackingLinksClean)
		trackingLinks.GET("/:id", opts.TrackingLinkHandler.GetTrackingLinkClean)
		trackingLinks.PUT("/:id", opts.TrackingLinkHandler.UpdateTrackingLinkClean)
		trackingLinks.DELETE("/:id", opts.TrackingLinkHandler.DeleteTrackingLinkClean)
	}

	// --- Dashboard Routes ---
	dashboard := v1.Group("/dashboard")
	dashboard.Use(profileMW()) // Load profile first to get user role
	dashboard.Use(rbacMW("AdvertiserManager", "AffiliateManager", "AgencyManager", "PlatformOwner", "Admin"))
	{
		// Main dashboard endpoint
		dashboard.GET("", opts.DashboardHandler.GetDashboard)

		// Offers endpoint
		dashboard.GET("/offers", opts.DashboardHandler.GetOffers)

		// Campaign detail endpoint
		dashboard.GET("/campaigns/:campaignId", opts.DashboardHandler.GetCampaignDetail)

		// Activity endpoints
		dashboard.GET("/activity", opts.DashboardHandler.GetRecentActivity)
		dashboard.POST("/activity", opts.DashboardHandler.TrackActivity)

		// System health endpoint (Platform Owner only)
		dashboard.GET("/system/health", rbacMW("PlatformOwner", "Admin"), opts.DashboardHandler.GetSystemHealth)

		// Cache management
		dashboard.POST("/cache/invalidate", opts.DashboardHandler.InvalidateCache)

		// Dashboard health check
		dashboard.GET("/health", opts.DashboardHandler.DashboardHealthCheck)
	}

	return r
}
