package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/everflow/affiliate"
	"github.com/affiliate-backend/internal/platform/provider"
	"github.com/affiliate-backend/internal/repository"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"time"
)

// AffiliateService defines the interface for affiliate operations
type AffiliateService interface {
	CreateAffiliate(ctx context.Context, affiliate *domain.Affiliate) (*domain.Affiliate, error)
	GetAffiliateByID(ctx context.Context, id int64) (*domain.Affiliate, error)
	UpdateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error
	ListAffiliatesByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]*domain.Affiliate, error)
	DeleteAffiliate(ctx context.Context, id int64) error

	// Provider sync methods
	SyncAffiliateToProvider(ctx context.Context, affiliateID int64) error
	SyncAffiliateFromProvider(ctx context.Context, affiliateID int64) error

	// Provider mapping methods
	CreateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) (*domain.AffiliateProviderMapping, error)
	GetAffiliateProviderMapping(ctx context.Context, affiliateID int64, providerType string) (*domain.AffiliateProviderMapping, error)
	UpdateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error
	DeleteAffiliateProviderMapping(ctx context.Context, mappingID int64) error
}

// affiliateService implements AffiliateService
type affiliateService struct {
	affiliateRepo       repository.AffiliateRepository
	providerMappingRepo repository.AffiliateProviderMappingRepository
	orgRepo             repository.OrganizationRepository
	integrationService  provider.IntegrationService
}

// NewAffiliateService creates a new affiliate service
func NewAffiliateService(
	affiliateRepo repository.AffiliateRepository,
	providerMappingRepo repository.AffiliateProviderMappingRepository,
	orgRepo repository.OrganizationRepository,
	integrationService provider.IntegrationService,
) AffiliateService {
	return &affiliateService{
		affiliateRepo:       affiliateRepo,
		providerMappingRepo: providerMappingRepo,
		orgRepo:             orgRepo,
		integrationService:  integrationService,
	}
}

// CreateAffiliate creates a new affiliate
func (s *affiliateService) CreateAffiliate(ctx context.Context, affiliate *domain.Affiliate) (*domain.Affiliate, error) {
	// Validate organization exists
	_, err := s.orgRepo.GetOrganizationByID(ctx, affiliate.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	if err := s.validateAffiliate(affiliate); err != nil {
		return nil, err
	}

	// Set default status if not provided
	if affiliate.Status == "" {
		affiliate.Status = "pending"
	}

	// Step 1: Insert local record
	if err := s.affiliateRepo.CreateAffiliate(ctx, affiliate); err != nil {
		return nil, fmt.Errorf("failed to create affiliate: %w", err)
	}

	s.addNetworksAffiliates()

	// Step 2: Call IntegrationService to create in provider
	// The integration service handles provider mapping creation internally
	_, err = s.integrationService.CreateAffiliate(ctx, *affiliate)
	if err != nil {
		// Log error but don't fail the operation since local creation succeeded
		fmt.Printf("Warning: failed to create affiliate in provider: %v\n", err)
		return affiliate, nil
	}

	return affiliate, nil
}

// GetAffiliateByID retrieves an affiliate by ID
func (s *affiliateService) GetAffiliateByID(ctx context.Context, id int64) (*domain.Affiliate, error) {
	return s.affiliateRepo.GetAffiliateByID(ctx, id)
}

// UpdateAffiliate updates an affiliate
func (s *affiliateService) UpdateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error {
	if err := s.validateAffiliate(affiliate); err != nil {
		return err
	}

	// Step 1: Update local record first
	if err := s.affiliateRepo.UpdateAffiliate(ctx, affiliate); err != nil {
		return fmt.Errorf("failed to update affiliate: %w", err)
	}

	// Step 2: Check if provider mapping exists
	_, err := s.providerMappingRepo.GetAffiliateProviderMapping(ctx, affiliate.AffiliateID, "everflow")
	if err != nil {
		// No provider mapping exists, skip provider sync
		return nil
	}

	// Step 3: Update in provider if mapping exists
	if err := s.integrationService.UpdateAffiliate(ctx, *affiliate); err != nil {
		// Log error but don't fail the operation since local update succeeded
		fmt.Printf("Warning: failed to update affiliate in provider: %v\n", err)
	}

	return nil
}

// ListAffiliatesByOrganization retrieves a list of affiliates for an organization with pagination
func (s *affiliateService) ListAffiliatesByOrganization(ctx context.Context, orgID int64, page, pageSize int) ([]*domain.Affiliate, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.affiliateRepo.ListAffiliatesByOrganization(ctx, orgID, pageSize, offset)
}

// DeleteAffiliate deletes an affiliate
func (s *affiliateService) DeleteAffiliate(ctx context.Context, id int64) error {
	return s.affiliateRepo.DeleteAffiliate(ctx, id)
}

// CreateAffiliateProviderMapping creates a new affiliate provider mapping
func (s *affiliateService) CreateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) (*domain.AffiliateProviderMapping, error) {
	// Validate affiliate exists
	_, err := s.affiliateRepo.GetAffiliateByID(ctx, mapping.AffiliateID)
	if err != nil {
		return nil, fmt.Errorf("affiliate not found: %w", err)
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

	if err := s.providerMappingRepo.CreateAffiliateProviderMapping(ctx, mapping); err != nil {
		return nil, fmt.Errorf("failed to create affiliate provider mapping: %w", err)
	}

	return mapping, nil
}

// GetAffiliateProviderMapping retrieves an affiliate provider mapping
func (s *affiliateService) GetAffiliateProviderMapping(ctx context.Context, affiliateID int64, providerType string) (*domain.AffiliateProviderMapping, error) {
	return s.providerMappingRepo.GetAffiliateProviderMapping(ctx, affiliateID, providerType)
}

// UpdateAffiliateProviderMapping updates an affiliate provider mapping
func (s *affiliateService) UpdateAffiliateProviderMapping(ctx context.Context, mapping *domain.AffiliateProviderMapping) error {
	// Validate provider config JSON if provided
	if mapping.ProviderConfig != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*mapping.ProviderConfig), &jsonData); err != nil {
			return fmt.Errorf("invalid provider config JSON: %w", err)
		}
	}

	return s.providerMappingRepo.UpdateAffiliateProviderMapping(ctx, mapping)
}

// DeleteAffiliateProviderMapping deletes an affiliate provider mapping
func (s *affiliateService) DeleteAffiliateProviderMapping(ctx context.Context, mappingID int64) error {
	return s.providerMappingRepo.DeleteAffiliateProviderMapping(ctx, mappingID)
}

// SyncAffiliateToProvider syncs an affiliate to the provider
func (s *affiliateService) SyncAffiliateToProvider(ctx context.Context, affiliateID int64) error {
	// Get local affiliate
	affiliate, err := s.affiliateRepo.GetAffiliateByID(ctx, affiliateID)
	if err != nil {
		return fmt.Errorf("failed to get affiliate: %w", err)
	}

	// Check if provider mapping exists
	_, err = s.providerMappingRepo.GetAffiliateProviderMapping(ctx, affiliateID, "everflow")
	if err != nil {
		// No mapping exists, create in provider
		return s.createAffiliateInProvider(ctx, affiliate)
	}

	// Mapping exists, update in provider
	if err := s.integrationService.UpdateAffiliate(ctx, *affiliate); err != nil {
		return fmt.Errorf("failed to sync affiliate to provider: %w", err)
	}

	return nil
}

// SyncAffiliateFromProvider syncs an affiliate from the provider
func (s *affiliateService) SyncAffiliateFromProvider(ctx context.Context, affiliateID int64) error {
	// Get provider mapping
	_, err := s.providerMappingRepo.GetAffiliateProviderMapping(ctx, affiliateID, "everflow")
	if err != nil {
		return fmt.Errorf("no provider mapping found for affiliate %d: %w", affiliateID, err)
	}

	// Convert affiliate ID to UUID for IntegrationService
	affiliateUUID := s.int64ToUUID(affiliateID)

	// Get affiliate from provider
	providerAffiliate, err := s.integrationService.GetAffiliate(ctx, affiliateUUID)
	if err != nil {
		return fmt.Errorf("failed to get affiliate from provider: %w", err)
	}

	// Update local affiliate with provider data
	localAffiliate, err := s.affiliateRepo.GetAffiliateByID(ctx, affiliateID)
	if err != nil {
		return fmt.Errorf("failed to get local affiliate: %w", err)
	}

	// Merge provider data into local affiliate
	s.mergeProviderDataIntoAffiliate(localAffiliate, &providerAffiliate)

	// Update local record
	return s.affiliateRepo.UpdateAffiliate(ctx, localAffiliate)
}

// Helper methods

// createAffiliateInProvider creates an affiliate in the provider when no mapping exists
func (s *affiliateService) createAffiliateInProvider(ctx context.Context, affiliate *domain.Affiliate) error {
	// Create in provider - integration service handles provider mapping creation
	_, err := s.integrationService.CreateAffiliate(ctx, *affiliate)
	if err != nil {
		return fmt.Errorf("failed to create affiliate in provider: %w", err)
	}

	return nil
}

// mergeProviderDataIntoAffiliate merges provider data into local affiliate
func (s *affiliateService) mergeProviderDataIntoAffiliate(local *domain.Affiliate, provider *domain.Affiliate) {
	// Merge relevant fields from provider into local
	// Provider-specific data like NetworkAffiliateID is now stored in provider mappings
	// This function can be used to merge general affiliate data if needed

	// Example: merge status if provider has updated status
	if provider.Status != "" {
		local.Status = provider.Status
	}
	// Add other general fields as needed based on what the provider returns
}

// validateAffiliate validates affiliate fields
func (s *affiliateService) validateAffiliate(affiliate *domain.Affiliate) error {
	if affiliate.Name == "" {
		return fmt.Errorf("affiliate name cannot be empty")
	}

	// Validate status
	validStatuses := map[string]bool{
		"active":   true,
		"pending":  true,
		"inactive": true,
		"rejected": true,
	}
	if !validStatuses[affiliate.Status] {
		return fmt.Errorf("invalid status: %s", affiliate.Status)
	}

	// Validate payment details JSON if provided
	if affiliate.PaymentDetails != nil {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*affiliate.PaymentDetails), &jsonData); err != nil {
			return fmt.Errorf("invalid payment details JSON: %w", err)
		}
	}

	return nil
}

// int64ToUUID converts int64 to UUID (copied from advertiser service)
func (s *affiliateService) int64ToUUID(id int64) uuid.UUID {
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
func (s *affiliateService) addNetworksAffiliates() {
	enable := false
	enables := true
	var name = "Test Affiliate"
	var x = "This is a test affiliate created using the API"
	var i int32 = 1
	var y int32 = 0
	var accountStatuss = "active"
	var defaultCurrencyIds = "USD"
	var str = "XXXXX"
	var none = "none"
	var monthly = "monthly"

	everflowReq := affiliate.CreateAffiliateRequest{
		Name:                         name,
		AccountStatus:                accountStatuss,
		InternalNotes:                &x,
		NetworkEmployeeId:            i,
		DefaultCurrencyId:            &defaultCurrencyIds,
		EnableMediaCostTrackingLinks: &enable,
		ReferrerId:                   &y,
		IsContactAddressEnabled:      &enables,
		Billing: &affiliate.BillingInfo{
			BillingFrequency: &monthly,
			PaymentType:      &none,
			TaxId:            &str,
			Details: &affiliate.BillingDetails{
				DayOfMonth: &i,
			},
		},
	}
	url := "https://api.eflow.team/v1/networks/affiliates" //Everflow创建联盟会员
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

func (s *affiliateService) getNetworksAffiliates() {
	var affiliateId = "5"
	url := "https://api.eflow.team/v1/networks/affiliates/" + affiliateId //按 ID 查询会员详情
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
	defer resp.Body.Close() // 确保关闭响应体[4](@ref)

	// 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}

	fmt.Println("状态码:", resp.StatusCode)
	fmt.Println("响应头:", resp.Header.Get("Content-Type"))

	// 5. 解析JSON数据（示例）
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("JSON解析失败:", err)
		return
	}
	fmt.Printf("解析结果: %+v\n", body)
}

func (s *affiliateService) updateNetworksAffiliates() {
	var affiliateId = "5"
	url := "https://api.eflow.team/v1/networks/affiliates/" + affiliateId //更新联盟会员

	enable := false
	enables := true
	var name = "Test Affiliate003"
	var x = "This is a test affiliate created using the API"
	var i int32 = 1
	var y int32 = 0
	var accountStatuss = "active"
	var defaultCurrencyIds = "USD"
	var str = "XXXXX"
	var none = "none"
	var monthly = "monthly"

	var address1 = "105 Ocean Avenue"
	var address2 = ""
	var city = "Los Angeles"
	var regionCode = "CA"
	var countryCode = "US"
	var zipPostalCode = "GHGHGH"

	everflowReq := affiliate.UpdateAffiliateRequest{
		Name:                         name,
		AccountStatus:                accountStatuss,
		InternalNotes:                &x,
		NetworkEmployeeId:            i,
		DefaultCurrencyId:            &defaultCurrencyIds,
		EnableMediaCostTrackingLinks: &enable,
		ReferrerId:                   &y,
		IsContactAddressEnabled:      &enables,
		ContactAddress: &affiliate.ContactAddress{
			Address1:      &address1,
			Address2:      &address2,
			City:          &city,
			RegionCode:    &regionCode,
			CountryCode:   &countryCode,
			ZipPostalCode: &zipPostalCode,
		},
		Billing: &affiliate.BillingInfo{
			BillingFrequency: &monthly,
			PaymentType:      &none,
			TaxId:            &str,
			Details: &affiliate.BillingDetails{
				DayOfMonth: &i,
			},
		},
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
