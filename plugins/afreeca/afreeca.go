package afreeca

import (
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"os"
	"strconv"
)

func Concurrent(user *[]string) {
	tools.Exists("downloads/Afreeca")

	if !CheckExists((*user)[0]) {
		(*user)[2] = "Not found"
		(*user)[3] = "Not found"
		(*user)[4] = "Not found"
		return
	}

	if ConcurrentCheck((*user)[0]) {
		nickname, broad_no := UserData((*user)[0])

		url := GetStream((*user)[0], broad_no, "")

		ConcurrentDownload(user, nickname, strconv.Itoa(broad_no), url)
	}
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

		Download(bjId, nickname, strconv.Itoa(broad_no), url)
	}
}

func Playlist(inputPlaylists []string) {
	tools.Exists("downloads/Afreeca/Playlists")

	playlists := ParsePlaylists(inputPlaylists)

	DownloadPlaylists(playlists)
}

func Vod(TitleNo string) {
	tools.Exists("downloads/Afreeca/Vods")

	sourceFiles := GetVodFiles(TitleNo)

	files := ParseVods(sourceFiles)

	DownloadVods(TitleNo, files)
}
