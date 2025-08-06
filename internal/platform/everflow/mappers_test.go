package everflow

import (
	"testing"

	"github.com/affiliate-backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestAffiliateMapping(t *testing.T) {
	t.Run("mapper methods handle nil inputs correctly", func(t *testing.T) {
		mapper := NewAffiliateProviderMapper()
		assert.NotNil(t, mapper)

		// Test that mapper methods exist and handle nil inputs properly
		_, err := mapper.MapAffiliateToEverflowRequest(nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "affiliate cannot be nil")

		err = mapper.MapEverflowResponseToAffiliate(nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "response and affiliate cannot be nil")

		_, err = mapper.MapEverflowResponseToProviderData(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "response cannot be nil")

		err = mapper.MapEverflowResponseToProviderMapping(nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid response or mapping")
	})
}

func TestStatusMapping(t *testing.T) {
	mapper := NewAffiliateProviderMapper()

	t.Run("mapDomainStatusToEverflowStatus", func(t *testing.T) {
		testCases := []struct {
			domainStatus   string
			expectedStatus string
		}{
			{"active", "active"},
			{"pending", "pending"},
			{"rejected", "rejected"},
			{"inactive", "inactive"},
			{"unknown", "pending"}, // Default case
		}

		for _, tc := range testCases {
			t.Run(tc.domainStatus, func(t *testing.T) {
				result := mapper.mapDomainStatusToEverflowStatus(tc.domainStatus)
				assert.Equal(t, tc.expectedStatus, result)
			})
		}
	})

	t.Run("getDefaultNetworkEmployeeID", func(t *testing.T) {
		t.Run("with provided data", func(t *testing.T) {
			everflowData := &domain.EverflowProviderData{
				NetworkEmployeeID: int32Ptr(5),
			}
			result := mapper.getDefaultNetworkEmployeeID(everflowData)
			assert.Equal(t, int32(5), result)
		})

		t.Run("with nil data", func(t *testing.T) {
			result := mapper.getDefaultNetworkEmployeeID(nil)
			assert.Equal(t, int32(1), result) // Default value
		})

		t.Run("with empty data", func(t *testing.T) {
			everflowData := &domain.EverflowProviderData{}
			result := mapper.getDefaultNetworkEmployeeID(everflowData)
			assert.Equal(t, int32(1), result) // Default value
		})
	})
}


