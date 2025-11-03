package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

func (db *Database) GetUsersFeed(req models.UsersFeedDbRequest) (models.UsersFeedDbResponse, error) {
	ctx := req.Ctx
	me := req.UserId
	n := req.Range
	if n <= 0 {
		n = 10
	}
	if n > 50 {
		n = 50
	}

	var feed []models.Feed
	need := n

	var inDM bool
	var cursorTs *time.Time // for DM segment
	var cursorName *string  // for alpha segment
	lastID := req.LastId

	// decide which segment we're resuming from
	if lastID > 0 {
		ok, ts, err := db.DMWithLastMsgAt(ctx, me, lastID)
		if err != nil {
			return models.UsersFeedDbResponse{}, err
		}
		inDM = ok
		if ok {
			cursorTs = ts // may be nil if no messages yet
		} else {
			if err := db.Pool.QueryRow(ctx, `SELECT LOWER(username) FROM users WHERE id = $1`, lastID).Scan(&cursorName); err != nil && !errors.Is(err, pgx.ErrNoRows) {
				return models.UsersFeedDbResponse{}, err
			}
		}
	}

	// 1) DM segment (on first page or if cursor is in DM)
	if lastID == 0 || inDM {
		rows, err := db.PageDMUsers(ctx, me, cursorTs, lastID, need+1) // +1 for hasMore within DM
		if err != nil {
			return models.UsersFeedDbResponse{}, err
		}

		if len(rows) > need {
			return models.UsersFeedDbResponse{
				UsersFeedResponse: models.UsersFeedResponse{
					Feed:     rows[:need],
					HaveMore: true, // still more in DM segment
				},
			}, nil
		}
		feed = append(feed, rows...)
		need -= len(rows)

		// after DM completes (or if there were none), alpha starts from the beginning
		cursorName, lastID = nil, 0
	}

	// 2) Alpha segment (top-up)
	if need > 0 {
		rows, err := db.PageAlphaUsers(ctx, me, cursorName, lastID, need+1)
		if err != nil {
			return models.UsersFeedDbResponse{}, err
		}

		haveMore := false
		if len(rows) > need {
			haveMore = true
			rows = rows[:need]
		}
		feed = append(feed, rows...)

		return models.UsersFeedDbResponse{
			UsersFeedResponse: models.UsersFeedResponse{
				Feed:     feed,
				HaveMore: haveMore,
			},
		}, nil
	}

	return models.UsersFeedDbResponse{
		UsersFeedResponse: models.UsersFeedResponse{
			Feed:     feed,
			HaveMore: false,
		},
	}, nil
}

func (db *Database) PageDMUsers(ctx context.Context, me int64, afterTs *time.Time, afterUserID int64, limit int) ([]models.Feed, error) {
	rows, err := db.Pool.Query(ctx, `
		WITH my_dms AS (
		SELECT c.id AS conv_id, other.user_id AS other_id
		FROM conversations c
		JOIN conversation_member me_cm ON me_cm.conversation_id = c.id AND me_cm.user_id = $1
		JOIN conversation_member other ON other.conversation_id = c.id AND other.user_id <> $1
		WHERE c.dm = TRUE
		),
		enriched AS (
		SELECT
			u.id AS user_id,
			u.username,
			COALESCE(u.avatar, '') AS avatar,
			me_cm.last_read_message_id,
			my_dms.conv_id,
			lm.created_at AS last_msg_at,
			COALESCE((
			SELECT COUNT(*) FROM messages m2
			WHERE m2.conversation_id = my_dms.conv_id
				AND m2.id > COALESCE(me_cm.last_read_message_id, 0)
				AND m2.sender <> $1
			), 0) AS unread_count
		FROM my_dms
		JOIN users u ON u.id = my_dms.other_id
		JOIN conversation_member me_cm ON me_cm.conversation_id = my_dms.conv_id AND me_cm.user_id = $1
		LEFT JOIN LATERAL (
			SELECT m.created_at
			FROM messages m
			WHERE m.conversation_id = my_dms.conv_id
			ORDER BY m.id DESC
			LIMIT 1
		) lm ON TRUE
		)
		SELECT
		user_id, username, avatar,
		conv_id,
		unread_count,
		last_read_message_id,
		last_msg_at
		FROM enriched
		WHERE ($2::timestamptz IS NULL OR (last_msg_at, user_id) < ($2, $3))
		ORDER BY last_msg_at DESC NULLS LAST, user_id DESC
		LIMIT $4
		`, me, afterTs, afterUserID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.Feed, 0, limit)
	for rows.Next() {
		var (
			uid           int64
			uname, avatar string
			convID        int64
			unread        int
			lastRead      sql.NullInt64
			_             *time.Time // last_msg_at, not needed in response body
		)
		var lastMsgAt sql.NullTime
		if err := rows.Scan(&uid, &uname, &avatar, &convID, &unread, &lastRead, &lastMsgAt); err != nil {
			return nil, err
		}
		var lrm *uint64
		if lastRead.Valid && lastRead.Int64 > 0 {
			v := uint64(lastRead.Int64)
			lrm = &v
		}
		out = append(out, models.Feed{
			User: models.UserOnFeed{Id: uid, Username: uname, Avatar: avatar, Status: "offline" /* hydrate presence later */},
			ConversationDetails: models.ConversationDetails{
				ConversationId:    uint64(convID),
				UnreadCount:       unread,
				LastReadMessageId: lrm,
			},
		})
	}
	return out, rows.Err()
}

func (db *Database) PageAlphaUsers(ctx context.Context, me int64, afterName *string, afterUserID int64, limit int) ([]models.Feed, error) {

	rows, err := db.Pool.Query(ctx, `
        WITH my_dms AS (
          SELECT cm2.user_id
          FROM conversations c
          JOIN conversation_member cm1 ON cm1.conversation_id = c.id AND cm1.user_id = $1
          JOIN conversation_member cm2 ON cm2.conversation_id = c.id AND cm2.user_id <> $1
          WHERE c.dm = TRUE
        ),
        base AS (
          SELECT u.id AS user_id, u.username, COALESCE(u.avatar, '') AS avatar
          FROM users u
          WHERE u.id <> $1
            AND u.id NOT IN (SELECT user_id FROM my_dms)
            AND ($2::text IS NULL OR (LOWER(u.username), u.id) > (LOWER($2::text), $3))
          ORDER BY LOWER(u.username) ASC, u.id ASC
          LIMIT $4
        )
        SELECT user_id, username, avatar FROM base
    `, me, afterName, afterUserID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.Feed, 0, limit)
	for rows.Next() {
		var (
			uid    int64
			uname  string
			avatar string
		)
		if err := rows.Scan(&uid, &uname, &avatar); err != nil {
			return nil, err
		}

		f := models.Feed{
			User: models.UserOnFeed{
				Id:       uid,
				Username: uname,
				Avatar:   avatar,
				Status:   "offline",
			},
			// ConversationDetails left zero / nil because we explicitly exclude users in conversations.
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (db *Database) DMWithLastMsgAt(ctx context.Context, me, other int64) (bool, *time.Time, error) {
	var t *time.Time
	err := db.Pool.QueryRow(ctx, `
		WITH dm AS (
		SELECT c.id
		FROM conversations c
		JOIN conversation_member cm1 ON cm1.conversation_id = c.id AND cm1.user_id = $1
		JOIN conversation_member cm2 ON cm2.conversation_id = c.id AND cm2.user_id = $2
		WHERE c.dm = TRUE
		LIMIT 1
		)
		SELECT m.created_at
		FROM messages m JOIN dm ON m.conversation_id = dm.id
		ORDER BY m.id DESC
		LIMIT 1
		`, me, other).Scan(&t)
	if err == nil {
		return true, t, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil, nil
	}
	return false, nil, err
}

func (db *Database) GetUserFeedById(ctx context.Context, userId int64, otherUserId int64) (models.Feed, error) {
	var feed models.Feed

	row := db.Pool.QueryRow(ctx, `
WITH my_dms AS (
    SELECT c.id AS conv_id, other.user_id AS other_id
    FROM conversations c
    JOIN conversation_member me_cm 
      ON me_cm.conversation_id = c.id 
     AND me_cm.user_id = $1
    JOIN conversation_member other 
      ON other.conversation_id = c.id 
     AND other.user_id <> $1
    WHERE c.dm = TRUE
      AND other.user_id = $2
),
enriched AS (
    SELECT
        u.id AS user_id,
        u.username,
        COALESCE(u.avatar, '') AS avatar,
        me_cm.last_read_message_id,
        my_dms.conv_id,
        COALESCE((
            SELECT COUNT(*) 
            FROM messages m2
            WHERE m2.conversation_id = my_dms.conv_id
              AND m2.id > COALESCE(me_cm.last_read_message_id, 0)
              AND m2.sender <> $1
        ), 0) AS unread_count
    FROM my_dms
    JOIN users u ON u.id = my_dms.other_id
    JOIN conversation_member me_cm 
      ON me_cm.conversation_id = my_dms.conv_id 
     AND me_cm.user_id = $1
)
SELECT
    user_id, username, avatar,
    conv_id,
    unread_count,
    last_read_message_id
FROM enriched
WHERE user_id = $2
LIMIT 1;
	`, userId, otherUserId)

	if err := row.Scan(
		&feed.User.Id,
		&feed.User.Username,
		&feed.User.Avatar,
		&feed.ConversationDetails.ConversationId,
		&feed.ConversationDetails.UnreadCount,
		&feed.ConversationDetails.LastReadMessageId,
	); err != nil {
		return models.Feed{}, err
	}

	return feed, nil
}
