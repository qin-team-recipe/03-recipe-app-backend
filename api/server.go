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
	usr := api.Group("/users")
	// /api/user/* リンクにアクセスしたとき、登録したハンドラ関数(ここではGetUserId)を実行する前に Authentication() を実行するようにする処理
	usr.Use(s.Authentication())
	// ユーザIDのみ返すAPI(検証用)
	usr.GET("/id", s.GetUserId)

	usr.GET("/", s.GetSelf)       // 自分を取得するAPI
	usr.PUT("/", s.UpdateSelf)    // 自分を更新するAPI
	usr.DELETE("/", s.DeleteSelf) // 自分を削除するAPI

	// ユーザー（= 一般シェフ）のレシピ関連
	//usr.GET("/users/recipe", s.GetUsrRecipe)
	usr.PUT("/recipe/:recipe_id", s.UpdateUserRecipe)    // 一般シェフのマイレシピを更新するAPI
	usr.DELETE("/recipe/:recipe_id", s.DeleteUserRecipe) // 一般シェフのマイレシピを削除するAPI
	usr.POST("/recipe", s.CreateUsrRecipe)               // 一般シェフのマイレシピを新規登録するAPI

	// 有名シェフのレシピ関連
	//api.GET("/chefs/:id/recipe", s.GetChefRecipe)
	api.PUT("/chefs/recipe/:recipe_id", s.UpdateChefRecipe)    // 有名シェフのレシピを更新するAPI
	api.DELETE("/chefs/recipe/:recipe_id", s.DeleteChefRecipe) // 有名シェフのレシピを削除するAPI
	api.POST("/chefs/:id/recipe", s.CreateChefRecipe)          // 有名シェフのレシピを新規登録するAPI
	api.GET("/chefs/:id/recipes", s.ListChefRecipe)            // 有名シェフのレシピ一覧を取得するAPI
	api.GET("/chefs/recipe/search", s.SearchChefRecipe)        // 有名シェフのレシピを全文検索するAPI

	// 有名シェフ関連
	api.GET("/chefs/:id", s.GetChef)               // 有名シェフを取得するAPI
	api.PUT("/chefs/:id", s.UpdateChef)            // 有名シェフを更新するAPI
	api.DELETE("/chefs/:id", s.DeleteChef)         // 有名シェフを削除するAPI
	api.GET("/chefs", s.ListChef)                  // 有名シェフを取得するAPI
	api.POST("/chefs", s.CreateChef)               // 有名シェフを新規登録するAPI
	api.GET("/chefs/featured", s.ListFeaturedChef) // 注目の有名シェフ一覧を取得するAPI
	api.GET("/chefs/search", s.SearchChef)         // 有名シェフを全文検索するAPI

	// レシピ関連
	api.GET("/recipes/:id", s.GetRecipe)         // レシピを取得するAPI
	api.GET("/recipes", s.ListRecipe)            // レシピ一覧を取得するAPI
	api.GET("/recipes/trend", s.ListTrendRecipe) // 話題のレシピ一覧を取得するAPI

	// 買い物リスト関連
	lists := api.Group("/lists")
	lists.Use(s.Authentication())
	lists.GET("/", s.ListShoppingList)                 // ユーザーの買い物リスト一覧を取得するAPI
	lists.GET("/recipe/:recipe_id", s.GetShoppingList) // 買い物リストを取得するAPI
	lists.PUT("/:id", s.UpdateShoppingList)            // 買い物リストを更新するAPI
	lists.DELETE("/:id", s.DeleteShoppingList)         // 買い物リストを削除するAPI
	lists.POST("/", s.CreateShoppingList)              // 買い物リストを新規登録するAPI

	// 画像関連
	api.GET("/images", s.GetImage)   // 画像を取得するAPI（webp限定）
	api.POST("/images", s.PostImage) // 画像を新規登録するAPI（webp形式に変換される）

	// フォローグループを作成
	follow := api.Group("/follow")
	follow.Use(s.Authentication())

	// 有名シェフフォロー関連
	follow.POST("/chefs/:id", s.CreateFollowChef)           // 有名シェフをフォローするAPI
	follow.DELETE("/chefs/:id", s.DeleteFollowChef)         // 有名シェフのフォローを解除するAPI
	follow.GET("/chefs/:id", s.ExistsFollowChef)            // 有名シェフをフォローしているか
	follow.GET("/chefs", s.ListFollowChef)                  // フォローしている有名シェフの一覧を取得するAPI
	follow.GET("/chefs/recipes", s.ListFollowChefNewRecipe) // フォローしているシェフの新着レシピ一覧を取得するAPI

	// 一般シェフフォロー関連
	follow.POST("/users/:id", s.CreateFollowUser)   // 一般シェフをフォローするAPI
	follow.DELETE("/users/:id", s.DeleteFollowUser) // 一般シェフのフォローを解除するAPI
	follow.GET("/users/:id", s.ExistsFollowUser)    // 一般シェフをフォローしているか
	follow.GET("/users", s.ListFollowUser)          // フォローしている一般シェフの一覧を取得するAPI

	// フォローグループを作成
	fav := api.Group("/favorite")
	fav.Use(s.Authentication())

	// お気に入りレシピ関連
	fav.POST("/recipes/:id", s.CreateFavoriteRecipe)   // お気に入りレシピ登録API
	fav.DELETE("/recipes/:id", s.DeleteFavoriteRecipe) // お気に入りレシピ解除API
	fav.GET("/recipes/:id", s.ExistsFavoriteRecipe)    // お気に入りレシピとして登録しているか確認API
	fav.GET("/recipes", s.ListFavoriteRecipe)          // お気に入りレシピの一覧を取得するAPI

	// 管理者グループを作成
	// admin := api.Group("/admin")
	// admin.Use() TODO: 管理者用のミドルウェアを作成する
	// admin.POST("/chefs", )
	// admin.PUT("/chefs", )
	// admin.DELETE("/chefs", )

	// admin.GET("/chefs/recipes", )
	// admin.POST("/chefs/recipes", )
	// admin.PUT("/chefs/recipes", )
	// admin.DELETE("/chefs/recipes", )

	// 開発者用
	dev := api.Group("/dev")
	dev.POST("/chefs", s.CreateChefData) // 有名シェフ＆レシピのテストデータを作成するAPI
	dev.POST("/users", s.CreateUserData) // 自分でないユーザー（一般シェフ）＆レシピのテストデータを作成するAPI
}

func (s *Server) Start(addr string) error {
	return s.r.Run(addr)
}
