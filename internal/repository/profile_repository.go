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
	UpdateProfile(ctx context.Context, profile *domain.Profile) error
	DeleteProfile(ctx context.Context, id uuid.UUID) error
	GetRoleByID(ctx context.Context, roleID int) (*domain.Role, error)
	UpsertProfile(ctx context.Context, profile *domain.Profile) error
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
	// First, get the role name for the given role_id
	roleQuery := `SELECT name FROM public.roles WHERE role_id = $1`
	err := r.db.QueryRow(ctx, roleQuery, profile.RoleID).Scan(&profile.RoleName)
	if err != nil {
		return fmt.Errorf("error getting role name for role_id %d: %w", profile.RoleID, err)
	}

	// Now insert the profile
	query := `INSERT INTO public.profiles (id, organization_id, role_id, email, first_name, last_name, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = r.db.Exec(ctx, query,
		profile.ID, profile.OrganizationID, profile.RoleID, profile.Email,
		profile.FirstName, profile.LastName, profile.CreatedAt, profile.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error creating profile: %w", err)
	}
	return nil
}

// GetProfileByID retrieves a profile by ID
func (r *pgxProfileRepository) GetProfileByID(ctx context.Context, id uuid.UUID) (*domain.Profile, error) {
	query := `SELECT p.id, p.organization_id, p.role_id, r.name as role_name, p.email, p.first_name, p.last_name, p.created_at, p.updated_at
              FROM public.profiles p
              JOIN public.roles r ON p.role_id = r.role_id
              WHERE p.id = $1`
	var p domain.Profile
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.OrganizationID, &p.RoleID, &p.RoleName, &p.Email,
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

// UpdateProfile updates an existing profile in the database
func (r *pgxProfileRepository) UpdateProfile(ctx context.Context, profile *domain.Profile) error {
	// First, get the role name for the given role_id
	roleQuery := `SELECT name FROM public.roles WHERE role_id = $1`
	err := r.db.QueryRow(ctx, roleQuery, profile.RoleID).Scan(&profile.RoleName)
	if err != nil {
		return fmt.Errorf("error getting role name for role_id %d: %w", profile.RoleID, err)
	}

	query := `UPDATE public.profiles 
	          SET organization_id = $1, role_id = $2, email = $3, first_name = $4, last_name = $5, updated_at = $6
	          WHERE id = $7`

	result, err := r.db.Exec(ctx, query,
		profile.OrganizationID, profile.RoleID, profile.Email,
		profile.FirstName, profile.LastName, profile.UpdatedAt,
		profile.ID)

	if err != nil {
		return fmt.Errorf("error updating profile: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("profile not found: %w", domain.ErrNotFound)
	}

	return nil
}

// DeleteProfile deletes a profile from the database
func (r *pgxProfileRepository) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM public.profiles WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting profile: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("profile not found: %w", domain.ErrNotFound)
	}

	return nil
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

// UpsertProfile creates or updates a profile in the database
func (r *pgxProfileRepository) UpsertProfile(ctx context.Context, profile *domain.Profile) error {
	// First, get the role name for the given role_id
	roleQuery := `SELECT name FROM public.roles WHERE role_id = $1`
	err := r.db.QueryRow(ctx, roleQuery, profile.RoleID).Scan(&profile.RoleName)
	if err != nil {
		return fmt.Errorf("error getting role name for role_id %d: %w", profile.RoleID, err)
	}

	query := `
		INSERT INTO public.profiles (id, organization_id, role_id, email, first_name, last_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) 
		DO UPDATE SET 
			organization_id = $2,
			role_id = $3,
			email = $4,
			first_name = $5,
			last_name = $6,
			updated_at = $8
	`
	_, err = r.db.Exec(ctx, query,
		profile.ID, profile.OrganizationID, profile.RoleID, profile.Email,
		profile.FirstName, profile.LastName, profile.CreatedAt, profile.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error upserting profile: %w", err)
	}
	return nil
}
