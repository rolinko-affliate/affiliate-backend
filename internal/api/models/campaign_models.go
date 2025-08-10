package models

import (
	"time"

	"github.com/affiliate-backend/internal/domain"
)

// CreateCampaignRequest represents the request to create a new campaign
type CreateCampaignRequest struct {
	OrganizationID int64      `json:"organization_id" binding:"required"`
	AdvertiserID   int64      `json:"advertiser_id" binding:"required"`
	Name           string     `json:"name" binding:"required"`
	Description    *string    `json:"description,omitempty"`
	Status         string     `json:"status" binding:"required,oneof=draft active paused archived"`
	StartDate      *time.Time `json:"start_date,omitempty"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	InternalNotes  *string    `json:"internal_notes,omitempty"`

	// Core campaign fields
	DestinationURL     *string `json:"destination_url,omitempty"`
	ThumbnailURL       *string `json:"thumbnail_url,omitempty"`
	PreviewURL         *string `json:"preview_url,omitempty"`
	Visibility         *string `json:"visibility,omitempty" binding:"omitempty,oneof=public require_approval private"`
	CurrencyID         *string `json:"currency_id,omitempty"`
	ConversionMethod   *string `json:"conversion_method,omitempty" binding:"omitempty,oneof=server_postback pixel"`
	SessionDefinition  *string `json:"session_definition,omitempty" binding:"omitempty,oneof=cookie ip fingerprint"`
	SessionDuration    *int32  `json:"session_duration,omitempty"`
	TermsAndConditions *string `json:"terms_and_conditions,omitempty"`

	// Caps and limits
	IsCapsEnabled        *bool `json:"is_caps_enabled,omitempty"`
	DailyConversionCap   *int  `json:"daily_conversion_cap,omitempty"`
	WeeklyConversionCap  *int  `json:"weekly_conversion_cap,omitempty"`
	MonthlyConversionCap *int  `json:"monthly_conversion_cap,omitempty"`
	GlobalConversionCap  *int  `json:"global_conversion_cap,omitempty"`
	DailyClickCap        *int  `json:"daily_click_cap,omitempty"`
	WeeklyClickCap       *int  `json:"weekly_click_cap,omitempty"`
	MonthlyClickCap      *int  `json:"monthly_click_cap,omitempty"`
	GlobalClickCap       *int  `json:"global_click_cap,omitempty"`

	// Simplified billing configuration
	FixedRevenue               *float64 `json:"fixed_revenue,omitempty" binding:"omitempty,min=0"`
	FixedClickAmount           *float64 `json:"fixed_click_amount,omitempty" binding:"omitempty,min=0"`
	FixedConversionAmount      *float64 `json:"fixed_conversion_amount,omitempty" binding:"omitempty,min=0"`
	PercentageConversionAmount *float64 `json:"percentage_conversion_amount,omitempty" binding:"omitempty,min=0,max=100"`
}

// UpdateCampaignRequest represents the request to update an existing campaign
type UpdateCampaignRequest struct {
	Name          string     `json:"name" binding:"required"`
	Description   *string    `json:"description,omitempty"`
	Status        string     `json:"status" binding:"required,oneof=draft active paused archived"`
	StartDate     *time.Time `json:"start_date,omitempty"`
	EndDate       *time.Time `json:"end_date,omitempty"`
	InternalNotes *string    `json:"internal_notes,omitempty"`

	// Core campaign fields
	DestinationURL     *string `json:"destination_url,omitempty"`
	ThumbnailURL       *string `json:"thumbnail_url,omitempty"`
	PreviewURL         *string `json:"preview_url,omitempty"`
	Visibility         *string `json:"visibility,omitempty" binding:"omitempty,oneof=public require_approval private"`
	CurrencyID         *string `json:"currency_id,omitempty"`
	ConversionMethod   *string `json:"conversion_method,omitempty" binding:"omitempty,oneof=server_postback pixel"`
	SessionDefinition  *string `json:"session_definition,omitempty" binding:"omitempty,oneof=cookie ip fingerprint"`
	SessionDuration    *int32  `json:"session_duration,omitempty"`
	TermsAndConditions *string `json:"terms_and_conditions,omitempty"`

	// Caps and limits
	IsCapsEnabled        *bool `json:"is_caps_enabled,omitempty"`
	DailyConversionCap   *int  `json:"daily_conversion_cap,omitempty"`
	WeeklyConversionCap  *int  `json:"weekly_conversion_cap,omitempty"`
	MonthlyConversionCap *int  `json:"monthly_conversion_cap,omitempty"`
	GlobalConversionCap  *int  `json:"global_conversion_cap,omitempty"`
	DailyClickCap        *int  `json:"daily_click_cap,omitempty"`
	WeeklyClickCap       *int  `json:"weekly_click_cap,omitempty"`
	MonthlyClickCap      *int  `json:"monthly_click_cap,omitempty"`
	GlobalClickCap       *int  `json:"global_click_cap,omitempty"`

	// Simplified billing configuration
	FixedRevenue               *float64 `json:"fixed_revenue,omitempty" binding:"omitempty,min=0"`
	FixedClickAmount           *float64 `json:"fixed_click_amount,omitempty" binding:"omitempty,min=0"`
	FixedConversionAmount      *float64 `json:"fixed_conversion_amount,omitempty" binding:"omitempty,min=0"`
	PercentageConversionAmount *float64 `json:"percentage_conversion_amount,omitempty" binding:"omitempty,min=0,max=100"`
}

// CampaignResponse represents the response for campaign operations
type CampaignResponse struct {
	CampaignID     int64      `json:"campaign_id"`
	OrganizationID int64      `json:"organization_id"`
	AdvertiserID   int64      `json:"advertiser_id"`
	Name           string     `json:"name"`
	Description    *string    `json:"description,omitempty"`
	Status         string     `json:"status"`
	StartDate      *time.Time `json:"start_date,omitempty"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	InternalNotes  *string    `json:"internal_notes,omitempty"`

	// Core campaign fields
	DestinationURL     *string `json:"destination_url,omitempty"`
	ThumbnailURL       *string `json:"thumbnail_url,omitempty"`
	PreviewURL         *string `json:"preview_url,omitempty"`
	Visibility         *string `json:"visibility,omitempty"`
	CurrencyID         *string `json:"currency_id,omitempty"`
	ConversionMethod   *string `json:"conversion_method,omitempty"`
	SessionDefinition  *string `json:"session_definition,omitempty"`
	SessionDuration    *int32  `json:"session_duration,omitempty"`
	TermsAndConditions *string `json:"terms_and_conditions,omitempty"`

	// Caps and limits
	IsCapsEnabled        *bool `json:"is_caps_enabled,omitempty"`
	DailyConversionCap   *int  `json:"daily_conversion_cap,omitempty"`
	WeeklyConversionCap  *int  `json:"weekly_conversion_cap,omitempty"`
	MonthlyConversionCap *int  `json:"monthly_conversion_cap,omitempty"`
	GlobalConversionCap  *int  `json:"global_conversion_cap,omitempty"`
	DailyClickCap        *int  `json:"daily_click_cap,omitempty"`
	WeeklyClickCap       *int  `json:"weekly_click_cap,omitempty"`
	MonthlyClickCap      *int  `json:"monthly_click_cap,omitempty"`
	GlobalClickCap       *int  `json:"global_click_cap,omitempty"`

	// Simplified billing configuration
	FixedRevenue               *float64 `json:"fixed_revenue,omitempty"`
	FixedClickAmount           *float64 `json:"fixed_click_amount,omitempty"`
	FixedConversionAmount      *float64 `json:"fixed_conversion_amount,omitempty"`
	PercentageConversionAmount *float64 `json:"percentage_conversion_amount,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CampaignListResponse represents the response for listing campaigns
type CampaignListResponse struct {
	Campaigns []CampaignResponse `json:"campaigns"`
	Total     int                `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
}

// ToCampaignDomain converts CreateCampaignRequest to domain.Campaign
func (r *CreateCampaignRequest) ToCampaignDomain() *domain.Campaign {
	return &domain.Campaign{
		OrganizationID: r.OrganizationID,
		AdvertiserID:   r.AdvertiserID,
		Name:           r.Name,
		Description:    r.Description,
		Status:         r.Status,
		StartDate:      r.StartDate,
		EndDate:        r.EndDate,
		InternalNotes:  r.InternalNotes,

		// Core campaign fields
		DestinationURL:     r.DestinationURL,
		ThumbnailURL:       r.ThumbnailURL,
		PreviewURL:         r.PreviewURL,
		Visibility:         r.Visibility,
		CurrencyID:         r.CurrencyID,
		ConversionMethod:   r.ConversionMethod,
		SessionDefinition:  r.SessionDefinition,
		SessionDuration:    r.SessionDuration,
		TermsAndConditions: r.TermsAndConditions,

		// Caps and limits
		IsCapsEnabled:        r.IsCapsEnabled,
		DailyConversionCap:   r.DailyConversionCap,
		WeeklyConversionCap:  r.WeeklyConversionCap,
		MonthlyConversionCap: r.MonthlyConversionCap,
		GlobalConversionCap:  r.GlobalConversionCap,
		DailyClickCap:        r.DailyClickCap,
		WeeklyClickCap:       r.WeeklyClickCap,
		MonthlyClickCap:      r.MonthlyClickCap,
		GlobalClickCap:       r.GlobalClickCap,

		// Simplified billing configuration
		FixedRevenue:               r.FixedRevenue,
		FixedClickAmount:           r.FixedClickAmount,
		FixedConversionAmount:      r.FixedConversionAmount,
		PercentageConversionAmount: r.PercentageConversionAmount,
	}
}

// UpdateCampaignDomain updates a domain.Campaign with UpdateCampaignRequest data
func (r *UpdateCampaignRequest) UpdateCampaignDomain(campaign *domain.Campaign) {
	campaign.Name = r.Name
	campaign.Description = r.Description
	campaign.Status = r.Status
	campaign.StartDate = r.StartDate
	campaign.EndDate = r.EndDate
	campaign.InternalNotes = r.InternalNotes

	// Core campaign fields
	campaign.DestinationURL = r.DestinationURL
	campaign.ThumbnailURL = r.ThumbnailURL
	campaign.PreviewURL = r.PreviewURL
	campaign.Visibility = r.Visibility
	campaign.CurrencyID = r.CurrencyID
	campaign.ConversionMethod = r.ConversionMethod
	campaign.SessionDefinition = r.SessionDefinition
	campaign.SessionDuration = r.SessionDuration
	campaign.TermsAndConditions = r.TermsAndConditions

	// Caps and limits
	campaign.IsCapsEnabled = r.IsCapsEnabled
	campaign.DailyConversionCap = r.DailyConversionCap
	campaign.WeeklyConversionCap = r.WeeklyConversionCap
	campaign.MonthlyConversionCap = r.MonthlyConversionCap
	campaign.GlobalConversionCap = r.GlobalConversionCap
	campaign.DailyClickCap = r.DailyClickCap
	campaign.WeeklyClickCap = r.WeeklyClickCap
	campaign.MonthlyClickCap = r.MonthlyClickCap
	campaign.GlobalClickCap = r.GlobalClickCap

	// Simplified billing configuration
	campaign.FixedRevenue = r.FixedRevenue
	campaign.FixedClickAmount = r.FixedClickAmount
	campaign.FixedConversionAmount = r.FixedConversionAmount
	campaign.PercentageConversionAmount = r.PercentageConversionAmount
}

// FromCampaignDomain converts domain.Campaign to CampaignResponse
func FromCampaignDomain(campaign *domain.Campaign) *CampaignResponse {
	return &CampaignResponse{
		CampaignID:     campaign.CampaignID,
		OrganizationID: campaign.OrganizationID,
		AdvertiserID:   campaign.AdvertiserID,
		Name:           campaign.Name,
		Description:    campaign.Description,
		Status:         campaign.Status,
		StartDate:      campaign.StartDate,
		EndDate:        campaign.EndDate,
		InternalNotes:  campaign.InternalNotes,

		// Core campaign fields
		DestinationURL:     campaign.DestinationURL,
		ThumbnailURL:       campaign.ThumbnailURL,
		PreviewURL:         campaign.PreviewURL,
		Visibility:         campaign.Visibility,
		CurrencyID:         campaign.CurrencyID,
		ConversionMethod:   campaign.ConversionMethod,
		SessionDefinition:  campaign.SessionDefinition,
		SessionDuration:    campaign.SessionDuration,
		TermsAndConditions: campaign.TermsAndConditions,

		// Caps and limits
		IsCapsEnabled:        campaign.IsCapsEnabled,
		DailyConversionCap:   campaign.DailyConversionCap,
		WeeklyConversionCap:  campaign.WeeklyConversionCap,
		MonthlyConversionCap: campaign.MonthlyConversionCap,
		GlobalConversionCap:  campaign.GlobalConversionCap,
		DailyClickCap:        campaign.DailyClickCap,
		WeeklyClickCap:       campaign.WeeklyClickCap,
		MonthlyClickCap:      campaign.MonthlyClickCap,
		GlobalClickCap:       campaign.GlobalClickCap,

		// Simplified billing configuration
		FixedRevenue:               campaign.FixedRevenue,
		FixedClickAmount:           campaign.FixedClickAmount,
		FixedConversionAmount:      campaign.FixedConversionAmount,
		PercentageConversionAmount: campaign.PercentageConversionAmount,

		CreatedAt: campaign.CreatedAt,
		UpdatedAt: campaign.UpdatedAt,
	}
}

// FromCampaignDomainList converts a list of domain.Campaign to CampaignListResponse
func FromCampaignDomainList(campaigns []*domain.Campaign, total, page, pageSize int) *CampaignListResponse {
	campaignResponses := make([]CampaignResponse, len(campaigns))
	for i, campaign := range campaigns {
		campaignResponses[i] = *FromCampaignDomain(campaign)
	}

	return &CampaignListResponse{
		Campaigns: campaignResponses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	}
}

// GetCampaignProviderMappingResponse represents the response for getting a campaign provider mapping
type GetCampaignProviderMappingResponse struct {
	ProviderMapping *CampaignProviderMappingResponse `json:"provider_mapping"`
}

// CampaignProviderMappingResponse represents a campaign provider mapping in API responses
type CampaignProviderMappingResponse struct {
	MappingID       int64      `json:"mapping_id"`
	CampaignID      int64      `json:"campaign_id"`
	ProviderType    string     `json:"provider_type"`
	ProviderOfferID *string    `json:"provider_offer_id,omitempty"`
	ProviderData    *string    `json:"provider_data,omitempty"`
	IsActiveOnProvider *bool   `json:"is_active_on_provider,omitempty"`
	LastSyncedAt    *time.Time `json:"last_synced_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// FromCampaignProviderMappingDomain converts domain.CampaignProviderMapping to GetCampaignProviderMappingResponse
func FromCampaignProviderMappingDomain(mapping *domain.CampaignProviderMapping) *GetCampaignProviderMappingResponse {
	if mapping == nil {
		return &GetCampaignProviderMappingResponse{
			ProviderMapping: nil,
		}
	}

	return &GetCampaignProviderMappingResponse{
		ProviderMapping: &CampaignProviderMappingResponse{
			MappingID:          mapping.MappingID,
			CampaignID:         mapping.CampaignID,
			ProviderType:       mapping.ProviderType,
			ProviderOfferID:    mapping.ProviderOfferID,
			ProviderData:       mapping.ProviderData,
			IsActiveOnProvider: mapping.IsActiveOnProvider,
			LastSyncedAt:       mapping.LastSyncedAt,
			CreatedAt:          mapping.CreatedAt,
			UpdatedAt:          mapping.UpdatedAt,
		},
	}
}
