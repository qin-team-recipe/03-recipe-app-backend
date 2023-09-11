package api

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	db "github.com/aopontann/gin-sqlc/db/sqlc"
	"github.com/aopontann/gin-sqlc/docs"
	"github.com/aopontann/gin-sqlc/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Server) CreateChefRecipe(c *gin.Context) {
	// リクエストボディを構造体にバインド
	reqb := docs.PostApiChefsChefIdRecipeJSONRequestBody{}
	if err := c.ShouldBind(&reqb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "リクエストボディを構造体にバインドできませんでした。", "error": err.Error()})
		return
	}

	// パスパラメータ取り出し
	chefId, err := utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "パスパラメータに異常があります。", "error": err.Error()})
	}

	// 構造体にchefIdを追加してJSONに変換
	type Alias docs.PostApiChefsChefIdRecipeJSONRequestBody
	jsn, err := json.Marshal(&struct {
		ChefId pgtype.UUID `json:"chefId"`
		*Alias
	}{
		ChefId: chefId,
		Alias:  (*Alias)(&reqb),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "構造体をJSONに変更できませんでした。", "error": err.Error()})
	}

	// 新規登録処理
	row, err := s.q.CreateRecipe(context.Background(), jsn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "SQLの処理が上手くいきませんでした。", "error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.VRecipe, docs.CreateChefRecipe](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "型のバリデーションが失敗しました。", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, row)
}

func (s *Server) UpdateChefRecipe(c *gin.Context) {
	var param db.UpdateRecipeParams
	var err error

	// パスパラメータ取り出し
	param.ID, err = utils.StrToUUID(c.Param("recipe_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "パスパラメータに異常があります。", "error": err.Error()})
	}

	// リクエストボディを構造体にバインド
	reqb := docs.PutApiChefsRecipeRecipeIdJSONRequestBody{}
	if err := c.ShouldBind(&reqb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "リクエストボディを構造体にバインドできませんでした。", "error": err.Error()})
		return
	}

	// 構造体からJSONに変換
	param.Data, err = json.Marshal(&reqb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "構造体からJSONに変換できませんでした。", "error": err.Error()})
	}

	// 更新処理
	row, err := s.q.UpdateRecipe(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "SQLの処理が上手くいきませんでした。", "error": err.Error()})
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

func (s *Server) DeleteChefRecipe(c *gin.Context) {
	// パスパラメータ取り出し
	id, err := utils.StrToUUID(c.Param("recipe_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "パスパラメータに異常があります。", "error": err.Error()})
	}

	// 問い合わせ処理
	row, err := s.q.DeleteChefRecipe(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "SQLの処理が上手くいきませんでした。", "error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.DeleteChefRecipeRow, docs.DeletedRecipe](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "型のバリデーションが失敗しました。", "error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, row)
}

func (s *Server) SearchChefRecipe(c *gin.Context) {
	type searchChefRecipeResponse struct {
		Data []db.SearchChefRecipeRow `json:"data"`
	}

	// クエリパラメータ取り出し
	query := c.Query("q")

	// 全文検索
	var response searchChefRecipeResponse
	var err error
	response.Data, err = s.q.SearchChefRecipe(context.Background(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "SQLの処理が上手くいきませんでした。", "error": err.Error()})
		return
	}

	if response.Data == nil || reflect.ValueOf(response.Data).IsNil() {
		response.Data = []db.SearchChefRecipeRow{}
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[searchChefRecipeResponse, docs.SearchChefRecipe](&response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "型のバリデーションが失敗しました。", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) ListChefRecipe(c *gin.Context) {
	// パスパラメータ取り出し
	id, err := utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// クエリパラメータ取り出し
	order := c.Query("order")

	if order == "new" {
		type listChefRecipeponse struct {
			Data []db.ListChefRecipeNewRow `json:"data"`
		}
		var response listChefRecipeponse

		response.Data, err = s.q.ListChefRecipeNew(context.Background(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"title": "SQLの処理が上手くいきませんでした。", "error": err.Error()})
		}

		if response.Data == nil || reflect.ValueOf(response.Data).IsNil() {
			response.Data = []db.ListChefRecipeNewRow{}
		}

		// レスポンス型バリデーション
		err = utils.ValidateStructTwoWay[listChefRecipeponse, docs.ListRecipe](&response)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"title": "型のバリデーションが失敗しました。", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	} else if order == "fav" {
		type listChefRecipeponse struct {
			Data []db.ListChefRecipeFavRow `json:"data"`
		}
		var response listChefRecipeponse

		response.Data, err = s.q.ListChefRecipeFav(context.Background(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"title": "SQLの処理が上手くいきませんでした。", "error": err.Error()})
		}

		if response.Data == nil || reflect.ValueOf(response.Data).IsNil() {
			response.Data = []db.ListChefRecipeFavRow{}
		}

		// レスポンス型バリデーション
		err = utils.ValidateStructTwoWay[listChefRecipeponse, docs.ListRecipe](&response)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"title": "型のバリデーションが失敗しました。", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"title": "orderがnew, fav以外のエラーが起きました。", "error": "Bad Request"})
		return
	}
}
