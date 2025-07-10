package domain

import (
	"fmt"
	"time"
	
	"github.com/shopspring/decimal"
)

// BillingFrequency represents the billing frequency options
type BillingFrequency string

const (
	BillingFrequencyWeekly     BillingFrequency = "weekly"
	BillingFrequencyBimonthly  BillingFrequency = "bimonthly"
	BillingFrequencyMonthly    BillingFrequency = "monthly"
	BillingFrequencyTwoMonths  BillingFrequency = "two_months"
	BillingFrequencyQuarterly  BillingFrequency = "quarterly"
	BillingFrequencyManual     BillingFrequency = "manual"
	BillingFrequencyOther      BillingFrequency = "other"
)

// PaymentType represents the payment method options
type PaymentType string

const (
	PaymentTypeWire         PaymentType = "wire"
	PaymentTypeACH          PaymentType = "ach"
	PaymentTypeCheck        PaymentType = "check"
	PaymentTypePayPal       PaymentType = "paypal"
	PaymentTypeCrypto       PaymentType = "crypto"
	PaymentTypeOther        PaymentType = "other"
)

// PaymentDetailsType represents the specific payment details type for validation
type PaymentDetailsType string

const (
	PaymentDetailsTypeWire   PaymentDetailsType = "wire"
	PaymentDetailsTypeACH    PaymentDetailsType = "ach"
	PaymentDetailsTypeCheck  PaymentDetailsType = "check"
	PaymentDetailsTypePayPal PaymentDetailsType = "paypal"
	PaymentDetailsTypeCrypto PaymentDetailsType = "crypto"
	PaymentDetailsTypeOther  PaymentDetailsType = "other"
)

// BillingSchedule represents billing schedule configuration based on frequency
type BillingSchedule struct {
	// For weekly billing
	DayOfWeek *int32 `json:"day_of_week,omitempty"` // 0=Sunday, 1=Monday, etc.
	
	// For monthly billing
	DayOfMonth    *int32 `json:"day_of_month,omitempty"`     // 1-31
	DayOfMonthOne *int32 `json:"day_of_month_one,omitempty"` // For bimonthly
	DayOfMonthTwo *int32 `json:"day_of_month_two,omitempty"` // For bimonthly
	
	// For quarterly/yearly billing
	StartingMonth *int32 `json:"starting_month,omitempty"` // 1-12
}

// Validate validates the billing schedule based on the frequency
func (bs *BillingSchedule) Validate(frequency BillingFrequency) error {
	switch frequency {
	case BillingFrequencyWeekly:
		if bs.DayOfWeek == nil || *bs.DayOfWeek < 0 || *bs.DayOfWeek > 6 {
			return fmt.Errorf("day_of_week must be between 0-6 for weekly billing")
		}
	case BillingFrequencyMonthly:
		if bs.DayOfMonth == nil || *bs.DayOfMonth < 1 || *bs.DayOfMonth > 31 {
			return fmt.Errorf("day_of_month must be between 1-31 for monthly billing")
		}
	case BillingFrequencyBimonthly:
		if bs.DayOfMonthOne == nil || *bs.DayOfMonthOne < 1 || *bs.DayOfMonthOne > 31 {
			return fmt.Errorf("day_of_month_one must be between 1-31 for bimonthly billing")
		}
		if bs.DayOfMonthTwo == nil || *bs.DayOfMonthTwo < 1 || *bs.DayOfMonthTwo > 31 {
			return fmt.Errorf("day_of_month_two must be between 1-31 for bimonthly billing")
		}
	case BillingFrequencyQuarterly:
		if bs.StartingMonth == nil || *bs.StartingMonth < 1 || *bs.StartingMonth > 12 {
			return fmt.Errorf("starting_month must be between 1-12 for quarterly billing")
		}
	}
	return nil
}

// PaymentDetails represents structured payment information
type PaymentDetails struct {
	// Payment type for validation and processing
	Type              *PaymentDetailsType `json:"type,omitempty"`
	
	// Wire transfer details
	BankName          *string `json:"bank_name,omitempty"`
	BankAddress       *string `json:"bank_address,omitempty"`
	AccountNumber     *string `json:"account_number,omitempty"`
	RoutingNumber     *string `json:"routing_number,omitempty"`
	SwiftCode         *string `json:"swift_code,omitempty"`
	IBAN              *string `json:"iban,omitempty"`
	
	// ACH details
	ACHAccountType    *string `json:"ach_account_type,omitempty"` // "checking", "savings"
	
	// PayPal details
	PayPalEmail       *string `json:"paypal_email,omitempty"`
	
	// Crypto details
	CryptoWalletType  *string `json:"crypto_wallet_type,omitempty"` // "bitcoin", "ethereum", etc.
	CryptoAddress     *string `json:"crypto_address,omitempty"`
	
	// Check details
	MailingAddress    *BillingAddress `json:"mailing_address,omitempty"`
	
	// Additional details for other payment types
	AdditionalDetails map[string]interface{} `json:"additional_details,omitempty"`
}

// Validate checks if the PaymentDetails has the required fields for its type
func (pd *PaymentDetails) Validate() error {
	if pd.Type == nil {
		return fmt.Errorf("payment details type is required")
	}

	switch *pd.Type {
	case PaymentDetailsTypeWire:
		if pd.BankName == nil || pd.AccountNumber == nil {
			return fmt.Errorf("wire transfer requires bank name and account number")
		}
		if pd.RoutingNumber == nil && pd.SwiftCode == nil && pd.IBAN == nil {
			return fmt.Errorf("wire transfer requires routing number, SWIFT code, or IBAN")
		}
	case PaymentDetailsTypeACH:
		if pd.BankName == nil || pd.AccountNumber == nil || pd.RoutingNumber == nil {
			return fmt.Errorf("ACH transfer requires bank name, account number, and routing number")
		}
	case PaymentDetailsTypePayPal:
		if pd.PayPalEmail == nil {
			return fmt.Errorf("PayPal payment requires PayPal email")
		}
	case PaymentDetailsTypeCrypto:
		if pd.CryptoAddress == nil || pd.CryptoWalletType == nil {
			return fmt.Errorf("crypto payment requires wallet address and wallet type")
		}
	case PaymentDetailsTypeCheck:
		if pd.MailingAddress == nil {
			return fmt.Errorf("check payment requires mailing address")
		}
	case PaymentDetailsTypeOther:
		// Other types are flexible, no specific validation
	default:
		return fmt.Errorf("unsupported payment details type: %s", *pd.Type)
	}

	return nil
}

// HasData returns true if the PaymentDetails has any meaningful data
func (pd *PaymentDetails) HasData() bool {
	if pd == nil {
		return false
	}
	
	return pd.Type != nil ||
		pd.BankName != nil ||
		pd.BankAddress != nil ||
		pd.AccountNumber != nil ||
		pd.RoutingNumber != nil ||
		pd.SwiftCode != nil ||
		pd.IBAN != nil ||
		pd.ACHAccountType != nil ||
		pd.PayPalEmail != nil ||
		pd.CryptoWalletType != nil ||
		pd.CryptoAddress != nil ||
		pd.MailingAddress != nil ||
		len(pd.AdditionalDetails) > 0
}

// BillingDetails represents structured billing information for an advertiser
type BillingDetails struct {
	Frequency                  *BillingFrequency       `json:"frequency,omitempty"`
	Schedule                   *BillingSchedule        `json:"schedule,omitempty"`
	PaymentType                *PaymentType            `json:"payment_type,omitempty"`
	PaymentDetails             *PaymentDetails         `json:"payment_details,omitempty"`
	TaxID                      *string                 `json:"tax_id,omitempty"`
	IsInvoiceCreationAuto      *bool                   `json:"is_invoice_creation_auto,omitempty"`
	InvoiceAmountThreshold     *float64                `json:"invoice_amount_threshold,omitempty"`
	AutoInvoiceStartDate       *string                 `json:"auto_invoice_start_date,omitempty"` // Format: "2019-06-01 00:00:00"
	DefaultInvoiceIsHidden     *bool                   `json:"default_invoice_is_hidden,omitempty"`
	InvoiceGenerationDaysDelay *int32                  `json:"invoice_generation_days_delay,omitempty"`
	DefaultPaymentTerms        *int                    `json:"default_payment_terms,omitempty"`
	Address                    *BillingAddress         `json:"address,omitempty"`
	AdditionalDetails          map[string]interface{}  `json:"additional_details,omitempty"`
}

// HasData returns true if any billing field has data
func (bd *BillingDetails) HasData() bool {
	if bd == nil {
		return false
	}
	return bd.Frequency != nil || bd.Schedule != nil || bd.PaymentType != nil || 
		bd.PaymentDetails != nil || bd.TaxID != nil ||
		bd.IsInvoiceCreationAuto != nil || bd.InvoiceAmountThreshold != nil ||
		bd.AutoInvoiceStartDate != nil || bd.DefaultInvoiceIsHidden != nil ||
		bd.InvoiceGenerationDaysDelay != nil || bd.DefaultPaymentTerms != nil ||
		bd.Address != nil || len(bd.AdditionalDetails) > 0
}

// BillingAddress represents billing address information
type BillingAddress struct {
	Line1       string  `json:"line1"`
	Line2       *string `json:"line2,omitempty"`
	City        string  `json:"city"`
	State       *string `json:"state,omitempty"`
	PostalCode  string  `json:"postal_code"`
	Country     string  `json:"country"`
	CompanyName *string `json:"company_name,omitempty"`
}

// PaymentMethod represents payment method information
type PaymentMethod struct {
	Type            string                 `json:"type"` // "bank_transfer", "check", "paypal", "wire", "other"
	BankAccountInfo *BankAccountInfo       `json:"bank_account_info,omitempty"`
	PaypalEmail     *string                `json:"paypal_email,omitempty"`
	CheckAddress    *BillingAddress        `json:"check_address,omitempty"`
	WireInfo        *WireTransferInfo      `json:"wire_info,omitempty"`
	OtherDetails    map[string]interface{} `json:"other_details,omitempty"`
}

// BankAccountInfo represents bank account information for payments
type BankAccountInfo struct {
	BankName      string  `json:"bank_name"`
	AccountNumber string  `json:"account_number"`
	RoutingNumber string  `json:"routing_number"`
	AccountType   string  `json:"account_type"` // "checking", "savings"
	AccountHolder string  `json:"account_holder"`
	SwiftCode     *string `json:"swift_code,omitempty"`
}

// WireTransferInfo represents wire transfer information
type WireTransferInfo struct {
	BankName           string  `json:"bank_name"`
	BankAddress        string  `json:"bank_address"`
	SwiftCode          string  `json:"swift_code"`
	AccountNumber      string  `json:"account_number"`
	AccountHolder      string  `json:"account_holder"`
	IntermediaryBank   *string `json:"intermediary_bank,omitempty"`
	IntermediarySwift  *string `json:"intermediary_swift,omitempty"`
	SpecialInstructions *string `json:"special_instructions,omitempty"`
}

// #############################################################################
// ## New Stripe-Based Billing System Models
// #############################################################################

// BillingMode represents the billing mode for an organization
type BillingMode string

const (
	BillingModePrepaid  BillingMode = "prepaid"  // Pay-as-you-go with balance
	BillingModePostpaid BillingMode = "postpaid" // Monthly invoicing
)

// IsValid checks if the billing mode is valid
func (bm BillingMode) IsValid() bool {
	switch bm {
	case BillingModePrepaid, BillingModePostpaid:
		return true
	default:
		return false
	}
}

// BillingAccountStatus represents the status of a billing account
type BillingAccountStatus string

const (
	BillingAccountStatusActive    BillingAccountStatus = "active"
	BillingAccountStatusSuspended BillingAccountStatus = "suspended"
	BillingAccountStatusClosed    BillingAccountStatus = "closed"
)

// BillingAccount represents organization-level billing configuration
type BillingAccount struct {
	BillingAccountID int64                `json:"billing_account_id" db:"billing_account_id"`
	OrganizationID   int64                `json:"organization_id" db:"organization_id"`
	
	// Stripe Integration
	StripeCustomerID *string `json:"stripe_customer_id,omitempty" db:"stripe_customer_id"`
	StripeAccountID  *string `json:"stripe_account_id,omitempty" db:"stripe_account_id"`
	
	// Billing Configuration
	BillingMode BillingMode `json:"billing_mode" db:"billing_mode"`
	Currency    string      `json:"currency" db:"currency"`
	
	// Balance and Credit
	Balance     decimal.Decimal `json:"balance" db:"balance"`
	CreditLimit decimal.Decimal `json:"credit_limit" db:"credit_limit"`
	
	// Payment Configuration
	DefaultPaymentMethodID *string         `json:"default_payment_method_id,omitempty" db:"default_payment_method_id"`
	AutoRechargeEnabled    bool            `json:"auto_recharge_enabled" db:"auto_recharge_enabled"`
	AutoRechargeThreshold  decimal.Decimal `json:"auto_recharge_threshold" db:"auto_recharge_threshold"`
	AutoRechargeAmount     decimal.Decimal `json:"auto_recharge_amount" db:"auto_recharge_amount"`
	
	// Invoice Configuration
	InvoiceDayOfMonth int `json:"invoice_day_of_month" db:"invoice_day_of_month"`
	PaymentTermsDays  int `json:"payment_terms_days" db:"payment_terms_days"`
	
	// Status and Metadata
	Status         BillingAccountStatus `json:"status" db:"status"`
	BillingEmail   *string              `json:"billing_email,omitempty" db:"billing_email"`
	BillingAddress map[string]interface{} `json:"billing_address,omitempty" db:"billing_address"`
	TaxInfo        map[string]interface{} `json:"tax_info,omitempty" db:"tax_info"`
	
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PaymentMethodStatus represents the status of a payment method
type PaymentMethodStatus string

const (
	PaymentMethodStatusActive   PaymentMethodStatus = "active"
	PaymentMethodStatusInactive PaymentMethodStatus = "inactive"
	PaymentMethodStatusExpired  PaymentMethodStatus = "expired"
	PaymentMethodStatusFailed   PaymentMethodStatus = "failed"
)

// StripePaymentMethod represents a Stripe payment method
type StripePaymentMethod struct {
	PaymentMethodID int64  `json:"payment_method_id" db:"payment_method_id"`
	OrganizationID  int64  `json:"organization_id" db:"organization_id"`
	BillingAccountID int64 `json:"billing_account_id" db:"billing_account_id"`
	
	// Stripe Integration
	StripePaymentMethodID string `json:"stripe_payment_method_id" db:"stripe_payment_method_id"`
	
	// Payment Method Details
	Type  string  `json:"type" db:"type"`
	Brand *string `json:"brand,omitempty" db:"brand"`
	Last4 *string `json:"last4,omitempty" db:"last4"`
	ExpMonth *int `json:"exp_month,omitempty" db:"exp_month"`
	ExpYear  *int `json:"exp_year,omitempty" db:"exp_year"`
	
	// Bank Account Details
	BankName           *string `json:"bank_name,omitempty" db:"bank_name"`
	AccountHolderType  *string `json:"account_holder_type,omitempty" db:"account_holder_type"`
	
	// Status and Configuration
	IsDefault bool                   `json:"is_default" db:"is_default"`
	Status    PaymentMethodStatus    `json:"status" db:"status"`
	
	// Metadata
	Nickname *string                `json:"nickname,omitempty" db:"nickname"`
	Metadata map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeRecharge        TransactionType = "recharge"
	TransactionTypeDebit          TransactionType = "debit"
	TransactionTypeCredit         TransactionType = "credit"
	TransactionTypeRefund         TransactionType = "refund"
	TransactionTypeChargeback     TransactionType = "chargeback"
	TransactionTypeInvoicePayment TransactionType = "invoice_payment"
	TransactionTypeUsageCharge    TransactionType = "usage_charge"
	TransactionTypeAffiliatePayout TransactionType = "affiliate_payout"
	TransactionTypePlatformFee    TransactionType = "platform_fee"
	TransactionTypeAdjustment     TransactionType = "adjustment"
	TransactionTypeTransfer       TransactionType = "transfer"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusCancelled TransactionStatus = "cancelled"
	TransactionStatusRefunded  TransactionStatus = "refunded"
)

// Transaction represents a billing/payment transaction
type Transaction struct {
	TransactionID    int64  `json:"transaction_id" db:"transaction_id"`
	OrganizationID   int64  `json:"organization_id" db:"organization_id"`
	BillingAccountID int64  `json:"billing_account_id" db:"billing_account_id"`
	
	// Transaction Details
	Type     TransactionType `json:"type" db:"type"`
	Amount   decimal.Decimal `json:"amount" db:"amount"`
	Currency string          `json:"currency" db:"currency"`
	
	// Balance Tracking
	BalanceBefore decimal.Decimal `json:"balance_before" db:"balance_before"`
	BalanceAfter  decimal.Decimal `json:"balance_after" db:"balance_after"`
	
	// References
	ReferenceType          *string `json:"reference_type,omitempty" db:"reference_type"`
	ReferenceID            *string `json:"reference_id,omitempty" db:"reference_id"`
	RelatedTransactionID   *int64  `json:"related_transaction_id,omitempty" db:"related_transaction_id"`
	
	// Stripe Integration
	StripePaymentIntentID *string `json:"stripe_payment_intent_id,omitempty" db:"stripe_payment_intent_id"`
	StripeInvoiceID       *string `json:"stripe_invoice_id,omitempty" db:"stripe_invoice_id"`
	StripeChargeID        *string `json:"stripe_charge_id,omitempty" db:"stripe_charge_id"`
	
	// Metadata
	Description *string                `json:"description,omitempty" db:"description"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	
	// Status
	Status TransactionStatus `json:"status" db:"status"`
	
	// Timestamps
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// InvoiceStatus represents the status of an invoice
type InvoiceStatus string

const (
	InvoiceStatusDraft         InvoiceStatus = "draft"
	InvoiceStatusOpen          InvoiceStatus = "open"
	InvoiceStatusPaid          InvoiceStatus = "paid"
	InvoiceStatusVoid          InvoiceStatus = "void"
	InvoiceStatusUncollectible InvoiceStatus = "uncollectible"
	InvoiceStatusOverdue       InvoiceStatus = "overdue"
)

// Invoice represents an invoice for postpaid billing
type Invoice struct {
	InvoiceID        int64  `json:"invoice_id" db:"invoice_id"`
	OrganizationID   int64  `json:"organization_id" db:"organization_id"`
	BillingAccountID int64  `json:"billing_account_id" db:"billing_account_id"`
	
	// Invoice Details
	InvoiceNumber string `json:"invoice_number" db:"invoice_number"`
	
	// Stripe Integration
	StripeInvoiceID *string `json:"stripe_invoice_id,omitempty" db:"stripe_invoice_id"`
	
	// Financial Details
	Subtotal    decimal.Decimal `json:"subtotal" db:"subtotal"`
	TaxAmount   decimal.Decimal `json:"tax_amount" db:"tax_amount"`
	TotalAmount decimal.Decimal `json:"total_amount" db:"total_amount"`
	AmountPaid  decimal.Decimal `json:"amount_paid" db:"amount_paid"`
	AmountDue   decimal.Decimal `json:"amount_due" db:"amount_due"`
	Currency    string          `json:"currency" db:"currency"`
	
	// Billing Period
	PeriodStart time.Time `json:"period_start" db:"period_start"`
	PeriodEnd   time.Time `json:"period_end" db:"period_end"`
	
	// Dates
	InvoiceDate time.Time  `json:"invoice_date" db:"invoice_date"`
	DueDate     time.Time  `json:"due_date" db:"due_date"`
	PaidAt      *time.Time `json:"paid_at,omitempty" db:"paid_at"`
	
	// Status
	Status InvoiceStatus `json:"status" db:"status"`
	
	// Invoice Content
	LineItems []InvoiceLineItem      `json:"line_items,omitempty" db:"line_items"`
	Notes     *string                `json:"notes,omitempty" db:"notes"`
	
	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// InvoiceLineItem represents a line item on an invoice
type InvoiceLineItem struct {
	Description string          `json:"description"`
	Quantity    decimal.Decimal `json:"quantity"`
	UnitPrice   decimal.Decimal `json:"unit_price"`
	Amount      decimal.Decimal `json:"amount"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// UsageRecordStatus represents the status of a usage record
type UsageRecordStatus string

const (
	UsageRecordStatusPending    UsageRecordStatus = "pending"
	UsageRecordStatusCalculated UsageRecordStatus = "calculated"
	UsageRecordStatusBilled     UsageRecordStatus = "billed"
	UsageRecordStatusPaid       UsageRecordStatus = "paid"
	UsageRecordStatusFailed     UsageRecordStatus = "failed"
)

// UsageRecord represents daily usage calculation for an organization
type UsageRecord struct {
	UsageRecordID    int64  `json:"usage_record_id" db:"usage_record_id"`
	OrganizationID   int64  `json:"organization_id" db:"organization_id"`
	BillingAccountID int64  `json:"billing_account_id" db:"billing_account_id"`
	
	// Usage Period
	UsageDate time.Time `json:"usage_date" db:"usage_date"`
	
	// Usage Metrics
	Clicks      int `json:"clicks" db:"clicks"`
	Conversions int `json:"conversions" db:"conversions"`
	Impressions int `json:"impressions" db:"impressions"`
	
	// Financial Metrics
	AdvertiserSpend  decimal.Decimal `json:"advertiser_spend" db:"advertiser_spend"`
	AffiliatePayout  decimal.Decimal `json:"affiliate_payout" db:"affiliate_payout"`
	PlatformRevenue  decimal.Decimal `json:"platform_revenue" db:"platform_revenue"`
	Currency         string          `json:"currency" db:"currency"`
	
	// Processing Status
	Status UsageRecordStatus `json:"status" db:"status"`
	
	// Allocation Details
	AllocatedAt *time.Time `json:"allocated_at,omitempty" db:"allocated_at"`
	BilledAt    *time.Time `json:"billed_at,omitempty" db:"billed_at"`
	
	// Metadata
	CampaignBreakdown  map[string]interface{} `json:"campaign_breakdown,omitempty" db:"campaign_breakdown"`
	AffiliateBreakdown map[string]interface{} `json:"affiliate_breakdown,omitempty" db:"affiliate_breakdown"`
	Metadata           map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// WebhookEventStatus represents the status of a webhook event
type WebhookEventStatus string

const (
	WebhookEventStatusPending   WebhookEventStatus = "pending"
	WebhookEventStatusProcessed WebhookEventStatus = "processed"
	WebhookEventStatusFailed    WebhookEventStatus = "failed"
	WebhookEventStatusIgnored   WebhookEventStatus = "ignored"
)

// WebhookEvent represents a Stripe webhook event
type WebhookEvent struct {
	WebhookEventID int64  `json:"webhook_event_id" db:"webhook_event_id"`
	
	// Stripe Event Details
	StripeEventID string `json:"stripe_event_id" db:"stripe_event_id"`
	EventType     string `json:"event_type" db:"event_type"`
	
	// Processing Status
	Status WebhookEventStatus `json:"status" db:"status"`
	
	// Event Data
	EventData map[string]interface{} `json:"event_data" db:"event_data"`
	
	// Processing Details
	ProcessedAt  *time.Time `json:"processed_at,omitempty" db:"processed_at"`
	ErrorMessage *string    `json:"error_message,omitempty" db:"error_message"`
	RetryCount   int        `json:"retry_count" db:"retry_count"`
	
	// Related Records
	OrganizationID *int64 `json:"organization_id,omitempty" db:"organization_id"`
	TransactionID  *int64 `json:"transaction_id,omitempty" db:"transaction_id"`
	
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// #############################################################################
// ## Request/Response Models for API
// #############################################################################

// CreatePaymentMethodRequest represents a request to create a payment method
type CreatePaymentMethodRequest struct {
	PaymentMethodID string  `json:"payment_method_id" binding:"required"` // Stripe Payment Method ID
	SetAsDefault    bool    `json:"set_as_default"`
	Nickname        *string `json:"nickname,omitempty"`
}

// RechargeRequest represents a request to recharge an account
type RechargeRequest struct {
	Amount          decimal.Decimal `json:"amount" binding:"required"`
	Currency        string          `json:"currency"`
	PaymentMethodID *string         `json:"payment_method_id,omitempty"`
	Description     *string         `json:"description,omitempty"`
}

// UpdateBillingConfigRequest represents a request to update billing configuration
type UpdateBillingConfigRequest struct {
	BillingMode           *BillingMode    `json:"billing_mode,omitempty"`
	AutoRechargeEnabled   *bool           `json:"auto_recharge_enabled,omitempty"`
	AutoRechargeThreshold *decimal.Decimal `json:"auto_recharge_threshold,omitempty"`
	AutoRechargeAmount    *decimal.Decimal `json:"auto_recharge_amount,omitempty"`
	BillingEmail          *string         `json:"billing_email,omitempty"`
	InvoiceDayOfMonth     *int            `json:"invoice_day_of_month,omitempty"`
	PaymentTermsDays      *int            `json:"payment_terms_days,omitempty"`
}

// BillingDashboardResponse represents billing dashboard data
type BillingDashboardResponse struct {
	BillingAccount   *BillingAccount       `json:"billing_account"`
	PaymentMethods   []StripePaymentMethod `json:"payment_methods"`
	RecentTransactions []Transaction       `json:"recent_transactions"`
	CurrentBalance   decimal.Decimal       `json:"current_balance"`
	MonthlySpend     decimal.Decimal       `json:"monthly_spend"`
	PendingInvoices  []Invoice             `json:"pending_invoices,omitempty"`
}