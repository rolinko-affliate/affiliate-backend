package service

import (
	"context"
	"time"

	"github.com/affiliate-backend/internal/platform/logger"
)

// CronService handles scheduled tasks
type CronService struct {
	usageCalculationService *UsageCalculationService
	stopChan                chan bool
}

// NewCronService creates a new cron service
func NewCronService(usageCalculationService *UsageCalculationService) *CronService {
	return &CronService{
		usageCalculationService: usageCalculationService,
		stopChan:                make(chan bool),
	}
}

// Start starts the cron service
func (s *CronService) Start() {
	logger.Info("Starting cron service")

	// Start daily usage calculation job
	go s.runDailyUsageCalculation()

	logger.Info("Cron service started")
}

// Stop stops the cron service
func (s *CronService) Stop() {
	logger.Info("Stopping cron service")
	close(s.stopChan)
	logger.Info("Cron service stopped")
}

// runDailyUsageCalculation runs the daily usage calculation job
func (s *CronService) runDailyUsageCalculation() {
	// Calculate the next midnight
	now := time.Now()
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

	// Wait until midnight for the first run
	timer := time.NewTimer(time.Until(nextMidnight))
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			// Run usage calculation for yesterday
			yesterday := time.Now().AddDate(0, 0, -1)
			yesterdayDate := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())

			logger.Info("Running daily usage calculation", "date", yesterdayDate.Format("2006-01-02"))

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			err := s.usageCalculationService.CalculateDailyUsage(ctx, yesterdayDate)
			cancel()

			if err != nil {
				logger.Error("Error in daily usage calculation", "error", err)
			} else {
				logger.Info("Daily usage calculation completed successfully", "date", yesterdayDate.Format("2006-01-02"))
			}

			// Reset timer for next day (24 hours from now)
			timer.Reset(24 * time.Hour)

		case <-s.stopChan:
			logger.Info("Daily usage calculation job stopped")
			return
		}
	}
}

// RunManualUsageCalculation runs usage calculation for a specific date manually
func (s *CronService) RunManualUsageCalculation(ctx context.Context, date time.Time) error {
	logger.Info("Running manual usage calculation", "date", date.Format("2006-01-02"))

	err := s.usageCalculationService.CalculateDailyUsage(ctx, date)
	if err != nil {
		logger.Error("Error in manual usage calculation", "error", err)
		return err
	}

	logger.Info("Manual usage calculation completed successfully", "date", date.Format("2006-01-02"))
	return nil
}
