package api

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"

	db "github.com/aopontann/gin-sqlc/db/sqlc"
	"github.com/aopontann/gin-sqlc/docs"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mattn/go-gimei"
)

type featuredChefResponse struct {
	Data []db.FakeListFeaturedChefRow `json:"data"`
}

func (s *Server) ListFeaturedChef(c *gin.Context) {
	const limit int32 = 10
	var response featuredChefResponse
	var err error
	response.Data, err = s.q.FakeListFeaturedChef(context.Background(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ダミーデータ作成（本番では消す）
	for i := 0; i < len(response.Data); i++ {
		response.Data[i].Name = gimei.NewName().String()
		response.Data[i].NumFollower = rand.Int31n(1000)
		response.Data[i].Score = rand.Int31n(100)
	}

	// レスポンス型バリデーション
	validate := validator.New()
	for i := 0; i < len(response.Data); i++ {
		// ListFeaturedChefRow型からJSONに変換
		jsn, err := json.Marshal(response.Data[i])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// JSONからdocs.FeaturedChef型に変換
		var obj docs.FeaturedChef
		if err := json.Unmarshal(jsn, &obj); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// バリデーション
		err = validate.Struct(obj)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, response)
}
