package flex

import (
	"fmt"
	"horsaen/afreeca-downloader/tools"
)

func Start(userId string) {
	tools.Exists("downloads/Flex")

	if DvrCheck(userId) {
		masterPlaylist, nickname, id := StreamData(userId)
		fmt.Println("User online.")
		url := ParseMaster(masterPlaylist)
		Download(url, nickname, id)
	}
}
