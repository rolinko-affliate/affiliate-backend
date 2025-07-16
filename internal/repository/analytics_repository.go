package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/affiliate-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AffiliatesSearchResult represents the result of affiliate search with total count
type AffiliatesSearchResult struct {
	Data  []*domain.AnalyticsPublisher
	Total int64
}

// AnalyticsRepository defines the interface for analytics data access
type AnalyticsRepository interface {
	// Advertiser methods
	CreateAdvertiser(ctx context.Context, advertiser *domain.AnalyticsAdvertiser) error
	GetAdvertiserByID(ctx context.Context, id int64) (*domain.AnalyticsAdvertiser, error)
	GetAdvertiserByDomain(ctx context.Context, domainName string) (*domain.AnalyticsAdvertiser, error)
	UpdateAdvertiser(ctx context.Context, advertiser *domain.AnalyticsAdvertiser) error
	DeleteAdvertiser(ctx context.Context, id int64) error

	// Publisher methods
	CreatePublisher(ctx context.Context, publisher *domain.AnalyticsPublisher) error
	GetPublisherByID(ctx context.Context, id int64) (*domain.AnalyticsPublisher, error)
	GetPublisherByDomain(ctx context.Context, domainName string) (*domain.AnalyticsPublisher, error)
	UpdatePublisher(ctx context.Context, publisher *domain.AnalyticsPublisher) error
	DeletePublisher(ctx context.Context, id int64) error

	// Search methods
	SearchAdvertisers(ctx context.Context, query string, limit int) ([]domain.AutocompleteResult, error)
	SearchPublishers(ctx context.Context, query string, limit int) ([]domain.AutocompleteResult, error)
	SearchBoth(ctx context.Context, query string, limit int) ([]domain.AutocompleteResult, error)
	AffiliatesSearch(ctx context.Context, domainFilter, country string, partnerDomains []string, verticals []string, limit int, offset int) (*AffiliatesSearchResult, error)
}

// analyticsRepository implements AnalyticsRepository
type analyticsRepository struct {
	db *pgxpool.Pool
}

func (r *analyticsRepository) AffiliatesSearch(ctx context.Context, domainFilter, country string, partnerDomains []string, verticals []string, limit int, offset int) (*AffiliatesSearchResult, error) {
	// Build WHERE conditions
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Domain filter - partial matching like auto-completion
	if domainFilter != "" {
		conditions = append(conditions, fmt.Sprintf(`domain ILIKE $%d`, argIndex))
		args = append(args, fmt.Sprintf("%%%s%%", domainFilter))
		argIndex++
	}

	// Country filter - search in country_rankings JSONB
	if country != "" {
		conditions = append(conditions, fmt.Sprintf(`
			country_rankings->'value' @> $%d`, argIndex))
		args = append(args, fmt.Sprintf(`[{"countryCode": "%s"}]`, strings.ToLower(country)))
		argIndex++
	}

	// Partner domains filter - search in partners JSONB array
	if len(partnerDomains) > 0 {
		partnerConditions := make([]string, len(partnerDomains))
		for i, partnerDomain := range partnerDomains {
			partnerConditions[i] = fmt.Sprintf(`partners->'value' @> $%d`, argIndex)
			args = append(args, fmt.Sprintf(`["%s"]`, partnerDomain))
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(partnerConditions, " OR ")))
	}

	// Verticals filter - search in verticalsV2 JSONB array
	if len(verticals) > 0 {
		verticalConditions := make([]string, len(verticals))
		for i, vertical := range verticals {
			verticalConditions[i] = fmt.Sprintf(`verticals_v2->'value' @> $%d`, argIndex)
			args = append(args, fmt.Sprintf(`[{"name": "%s"}]`, vertical))
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(verticalConditions, " OR ")))
	}

	// Build WHERE clause
	whereClause := "WHERE 1=1"
	if len(conditions) > 0 {
		whereClause += " AND " + strings.Join(conditions, " AND ")
	}

	// First query: Get total count
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM analytics_publishers %s`, whereClause)

	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error counting affiliates: %w", err)
	}

	// Second query: Get paginated data
	dataQuery := fmt.Sprintf(`
        SELECT 
            id, domain, description, favicon_image_url, screenshot_image_url,
            affiliate_networks, country_rankings, keywords, verticals, verticals_v2,
            partner_information, partners, related_publishers, social_media, live_urls,
            known, relevance, traffic_score, promotype, additional_data,
            created_at, updated_at
        FROM analytics_publishers 
        %s
        ORDER BY 
            COALESCE((partners->'count')::int, 0) DESC
        LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	// Add pagination args
	dataArgs := append(args, limit, offset)

	rows, err := r.db.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, fmt.Errorf("error searching affiliates: %w", err)
	}
	defer rows.Close()

	// Initialize as empty slice to ensure we never return null
	publishers := make([]*domain.AnalyticsPublisher, 0)

	for rows.Next() {
		var p domain.AnalyticsPublisher

		err := rows.Scan(
			&p.ID, &p.Domain,
			&p.Description, &p.FaviconImageURL, &p.ScreenshotImageURL,
			&p.AffiliateNetworks, &p.CountryRankings, &p.Keywords, &p.Verticals, &p.VerticalsV2,
			&p.PartnerInformation, &p.Partners, &p.RelatedPublishers, &p.SocialMedia, &p.LiveURLs,
			&p.Known, &p.Relevance, &p.TrafficScore, &p.Promotype, &p.AdditionalData,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning publisher row: %w", err)
		}

		publishers = append(publishers, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return &AffiliatesSearchResult{
		Data:  publishers,
		Total: total,
	}, nil
}

// NewAnalyticsRepository creates a new analytics repository
func NewAnalyticsRepository(db *pgxpool.Pool) AnalyticsRepository {
	return &analyticsRepository{db: db}
}

// Advertiser methods

func (r *analyticsRepository) CreateAdvertiser(ctx context.Context, advertiser *domain.AnalyticsAdvertiser) error {
	query := `
		INSERT INTO analytics_advertisers (
			domain, description, favicon_image_url, screenshot_image_url,
			affiliate_networks, contact_emails, keywords, verticals,
			partner_information, related_advertisers, social_media, backlinks, additional_data
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		advertiser.Domain, advertiser.Description, advertiser.FaviconImageURL, advertiser.ScreenshotImageURL,
		advertiser.AffiliateNetworks, advertiser.ContactEmails, advertiser.Keywords, advertiser.Verticals,
		advertiser.PartnerInformation, advertiser.RelatedAdvertisers, advertiser.SocialMedia,
		advertiser.Backlinks, advertiser.AdditionalData,
	).Scan(&advertiser.ID, &advertiser.CreatedAt, &advertiser.UpdatedAt)

	return err
}

func (r *analyticsRepository) GetAdvertiserByID(ctx context.Context, id int64) (*domain.AnalyticsAdvertiser, error) {
	query := `
		SELECT id, domain, description, favicon_image_url, screenshot_image_url,
			   affiliate_networks, contact_emails, keywords, verticals,
			   partner_information, related_advertisers, social_media, backlinks, additional_data,
			   created_at, updated_at
		FROM analytics_advertisers 
		WHERE id = $1`

	var advertiser domain.AnalyticsAdvertiser
	err := r.db.QueryRow(ctx, query, id).Scan(
		&advertiser.ID, &advertiser.Domain, &advertiser.Description, &advertiser.FaviconImageURL, &advertiser.ScreenshotImageURL,
		&advertiser.AffiliateNetworks, &advertiser.ContactEmails, &advertiser.Keywords, &advertiser.Verticals,
		&advertiser.PartnerInformation, &advertiser.RelatedAdvertisers, &advertiser.SocialMedia,
		&advertiser.Backlinks, &advertiser.AdditionalData,
		&advertiser.CreatedAt, &advertiser.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &advertiser, nil
}

func (r *analyticsRepository) GetAdvertiserByDomain(ctx context.Context, domainName string) (*domain.AnalyticsAdvertiser, error) {
	query := `
		SELECT id, domain, description, favicon_image_url, screenshot_image_url,
			   affiliate_networks, contact_emails, keywords, verticals,
			   partner_information, related_advertisers, social_media, backlinks, additional_data,
			   created_at, updated_at
		FROM analytics_advertisers 
		WHERE domain = $1`

	var advertiser domain.AnalyticsAdvertiser
	err := r.db.QueryRow(ctx, query, domainName).Scan(
		&advertiser.ID, &advertiser.Domain, &advertiser.Description, &advertiser.FaviconImageURL, &advertiser.ScreenshotImageURL,
		&advertiser.AffiliateNetworks, &advertiser.ContactEmails, &advertiser.Keywords, &advertiser.Verticals,
		&advertiser.PartnerInformation, &advertiser.RelatedAdvertisers, &advertiser.SocialMedia,
		&advertiser.Backlinks, &advertiser.AdditionalData,
		&advertiser.CreatedAt, &advertiser.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &advertiser, nil
}

func (r *analyticsRepository) UpdateAdvertiser(ctx context.Context, advertiser *domain.AnalyticsAdvertiser) error {
	query := `
		UPDATE analytics_advertisers SET
			domain = $2, description = $3, favicon_image_url = $4, screenshot_image_url = $5,
			affiliate_networks = $6, contact_emails = $7, keywords = $8, verticals = $9,
			partner_information = $10, related_advertisers = $11, social_media = $12, 
			backlinks = $13, additional_data = $14, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query,
		advertiser.ID, advertiser.Domain, advertiser.Description, advertiser.FaviconImageURL, advertiser.ScreenshotImageURL,
		advertiser.AffiliateNetworks, advertiser.ContactEmails, advertiser.Keywords, advertiser.Verticals,
		advertiser.PartnerInformation, advertiser.RelatedAdvertisers, advertiser.SocialMedia,
		advertiser.Backlinks, advertiser.AdditionalData,
	).Scan(&advertiser.UpdatedAt)

	return err
}

func (r *analyticsRepository) DeleteAdvertiser(ctx context.Context, id int64) error {
	query := `DELETE FROM analytics_advertisers WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Publisher methods

func (r *analyticsRepository) CreatePublisher(ctx context.Context, publisher *domain.AnalyticsPublisher) error {
	query := `
		INSERT INTO analytics_publishers (
			domain, description, favicon_image_url, screenshot_image_url,
			affiliate_networks, country_rankings, keywords, verticals, verticals_v2,
			partner_information, partners, related_publishers, social_media, live_urls,
			known, relevance, traffic_score, promotype, additional_data
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		publisher.Domain, publisher.Description, publisher.FaviconImageURL, publisher.ScreenshotImageURL,
		publisher.AffiliateNetworks, publisher.CountryRankings, publisher.Keywords, publisher.Verticals, publisher.VerticalsV2,
		publisher.PartnerInformation, publisher.Partners, publisher.RelatedPublishers, publisher.SocialMedia, publisher.LiveURLs,
		publisher.Known, publisher.Relevance, publisher.TrafficScore, publisher.Promotype, publisher.AdditionalData,
	).Scan(&publisher.ID, &publisher.CreatedAt, &publisher.UpdatedAt)

	return err
}

func (r *analyticsRepository) GetPublisherByID(ctx context.Context, id int64) (*domain.AnalyticsPublisher, error) {
	query := `
		SELECT id, domain, description, favicon_image_url, screenshot_image_url,
			   affiliate_networks, country_rankings, keywords, verticals, verticals_v2,
			   partner_information, partners, related_publishers, social_media, live_urls,
			   known, relevance, traffic_score, promotype, additional_data,
			   created_at, updated_at
		FROM analytics_publishers 
		WHERE id = $1`

	var publisher domain.AnalyticsPublisher
	err := r.db.QueryRow(ctx, query, id).Scan(
		&publisher.ID, &publisher.Domain, &publisher.Description, &publisher.FaviconImageURL, &publisher.ScreenshotImageURL,
		&publisher.AffiliateNetworks, &publisher.CountryRankings, &publisher.Keywords, &publisher.Verticals, &publisher.VerticalsV2,
		&publisher.PartnerInformation, &publisher.Partners, &publisher.RelatedPublishers, &publisher.SocialMedia, &publisher.LiveURLs,
		&publisher.Known, &publisher.Relevance, &publisher.TrafficScore, &publisher.Promotype, &publisher.AdditionalData,
		&publisher.CreatedAt, &publisher.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &publisher, nil
}

func (r *analyticsRepository) GetPublisherByDomain(ctx context.Context, domainName string) (*domain.AnalyticsPublisher, error) {
	query := `
		SELECT id, domain, description, favicon_image_url, screenshot_image_url,
			   affiliate_networks, country_rankings, keywords, verticals, verticals_v2,
			   partner_information, partners, related_publishers, social_media, live_urls,
			   known, relevance, traffic_score, promotype, additional_data,
			   created_at, updated_at
		FROM analytics_publishers 
		WHERE domain = $1`

	var publisher domain.AnalyticsPublisher
	err := r.db.QueryRow(ctx, query, domainName).Scan(
		&publisher.ID, &publisher.Domain, &publisher.Description, &publisher.FaviconImageURL, &publisher.ScreenshotImageURL,
		&publisher.AffiliateNetworks, &publisher.CountryRankings, &publisher.Keywords, &publisher.Verticals, &publisher.VerticalsV2,
		&publisher.PartnerInformation, &publisher.Partners, &publisher.RelatedPublishers, &publisher.SocialMedia, &publisher.LiveURLs,
		&publisher.Known, &publisher.Relevance, &publisher.TrafficScore, &publisher.Promotype, &publisher.AdditionalData,
		&publisher.CreatedAt, &publisher.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &publisher, nil
}

func (r *analyticsRepository) UpdatePublisher(ctx context.Context, publisher *domain.AnalyticsPublisher) error {
	query := `
		UPDATE analytics_publishers SET
			domain = $2, description = $3, favicon_image_url = $4, screenshot_image_url = $5,
			affiliate_networks = $6, country_rankings = $7, keywords = $8, verticals = $9, verticals_v2 = $10,
			partner_information = $11, partners = $12, related_publishers = $13, social_media = $14, live_urls = $15,
			known = $16, relevance = $17, traffic_score = $18, promotype = $19, additional_data = $20,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query,
		publisher.ID, publisher.Domain, publisher.Description, publisher.FaviconImageURL, publisher.ScreenshotImageURL,
		publisher.AffiliateNetworks, publisher.CountryRankings, publisher.Keywords, publisher.Verticals, publisher.VerticalsV2,
		publisher.PartnerInformation, publisher.Partners, publisher.RelatedPublishers, publisher.SocialMedia, publisher.LiveURLs,
		publisher.Known, publisher.Relevance, publisher.TrafficScore, publisher.Promotype, publisher.AdditionalData,
	).Scan(&publisher.UpdatedAt)

	return err
}

func (r *analyticsRepository) DeletePublisher(ctx context.Context, id int64) error {
	query := `DELETE FROM analytics_publishers WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Search methods

func (r *analyticsRepository) SearchAdvertisers(ctx context.Context, query string, limit int) ([]domain.AutocompleteResult, error) {
	searchQuery := `
		SELECT id, domain
		FROM analytics_advertisers 
		WHERE domain ILIKE $1
		ORDER BY 
			CASE WHEN domain ILIKE $2 THEN 1 ELSE 2 END,
			LENGTH(domain),
			domain
		LIMIT $3`

	likePattern := "%" + strings.ToLower(query) + "%"
	startsWithPattern := strings.ToLower(query) + "%"

	rows, err := r.db.Query(ctx, searchQuery, likePattern, startsWithPattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.AutocompleteResult
	for rows.Next() {
		var result domain.AutocompleteResult
		err := rows.Scan(&result.ID, &result.Domain)
		if err != nil {
			return nil, err
		}
		result.Type = "advertiser"
		result.Name = result.Domain
		results = append(results, result)
	}

	return results, rows.Err()
}

func (r *analyticsRepository) SearchPublishers(ctx context.Context, query string, limit int) ([]domain.AutocompleteResult, error) {
	searchQuery := `
		SELECT id, domain
		FROM analytics_publishers 
		WHERE domain ILIKE $1
		ORDER BY 
			CASE WHEN domain ILIKE $2 THEN 1 ELSE 2 END,
			LENGTH(domain),
			domain
		LIMIT $3`

	likePattern := "%" + strings.ToLower(query) + "%"
	startsWithPattern := strings.ToLower(query) + "%"

	rows, err := r.db.Query(ctx, searchQuery, likePattern, startsWithPattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.AutocompleteResult
	for rows.Next() {
		var result domain.AutocompleteResult
		err := rows.Scan(&result.ID, &result.Domain)
		if err != nil {
			return nil, err
		}
		result.Type = "publisher"
		result.Name = result.Domain
		results = append(results, result)
	}

	return results, rows.Err()
}

func (r *analyticsRepository) SearchBoth(ctx context.Context, query string, limit int) ([]domain.AutocompleteResult, error) {
	searchQuery := `
		SELECT id, domain, type FROM (
			SELECT id, domain, 'advertiser' as type
			FROM analytics_advertisers 
			WHERE domain ILIKE $1
			UNION ALL
			SELECT id, domain, 'publisher' as type
			FROM analytics_publishers 
			WHERE domain ILIKE $1
		) combined
		ORDER BY domain
		LIMIT $2`

	likePattern := "%" + strings.ToLower(query) + "%"

	rows, err := r.db.Query(ctx, searchQuery, likePattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.AutocompleteResult
	for rows.Next() {
		var result domain.AutocompleteResult
		err := rows.Scan(&result.ID, &result.Domain, &result.Type)
		if err != nil {
			return nil, err
		}
		result.Name = result.Domain
		results = append(results, result)
	}

	return results, rows.Err()
}
