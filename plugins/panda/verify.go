package panda

import (
	"encoding/json"
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func CheckExists(userId string) bool {
	client := &http.Client{}

	url := "https://api.pandalive.co.kr/v1/member/bj"
	cookies := tools.LoadCookies("panda")

	payload := strings.NewReader("userId=" + userId + "&info=media%20fanGrade")

	req, err := http.NewRequest("POST", url, payload)

	req.Header.Add("Cookie", "sessKey="+cookies[0])
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		log.Fatal(err)
	}

	resp, _ := client.Do(req)

	bodyText, _ := io.ReadAll(resp.Body)

	var user User

	json.Unmarshal(bodyText, &user)

	if user.Message == "유저 정보가 없습니다." {
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

	url := "https://api.pandalive.co.kr/v1/live/play"
	cookies := tools.LoadCookies("panda")

	payload := strings.NewReader("action=watch&userId=" + userId + "&password=&width=48&height=48&imageResize=crop&fanIconWidth=44&fanIconHeight=44&fanIconImageResize=crop")

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "sessKey="+cookies[0])

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	bodyText, _ := io.ReadAll(resp.Body)

	var user User

	json.Unmarshal(bodyText, &user)

	if user.ErrorData.Code == "needAdult" {
		fmt.Println("Stream is 19+ and unable to retrieve stream URL, input a valid sessKey in panda-cookies.")
		os.Exit(2)
	}

	if user.ErrorData.Code == "castEnd" {
		return false
	}

	return true
}
