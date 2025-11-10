
-- Conversations table 
--TODO check if better as two separate tables
CREATE TABLE conversations (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    dm BOOLEAN NOT NULL,
    group_id BIGINT REFERENCES groups(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT dm_group_constraint
        CHECK (
            (dm = TRUE AND group_id IS NULL) OR
            (dm = FALSE AND group_id IS NOT NULL)
        )
);

-- Conversation members table
CREATE TABLE conversation_members (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_read_message_id BIGINT REFERENCES messages(id) ON DELETE SET NULL,
    conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    CONSTRAINT conversation_member_unique UNIQUE (conversation_id, user_id)
);
CREATE INDEX idx_conversation_member_conversation ON conversation_member(conversation_id);
CREATE INDEX idx_conversation_member_user ON conversation_member(user_id);


-- Messages table
CREATE TABLE messages (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    message_text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    delivered BOOLEAN NOT NULL DEFAULT TRUE,
    edited_at TIMESTAMPTZ
);


-- Reactions table
CREATE TABLE reactions (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    content_id BIGINT NOT NULL REFERENCES master_index(id) ON DELETE CASCADE,
    reaction_type TEXT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_reaction_per_content UNIQUE (user_id, content_id, reaction_type)
);

-- Reaction details table
CREATE TABLE reaction_details (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
