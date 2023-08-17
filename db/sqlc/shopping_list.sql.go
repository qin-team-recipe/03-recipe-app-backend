// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: shopping_list.sql

package db

import (
	"context"

	dto "github.com/aopontann/gin-sqlc/db/dto"
	"github.com/jackc/pgx/v5/pgtype"
)

const createShoppingList = `-- name: CreateShoppingList :one
INSERT INTO shopping_list (
    usr_id,
    recipe_id,
    r_idx,
    description,
    is_fair_copy
)
VALUES
(
    $1,
    $2,
    (SELECT
        COALESCE(MAX(r_idx) + 1, 1)
    FROM
        shopping_list
    WHERE
        usr_id = $1
    ),
    $3,
    $4
)
RETURNING
    id,
    usr_id,
    recipe_id,
    r_idx,
    description,
    is_fair_copy,
    created_at,
    updated_at
`

type CreateShoppingListParams struct {
	UsrID       pgtype.UUID `json:"usrId"`
	RecipeID    pgtype.UUID `json:"recipeId"`
	Description pgtype.Text `json:"description"`
	IsFairCopy  bool        `json:"isFairCopy"`
}

func (q *Queries) CreateShoppingList(ctx context.Context, arg CreateShoppingListParams) (ShoppingList, error) {
	row := q.db.QueryRow(ctx, createShoppingList,
		arg.UsrID,
		arg.RecipeID,
		arg.Description,
		arg.IsFairCopy,
	)
	var i ShoppingList
	err := row.Scan(
		&i.ID,
		&i.UsrID,
		&i.RecipeID,
		&i.RIdx,
		&i.Description,
		&i.IsFairCopy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteShoppingList = `-- name: DeleteShoppingList :one
DELETE FROM
    shopping_list
WHERE
    usr_id = $1
AND
    id = $2
RETURNING
    id,
    usr_id,
    recipe_id,
    r_idx,
    description,
    is_fair_copy,
    created_at,
    updated_at
`

type DeleteShoppingListParams struct {
	UsrID pgtype.UUID `json:"usrId"`
	ID    pgtype.UUID `json:"id"`
}

func (q *Queries) DeleteShoppingList(ctx context.Context, arg DeleteShoppingListParams) (ShoppingList, error) {
	row := q.db.QueryRow(ctx, deleteShoppingList, arg.UsrID, arg.ID)
	var i ShoppingList
	err := row.Scan(
		&i.ID,
		&i.UsrID,
		&i.RecipeID,
		&i.RIdx,
		&i.Description,
		&i.IsFairCopy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getShoppingList = `-- name: GetShoppingList :one
SELECT
    id,
    usr_id,
    recipe_id,
    r_idx,
    recipe_name,
    chef_name,
    general_chef_name,
    description,
    is_fair_copy,
    created_at,
    updated_at,
    item
FROM
    v_shopping_list
WHERE
    usr_id = $1
AND
    recipe_id = $2
`

type GetShoppingListParams struct {
	UsrID    pgtype.UUID `json:"usrId"`
	RecipeID pgtype.UUID `json:"recipeId"`
}

type GetShoppingListRow struct {
	ID              pgtype.UUID              `json:"id"`
	UsrID           pgtype.UUID              `json:"usrId"`
	RecipeID        pgtype.UUID              `json:"recipeId"`
	RIdx            int32                    `json:"rIdx"`
	RecipeName      string                   `json:"recipeName"`
	ChefName        pgtype.Text              `json:"chefName"`
	GeneralChefName pgtype.Text              `json:"generalChefName"`
	Description     pgtype.Text              `json:"description"`
	IsFairCopy      bool                     `json:"isFairCopy"`
	CreatedAt       pgtype.Timestamptz       `json:"createdAt"`
	UpdatedAt       pgtype.Timestamptz       `json:"updatedAt"`
	Item            dto.ShoppingItemArrayDto `json:"item"`
}

func (q *Queries) GetShoppingList(ctx context.Context, arg GetShoppingListParams) (GetShoppingListRow, error) {
	row := q.db.QueryRow(ctx, getShoppingList, arg.UsrID, arg.RecipeID)
	var i GetShoppingListRow
	err := row.Scan(
		&i.ID,
		&i.UsrID,
		&i.RecipeID,
		&i.RIdx,
		&i.RecipeName,
		&i.ChefName,
		&i.GeneralChefName,
		&i.Description,
		&i.IsFairCopy,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Item,
	)
	return i, err
}

const innerCreateShoppingItem = `-- name: InnerCreateShoppingItem :one
INSERT INTO shopping_item (
    shopping_list_id,
    ingredient_id,
    idx,
    name,
    supplement
)
VALUES
(
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING
    id,
    ingredient_id,
    name,
    supplement
`

type InnerCreateShoppingItemParams struct {
	ShoppingListID pgtype.UUID `json:"shoppingListId"`
	IngredientID   pgtype.UUID `json:"ingredientId"`
	Idx            int32       `json:"idx"`
	Name           string      `json:"name"`
	Supplement     pgtype.Text `json:"supplement"`
}

type InnerCreateShoppingItemRow struct {
	ID           pgtype.UUID `json:"id"`
	IngredientID pgtype.UUID `json:"ingredientId"`
	Name         string      `json:"name"`
	Supplement   pgtype.Text `json:"supplement"`
}

func (q *Queries) InnerCreateShoppingItem(ctx context.Context, arg InnerCreateShoppingItemParams) (InnerCreateShoppingItemRow, error) {
	row := q.db.QueryRow(ctx, innerCreateShoppingItem,
		arg.ShoppingListID,
		arg.IngredientID,
		arg.Idx,
		arg.Name,
		arg.Supplement,
	)
	var i InnerCreateShoppingItemRow
	err := row.Scan(
		&i.ID,
		&i.IngredientID,
		&i.Name,
		&i.Supplement,
	)
	return i, err
}

const innerDeleteNotAnyShoppingItem = `-- name: InnerDeleteNotAnyShoppingItem :exec
DELETE FROM
    shopping_item
WHERE
    shopping_list_id = $1
AND
    NOT (id = ANY ($2::UUID[]))
`

type InnerDeleteNotAnyShoppingItemParams struct {
	ShoppingListID pgtype.UUID   `json:"shoppingListId"`
	ID             []pgtype.UUID `json:"id"`
}

func (q *Queries) InnerDeleteNotAnyShoppingItem(ctx context.Context, arg InnerDeleteNotAnyShoppingItemParams) error {
	_, err := q.db.Exec(ctx, innerDeleteNotAnyShoppingItem, arg.ShoppingListID, arg.ID)
	return err
}

const innerUpdateShoppingItem = `-- name: InnerUpdateShoppingItem :one
UPDATE shopping_item SET
    idx           = $1,
    name          = $2,
    supplement    = $3
WHERE
    (ingredient_id = $4 OR $4 IS NULL)
AND
    shopping_list_id = $5
AND
    id = $6
RETURNING
    id,
    ingredient_id,
    name,
    supplement
`

type InnerUpdateShoppingItemParams struct {
	Idx            int32       `json:"idx"`
	Name           string      `json:"name"`
	Supplement     pgtype.Text `json:"supplement"`
	IngredientID   pgtype.UUID `json:"ingredientId"`
	ShoppingListID pgtype.UUID `json:"shoppingListId"`
	ID             pgtype.UUID `json:"id"`
}

type InnerUpdateShoppingItemRow struct {
	ID           pgtype.UUID `json:"id"`
	IngredientID pgtype.UUID `json:"ingredientId"`
	Name         string      `json:"name"`
	Supplement   pgtype.Text `json:"supplement"`
}

func (q *Queries) InnerUpdateShoppingItem(ctx context.Context, arg InnerUpdateShoppingItemParams) (InnerUpdateShoppingItemRow, error) {
	row := q.db.QueryRow(ctx, innerUpdateShoppingItem,
		arg.Idx,
		arg.Name,
		arg.Supplement,
		arg.IngredientID,
		arg.ShoppingListID,
		arg.ID,
	)
	var i InnerUpdateShoppingItemRow
	err := row.Scan(
		&i.ID,
		&i.IngredientID,
		&i.Name,
		&i.Supplement,
	)
	return i, err
}

const listShoppingList = `-- name: ListShoppingList :many
SELECT
    id,
    usr_id,
    recipe_id,
    r_idx,
    recipe_name,
    chef_name,
    general_chef_name,
    description,
    is_fair_copy,
    created_at,
    updated_at,
    item
FROM
    v_shopping_list
WHERE
    usr_id = $1
ORDER BY
    r_idx DESC
`

type ListShoppingListRow struct {
	ID              pgtype.UUID              `json:"id"`
	UsrID           pgtype.UUID              `json:"usrId"`
	RecipeID        pgtype.UUID              `json:"recipeId"`
	RIdx            int32                    `json:"rIdx"`
	RecipeName      string                   `json:"recipeName"`
	ChefName        pgtype.Text              `json:"chefName"`
	GeneralChefName pgtype.Text              `json:"generalChefName"`
	Description     pgtype.Text              `json:"description"`
	IsFairCopy      bool                     `json:"isFairCopy"`
	CreatedAt       pgtype.Timestamptz       `json:"createdAt"`
	UpdatedAt       pgtype.Timestamptz       `json:"updatedAt"`
	Item            dto.ShoppingItemArrayDto `json:"item"`
}

func (q *Queries) ListShoppingList(ctx context.Context, usrID pgtype.UUID) ([]ListShoppingListRow, error) {
	rows, err := q.db.Query(ctx, listShoppingList, usrID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListShoppingListRow
	for rows.Next() {
		var i ListShoppingListRow
		if err := rows.Scan(
			&i.ID,
			&i.UsrID,
			&i.RecipeID,
			&i.RIdx,
			&i.RecipeName,
			&i.ChefName,
			&i.GeneralChefName,
			&i.Description,
			&i.IsFairCopy,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Item,
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

const updateShoppingList = `-- name: UpdateShoppingList :one
UPDATE shopping_list SET
    r_idx        = $1,
    description  = $2,
    is_fair_copy = $3
WHERE
    usr_id = $4
AND
    id = $5
RETURNING
    id,
    usr_id,
    recipe_id,
    r_idx,
    description,
    is_fair_copy,
    created_at,
    updated_at
`

type UpdateShoppingListParams struct {
	RIdx        int32       `json:"rIdx"`
	Description pgtype.Text `json:"description"`
	IsFairCopy  bool        `json:"isFairCopy"`
	UsrID       pgtype.UUID `json:"usrId"`
	ID          pgtype.UUID `json:"id"`
}

func (q *Queries) UpdateShoppingList(ctx context.Context, arg UpdateShoppingListParams) (ShoppingList, error) {
	row := q.db.QueryRow(ctx, updateShoppingList,
		arg.RIdx,
		arg.Description,
		arg.IsFairCopy,
		arg.UsrID,
		arg.ID,
	)
	var i ShoppingList
	err := row.Scan(
		&i.ID,
		&i.UsrID,
		&i.RecipeID,
		&i.RIdx,
		&i.Description,
		&i.IsFairCopy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
