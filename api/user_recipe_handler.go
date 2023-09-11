package api

import (
	"context"
	"encoding/json"
	"net/http"

	db "github.com/aopontann/gin-sqlc/db/sqlc"
	"github.com/aopontann/gin-sqlc/docs"
	"github.com/aopontann/gin-sqlc/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Server) CreateUsrRecipe(c *gin.Context) {
	// リクエストボディを構造体にバインド
	reqb := docs.PostApiUserUsersRecipeJSONRequestBody{}
	if err := c.ShouldBind(&reqb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "リクエストボディを構造体にバインドが失敗しました。", "error": err.Error()})
		return
	}

	// Authentication()でセットしたUsrIDを取得
	rv := c.MustGet("rv").(redisValue)
	usrId, err := utils.StrToUUID(rv.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"title": "認証が失敗しました。", "error": err.Error()})
		return
	}

	// 構造体にuserIdを追加してJSONに変換
	type Alias docs.PostApiUserUsersRecipeJSONRequestBody
	jsn, err := json.Marshal(&struct {
		UsrId pgtype.UUID `json:"usrId"`
		*Alias
	}{
		UsrId: usrId,
		Alias: (*Alias)(&reqb),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "構造体をJSONに変換するのが失敗しました。", "error": err.Error()})
	}

	// 新規登録処理
	row, err := s.q.CreateRecipe(context.Background(), jsn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "SQLの処理に失敗しました。", "error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.VRecipe, docs.CreateUsrRecipe](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "型のバリデーションが失敗しました。", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, row)
}

func (s *Server) UpdateUserRecipe(c *gin.Context) {
	var param db.UpdateRecipeParams
	var err error

	// パスパラメータ取り出し
	param.ID, err = utils.StrToUUID(c.Param("recipe_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "パスパラメータの取得に失敗しました。", "error": err.Error()})
	}

	// リクエストボディを構造体にバインド
	reqb := docs.PutApiUserUsersRecipeRecipeIdJSONRequestBody{}
	if err := c.ShouldBind(&reqb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "リクエストボディを構造体に変換するのが失敗しました。", "error": err.Error()})
		return
	}

	// Authentication()でセットしたUsrIDを取得
	rv := c.MustGet("rv").(redisValue)
	usrId, err := utils.StrToUUID(rv.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"title": "認証が失敗しました。", "error": err.Error()})
		return
	}

	// 構造体にuserIdを追加してJSONに変換
	type Alias docs.PutApiUserUsersRecipeRecipeIdJSONRequestBody
	param.Data, err = json.Marshal(&struct {
		UsrId pgtype.UUID `json:"usrId"`
		*Alias
	}{
		UsrId: usrId,
		Alias: (*Alias)(&reqb),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "構造体をJSONに変換するのが失敗しました。", "error": err.Error()})
	}

	// 更新処理
	row, err := s.q.UpdateRecipe(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "SQLの処理に失敗しました。", "error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.VRecipe, docs.UpdateRecipe](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "型のバリデーションが失敗しました。", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}

func (s *Server) DeleteUserRecipe(c *gin.Context) {
	var param db.DeleteUserRecipeParams
	var err error

	// パスパラメータ取り出し
	param.ID, err = utils.StrToUUID(c.Param("recipe_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "パスパラメータの取得に失敗しました。", "error": err.Error()})
	}

	// Authentication()でセットしたUsrIDを取得
	rv := c.MustGet("rv").(redisValue)
	param.UsrID, err = utils.StrToUUID(rv.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"title": "認証が失敗しました。", "error": err.Error()})
		return
	}

	// 問い合わせ処理
	row, err := s.q.DeleteUserRecipe(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "SQLの処理に失敗しました。", "error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.DeleteUserRecipeRow, docs.DeletedRecipe](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "型のバリデーションが失敗しました。", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}
