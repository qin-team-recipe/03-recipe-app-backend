package dto

import "github.com/jackc/pgx/v5/pgtype"

type ChefLinkDto struct {
	Label pgtype.Text `json:"label"`
	Url   pgtype.Text `json:"url"`
}

type ChefLinkArrayDto []ChefLinkDto

type RecipeMethodDto struct {
	Html       pgtype.Text `json:"html"`
	Supplement any         `json:"supplement"`
}

type RecipeMethodArrayDto []RecipeMethodDto
