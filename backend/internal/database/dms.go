package database

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

func (db *Database) CreateMessage(ctx context.Context, conversationID, senderID, receiverID int64, text string) (convID, msgID int64, createdAt time.Time, err error) {
	if strings.TrimSpace(text) == "" {
		return 0, 0, time.Time{}, fmt.Errorf("empty message")
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return 0, 0, time.Time{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var userExists bool
	if err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, receiverID).Scan(&userExists); err != nil {
		return 0, 0, time.Time{}, err
	}
	if !userExists {
		return 0, 0, time.Time{}, fmt.Errorf("receiver not found")
	}

	// Verify sender is a member of the conversation (covers both DM & non-DM)
	var isMember bool
	if err = tx.QueryRow(ctx, `
        SELECT EXISTS (
            SELECT 1 FROM conversation_member
            WHERE conversation_id = $1 AND user_id = $2
        )
    `, conversationID, senderID).Scan(&isMember); err != nil {
		return 0, 0, time.Time{}, err
	}
	if !isMember && conversationID != 0 {
		return 0, 0, time.Time{}, fmt.Errorf("sender is not in the conversation")
	}

	// Allocate message id
	if err = tx.QueryRow(ctx, `
        INSERT INTO master_index (content_type) VALUES ('message') RETURNING id
    `).Scan(&msgID); err != nil {
		return 0, 0, time.Time{}, err
	}

	// If no conversationID is provided, resolve/create a DM between sender and receiver.
	var cmID1, cmID2 int64
	if conversationID == 0 {
		if receiverID <= 0 || receiverID == senderID {
			return 0, 0, time.Time{}, fmt.Errorf("invalid receiver")
		}

		err = tx.QueryRow(ctx, `
            SELECT c.id
            FROM conversations c
            JOIN conversation_member cm1 ON cm1.conversation_id = c.id AND cm1.user_id = $1
            JOIN conversation_member cm2 ON cm2.conversation_id = c.id AND cm2.user_id = $2
            WHERE c.dm = TRUE
            LIMIT 1
        `, senderID, receiverID).Scan(&conversationID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				if err = tx.QueryRow(ctx, `
                    INSERT INTO master_index (content_type) VALUES ('conversation') RETURNING id
                `).Scan(&conversationID); err != nil {
					return 0, 0, time.Time{}, err
				}

				if _, err = tx.Exec(ctx, `
                    INSERT INTO conversations (id, dm)
                    VALUES ($1, TRUE)
                `, conversationID); err != nil {
					return 0, 0, time.Time{}, err
				}

				if err = tx.QueryRow(ctx, `INSERT INTO master_index (content_type) VALUES ('conversation_member') RETURNING id`).Scan(&cmID1); err != nil {
					return 0, 0, time.Time{}, err
				}
				if err = tx.QueryRow(ctx, `INSERT INTO master_index (content_type) VALUES ('conversation_member') RETURNING id`).Scan(&cmID2); err != nil {
					return 0, 0, time.Time{}, err
				}

			} else {
				return 0, 0, time.Time{}, err
			}
		}
	}

	if err = tx.QueryRow(ctx, `
		INSERT INTO messages (id, conversation_id, sender, message_text)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`, msgID, conversationID, senderID, text).Scan(&createdAt); err != nil {
		return 0, 0, time.Time{}, err
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO conversation_member (id, user_id, conversation_id, last_read_message_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, conversation_id)
		DO UPDATE SET last_read_message_id = GREATEST(conversation_member.last_read_message_id, EXCLUDED.last_read_message_id);
		`, cmID1, senderID, conversationID, msgID); err != nil {
		return 0, 0, time.Time{}, err
	}

	if _, err = tx.Exec(ctx, `
		INSERT INTO conversation_member (id, user_id, conversation_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, conversation_id)
		DO NOTHING;
		`, cmID2, receiverID, conversationID); err != nil {
		return 0, 0, time.Time{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, 0, time.Time{}, err
	}
	fmt.Println("Create Msg", err)
	return conversationID, msgID, createdAt, nil
}

func (db *Database) GetConversationDms(req models.MessagesDbRequest) (res models.MessagesDbResponse, err error) {
	ctx := req.Ctx

	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{}) // or sql.TxOptions{}
	if err != nil {
		return res, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	// --- BEFORE (including LastId) ---
	beforeRows, err := tx.Query(ctx, `
		WITH user_reactions_agg AS (
			SELECT content_id, array_agg(reaction_type) AS user_reactions
			FROM reactions
			WHERE user_id = $1
			GROUP BY content_id
		),
		reaction_counts AS (
			SELECT content_id, jsonb_object_agg(reaction_type, count) AS reactions
			FROM (
				SELECT content_id, reaction_type, COUNT(*) AS count
				FROM reactions
				GROUP BY content_id, reaction_type
			) t
			GROUP BY content_id
		)
		SELECT m.id, m.conversation_id, m.sender, m.message_text, m.created_at,
			COALESCE(ura.user_reactions, '{}') AS user_reactions,
			COALESCE(rc.reactions, '{}') AS reactions
		FROM messages m
		JOIN conversation_member cm
		ON cm.conversation_id = m.conversation_id
		AND cm.user_id = $1
		LEFT JOIN user_reactions_agg ura
		ON ura.content_id = m.id
		LEFT JOIN reaction_counts rc
		ON rc.content_id = m.id
		WHERE m.conversation_id = $2 AND m.id <= $3
		ORDER BY m.id DESC
		LIMIT $4;
    `, req.UserId, req.ConversationId, req.LastId, req.Before+1)
	if err != nil {
		return res, fmt.Errorf("failed to fetch messages before: %w", err)
	}
	defer beforeRows.Close()

	var before []models.Message
	for beforeRows.Next() {
		var m models.Message
		if err := beforeRows.Scan(&m.Id, &m.ConversationId, &m.SenderId, &m.Body, &m.CreatedAt, &m.CurrentUserReactions, &m.ReactionCount); err != nil {
			return res, fmt.Errorf("failed to scan before message: %w", err)
		}
		before = append(before, m)
	}
	res.HaveMoreBefore = len(before) > int(req.Before)
	if res.HaveMoreBefore {
		before = before[:len(before)-1] // drop extra
	}
	// reverse for chronological order
	for i, j := 0, len(before)-1; i < j; i, j = i+1, j-1 {
		before[i], before[j] = before[j], before[i]
	}

	// --- AFTER ---
	afterRows, err := tx.Query(ctx, `
		WITH user_reactions_agg AS (
			SELECT content_id, array_agg(reaction_type) AS user_reactions
			FROM reactions
			WHERE user_id = $1
			GROUP BY content_id
		),
		reaction_counts AS (
			SELECT content_id, jsonb_object_agg(reaction_type, count) AS reactions
			FROM (
				SELECT content_id, reaction_type, COUNT(*) AS count
				FROM reactions
				GROUP BY content_id, reaction_type
			) t
			GROUP BY content_id
		)
		SELECT m.id, m.conversation_id, m.sender, m.message_text, m.created_at,
			COALESCE(ura.user_reactions, '{}') AS user_reactions,
			COALESCE(rc.reactions, '{}') AS reactions
		FROM messages m
		JOIN conversation_member cm
		ON cm.conversation_id = m.conversation_id
		AND cm.user_id = $1
		LEFT JOIN user_reactions_agg ura
		ON ura.content_id = m.id
		LEFT JOIN reaction_counts rc
		ON rc.content_id = m.id
		WHERE m.conversation_id = $2 AND m.id > $3
		ORDER BY m.id ASC
		LIMIT $4;
    `, req.UserId, req.ConversationId, req.LastId, req.After+1)
	if err != nil {
		return res, fmt.Errorf("failed to fetch messages after: %w", err)
	}
	defer afterRows.Close()

	var after []models.Message
	for afterRows.Next() {
		var m models.Message
		if err := afterRows.Scan(&m.Id, &m.ConversationId, &m.SenderId, &m.Body, &m.CreatedAt, &m.CurrentUserReactions, &m.ReactionCount); err != nil {
			return res, fmt.Errorf("failed to scan after message: %w", err)
		}
		after = append(after, m)
	}
	res.HaveMoreAfter = len(after) > int(req.After)
	if res.HaveMoreAfter {
		after = after[:len(after)-1] // drop extra
	}

	// Combine before+after into final response
	res.Messages = append(before, after...)
	return res, nil
}
