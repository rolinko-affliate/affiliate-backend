package everflow

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/affiliate-backend/internal/platform/crypto"
	"github.com/affiliate-backend/internal/repository"
)

// Config represents the configuration for the Everflow service
type Config struct {
	APIKey string `json:"api_key"`
}

// NewEverflowServiceFromEnv creates a new Everflow service using environment variables
func NewEverflowServiceFromEnv(
	advertiserRepo repository.AdvertiserRepository,
	campaignRepo repository.CampaignRepository,
	cryptoService crypto.Service,
) (*Service, error) {
	// Check if EVERFLOW_API_KEY is set directly
	apiKey := os.Getenv("EVERFLOW_API_KEY")
	if apiKey != "" {
		log.Println("Creating Everflow service with API key from environment variable")
		return NewService(apiKey, advertiserRepo, campaignRepo, cryptoService), nil
	}

	// Check if EVERFLOW_CONFIG is set (JSON string)
	configJSON := os.Getenv("EVERFLOW_CONFIG")
	if configJSON != "" {
		var config Config
		if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
			return nil, fmt.Errorf("failed to parse EVERFLOW_CONFIG: %w", err)
		}
		
		if config.APIKey == "" {
			return nil, fmt.Errorf("EVERFLOW_CONFIG is missing api_key")
		}
		
		log.Println("Creating Everflow service with API key from EVERFLOW_CONFIG")
		return NewService(config.APIKey, advertiserRepo, campaignRepo, cryptoService), nil
	}

	// If no configuration is found, return nil without error
	// This allows the application to run without Everflow integration
	log.Println("No Everflow configuration found, Everflow integration will be disabled")
	return nil, nil
}