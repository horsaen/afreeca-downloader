package bigo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
)

func CheckExists(userId string) bool {
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

	if response.Data.Sid != 0 {
		return true
	} else {
		return false
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

	if response.Data.HlsSrc != "" {
		return true
	} else {
		return false
	}
}
