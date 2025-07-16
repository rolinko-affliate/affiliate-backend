package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/service"
)

// PublisherMessagingHandler handles HTTP requests for publisher messaging
type PublisherMessagingHandler struct {
	messagingService service.PublisherMessagingService
}

// NewPublisherMessagingHandler creates a new publisher messaging handler
func NewPublisherMessagingHandler(messagingService service.PublisherMessagingService) *PublisherMessagingHandler {
	return &PublisherMessagingHandler{
		messagingService: messagingService,
	}
}

// CreateConversation creates a new conversation with a publisher
// @Summary Create a new conversation with a publisher
// @Description Initiates a new conversation with a publisher from a favorite list
// @Tags Publisher Messaging
// @Accept json
// @Produce json
// @Param request body domain.CreateConversationRequest true "Conversation creation request"
// @Success 201 {object} domain.PublisherConversation "Conversation created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body or validation error"
// @Failure 401 {object} ErrorResponse "Organization ID not found in context"
// @Failure 404 {object} ErrorResponse "Publisher or favorite list not found"
// @Failure 409 {object} ErrorResponse "Active conversation with publisher already exists"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/publisher-messaging/conversations [post]
func (h *PublisherMessagingHandler) CreateConversation(c *gin.Context) {
	userOrgID, exists := c.Get("organizationID")
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Organization ID not found in context",
			Details: "Please ensure you are properly authenticated",
		})
		return
	}
	organizationID := userOrgID.(int64)

	var req domain.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	conversation, err := h.messagingService.CreateConversation(c.Request.Context(), organizationID, &req)
	if err != nil {
		switch err.Error() {
		case "publisher not found":
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Publisher not found",
				Details: "The specified publisher domain does not exist",
			})
		case "active conversation with publisher already exists":
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "Active conversation already exists",
				Details: "There is already an active conversation with this publisher",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to create conversation",
				Details: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, conversation)
}

// GetConversations retrieves conversations for the organization
// @Summary Get conversations for organization
// @Description Retrieves a paginated list of conversations for the organization
// @Tags Publisher Messaging
// @Produce json
// @Param status query string false "Filter by conversation status (active, closed)"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20)"
// @Success 200 {object} domain.ConversationListResponse "Conversations retrieved successfully"
// @Failure 401 {object} ErrorResponse "Organization ID not found in context"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/publisher-messaging/conversations [get]
func (h *PublisherMessagingHandler) GetConversations(c *gin.Context) {
	userOrgID, exists := c.Get("organizationID")
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Organization ID not found in context",
			Details: "Please ensure you are properly authenticated",
		})
		return
	}
	organizationID := userOrgID.(int64)

	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	response, err := h.messagingService.GetConversations(c.Request.Context(), organizationID, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve conversations",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetConversation retrieves a specific conversation with messages
// @Summary Get conversation with messages
// @Description Retrieves a specific conversation along with its messages
// @Tags Publisher Messaging
// @Produce json
// @Param conversation_id path int true "Conversation ID"
// @Success 200 {object} domain.ConversationWithMessagesResponse "Conversation retrieved successfully"
// @Failure 400 {object} ErrorResponse "Invalid conversation ID"
// @Failure 401 {object} ErrorResponse "Organization ID not found in context"
// @Failure 404 {object} ErrorResponse "Conversation not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/publisher-messaging/conversations/{conversation_id} [get]
func (h *PublisherMessagingHandler) GetConversation(c *gin.Context) {
	userOrgID, exists := c.Get("organizationID")
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Organization ID not found in context",
			Details: "Please ensure you are properly authenticated",
		})
		return
	}
	organizationID := userOrgID.(int64)

	conversationIDStr := c.Param("conversation_id")
	conversationID, err := strconv.ParseInt(conversationIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid conversation ID",
			Details: "Conversation ID must be a valid integer",
		})
		return
	}

	response, err := h.messagingService.GetConversationMessages(c.Request.Context(), organizationID, conversationID, 1, 100)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Conversation not found",
				Details: "The specified conversation does not exist or you don't have access to it",
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to retrieve conversation",
				Details: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// AddMessage adds a message to a conversation
// @Summary Add message to conversation
// @Description Adds a new message to an existing conversation
// @Tags Publisher Messaging
// @Accept json
// @Produce json
// @Param conversation_id path int true "Conversation ID"
// @Param request body domain.SendMessageRequest true "Message request"
// @Success 201 {object} domain.PublisherMessage "Message added successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body or conversation ID"
// @Failure 401 {object} ErrorResponse "Organization ID not found in context"
// @Failure 404 {object} ErrorResponse "Conversation not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/publisher-messaging/conversations/{conversation_id}/messages [post]
func (h *PublisherMessagingHandler) AddMessage(c *gin.Context) {
	userOrgID, exists := c.Get("organizationID")
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Organization ID not found in context",
			Details: "Please ensure you are properly authenticated",
		})
		return
	}
	organizationID := userOrgID.(int64)

	conversationIDStr := c.Param("conversation_id")
	conversationID, err := strconv.ParseInt(conversationIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid conversation ID",
			Details: "Conversation ID must be a valid integer",
		})
		return
	}

	var req domain.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	message, err := h.messagingService.SendMessage(c.Request.Context(), organizationID, conversationID, &req)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Conversation not found",
				Details: "The specified conversation does not exist or you don't have access to it",
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to add message",
				Details: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, message)
}

// UpdateConversationStatus updates the status of a conversation
// @Summary Update conversation status
// @Description Updates the status of a conversation (e.g., close conversation)
// @Tags Publisher Messaging
// @Accept json
// @Produce json
// @Param conversation_id path int true "Conversation ID"
// @Param request body domain.UpdateConversationStatusRequest true "Status update request"
// @Success 200 {object} domain.PublisherConversation "Conversation status updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body or conversation ID"
// @Failure 401 {object} ErrorResponse "Organization ID not found in context"
// @Failure 404 {object} ErrorResponse "Conversation not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/publisher-messaging/conversations/{conversation_id}/status [put]
func (h *PublisherMessagingHandler) UpdateConversationStatus(c *gin.Context) {
	userOrgID, exists := c.Get("organizationID")
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Organization ID not found in context",
			Details: "Please ensure you are properly authenticated",
		})
		return
	}
	organizationID := userOrgID.(int64)

	conversationIDStr := c.Param("conversation_id")
	conversationID, err := strconv.ParseInt(conversationIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid conversation ID",
			Details: "Conversation ID must be a valid integer",
		})
		return
	}

	var req domain.UpdateConversationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	err = h.messagingService.UpdateConversationStatus(c.Request.Context(), organizationID, conversationID, &req)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Conversation not found",
				Details: "The specified conversation does not exist or you don't have access to it",
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to update conversation status",
				Details: err.Error(),
			})
		}
		return
	}

	// Get updated conversation
	conversation, err := h.messagingService.GetConversation(c.Request.Context(), organizationID, conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve updated conversation",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, conversation)
}

// AddExternalMessage allows external services to add messages to conversations
// @Summary Add external message to conversation
// @Description Allows external services to add messages to existing conversations (e.g., publisher replies)
// @Tags Publisher Messaging
// @Accept json
// @Produce json
// @Param conversation_id path int true "Conversation ID"
// @Param request body domain.AddExternalMessageRequest true "External message request"
// @Success 201 {object} domain.PublisherMessage "Message added successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body or conversation ID"
// @Failure 404 {object} ErrorResponse "Conversation not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/publisher-messaging/conversations/{conversation_id}/external-messages [post]
func (h *PublisherMessagingHandler) AddExternalMessage(c *gin.Context) {
	conversationIDStr := c.Param("conversation_id")
	conversationID, err := strconv.ParseInt(conversationIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid conversation ID",
			Details: "Conversation ID must be a valid integer",
		})
		return
	}

	var req domain.AddExternalMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Set conversation ID from URL parameter
	req.ConversationID = conversationID

	message, err := h.messagingService.AddExternalMessage(c.Request.Context(), &req)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Conversation not found",
				Details: "The specified conversation does not exist",
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to add external message",
				Details: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, message)
}

// DeleteConversation deletes a conversation and all its messages
// @Summary Delete conversation
// @Description Deletes a conversation and all associated messages
// @Tags Publisher Messaging
// @Produce json
// @Param conversation_id path int true "Conversation ID"
// @Success 204 "Conversation deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid conversation ID"
// @Failure 401 {object} ErrorResponse "Organization ID not found in context"
// @Failure 404 {object} ErrorResponse "Conversation not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/publisher-messaging/conversations/{conversation_id} [delete]
func (h *PublisherMessagingHandler) DeleteConversation(c *gin.Context) {
	userOrgID, exists := c.Get("organizationID")
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Organization ID not found in context",
			Details: "Please ensure you are properly authenticated",
		})
		return
	}
	organizationID := userOrgID.(int64)

	conversationIDStr := c.Param("conversation_id")
	conversationID, err := strconv.ParseInt(conversationIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid conversation ID",
			Details: "Conversation ID must be a valid integer",
		})
		return
	}

	err = h.messagingService.DeleteConversation(c.Request.Context(), organizationID, conversationID)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Conversation not found",
				Details: "The specified conversation does not exist or you don't have access to it",
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to delete conversation",
				Details: err.Error(),
			})
		}
		return
	}

	c.Status(http.StatusNoContent)
}