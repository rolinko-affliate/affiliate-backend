package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/stripe"
	"github.com/affiliate-backend/internal/repository"
	"github.com/shopspring/decimal"
	stripeLib "github.com/stripe/stripe-go/v76"
)

// BillingService provides billing and payment functionality
type BillingService struct {
	billingAccountRepo repository.BillingAccountRepository
	paymentMethodRepo  repository.PaymentMethodRepository
	transactionRepo    repository.TransactionRepository
	organizationRepo   repository.OrganizationRepository
	stripeService      *stripe.Service
}

// NewBillingService creates a new billing service
func NewBillingService(
	billingAccountRepo repository.BillingAccountRepository,
	paymentMethodRepo repository.PaymentMethodRepository,
	transactionRepo repository.TransactionRepository,
	organizationRepo repository.OrganizationRepository,
	stripeService *stripe.Service,
) *BillingService {
	return &BillingService{
		billingAccountRepo: billingAccountRepo,
		paymentMethodRepo:  paymentMethodRepo,
		transactionRepo:    transactionRepo,
		organizationRepo:   organizationRepo,
		stripeService:      stripeService,
	}
}

// GetOrCreateBillingAccount gets or creates a billing account for an organization
func (s *BillingService) GetOrCreateBillingAccount(ctx context.Context, organizationID int64) (*domain.BillingAccount, error) {
	// Try to get existing billing account
	account, err := s.billingAccountRepo.GetByOrganizationID(ctx, organizationID)
	if err == nil {
		return account, nil
	}

	// If not found, create a new one
	org, err := s.organizationRepo.GetOrganizationByID(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	// Create Stripe customer
	billingEmail := fmt.Sprintf("billing@%s.com", org.Name) // Default email, should be configurable
	stripeCustomer, err := s.stripeService.CreateCustomer(ctx, org, billingEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe customer: %w", err)
	}

	// Create billing account
	account = &domain.BillingAccount{
		OrganizationID:        organizationID,
		StripeCustomerID:      &stripeCustomer.ID,
		BillingMode:           domain.BillingModePrepaid,
		Currency:              "USD",
		Balance:               decimal.Zero,
		CreditLimit:           decimal.Zero,
		AutoRechargeEnabled:   false,
		AutoRechargeThreshold: decimal.Zero,
		AutoRechargeAmount:    decimal.Zero,
		InvoiceDayOfMonth:     1,
		PaymentTermsDays:      30,
		Status:                domain.BillingAccountStatusActive,
		BillingEmail:          &billingEmail,
		BillingAddress:        make(map[string]interface{}),
		TaxInfo:               make(map[string]interface{}),
	}

	err = s.billingAccountRepo.Create(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("failed to create billing account: %w", err)
	}

	log.Printf("Created billing account %d for organization %d", account.BillingAccountID, organizationID)
	return account, nil
}

// GetBillingDashboard returns billing dashboard data for an organization
func (s *BillingService) GetBillingDashboard(ctx context.Context, organizationID int64) (*domain.BillingDashboardResponse, error) {
	// Get billing account
	account, err := s.GetOrCreateBillingAccount(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing account: %w", err)
	}

	// Get payment methods
	paymentMethods, err := s.paymentMethodRepo.GetByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}

	// Get recent transactions
	recentTransactions, err := s.transactionRepo.GetByOrganizationID(ctx, organizationID, 10, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent transactions: %w", err)
	}

	// Calculate monthly spend
	now := time.Now()
	monthlySpend, err := s.transactionRepo.GetMonthlySpend(ctx, organizationID, now.Year(), int(now.Month()))
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly spend: %w", err)
	}

	dashboard := &domain.BillingDashboardResponse{
		BillingAccount:     account,
		PaymentMethods:     paymentMethods,
		RecentTransactions: recentTransactions,
		CurrentBalance:     account.Balance,
		MonthlySpend:       monthlySpend,
		PendingInvoices:    []domain.Invoice{}, // TODO: Implement invoice retrieval
	}

	return dashboard, nil
}

// AddPaymentMethod adds a new payment method for an organization
func (s *BillingService) AddPaymentMethod(ctx context.Context, organizationID int64, req *domain.CreatePaymentMethodRequest) (*domain.StripePaymentMethod, error) {
	// Get billing account
	account, err := s.GetOrCreateBillingAccount(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing account: %w", err)
	}

	if account.StripeCustomerID == nil {
		return nil, fmt.Errorf("billing account has no Stripe customer ID")
	}

	// Attach payment method to Stripe customer
	stripePM, err := s.stripeService.AttachPaymentMethod(ctx, req.PaymentMethodID, *account.StripeCustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to attach payment method to Stripe customer: %w", err)
	}

	// Convert to domain model
	paymentMethod := s.stripeService.ConvertStripePaymentMethodToDomain(stripePM, organizationID, account.BillingAccountID)
	paymentMethod.IsDefault = req.SetAsDefault
	paymentMethod.Nickname = req.Nickname

	// Save to database
	err = s.paymentMethodRepo.Create(ctx, paymentMethod)
	if err != nil {
		return nil, fmt.Errorf("failed to save payment method: %w", err)
	}

	// Set as default if requested
	if req.SetAsDefault {
		err = s.paymentMethodRepo.SetAsDefault(ctx, paymentMethod.PaymentMethodID, organizationID)
		if err != nil {
			return nil, fmt.Errorf("failed to set payment method as default: %w", err)
		}
	}

	log.Printf("Added payment method %d for organization %d", paymentMethod.PaymentMethodID, organizationID)
	return paymentMethod, nil
}

// RemovePaymentMethod removes a payment method for an organization
func (s *BillingService) RemovePaymentMethod(ctx context.Context, organizationID int64, paymentMethodID int64) error {
	// Get payment method
	paymentMethod, err := s.paymentMethodRepo.GetByID(ctx, paymentMethodID)
	if err != nil {
		return fmt.Errorf("failed to get payment method: %w", err)
	}

	// Verify ownership
	if paymentMethod.OrganizationID != organizationID {
		return fmt.Errorf("payment method does not belong to organization")
	}

	// Detach from Stripe
	_, err = s.stripeService.DetachPaymentMethod(ctx, paymentMethod.StripePaymentMethodID)
	if err != nil {
		log.Printf("Warning: failed to detach payment method from Stripe: %v", err)
		// Continue with local deletion even if Stripe fails
	}

	// Delete from database (soft delete)
	err = s.paymentMethodRepo.Delete(ctx, paymentMethodID)
	if err != nil {
		return fmt.Errorf("failed to delete payment method: %w", err)
	}

	log.Printf("Removed payment method %d for organization %d", paymentMethodID, organizationID)
	return nil
}

// Recharge adds funds to an organization's account
func (s *BillingService) Recharge(ctx context.Context, organizationID int64, req *domain.RechargeRequest) (*domain.Transaction, error) {
	// Get billing account
	account, err := s.GetOrCreateBillingAccount(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing account: %w", err)
	}

	if account.StripeCustomerID == nil {
		return nil, fmt.Errorf("billing account has no Stripe customer ID")
	}

	// Validate amount
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("recharge amount must be positive")
	}

	// Get payment method
	var paymentMethodID *string
	if req.PaymentMethodID != nil {
		paymentMethodID = req.PaymentMethodID
	} else {
		// Use default payment method
		defaultPM, err := s.paymentMethodRepo.GetDefaultByOrganizationID(ctx, organizationID)
		if err != nil {
			return nil, fmt.Errorf("no payment method specified and no default payment method found: %w", err)
		}
		paymentMethodID = &defaultPM.StripePaymentMethodID
	}

	// Create payment intent
	currency := req.Currency
	if currency == "" {
		currency = account.Currency
	}

	paymentIntent, err := s.stripeService.CreatePaymentIntent(
		ctx,
		req.Amount,
		currency,
		*account.StripeCustomerID,
		paymentMethodID,
		req.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	// Create transaction record
	transaction := &domain.Transaction{
		OrganizationID:        organizationID,
		BillingAccountID:      account.BillingAccountID,
		Type:                  domain.TransactionTypeRecharge,
		Amount:                req.Amount,
		Currency:              currency,
		BalanceBefore:         account.Balance,
		BalanceAfter:          account.Balance.Add(req.Amount),
		ReferenceType:         stringPtr("stripe_payment_intent"),
		ReferenceID:           &paymentIntent.ID,
		StripePaymentIntentID: &paymentIntent.ID,
		Description:           req.Description,
		Status:                s.convertStripePaymentIntentStatus(paymentIntent.Status),
		ProcessedAt:           time.Now(),
		Metadata:              make(map[string]interface{}),
	}

	// Add metadata
	transaction.Metadata["stripe_payment_intent_id"] = paymentIntent.ID
	transaction.Metadata["payment_method_id"] = *paymentMethodID

	err = s.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Update balance if payment succeeded
	if paymentIntent.Status == stripeLib.PaymentIntentStatusSucceeded {
		err = s.billingAccountRepo.UpdateBalance(ctx, account.BillingAccountID, transaction.BalanceAfter)
		if err != nil {
			return nil, fmt.Errorf("failed to update account balance: %w", err)
		}
	}

	log.Printf("Created recharge transaction %d for organization %d, amount: %s",
		transaction.TransactionID, organizationID, req.Amount.String())
	return transaction, nil
}

// DebitAccount debits an amount from an organization's account
func (s *BillingService) DebitAccount(ctx context.Context, organizationID int64, amount decimal.Decimal, description string, referenceType, referenceID *string) (*domain.Transaction, error) {
	// Get billing account
	account, err := s.GetOrCreateBillingAccount(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing account: %w", err)
	}

	// Check if sufficient balance (for prepaid accounts)
	if account.BillingMode == domain.BillingModePrepaid {
		if account.Balance.LessThan(amount) {
			return nil, fmt.Errorf("insufficient balance: current balance %s, required %s",
				account.Balance.String(), amount.String())
		}
	}

	// Create debit transaction
	transaction := &domain.Transaction{
		OrganizationID:   organizationID,
		BillingAccountID: account.BillingAccountID,
		Type:             domain.TransactionTypeDebit,
		Amount:           amount.Neg(), // Negative for debit
		Currency:         account.Currency,
		BalanceBefore:    account.Balance,
		BalanceAfter:     account.Balance.Sub(amount),
		ReferenceType:    referenceType,
		ReferenceID:      referenceID,
		Description:      &description,
		Status:           domain.TransactionStatusCompleted,
		ProcessedAt:      time.Now(),
		Metadata:         make(map[string]interface{}),
	}

	err = s.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create debit transaction: %w", err)
	}

	// Update balance
	err = s.billingAccountRepo.UpdateBalance(ctx, account.BillingAccountID, transaction.BalanceAfter)
	if err != nil {
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	log.Printf("Created debit transaction %d for organization %d, amount: %s",
		transaction.TransactionID, organizationID, amount.String())
	return transaction, nil
}

// UpdateBillingConfig updates billing configuration for an organization
func (s *BillingService) UpdateBillingConfig(ctx context.Context, organizationID int64, req *domain.UpdateBillingConfigRequest) (*domain.BillingAccount, error) {
	// Get billing account
	account, err := s.GetOrCreateBillingAccount(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing account: %w", err)
	}

	// Update fields if provided
	if req.BillingMode != nil {
		account.BillingMode = *req.BillingMode
	}
	if req.AutoRechargeEnabled != nil {
		account.AutoRechargeEnabled = *req.AutoRechargeEnabled
	}
	if req.AutoRechargeThreshold != nil {
		account.AutoRechargeThreshold = *req.AutoRechargeThreshold
	}
	if req.AutoRechargeAmount != nil {
		account.AutoRechargeAmount = *req.AutoRechargeAmount
	}
	if req.BillingEmail != nil {
		account.BillingEmail = req.BillingEmail
	}
	if req.InvoiceDayOfMonth != nil {
		account.InvoiceDayOfMonth = *req.InvoiceDayOfMonth
	}
	if req.PaymentTermsDays != nil {
		account.PaymentTermsDays = *req.PaymentTermsDays
	}

	// Save changes
	err = s.billingAccountRepo.Update(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("failed to update billing account: %w", err)
	}

	log.Printf("Updated billing config for organization %d", organizationID)
	return account, nil
}

// GetTransactionHistory returns transaction history for an organization
func (s *BillingService) GetTransactionHistory(ctx context.Context, organizationID int64, limit, offset int) ([]domain.Transaction, error) {
	transactions, err := s.transactionRepo.GetByOrganizationID(ctx, organizationID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction history: %w", err)
	}

	return transactions, nil
}

// Helper functions

func (s *BillingService) convertStripePaymentIntentStatus(status stripeLib.PaymentIntentStatus) domain.TransactionStatus {
	switch status {
	case stripeLib.PaymentIntentStatusSucceeded:
		return domain.TransactionStatusCompleted
	case stripeLib.PaymentIntentStatusProcessing:
		return domain.TransactionStatusPending
	case stripeLib.PaymentIntentStatusRequiresPaymentMethod,
		stripeLib.PaymentIntentStatusRequiresConfirmation,
		stripeLib.PaymentIntentStatusRequiresAction:
		return domain.TransactionStatusPending
	case stripeLib.PaymentIntentStatusCanceled:
		return domain.TransactionStatusCancelled
	default:
		return domain.TransactionStatusFailed
	}
}
