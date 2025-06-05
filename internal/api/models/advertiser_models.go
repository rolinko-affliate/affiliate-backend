package models

import (
	"time"

	"github.com/affiliate-backend/internal/domain"
)

// Helper functions for type conversion
func convertIntToInt32Ptr(i *int) *int32 {
	if i == nil {
		return nil
	}
	val := int32(*i)
	return &val
}

func convertInt32ToIntPtr(i *int32) *int {
	if i == nil {
		return nil
	}
	val := int(*i)
	return &val
}


type CreateAdvertiserRequest struct {
	Name                         string                  `json:"name" binding:"required"`
	Status                       string                  `json:"status,omitempty"`
	OrganizationID               int64                   `json:"organization_id" binding:"required"`
	ContactEmail                 *string                 `json:"contact_email,omitempty"`
	BillingDetails               *domain.BillingDetails  `json:"billing_details,omitempty"`
	InternalNotes                *string                 `json:"internal_notes,omitempty"`
	DefaultCurrencyID            *string                 `json:"default_currency_id,omitempty"`
	PlatformName                 *string                 `json:"platform_name,omitempty"`
	PlatformURL                  *string                 `json:"platform_url,omitempty"`
	PlatformUsername             *string                 `json:"platform_username,omitempty"`
	AccountingContactEmail       *string                 `json:"accounting_contact_email,omitempty"`
	OfferIDMacro                 *string                 `json:"offer_id_macro,omitempty"`
	AffiliateIDMacro             *string                 `json:"affiliate_id_macro,omitempty"`
	AttributionMethod            *string                 `json:"attribution_method,omitempty"`
	EmailAttributionMethod       *string                 `json:"email_attribution_method,omitempty"`
	AttributionPriority          *string                 `json:"attribution_priority,omitempty"`
	ReportingTimezoneID          *int                    `json:"reporting_timezone_id,omitempty"`
}

type UpdateAdvertiserRequest struct {
	Name                         string                  `json:"name" binding:"required"`
	Status                       string                  `json:"status" binding:"required"`
	ContactEmail                 *string                 `json:"contact_email,omitempty"`
	BillingDetails               *domain.BillingDetails  `json:"billing_details,omitempty"`
	InternalNotes                *string                 `json:"internal_notes,omitempty"`
	DefaultCurrencyID            *string                 `json:"default_currency_id,omitempty"`
	PlatformName                 *string                 `json:"platform_name,omitempty"`
	PlatformURL                  *string                 `json:"platform_url,omitempty"`
	PlatformUsername             *string                 `json:"platform_username,omitempty"`
	AccountingContactEmail       *string                 `json:"accounting_contact_email,omitempty"`
	OfferIDMacro                 *string                 `json:"offer_id_macro,omitempty"`
	AffiliateIDMacro             *string                 `json:"affiliate_id_macro,omitempty"`
	AttributionMethod            *string                 `json:"attribution_method,omitempty"`
	EmailAttributionMethod       *string                 `json:"email_attribution_method,omitempty"`
	AttributionPriority          *string                 `json:"attribution_priority,omitempty"`
	ReportingTimezoneID          *int                    `json:"reporting_timezone_id,omitempty"`
}

type AdvertiserResponse struct {
	AdvertiserID                 int64                   `json:"advertiser_id"`
	Name                         string                  `json:"name"`
	Status                       string                  `json:"status"`
	OrganizationID               int64                   `json:"organization_id"`
	ContactEmail                 *string                 `json:"contact_email,omitempty"`
	BillingDetails               *domain.BillingDetails  `json:"billing_details,omitempty"`
	InternalNotes                *string                 `json:"internal_notes,omitempty"`
	DefaultCurrencyID            *string                 `json:"default_currency_id,omitempty"`
	PlatformName                 *string                 `json:"platform_name,omitempty"`
	PlatformURL                  *string                 `json:"platform_url,omitempty"`
	PlatformUsername             *string                 `json:"platform_username,omitempty"`
	AccountingContactEmail       *string                 `json:"accounting_contact_email,omitempty"`
	OfferIDMacro                 *string                 `json:"offer_id_macro,omitempty"`
	AffiliateIDMacro             *string                 `json:"affiliate_id_macro,omitempty"`
	AttributionMethod            *string                 `json:"attribution_method,omitempty"`
	EmailAttributionMethod       *string                 `json:"email_attribution_method,omitempty"`
	AttributionPriority          *string                 `json:"attribution_priority,omitempty"`
	ReportingTimezoneID          *int                    `json:"reporting_timezone_id,omitempty"`
	CreatedAt                    time.Time               `json:"created_at"`
	UpdatedAt                    time.Time               `json:"updated_at"`
}

type AdvertiserWithEverflowResponse struct {
	Advertiser    *AdvertiserResponse                `json:"advertiser"`
	EverflowData  interface{}                        `json:"everflow_data,omitempty"`
	SyncStatus    string                             `json:"sync_status"`
	Discrepancies []domain.AdvertiserDiscrepancy     `json:"discrepancies,omitempty"`
}

type ListAdvertisersResponse struct {
	Advertisers []*AdvertiserResponse `json:"advertisers"`
	Page        int                   `json:"page"`
	PageSize    int                   `json:"page_size"`
	Total       int                   `json:"total"`
}

type CreateProviderMappingRequest struct {
	AdvertiserID         int64   `json:"advertiser_id" binding:"required"`
	ProviderType         string  `json:"provider_type" binding:"required"`
	ProviderAdvertiserID *string `json:"provider_advertiser_id,omitempty"`
	ProviderConfig       *string `json:"provider_config,omitempty"`
}

type UpdateProviderMappingRequest struct {
	ProviderAdvertiserID *string `json:"provider_advertiser_id,omitempty"`
	ProviderConfig       *string `json:"provider_config,omitempty"`
}

type ProviderMappingResponse struct {
	MappingID            int64     `json:"mapping_id"`
	AdvertiserID         int64     `json:"advertiser_id"`
	ProviderType         string    `json:"provider_type"`
	ProviderAdvertiserID *string   `json:"provider_advertiser_id,omitempty"`
	ProviderConfig       *string   `json:"provider_config,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type CreateAdvertiserProviderMappingRequest struct {
	ProviderMapping domain.AdvertiserProviderMapping `json:"provider_mapping" binding:"required"`
}

type CreateAdvertiserProviderMappingResponse struct {
	ProviderMapping domain.AdvertiserProviderMapping `json:"provider_mapping"`
}

type GetAdvertiserProviderMappingResponse struct {
	ProviderMapping domain.AdvertiserProviderMapping `json:"provider_mapping"`
}

type UpdateAdvertiserProviderMappingRequest struct {
	ProviderMapping domain.AdvertiserProviderMapping `json:"provider_mapping" binding:"required"`
}

func (req *CreateAdvertiserRequest) ToDomain() *domain.Advertiser {
	return &domain.Advertiser{
		Name:                       req.Name,
		Status:                     req.Status,
		OrganizationID:             req.OrganizationID,
		ContactEmail:               req.ContactEmail,
		BillingDetails:             req.BillingDetails,
		InternalNotes:              req.InternalNotes,
		DefaultCurrencyID:          req.DefaultCurrencyID,
		PlatformName:               req.PlatformName,
		PlatformURL:                req.PlatformURL,
		PlatformUsername:           req.PlatformUsername,
		AccountingContactEmail:     req.AccountingContactEmail,
		OfferIDMacro:               req.OfferIDMacro,
		AffiliateIDMacro:           req.AffiliateIDMacro,
		AttributionMethod:          req.AttributionMethod,
		EmailAttributionMethod:     req.EmailAttributionMethod,
		AttributionPriority:        req.AttributionPriority,
		ReportingTimezoneID:        convertIntToInt32Ptr(req.ReportingTimezoneID),
	}
}

func (req *UpdateAdvertiserRequest) ToDomain(advertiserID int64, orgID int64) *domain.Advertiser {
	return &domain.Advertiser{
		AdvertiserID:               advertiserID,
		Name:                       req.Name,
		Status:                     req.Status,
		OrganizationID:             orgID,
		ContactEmail:               req.ContactEmail,
		BillingDetails:             req.BillingDetails,
		InternalNotes:              req.InternalNotes,
		DefaultCurrencyID:          req.DefaultCurrencyID,
		PlatformName:               req.PlatformName,
		PlatformURL:                req.PlatformURL,
		PlatformUsername:           req.PlatformUsername,
		AccountingContactEmail:     req.AccountingContactEmail,
		OfferIDMacro:               req.OfferIDMacro,
		AffiliateIDMacro:           req.AffiliateIDMacro,
		AttributionMethod:          req.AttributionMethod,
		EmailAttributionMethod:     req.EmailAttributionMethod,
		AttributionPriority:        req.AttributionPriority,
		ReportingTimezoneID:        convertIntToInt32Ptr(req.ReportingTimezoneID),
	}
}

func ToAdvertiserResponse(advertiser *domain.Advertiser) *AdvertiserResponse {
	return &AdvertiserResponse{
		AdvertiserID:               advertiser.AdvertiserID,
		Name:                       advertiser.Name,
		Status:                     advertiser.Status,
		OrganizationID:             advertiser.OrganizationID,
		ContactEmail:               advertiser.ContactEmail,
		BillingDetails:             advertiser.BillingDetails,
		InternalNotes:              advertiser.InternalNotes,
		DefaultCurrencyID:          advertiser.DefaultCurrencyID,
		PlatformName:               advertiser.PlatformName,
		PlatformURL:                advertiser.PlatformURL,
		PlatformUsername:           advertiser.PlatformUsername,
		AccountingContactEmail:     advertiser.AccountingContactEmail,
		OfferIDMacro:               advertiser.OfferIDMacro,
		AffiliateIDMacro:           advertiser.AffiliateIDMacro,
		AttributionMethod:          advertiser.AttributionMethod,
		EmailAttributionMethod:     advertiser.EmailAttributionMethod,
		AttributionPriority:        advertiser.AttributionPriority,
		ReportingTimezoneID:        convertInt32ToIntPtr(advertiser.ReportingTimezoneID),

		CreatedAt:                  advertiser.CreatedAt,
		UpdatedAt:                  advertiser.UpdatedAt,
	}
}

func ToAdvertiserWithEverflowResponse(data *domain.AdvertiserWithProviderData) *AdvertiserWithEverflowResponse {
	return &AdvertiserWithEverflowResponse{
		Advertiser:    ToAdvertiserResponse(data.Advertiser),
		EverflowData:  data.ProviderData,
		SyncStatus:    data.SyncStatus,
		Discrepancies: data.Discrepancies,
	}
}

func (req *CreateProviderMappingRequest) ToDomain() *domain.AdvertiserProviderMapping {
	return &domain.AdvertiserProviderMapping{
		AdvertiserID:         req.AdvertiserID,
		ProviderType:         req.ProviderType,
		ProviderAdvertiserID: req.ProviderAdvertiserID,
		ProviderConfig:       req.ProviderConfig,
	}
}

func ToProviderMappingResponse(mapping *domain.AdvertiserProviderMapping) *ProviderMappingResponse {
	return &ProviderMappingResponse{
		MappingID:            mapping.MappingID,
		AdvertiserID:         mapping.AdvertiserID,
		ProviderType:         mapping.ProviderType,
		ProviderAdvertiserID: mapping.ProviderAdvertiserID,
		ProviderConfig:       mapping.ProviderConfig,
		CreatedAt:            mapping.CreatedAt,
		UpdatedAt:            mapping.UpdatedAt,
	}
}