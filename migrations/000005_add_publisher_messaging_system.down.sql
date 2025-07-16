-- Drop the messaging system tables and related objects

-- Drop triggers first
DROP TRIGGER IF EXISTS update_conversation_last_message_trigger ON public.publisher_messages;
DROP TRIGGER IF EXISTS set_publisher_conversations_timestamp ON public.publisher_conversations;

-- Drop functions
DROP FUNCTION IF EXISTS update_conversation_last_message();

-- Drop tables (in reverse order of creation due to foreign key constraints)
DROP TABLE IF EXISTS public.publisher_messages;
DROP TABLE IF EXISTS public.publisher_conversations;