package bigo

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
)

type Response struct {
	Code int    `json:"code"`
	Data Data   `json:"data"`
	Msg  string `json:"msg"`
	Flag *int   `json:"flag"`
}

type Data struct {
	Sid          int64    `json:"sid"`
	SiteId       string   `json:"siteId"`
	Uid          int      `json:"uid"`
	Avatar       string   `json:"avatar"`
	NickName     string   `json:"nick_name"`
	CountryCode  string   `json:"country_code"`
	GameTitle    string   `json:"gameTitle"`
	GameId       int      `json:"gameId"`
	RoomTopic    string   `json:"roomTopic"`
	Snapshot     string   `json:"snapshot"`
	Alive        int      `json:"alive"`
	RoomId       string   `json:"roomId"`
	RoomStatus   int      `json:"roomStatus"`
	Covers       *string  `json:"covers"`
	ClientIp     string   `json:"client_ip"`
	HlsSrc       string   `json:"hls_src"`
	CdnSrc       []string `json:"cdn_src"`
	PassRoom     bool     `json:"passRoom"`
	Reserver     int      `json:"reserver"`
	ClientBigoId string   `json:"clientBigoId"`
	RoomType     string   `json:"roomType"`
}

func GetStreamData(userId string) (username string, nickname string, hlsUrl string) {
	client := &http.Client{}

	payload := strings.NewReader("siteId=" + userId + "&=verify%3D")

	req, err := http.NewRequest("POST", "https://ta.bigo.tv/official_website/studio/getInternalStudioInfo", payload)

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/119.0")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Accept-Encoding", "br")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Origin", "https://www.bigo.tv")
	req.Header.Add("DNT", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Referer", "https://www.bigo.tv/")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-site")
	req.Header.Add("Sec-GPC", "1")
	req.Header.Add("TE", "trailers")

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	reader := brotli.NewReader(resp.Body)

	bodyText, err := io.ReadAll(reader)

	if err != nil {
		log.Fatal(err)
	}

	var response Response

	err = json.Unmarshal(bodyText, &response)
	if err != nil {
		log.Fatal(err)
	}

	return response.Data.SiteId, response.Data.NickName, response.Data.HlsSrc
}
