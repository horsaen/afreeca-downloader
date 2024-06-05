package chzzk

import (
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"os"
	"os/exec"
	"time"
)

func Start(userId string) {
	tools.Exists("downloads/Chzzk")
	if DvrCheck(userId) {
		fmt.Println("User online.")
		username, url := GetData(userId)
		var start_time = time.Now()
		filename := username + "-" + userId + "-" + time.Now().Format("200601021504") + "-chzzk.mp4"

		cmd := exec.Command("yt-dlp", "-o", "downloads/Chzzk/"+username+"/"+filename, url)

		go func() {
			for {
				elapsed_time := time.Since(start_time)
				filePath := "downloads/Chzzk/" + username + "/" + filename + ".part"
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					continue
				}
				filesize, _ := os.Stat(filePath)

				fmt.Printf("\rDownloading to %s || %s @ %s      \x1b[?25l", filename, tools.FormatTime(elapsed_time), tools.FormatBytes(filesize.Size()))
				time.Sleep(1 * time.Second)
			}
		}()

		cmd.Run()

		fmt.Println("User" + userId + "offline.")
		Start(userId)
	}
}
