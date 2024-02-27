package bigo

import (
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"os"
)

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
