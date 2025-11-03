package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

// CreateUser is only for admin/imports not for public registrations
func (db *Database) CreateUser(req models.RegisterDbRequest) (res models.RegisterDbResponse, err error) {
	res = models.RegisterDbResponse{}
	tx, err := db.Pool.Begin(req.Ctx)
	if err != nil {
		return res, err
	}
	defer tx.Rollback(req.Ctx)

	var id int64
	if err := tx.QueryRow(req.Ctx,
		`INSERT INTO master_index (content_type) VALUES ('user') RETURNING id`,
	).Scan(&id); err != nil {
		fmt.Println("error inserting into master_index:", err)
		return res, err
	}
	_, err = tx.Exec(req.Ctx, `
        INSERT INTO users (id, username, gender, first_name, last_name, email, age, avatar)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, id, req.Username, req.Gender, req.FirstName, req.LastName, req.Email, req.Age, req.Avatar)
	if err != nil {
		fmt.Println("error inserting into users:", err)
		return res, err
	}

	if err := tx.Commit(req.Ctx); err != nil {
		return res, err
	}
	return res, nil
}

// // CreateConversation is not needed as is happening in the  create Message
// func (db *Database) CreateConversation(ctx context.Context, members []int, dm bool) (convID int, joinedAt time.Time, err error) {
// 	if len(members) == 0 {
// 		return 0, time.Time{}, fmt.Errorf("members required")
// 	}

// 	tx, err := db.Pool.Begin(ctx)
// 	if err != nil {
// 		log.Println(err)
// 		return 0, time.Time{}, err
// 	}
// 	defer func() { _ = tx.Rollback(ctx) }()

// 	var convId int
// 	if err := tx.QueryRow(ctx,
// 		`INSERT INTO master_index (content_type) VALUES ('conversation') RETURNING id`,
// 	).Scan(&convId); err != nil {
// 		log.Println(err)
// 		return 0, time.Time{}, err
// 	}

// 	if _, err := tx.Exec(ctx,
// 		`INSERT INTO conversations (id, dm) VALUES ($1, $2)`,
// 		convId, dm,
// 	); err != nil {
// 		log.Println(err)
// 		return 0, time.Time{}, err
// 	}

// 	for _, userId := range members {
// 		var ja time.Time
// 		err := tx.QueryRow(ctx, `
//         INSERT INTO conversation_member (conversation_id, user_id, last_message_id, joined_at)
//         VALUES ($1, $2, 0)
//         RETURNING joined_at
//     `, convId, userId).Scan(&joinedAt)
// 		if err != nil {
// 			log.Println(err)
// 			return 0, time.Time{}, err
// 		}
// 		if joinedAt.IsZero() {
// 			joinedAt = ja
// 		}
// 	}
// 	// Todo lastSendMessage = 0 and joined_at ie created_at
// 	if err := tx.Commit(ctx); err != nil {
// 		log.Println(err)
// 		return 0, time.Time{}, err
// 	}
// 	return convId, joinedAt, nil
// }

func (db *Database) UpdateUser(ctx context.Context, userID int64, username, gender, firstName, lastName, email, age, avatar string) (int64, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	cmdTag, err := tx.Exec(ctx, `
		UPDATE users 
		SET username = $1, gender = $2, first_name = $3, last_name = $4, email = $5, age = $ 6, avatar = $7, updated_at = NOW()
		WHERE id = $8`, username, gender, firstName, lastName, email, age, avatar, userID)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	if cmdTag.RowsAffected() == 0 {
		return 0, fmt.Errorf("user %d not found", userID)
	}
	return userID, nil
}

func (db *Database) EditPost(ctx context.Context, postID int64, postTitle, postBody string) (int64, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	cmdTag, err := tx.Exec(ctx, ` 
			UPDATE posts
			SET post_title = $1, post_body = $2, updated_at = NOW()
			WHERE id = $3`, postTitle, postBody, postID)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	if cmdTag.RowsAffected() == 0 {
		return 0, fmt.Errorf("post %d not found", postID)
	}
	if err := tx.Commit(ctx); err != nil {
		log.Println(err)
		return 0, err
	}
	return postID, nil
}

// TODO getPostByRange will receive lastID(after lastID and over within range) range category and retrun hasMore for pagination
// TODO get userReactions func
//TODO userFeed(ctx, lastId, range, withConversation bool)
// IF withCon == true. Return users that have con with me shorted by last interaction. The rest to be returned alphabetical. Auto tha exei mesa ena feed me user name id and avatar. an einai withConv ta conv id, unread count & lastread message id

func parseCategoryIDsCSV(parts []string) ([]int32, bool, error) {
	if len(parts) == 0 {
		return nil, true, nil // no filter
	}

	out := make([]int32, 0, len(parts))
	for _, s := range parts {
		s = strings.TrimSpace(s)
		if s == "" { // skip empties
			continue
		}
		n, err := strconv.Atoi(s)
		if err != nil || n <= 0 {
			return nil, false, fmt.Errorf("bad cat id value: %q", s)
		}
		out = append(out, int32(n))
	}

	if len(out) == 0 {
		return nil, true, nil // effectively no filter
	}
	return out, false, nil
}

//TODO userFeed(ctx, lastId, range, withConversation bool)
// IF withCon == true. Return users that have con with me shorted by last interaction. The rest to be returned alphabetical. Auto tha exei mesa ena feed me user name id and avatar. an einai withConv ta conv id, unread count & lastread message id
