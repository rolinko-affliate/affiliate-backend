package everflow

import (
	"encoding/json"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/everflow/affiliate"
)

// AffiliateProviderMapper handles mapping between domain models and Everflow API models
type AffiliateProviderMapper struct{}

// NewAffiliateProviderMapper creates a new affiliate provider mapper
func NewAffiliateProviderMapper() *AffiliateProviderMapper {
	return &AffiliateProviderMapper{}
}

// MapAffiliateToEverflowRequest creates an Everflow affiliate request from domain affiliate and provider mapping
func (m *AffiliateProviderMapper) MapAffiliateToEverflowRequest(
	aff *domain.Affiliate, 
	mapping *domain.AffiliateProviderMapping,
) (*affiliate.CreateAffiliateRequest, error) {
	// TODO: Implement full Everflow mapping when Everflow API is properly integrated
	return nil, fmt.Errorf("Everflow affiliate mapping not implemented")
}

// MapEverflowResponseToProviderData converts Everflow response to provider data
// Note: This now only maps Everflow-specific fields, as general purpose fields
// are handled separately and stored in the main affiliate model
func (m *AffiliateProviderMapper) MapEverflowResponseToProviderData(
	resp *affiliate.Affiliate,
) (*domain.EverflowProviderData, error) {
	// TODO: Implement full Everflow provider data mapping when Everflow API is properly integrated
	return &domain.EverflowProviderData{}, fmt.Errorf("Everflow provider data mapping not implemented")
}

// MapEverflowResponseToAffiliate converts Everflow response to main affiliate fields
// This handles the general purpose fields that were moved to the main affiliate model
func (m *AffiliateProviderMapper) MapEverflowResponseToAffiliate(
	resp *affiliate.Affiliate,
	aff *domain.Affiliate,
) error {
	// TODO: Implement full Everflow response mapping when Everflow API is properly integrated
	return fmt.Errorf("Everflow response mapping not implemented")
}

// Helper methods
func (m *AffiliateProviderMapper) mapDomainStatusToEverflowStatus(status string) string {
	switch status {
	case "active":
		return "active"
	case "pending":
		return "pending"
	case "rejected":
		return "rejected"
	case "inactive":
		return "inactive"
	default:
		return "pending"
	}
}

func (m *AffiliateProviderMapper) getDefaultNetworkEmployeeID(everflowData *domain.EverflowProviderData) int32 {
	if everflowData != nil && everflowData.NetworkEmployeeID != nil {
		return *everflowData.NetworkEmployeeID
	}
	return 1 // Default network employee ID
}

// MapEverflowResponseToProviderMapping updates provider mapping with Everflow response data
func (m *AffiliateProviderMapper) MapEverflowResponseToProviderMapping(
	resp interface{}, 
	mapping *domain.AffiliateProviderMapping,
) error {
	if resp == nil || mapping == nil {
		return fmt.Errorf("invalid response or mapping")
	}

	// TODO: Implement proper response parsing when affiliate sync functionality is needed
	// For now, just ensure provider data is properly initialized
	
	// Create or update provider data with Everflow-specific fields
	everflowData := &domain.EverflowProviderData{}
	
	// Unmarshal existing provider data if it exists
	if mapping.ProviderData != nil {
		if err := json.Unmarshal([]byte(*mapping.ProviderData), everflowData); err != nil {
			// If unmarshal fails, start with empty data
			everflowData = &domain.EverflowProviderData{}
		}
	}

	// TODO: Update Everflow-specific fields from response when needed

	// Marshal updated provider data
	providerDataBytes, err := json.Marshal(everflowData)
	if err != nil {
		return fmt.Errorf("error marshaling provider data: %w", err)
	}
	
	providerDataStr := string(providerDataBytes)
	mapping.ProviderData = &providerDataStr

	return nil
}