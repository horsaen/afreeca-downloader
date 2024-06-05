package chzzk

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type User struct {
	Content struct {
		Channel struct {
			Name string `json:"channelName"`
		} `json:"channel"`
		Stream string `json:"livePlaybackJson"`
	} `json:"content"`
}

type Stream struct {
	Media []struct {
		Playlist string `json:"path"`
	} `json:"media"`
}

func GetData(userId string) (string, string) {
	url := "https://api.chzzk.naver.com/service/v2/channels/" + userId + "/live-detail"

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/117.0")
	req.Header.Add("Accept", "application/json, text/plain, */*")

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	bodyText, _ := io.ReadAll(resp.Body)

	var user User
	json.Unmarshal(bodyText, &user)

	var stream Stream
	json.Unmarshal([]byte(user.Content.Stream), &stream)

	return user.Content.Channel.Name, stream.Media[0].Playlist
}
