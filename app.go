package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aopontann/gin-sqlc/api"

	"github.com/jackc/pgx/v5"
)

func main() {

	// DATABASE_URL := "postgres://username:password@localhost:5432/database_name"
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	server := api.NewServer(conn)

	server.MountHandlers()

	err = server.Start(":8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start up API server: %v\n", err)
		os.Exit(1)
	}
}
