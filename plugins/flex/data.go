package flex

import (
	"bufio"
	"encoding/json"
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"io"
	"log"
	"net/http"
	"strings"
)

type User struct {
	Id    int `json:"id"`
	Owner struct {
		LoginId  string `json:"loginId"`
		Nickname string `json:"nickname"`
	} `json:"owner"`
	Sources []struct {
		Resolution int    `json:"resolution"`
		URL        string `json:"url"`
	}
}

func StreamData(userId string) (string, string, string) {
	client := &http.Client{}

	url := "https://api.flextv.co.kr/api/channels/" + userId + "/stream?option=all"
	cookies := tools.LoadCookies("flex")

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("Cookie", "flx_oauth_access="+cookies[0]+";")

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

	return user.Sources[0].URL, user.Owner.Nickname, user.Owner.LoginId
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
