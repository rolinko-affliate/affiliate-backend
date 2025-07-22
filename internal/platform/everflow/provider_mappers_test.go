package everflow

import (
	"testing"

	"github.com/affiliate-backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestAffiliateProviderMapper(t *testing.T) {
	mapper := NewAffiliateProviderMapper()

	t.Run("MapAffiliateToEverflowRequest", func(t *testing.T) {
		t.Run("returns not implemented error", func(t *testing.T) {
			// Test that the stubbed method returns the expected error
			_, err := mapper.MapAffiliateToEverflowRequest(nil, nil)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "not implemented")
		})
	})

	t.Run("MapEverflowResponseToProviderMapping", func(t *testing.T) {
		t.Run("handles nil inputs", func(t *testing.T) {
			err := mapper.MapEverflowResponseToProviderMapping(nil, nil)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid response or mapping")
		})

		t.Run("initializes provider data", func(t *testing.T) {
			mapping := &domain.AffiliateProviderMapping{
				AffiliateID:  1,
				ProviderType: "everflow",
			}

			err := mapper.MapEverflowResponseToProviderMapping(map[string]interface{}{}, mapping)
			assert.NoError(t, err)
			assert.NotNil(t, mapping.ProviderData)
		})
	})

	t.Run("MapEverflowResponseToAffiliate", func(t *testing.T) {
		t.Run("returns not implemented error", func(t *testing.T) {
			err := mapper.MapEverflowResponseToAffiliate(nil, nil)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "not implemented")
		})
	})

	t.Run("MapEverflowResponseToProviderData", func(t *testing.T) {
		t.Run("returns not implemented error", func(t *testing.T) {
			_, err := mapper.MapEverflowResponseToProviderData(nil)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "not implemented")
		})
	})
}
