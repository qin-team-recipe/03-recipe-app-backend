package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aopontann/gin-sqlc/api"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

	server := api.NewServer(conn, rdb)

	server.MountHandlers()

	err = server.Start(":8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start up API server: %v\n", err)
		os.Exit(1)
	}
}
