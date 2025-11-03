package database

import (
	"context"
	"log"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

func (db *Database) AllCategories(ctx context.Context) ([]models.Category, error) {
	query := `SELECT id, name FROM categories ORDER BY name`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer rows.Close()

	var categ []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(
			&c.Id,
			&c.Name,
		); err != nil {
			return nil, err
		}

		categ = append(categ, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categ, nil
}

func (db *Database) CategoriesById(ctx context.Context, id int64) (*models.Category, error) {
	query := `SELECT id, name FROM categories WHERE id = $1`

	var categ models.Category
	err := db.Pool.QueryRow(ctx, query, id).Scan(&categ.Id, &categ.Name)
	if err != nil {
		return nil, err
	}
	return &categ, nil
}
