package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AffiliateRepository defines the interface for affiliate data operations
type AffiliateRepository interface {
	CreateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error
	GetAffiliateByID(ctx context.Context, affiliateID int64) (*domain.Affiliate, error)
	UpdateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error
	DeleteAffiliate(ctx context.Context, affiliateID int64) error
	GetAffiliatesByOrganization(ctx context.Context, organizationID int64) ([]*domain.Affiliate, error)
	ListAffiliatesByOrganization(ctx context.Context, organizationID int64, limit, offset int) ([]*domain.Affiliate, error)
	GetAffiliateByEmail(ctx context.Context, email string) (*domain.Affiliate, error)
}

// pgxAffiliateRepository implements AffiliateRepository using pgx
type pgxAffiliateRepository struct {
	db *pgxpool.Pool
}

// NewPgxAffiliateRepository creates a new affiliate repository
func NewPgxAffiliateRepository(db *pgxpool.Pool) AffiliateRepository {
	return &pgxAffiliateRepository{db: db}
}

// CreateAffiliate creates a new affiliate in the database (clean domain model)
func (r *pgxAffiliateRepository) CreateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error {
	query := `INSERT INTO public.affiliates (
		organization_id, name, contact_email, payment_details, status,
		internal_notes, default_currency_id, contact_address, billing_info, labels,
		invoice_amount_threshold, default_payment_terms,
		created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	RETURNING affiliate_id, created_at, updated_at`

	// Handle nullable fields
	var contactEmail, paymentDetails, internalNotes, defaultCurrencyID, contactAddress, billingInfo, labels sql.NullString
	var invoiceAmountThreshold sql.NullFloat64
	var defaultPaymentTerms sql.NullInt32

	if affiliate.ContactEmail != nil {
		contactEmail = sql.NullString{String: *affiliate.ContactEmail, Valid: true}
	}
	if affiliate.PaymentDetails != nil {
		paymentDetails = sql.NullString{String: *affiliate.PaymentDetails, Valid: true}
	}
	if affiliate.InternalNotes != nil {
		internalNotes = sql.NullString{String: *affiliate.InternalNotes, Valid: true}
	}
	if affiliate.DefaultCurrencyID != nil {
		defaultCurrencyID = sql.NullString{String: *affiliate.DefaultCurrencyID, Valid: true}
	}
	if affiliate.ContactAddress != nil {
		contactAddress = sql.NullString{String: *affiliate.ContactAddress, Valid: true}
	}
	if affiliate.BillingInfo != nil {
		billingInfo = sql.NullString{String: *affiliate.BillingInfo, Valid: true}
	}
	if affiliate.Labels != nil {
		labels = sql.NullString{String: *affiliate.Labels, Valid: true}
	}
	if affiliate.InvoiceAmountThreshold != nil {
		invoiceAmountThreshold = sql.NullFloat64{Float64: *affiliate.InvoiceAmountThreshold, Valid: true}
	}
	if affiliate.DefaultPaymentTerms != nil {
		defaultPaymentTerms = sql.NullInt32{Int32: *affiliate.DefaultPaymentTerms, Valid: true}
	}

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		affiliate.OrganizationID,
		affiliate.Name,
		contactEmail,
		paymentDetails,
		affiliate.Status,
		internalNotes,
		defaultCurrencyID,
		contactAddress,
		billingInfo,
		labels,
		invoiceAmountThreshold,
		defaultPaymentTerms,
		now,
		now,
	).Scan(&affiliate.AffiliateID, &affiliate.CreatedAt, &affiliate.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating affiliate: %w", err)
	}

	return nil
}

// GetAffiliateByID retrieves an affiliate by ID
func (r *pgxAffiliateRepository) GetAffiliateByID(ctx context.Context, affiliateID int64) (*domain.Affiliate, error) {
	query := `SELECT 
		affiliate_id, organization_id, name, contact_email, payment_details, status,
		internal_notes, default_currency_id, contact_address, billing_info, labels,
		invoice_amount_threshold, default_payment_terms,
		created_at, updated_at
	FROM public.affiliates 
	WHERE affiliate_id = $1`

	affiliate := &domain.Affiliate{}
	var contactEmail, paymentDetails, internalNotes, defaultCurrencyID, contactAddress, billingInfo, labels sql.NullString
	var invoiceAmountThreshold sql.NullFloat64
	var defaultPaymentTerms sql.NullInt32

	err := r.db.QueryRow(ctx, query, affiliateID).Scan(
		&affiliate.AffiliateID,
		&affiliate.OrganizationID,
		&affiliate.Name,
		&contactEmail,
		&paymentDetails,
		&affiliate.Status,
		&internalNotes,
		&defaultCurrencyID,
		&contactAddress,
		&billingInfo,
		&labels,
		&invoiceAmountThreshold,
		&defaultPaymentTerms,
		&affiliate.CreatedAt,
		&affiliate.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error getting affiliate by ID: %w", err)
	}

	// Handle nullable fields
	if contactEmail.Valid {
		affiliate.ContactEmail = &contactEmail.String
	}
	if paymentDetails.Valid {
		affiliate.PaymentDetails = &paymentDetails.String
	}
	if internalNotes.Valid {
		affiliate.InternalNotes = &internalNotes.String
	}
	if defaultCurrencyID.Valid {
		affiliate.DefaultCurrencyID = &defaultCurrencyID.String
	}
	if contactAddress.Valid {
		affiliate.ContactAddress = &contactAddress.String
	}
	if billingInfo.Valid {
		affiliate.BillingInfo = &billingInfo.String
	}
	if labels.Valid {
		affiliate.Labels = &labels.String
	}
	if invoiceAmountThreshold.Valid {
		affiliate.InvoiceAmountThreshold = &invoiceAmountThreshold.Float64
	}
	if defaultPaymentTerms.Valid {
		affiliate.DefaultPaymentTerms = &defaultPaymentTerms.Int32
	}

	return affiliate, nil
}

// UpdateAffiliate updates an existing affiliate
func (r *pgxAffiliateRepository) UpdateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error {
	query := `UPDATE public.affiliates SET 
		organization_id = $2, name = $3, contact_email = $4, payment_details = $5, 
		status = $6, internal_notes = $7, default_currency_id = $8, contact_address = $9,
		billing_info = $10, labels = $11, invoice_amount_threshold = $12, 
		default_payment_terms = $13, updated_at = $14
	WHERE affiliate_id = $1`

	// Handle nullable fields
	var contactEmail, paymentDetails, internalNotes, defaultCurrencyID, contactAddress, billingInfo, labels sql.NullString
	var invoiceAmountThreshold sql.NullFloat64
	var defaultPaymentTerms sql.NullInt32

	if affiliate.ContactEmail != nil {
		contactEmail = sql.NullString{String: *affiliate.ContactEmail, Valid: true}
	}
	if affiliate.PaymentDetails != nil {
		paymentDetails = sql.NullString{String: *affiliate.PaymentDetails, Valid: true}
	}
	if affiliate.InternalNotes != nil {
		internalNotes = sql.NullString{String: *affiliate.InternalNotes, Valid: true}
	}
	if affiliate.DefaultCurrencyID != nil {
		defaultCurrencyID = sql.NullString{String: *affiliate.DefaultCurrencyID, Valid: true}
	}
	if affiliate.ContactAddress != nil {
		contactAddress = sql.NullString{String: *affiliate.ContactAddress, Valid: true}
	}
	if affiliate.BillingInfo != nil {
		billingInfo = sql.NullString{String: *affiliate.BillingInfo, Valid: true}
	}
	if affiliate.Labels != nil {
		labels = sql.NullString{String: *affiliate.Labels, Valid: true}
	}
	if affiliate.InvoiceAmountThreshold != nil {
		invoiceAmountThreshold = sql.NullFloat64{Float64: *affiliate.InvoiceAmountThreshold, Valid: true}
	}
	if affiliate.DefaultPaymentTerms != nil {
		defaultPaymentTerms = sql.NullInt32{Int32: *affiliate.DefaultPaymentTerms, Valid: true}
	}

	now := time.Now()
	_, err := r.db.Exec(ctx, query,
		affiliate.AffiliateID,
		affiliate.OrganizationID,
		affiliate.Name,
		contactEmail,
		paymentDetails,
		affiliate.Status,
		internalNotes,
		defaultCurrencyID,
		contactAddress,
		billingInfo,
		labels,
		invoiceAmountThreshold,
		defaultPaymentTerms,
		now,
	)

	if err != nil {
		return fmt.Errorf("error updating affiliate: %w", err)
	}

	affiliate.UpdatedAt = now
	return nil
}

// DeleteAffiliate deletes an affiliate by ID
func (r *pgxAffiliateRepository) DeleteAffiliate(ctx context.Context, affiliateID int64) error {
	query := `DELETE FROM public.affiliates WHERE affiliate_id = $1`

	_, err := r.db.Exec(ctx, query, affiliateID)
	if err != nil {
		return fmt.Errorf("error deleting affiliate: %w", err)
	}

	return nil
}

// GetAffiliatesByOrganization retrieves all affiliates for an organization
func (r *pgxAffiliateRepository) GetAffiliatesByOrganization(ctx context.Context, organizationID int64) ([]*domain.Affiliate, error) {
	query := `SELECT 
		affiliate_id, organization_id, name, contact_email, payment_details, status,
		internal_notes, default_currency_id, contact_address, billing_info, labels,
		invoice_amount_threshold, default_payment_terms,
		created_at, updated_at
	FROM public.affiliates 
	WHERE organization_id = $1
	ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, organizationID)
	if err != nil {
		return nil, fmt.Errorf("error querying affiliates by organization: %w", err)
	}
	defer rows.Close()

	var affiliates []*domain.Affiliate

	for rows.Next() {
		affiliate := &domain.Affiliate{}
		var contactEmail, paymentDetails, internalNotes, defaultCurrencyID, contactAddress, billingInfo, labels sql.NullString
		var invoiceAmountThreshold sql.NullFloat64
		var defaultPaymentTerms sql.NullInt32

		err := rows.Scan(
			&affiliate.AffiliateID,
			&affiliate.OrganizationID,
			&affiliate.Name,
			&contactEmail,
			&paymentDetails,
			&affiliate.Status,
			&internalNotes,
			&defaultCurrencyID,
			&contactAddress,
			&billingInfo,
			&labels,
			&invoiceAmountThreshold,
			&defaultPaymentTerms,
			&affiliate.CreatedAt,
			&affiliate.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning affiliate row: %w", err)
		}

		// Handle nullable fields
		if contactEmail.Valid {
			affiliate.ContactEmail = &contactEmail.String
		}
		if paymentDetails.Valid {
			affiliate.PaymentDetails = &paymentDetails.String
		}
		if internalNotes.Valid {
			affiliate.InternalNotes = &internalNotes.String
		}
		if defaultCurrencyID.Valid {
			affiliate.DefaultCurrencyID = &defaultCurrencyID.String
		}
		if contactAddress.Valid {
			affiliate.ContactAddress = &contactAddress.String
		}
		if billingInfo.Valid {
			affiliate.BillingInfo = &billingInfo.String
		}
		if labels.Valid {
			affiliate.Labels = &labels.String
		}
		if invoiceAmountThreshold.Valid {
			affiliate.InvoiceAmountThreshold = &invoiceAmountThreshold.Float64
		}
		if defaultPaymentTerms.Valid {
			affiliate.DefaultPaymentTerms = &defaultPaymentTerms.Int32
		}

		affiliates = append(affiliates, affiliate)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating affiliate rows: %w", err)
	}

	return affiliates, nil
}

// GetAffiliateByEmail retrieves an affiliate by email
func (r *pgxAffiliateRepository) GetAffiliateByEmail(ctx context.Context, email string) (*domain.Affiliate, error) {
	query := `SELECT 
		affiliate_id, organization_id, name, contact_email, payment_details, status,
		internal_notes, default_currency_id, contact_address, billing_info, labels,
		invoice_amount_threshold, default_payment_terms,
		created_at, updated_at
	FROM public.affiliates 
	WHERE contact_email = $1`

	affiliate := &domain.Affiliate{}
	var contactEmail, paymentDetails, internalNotes, defaultCurrencyID, contactAddress, billingInfo, labels sql.NullString
	var invoiceAmountThreshold sql.NullFloat64
	var defaultPaymentTerms sql.NullInt32

	err := r.db.QueryRow(ctx, query, email).Scan(
		&affiliate.AffiliateID,
		&affiliate.OrganizationID,
		&affiliate.Name,
		&contactEmail,
		&paymentDetails,
		&affiliate.Status,
		&internalNotes,
		&defaultCurrencyID,
		&contactAddress,
		&billingInfo,
		&labels,
		&invoiceAmountThreshold,
		&defaultPaymentTerms,
		&affiliate.CreatedAt,
		&affiliate.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error getting affiliate by email: %w", err)
	}

	// Handle nullable fields
	if contactEmail.Valid {
		affiliate.ContactEmail = &contactEmail.String
	}
	if paymentDetails.Valid {
		affiliate.PaymentDetails = &paymentDetails.String
	}
	if internalNotes.Valid {
		affiliate.InternalNotes = &internalNotes.String
	}
	if defaultCurrencyID.Valid {
		affiliate.DefaultCurrencyID = &defaultCurrencyID.String
	}
	if contactAddress.Valid {
		affiliate.ContactAddress = &contactAddress.String
	}
	if billingInfo.Valid {
		affiliate.BillingInfo = &billingInfo.String
	}
	if labels.Valid {
		affiliate.Labels = &labels.String
	}
	if invoiceAmountThreshold.Valid {
		affiliate.InvoiceAmountThreshold = &invoiceAmountThreshold.Float64
	}
	if defaultPaymentTerms.Valid {
		affiliate.DefaultPaymentTerms = &defaultPaymentTerms.Int32
	}

	return affiliate, nil
}

// ListAffiliatesByOrganization retrieves affiliates by organization with pagination
func (r *pgxAffiliateRepository) ListAffiliatesByOrganization(ctx context.Context, organizationID int64, limit, offset int) ([]*domain.Affiliate, error) {
	query := `SELECT 
		affiliate_id, organization_id, name, contact_email, payment_details, status,
		internal_notes, default_currency_id, contact_address, billing_info, labels,
		invoice_amount_threshold, default_payment_terms,
		created_at, updated_at
	FROM public.affiliates 
	WHERE organization_id = $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying affiliates by organization: %w", err)
	}
	defer rows.Close()

	var affiliates []*domain.Affiliate
	for rows.Next() {
		affiliate := &domain.Affiliate{}
		var contactEmail, paymentDetails, internalNotes, defaultCurrencyID, contactAddress, billingInfo, labels sql.NullString
		var invoiceAmountThreshold sql.NullFloat64
		var defaultPaymentTerms sql.NullInt32

		err := rows.Scan(
			&affiliate.AffiliateID,
			&affiliate.OrganizationID,
			&affiliate.Name,
			&contactEmail,
			&paymentDetails,
			&affiliate.Status,
			&internalNotes,
			&defaultCurrencyID,
			&contactAddress,
			&billingInfo,
			&labels,
			&invoiceAmountThreshold,
			&defaultPaymentTerms,
			&affiliate.CreatedAt,
			&affiliate.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning affiliate row: %w", err)
		}

		// Handle nullable fields
		if contactEmail.Valid {
			affiliate.ContactEmail = &contactEmail.String
		}
		if paymentDetails.Valid {
			affiliate.PaymentDetails = &paymentDetails.String
		}
		if internalNotes.Valid {
			affiliate.InternalNotes = &internalNotes.String
		}
		if defaultCurrencyID.Valid {
			affiliate.DefaultCurrencyID = &defaultCurrencyID.String
		}
		if contactAddress.Valid {
			affiliate.ContactAddress = &contactAddress.String
		}
		if billingInfo.Valid {
			affiliate.BillingInfo = &billingInfo.String
		}
		if labels.Valid {
			affiliate.Labels = &labels.String
		}
		if invoiceAmountThreshold.Valid {
			affiliate.InvoiceAmountThreshold = &invoiceAmountThreshold.Float64
		}
		if defaultPaymentTerms.Valid {
			affiliate.DefaultPaymentTerms = &defaultPaymentTerms.Int32
		}

		affiliates = append(affiliates, affiliate)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating affiliate rows: %w", err)
	}

	return affiliates, nil
}
