package database

import (
	"database/sql"
	"errors"
	"fmt"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

func (db *Database) MarkSeen(req models.MarkSeenDbRequest) (int64, error) {
	tx, err := db.Pool.Begin(req.Ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(req.Ctx)
	var otherMemberId int64
	err = tx.QueryRow(req.Ctx, `
		WITH updated AS (
		UPDATE conversation_member cm
		SET last_read_message_id = $1
		WHERE cm.user_id = $2
			AND cm.conversation_id = $3
			AND (cm.last_read_message_id IS NULL OR cm.last_read_message_id < $1)
		RETURNING cm.conversation_id
		)
		SELECT r.user_id
		FROM conversation_member r
		JOIN updated u ON r.conversation_id = u.conversation_id
		WHERE r.user_id <> $2;
	`, req.MarkSeenRequest.DmId, req.UserId, req.MarkSeenRequest.ConversationId).Scan(&otherMemberId)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("user is not a member or new last_read_message_id not greater")
	} else if err != nil {
		return 0, fmt.Errorf("failed to update last_read_message_id: %w", err)
	}
	if err := tx.Commit(req.Ctx); err != nil {
		return 0, err
	}
	fmt.Println(otherMemberId)
	return otherMemberId, nil
}
