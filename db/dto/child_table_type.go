package dto

import "github.com/jackc/pgx/v5/pgtype"

type RecipeIngredientDto struct {
	ID pgtype.UUID `json:"id"`
	//RecipeID   pgtype.UUID        `json:"recipeId"`
	//Idx        int32              `json:"idx"`
	Name       string      `json:"name"`
	Supplement pgtype.Text `json:"supplement"`
	//CreatedAt  pgtype.Timestamptz `json:"createdAt"`
	//UpdatedAt  pgtype.Timestamptz `json:"updatedAt"`
}

type RecipeIngredientArrayDto []RecipeIngredientDto

type ShoppingItemDto struct {
	ID           pgtype.UUID `json:"id"`
	IngredientId pgtype.UUID `json:"ingredientId"`
	Name         string      `json:"name"`
	Supplement   pgtype.Text `json:"supplement"`
	//CreatedAt  pgtype.Timestamptz `json:"createdAt"`
	//UpdatedAt  pgtype.Timestamptz `json:"updatedAt"`
}

type ShoppingItemArrayDto []ShoppingItemDto
