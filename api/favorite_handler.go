package api

import (
	"context"
	"net/http"
	"reflect"

	"github.com/aopontann/gin-sqlc/docs"
	"github.com/aopontann/gin-sqlc/utils"

	db "github.com/aopontann/gin-sqlc/db/sqlc"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateFavoriteRecipe(c *gin.Context) {
	// パスパラメータ取り出し
	param := db.CreateFavoriteRecipeParams{}
	var err error
	param.RecipeID, err = utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "パスパラメータが取得できませんでした。", "error": err.Error()})
	}

	// Authentication()でセットしたUsrIDを取得
	rv := c.MustGet("rv").(redisValue)
	param.UsrID, err = utils.StrToUUID(rv.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"title": "認証が許可されていません。", "error": err.Error()})
		return
	}

	// 新規登録処理
	row, err := s.q.CreateFavoriteRecipe(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "SQLの処理が上手くいきませんでした。", "error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.Favoring, docs.CreateFavoriteRecipe](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "型のバリデーションが失敗しました。", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}

func (s *Server) DeleteFavoriteRecipe(c *gin.Context) {
	// パスパラメータ取り出し
	param := db.DeleteFavoriteRecipeParams{}
	var err error
	param.RecipeID, err = utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Authentication()でセットしたUsrIDを取得
	rv := c.MustGet("rv").(redisValue)
	param.UsrID, err = utils.StrToUUID(rv.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 削除処理
	row, err := s.q.DeleteFavoriteRecipe(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.Favoring, docs.DeleteFavoriteRecipe](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}

func (s *Server) ExistsFavoriteRecipe(c *gin.Context) {
	type existsFavoriteResponse struct {
		Exists bool `json:"exists"`
	}
	var response existsFavoriteResponse

	// パスパラメータ取り出し
	param := db.ExistsFavoriteRecipeParams{}
	var err error
	param.RecipeID, err = utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Authentication()でセットしたUsrIDを取得
	rv := c.MustGet("rv").(redisValue)
	param.UsrID, err = utils.StrToUUID(rv.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 確認処理
	response.Exists, err = s.q.ExistsFavoriteRecipe(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[existsFavoriteResponse, docs.ExistsFavoriteRecipe](&response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "型のバリデーションが失敗しました。", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) ListFavoriteRecipe(c *gin.Context) {
	type listFavoriteRecipeResponse struct {
		Data []db.ListFavoriteRecipeRow `json:"data"`
	}
	var response listFavoriteRecipeResponse

	// Authentication()でセットしたUsrIDを取得
	rv := c.MustGet("rv").(redisValue)
	usrId, err := utils.StrToUUID(rv.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// リスト化
	response.Data, err = s.q.ListFavoriteRecipe(context.Background(), usrId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.Data == nil || reflect.ValueOf(response.Data).IsNil() {
		response.Data = []db.ListFavoriteRecipeRow{}
	}

	// レスポンス型バリデーション
	// UUID が nil の時はエラーが出る
	err = utils.ValidateStructTwoWay[listFavoriteRecipeResponse, docs.ListFavoriteRecipe](&response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "型のバリデーションが失敗しました。", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
