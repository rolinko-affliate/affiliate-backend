package everflow

import (
	"encoding/json"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/everflow/affiliate"
)

// AffiliateProviderMapper handles mapping between domain models and Everflow API models
type AffiliateProviderMapper struct{}

// NewAffiliateProviderMapper creates a new affiliate provider mapper
func NewAffiliateProviderMapper() *AffiliateProviderMapper {
	return &AffiliateProviderMapper{}
}

// MapAffiliateToEverflowRequest creates an Everflow affiliate request from domain affiliate and provider mapping
func (m *AffiliateProviderMapper) MapAffiliateToEverflowRequest(
	aff *domain.Affiliate,
	mapping *domain.AffiliateProviderMapping,
) (*affiliate.CreateAffiliateRequest, error) {
	if aff == nil {
		return nil, fmt.Errorf("affiliate cannot be nil")
	}

	// Get Everflow-specific data from provider mapping if available
	var everflowData *domain.EverflowProviderData
	if mapping != nil && mapping.ProviderData != nil {
		everflowData = &domain.EverflowProviderData{}
		if err := json.Unmarshal([]byte(*mapping.ProviderData), everflowData); err != nil {
			// If unmarshal fails, use default values
			everflowData = &domain.EverflowProviderData{}
		}
	} else {
		everflowData = &domain.EverflowProviderData{}
	}

	// Map account status
	accountStatus := m.mapDomainStatusToEverflowStatus(aff.Status)
	
	// Get network employee ID (required field)
	networkEmployeeID := m.getDefaultNetworkEmployeeID(everflowData)

	// Create the base request with required fields
	req := affiliate.NewCreateAffiliateRequest(aff.Name, accountStatus, networkEmployeeID)

	// Set optional fields
	if aff.InternalNotes != nil {
		req.SetInternalNotes(*aff.InternalNotes)
	}

	// Set default currency ID (required by Everflow)
	if aff.DefaultCurrencyID != nil {
		req.SetDefaultCurrencyId(*aff.DefaultCurrencyID)
	} else {
		// Default to USD if not specified
		req.SetDefaultCurrencyId("USD")
	}

	// Set Everflow-specific fields
	if everflowData.EnableMediaCostTrackingLinks != nil {
		req.SetEnableMediaCostTrackingLinks(*everflowData.EnableMediaCostTrackingLinks)
	}

	if everflowData.ReferrerID != nil {
		req.SetReferrerId(*everflowData.ReferrerID)
	}

	if everflowData.IsContactAddressEnabled != nil {
		req.SetIsContactAddressEnabled(*everflowData.IsContactAddressEnabled)
	}

	if everflowData.NetworkAffiliateTierID != nil {
		req.SetNetworkAffiliateTierId(*everflowData.NetworkAffiliateTierID)
	}

	// Map contact address if available
	if aff.ContactAddress != nil {
		var contactAddr domain.ContactAddress
		if err := json.Unmarshal([]byte(*aff.ContactAddress), &contactAddr); err == nil && contactAddr.HasData() {
			everflowAddr := affiliate.NewContactAddress()
			if contactAddr.Address1 != nil {
				everflowAddr.SetAddress1(*contactAddr.Address1)
			}
			if contactAddr.Address2 != nil {
				everflowAddr.SetAddress2(*contactAddr.Address2)
			}
			if contactAddr.City != nil {
				everflowAddr.SetCity(*contactAddr.City)
			}
			if contactAddr.RegionCode != nil {
				everflowAddr.SetRegionCode(*contactAddr.RegionCode)
			}
			if contactAddr.CountryCode != nil {
				everflowAddr.SetCountryCode(*contactAddr.CountryCode)
			}
			if contactAddr.ZipPostalCode != nil {
				everflowAddr.SetZipPostalCode(*contactAddr.ZipPostalCode)
			}
			req.SetContactAddress(*everflowAddr)
		}
	}

	// Map labels if available
	if aff.Labels != nil {
		var labels []string
		if err := json.Unmarshal([]byte(*aff.Labels), &labels); err == nil && len(labels) > 0 {
			req.SetLabels(labels)
		}
	}

	// Map billing info (required by Everflow)
	if aff.BillingInfo != nil {
		billingInfo, err := m.mapDomainBillingToEverflowBilling(*aff.BillingInfo, aff)
		if err == nil && billingInfo != nil {
			req.SetBilling(*billingInfo)
		}
	} else {
		// Create default billing info if not provided
		billingInfo := m.createDefaultBillingInfo(aff)
		req.SetBilling(*billingInfo)
	}

	// Map users if available from Everflow data
	if everflowData.Users != nil && len(*everflowData.Users) > 0 {
		users, err := m.mapUsersFromEverflowData(*everflowData.Users)
		if err == nil && len(users) > 0 {
			req.SetUsers(users)
		}
	}

	return req, nil
}

// MapEverflowResponseToProviderData converts Everflow response to provider data
// Note: This now only maps Everflow-specific fields, as general purpose fields
// are handled separately and stored in the main affiliate model
func (m *AffiliateProviderMapper) MapEverflowResponseToProviderData(
	resp *affiliate.Affiliate,
) (*domain.EverflowProviderData, error) {
	if resp == nil {
		return nil, fmt.Errorf("response cannot be nil")
	}

	everflowData := &domain.EverflowProviderData{}

	// Map Everflow-specific fields
	if resp.HasNetworkAffiliateId() {
		networkAffiliateID := resp.GetNetworkAffiliateId()
		everflowData.NetworkAffiliateID = &networkAffiliateID
	}

	if resp.HasNetworkEmployeeId() {
		networkEmployeeID := resp.GetNetworkEmployeeId()
		everflowData.NetworkEmployeeID = &networkEmployeeID
	}

	if resp.HasEnableMediaCostTrackingLinks() {
		enableMediaCost := resp.GetEnableMediaCostTrackingLinks()
		everflowData.EnableMediaCostTrackingLinks = &enableMediaCost
	}

	if resp.HasReferrerId() {
		referrerID := resp.GetReferrerId()
		everflowData.ReferrerID = &referrerID
	}

	if resp.HasIsContactAddressEnabled() {
		isContactAddressEnabled := resp.GetIsContactAddressEnabled()
		everflowData.IsContactAddressEnabled = &isContactAddressEnabled
	}

	// Initialize additional fields map if needed
	everflowData.AdditionalFields = make(map[string]interface{})

	// Store additional Everflow fields that might be useful
	if resp.HasNetworkId() {
		everflowData.AdditionalFields["network_id"] = resp.GetNetworkId()
	}

	if resp.HasAccountManagerId() {
		everflowData.AdditionalFields["account_manager_id"] = resp.GetAccountManagerId()
	}

	if resp.HasAccountManagerName() {
		everflowData.AdditionalFields["account_manager_name"] = resp.GetAccountManagerName()
	}

	if resp.HasTimeCreated() {
		everflowData.AdditionalFields["time_created"] = resp.GetTimeCreated()
	}

	if resp.HasTimeSaved() {
		everflowData.AdditionalFields["time_saved"] = resp.GetTimeSaved()
	}

	return everflowData, nil
}

// MapEverflowResponseToAffiliate converts Everflow response to main affiliate fields
// This handles the general purpose fields that were moved to the main affiliate model
func (m *AffiliateProviderMapper) MapEverflowResponseToAffiliate(
	resp *affiliate.Affiliate,
	aff *domain.Affiliate,
) error {
	if resp == nil || aff == nil {
		return fmt.Errorf("response and affiliate cannot be nil")
	}

	// Update general purpose fields from Everflow response
	if resp.HasName() {
		aff.Name = resp.GetName()
	}

	if resp.HasAccountStatus() {
		aff.Status = m.mapEverflowStatusToDomainStatus(resp.GetAccountStatus())
	}

	if resp.HasInternalNotes() {
		notes := resp.GetInternalNotes()
		aff.InternalNotes = &notes
	}

	if resp.HasDefaultCurrencyId() {
		currencyID := resp.GetDefaultCurrencyId()
		aff.DefaultCurrencyID = &currencyID
	}

	// Update labels if available
	if len(resp.Labels) > 0 {
		labelsJSON, err := json.Marshal(resp.Labels)
		if err == nil {
			labelsStr := string(labelsJSON)
			aff.Labels = &labelsStr
		}
	}

	return nil
}

// mapEverflowStatusToDomainStatus converts Everflow status to domain status
func (m *AffiliateProviderMapper) mapEverflowStatusToDomainStatus(status string) string {
	switch status {
	case "active":
		return "active"
	case "pending":
		return "pending"
	case "rejected":
		return "rejected"
	case "inactive":
		return "inactive"
	default:
		return "pending"
	}
}

// Helper methods
func (m *AffiliateProviderMapper) mapDomainStatusToEverflowStatus(status string) string {
	switch status {
	case "active":
		return "active"
	case "pending":
		return "pending"
	case "rejected":
		return "rejected"
	case "inactive":
		return "inactive"
	default:
		return "pending"
	}
}

func (m *AffiliateProviderMapper) getDefaultNetworkEmployeeID(everflowData *domain.EverflowProviderData) int32 {
	if everflowData != nil && everflowData.NetworkEmployeeID != nil {
		return *everflowData.NetworkEmployeeID
	}
	return 1 // Default network employee ID
}

// mapDomainBillingToEverflowBilling converts domain billing info to Everflow billing structure
func (m *AffiliateProviderMapper) mapDomainBillingToEverflowBilling(billingInfoJSON string, aff *domain.Affiliate) (*affiliate.BillingInfo, error) {
	// Parse the billing info JSON
	var billingData map[string]interface{}
	if err := json.Unmarshal([]byte(billingInfoJSON), &billingData); err != nil {
		return nil, fmt.Errorf("failed to parse billing info: %w", err)
	}

	billingInfo := affiliate.NewBillingInfo()

	// Map billing frequency
	if freq, ok := billingData["billing_frequency"].(string); ok {
		billingInfo.SetBillingFrequency(freq)
	} else {
		billingInfo.SetBillingFrequency("monthly") // Default
	}

	// Map payment type
	if paymentType, ok := billingData["payment_type"].(string); ok {
		billingInfo.SetPaymentType(paymentType)
	} else {
		billingInfo.SetPaymentType("none") // Default
	}

	// Map tax ID
	if taxID, ok := billingData["tax_id"].(string); ok {
		billingInfo.SetTaxId(taxID)
	}

	// Map invoice amount threshold from affiliate
	if aff.InvoiceAmountThreshold != nil {
		billingInfo.SetInvoiceAmountThreshold(*aff.InvoiceAmountThreshold)
	}

	// Map default payment terms from affiliate
	if aff.DefaultPaymentTerms != nil {
		billingInfo.SetDefaultPaymentTerms(*aff.DefaultPaymentTerms)
	}

	// Map billing details (day of month, etc.)
	if detailsData, ok := billingData["details"].(map[string]interface{}); ok {
		details := affiliate.NewBillingDetails()
		if dayOfMonth, ok := detailsData["day_of_month"].(float64); ok {
			details.SetDayOfMonth(int32(dayOfMonth))
		} else {
			details.SetDayOfMonth(1) // Default to 1st of month
		}
		billingInfo.SetDetails(*details)
	} else {
		// Set default billing details
		details := affiliate.NewBillingDetails()
		details.SetDayOfMonth(1)
		billingInfo.SetDetails(*details)
	}

	// Map payment details if available
	if aff.PaymentDetails != nil {
		paymentDetails, err := m.mapDomainPaymentToEverflowPayment(*aff.PaymentDetails)
		if err == nil && paymentDetails != nil {
			billingInfo.SetPayment(*paymentDetails)
		}
	}

	return billingInfo, nil
}

// mapDomainPaymentToEverflowPayment converts domain payment details to Everflow payment structure
func (m *AffiliateProviderMapper) mapDomainPaymentToEverflowPayment(paymentDetailsJSON string) (*affiliate.PaymentDetails, error) {
	var paymentData map[string]interface{}
	if err := json.Unmarshal([]byte(paymentDetailsJSON), &paymentData); err != nil {
		return nil, fmt.Errorf("failed to parse payment details: %w", err)
	}

	paymentDetails := affiliate.NewPaymentDetails()

	// Map bank account details
	if bankAccount, ok := paymentData["bank_account"].(string); ok {
		paymentDetails.SetAccountNumber(bankAccount)
	}

	if routingNumber, ok := paymentData["routing_number"].(string); ok {
		paymentDetails.SetRoutingNumber(routingNumber)
	}

	if accountName, ok := paymentData["account_name"].(string); ok {
		paymentDetails.SetAccountName(accountName)
	}

	if bankName, ok := paymentData["bank_name"].(string); ok {
		paymentDetails.SetBankName(bankName)
	}

	if bankAddress, ok := paymentData["bank_address"].(string); ok {
		paymentDetails.SetBankAddress(bankAddress)
	}

	if swiftCode, ok := paymentData["swift_code"].(string); ok {
		paymentDetails.SetSwiftCode(swiftCode)
	}

	// Map other payment methods
	if email, ok := paymentData["email"].(string); ok {
		paymentDetails.SetEmail(email)
	}

	if paxumId, ok := paymentData["paxum_id"].(string); ok {
		paymentDetails.SetPaxumId(paxumId)
	}

	return paymentDetails, nil
}

// mapUsersFromEverflowData converts user data from Everflow provider data
func (m *AffiliateProviderMapper) mapUsersFromEverflowData(usersData []interface{}) ([]affiliate.AffiliateUser, error) {
	var users []affiliate.AffiliateUser

	for _, userData := range usersData {
		userMap, ok := userData.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract required fields with defaults
		firstName := ""
		if fn, ok := userMap["first_name"].(string); ok {
			firstName = fn
		}

		lastName := ""
		if ln, ok := userMap["last_name"].(string); ok {
			lastName = ln
		}

		email := ""
		if e, ok := userMap["email"].(string); ok {
			email = e
		}

		accountStatus := "active" // Default status
		if status, ok := userMap["account_status"].(string); ok {
			accountStatus = status
		}

		// Skip if required fields are missing
		if firstName == "" || lastName == "" || email == "" {
			continue
		}

		user := affiliate.NewAffiliateUser(firstName, lastName, email, accountStatus)

		// Set optional fields
		if title, ok := userMap["title"].(string); ok {
			user.SetTitle(title)
		}

		if workPhone, ok := userMap["work_phone"].(string); ok {
			user.SetWorkPhone(workPhone)
		}

		if cellPhone, ok := userMap["cell_phone"].(string); ok {
			user.SetCellPhone(cellPhone)
		}

		if password, ok := userMap["initial_password"].(string); ok {
			user.SetInitialPassword(password)
		}

		users = append(users, *user)
	}

	return users, nil
}

// MapEverflowResponseToProviderMapping updates provider mapping with Everflow response data
func (m *AffiliateProviderMapper) MapEverflowResponseToProviderMapping(
	resp interface{},
	mapping *domain.AffiliateProviderMapping,
) error {
	if resp == nil || mapping == nil {
		return fmt.Errorf("invalid response or mapping")
	}

	// Try to cast response to Affiliate type
	affiliateResp, ok := resp.(*affiliate.Affiliate)
	if !ok {
		return fmt.Errorf("invalid response type, expected *affiliate.Affiliate")
	}

	// Create or update provider data with Everflow-specific fields
	everflowData := &domain.EverflowProviderData{}

	// Unmarshal existing provider data if it exists
	if mapping.ProviderData != nil {
		if err := json.Unmarshal([]byte(*mapping.ProviderData), everflowData); err != nil {
			// If unmarshal fails, start with empty data
			everflowData = &domain.EverflowProviderData{}
		}
	}

	// Update Everflow-specific fields from response
	if affiliateResp.HasNetworkAffiliateId() {
		networkAffiliateID := affiliateResp.GetNetworkAffiliateId()
		everflowData.NetworkAffiliateID = &networkAffiliateID
	}

	if affiliateResp.HasNetworkEmployeeId() {
		networkEmployeeID := affiliateResp.GetNetworkEmployeeId()
		everflowData.NetworkEmployeeID = &networkEmployeeID
	}

	if affiliateResp.HasEnableMediaCostTrackingLinks() {
		enableMediaCost := affiliateResp.GetEnableMediaCostTrackingLinks()
		everflowData.EnableMediaCostTrackingLinks = &enableMediaCost
	}

	if affiliateResp.HasReferrerId() {
		referrerID := affiliateResp.GetReferrerId()
		everflowData.ReferrerID = &referrerID
	}

	if affiliateResp.HasIsContactAddressEnabled() {
		isContactAddressEnabled := affiliateResp.GetIsContactAddressEnabled()
		everflowData.IsContactAddressEnabled = &isContactAddressEnabled
	}

	// Set provider affiliate ID in mapping
	if affiliateResp.HasNetworkAffiliateId() {
		providerAffiliateID := fmt.Sprintf("%d", affiliateResp.GetNetworkAffiliateId())
		mapping.ProviderAffiliateID = &providerAffiliateID
	}

	// Marshal updated provider data
	providerDataBytes, err := json.Marshal(everflowData)
	if err != nil {
		return fmt.Errorf("error marshaling provider data: %w", err)
	}

	providerDataStr := string(providerDataBytes)
	mapping.ProviderData = &providerDataStr

	return nil
}

// createDefaultBillingInfo creates default billing information for Everflow
func (m *AffiliateProviderMapper) createDefaultBillingInfo(aff *domain.Affiliate) *affiliate.BillingInfo {
	billingInfo := affiliate.NewBillingInfo()
	
	// Set default billing frequency
	billingInfo.SetBillingFrequency("monthly")
	
	// Set default payment type
	billingInfo.SetPaymentType("none")
	
	// Set default tax ID
	billingInfo.SetTaxId("XXXXX")
	
	// Set default billing details (day of month)
	details := affiliate.NewBillingDetails()
	details.SetDayOfMonth(1)
	billingInfo.SetDetails(*details)
	
	// Use affiliate's invoice amount threshold if available
	if aff.InvoiceAmountThreshold != nil {
		billingInfo.SetInvoiceAmountThreshold(*aff.InvoiceAmountThreshold)
	}
	
	// Use affiliate's default payment terms if available
	if aff.DefaultPaymentTerms != nil {
		billingInfo.SetDefaultPaymentTerms(*aff.DefaultPaymentTerms)
	}
	
	return billingInfo
}
