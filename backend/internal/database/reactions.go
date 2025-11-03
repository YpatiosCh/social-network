package database

import (
	"context"
	"fmt"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

// CreateReaction inserts a new reaction into the database with the given userId, and type
func (db *Database) CreateReaction(ctx context.Context, userID, contentID int64, reactionType string) (reactionID, total int64, err error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Allocate global id
	if err = tx.QueryRow(ctx, `INSERT INTO master_index (content_type) VALUES ('reaction') RETURNING id`).Scan(&reactionID); err != nil {
		return 0, 0, err
	}

	// Insert reaction
	if _, err = tx.Exec(ctx, `
		INSERT INTO reactions (id, user_id, content_id, reaction_type)
		VALUES ($1, $2, $3, $4)
	`, reactionID, userID, contentID, reactionType); err != nil {
		return 0, 0, err
	}

	// Total **unique users** who reacted to this content
	if err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM reactions WHERE content_id = $1`, contentID).Scan(&total); err != nil {
		return 0, 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, 0, err
	}
	return reactionID, total, nil
}

func (db *Database) DeleteReaction(ctx context.Context, userID, contentID int64, reactionType string) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	cmd, err := tx.Exec(ctx, `
        DELETE FROM reactions
        WHERE user_id = $1 AND content_id = $2 AND reaction_type = $3
    `, userID, contentID, reactionType)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("no reaction found to delete")
	}

	return tx.Commit(ctx)
}

func (db *Database) ListUsersReacted(ctx context.Context, contentID int64) ([]models.UserReacted, error) {
	rows, err := db.Pool.Query(ctx, `
        SELECT u.id, u.username
        FROM reactions r
        JOIN users u ON u.id = r.user_id
        WHERE r.content_id = $1
        GROUP BY u.id, u.username
        ORDER BY MAX(r.created_at) DESC
    `, contentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserReacted
	for rows.Next() {
		var u models.UserReacted
		if err := rows.Scan(&u.Id, &u.Username); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (db *Database) ReactionCountsByType(ctx context.Context, contentID int64) (map[string]int, int64, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT reaction_type, COUNT(*)::int
		FROM reactions
		WHERE content_id = $1
		GROUP BY reaction_type
	`, contentID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	var total int64
	for rows.Next() {
		var t string
		var c int
		if err := rows.Scan(&t, &c); err != nil {
			return nil, 0, err
		}
		counts[t] = c
		total += int64(c)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return counts, total, nil
}

func (db *Database) CurrentUserReactions(ctx context.Context, userID, contentID int64) ([]string, error) {

	const myReactionsQ = `
		SELECT reaction_type
		FROM reactions
		WHERE content_id = $1 AND user_id = $2
		ORDER BY reaction_type
	`
	var reactions []string
	if rows, err := db.Pool.Query(ctx, myReactionsQ, contentID, userID); err == nil {
		defer rows.Close()
		for rows.Next() {
			var r string
			if scanErr := rows.Scan(&r); scanErr != nil {
				return nil, scanErr
			}
			reactions = append(reactions, r)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	return reactions, nil
}
