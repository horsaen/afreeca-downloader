package main

import (
	"flag"
	"fmt"
	"horsaen/afreeca-downloader/plugins/afreeca"
	"horsaen/afreeca-downloader/plugins/bigo"
	"horsaen/afreeca-downloader/plugins/flex"
	"horsaen/afreeca-downloader/plugins/kick"
	"horsaen/afreeca-downloader/plugins/panda"
	"horsaen/afreeca-downloader/plugins/tiktok"
	"horsaen/afreeca-downloader/tools"
	"os"
)

func main() {
	tools.InitConfDir()

	var username string

	mode := flag.String("mode", "afreeca", "Mode")
	userArg := flag.String("username", "", "Streamer username")
	playlist := flag.Bool("playlist", false, "Download playlist")
	vod := flag.Bool("vod", false, "Download vod")
	version := flag.Bool("version", false, "Print version")

	flag.Parse()

	if *version {
		fmt.Println("https://github.com/horsaen/afreeca-downloader")
		fmt.Println("v2.0.2")
		os.Exit(0)
	}

	if *userArg != "" || *playlist || *vod {
		username = *userArg
	} else {
		fmt.Println("Enter username:")
		fmt.Scan(&username)
	}

	tools.Exists("downloads")

	switch *mode {
	case "afreeca":
		tools.Exists("downloads/Afreeca")
		if *playlist {
			afreeca.Playlist()
		} else if *vod {
			var TitleNo string
			fmt.Println("Enter title number:")
			fmt.Scan(&TitleNo)
			afreeca.Vod(TitleNo)
		} else {
			afreeca.Start(username)
		}
	case "bigo":
		bigo.Start(username)
	case "flex":
		flex.Start(username)
	case "kick":
		kick.Start(username)
	case "panda":
		panda.Start(username)
	case "tiktok":
		tiktok.Start(username)
	default:
		fmt.Println("Mode not supported, exiting.")
		os.Exit(1)
	}
}
