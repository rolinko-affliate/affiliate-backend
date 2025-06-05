package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OrganizationRepository defines the interface for organization operations
type OrganizationRepository interface {
	CreateOrganization(ctx context.Context, org *domain.Organization) error
	GetOrganizationByID(ctx context.Context, id int64) (*domain.Organization, error)
	UpdateOrganization(ctx context.Context, org *domain.Organization) error
	ListOrganizations(ctx context.Context, limit, offset int) ([]*domain.Organization, error)
	DeleteOrganization(ctx context.Context, id int64) error
}

// pgxOrganizationRepository implements OrganizationRepository using pgx
type pgxOrganizationRepository struct {
	db *pgxpool.Pool
}

// NewPgxOrganizationRepository creates a new organization repository
func NewPgxOrganizationRepository(db *pgxpool.Pool) OrganizationRepository {
	return &pgxOrganizationRepository{db: db}
}

// CreateOrganization creates a new organization in the database
func (r *pgxOrganizationRepository) CreateOrganization(ctx context.Context, org *domain.Organization) error {
	query := `INSERT INTO public.organizations (name, type, created_at, updated_at)
              VALUES ($1, $2, $3, $4)
              RETURNING organization_id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(ctx, query, org.Name, org.Type, now, now).Scan(
		&org.OrganizationID, &org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creating organization: %w", err)
	}
	return nil
}

// GetOrganizationByID retrieves an organization by ID
func (r *pgxOrganizationRepository) GetOrganizationByID(ctx context.Context, id int64) (*domain.Organization, error) {
	query := `SELECT organization_id, name, type, created_at, updated_at
              FROM public.organizations WHERE organization_id = $1`

	var org domain.Organization
	err := r.db.QueryRow(ctx, query, id).Scan(
		&org.OrganizationID, &org.Name, &org.Type, &org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("organization not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("error getting organization by ID: %w", err)
	}
	return &org, nil
}

// UpdateOrganization updates an organization in the database
func (r *pgxOrganizationRepository) UpdateOrganization(ctx context.Context, org *domain.Organization) error {
	query := `UPDATE public.organizations
              SET name = $1, type = $2
              WHERE organization_id = $3
              RETURNING updated_at`

	err := r.db.QueryRow(ctx, query, org.Name, org.Type, org.OrganizationID).Scan(&org.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("organization not found: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("error updating organization: %w", err)
	}
	return nil
}

// ListOrganizations retrieves a list of organizations with pagination
func (r *pgxOrganizationRepository) ListOrganizations(ctx context.Context, limit, offset int) ([]*domain.Organization, error) {
	query := `SELECT organization_id, name, type, created_at, updated_at
              FROM public.organizations
              ORDER BY organization_id
              LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing organizations: %w", err)
	}
	defer rows.Close()

	var organizations []*domain.Organization
	for rows.Next() {
		var org domain.Organization
		if err := rows.Scan(&org.OrganizationID, &org.Name, &org.Type, &org.CreatedAt, &org.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning organization row: %w", err)
		}
		organizations = append(organizations, &org)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating organization rows: %w", err)
	}

	return organizations, nil
}

// DeleteOrganization deletes an organization from the database
func (r *pgxOrganizationRepository) DeleteOrganization(ctx context.Context, id int64) error {
	query := `DELETE FROM public.organizations WHERE organization_id = $1`

	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting organization: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("organization not found: %w", domain.ErrNotFound)
	}

	return nil
}
