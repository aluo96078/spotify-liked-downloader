package processer

import (
	"fmt"
	"go_spotdl/library"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/miyago9267/ytvser/pkg/model"
)

type DownloadResult struct {
	Status string
	Data   model.Video
}

var outputPath string          // 下載輸出目錄
var maxConcurrentDownloads = 5 // 最大同時下載數
var DownloadResultChan = make(chan DownloadResult)

func loadSetting() {
	if maxThread := os.Getenv("MAX_THREAD"); maxThread != "" {
		threads, err := strconv.Atoi(maxThread)
		if err == nil {
			maxConcurrentDownloads = threads
		}
	}
	outputPath = os.Getenv("OUTPUT_PATH")
}

func MusicDownloader(accessToken string) {
	loadSetting()
	// 創建輸出目錄
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		os.Mkdir(outputPath, 0755)
	}

	// 獲取歌曲列表
	songList, err := library.CreateLikedSongList(accessToken)
	if err != nil {
		log.Fatalf("無法獲取歌曲列表: %v", err)
		return
	}
	if len(songList) <= maxConcurrentDownloads {
		maxConcurrentDownloads = len(songList)
	}
	// 協程同步控制
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrentDownloads) // 限制同時協程數量

	for _, song := range songList {
		video, err := library.FindMusicFromYoutube(song["name"])
		if err != nil {
			DownloadResultChan <- DownloadResult{"failed", *video}
			log.Printf("無法獲取 %s 的 YouTube 鏈接: %v", song["name"], err)
			continue
		}

		wg.Add(1)         // 增加等待組計數
		sem <- struct{}{} // 阻塞直到有可用的槽位

		go func(video *model.Video) {
			defer wg.Done()          // 協程完成後減少等待組計數
			defer func() { <-sem }() // 釋放槽位
			downloadMusic(video)
		}(video)
	}

	// 等待所有下載完成
	wg.Wait()
	close(DownloadResultChan)
}

func downloadMusic(video *model.Video) {
	// 執行下載
	cmd := exec.Command("yt-dlp",
		"--output", fmt.Sprintf("%s/%s.mp3", outputPath, video.ID),
		"--embed-thumbnail", "--add-metadata",
		"--extract-audio", "--audio-format", "mp3",
		"--audio-quality", "320K", video.URL,
	)

	if os.Getenv("COOKIES_FILE_PATH") != "" {
		cmd.Args = append(cmd.Args, "--cookies", os.Getenv("COOKIES_FILE_PATH"))
	}

	if err := cmd.Run(); err != nil {
		DownloadResultChan <- DownloadResult{"failed", *video}
		log.Printf("無法下載 %s - %s: %v", video.ID, video.URL, err)
		return
	}
	newName := fmt.Sprintf("%s/%s.mp3", outputPath, sanitizeFileName(video.Title))
	err := os.Rename(fmt.Sprintf("%s/%s.mp3", outputPath, video.ID), newName)
	if err != nil {
		log.Printf("無法重命名 %s.mp3: %v", video.ID, err)
	}
	DownloadResultChan <- DownloadResult{"success", *video}
}

func sanitizeFileName(name string) string {
	// 定義需要替換的非法字符
	illegalChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range illegalChars {
		name = strings.ReplaceAll(name, char, "_")
	}
	return name
}
