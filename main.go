package main

import (
	"flag"
	"fmt"
	"horsaen/afreeca-downloader/plugins/afreeca"
	"horsaen/afreeca-downloader/plugins/bigo"
	"horsaen/afreeca-downloader/plugins/concurrent"
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
	playlist := flag.Bool("playlist", false, "Download bot playlists")
	concurrently := flag.Bool("concurrent", false, "Download streams concurrently")
	vod := flag.Bool("vod", false, "Download Afreeca vod")
	version := flag.Bool("version", false, "Print version")

	flag.Parse()

	if *version {
		tools.Version()
	}

	tools.Exists("downloads")

	if *concurrently {
		concurrent.Start()
	}

	if *userArg != "" || *playlist || *vod || *concurrently {
		username = *userArg
	} else {
		fmt.Println("Enter username:")
		fmt.Scan(&username)
	}

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
