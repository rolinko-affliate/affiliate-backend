package offer

type RelationshipInfo struct {
	Category                      CategoryInfo `json:"category,omitempty"`
	Labels                        Details      `json:"labels"`
	PayoutRevenue                 Details      `json:"payout_revenue"`
	EncodedValue                  string       `json:"encoded_value"`
	IsLockedCurrency              bool         `json:"is_locked_currency"`
	Channels                      Details      `json:"channels"`
	IsLockedCapsTimezone          bool         `json:"is_locked_caps_timezone"`
	RequirementKpis               Details      `json:"requirement_kpis"`
	RequirementTrackingParameters Details      `json:"requirement_tracking_parameters"`
}
