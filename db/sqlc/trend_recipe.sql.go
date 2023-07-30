// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: trend_recipe.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const fakeListTrendRecipe = `-- name: FakeListTrendRecipe :many
WITH
RECURSIVE generate_index (ii) AS (
    SELECT 0
    UNION ALL
    SELECT ii + 1
    FROM generate_index
)
SELECT
    GEN_RANDOM_UUID() AS recipe_id,
    '' AS title,
    '' AS introduction,
    'https://source.unsplash.com/random/300x300?v=1' AS image_url,
    0 AS num_fav,
    0 AS score
FROM
    generate_index
LIMIT $1
`

type FakeListTrendRecipeRow struct {
	RecipeID     pgtype.UUID `json:"recipeId"`
	Title        string      `json:"title"`
	Introduction string      `json:"introduction"`
	ImageUrl     string      `json:"imageUrl"`
	NumFav       int32       `json:"numFav"`
	Score        int32       `json:"score"`
}

func (q *Queries) FakeListTrendRecipe(ctx context.Context, lim int32) ([]FakeListTrendRecipeRow, error) {
	rows, err := q.db.Query(ctx, fakeListTrendRecipe, lim)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FakeListTrendRecipeRow
	for rows.Next() {
		var i FakeListTrendRecipeRow
		if err := rows.Scan(
			&i.RecipeID,
			&i.Title,
			&i.Introduction,
			&i.ImageUrl,
			&i.NumFav,
			&i.Score,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTrendRecipe = `-- name: ListTrendRecipe :many
WITH
history AS (
    SELECT
        SUM(CASE WHEN is_fav THEN 1 ELSE 0 END) - SUM(CASE WHEN is_fav THEN 0 ELSE 1 END) AS score,
        recipe_id
    FROM
        fav_history
    WHERE
        CURRENT_TIMESTAMP - INTERVAL '3 days' <= created_at
    GROUP BY
        recipe_id
)
SELECT
    history.recipe_id,
    recipe.title,
    recipe.introduction,
    recipe.image_url,
    recipe.num_fav,
    history.score
FROM
    history
INNER JOIN
    recipe
ON
    history.recipe_id = recipe.id
ORDER BY
    score DESC
LIMIT $1
`

type ListTrendRecipeRow struct {
	RecipeID     pgtype.UUID `json:"recipeId"`
	Title        string      `json:"title"`
	Introduction string      `json:"introduction"`
	ImageUrl     pgtype.Text `json:"imageUrl"`
	NumFav       int32       `json:"numFav"`
	Score        int32       `json:"score"`
}

func (q *Queries) ListTrendRecipe(ctx context.Context, lim int32) ([]ListTrendRecipeRow, error) {
	rows, err := q.db.Query(ctx, listTrendRecipe, lim)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListTrendRecipeRow
	for rows.Next() {
		var i ListTrendRecipeRow
		if err := rows.Scan(
			&i.RecipeID,
			&i.Title,
			&i.Introduction,
			&i.ImageUrl,
			&i.NumFav,
			&i.Score,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}