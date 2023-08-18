package api

import (
	"bytes"
	"github.com/aopontann/gin-sqlc/docs"
	"github.com/aopontann/gin-sqlc/utils"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gocv.io/x/gocv"
)

func (s *Server) GetImage(c *gin.Context) {
	// 画像入出力先（TODO：外から与える）
	dir := "./"

	// クエリパラメータ
	path := dir + c.Query("path")

	// 画像ファイルopen
	reader, err := os.Open(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer reader.Close()

	// 画像ファイルサイズ取得
	stats, statsErr := reader.Stat()
	if statsErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	contentLength := stats.Size()

	c.DataFromReader(http.StatusOK, contentLength, "image/webp", reader, nil)
}

func (s *Server) PostImage(c *gin.Context) {
	// 画像入出力先（TODO：外から与える）
	dir := "./"

	// レスポンス
	type postWebpResponse struct {
		Path string `form:"path"`
	}
	var response postWebpResponse

	// リクエストボディ
	fi, _, err := c.Request.FormFile("uploadFile")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer fi.Close()

	// ファイルをバイナリ形式で読み込み
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, fi); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// バイナリを画像形式に変換
	img, err := gocv.IMDecode(buf.Bytes(), gocv.IMReadColor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 画像をWebP形式に変換。圧縮率80にした。
	webp, err := gocv.IMEncodeWithParams(".webp", img, []int{gocv.IMWriteWebpQuality, 80})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 画像ファイル名（TODO：一応ファイルの存在チェック）
	u, err := uuid.NewRandom()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.Path = u.String() + ".webp"

	// 出力ファイルopen
	fo, err := os.Create(dir + response.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer fo.Close()

	// WebP画像をwrite
	if _, err := fo.Write(webp.GetBytes()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンス型バリデーション
	err = utils.ValidateStructTwoWay[postWebpResponse, docs.PostImage](&response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
