package kick

import (
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"os"
)

func Concurrent(user *[]string) {
	tools.Exists("downloads/Kick")

	if !CheckExists((*user)[0]) {
		(*user)[2] = "Not found"
		(*user)[3] = "Not found"
		(*user)[4] = "Not found"
		return
	}

	if ConcurrentCheck((*user)[0]) {
		masterPlaylist := GetMaster((*user)[0])

		playlist := ParseMaster(masterPlaylist)

		ConcurrentDownload(user, playlist)
	}
}

func Start(userId string) {
	tools.Exists("downloads/Kick")

	if !CheckExists(userId) {
		fmt.Printf("User %s not found.\n", userId)
		os.Exit(1)
	}

	if DvrCheck(userId) {
		fmt.Printf("User %s online.\n", userId)

		masterPlaylist := GetMaster(userId)

		playlist := ParseMaster(masterPlaylist)

		Download(playlist, userId)
	}
}
