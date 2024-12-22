package library

import (
	"encoding/json"
	"io"
	"net/http"
)

const apiURL = "https://api.spotify.com/v1/me/tracks"

func CreateLikedSongList(accessToken string) ([]map[string]string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析按讚歌曲
	var result struct {
		Items []struct {
			Track struct {
				Name    string `json:"name"`
				Artists []struct {
					Name string `json:"name"`
				} `json:"artists"`
				ExternalURLs map[string]string `json:"external_urls"`
			} `json:"track"`
		} `json:"items"`
	}
	json.Unmarshal(body, &result)

	var likedSongs []map[string]string
	for _, item := range result.Items {
		likedSongs = append(likedSongs, map[string]string{
			"name":   item.Track.Name,
			"artist": item.Track.Artists[0].Name,
			"url":    item.Track.ExternalURLs["spotify"],
		})
	}
	return likedSongs, nil
}
