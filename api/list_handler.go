package api

import (
	"context"
	"net/http"
	"reflect"

	db "github.com/aopontann/gin-sqlc/db/sqlc"
	"github.com/aopontann/gin-sqlc/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type shoppingListResponse struct {
	Data []db.ListShoppingListRow `json:"data"`
}

func (s *Server) ListShoppingList(c *gin.Context) {
	// usrIdを取得
	email := c.MustGet("email").(string)
	usrId, err := s.q.GetUserId(context.Background(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	//// レスポンス型バリデーション
	//err = utils.ValidateStructTwoWay[shoppingListResponse, docs.TrendRecipe](&response)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	c.JSON(http.StatusOK, response)
}

func (s *Server) GetShoppingList(c *gin.Context) {
	var param db.GetShoppingListParams
	var err error

	// usrIdを取得
	email := c.MustGet("email").(string)
	param.UsrID, err = s.q.GetUserId(context.Background(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	//// レスポンス型バリデーション
	//err = utils.ValidateStructTwoWay[db.GetShoppingListRow, docs.GetUsr](&row)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	c.JSON(http.StatusOK, row)
}

func (s *Server) CreateShoppingList(c *gin.Context) {
	// CreateShoppingListParams構造体に[]CreateShoppingItemParamsを追加
	type Alias1 = db.CreateShoppingListParams
	type requestBody struct {
		*Alias1
		Item []db.CreateShoppingItemParams `json:"item"`
	}
	reqb := requestBody{}

	// CreateShoppingListRow構造体に[]CreateShoppingItemRowを追加
	type Alias2 = db.CreateShoppingListRow
	type response struct {
		*Alias2
		Item []db.CreateShoppingItemRow `json:"item"`
	}
	resp := response{}
	resp.Alias2 = &Alias2{}
	resp.Item = []db.CreateShoppingItemRow{}

	// リクエストボディを構造体にバインド
	if err := c.ShouldBind(&reqb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// usrIdを取得して設定
	var err error
	email := c.MustGet("email").(string)
	reqb.UsrID, err = s.q.GetUserId(context.Background(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		reqb.Item[i].ShoppingListID = resp.Alias2.ID
		reqb.Item[i].Idx = int32(i + 1)
		itemRow, err := qtx.CreateShoppingItem(context.Background(), reqb.Item[i])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		resp.Item = append(resp.Item, itemRow)
	}

	//// レスポンス型バリデーション
	//err = utils.ValidateStructTwoWay[db.GetShoppingListRow, docs.GetUsr](&row)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	// トランザクション終了
	tx.Commit(c)

	c.JSON(http.StatusOK, resp)
}

func (s *Server) UpdateShoppingList(c *gin.Context) {
	// パスパラメータ取り出し
	id, err := utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// UpdateShoppingListParams構造体に[]UpdateShoppingItemParamsを追加
	type Alias1 = db.UpdateShoppingListParams
	type requestBody struct {
		*Alias1
		Item []db.UpdateShoppingItemParams `json:"item"`
	}
	reqb := requestBody{}

	// UpdateShoppingListRow構造体に[]UpdateShoppingItemRowを追加
	type Alias2 = db.UpdateShoppingListRow
	type response struct {
		*Alias2
		Item []db.UpdateShoppingItemRow `json:"item"`
	}
	resp := response{}
	resp.Alias2 = &Alias2{}
	resp.Item = []db.UpdateShoppingItemRow{}

	// リクエストボディを構造体にバインド
	if err := c.ShouldBind(&reqb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reqb.ID = id

	//// usrIdを取得して設定
	//email := c.MustGet("email").(string)
	//reqb.UsrID, err = s.q.GetUserId(context.Background(), email)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	// 削除対象外の買い物明細IDを抽出
	var deleteParam db.DeleteNotAnyShoppingItemParams
	deleteParam.ShoppingListID = id
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
	err = qtx.DeleteNotAnyShoppingItem(context.Background(), deleteParam)
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
		reqb.Item[i].ShoppingListID = resp.Alias2.ID
		reqb.Item[i].Idx = int32(i + 1)
		if reqb.Item[i].ID.Valid {
			// 買い物明細テーブルの更新
			itemRow, err := qtx.UpdateShoppingItem(context.Background(), reqb.Item[i])
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			resp.Item = append(resp.Item, itemRow)
		} else {
			// 買い物明細テーブルへの新規登録
			var param db.CreateShoppingItemParams
			if err := copier.Copy(&param, &reqb.Item[i]); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			itemRow, err := qtx.CreateShoppingItem(context.Background(), param)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			var itemRow2 db.UpdateShoppingItemRow
			if err := copier.Copy(&itemRow2, &itemRow); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			resp.Item = append(resp.Item, itemRow2)
		}
	}

	//// レスポンス型バリデーション
	//err = utils.ValidateStructTwoWay[db.GetShoppingListRow, docs.GetUsr](&row)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	// トランザクション終了
	tx.Commit(c)

	c.JSON(http.StatusOK, resp)
}
