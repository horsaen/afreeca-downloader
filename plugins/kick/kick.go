package kick

import (
	"fmt"
	"horsaen/afreeca-downloader/tools"
)

func Start(userId string) {
	tools.Exists("downloads/Kick")

	if DvrCheck(userId) {
		fmt.Println("User online.")
		masterPlaylist := GetMaster(userId)

		playlist := ParseMaster(masterPlaylist)

		Download(playlist, userId)
	}
}
