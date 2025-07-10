package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PaymentMethodRepository defines the interface for payment method operations
type PaymentMethodRepository interface {
	Create(ctx context.Context, paymentMethod *domain.StripePaymentMethod) error
	GetByID(ctx context.Context, paymentMethodID int64) (*domain.StripePaymentMethod, error)
	GetByStripePaymentMethodID(ctx context.Context, stripePaymentMethodID string) (*domain.StripePaymentMethod, error)
	GetByOrganizationID(ctx context.Context, organizationID int64) ([]domain.StripePaymentMethod, error)
	GetDefaultByOrganizationID(ctx context.Context, organizationID int64) (*domain.StripePaymentMethod, error)
	Update(ctx context.Context, paymentMethod *domain.StripePaymentMethod) error
	SetAsDefault(ctx context.Context, paymentMethodID int64, organizationID int64) error
	Delete(ctx context.Context, paymentMethodID int64) error
	List(ctx context.Context, limit, offset int) ([]domain.StripePaymentMethod, error)
}

// PgxPaymentMethodRepository implements PaymentMethodRepository using pgx
type PgxPaymentMethodRepository struct {
	db *pgxpool.Pool
}

// NewPgxPaymentMethodRepository creates a new PgxPaymentMethodRepository
func NewPgxPaymentMethodRepository(db *pgxpool.Pool) PaymentMethodRepository {
	return &PgxPaymentMethodRepository{db: db}
}

// Create creates a new payment method
func (r *PgxPaymentMethodRepository) Create(ctx context.Context, paymentMethod *domain.StripePaymentMethod) error {
	query := `
		INSERT INTO payment_methods (
			organization_id, billing_account_id, stripe_payment_method_id, type, brand,
			last4, exp_month, exp_year, bank_name, account_holder_type, is_default,
			status, nickname, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		) RETURNING payment_method_id, created_at, updated_at`

	metadataJSON, err := json.Marshal(paymentMethod.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	err = r.db.QueryRow(ctx, query,
		paymentMethod.OrganizationID,
		paymentMethod.BillingAccountID,
		paymentMethod.StripePaymentMethodID,
		paymentMethod.Type,
		paymentMethod.Brand,
		paymentMethod.Last4,
		paymentMethod.ExpMonth,
		paymentMethod.ExpYear,
		paymentMethod.BankName,
		paymentMethod.AccountHolderType,
		paymentMethod.IsDefault,
		paymentMethod.Status,
		paymentMethod.Nickname,
		metadataJSON,
	).Scan(&paymentMethod.PaymentMethodID, &paymentMethod.CreatedAt, &paymentMethod.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create payment method: %w", err)
	}

	return nil
}

// GetByID retrieves a payment method by ID
func (r *PgxPaymentMethodRepository) GetByID(ctx context.Context, paymentMethodID int64) (*domain.StripePaymentMethod, error) {
	query := `
		SELECT payment_method_id, organization_id, billing_account_id, stripe_payment_method_id,
			   type, brand, last4, exp_month, exp_year, bank_name, account_holder_type,
			   is_default, status, nickname, metadata, created_at, updated_at
		FROM payment_methods
		WHERE payment_method_id = $1`

	paymentMethod := &domain.StripePaymentMethod{}
	var metadataJSON []byte

	err := r.db.QueryRow(ctx, query, paymentMethodID).Scan(
		&paymentMethod.PaymentMethodID,
		&paymentMethod.OrganizationID,
		&paymentMethod.BillingAccountID,
		&paymentMethod.StripePaymentMethodID,
		&paymentMethod.Type,
		&paymentMethod.Brand,
		&paymentMethod.Last4,
		&paymentMethod.ExpMonth,
		&paymentMethod.ExpYear,
		&paymentMethod.BankName,
		&paymentMethod.AccountHolderType,
		&paymentMethod.IsDefault,
		&paymentMethod.Status,
		&paymentMethod.Nickname,
		&metadataJSON,
		&paymentMethod.CreatedAt,
		&paymentMethod.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("payment method not found")
		}
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}

	// Unmarshal metadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &paymentMethod.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return paymentMethod, nil
}

// GetByStripePaymentMethodID retrieves a payment method by Stripe payment method ID
func (r *PgxPaymentMethodRepository) GetByStripePaymentMethodID(ctx context.Context, stripePaymentMethodID string) (*domain.StripePaymentMethod, error) {
	query := `
		SELECT payment_method_id, organization_id, billing_account_id, stripe_payment_method_id,
			   type, brand, last4, exp_month, exp_year, bank_name, account_holder_type,
			   is_default, status, nickname, metadata, created_at, updated_at
		FROM payment_methods
		WHERE stripe_payment_method_id = $1`

	paymentMethod := &domain.StripePaymentMethod{}
	var metadataJSON []byte

	err := r.db.QueryRow(ctx, query, stripePaymentMethodID).Scan(
		&paymentMethod.PaymentMethodID,
		&paymentMethod.OrganizationID,
		&paymentMethod.BillingAccountID,
		&paymentMethod.StripePaymentMethodID,
		&paymentMethod.Type,
		&paymentMethod.Brand,
		&paymentMethod.Last4,
		&paymentMethod.ExpMonth,
		&paymentMethod.ExpYear,
		&paymentMethod.BankName,
		&paymentMethod.AccountHolderType,
		&paymentMethod.IsDefault,
		&paymentMethod.Status,
		&paymentMethod.Nickname,
		&metadataJSON,
		&paymentMethod.CreatedAt,
		&paymentMethod.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("payment method not found")
		}
		return nil, fmt.Errorf("failed to get payment method by Stripe ID: %w", err)
	}

	// Unmarshal metadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &paymentMethod.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return paymentMethod, nil
}

// GetByOrganizationID retrieves all payment methods for an organization
func (r *PgxPaymentMethodRepository) GetByOrganizationID(ctx context.Context, organizationID int64) ([]domain.StripePaymentMethod, error) {
	query := `
		SELECT payment_method_id, organization_id, billing_account_id, stripe_payment_method_id,
			   type, brand, last4, exp_month, exp_year, bank_name, account_holder_type,
			   is_default, status, nickname, metadata, created_at, updated_at
		FROM payment_methods
		WHERE organization_id = $1 AND status = 'active'
		ORDER BY is_default DESC, created_at DESC`

	rows, err := r.db.Query(ctx, query, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods by organization: %w", err)
	}
	defer rows.Close()

	var paymentMethods []domain.StripePaymentMethod
	for rows.Next() {
		paymentMethod := domain.StripePaymentMethod{}
		var metadataJSON []byte

		err := rows.Scan(
			&paymentMethod.PaymentMethodID,
			&paymentMethod.OrganizationID,
			&paymentMethod.BillingAccountID,
			&paymentMethod.StripePaymentMethodID,
			&paymentMethod.Type,
			&paymentMethod.Brand,
			&paymentMethod.Last4,
			&paymentMethod.ExpMonth,
			&paymentMethod.ExpYear,
			&paymentMethod.BankName,
			&paymentMethod.AccountHolderType,
			&paymentMethod.IsDefault,
			&paymentMethod.Status,
			&paymentMethod.Nickname,
			&metadataJSON,
			&paymentMethod.CreatedAt,
			&paymentMethod.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment method: %w", err)
		}

		// Unmarshal metadata
		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &paymentMethod.Metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		paymentMethods = append(paymentMethods, paymentMethod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating payment methods: %w", err)
	}

	return paymentMethods, nil
}

// GetDefaultByOrganizationID retrieves the default payment method for an organization
func (r *PgxPaymentMethodRepository) GetDefaultByOrganizationID(ctx context.Context, organizationID int64) (*domain.StripePaymentMethod, error) {
	query := `
		SELECT payment_method_id, organization_id, billing_account_id, stripe_payment_method_id,
			   type, brand, last4, exp_month, exp_year, bank_name, account_holder_type,
			   is_default, status, nickname, metadata, created_at, updated_at
		FROM payment_methods
		WHERE organization_id = $1 AND is_default = true AND status = 'active'`

	paymentMethod := &domain.StripePaymentMethod{}
	var metadataJSON []byte

	err := r.db.QueryRow(ctx, query, organizationID).Scan(
		&paymentMethod.PaymentMethodID,
		&paymentMethod.OrganizationID,
		&paymentMethod.BillingAccountID,
		&paymentMethod.StripePaymentMethodID,
		&paymentMethod.Type,
		&paymentMethod.Brand,
		&paymentMethod.Last4,
		&paymentMethod.ExpMonth,
		&paymentMethod.ExpYear,
		&paymentMethod.BankName,
		&paymentMethod.AccountHolderType,
		&paymentMethod.IsDefault,
		&paymentMethod.Status,
		&paymentMethod.Nickname,
		&metadataJSON,
		&paymentMethod.CreatedAt,
		&paymentMethod.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("default payment method not found")
		}
		return nil, fmt.Errorf("failed to get default payment method: %w", err)
	}

	// Unmarshal metadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &paymentMethod.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return paymentMethod, nil
}

// Update updates a payment method
func (r *PgxPaymentMethodRepository) Update(ctx context.Context, paymentMethod *domain.StripePaymentMethod) error {
	query := `
		UPDATE payment_methods SET
			type = $2,
			brand = $3,
			last4 = $4,
			exp_month = $5,
			exp_year = $6,
			bank_name = $7,
			account_holder_type = $8,
			is_default = $9,
			status = $10,
			nickname = $11,
			metadata = $12,
			updated_at = NOW()
		WHERE payment_method_id = $1
		RETURNING updated_at`

	metadataJSON, err := json.Marshal(paymentMethod.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	err = r.db.QueryRow(ctx, query,
		paymentMethod.PaymentMethodID,
		paymentMethod.Type,
		paymentMethod.Brand,
		paymentMethod.Last4,
		paymentMethod.ExpMonth,
		paymentMethod.ExpYear,
		paymentMethod.BankName,
		paymentMethod.AccountHolderType,
		paymentMethod.IsDefault,
		paymentMethod.Status,
		paymentMethod.Nickname,
		metadataJSON,
	).Scan(&paymentMethod.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update payment method: %w", err)
	}

	return nil
}

// SetAsDefault sets a payment method as the default for an organization
func (r *PgxPaymentMethodRepository) SetAsDefault(ctx context.Context, paymentMethodID int64, organizationID int64) error {
	// Start a transaction to ensure atomicity
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// First, unset all other payment methods as default for this organization
	unsetQuery := `
		UPDATE payment_methods SET
			is_default = false,
			updated_at = NOW()
		WHERE organization_id = $1 AND is_default = true`

	_, err = tx.Exec(ctx, unsetQuery, organizationID)
	if err != nil {
		return fmt.Errorf("failed to unset default payment methods: %w", err)
	}

	// Then set the specified payment method as default
	setQuery := `
		UPDATE payment_methods SET
			is_default = true,
			updated_at = NOW()
		WHERE payment_method_id = $1 AND organization_id = $2`

	result, err := tx.Exec(ctx, setQuery, paymentMethodID, organizationID)
	if err != nil {
		return fmt.Errorf("failed to set payment method as default: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("payment method not found or does not belong to organization")
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Delete deletes a payment method (soft delete by setting status to inactive)
func (r *PgxPaymentMethodRepository) Delete(ctx context.Context, paymentMethodID int64) error {
	query := `
		UPDATE payment_methods SET
			status = 'inactive',
			is_default = false,
			updated_at = NOW()
		WHERE payment_method_id = $1`

	result, err := r.db.Exec(ctx, query, paymentMethodID)
	if err != nil {
		return fmt.Errorf("failed to delete payment method: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("payment method not found")
	}

	return nil
}

// List retrieves a list of payment methods with pagination
func (r *PgxPaymentMethodRepository) List(ctx context.Context, limit, offset int) ([]domain.StripePaymentMethod, error) {
	query := `
		SELECT payment_method_id, organization_id, billing_account_id, stripe_payment_method_id,
			   type, brand, last4, exp_month, exp_year, bank_name, account_holder_type,
			   is_default, status, nickname, metadata, created_at, updated_at
		FROM payment_methods
		WHERE status = 'active'
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list payment methods: %w", err)
	}
	defer rows.Close()

	var paymentMethods []domain.StripePaymentMethod
	for rows.Next() {
		paymentMethod := domain.StripePaymentMethod{}
		var metadataJSON []byte

		err := rows.Scan(
			&paymentMethod.PaymentMethodID,
			&paymentMethod.OrganizationID,
			&paymentMethod.BillingAccountID,
			&paymentMethod.StripePaymentMethodID,
			&paymentMethod.Type,
			&paymentMethod.Brand,
			&paymentMethod.Last4,
			&paymentMethod.ExpMonth,
			&paymentMethod.ExpYear,
			&paymentMethod.BankName,
			&paymentMethod.AccountHolderType,
			&paymentMethod.IsDefault,
			&paymentMethod.Status,
			&paymentMethod.Nickname,
			&metadataJSON,
			&paymentMethod.CreatedAt,
			&paymentMethod.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment method: %w", err)
		}

		// Unmarshal metadata
		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &paymentMethod.Metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		paymentMethods = append(paymentMethods, paymentMethod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating payment methods: %w", err)
	}

	return paymentMethods, nil
}