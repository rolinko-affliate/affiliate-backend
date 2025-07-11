package offer

type EntriesInfo []struct {
	NetworkOfferPayoutRevenueId    int     `json:"network_offer_payout_revenue_id"`
	NetworkId                      int     `json:"network_id"`
	NetworkOfferId                 int     `json:"network_offer_id"`
	EntryName                      string  `json:"entry_name"`
	PayoutType                     string  `json:"payout_type"`
	PayoutAmount                   float64 `json:"payout_amount"`
	PayoutPercentage               int     `json:"payout_percentage"`
	RevenueType                    string  `json:"revenue_type"`
	RevenueAmount                  int     `json:"revenue_amount"`
	RevenuePercentage              int     `json:"revenue_percentage"`
	IsDefault                      bool    `json:"is_default"`
	IsPrivate                      bool    `json:"is_private"`
	IsPostbackDisabled             bool    `json:"is_postback_disabled"`
	IsEnforceCaps                  bool    `json:"is_enforce_caps"`
	TimeCreated                    int     `json:"time_created"`
	GlobalAdvertiserEventId        int     `json:"global_advertiser_event_id"`
	IsMustApproveConversion        bool    `json:"is_must_approve_conversion"`
	IsAllowDuplicateConversion     bool    `json:"is_allow_duplicate_conversion"`
	IsEmailAttributionDefaultEvent bool    `json:"is_email_attribution_default_event"`
}
