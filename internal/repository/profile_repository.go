package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProfileRepository defines the interface for profile operations
type ProfileRepository interface {
	CreateProfile(ctx context.Context, profile *domain.Profile) error
	GetProfileByID(ctx context.Context, id uuid.UUID) (*domain.Profile, error)
	GetRoleByID(ctx context.Context, roleID int) (*domain.Role, error)
	// Add other methods as needed
}

// pgxProfileRepository implements ProfileRepository using pgx
type pgxProfileRepository struct {
	db *pgxpool.Pool
}

// NewPgxProfileRepository creates a new profile repository
func NewPgxProfileRepository(db *pgxpool.Pool) ProfileRepository {
	return &pgxProfileRepository{db: db}
}

// CreateProfile creates a new profile in the database
func (r *pgxProfileRepository) CreateProfile(ctx context.Context, profile *domain.Profile) error {
	query := `INSERT INTO public.profiles (id, organization_id, role_id, email, first_name, last_name, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query,
		profile.ID, profile.OrganizationID, profile.RoleID, profile.Email,
		profile.FirstName, profile.LastName, profile.CreatedAt, profile.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error creating profile: %w", err)
	}
	return nil
}

// GetProfileByID retrieves a profile by ID
func (r *pgxProfileRepository) GetProfileByID(ctx context.Context, id uuid.UUID) (*domain.Profile, error) {
	query := `SELECT id, organization_id, role_id, email, first_name, last_name, created_at, updated_at
              FROM public.profiles WHERE id = $1`
	var p domain.Profile
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.OrganizationID, &p.RoleID, &p.Email,
		&p.FirstName, &p.LastName, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("profile not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("error getting profile by ID: %w", err)
	}
	return &p, nil
}

// GetRoleByID retrieves a role by ID
func (r *pgxProfileRepository) GetRoleByID(ctx context.Context, roleID int) (*domain.Role, error) {
	query := `SELECT role_id, name, description FROM public.roles WHERE role_id = $1`
	var role domain.Role
	var description sql.NullString
	err := r.db.QueryRow(ctx, query, roleID).Scan(&role.RoleID, &role.Name, &description)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("role not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("error getting role by ID: %w", err)
	}
	
	if description.Valid {
		desc := description.String
		role.Description = &desc
	}
	
	return &role, nil
}