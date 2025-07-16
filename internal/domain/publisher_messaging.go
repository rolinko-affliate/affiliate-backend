package domain

import (
	"encoding/json"
	"time"
)

// Conversation status constants
const (
	ConversationStatusActive   = "active"
	ConversationStatusClosed   = "closed"
	ConversationStatusArchived = "archived"
)

// Message sender type constants
const (
	SenderTypeOrganization = "organization"
	SenderTypePublisher    = "publisher"
	SenderTypeSystem       = "system"
)

// Message type constants
const (
	MessageTypeText         = "text"
	MessageTypeSystem       = "system"
	MessageTypeNotification = "notification"
)

// PublisherConversation represents a conversation session between an organization and a publisher
type PublisherConversation struct {
	ConversationID  int64      `json:"conversation_id" db:"conversation_id"`
	OrganizationID  int64      `json:"organization_id" db:"organization_id"`
	PublisherDomain string     `json:"publisher_domain" db:"publisher_domain"`
	ListID          *int64     `json:"list_id,omitempty" db:"list_id"`
	Subject         string     `json:"subject" db:"subject"`
	Status          string     `json:"status" db:"status"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	LastMessageAt   time.Time  `json:"last_message_at" db:"last_message_at"`

	// Optional: Include messages when fetching with details
	Messages []PublisherMessage `json:"messages,omitempty" db:"-"`
	
	// Optional: Include publisher details
	Publisher *AnalyticsPublisher `json:"publisher,omitempty" db:"-"`
	
	// Optional: Include list details
	List *FavoritePublisherList `json:"list,omitempty" db:"-"`
	
	// Computed fields
	MessageCount int `json:"message_count,omitempty" db:"-"`
}

// PublisherMessage represents an individual message within a conversation
type PublisherMessage struct {
	MessageID         int64                  `json:"message_id" db:"message_id"`
	ConversationID    int64                  `json:"conversation_id" db:"conversation_id"`
	SenderType        string                 `json:"sender_type" db:"sender_type"`
	SenderID          *string                `json:"sender_id,omitempty" db:"sender_id"`
	Content           string                 `json:"content" db:"content"`
	MessageType       string                 `json:"message_type" db:"message_type"`
	ExternalMessageID *string                `json:"external_message_id,omitempty" db:"external_message_id"`
	Metadata          map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	SentAt            time.Time              `json:"sent_at" db:"sent_at"`
}

// Request/Response models

// CreateConversationRequest represents the request to start a new conversation with a publisher
type CreateConversationRequest struct {
	PublisherDomain string  `json:"publisher_domain" binding:"required,min=1,max=255"`
	ListID          *int64  `json:"list_id,omitempty"`
	Subject         string  `json:"subject" binding:"required,min=1,max=500"`
	InitialMessage  string  `json:"initial_message" binding:"required,min=1,max=5000"`
}

// SendMessageRequest represents the request to send a message in a conversation
type SendMessageRequest struct {
	Content     string                 `json:"content" binding:"required,min=1,max=5000"`
	MessageType *string                `json:"message_type,omitempty" binding:"omitempty,oneof=text system notification"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// AddExternalMessageRequest represents the request to add a message from an external service
type AddExternalMessageRequest struct {
	ConversationID    int64                  `json:"conversation_id" binding:"required"`
	SenderType        string                 `json:"sender_type" binding:"required,oneof=organization publisher system"`
	SenderID          *string                `json:"sender_id,omitempty"`
	Content           string                 `json:"content" binding:"required,min=1,max=5000"`
	MessageType       *string                `json:"message_type,omitempty" binding:"omitempty,oneof=text system notification"`
	ExternalMessageID *string                `json:"external_message_id,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	SentAt            *time.Time             `json:"sent_at,omitempty"`
}

// UpdateConversationStatusRequest represents the request to update conversation status
type UpdateConversationStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active closed archived"`
}

// ConversationListResponse represents the response for listing conversations
type ConversationListResponse struct {
	Conversations []PublisherConversation `json:"conversations"`
	Total         int                     `json:"total"`
	Page          int                     `json:"page"`
	PageSize      int                     `json:"page_size"`
}

// ConversationWithMessagesResponse represents a conversation with its messages
type ConversationWithMessagesResponse struct {
	Conversation PublisherConversation `json:"conversation"`
	Messages     []PublisherMessage    `json:"messages"`
	Total        int                   `json:"total"`
	Page         int                   `json:"page"`
	PageSize     int                   `json:"page_size"`
}

// Validation methods

// Validate validates the CreateConversationRequest
func (r *CreateConversationRequest) Validate() error {
	if err := validateStringLength(r.PublisherDomain, 1, 255); err != nil {
		return err
	}
	if err := validateStringLength(r.Subject, 1, 500); err != nil {
		return err
	}
	return validateStringLength(r.InitialMessage, 1, 5000)
}

// Validate validates the SendMessageRequest
func (r *SendMessageRequest) Validate() error {
	if err := validateStringLength(r.Content, 1, 5000); err != nil {
		return err
	}
	if r.MessageType != nil {
		return validateMessageType(*r.MessageType)
	}
	return nil
}

// Validate validates the AddExternalMessageRequest
func (r *AddExternalMessageRequest) Validate() error {
	if r.ConversationID <= 0 {
		return ErrInvalidInput
	}
	if err := validateSenderType(r.SenderType); err != nil {
		return err
	}
	if err := validateStringLength(r.Content, 1, 5000); err != nil {
		return err
	}
	if r.MessageType != nil {
		return validateMessageType(*r.MessageType)
	}
	return nil
}

// Validate validates the UpdateConversationStatusRequest
func (r *UpdateConversationStatusRequest) Validate() error {
	return validateConversationStatus(r.Status)
}

// Helper validation functions

func validateConversationStatus(status string) error {
	switch status {
	case ConversationStatusActive, ConversationStatusClosed, ConversationStatusArchived:
		return nil
	default:
		return ErrInvalidInput
	}
}

func validateSenderType(senderType string) error {
	switch senderType {
	case SenderTypeOrganization, SenderTypePublisher, SenderTypeSystem:
		return nil
	default:
		return ErrInvalidInput
	}
}

func validateMessageType(messageType string) error {
	switch messageType {
	case MessageTypeText, MessageTypeSystem, MessageTypeNotification:
		return nil
	default:
		return ErrInvalidInput
	}
}

// Custom JSON marshaling for metadata field
func (m *PublisherMessage) MarshalJSON() ([]byte, error) {
	type Alias PublisherMessage
	aux := &struct {
		*Alias
		Metadata json.RawMessage `json:"metadata,omitempty"`
	}{
		Alias: (*Alias)(m),
	}
	
	if m.Metadata != nil {
		metadataBytes, err := json.Marshal(m.Metadata)
		if err != nil {
			return nil, err
		}
		aux.Metadata = metadataBytes
	}
	
	return json.Marshal(aux)
}

// Custom JSON unmarshaling for metadata field
func (m *PublisherMessage) UnmarshalJSON(data []byte) error {
	type Alias PublisherMessage
	aux := &struct {
		*Alias
		Metadata json.RawMessage `json:"metadata,omitempty"`
	}{
		Alias: (*Alias)(m),
	}
	
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	
	if aux.Metadata != nil {
		if err := json.Unmarshal(aux.Metadata, &m.Metadata); err != nil {
			return err
		}
	}
	
	return nil
}