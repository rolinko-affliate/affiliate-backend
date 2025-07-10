package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// TransactionRepository defines the interface for transaction operations
type TransactionRepository interface {
	Create(ctx context.Context, transaction *domain.Transaction) error
	GetByID(ctx context.Context, transactionID int64) (*domain.Transaction, error)
	GetByOrganizationID(ctx context.Context, organizationID int64, limit, offset int) ([]domain.Transaction, error)
	GetByBillingAccountID(ctx context.Context, billingAccountID int64, limit, offset int) ([]domain.Transaction, error)
	GetByStripePaymentIntentID(ctx context.Context, stripePaymentIntentID string) (*domain.Transaction, error)
	GetByDateRange(ctx context.Context, organizationID int64, startDate, endDate time.Time) ([]domain.Transaction, error)
	GetMonthlySpend(ctx context.Context, organizationID int64, year int, month int) (decimal.Decimal, error)
	Update(ctx context.Context, transaction *domain.Transaction) error
	List(ctx context.Context, limit, offset int) ([]domain.Transaction, error)
	GetBalance(ctx context.Context, billingAccountID int64) (decimal.Decimal, error)
}

// PgxTransactionRepository implements TransactionRepository using pgx
type PgxTransactionRepository struct {
	db *pgxpool.Pool
}

// NewPgxTransactionRepository creates a new PgxTransactionRepository
func NewPgxTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &PgxTransactionRepository{db: db}
}

// Create creates a new transaction
func (r *PgxTransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	query := `
		INSERT INTO transactions (
			organization_id, billing_account_id, type, amount, currency, balance_before,
			balance_after, reference_type, reference_id, related_transaction_id,
			stripe_payment_intent_id, stripe_invoice_id, stripe_charge_id, description,
			metadata, status, processed_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		) RETURNING transaction_id, created_at, updated_at`

	metadataJSON, err := json.Marshal(transaction.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	err = r.db.QueryRow(ctx, query,
		transaction.OrganizationID,
		transaction.BillingAccountID,
		transaction.Type,
		transaction.Amount,
		transaction.Currency,
		transaction.BalanceBefore,
		transaction.BalanceAfter,
		transaction.ReferenceType,
		transaction.ReferenceID,
		transaction.RelatedTransactionID,
		transaction.StripePaymentIntentID,
		transaction.StripeInvoiceID,
		transaction.StripeChargeID,
		transaction.Description,
		metadataJSON,
		transaction.Status,
		transaction.ProcessedAt,
	).Scan(&transaction.TransactionID, &transaction.CreatedAt, &transaction.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

// GetByID retrieves a transaction by ID
func (r *PgxTransactionRepository) GetByID(ctx context.Context, transactionID int64) (*domain.Transaction, error) {
	query := `
		SELECT transaction_id, organization_id, billing_account_id, type, amount, currency,
			   balance_before, balance_after, reference_type, reference_id, related_transaction_id,
			   stripe_payment_intent_id, stripe_invoice_id, stripe_charge_id, description,
			   metadata, status, processed_at, created_at, updated_at
		FROM transactions
		WHERE transaction_id = $1`

	transaction := &domain.Transaction{}
	var metadataJSON []byte

	err := r.db.QueryRow(ctx, query, transactionID).Scan(
		&transaction.TransactionID,
		&transaction.OrganizationID,
		&transaction.BillingAccountID,
		&transaction.Type,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.BalanceBefore,
		&transaction.BalanceAfter,
		&transaction.ReferenceType,
		&transaction.ReferenceID,
		&transaction.RelatedTransactionID,
		&transaction.StripePaymentIntentID,
		&transaction.StripeInvoiceID,
		&transaction.StripeChargeID,
		&transaction.Description,
		&metadataJSON,
		&transaction.Status,
		&transaction.ProcessedAt,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Unmarshal metadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &transaction.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return transaction, nil
}

// GetByOrganizationID retrieves transactions for an organization with pagination
func (r *PgxTransactionRepository) GetByOrganizationID(ctx context.Context, organizationID int64, limit, offset int) ([]domain.Transaction, error) {
	query := `
		SELECT transaction_id, organization_id, billing_account_id, type, amount, currency,
			   balance_before, balance_after, reference_type, reference_id, related_transaction_id,
			   stripe_payment_intent_id, stripe_invoice_id, stripe_charge_id, description,
			   metadata, status, processed_at, created_at, updated_at
		FROM transactions
		WHERE organization_id = $1
		ORDER BY processed_at DESC, created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by organization: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// GetByBillingAccountID retrieves transactions for a billing account with pagination
func (r *PgxTransactionRepository) GetByBillingAccountID(ctx context.Context, billingAccountID int64, limit, offset int) ([]domain.Transaction, error) {
	query := `
		SELECT transaction_id, organization_id, billing_account_id, type, amount, currency,
			   balance_before, balance_after, reference_type, reference_id, related_transaction_id,
			   stripe_payment_intent_id, stripe_invoice_id, stripe_charge_id, description,
			   metadata, status, processed_at, created_at, updated_at
		FROM transactions
		WHERE billing_account_id = $1
		ORDER BY processed_at DESC, created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, billingAccountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by billing account: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// GetByStripePaymentIntentID retrieves a transaction by Stripe payment intent ID
func (r *PgxTransactionRepository) GetByStripePaymentIntentID(ctx context.Context, stripePaymentIntentID string) (*domain.Transaction, error) {
	query := `
		SELECT transaction_id, organization_id, billing_account_id, type, amount, currency,
			   balance_before, balance_after, reference_type, reference_id, related_transaction_id,
			   stripe_payment_intent_id, stripe_invoice_id, stripe_charge_id, description,
			   metadata, status, processed_at, created_at, updated_at
		FROM transactions
		WHERE stripe_payment_intent_id = $1`

	transaction := &domain.Transaction{}
	var metadataJSON []byte

	err := r.db.QueryRow(ctx, query, stripePaymentIntentID).Scan(
		&transaction.TransactionID,
		&transaction.OrganizationID,
		&transaction.BillingAccountID,
		&transaction.Type,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.BalanceBefore,
		&transaction.BalanceAfter,
		&transaction.ReferenceType,
		&transaction.ReferenceID,
		&transaction.RelatedTransactionID,
		&transaction.StripePaymentIntentID,
		&transaction.StripeInvoiceID,
		&transaction.StripeChargeID,
		&transaction.Description,
		&metadataJSON,
		&transaction.Status,
		&transaction.ProcessedAt,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction by Stripe payment intent ID: %w", err)
	}

	// Unmarshal metadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &transaction.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return transaction, nil
}

// GetByDateRange retrieves transactions for an organization within a date range
func (r *PgxTransactionRepository) GetByDateRange(ctx context.Context, organizationID int64, startDate, endDate time.Time) ([]domain.Transaction, error) {
	query := `
		SELECT transaction_id, organization_id, billing_account_id, type, amount, currency,
			   balance_before, balance_after, reference_type, reference_id, related_transaction_id,
			   stripe_payment_intent_id, stripe_invoice_id, stripe_charge_id, description,
			   metadata, status, processed_at, created_at, updated_at
		FROM transactions
		WHERE organization_id = $1 AND processed_at >= $2 AND processed_at <= $3
		ORDER BY processed_at DESC`

	rows, err := r.db.Query(ctx, query, organizationID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by date range: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// GetMonthlySpend calculates the total spend for an organization in a specific month
func (r *PgxTransactionRepository) GetMonthlySpend(ctx context.Context, organizationID int64, year int, month int) (decimal.Decimal, error) {
	query := `
		SELECT COALESCE(SUM(ABS(amount)), 0) as total_spend
		FROM transactions
		WHERE organization_id = $1 
		AND type IN ('usage_charge', 'debit')
		AND status = 'completed'
		AND EXTRACT(YEAR FROM processed_at) = $2
		AND EXTRACT(MONTH FROM processed_at) = $3`

	var totalSpend decimal.Decimal
	err := r.db.QueryRow(ctx, query, organizationID, year, month).Scan(&totalSpend)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get monthly spend: %w", err)
	}

	return totalSpend, nil
}

// Update updates a transaction
func (r *PgxTransactionRepository) Update(ctx context.Context, transaction *domain.Transaction) error {
	query := `
		UPDATE transactions SET
			type = $2,
			amount = $3,
			currency = $4,
			balance_before = $5,
			balance_after = $6,
			reference_type = $7,
			reference_id = $8,
			related_transaction_id = $9,
			stripe_payment_intent_id = $10,
			stripe_invoice_id = $11,
			stripe_charge_id = $12,
			description = $13,
			metadata = $14,
			status = $15,
			processed_at = $16,
			updated_at = NOW()
		WHERE transaction_id = $1
		RETURNING updated_at`

	metadataJSON, err := json.Marshal(transaction.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	err = r.db.QueryRow(ctx, query,
		transaction.TransactionID,
		transaction.Type,
		transaction.Amount,
		transaction.Currency,
		transaction.BalanceBefore,
		transaction.BalanceAfter,
		transaction.ReferenceType,
		transaction.ReferenceID,
		transaction.RelatedTransactionID,
		transaction.StripePaymentIntentID,
		transaction.StripeInvoiceID,
		transaction.StripeChargeID,
		transaction.Description,
		metadataJSON,
		transaction.Status,
		transaction.ProcessedAt,
	).Scan(&transaction.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	return nil
}

// List retrieves a list of transactions with pagination
func (r *PgxTransactionRepository) List(ctx context.Context, limit, offset int) ([]domain.Transaction, error) {
	query := `
		SELECT transaction_id, organization_id, billing_account_id, type, amount, currency,
			   balance_before, balance_after, reference_type, reference_id, related_transaction_id,
			   stripe_payment_intent_id, stripe_invoice_id, stripe_charge_id, description,
			   metadata, status, processed_at, created_at, updated_at
		FROM transactions
		ORDER BY processed_at DESC, created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// GetBalance calculates the current balance for a billing account based on transactions
func (r *PgxTransactionRepository) GetBalance(ctx context.Context, billingAccountID int64) (decimal.Decimal, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0) as balance
		FROM transactions
		WHERE billing_account_id = $1 AND status = 'completed'`

	var balance decimal.Decimal
	err := r.db.QueryRow(ctx, query, billingAccountID).Scan(&balance)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to calculate balance: %w", err)
	}

	return balance, nil
}

// scanTransactions is a helper method to scan multiple transactions from rows
func (r *PgxTransactionRepository) scanTransactions(rows pgx.Rows) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	for rows.Next() {
		transaction := domain.Transaction{}
		var metadataJSON []byte

		err := rows.Scan(
			&transaction.TransactionID,
			&transaction.OrganizationID,
			&transaction.BillingAccountID,
			&transaction.Type,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.BalanceBefore,
			&transaction.BalanceAfter,
			&transaction.ReferenceType,
			&transaction.ReferenceID,
			&transaction.RelatedTransactionID,
			&transaction.StripePaymentIntentID,
			&transaction.StripeInvoiceID,
			&transaction.StripeChargeID,
			&transaction.Description,
			&metadataJSON,
			&transaction.Status,
			&transaction.ProcessedAt,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		// Unmarshal metadata
		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &transaction.Metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}