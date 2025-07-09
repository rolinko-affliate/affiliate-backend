package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/repository"
)

// AnalyticsService defines the interface for analytics business logic
type AnalyticsService interface {
	// Autocompletion
	SearchOrganizations(ctx context.Context, query string, orgType string, limit int) ([]domain.AutocompleteResult, error)

	// Advertiser methods
	GetAdvertiserByID(ctx context.Context, id int64) (*domain.AnalyticsAdvertiserResponse, error)
	CreateAdvertiser(ctx context.Context, advertiser *domain.AnalyticsAdvertiser) error
	UpdateAdvertiser(ctx context.Context, advertiser *domain.AnalyticsAdvertiser) error
	DeleteAdvertiser(ctx context.Context, id int64) error

	// Publisher methods  
	GetPublisherByID(ctx context.Context, id int64) (*domain.AnalyticsPublisherResponse, error)
	CreatePublisher(ctx context.Context, publisher *domain.AnalyticsPublisher) error
	UpdatePublisher(ctx context.Context, publisher *domain.AnalyticsPublisher) error
	DeletePublisher(ctx context.Context, id int64) error
}

// analyticsService implements AnalyticsService
type analyticsService struct {
	analyticsRepo repository.AnalyticsRepository
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(analyticsRepo repository.AnalyticsRepository) AnalyticsService {
	return &analyticsService{
		analyticsRepo: analyticsRepo,
	}
}

// SearchOrganizations performs autocompletion search
func (s *analyticsService) SearchOrganizations(ctx context.Context, query string, orgType string, limit int) ([]domain.AutocompleteResult, error) {
	if len(query) < 3 {
		return nil, fmt.Errorf("query must be at least 3 characters long")
	}

	if limit <= 0 || limit > 50 {
		limit = 10 // Default limit
	}

	switch orgType {
	case "advertiser":
		return s.analyticsRepo.SearchAdvertisers(ctx, query, limit)
	case "publisher":
		return s.analyticsRepo.SearchPublishers(ctx, query, limit)
	case "both", "":
		return s.analyticsRepo.SearchBoth(ctx, query, limit)
	default:
		return nil, fmt.Errorf("invalid organization type: %s", orgType)
	}
}

// GetAdvertiserByID retrieves advertiser data and formats it for API response
func (s *analyticsService) GetAdvertiserByID(ctx context.Context, id int64) (*domain.AnalyticsAdvertiserResponse, error) {
	advertiser, err := s.analyticsRepo.GetAdvertiserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.buildAdvertiserResponse(advertiser)
}

// GetPublisherByID retrieves publisher data and formats it for API response
func (s *analyticsService) GetPublisherByID(ctx context.Context, id int64) (*domain.AnalyticsPublisherResponse, error) {
	publisher, err := s.analyticsRepo.GetPublisherByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.buildPublisherResponse(publisher)
}

// CreateAdvertiser creates a new advertiser
func (s *analyticsService) CreateAdvertiser(ctx context.Context, advertiser *domain.AnalyticsAdvertiser) error {
	return s.analyticsRepo.CreateAdvertiser(ctx, advertiser)
}

// UpdateAdvertiser updates an existing advertiser
func (s *analyticsService) UpdateAdvertiser(ctx context.Context, advertiser *domain.AnalyticsAdvertiser) error {
	return s.analyticsRepo.UpdateAdvertiser(ctx, advertiser)
}

// DeleteAdvertiser deletes an advertiser
func (s *analyticsService) DeleteAdvertiser(ctx context.Context, id int64) error {
	return s.analyticsRepo.DeleteAdvertiser(ctx, id)
}

// CreatePublisher creates a new publisher
func (s *analyticsService) CreatePublisher(ctx context.Context, publisher *domain.AnalyticsPublisher) error {
	return s.analyticsRepo.CreatePublisher(ctx, publisher)
}

// UpdatePublisher updates an existing publisher
func (s *analyticsService) UpdatePublisher(ctx context.Context, publisher *domain.AnalyticsPublisher) error {
	return s.analyticsRepo.UpdatePublisher(ctx, publisher)
}

// DeletePublisher deletes a publisher
func (s *analyticsService) DeletePublisher(ctx context.Context, id int64) error {
	return s.analyticsRepo.DeletePublisher(ctx, id)
}

// buildAdvertiserResponse constructs the API response for advertiser data
func (s *analyticsService) buildAdvertiserResponse(advertiser *domain.AnalyticsAdvertiser) (*domain.AnalyticsAdvertiserResponse, error) {
	response := &domain.AnalyticsAdvertiserResponse{}
	response.Advertiser.Domain = advertiser.Domain

	// Build metadata
	if advertiser.Description != nil || advertiser.FaviconImageURL != nil || advertiser.ScreenshotImageURL != nil {
		response.Advertiser.MetaData = &domain.MetaData{}
		if advertiser.Description != nil {
			response.Advertiser.MetaData.Description = *advertiser.Description
		}
		if advertiser.FaviconImageURL != nil {
			response.Advertiser.MetaData.FaviconImageURL = *advertiser.FaviconImageURL
		}
		if advertiser.ScreenshotImageURL != nil {
			response.Advertiser.MetaData.ScreenshotImageURL = *advertiser.ScreenshotImageURL
		}
	}

	// Parse and set affiliate networks
	if affiliateNetworks, err := advertiser.GetAffiliateNetworks(); err == nil && affiliateNetworks != nil {
		response.Advertiser.AffiliateNetworks = affiliateNetworks
	}

	// Parse and set contact emails
	if contactEmails, err := advertiser.GetContactEmails(); err == nil && contactEmails != nil {
		response.Advertiser.ContactEmails = contactEmails
	}

	// Parse and set keywords
	if keywords, err := advertiser.GetKeywords(); err == nil && keywords != nil {
		response.Advertiser.Keywords = keywords
	}

	// Parse and set verticals
	if verticals, err := advertiser.GetVerticals(); err == nil && verticals != nil {
		response.Advertiser.Verticals = verticals
	}

	// Parse and set social media
	if socialMedia, err := advertiser.GetSocialMedia(); err == nil && socialMedia != nil {
		response.Advertiser.SocialMedia = socialMedia
	}

	// Parse and set other complex fields
	if advertiser.PartnerInformation != nil {
		var partnerInfo interface{}
		if err := json.Unmarshal([]byte(*advertiser.PartnerInformation), &partnerInfo); err == nil {
			response.Advertiser.PartnerInformation = partnerInfo
		}
	}

	if advertiser.RelatedAdvertisers != nil {
		var relatedAdvs interface{}
		if err := json.Unmarshal([]byte(*advertiser.RelatedAdvertisers), &relatedAdvs); err == nil {
			response.Advertiser.RelatedAdvertisers = relatedAdvs
		}
	}

	if advertiser.Backlinks != nil {
		var backlinks interface{}
		if err := json.Unmarshal([]byte(*advertiser.Backlinks), &backlinks); err == nil {
			response.Advertiser.Backlinks = backlinks
		}
	}

	return response, nil
}

// buildPublisherResponse constructs the API response for publisher data
func (s *analyticsService) buildPublisherResponse(publisher *domain.AnalyticsPublisher) (*domain.AnalyticsPublisherResponse, error) {
	response := &domain.AnalyticsPublisherResponse{}
	response.Publisher.Domain = publisher.Domain
	response.Publisher.Relevance = publisher.Relevance
	response.Publisher.TrafficScore = publisher.TrafficScore

	// Build metadata
	if publisher.Description != nil || publisher.FaviconImageURL != nil || publisher.ScreenshotImageURL != nil {
		response.Publisher.MetaData = &domain.MetaData{}
		if publisher.Description != nil {
			response.Publisher.MetaData.Description = *publisher.Description
		}
		if publisher.FaviconImageURL != nil {
			response.Publisher.MetaData.FaviconImageURL = *publisher.FaviconImageURL
		}
		if publisher.ScreenshotImageURL != nil {
			response.Publisher.MetaData.ScreenshotImageURL = *publisher.ScreenshotImageURL
		}
	}

	// Set known flag
	response.Publisher.Known = &struct {
		Value bool `json:"value"`
	}{Value: publisher.Known}

	// Set promotype
	if publisher.Promotype != nil {
		response.Publisher.Promotype = &struct {
			Value *string `json:"value"`
		}{Value: publisher.Promotype}
	}

	// Parse and set affiliate networks
	if affiliateNetworks, err := publisher.GetAffiliateNetworks(); err == nil && affiliateNetworks != nil {
		response.Publisher.AffiliateNetworks = affiliateNetworks
	}

	// Parse and set country rankings
	if countryRankings, err := publisher.GetCountryRankings(); err == nil && countryRankings != nil {
		response.Publisher.CountryRankings = countryRankings
	}

	// Parse and set keywords
	if keywords, err := publisher.GetKeywords(); err == nil && keywords != nil {
		response.Publisher.Keywords = keywords
	}

	// Parse and set enhanced verticals
	if verticalsV2, err := publisher.GetVerticalsV2(); err == nil && verticalsV2 != nil {
		response.Publisher.VerticalsV2 = verticalsV2
	}

	// Parse and set social media
	if socialMedia, err := publisher.GetSocialMedia(); err == nil && socialMedia != nil {
		response.Publisher.SocialMedia = socialMedia
	}

	// Parse and set other complex fields
	if publisher.Verticals != nil {
		var verticals interface{}
		if err := json.Unmarshal([]byte(*publisher.Verticals), &verticals); err == nil {
			response.Publisher.Verticals = verticals
		}
	}

	if publisher.PartnerInformation != nil {
		var partnerInfo interface{}
		if err := json.Unmarshal([]byte(*publisher.PartnerInformation), &partnerInfo); err == nil {
			response.Publisher.PartnerInformation = partnerInfo
		}
	}

	if publisher.Partners != nil {
		var partners interface{}
		if err := json.Unmarshal([]byte(*publisher.Partners), &partners); err == nil {
			response.Publisher.Partners = partners
		}
	}

	if publisher.RelatedPublishers != nil {
		var relatedPubs interface{}
		if err := json.Unmarshal([]byte(*publisher.RelatedPublishers), &relatedPubs); err == nil {
			response.Publisher.RelatedPublishers = relatedPubs
		}
	}

	if publisher.LiveURLs != nil {
		var liveURLs interface{}
		if err := json.Unmarshal([]byte(*publisher.LiveURLs), &liveURLs); err == nil {
			response.Publisher.LiveURLs = liveURLs
		}
	}

	return response, nil
}