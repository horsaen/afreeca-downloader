package tiktok

import (
	"fmt"
	"horsaen/afreeca-downloader/tools"
)

func Start(userId string) {
	tools.Exists("downloads/Tiktok")

	roomId := RoomId(userId)

	if DvrCheck(roomId) {
		fmt.Println("User online.")

		url := GetPlaylist(roomId)

		Download(url, userId)
	}
}
