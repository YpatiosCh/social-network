package dbservice

const (
	// MESSAGES
	createMessage = `
	INSERT INTO messages (conversation_id, sender_id, message_text)
	SELECT $1, $2, $3
	FROM conversation_members
	WHERE conversation_id = $1
	AND user_id = $2
	AND deleted_at IS NULL
	RETURNING id, conversation_id, sender_id, message_text, created_at, updated_at, deleted_at
`
	getPrevMessages = `
	SELECT 
		m.id,
		m.conversation_id,
		m.sender_id,
		m.message_text,
		m.created_at,
		m.updated_at,
		m.deleted_at,
		c.first_message_id
	FROM messages m
	JOIN conversations c
		ON c.id = m.conversation_id
	JOIN conversation_members cm
		ON cm.conversation_id = m.conversation_id
	WHERE m.conversation_id = $2
	AND cm.user_id = $3
	AND m.deleted_at IS NULL
	AND m.id < COALESCE(
			$1,
			c.last_message_id
		)
	AND COALESCE(
			$1,
			c.last_message_id
		) > c.first_message_id
	ORDER BY m.id DESC
	LIMIT $4 OFFSET $5;
`
	getNextMessages = `
	SELECT 
		m.id,
		m.conversation_id,
		m.sender_id,
		m.message_text,
		m.created_at,
		m.updated_at,
		m.deleted_at,
		c.last_message_id
	FROM messages m
	JOIN conversations c
		ON c.id = m.conversation_id
	JOIN conversation_members cm
		ON cm.conversation_id = m.conversation_id
	WHERE m.conversation_id = $2
	AND cm.user_id = $3
	AND m.deleted_at IS NULL
	AND m.id > COALESCE(
			$1,
			c.first_message_id
		)
	AND COALESCE(
			$1,
			c.first_message_id
		) < c.last_message_id
	ORDER BY m.id ASC
	LIMIT $4 OFFSET $5;
`
	updateLastReadMessage = `
	UPDATE conversation_members cm
	SET last_read_message_id = $3
	WHERE cm.conversation_id = $1
	AND cm.user_id = $2
	AND cm.deleted_at IS NULL
	RETURNING conversation_id, user_id, last_read_message_id, joined_at, deleted_at
`
)
