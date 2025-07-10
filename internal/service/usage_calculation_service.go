package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/repository"
	"github.com/shopspring/decimal"
)

// UsageCalculationService handles daily usage calculation and billing
type UsageCalculationService struct {
	usageRecordRepo    repository.UsageRecordRepository
	billingAccountRepo repository.BillingAccountRepository
	transactionRepo    repository.TransactionRepository
	campaignRepo       repository.CampaignRepository
	affiliateRepo      repository.AffiliateRepository
	billingService     *BillingService
}

// NewUsageCalculationService creates a new usage calculation service
func NewUsageCalculationService(
	usageRecordRepo repository.UsageRecordRepository,
	billingAccountRepo repository.BillingAccountRepository,
	transactionRepo repository.TransactionRepository,
	campaignRepo repository.CampaignRepository,
	affiliateRepo repository.AffiliateRepository,
	billingService *BillingService,
) *UsageCalculationService {
	return &UsageCalculationService{
		usageRecordRepo:    usageRecordRepo,
		billingAccountRepo: billingAccountRepo,
		transactionRepo:    transactionRepo,
		campaignRepo:       campaignRepo,
		affiliateRepo:      affiliateRepo,
		billingService:     billingService,
	}
}

// CalculateDailyUsage calculates usage for all organizations for a specific date
func (s *UsageCalculationService) CalculateDailyUsage(ctx context.Context, date time.Time) error {
	log.Printf("Starting daily usage calculation for date: %s", date.Format("2006-01-02"))

	// Get all active billing accounts
	billingAccounts, err := s.billingAccountRepo.List(ctx, 1000, 0) // TODO: Implement pagination
	if err != nil {
		return fmt.Errorf("failed to get billing accounts: %w", err)
	}

	for _, account := range billingAccounts {
		if account.Status != domain.BillingAccountStatusActive {
			continue
		}

		err := s.calculateUsageForOrganization(ctx, account.OrganizationID, date)
		if err != nil {
			log.Printf("Error calculating usage for organization %d: %v", account.OrganizationID, err)
			// Continue with other organizations
		}
	}

	log.Printf("Completed daily usage calculation for date: %s", date.Format("2006-01-02"))
	return nil
}

// calculateUsageForOrganization calculates usage for a specific organization and date
func (s *UsageCalculationService) calculateUsageForOrganization(ctx context.Context, organizationID int64, date time.Time) error {
	log.Printf("Calculating usage for organization %d on %s", organizationID, date.Format("2006-01-02"))

	// Check if usage record already exists
	existingRecord, err := s.usageRecordRepo.GetByOrganizationAndDate(ctx, organizationID, date)
	if err == nil && existingRecord.Status != domain.UsageRecordStatusPending {
		log.Printf("Usage already calculated for organization %d on %s", organizationID, date.Format("2006-01-02"))
		return nil
	}

	// Get billing account
	billingAccount, err := s.billingAccountRepo.GetByOrganizationID(ctx, organizationID)
	if err != nil {
		return fmt.Errorf("failed to get billing account: %w", err)
	}

	// Calculate usage metrics
	usageMetrics, err := s.calculateUsageMetrics(ctx, organizationID, date)
	if err != nil {
		return fmt.Errorf("failed to calculate usage metrics: %w", err)
	}

	// Calculate financial metrics
	financialMetrics, err := s.calculateFinancialMetrics(ctx, organizationID, date, usageMetrics)
	if err != nil {
		return fmt.Errorf("failed to calculate financial metrics: %w", err)
	}

	// Create or update usage record
	usageRecord := &domain.UsageRecord{
		OrganizationID:     organizationID,
		BillingAccountID:   billingAccount.BillingAccountID,
		UsageDate:          date,
		Clicks:             usageMetrics.Clicks,
		Conversions:        usageMetrics.Conversions,
		Impressions:        usageMetrics.Impressions,
		AdvertiserSpend:    financialMetrics.AdvertiserSpend,
		AffiliatePayout:    financialMetrics.AffiliatePayout,
		PlatformRevenue:    financialMetrics.PlatformRevenue,
		Currency:           billingAccount.Currency,
		Status:             domain.UsageRecordStatusCalculated,
		CampaignBreakdown:  financialMetrics.CampaignBreakdown,
		AffiliateBreakdown: financialMetrics.AffiliateBreakdown,
		Metadata:           make(map[string]interface{}),
	}

	if existingRecord != nil {
		usageRecord.UsageRecordID = existingRecord.UsageRecordID
		err = s.usageRecordRepo.Update(ctx, usageRecord)
	} else {
		err = s.usageRecordRepo.Create(ctx, usageRecord)
	}

	if err != nil {
		return fmt.Errorf("failed to save usage record: %w", err)
	}

	// Process billing based on billing mode
	err = s.processBilling(ctx, billingAccount, usageRecord)
	if err != nil {
		return fmt.Errorf("failed to process billing: %w", err)
	}

	log.Printf("Successfully calculated usage for organization %d: spend=%s, payout=%s, revenue=%s",
		organizationID,
		usageRecord.AdvertiserSpend.String(),
		usageRecord.AffiliatePayout.String(),
		usageRecord.PlatformRevenue.String())

	return nil
}

// UsageMetrics represents calculated usage metrics
type UsageMetrics struct {
	Clicks      int
	Conversions int
	Impressions int
}

// FinancialMetrics represents calculated financial metrics
type FinancialMetrics struct {
	AdvertiserSpend    decimal.Decimal
	AffiliatePayout    decimal.Decimal
	PlatformRevenue    decimal.Decimal
	CampaignBreakdown  map[string]interface{}
	AffiliateBreakdown map[string]interface{}
}

// calculateUsageMetrics calculates usage metrics for an organization and date
func (s *UsageCalculationService) calculateUsageMetrics(ctx context.Context, organizationID int64, date time.Time) (*UsageMetrics, error) {
	// TODO: Implement actual usage calculation from tracking data
	// This would typically query tracking_links, campaigns, and analytics data
	// For now, returning mock data

	metrics := &UsageMetrics{
		Clicks:      100,  // Mock data
		Conversions: 10,   // Mock data
		Impressions: 1000, // Mock data
	}

	return metrics, nil
}

// calculateFinancialMetrics calculates financial metrics based on usage
func (s *UsageCalculationService) calculateFinancialMetrics(ctx context.Context, organizationID int64, date time.Time, usage *UsageMetrics) (*FinancialMetrics, error) {
	// Get campaigns for the organization
	campaigns, err := s.campaignRepo.ListCampaignsByOrganization(ctx, organizationID, 100, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaigns: %w", err)
	}

	var totalAdvertiserSpend decimal.Decimal
	var totalAffiliatePayout decimal.Decimal
	var totalPlatformRevenue decimal.Decimal

	campaignBreakdown := make(map[string]interface{})
	affiliateBreakdown := make(map[string]interface{})

	// Calculate spend per campaign
	for _, campaign := range campaigns {
		// Calculate based on billing model and payout structure
		var campaignSpend decimal.Decimal
		var campaignPayout decimal.Decimal

		if campaign.BillingModel != nil && *campaign.BillingModel == "click" {
			// Charge per click
			if campaign.RevenueAmount != nil {
				campaignSpend = decimal.NewFromFloat(*campaign.RevenueAmount).Mul(decimal.NewFromInt(int64(usage.Clicks)))
			}
			if campaign.PayoutAmount != nil {
				campaignPayout = decimal.NewFromFloat(*campaign.PayoutAmount).Mul(decimal.NewFromInt(int64(usage.Clicks)))
			}
		} else if campaign.BillingModel != nil && *campaign.BillingModel == "conversion" {
			// Charge per conversion
			if campaign.RevenueAmount != nil {
				campaignSpend = decimal.NewFromFloat(*campaign.RevenueAmount).Mul(decimal.NewFromInt(int64(usage.Conversions)))
			}
			if campaign.PayoutAmount != nil {
				campaignPayout = decimal.NewFromFloat(*campaign.PayoutAmount).Mul(decimal.NewFromInt(int64(usage.Conversions)))
			}
		}

		totalAdvertiserSpend = totalAdvertiserSpend.Add(campaignSpend)
		totalAffiliatePayout = totalAffiliatePayout.Add(campaignPayout)

		// Store campaign breakdown
		campaignBreakdown[fmt.Sprintf("campaign_%d", campaign.CampaignID)] = map[string]interface{}{
			"name":             campaign.Name,
			"spend":            campaignSpend,
			"payout":           campaignPayout,
			"clicks":           usage.Clicks,   // Simplified - should be per campaign
			"conversions":      usage.Conversions, // Simplified - should be per campaign
			"billing_model":    campaign.BillingModel,
			"revenue_amount":   campaign.RevenueAmount,
			"payout_amount":    campaign.PayoutAmount,
		}
	}

	// Platform revenue is the difference between spend and payout
	totalPlatformRevenue = totalAdvertiserSpend.Sub(totalAffiliatePayout)

	// TODO: Calculate affiliate breakdown
	// This would require tracking which affiliates generated which conversions/clicks

	metrics := &FinancialMetrics{
		AdvertiserSpend:    totalAdvertiserSpend,
		AffiliatePayout:    totalAffiliatePayout,
		PlatformRevenue:    totalPlatformRevenue,
		CampaignBreakdown:  campaignBreakdown,
		AffiliateBreakdown: affiliateBreakdown,
	}

	return metrics, nil
}

// processBilling processes billing based on the billing mode
func (s *UsageCalculationService) processBilling(ctx context.Context, billingAccount *domain.BillingAccount, usageRecord *domain.UsageRecord) error {
	if usageRecord.AdvertiserSpend.IsZero() {
		// No spend to bill
		usageRecord.Status = domain.UsageRecordStatusBilled
		usageRecord.BilledAt = timePtr(time.Now())
		return s.usageRecordRepo.Update(ctx, usageRecord)
	}

	switch billingAccount.BillingMode {
	case domain.BillingModePrepaid:
		return s.processPrepaidBilling(ctx, billingAccount, usageRecord)
	case domain.BillingModePostpaid:
		return s.processPostpaidBilling(ctx, billingAccount, usageRecord)
	default:
		return fmt.Errorf("unknown billing mode: %s", billingAccount.BillingMode)
	}
}

// processPrepaidBilling processes billing for prepaid accounts
func (s *UsageCalculationService) processPrepaidBilling(ctx context.Context, billingAccount *domain.BillingAccount, usageRecord *domain.UsageRecord) error {
	// Debit the advertiser spend from the account
	description := fmt.Sprintf("Daily usage charge for %s", usageRecord.UsageDate.Format("2006-01-02"))
	referenceType := "usage_record"
	referenceID := fmt.Sprintf("%d", usageRecord.UsageRecordID)

	_, err := s.billingService.DebitAccount(
		ctx,
		billingAccount.OrganizationID,
		usageRecord.AdvertiserSpend,
		description,
		&referenceType,
		&referenceID,
	)

	if err != nil {
		// If insufficient funds, mark as failed
		usageRecord.Status = domain.UsageRecordStatusFailed
		usageRecord.Metadata["error"] = err.Error()
		return s.usageRecordRepo.Update(ctx, usageRecord)
	}

	// Mark as billed
	usageRecord.Status = domain.UsageRecordStatusBilled
	usageRecord.BilledAt = timePtr(time.Now())
	return s.usageRecordRepo.Update(ctx, usageRecord)
}

// processPostpaidBilling processes billing for postpaid accounts
func (s *UsageCalculationService) processPostpaidBilling(ctx context.Context, billingAccount *domain.BillingAccount, usageRecord *domain.UsageRecord) error {
	// For postpaid accounts, we accumulate charges and bill monthly
	// Mark as billed (will be included in next invoice)
	usageRecord.Status = domain.UsageRecordStatusBilled
	usageRecord.BilledAt = timePtr(time.Now())
	return s.usageRecordRepo.Update(ctx, usageRecord)
}

// ProcessAffiliatePayout processes payout to affiliates
func (s *UsageCalculationService) ProcessAffiliatePayout(ctx context.Context, usageRecord *domain.UsageRecord) error {
	// TODO: Implement affiliate payout logic
	// This would:
	// 1. Calculate individual affiliate payouts from the breakdown
	// 2. Create payout transactions
	// 3. Integrate with payment providers (Stripe Connect, PayPal, etc.)
	// 4. Update affiliate balances

	log.Printf("Processing affiliate payout for usage record %d: %s",
		usageRecord.UsageRecordID, usageRecord.AffiliatePayout.String())

	// Mark usage record as paid
	usageRecord.Status = domain.UsageRecordStatusPaid
	usageRecord.AllocatedAt = timePtr(time.Now())
	return s.usageRecordRepo.Update(ctx, usageRecord)
}

// Helper function
func timePtr(t time.Time) *time.Time {
	return &t
}