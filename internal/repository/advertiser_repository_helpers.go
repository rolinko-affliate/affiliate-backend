package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
)

func marshalBillingDetails(billing *domain.BillingDetails) (sql.NullString, error) {
	if billing == nil {
		return sql.NullString{}, nil
	}
	
	billingBytes, err := json.Marshal(billing)
	if err != nil {
		return sql.NullString{}, fmt.Errorf("failed to marshal billing details: %w", err)
	}
	
	return sql.NullString{String: string(billingBytes), Valid: true}, nil
}

func scanNullableFields(
	advertiser *domain.Advertiser,
	contactEmail, billingDetails sql.NullString,
	internalNotes, defaultCurrencyID, platformName, platformURL, platformUsername sql.NullString,
	accountingContactEmail, offerIDMacro, affiliateIDMacro sql.NullString,
	attributionMethod, emailAttributionMethod, attributionPriority sql.NullString,
	reportingTimezoneID sql.NullInt32,
	isExposePublisherReporting sql.NullBool,
	everflowSyncStatus, everflowSyncError sql.NullString,
	lastEverflowSyncAt sql.NullTime,
) error {
	if contactEmail.Valid {
		advertiser.ContactEmail = &contactEmail.String
	}
	
	if billingDetails.Valid {
		var billing domain.BillingDetails
		if err := json.Unmarshal([]byte(billingDetails.String), &billing); err != nil {
			return fmt.Errorf("failed to unmarshal billing details: %w", err)
		}
		advertiser.BillingDetails = &billing
	}
	
	if internalNotes.Valid {
		advertiser.InternalNotes = &internalNotes.String
	}
	if defaultCurrencyID.Valid {
		advertiser.DefaultCurrencyID = &defaultCurrencyID.String
	}
	if platformName.Valid {
		advertiser.PlatformName = &platformName.String
	}
	if platformURL.Valid {
		advertiser.PlatformURL = &platformURL.String
	}
	if platformUsername.Valid {
		advertiser.PlatformUsername = &platformUsername.String
	}
	if accountingContactEmail.Valid {
		advertiser.AccountingContactEmail = &accountingContactEmail.String
	}
	if offerIDMacro.Valid {
		advertiser.OfferIDMacro = &offerIDMacro.String
	}
	if affiliateIDMacro.Valid {
		advertiser.AffiliateIDMacro = &affiliateIDMacro.String
	}
	if attributionMethod.Valid {
		advertiser.AttributionMethod = &attributionMethod.String
	}
	if emailAttributionMethod.Valid {
		advertiser.EmailAttributionMethod = &emailAttributionMethod.String
	}
	if attributionPriority.Valid {
		advertiser.AttributionPriority = &attributionPriority.String
	}
	if reportingTimezoneID.Valid {
		timezoneID := int(reportingTimezoneID.Int32)
		advertiser.ReportingTimezoneID = &timezoneID
	}
	if isExposePublisherReporting.Valid {
		advertiser.IsExposePublisherReporting = &isExposePublisherReporting.Bool
	}
	if everflowSyncStatus.Valid {
		advertiser.EverflowSyncStatus = &everflowSyncStatus.String
	}
	if lastEverflowSyncAt.Valid {
		advertiser.LastEverflowSyncAt = &lastEverflowSyncAt.Time
	}
	if everflowSyncError.Valid {
		advertiser.EverflowSyncError = &everflowSyncError.String
	}
	
	return nil
}

const advertiserSelectFields = `
	advertiser_id, organization_id, name, contact_email, billing_details, status,
	internal_notes, default_currency_id, platform_name, platform_url, platform_username,
	accounting_contact_email, offer_id_macro, affiliate_id_macro, attribution_method,
	email_attribution_method, attribution_priority, reporting_timezone_id, is_expose_publisher_reporting,
	everflow_sync_status, last_everflow_sync_at, everflow_sync_error,
	created_at, updated_at`