package everflow

import (
	"testing"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdvertiserProviderMapper_MapAdvertiserToEverflowRequestWithContext(t *testing.T) {
	mapper := NewAdvertiserProviderMapper()

	// Create test advertiser
	advertiser := &domain.Advertiser{
		AdvertiserID: 12345,
		Name:         "Test Advertiser",
		Status:       "active",
	}

	// Create test context with organization and user
	userID := "user-123"
	ctx := &provider.AdvertiserMappingContext{
		Organization: &domain.Organization{
			OrganizationID: 67890,
			Name:           "Test Organization",
		},
		UserID: &userID,
	}

	// Test mapping with context
	req, err := mapper.MapAdvertiserToEverflowRequestWithContext(advertiser, nil, ctx)
	require.NoError(t, err)
	require.NotNil(t, req)

	// Verify labels contain organization ID and name
	labels := req.GetLabels()
	assert.Len(t, labels, 2)
	assert.Contains(t, labels, "Org ID: 67890")
	assert.Contains(t, labels, "Org Name: Test Organization")

	// Verify internal notes contain advertiser ID and user ID
	internalNotes := req.GetInternalNotes()
	assert.Contains(t, internalNotes, "Advertiser ID: 12345")
	assert.Contains(t, internalNotes, "User ID: user-123")

	// Verify billing frequency is set to manual
	if req.HasBilling() {
		billing := req.GetBilling()
		assert.Equal(t, "manual", billing.GetBillingFrequency())
	}
}

func TestAdvertiserProviderMapper_MapAdvertiserToEverflowRequestWithContext_NoContext(t *testing.T) {
	mapper := NewAdvertiserProviderMapper()

	// Create test advertiser
	advertiser := &domain.Advertiser{
		AdvertiserID: 12345,
		Name:         "Test Advertiser",
		Status:       "active",
	}

	// Test mapping without context
	req, err := mapper.MapAdvertiserToEverflowRequestWithContext(advertiser, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, req)

	// Verify default labels are used
	labels := req.GetLabels()
	assert.Len(t, labels, 1)
	assert.Contains(t, labels, "DTC Brand")

	// Verify default internal notes are used
	internalNotes := req.GetInternalNotes()
	assert.Equal(t, "Some notes not visible to the advertiser", internalNotes)

	// Verify billing frequency is still set to manual
	if req.HasBilling() {
		billing := req.GetBilling()
		assert.Equal(t, "manual", billing.GetBillingFrequency())
	}
}

func TestAdvertiserProviderMapper_MapAdvertiserToEverflowRequest_ManualBilling(t *testing.T) {
	mapper := NewAdvertiserProviderMapper()

	// Create test advertiser
	advertiser := &domain.Advertiser{
		AdvertiserID: 12345,
		Name:         "Test Advertiser",
		Status:       "active",
	}

	// Test mapping (backward compatibility method)
	req, err := mapper.MapAdvertiserToEverflowRequest(advertiser, nil)
	require.NoError(t, err)
	require.NotNil(t, req)

	// Verify billing frequency is set to manual
	if req.HasBilling() {
		billing := req.GetBilling()
		assert.Equal(t, "manual", billing.GetBillingFrequency())
	}
}