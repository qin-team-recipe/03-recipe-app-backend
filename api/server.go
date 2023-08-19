package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"

	db "github.com/aopontann/gin-sqlc/db/sqlc"
)

type Server struct {
	r   *gin.Engine
	db  *pgx.Conn
	q   *db.Queries
	rbd *redis.Client
}

func NewServer(conn *pgx.Conn, rdb *redis.Client) *Server {
	engine := gin.Default()
	server := &Server{
		r:   engine,
		db:  conn,
		q:   db.New(conn),
		rbd: rdb,
	}
	return server
}

func (s *Server) MountHandlers() {
	auth := s.r.Group("/auth")
	api := s.r.Group("/api")

	// api.POST("/users", s.RegisterUser)
	// api.POST("/users/login", s.LoginUser)
	// user.Use(AuthMiddleware())

	// 認証認可関連
	auth.GET("/google/login", s.OauthGoogleLogin)
	auth.GET("/google/callback", s.OauthGoogleCallback)

	// ユーザー関連
	api.POST("/users", s.CreateUser) // ユーザーを新規登録するAPI
	api.GET("/users/:id", s.GetUser) // ユーザーを取得するAPI

	//// 仮で作成　セッションの説明用 ////
	// グループを作成
	usr := api.Group("/user")
	// /api/user/* リンクにアクセスしたとき、登録したハンドラ関数(ここではGetUserId)を実行する前に Authentication() を実行するようにする処理
	usr.Use(s.Authentication())
	// ユーザIDのみ返すAPI
	usr.GET("/id", s.GetUserId)

	usr.GET("/users", s.GetSelf)       // 自分を取得するAPI
	usr.PUT("/users", s.UpdateSelf)    // 自分を更新するAPI
	usr.DELETE("/users", s.DeleteSelf) // 自分を削除するAPI

	// ユーザー（= 一般シェフ）のレシピ関連
	//usr.GET("/users/recipe", s.GetUsrRecipe)
	usr.PUT("/users/recipe/:recipe_id", s.UpdateUserRecipe)    // 一般シェフのマイレシピを更新するAPI
	usr.DELETE("/users/recipe/:recipe_id", s.DeleteUserRecipe) // 一般シェフのマイレシピを削除するAPI
	usr.POST("/users/recipe", s.CreateUsrRecipe)               // 一般シェフのマイレシピを新規登録するAPI

	// 有名シェフのレシピ関連
	//api.GET("/chefs/:id/recipe", s.GetChefRecipe)
	api.PUT("/chefs/recipe/:recipe_id", s.UpdateChefRecipe)    // 有名シェフのレシピを更新するAPI
	api.DELETE("/chefs/recipe/:recipe_id", s.DeleteChefRecipe) // 有名シェフのレシピを削除するAPI
	api.POST("/chefs/:id/recipe", s.CreateChefRecipe)          // 有名シェフのレシピを新規登録するAPI
	api.GET("/chefs/recipe/search", s.SearchChefRecipe)        // 有名シェフのレシピを全文検索するAPI

	// 有名シェフ関連
	api.GET("/chefs/:id", s.GetChef)               // 有名シェフを取得するAPI
	api.PUT("/chefs/:id", s.UpdateChef)            // 有名シェフを更新するAPI
	api.DELETE("/chefs/:id", s.DeleteChef)         // 有名シェフを削除するAPI
	api.POST("/chefs", s.CreateChef)               // 有名シェフを新規登録するAPI
	api.GET("/chefs/featured", s.ListFeaturedChef) // 注目の有名シェフ一覧を取得するAPI
	api.GET("/chefs/search", s.SearchChef)         // 有名シェフを全文検索するAPI

	// レシピ関連
	api.GET("/recipes/:id", s.GetRecipe)         // レシピを取得するAPI
	api.GET("/recipes/trend", s.ListTrendRecipe) // 話題のレシピ一覧を取得するAPI

	// 買い物リスト関連
	usr.GET("/lists", s.ListShoppingList)                  // ユーザーの買い物リスト一覧を取得するAPI
	usr.GET("/lists/recipe/:recipe_id", s.GetShoppingList) // 買い物リストを取得するAPI
	usr.PUT("/lists/:id", s.UpdateShoppingList)            // 買い物リストを更新するAPI
	usr.DELETE("/lists/:id", s.DeleteShoppingList)         // 買い物リストを削除するAPI
	usr.POST("/lists", s.CreateShoppingList)               // 買い物リストを新規登録するAPI

	// 画像関連
	api.GET("/images", s.GetImage)   // 画像を取得するAPI（webp限定）
	api.POST("/images", s.PostImage) // 画像を新規登録するAPI（webp形式に変換される）

	// 有名シェフフォロー関連
	usr.POST("/follow/chefs/:id", s.CreateFollowChef)           // 有名シェフをフォローするAPI
	usr.DELETE("/follow/chefs/:id", s.DeleteFollowChef)         // 有名シェフのフォローを解除するAPI
	usr.GET("/follow/chefs/:id", s.ExistsFollowChef)            // 有名シェフをフォローしているか
	usr.GET("/follow/chefs", s.ListFollowChef)                  // フォローしている有名シェフの一覧を取得するAPI
	usr.GET("/follow/chefs/recipes", s.ListFollowChefNewRecipe) // フォローしているシェフの新着レシピ一覧を取得するAPI

	// 一般シェフフォロー関連
	usr.POST("/follow/users/:id", s.CreateFollowUser)   // 一般シェフをフォローするAPI
	usr.DELETE("/follow/users/:id", s.DeleteFollowUser) // 一般シェフのフォローを解除するAPI
	usr.GET("/follow/users/:id", s.ExistsFollowUser)    // 一般シェフをフォローしているか
	usr.GET("/follow/users", s.ListFollowUser)          // フォローしている一般シェフの一覧を取得するAPI

	// お気に入りレシピ関連
	usr.POST("/favorite/recipes/:id", s.CreateFavoriteRecipe)   // お気に入りレシピ登録API
	usr.DELETE("/favorite/recipes/:id", s.DeleteFavoriteRecipe) // お気に入りレシピ解除API
	usr.GET("/favorite/recipes/:id", s.ExistsFavoriteRecipe)    // お気に入りレシピとして登録しているか確認API
	usr.GET("/favorite/recipes", s.ListFavoriteRecipe)          // お気に入りレシピの一覧を取得するAPI
}

func (s *Server) Start(addr string) error {
	return s.r.Run(addr)
}
