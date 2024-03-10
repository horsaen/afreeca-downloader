package afreeca

import (
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"os"
)

func Playlist() {
	inputPlaylists := GetPlaylists()

	playlists := ParsePlaylists(inputPlaylists)

	DownloadPlaylists(playlists)
}

func Start(bjId string) {
	tools.Exists("downloads/Afreeca")

	if !CheckExists(bjId) {
		fmt.Printf("User %s not found.\n", bjId)
		os.Exit(1)
	}

	if DvrCheck(bjId) {
		nickname, broad_no := UserData(bjId)

		url := GetStream(bjId, broad_no, "")

		Download(bjId, nickname, url)
	}
}
