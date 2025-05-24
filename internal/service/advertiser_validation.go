package service

import (
	"fmt"

	"github.com/affiliate-backend/internal/domain"
)

var (
	validStatuses = map[string]bool{
		"active":   true,
		"pending":  true,
		"inactive": true,
		"rejected": true,
	}

	validBillingFrequencies = map[string]bool{
		"weekly":     true,
		"bimonthly":  true,
		"monthly":    true,
		"two_months": true,
		"quarterly":  true,
		"manual":     true,
		"other":      true,
	}
)

func validateAdvertiserBasics(advertiser *domain.Advertiser) error {
	if advertiser.Name == "" {
		return fmt.Errorf("advertiser name cannot be empty")
	}

	if !validStatuses[advertiser.Status] {
		return fmt.Errorf("invalid status: %s", advertiser.Status)
	}

	return nil
}

func validateBillingDetails(billing *domain.BillingDetails) error {
	if billing == nil {
		return nil
	}

	if billing.BillingFrequency != "" && !validBillingFrequencies[billing.BillingFrequency] {
		return fmt.Errorf("invalid billing frequency: %s", billing.BillingFrequency)
	}

	return nil
}

func validateAdvertiser(advertiser *domain.Advertiser) error {
	if err := validateAdvertiserBasics(advertiser); err != nil {
		return err
	}

	return validateBillingDetails(advertiser.BillingDetails)
}

func setDefaultStatus(advertiser *domain.Advertiser) {
	if advertiser.Status == "" {
		advertiser.Status = "pending"
	}
}