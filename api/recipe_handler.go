package api

import (
	"context"
	"net/http"
	"reflect"

	db "github.com/aopontann/gin-sqlc/db/sqlc"
	"github.com/aopontann/gin-sqlc/docs"
	"github.com/aopontann/gin-sqlc/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) ListTrendRecipe(c *gin.Context) {
	type trendRecipeResponse struct {
		Data []db.ListTrendRecipeRow `json:"data"`
	}

	const limit int32 = 10
	var response trendRecipeResponse
	var err error
	response.Data, err = s.q.ListTrendRecipe(context.Background(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "SQLの実行中に失敗しました。", "error": err.Error()})
		return
	}

	if response.Data == nil || reflect.ValueOf(response.Data).IsNil() {
		response.Data = []db.ListTrendRecipeRow{}
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[trendRecipeResponse, docs.TrendRecipe](&response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "型のバリデーションが失敗しました。", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) GetRecipe(c *gin.Context) {
	// パスパラメータ取り出し
	id, err := utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "パスパラメータの取り出しに失敗しました。", "error": err.Error()})
	}

	// 問い合わせ処理
	row, err := s.q.GetRecipe(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "SQLの実行中に失敗しました。", "error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.VRecipe, docs.GetRecipe](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "型のバリデーションが失敗しました。", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}

func (s *Server) ListRecipe(c *gin.Context) {
	type listRecipeponse struct {
		Data []db.ListRecipeRow `json:"data"`
	}
	var response listRecipeponse

	var err error
	response.Data, err = s.q.ListRecipe(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "SQLの実行中に失敗しました。", "error": err.Error()})
	}

	if response.Data == nil || reflect.ValueOf(response.Data).IsNil() {
		response.Data = []db.ListRecipeRow{}
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[listRecipeponse, docs.ListRecipe](&response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "型のバリデーションが失敗しました。", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
