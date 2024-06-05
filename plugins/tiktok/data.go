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

type LiveV2 struct {
	Data struct {
		Owner struct {
			Nickname string `json:"nickname"`
		} `json:"owner"`
		StreamUrl struct {
			FlvPullUrl struct {
				FullHD string `json:"FULL_HD1"`
			} `json:"flv_pull_url"`
			HlsPullUrl string `json:"hls_pull_url"`
		} `json:"stream_url"`
	} `json:"data"`
	Status int `json:"status_code"`
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

func GetPlaylist(streamId string) (string, string) {
	// res, err := http.Get("https://www.tiktok.com/api/live/detail/?aid=1988&roomID=" + streamId)
	res, err := http.Get("https://webcast.tiktok.com/webcast/room/info/?aid=1988&room_id=" + streamId)

	if err != nil {
		log.Fatal(err)
	}

	bodyText, _ := io.ReadAll(res.Body)

	var live LiveV2

	json.Unmarshal(bodyText, &live)

	if live.Data.StreamUrl.HlsPullUrl == "" {
		return "flv", live.Data.StreamUrl.FlvPullUrl.FullHD
	} else {
		return "hls", live.Data.StreamUrl.HlsPullUrl
	}
}
