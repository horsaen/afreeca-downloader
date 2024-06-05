package bigo

import (
	"bufio"
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func ConcurrentDownload(user *[]string, nickname string, playlist string) {
	userId := (*user)[0]
	tools.Exists("downloads/Bigo/" + userId)

	client := &http.Client{}

	req, _ := http.NewRequest("GET", playlist, nil)

	parsedUrl, _ := url.Parse(playlist)

	filename := nickname + "-" + userId + "-" + time.Now().Format("200601021504") + "-bigo.ts"

	out, _ := os.Create("downloads/Bigo/" + userId + "/" + filename)

	playlistUrls := make(map[string]bool)

	var bytes int64 = 0
	var start_time = time.Now()

	for {
		resp, err := client.Do(req)

		if err != nil {
			(*user)[2] = "ERROR"
			(*user)[3] = "RETRYING"
			(*user)[4] = err.Error()
			ConcurrentDownload(user, nickname, playlist)
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
			if strings.HasSuffix(line, ".ts") {

				playlistUrl := fmt.Sprintf("%s://%s/%s", parsedUrl.Scheme, parsedUrl.Host, line)

				if !playlistUrls[playlistUrl] {

					resp, err := http.Get(playlistUrl)

					if err != nil {
						(*user)[2] = "ERROR"
						(*user)[3] = "RETRYING"
						(*user)[4] = err.Error()
						ConcurrentDownload(user, nickname, playlist)
					}

					bytes += resp.ContentLength
					elapsed_time := time.Since(start_time)

					(*user)[2] = tools.FormatBytes(bytes)
					(*user)[3] = tools.FormatTime(elapsed_time)
					(*user)[4] = filename

					_, err = io.Copy(out, resp.Body)

					if err != nil {
						(*user)[2] = "ERROR"
						(*user)[3] = "RETRYING"
						(*user)[4] = err.Error()
						ConcurrentDownload(user, nickname, playlist)
					}

					playlistUrls[playlistUrl] = true
				}
			}
		}

		time.Sleep(3 * time.Second)
	}
}

func Download(playlist string, nickname string, userId string) bool {
	tools.Exists("downloads/Bigo/" + userId)

	client := &http.Client{}

	req, _ := http.NewRequest("GET", playlist, nil)

	parsedUrl, _ := url.Parse(playlist)

	filename := nickname + "-" + userId + "-" + time.Now().Format("200601021504") + "-bigo.ts"

	out, _ := os.Create("downloads/Bigo/" + userId + "/" + filename)

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
			if strings.HasSuffix(line, ".ts") {

				playlistUrl := fmt.Sprintf("%s://%s/%s", parsedUrl.Scheme, parsedUrl.Host, line)

				if !playlistUrls[playlistUrl] {

					resp, err := http.Get(playlistUrl)

					if err != nil {
						fmt.Println(err)
					}

					bytes += resp.ContentLength
					elapsed_time := time.Since(start_time)

					fmt.Printf("\rDownloading to %s || %s @ %s      \x1b[?25l", filename, tools.FormatTime(elapsed_time), tools.FormatBytes(bytes))

					_, err = io.Copy(out, resp.Body)

					if err != nil {
						fmt.Println(err)
					}

					playlistUrls[playlistUrl] = true
				}
			}
		}

		time.Sleep(3 * time.Second)
	}
}
