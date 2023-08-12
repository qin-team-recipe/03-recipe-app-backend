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

	// シェフ関連
	api.GET("/chefs/:id", s.GetChef)
	api.PUT("/chefs/:id", s.UpdateChef)
	api.DELETE("/chefs/:id", s.DeleteChef)
	api.POST("/chefs", s.CreateChef)
	api.GET("/chefs/featured", s.ListFeaturedChef)

	// ユーザー関連
	api.POST("/users", s.CreateUser)
	api.GET("/users/:id", s.GetUser)

	//// 仮で作成　セッションの説明用 ////
	// グループを作成
	usr := api.Group("/user")
	// /api/user/* リンクにアクセスしたとき、登録したハンドラ関数(ここではGetUserId)を実行する前に Authentication() を実行するようにする処理
	usr.Use(s.Authentication())
	// ユーザIDのみ返すAPI
	usr.GET("/id", s.GetUserId)

	usr.GET("/users", s.GetSelf)
	usr.PUT("/users", s.UpdateSelf)
	usr.DELETE("/users", s.DeleteSelf)

	// レシピ関連
	api.GET("/recipes/:id", s.GetRecipe)
	api.PUT("/recipes/:id", s.UpdateRecipe)
	api.DELETE("/recipes/:id", s.DeleteRecipe)
	api.POST("/recipes/chef", s.CreateChefRecipe)
	usr.POST("/recipes/user", s.CreateUsrRecipe)
	api.GET("/recipes/trend", s.ListTrendRecipe)
	//api.GET("/recipes/chef/:chef_id", s.)
	//api.GET("/recipes/user/:usr_id", s.)

	// ショッピングリスト関連
	usr.GET("/lists", s.ListShoppingList)
	usr.GET("/lists/:recipe_id", s.GetShoppingList)
	usr.PUT("/lists/:id", s.UpdateShoppingList)
	usr.DELETE("/lists/:id", s.DeleteShoppingList)
	usr.POST("/lists", s.CreateShoppingList)
}

func (s *Server) Start(addr string) error {
	return s.r.Run(addr)
}
