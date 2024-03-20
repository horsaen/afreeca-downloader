package flex

import (
	"fmt"
	"horsaen/afreeca-downloader/tools"
)

func Concurrent(user *[]string) {
	tools.Exists("downloads/Flex")

	if ConcurrentCheck((*user)[0]) {
		masterPlaylist, nickname, id := StreamData((*user)[0])

		url := ParseMaster(masterPlaylist)

		ConcurrentDownload(user, url, nickname, id)
	}
}

func Start(userId string) {
	tools.Exists("downloads/Flex")

	if DvrCheck(userId) {
		masterPlaylist, nickname, id := StreamData(userId)
		fmt.Println("User online.")
		url := ParseMaster(masterPlaylist)
		Download(url, nickname, id)
	}
}
