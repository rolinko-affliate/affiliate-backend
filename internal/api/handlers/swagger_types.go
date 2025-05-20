package handlers

import (
	"encoding/json"
)

// This file contains type definitions to help Swagger generate proper documentation

// SwaggerRawMessage is used to represent JSON raw message in Swagger docs
// @Description JSON raw message
type SwaggerRawMessage struct {
	// Raw JSON content
	RawJSON interface{} `json:"raw_json"`
}

// Define aliases for json.RawMessage to help Swagger
type RawMessage json.RawMessage

// Define Swagger examples for common JSON fields

// BillingDetailsExample provides an example of billing details
// @Description Billing details in JSON format
// swagger:model
type BillingDetailsExample struct {
	// Billing address
	Address string `json:"address" example:"123 Main St, City, Country"`
	// Payment terms
	PaymentTerms string `json:"payment_terms" example:"Net 30"`
	// Tax ID
	TaxID string `json:"tax_id" example:"TAX-12345"`
}

// PaymentDetailsExample provides an example of payment details
// @Description Payment details in JSON format
// swagger:model
type PaymentDetailsExample struct {
	// Payment method
	Method string `json:"method" example:"bank_transfer"`
	// Bank account details
	BankAccount string `json:"bank_account" example:"IBAN: XX00 0000 0000 0000"`
	// Payment currency
	Currency string `json:"currency" example:"USD"`
}

// ProviderConfigExample provides an example of provider configuration
// @Description Provider configuration in JSON format
// swagger:model
type ProviderConfigExample struct {
	// API endpoint
	Endpoint string `json:"endpoint" example:"https://api.provider.com/v1"`
	// API version
	Version string `json:"version" example:"1.0"`
	// Additional settings
	Settings map[string]string `json:"settings"`
}

// ProviderOfferConfigExample provides an example of provider offer configuration
// @Description Provider offer configuration in JSON format
// swagger:model
type ProviderOfferConfigExample struct {
	// Offer ID in the provider system
	ProviderOfferID string `json:"provider_offer_id" example:"OFF-12345"`
	// Tracking URL template
	TrackingURLTemplate string `json:"tracking_url_template" example:"https://track.provider.com/{offer_id}/{affiliate_id}"`
	// Payout details
	Payout map[string]interface{} `json:"payout"`
}