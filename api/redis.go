package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aopontann/gin-sqlc/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Server) GetRedisInfo(c *gin.Context) (pgtype.UUID, string, int, error) {
	// CookieにセットされたセッションIDを使い、redisからユーザ情報を取得する
	sid, _ := c.Cookie("session_id")
	data, err := s.rbd.Get(context.Background(), sid).Result()
	if err != nil {
		return pgtype.UUID{}, "", http.StatusUnauthorized, err
	}

	// redisから取得したユーザ情報を構造体に変換する
	var rv redisValue
	if err := json.Unmarshal([]byte(data), &rv); err != nil {
		return pgtype.UUID{}, "", http.StatusInternalServerError, err
	}

	// ユーザーIDを取得
	usrId, err := utils.StrToUUID(rv.ID)
	if err != nil {
		return pgtype.UUID{}, "", http.StatusInternalServerError, err
	}

	return usrId, rv.Email, http.StatusOK, nil
}
