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
	campaignRepo := repository.NewPgxCampaignRepository(db)
	orgRepo := repository.NewPgxOrganizationRepository(db)

	// Initialize crypto service
	cryptoService := crypto.NewServiceFromConfig()

	// Initialize Everflow service
	everflowService, err := everflow.NewEverflowServiceFromEnv(
		advertiserRepo,
		campaignRepo,
		cryptoService,
	)
	if err != nil {
		log.Printf("Warning: Failed to initialize Everflow service: %v", err)
	}

	// Initialize services
	advertiserService := service.NewAdvertiserService(
		advertiserRepo,
		orgRepo,
		everflowService,
		cryptoService,
	)

	campaignService := service.NewCampaignService(
		campaignRepo,
		advertiserRepo,
		orgRepo,
		everflowService,
		cryptoService,
	)
	
	fmt.Println("Advertiser Service:", advertiserService)
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