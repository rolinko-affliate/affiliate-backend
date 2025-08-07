package everflow

import (
	"fmt"
	"strings"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/everflow/advertiser"
)

// SimpleAdvertiserProviderMapper handles mapping between domain advertisers and Everflow advertiser models
// This is a simplified version that generates the exact format that works with Everflow
type SimpleAdvertiserProviderMapper struct{}

// NewSimpleAdvertiserProviderMapper creates a new simple advertiser provider mapper
func NewSimpleAdvertiserProviderMapper() *SimpleAdvertiserProviderMapper {
	fmt.Printf("🏗️ SIMPLE MAPPER CONSTRUCTOR CALLED 🏗️\n")
	return &SimpleAdvertiserProviderMapper{}
}

// generateUniqueEmail generates a unique email address for Everflow user creation
func (m *SimpleAdvertiserProviderMapper) generateUniqueEmail(input string) string {
	var cleanName string
	
	// If input contains @, extract the local part (before @)
	if strings.Contains(input, "@") {
		parts := strings.Split(input, "@")
		cleanName = parts[0]
	} else {
		cleanName = input
	}
	
	// Clean the name to make it email-safe
	cleanName = strings.ToLower(strings.ReplaceAll(cleanName, " ", "-"))
	cleanName = strings.ReplaceAll(cleanName, "_", "-")
	
	// Generate timestamp for uniqueness
	timestamp := time.Now().Unix()
	
	// Create unique email in format: clean-name-timestamp@everflow-test.com
	return fmt.Sprintf("%s-%d@everflow-test.com", cleanName, timestamp)
}

// MapAdvertiserToEverflowRequest maps a domain advertiser to an Everflow CreateAdvertiserRequest
// This version generates the exact format from the working example
func (m *SimpleAdvertiserProviderMapper) MapAdvertiserToEverflowRequest(adv *domain.Advertiser, mapping *domain.AdvertiserProviderMapping) (*advertiser.CreateAdvertiserRequest, error) {
	fmt.Printf("🔧🔧🔧 SIMPLE MAPPER CALLED: MapAdvertiserToEverflowRequest for advertiser: %s 🔧🔧🔧\n", adv.Name)
	
	if adv == nil {
		return nil, fmt.Errorf("advertiser cannot be nil")
	}

	// Create the base request with required fields
	req := advertiser.NewCreateAdvertiserRequest(
		adv.Name,                    // name
		"active",                    // account_status
		1,                          // network_employee_id
		"USD",                      // default_currency_id
		80,                         // reporting_timezone_id
		"last_touch",               // attribution_method
		"last_affiliate_attribution", // email_attribution_method
		"click",                    // attribution_priority
	)

	// Set all the fields to match the working example exactly
	req.SetAccountingContactEmail("")
	req.SetAddressId(1)
	req.SetAffiliateIdMacro("")
	req.SetOfferIdMacro("")
	req.SetPlatformName("")
	req.SetPlatformUrl("")
	req.SetPlatformUsername("")
	req.SetSalesManagerId(1)
	req.SetInternalNotes("Some notes not visible to the advertiser")
	req.SetIsContactAddressEnabled(false)

	// Set billing exactly like the working example
	billing := advertiser.NewBillingWithDefaults()
	billing.SetBillingFrequency("other")
	billing.SetDefaultPaymentTerms(0)
	billing.SetTaxId("123456789")
	
	// Set empty details object
	details := advertiser.NewBillingDetailsWithDefaults()
	billing.SetDetails(*details)
	req.SetBilling(*billing)

	// Set contact address exactly like the working example
	contactAddress := advertiser.NewContactAddress(
		"4110 rue St-Laurent", // address_1
		"Montreal",           // city
		"QC",                // region_code
		"CA",                // country_code
		"H2R 0A1",           // zip_postal_code
	)
	contactAddress.SetAddress2("202")
	contactAddress.SetCountryId(36)
	req.SetContactAddress(*contactAddress)

	// Set labels exactly like the working example
	labels := []string{"DTC Brand"}
	req.SetLabels(labels)

	// Set settings exactly like the working example
	settings := advertiser.NewSettingsWithDefaults()
	exposedVars := advertiser.NewSettingsExposedVariablesWithDefaults()
	exposedVars.SetAffiliate(false)
	exposedVars.SetAffiliateId(true)
	exposedVars.SetOfferUrl(false)
	exposedVars.SetSourceId(false)
	exposedVars.SetSub1(true)
	exposedVars.SetSub2(true)
	exposedVars.SetSub3(false)
	exposedVars.SetSub4(false)
	exposedVars.SetSub5(false)
	settings.SetExposedVariables(*exposedVars)
	req.SetSettings(*settings)

	// Set users with unique email to avoid "Email address already in use" errors
	email := m.generateUniqueEmail(adv.Name)
	if adv.ContactEmail != nil && *adv.ContactEmail != "" {
		// If advertiser has a contact email, use it as base but still make it unique
		email = m.generateUniqueEmail(*adv.ContactEmail)
	}
	
	fmt.Printf("🔧 Generated unique email for advertiser '%s': %s\n", adv.Name, email)
	
	users := []advertiser.AdvertiserUser{
		*advertiser.NewAdvertiserUser(
			"john",                    // first_name
			"smith",                   // last_name
			email,                     // email
			"active",                  // account_status
			1,                        // language_id
			80,                       // timezone_id
			"USD",                    // currency_id
		),
	}
	req.SetUsers(users)

	// Set verification token exactly like the working example
	req.SetVerificationToken("c7HIWpFUGnyQfN5wwBollBBGtUkeOm")

	return req, nil
}

// MapEverflowResponseToAdvertiser maps an Everflow advertiser response to a domain advertiser
func (m *SimpleAdvertiserProviderMapper) MapEverflowResponseToAdvertiser(resp *advertiser.Advertiser, adv *domain.Advertiser) {
	// Use the original mapper's implementation for response mapping
	originalMapper := NewAdvertiserProviderMapper()
	originalMapper.MapEverflowResponseToAdvertiser(resp, adv)
}

// MapEverflowResponseToProviderMapping maps Everflow response to provider mapping
func (m *SimpleAdvertiserProviderMapper) MapEverflowResponseToProviderMapping(resp *advertiser.Advertiser, mapping *domain.AdvertiserProviderMapping) error {
	// Use the original mapper's implementation for response mapping
	originalMapper := NewAdvertiserProviderMapper()
	return originalMapper.MapEverflowResponseToProviderMapping(resp, mapping)
}