package api_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aopontann/gin-sqlc/api"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestGetChef(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	// プログラムが終了すると、開いていたコネクションはクローズされる
	// defer conn.Close(context.Background())

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	server := api.NewServer(conn, rdb)
	server.MountHandlers()
	r := httptest.NewRequest(http.MethodGet, "/api/chefs/c73fe8ae-22e8-45b6-a257-8595db1b951d", nil)
	w := httptest.NewRecorder()
	server.R.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)
}
