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
	api := s.r.Group("/api")
	auth := s.r.Group("/auth")
	auth.GET("/google/login", s.OauthGoogleLogin)
	auth.GET("/google/callback", s.OauthGoogleCallback)
	// api.POST("/users", s.RegisterUser)
	// api.POST("/users/login", s.LoginUser)

	user := api.Group("/author")
	// user.Use(AuthMiddleware())
	user.GET("", s.ListAuthors)
}

func (s *Server) Start(addr string) error {
	return s.r.Run(addr)
}
