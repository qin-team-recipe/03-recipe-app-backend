package api

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"reflect"

	db "github.com/aopontann/gin-sqlc/db/sqlc"
	"github.com/aopontann/gin-sqlc/docs"
	"github.com/aopontann/gin-sqlc/utils"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-gimei"
)

func (s *Server) ListFeaturedChef(c *gin.Context) {
	type featuredChefResponse struct {
		Data []db.FakeListFeaturedChefRow `json:"data"`
	}

	const limit int32 = 10
	var response featuredChefResponse
	var err error
	response.Data, err = s.q.FakeListFeaturedChef(context.Background(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.Data == nil || reflect.ValueOf(response.Data).IsNil() {
		response.Data = []db.FakeListFeaturedChefRow{}
	}

	// ダミーデータ作成（本番では消す）
	for i := 0; i < len(response.Data); i++ {
		response.Data[i].Name = gimei.NewName().String()
		response.Data[i].NumFollower = rand.Int31n(1000)
		response.Data[i].Score = rand.Int31n(100)
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[featuredChefResponse, docs.FeaturedChef](&response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) GetChef(c *gin.Context) {
	// パスパラメータ取り出し
	id, err := utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// 問い合わせ処理
	row, err := s.q.GetChef(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.GetChefRow, docs.GetChef](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}

func (s *Server) CreateChef(c *gin.Context) {
	// リクエストボディを構造体にバインド
	reqb := docs.PostApiChefsJSONRequestBody{}
	if err := c.ShouldBind(&reqb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 構造体からJSONに変換
	jsn, err := json.Marshal(&reqb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// 新規登録処理
	row, err := s.q.CreateChef(context.Background(), jsn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.CreateChefRow, docs.CreateChef](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}

func (s *Server) UpdateChef(c *gin.Context) {
	var param db.UpdateChefParams
	var err error

	// パスパラメータ取り出し
	param.ID, err = utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// リクエストボディを構造体にバインド
	reqb := docs.PutApiChefsIdJSONRequestBody{}
	if err := c.ShouldBind(&reqb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 構造体からJSONに変換
	param.Data, err = json.Marshal(&reqb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// 更新処理
	row, err := s.q.UpdateChef(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.UpdateChefRow, docs.UpdateChef](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}

func (s *Server) DeleteChef(c *gin.Context) {
	// パスパラメータ取り出し
	id, err := utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// 問い合わせ処理
	row, err := s.q.DeleteChef(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.DeleteChefRow, docs.DeletedChef](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}

func (s *Server) SearchChef(c *gin.Context) {
	type searchChefResponse struct {
		Data []db.SearchChefRow `json:"data"`
	}

	// クエリパラメータ取り出し
	query := c.Query("q")

	// 全文検索
	var response searchChefResponse
	var err error
	response.Data, err = s.q.SearchChef(context.Background(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.Data == nil || reflect.ValueOf(response.Data).IsNil() {
		response.Data = []db.SearchChefRow{}
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[searchChefResponse, docs.SearchChef](&response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
