package repository

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PublisherMessagingRepository defines the interface for publisher messaging operations
type PublisherMessagingRepository interface {
	// Conversation operations
	CreateConversation(ctx context.Context, conversation *domain.PublisherConversation) error
	GetConversationByID(ctx context.Context, conversationID int64) (*domain.PublisherConversation, error)
	GetConversationsByOrganization(ctx context.Context, organizationID int64, status string, limit, offset int) ([]domain.PublisherConversation, int, error)
	GetConversationByPublisher(ctx context.Context, organizationID int64, publisherDomain string, status string) (*domain.PublisherConversation, error)
	UpdateConversationStatus(ctx context.Context, conversationID int64, status string) error
	DeleteConversation(ctx context.Context, conversationID int64) error

	// Message operations
	CreateMessage(ctx context.Context, message *domain.PublisherMessage) error
	GetMessagesByConversation(ctx context.Context, conversationID int64, limit, offset int) ([]domain.PublisherMessage, int, error)
	GetMessageByExternalID(ctx context.Context, externalMessageID string) (*domain.PublisherMessage, error)
	DeleteMessage(ctx context.Context, messageID int64) error

	// Combined operations
	GetConversationWithMessages(ctx context.Context, conversationID int64, messageLimit, messageOffset int) (*domain.PublisherConversation, []domain.PublisherMessage, int, error)
}

// publisherMessagingRepository implements PublisherMessagingRepository
type publisherMessagingRepository struct {
	db *pgxpool.Pool
}

// NewPublisherMessagingRepository creates a new publisher messaging repository
func NewPublisherMessagingRepository(db *pgxpool.Pool) PublisherMessagingRepository {
	return &publisherMessagingRepository{db: db}
}

// JSONB type for handling PostgreSQL JSONB columns
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into JSONB", value)
	}
	
	return json.Unmarshal(bytes, j)
}

// Conversation operations

func (r *publisherMessagingRepository) CreateConversation(ctx context.Context, conversation *domain.PublisherConversation) error {
	query := `
		INSERT INTO publisher_conversations (organization_id, publisher_domain, list_id, subject, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING conversation_id, created_at, updated_at, last_message_at`

	err := r.db.QueryRow(ctx, query,
		conversation.OrganizationID,
		conversation.PublisherDomain,
		conversation.ListID,
		conversation.Subject,
		conversation.Status,
	).Scan(
		&conversation.ConversationID,
		&conversation.CreatedAt,
		&conversation.UpdatedAt,
		&conversation.LastMessageAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	return nil
}

func (r *publisherMessagingRepository) GetConversationByID(ctx context.Context, conversationID int64) (*domain.PublisherConversation, error) {
	query := `
		SELECT conversation_id, organization_id, publisher_domain, list_id, subject, status,
		       created_at, updated_at, last_message_at
		FROM publisher_conversations
		WHERE conversation_id = $1`

	var conversation domain.PublisherConversation
	err := r.db.QueryRow(ctx, query, conversationID).Scan(
		&conversation.ConversationID,
		&conversation.OrganizationID,
		&conversation.PublisherDomain,
		&conversation.ListID,
		&conversation.Subject,
		&conversation.Status,
		&conversation.CreatedAt,
		&conversation.UpdatedAt,
		&conversation.LastMessageAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return &conversation, nil
}

func (r *publisherMessagingRepository) GetConversationsByOrganization(ctx context.Context, organizationID int64, status string, limit, offset int) ([]domain.PublisherConversation, int, error) {
	// Build query with optional status filter
	whereClause := "WHERE organization_id = $1"
	args := []interface{}{organizationID}
	argIndex := 2

	if status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM publisher_conversations %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count conversations: %w", err)
	}

	// Main query with pagination
	query := fmt.Sprintf(`
		SELECT conversation_id, organization_id, publisher_domain, list_id, subject, status,
		       created_at, updated_at, last_message_at,
		       (SELECT COUNT(*) FROM publisher_messages WHERE conversation_id = pc.conversation_id) as message_count
		FROM publisher_conversations pc
		%s
		ORDER BY last_message_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get conversations: %w", err)
	}
	defer rows.Close()

	var conversations []domain.PublisherConversation
	for rows.Next() {
		var conversation domain.PublisherConversation
		err := rows.Scan(
			&conversation.ConversationID,
			&conversation.OrganizationID,
			&conversation.PublisherDomain,
			&conversation.ListID,
			&conversation.Subject,
			&conversation.Status,
			&conversation.CreatedAt,
			&conversation.UpdatedAt,
			&conversation.LastMessageAt,
			&conversation.MessageCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan conversation: %w", err)
		}
		conversations = append(conversations, conversation)
	}

	return conversations, total, nil
}

func (r *publisherMessagingRepository) GetConversationByPublisher(ctx context.Context, organizationID int64, publisherDomain string, status string) (*domain.PublisherConversation, error) {
	query := `
		SELECT conversation_id, organization_id, publisher_domain, list_id, subject, status,
		       created_at, updated_at, last_message_at
		FROM publisher_conversations
		WHERE organization_id = $1 AND publisher_domain = $2 AND status = $3`

	var conversation domain.PublisherConversation
	err := r.db.QueryRow(ctx, query, organizationID, publisherDomain, status).Scan(
		&conversation.ConversationID,
		&conversation.OrganizationID,
		&conversation.PublisherDomain,
		&conversation.ListID,
		&conversation.Subject,
		&conversation.Status,
		&conversation.CreatedAt,
		&conversation.UpdatedAt,
		&conversation.LastMessageAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get conversation by publisher: %w", err)
	}

	return &conversation, nil
}

func (r *publisherMessagingRepository) UpdateConversationStatus(ctx context.Context, conversationID int64, status string) error {
	query := `
		UPDATE publisher_conversations
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE conversation_id = $2`

	result, err := r.db.Exec(ctx, query, status, conversationID)
	if err != nil {
		return fmt.Errorf("failed to update conversation status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *publisherMessagingRepository) DeleteConversation(ctx context.Context, conversationID int64) error {
	query := `DELETE FROM publisher_conversations WHERE conversation_id = $1`

	result, err := r.db.Exec(ctx, query, conversationID)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Message operations

func (r *publisherMessagingRepository) CreateMessage(ctx context.Context, message *domain.PublisherMessage) error {
	var metadataJSON JSONB
	if message.Metadata != nil {
		metadataJSON = JSONB(message.Metadata)
	}

	query := `
		INSERT INTO publisher_messages (conversation_id, sender_type, sender_id, content, message_type, external_message_id, metadata, sent_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, COALESCE($8, CURRENT_TIMESTAMP))
		RETURNING message_id, sent_at`

	var sentAt interface{}
	if !message.SentAt.IsZero() {
		sentAt = message.SentAt
	}

	err := r.db.QueryRow(ctx, query,
		message.ConversationID,
		message.SenderType,
		message.SenderID,
		message.Content,
		message.MessageType,
		message.ExternalMessageID,
		metadataJSON,
		sentAt,
	).Scan(
		&message.MessageID,
		&message.SentAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}

func (r *publisherMessagingRepository) GetMessagesByConversation(ctx context.Context, conversationID int64, limit, offset int) ([]domain.PublisherMessage, int, error) {
	// Count query
	countQuery := `SELECT COUNT(*) FROM publisher_messages WHERE conversation_id = $1`
	var total int
	err := r.db.QueryRow(ctx, countQuery, conversationID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count messages: %w", err)
	}

	// Main query with pagination
	query := `
		SELECT message_id, conversation_id, sender_type, sender_id, content, message_type,
		       external_message_id, metadata, sent_at
		FROM publisher_messages
		WHERE conversation_id = $1
		ORDER BY sent_at ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, conversationID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []domain.PublisherMessage
	for rows.Next() {
		var message domain.PublisherMessage
		var metadataJSON JSONB

		err := rows.Scan(
			&message.MessageID,
			&message.ConversationID,
			&message.SenderType,
			&message.SenderID,
			&message.Content,
			&message.MessageType,
			&message.ExternalMessageID,
			&metadataJSON,
			&message.SentAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan message: %w", err)
		}

		if metadataJSON != nil {
			message.Metadata = map[string]interface{}(metadataJSON)
		}

		messages = append(messages, message)
	}

	return messages, total, nil
}

func (r *publisherMessagingRepository) GetMessageByExternalID(ctx context.Context, externalMessageID string) (*domain.PublisherMessage, error) {
	query := `
		SELECT message_id, conversation_id, sender_type, sender_id, content, message_type,
		       external_message_id, metadata, sent_at
		FROM publisher_messages
		WHERE external_message_id = $1`

	var message domain.PublisherMessage
	var metadataJSON JSONB

	err := r.db.QueryRow(ctx, query, externalMessageID).Scan(
		&message.MessageID,
		&message.ConversationID,
		&message.SenderType,
		&message.SenderID,
		&message.Content,
		&message.MessageType,
		&message.ExternalMessageID,
		&metadataJSON,
		&message.SentAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get message by external ID: %w", err)
	}

	if metadataJSON != nil {
		message.Metadata = map[string]interface{}(metadataJSON)
	}

	return &message, nil
}

func (r *publisherMessagingRepository) DeleteMessage(ctx context.Context, messageID int64) error {
	query := `DELETE FROM publisher_messages WHERE message_id = $1`

	result, err := r.db.Exec(ctx, query, messageID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Combined operations

func (r *publisherMessagingRepository) GetConversationWithMessages(ctx context.Context, conversationID int64, messageLimit, messageOffset int) (*domain.PublisherConversation, []domain.PublisherMessage, int, error) {
	// Get conversation
	conversation, err := r.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, nil, 0, err
	}

	// Get messages
	messages, total, err := r.GetMessagesByConversation(ctx, conversationID, messageLimit, messageOffset)
	if err != nil {
		return nil, nil, 0, err
	}

	return conversation, messages, total, nil
}