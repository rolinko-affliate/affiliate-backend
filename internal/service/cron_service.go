package service

import (
	"context"
	"log"
	"time"
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
	log.Println("Starting cron service...")

	// Start daily usage calculation job
	go s.runDailyUsageCalculation()

	log.Println("Cron service started")
}

// Stop stops the cron service
func (s *CronService) Stop() {
	log.Println("Stopping cron service...")
	close(s.stopChan)
	log.Println("Cron service stopped")
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

			log.Printf("Running daily usage calculation for %s", yesterdayDate.Format("2006-01-02"))

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			err := s.usageCalculationService.CalculateDailyUsage(ctx, yesterdayDate)
			cancel()

			if err != nil {
				log.Printf("Error in daily usage calculation: %v", err)
			} else {
				log.Printf("Daily usage calculation completed successfully for %s", yesterdayDate.Format("2006-01-02"))
			}

			// Reset timer for next day (24 hours from now)
			timer.Reset(24 * time.Hour)

		case <-s.stopChan:
			log.Println("Daily usage calculation job stopped")
			return
		}
	}
}

// RunManualUsageCalculation runs usage calculation for a specific date manually
func (s *CronService) RunManualUsageCalculation(ctx context.Context, date time.Time) error {
	log.Printf("Running manual usage calculation for %s", date.Format("2006-01-02"))

	err := s.usageCalculationService.CalculateDailyUsage(ctx, date)
	if err != nil {
		log.Printf("Error in manual usage calculation: %v", err)
		return err
	}

	log.Printf("Manual usage calculation completed successfully for %s", date.Format("2006-01-02"))
	return nil
}
