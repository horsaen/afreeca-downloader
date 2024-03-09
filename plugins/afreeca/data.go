package afreeca

import (
	"bufio"
	"encoding/json"
	"horsaen/afreeca-downloader/tools"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type User struct {
	Broad struct {
		BroadNo int `json:"broad_no"`
	} `json:"broad"`
	Station struct {
		Nickname string `json:"user_nick"`
		UserId   string `json:"user_id"`
	} `json:"station"`
	Code int `json:"code"`
}

type Base struct {
	ViewUrl string `json:"view_url"`
}

type Master struct {
	Channel struct {
		Aid string `json:"AID"`
	} `json:"CHANNEL"`
}

func UserData(bjId string) (string, int) {
	url := "https://bjapi.afreecatv.com/api/" + bjId + "/station"

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

	err = json.Unmarshal(bodyText, &user)

	if err != nil {
		log.Fatal(err)
	}

	return user.Station.Nickname, user.Broad.BroadNo
}

func PlaylistBase(broad_no int) string {
	url := "https://livestream-manager.afreecatv.com/broad_stream_assign.html?return_type=gcp_cdn&cors_origin_url=play.afreecatv.com&broad_key=" + strconv.Itoa(broad_no) + "-common-master-hls"

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	bodyText, _ := io.ReadAll(resp.Body)

	var base Base

	err = json.Unmarshal(bodyText, &base)

	if err != nil {
		log.Fatal(err)
	}

	return base.ViewUrl
}

func MasterPlaylist(bjId string, broad_no int, password string) string {
	url := "https://live.afreecatv.com/afreeca/player_live_api.php"

	payload := strings.NewReader("bid=" + bjId + "&bno=" + strconv.Itoa(broad_no) + "&type=aid&pwd=" + password + "&player_type=html5&stream_type=common&quality=master&mode=landing&from_api=0")

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, payload)

	cookies := tools.LoadCookies("afreeca")

	req.Header.Add("Cookie", "PdboxTicket="+cookies[0])
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/118.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Origin", "https://play.afreecatv.com")
	req.Header.Add("DNT", "1")
	req.Header.Add("Alt-Used", "live.afreecatv.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-site")
	req.Header.Add("Sec-GPC", "1")

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	bodyText, _ := io.ReadAll(resp.Body)

	var master Master

	err = json.Unmarshal(bodyText, &master)

	if err != nil {
		log.Fatal(err)
	}

	return master.Channel.Aid
}

func StreamList(baseUrl string, masterUrl string) string {
	url := baseUrl + "?aid=" + masterUrl

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	playlists := make([]string, 0)

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "auth_playlist") {
			playlists = append(playlists, line)
		}
	}

	return playlists[0]

}

func GetStream(bjId string, broad_no int, password string) string {
	base := PlaylistBase(broad_no)

	master := MasterPlaylist(bjId, broad_no, "")

	return base + "/" + StreamList(base, master)

}
