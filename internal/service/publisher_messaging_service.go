package service

import (
	"context"
	"fmt"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/repository"
)

// PublisherMessagingService defines the interface for publisher messaging business logic
type PublisherMessagingService interface {
	// Conversation operations
	CreateConversation(ctx context.Context, organizationID int64, req *domain.CreateConversationRequest) (*domain.PublisherConversation, error)
	GetConversation(ctx context.Context, organizationID int64, conversationID int64) (*domain.PublisherConversation, error)
	GetConversations(ctx context.Context, organizationID int64, status string, page, pageSize int) (*domain.ConversationListResponse, error)
	UpdateConversationStatus(ctx context.Context, organizationID int64, conversationID int64, req *domain.UpdateConversationStatusRequest) error
	DeleteConversation(ctx context.Context, organizationID int64, conversationID int64) error

	// Message operations
	SendMessage(ctx context.Context, organizationID int64, conversationID int64, req *domain.SendMessageRequest) (*domain.PublisherMessage, error)
	GetConversationMessages(ctx context.Context, organizationID int64, conversationID int64, page, pageSize int) (*domain.ConversationWithMessagesResponse, error)
	AddExternalMessage(ctx context.Context, req *domain.AddExternalMessageRequest) (*domain.PublisherMessage, error)

	// Utility operations
	FindOrCreateConversation(ctx context.Context, organizationID int64, publisherDomain string, subject string) (*domain.PublisherConversation, error)
}

// publisherMessagingService implements PublisherMessagingService
type publisherMessagingService struct {
	messagingRepo repository.PublisherMessagingRepository
	analyticsRepo repository.AnalyticsRepository
	favListRepo   repository.FavoritePublisherListRepository
}

// NewPublisherMessagingService creates a new publisher messaging service
func NewPublisherMessagingService(
	messagingRepo repository.PublisherMessagingRepository,
	analyticsRepo repository.AnalyticsRepository,
	favListRepo repository.FavoritePublisherListRepository,
) PublisherMessagingService {
	return &publisherMessagingService{
		messagingRepo: messagingRepo,
		analyticsRepo: analyticsRepo,
		favListRepo:   favListRepo,
	}
}

// Conversation operations

func (s *publisherMessagingService) CreateConversation(ctx context.Context, organizationID int64, req *domain.CreateConversationRequest) (*domain.PublisherConversation, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Check if publisher exists in analytics data
	_, err := s.analyticsRepo.GetPublisherByDomain(ctx, req.PublisherDomain)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, fmt.Errorf("publisher with domain %s not found", req.PublisherDomain)
		}
		return nil, fmt.Errorf("failed to verify publisher: %w", err)
	}

	// Validate list_id if provided
	if req.ListID != nil {
		_, err := s.favListRepo.GetListByID(ctx, *req.ListID)
		if err != nil {
			if err == domain.ErrNotFound {
				return nil, fmt.Errorf("favorite list with ID %d not found", *req.ListID)
			}
			return nil, fmt.Errorf("failed to verify favorite list: %w", err)
		}
	}

	// Check if there's already an active conversation with this publisher
	existingConv, err := s.messagingRepo.GetConversationByPublisher(ctx, organizationID, req.PublisherDomain, domain.ConversationStatusActive)
	if err != nil && err != domain.ErrNotFound {
		return nil, fmt.Errorf("failed to check existing conversation: %w", err)
	}
	if existingConv != nil {
		return nil, fmt.Errorf("active conversation with publisher %s already exists", req.PublisherDomain)
	}

	// Create conversation
	conversation := &domain.PublisherConversation{
		OrganizationID:  organizationID,
		PublisherDomain: req.PublisherDomain,
		ListID:          req.ListID,
		Subject:         req.Subject,
		Status:          domain.ConversationStatusActive,
	}

	err = s.messagingRepo.CreateConversation(ctx, conversation)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	// Create initial message
	initialMessage := &domain.PublisherMessage{
		ConversationID: conversation.ConversationID,
		SenderType:     domain.SenderTypeOrganization,
		Content:        req.InitialMessage,
		MessageType:    domain.MessageTypeText,
	}

	err = s.messagingRepo.CreateMessage(ctx, initialMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to create initial message: %w", err)
	}

	return conversation, nil
}

func (s *publisherMessagingService) GetConversation(ctx context.Context, organizationID int64, conversationID int64) (*domain.PublisherConversation, error) {
	conversation, err := s.messagingRepo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	// Verify organization ownership
	if conversation.OrganizationID != organizationID {
		return nil, domain.ErrNotFound
	}

	// Optionally load publisher details
	if publisher, err := s.analyticsRepo.GetPublisherByDomain(ctx, conversation.PublisherDomain); err == nil {
		conversation.Publisher = publisher
	}

	// Optionally load list details
	if conversation.ListID != nil {
		if list, err := s.favListRepo.GetListByID(ctx, *conversation.ListID); err == nil {
			conversation.List = list
		}
	}

	return conversation, nil
}

func (s *publisherMessagingService) GetConversations(ctx context.Context, organizationID int64, status string, page, pageSize int) (*domain.ConversationListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	conversations, total, err := s.messagingRepo.GetConversationsByOrganization(ctx, organizationID, status, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}

	// Optionally load publisher details for each conversation
	for i := range conversations {
		if publisher, err := s.analyticsRepo.GetPublisherByDomain(ctx, conversations[i].PublisherDomain); err == nil {
			conversations[i].Publisher = publisher
		}
	}

	return &domain.ConversationListResponse{
		Conversations: conversations,
		Total:         total,
		Page:          page,
		PageSize:      pageSize,
	}, nil
}

func (s *publisherMessagingService) UpdateConversationStatus(ctx context.Context, organizationID int64, conversationID int64, req *domain.UpdateConversationStatusRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	// Verify conversation exists and belongs to organization
	conversation, err := s.messagingRepo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return err
	}

	if conversation.OrganizationID != organizationID {
		return domain.ErrNotFound
	}

	return s.messagingRepo.UpdateConversationStatus(ctx, conversationID, req.Status)
}

func (s *publisherMessagingService) DeleteConversation(ctx context.Context, organizationID int64, conversationID int64) error {
	// Verify conversation exists and belongs to organization
	conversation, err := s.messagingRepo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return err
	}

	if conversation.OrganizationID != organizationID {
		return domain.ErrNotFound
	}

	return s.messagingRepo.DeleteConversation(ctx, conversationID)
}

// Message operations

func (s *publisherMessagingService) SendMessage(ctx context.Context, organizationID int64, conversationID int64, req *domain.SendMessageRequest) (*domain.PublisherMessage, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Verify conversation exists and belongs to organization
	conversation, err := s.messagingRepo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	if conversation.OrganizationID != organizationID {
		return nil, domain.ErrNotFound
	}

	// Check if conversation is active
	if conversation.Status != domain.ConversationStatusActive {
		return nil, fmt.Errorf("cannot send message to %s conversation", conversation.Status)
	}

	messageType := domain.MessageTypeText
	if req.MessageType != nil {
		messageType = *req.MessageType
	}

	message := &domain.PublisherMessage{
		ConversationID: conversationID,
		SenderType:     domain.SenderTypeOrganization,
		Content:        req.Content,
		MessageType:    messageType,
		Metadata:       req.Metadata,
	}

	err = s.messagingRepo.CreateMessage(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return message, nil
}

func (s *publisherMessagingService) GetConversationMessages(ctx context.Context, organizationID int64, conversationID int64, page, pageSize int) (*domain.ConversationWithMessagesResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	// Verify conversation exists and belongs to organization
	conversation, err := s.messagingRepo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	if conversation.OrganizationID != organizationID {
		return nil, domain.ErrNotFound
	}

	offset := (page - 1) * pageSize
	messages, total, err := s.messagingRepo.GetMessagesByConversation(ctx, conversationID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	// Load additional conversation details
	if publisher, err := s.analyticsRepo.GetPublisherByDomain(ctx, conversation.PublisherDomain); err == nil {
		conversation.Publisher = publisher
	}

	if conversation.ListID != nil {
		if list, err := s.favListRepo.GetListByID(ctx, *conversation.ListID); err == nil {
			conversation.List = list
		}
	}

	return &domain.ConversationWithMessagesResponse{
		Conversation: *conversation,
		Messages:     messages,
		Total:        total,
		Page:         page,
		PageSize:     pageSize,
	}, nil
}

func (s *publisherMessagingService) AddExternalMessage(ctx context.Context, req *domain.AddExternalMessageRequest) (*domain.PublisherMessage, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Verify conversation exists
	_, err := s.messagingRepo.GetConversationByID(ctx, req.ConversationID)
	if err != nil {
		return nil, err
	}

	// Check for duplicate external message ID
	if req.ExternalMessageID != nil {
		existing, err := s.messagingRepo.GetMessageByExternalID(ctx, *req.ExternalMessageID)
		if err != nil && err != domain.ErrNotFound {
			return nil, fmt.Errorf("failed to check existing message: %w", err)
		}
		if existing != nil {
			return existing, nil // Return existing message if duplicate
		}
	}

	messageType := domain.MessageTypeText
	if req.MessageType != nil {
		messageType = *req.MessageType
	}

	message := &domain.PublisherMessage{
		ConversationID:    req.ConversationID,
		SenderType:        req.SenderType,
		SenderID:          req.SenderID,
		Content:           req.Content,
		MessageType:       messageType,
		ExternalMessageID: req.ExternalMessageID,
		Metadata:          req.Metadata,
	}

	if req.SentAt != nil {
		message.SentAt = *req.SentAt
	}

	err = s.messagingRepo.CreateMessage(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("failed to create external message: %w", err)
	}

	return message, nil
}

// Utility operations

func (s *publisherMessagingService) FindOrCreateConversation(ctx context.Context, organizationID int64, publisherDomain string, subject string) (*domain.PublisherConversation, error) {
	// Try to find existing active conversation
	conversation, err := s.messagingRepo.GetConversationByPublisher(ctx, organizationID, publisherDomain, domain.ConversationStatusActive)
	if err != nil && err != domain.ErrNotFound {
		return nil, fmt.Errorf("failed to check existing conversation: %w", err)
	}

	if conversation != nil {
		return conversation, nil
	}

	// Create new conversation if none exists
	req := &domain.CreateConversationRequest{
		PublisherDomain: publisherDomain,
		Subject:         subject,
		InitialMessage:  fmt.Sprintf("Starting conversation with %s", publisherDomain),
	}

	return s.CreateConversation(ctx, organizationID, req)
}