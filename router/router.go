package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go_spotdl/processer"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	redirectURI = "http://localhost:8888/callback"
	authURL     = "https://accounts.spotify.com/authorize"
	tokenURL    = "https://accounts.spotify.com/api/token"
	apiURL      = "https://api.spotify.com/v1/me/tracks"
	state       = "random_state_string" // 防止 CSRF 攻擊
)

var (
	clientID     string
	clientSecret string
)

func loadKeys() {
	clientID = os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
}

var accessToken string

func Routering() {
	loadKeys()
	r := gin.Default()

	// 設定路由
	r.GET("/login", login)
	r.GET("/callback", callback)

	// 啟動伺服器
	log.Println("伺服器運行於 http://localhost:8888")
	if err := r.Run(":8888"); err != nil {
		log.Fatalf("伺服器啟動失敗: %v", err)
	}
}

// login 路由，生成 Spotify 認證 URL 並引導使用者登入
func login(c *gin.Context) {
	authEndpoint := fmt.Sprintf("%s?client_id=%s&response_type=code&redirect_uri=%s&scope=user-library-read&state=%s",
		authURL, clientID, redirectURI, state)
	c.Redirect(http.StatusFound, authEndpoint)
}

// callback 路由，處理 Spotify 的回調並交換 access token
func callback(c *gin.Context) {
	// 獲取 code 和 state
	code := c.Query("code")
	receivedState := c.Query("state")
	if receivedState != state {
		c.JSON(http.StatusBadRequest, gin.H{"error": "state 不匹配，可能存在 CSRF 攻擊"})
		return
	}

	// 發送請求交換 token
	data := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s", code, redirectURI)
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBuffer([]byte(data)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立請求"})
		return
	}
	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "交換 token 失敗"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "讀取 token 響應失敗"})
		return
	}

	// 解析 token 響應
	var tokenData map[string]interface{}
	json.Unmarshal(body, &tokenData)
	accessToken = tokenData["access_token"].(string)
	go processer.MusicDownloader(accessToken)
	c.Data(http.StatusOK, "text/html", []byte(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>完成</title>
			<script>
				window.onload = function() {
					window.close(); // 嘗試關閉標籤頁
				};
			</script>
		</head>
		<body>
			<h1>如果標籤頁沒有自動關閉，您可以手動關閉它。</h1>
		</body>
		</html>
	`))
	// c.JSON(http.StatusOK, gin.H{"message": "認證成功，您可以訪問 /liked-songs 來獲取按讚歌曲"})
	// c.Redirect(http.StatusFound, "/liked-songs")
}
