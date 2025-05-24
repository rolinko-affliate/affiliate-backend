package domain

// BillingDetails represents structured billing information for an advertiser
type BillingDetails struct {
	BillingFrequency           string                 `json:"billing_frequency"` // "weekly", "bimonthly", "monthly", "two_months", "quarterly", "manual", "other"
	TaxID                      *string                `json:"tax_id,omitempty"`
	IsInvoiceCreationAuto      *bool                  `json:"is_invoice_creation_auto,omitempty"`
	InvoiceAmountThreshold     *float64               `json:"invoice_amount_threshold,omitempty"`
	AutoInvoiceStartDate       *string                `json:"auto_invoice_start_date,omitempty"` // Format: "2019-06-01 00:00:00"
	DefaultInvoiceIsHidden     *bool                  `json:"default_invoice_is_hidden,omitempty"`
	InvoiceGenerationDaysDelay *int                   `json:"invoice_generation_days_delay,omitempty"`
	DefaultPaymentTerms        *int                   `json:"default_payment_terms,omitempty"`
	Address                    *BillingAddress        `json:"address,omitempty"`
	PaymentMethod              *PaymentMethod         `json:"payment_method,omitempty"`
	AdditionalDetails          map[string]interface{} `json:"additional_details,omitempty"`
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