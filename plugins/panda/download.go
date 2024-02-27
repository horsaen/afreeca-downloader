package panda

import (
	"bufio"
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func Download(userId string, playlist string) bool {
	tools.Exists("downloads/Panda/" + userId)

	client := &http.Client{}

	req, err := http.NewRequest("GET", playlist, nil)

	if err != nil {
		log.Fatal(err)
	}

	filename := userId + "-" + time.Now().Format("200601021504") + "-panda.ts"

	out, _ := os.Create("downloads/Panda/" + userId + "/" + filename)

	playlistUrls := make(map[string]bool)

	var bytes int64 = 0
	var start_time = time.Now()

	for {
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
		}

		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			fmt.Println(err)
		}

		bodyText := string(bodyBytes)

		if !strings.Contains(bodyText, ".ts") {
			Start(userId)
		}

		scanner := bufio.NewScanner(strings.NewReader(bodyText))

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "https://") {
				if !playlistUrls[line] {

					resp, err := http.Get(line)

					if err != nil {
						fmt.Println(err)
					}

					// since panda doesn't actually return a content-length header for go just use this instead of having to read the segments in ram
					filesize, _ := os.Stat("downloads/Panda/" + userId + "/" + filename)
					bytes = filesize.Size()
					elapsed_time := time.Since(start_time)

					fmt.Printf("\rDownloading to %s || %s @ %s      \x1b[?25l", filename, tools.FormatTime(elapsed_time), tools.FormatBytes(bytes))

					_, err = io.Copy(out, resp.Body)

					if err != nil {
						fmt.Println(err)
					}

					playlistUrls[line] = true
				}
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}

		time.Sleep(3 * time.Second)
	}
}
