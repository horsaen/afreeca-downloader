package flex

import (
	"encoding/json"
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"io"
	"log"
	"net/http"
	"time"
)

func CheckExists(userId string) bool {
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

	if user.Message == "채널 정보를 확인할 수 없습니다." {
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

	return len(user.Sources) != 0
}
