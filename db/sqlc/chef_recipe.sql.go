// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: chef_recipe.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const deleteChefRecipe = `-- name: DeleteChefRecipe :one
DELETE FROM
    recipe
WHERE
    id = $1
RETURNING
    id,
    chef_id,
    usr_id,
    name,
    servings,
    image_url,
    introduction,
    link,
    access_level,
    created_at,
    updated_at,
    num_fav
`

type DeleteChefRecipeRow struct {
	ID           pgtype.UUID        `json:"id"`
	ChefID       pgtype.UUID        `json:"chefId"`
	UsrID        pgtype.UUID        `json:"usrId"`
	Name         string             `json:"name"`
	Servings     int32              `json:"servings"`
	ImageUrl     pgtype.Text        `json:"imageUrl"`
	Introduction pgtype.Text        `json:"introduction"`
	Link         []string           `json:"link"`
	AccessLevel  int32              `json:"accessLevel"`
	CreatedAt    pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt    pgtype.Timestamptz `json:"updatedAt"`
	NumFav       int32              `json:"numFav"`
}

func (q *Queries) DeleteChefRecipe(ctx context.Context, id pgtype.UUID) (DeleteChefRecipeRow, error) {
	row := q.db.QueryRow(ctx, deleteChefRecipe, id)
	var i DeleteChefRecipeRow
	err := row.Scan(
		&i.ID,
		&i.ChefID,
		&i.UsrID,
		&i.Name,
		&i.Servings,
		&i.ImageUrl,
		&i.Introduction,
		&i.Link,
		&i.AccessLevel,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.NumFav,
	)
	return i, err
}

const getChefRecipes = `-- name: GetChefRecipes :many
SELECT
    id, chef_id, usr_id, name, servings, ingredient, method, image_url, introduction, link, access_level, created_at, updated_at, num_fav
FROM
    v_recipe
WHERE
    chef_id = $1
ORDER BY
    created_at DESC
`

func (q *Queries) GetChefRecipes(ctx context.Context, chefID pgtype.UUID) ([]VRecipe, error) {
	rows, err := q.db.Query(ctx, getChefRecipes, chefID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []VRecipe
	for rows.Next() {
		var i VRecipe
		if err := rows.Scan(
			&i.ID,
			&i.ChefID,
			&i.UsrID,
			&i.Name,
			&i.Servings,
			&i.Ingredient,
			&i.Method,
			&i.ImageUrl,
			&i.Introduction,
			&i.Link,
			&i.AccessLevel,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.NumFav,
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

const searchChefRecipe = `-- name: SearchChefRecipe :many
SELECT
    id,
    chef_id,
    name,
    servings,
    image_url,
    introduction,
    access_level,
    created_at,
    updated_at,
    num_fav
FROM
    recipe
WHERE
    chef_id IS NOT NULL
AND (
    name &@~ $1
OR
    introduction &@~ $1
)
ORDER BY
    pgroonga_score(tableoid, ctid) DESC,
    num_fav DESC
`

type SearchChefRecipeRow struct {
	ID           pgtype.UUID        `json:"id"`
	ChefID       pgtype.UUID        `json:"chefId"`
	Name         string             `json:"name"`
	Servings     int32              `json:"servings"`
	ImageUrl     pgtype.Text        `json:"imageUrl"`
	Introduction pgtype.Text        `json:"introduction"`
	AccessLevel  int32              `json:"accessLevel"`
	CreatedAt    pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt    pgtype.Timestamptz `json:"updatedAt"`
	NumFav       int32              `json:"numFav"`
}

func (q *Queries) SearchChefRecipe(ctx context.Context, query string) ([]SearchChefRecipeRow, error) {
	rows, err := q.db.Query(ctx, searchChefRecipe, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchChefRecipeRow
	for rows.Next() {
		var i SearchChefRecipeRow
		if err := rows.Scan(
			&i.ID,
			&i.ChefID,
			&i.Name,
			&i.Servings,
			&i.ImageUrl,
			&i.Introduction,
			&i.AccessLevel,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.NumFav,
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