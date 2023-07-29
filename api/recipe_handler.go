package api

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/aopontann/gin-sqlc/common"
	db "github.com/aopontann/gin-sqlc/db/sqlc"
	"github.com/aopontann/gin-sqlc/docs"

	"github.com/gin-gonic/gin"
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
	err = common.ValidateStructTwoWay[trendRecipeResponse, docs.TrendRecipe](&response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) CreateChefRecipe(c *gin.Context) {
	// リクエストボディーをJSONとして取り出し
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// JSONのバリデーション
	err = common.ValidateStruct[docs.PostApiCreateChefRecipeJSONRequestBody](data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 新規登録
	row, err := s.q.CreateRecipe(context.Background(), data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = common.ValidateStructTwoWay[db.VRecipe, docs.CreateChefRecipe](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}

func (s *Server) CreateUsrRecipe(c *gin.Context) {
	// リクエストボディーをJSONとして取り出し
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// JSONのバリデーション
	err = common.ValidateStruct[docs.PostApiCreateUsrRecipeJSONRequestBody](data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 新規登録
	row, err := s.q.CreateRecipe(context.Background(), data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = common.ValidateStructTwoWay[db.VRecipe, docs.CreateUsrRecipe](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}
