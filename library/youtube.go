package library

import (
	"os/exec"

	"github.com/miyago9267/ytvser/pkg/model"
	"github.com/miyago9267/ytvser/pkg/searcher"
)

func FindMusicFromYoutube(songName string) (*model.Video, error) {
	s := searcher.NewSearcher()
	videos, err := s.YoutubeSearch(songName, 10)
	if err != nil {
		return nil, err
	}
	return &videos[0], nil
}

func YtDlpInstalled() bool {
	cmd := exec.Command("yt-dlp", "--version")
	err := cmd.Run()
	return err == nil
}
