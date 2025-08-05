package repository

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// BillingAccountRepository defines the interface for billing account operations
type BillingAccountRepository interface {
	Create(ctx context.Context, account *domain.BillingAccount) error
	GetByID(ctx context.Context, billingAccountID int64) (*domain.BillingAccount, error)
	GetByOrganizationID(ctx context.Context, organizationID int64) (*domain.BillingAccount, error)
	GetByStripeCustomerID(ctx context.Context, stripeCustomerID string) (*domain.BillingAccount, error)
	Update(ctx context.Context, account *domain.BillingAccount) error
	UpdateBalance(ctx context.Context, billingAccountID int64, newBalance decimal.Decimal) error
	List(ctx context.Context, limit, offset int) ([]domain.BillingAccount, error)
	Delete(ctx context.Context, billingAccountID int64) error
}

// PgxBillingAccountRepository implements BillingAccountRepository using pgx
type PgxBillingAccountRepository struct {
	db *pgxpool.Pool
}

// NewPgxBillingAccountRepository creates a new PgxBillingAccountRepository
func NewPgxBillingAccountRepository(db *pgxpool.Pool) BillingAccountRepository {
	return &PgxBillingAccountRepository{db: db}
}

// Create creates a new billing account
func (r *PgxBillingAccountRepository) Create(ctx context.Context, account *domain.BillingAccount) error {
	query := `
		INSERT INTO billing_accounts (
			organization_id, stripe_customer_id, stripe_account_id, billing_mode, currency,
			balance, credit_limit, default_payment_method_id, auto_recharge_enabled,
			auto_recharge_threshold, auto_recharge_amount, invoice_day_of_month,
			payment_terms_days, status, billing_email, billing_address, tax_info
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		) RETURNING billing_account_id, created_at, updated_at`

	billingAddressJSON, err := json.Marshal(account.BillingAddress)
	if err != nil {
		return fmt.Errorf("failed to marshal billing address: %w", err)
	}

	taxInfoJSON, err := json.Marshal(account.TaxInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal tax info: %w", err)
	}

	err = r.db.QueryRow(ctx, query,
		account.OrganizationID,
		account.StripeCustomerID,
		account.StripeAccountID,
		account.BillingMode,
		account.Currency,
		account.Balance,
		account.CreditLimit,
		account.DefaultPaymentMethodID,
		account.AutoRechargeEnabled,
		account.AutoRechargeThreshold,
		account.AutoRechargeAmount,
		account.InvoiceDayOfMonth,
		account.PaymentTermsDays,
		account.Status,
		account.BillingEmail,
		billingAddressJSON,
		taxInfoJSON,
	).Scan(&account.BillingAccountID, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create billing account: %w", err)
	}

	return nil
}

// GetByID retrieves a billing account by ID
func (r *PgxBillingAccountRepository) GetByID(ctx context.Context, billingAccountID int64) (*domain.BillingAccount, error) {
	query := `
		SELECT billing_account_id, organization_id, stripe_customer_id, stripe_account_id,
			   billing_mode, currency, balance, credit_limit, default_payment_method_id,
			   auto_recharge_enabled, auto_recharge_threshold, auto_recharge_amount,
			   invoice_day_of_month, payment_terms_days, status, billing_email,
			   billing_address, tax_info, created_at, updated_at
		FROM billing_accounts
		WHERE billing_account_id = $1`

	account := &domain.BillingAccount{}
	var billingAddressJSON, taxInfoJSON []byte

	err := r.db.QueryRow(ctx, query, billingAccountID).Scan(
		&account.BillingAccountID,
		&account.OrganizationID,
		&account.StripeCustomerID,
		&account.StripeAccountID,
		&account.BillingMode,
		&account.Currency,
		&account.Balance,
		&account.CreditLimit,
		&account.DefaultPaymentMethodID,
		&account.AutoRechargeEnabled,
		&account.AutoRechargeThreshold,
		&account.AutoRechargeAmount,
		&account.InvoiceDayOfMonth,
		&account.PaymentTermsDays,
		&account.Status,
		&account.BillingEmail,
		&billingAddressJSON,
		&taxInfoJSON,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("billing account not found")
		}
		return nil, fmt.Errorf("failed to get billing account: %w", err)
	}

	// Unmarshal JSON fields
	if len(billingAddressJSON) > 0 {
		if err := json.Unmarshal(billingAddressJSON, &account.BillingAddress); err != nil {
			return nil, fmt.Errorf("failed to unmarshal billing address: %w", err)
		}
	}

	if len(taxInfoJSON) > 0 {
		if err := json.Unmarshal(taxInfoJSON, &account.TaxInfo); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tax info: %w", err)
		}
	}

	return account, nil
}

// GetByOrganizationID retrieves a billing account by organization ID
func (r *PgxBillingAccountRepository) GetByOrganizationID(ctx context.Context, organizationID int64) (*domain.BillingAccount, error) {
	query := `
		SELECT billing_account_id, organization_id, stripe_customer_id, stripe_account_id,
			   billing_mode, currency, balance, credit_limit, default_payment_method_id,
			   auto_recharge_enabled, auto_recharge_threshold, auto_recharge_amount,
			   invoice_day_of_month, payment_terms_days, status, billing_email,
			   billing_address, tax_info, created_at, updated_at
		FROM billing_accounts
		WHERE organization_id = $1`

	account := &domain.BillingAccount{}
	var billingAddressJSON, taxInfoJSON []byte

	err := r.db.QueryRow(ctx, query, organizationID).Scan(
		&account.BillingAccountID,
		&account.OrganizationID,
		&account.StripeCustomerID,
		&account.StripeAccountID,
		&account.BillingMode,
		&account.Currency,
		&account.Balance,
		&account.CreditLimit,
		&account.DefaultPaymentMethodID,
		&account.AutoRechargeEnabled,
		&account.AutoRechargeThreshold,
		&account.AutoRechargeAmount,
		&account.InvoiceDayOfMonth,
		&account.PaymentTermsDays,
		&account.Status,
		&account.BillingEmail,
		&billingAddressJSON,
		&taxInfoJSON,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("billing account not found for organization")
		}
		return nil, fmt.Errorf("failed to get billing account by organization: %w", err)
	}

	// Unmarshal JSON fields
	if len(billingAddressJSON) > 0 {
		if err := json.Unmarshal(billingAddressJSON, &account.BillingAddress); err != nil {
			return nil, fmt.Errorf("failed to unmarshal billing address: %w", err)
		}
	}

	if len(taxInfoJSON) > 0 {
		if err := json.Unmarshal(taxInfoJSON, &account.TaxInfo); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tax info: %w", err)
		}
	}

	return account, nil
}

// GetByStripeCustomerID retrieves a billing account by Stripe customer ID
func (r *PgxBillingAccountRepository) GetByStripeCustomerID(ctx context.Context, stripeCustomerID string) (*domain.BillingAccount, error) {
	query := `
		SELECT billing_account_id, organization_id, stripe_customer_id, stripe_account_id,
			   billing_mode, currency, balance, credit_limit, default_payment_method_id,
			   auto_recharge_enabled, auto_recharge_threshold, auto_recharge_amount,
			   invoice_day_of_month, payment_terms_days, status, billing_email,
			   billing_address, tax_info, created_at, updated_at
		FROM billing_accounts
		WHERE stripe_customer_id = $1`

	account := &domain.BillingAccount{}
	var billingAddressJSON, taxInfoJSON []byte

	err := r.db.QueryRow(ctx, query, stripeCustomerID).Scan(
		&account.BillingAccountID,
		&account.OrganizationID,
		&account.StripeCustomerID,
		&account.StripeAccountID,
		&account.BillingMode,
		&account.Currency,
		&account.Balance,
		&account.CreditLimit,
		&account.DefaultPaymentMethodID,
		&account.AutoRechargeEnabled,
		&account.AutoRechargeThreshold,
		&account.AutoRechargeAmount,
		&account.InvoiceDayOfMonth,
		&account.PaymentTermsDays,
		&account.Status,
		&account.BillingEmail,
		&billingAddressJSON,
		&taxInfoJSON,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("billing account not found for Stripe customer")
		}
		return nil, fmt.Errorf("failed to get billing account by Stripe customer ID: %w", err)
	}

	// Unmarshal JSON fields
	if len(billingAddressJSON) > 0 {
		if err := json.Unmarshal(billingAddressJSON, &account.BillingAddress); err != nil {
			return nil, fmt.Errorf("failed to unmarshal billing address: %w", err)
		}
	}

	if len(taxInfoJSON) > 0 {
		if err := json.Unmarshal(taxInfoJSON, &account.TaxInfo); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tax info: %w", err)
		}
	}

	return account, nil
}

// Update updates a billing account
func (r *PgxBillingAccountRepository) Update(ctx context.Context, account *domain.BillingAccount) error {
	query := `
		UPDATE billing_accounts SET
			stripe_customer_id = $2,
			stripe_account_id = $3,
			billing_mode = $4,
			currency = $5,
			balance = $6,
			credit_limit = $7,
			default_payment_method_id = $8,
			auto_recharge_enabled = $9,
			auto_recharge_threshold = $10,
			auto_recharge_amount = $11,
			invoice_day_of_month = $12,
			payment_terms_days = $13,
			status = $14,
			billing_email = $15,
			billing_address = $16,
			tax_info = $17,
			updated_at = NOW()
		WHERE billing_account_id = $1
		RETURNING updated_at`

	billingAddressJSON, err := json.Marshal(account.BillingAddress)
	if err != nil {
		return fmt.Errorf("failed to marshal billing address: %w", err)
	}

	taxInfoJSON, err := json.Marshal(account.TaxInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal tax info: %w", err)
	}

	err = r.db.QueryRow(ctx, query,
		account.BillingAccountID,
		account.StripeCustomerID,
		account.StripeAccountID,
		account.BillingMode,
		account.Currency,
		account.Balance,
		account.CreditLimit,
		account.DefaultPaymentMethodID,
		account.AutoRechargeEnabled,
		account.AutoRechargeThreshold,
		account.AutoRechargeAmount,
		account.InvoiceDayOfMonth,
		account.PaymentTermsDays,
		account.Status,
		account.BillingEmail,
		billingAddressJSON,
		taxInfoJSON,
	).Scan(&account.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update billing account: %w", err)
	}

	return nil
}

// UpdateBalance updates only the balance of a billing account
func (r *PgxBillingAccountRepository) UpdateBalance(ctx context.Context, billingAccountID int64, newBalance decimal.Decimal) error {
	query := `
		UPDATE billing_accounts SET
			balance = $2,
			updated_at = NOW()
		WHERE billing_account_id = $1`

	result, err := r.db.Exec(ctx, query, billingAccountID, newBalance)
	if err != nil {
		return fmt.Errorf("failed to update billing account balance: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("billing account not found")
	}

	return nil
}

// List retrieves a list of billing accounts with pagination
func (r *PgxBillingAccountRepository) List(ctx context.Context, limit, offset int) ([]domain.BillingAccount, error) {
	query := `
		SELECT billing_account_id, organization_id, stripe_customer_id, stripe_account_id,
			   billing_mode, currency, balance, credit_limit, default_payment_method_id,
			   auto_recharge_enabled, auto_recharge_threshold, auto_recharge_amount,
			   invoice_day_of_month, payment_terms_days, status, billing_email,
			   billing_address, tax_info, created_at, updated_at
		FROM billing_accounts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list billing accounts: %w", err)
	}
	defer rows.Close()

	accounts := make([]domain.BillingAccount, 0)
	for rows.Next() {
		account := domain.BillingAccount{}
		var billingAddressJSON, taxInfoJSON []byte

		err := rows.Scan(
			&account.BillingAccountID,
			&account.OrganizationID,
			&account.StripeCustomerID,
			&account.StripeAccountID,
			&account.BillingMode,
			&account.Currency,
			&account.Balance,
			&account.CreditLimit,
			&account.DefaultPaymentMethodID,
			&account.AutoRechargeEnabled,
			&account.AutoRechargeThreshold,
			&account.AutoRechargeAmount,
			&account.InvoiceDayOfMonth,
			&account.PaymentTermsDays,
			&account.Status,
			&account.BillingEmail,
			&billingAddressJSON,
			&taxInfoJSON,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan billing account: %w", err)
		}

		// Unmarshal JSON fields
		if len(billingAddressJSON) > 0 {
			if err := json.Unmarshal(billingAddressJSON, &account.BillingAddress); err != nil {
				return nil, fmt.Errorf("failed to unmarshal billing address: %w", err)
			}
		}

		if len(taxInfoJSON) > 0 {
			if err := json.Unmarshal(taxInfoJSON, &account.TaxInfo); err != nil {
				return nil, fmt.Errorf("failed to unmarshal tax info: %w", err)
			}
		}

		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating billing accounts: %w", err)
	}

	return accounts, nil
}

// Delete deletes a billing account
func (r *PgxBillingAccountRepository) Delete(ctx context.Context, billingAccountID int64) error {
	query := `DELETE FROM billing_accounts WHERE billing_account_id = $1`

	result, err := r.db.Exec(ctx, query, billingAccountID)
	if err != nil {
		return fmt.Errorf("failed to delete billing account: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("billing account not found")
	}

	return nil
}

// JSONMap is a custom type for JSON fields
type JSONMap map[string]interface{}

// Value implements driver.Valuer for JSON fields
func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}
