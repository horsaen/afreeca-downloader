package panda

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"io"
	"log"
	"net/http"
	"strings"
)

type User struct {
	Message   string `json:"message"`
	ErrorData struct {
		Code string `json:"code"`
	} `json:"errorData"`
	Playlist struct {
		Hls []struct {
			Url string `json:"url"`
		} `json:"hls"`
	} `json:"PlayList"`
}

func GetMaster(userId string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	url := "https://api.pandalive.co.kr/v1/live/play"
	cookies := tools.LoadCookies("panda")

	payload := strings.NewReader("action=watch&userId=" + userId + "&password=&width=48&height=48&imageResize=crop&fanIconWidth=44&fanIconHeight=44&fanIconImageResize=crop")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Cookie", "sessKey="+cookies[0])
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Origin", "https://www.pandalive.co.kr")
	req.Header.Add("DNT", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Referer", "https://www.pandalive.co.kr/")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-site")
	req.Header.Add("Sec-GPC", "1")
	req.Header.Add("TE", "trailers")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	bodyText, _ := io.ReadAll(resp.Body)

	var user User

	json.Unmarshal(bodyText, &user)

	return user.Playlist.Hls[0].Url
}

func ParseMaster(masterPlaylist string) string {
	client := &http.Client{}

	cookies := tools.LoadCookies("panda")

	req, _ := http.NewRequest("GET", masterPlaylist, nil)

	req.Header.Add("Cookie", "sessKey="+cookies[0])
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Origin", "https://www.pandalive.co.kr")
	req.Header.Add("DNT", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Referer", "https://www.pandalive.co.kr/")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-site")
	req.Header.Add("Sec-GPC", "1")
	req.Header.Add("TE", "trailers")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	playlists := make([]string, 0)

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "https://") {
			playlists = append(playlists, line)
		}
	}

	return playlists[0]
}
