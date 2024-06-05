package flex

import (
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"os"
)

func Concurrent(user *[]string) {
	tools.Exists("downloads/Flex")

	if !CheckExists((*user)[0]) {
		(*user)[2] = "Not found"
		(*user)[3] = "Not found"
		(*user)[4] = "Not found"
		return
	}

	if ConcurrentCheck((*user)[0]) {
		masterPlaylist, nickname, id := StreamData((*user)[0])

		url := ParseMaster(masterPlaylist)

		ConcurrentDownload(user, url, nickname, id)
	}
}

func Start(userId string) {
	tools.Exists("downloads/Flex")

	if !CheckExists(userId) {
		fmt.Printf("User %s not found.\n", userId)
		os.Exit(1)
	}

	if DvrCheck(userId) {
		fmt.Printf("User %s online.", userId)

		masterPlaylist, nickname, id := StreamData(userId)

		url := ParseMaster(masterPlaylist)

		Download(url, nickname, id)
	}
}
