package soop

import (
	"bufio"
	"encoding/json"
	"horsaen/afreeca-downloader/tools"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Station struct {
	Code  int `json:"code"`
	Broad struct {
		BroadNo    int  `json:"broad_no"`
		IsPassword bool `json:"is_password"`
	} `json:"broad"`
}
type UserData struct {
	Channel struct {
		BNo    string `json:"BNO"`
		BJNick string `json:"BJNICK"`
	} `json:"CHANNEL"`
}
type Master struct {
	Channel struct {
		Aid string `json:"AID"`
	} `json:"CHANNEL"`
}
type Base struct {
	ViewUrl string `json:"view_url"`
}

func GetUserData(user string) (bool, bool, bool, string) {
	// moved to one function for effiency, also doensn't need cookies
	url := "https://chapi.sooplive.co.kr/api/" + user + "/station/"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:132.0) Gecko/20100101 Firefox/132.0")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")

	resp, _ := http.DefaultClient.Do(req)

	body, _ := io.ReadAll(resp.Body)

	var station Station

	json.Unmarshal(body, &station)

	exists := station.Code != 9000
	pwd := station.Broad.IsPassword
	online := station.Broad.BroadNo != 0
	broadNo := station.Broad.BroadNo

	return exists, pwd, online, strconv.Itoa(broadNo)
}

func GetMaster(user string) string {
	url := "https://live.sooplive.co.kr/afreeca/player_live_api.php?bjid=" + user

	payload := strings.NewReader("bid=" + user + "&type=aid&pwd=&player_type=html5&stream_type=common&quality=master&mode=landing&from_api=0&is_revive=false")

	req, _ := http.NewRequest("POST", url, payload)

	cookies := tools.LoadCookies("soop")

	req.Header.Add("Cookie", "AuthTicket="+cookies[0])
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:132.0) Gecko/20100101 Firefox/132.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := http.DefaultClient.Do(req)

	bodyText, _ := io.ReadAll(resp.Body)

	var master Master

	json.Unmarshal(bodyText, &master)

	return master.Channel.Aid
}

func GetBase(broad_key string) string {
	url := "https://livestream-manager.sooplive.co.kr/broad_stream_assign.html?return_type=gcp_cdn&use_cors=true&cors_origin_url=play.sooplive.co.kr&broad_key=" + broad_key + "-common-master-hls"

	req, _ := http.NewRequest("POST", url, nil)

	req.Header.Add("cookie", "AbroadChk=OK")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:132.0) Gecko/20100101 Firefox/132.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := http.DefaultClient.Do(req)

	bodyText, _ := io.ReadAll(resp.Body)

	var base Base

	json.Unmarshal(bodyText, &base)

	return strings.Replace(base.ViewUrl, "auth_master_playlist.m3u8", "", -1)
}

func GetStream(base string, master string) string {
	url := base + "auth_master_playlist.m3u8?aid=" + master

	req, _ := http.NewRequest("GET", url, nil)

	resp, _ := http.DefaultClient.Do(req)

	playlists := make([]string, 0)

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "auth_playlist") {
			playlists = append(playlists, line)
		}
	}

	return base + playlists[0]
}
