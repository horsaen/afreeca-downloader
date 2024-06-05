package kick

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"io"
	"net/http"
	"time"
)

func CheckExists(userId string) bool {
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

	if res.StatusCode == 404 {
		return false
	} else {
		return true
	}
}

func ConcurrentCheck(userId string) bool {
	for {
		if CheckOnline(userId) {
			return true
		}
		time.Sleep(3 * time.Minute)
	}
}

func DvrCheck(userId string) bool {
	for {
		if CheckOnline(userId) {
			return true
		}
		fmt.Printf("User %s offline, rechecking in 3 minutes.\n", userId)
		time.Sleep(3 * time.Minute)
	}
}

func CheckOnline(userId string) bool {
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

	return user.Livestream.IsLive
}
