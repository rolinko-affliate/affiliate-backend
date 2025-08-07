package stripe

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/affiliate-backend/internal/domain"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/invoice"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/paymentmethod"
	"github.com/stripe/stripe-go/v76/setupintent"
)

// Config holds Stripe configuration
type Config struct {
	SecretKey      string
	PublishableKey string
	WebhookSecret  string
	Environment    string // "test" or "live"
}

// Service provides Stripe integration functionality
type Service struct {
	config Config
}

// NewService creates a new Stripe service
func NewService(config Config) *Service {
	stripe.Key = config.SecretKey
	return &Service{config: config}
}

// CreateCustomer creates a new Stripe customer for an organization
func (s *Service) CreateCustomer(ctx context.Context, org *domain.Organization, billingEmail string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Name:  stripe.String(org.Name),
		Email: stripe.String(billingEmail),
		Metadata: map[string]string{
			"organization_id":   fmt.Sprintf("%d", org.OrganizationID),
			"organization_type": string(org.Type),
		},
	}

	customer, err := customer.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe customer: %w", err)
	}

	slog.Info("Created Stripe customer",
		"customer_id", customer.ID,
		"organization_id", org.OrganizationID)
	return customer, nil
}

// UpdateCustomer updates a Stripe customer
func (s *Service) UpdateCustomer(ctx context.Context, customerID string, params *stripe.CustomerParams) (*stripe.Customer, error) {
	customer, err := customer.Update(customerID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update Stripe customer: %w", err)
	}

	return customer, nil
}

// GetCustomer retrieves a Stripe customer
func (s *Service) GetCustomer(ctx context.Context, customerID string) (*stripe.Customer, error) {
	customer, err := customer.Get(customerID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get Stripe customer: %w", err)
	}

	return customer, nil
}

// CreateSetupIntent creates a setup intent for adding payment methods
func (s *Service) CreateSetupIntent(ctx context.Context, customerID string) (*stripe.SetupIntent, error) {
	params := &stripe.SetupIntentParams{
		Customer: stripe.String(customerID),
		Usage:    stripe.String("off_session"),
	}

	intent, err := setupintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create setup intent: %w", err)
	}

	return intent, nil
}

// AttachPaymentMethod attaches a payment method to a customer
func (s *Service) AttachPaymentMethod(ctx context.Context, paymentMethodID, customerID string) (*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customerID),
	}

	pm, err := paymentmethod.Attach(paymentMethodID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to attach payment method: %w", err)
	}

	return pm, nil
}

// DetachPaymentMethod detaches a payment method from a customer
func (s *Service) DetachPaymentMethod(ctx context.Context, paymentMethodID string) (*stripe.PaymentMethod, error) {
	pm, err := paymentmethod.Detach(paymentMethodID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to detach payment method: %w", err)
	}

	return pm, nil
}

// GetPaymentMethod retrieves a payment method
func (s *Service) GetPaymentMethod(ctx context.Context, paymentMethodID string) (*stripe.PaymentMethod, error) {
	pm, err := paymentmethod.Get(paymentMethodID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}

	return pm, nil
}

// ListPaymentMethods lists payment methods for a customer
func (s *Service) ListPaymentMethods(ctx context.Context, customerID string) ([]*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(customerID),
		Type:     stripe.String("card"),
	}

	iter := paymentmethod.List(params)
	var paymentMethods []*stripe.PaymentMethod

	for iter.Next() {
		paymentMethods = append(paymentMethods, iter.PaymentMethod())
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to list payment methods: %w", err)
	}

	return paymentMethods, nil
}

// CreatePaymentIntent creates a payment intent for charging a customer
func (s *Service) CreatePaymentIntent(ctx context.Context, amount decimal.Decimal, currency, customerID string, paymentMethodID *string, description *string) (*stripe.PaymentIntent, error) {
	// Convert decimal amount to cents (Stripe expects amounts in smallest currency unit)
	amountCents := amount.Mul(decimal.NewFromInt(100)).IntPart()

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountCents),
		Currency: stripe.String(currency),
		Customer: stripe.String(customerID),
		Metadata: map[string]string{
			"type": "recharge",
		},
	}

	if paymentMethodID != nil {
		params.PaymentMethod = stripe.String(*paymentMethodID)
		params.ConfirmationMethod = stripe.String("manual")
		params.Confirm = stripe.Bool(true)
	}

	if description != nil {
		params.Description = stripe.String(*description)
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	return intent, nil
}

// ConfirmPaymentIntent confirms a payment intent
func (s *Service) ConfirmPaymentIntent(ctx context.Context, paymentIntentID string, paymentMethodID *string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentConfirmParams{}

	if paymentMethodID != nil {
		params.PaymentMethod = stripe.String(*paymentMethodID)
	}

	intent, err := paymentintent.Confirm(paymentIntentID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment intent: %w", err)
	}

	return intent, nil
}

// GetPaymentIntent retrieves a payment intent
func (s *Service) GetPaymentIntent(ctx context.Context, paymentIntentID string) (*stripe.PaymentIntent, error) {
	intent, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}

	return intent, nil
}

// CreateInvoice creates a Stripe invoice for postpaid billing
func (s *Service) CreateInvoice(ctx context.Context, customerID string, lineItems []domain.InvoiceLineItem, dueDate int64) (*stripe.Invoice, error) {
	params := &stripe.InvoiceParams{
		Customer: stripe.String(customerID),
		DueDate:  stripe.Int64(dueDate),
		Metadata: map[string]string{
			"type": "usage_billing",
		},
	}

	inv, err := invoice.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	return inv, nil
}

// FinalizeInvoice finalizes a draft invoice
func (s *Service) FinalizeInvoice(ctx context.Context, invoiceID string) (*stripe.Invoice, error) {
	params := &stripe.InvoiceFinalizeInvoiceParams{}

	inv, err := invoice.FinalizeInvoice(invoiceID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to finalize invoice: %w", err)
	}

	return inv, nil
}

// SendInvoice sends an invoice to the customer
func (s *Service) SendInvoice(ctx context.Context, invoiceID string) (*stripe.Invoice, error) {
	params := &stripe.InvoiceSendInvoiceParams{}

	inv, err := invoice.SendInvoice(invoiceID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to send invoice: %w", err)
	}

	return inv, nil
}

// GetInvoice retrieves an invoice
func (s *Service) GetInvoice(ctx context.Context, invoiceID string) (*stripe.Invoice, error) {
	inv, err := invoice.Get(invoiceID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	return inv, nil
}

// PayInvoice pays an invoice using the customer's default payment method
func (s *Service) PayInvoice(ctx context.Context, invoiceID string) (*stripe.Invoice, error) {
	params := &stripe.InvoicePayParams{}

	inv, err := invoice.Pay(invoiceID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to pay invoice: %w", err)
	}

	return inv, nil
}

// ConvertStripePaymentMethodToDomain converts a Stripe payment method to domain model
func (s *Service) ConvertStripePaymentMethodToDomain(stripePM *stripe.PaymentMethod, organizationID, billingAccountID int64) *domain.StripePaymentMethod {
	pm := &domain.StripePaymentMethod{
		OrganizationID:        organizationID,
		BillingAccountID:      billingAccountID,
		StripePaymentMethodID: stripePM.ID,
		Type:                  string(stripePM.Type),
		Status:                domain.PaymentMethodStatusActive,
		Metadata:              make(map[string]interface{}),
	}

	// Copy metadata
	for k, v := range stripePM.Metadata {
		pm.Metadata[k] = v
	}

	// Handle card details
	if stripePM.Card != nil {
		brand := string(stripePM.Card.Brand)
		pm.Brand = &brand
		pm.Last4 = &stripePM.Card.Last4
		expMonth := int(stripePM.Card.ExpMonth)
		expYear := int(stripePM.Card.ExpYear)
		pm.ExpMonth = &expMonth
		pm.ExpYear = &expYear
	}

	// Handle bank account details
	if stripePM.USBankAccount != nil {
		pm.BankName = &stripePM.USBankAccount.BankName
		pm.Last4 = &stripePM.USBankAccount.Last4
		accountHolderType := string(stripePM.USBankAccount.AccountHolderType)
		pm.AccountHolderType = &accountHolderType
	}

	return pm
}

// ConvertAmountFromCents converts Stripe amount (in cents) to decimal
func (s *Service) ConvertAmountFromCents(amountCents int64) decimal.Decimal {
	return decimal.NewFromInt(amountCents).Div(decimal.NewFromInt(100))
}

// ConvertAmountToCents converts decimal amount to Stripe amount (in cents)
func (s *Service) ConvertAmountToCents(amount decimal.Decimal) int64 {
	return amount.Mul(decimal.NewFromInt(100)).IntPart()
}
