package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	db "github.com/aopontann/gin-sqlc/db/sqlc"
	"github.com/aopontann/gin-sqlc/docs"
	"github.com/aopontann/gin-sqlc/utils"
	"github.com/gin-gonic/gin"
)

type createUserReqBody struct {
	Name string `json:"name" binding:"required"`
}

func (s *Server) CreateUser(c *gin.Context) {
	// リクエストボディを構造体にバインド
	reqb := createUserReqBody{}
	if err := c.ShouldBind(&reqb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "リクエストボディを構造体にバインドするのに失敗しました。", "error": err.Error()})
		return
	}

	// CookieにセットされたセッションIDを使い、redisからユーザ情報を取得する
	sid, _ := c.Cookie("session_id")
	data, err := s.rbd.Get(context.Background(), sid).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"title": "redis から情報を引き出すのに失敗しました。", "error": err.Error()})
		return
	}

	// redisから取得したユーザ情報を構造体に変換する
	var uInfo userInfo
	if err := json.Unmarshal([]byte(data), &uInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "ユーザ情報を構造体に変換するのに失敗しました。", "error": err.Error()})
		return
	}

	// メールアドレスが google かチェック
	var authServer string
	re := regexp.MustCompile(`@gmail.com$`)
	if re.MatchString(uInfo.Email) {
		authServer = "google"
	} else {
		// google 以外のメールアドレスだった場合
		c.JSON(http.StatusBadRequest, gin.H{"title": "gmail 以外のメールが入力されました。", "error": err.Error()})
		return
	}

	// DBにユーザ情報を保存する
	res, err := s.q.CreateUser(context.Background(), db.CreateUserParams{Name: reqb.Name, Email: uInfo.Email, AuthServer: authServer, AuthUserinfo: []byte(data)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "SQLの処理に失敗しました。", "error": err.Error()})
		return
	}

	strID := fmt.Sprintf("%x-%x-%x-%x-%x", res.ID.Bytes[0:4], res.ID.Bytes[4:6], res.ID.Bytes[6:8], res.ID.Bytes[8:10], res.ID.Bytes[10:16])
	// 構造体のまま redis に保存できないため、byte型に変換する
	b, err := json.Marshal(redisValue{ID: strID, Email: res.Email})
	if err != nil {
		return
	}
	s.rbd.Set(context.Background(), sid, b, 24*time.Hour)

	c.JSON(http.StatusCreated, gin.H{"id": res.ID, "name": res.Name, "email": res.Email})
}

// 仮で作成したハンドラ関数
func (s *Server) GetUserId(c *gin.Context) {
	// Authentication()でセットしたメールアドレスを取得
	email := c.MustGet("email").(string)
	uid, err := s.q.GetUserId(context.Background(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "認証に失敗しました。", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": uid, "email": email})
}

func (s *Server) GetUser(c *gin.Context) {
	// パスパラメータ取り出し
	id, err := utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "パスパラメータ取り出しに失敗しました。", "error": err.Error()})
	}

	// 問い合わせ処理
	row, err := s.q.GetUser(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "SQLの処理に失敗しました。", "error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.GetUserRow, docs.GetUsr](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "型のバリデーションが失敗しました。", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}

func (s *Server) GetSelf(c *gin.Context) {
	// Authentication()でセットしたメールアドレスを取得
	email := c.MustGet("email").(string)

	// 問い合わせ処理
	row, err := s.q.GetSelf(context.Background(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "SQLの処理に失敗しました。", "error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.GetSelfRow, docs.GetUsr](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "型のバリデーションが失敗しました。", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}

func (s *Server) UpdateSelf(c *gin.Context) {
	// emailを取得
	var param db.UpdateUserParams

	// Authentication()でセットしたメールアドレスを取得
	param.Email = c.MustGet("email").(string)

	// リクエストボディを構造体にバインド
	reqb := docs.PutApiUserUsersJSONRequestBody{}
	if err := c.ShouldBind(&reqb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "リクエストボディを構造体に変換を失敗しました。", "error": err.Error()})
		return
	}

	// 構造体からJSONに変換
	var err error
	param.Data, err = json.Marshal(&reqb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"title": "構造体からJSONに変換を失敗しました。", "error": err.Error()})
	}

	// 更新処理
	row, err := s.q.UpdateUser(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "SQLの処理に失敗しました。", "error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.UpdateUserRow, docs.UpdateUsr](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "型のバリデーションが失敗しました。", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, row)
}

func (s *Server) DeleteSelf(c *gin.Context) {
	// Authentication()でセットしたメールアドレスを取得
	email := c.MustGet("email").(string)

	// 削除処理
	row, err := s.q.DeleteUser(context.Background(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "SQLの処理に失敗しました。", "error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[db.DeleteUserRow, docs.DeletedUsr](&row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"title": "型のバリデーションが失敗しました。", "error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, row)
}
