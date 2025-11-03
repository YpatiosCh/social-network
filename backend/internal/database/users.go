package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

func (db *Database) AddAuthUser(req models.RegisterDbRequest) (res models.RegisterDbResponse, err error) {
	res = models.RegisterDbResponse{}
	tx, err := db.Pool.Begin(req.Ctx)
	if err != nil {
		return res, err
	}
	defer func() { _ = tx.Rollback(req.Ctx) }()

	if err := tx.QueryRow(req.Ctx,
		`INSERT INTO master_index (content_type) VALUES ('user') RETURNING id`,
	).Scan(&res.Id); err != nil {
		fmt.Println("error inserting into master_index:", err)
		return res, err
	}

	if err := tx.QueryRow(req.Ctx, `
        INSERT INTO users (id, username, gender, first_name, last_name, email, age, avatar)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING username, avatar
    `, res.Id, req.Username, req.Gender, req.FirstName, req.LastName, req.Email, req.Age, req.Avatar).
		Scan(&res.RegisterResponseUser.UserName, &res.Avatar); err != nil {
		fmt.Println("error inserting into users:", err)
		return res, err
	}
	// insert credentials (fail if username already taken)
	_, err = tx.Exec(req.Ctx, `
		INSERT INTO auth_user (user_id, identifier, password_hash)
		VALUES ($1, $2, $3)
	`, res.Id, req.Identifier, req.PasswordHash)
	if err != nil {
		// map unique violations to a nicer error (username already in use)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return res, fmt.Errorf("username already taken")
		}
		return res, err
	}
	if err := tx.Commit(req.Ctx); err != nil {
		return res, err
	}
	return res, nil
}

func (db *Database) Authenticate(ctx context.Context, identifier, password string) (int64, string, string, error) {
	const query = `
		SELECT u.id, u.username, COALESCE(u.avatar, ''), au.password_hash
		FROM auth_user AS au
		JOIN users     AS u ON u.id = au.user_id
		WHERE u.username = $1 OR u.email = $1
		LIMIT 1;`

	var id int64
	var hash string
	var username string
	var avatar string

	if err := db.Pool.QueryRow(ctx, query, identifier).Scan(&id, &username, &avatar, &hash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, "", "", fmt.Errorf("invalid credentials")
		}
		return 0, "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return 0, "", "", fmt.Errorf("invalid credentials")
	}

	return id, username, avatar, nil
}

func (db *Database) GetUserInfos(ctx context.Context, userID int64) (*Users, error) {
	const query = `
        SELECT id, username, gender, first_name, last_name, email, age, avatar
        FROM users
        WHERE id = $1
    `

	var u Users
	err := db.Pool.QueryRow(ctx, query, userID).Scan(
		&u.Id,
		&u.Username,
		&u.Gender,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Age,
		&u.Avatar,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
