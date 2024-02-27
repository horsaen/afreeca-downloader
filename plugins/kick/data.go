package kick

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"io"
	"net/http"
	"strings"
)

type User struct {
	PlaybackUrl string `json:"playback_url"`
	Livestream  struct {
		IsLive bool `json:"is_live"`
	} `json:"livestream"`
}

func GetMaster(userId string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	url := "https://kick.com/api/v1/channels/" + userId
	cookies := tools.LoadCookies("kick")

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Cookie", "__cf_bm="+cookies[0]+"; cf_clearance="+cookies[1]+"; kick_session="+cookies[2]+";")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Language", "en-US")
	req.Header.Add("DNT", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-GPC", "1")
	req.Header.Add("TE", "trailers")

	res, _ := client.Do(req)

	bodyText, _ := io.ReadAll(res.Body)

	var user User

	json.Unmarshal(bodyText, &user)

	return user.PlaybackUrl
}

func ParseMaster(masterPlaylist string) string {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", masterPlaylist, nil)

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
