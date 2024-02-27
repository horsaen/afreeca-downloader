package panda

import (
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"os"
)

func Start(userId string) {
	tools.Exists("downloads/Panda")

	if !CheckExists(userId) {
		fmt.Printf("User %s not found.\n", userId)
		os.Exit(1)
	}

	if DvrCheck(userId) {
		fmt.Println("User online.")

		masterPlaylist := GetMaster(userId)

		url := ParseMaster(masterPlaylist)

		Download(userId, url)
	}
}
