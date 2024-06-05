package soop

import (
	"encoding/json"
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"io"
	"net/http"
	"time"
)

func CheckExists(channelId string) bool {
	cookies := tools.LoadCookies("soop")

	client := &http.Client{}

	url := "https://api.sooplive.com/stream/info/" + channelId

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("client-id", cookies[0])

	resp, _ := client.Do(req)

	return resp.StatusCode != 400
}

func DvrCheck(channelId string) bool {
	for {
		if CheckOnline(channelId) {
			return true
		}
		fmt.Printf("User %s offline, rechecking in 3 minutes.\n", channelId)
		time.Sleep(3 * time.Minute)
	}
}

func CheckOnline(channelId string) bool {
	cookies := tools.LoadCookies("soop")

	client := &http.Client{}

	url := "https://api.sooplive.com/stream/info/" + channelId

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("client-id", cookies[0])

	resp, _ := client.Do(req)

	bodyText, _ := io.ReadAll(resp.Body)

	var user User

	json.Unmarshal(bodyText, &user)

	return user.IsStream

}
