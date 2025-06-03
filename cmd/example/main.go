package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/affiliate-backend/internal/config"
	"github.com/affiliate-backend/internal/platform/crypto"
	"github.com/affiliate-backend/internal/platform/everflow"
	"github.com/affiliate-backend/internal/repository"
	"github.com/affiliate-backend/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database connection
	db, err := pgxpool.New(context.Background(), config.AppConfig.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	advertiserRepo := repository.NewPgxAdvertiserRepository(db)
	advertiserProviderMappingRepo := repository.NewPgxAdvertiserProviderMappingRepository(db)
	affiliateRepo := repository.NewPgxAffiliateRepository(db)
	affiliateProviderMappingRepo := repository.NewPgxAffiliateProviderMappingRepository(db)
	campaignRepo := repository.NewPgxCampaignRepository(db)
	campaignProviderMappingRepo := repository.NewPgxCampaignProviderMappingRepository(db)
	orgRepo := repository.NewPgxOrganizationRepository(db)

	// Initialize crypto service
	cryptoService := crypto.NewServiceFromConfig()

	// Initialize integration service with Everflow configuration
	everflowConfig := everflow.Config{
		BaseURL: "https://api.eflow.team",
		APIKey:  "your-api-key-here", // TODO: Load from environment
	}
	integrationService := everflow.NewIntegrationServiceWithClients(
		everflowConfig,
		advertiserRepo,
		affiliateRepo,
		campaignRepo,
		advertiserProviderMappingRepo,
		affiliateProviderMappingRepo,
		campaignProviderMappingRepo,
	)

	// Initialize services
	advertiserService := service.NewAdvertiserService(
		advertiserRepo,
		advertiserProviderMappingRepo,
		orgRepo,
		cryptoService,
		integrationService,
	)

	affiliateService := service.NewAffiliateService(
		affiliateRepo,
		affiliateProviderMappingRepo,
		orgRepo,
		integrationService,
	)

	campaignService := service.NewCampaignService(
		campaignRepo,
		campaignProviderMappingRepo,
		advertiserRepo,
		orgRepo,
		cryptoService,
		integrationService,
	)
	
	fmt.Println("Advertiser Service:", advertiserService)
	fmt.Println("Affiliate Service:", affiliateService)
	fmt.Println("Campaign Service:", campaignService)

	// Use the services...
	log.Println("Services initialized successfully")

	// Example: Set Everflow API key via environment variable
	if len(os.Args) > 1 && os.Args[1] == "set-everflow-key" {
		if len(os.Args) < 3 {
			log.Fatal("Usage: example set-everflow-key <api-key>")
		}
		os.Setenv("EVERFLOW_API_KEY", os.Args[2])
		log.Println("Everflow API key set")
	}
}