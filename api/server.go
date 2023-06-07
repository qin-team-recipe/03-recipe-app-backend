package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	db "github.com/aopontann/gin-sqlc/db/sqlc"
)

type Server struct {
	r  *gin.Engine
	db *pgx.Conn
	q  *db.Queries
}

func NewServer(conn *pgx.Conn) *Server {
	engine := gin.Default()
	server := &Server{
		r:  engine,
		db: conn,
		q:  db.New(conn),
	}
	return server
}

func (s *Server) MountHandlers() {
	api := s.r.Group("/api")
	// api.POST("/users", s.RegisterUser)
	// api.POST("/users/login", s.LoginUser)

	user := api.Group("/author")
	// user.Use(AuthMiddleware())
	user.GET("", s.ListAuthors)
}

func (s *Server) Start(addr string) error {
	return s.r.Run(addr)
}
