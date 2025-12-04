package main

import (
	"flag"
	"fmt"
	"horsaen/afreeca-downloader/plugins/bigo"
	"horsaen/afreeca-downloader/plugins/chzzk"
	"horsaen/afreeca-downloader/plugins/flex"
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
	var username string

	mode := flag.String("mode", "tui", "Mode")
	userArg := flag.String("username", "", "Streamer username")
	concurrently := flag.Bool("concurrent", false, "Download streams concurrently")
	version := flag.Bool("version", false, "Print version")

	// soop stuff
	playlist := flag.Bool("playlist", false, "Download bot playlists")
	userVod := flag.Bool("uservod", false, "Download soop user vods")

	flag.Parse()

	if *version {
		tools.Version()
	}

	// if *concurrently {
	// 	concurrent.Start()
	// }

	if *userArg != "" || *playlist || *userVod || *concurrently || *mode == "tui" {
		username = *userArg
	} else {
		fmt.Println("Enter username:")
		fmt.Scan(&username)
	}

	switch *mode {
	case "soop":
		if *playlist {
			var playlist string
			fmt.Println("Enter playlist url:")
			fmt.Scan(&playlist)
			soop.DownloadPlaylist(playlist)
		} else if *userVod {
			soop.UserVods(username)
		} else {
			soop.Start(username)
		}
	case "bigo":
		bigo.Start(username)
	case "chzzk":
		chzzk.Start(username)
	case "flex":
		flex.Start(username)
	case "panda":
		panda.Start(username)
	case "tiktok":
		tiktok.Start(username)
	case "kick":
	case "afreeca":
		fmt.Println("DEPRICATED")
		os.Exit(1)
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
