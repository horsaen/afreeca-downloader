package soop

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func CheckOnline(bjid string) bool {
	res, _ := http.Get("https://api-channel.sooplive.co.kr/v1.1/channel/" + bjid + "/home/section/broad")

	body, _ := io.ReadAll(res.Body)

	return len(body) > 0
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

// todo
// func CheckExists(bjid string) bool {

// }
