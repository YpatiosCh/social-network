-- Drop triggers
DROP TRIGGER IF EXISTS trg_update_group_conversations_modtime ON group_conversations;
DROP TRIGGER IF EXISTS trg_update_group_messages_modtime ON group_messages;
DROP TRIGGER IF EXISTS trg_update_private_conversations_modtime ON private_conversations;
DROP TRIGGER IF EXISTS trg_update_private_messages_modtime ON private_messages;

-- Drop trigger functions
DROP FUNCTION IF EXISTS update_timestamp();

-- Drop tables
DROP TABLE IF EXISTS group_messages;
DROP TABLE IF EXISTS private_messages;
DROP TABLE IF EXISTS group_conversations;
DROP TABLE IF EXISTS private_conversations;
