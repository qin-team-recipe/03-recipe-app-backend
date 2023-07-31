package api

import (
	"context"
	"math/rand"
	"net/http"

	db "github.com/aopontann/gin-sqlc/db/sqlc"
	"github.com/aopontann/gin-sqlc/docs"
	"github.com/aopontann/gin-sqlc/utils"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-gimei"
)

type featuredChefResponse struct {
	Data []db.FakeListFeaturedChefRow `json:"data"`
}

func (s *Server) ListFeaturedChef(c *gin.Context) {
	const limit int32 = 10
	var response featuredChefResponse
	var err error
	response.Data, err = s.q.FakeListFeaturedChef(context.Background(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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
