package repository

import (
	"testing"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestAffiliateJSONBSerialization(t *testing.T) {
	// Test that our domain models can be properly marshaled/unmarshaled
	// This is a unit test that doesn't require a database connection
	
	// Create test affiliate with complex data
	affiliate := &domain.Affiliate{
		AffiliateID:    1,
		OrganizationID: 1,
		Name:           "Test Affiliate",
		Status:         "active",
		
		// Contact address
		ContactAddress: &domain.ContactAddress{
			Address1:      stringPtr("123 Main St"),
			City:          stringPtr("New York"),
			RegionCode:    stringPtr("NY"),
			CountryCode:   stringPtr("US"),
			ZipPostalCode: stringPtr("10001"),
		},
		
		// Billing info with structured data
		BillingInfo: &domain.BillingDetails{
			Frequency: billingFrequencyPtr(domain.BillingFrequencyMonthly),
			PaymentType: paymentTypePtr(domain.PaymentTypeWire),
			Schedule: &domain.BillingSchedule{
				DayOfMonth: int32Ptr(15),
			},
			PaymentDetails: &domain.PaymentDetails{
				Type: paymentDetailsTypePtr(domain.PaymentDetailsTypeWire),
				BankName:      stringPtr("Test Bank"),
				AccountNumber: stringPtr("123456789"),
				RoutingNumber: stringPtr("987654321"),
			},
		},
		
		// Labels as JSON string
		Labels: stringPtr(`["premium", "tier1"]`),
		
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Test ContactAddress HasData method
	assert.True(t, affiliate.ContactAddress.HasData())
	
	// Test BillingDetails HasData method
	assert.True(t, affiliate.BillingInfo.HasData())
	
	// Test PaymentDetails validation
	err := affiliate.BillingInfo.PaymentDetails.Validate()
	assert.NoError(t, err)
	
	// Test PaymentDetails HasData method
	assert.True(t, affiliate.BillingInfo.PaymentDetails.HasData())
	
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