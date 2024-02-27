package tiktok

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func DvrCheck(streamId string) bool {
	for {
		if CheckOnline(streamId) {
			return true
		}
		fmt.Printf("User %s offline, rechecking in 3 minutes.\n", streamId)
		time.Sleep(3 * time.Minute)
	}
}

func CheckOnline(streamId string) bool {
	res, err := http.Get("https://www.tiktok.com/api/live/detail/?aid=1988&roomID=" + streamId)

	if err != nil {
		log.Fatal(err)
	}

	bodyText, _ := io.ReadAll(res.Body)

	var live Live

	json.Unmarshal(bodyText, &live)

	return live.LiveRoomInfo.LiveUrl != ""
}
