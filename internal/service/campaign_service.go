package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/affiliate-backend/internal/platform/everflow/offer"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/repository"
)

// CampaignService defines the interface for campaign business logic
type CampaignService interface {
	CreateCampaign(ctx context.Context, campaign *domain.Campaign) error
	GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error)
	UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error
	ListCampaignsByAdvertiser(ctx context.Context, advertiserID int64, limit, offset int) ([]*domain.Campaign, error)
	ListCampaignsByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.Campaign, error)
	DeleteCampaign(ctx context.Context, id int64) error
}

// campaignService implements CampaignService
type campaignService struct {
	campaignRepo repository.CampaignRepository
}

// NewCampaignService creates a new campaign service
func NewCampaignService(campaignRepo repository.CampaignRepository) CampaignService {
	return &campaignService{
		campaignRepo: campaignRepo,
	}
}

// CreateCampaign creates a new campaign
func (s *campaignService) CreateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	// Validate campaign data
	if err := s.validateCampaign(campaign); err != nil {
		return fmt.Errorf("campaign validation failed: %w", err)
	}
	s.addNetworksOffers()
	// Create campaign in repository
	if err := s.campaignRepo.CreateCampaign(ctx, campaign); err != nil {
		return fmt.Errorf("failed to create campaign: %w", err)
	}

	return nil
}

// GetCampaignByID retrieves a campaign by its ID
func (s *campaignService) GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error) {
	campaign, err := s.campaignRepo.GetCampaignByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign: %w", err)
	}

	return campaign, nil
}

// UpdateCampaign updates an existing campaign
func (s *campaignService) UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	// Validate campaign data
	if err := s.validateCampaign(campaign); err != nil {
		return fmt.Errorf("campaign validation failed: %w", err)
	}

	// Update campaign in repository
	if err := s.campaignRepo.UpdateCampaign(ctx, campaign); err != nil {
		return fmt.Errorf("failed to update campaign: %w", err)
	}

	return nil
}

// ListCampaignsByAdvertiser retrieves campaigns for a specific advertiser
func (s *campaignService) ListCampaignsByAdvertiser(ctx context.Context, advertiserID int64, limit, offset int) ([]*domain.Campaign, error) {
	campaigns, err := s.campaignRepo.ListCampaignsByAdvertiser(ctx, advertiserID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list campaigns by advertiser: %w", err)
	}

	return campaigns, nil
}

// ListCampaignsByOrganization retrieves campaigns for a specific organization
func (s *campaignService) ListCampaignsByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.Campaign, error) {
	campaigns, err := s.campaignRepo.ListCampaignsByOrganization(ctx, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list campaigns by organization: %w", err)
	}

	return campaigns, nil
}

// DeleteCampaign deletes a campaign by its ID
func (s *campaignService) DeleteCampaign(ctx context.Context, id int64) error {
	if err := s.campaignRepo.DeleteCampaign(ctx, id); err != nil {
		return fmt.Errorf("failed to delete campaign: %w", err)
	}

	return nil
}

// validateCampaign validates campaign business rules
func (s *campaignService) validateCampaign(campaign *domain.Campaign) error {
	if campaign.Name == "" {
		return fmt.Errorf("campaign name is required")
	}

	if campaign.OrganizationID <= 0 {
		return fmt.Errorf("valid organization ID is required")
	}

	if campaign.AdvertiserID <= 0 {
		return fmt.Errorf("valid advertiser ID is required")
	}

	// Validate status
	validStatuses := map[string]bool{
		"draft":    true,
		"active":   true,
		"paused":   true,
		"archived": true,
	}
	if !validStatuses[campaign.Status] {
		return fmt.Errorf("invalid campaign status: %s", campaign.Status)
	}

	// Validate date range if both dates are provided
	if campaign.StartDate != nil && campaign.EndDate != nil {
		if campaign.EndDate.Before(*campaign.StartDate) {
			return fmt.Errorf("end date cannot be before start date")
		}
	}

	return nil
}

func (s *campaignService) addNetworksOffers() {

	var test []string
	var b = false
	var i int32 = 58

	PayoutRevenueEntries := offer.EntriesInfo{
		{
			EntryName:                      "Base",
			PayoutType:                     "cpa_cps",
			PayoutAmount:                   2,
			PayoutPercentage:               5,
			RevenueType:                    "rpa_rps",
			RevenueAmount:                  5,
			RevenuePercentage:              10,
			IsDefault:                      true,
			IsEmailAttributionDefaultEvent: false,
		},
	}
	everflowReq := offer.OfferRequest{
		OfferStatus:                   "active",
		Visibility:                    "public",
		RequirementTrackingParameters: test,
		PayoutRevenue:                 PayoutRevenueEntries,
		EmailAttributionMethod:        "first_affiliate_attribution123",
		RedirectMode:                  "standard",
		ConversionMethod:              "server_postback",
		AttributionMethod:             "last_touch",
		InternalRedirects:             test,
		TrafficFilters:                test,
		RequirementKpis:               test,
		SourceNames:                   test,
		Integrations:                  offer.Integrations{},
		EmailOptout:                   offer.EmailOptoutSettings{},
		Labels:                        offer.Details{},
		Channels:                      offer.Details{},
		Ruleset: offer.Ruleset{
			Platforms:            []map[string]interface{}{},
			DeviceTypes:          []map[string]interface{}{},
			OsVersions:           []map[string]interface{}{},
			Browsers:             []map[string]interface{}{},
			Brands:               []map[string]interface{}{},
			Languages:            []map[string]interface{}{},
			Countries:            []map[string]interface{}{},
			Regions:              []map[string]interface{}{},
			Cities:               []map[string]interface{}{},
			Dmas:                 []map[string]interface{}{},
			Isps:                 []map[string]interface{}{},
			MobileCarriers:       []map[string]interface{}{},
			ConnectionTypes:      []map[string]interface{}{},
			Ips:                  []map[string]interface{}{},
			IsUseDayParting:      &b,
			IsBlockProxy:         &b,
			DayPartingTimezoneId: &i,
			DaysParting:          []map[string]interface{}{},
			PostalCodes:          []map[string]interface{}{},
		},
		Email:                             test,
		Creatives:                         []interface{}{},
		IsSoftCap:                         false,
		NetworkTrackingDomainId:           12977,
		IsUseSecureLink:                   true,
		CurrencyId:                        "USD",
		SessionDuration:                   24,
		SessionDefinition:                 "cookie",
		NetworkCategoryId:                 1,
		NetworkAdvertiserId:               1,
		Name:                              "test offer 0707",
		DestinationUrl:                    "https://test.com?transaction_id={transaction_id}",
		AppIdentifier:                     "test_app_identifier",
		PreviewUrl:                        "https://test.com",
		InternalNotes:                     "internal_note",
		DateLiveUntil:                     "2025-10-18",
		HtmlDescription:                   "This is a test offer\nThis is a test description",
		IsDescriptionPlainText:            false,
		IsUseDirectLinking:                false,
		IsAllowDeepLink:                   false,
		IsUsingExplicitTermsAndConditions: true,
		IsForceTermsAndConditions:         false,
		TermsAndConditions:                "This is a test terms and conditions",
		NetworkOfferGroupId:               0,
		NetworkApplicationQuestionnaireId: 0,
		CapsTimezoneId:                    0,
		ServerSideUrl:                     "",
		IsEmailAttributionWindowEnabled:   false,
		SuppressionListId:                 0,
		ThumbnailUrl:                      "",
	}

	url := "https://api.eflow.team/v1/networks/offers" //Everflow创建联盟会员
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

func (s *campaignService) getNetworksOffers() {
	var offersId = "1"
	url := "https://api.eflow.team/v1/networks/offers/" + offersId //按 ID 查询会员详情

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
func (s *campaignService) updateNetworksOffers() {
	var offersId = "1"
	url := "https://api.eflow.team/v1/networks/offers/" + offersId //更新联盟会员

	var test []string
	var b = false
	var i int32 = 58

	PayoutRevenueEntries := offer.EntriesInfo{
		{
			EntryName:                      "Base",
			PayoutType:                     "cpa_cps",
			PayoutAmount:                   2,
			PayoutPercentage:               5,
			RevenueType:                    "rpa_rps",
			RevenueAmount:                  5,
			RevenuePercentage:              10,
			IsDefault:                      true,
			IsEmailAttributionDefaultEvent: false,
		},
	}
	everflowReq := offer.OfferRequest{
		OfferStatus:                   "active",
		Visibility:                    "public",
		RequirementTrackingParameters: test,
		PayoutRevenue:                 PayoutRevenueEntries,
		EmailAttributionMethod:        "first_affiliate_attribution",
		RedirectMode:                  "standard",
		ConversionMethod:              "server_postback",
		AttributionMethod:             "last_touch",
		InternalRedirects:             test,
		TrafficFilters:                test,
		RequirementKpis:               test,
		SourceNames:                   test,
		Integrations:                  offer.Integrations{},
		EmailOptout:                   offer.EmailOptoutSettings{},
		Labels:                        offer.Details{},
		Channels:                      offer.Details{},
		Ruleset: offer.Ruleset{
			Platforms:            []map[string]interface{}{},
			DeviceTypes:          []map[string]interface{}{},
			OsVersions:           []map[string]interface{}{},
			Browsers:             []map[string]interface{}{},
			Brands:               []map[string]interface{}{},
			Languages:            []map[string]interface{}{},
			Countries:            []map[string]interface{}{},
			Regions:              []map[string]interface{}{},
			Cities:               []map[string]interface{}{},
			Dmas:                 []map[string]interface{}{},
			Isps:                 []map[string]interface{}{},
			MobileCarriers:       []map[string]interface{}{},
			ConnectionTypes:      []map[string]interface{}{},
			Ips:                  []map[string]interface{}{},
			IsUseDayParting:      &b,
			IsBlockProxy:         &b,
			DayPartingTimezoneId: &i,
			DaysParting:          []map[string]interface{}{},
			PostalCodes:          []map[string]interface{}{},
		},
		Email:                             test,
		Creatives:                         []interface{}{},
		IsSoftCap:                         false,
		NetworkTrackingDomainId:           12977,
		IsUseSecureLink:                   true,
		CurrencyId:                        "USD",
		SessionDuration:                   24,
		SessionDefinition:                 "cookie",
		NetworkCategoryId:                 1,
		NetworkAdvertiserId:               1,
		Name:                              "test offer 222",
		DestinationUrl:                    "https://test.com?transaction_id={transaction_id}",
		AppIdentifier:                     "test_app_identifier",
		PreviewUrl:                        "https://test.com",
		InternalNotes:                     "internal_note",
		DateLiveUntil:                     "2025-10-18",
		HtmlDescription:                   "This is a test offer\nThis is a test description",
		IsDescriptionPlainText:            false,
		IsUseDirectLinking:                false,
		IsAllowDeepLink:                   false,
		IsUsingExplicitTermsAndConditions: true,
		IsForceTermsAndConditions:         false,
		TermsAndConditions:                "This is a test terms and conditions",
		NetworkOfferGroupId:               0,
		NetworkApplicationQuestionnaireId: 0,
		CapsTimezoneId:                    0,
		ServerSideUrl:                     "",
		IsEmailAttributionWindowEnabled:   false,
		SuppressionListId:                 0,
		ThumbnailUrl:                      "",
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
