package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/crypto"
	"github.com/affiliate-backend/internal/platform/everflow/advertiser"
	"github.com/affiliate-backend/internal/platform/provider"
	"github.com/affiliate-backend/internal/repository"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

type AdvertiserService interface {
	CreateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) (*domain.Advertiser, error)
	GetAdvertiserByID(ctx context.Context, id int64) (*domain.Advertiser, error)
	GetAdvertiserWithProviderData(ctx context.Context, id int64) (*domain.AdvertiserWithProviderData, error)
	UpdateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error
	ListAdvertisersByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]*domain.Advertiser, error)
	DeleteAdvertiser(ctx context.Context, id int64) error

	SyncAdvertiserToProvider(ctx context.Context, advertiserID int64) error
	SyncAdvertiserFromProvider(ctx context.Context, advertiserID int64) error
	CompareAdvertiserWithProvider(ctx context.Context, advertiserID int64) ([]domain.AdvertiserDiscrepancy, error)

	CreateProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) (*domain.AdvertiserProviderMapping, error)
	GetProviderMapping(ctx context.Context, advertiserID int64, providerType string) (*domain.AdvertiserProviderMapping, error)
	GetProviderMappings(ctx context.Context, advertiserID int64) ([]*domain.AdvertiserProviderMapping, error)
	UpdateProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error
	DeleteProviderMapping(ctx context.Context, mappingID int64) error
}

type advertiserService struct {
	advertiserRepo      repository.AdvertiserRepository
	providerMappingRepo repository.AdvertiserProviderMappingRepository
	orgRepo             repository.OrganizationRepository
	cryptoService       crypto.Service
	integrationService  provider.IntegrationService
	trackingLinkService trackingLinkService
}

func NewAdvertiserService(
	advertiserRepo repository.AdvertiserRepository,
	providerMappingRepo repository.AdvertiserProviderMappingRepository,
	orgRepo repository.OrganizationRepository,
	cryptoService crypto.Service,
	integrationService provider.IntegrationService,
) AdvertiserService {
	return &advertiserService{
		advertiserRepo:      advertiserRepo,
		providerMappingRepo: providerMappingRepo,
		orgRepo:             orgRepo,
		cryptoService:       cryptoService,
		integrationService:  integrationService,
	}
}

func (s *advertiserService) CreateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) (*domain.Advertiser, error) {
	// Validate organization exists
	_, err := s.orgRepo.GetOrganizationByID(ctx, advertiser.OrganizationID)
	s.addNetworksAdvertiser()
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	setDefaultStatus(advertiser)

	if err := validateAdvertiser(advertiser); err != nil {
		return nil, err
	}

	// Step 1: Insert local record
	if err := s.advertiserRepo.CreateAdvertiser(ctx, advertiser); err != nil {
		return nil, fmt.Errorf("failed to create advertiser: %w", err)
	}

	// Step 2: Create provider mapping with "pending" status
	now := time.Now()
	mapping := &domain.AdvertiserProviderMapping{
		AdvertiserID:   advertiser.AdvertiserID,
		ProviderType:   "everflow",
		SyncStatus:     stringPtr("pending"),
		LastSyncAt:     &now,
		APICredentials: nil, // Set during configuration
		ProviderConfig: nil, // Set by IntegrationService with full payload
	}

	err = s.providerMappingRepo.CreateMapping(ctx, mapping)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider mapping: %w", err)
	}

	// Step 3: Call IntegrationService to create in provider
	providerAdvertiser, err := s.integrationService.CreateAdvertiser(ctx, *advertiser)
	if err != nil {
		// Update mapping status to "failed"
		mapping.SyncStatus = stringPtr("failed")
		mapping.SyncError = stringPtr(err.Error())
		mapping.LastSyncAt = &now
		s.providerMappingRepo.UpdateMapping(ctx, mapping)
		return nil, fmt.Errorf("failed to create advertiser in provider: %w", err)
	}

	// Step 4: Update mapping with provider ID and "synced" status
	// For now, we'll use the advertiser ID as the provider ID since the integration service
	// doesn't return provider-specific IDs in the mock implementation
	providerID := fmt.Sprintf("%d", providerAdvertiser.AdvertiserID)
	mapping.ProviderAdvertiserID = &providerID
	mapping.SyncStatus = stringPtr("synced")
	mapping.SyncError = nil
	mapping.LastSyncAt = &now
	if err := s.providerMappingRepo.UpdateMapping(ctx, mapping); err != nil {
		// Log error but don't fail since advertiser was created successfully
		fmt.Printf("Warning: failed to update provider mapping status to active: %v\n", err)
	}

	return advertiser, nil
}

// GetAdvertiserByID retrieves an advertiser by ID
func (s *advertiserService) GetAdvertiserByID(ctx context.Context, id int64) (*domain.Advertiser, error) {
	return s.advertiserRepo.GetAdvertiserByID(ctx, id)
}

// UpdateAdvertiser updates an advertiser with Everflow synchronization
func (s *advertiserService) UpdateAdvertiser(ctx context.Context, advertiser *domain.Advertiser) error {
	if err := validateAdvertiser(advertiser); err != nil {
		return err
	}

	// Step 1: Update local record first
	if err := s.advertiserRepo.UpdateAdvertiser(ctx, advertiser); err != nil {
		return fmt.Errorf("failed to update advertiser: %w", err)
	}

	// Step 2: Check if provider mapping exists
	mapping, err := s.providerMappingRepo.GetMappingByAdvertiserAndProvider(ctx, advertiser.AdvertiserID, "everflow")
	if err != nil {
		// No provider mapping exists, skip provider sync
		return nil
	}

	// Step 3: Update in provider if mapping exists
	now := time.Now()
	if err := s.integrationService.UpdateAdvertiser(ctx, *advertiser); err != nil {
		// Log error but don't fail the operation since local update succeeded
		fmt.Printf("Warning: failed to update advertiser in provider: %v\n", err)

		// Update mapping sync status to indicate sync failure
		mapping.SyncStatus = stringPtr("failed")
		mapping.SyncError = stringPtr(err.Error())
		mapping.LastSyncAt = &now
		s.providerMappingRepo.UpdateMapping(ctx, mapping)
	} else {
		// Update mapping sync status to indicate successful sync
		mapping.SyncStatus = stringPtr("synced")
		mapping.SyncError = nil
		mapping.LastSyncAt = &now
		s.providerMappingRepo.UpdateMapping(ctx, mapping)
	}

	return nil
}

// ListAdvertisersByOrganization retrieves a list of advertisers for an organization with pagination
func (s *advertiserService) ListAdvertisersByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]*domain.Advertiser, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.advertiserRepo.ListAdvertisersByOrganization(ctx, orgID, pageSize, offset)
}

// DeleteAdvertiser deletes an advertiser
func (s *advertiserService) DeleteAdvertiser(ctx context.Context, id int64) error {
	return s.advertiserRepo.DeleteAdvertiser(ctx, id)
}

// CreateAdvertiserProviderMapping creates a new advertiser provider mapping
func (s *advertiserService) CreateProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) (*domain.AdvertiserProviderMapping, error) {
	// Validate advertiser exists
	_, err := s.advertiserRepo.GetAdvertiserByID(ctx, mapping.AdvertiserID)
	if err != nil {
		return nil, fmt.Errorf("advertiser not found: %w", err)
	}

	// Validate provider type
	if mapping.ProviderType != "everflow" {
		return nil, fmt.Errorf("invalid provider type: %s", mapping.ProviderType)
	}

	// Validate provider config JSON if provided
	if mapping.ProviderConfig != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*mapping.ProviderConfig), &jsonData); err != nil {
			return nil, fmt.Errorf("invalid provider config JSON: %w", err)
		}
	}

	// Validate API credentials JSON if provided
	if mapping.APICredentials != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*mapping.APICredentials), &jsonData); err != nil {
			return nil, fmt.Errorf("invalid API credentials JSON: %w", err)
		}
	}

	// Validate provider data JSON if provided
	if mapping.ProviderData != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*mapping.ProviderData), &jsonData); err != nil {
			return nil, fmt.Errorf("invalid provider data JSON: %w", err)
		}
	}

	if err := s.providerMappingRepo.CreateMapping(ctx, mapping); err != nil {
		return nil, fmt.Errorf("failed to create advertiser provider mapping: %w", err)
	}

	return mapping, nil
}

// GetAdvertiserProviderMapping retrieves an advertiser provider mapping
func (s *advertiserService) GetProviderMapping(ctx context.Context, advertiserID int64, providerType string) (*domain.AdvertiserProviderMapping, error) {
	return s.providerMappingRepo.GetMappingByAdvertiserAndProvider(ctx, advertiserID, providerType)
}

func (s *advertiserService) GetProviderMappings(ctx context.Context, advertiserID int64) ([]*domain.AdvertiserProviderMapping, error) {
	return s.providerMappingRepo.GetMappingsByAdvertiserID(ctx, advertiserID)
}

// UpdateAdvertiserProviderMapping updates an advertiser provider mapping
func (s *advertiserService) UpdateProviderMapping(ctx context.Context, mapping *domain.AdvertiserProviderMapping) error {
	// Validate provider config JSON if provided
	if mapping.ProviderConfig != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*mapping.ProviderConfig), &jsonData); err != nil {
			return fmt.Errorf("invalid provider config JSON: %w", err)
		}
	}

	// Validate API credentials JSON if provided
	if mapping.APICredentials != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*mapping.APICredentials), &jsonData); err != nil {
			return fmt.Errorf("invalid API credentials JSON: %w", err)
		}
	}

	// Validate provider data JSON if provided
	if mapping.ProviderData != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*mapping.ProviderData), &jsonData); err != nil {
			return fmt.Errorf("invalid provider data JSON: %w", err)
		}
	}

	return s.providerMappingRepo.UpdateMapping(ctx, mapping)
}

// DeleteAdvertiserProviderMapping deletes an advertiser provider mapping
func (s *advertiserService) DeleteProviderMapping(ctx context.Context, mappingID int64) error {
	return s.providerMappingRepo.DeleteMapping(ctx, mappingID)
}

func (s *advertiserService) GetAdvertiserWithProviderData(ctx context.Context, id int64) (*domain.AdvertiserWithProviderData, error) {
	advertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get advertiser: %w", err)
	}
	return &domain.AdvertiserWithProviderData{
		Advertiser: advertiser,
		SyncStatus: "not_implemented",
	}, nil
}

func (s *advertiserService) SyncAdvertiserToProvider(ctx context.Context, advertiserID int64) error {
	// Get local advertiser
	advertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, advertiserID)
	if err != nil {
		return fmt.Errorf("failed to get advertiser: %w", err)
	}

	// Check if provider mapping exists
	mapping, err := s.providerMappingRepo.GetMappingByAdvertiserAndProvider(ctx, advertiserID, "everflow")
	if err != nil {
		// No mapping exists, create in provider
		return s.createAdvertiserInProvider(ctx, advertiser)
	}

	// Mapping exists, update in provider
	now := time.Now()
	if err := s.integrationService.UpdateAdvertiser(ctx, *advertiser); err != nil {
		mapping.SyncStatus = stringPtr("failed")
		mapping.SyncError = stringPtr(err.Error())
		mapping.LastSyncAt = &now
		s.providerMappingRepo.UpdateMapping(ctx, mapping)
		return fmt.Errorf("failed to sync advertiser to provider: %w", err)
	}

	// Update mapping sync status
	mapping.SyncStatus = stringPtr("synced")
	mapping.SyncError = nil
	mapping.LastSyncAt = &now
	return s.providerMappingRepo.UpdateMapping(ctx, mapping)
}

func (s *advertiserService) SyncAdvertiserFromProvider(ctx context.Context, advertiserID int64) error {
	// Get provider mapping
	mapping, err := s.providerMappingRepo.GetMappingByAdvertiserAndProvider(ctx, advertiserID, "everflow")
	if err != nil {
		return fmt.Errorf("no provider mapping found for advertiser %d: %w", advertiserID, err)
	}

	// Convert advertiser ID to UUID for IntegrationService
	advertiserUUID := int64ToUUID(advertiserID)

	// Get advertiser from provider
	providerAdvertiser, err := s.integrationService.GetAdvertiser(ctx, advertiserUUID)
	if err != nil {
		now := time.Now()
		mapping.SyncStatus = stringPtr("failed")
		mapping.SyncError = stringPtr(err.Error())
		mapping.LastSyncAt = &now
		s.providerMappingRepo.UpdateMapping(ctx, mapping)
		return fmt.Errorf("failed to get advertiser from provider: %w", err)
	}

	// Update local advertiser with provider data
	localAdvertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, advertiserID)
	if err != nil {
		return fmt.Errorf("failed to get local advertiser: %w", err)
	}

	// Merge provider data into local advertiser
	s.mergeProviderDataIntoAdvertiser(localAdvertiser, &providerAdvertiser)

	// Update local record
	if err := s.advertiserRepo.UpdateAdvertiser(ctx, localAdvertiser); err != nil {
		return fmt.Errorf("failed to update local advertiser: %w", err)
	}

	// Update mapping sync status
	now := time.Now()
	mapping.SyncStatus = stringPtr("synced")
	mapping.SyncError = nil
	mapping.LastSyncAt = &now
	return s.providerMappingRepo.UpdateMapping(ctx, mapping)
}

func (s *advertiserService) CompareAdvertiserWithProvider(ctx context.Context, advertiserID int64) ([]domain.AdvertiserDiscrepancy, error) {
	// Get local advertiser
	localAdvertiser, err := s.advertiserRepo.GetAdvertiserByID(ctx, advertiserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get local advertiser: %w", err)
	}

	// Get provider mapping
	_, err = s.providerMappingRepo.GetMappingByAdvertiserAndProvider(ctx, advertiserID, "everflow")
	if err != nil {
		return []domain.AdvertiserDiscrepancy{{
			Field:         "provider_mapping",
			LocalValue:    "exists",
			ProviderValue: "missing",
			Severity:      "critical",
		}}, nil
	}

	// Convert advertiser ID to UUID for IntegrationService
	advertiserUUID := int64ToUUID(advertiserID)

	// Get advertiser from provider
	providerAdvertiser, err := s.integrationService.GetAdvertiser(ctx, advertiserUUID)
	if err != nil {
		return []domain.AdvertiserDiscrepancy{{
			Field:         "provider_record",
			LocalValue:    "exists",
			ProviderValue: "missing",
			Severity:      "critical",
		}}, nil
	}

	// Compare fields and return discrepancies
	return s.compareAdvertiserFields(localAdvertiser, &providerAdvertiser), nil
}

// Helper methods

// createAdvertiserInProvider creates an advertiser in the provider when no mapping exists
func (s *advertiserService) createAdvertiserInProvider(ctx context.Context, advertiser *domain.Advertiser) error {
	// Create provider mapping with "pending" status
	now := time.Now()
	mapping := &domain.AdvertiserProviderMapping{
		AdvertiserID:   advertiser.AdvertiserID,
		ProviderType:   "everflow",
		SyncStatus:     stringPtr("pending"),
		LastSyncAt:     &now,
		APICredentials: nil, // Set by IntegrationService
		ProviderConfig: nil, // Set by IntegrationService with full payload
	}

	if err := s.providerMappingRepo.CreateMapping(ctx, mapping); err != nil {
		return fmt.Errorf("failed to create provider mapping: %w", err)
	}

	// Create in provider
	providerAdvertiser, err := s.integrationService.CreateAdvertiser(ctx, *advertiser)
	if err != nil {
		mapping.SyncStatus = stringPtr("failed")
		mapping.SyncError = stringPtr(err.Error())
		mapping.LastSyncAt = &now
		s.providerMappingRepo.UpdateMapping(ctx, mapping)
		return fmt.Errorf("failed to create advertiser in provider: %w", err)
	}

	// Update mapping with provider ID and "synced" status
	providerID := fmt.Sprintf("%d", providerAdvertiser.AdvertiserID)
	mapping.ProviderAdvertiserID = &providerID
	mapping.SyncStatus = stringPtr("synced")
	mapping.SyncError = nil
	mapping.LastSyncAt = &now
	return s.providerMappingRepo.UpdateMapping(ctx, mapping)
}

// mergeProviderDataIntoAdvertiser merges provider data into local advertiser
func (s *advertiserService) mergeProviderDataIntoAdvertiser(local *domain.Advertiser, provider *domain.Advertiser) {
	// Merge relevant fields from provider into local
	// For now, we only merge basic fields since provider-specific fields
	// are handled through the provider mapping
	if provider.Name != "" {
		local.Name = provider.Name
	}
	if provider.Status != "" {
		local.Status = provider.Status
	}
	// Add other fields as needed based on what the provider returns
}

// compareAdvertiserFields compares local and provider advertiser fields
func (s *advertiserService) compareAdvertiserFields(local *domain.Advertiser, provider *domain.Advertiser) []domain.AdvertiserDiscrepancy {
	var discrepancies []domain.AdvertiserDiscrepancy

	// Compare name
	if local.Name != provider.Name {
		discrepancies = append(discrepancies, domain.AdvertiserDiscrepancy{
			Field:         "name",
			LocalValue:    local.Name,
			ProviderValue: provider.Name,
			Severity:      "medium",
		})
	}

	// Compare contact email
	if (local.ContactEmail == nil) != (provider.ContactEmail == nil) ||
		(local.ContactEmail != nil && provider.ContactEmail != nil && *local.ContactEmail != *provider.ContactEmail) {
		discrepancies = append(discrepancies, domain.AdvertiserDiscrepancy{
			Field:         "contact_email",
			LocalValue:    local.ContactEmail,
			ProviderValue: provider.ContactEmail,
			Severity:      "medium",
		})
	}

	// Compare status
	if local.Status != provider.Status {
		discrepancies = append(discrepancies, domain.AdvertiserDiscrepancy{
			Field:         "status",
			LocalValue:    local.Status,
			ProviderValue: provider.Status,
			Severity:      "high",
		})
	}

	// Compare status
	if local.Status != provider.Status {
		discrepancies = append(discrepancies, domain.AdvertiserDiscrepancy{
			Field:         "status",
			LocalValue:    local.Status,
			ProviderValue: provider.Status,
			Severity:      "medium",
		})
	}

	return discrepancies
}

// UUID conversion helpers (copied from IntegrationService)
func int64ToUUID(id int64) uuid.UUID {
	// Convert int64 back to UUID format
	// This is a simplified approach - in production you might want a more sophisticated mapping
	hex := fmt.Sprintf("%015x", id)
	// Pad to 32 characters
	for len(hex) < 32 {
		hex = "0" + hex
	}
	// Format as UUID
	uuidStr := fmt.Sprintf("%s-%s-%s-%s-%s", hex[:8], hex[8:12], hex[12:16], hex[16:20], hex[20:32])
	parsed, _ := uuid.Parse(uuidStr)
	return parsed
}

// setDefaultStatus sets default status for advertiser if not provided
func setDefaultStatus(advertiser *domain.Advertiser) {
	if advertiser.Status == "" {
		advertiser.Status = "pending"
	}
}

// validateAdvertiser validates advertiser fields
func validateAdvertiser(advertiser *domain.Advertiser) error {
	if advertiser.Name == "" {
		return fmt.Errorf("advertiser name cannot be empty")
	}

	if advertiser.ContactEmail == nil || *advertiser.ContactEmail == "" {
		return fmt.Errorf("advertiser contact email cannot be empty")
	}

	// Validate status
	validStatuses := map[string]bool{
		"active":   true,
		"pending":  true,
		"inactive": true,
		"rejected": true,
	}
	if !validStatuses[advertiser.Status] {
		return fmt.Errorf("invalid status: %s", advertiser.Status)
	}

	return nil
}

// stringPtr returns a pointer to the given string
func stringPtr(s string) *string {
	return &s
}

func (s *advertiserService) addNetworksAdvertiser() {

	var str []string = []string{"DTC Brand"}
	var x, y = false, true
	var Address2 = "202"
	var CountryId int32 = 36
	var DefaultPaymentTerms int32 = 0
	var BillingFrequency string = "other"
	var TaxId string = "123456789"
	var InternalNotes string = "Some notes not visible to the advertiser"
	var AddressId int32 = 1
	var SalesManagerId int32 = 1
	var null string = ""

	var VerificationToken string = "c7HIWpFUGnyQfN5wwBollBBGtUkeOm"

	var ContactAddress = advertiser.ContactAddress{
		Address1:      "4110 rue St-Laurent",
		Address2:      &Address2,
		City:          "Montreal",
		ZipPostalCode: "H2R 0A1",
		CountryId:     &CountryId,
		CountryCode:   "CA",
		RegionCode:    "QC",
	}
	var Details = advertiser.BillingDetails{}
	var Billings = advertiser.Billing{
		Details:             &Details,
		DefaultPaymentTerms: &DefaultPaymentTerms,
		BillingFrequency:    &BillingFrequency,
		TaxId:               &TaxId,
	}
	var SettingsExposedVariables = advertiser.SettingsExposedVariables{
		AffiliateId: &y,
		Affiliate:   &x,
		Sub1:        &y,
		Sub2:        &y,
		Sub3:        &x,
		Sub4:        &x,
		Sub5:        &x,
		SourceId:    &x,
		OfferUrl:    &x,
	}
	var Settings = advertiser.Settings{
		ExposedVariables: &SettingsExposedVariables,
	}

	everflowReq := advertiser.CreateAdvertiserRequest{
		Name:                    "test",
		AccountStatus:           "active",
		NetworkEmployeeId:       1,
		InternalNotes:           &InternalNotes,
		AddressId:               &AddressId,
		IsContactAddressEnabled: &x,
		SalesManagerId:          &SalesManagerId,
		DefaultCurrencyId:       "USD",
		PlatformName:            &null,
		PlatformUrl:             &null,
		PlatformUsername:        &null,
		ReportingTimezoneId:     80,
		AttributionMethod:       "last_touch",
		EmailAttributionMethod:  "last_affiliate_attribution",
		AttributionPriority:     "click",
		AccountingContactEmail:  &null,
		VerificationToken:       &VerificationToken,
		OfferIdMacro:            &null,
		AffiliateIdMacro:        &null,
		Labels:                  str,
		Users: []advertiser.AdvertiserUser{
			{
				AccountStatus:   "active",
				LanguageId:      1,
				TimezoneId:      80,
				CurrencyId:      "USD",
				FirstName:       "john",
				LastName:        "smith",
				Email:           "john.doe@example.com",
				InitialPassword: &null,
			},
		},
		ContactAddress: &ContactAddress,
		Billing:        &Billings,
		Settings:       &Settings,
	}

	url := "https://api.eflow.team/v1/networks/advertisers" //Everflow广告
	jsonBody, err := json.Marshal(everflowReq)

	// 包装为 io.Reader
	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Eflow-API-Key", "GReOQMUkSWOvtQnJ1AnWzw")
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("请求失败:", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 处理响应
	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("状态码:", resp.StatusCode)
	fmt.Println("响应内容:", string(bodyBytes))
}

func (s *advertiserService) getNetworksAdvertiser() {
	var AdvertiserId = "3"
	url := "https://api.eflow.team/v1/networks/advertisers/" + AdvertiserId //按 ID 查询会员详情

	//  创建带上下文的请求（支持超时控制）
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}

	//请求头（
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Eflow-API-Key", "GReOQMUkSWOvtQnJ1AnWzw") // Everflow认证头[1,6](@ref)

	// 发送请求并处理响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}

	fmt.Println("状态码:", resp.StatusCode)
	fmt.Println("响应头:", resp.Header.Get("Content-Type"))

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("JSON解析失败:", err)
		return
	}
	fmt.Printf("解析结果: %+v\n", body)
}
func (s *advertiserService) updateNetworksAdvertiser() {
	var AdvertiserId = "3"
	url := "https://api.eflow.team/v1/networks/advertisers/" + AdvertiserId //更新联盟会员

	var str []string = []string{"DTC Brand"}
	var x, y = false, true
	var Address2 = "202"
	var CountryId int32 = 36
	var DefaultPaymentTerms int32 = 0
	var BillingFrequency string = "other"
	var TaxId string = "123456789"
	var InternalNotes string = "Some notes not visible to the advertiser"
	var AddressId int32 = 1
	var SalesManagerId int32 = 1
	var null string = ""

	var VerificationToken string = "c7HIWpFUGnyQfN5wwBollBBGtUkeOm"

	var ContactAddress = advertiser.ContactAddress{
		Address1:      "4110 rue St-Laurent",
		Address2:      &Address2,
		City:          "Montreal",
		ZipPostalCode: "H2R 0A1",
		CountryId:     &CountryId,
		CountryCode:   "CA",
		RegionCode:    "QC",
	}
	var Details = advertiser.BillingDetails{}
	var Billings = advertiser.Billing{
		Details:             &Details,
		DefaultPaymentTerms: &DefaultPaymentTerms,
		BillingFrequency:    &BillingFrequency,
		TaxId:               &TaxId,
	}
	var SettingsExposedVariables = advertiser.SettingsExposedVariables{
		AffiliateId: &y,
		Affiliate:   &x,
		Sub1:        &y,
		Sub2:        &y,
		Sub3:        &x,
		Sub4:        &x,
		Sub5:        &x,
		SourceId:    &x,
		OfferUrl:    &x,
	}
	var Settings = advertiser.Settings{
		ExposedVariables: &SettingsExposedVariables,
	}

	everflowReq := advertiser.CreateAdvertiserRequest{
		Name:                    "test",
		AccountStatus:           "active",
		NetworkEmployeeId:       1,
		InternalNotes:           &InternalNotes,
		AddressId:               &AddressId,
		IsContactAddressEnabled: &x,
		SalesManagerId:          &SalesManagerId,
		DefaultCurrencyId:       "USD",
		PlatformName:            &null,
		PlatformUrl:             &null,
		PlatformUsername:        &null,
		ReportingTimezoneId:     80,
		AttributionMethod:       "last_touch",
		EmailAttributionMethod:  "last_affiliate_attribution",
		AttributionPriority:     "click",
		AccountingContactEmail:  &null,
		VerificationToken:       &VerificationToken,
		OfferIdMacro:            &null,
		AffiliateIdMacro:        &null,
		Labels:                  str,
		Users: []advertiser.AdvertiserUser{
			{
				AccountStatus:   "active",
				LanguageId:      1,
				TimezoneId:      80,
				CurrencyId:      "USD",
				FirstName:       "john",
				LastName:        "smith",
				Email:           "john.doe@example.com",
				InitialPassword: &null,
			},
		},
		ContactAddress: &ContactAddress,
		Billing:        &Billings,
		Settings:       &Settings,
	}
	jsonBody, err := json.Marshal(everflowReq)

	// 包装为 io.Reader
	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodPut, url, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Eflow-API-Key", "GReOQMUkSWOvtQnJ1AnWzw")
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("请求失败:", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 处理响应
	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("状态码:", resp.StatusCode)
	fmt.Println("响应内容:", string(bodyBytes))
}
