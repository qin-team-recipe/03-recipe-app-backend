package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	db "github.com/aopontann/gin-sqlc/db/sqlc"
	"github.com/aopontann/gin-sqlc/docs"
	"github.com/aopontann/gin-sqlc/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

func (s *Server) ListShoppingList(c *gin.Context) {
	type shoppingListResponse struct {
		Data []db.ListShoppingListRow `json:"data"`
	}

	// Authentication()でセットしたUsrIDを取得
	rv := c.MustGet("rv").(redisValue)
	usrId, err := utils.StrToUUID(rv.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var response shoppingListResponse
	response.Data, err = s.q.ListShoppingList(context.Background(), usrId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.Data == nil || reflect.ValueOf(response.Data).IsNil() {
		response.Data = []db.ListShoppingListRow{}
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[shoppingListResponse, docs.GetShoppingLists](&response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) GetShoppingList(c *gin.Context) {
	var param db.GetShoppingListParams
	var err error

	// Authentication()でセットしたUsrIDを取得
	rv := c.MustGet("rv").(redisValue)
	param.UsrID, err = utils.StrToUUID(rv.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// パスパラメータ取り出し
	param.RecipeID, err = utils.StrToUUID(c.Param("recipe_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// 問い合わせ処理
	row, err := s.q.GetShoppingList(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.GetShoppingListRow, docs.GetShoppingList](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}

func (s *Server) CreateShoppingList(c *gin.Context) {
	// リクエストボディをJSONに変換
	jsn, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// リクエストボディJSONのバリデーション
	err = utils.ValidateStruct[docs.PostApiUserListsJSONRequestBody](jsn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// CreateShoppingListParams構造体に[]InnerCreateShoppingItemParamsを追加
	type Alias1 = db.CreateShoppingListParams
	type requestBody struct {
		*Alias1
		Item []db.InnerCreateShoppingItemParams `json:"item"`
	}
	reqb := requestBody{}

	// CreateShoppingListRow構造体に[]InnerCreateShoppingItemRowを追加
	type Alias2 = db.ShoppingList
	type response struct {
		*Alias2
		Item []db.InnerCreateShoppingItemRow `json:"item"`
	}
	resp := response{}
	resp.Alias2 = &Alias2{}
	resp.Item = []db.InnerCreateShoppingItemRow{}

	// バリデーション済のJSONをreqb変数に変換
	decoder := json.NewDecoder(bytes.NewReader(jsn))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&reqb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Authentication()でセットしたUsrIDを取得
	rv := c.MustGet("rv").(redisValue)
	reqb.UsrID, err = utils.StrToUUID(rv.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// トランザクション開始
	tx, err := s.db.Begin(c)
	if err != nil {
		return
	}
	defer tx.Rollback(c)
	qtx := s.q.WithTx(tx)

	// 買い物リストテーブルへの新規登録
	*resp.Alias2, err = qtx.CreateShoppingList(context.Background(), *reqb.Alias1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 買い物明細テーブルへの新規登録
	for i := 0; i < len(reqb.Item); i++ {
		reqb.Item[i].ShoppingListID = resp.ID
		reqb.Item[i].Idx = int32(i + 1)
		itemRow, err := qtx.InnerCreateShoppingItem(context.Background(), reqb.Item[i])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		resp.Item = append(resp.Item, itemRow)
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[response, docs.CreateShoppingList](&resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// トランザクション終了
	err = tx.Commit(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Server) UpdateShoppingList(c *gin.Context) {
	// リクエストボディをJSONに変換
	jsn, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// リクエストボディJSONのバリデーション
	err = utils.ValidateStruct[docs.PutApiUserListsIdJSONRequestBody](jsn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// UpdateShoppingListParams構造体に[]InnerUpdateShoppingItemParamsを追加
	type Alias1 = db.UpdateShoppingListParams
	type requestBody struct {
		*Alias1
		Item []db.InnerUpdateShoppingItemParams `json:"item"`
	}
	reqb := requestBody{}

	// UpdateShoppingListRow構造体に[]InnerUpdateShoppingItemRowを追加
	type Alias2 = db.ShoppingList
	type response struct {
		*Alias2
		Item []db.InnerUpdateShoppingItemRow `json:"item"`
	}
	resp := response{}
	resp.Alias2 = &Alias2{}
	resp.Item = []db.InnerUpdateShoppingItemRow{}

	// バリデーション済のJSONをreqb変数に変換
	decoder := json.NewDecoder(bytes.NewReader(jsn))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&reqb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// パスパラメータ取り出し
	reqb.ID, err = utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Authentication()でセットしたUsrIDを取得
	rv := c.MustGet("rv").(redisValue)
	reqb.UsrID, err = utils.StrToUUID(rv.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 削除対象外の買い物明細IDを抽出
	var deleteParam db.InnerDeleteNotAnyShoppingItemParams
	deleteParam.ShoppingListID = reqb.ID
	for i := 0; i < len(reqb.Item); i++ {
		if reqb.Item[i].ID.Valid {
			deleteParam.ID = append(deleteParam.ID, reqb.Item[i].ID)
		}
	}

	// トランザクション開始
	tx, err := s.db.Begin(c)
	if err != nil {
		return
	}
	defer tx.Rollback(c)
	qtx := s.q.WithTx(tx)

	// 削除対象の買い物明細を削除
	err = qtx.InnerDeleteNotAnyShoppingItem(context.Background(), deleteParam)
	if err != nil {
		return
	}

	// 買い物リストテーブルの更新
	*resp.Alias2, err = qtx.UpdateShoppingList(context.Background(), *reqb.Alias1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := 0; i < len(reqb.Item); i++ {
		reqb.Item[i].ShoppingListID = resp.ID
		reqb.Item[i].Idx = int32(i + 1)
		if reqb.Item[i].ID.Valid {
			// 買い物明細テーブルの更新
			itemRow, err := qtx.InnerUpdateShoppingItem(context.Background(), reqb.Item[i])
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			resp.Item = append(resp.Item, itemRow)
		} else {
			// 買い物明細テーブルへの新規登録
			var param db.InnerCreateShoppingItemParams
			if err := copier.Copy(&param, &reqb.Item[i]); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			itemRow, err := qtx.InnerCreateShoppingItem(context.Background(), param)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			var itemRow2 db.InnerUpdateShoppingItemRow
			if err := copier.Copy(&itemRow2, &itemRow); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			resp.Item = append(resp.Item, itemRow2)
		}
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[response, docs.UpdateShoppingList](&resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// トランザクション終了
	err = tx.Commit(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Server) DeleteShoppingList(c *gin.Context) {
	var param db.DeleteShoppingListParams
	var err error

	// Authentication()でセットしたUsrIDを取得
	rv := c.MustGet("rv").(redisValue)
	param.UsrID, err = utils.StrToUUID(rv.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// パスパラメータ取り出し
	param.ID, err = utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// 問い合わせ処理
	row, err := s.q.DeleteShoppingList(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.ShoppingList, docs.DeletedShoppingList](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}
