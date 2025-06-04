package domain

import "fmt"

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