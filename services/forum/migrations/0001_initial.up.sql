

-- Master index table
CREATE TABLE master_index (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    content_type TEXT NOT NULL CHECK (content_type IN ('user', 'post', 'comment', 'group', 'message', 'event')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_master_type ON master_index(content_type);


CREATE TYPE post_visibility AS ENUM ('public', 'almost_private', 'private');

-- Posts table
CREATE TABLE posts (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    post_title TEXT NOT NULL,
    post_body TEXT NOT NULL,
    post_creator BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    group_id BIGINT REFERENCES groups(id) ON DELETE SET NULL,
    visibility post_visibility NOT NULL DEFAULT 'public',
    image_id BIGINT REFERENCES images(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Comments table
CREATE TABLE comments (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    comment_creator_id BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    parent_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    comment_body TEXT NOT NULL,
    group_id BIGINT REFERENCES groups(id) ON DELETE SET NULL,
    image_id BIGINT REFERENCES images(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- Events table
CREATE TABLE events (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    event_title TEXT NOT NULL,
    event_body TEXT NOT NULL,
    event_creator BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    group_id BIGINT REFERENCES groups(id) ON DELETE SET NULL,
    event_date DATE NOT NULL,
    still_valid BOOLEAN NOT NULL DEFAULT TRUE,
    image_id BIGINT REFERENCES images(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Event response table
CREATE TABLE event_response (
     id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
     event_id BIGINT REFERENCES events(id) ON DELETE CASCADE,
     user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
     going BOOLEAN,
     created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
     CONSTRAINT ux_event_user UNIQUE (event_id, user_id)
);

-- Images table
CREATE TABLE images (
     id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
     file_name TEXT,
     entity_id BIGINT NOT NULL REFERENCES master_index(id) ON DELETE CASCADE,
     created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
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
