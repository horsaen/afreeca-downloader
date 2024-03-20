package bigo

import (
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"os"
)

func Concurrent(user *[]string) {
	tools.Exists("downloads/Bigo")

	if !CheckExists((*user)[0]) {
		(*user)[2] = "Not found"
		(*user)[3] = "Not found"
		(*user)[4] = "Not found"
		return
	}

	if ConcurrentCheck((*user)[0]) {
		_, nickname, url := GetStreamData((*user)[0])
		ConcurrentDownload(user, nickname, url)
	}
}

func Start(userId string) {
	tools.Exists("downloads/Bigo")

	if !CheckExists(userId) {
		fmt.Printf("User %s not found.\n", userId)
		os.Exit(1)
	}

	if DvrCheck(userId) {
		fmt.Println("User online.")
		_, nickname, url := GetStreamData(userId)
		Download(url, nickname, userId)
	}
}
