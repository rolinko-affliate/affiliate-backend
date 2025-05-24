package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/everflow"
	"github.com/affiliate-backend/internal/repository"
)

type AdvertiserSyncService struct {
	advertiserRepo         repository.AdvertiserRepository
	providerMappingRepo    repository.AdvertiserProviderMappingRepository
	everflowService        *everflow.Service
}

func NewAdvertiserSyncService(
	advertiserRepo repository.AdvertiserRepository,
	providerMappingRepo repository.AdvertiserProviderMappingRepository,
	everflowService *everflow.Service,
) *AdvertiserSyncService {
	return &AdvertiserSyncService{
		advertiserRepo:      advertiserRepo,
		providerMappingRepo: providerMappingRepo,
		everflowService:     everflowService,
	}
}

func (s *AdvertiserSyncService) SyncToEverflow(ctx context.Context, advertiserID int64) error {
	if s.everflowService == nil {
		return fmt.Errorf("Everflow service not available")
	}

	advertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, advertiserID)
	if err != nil {
		return fmt.Errorf("failed to get advertiser: %w", err)
	}

	mapping, err := s.providerMappingRepo.GetAdvertiserProviderMapping(ctx, advertiserID, "everflow")
	if err != nil {
		return s.everflowService.CreateAdvertiserInEverflow(ctx, advertiser)
	}

	if mapping.ProviderAdvertiserID == nil {
		return fmt.Errorf("advertiser mapping exists but has no provider advertiser ID")
	}

	updateReq, err := s.mapAdvertiserToEverflowUpdateRequest(advertiser)
	if err != nil {
		return fmt.Errorf("failed to map advertiser to Everflow update request: %w", err)
	}

	_, err = s.everflowService.UpdateAdvertiserInEverflowByMapping(ctx, advertiserID, *updateReq)
	return err
}

func (s *AdvertiserSyncService) SyncFromEverflow(ctx context.Context, advertiserID int64) error {
	if s.everflowService == nil {
		return fmt.Errorf("Everflow service not available")
	}

	everflowAdvertiser, err := s.everflowService.GetAdvertiserFromEverflowByMapping(ctx, advertiserID, []string{"billing", "settings"})
	if err != nil {
		return fmt.Errorf("failed to get advertiser from Everflow: %w", err)
	}

	localAdvertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, advertiserID)
	if err != nil {
		return fmt.Errorf("failed to get local advertiser: %w", err)
	}

	s.mapEverflowToLocalAdvertiser(everflowAdvertiser, localAdvertiser)
	return s.advertiserRepo.UpdateAdvertiser(ctx, localAdvertiser)
}

func (s *AdvertiserSyncService) CompareWithEverflow(ctx context.Context, advertiserID int64) ([]domain.AdvertiserDiscrepancy, error) {
	var discrepancies []domain.AdvertiserDiscrepancy

	if s.everflowService == nil {
		return discrepancies, fmt.Errorf("Everflow service not available")
	}

	localAdvertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, advertiserID)
	if err != nil {
		return discrepancies, fmt.Errorf("failed to get local advertiser: %w", err)
	}

	everflowAdvertiser, err := s.everflowService.GetAdvertiserFromEverflowByMapping(ctx, advertiserID, []string{"billing", "settings"})
	if err != nil {
		return discrepancies, fmt.Errorf("failed to get Everflow advertiser: %w", err)
	}

	return s.compareAdvertisers(localAdvertiser, everflowAdvertiser), nil
}

func (s *AdvertiserSyncService) AsyncSyncToEverflow(ctx context.Context, advertiser *domain.Advertiser) {
	go func() {
		bgCtx := context.Background()
		
		pendingStatus := "pending"
		advertiser.EverflowSyncStatus = &pendingStatus
		s.advertiserRepo.UpdateAdvertiser(bgCtx, advertiser)

		if err := s.everflowService.CreateAdvertiserInEverflow(bgCtx, advertiser); err != nil {
			failedStatus := "failed"
			errorMsg := err.Error()
			advertiser.EverflowSyncStatus = &failedStatus
			advertiser.EverflowSyncError = &errorMsg
			s.advertiserRepo.UpdateAdvertiser(bgCtx, advertiser)
			
			log.Printf("Error creating advertiser %d in Everflow: %v", advertiser.AdvertiserID, err)
		} else {
			syncedStatus := "synced"
			now := time.Now()
			advertiser.EverflowSyncStatus = &syncedStatus
			advertiser.LastEverflowSyncAt = &now
			advertiser.EverflowSyncError = nil
			s.advertiserRepo.UpdateAdvertiser(bgCtx, advertiser)
			
			log.Printf("Successfully synced advertiser %d to Everflow", advertiser.AdvertiserID)
		}
	}()
}

func (s *AdvertiserSyncService) AsyncSyncUpdateToEverflow(ctx context.Context, advertiser *domain.Advertiser) {
	if advertiser.EverflowSyncStatus == nil || *advertiser.EverflowSyncStatus != "synced" {
		return
	}

	go func() {
		bgCtx := context.Background()
		
		pendingStatus := "pending"
		advertiser.EverflowSyncStatus = &pendingStatus
		s.advertiserRepo.UpdateAdvertiser(bgCtx, advertiser)
		
		if err := s.SyncToEverflow(bgCtx, advertiser.AdvertiserID); err != nil {
			failedStatus := "failed"
			errorMsg := err.Error()
			advertiser.EverflowSyncStatus = &failedStatus
			advertiser.EverflowSyncError = &errorMsg
			s.advertiserRepo.UpdateAdvertiser(bgCtx, advertiser)
			
			log.Printf("Error updating advertiser %d in Everflow: %v", advertiser.AdvertiserID, err)
		} else {
			syncedStatus := "synced"
			now := time.Now()
			advertiser.EverflowSyncStatus = &syncedStatus
			advertiser.LastEverflowSyncAt = &now
			advertiser.EverflowSyncError = nil
			s.advertiserRepo.UpdateAdvertiser(bgCtx, advertiser)
			
			log.Printf("Successfully updated advertiser %d in Everflow", advertiser.AdvertiserID)
		}
	}()
}

func (s *AdvertiserSyncService) GetAdvertiserWithEverflowData(ctx context.Context, id int64) (*domain.AdvertiserWithEverflowData, error) {
	advertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get advertiser: %w", err)
	}

	result := &domain.AdvertiserWithEverflowData{
		Advertiser: advertiser,
		SyncStatus: "not_synced",
	}

	if s.everflowService == nil {
		result.SyncStatus = "service_unavailable"
		return result, nil
	}

	everflowAdvertiser, err := s.everflowService.GetAdvertiserFromEverflowByMapping(ctx, id, []string{"billing", "settings"})
	if err != nil {
		if advertiser.EverflowSyncStatus == nil || *advertiser.EverflowSyncStatus == "not_synced" {
			result.SyncStatus = "not_synced"
		} else {
			result.SyncStatus = "error"
			log.Printf("Error fetching advertiser %d from Everflow: %v", id, err)
		}
		return result, nil
	}

	result.EverflowData = everflowAdvertiser

	discrepancies, err := s.CompareWithEverflow(ctx, id)
	if err != nil {
		log.Printf("Error comparing advertiser %d with Everflow: %v", id, err)
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

// Helper methods that would be moved from the main service
func (s *AdvertiserSyncService) mapAdvertiserToEverflowUpdateRequest(advertiser *domain.Advertiser) (*everflow.EverflowUpdateAdvertiserRequest, error) {
	// Implementation would be moved from advertiser_service.go
	// This is a placeholder - the actual implementation exists in the main service
	return nil, fmt.Errorf("not implemented - to be moved from main service")
}

func (s *AdvertiserSyncService) mapEverflowToLocalAdvertiser(everflowAdvertiser *everflow.Advertiser, localAdvertiser *domain.Advertiser) {
	// Implementation would be moved from advertiser_service.go
	// This is a placeholder - the actual implementation exists in the main service
}

func (s *AdvertiserSyncService) compareAdvertisers(local *domain.Advertiser, everflow *everflow.Advertiser) []domain.AdvertiserDiscrepancy {
	// Implementation would be moved from advertiser_service.go
	// This is a placeholder - the actual implementation exists in the main service
	return []domain.AdvertiserDiscrepancy{}
}