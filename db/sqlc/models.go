// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	dto "github.com/aopontann/gin-sqlc/db/dto"
	"github.com/jackc/pgx/v5/pgtype"
)

// シェフ
type Chef struct {
	ID pgtype.UUID `json:"id"`
	// ログインemail
	Email pgtype.Text `json:"email"`
	// 登録名
	Name string `json:"name"`
	// プロフィール画像
	ImageUrl pgtype.Text `json:"imageUrl"`
	// シェフ紹介
	Profile pgtype.Text `json:"profile"`
	// リンク
	Link []string `json:"link"`
	// google／apple
	AuthServer pgtype.Text `json:"authServer"`
	// 認証USERINFO
	AuthUserinfo []byte             `json:"authUserinfo"`
	CreatedAt    pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt    pgtype.Timestamptz `json:"updatedAt"`
	// レシピ数
	NumRecipe int32 `json:"numRecipe"`
	// フォロワー数
	NumFollower int32 `json:"numFollower"`
}

// ファボ／リム履歴
type FavHistory struct {
	RecipeID pgtype.UUID `json:"recipeID"`
	UsrID    pgtype.UUID `json:"usrID"`
	// ファボ／リム
	IsFav     bool               `json:"isFav"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

// ファボ中
type Favoring struct {
	ID        pgtype.UUID        `json:"id"`
	RecipeID  pgtype.UUID        `json:"recipeID"`
	UsrID     pgtype.UUID        `json:"usrID"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

// フォロー／リム履歴
type FollowChefHistory struct {
	ChefID pgtype.UUID `json:"chefID"`
	UsrID  pgtype.UUID `json:"usrID"`
	// フォロー／リム
	IsFollow  bool               `json:"isFollow"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

// フォロー中シェフ:Table to show favorite chefs of a user
type FollowingChef struct {
	ID        pgtype.UUID        `json:"id"`
	ChefID    pgtype.UUID        `json:"chefID"`
	UsrID     pgtype.UUID        `json:"usrID"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

// フォロー中ユーザー:Table to show favorite chefs of a user
type FollowingUser struct {
	ID         pgtype.UUID        `json:"id"`
	FolloweeID pgtype.UUID        `json:"followeeID"`
	FollowerID pgtype.UUID        `json:"followerID"`
	CreatedAt  pgtype.Timestamptz `json:"createdAt"`
}

// 材料
type Ingredient struct {
	ID       pgtype.UUID `json:"id"`
	RecipeID pgtype.UUID `json:"recipeID"`
	// インデックス
	Idx int32 `json:"idx"`
	// 材料名
	Name string `json:"name"`
	// 補足
	Supplement pgtype.Text        `json:"supplement"`
	CreatedAt  pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt  pgtype.Timestamptz `json:"updatedAt"`
}

// シェフのレシピ＆マイレシピ
type Recipe struct {
	ID     pgtype.UUID `json:"id"`
	ChefID pgtype.UUID `json:"chefID"`
	UsrID  pgtype.UUID `json:"usrID"`
	// レシピタイトル
	Title string `json:"title"`
	// ＊人前
	Servings int32 `json:"servings"`
	// 作り方
	Method []string `json:"method"`
	// 画像
	ImageUrl pgtype.Text `json:"imageUrl"`
	// レシピの紹介文
	Introduction string `json:"introduction"`
	// 公開等:公開、限定公開、非公開、下書き
	AccessLevel int32              `json:"accessLevel"`
	CreatedAt   pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt   pgtype.Timestamptz `json:"updatedAt"`
	// ファボられ数
	NumFav int32 `json:"numFav"`
}

// 買い物明細
type ShoppingItem struct {
	ID             pgtype.UUID `json:"id"`
	ShoppingListID pgtype.UUID `json:"shoppingListID"`
	IngredientID   pgtype.UUID `json:"ingredientID"`
	// インデックス
	Idx       int32              `json:"idx"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt"`
}

// 買い物リスト
type ShoppingList struct {
	ID    pgtype.UUID `json:"id"`
	UsrID pgtype.UUID `json:"usrID"`
	// NULL：メモリスト／削除レシピ
	RecipeID pgtype.UUID `json:"recipeID"`
	// 「*人前」「メモリスト」
	Description pgtype.Text `json:"description"`
	// 清書or下書き
	IsFairCopy bool               `json:"isFairCopy"`
	CreatedAt  pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt  pgtype.Timestamptz `json:"updatedAt"`
}

// ユーザー（一般シェフ）
type Usr struct {
	ID pgtype.UUID `json:"id"`
	// ログインemail
	Email string `json:"email"`
	// ニックネーム
	Name string `json:"name"`
	// プロフィール画像（任意）
	ImageUrl pgtype.Text `json:"imageUrl"`
	// 自己紹介（任意）
	Profile pgtype.Text `json:"profile"`
	// リンク（任意）
	Link []string `json:"link"`
	// google／apple
	AuthServer string `json:"authServer"`
	// 認証USERINFO
	AuthUserinfo []byte             `json:"authUserinfo"`
	CreatedAt    pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt    pgtype.Timestamptz `json:"updatedAt"`
	// マイレシピ数
	NumRecipe int32 `json:"numRecipe"`
}

type VChef struct {
	ID           pgtype.UUID          `json:"id"`
	Email        pgtype.Text          `json:"email"`
	Name         string               `json:"name"`
	ImageUrl     pgtype.Text          `json:"imageUrl"`
	Profile      pgtype.Text          `json:"profile"`
	Link         dto.ChefLinkArrayDto `json:"link"`
	AuthServer   pgtype.Text          `json:"authServer"`
	AuthUserinfo []byte               `json:"authUserinfo"`
	CreatedAt    pgtype.Timestamptz   `json:"createdAt"`
	UpdatedAt    pgtype.Timestamptz   `json:"updatedAt"`
	NumRecipe    int32                `json:"numRecipe"`
	NumFollower  int32                `json:"numFollower"`
}

type VRecipe struct {
	ID           pgtype.UUID              `json:"id"`
	ChefID       pgtype.UUID              `json:"chefID"`
	UsrID        pgtype.UUID              `json:"usrID"`
	Title        string                   `json:"title"`
	Servings     int32                    `json:"servings"`
	Method       dto.RecipeMethodArrayDto `json:"method"`
	ImageUrl     pgtype.Text              `json:"imageUrl"`
	Introduction string                   `json:"introduction"`
	AccessLevel  int32                    `json:"accessLevel"`
	CreatedAt    pgtype.Timestamptz       `json:"createdAt"`
	UpdatedAt    pgtype.Timestamptz       `json:"updatedAt"`
}

type VUsr struct {
	ID           pgtype.UUID          `json:"id"`
	Email        string               `json:"email"`
	Name         string               `json:"name"`
	ImageUrl     pgtype.Text          `json:"imageUrl"`
	Profile      pgtype.Text          `json:"profile"`
	Link         dto.ChefLinkArrayDto `json:"link"`
	AuthServer   string               `json:"authServer"`
	AuthUserinfo []byte               `json:"authUserinfo"`
	CreatedAt    pgtype.Timestamptz   `json:"createdAt"`
	UpdatedAt    pgtype.Timestamptz   `json:"updatedAt"`
	NumRecipe    int32                `json:"numRecipe"`
}
