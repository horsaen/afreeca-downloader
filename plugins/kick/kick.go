package kick

import (
	"fmt"
	"horsaen/afreeca-downloader/tools"
)

func Concurrent(user *[]string) {
	tools.Exists("downloads/Kick")

	if ConcurrentCheck((*user)[0]) {
		masterPlaylist := GetMaster((*user)[0])

		playlist := ParseMaster(masterPlaylist)

		ConcurrentDownload(user, playlist)
	}
}

func Start(userId string) {
	tools.Exists("downloads/Kick")

	if DvrCheck(userId) {
		fmt.Println("User online.")
		masterPlaylist := GetMaster(userId)

		playlist := ParseMaster(masterPlaylist)

		Download(playlist, userId)
	}
}
