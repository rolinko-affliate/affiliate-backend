-- #############################################################################
-- ## Publisher Messaging System Migration
-- ## This migration adds support for messaging with publishers in favorite lists
-- ## 
-- ## Features:
-- ## - Conversation sessions between organizations and publishers
-- ## - Messages within conversations with sender identification
-- ## - Integration with favorite publisher lists
-- ## - Support for external service message integration
-- #############################################################################

-- publisher_conversations: Stores conversation sessions between organizations and publishers
CREATE TABLE public.publisher_conversations (
    conversation_id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES public.organizations(organization_id) ON DELETE CASCADE,
    publisher_domain VARCHAR(255) NOT NULL,
    list_id BIGINT REFERENCES public.favorite_publisher_lists(list_id) ON DELETE SET NULL,
    subject VARCHAR(500) NOT NULL,
    status VARCHAR(20) DEFAULT 'active' NOT NULL CHECK (status IN ('active', 'closed', 'archived')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    last_message_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Ensure unique active conversation per organization-publisher pair
    CONSTRAINT unique_active_conversation UNIQUE (organization_id, publisher_domain, status) DEFERRABLE INITIALLY DEFERRED
);

-- Add trigger for automatic timestamp updates
CREATE TRIGGER set_publisher_conversations_timestamp
BEFORE UPDATE ON public.publisher_conversations
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- publisher_messages: Stores individual messages within conversations
CREATE TABLE public.publisher_messages (
    message_id BIGSERIAL PRIMARY KEY,
    conversation_id BIGINT NOT NULL REFERENCES public.publisher_conversations(conversation_id) ON DELETE CASCADE,
    sender_type VARCHAR(20) NOT NULL CHECK (sender_type IN ('organization', 'publisher', 'system')),
    sender_id VARCHAR(255), -- Can be user_id for organization, email for publisher, or system identifier
    content TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'text' NOT NULL CHECK (message_type IN ('text', 'system', 'notification')),
    external_message_id VARCHAR(255), -- For tracking messages from external services
    metadata JSONB, -- Additional metadata (e.g., email headers, external service data)
    sent_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Index for external message tracking
    CONSTRAINT unique_external_message_id UNIQUE (external_message_id) DEFERRABLE INITIALLY DEFERRED
);

-- Indexes for performance
CREATE INDEX idx_publisher_conversations_organization_id ON public.publisher_conversations(organization_id);
CREATE INDEX idx_publisher_conversations_publisher_domain ON public.publisher_conversations(publisher_domain);
CREATE INDEX idx_publisher_conversations_list_id ON public.publisher_conversations(list_id);
CREATE INDEX idx_publisher_conversations_status ON public.publisher_conversations(status);
CREATE INDEX idx_publisher_conversations_created_at ON public.publisher_conversations(created_at);
CREATE INDEX idx_publisher_conversations_last_message_at ON public.publisher_conversations(last_message_at);

CREATE INDEX idx_publisher_messages_conversation_id ON public.publisher_messages(conversation_id);
CREATE INDEX idx_publisher_messages_sender_type ON public.publisher_messages(sender_type);
CREATE INDEX idx_publisher_messages_sent_at ON public.publisher_messages(sent_at);
CREATE INDEX idx_publisher_messages_external_id ON public.publisher_messages(external_message_id);

-- Composite indexes for common queries
CREATE INDEX idx_conversations_org_status_updated ON public.publisher_conversations(organization_id, status, last_message_at DESC);
CREATE INDEX idx_messages_conversation_sent ON public.publisher_messages(conversation_id, sent_at DESC);

-- Add comments for documentation
COMMENT ON TABLE public.publisher_conversations IS 'Stores conversation sessions between organizations and publishers';
COMMENT ON TABLE public.publisher_messages IS 'Stores individual messages within publisher conversations';

COMMENT ON COLUMN public.publisher_conversations.organization_id IS 'Reference to the organization initiating the conversation';
COMMENT ON COLUMN public.publisher_conversations.publisher_domain IS 'Domain name of the publisher being contacted';
COMMENT ON COLUMN public.publisher_conversations.list_id IS 'Optional reference to the favorite list where this publisher was found';
COMMENT ON COLUMN public.publisher_conversations.subject IS 'Subject line of the conversation';
COMMENT ON COLUMN public.publisher_conversations.status IS 'Current status of the conversation';
COMMENT ON COLUMN public.publisher_conversations.last_message_at IS 'Timestamp of the last message in this conversation';

COMMENT ON COLUMN public.publisher_messages.conversation_id IS 'Reference to the conversation this message belongs to';
COMMENT ON COLUMN public.publisher_messages.sender_type IS 'Type of sender: organization, publisher, or system';
COMMENT ON COLUMN public.publisher_messages.sender_id IS 'Identifier of the sender (user_id, email, etc.)';
COMMENT ON COLUMN public.publisher_messages.content IS 'The message content';
COMMENT ON COLUMN public.publisher_messages.message_type IS 'Type of message: text, system notification, etc.';
COMMENT ON COLUMN public.publisher_messages.external_message_id IS 'External service message ID for tracking';
COMMENT ON COLUMN public.publisher_messages.metadata IS 'Additional message metadata in JSON format';

-- Function to update conversation last_message_at when new messages are added
CREATE OR REPLACE FUNCTION update_conversation_last_message()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE public.publisher_conversations 
    SET last_message_at = NEW.sent_at,
        updated_at = CURRENT_TIMESTAMP
    WHERE conversation_id = NEW.conversation_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically update conversation timestamps
CREATE TRIGGER update_conversation_last_message_trigger
    AFTER INSERT ON public.publisher_messages
    FOR EACH ROW
    EXECUTE FUNCTION update_conversation_last_message();