package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

func (db *Database) GetPostById(ctx context.Context, postID, userID int64) (*models.Post, error) {
	const postQ = `
		SELECT
			p.id,
			p.post_title,
			p.post_body,
			p.created_at,
			u.id       AS creator_id,
			u.username AS creator_username,
			u.avatar   AS creator_avatar
		FROM posts p
		JOIN users u ON u.id = p.post_creator
		WHERE p.id = $1
	`

	var p models.Post
	var createdAt time.Time

	// 1) Basic post + creator
	if err := db.Pool.QueryRow(ctx, postQ, postID).Scan(
		&p.Id,
		&p.Title,
		&p.Body,
		&createdAt,
		&p.Creator.Id,
		&p.Creator.Username,
		&p.Creator.Avatar,
	); err != nil {
		return nil, err
	}
	p.CreatedAt = &createdAt

	const commentsQ = `
		SELECT COUNT(*)::INT,
       COALESCE((
           SELECT id
           FROM comments
           WHERE parent_id = $1
           ORDER BY created_at ASC, id ASC
           LIMIT 1
       		), 0)
		FROM comments
		WHERE parent_id = $1;                               
	`
	var firstCommentID int64
	if err := db.Pool.QueryRow(ctx, commentsQ, postID).Scan(&p.CommentCount, &firstCommentID); err != nil {
		return nil, err
	}
	p.FirstCommentId = firstCommentID

	const reactCountsQ = `
		SELECT reaction_type, COUNT(*)::INT
		FROM reactions
		WHERE content_id = $1
		GROUP BY reaction_type
	`
	p.ReactionCount = make(map[string]int)
	if rows, err := db.Pool.Query(ctx, reactCountsQ, postID); err == nil {
		defer rows.Close()
		for rows.Next() {
			var t string
			var c int
			if scanErr := rows.Scan(&t, &c); scanErr != nil {
				return nil, scanErr
			}
			// p.ReactionCount = p.ReactionCount // (no-op, just clarity)
			p.ReactionCount[t] = c
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	const myReactionsQ = `
		SELECT reaction_type
		FROM reactions
		WHERE content_id = $1 AND user_id = $2
		ORDER BY reaction_type
	`
	p.CurrentUserReactions = nil
	if rows, err := db.Pool.Query(ctx, myReactionsQ, postID, userID); err == nil {
		defer rows.Close()
		for rows.Next() {
			var t string
			if scanErr := rows.Scan(&t); scanErr != nil {
				return nil, scanErr
			}
			p.CurrentUserReactions = append(p.CurrentUserReactions, t)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	const catsQ = `
		SELECT c.id, c.name
		FROM post_categories pc
		JOIN categories c ON c.id = pc.category_id
		WHERE pc.post_id = $1
		ORDER BY c.name
	`
	p.Categories = nil

	if rows, err := db.Pool.Query(ctx, catsQ, postID); err == nil {
		defer rows.Close()
		for rows.Next() {
			var cr models.Category // adjust to your actual fields
			if scanErr := rows.Scan(&cr.Id, &cr.Name); scanErr != nil {
				return nil, scanErr
			}
			p.Categories = append(p.Categories, cr)
		}
		// p.Categories = []models.CategoriesResponse{
		// 	{Categories: cats},
		// }
		if err := rows.Err(); err != nil {
			return nil, err

		}
	} else {
		return nil, err
	}

	p.Truncated = false

	if p.CurrentUserReactions == nil {
		p.CurrentUserReactions = []string{}
	}

	return &p, nil
}

// CreatePost inserts a new post into the database with the given title, content, and author
func (db *Database) CreatePost(ctx context.Context, postTitle, postBody string, postCreator int64, categories []int) (id int64, createdAt time.Time, err error) {
	if strings.TrimSpace(postTitle) == "" || strings.TrimSpace(postBody) == "" {
		return 0, time.Time{}, fmt.Errorf("title and body are required")
	}
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		log.Println(err)
		return 0, time.Time{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err = tx.QueryRow(ctx,
		`INSERT INTO master_index (content_type) VALUES ('post') RETURNING id`,
	).Scan(&id); err != nil {
		log.Println(err)
		return 0, time.Time{}, err
	}
	if err = tx.QueryRow(ctx,
		`INSERT INTO posts (id, post_title, post_body, post_creator) VALUES($1, $2, $3, $4 ) RETURNING created_at `,
		id, postTitle, postBody, postCreator).Scan(&createdAt); err != nil {
		log.Println(err)
		return 0, time.Time{}, err
	}

	if len(categories) > 0 {
		for _, c := range categories {
			if _, err := tx.Exec(ctx, `
                INSERT INTO post_categories (post_id, category_id)
                VALUES ($1, $2) ON CONFLICT DO NOTHING`, id, c); err != nil {
				return 0, time.Time{}, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println(err)
		return 0, time.Time{}, err
	}
	return id, createdAt, nil
}

func (db *Database) GetPosts(req models.PostsFeedDbRequest) (res models.PostsFeedDbResponse, err error) {
	ctx := req.Ctx

	rows, err := db.Pool.Query(ctx, `
	WITH filtered_posts AS (
		SELECT p.*
		FROM posts p
		WHERE ($1::bigint[] IS NULL OR EXISTS (
			SELECT 1
			FROM post_categories pc
			WHERE pc.post_id = p.id
			AND pc.category_id = ANY($1)
		))
		AND ($2::bigint IS NULL OR p.id < $2)
		ORDER BY p.id DESC
		LIMIT $3
	),
	post_reactions AS (
		SELECT r.content_id AS post_id,
			jsonb_object_agg(r.reaction_type, cnt) AS reactions
		FROM (
			SELECT content_id, reaction_type, COUNT(*) AS cnt
			FROM reactions
			WHERE content_id IN (SELECT id FROM filtered_posts)
			GROUP BY content_id, reaction_type
		) r
		GROUP BY r.content_id
	),
	user_reactions AS (
		SELECT content_id AS post_id,
			array_agg(reaction_type) AS current_user_reactions
		FROM reactions
		WHERE content_id IN (SELECT id FROM filtered_posts)
		AND user_id = $4
		GROUP BY content_id
	),
	comment_counts AS (
		SELECT parent_id AS post_id,
			COUNT(*) AS comment_count,
			MIN(id) AS first_comment_id
		FROM comments
		WHERE parent_id IN (SELECT id FROM filtered_posts)
		GROUP BY parent_id
	),
	post_categories_agg AS (
		SELECT pc.post_id,
			jsonb_agg(jsonb_build_object('id', c.id, 'name', c.name)) AS categories
		FROM post_categories pc
		JOIN categories c ON c.id = pc.category_id
		WHERE pc.post_id IN (SELECT id FROM filtered_posts)
		AND (
				array_length($1::bigint[], 1) IS NULL -- no filter provided (NULL)
			OR array_length($1::bigint[], 1) = 0    -- empty array (no category filter)
			OR pc.category_id = ANY($1)             -- filter by given categories
		)
		GROUP BY pc.post_id
	)
	SELECT 
		p.id,
		p.created_at,
		p.post_title,
		p.post_body,
		u.id, u.username, u.avatar,
		COALESCE(pr.reactions, '{}'::jsonb),
		COALESCE(ur.current_user_reactions, '{}'),
		COALESCE(cc.comment_count, 0),
		COALESCE(pca.categories, '[]'::jsonb),
		COALESCE(cc.first_comment_id, 0)
	FROM filtered_posts p
	JOIN users u ON u.id = p.post_creator
	LEFT JOIN post_reactions pr ON pr.post_id = p.id
	LEFT JOIN user_reactions ur ON ur.post_id = p.id
	LEFT JOIN comment_counts cc ON cc.post_id = p.id
	LEFT JOIN post_categories_agg pca ON pca.post_id = p.id
	ORDER BY p.id DESC
`, req.Categories, nullableBigInt(req.LastSeenId), req.Range+1, req.UserId)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var (
			post                 models.Post
			creator              models.PostCreator
			reactionCountJSON    []byte
			currentUserReactions []string
			categoriesJSON       []byte
		)

		err = rows.Scan(
			&post.Id,
			&post.CreatedAt,
			&post.Title,
			&post.Body,
			&creator.Id,
			&creator.Username,
			&creator.Avatar,
			&reactionCountJSON,
			&currentUserReactions,
			&post.CommentCount,
			&categoriesJSON,
			&post.FirstCommentId,
		)
		if err != nil {
			return res, err
		}

		// decode JSON -> Go structs
		if err := json.Unmarshal(reactionCountJSON, &post.ReactionCount); err != nil {
			post.ReactionCount = map[string]int{}
		}
		if err := json.Unmarshal(categoriesJSON, &post.Categories); err != nil {
			post.Categories = []models.Category{}
		}
		post.Creator = creator
		post.CurrentUserReactions = currentUserReactions

		posts = append(posts, post)
	}

	res.Posts = posts
	res.HaveMore = len(posts) > req.Range
	return res, rows.Err()
}

// Helper for nullable last_id
func nullableBigInt(id int64) *int64 {
	if id == 0 {
		return nil
	}
	return &id
}
