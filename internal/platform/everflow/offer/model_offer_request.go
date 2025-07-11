package offer

type OfferRequest struct {
	OfferStatus                   string              `json:"offer_status"`
	Visibility                    string              `json:"visibility"`
	RequirementTrackingParameters []string            `json:"requirement_tracking_parameters"`
	PayoutRevenue                 EntriesInfo         `json:"payout_revenue"`
	RedirectMode                  string              `json:"redirect_mode"`
	ConversionMethod              string              `json:"conversion_method"`
	Ruleset                       Ruleset             `json:"ruleset,omitempty"`
	EmailAttributionMethod        string              `json:"email_attribution_method"`
	AttributionMethod             string              `json:"attribution_method"`
	EmailOptout                   EmailOptoutSettings `json:"email_optout"`
	Labels                        Details             `json:"labels"`
	Channels                      Details             `json:"channels"`
	Integrations                  Integrations        `json:"integrations"`
	InternalRedirects             []string            `json:"internal_redirects"`
	TrafficFilters                []string            `json:"traffic_filters"`
	RequirementKpis               []string            `json:"requirement_kpis"`
	SourceNames                   []string            `json:"source_names"`
	Relationship                  RelationshipInfo    `json:"relationship"`
	Email                         []string            `json:"email"`

	Creatives                         []interface{} `json:"creatives"`
	IsSoftCap                         bool          `json:"is_soft_cap"`
	NetworkTrackingDomainId           int           `json:"network_tracking_domain_id"`
	IsUseSecureLink                   bool          `json:"is_use_secure_link"`
	CurrencyId                        string        `json:"currency_id"`
	SessionDuration                   int           `json:"session_duration"`
	SessionDefinition                 string        `json:"session_definition"`
	NetworkCategoryId                 int           `json:"network_category_id"`
	NetworkAdvertiserId               int           `json:"network_advertiser_id"`
	Name                              string        `json:"name"`
	DestinationUrl                    string        `json:"destination_url"`
	AppIdentifier                     string        `json:"app_identifier"`
	PreviewUrl                        string        `json:"preview_url"`
	InternalNotes                     string        `json:"internal_notes"`
	DateLiveUntil                     string        `json:"date_live_until"`
	HtmlDescription                   string        `json:"html_description"`
	IsDescriptionPlainText            bool          `json:"is_description_plain_text"`
	IsUseDirectLinking                bool          `json:"is_use_direct_linking"`
	IsAllowDeepLink                   bool          `json:"is_allow_deep_link"`
	IsUsingExplicitTermsAndConditions bool          `json:"is_using_explicit_terms_and_conditions"`
	IsForceTermsAndConditions         bool          `json:"is_force_terms_and_conditions"`
	TermsAndConditions                string        `json:"terms_and_conditions"`
	NetworkOfferGroupId               int           `json:"network_offer_group_id"`
	NetworkApplicationQuestionnaireId int           `json:"network_application_questionnaire_id"`
	CapsTimezoneId                    int           `json:"caps_timezone_id"`
	ServerSideUrl                     string        `json:"server_side_url"`
	IsEmailAttributionWindowEnabled   bool          `json:"is_email_attribution_window_enabled"`
	SuppressionListId                 int           `json:"suppression_list_id"`
	ThumbnailUrl                      string        `json:"thumbnail_url"`
}
