package everflow

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/everflow/advertiser"
)

// AdvertiserProviderMapper handles mapping between domain advertisers and Everflow advertiser models
type AdvertiserProviderMapper struct{}

// NewAdvertiserProviderMapper creates a new advertiser provider mapper
func NewAdvertiserProviderMapper() *AdvertiserProviderMapper {
	return &AdvertiserProviderMapper{}
}

// MapAdvertiserToEverflowRequest maps a domain advertiser to an Everflow CreateAdvertiserRequest
func (m *AdvertiserProviderMapper) MapAdvertiserToEverflowRequest(adv *domain.Advertiser, mapping *domain.AdvertiserProviderMapping) (*advertiser.CreateAdvertiserRequest, error) {
	if adv == nil {
		return nil, fmt.Errorf("advertiser cannot be nil")
	}

	// Required fields with defaults
	accountStatus := "active"
	if adv.Status != "" && adv.Status != "pending" {
		accountStatus = adv.Status
	}

	// Default values for required fields
	networkEmployeeId := int32(1) // Default employee ID
	defaultCurrencyId := "USD"
	if adv.DefaultCurrencyID != nil && *adv.DefaultCurrencyID != "" {
		defaultCurrencyId = *adv.DefaultCurrencyID
	}

	reportingTimezoneId := int32(80) // Default timezone (EST)
	if adv.ReportingTimezoneID != nil {
		reportingTimezoneId = *adv.ReportingTimezoneID
	}

	attributionMethod := "last_touch"
	if adv.AttributionMethod != nil && *adv.AttributionMethod != "" {
		attributionMethod = *adv.AttributionMethod
	}

	emailAttributionMethod := "last_affiliate_attribution"
	if adv.EmailAttributionMethod != nil && *adv.EmailAttributionMethod != "" {
		emailAttributionMethod = *adv.EmailAttributionMethod
	}

	attributionPriority := "click"
	if adv.AttributionPriority != nil && *adv.AttributionPriority != "" {
		attributionPriority = *adv.AttributionPriority
	}

	// Create the base request
	req := advertiser.NewCreateAdvertiserRequest(
		adv.Name,
		accountStatus,
		networkEmployeeId,
		defaultCurrencyId,
		reportingTimezoneId,
		attributionMethod,
		emailAttributionMethod,
		attributionPriority,
	)

	// Set optional fields
	if adv.InternalNotes != nil {
		req.SetInternalNotes(*adv.InternalNotes)
	}

	if adv.PlatformName != nil {
		req.SetPlatformName(*adv.PlatformName)
	}

	if adv.PlatformURL != nil {
		req.SetPlatformUrl(*adv.PlatformURL)
	}

	if adv.PlatformUsername != nil {
		req.SetPlatformUsername(*adv.PlatformUsername)
	}

	if adv.AccountingContactEmail != nil {
		req.SetAccountingContactEmail(*adv.AccountingContactEmail)
	}

	if adv.OfferIDMacro != nil {
		req.SetOfferIdMacro(*adv.OfferIDMacro)
	}

	if adv.AffiliateIDMacro != nil {
		req.SetAffiliateIdMacro(*adv.AffiliateIDMacro)
	}

	// Map billing details if present
	if adv.BillingDetails != nil && adv.BillingDetails.HasData() {
		billing, err := m.mapBillingDetails(adv.BillingDetails)
		if err != nil {
			return nil, fmt.Errorf("failed to map billing details: %w", err)
		}
		req.SetBilling(*billing)
	}

	// Create user if contact email is provided
	if adv.ContactEmail != nil && *adv.ContactEmail != "" {
		users := []advertiser.AdvertiserUser{
			*m.createAdvertiserUser(*adv.ContactEmail, defaultCurrencyId, reportingTimezoneId),
		}
		req.SetUsers(users)
	}

	// Map contact address if billing details contain address information
	if adv.BillingDetails != nil && adv.BillingDetails.Address != nil {
		contactAddress, err := m.mapContactAddress(adv.BillingDetails.Address)
		if err != nil {
			return nil, fmt.Errorf("failed to map contact address: %w", err)
		}
		req.SetContactAddress(*contactAddress)
		req.SetIsContactAddressEnabled(true)
	}

	// Set empty string values for optional fields to match Everflow format
	emptyString := ""
	req.SetAccountingContactEmail(emptyString)
	req.SetOfferIdMacro(emptyString)
	req.SetAffiliateIdMacro(emptyString)
	req.SetPlatformName(emptyString)
	req.SetPlatformUrl(emptyString)
	req.SetPlatformUsername(emptyString)
	
	// Set required fields that might be missing
	req.SetInternalNotes("Some notes not visible to the advertiser")
	salesManagerId := int32(1)
	req.SetSalesManagerId(salesManagerId)
	
	// Set labels array
	labels := []string{"DTC Brand"}
	req.SetLabels(labels)

	// Set default settings
	settings := m.createDefaultSettings()
	req.SetSettings(*settings)

	return req, nil
}

// MapEverflowResponseToAdvertiser maps an Everflow advertiser response to a domain advertiser
func (m *AdvertiserProviderMapper) MapEverflowResponseToAdvertiser(resp *advertiser.Advertiser, adv *domain.Advertiser) {
	if resp == nil || adv == nil {
		return
	}

	// Update fields from Everflow response
	if resp.HasName() {
		adv.Name = resp.GetName()
	}

	if resp.HasAccountStatus() {
		adv.Status = resp.GetAccountStatus()
	}

	if resp.HasDefaultCurrencyId() {
		currency := resp.GetDefaultCurrencyId()
		adv.DefaultCurrencyID = &currency
	}

	if resp.HasInternalNotes() {
		notes := resp.GetInternalNotes()
		adv.InternalNotes = &notes
	}

	if resp.HasPlatformName() {
		platform := resp.GetPlatformName()
		adv.PlatformName = &platform
	}

	if resp.HasPlatformUrl() {
		url := resp.GetPlatformUrl()
		adv.PlatformURL = &url
	}

	if resp.HasPlatformUsername() {
		username := resp.GetPlatformUsername()
		adv.PlatformUsername = &username
	}

	if resp.HasAccountingContactEmail() {
		email := resp.GetAccountingContactEmail()
		adv.AccountingContactEmail = &email
	}

	if resp.HasOfferIdMacro() {
		macro := resp.GetOfferIdMacro()
		adv.OfferIDMacro = &macro
	}

	if resp.HasAffiliateIdMacro() {
		macro := resp.GetAffiliateIdMacro()
		adv.AffiliateIDMacro = &macro
	}

	if resp.HasAttributionMethod() {
		method := resp.GetAttributionMethod()
		adv.AttributionMethod = &method
	}

	if resp.HasEmailAttributionMethod() {
		method := resp.GetEmailAttributionMethod()
		adv.EmailAttributionMethod = &method
	}

	if resp.HasAttributionPriority() {
		priority := resp.GetAttributionPriority()
		adv.AttributionPriority = &priority
	}

	if resp.HasReportingTimezoneId() {
		timezone := resp.GetReportingTimezoneId()
		adv.ReportingTimezoneID = &timezone
	}

	// Update timestamp
	adv.UpdatedAt = time.Now()
}

// MapEverflowResponseToProviderMapping maps Everflow response to provider mapping
func (m *AdvertiserProviderMapper) MapEverflowResponseToProviderMapping(resp *advertiser.Advertiser, mapping *domain.AdvertiserProviderMapping) error {
	if resp == nil || mapping == nil {
		return fmt.Errorf("response and mapping cannot be nil")
	}

	// Set provider advertiser ID
	if resp.HasNetworkAdvertiserId() {
		advertiserId := strconv.Itoa(int(resp.GetNetworkAdvertiserId()))
		mapping.ProviderAdvertiserID = &advertiserId
	}

	// Store Everflow-specific data in ProviderData
	everflowData := domain.EverflowAdvertiserProviderData{}

	if resp.HasNetworkAdvertiserId() {
		networkId := resp.GetNetworkAdvertiserId()
		everflowData.NetworkAdvertiserID = &networkId
	}

	if resp.HasNetworkEmployeeId() {
		employeeId := resp.GetNetworkEmployeeId()
		everflowData.NetworkEmployeeID = &employeeId
	}

	if resp.HasSalesManagerId() {
		salesManagerId := resp.GetSalesManagerId()
		everflowData.SalesManagerID = &salesManagerId
	}

	if resp.HasAddressId() {
		addressId := resp.GetAddressId()
		everflowData.AddressID = &addressId
	}

	if resp.HasIsContactAddressEnabled() {
		isEnabled := resp.GetIsContactAddressEnabled()
		everflowData.IsContactAddressEnabled = &isEnabled
	}

	if resp.HasVerificationToken() {
		token := resp.GetVerificationToken()
		everflowData.VerificationToken = &token
	}

	// Store relationship data if present
	if resp.Relationship != nil {
		// Store labels if present
		if resp.Relationship.HasLabels() {
			labels := resp.Relationship.GetLabels()
			var labelsInterface interface{} = labels
			everflowData.Labels = &labelsInterface
		}

		// Store settings if present
		if resp.Relationship.HasSettings() {
			settings := resp.Relationship.GetSettings()
			var settingsInterface interface{} = settings
			everflowData.Settings = &settingsInterface
		}

		// Store billing if present
		if resp.Relationship.HasBilling() {
			billing := resp.Relationship.GetBilling()
			var billingInterface interface{} = billing
			everflowData.Billing = &billingInterface
		}
	}

	// Note: Users are typically only in CreateAdvertiserRequest, not in response
	// ContactAddress is handled via IsContactAddressEnabled flag

	// Serialize and store provider data
	providerDataJSON, err := json.Marshal(everflowData)
	if err != nil {
		return fmt.Errorf("failed to marshal provider data: %w", err)
	}

	providerDataStr := string(providerDataJSON)
	mapping.ProviderData = &providerDataStr

	return nil
}

// mapBillingDetails maps domain billing details to Everflow billing
func (m *AdvertiserProviderMapper) mapBillingDetails(bd *domain.BillingDetails) (*advertiser.Billing, error) {
	billing := advertiser.NewBillingWithDefaults()

	// Set default billing frequency to match Everflow example
	billing.SetBillingFrequency("other")
	
	if bd.Frequency != nil {
		frequency := string(*bd.Frequency)
		billing.SetBillingFrequency(frequency)
	}

	// Set default tax ID if not provided
	taxId := "123456789"
	if bd.TaxID != nil && *bd.TaxID != "" {
		taxId = *bd.TaxID
	}
	billing.SetTaxId(taxId)

	if bd.IsInvoiceCreationAuto != nil {
		billing.SetIsInvoiceCreationAuto(*bd.IsInvoiceCreationAuto)
	}

	if bd.AutoInvoiceStartDate != nil {
		billing.SetAutoInvoiceStartDate(*bd.AutoInvoiceStartDate)
	}

	if bd.DefaultInvoiceIsHidden != nil {
		billing.SetDefaultInvoiceIsHidden(*bd.DefaultInvoiceIsHidden)
	}

	if bd.InvoiceGenerationDaysDelay != nil {
		billing.SetInvoiceGenerationDaysDelay(*bd.InvoiceGenerationDaysDelay)
	}

	if bd.DefaultPaymentTerms != nil {
		terms := int32(*bd.DefaultPaymentTerms)
		billing.SetDefaultPaymentTerms(terms)
	}

	if bd.InvoiceAmountThreshold != nil {
		billing.SetInvoiceAmountThreshold(*bd.InvoiceAmountThreshold)
	}

	// Always set empty details object to match Everflow format
	details := advertiser.NewBillingDetailsWithDefaults()
	
	// Map billing schedule details if available
	if bd.Schedule != nil {
		if bd.Schedule.DayOfWeek != nil {
			details.SetDayOfWeek(*bd.Schedule.DayOfWeek)
		}

		if bd.Schedule.DayOfMonth != nil {
			details.SetDayOfMonth(*bd.Schedule.DayOfMonth)
		}

		if bd.Schedule.DayOfMonthOne != nil {
			details.SetDayOfMonthOne(*bd.Schedule.DayOfMonthOne)
		}

		if bd.Schedule.DayOfMonthTwo != nil {
			details.SetDayOfMonthTwo(*bd.Schedule.DayOfMonthTwo)
		}

		if bd.Schedule.StartingMonth != nil {
			details.SetStartingMonth(*bd.Schedule.StartingMonth)
		}
	}
	
	billing.SetDetails(*details)

	return billing, nil
}

// createAdvertiserUser creates an advertiser user from contact email
func (m *AdvertiserProviderMapper) createAdvertiserUser(email, currencyId string, timezoneId int32) *advertiser.AdvertiserUser {
	// Extract name from email (simple approach)
	firstName := "User"
	lastName := "Account"

	user := advertiser.NewAdvertiserUser(
		firstName,
		lastName,
		email,
		"active",
		1, // Language ID (English)
		timezoneId,
		currencyId,
	)

	return user
}

// mapContactAddress maps domain billing address to Everflow contact address
func (m *AdvertiserProviderMapper) mapContactAddress(addr *domain.BillingAddress) (*advertiser.ContactAddress, error) {
	if addr == nil {
		return nil, fmt.Errorf("address cannot be nil")
	}

	// Map state to region code (simple mapping for US states)
	regionCode := "NY" // Default to NY if not provided
	if addr.State != nil && *addr.State != "" {
		regionCode = *addr.State
	}

	// Map country to country code
	countryCode := "US"
	if addr.Country != "" {
		countryCode = addr.Country
	}

	contactAddr := advertiser.NewContactAddress(
		addr.Line1,
		addr.City,
		regionCode,
		countryCode,
		addr.PostalCode,
	)

	// Set optional address line 2 if available, otherwise set empty string
	if addr.Line2 != nil && *addr.Line2 != "" {
		contactAddr.SetAddress2(*addr.Line2)
	} else {
		contactAddr.SetAddress2("") // Ensure address_2 is always present
	}

	// Set country ID based on country code (simplified mapping)
	countryId := int32(1) // US country ID (try 1 instead of 840)
	if countryCode == "CA" {
		countryId = 36 // Canada
	}
	contactAddr.SetCountryId(countryId)

	return contactAddr, nil
}

// createDefaultSettings creates default settings for the advertiser
func (m *AdvertiserProviderMapper) createDefaultSettings() *advertiser.Settings {
	settings := advertiser.NewSettingsWithDefaults()

	// Set default exposed variables
	exposedVars := advertiser.NewSettingsExposedVariablesWithDefaults()
	exposedVars.SetAffiliateId(true)
	exposedVars.SetAffiliate(false)
	exposedVars.SetSub1(true)
	exposedVars.SetSub2(true)
	exposedVars.SetSub3(false)
	exposedVars.SetSub4(false)
	exposedVars.SetSub5(false)
	exposedVars.SetSourceId(false)
	exposedVars.SetOfferUrl(false)

	settings.SetExposedVariables(*exposedVars)

	return settings
}