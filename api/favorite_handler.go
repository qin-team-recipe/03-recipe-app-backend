package api

import (
	"context"
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// usrIdを取得
	email := c.MustGet("email").(string)
	param.UsrID, err = s.q.GetUserId(context.Background(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 新規登録処理
	row, err := s.q.CreateFavoriteRecipe(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.Favoring, docs.CreateFavoriteRecipe](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	// usrIdを取得
	email := c.MustGet("email").(string)
	param.UsrID, err = s.q.GetUserId(context.Background(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
