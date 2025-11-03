package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

func (db *Database) CreateComment(ctx context.Context, commentBody string, parentID int64, commentCreatorId int64) (id int64, createdAt time.Time, totalComments int64, err error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		log.Println(err)
		return 0, time.Time{}, 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Step 1: Insert into master_index
	if err := tx.QueryRow(ctx,
		`INSERT INTO master_index (content_type) VALUES ('comment') RETURNING id`,
	).Scan(&id); err != nil {
		log.Println(err)
		return 0, time.Time{}, 0, err
	}

	// Step 2: Insert into comments
	if err := tx.QueryRow(ctx,
		`INSERT INTO comments (id, comment_creator_id, parent_id, comment_body)
         VALUES ($1, $2, $3, $4) RETURNING created_at`,
		id, commentCreatorId, parentID, commentBody,
	).Scan(&createdAt); err != nil {
		log.Println("insert to comments", err)
		return 0, time.Time{}, 0, err
	}

	// Step 3: Count all comments for this parent
	if err := tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM comments WHERE parent_id = $1`,
		parentID,
	).Scan(&totalComments); err != nil {
		log.Println("count comments", err)
		return 0, time.Time{}, 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println(err)
		return 0, time.Time{}, 0, err
	}

	return id, createdAt, totalComments, nil
}

// get userId lastcommentId range and fix as GetpostbyId and pagination and return has more
// @Stam: You need to join user as c.Creator DONE!! and join reaction of current user as CurrentUserReactions
func (db *Database) GetCommentsByPostID(ctx context.Context, userID, parentId, lastCommentId int64, pageSize int) (comments []models.Comment, hasMore bool, err error) {
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}
	limit := pageSize + 1

	// const q = `
	//     SELECT
	//         c.id,
	//         c.parent_id,
	//         u.id,
	// 		u.username,
	// 		u.avatar,
	// 		u.status,
	//         c.comment_body,
	//         c.created_at
	//     FROM comments c
	// 	JOIN users u ON u.id = c.comment_creator_id
	//     WHERE c.parent_id = $1
	//       AND c.id > $2
	//     ORDER BY c.id ASC
	//     LIMIT $3
	// `

	const q = `
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
		SELECT
            c.id,
            c.parent_id,
            u.id, 
			u.username, 
			u.avatar,
            c.comment_body,
            c.created_at,
			COALESCE(ura.user_reactions, '{}') AS user_reactions,
			COALESCE(rc.reactions, '{}') AS reactions
        FROM comments c
		JOIN users u ON u.id = c.comment_creator_id
        LEFT JOIN user_reactions_agg ura
			ON ura.content_id = c.id
		LEFT JOIN reaction_counts rc
			ON rc.content_id = c.id
		WHERE c.parent_id = $2
          AND c.id > $3
        ORDER BY c.id ASC
        LIMIT $4
	`
	rows, err := db.Pool.Query(ctx, q, userID, parentId, lastCommentId, limit)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	out := make([]models.Comment, 0, pageSize)
	ids := make([]int64, 0, pageSize)
	for rows.Next() {
		var c models.Comment
		var created time.Time
		if err := rows.Scan(
			&c.Id,
			&c.ParentId,
			&c.Creator.Id,
			&c.Creator.Username,
			&c.Creator.Avatar,
			// &c.Creator.Status,
			&c.Body,
			&created,
			&c.CurrentUserReactions,
			&c.ReactionCount,
		); err != nil {
			return nil, false, err
		}
		c.Creator.Status = "offline"
		c.CreatedAt = &created
		c.ReactionCount = make(map[string]int)
		c.CurrentUserReactions = []string{}
		out = append(out, c)
		ids = append(ids, c.Id)
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}
	if len(out) == 0 {
		return out, false, nil
	}
	{
		countsRows, err := db.Pool.Query(ctx, `
            SELECT content_id, reaction_type, COUNT(*)::int
            FROM reactions
            WHERE content_id = ANY($1)
            GROUP BY content_id, reaction_type
        `, ids)
		if err != nil {
			return nil, false, err
		}
		defer countsRows.Close()
		idx := make(map[int64]int, len(out))
		for i := range out {
			idx[out[i].Id] = i
		}

		for countsRows.Next() {
			var cid int64
			var typ string
			var cnt int
			if err := countsRows.Scan(&cid, &typ, &cnt); err != nil {
				return nil, false, err
			}
			if i, ok := idx[cid]; ok {
				out[i].ReactionCount[typ] = cnt
			}
		}
		if err := countsRows.Err(); err != nil {
			return nil, false, err
		}
	}
	{
		urRows, err := db.Pool.Query(ctx, `
            SELECT content_id, reaction_type
            FROM reactions
            WHERE user_id = $2
              AND content_id = ANY($1)
            ORDER BY content_id
        `, ids, userID)
		if err != nil {
			return nil, false, err
		}
		defer urRows.Close()

		idx := make(map[int64]int, len(out))
		for i := range out {
			idx[out[i].Id] = i
		}
		for urRows.Next() {
			var cid int64
			var typ string
			if err := urRows.Scan(&cid, &typ); err != nil {
				return nil, false, err
			}
			if i, ok := idx[cid]; ok {
				out[i].CurrentUserReactions = append(out[i].CurrentUserReactions, typ)
			}
		}
		if err := urRows.Err(); err != nil {
			return nil, false, err
		}
	}
	if len(out) > pageSize {
		hasMore = true
		out = out[:pageSize]
	}

	return out, hasMore, nil
}

func (db *Database) EditComment(ctx context.Context, commentid int, commentBody string) (int, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	cmdTag, err := tx.Exec(ctx, `
					UPDATE comments
					SET comment_body = $1, updated_at = NOW()
					WHERE id = $2`, commentBody, commentid)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	if cmdTag.RowsAffected() == 0 {
		return 0, fmt.Errorf("comment %d not found", commentid)
	}
	if err := tx.Commit(ctx); err != nil {
		log.Println(err)
		return 0, err
	}
	return commentid, nil
}

func (db *Database) GetCommentById(ctx context.Context, commentId int64, senderId int64) (*models.Comment, error) {

	const q = `
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
		SELECT
            c.id,
            c.parent_id,
            u.id, 
			u.username, 
			u.avatar,
            c.comment_body,
            c.created_at,
			COALESCE(ura.user_reactions, '{}') AS user_reactions,
			COALESCE(rc.reactions, '{}') AS reactions
        FROM comments c
		JOIN users u ON u.id = c.comment_creator_id
        LEFT JOIN user_reactions_agg ura
			ON ura.content_id = c.id
		LEFT JOIN reaction_counts rc
			ON rc.content_id = c.id
          WHERE c.id = $2
	`

	var c models.Comment

	if err := db.Pool.QueryRow(ctx, q, senderId, commentId).Scan(
		&c.Id,
		&c.ParentId,
		&c.Creator.Id,
		&c.Creator.Username,
		&c.Creator.Avatar,
		&c.Body,
		&c.CreatedAt,
		&c.CurrentUserReactions,
		&c.ReactionCount,
	); err != nil {
		return nil, err
	}

	if c.CurrentUserReactions == nil {
		c.CurrentUserReactions = []string{}
	}

	return &c, nil
}
