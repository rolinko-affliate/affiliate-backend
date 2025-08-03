package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WebhookEventRepository defines the interface for webhook event operations
type WebhookEventRepository interface {
	Create(ctx context.Context, event *domain.WebhookEvent) error
	GetByID(ctx context.Context, webhookEventID int64) (*domain.WebhookEvent, error)
	GetByStripeEventID(ctx context.Context, stripeEventID string) (*domain.WebhookEvent, error)
	Update(ctx context.Context, event *domain.WebhookEvent) error
	List(ctx context.Context, limit, offset int) ([]domain.WebhookEvent, error)
	GetPendingEvents(ctx context.Context, limit int) ([]domain.WebhookEvent, error)
	GetFailedEvents(ctx context.Context, limit int) ([]domain.WebhookEvent, error)
	IncrementRetryCount(ctx context.Context, webhookEventID int64) error
}

// PgxWebhookEventRepository implements WebhookEventRepository using pgx
type PgxWebhookEventRepository struct {
	db *pgxpool.Pool
}

// NewPgxWebhookEventRepository creates a new PgxWebhookEventRepository
func NewPgxWebhookEventRepository(db *pgxpool.Pool) WebhookEventRepository {
	return &PgxWebhookEventRepository{db: db}
}

// Create creates a new webhook event
func (r *PgxWebhookEventRepository) Create(ctx context.Context, event *domain.WebhookEvent) error {
	query := `
		INSERT INTO webhook_events (
			stripe_event_id, event_type, status, event_data, processed_at,
			error_message, retry_count, organization_id, transaction_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) RETURNING webhook_event_id, created_at, updated_at`

	eventDataJSON, err := json.Marshal(event.EventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	err = r.db.QueryRow(ctx, query,
		event.StripeEventID,
		event.EventType,
		event.Status,
		eventDataJSON,
		event.ProcessedAt,
		event.ErrorMessage,
		event.RetryCount,
		event.OrganizationID,
		event.TransactionID,
	).Scan(&event.WebhookEventID, &event.CreatedAt, &event.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create webhook event: %w", err)
	}

	return nil
}

// GetByID retrieves a webhook event by ID
func (r *PgxWebhookEventRepository) GetByID(ctx context.Context, webhookEventID int64) (*domain.WebhookEvent, error) {
	query := `
		SELECT webhook_event_id, stripe_event_id, event_type, status, event_data,
			   processed_at, error_message, retry_count, organization_id, transaction_id,
			   created_at, updated_at
		FROM webhook_events
		WHERE webhook_event_id = $1`

	event := &domain.WebhookEvent{}
	var eventDataJSON []byte

	err := r.db.QueryRow(ctx, query, webhookEventID).Scan(
		&event.WebhookEventID,
		&event.StripeEventID,
		&event.EventType,
		&event.Status,
		&eventDataJSON,
		&event.ProcessedAt,
		&event.ErrorMessage,
		&event.RetryCount,
		&event.OrganizationID,
		&event.TransactionID,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("webhook event not found")
		}
		return nil, fmt.Errorf("failed to get webhook event: %w", err)
	}

	// Unmarshal event data
	if len(eventDataJSON) > 0 {
		if err := json.Unmarshal(eventDataJSON, &event.EventData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
		}
	}

	return event, nil
}

// GetByStripeEventID retrieves a webhook event by Stripe event ID
func (r *PgxWebhookEventRepository) GetByStripeEventID(ctx context.Context, stripeEventID string) (*domain.WebhookEvent, error) {
	query := `
		SELECT webhook_event_id, stripe_event_id, event_type, status, event_data,
			   processed_at, error_message, retry_count, organization_id, transaction_id,
			   created_at, updated_at
		FROM webhook_events
		WHERE stripe_event_id = $1`

	event := &domain.WebhookEvent{}
	var eventDataJSON []byte

	err := r.db.QueryRow(ctx, query, stripeEventID).Scan(
		&event.WebhookEventID,
		&event.StripeEventID,
		&event.EventType,
		&event.Status,
		&eventDataJSON,
		&event.ProcessedAt,
		&event.ErrorMessage,
		&event.RetryCount,
		&event.OrganizationID,
		&event.TransactionID,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("webhook event not found")
		}
		return nil, fmt.Errorf("failed to get webhook event by Stripe event ID: %w", err)
	}

	// Unmarshal event data
	if len(eventDataJSON) > 0 {
		if err := json.Unmarshal(eventDataJSON, &event.EventData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
		}
	}

	return event, nil
}

// Update updates a webhook event
func (r *PgxWebhookEventRepository) Update(ctx context.Context, event *domain.WebhookEvent) error {
	query := `
		UPDATE webhook_events SET
			event_type = $2,
			status = $3,
			event_data = $4,
			processed_at = $5,
			error_message = $6,
			retry_count = $7,
			organization_id = $8,
			transaction_id = $9,
			updated_at = NOW()
		WHERE webhook_event_id = $1
		RETURNING updated_at`

	eventDataJSON, err := json.Marshal(event.EventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	err = r.db.QueryRow(ctx, query,
		event.WebhookEventID,
		event.EventType,
		event.Status,
		eventDataJSON,
		event.ProcessedAt,
		event.ErrorMessage,
		event.RetryCount,
		event.OrganizationID,
		event.TransactionID,
	).Scan(&event.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update webhook event: %w", err)
	}

	return nil
}

// List retrieves a list of webhook events with pagination
func (r *PgxWebhookEventRepository) List(ctx context.Context, limit, offset int) ([]domain.WebhookEvent, error) {
	query := `
		SELECT webhook_event_id, stripe_event_id, event_type, status, event_data,
			   processed_at, error_message, retry_count, organization_id, transaction_id,
			   created_at, updated_at
		FROM webhook_events
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list webhook events: %w", err)
	}
	defer rows.Close()

	return r.scanWebhookEvents(rows)
}

// GetPendingEvents retrieves pending webhook events for retry processing
func (r *PgxWebhookEventRepository) GetPendingEvents(ctx context.Context, limit int) ([]domain.WebhookEvent, error) {
	query := `
		SELECT webhook_event_id, stripe_event_id, event_type, status, event_data,
			   processed_at, error_message, retry_count, organization_id, transaction_id,
			   created_at, updated_at
		FROM webhook_events
		WHERE status = 'pending' AND retry_count < 5
		ORDER BY created_at ASC
		LIMIT $1`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending webhook events: %w", err)
	}
	defer rows.Close()

	return r.scanWebhookEvents(rows)
}

// GetFailedEvents retrieves failed webhook events
func (r *PgxWebhookEventRepository) GetFailedEvents(ctx context.Context, limit int) ([]domain.WebhookEvent, error) {
	query := `
		SELECT webhook_event_id, stripe_event_id, event_type, status, event_data,
			   processed_at, error_message, retry_count, organization_id, transaction_id,
			   created_at, updated_at
		FROM webhook_events
		WHERE status = 'failed' OR retry_count >= 5
		ORDER BY created_at DESC
		LIMIT $1`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get failed webhook events: %w", err)
	}
	defer rows.Close()

	return r.scanWebhookEvents(rows)
}

// IncrementRetryCount increments the retry count for a webhook event
func (r *PgxWebhookEventRepository) IncrementRetryCount(ctx context.Context, webhookEventID int64) error {
	query := `
		UPDATE webhook_events SET
			retry_count = retry_count + 1,
			updated_at = NOW()
		WHERE webhook_event_id = $1`

	result, err := r.db.Exec(ctx, query, webhookEventID)
	if err != nil {
		return fmt.Errorf("failed to increment retry count: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("webhook event not found")
	}

	return nil
}

// scanWebhookEvents is a helper method to scan multiple webhook events from rows
func (r *PgxWebhookEventRepository) scanWebhookEvents(rows pgx.Rows) ([]domain.WebhookEvent, error) {
	events := make([]domain.WebhookEvent, 0)
	for rows.Next() {
		event := domain.WebhookEvent{}
		var eventDataJSON []byte

		err := rows.Scan(
			&event.WebhookEventID,
			&event.StripeEventID,
			&event.EventType,
			&event.Status,
			&eventDataJSON,
			&event.ProcessedAt,
			&event.ErrorMessage,
			&event.RetryCount,
			&event.OrganizationID,
			&event.TransactionID,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan webhook event: %w", err)
		}

		// Unmarshal event data
		if len(eventDataJSON) > 0 {
			if err := json.Unmarshal(eventDataJSON, &event.EventData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
			}
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating webhook events: %w", err)
	}

	return events, nil
}
