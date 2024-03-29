package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type userInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

type redisValue struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// Scopes: OAuth 2.0 scopes provide a way to limit the amount of access that is granted to an access token.
var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func (s *Server) OauthGoogleLogin(c *gin.Context) {

	// Create oauthState cookie
	oauthState := generateStateOauthCookie(c)
	log.Printf("oauth google state, %s", oauthState)

	/*
		AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
		validate that it matches the the state query parameter on your redirect callback.
	*/
	u := googleOauthConfig.AuthCodeURL(oauthState)
	c.Redirect(http.StatusTemporaryRedirect, u)
}

func (s *Server) OauthGoogleCallback(c *gin.Context) {
	// Read oauthState from Cookie
	oauthState, _ := c.Cookie("oauthstate")
	if c.Query("state") != oauthState {
		log.Println("invalid oauth google state")
		// c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	data, err := getUserDataFromGoogle(c.Query("code"))
	if err != nil {
		log.Println(err.Error())
		// c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	// セッションIDを生成
	guid := xid.New()

	c.SetCookie("session_id", guid.String(), 3600, "/", "localhost", false, true)

	var uInfo userInfo
	if err := json.Unmarshal(data, &uInfo); err != nil {
		return
	}

	uid, err := s.q.GetUserId(context.Background(), uInfo.Email)
	if !uid.Valid {
		s.rbd.Set(context.Background(), guid.String(), data, 3600*time.Second)
		afterLoginUrl, err := c.Cookie("after_login_url")
		if (err != nil) || (afterLoginUrl == "") {
			c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/create-user")
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, afterLoginUrl)
		return
	}
	if err != nil {
		fmt.Printf("ExistsUser Error: %s\n", err.Error())
		return
	}

	strID := fmt.Sprintf("%x-%x-%x-%x-%x", uid.Bytes[0:4], uid.Bytes[4:6], uid.Bytes[6:8], uid.Bytes[8:10], uid.Bytes[10:16])
	b, err := json.Marshal(redisValue{ID: strID, Email: uInfo.Email})
	if err != nil {
		return
	}
	s.rbd.Set(context.Background(), guid.String(), b, 3600*time.Second)
	afterLoginUrl, err := c.Cookie("after_login_url")
	if (err != nil) || (afterLoginUrl == "") {
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000")
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, afterLoginUrl)
}

func generateStateOauthCookie(c *gin.Context) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	c.SetCookie("oauthstate", state, 3600, "/", "localhost", false, true)

	return state
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	// Use code to get token and get user info from Google.

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}

func (s *Server) Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// CookieにセットされたセッションIDを使い、redisからユーザのメールアドレスを取得する
		sid, _ := c.Cookie("session_id")
		data, err := s.rbd.Get(context.Background(), sid).Result()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
		}

		// redisから取得したユーザ情報を構造体に変換する
		var rv redisValue
		if err := json.Unmarshal([]byte(data), &rv); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// 次に実行されるハンドラ関数で取得したメールアドレスを使うための処理
		c.Set("email", rv.Email)
		c.Set("id", rv.ID)
		c.Set("rv", rv)
	}
}
