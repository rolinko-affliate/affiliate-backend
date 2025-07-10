package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/stripe"
	"github.com/affiliate-backend/internal/repository"
	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
	stripeLib "github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

// WebhookHandler handles Stripe webhook events
type WebhookHandler struct {
	stripeService      *stripe.Service
	billingService     *service.BillingService
	webhookEventRepo   repository.WebhookEventRepository
	billingAccountRepo repository.BillingAccountRepository
	transactionRepo    repository.TransactionRepository
	webhookSecret      string
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(
	stripeService *stripe.Service,
	billingService *service.BillingService,
	webhookEventRepo repository.WebhookEventRepository,
	billingAccountRepo repository.BillingAccountRepository,
	transactionRepo repository.TransactionRepository,
	webhookSecret string,
) *WebhookHandler {
	return &WebhookHandler{
		stripeService:      stripeService,
		billingService:     billingService,
		webhookEventRepo:   webhookEventRepo,
		billingAccountRepo: billingAccountRepo,
		transactionRepo:    transactionRepo,
		webhookSecret:      webhookSecret,
	}
}

// HandleStripeWebhook handles incoming Stripe webhook events
func (h *WebhookHandler) HandleStripeWebhook(c *gin.Context) {
	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading webhook body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading request body"})
		return
	}

	// Verify the webhook signature
	event, err := webhook.ConstructEvent(body, c.GetHeader("Stripe-Signature"), h.webhookSecret)
	if err != nil {
		log.Printf("Error verifying webhook signature: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}

	// Store the webhook event
	webhookEvent := &domain.WebhookEvent{
		StripeEventID: event.ID,
		EventType:     string(event.Type),
		Status:        domain.WebhookEventStatusPending,
		EventData:     make(map[string]interface{}),
		RetryCount:    0,
	}

	// Convert event data to map
	eventDataBytes, _ := json.Marshal(event.Data)
	json.Unmarshal(eventDataBytes, &webhookEvent.EventData)

	// Save webhook event to database
	err = h.webhookEventRepo.Create(c.Request.Context(), webhookEvent)
	if err != nil {
		log.Printf("Error storing webhook event: %v", err)
		// Continue processing even if storage fails
	}

	// Process the event
	err = h.processWebhookEvent(c.Request.Context(), &event, webhookEvent)
	if err != nil {
		log.Printf("Error processing webhook event %s: %v", event.ID, err)
		
		// Update webhook event status to failed
		webhookEvent.Status = domain.WebhookEventStatusFailed
		webhookEvent.ErrorMessage = stringPtr(err.Error())
		webhookEvent.ProcessedAt = timePtr(time.Now())
		h.webhookEventRepo.Update(c.Request.Context(), webhookEvent)
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing webhook"})
		return
	}

	// Update webhook event status to processed
	webhookEvent.Status = domain.WebhookEventStatusProcessed
	webhookEvent.ProcessedAt = timePtr(time.Now())
	h.webhookEventRepo.Update(c.Request.Context(), webhookEvent)

	c.JSON(http.StatusOK, gin.H{"received": true})
}

// processWebhookEvent processes different types of Stripe webhook events
func (h *WebhookHandler) processWebhookEvent(ctx context.Context, event *stripeLib.Event, webhookEvent *domain.WebhookEvent) error {
	switch event.Type {
	case "payment_intent.succeeded":
		return h.handlePaymentIntentSucceeded(ctx, event, webhookEvent)
	case "payment_intent.payment_failed":
		return h.handlePaymentIntentFailed(ctx, event, webhookEvent)
	case "invoice.payment_succeeded":
		return h.handleInvoicePaymentSucceeded(ctx, event, webhookEvent)
	case "invoice.payment_failed":
		return h.handleInvoicePaymentFailed(ctx, event, webhookEvent)
	case "customer.subscription.created":
		return h.handleSubscriptionCreated(ctx, event, webhookEvent)
	case "customer.subscription.updated":
		return h.handleSubscriptionUpdated(ctx, event, webhookEvent)
	case "customer.subscription.deleted":
		return h.handleSubscriptionDeleted(ctx, event, webhookEvent)
	default:
		log.Printf("Unhandled webhook event type: %s", event.Type)
		webhookEvent.Status = domain.WebhookEventStatusIgnored
		return nil
	}
}

// handlePaymentIntentSucceeded handles successful payment intents
func (h *WebhookHandler) handlePaymentIntentSucceeded(ctx context.Context, event *stripeLib.Event, webhookEvent *domain.WebhookEvent) error {
	var paymentIntent stripeLib.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		return fmt.Errorf("error parsing payment intent: %w", err)
	}

	log.Printf("Processing successful payment intent: %s", paymentIntent.ID)

	// Find the corresponding transaction
	transaction, err := h.transactionRepo.GetByStripePaymentIntentID(ctx, paymentIntent.ID)
	if err != nil {
		return fmt.Errorf("transaction not found for payment intent %s: %w", paymentIntent.ID, err)
	}

	// Update transaction status
	transaction.Status = domain.TransactionStatusCompleted
	transaction.ProcessedAt = time.Now()

	// Add Stripe charge ID if available
	if paymentIntent.LatestCharge != nil {
		transaction.StripeChargeID = &paymentIntent.LatestCharge.ID
	}

	err = h.transactionRepo.Update(ctx, transaction)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	// Update billing account balance if this was a recharge
	if transaction.Type == domain.TransactionTypeRecharge {
		err = h.billingAccountRepo.UpdateBalance(ctx, transaction.BillingAccountID, transaction.BalanceAfter)
		if err != nil {
			return fmt.Errorf("failed to update billing account balance: %w", err)
		}
	}

	// Update webhook event with related records
	webhookEvent.OrganizationID = &transaction.OrganizationID
	webhookEvent.TransactionID = &transaction.TransactionID

	log.Printf("Successfully processed payment intent %s for organization %d", paymentIntent.ID, transaction.OrganizationID)
	return nil
}

// handlePaymentIntentFailed handles failed payment intents
func (h *WebhookHandler) handlePaymentIntentFailed(ctx context.Context, event *stripeLib.Event, webhookEvent *domain.WebhookEvent) error {
	var paymentIntent stripeLib.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		return fmt.Errorf("error parsing payment intent: %w", err)
	}

	log.Printf("Processing failed payment intent: %s", paymentIntent.ID)

	// Find the corresponding transaction
	transaction, err := h.transactionRepo.GetByStripePaymentIntentID(ctx, paymentIntent.ID)
	if err != nil {
		return fmt.Errorf("transaction not found for payment intent %s: %w", paymentIntent.ID, err)
	}

	// Update transaction status
	transaction.Status = domain.TransactionStatusFailed
	transaction.ProcessedAt = time.Now()

	// Add failure reason to metadata
	if transaction.Metadata == nil {
		transaction.Metadata = make(map[string]interface{})
	}
	if paymentIntent.LastPaymentError != nil {
		transaction.Metadata["failure_reason"] = paymentIntent.LastPaymentError.DeclineCode
		transaction.Metadata["failure_code"] = paymentIntent.LastPaymentError.Code
	}

	err = h.transactionRepo.Update(ctx, transaction)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	// Update webhook event with related records
	webhookEvent.OrganizationID = &transaction.OrganizationID
	webhookEvent.TransactionID = &transaction.TransactionID

	log.Printf("Successfully processed failed payment intent %s for organization %d", paymentIntent.ID, transaction.OrganizationID)
	return nil
}

// handleInvoicePaymentSucceeded handles successful invoice payments
func (h *WebhookHandler) handleInvoicePaymentSucceeded(ctx context.Context, event *stripeLib.Event, webhookEvent *domain.WebhookEvent) error {
	var invoice stripeLib.Invoice
	err := json.Unmarshal(event.Data.Raw, &invoice)
	if err != nil {
		return fmt.Errorf("error parsing invoice: %w", err)
	}

	log.Printf("Processing successful invoice payment: %s", invoice.ID)

	// Find billing account by Stripe customer ID
	billingAccount, err := h.billingAccountRepo.GetByStripeCustomerID(ctx, invoice.Customer.ID)
	if err != nil {
		return fmt.Errorf("billing account not found for customer %s: %w", invoice.Customer.ID, err)
	}

	// Create a transaction for the invoice payment
	amount := h.stripeService.ConvertAmountFromCents(invoice.AmountPaid)
	
	transaction := &domain.Transaction{
		OrganizationID:   billingAccount.OrganizationID,
		BillingAccountID: billingAccount.BillingAccountID,
		Type:             domain.TransactionTypeInvoicePayment,
		Amount:           amount,
		Currency:         string(invoice.Currency),
		BalanceBefore:    billingAccount.Balance,
		BalanceAfter:     billingAccount.Balance.Add(amount),
		ReferenceType:    stringPtr("stripe_invoice"),
		ReferenceID:      &invoice.ID,
		StripeInvoiceID:  &invoice.ID,
		Description:      stringPtr(fmt.Sprintf("Invoice payment for %s", invoice.Number)),
		Status:           domain.TransactionStatusCompleted,
		ProcessedAt:      time.Now(),
		Metadata:         make(map[string]interface{}),
	}

	transaction.Metadata["stripe_invoice_id"] = invoice.ID
	transaction.Metadata["invoice_number"] = invoice.Number

	err = h.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return fmt.Errorf("failed to create invoice payment transaction: %w", err)
	}

	// Update billing account balance for postpaid accounts
	if billingAccount.BillingMode == domain.BillingModePostpaid {
		err = h.billingAccountRepo.UpdateBalance(ctx, billingAccount.BillingAccountID, transaction.BalanceAfter)
		if err != nil {
			return fmt.Errorf("failed to update billing account balance: %w", err)
		}
	}

	// Update webhook event with related records
	webhookEvent.OrganizationID = &billingAccount.OrganizationID
	webhookEvent.TransactionID = &transaction.TransactionID

	log.Printf("Successfully processed invoice payment %s for organization %d", invoice.ID, billingAccount.OrganizationID)
	return nil
}

// handleInvoicePaymentFailed handles failed invoice payments
func (h *WebhookHandler) handleInvoicePaymentFailed(ctx context.Context, event *stripeLib.Event, webhookEvent *domain.WebhookEvent) error {
	var invoice stripeLib.Invoice
	err := json.Unmarshal(event.Data.Raw, &invoice)
	if err != nil {
		return fmt.Errorf("error parsing invoice: %w", err)
	}

	log.Printf("Processing failed invoice payment: %s", invoice.ID)

	// Find billing account by Stripe customer ID
	billingAccount, err := h.billingAccountRepo.GetByStripeCustomerID(ctx, invoice.Customer.ID)
	if err != nil {
		return fmt.Errorf("billing account not found for customer %s: %w", invoice.Customer.ID, err)
	}

	// Update webhook event with related records
	webhookEvent.OrganizationID = &billingAccount.OrganizationID

	// TODO: Implement invoice payment failure handling
	// - Send notification to organization
	// - Update invoice status
	// - Potentially suspend account if multiple failures

	log.Printf("Successfully processed failed invoice payment %s for organization %d", invoice.ID, billingAccount.OrganizationID)
	return nil
}

// handleSubscriptionCreated handles subscription creation events
func (h *WebhookHandler) handleSubscriptionCreated(ctx context.Context, event *stripeLib.Event, webhookEvent *domain.WebhookEvent) error {
	// TODO: Implement subscription handling if needed
	log.Printf("Subscription created event received: %s", event.ID)
	webhookEvent.Status = domain.WebhookEventStatusIgnored
	return nil
}

// handleSubscriptionUpdated handles subscription update events
func (h *WebhookHandler) handleSubscriptionUpdated(ctx context.Context, event *stripeLib.Event, webhookEvent *domain.WebhookEvent) error {
	// TODO: Implement subscription handling if needed
	log.Printf("Subscription updated event received: %s", event.ID)
	webhookEvent.Status = domain.WebhookEventStatusIgnored
	return nil
}

// handleSubscriptionDeleted handles subscription deletion events
func (h *WebhookHandler) handleSubscriptionDeleted(ctx context.Context, event *stripeLib.Event, webhookEvent *domain.WebhookEvent) error {
	// TODO: Implement subscription handling if needed
	log.Printf("Subscription deleted event received: %s", event.ID)
	webhookEvent.Status = domain.WebhookEventStatusIgnored
	return nil
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}