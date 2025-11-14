------------------------------------------
-- Stub tables for external references
------------------------------------------
CREATE TABLE ext_users (
    id BIGINT PRIMARY KEY,
    username TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ext_groups (
    id BIGINT PRIMARY KEY,
    title TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

------------------------------------------
-- Conversations
------------------------------------------
CREATE TABLE conversations (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    dm BOOLEAN NOT NULL,
    group_id BIGINT NOT NULL REFERENCES ext_groups(id), 
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

------------------------------------------
-- Messages
------------------------------------------
CREATE TABLE messages (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id BIGINT NOT NULL REFERENCES ext_users(id) ON DELETE NO ACTION,
    message_text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_messages_conversation ON messages(conversation_id);
CREATE INDEX idx_messages_sender ON messages(sender_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);

------------------------------------------
-- Conversation Members
------------------------------------------
CREATE TABLE conversation_members (
    conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES ext_users(id) ON DELETE CASCADE,
    last_read_message_id BIGINT REFERENCES messages(id) ON DELETE SET NULL,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (conversation_id, user_id)
);

-- Indexes for fast lookups
CREATE INDEX idx_conversation_members_user ON conversation_members(user_id);
CREATE INDEX idx_conversation_members_last_read ON conversation_members(last_read_message_id);


------------------------------------------
-- Reactions
------------------------------------------
CREATE TABLE reactions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    message_id BIGINT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    reaction_type TEXT NOT NULL,
    user_id BIGINT NOT NULL, -- in users service
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_reaction_per_content UNIQUE (user_id, message_id, reaction_type)
);

CREATE INDEX idx_reactions_user ON reactions(user_id);