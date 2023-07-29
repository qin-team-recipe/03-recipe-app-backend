// Package docs provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.0 DO NOT EDIT.
package docs

import (
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

// CreateChefRecipe defines model for createChefRecipe.
type CreateChefRecipe struct {
	// AccessLevel 公開、限定公開、非公開、下書き
	AccessLevel int `json:"accessLevel"`

	// ChefId シェフID
	ChefId    openapi_types.UUID `json:"chefId"`
	CreatedAt string             `json:"createdAt"`

	// Id レシピID
	Id openapi_types.UUID `json:"id"`

	// ImageUrl 画像
	ImageUrl *string `json:"imageUrl,omitempty"`

	// Introduction レシピの紹介文
	Introduction *string `json:"introduction,omitempty"`

	// Link リンク
	Link *[]string `json:"link,omitempty"`

	// Method 作り方
	Method []struct {
		Html *string `json:"html,omitempty"`

		// Supplement 補足
		Supplement *map[string]interface{} `json:"supplement,omitempty"`
	} `json:"method"`

	// Name レシピ名
	Name string `json:"name"`

	// Servings ＊人前
	Servings  int    `json:"servings"`
	UpdatedAt string `json:"updatedAt"`

	// UsrId ユーザーID
	UsrId *openapi_types.UUID `json:"usrId,omitempty"`
}

// CreateUsrRecipe defines model for createUsrRecipe.
type CreateUsrRecipe struct {
	// AccessLevel 公開、限定公開、非公開、下書き
	AccessLevel int `json:"accessLevel"`

	// ChefId シェフID
	ChefId    *openapi_types.UUID `json:"chefId,omitempty"`
	CreatedAt string              `json:"createdAt"`

	// Id レシピID
	Id openapi_types.UUID `json:"id"`

	// ImageUrl 画像
	ImageUrl *string `json:"imageUrl,omitempty"`

	// Introduction レシピの紹介文
	Introduction *string `json:"introduction,omitempty"`

	// Link リンク
	Link *[]string `json:"link,omitempty"`

	// Method 作り方
	Method []struct {
		Html *string `json:"html,omitempty"`

		// Supplement 補足
		Supplement *map[string]interface{} `json:"supplement,omitempty"`
	} `json:"method"`

	// Name レシピ名
	Name string `json:"name"`

	// Servings ＊人前
	Servings  int    `json:"servings"`
	UpdatedAt string `json:"updatedAt"`

	// UsrId ユーザーID
	UsrId openapi_types.UUID `json:"usrId"`
}

// FeaturedChef defines model for featuredChef.
type FeaturedChef struct {
	Data []struct {
		// ChefId シェフID
		ChefId openapi_types.UUID `json:"chefId"`

		// ImageUrl プロフィール画像
		ImageUrl *string `json:"imageUrl,omitempty"`

		// Name シェフ登録名
		Name string `json:"name"`

		// NumFollower フォロワー数
		NumFollower int `json:"numFollower"`

		// Score 注目度→実際のデータはscoreの降順でソートされる
		Score int `json:"score"`
	} `json:"data"`
}

// TrendRecipe defines model for trendRecipe.
type TrendRecipe struct {
	Data []struct {
		// ImageUrl レシピ画像
		ImageUrl *string `json:"imageUrl,omitempty"`

		// Introduction レシピの紹介文
		Introduction *string `json:"introduction,omitempty"`

		// Name レシピ名
		Name string `json:"name"`

		// NumFav ファボられ数
		NumFav int `json:"numFav"`

		// RecipeId レシピID
		RecipeId openapi_types.UUID `json:"recipeId"`

		// Score 話題度→実際のデータはscoreの降順でソートされる
		Score int `json:"score"`
	} `json:"data"`
}

// PostApiCreateChefRecipeJSONBody defines parameters for PostApiCreateChefRecipe.
type PostApiCreateChefRecipeJSONBody struct {
	// AccessLevel 公開、限定公開、非公開、下書き
	AccessLevel int `json:"accessLevel"`

	// ChefId シェフID
	ChefId openapi_types.UUID `json:"chefId"`

	// ImageUrl 画像
	ImageUrl *string `json:"imageUrl,omitempty"`

	// Introduction レシピの紹介文
	Introduction *string `json:"introduction,omitempty"`

	// Link リンク
	Link *[]string `json:"link,omitempty"`

	// Method 作り方
	Method []struct {
		Html *string `json:"html,omitempty"`

		// Supplement 補足
		Supplement *map[string]interface{} `json:"supplement,omitempty"`
	} `json:"method"`

	// Name レシピ名
	Name string `json:"name"`

	// Servings ＊人前
	Servings int `json:"servings"`
}

// PostApiCreateUsrRecipeJSONBody defines parameters for PostApiCreateUsrRecipe.
type PostApiCreateUsrRecipeJSONBody struct {
	// AccessLevel 公開、限定公開、非公開、下書き
	AccessLevel int `json:"accessLevel"`

	// ImageUrl 画像
	ImageUrl *string `json:"imageUrl,omitempty"`

	// Introduction レシピの紹介文
	Introduction *string `json:"introduction,omitempty"`

	// Link リンク
	Link *[]string `json:"link,omitempty"`

	// Method 作り方
	Method []struct {
		Html *string `json:"html,omitempty"`

		// Supplement 補足
		Supplement *map[string]interface{} `json:"supplement,omitempty"`
	} `json:"method"`

	// Name レシピ名
	Name string `json:"name"`

	// Servings ＊人前
	Servings int `json:"servings"`

	// UsrId ユーザーID
	UsrId openapi_types.UUID `json:"usrId"`
}

// PostApiCreateChefRecipeJSONRequestBody defines body for PostApiCreateChefRecipe for application/json ContentType.
type PostApiCreateChefRecipeJSONRequestBody PostApiCreateChefRecipeJSONBody

// PostApiCreateUsrRecipeJSONRequestBody defines body for PostApiCreateUsrRecipe for application/json ContentType.
type PostApiCreateUsrRecipeJSONRequestBody PostApiCreateUsrRecipeJSONBody
