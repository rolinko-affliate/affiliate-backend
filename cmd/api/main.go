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
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/affiliate-backend/internal/api"
	"github.com/affiliate-backend/internal/api/handlers"
	"github.com/affiliate-backend/internal/config"
	"github.com/affiliate-backend/internal/platform/crypto"
	"github.com/affiliate-backend/internal/platform/everflow"
	"github.com/affiliate-backend/internal/platform/provider"
	"github.com/affiliate-backend/internal/repository"
	"github.com/affiliate-backend/internal/service"
)

// checkDatabaseMigrations verifies that the database schema is up to date
// If autoMigrate is true, it will attempt to run pending migrations
// Skips database checks when in mock mode
func checkDatabaseMigrations(cfg *config.Config, autoMigrate bool) error {
	
	if cfg.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}

	log.Println("Checking database connection and migration status...")
	
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
		log.Println("Migration table does not exist. Database has not been initialized with migrations.")
		if autoMigrate {
			log.Println("Auto-migrate flag is set. Attempting to run migrations...")
			
			// Execute the migrate command
			cmd := exec.Command("go", "run", "./cmd/migrate/main.go", "up")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to run migrations: %v", err)
			}
			
			log.Println("Migrations applied successfully")
		} else {
			return fmt.Errorf("database migrations are required. Run 'make migrate-up' or start with --auto-migrate flag")
		}
	} else {
		// Check if there are pending migrations by querying the schema_migrations table
		var version int64
		var dirty bool
		err = db.QueryRow(context.Background(), 
			"SELECT version, dirty FROM schema_migrations LIMIT 1").Scan(&version, &dirty)
		
		if err != nil {
			// If the table exists but we can't query it, something is wrong
			return fmt.Errorf("failed to check migration status: %v", err)
		}
		
		if dirty {
			return fmt.Errorf("database schema is in a dirty state (version: %d). Manual intervention required", version)
		}
		
		log.Printf("Database schema version: %d\n", version)
		log.Println("Database schema appears to be up to date")
	}
	
	if autoMigrate {
		log.Println("Note: For full migration functionality, install the required packages:")
		log.Println("  go get github.com/golang-migrate/migrate/v4")
		log.Println("  go get github.com/golang-migrate/migrate/v4/database/postgres")
		log.Println("  go get github.com/golang-migrate/migrate/v4/source/file")
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
			fmt.Println("  --mock-mode       Run in mock mode (uses mock integration service instead of real provider)")
			fmt.Println("")
			fmt.Println("Environment Variables:")
			fmt.Println("  DATABASE_URL      PostgreSQL connection string")
			fmt.Println("  PORT              Server port (default: 8080)")
			fmt.Println("  MOCK_MODE         Enable mock mode (true/false, default: false)")
			fmt.Println("")
			return
		}
	}

	// Load Configuration
	config.LoadConfig()
	appConf := config.AppConfig

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
		log.Println("🔧 Mock mode enabled via command line flag")
	}

	// Check database migration status
	if err := checkDatabaseMigrations(&appConf, autoMigrate); err != nil {
		log.Fatalf("Database migration check failed: %v", err)
	}

	// Initialize Database and Repositories
	repository.InitDB(&appConf)
	defer repository.CloseDB()

	// Initialize Repositories
	profileRepo := repository.NewPgxProfileRepository(repository.DB)
	organizationRepo := repository.NewPgxOrganizationRepository(repository.DB)
	advertiserRepo := repository.NewPgxAdvertiserRepository(repository.DB)
	advertiserProviderMappingRepo := repository.NewPgxAdvertiserProviderMappingRepository(repository.DB)
	affiliateRepo := repository.NewPgxAffiliateRepository(repository.DB)
	affiliateProviderMappingRepo := repository.NewPgxAffiliateProviderMappingRepository(repository.DB)
	campaignRepo := repository.NewPgxCampaignRepository(repository.DB)
	campaignProviderMappingRepo := repository.NewPgxCampaignProviderMappingRepository(repository.DB)

	// Initialize Platform Services
	cryptoService := crypto.NewServiceFromConfig()
	
	// Initialize integration service based on configuration
	var integrationService provider.IntegrationService
	if appConf.IsMockMode() {
		log.Println("🔧 Starting in MOCK MODE - using LoggingMockIntegrationService")
		integrationService = provider.NewLoggingMockIntegrationService()
	} else {
		log.Println("🔧 Starting in PRODUCTION MODE - using real Everflow integration")
		// Initialize integration service with Everflow configuration
		everflowConfig := everflow.Config{
			BaseURL: "https://api.eflow.team",
			APIKey:  "your-api-key-here", // TODO: Load from environment
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
	
	// Initialize Domain Services
	profileService := service.NewProfileService(profileRepo)
	organizationService := service.NewOrganizationService(organizationRepo)
	advertiserService := service.NewAdvertiserService(advertiserRepo, advertiserProviderMappingRepo, organizationRepo, cryptoService, integrationService)
	affiliateService := service.NewAffiliateService(affiliateRepo, affiliateProviderMappingRepo, organizationRepo, integrationService)
	campaignService := service.NewCampaignService(campaignRepo, campaignProviderMappingRepo, advertiserRepo, organizationRepo, cryptoService, integrationService)

	// Initialize Handlers
	profileHandler := handlers.NewProfileHandler(profileService)
	organizationHandler := handlers.NewOrganizationHandler(organizationService, profileService)
	advertiserHandler := handlers.NewAdvertiserHandler(advertiserService, profileService)
	affiliateHandler := handlers.NewAffiliateHandler(affiliateService, profileService)
	campaignHandler := handlers.NewCampaignHandler(campaignService)

	// Setup Router
	router := api.SetupRouter(api.RouterOptions{
		ProfileHandler:      profileHandler,
		ProfileService:      profileService,
		OrganizationHandler: organizationHandler,
		AdvertiserHandler:   advertiserHandler,
		AffiliateHandler:    affiliateHandler,
		CampaignHandler:     campaignHandler,
	})

	// Start Server
	srv := &http.Server{
		Addr:    ":" + appConf.Port,
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server starting on port %s\n", appConf.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s\n", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}