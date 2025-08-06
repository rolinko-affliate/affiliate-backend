package domain

import (
	"testing"
)

func TestOrganizationType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		orgType  OrganizationType
		expected bool
	}{
		{
			name:     "advertiser is valid",
			orgType:  OrganizationTypeAdvertiser,
			expected: true,
		},
		{
			name:     "affiliate is valid",
			orgType:  OrganizationTypeAffiliate,
			expected: true,
		},
		{
			name:     "platform_owner is valid",
			orgType:  OrganizationTypePlatformOwner,
			expected: true,
		},
		{
			name:     "agency is valid",
			orgType:  OrganizationTypeAgency,
			expected: true,
		},
		{
			name:     "invalid type is not valid",
			orgType:  OrganizationType("invalid"),
			expected: false,
		},
		{
			name:     "empty type is not valid",
			orgType:  OrganizationType(""),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.orgType.IsValid()
			if result != tt.expected {
				t.Errorf("OrganizationType.IsValid() = %v, expected %v for type %s", result, tt.expected, tt.orgType)
			}
		})
	}
}

func TestGetValidOrganizationTypes(t *testing.T) {
	validTypes := GetValidOrganizationTypes()
	
	expectedTypes := []OrganizationType{
		OrganizationTypeAdvertiser,
		OrganizationTypeAffiliate,
		OrganizationTypePlatformOwner,
		OrganizationTypeAgency,
	}

	if len(validTypes) != len(expectedTypes) {
		t.Errorf("GetValidOrganizationTypes() returned %d types, expected %d", len(validTypes), len(expectedTypes))
	}

	// Check that all expected types are present
	typeMap := make(map[OrganizationType]bool)
	for _, orgType := range validTypes {
		typeMap[orgType] = true
	}

	for _, expectedType := range expectedTypes {
		if !typeMap[expectedType] {
			t.Errorf("GetValidOrganizationTypes() missing expected type: %s", expectedType)
		}
	}

	// Verify all returned types are valid
	for _, orgType := range validTypes {
		if !orgType.IsValid() {
			t.Errorf("GetValidOrganizationTypes() returned invalid type: %s", orgType)
		}
	}
}

func TestOrganizationType_String(t *testing.T) {
	tests := []struct {
		name     string
		orgType  OrganizationType
		expected string
	}{
		{
			name:     "advertiser string representation",
			orgType:  OrganizationTypeAdvertiser,
			expected: "advertiser",
		},
		{
			name:     "affiliate string representation",
			orgType:  OrganizationTypeAffiliate,
			expected: "affiliate",
		},
		{
			name:     "platform_owner string representation",
			orgType:  OrganizationTypePlatformOwner,
			expected: "platform_owner",
		},
		{
			name:     "agency string representation",
			orgType:  OrganizationTypeAgency,
			expected: "agency",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.orgType.String()
			if result != tt.expected {
				t.Errorf("OrganizationType.String() = %v, expected %v", result, tt.expected)
			}
		})
	}
}