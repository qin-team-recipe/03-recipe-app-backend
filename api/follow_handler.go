package api

import (
	"context"
	"github.com/aopontann/gin-sqlc/docs"
	"github.com/aopontann/gin-sqlc/utils"
	"net/http"
	"reflect"

	db "github.com/aopontann/gin-sqlc/db/sqlc"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateFollowChef(c *gin.Context) {
	// パスパラメータ取り出し
	param := db.CreateFollowChefParams{}
	var err error
	param.ChefID, err = utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// usrIdを取得
	var status int
	param.UsrID, _, err, status = GetRedisInfo(c, s)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// 新規登録処理
	row, err := s.q.CreateFollowChef(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.FollowingChef, docs.CreateFollowChef](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}

func (s *Server) DeleteFollowChef(c *gin.Context) {
	// パスパラメータ取り出し
	param := db.DeleteFollowChefParams{}
	var err error
	param.ChefID, err = utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// usrIdを取得
	var status int
	param.UsrID, _, err, status = GetRedisInfo(c, s)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// 登録解除処理
	row, err := s.q.DeleteFollowChef(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.FollowingChef, docs.DeletedFollowChef](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}

func (s *Server) ExistsFollowChef(c *gin.Context) {
	type existsFollowResponse struct {
		Exists bool `json:"exists"`
	}
	var response existsFollowResponse

	// パスパラメータ取り出し
	param := db.ExistsFollowChefParams{}
	var err error
	param.ChefID, err = utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// usrIdを取得
	var status int
	param.UsrID, _, err, status = GetRedisInfo(c, s)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// お気に入りしているか
	response.Exists, err = s.q.ExistsFollowChef(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[existsFollowResponse, docs.ExistsFollowChef](&response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) ListFollowChef(c *gin.Context) {
	type listFollowChefResponse struct {
		Data []db.ListFollowChefRow `json:"data"`
	}
	var response listFollowChefResponse

	// usrIdを取得
	usrID, _, err, status := GetRedisInfo(c, s)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// フォローしている有名シェフ一覧を取得
	response.Data, err = s.q.ListFollowChef(context.Background(), usrID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.Data == nil || reflect.ValueOf(response.Data).IsNil() {
		response.Data = []db.ListFollowChefRow{}
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[listFollowChefResponse, docs.ListFollowChef](&response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) ListFollowChefNewRecipe(c *gin.Context) {
	type followNewRecipeResponse struct {
		Data []db.ListFollowChefNewRecipeRow `json:"data"`
	}
	var response followNewRecipeResponse

	// usrIdを取得
	usrID, _, err, status := GetRedisInfo(c, s)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// フォローしている有名シェフの新着レシピ一覧を取得
	response.Data, err = s.q.ListFollowChefNewRecipe(context.Background(), usrID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.Data == nil || reflect.ValueOf(response.Data).IsNil() {
		response.Data = []db.ListFollowChefNewRecipeRow{}
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[followNewRecipeResponse, docs.ListFollowChefRecipe](&response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) CreateFollowUser(c *gin.Context) {
	// パスパラメータ取り出し
	param := db.CreateFollowUserParams{}
	var err error
	param.FolloweeID, err = utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// usrIdを取得
	var status int
	param.FollowerID, _, err, status = GetRedisInfo(c, s)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// FolloweeIDとFollowerIDが同じのときエラー
	if param.FolloweeID == param.FollowerID {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "自分をフォローすることはできません"})
		return
	}

	// 新規登録処理
	row, err := s.q.CreateFollowUser(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.FollowingUser, docs.CreateFollowUser](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}

func (s *Server) DeleteFollowUser(c *gin.Context) {
	// パスパラメータ取り出し
	param := db.DeleteFollowUserParams{}
	var err error
	param.FolloweeID, err = utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// usrIdを取得
	var status int
	param.FollowerID, _, err, status = GetRedisInfo(c, s)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// 登録解除処理
	row, err := s.q.DeleteFollowUser(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.FollowingUser, docs.DeletedFollowUser](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}

func (s *Server) ExistsFollowUser(c *gin.Context) {
	type existsFollowResponse struct {
		Exists bool `json:"exists"`
	}
	var response existsFollowResponse

	// パスパラメータ取り出し
	param := db.ExistsFollowUserParams{}
	var err error
	param.FolloweeID, err = utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// usrIdを取得
	var status int
	param.FollowerID, _, err, status = GetRedisInfo(c, s)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// お気に入りしているか
	response.Exists, err = s.q.ExistsFollowUser(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[existsFollowResponse, docs.ExistsFollowUser](&response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) ListFollowUser(c *gin.Context) {
	type listFollowUserResponse struct {
		Data []db.ListFollowUserRow `json:"data"`
	}
	var response listFollowUserResponse

	// usrIdを取得
	followerID, _, err, status := GetRedisInfo(c, s)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// フォローしている有名シェフ一覧を取得
	response.Data, err = s.q.ListFollowUser(context.Background(), followerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.Data == nil || reflect.ValueOf(response.Data).IsNil() {
		response.Data = []db.ListFollowUserRow{}
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[listFollowUserResponse, docs.ListFollowUser](&response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
