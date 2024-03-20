package afreeca

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func CheckExists(bjId string) bool {
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

	if user.Code == 9000 {
		return false
	} else {
		return true
	}
}

func ConcurrentCheck(bjId string) bool {
	for {
		if CheckOnline(bjId) {
			return true
		}
		time.Sleep(3 * time.Minute)
	}
}

func DvrCheck(bjId string) bool {
	for {
		if CheckOnline(bjId) {
			return true
		}
		fmt.Printf("User %s offline, rechecking in 3 minutes.\n", bjId)
		time.Sleep(3 * time.Minute)
	}
}

func CheckOnline(bjId string) bool {
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

	if user.Broad.BroadNo != 0 {
		return true
	} else {
		return false
	}
}
