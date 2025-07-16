package models

import (
	"fmt"
	"time"

	"github.com/affiliate-backend/internal/domain"
)

// TrackingLinkRequest represents the request to create a tracking link
type TrackingLinkRequest struct {
	Name                string  `json:"name" binding:"required" example:"Facebook Campaign Link"`
	Description         *string `json:"description,omitempty" example:"Tracking link for Facebook traffic"`
	Status              string  `json:"status" binding:"required" example:"active"`
	CampaignID          int64   `json:"campaign_id" binding:"required" example:"1"`
	AffiliateID         int64   `json:"affiliate_id" binding:"required" example:"1"`
	SourceID            *string `json:"source_id,omitempty" example:"facebook"`
	Sub1                *string `json:"sub1,omitempty" example:"campaign_123"`
	Sub2                *string `json:"sub2,omitempty" example:"adset_456"`
	Sub3                *string `json:"sub3,omitempty" example:"ad_789"`
	Sub4                *string `json:"sub4,omitempty" example:"placement_mobile"`
	Sub5                *string `json:"sub5,omitempty" example:"audience_lookalike"`
	IsEncryptParameters *bool   `json:"is_encrypt_parameters,omitempty" example:"false"`
	IsRedirectLink      *bool   `json:"is_redirect_link,omitempty" example:"true"`
	InternalNotes       *string `json:"internal_notes,omitempty" example:"High-performing traffic source"`
	Tags                *string `json:"tags,omitempty" example:"facebook,mobile,lookalike"`
}

// TrackingLinkUpdateRequest represents the request to update a tracking link
type TrackingLinkUpdateRequest struct {
	Name                *string `json:"name,omitempty" example:"Facebook Campaign Link"`
	Description         *string `json:"description,omitempty" example:"Tracking link for Facebook traffic"`
	Status              *string `json:"status,omitempty" example:"active"`
	SourceID            *string `json:"source_id,omitempty" example:"facebook"`
	Sub1                *string `json:"sub1,omitempty" example:"campaign_123"`
	Sub2                *string `json:"sub2,omitempty" example:"adset_456"`
	Sub3                *string `json:"sub3,omitempty" example:"ad_789"`
	Sub4                *string `json:"sub4,omitempty" example:"placement_mobile"`
	Sub5                *string `json:"sub5,omitempty" example:"audience_lookalike"`
	IsEncryptParameters *bool   `json:"is_encrypt_parameters,omitempty" example:"false"`
	IsRedirectLink      *bool   `json:"is_redirect_link,omitempty" example:"true"`
	InternalNotes       *string `json:"internal_notes,omitempty" example:"High-performing traffic source"`
	Tags                *string `json:"tags,omitempty" example:"facebook,mobile,lookalike"`
}

// TrackingLinkGenerationRequest represents the request to generate a tracking link
type TrackingLinkGenerationRequest struct {
	Name                    string  `json:"name" binding:"required" example:"Facebook Campaign Link"`
	Description             *string `json:"description,omitempty" example:"Tracking link for Facebook traffic"`
	CampaignID              int64   `json:"campaign_id" binding:"required" example:"1"`
	AffiliateID             int64   `json:"affiliate_id" binding:"required" example:"1"`
	SourceID                *string `json:"source_id,omitempty" example:"facebook"`
	Sub1                    *string `json:"sub1,omitempty" example:"campaign_123"`
	Sub2                    *string `json:"sub2,omitempty" example:"adset_456"`
	Sub3                    *string `json:"sub3,omitempty" example:"ad_789"`
	Sub4                    *string `json:"sub4,omitempty" example:"placement_mobile"`
	Sub5                    *string `json:"sub5,omitempty" example:"audience_lookalike"`
	IsEncryptParameters     *bool   `json:"is_encrypt_parameters,omitempty" example:"false"`
	IsRedirectLink          *bool   `json:"is_redirect_link,omitempty" example:"true"`
	NetworkTrackingDomainID *int32  `json:"network_tracking_domain_id,omitempty" example:"1"`
	NetworkOfferURLID       *int32  `json:"network_offer_url_id,omitempty" example:"1"`
	CreativeID              *int32  `json:"creative_id,omitempty" example:"1"`
	NetworkTrafficSourceID  *int32  `json:"network_traffic_source_id,omitempty" example:"1"`
	InternalNotes           *string `json:"internal_notes,omitempty" example:"High-performing traffic source"`
	Tags                    *string `json:"tags,omitempty" example:"facebook,mobile,lookalike"`
}

// TrackingLinkResponse represents the response for tracking link operations
type TrackingLinkResponse struct {
	TrackingLinkID      int64   `json:"tracking_link_id" example:"1"`
	OrganizationID      int64   `json:"organization_id" example:"1"`
	CampaignID          int64   `json:"campaign_id" example:"1"`
	AffiliateID         int64   `json:"affiliate_id" example:"1"`
	Name                string  `json:"name" example:"Facebook Campaign Link"`
	Description         *string `json:"description,omitempty" example:"Tracking link for Facebook traffic"`
	Status              string  `json:"status" example:"active"`
	TrackingURL         *string `json:"tracking_url,omitempty" example:"https://tracking.example.com/ABC123/DEF456/?sub1=campaign_123"`
	SourceID            *string `json:"source_id,omitempty" example:"facebook"`
	Sub1                *string `json:"sub1,omitempty" example:"campaign_123"`
	Sub2                *string `json:"sub2,omitempty" example:"adset_456"`
	Sub3                *string `json:"sub3,omitempty" example:"ad_789"`
	Sub4                *string `json:"sub4,omitempty" example:"placement_mobile"`
	Sub5                *string `json:"sub5,omitempty" example:"audience_lookalike"`
	IsEncryptParameters *bool   `json:"is_encrypt_parameters,omitempty" example:"false"`
	IsRedirectLink      *bool   `json:"is_redirect_link,omitempty" example:"true"`

	InternalNotes *string   `json:"internal_notes,omitempty" example:"High-performing traffic source"`
	Tags          *string   `json:"tags,omitempty" example:"facebook,mobile,lookalike"`
	CreatedAt     time.Time `json:"created_at" example:"2023-12-01T10:00:00Z"`
	UpdatedAt     time.Time `json:"updated_at" example:"2023-12-01T10:30:00Z"`
}

// TrackingLinkGenerationResponse represents the response for tracking link generation
type TrackingLinkGenerationResponse struct {
	TrackingLink *TrackingLinkResponse `json:"tracking_link"`
	GeneratedURL string                `json:"generated_url" example:"https://tracking.example.com/ABC123/DEF456/?sub1=campaign_123"`
	QRCodeURL    *string               `json:"qr_code_url,omitempty" example:"https://api.example.com/tracking-links/1/qr"`
}

// TrackingLinkListResponse represents the response for listing tracking links
type TrackingLinkListResponse struct {
	TrackingLinks []*TrackingLinkResponse `json:"tracking_links"`
	Total         int                     `json:"total" example:"150"`
	Page          int                     `json:"page" example:"1"`
	PageSize      int                     `json:"page_size" example:"20"`
	TotalPages    int                     `json:"total_pages" example:"8"`
}

// TrackingLinkProviderMappingResponse represents the response for tracking link provider mappings
type TrackingLinkProviderMappingResponse struct {
	MappingID              int64      `json:"mapping_id" example:"1"`
	TrackingLinkID         int64      `json:"tracking_link_id" example:"1"`
	ProviderType           string     `json:"provider_type" example:"everflow"`
	ProviderTrackingLinkID *string    `json:"provider_tracking_link_id,omitempty" example:"abc123"`
	ProviderData           *string    `json:"provider_data,omitempty"`
	SyncStatus             *string    `json:"sync_status,omitempty" example:"synced"`
	LastSyncAt             *time.Time `json:"last_sync_at,omitempty" example:"2023-12-01T10:30:00Z"`
	SyncError              *string    `json:"sync_error,omitempty"`
	CreatedAt              time.Time  `json:"created_at" example:"2023-12-01T10:00:00Z"`
	UpdatedAt              time.Time  `json:"updated_at" example:"2023-12-01T10:30:00Z"`
}

// Conversion functions

// ToTrackingLinkDomain converts API request to domain model
func (req *TrackingLinkRequest) ToTrackingLinkDomain(organizationID int64) *domain.TrackingLink {
	return &domain.TrackingLink{
		OrganizationID:      organizationID,
		CampaignID:          req.CampaignID,
		AffiliateID:         req.AffiliateID,
		Name:                req.Name,
		Description:         req.Description,
		Status:              req.Status,
		SourceID:            req.SourceID,
		Sub1:                req.Sub1,
		Sub2:                req.Sub2,
		Sub3:                req.Sub3,
		Sub4:                req.Sub4,
		Sub5:                req.Sub5,
		IsEncryptParameters: req.IsEncryptParameters,
		IsRedirectLink:      req.IsRedirectLink,
		InternalNotes:       req.InternalNotes,
		Tags:                req.Tags,
	}
}

// UpdateTrackingLinkDomain updates a domain model with non-nil fields from the update request
func (req *TrackingLinkUpdateRequest) UpdateTrackingLinkDomain(trackingLink *domain.TrackingLink) {
	if req.Name != nil {
		trackingLink.Name = *req.Name
	}
	if req.Description != nil {
		trackingLink.Description = req.Description
	}
	if req.Status != nil {
		trackingLink.Status = *req.Status
	}
	if req.SourceID != nil {
		trackingLink.SourceID = req.SourceID
	}
	if req.Sub1 != nil {
		trackingLink.Sub1 = req.Sub1
	}
	if req.Sub2 != nil {
		trackingLink.Sub2 = req.Sub2
	}
	if req.Sub3 != nil {
		trackingLink.Sub3 = req.Sub3
	}
	if req.Sub4 != nil {
		trackingLink.Sub4 = req.Sub4
	}
	if req.Sub5 != nil {
		trackingLink.Sub5 = req.Sub5
	}
	if req.IsEncryptParameters != nil {
		trackingLink.IsEncryptParameters = req.IsEncryptParameters
	}
	if req.IsRedirectLink != nil {
		trackingLink.IsRedirectLink = req.IsRedirectLink
	}
	if req.InternalNotes != nil {
		trackingLink.InternalNotes = req.InternalNotes
	}
	if req.Tags != nil {
		trackingLink.Tags = req.Tags
	}
}

// ToTrackingLinkGenerationDomain converts API generation request to domain model
func (req *TrackingLinkGenerationRequest) ToTrackingLinkGenerationDomain() *domain.TrackingLinkGenerationRequest {
	return &domain.TrackingLinkGenerationRequest{
		Name:                    req.Name,
		Description:             req.Description,
		CampaignID:              req.CampaignID,
		AffiliateID:             req.AffiliateID,
		SourceID:                req.SourceID,
		Sub1:                    req.Sub1,
		Sub2:                    req.Sub2,
		Sub3:                    req.Sub3,
		Sub4:                    req.Sub4,
		Sub5:                    req.Sub5,
		IsEncryptParameters:     req.IsEncryptParameters,
		IsRedirectLink:          req.IsRedirectLink,
		NetworkTrackingDomainID: req.NetworkTrackingDomainID,
		NetworkOfferURLID:       req.NetworkOfferURLID,
		CreativeID:              req.CreativeID,
		NetworkTrafficSourceID:  req.NetworkTrafficSourceID,
	}
}

// FromTrackingLinkDomain converts domain model to API response
func FromTrackingLinkDomain(trackingLink *domain.TrackingLink) *TrackingLinkResponse {
	return &TrackingLinkResponse{
		TrackingLinkID:      trackingLink.TrackingLinkID,
		OrganizationID:      trackingLink.OrganizationID,
		CampaignID:          trackingLink.CampaignID,
		AffiliateID:         trackingLink.AffiliateID,
		Name:                trackingLink.Name,
		Description:         trackingLink.Description,
		Status:              trackingLink.Status,
		TrackingURL:         trackingLink.TrackingURL,
		SourceID:            trackingLink.SourceID,
		Sub1:                trackingLink.Sub1,
		Sub2:                trackingLink.Sub2,
		Sub3:                trackingLink.Sub3,
		Sub4:                trackingLink.Sub4,
		Sub5:                trackingLink.Sub5,
		IsEncryptParameters: trackingLink.IsEncryptParameters,
		IsRedirectLink:      trackingLink.IsRedirectLink,

		InternalNotes: trackingLink.InternalNotes,
		Tags:          trackingLink.Tags,
		CreatedAt:     trackingLink.CreatedAt,
		UpdatedAt:     trackingLink.UpdatedAt,
	}
}

// FromTrackingLinkGenerationDomain converts domain generation response to API response
func FromTrackingLinkGenerationDomain(response *domain.TrackingLinkGenerationResponse, baseURL string) *TrackingLinkGenerationResponse {
	apiResponse := &TrackingLinkGenerationResponse{
		GeneratedURL: response.GeneratedURL,
	}

	if response.TrackingLink != nil {
		apiResponse.TrackingLink = FromTrackingLinkDomain(response.TrackingLink)
		// Add QR code URL if tracking link ID is available
		if response.TrackingLink.TrackingLinkID > 0 {
			qrURL := fmt.Sprintf("%s/api/v1/organizations/%d/tracking-links/%d/qr",
				baseURL, response.TrackingLink.OrganizationID, response.TrackingLink.TrackingLinkID)
			apiResponse.QRCodeURL = &qrURL
		}
	}

	return apiResponse
}

// FromTrackingLinkProviderMappingDomain converts domain provider mapping to API response
func FromTrackingLinkProviderMappingDomain(mapping *domain.TrackingLinkProviderMapping) *TrackingLinkProviderMappingResponse {
	return &TrackingLinkProviderMappingResponse{
		MappingID:              mapping.MappingID,
		TrackingLinkID:         mapping.TrackingLinkID,
		ProviderType:           mapping.ProviderType,
		ProviderTrackingLinkID: mapping.ProviderTrackingLinkID,
		ProviderData:           mapping.ProviderData,
		SyncStatus:             mapping.SyncStatus,
		LastSyncAt:             mapping.LastSyncAt,
		SyncError:              mapping.SyncError,
		CreatedAt:              mapping.CreatedAt,
		UpdatedAt:              mapping.UpdatedAt,
	}
}
