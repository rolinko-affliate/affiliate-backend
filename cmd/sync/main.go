package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/affiliate-backend/internal/config"
	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/everflow"
	"github.com/affiliate-backend/internal/platform/everflow/advertiser"
	"github.com/affiliate-backend/internal/platform/everflow/affiliate"
	"github.com/affiliate-backend/internal/platform/everflow/offer"
	"github.com/affiliate-backend/internal/platform/everflow/tracking"
	"github.com/affiliate-backend/internal/platform/logger"
	"github.com/affiliate-backend/internal/repository"
)

// SyncReport represents the results of a synchronization operation
type SyncReport struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  string    `json:"duration"`

	// Entity counts
	AdvertisersFound    int `json:"advertisers_found"`
	AdvertisersSynced   int `json:"advertisers_synced"`
	AdvertisersFailed   int `json:"advertisers_failed"`
	
	AffiliatesFound     int `json:"affiliates_found"`
	AffiliatesSynced    int `json:"affiliates_synced"`
	AffiliatesFailed    int `json:"affiliates_failed"`
	
	CampaignsFound      int `json:"campaigns_found"`
	CampaignsSynced     int `json:"campaigns_synced"`
	CampaignsFailed     int `json:"campaigns_failed"`

	// Error details
	Errors []SyncError `json:"errors,omitempty"`
}

// SyncError represents an error that occurred during sync
type SyncError struct {
	EntityType string `json:"entity_type"`
	EntityID   int64  `json:"entity_id"`
	EntityName string `json:"entity_name"`
	Error      string `json:"error"`
}

// SyncOptions contains configuration for the sync operation
type SyncOptions struct {
	DryRun          bool
	EntityTypes     []string // "advertisers", "affiliates", "campaigns"
	MaxEntities     int      // Maximum entities to sync per type
	IncludePending  bool     // Include entities with pending sync status
	IncludeFailed   bool     // Include entities with failed sync status
	IncludeUnsynced bool     // Include entities without provider mappings
	OutputFile      string   // File to write sync report
}

func main() {
	// Parse command line flags
	var opts SyncOptions
	flag.BoolVar(&opts.DryRun, "dry-run", false, "Perform a dry run without actually syncing")
	flag.IntVar(&opts.MaxEntities, "max-entities", 100, "Maximum entities to sync per type")
	flag.BoolVar(&opts.IncludePending, "include-pending", true, "Include entities with pending sync status")
	flag.BoolVar(&opts.IncludeFailed, "include-failed", true, "Include entities with failed sync status")
	flag.BoolVar(&opts.IncludeUnsynced, "include-unsynced", true, "Include entities without provider mappings")
	flag.StringVar(&opts.OutputFile, "output", "", "File to write sync report (JSON format)")
	
	var entityTypesFlag string
	flag.StringVar(&entityTypesFlag, "entities", "advertisers,affiliates,campaigns", "Comma-separated list of entity types to sync")
	
	flag.Parse()

	// Parse entity types
	if entityTypesFlag != "" {
		opts.EntityTypes = parseEntityTypes(entityTypesFlag)
	} else {
		opts.EntityTypes = []string{"advertisers", "affiliates", "campaigns"}
	}

	// Load configuration
	config.LoadConfig()
	cfg := config.AppConfig

	// Initialize logger
	loggerConfig := logger.Config{
		Level:     logger.LogLevel(cfg.LogLevel),
		Format:    cfg.LogFormat,
		Output:    cfg.LogOutput,
		AddSource: cfg.LogAddSource,
	}
	
	loggerInstance := logger.NewLogger(loggerConfig)

	// Set global logger
	slog.SetDefault(loggerInstance.Logger)

	loggerInstance.Info("Starting Everflow synchronization", 
		"dry_run", opts.DryRun,
		"entity_types", opts.EntityTypes,
		"max_entities", opts.MaxEntities,
		"include_pending", opts.IncludePending,
		"include_failed", opts.IncludeFailed,
		"include_unsynced", opts.IncludeUnsynced)

	// Initialize database
	repository.InitDB(&cfg)
	defer repository.CloseDB()
	db := repository.DB

	// Initialize repositories
	advertiserRepo := repository.NewPgxAdvertiserRepository(db)
	affiliateRepo := repository.NewPgxAffiliateRepository(db)
	campaignRepo := repository.NewPgxCampaignRepository(db)
	
	advertiserProviderMappingRepo := repository.NewAdvertiserProviderMappingRepository(db)
	affiliateProviderMappingRepo := repository.NewAffiliateProviderMappingRepository(db)
	campaignProviderMappingRepo := repository.NewCampaignProviderMappingRepository(db)

	// Initialize Everflow integration service
	// Create Everflow API clients
	advertiserConfig := advertiser.NewConfiguration()
	advertiserConfig.DefaultHeader["X-Eflow-API-Key"] = cfg.EverflowAPIKey
	advertiserClient := advertiser.NewAPIClient(advertiserConfig)
	
	affiliateConfig := affiliate.NewConfiguration()
	affiliateConfig.DefaultHeader["X-Eflow-API-Key"] = cfg.EverflowAPIKey
	affiliateClient := affiliate.NewAPIClient(affiliateConfig)
	
	offerConfig := offer.NewConfiguration()
	offerConfig.DefaultHeader["X-Eflow-API-Key"] = cfg.EverflowAPIKey
	offerClient := offer.NewAPIClient(offerConfig)
	
	trackingConfig := tracking.NewConfiguration()
	trackingConfig.DefaultHeader["X-Eflow-API-Key"] = cfg.EverflowAPIKey
	trackingClient := tracking.NewAPIClient(trackingConfig)
	
	everflowService := everflow.NewIntegrationService(
		advertiserClient,
		affiliateClient,
		offerClient,
		trackingClient,
		advertiserRepo,
		affiliateRepo,
		campaignRepo,
		advertiserProviderMappingRepo,
		affiliateProviderMappingRepo,
		campaignProviderMappingRepo,
	)

	// Create sync service
	syncService := &SyncService{
		logger:                            loggerInstance.Logger,
		everflowService:                   everflowService,
		advertiserRepo:                    advertiserRepo,
		affiliateRepo:                     affiliateRepo,
		campaignRepo:                      campaignRepo,
		advertiserProviderMappingRepo:     advertiserProviderMappingRepo,
		affiliateProviderMappingRepo:      affiliateProviderMappingRepo,
		campaignProviderMappingRepo:       campaignProviderMappingRepo,
	}

	// Perform synchronization
	ctx := context.Background()
	report, err := syncService.SyncEntities(ctx, opts)
	if err != nil {
		logger.Error("Synchronization failed", "error", err)
		os.Exit(1)
	}

	// Output report
	if err := outputReport(report, opts.OutputFile); err != nil {
		logger.Error("Failed to output report", "error", err)
		os.Exit(1)
	}

	logger.Info("Synchronization completed",
		"duration", report.Duration,
		"advertisers_synced", report.AdvertisersSynced,
		"affiliates_synced", report.AffiliatesSynced,
		"campaigns_synced", report.CampaignsSynced,
		"total_errors", len(report.Errors))
}

// SyncService handles the synchronization logic
type SyncService struct {
	logger                            *slog.Logger
	everflowService                   *everflow.IntegrationService
	advertiserRepo                    repository.AdvertiserRepository
	affiliateRepo                     repository.AffiliateRepository
	campaignRepo                      repository.CampaignRepository
	advertiserProviderMappingRepo     repository.AdvertiserProviderMappingRepository
	affiliateProviderMappingRepo      repository.AffiliateProviderMappingRepository
	campaignProviderMappingRepo       repository.CampaignProviderMappingRepository
}

// SyncEntities performs the main synchronization logic
func (s *SyncService) SyncEntities(ctx context.Context, opts SyncOptions) (*SyncReport, error) {
	report := &SyncReport{
		StartTime: time.Now(),
		Errors:    make([]SyncError, 0),
	}

	// Sync advertisers
	if contains(opts.EntityTypes, "advertisers") {
		s.logger.Info("Starting advertiser synchronization")
		if err := s.syncAdvertisers(ctx, opts, report); err != nil {
			s.logger.Error("Advertiser synchronization failed", "error", err)
			return nil, err
		}
	}

	// Sync affiliates
	if contains(opts.EntityTypes, "affiliates") {
		s.logger.Info("Starting affiliate synchronization")
		if err := s.syncAffiliates(ctx, opts, report); err != nil {
			s.logger.Error("Affiliate synchronization failed", "error", err)
			return nil, err
		}
	}

	// Sync campaigns
	if contains(opts.EntityTypes, "campaigns") {
		s.logger.Info("Starting campaign synchronization")
		if err := s.syncCampaigns(ctx, opts, report); err != nil {
			s.logger.Error("Campaign synchronization failed", "error", err)
			return nil, err
		}
	}

	report.EndTime = time.Now()
	report.Duration = report.EndTime.Sub(report.StartTime).String()

	return report, nil
}

// syncAdvertisers handles advertiser synchronization
func (s *SyncService) syncAdvertisers(ctx context.Context, opts SyncOptions, report *SyncReport) error {
	var advertisersToSync []*domain.Advertiser

	// Get advertisers without provider mappings
	if opts.IncludeUnsynced {
		unsynced, err := s.advertiserRepo.ListAdvertisersWithoutProviderMapping(ctx, "everflow", opts.MaxEntities, 0)
		if err != nil {
			return fmt.Errorf("failed to get unsynced advertisers: %w", err)
		}
		advertisersToSync = append(advertisersToSync, unsynced...)
		s.logger.Info("Found unsynced advertisers", "count", len(unsynced))
	}

	// Get advertisers with failed/pending sync status
	if opts.IncludeFailed || opts.IncludePending {
		failed, err := s.getAdvertisersWithSyncStatus(ctx, opts)
		if err != nil {
			return fmt.Errorf("failed to get advertisers with sync status: %w", err)
		}
		advertisersToSync = append(advertisersToSync, failed...)
		s.logger.Info("Found advertisers with sync issues", "count", len(failed))
	}

	report.AdvertisersFound = len(advertisersToSync)

	// Sync each advertiser
	for _, advertiser := range advertisersToSync {
		if report.AdvertisersSynced >= opts.MaxEntities {
			s.logger.Warn("Reached maximum advertiser sync limit", "limit", opts.MaxEntities)
			break
		}

		s.logger.Info("Syncing advertiser", "advertiser_id", advertiser.AdvertiserID, "name", advertiser.Name)

		if opts.DryRun {
			s.logger.Info("DRY RUN: Would sync advertiser", "advertiser_id", advertiser.AdvertiserID)
			report.AdvertisersSynced++
			continue
		}

		if err := s.syncSingleAdvertiser(ctx, advertiser); err != nil {
			s.logger.Error("Failed to sync advertiser", 
				"advertiser_id", advertiser.AdvertiserID, 
				"name", advertiser.Name, 
				"error", err)
			
			report.Errors = append(report.Errors, SyncError{
				EntityType: "advertiser",
				EntityID:   advertiser.AdvertiserID,
				EntityName: advertiser.Name,
				Error:      err.Error(),
			})
			report.AdvertisersFailed++
		} else {
			s.logger.Info("Successfully synced advertiser", 
				"advertiser_id", advertiser.AdvertiserID, 
				"name", advertiser.Name)
			report.AdvertisersSynced++
		}
	}

	return nil
}

// syncSingleAdvertiser syncs a single advertiser to Everflow
func (s *SyncService) syncSingleAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error {
	// Check if mapping already exists
	mapping, err := s.advertiserProviderMappingRepo.GetMappingByAdvertiserAndProvider(ctx, advertiser.AdvertiserID, "everflow")
	if err != nil && err.Error() != "no rows in result set" {
		return fmt.Errorf("failed to check existing mapping: %w", err)
	}

	if mapping == nil {
		// Create new mapping
		mapping = &domain.AdvertiserProviderMapping{
			AdvertiserID: advertiser.AdvertiserID,
			ProviderType: "everflow",
			SyncStatus:   stringPtr("pending"),
		}
		if err := s.advertiserProviderMappingRepo.CreateMapping(ctx, mapping); err != nil {
			return fmt.Errorf("failed to create advertiser mapping: %w", err)
		}
	}

	// Use the Everflow integration service to create the advertiser
	providerAdvertiser, err := s.everflowService.CreateAdvertiser(ctx, *advertiser)
	if err != nil {
		// Update sync status to failed
		s.advertiserProviderMappingRepo.UpdateSyncStatus(ctx, mapping.MappingID, "failed", stringPtr(err.Error()))
		return fmt.Errorf("failed to create advertiser in Everflow: %w", err)
	}

	// Update mapping with success
	mapping.ProviderAdvertiserID = stringPtr(fmt.Sprintf("%d", providerAdvertiser.AdvertiserID))
	mapping.SyncStatus = stringPtr("synced")
	mapping.LastSyncAt = timePtr(time.Now())
	mapping.SyncError = nil

	if err := s.advertiserProviderMappingRepo.UpdateMapping(ctx, mapping); err != nil {
		return fmt.Errorf("failed to update advertiser mapping: %w", err)
	}

	return nil
}

// getAdvertisersWithSyncStatus gets advertisers with failed or pending sync status
func (s *SyncService) getAdvertisersWithSyncStatus(ctx context.Context, opts SyncOptions) ([]*domain.Advertiser, error) {
	// This would need to be implemented in the repository
	// For now, return empty slice
	return []*domain.Advertiser{}, nil
}

// syncAffiliates handles affiliate synchronization
func (s *SyncService) syncAffiliates(ctx context.Context, opts SyncOptions, report *SyncReport) error {
	// Similar implementation to advertisers
	// For now, just log that it's not implemented
	s.logger.Warn("Affiliate synchronization not yet implemented")
	return nil
}

// syncCampaigns handles campaign synchronization
func (s *SyncService) syncCampaigns(ctx context.Context, opts SyncOptions, report *SyncReport) error {
	// Similar implementation to advertisers
	// For now, just log that it's not implemented
	s.logger.Warn("Campaign synchronization not yet implemented")
	return nil
}

// Helper functions
func parseEntityTypes(entityTypesFlag string) []string {
	// Simple comma-separated parsing
	var types []string
	current := ""
	for _, char := range entityTypesFlag {
		if char == ',' {
			if current != "" {
				types = append(types, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		types = append(types, current)
	}
	return types
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func outputReport(report *SyncReport, outputFile string) error {
	reportJSON, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	if outputFile != "" {
		if err := os.WriteFile(outputFile, reportJSON, 0644); err != nil {
			return fmt.Errorf("failed to write report to file: %w", err)
		}
		fmt.Printf("Sync report written to: %s\n", outputFile)
	} else {
		fmt.Println(string(reportJSON))
	}

	return nil
}