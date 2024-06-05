package soop

import (
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"os"
	"os/exec"
	"time"
)

func Start(channelId string) {
	tools.Exists("downloads/Soop")

	if !CheckExists(channelId) {
		fmt.Printf("Channel %s not found.\n", channelId)
		os.Exit(1)
	}

	if DvrCheck(channelId) {
		fmt.Printf("Channel %s online.\n", channelId)

		masterPlaylist := "https://api.sooplive.com/media/live/" + channelId + "/master.m3u8"

		var start_time = time.Now()
		filename := channelId + "-" + time.Now().Format("200601021504") + "-soop.mp4"

		cmd := exec.Command("yt-dlp", "-o", "downloads/Soop/"+channelId+"/"+filename, masterPlaylist)

		go func() {
			for {
				elapsed_time := time.Since(start_time)
				filePath := "downloads/Soop/" + channelId + "/" + filename + ".part"

				filesize, err := os.Stat(filePath)

				if err != nil {
					continue
				}

				fmt.Printf("\rDownloading to %s || %s @ %s      \x1b[?25l", filename, tools.FormatTime(elapsed_time), tools.FormatBytes(filesize.Size()))
				time.Sleep(1 * time.Second)
			}
		}()

		cmd.Run()

		fmt.Println("User" + channelId + "offline.")
		Start(channelId)
	}
}
