-- Enable citext extension for case-insensitive text
CREATE EXTENSION IF NOT EXISTS citext;

CREATE COLLATION IF NOT EXISTS case_insensitive_ai (
  provider = icu,
  locale = 'und-u-ks-level1',
  deterministic = false
);

CREATE COLLATION IF NOT EXISTS case_insensitive_ai (
  provider = icu,
  locale = 'und-u-ks-level1',
  deterministic = false
);

-- Master index table
CREATE TABLE master_index (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    content_type TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_master_type ON master_index(content_type);

-- Users table
CREATE TABLE users (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    username CITEXT COLLATE case_insensitive_ai UNIQUE,
    email CITEXT COLLATE case_insensitive_ai UNIQUE,
    gender VARCHAR(255),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    -- active INTEGER NOT NULL DEFAULT 1,
    age VARCHAR(255),
    avatar VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    CONSTRAINT email_or_username_required CHECK (
        (username IS NOT NULL AND username <> '')
        OR (email IS NOT NULL AND email <> '')
    )
);

-- Auth users table
CREATE TABLE auth_user (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    identifier CITEXT COLLATE case_insensitive_ai UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    failed_attempts INTEGER NOT NULL DEFAULT 0,
    last_login_at TIMESTAMPTZ,
    CONSTRAINT ux_auth_user_user UNIQUE (user_id),
    CONSTRAINT ux_auth_user_identifier UNIQUE (Identifier),
    CHECK (Identifier = btrim(Identifier))
);


-- Categories table
CREATE TABLE categories (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- Conversations table
CREATE TABLE conversations (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    dm BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

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

-- Conversation members table
CREATE TABLE conversation_member (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_read_message_id BIGINT REFERENCES messages(id) ON DELETE SET NULL,
    conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    CONSTRAINT conversation_member_unique UNIQUE (conversation_id, user_id)
);
CREATE INDEX idx_conversation_member_conversation ON conversation_member(conversation_id);
CREATE INDEX idx_conversation_member_user ON conversation_member(user_id);

-- Schema migrations table
CREATE TABLE IF NOT EXISTS schema_migrations (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    version TEXT NOT NULL UNIQUE,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Posts table
CREATE TABLE posts (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    post_title TEXT NOT NULL,
    post_body TEXT NOT NULL,
    post_creator BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Comments table
CREATE TABLE comments (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    comment_creator_id BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    parent_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    comment_body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Post categories table (many-to-many)
CREATE TABLE post_categories (
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, category_id)
);
CREATE INDEX idx_post_categories_category ON post_categories(category_id);

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
