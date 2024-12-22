package main

import (
	"go_spotdl/library"
	"go_spotdl/processer"
	"go_spotdl/router"
	"log"
)

func main() {
	library.LoadEnv()
	if !library.YtDlpInstalled() {
		log.Fatalf("請先安裝 yt-dlp 並確保環境變數設置完成！")
		return
	}
	// 啟動 OAuth 路由
	go router.Routering()
	// 自動引導使用者登入
	library.OpenBrowser("http://localhost:8888/login")
	for result := range processer.DownloadResultChan {
		log.Printf("download music '%s' %s！", result.Data.Title, result.Status)
	}
	log.Println("下載完成！")
}
