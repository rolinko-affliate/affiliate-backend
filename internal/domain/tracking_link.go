package domain

import (
	"encoding/json"
	"time"
)

// TrackingLink represents a clean tracking link entity following clean architecture principles
type TrackingLink struct {
	TrackingLinkID int64   `json:"tracking_link_id" db:"tracking_link_id"`
	OrganizationID int64   `json:"organization_id" db:"organization_id"`
	CampaignID     int64   `json:"campaign_id" db:"campaign_id"`
	AffiliateID    int64   `json:"affiliate_id" db:"affiliate_id"`
	Name           string  `json:"name" db:"name"`
	Description    *string `json:"description,omitempty" db:"description"`
	Status         string  `json:"status" db:"status"` // 'active', 'paused', 'archived'

	// Core tracking link fields (provider-agnostic)
	TrackingURL *string `json:"tracking_url,omitempty" db:"tracking_url"`
	SourceID    *string `json:"source_id,omitempty" db:"source_id"`
	Sub1        *string `json:"sub1,omitempty" db:"sub1"`
	Sub2        *string `json:"sub2,omitempty" db:"sub2"`
	Sub3        *string `json:"sub3,omitempty" db:"sub3"`
	Sub4        *string `json:"sub4,omitempty" db:"sub4"`
	Sub5        *string `json:"sub5,omitempty" db:"sub5"`

	// Link configuration
	IsEncryptParameters *bool `json:"is_encrypt_parameters,omitempty" db:"is_encrypt_parameters"`
	IsRedirectLink      *bool `json:"is_redirect_link,omitempty" db:"is_redirect_link"`

	// General purpose fields
	InternalNotes *string `json:"internal_notes,omitempty" db:"internal_notes"`
	Tags          *string `json:"tags,omitempty" db:"tags"` // JSONB stored as string (array of strings)

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TrackingLinkProviderMapping represents a mapping between a tracking link and a provider
type TrackingLinkProviderMapping struct {
	MappingID              int64   `json:"mapping_id" db:"mapping_id"`
	TrackingLinkID         int64   `json:"tracking_link_id" db:"tracking_link_id"`
	ProviderType           string  `json:"provider_type" db:"provider_type"` // 'everflow' for MVP
	ProviderTrackingLinkID *string `json:"provider_tracking_link_id,omitempty" db:"provider_tracking_link_id"`

	// Provider-specific data stored as JSONB
	ProviderData *string `json:"provider_data,omitempty" db:"provider_data"`

	// Synchronization metadata
	SyncStatus *string    `json:"sync_status,omitempty" db:"sync_status"`
	LastSyncAt *time.Time `json:"last_sync_at,omitempty" db:"last_sync_at"`
	SyncError  *string    `json:"sync_error,omitempty" db:"sync_error"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// EverflowTrackingLinkProviderData represents Everflow-specific tracking link data
type EverflowTrackingLinkProviderData struct {
	NetworkOfferID          *int32 `json:"network_offer_id,omitempty"`
	NetworkCampaignID       *int32 `json:"network_campaign_id,omitempty"`
	NetworkAffiliateID      *int32 `json:"network_affiliate_id,omitempty"`
	NetworkTrackingDomainID *int32 `json:"network_tracking_domain_id,omitempty"`
	NetworkOfferURLID       *int32 `json:"network_offer_url_id,omitempty"`
	CreativeID              *int32 `json:"creative_id,omitempty"`
	NetworkTrafficSourceID  *int32 `json:"network_traffic_source_id,omitempty"`

	// Generated tracking URL from Everflow
	GeneratedURL *string `json:"generated_url,omitempty"`

	// Additional Everflow-specific fields
	CanAffiliateRunAllOffers *bool `json:"can_affiliate_run_all_offers,omitempty"`

	// Additional fields for extensibility
	AdditionalFields map[string]interface{} `json:"additional_fields,omitempty"`
}

// ToJSON converts EverflowTrackingLinkProviderData to JSON string
func (e *EverflowTrackingLinkProviderData) ToJSON() (string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON populates EverflowTrackingLinkProviderData from JSON string
func (e *EverflowTrackingLinkProviderData) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), e)
}

// TrackingLinkGenerationRequest represents a request to generate a tracking link
type TrackingLinkGenerationRequest struct {
	CampaignID  int64   `json:"campaign_id"`
	AffiliateID int64   `json:"affiliate_id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`

	// Optional tracking parameters
	SourceID *string `json:"source_id,omitempty"`
	Sub1     *string `json:"sub1,omitempty"`
	Sub2     *string `json:"sub2,omitempty"`
	Sub3     *string `json:"sub3,omitempty"`
	Sub4     *string `json:"sub4,omitempty"`
	Sub5     *string `json:"sub5,omitempty"`

	// Link configuration
	IsEncryptParameters *bool `json:"is_encrypt_parameters,omitempty"`
	IsRedirectLink      *bool `json:"is_redirect_link,omitempty"`

	// Provider-specific options
	NetworkTrackingDomainID *int32 `json:"network_tracking_domain_id,omitempty"`
	NetworkOfferURLID       *int32 `json:"network_offer_url_id,omitempty"`
	CreativeID              *int32 `json:"creative_id,omitempty"`
	NetworkTrafficSourceID  *int32 `json:"network_traffic_source_id,omitempty"`
}

// TrackingLinkGenerationResponse represents the response from generating a tracking link
type TrackingLinkGenerationResponse struct {
	TrackingLink *TrackingLink `json:"tracking_link"`
	GeneratedURL string        `json:"generated_url"`
	ProviderData *string       `json:"provider_data,omitempty"`
}

type AddTrackingLink struct {
	NetworkOfferId          int    `json:"network_offer_id"`
	NetworkAffiliateId      int    `json:"network_affiliate_id"`
	NetworkTrackingDomainId int    `json:"network_tracking_domain_id"`
	NetworkOfferUrlId       int    `json:"network_offer_url_id"`
	CreativeId              int    `json:"creative_id"`
	NetworkTrafficSourceId  int    `json:"network_traffic_source_id"`
	SourceId                string `json:"source_id"`
	Sub1                    string `json:"sub1"`
	Sub2                    string `json:"sub2"`
	Sub3                    string `json:"sub3"`
	Sub4                    string `json:"sub4"`
	Sub5                    string `json:"sub5"`
	IsEncryptParameters     bool   `json:"is_encrypt_parameters"`
	IsRedirectLink          bool   `json:"is_redirect_link"`
}
