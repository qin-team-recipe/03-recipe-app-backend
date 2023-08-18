package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aopontann/gin-sqlc/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetRedisInfo(c *gin.Context, s *Server) (pgtype.UUID, string, error, int) {
	// CookieにセットされたセッションIDを使い、redisからユーザ情報を取得する
	sid, _ := c.Cookie("session_id")
	data, err := s.rbd.Get(context.Background(), sid).Result()
	if err != nil {
		return pgtype.UUID{}, "", err, http.StatusUnauthorized
	}

	// redisから取得したユーザ情報を構造体に変換する
	var rv redisValue
	if err := json.Unmarshal([]byte(data), &rv); err != nil {
		return pgtype.UUID{}, "", err, http.StatusInternalServerError
	}

	// ユーザーIDを取得
	usrId, err := utils.StrToUUID(rv.ID)
	if err != nil {
		return pgtype.UUID{}, "", err, http.StatusInternalServerError
	}

	return usrId, rv.Email, nil, http.StatusOK
}
