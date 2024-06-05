package chzzk

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func DvrCheck(userId string) bool {
	for {
		if CheckOnline(userId) {
			return true
		}
		fmt.Printf("User %s offline, rechecking in 3 minutes.\n", userId)
		time.Sleep(3 * time.Minute)
	}
}

func CheckOnline(userId string) bool {
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

	if stream.Media[0].Playlist != "" {
		return true
	} else {
		return false
	}
}
