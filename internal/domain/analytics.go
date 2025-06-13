package domain

import (
	"encoding/json"
	"time"
)

// AnalyticsAdvertiser represents advertiser analytics data
type AnalyticsAdvertiser struct {
	ID                   int64     `json:"id" db:"id"`
	Domain               string    `json:"domain" db:"domain"`
	Description          *string   `json:"description,omitempty" db:"description"`
	FaviconImageURL      *string   `json:"favicon_image_url,omitempty" db:"favicon_image_url"`
	ScreenshotImageURL   *string   `json:"screenshot_image_url,omitempty" db:"screenshot_image_url"`
	AffiliateNetworks    *string   `json:"-" db:"affiliate_networks"` // JSONB stored as string
	ContactEmails        *string   `json:"-" db:"contact_emails"`     // JSONB stored as string
	Keywords             *string   `json:"-" db:"keywords"`           // JSONB stored as string
	Verticals            *string   `json:"-" db:"verticals"`          // JSONB stored as string
	PartnerInformation   *string   `json:"-" db:"partner_information"` // JSONB stored as string
	RelatedAdvertisers   *string   `json:"-" db:"related_advertisers"` // JSONB stored as string
	SocialMedia          *string   `json:"-" db:"social_media"`       // JSONB stored as string
	Backlinks            *string   `json:"-" db:"backlinks"`          // JSONB stored as string
	AdditionalData       *string   `json:"-" db:"additional_data"`    // JSONB stored as string
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// AnalyticsPublisher represents publisher analytics data
type AnalyticsPublisher struct {
	ID                   int64     `json:"id" db:"id"`
	Domain               string    `json:"domain" db:"domain"`
	Description          *string   `json:"description,omitempty" db:"description"`
	FaviconImageURL      *string   `json:"favicon_image_url,omitempty" db:"favicon_image_url"`
	ScreenshotImageURL   *string   `json:"screenshot_image_url,omitempty" db:"screenshot_image_url"`
	AffiliateNetworks    *string   `json:"-" db:"affiliate_networks"`    // JSONB stored as string
	CountryRankings      *string   `json:"-" db:"country_rankings"`      // JSONB stored as string
	Keywords             *string   `json:"-" db:"keywords"`              // JSONB stored as string
	Verticals            *string   `json:"-" db:"verticals"`             // JSONB stored as string
	VerticalsV2          *string   `json:"-" db:"verticals_v2"`          // JSONB stored as string
	PartnerInformation   *string   `json:"-" db:"partner_information"`   // JSONB stored as string
	Partners             *string   `json:"-" db:"partners"`              // JSONB stored as string
	RelatedPublishers    *string   `json:"-" db:"related_publishers"`    // JSONB stored as string
	SocialMedia          *string   `json:"-" db:"social_media"`          // JSONB stored as string
	LiveURLs             *string   `json:"-" db:"live_urls"`             // JSONB stored as string
	Known                bool      `json:"known" db:"known"`
	Relevance            float64   `json:"relevance" db:"relevance"`
	TrafficScore         float64   `json:"traffic_score" db:"traffic_score"`
	Promotype            *string   `json:"promotype,omitempty" db:"promotype"`
	AdditionalData       *string   `json:"-" db:"additional_data"`       // JSONB stored as string
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// AutocompleteResult represents a search result for autocompletion
type AutocompleteResult struct {
	ID     int64  `json:"id"`
	Domain string `json:"domain"`
	Type   string `json:"type"` // "advertiser" or "publisher"
	Name   string `json:"name"` // Display name (domain for now)
}

// Structured data types for API responses

// AffiliateNetworkData represents affiliate network information
type AffiliateNetworkData struct {
	Count       int      `json:"count"`
	SampleValue []string `json:"sampleValue,omitempty"`
	Value       []string `json:"value"`
}

// ContactEmailData represents contact email information
type ContactEmailData struct {
	Count int `json:"count"`
	Error bool `json:"error"`
	Value []struct {
		Department *string `json:"department"`
		Value      string  `json:"value"`
	} `json:"value"`
}

// KeywordData represents keyword information
type KeywordData struct {
	Count       int `json:"count"`
	SampleValue []struct {
		Score float64 `json:"score"`
		Value string  `json:"value"`
	} `json:"sampleValue,omitempty"`
	Value []struct {
		Score float64 `json:"score"`
		Value string  `json:"value"`
	} `json:"value"`
}

// VerticalData represents vertical/industry information
type VerticalData struct {
	Count       int `json:"count"`
	SampleValue *struct {
		Name  string `json:"name"`
		Rank  int    `json:"rank"`
		Score int    `json:"score"`
	} `json:"sampleValue,omitempty"`
	Value []struct {
		Name  string `json:"name"`
		Rank  int    `json:"rank"`
		Score int    `json:"score"`
	} `json:"value"`
}

// MetaData represents metadata information
type MetaData struct {
	Description        string `json:"description"`
	FaviconImageURL    string `json:"faviconImageUrl"`
	ScreenshotImageURL string `json:"screenshotImageUrl"`
}

// SocialMediaData represents social media information
type SocialMediaData struct {
	Count             int                `json:"count"`
	SocialsAvailable  []string           `json:"socialsAvailable"`
	Value             map[string]string  `json:"value"`
}

// CountryRankingData represents country ranking information (for publishers)
type CountryRankingData struct {
	Count        int `json:"count"`
	HighestValue *struct {
		CountryCode string  `json:"countryCode"`
		Score       float64 `json:"score"`
	} `json:"highestValue,omitempty"`
	SampleValue []struct {
		CountryCode string  `json:"countryCode"`
		Score       float64 `json:"score"`
	} `json:"sampleValue,omitempty"`
	Value []struct {
		CountryCode string  `json:"countryCode"`
		Score       float64 `json:"score"`
	} `json:"value"`
}

// AnalyticsAdvertiserResponse represents the API response for advertiser data
type AnalyticsAdvertiserResponse struct {
	Advertiser struct {
		AffiliateNetworks    *AffiliateNetworkData `json:"affiliateNetworks,omitempty"`
		Backlinks            interface{}           `json:"backlinks,omitempty"`
		ContactEmails        *ContactEmailData     `json:"contactEmails,omitempty"`
		Domain               string                `json:"domain"`
		Keywords             *KeywordData          `json:"keywords,omitempty"`
		MetaData             *MetaData             `json:"metaData,omitempty"`
		PartnerInformation   interface{}           `json:"partnerInformation,omitempty"`
		RelatedAdvertisers   interface{}           `json:"relatedAdvertisers,omitempty"`
		SocialMedia          *SocialMediaData      `json:"socialMedia,omitempty"`
		Verticals            *VerticalData         `json:"verticals,omitempty"`
	} `json:"advertiser"`
}

// AnalyticsPublisherResponse represents the API response for publisher data
type AnalyticsPublisherResponse struct {
	Publisher struct {
		AffiliateNetworks    *AffiliateNetworkData  `json:"affiliateNetworks,omitempty"`
		CountryRankings      *CountryRankingData    `json:"countryRankings,omitempty"`
		Domain               string                 `json:"domain"`
		Keywords             *KeywordData           `json:"keywords,omitempty"`
		Known                *struct {
			Value bool `json:"value"`
		} `json:"known,omitempty"`
		LiveURLs             interface{}            `json:"liveUrls,omitempty"`
		MetaData             *MetaData              `json:"metaData,omitempty"`
		PartnerInformation   interface{}            `json:"partnerInformation,omitempty"`
		Partners             interface{}            `json:"partners,omitempty"`
		Promotype            *struct {
			Value *string `json:"value"`
		} `json:"promotype,omitempty"`
		RelatedPublishers    interface{}            `json:"relatedPublishers,omitempty"`
		Relevance            float64                `json:"relevance"`
		SocialMedia          *SocialMediaData       `json:"socialMedia,omitempty"`
		TrafficScore         float64                `json:"trafficScore"`
		Verticals            interface{}            `json:"verticals,omitempty"`
		VerticalsV2          *VerticalData          `json:"verticalsV2,omitempty"`
	} `json:"publisher"`
}

// Helper methods for AnalyticsAdvertiser

// GetAffiliateNetworks parses and returns affiliate networks data
func (a *AnalyticsAdvertiser) GetAffiliateNetworks() (*AffiliateNetworkData, error) {
	if a.AffiliateNetworks == nil {
		return nil, nil
	}
	var data AffiliateNetworkData
	err := json.Unmarshal([]byte(*a.AffiliateNetworks), &data)
	return &data, err
}

// GetContactEmails parses and returns contact emails data
func (a *AnalyticsAdvertiser) GetContactEmails() (*ContactEmailData, error) {
	if a.ContactEmails == nil {
		return nil, nil
	}
	var data ContactEmailData
	err := json.Unmarshal([]byte(*a.ContactEmails), &data)
	return &data, err
}

// GetKeywords parses and returns keywords data
func (a *AnalyticsAdvertiser) GetKeywords() (*KeywordData, error) {
	if a.Keywords == nil {
		return nil, nil
	}
	var data KeywordData
	err := json.Unmarshal([]byte(*a.Keywords), &data)
	return &data, err
}

// GetVerticals parses and returns verticals data
func (a *AnalyticsAdvertiser) GetVerticals() (*VerticalData, error) {
	if a.Verticals == nil {
		return nil, nil
	}
	var data VerticalData
	err := json.Unmarshal([]byte(*a.Verticals), &data)
	return &data, err
}

// GetSocialMedia parses and returns social media data
func (a *AnalyticsAdvertiser) GetSocialMedia() (*SocialMediaData, error) {
	if a.SocialMedia == nil {
		return nil, nil
	}
	var data SocialMediaData
	err := json.Unmarshal([]byte(*a.SocialMedia), &data)
	return &data, err
}

// Helper methods for AnalyticsPublisher

// GetAffiliateNetworks parses and returns affiliate networks data
func (p *AnalyticsPublisher) GetAffiliateNetworks() (*AffiliateNetworkData, error) {
	if p.AffiliateNetworks == nil {
		return nil, nil
	}
	var data AffiliateNetworkData
	err := json.Unmarshal([]byte(*p.AffiliateNetworks), &data)
	return &data, err
}

// GetCountryRankings parses and returns country rankings data
func (p *AnalyticsPublisher) GetCountryRankings() (*CountryRankingData, error) {
	if p.CountryRankings == nil {
		return nil, nil
	}
	var data CountryRankingData
	err := json.Unmarshal([]byte(*p.CountryRankings), &data)
	return &data, err
}

// GetKeywords parses and returns keywords data
func (p *AnalyticsPublisher) GetKeywords() (*KeywordData, error) {
	if p.Keywords == nil {
		return nil, nil
	}
	var data KeywordData
	err := json.Unmarshal([]byte(*p.Keywords), &data)
	return &data, err
}

// GetVerticalsV2 parses and returns enhanced verticals data
func (p *AnalyticsPublisher) GetVerticalsV2() (*VerticalData, error) {
	if p.VerticalsV2 == nil {
		return nil, nil
	}
	var data VerticalData
	err := json.Unmarshal([]byte(*p.VerticalsV2), &data)
	return &data, err
}

// GetSocialMedia parses and returns social media data
func (p *AnalyticsPublisher) GetSocialMedia() (*SocialMediaData, error) {
	if p.SocialMedia == nil {
		return nil, nil
	}
	var data SocialMediaData
	err := json.Unmarshal([]byte(*p.SocialMedia), &data)
	return &data, err
}