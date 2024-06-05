package tiktok

import (
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"os"
)

func Start(userId string) {
	tools.Exists("downloads/Tiktok")

	roomId := RoomId(userId)

	if !CheckExists(roomId) {
		fmt.Printf("User %s not found.\n", userId)
		os.Exit(1)
	}

	if DvrCheck(roomId) {
		fmt.Printf("User %s online.\n", userId)

		mode, url := GetPlaylist(roomId)

		if mode == "hls" {
			Download(url, userId)
		} else {
			fmt.Println("No HLS Stream found, switching to FLV.")
			DownloadFlv(url, userId)
		}
	}
}
