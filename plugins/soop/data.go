package soop

import (
	"bufio"
	"encoding/json"
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Channel struct {
	Channel struct {
		TK  string `json:"TK"`
		AID string `json:"AID"`
	} `json:"CHANNEL"`
}

type Server struct {
	ViewUrl string `json:"view_url"`
}

type User struct {
	BroadNo int `json:"broadNo"`
}

func GetBroadNo(bjid string) string {
	res, _ := http.Get("https://api-channel.sooplive.co.kr/v1.1/channel/" + bjid + "/home/section/broad")

	body, _ := io.ReadAll(res.Body)

	var user User
	json.Unmarshal(body, &user)

	return strconv.Itoa(user.BroadNo)
}

func GetStreamTk(bjid, broad_no string) string {
	url := "https://live.sooplive.co.kr/afreeca/player_live_api.php?bjid=" + bjid

	payload := fmt.Sprintf(
		"bid=%s&broad_no=%s&type=live&pwd=&player_type=html5&stream_type=common&quality=HD&mode=landing&from_api=0&is_revive=false",
		bjid, broad_no,
	)

	req, _ := http.NewRequest("POST", url, strings.NewReader(payload))

	cookies := tools.LoadCookies("soop")

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:144.0) Gecko/20100101 Firefox/144.0")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "AuthTicket="+cookies[0])

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var channel Channel
	json.Unmarshal(body, &channel)

	return channel.Channel.TK
}

func GetStreamAid(bjid, broad_no string) string {
	url := "https://live.sooplive.co.kr/afreeca/player_live_api.php?bjid=" + bjid

	payload := fmt.Sprintf(
		"bid=%s&broad_no=%s&type=aid&pwd=&player_type=html5&stream_type=common&quality=master&mode=landing&from_api=0&is_revive=false",
		bjid, broad_no,
	)

	req, _ := http.NewRequest("POST", url, strings.NewReader(payload))

	cookies := tools.LoadCookies("soop")

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:144.0) Gecko/20100101 Firefox/144.0")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "AuthTicket="+cookies[0])

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var channel Channel
	json.Unmarshal(body, &channel)

	return channel.Channel.AID
}

func GetStreamServer(broad_no string) string {
	url := fmt.Sprintf("https://livestream-manager.sooplive.co.kr/broad_stream_assign.html?return_type=gcp_cdn&use_cors=true&cors_origin_url=play.sooplive.co.kr&broad_key=%s-common-master-hls&player_mode=landing", broad_no)

	req, _ := http.NewRequest("GET", url, nil)

	cookies := tools.LoadCookies("soop")

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:144.0) Gecko/20100101 Firefox/144.0")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "AuthTicket="+cookies[0])

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var server Server
	json.Unmarshal(body, &server)

	return server.ViewUrl
}

func GetMasterPlist(server, aid string) string {
	return fmt.Sprintf("%s?aid=%s", server, aid)
}

func GetStreamQualities(master string, quality int) string {
	res, _ := http.Get(master)

	playlists := make([]string, 0)

	scanner := bufio.NewScanner(res.Body)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "auth_playlist") {
			playlists = append(playlists, line)
		}
	}

	return playlists[quality]
}

func GetStream(server, playlist string) (string, string) {
	parsedUrl, _ := url.Parse(server)

	path := parsedUrl.ResolveReference(&url.URL{Path: "./"}).String()

	return path, path + playlist
}
