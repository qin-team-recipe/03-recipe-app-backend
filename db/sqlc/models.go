// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// シェフ
type Chef struct {
	ID string `json:"id"`
	// ログインemail
	Email pgtype.Text `json:"email"`
	// 登録名
	Name string `json:"name"`
	// 登録画像
	ImageUrl pgtype.Text `json:"imageUrl"`
	// シェフコメント
	Comment   pgtype.Text        `json:"comment"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt"`
}

// ファボ／リム履歴
type FavHistory struct {
	RecipeID string `json:"recipeID"`
	UsrID    string `json:"usrID"`
	// ファボ／リム
	IsFav     bool               `json:"isFav"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

// ファボ中
type Favoring struct {
	ID        string             `json:"id"`
	RecipeID  string             `json:"recipeID"`
	UsrID     string             `json:"usrID"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

// フォロー／リム履歴
type FollowHistory struct {
	ChefID string `json:"chefID"`
	UsrID  string `json:"usrID"`
	// フォロー／リム
	IsFollow  bool               `json:"isFollow"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

// フォロー中:Table to show favorite chefs of a user
type Following struct {
	ID        string             `json:"id"`
	ChefID    string             `json:"chefID"`
	UsrID     string             `json:"usrID"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

// シェフのレシピ＆マイレシピ
type Recipe struct {
	ID     string      `json:"id"`
	ChefID pgtype.Text `json:"chefID"`
	UsrID  pgtype.Text `json:"usrID"`
	// レシピタイトル
	Title string `json:"title"`
	// レシピコメント
	Comment string `json:"comment"`
	// ＊人前
	Servings int32 `json:"servings"`
	// 材料
	Ingredient []string `json:"ingredient"`
	// 作り方
	Method []string `json:"method"`
	// 画像
	ImageUrl pgtype.Text `json:"imageUrl"`
	// リンク
	Link []string `json:"link"`
	// 公開？:公開、限定公開、非公開、下書き
	AccessLevel int32              `json:"accessLevel"`
	CreatedAt   pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt   pgtype.Timestamptz `json:"updatedAt"`
}

// 買い物リスト
type ShoppingList struct {
	ID    string `json:"id"`
	UsrID string `json:"usrID"`
	// NULL=メモ
	RecipeID pgtype.Text `json:"recipeID"`
	// 清書or下書き
	IsFairCopy bool `json:"isFairCopy"`
	// 明細
	ShoppingItem []string           `json:"shoppingItem"`
	CreatedAt    pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt    pgtype.Timestamptz `json:"updatedAt"`
}

// sns
type Sn struct {
	ID           string      `json:"id"`
	ChefID       string      `json:"chefID"`
	Name         string      `json:"name"`
	AccountName  string      `json:"accountName"`
	NumFollowers int32       `json:"numFollowers"`
	Link         pgtype.Text `json:"link"`
}

// ユーザー
type Usr struct {
	ID string `json:"id"`
	// ログインemail
	Email string `json:"email"`
	// 登録名
	Name string `json:"name"`
	// 登録画像
	ImageUrl  pgtype.Text        `json:"imageUrl"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt"`
}
