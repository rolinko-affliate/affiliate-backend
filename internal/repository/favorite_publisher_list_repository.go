package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FavoritePublisherListRepository defines the interface for favorite publisher list data access
type FavoritePublisherListRepository interface {
	// List management
	CreateList(ctx context.Context, list *domain.FavoritePublisherList) error
	GetListByID(ctx context.Context, listID int64) (*domain.FavoritePublisherList, error)
	GetListsByOrganization(ctx context.Context, organizationID int64) ([]*domain.FavoritePublisherListWithStats, error)
	UpdateList(ctx context.Context, list *domain.FavoritePublisherList) error
	DeleteList(ctx context.Context, listID int64) error

	// List item management
	AddPublisherToList(ctx context.Context, item *domain.FavoritePublisherListItem) error
	RemovePublisherFromList(ctx context.Context, listID int64, publisherDomain string) error
	GetListItems(ctx context.Context, listID int64) ([]*domain.FavoritePublisherListItem, error)
	GetListItemsWithPublisherDetails(ctx context.Context, listID int64) ([]*domain.FavoritePublisherListItem, error)
	UpdatePublisherInList(ctx context.Context, listID int64, publisherDomain string, notes *string) error
	UpdatePublisherStatus(ctx context.Context, listID int64, publisherDomain string, status string) error
	GetPublisherFromList(ctx context.Context, listID int64, publisherDomain string) (*domain.FavoritePublisherListItem, error)

	// Utility methods
	IsPublisherInList(ctx context.Context, listID int64, publisherDomain string) (bool, error)
	GetListsContainingPublisher(ctx context.Context, organizationID int64, publisherDomain string) ([]*domain.FavoritePublisherList, error)
}

// favoritePublisherListRepository implements FavoritePublisherListRepository
type favoritePublisherListRepository struct {
	db *pgxpool.Pool
}

// NewFavoritePublisherListRepository creates a new favorite publisher list repository
func NewFavoritePublisherListRepository(db *pgxpool.Pool) FavoritePublisherListRepository {
	return &favoritePublisherListRepository{db: db}
}

// SQL query constants
const (
	selectListFields = "list_id, organization_id, name, description, created_at, updated_at"
	selectItemFields = "item_id, list_id, publisher_domain, notes, status, added_at"
)

// Helper function to scan a list row
func (r *favoritePublisherListRepository) scanList(row pgx.Row) (*domain.FavoritePublisherList, error) {
	list := &domain.FavoritePublisherList{}
	err := row.Scan(&list.ListID, &list.OrganizationID, &list.Name, &list.Description,
		&list.CreatedAt, &list.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan favorite publisher list: %w", err)
	}
	return list, nil
}

// Helper function to scan a list item row
func (r *favoritePublisherListRepository) scanListItem(row pgx.Row) (*domain.FavoritePublisherListItem, error) {
	item := &domain.FavoritePublisherListItem{}
	err := row.Scan(&item.ItemID, &item.ListID, &item.PublisherDomain, &item.Notes, &item.Status, &item.AddedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan favorite publisher list item: %w", err)
	}
	return item, nil
}

// CreateList creates a new favorite publisher list
func (r *favoritePublisherListRepository) CreateList(ctx context.Context, list *domain.FavoritePublisherList) error {
	query := `
		INSERT INTO favorite_publisher_lists (organization_id, name, description)
		VALUES ($1, $2, $3)
		RETURNING list_id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, list.OrganizationID, list.Name, list.Description).Scan(
		&list.ListID, &list.CreatedAt, &list.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create favorite publisher list: %w", err)
	}

	return nil
}

// GetListByID retrieves a favorite publisher list by ID
func (r *favoritePublisherListRepository) GetListByID(ctx context.Context, listID int64) (*domain.FavoritePublisherList, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM favorite_publisher_lists
		WHERE list_id = $1`, selectListFields)

	return r.scanList(r.db.QueryRow(ctx, query, listID))
}

// GetListsByOrganization retrieves all favorite publisher lists for an organization with stats
func (r *favoritePublisherListRepository) GetListsByOrganization(ctx context.Context, organizationID int64) ([]*domain.FavoritePublisherListWithStats, error) {
	query := `
		SELECT 
			fpl.list_id, fpl.organization_id, fpl.name, fpl.description, 
			fpl.created_at, fpl.updated_at,
			COALESCE(COUNT(fpli.item_id), 0) as publisher_count
		FROM favorite_publisher_lists fpl
		LEFT JOIN favorite_publisher_list_items fpli ON fpl.list_id = fpli.list_id
		WHERE fpl.organization_id = $1
		GROUP BY fpl.list_id, fpl.organization_id, fpl.name, fpl.description, fpl.created_at, fpl.updated_at
		ORDER BY fpl.created_at DESC`

	rows, err := r.db.Query(ctx, query, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get favorite publisher lists: %w", err)
	}
	defer rows.Close()

	var lists []*domain.FavoritePublisherListWithStats
	for rows.Next() {
		list := &domain.FavoritePublisherListWithStats{}
		err := rows.Scan(
			&list.ListID, &list.OrganizationID, &list.Name, &list.Description,
			&list.CreatedAt, &list.UpdatedAt, &list.PublisherCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan favorite publisher list: %w", err)
		}
		lists = append(lists, list)
	}

	return lists, nil
}

// UpdateList updates a favorite publisher list
func (r *favoritePublisherListRepository) UpdateList(ctx context.Context, list *domain.FavoritePublisherList) error {
	query := `
		UPDATE favorite_publisher_lists
		SET name = $2, description = $3, updated_at = CURRENT_TIMESTAMP
		WHERE list_id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query, list.ListID, list.Name, list.Description).Scan(&list.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to update favorite publisher list: %w", err)
	}

	return nil
}

// DeleteList deletes a favorite publisher list
func (r *favoritePublisherListRepository) DeleteList(ctx context.Context, listID int64) error {
	query := `DELETE FROM favorite_publisher_lists WHERE list_id = $1`

	result, err := r.db.Exec(ctx, query, listID)
	if err != nil {
		return fmt.Errorf("failed to delete favorite publisher list: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// AddPublisherToList adds a publisher to a favorite list
func (r *favoritePublisherListRepository) AddPublisherToList(ctx context.Context, item *domain.FavoritePublisherListItem) error {
	query := `
		INSERT INTO favorite_publisher_list_items (list_id, publisher_domain, notes, status)
		VALUES ($1, $2, $3, $4)
		RETURNING item_id, added_at`

	err := r.db.QueryRow(ctx, query, item.ListID, item.PublisherDomain, item.Notes, item.Status).Scan(
		&item.ItemID, &item.AddedAt)
	if err != nil {
		return fmt.Errorf("failed to add publisher to list: %w", err)
	}

	return nil
}

// RemovePublisherFromList removes a publisher from a favorite list
func (r *favoritePublisherListRepository) RemovePublisherFromList(ctx context.Context, listID int64, publisherDomain string) error {
	query := `DELETE FROM favorite_publisher_list_items WHERE list_id = $1 AND publisher_domain = $2`

	result, err := r.db.Exec(ctx, query, listID, publisherDomain)
	if err != nil {
		return fmt.Errorf("failed to remove publisher from list: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// GetListItems retrieves all items in a favorite list
func (r *favoritePublisherListRepository) GetListItems(ctx context.Context, listID int64) ([]*domain.FavoritePublisherListItem, error) {
	query := `
		SELECT item_id, list_id, publisher_domain, notes, status, added_at
		FROM favorite_publisher_list_items
		WHERE list_id = $1
		ORDER BY added_at DESC`

	rows, err := r.db.Query(ctx, query, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list items: %w", err)
	}
	defer rows.Close()

	var items []*domain.FavoritePublisherListItem
	for rows.Next() {
		item := &domain.FavoritePublisherListItem{}
		err := rows.Scan(&item.ItemID, &item.ListID, &item.PublisherDomain, &item.Notes, &item.Status, &item.AddedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan list item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// GetListItemsWithPublisherDetails retrieves all items in a favorite list with publisher details
func (r *favoritePublisherListRepository) GetListItemsWithPublisherDetails(ctx context.Context, listID int64) ([]*domain.FavoritePublisherListItem, error) {
	query := `
		SELECT 
			fpli.item_id, fpli.list_id, fpli.publisher_domain, fpli.notes, fpli.status, fpli.added_at,
			ap.id, ap.domain, ap.description, ap.favicon_image_url, ap.screenshot_image_url,
			ap.known, ap.relevance, ap.traffic_score, ap.promotype, ap.created_at, ap.updated_at
		FROM favorite_publisher_list_items fpli
		LEFT JOIN analytics_publishers ap ON fpli.publisher_domain = ap.domain
		WHERE fpli.list_id = $1
		ORDER BY fpli.added_at DESC`

	rows, err := r.db.Query(ctx, query, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list items with details: %w", err)
	}
	defer rows.Close()

	var items []*domain.FavoritePublisherListItem
	for rows.Next() {
		item := &domain.FavoritePublisherListItem{}
		var publisher domain.AnalyticsPublisher
		var publisherID sql.NullInt64
		var publisherDomain sql.NullString
		var publisherDescription sql.NullString
		var publisherFavicon sql.NullString
		var publisherScreenshot sql.NullString
		var publisherKnown sql.NullBool
		var publisherRelevance sql.NullFloat64
		var publisherTrafficScore sql.NullFloat64
		var publisherPromotype sql.NullString
		var publisherCreatedAt sql.NullTime
		var publisherUpdatedAt sql.NullTime

		err := rows.Scan(
			&item.ItemID, &item.ListID, &item.PublisherDomain, &item.Notes, &item.Status, &item.AddedAt,
			&publisherID, &publisherDomain, &publisherDescription, &publisherFavicon, &publisherScreenshot,
			&publisherKnown, &publisherRelevance, &publisherTrafficScore, &publisherPromotype,
			&publisherCreatedAt, &publisherUpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan list item with details: %w", err)
		}

		// If publisher data exists, populate the Publisher field
		if publisherID.Valid {
			publisher.ID = publisherID.Int64
			publisher.Domain = publisherDomain.String
			if publisherDescription.Valid {
				publisher.Description = &publisherDescription.String
			}
			if publisherFavicon.Valid {
				publisher.FaviconImageURL = &publisherFavicon.String
			}
			if publisherScreenshot.Valid {
				publisher.ScreenshotImageURL = &publisherScreenshot.String
			}
			publisher.Known = publisherKnown.Bool
			publisher.Relevance = publisherRelevance.Float64
			publisher.TrafficScore = publisherTrafficScore.Float64
			if publisherPromotype.Valid {
				publisher.Promotype = &publisherPromotype.String
			}
			if publisherCreatedAt.Valid {
				publisher.CreatedAt = publisherCreatedAt.Time
			}
			if publisherUpdatedAt.Valid {
				publisher.UpdatedAt = publisherUpdatedAt.Time
			}
			item.Publisher = &publisher
		}

		items = append(items, item)
	}

	return items, nil
}

// UpdatePublisherInList updates the notes for a publisher in a list
func (r *favoritePublisherListRepository) UpdatePublisherInList(ctx context.Context, listID int64, publisherDomain string, notes *string) error {
	query := `
		UPDATE favorite_publisher_list_items
		SET notes = $3
		WHERE list_id = $1 AND publisher_domain = $2`

	result, err := r.db.Exec(ctx, query, listID, publisherDomain, notes)
	if err != nil {
		return fmt.Errorf("failed to update publisher in list: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// IsPublisherInList checks if a publisher is already in a list
func (r *favoritePublisherListRepository) IsPublisherInList(ctx context.Context, listID int64, publisherDomain string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM favorite_publisher_list_items WHERE list_id = $1 AND publisher_domain = $2)`

	var exists bool
	err := r.db.QueryRow(ctx, query, listID, publisherDomain).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if publisher is in list: %w", err)
	}

	return exists, nil
}

// GetListsContainingPublisher retrieves all lists that contain a specific publisher for an organization
func (r *favoritePublisherListRepository) GetListsContainingPublisher(ctx context.Context, organizationID int64, publisherDomain string) ([]*domain.FavoritePublisherList, error) {
	query := `
		SELECT fpl.list_id, fpl.organization_id, fpl.name, fpl.description, fpl.created_at, fpl.updated_at
		FROM favorite_publisher_lists fpl
		INNER JOIN favorite_publisher_list_items fpli ON fpl.list_id = fpli.list_id
		WHERE fpl.organization_id = $1 AND fpli.publisher_domain = $2
		ORDER BY fpl.name`

	rows, err := r.db.Query(ctx, query, organizationID, publisherDomain)
	if err != nil {
		return nil, fmt.Errorf("failed to get lists containing publisher: %w", err)
	}
	defer rows.Close()

	var lists []*domain.FavoritePublisherList
	for rows.Next() {
		list := &domain.FavoritePublisherList{}
		err := rows.Scan(&list.ListID, &list.OrganizationID, &list.Name, &list.Description, &list.CreatedAt, &list.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan list: %w", err)
		}
		lists = append(lists, list)
	}

	return lists, nil
}

// UpdatePublisherStatus updates the status of a publisher in a favorite list
func (r *favoritePublisherListRepository) UpdatePublisherStatus(ctx context.Context, listID int64, publisherDomain string, status string) error {
	query := `UPDATE favorite_publisher_list_items SET status = $1 WHERE list_id = $2 AND publisher_domain = $3`

	result, err := r.db.Exec(ctx, query, status, listID, publisherDomain)
	if err != nil {
		return fmt.Errorf("failed to update publisher status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// GetPublisherFromList retrieves a specific publisher from a favorite list
func (r *favoritePublisherListRepository) GetPublisherFromList(ctx context.Context, listID int64, publisherDomain string) (*domain.FavoritePublisherListItem, error) {
	query := fmt.Sprintf(`SELECT %s FROM favorite_publisher_list_items WHERE list_id = $1 AND publisher_domain = $2`, selectItemFields)

	row := r.db.QueryRow(ctx, query, listID, publisherDomain)
	return r.scanListItem(row)
}
