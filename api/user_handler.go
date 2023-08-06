package api

import (
	"context"
	"encoding/json"
	"github.com/aopontann/gin-sqlc/utils"
	"net/http"
	"regexp"
	"time"

	db "github.com/aopontann/gin-sqlc/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createUserReqBody struct {
	Name string `json:"name" binding:"required"`
}

func (s *Server) CreateUser(c *gin.Context) {
	// リクエストボディを構造体にバインド
	reqb := createUserReqBody{}
	if err := c.ShouldBind(&reqb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// CookieにセットされたセッションIDを使い、redisからユーザ情報を取得する
	sid, _ := c.Cookie("session_id")
	data, err := s.rbd.Get(context.Background(), sid).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// redisから取得したユーザ情報を構造体に変換する
	var uInfo userInfo
	if err := json.Unmarshal([]byte(data), &uInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// redisに保存されているユーザ情報をメールアドレスで上書きする
	s.rbd.Set(context.Background(), sid, uInfo.Email, 24*time.Hour)

	// メールアドレスが google かチェック
	var authServer string
	re := regexp.MustCompile(`@gmail.com$`)
	if re.MatchString(uInfo.Email) {
		authServer = "google"
	} else {
		// google 以外のメールアドレスだった場合
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// DBにユーザ情報を保存する
	res, err := s.q.CreateUser(context.Background(), db.CreateUserParams{Name: reqb.Name, Email: uInfo.Email, AuthServer: authServer, AuthUserinfo: []byte(data)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": res.ID, "name": res.Name, "email": res.Email})
}

// 仮で作成したハンドラ関数
func (s *Server) GetUserId(c *gin.Context) {
	// Authentication()でセットしたメールアドレスを取得
	email := c.MustGet("email").(string)
	uid, err := s.q.GetUserId(context.Background(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": uid, "email": email})
}

func (s *Server) GetUser(c *gin.Context) {
	// パスパラメータ取り出し
	id, err := utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// 問い合わせ処理
	row, err := s.q.GetUser(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//// レスポンス型バリデーション
	//err = utils.ValidateStructTwoWay[db.VUsr, docs.GetUser](&row)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	c.JSON(http.StatusOK, row)
}

func (s *Server) UpdateUser(c *gin.Context) {
	var param db.UpdateUserParams
	var err error

	// パスパラメータ取り出し
	param.ID, err = utils.StrToUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	//// リクエストボディを構造体にバインド
	//reqb := docs.PutApiUpdateUserJSONRequestBody{}
	//if err := c.ShouldBind(&reqb); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//
	//// 構造体からJSONに変換
	//param.Data, err = json.Marshal(&reqb)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//}
	param.Data, err = c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新処理
	row, err := s.q.UpdateUser(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//// レスポンス型バリデーション
	//err = utils.ValidateStructTwoWay[db.VUser, docs.UpdateUser](&row)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	c.JSON(http.StatusOK, row)
}
