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
	api.GET("/chef/:id", s.GetChef)
	api.GET("/featured-chefs", s.ListFeaturedChef)
	api.POST("/chef", s.CreateChef)
	api.PUT("/chef/:id", s.UpdateChef)

	// ユーザー関連
	api.GET("/getUser/:id", s.GetUser)
	api.POST("/createUser", s.CreateUser)
	api.PUT("/updateUser/:id", s.UpdateUser)

	//// 仮で作成　セッションの説明用 ////
	// グループを作成
	usr := api.Group("/user")
	// /api/user/* リンクにアクセスしたとき、登録したハンドラ関数(ここではGetUserId)を実行する前に Authentication() を実行するようにする処理
	usr.Use(s.Authentication())
	// ユーザIDのみ返すAPI
	usr.GET("/id", s.GetUserId)

	// レシピ関連
	api.GET("/trend-recipes", s.ListTrendRecipe)
	api.GET("/recipe/:id", s.GetRecipe)
	//api.GET("/chef-recipes/:chef_id", s.)
	//api.GET("/user-recipes/:usr_id", s.)
	api.POST("/chef-recipe", s.CreateChefRecipe)
	api.POST("/user-recipe", s.CreateUsrRecipe)
	api.PUT("/recipe/:id", s.UpdateRecipe)
}

func (s *Server) Start(addr string) error {
	return s.r.Run(addr)
}
