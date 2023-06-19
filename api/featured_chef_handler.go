package api

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mattn/go-gimei"
)

func (s *Server) ListFeaturedChef(c *gin.Context) {
	list, err := s.q.FakeListFeaturedChef(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := 0; i < len(list); i++ {
		list[i].Name = gimei.NewName().String()
		list[i].NumFollower = rand.Int31n(1000)
		list[i].KpiFeatured = rand.Int31n(100)
	}
	c.JSON(http.StatusOK, list)
}
