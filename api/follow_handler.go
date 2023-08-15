package api

import (
	"context"
	"github.com/aopontann/gin-sqlc/utils"
	"net/http"

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
	email := c.MustGet("email").(string)
	param.UsrID, err = s.q.GetUserId(context.Background(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 新規登録処理
	row, err := s.q.CreateFollowChef(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//// レスポンス型バリデーション
	//err = utils.ValidateStructTwoWay[db.VRecipe, docs.CreateUsrRecipe](&row)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

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
	email := c.MustGet("email").(string)
	param.UsrID, err = s.q.GetUserId(context.Background(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 登録解除処理
	row, err := s.q.DeleteFollowChef(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//// レスポンス型バリデーション
	//err = utils.ValidateStructTwoWay[db.VRecipe, docs.CreateUsrRecipe](&row)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

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
	email := c.MustGet("email").(string)
	param.UsrID, err = s.q.GetUserId(context.Background(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// お気に入りしているか
	response.Exists, err = s.q.ExistsFollowChef(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//// レスポンス型バリデーション
	//err = utils.ValidateStructTwoWay[db.VRecipe, docs.CreateUsrRecipe](&row)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	c.JSON(http.StatusOK, response)
}

func (s *Server) GetFollowChef(c *gin.Context) {
	// usrIdを取得
	email := c.MustGet("email").(string)
	usrID, err := s.q.GetUserId(context.Background(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// お気に入りしているか
	list, err := s.q.GetFollowChef(context.Background(), usrID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//// レスポンス型バリデーション
	//err = utils.ValidateStructTwoWay[db.VRecipe, docs.CreateUsrRecipe](&row)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	c.JSON(http.StatusOK, list)
}
