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
	// Add other handlers and services as needed
}

// SetupRouter sets up the API router
func SetupRouter(opts RouterOptions) *gin.Engine {
	r := gin.Default() // Starts with Logger and Recovery middleware

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

	// --- Profile Routes ---
	// Create RBAC middleware factory
	rbacMW := func(allowedRoles ...string) gin.HandlerFunc {
		return middleware.RBACMiddleware(opts.ProfileService, allowedRoles...)
	}

	v1.GET("/users/me", rbacMW("Affiliate", "AdvertiserManager", "AffiliateManager", "Admin"), 
		opts.ProfileHandler.GetMyProfile)

	// --- Organization Routes (Example - Admin only) ---
	// Add organization routes here

	// --- Advertiser Routes (Example - AdvertiserManager & Admin) ---
	// Add advertiser routes here

	// --- Campaign & Offer Routes (Example - AdvertiserManager & Admin) ---
	// Add campaign and offer routes here

	// --- Affiliate Routes (Example - Affiliate & AffiliateManager, Admin) ---
	// Add affiliate routes here

	return r
}