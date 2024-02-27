package tiktok

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
)

type Live struct {
	LiveRoomInfo struct {
		// OwnerInfo struct {
		// 	UniqueId string `json:"uniqueId"`
		// } `json:"ownerInfo"`
		LiveUrl string `json:"liveUrl"`
	} `json:"LiveRoomInfo"`
}

func RoomId(userId string) string {
	client := &http.Client{}

	url := "https://www.tiktok.com/@" + userId + "/live"

	req, _ := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	pattern := `room_id=(\d{19})"`

	bodyText, _ := io.ReadAll(resp.Body)

	re := regexp.MustCompile(pattern)

	match := re.FindStringSubmatch(string(bodyText))

	if len(match) > 1 {
		roomID := match[1]
		return roomID
	} else {
		return ""
	}
}

func GetPlaylist(streamId string) string {
	res, err := http.Get("https://www.tiktok.com/api/live/detail/?aid=1988&roomID=" + streamId)

	if err != nil {
		log.Fatal(err)
	}

	bodyText, _ := io.ReadAll(res.Body)

	var live Live

	json.Unmarshal(bodyText, &live)

	return live.LiveRoomInfo.LiveUrl
}
