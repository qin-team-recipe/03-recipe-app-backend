package api

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/aopontann/gin-sqlc/docs"
	"github.com/gin-gonic/gin"

	"github.com/mattn/go-gimei"
)

func (s *Server) ListFeaturedChef(c *gin.Context) {
	const limit int32 = 10
	list, err := s.q.FakeListFeaturedChef(context.Background(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if chefs, ok := list.([]docs.FeaturedChef); ok {
		for i := 0; i < len(chefs); i++ {
			chefs[i].Name = gimei.NewName().String()
			chefs[i].NumFollower = rand.Int31n(1000)
			chefs[i].Score = rand.Int31n(100)
		}

		c.JSON(http.StatusOK, chefs)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list is not of type []docs.FeaturedChef"})
	}
}
