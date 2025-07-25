package repository

import (
	"database/sql"
	"encoding/json"

	"github.com/affiliate-backend/internal/domain"
)

func marshalBillingDetails(billing *domain.BillingDetails) (sql.NullString, error) {
	if billing == nil {
		return sql.NullString{}, nil
	}

	jsonBytes, err := json.Marshal(billing)
	if err != nil {
		return sql.NullString{}, err
	}

	return sql.NullString{String: string(jsonBytes), Valid: true}, nil
}

func scanNullableFields(
	advertiser *domain.Advertiser,
	contactEmail, billingDetails sql.NullString,
	internalNotes, defaultCurrencyID, platformName, platformURL, platformUsername sql.NullString,
	accountingContactEmail, offerIDMacro, affiliateIDMacro sql.NullString,
	attributionMethod, emailAttributionMethod, attributionPriority sql.NullString,
	reportingTimezoneID sql.NullInt32,
) error {
	if contactEmail.Valid {
		advertiser.ContactEmail = &contactEmail.String
	}

	if billingDetails.Valid {
		var bd domain.BillingDetails
		if err := json.Unmarshal([]byte(billingDetails.String), &bd); err == nil {
			advertiser.BillingDetails = &bd
		}
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
		timezoneID := reportingTimezoneID.Int32
		advertiser.ReportingTimezoneID = &timezoneID
	}

	return nil
}

const advertiserSelectFields = `
	advertiser_id, organization_id, name, contact_email, billing_details, status,
	internal_notes, default_currency_id, platform_name, platform_url, platform_username,
	accounting_contact_email, offer_id_macro, affiliate_id_macro, attribution_method,
	email_attribution_method, attribution_priority, reporting_timezone_id,
	created_at, updated_at`
