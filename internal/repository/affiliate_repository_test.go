package repository

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestAffiliateJSONBSerialization(t *testing.T) {
	// Test that our domain models can be properly marshaled/unmarshaled
	// This is a unit test that doesn't require a database connection

	// Create test contact address and billing info structures
	contactAddress := &domain.ContactAddress{
		Address1:      stringPtr("123 Main St"),
		City:          stringPtr("New York"),
		RegionCode:    stringPtr("NY"),
		CountryCode:   stringPtr("US"),
		ZipPostalCode: stringPtr("10001"),
	}

	billingInfo := &domain.BillingDetails{
		Frequency:   billingFrequencyPtr(domain.BillingFrequencyMonthly),
		PaymentType: paymentTypePtr(domain.PaymentTypeWire),
		Schedule: &domain.BillingSchedule{
			DayOfMonth: int32Ptr(15),
		},
		PaymentDetails: &domain.PaymentDetails{
			Type:          paymentDetailsTypePtr(domain.PaymentDetailsTypeWire),
			BankName:      stringPtr("Test Bank"),
			AccountNumber: stringPtr("123456789"),
			RoutingNumber: stringPtr("987654321"),
		},
	}

	// Convert structs to JSON strings as expected by the domain model
	contactAddressJSON, err := json.Marshal(contactAddress)
	assert.NoError(t, err)

	billingInfoJSON, err := json.Marshal(billingInfo)
	assert.NoError(t, err)

	// Create test affiliate with JSON string fields
	affiliate := &domain.Affiliate{
		AffiliateID:    1,
		OrganizationID: 1,
		Name:           "Test Affiliate",
		Status:         "active",

		// Contact address as JSON string
		ContactAddress: stringPtr(string(contactAddressJSON)),

		// Billing info as JSON string
		BillingInfo: stringPtr(string(billingInfoJSON)),

		// Labels as JSON string
		Labels: stringPtr(`["premium", "tier1"]`),

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test that JSON strings were created properly
	assert.NotNil(t, affiliate.ContactAddress)
	assert.Contains(t, *affiliate.ContactAddress, "123 Main St")
	assert.Contains(t, *affiliate.ContactAddress, "New York")

	assert.NotNil(t, affiliate.BillingInfo)
	assert.Contains(t, *affiliate.BillingInfo, "Test Bank")
	assert.Contains(t, *affiliate.BillingInfo, "monthly")

	// Test ContactAddress struct methods directly
	assert.True(t, contactAddress.HasData())

	// Test BillingDetails struct methods directly
	assert.True(t, billingInfo.HasData())

	// Test PaymentDetails validation directly
	err = billingInfo.PaymentDetails.Validate()
	assert.NoError(t, err)

	// Test PaymentDetails HasData method directly
	assert.True(t, billingInfo.PaymentDetails.HasData())

	// Test that empty structures return false for HasData
	emptyContactAddress := &domain.ContactAddress{}
	assert.False(t, emptyContactAddress.HasData())

	emptyBillingDetails := &domain.BillingDetails{}
	assert.False(t, emptyBillingDetails.HasData())

	emptyPaymentDetails := &domain.PaymentDetails{}
	assert.False(t, emptyPaymentDetails.HasData())
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func int32Ptr(i int32) *int32 {
	return &i
}

func billingFrequencyPtr(bf domain.BillingFrequency) *domain.BillingFrequency {
	return &bf
}

func paymentTypePtr(pt domain.PaymentType) *domain.PaymentType {
	return &pt
}

func paymentDetailsTypePtr(pdt domain.PaymentDetailsType) *domain.PaymentDetailsType {
	return &pdt
}
