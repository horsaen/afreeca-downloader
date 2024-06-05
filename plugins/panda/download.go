package panda

import (
	"bufio"
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func ConcurrentDownload(user *[]string, playlist string) {
	userId := (*user)[0]
	tools.Exists("downloads/Panda/" + userId)

	client := &http.Client{}

	req, _ := http.NewRequest("GET", playlist, nil)

	filename := userId + "-" + time.Now().Format("200601021504") + "-panda.ts"

	out, _ := os.Create("downloads/Panda/" + userId + "/" + filename)

	playlistUrls := make(map[string]bool)

	var bytes int64 = 0
	var start_time = time.Now()

	for {
		resp, err := client.Do(req)

		if err != nil {
			(*user)[2] = "ERROR"
			(*user)[3] = "RETRYING"
			(*user)[4] = err.Error()
			ConcurrentDownload(user, playlist)
		}

		bodyBytes, _ := io.ReadAll(resp.Body)

		bodyText := string(bodyBytes)

		if !strings.Contains(bodyText, ".ts") {
			(*user)[2] = "Offline"
			(*user)[3] = "Offline"
			(*user)[4] = "Offline"
			Concurrent(user)
		}

		scanner := bufio.NewScanner(strings.NewReader(bodyText))

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "https://") {
				if !playlistUrls[line] {

					resp, err := http.Get(line)

					if err != nil {
						(*user)[2] = "ERROR"
						(*user)[3] = "RETRYING"
						(*user)[4] = err.Error()
						ConcurrentDownload(user, playlist)
					}

					// since panda doesn't actually return a content-length header for go just use this instead of having to read the segments in ram
					filesize, _ := os.Stat("downloads/Panda/" + userId + "/" + filename)
					bytes = filesize.Size()
					elapsed_time := time.Since(start_time)

					(*user)[2] = tools.FormatBytes(bytes)
					(*user)[3] = tools.FormatTime(elapsed_time)
					(*user)[4] = filename

					_, err = io.Copy(out, resp.Body)

					if err != nil {
						(*user)[2] = "ERROR"
						(*user)[3] = "RETRYING"
						(*user)[4] = err.Error()
						ConcurrentDownload(user, playlist)
					}

					playlistUrls[line] = true
				}
			}
		}

		time.Sleep(3 * time.Second)
	}
}

func Download(userId string, playlist string) bool {
	tools.Exists("downloads/Panda/" + userId)

	client := &http.Client{}

	req, _ := http.NewRequest("GET", playlist, nil)

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

		bodyBytes, _ := io.ReadAll(resp.Body)

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

		time.Sleep(3 * time.Second)
	}
}
