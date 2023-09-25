package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aopontann/gin-sqlc/api"
	db "github.com/aopontann/gin-sqlc/db/sqlc"
	"github.com/aopontann/gin-sqlc/docs"
	"github.com/gin-gonic/gin"
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
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	server := api.NewServer(conn, rdb)
	query := db.New(conn)

	reqb := docs.PostApiChefsJSONRequestBody{
		Name: "test1-chef",
	}

	// 構造体からJSONに変換
	jsn, err := json.Marshal(&reqb)
	if err != nil {
		t.Error(err)
	}

	// 新規登録処理
	row, err := query.CreateChef(context.Background(), jsn)
	if err != nil {
		t.Error(err)
		return
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: string(row.ID.Bytes[:])}}

	server.GetChef(c)
	assert.Equal(t, 200, w.Code)

	_, err = query.DeleteChef(context.Background(), row.ID)
	if err != nil {
		t.Error(err)
	}
}
