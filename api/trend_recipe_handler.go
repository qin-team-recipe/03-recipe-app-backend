package api

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/aopontann/gin-sqlc/docs"

	"github.com/gin-gonic/gin"

	"github.com/mattn/go-gimei"
)

func (s *Server) ListTrendRecipe(c *gin.Context) {
	const limit int32 = 10
	list, err := s.q.FakeListTrendRecipe(context.Background(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := 0; i < len(list); i++ {
		list[i].Title = gimei.NewName().First.Katakana()
		list[i].Introduction = gimei.NewAddress().String() + "。" + gimei.NewAddress().String() + "。" + gimei.NewAddress().String() + "。" + gimei.NewAddress().String() + "。" + gimei.NewAddress().String() + "。"
		list[i].NumFav = rand.Int31n(1000)
		list[i].Score = rand.Int31n(100)
	}

	s := list.([]docs.TrendRecipe)

	c.JSON(http.StatusOK, list)
}
