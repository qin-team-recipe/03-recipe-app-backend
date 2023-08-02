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
	api.GET("/featuredChef", s.ListFeaturedChef)

	// ユーザー関連
	api.POST("/createUser", s.CreateUser)

	//// 仮で作成　セッションの説明用 ////
	// グループを作成
	usr := api.Group("/user")
	// /api/user/* リンクにアクセスしたとき、登録したハンドラ関数(ここではGetUserId)を実行する前に Authentication() を実行するようにする処理
	usr.Use(s.Authentication())
	// ユーザIDのみ返すAPI
	usr.GET("/id", s.GetUserId)

	// レシピ関連
	api.GET("/trendRecipe", s.ListTrendRecipe)
	api.POST("/createChefRecipe", s.CreateChefRecipe)
	api.POST("/createUsrRecipe", s.CreateUsrRecipe)
	api.PUT("/updateRecipe/:id", s.UpdateRecipe)
}

func (s *Server) Start(addr string) error {
	return s.r.Run(addr)
}
