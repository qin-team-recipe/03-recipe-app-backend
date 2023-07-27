package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) CreateRecipe(c *gin.Context) {
	// TODO: リクエストボディのバリデーション

	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	row, err := s.q.CreateRecipe(context.Background(), data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: レスポンスのバリデーション

	c.JSON(http.StatusOK, row)
}
