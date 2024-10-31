package main

import (
	"flag"
	"fmt"
	"horsaen/afreeca-downloader/plugins/bigo"
	"horsaen/afreeca-downloader/plugins/chzzk"
	"horsaen/afreeca-downloader/plugins/concurrent"
	"horsaen/afreeca-downloader/plugins/flex"
	"horsaen/afreeca-downloader/plugins/kick"
	"horsaen/afreeca-downloader/plugins/panda"
	"horsaen/afreeca-downloader/plugins/soop"
	"horsaen/afreeca-downloader/plugins/tiktok"
	"horsaen/afreeca-downloader/tools"
	"horsaen/afreeca-downloader/views"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	tools.InitConfDir()
	tools.Exists("downloads")

	var username string

	mode := flag.String("mode", "tui", "Mode")
	userArg := flag.String("username", "", "Streamer username")
	playlist := flag.Bool("playlist", false, "Download bot playlists")
	concurrently := flag.Bool("concurrent", false, "Download streams concurrently")
	vod := flag.Bool("vod", false, "Download Afreeca vod")
	version := flag.Bool("version", false, "Print version")

	flag.Parse()

	if *version {
		tools.Version()
	}

	if *concurrently {
		concurrent.Start()
	}

	if *userArg != "" || *playlist || *vod || *concurrently || *mode == "tui" {
		username = *userArg
	} else {
		fmt.Println("Enter username:")
		fmt.Scan(&username)
	}

	switch *mode {
	// case "afreeca":
	// 	if *playlist {
	// 		inputPlaylists := afreeca.GetPlaylists()
	// 		afreeca.Playlist(inputPlaylists)
	// 	} else if *vod {
	// 		var TitleNo string
	// 		fmt.Println("Enter title number:")
	// 		fmt.Scan(&TitleNo)
	// 		afreeca.Vod(TitleNo)
	// 	} else {
	// 		afreeca.Start(username)
	// 	}
	case "soop":
		soop.Start(username)
	case "bigo":
		bigo.Start(username)
	case "chzzk":
		chzzk.Start(username)
	case "flex":
		flex.Start(username)
	case "kick":
		kick.Start(username)
	case "panda":
		panda.Start(username)
	case "tiktok":
		tiktok.Start(username)
	case "tui":
		p := tea.NewProgram(views.InitialModel())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("Mode not supported, exiting.")
		os.Exit(1)
	}
}
