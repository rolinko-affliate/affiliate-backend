package main

// @title           Affiliate Backend API
// @version         1.0
// @description     API Server for Affiliate Backend Application
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	_ "github.com/affiliate-backend/docs" // Import for swagger docs
	"github.com/affiliate-backend/internal/api"
	"github.com/affiliate-backend/internal/api/handlers"
	"github.com/affiliate-backend/internal/config"
	"github.com/affiliate-backend/internal/platform/crypto"
	"github.com/affiliate-backend/internal/platform/everflow"
	"github.com/affiliate-backend/internal/platform/logger"
	"github.com/affiliate-backend/internal/platform/provider"
	"github.com/affiliate-backend/internal/platform/stripe"
	"github.com/affiliate-backend/internal/repository"
	"github.com/affiliate-backend/internal/service"
	"github.com/redis/go-redis/v9"
)

// getLatestMigrationVersion scans the migrations directory to find the highest version number
func getLatestMigrationVersion() (int64, error) {
	files, err := filepath.Glob("migrations/*.up.sql")
	if err != nil {
		return 0, fmt.Errorf("failed to read migrations directory: %v", err)
	}

	var latestVersion int64 = 0
	for _, file := range files {
		// Extract version number from filename like "000002_create_analytics_tables.up.sql"
		basename := filepath.Base(file)
		versionStr := strings.Split(basename, "_")[0]

		version, err := strconv.ParseInt(versionStr, 10, 64)
		if err != nil {
			continue // Skip files that don't follow the expected naming convention
		}

		if version > latestVersion {
			latestVersion = version
		}
	}

	return latestVersion, nil
}

// checkDatabaseMigrations verifies that the database schema is up to date
// If autoMigrate is true, it will attempt to run pending migrations
func checkDatabaseMigrations(cfg *config.Config, autoMigrate bool) error {

	if cfg.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}

	logger.Info("Checking database connection and migration status...")

	// Initialize the database connection
	db, err := repository.InitDBConnection(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Check if the schema_migrations table exists
	var exists bool
	err = db.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'schema_migrations')").Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if schema_migrations table exists: %v", err)
	}

	if !exists {
		logger.Warn("Migration table does not exist. Database has not been initialized with migrations.")
		if autoMigrate {
			logger.Info("Auto-migrate flag is set. Attempting to run migrations...")

			// Execute the migrate command
			cmd := exec.Command("go", "run", "./cmd/migrate/main.go", "up")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to run migrations: %v", err)
			}

			logger.Info("Migrations applied successfully")
		} else {
			return fmt.Errorf("database migrations are required. Run 'make migrate-up' or start with --auto-migrate flag")
		}
	} else {
		// Check if there are pending migrations by querying the schema_migrations table
		var version int64
		var dirty bool
		err = db.QueryRow(context.Background(),
			"SELECT MAX(version) as version, bool_or(dirty) as dirty FROM schema_migrations").Scan(&version, &dirty)

		if err != nil {
			// If the table exists but we can't query it, something is wrong
			return fmt.Errorf("failed to check migration status: %v", err)
		}

		if dirty {
			return fmt.Errorf("database schema is in a dirty state (version: %d). Manual intervention required", version)
		}

		logger.Info("Database schema version", "version", version)

		// Check if there are newer migration files available
		latestVersion, err := getLatestMigrationVersion()
		if err != nil {
			logger.Warn("Could not determine latest migration version", "error", err)
			logger.Info("Database schema appears to be up to date (unable to verify)")
		} else if version < latestVersion {
			logger.Warn("Database schema is outdated", "current", version, "latest", latestVersion)
			if autoMigrate {
				logger.Info("Auto-migrate flag is set. Attempting to run migrations...")

				// Execute the migrate command
				cmd := exec.Command("go", "run", "./cmd/migrate/main.go", "up")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					return fmt.Errorf("failed to run migrations: %v", err)
				}

				logger.Info("Migrations applied successfully")
			} else {
				return fmt.Errorf("database migrations are required. Current version: %d, Latest: %d. Run 'make migrate-up' or start with --auto-migrate flag", version, latestVersion)
			}
		} else {
			logger.Info("Database schema is up to date")
		}
	}

	if autoMigrate {
		logger.Debug("Note: For full migration functionality, install the required packages:")
		logger.Debug("  go get github.com/golang-migrate/migrate/v4")
		logger.Debug("  go get github.com/golang-migrate/migrate/v4/database/postgres")
		logger.Debug("  go get github.com/golang-migrate/migrate/v4/source/file")
	}

	return nil
}

func main() {
	// Check for help flag first
	for _, arg := range os.Args {
		if arg == "--help" || arg == "-h" {
			fmt.Println("Affiliate Backend API Server")
			fmt.Println("")
			fmt.Println("Usage:")
			fmt.Println("  affiliate-backend [flags]")
			fmt.Println("")
			fmt.Println("Flags:")
			fmt.Println("  --help, -h        Show this help message")
			fmt.Println("  --auto-migrate    Automatically run database migrations if needed")
			fmt.Println("  --mock-mode       Run in mock mode (uses mock integration service instead of real Everflow provider)")
			fmt.Println("")
			fmt.Println("Environment Variables:")
			fmt.Println("  DATABASE_URL      PostgreSQL connection string")
			fmt.Println("  PORT              Server port (default: 8080)")
			fmt.Println("  MOCK_MODE         Enable mock mode (true/false, default: false)")
			fmt.Println("  LOG_LEVEL         Logging level (DEBUG, INFO, WARN, ERROR, default: INFO)")
			fmt.Println("  LOG_FORMAT        Log format (text, json, default: text)")
			fmt.Println("                    Note: Mock mode only replaces the integration service, database is still required")
			fmt.Println("")
			return
		}
	}

	// Load Configuration
	config.LoadConfig()
	appConf := config.AppConfig

	// Initialize Logger
	loggerConfig := logger.Config{
		Level:     logger.LogLevel(appConf.LogLevel),
		Format:    appConf.LogFormat,
		Output:    appConf.LogOutput,
		AddSource: appConf.LogAddSource,
	}
	logger.InitDefault(loggerConfig)

	// Check command line flags
	autoMigrate := false
	mockMode := false
	for _, arg := range os.Args {
		if arg == "--auto-migrate" {
			autoMigrate = true
		}
		if arg == "--mock-mode" {
			mockMode = true
		}
	}

	// Override config with command line flag if provided
	if mockMode {
		appConf.MockMode = true
		logger.Info("Mock mode enabled via command line flag")
	}

	// Check database migration status
	if err := checkDatabaseMigrations(&appConf, autoMigrate); err != nil {
		logger.Fatal("Database migration check failed", "error", err)
	}

	// Initialize Database and Repositories
	repository.InitDB(&appConf)
	defer repository.CloseDB()

	// Initialize Repositories
	profileRepo := repository.NewPgxProfileRepository(repository.DB)
	organizationRepo := repository.NewPgxOrganizationRepository(repository.DB)
	organizationAssociationRepo := repository.NewPgxOrganizationAssociationRepository(repository.DB)
	advertiserAssociationInvitationRepo := repository.NewPgxAdvertiserAssociationInvitationRepository(repository.DB)
	agencyDelegationRepo := repository.NewPgxAgencyDelegationRepository(repository.DB)
	advertiserRepo := repository.NewPgxAdvertiserRepository(repository.DB)
	advertiserProviderMappingRepo := repository.NewAdvertiserProviderMappingRepository(repository.DB)
	affiliateRepo := repository.NewPgxAffiliateRepository(repository.DB)
	affiliateProviderMappingRepo := repository.NewPgxAffiliateProviderMappingRepository(repository.DB)
	campaignRepo := repository.NewPgxCampaignRepository(repository.DB)
	campaignProviderMappingRepo := repository.NewPgxCampaignProviderMappingRepository(repository.DB)
	trackingLinkRepo := repository.NewTrackingLinkRepository(repository.DB)
	trackingLinkProviderMappingRepo := repository.NewTrackingLinkProviderMappingRepository(repository.DB)
	analyticsRepo := repository.NewAnalyticsRepository(repository.DB)
	favoritePublisherListRepo := repository.NewFavoritePublisherListRepository(repository.DB)
	publisherMessagingRepo := repository.NewPublisherMessagingRepository(repository.DB)
	reportingRepo := repository.NewReportingRepository(repository.DB)

	// Initialize Billing Repositories
	billingAccountRepo := repository.NewPgxBillingAccountRepository(repository.DB)
	paymentMethodRepo := repository.NewPgxPaymentMethodRepository(repository.DB)
	transactionRepo := repository.NewPgxTransactionRepository(repository.DB)
	usageRecordRepo := repository.NewPgxUsageRecordRepository(repository.DB)
	webhookEventRepo := repository.NewPgxWebhookEventRepository(repository.DB)

	// Initialize Platform Services
	cryptoService := crypto.NewServiceFromConfig()

	// Initialize Stripe service
	stripeConfig := stripe.Config{
		SecretKey:      os.Getenv("STRIPE_SECRET_KEY"),
		PublishableKey: os.Getenv("STRIPE_PUBLISHABLE_KEY"),
		WebhookSecret:  os.Getenv("STRIPE_WEBHOOK_SECRET"),
		Environment:    os.Getenv("STRIPE_ENVIRONMENT"), // "test" or "live"
	}

	// Use test keys if not provided
	if stripeConfig.SecretKey == "" {
		stripeConfig.SecretKey = "sk_test_..." // Default test key
		logger.Warn("Using default Stripe test secret key")
	}
	if stripeConfig.WebhookSecret == "" {
		stripeConfig.WebhookSecret = "whsec_..." // Default test webhook secret
		logger.Warn("Using default Stripe webhook secret")
	}
	if stripeConfig.Environment == "" {
		stripeConfig.Environment = "test"
	}

	stripeService := stripe.NewService(stripeConfig)

	// Initialize integration service based on configuration
	var integrationService provider.IntegrationService
	if appConf.IsMockMode() {
		logger.Info("Starting in MOCK MODE - using LoggingMockIntegrationService")
		integrationService = provider.NewLoggingMockIntegrationService()
	} else {
		logger.Info("Starting in PRODUCTION MODE - using real Everflow integration")
		// Initialize integration service with Everflow configuration
		everflowConfig := everflow.Config{
			BaseURL: "https://api.eflow.team/v1",
			APIKey:  appConf.EverflowAPIKey,
		}
		integrationService = everflow.NewIntegrationServiceWithClients(
			everflowConfig,
			advertiserRepo,
			affiliateRepo,
			campaignRepo,
			advertiserProviderMappingRepo,
			affiliateProviderMappingRepo,
			campaignProviderMappingRepo,
		)
	}

	// Initialize Reporting Client and Service
	var reportingService service.ReportingService
	if appConf.MockMode {
		// In mock mode, we could create a mock reporting service
		// For now, we'll use the real service but it will fail gracefully
		reportingClient := everflow.NewReportingClient(everflow.Config{
			BaseURL: "https://api.eflow.team/v1",
			APIKey:  "mock-key",
		})
		reportingService = service.NewReportingService(reportingClient, reportingRepo, campaignRepo)
	} else {
		everflowConfig := everflow.Config{
			BaseURL: "https://api.eflow.team/v1",
			APIKey:  appConf.EverflowAPIKey,
		}
		reportingClient := everflow.NewReportingClient(everflowConfig)
		reportingService = service.NewReportingService(reportingClient, reportingRepo, campaignRepo)
	}

	// Initialize Domain Services
	profileService := service.NewProfileService(profileRepo)
	organizationService := service.NewOrganizationService(organizationRepo, advertiserRepo, affiliateRepo)
	organizationAssociationService := service.NewOrganizationAssociationService(organizationAssociationRepo, organizationRepo, profileRepo, affiliateRepo, campaignRepo)
	advertiserAssociationInvitationService := service.NewAdvertiserAssociationInvitationService(advertiserAssociationInvitationRepo, organizationAssociationRepo, organizationRepo, profileRepo, organizationAssociationService)
	agencyDelegationService := service.NewAgencyDelegationService(agencyDelegationRepo, organizationRepo, profileRepo)
	advertiserService := service.NewAdvertiserService(advertiserRepo, advertiserProviderMappingRepo, organizationRepo, cryptoService, integrationService)
	affiliateService := service.NewAffiliateService(affiliateRepo, affiliateProviderMappingRepo, organizationRepo, integrationService)
	campaignService := service.NewCampaignService(campaignRepo, campaignProviderMappingRepo, integrationService)
	trackingLinkService := service.NewTrackingLinkService(trackingLinkRepo, trackingLinkProviderMappingRepo, campaignRepo, affiliateRepo, campaignProviderMappingRepo, affiliateProviderMappingRepo, integrationService, organizationAssociationService)
	analyticsService := service.NewAnalyticsService(analyticsRepo)
	favoritePublisherListService := service.NewFavoritePublisherListService(favoritePublisherListRepo, analyticsRepo)
	publisherMessagingService := service.NewPublisherMessagingService(publisherMessagingRepo, analyticsRepo, favoritePublisherListRepo)

	// Initialize Billing Services
	billingService := service.NewBillingService(billingAccountRepo, paymentMethodRepo, transactionRepo, organizationRepo, stripeService)
	usageCalculationService := service.NewUsageCalculationService(usageRecordRepo, billingAccountRepo, transactionRepo, campaignRepo, affiliateRepo, billingService)
	cronService := service.NewCronService(usageCalculationService)

	// Initialize Handlers
	profileHandler := handlers.NewProfileHandler(profileService)
	organizationHandler := handlers.NewOrganizationHandler(organizationService, profileService)
	organizationAssociationHandler := handlers.NewOrganizationAssociationHandler(organizationAssociationService)
	advertiserAssociationInvitationHandler := handlers.NewAdvertiserAssociationInvitationHandler(advertiserAssociationInvitationService)
	agencyDelegationHandler := handlers.NewAgencyDelegationHandler(agencyDelegationService)
	advertiserHandler := handlers.NewAdvertiserHandler(advertiserService, profileService)
	affiliateHandler := handlers.NewAffiliateHandler(affiliateService, profileService, analyticsService)
	campaignHandler := handlers.NewCampaignHandler(campaignService)
	trackingLinkHandler := handlers.NewTrackingLinkHandler(trackingLinkService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	favoritePublisherListHandler := handlers.NewFavoritePublisherListHandler(favoritePublisherListService)
	publisherMessagingHandler := handlers.NewPublisherMessagingHandler(publisherMessagingService)

	// Initialize Billing Handlers
	billingHandler := handlers.NewBillingHandler(billingService, profileService)
	webhookHandler := handlers.NewWebhookHandler(stripeService, billingService, webhookEventRepo, billingAccountRepo, transactionRepo, stripeConfig.WebhookSecret)

	// Initialize Reporting Handler
	reportingHandler := handlers.NewReportingHandler(reportingService)

	// Initialize Redis client for caching (optional - can be nil for now)
	var redisClient *redis.Client
	if appConf.RedisURL != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr: appConf.RedisURL,
		})
		// Test connection
		if err := redisClient.Ping(context.Background()).Err(); err != nil {
			logger.Warn("Redis connection failed, continuing without cache", "error", err)
			redisClient = nil
		}
	}

	// Initialize Dashboard Service and Handler
	dashboardCacheRepo := repository.NewDashboardCacheRepository(redisClient)
	everflowRepo := repository.NewEverflowRepository(appConf.EverflowAPIURL, appConf.EverflowAPIKey, redisClient, logger.GetDefault())
	dashboardService := service.NewDashboardService(dashboardCacheRepo, everflowRepo, reportingService, profileService, organizationService, logger.GetDefault())
	dashboardHandler := handlers.NewDashboardHandler(dashboardService, logger.GetDefault())

	// Setup Router
	router := api.SetupRouter(api.RouterOptions{
		ProfileHandler:                         profileHandler,
		ProfileService:                         profileService,
		OrganizationHandler:                    organizationHandler,
		OrganizationAssociationHandler:         organizationAssociationHandler,
		AdvertiserAssociationInvitationHandler: advertiserAssociationInvitationHandler,
		AgencyDelegationHandler:                agencyDelegationHandler,
		AdvertiserHandler:                      advertiserHandler,
		AffiliateHandler:                       affiliateHandler,
		CampaignHandler:                        campaignHandler,
		TrackingLinkHandler:                    trackingLinkHandler,
		AnalyticsHandler:                       analyticsHandler,
		FavoritePublisherListHandler:           favoritePublisherListHandler,
		PublisherMessagingHandler:              publisherMessagingHandler,
		BillingHandler:                         billingHandler,
		WebhookHandler:                         webhookHandler,
		ReportingHandler:                       reportingHandler,
		DashboardHandler:                       dashboardHandler,
	})

	// Start Server
	srv := &http.Server{
		Addr:    ":" + appConf.Port,
		Handler: router,
	}

	// Start cron service
	cronService.Start()
	defer cronService.Stop()

	// Start the server in a goroutine
	go func() {
		logger.Info("Server starting", "port", appConf.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", "error", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exiting")
}
