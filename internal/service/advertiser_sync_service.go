package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/provider"
	"github.com/affiliate-backend/internal/repository"
)

type AdvertiserSyncService struct {
	advertiserRepo         repository.AdvertiserRepository
	providerMappingRepo    repository.AdvertiserProviderMappingRepository
	providerAdvertiserSvc  provider.ProviderAdvertiserService
}

func NewAdvertiserSyncService(
	advertiserRepo repository.AdvertiserRepository,
	providerMappingRepo repository.AdvertiserProviderMappingRepository,
	providerAdvertiserSvc provider.ProviderAdvertiserService,
) *AdvertiserSyncService {
	return &AdvertiserSyncService{
		advertiserRepo:        advertiserRepo,
		providerMappingRepo:   providerMappingRepo,
		providerAdvertiserSvc: providerAdvertiserSvc,
	}
}

func (s *AdvertiserSyncService) SyncToProvider(ctx context.Context, advertiserID int64) error {
	if s.providerAdvertiserSvc == nil {
		return fmt.Errorf("Provider advertiser service not available")
	}

	advertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, advertiserID)
	if err != nil {
		return fmt.Errorf("failed to get advertiser: %w", err)
	}

	mapping, err := s.providerMappingRepo.GetAdvertiserProviderMapping(ctx, advertiserID, "everflow")
	if err != nil {
		return s.providerAdvertiserSvc.CreateAdvertiserInProvider(ctx, advertiser)
	}

	if mapping.ProviderAdvertiserID == nil {
		return fmt.Errorf("advertiser mapping exists but has no provider advertiser ID")
	}

	_, err = s.providerAdvertiserSvc.UpdateAdvertiserInProvider(ctx, advertiserID, advertiser)
	return err
}

func (s *AdvertiserSyncService) SyncFromProvider(ctx context.Context, advertiserID int64) error {
	if s.providerAdvertiserSvc == nil {
		return fmt.Errorf("Provider advertiser service not available")
	}

	providerAdvertiser, err := s.providerAdvertiserSvc.GetAdvertiserFromProvider(ctx, advertiserID, []string{"billing", "settings"})
	if err != nil {
		return fmt.Errorf("failed to get advertiser from provider: %w", err)
	}

	// The provider service already returns a domain.Advertiser with merged data
	return s.advertiserRepo.UpdateAdvertiser(ctx, providerAdvertiser)
}

func (s *AdvertiserSyncService) CompareWithProvider(ctx context.Context, advertiserID int64) ([]domain.AdvertiserDiscrepancy, error) {
	var discrepancies []domain.AdvertiserDiscrepancy

	if s.providerAdvertiserSvc == nil {
		return discrepancies, fmt.Errorf("Provider advertiser service not available")
	}

	localAdvertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, advertiserID)
	if err != nil {
		return discrepancies, fmt.Errorf("failed to get local advertiser: %w", err)
	}

	providerAdvertiser, err := s.providerAdvertiserSvc.GetAdvertiserFromProvider(ctx, advertiserID, []string{"billing", "settings"})
	if err != nil {
		return discrepancies, fmt.Errorf("failed to get provider advertiser: %w", err)
	}

	return s.compareAdvertisers(localAdvertiser, providerAdvertiser), nil
}

func (s *AdvertiserSyncService) AsyncSyncToProvider(ctx context.Context, advertiser *domain.Advertiser) {
	go func() {
		bgCtx := context.Background()
		
		pendingStatus := "pending"
		advertiser.EverflowSyncStatus = &pendingStatus
		s.advertiserRepo.UpdateAdvertiser(bgCtx, advertiser)

		if err := s.providerAdvertiserSvc.CreateAdvertiserInProvider(bgCtx, advertiser); err != nil {
			failedStatus := "failed"
			errorMsg := err.Error()
			advertiser.EverflowSyncStatus = &failedStatus
			advertiser.EverflowSyncError = &errorMsg
			s.advertiserRepo.UpdateAdvertiser(bgCtx, advertiser)
			
			log.Printf("Error creating advertiser %d in provider: %v", advertiser.AdvertiserID, err)
		} else {
			syncedStatus := "synced"
			now := time.Now()
			advertiser.EverflowSyncStatus = &syncedStatus
			advertiser.LastEverflowSyncAt = &now
			advertiser.EverflowSyncError = nil
			s.advertiserRepo.UpdateAdvertiser(bgCtx, advertiser)
			
			log.Printf("Successfully synced advertiser %d to provider", advertiser.AdvertiserID)
		}
	}()
}

func (s *AdvertiserSyncService) AsyncSyncUpdateToProvider(ctx context.Context, advertiser *domain.Advertiser) {
	if advertiser.EverflowSyncStatus == nil || *advertiser.EverflowSyncStatus != "synced" {
		return
	}

	go func() {
		bgCtx := context.Background()
		
		pendingStatus := "pending"
		advertiser.EverflowSyncStatus = &pendingStatus
		s.advertiserRepo.UpdateAdvertiser(bgCtx, advertiser)
		
		if err := s.SyncToProvider(bgCtx, advertiser.AdvertiserID); err != nil {
			failedStatus := "failed"
			errorMsg := err.Error()
			advertiser.EverflowSyncStatus = &failedStatus
			advertiser.EverflowSyncError = &errorMsg
			s.advertiserRepo.UpdateAdvertiser(bgCtx, advertiser)
			
			log.Printf("Error updating advertiser %d in provider: %v", advertiser.AdvertiserID, err)
		} else {
			syncedStatus := "synced"
			now := time.Now()
			advertiser.EverflowSyncStatus = &syncedStatus
			advertiser.LastEverflowSyncAt = &now
			advertiser.EverflowSyncError = nil
			s.advertiserRepo.UpdateAdvertiser(bgCtx, advertiser)
			
			log.Printf("Successfully updated advertiser %d in provider", advertiser.AdvertiserID)
		}
	}()
}

func (s *AdvertiserSyncService) GetAdvertiserWithProviderData(ctx context.Context, id int64) (*domain.AdvertiserWithProviderData, error) {
	advertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get advertiser: %w", err)
	}

	result := &domain.AdvertiserWithProviderData{
		Advertiser: advertiser,
		SyncStatus: "not_synced",
	}

	if s.providerAdvertiserSvc == nil {
		result.SyncStatus = "service_unavailable"
		return result, nil
	}

	providerAdvertiser, err := s.providerAdvertiserSvc.GetAdvertiserFromProvider(ctx, id, []string{"billing", "settings"})
	if err != nil {
		if advertiser.EverflowSyncStatus == nil || *advertiser.EverflowSyncStatus == "not_synced" {
			result.SyncStatus = "not_synced"
		} else {
			result.SyncStatus = "error"
			log.Printf("Error fetching advertiser %d from provider: %v", id, err)
		}
		return result, nil
	}

	// Store the provider data in the ProviderData field
	result.ProviderData = providerAdvertiser

	discrepancies, err := s.CompareWithProvider(ctx, id)
	if err != nil {
		log.Printf("Error comparing advertiser %d with provider: %v", id, err)
		result.SyncStatus = "comparison_error"
	} else {
		result.Discrepancies = discrepancies
		if len(discrepancies) == 0 {
			result.SyncStatus = "synced"
		} else {
			result.SyncStatus = "out_of_sync"
		}
	}

	return result, nil
}

// Helper methods for comparing advertisers
func (s *AdvertiserSyncService) compareAdvertisers(local *domain.Advertiser, provider *domain.Advertiser) []domain.AdvertiserDiscrepancy {
	var discrepancies []domain.AdvertiserDiscrepancy

	// Compare basic fields
	if local.Name != provider.Name {
		discrepancies = append(discrepancies, domain.AdvertiserDiscrepancy{
			Field:         "name",
			LocalValue:    local.Name,
			ProviderValue: provider.Name,
		})
	}

	if local.Status != provider.Status {
		discrepancies = append(discrepancies, domain.AdvertiserDiscrepancy{
			Field:         "status",
			LocalValue:    local.Status,
			ProviderValue: provider.Status,
		})
	}

	// Compare optional fields
	if (local.DefaultCurrencyID == nil) != (provider.DefaultCurrencyID == nil) ||
		(local.DefaultCurrencyID != nil && provider.DefaultCurrencyID != nil && *local.DefaultCurrencyID != *provider.DefaultCurrencyID) {
		localVal := ""
		providerVal := ""
		if local.DefaultCurrencyID != nil {
			localVal = *local.DefaultCurrencyID
		}
		if provider.DefaultCurrencyID != nil {
			providerVal = *provider.DefaultCurrencyID
		}
		discrepancies = append(discrepancies, domain.AdvertiserDiscrepancy{
			Field:         "default_currency_id",
			LocalValue:    localVal,
			ProviderValue: providerVal,
		})
	}

	// Add more field comparisons as needed...

	return discrepancies
}