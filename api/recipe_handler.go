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

type trendRecipeResponse struct {
	Data []db.FakeListTrendRecipeRow `json:"data"`
}

func (s *Server) ListTrendRecipe(c *gin.Context) {
	const limit int32 = 10
	var response trendRecipeResponse
	var err error
	response.Data, err = s.q.FakeListTrendRecipe(context.Background(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ダミーデータ作成（本番では消す）
	for i := 0; i < len(response.Data); i++ {
		response.Data[i].Name = gimei.NewName().First.Katakana()
		response.Data[i].Introduction = gimei.NewAddress().String() + "。" + gimei.NewAddress().String() + "。" + gimei.NewAddress().String() + "。" + gimei.NewAddress().String() + "。" + gimei.NewAddress().String() + "。"
		response.Data[i].NumFav = rand.Int31n(1000)
		response.Data[i].Score = rand.Int31n(100)
	}

	// レスポンス型バリデーション
	validate := validator.New()
	for i := 0; i < len(response.Data); i++ {
		// ListTrendRecipeRow型からJSONに変換
		jsn, err := json.Marshal(response.Data[i])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// JSONからdocs.TrendRecipe型に変換
		var obj docs.TrendRecipe
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

func (s *Server) CreateRecipe(c *gin.Context) {
	// TODO: リクエストボディのバリデーション

	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	row, err := s.q.CreateRecipe(context.Background(), data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: レスポンスのバリデーション

	c.JSON(http.StatusOK, row)
}
