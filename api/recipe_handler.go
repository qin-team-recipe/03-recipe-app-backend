package api

import (
	"context"
	"encoding/json"
	"github.com/aopontann/gin-sqlc/common"
	db "github.com/aopontann/gin-sqlc/db/sqlc"
	"github.com/aopontann/gin-sqlc/docs"
	"math/rand"
	"net/http"

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
	// リクエストボディを構造体にバインド
	reqb := docs.PostApiCreateChefRecipeJSONRequestBody{}
	if err := c.ShouldBind(&reqb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 構造体からJSONに変換
	jsn, err := json.Marshal(&reqb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// 新規登録
	row, err := s.q.CreateRecipe(context.Background(), jsn)
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
	// リクエストボディを構造体にバインド
	reqb := docs.PostApiCreateUsrRecipeJSONRequestBody{}
	if err := c.ShouldBind(&reqb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 構造体からJSONに変換
	jsn, err := json.Marshal(&reqb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// 新規登録
	row, err := s.q.CreateRecipe(context.Background(), jsn)
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

func (s *Server) UpdateRecipe(c *gin.Context) {
	var param db.UpdateRecipeParams
	var err error
	param.ID, err = common.StrToUUID(c.Param("id"))

	//// リクエストボディを構造体にバインド
	//reqb := docs.PostApiUpdateRecipeJSONRequestBody{}
	//if err := c.ShouldBind(&reqb); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//
	//// 構造体からJSONに変換
	//jsn, err := json.Marshal(&reqb)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//}

	param.Data, err = c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新
	row, err := s.q.UpdateRecipe(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//// レスポンス型バリデーション
	//err = common.ValidateStructTwoWay[db.VRecipe, docs.UpdateRecipe](&row)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	c.JSON(http.StatusOK, row)
}