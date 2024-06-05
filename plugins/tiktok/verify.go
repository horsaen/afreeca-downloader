package tiktok

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func CheckExists(roomId string) bool {
	// res, err := http.Get("https://www.tiktok.com/api/live/detail/?aid=1988&roomID=" + streamId)
	res, err := http.Get("https://webcast.tiktok.com/webcast/room/info/?aid=1988&room_id=" + roomId)

	if err != nil {
		log.Fatal(err)
	}

	bodyText, _ := io.ReadAll(res.Body)

	var live LiveV2

	json.Unmarshal(bodyText, &live)

	if live.Status == 10011 {
		return false
	} else {
		return true
	}
}

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
	// res, err := http.Get("https://www.tiktok.com/api/live/detail/?aid=1988&roomID=" + streamId)
	res, err := http.Get("https://webcast.tiktok.com/webcast/room/info/?aid=1988&room_id=" + streamId)

	if err != nil {
		log.Fatal(err)
	}

	bodyText, _ := io.ReadAll(res.Body)

	var live LiveV2

	json.Unmarshal(bodyText, &live)

	if live.Status == 0 {
		return true
	} else {
		return false
	}
}
