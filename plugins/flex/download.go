package flex

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

func ConcurrentDownload(user *[]string, playlist string, nickname string, userId string) {
	tools.Exists("downloads/Flex/" + userId)

	tr := http.Transport{
		DisableCompression: true,
	}

	client := &http.Client{Transport: &tr}

	req, _ := http.NewRequest("GET", playlist, nil)

	filename := nickname + "-" + userId + "-" + time.Now().Format("200601021504") + "-flex.ts"

	out, _ := os.Create("downloads/Flex/" + userId + "/" + filename)

	playlistUrls := make(map[string]bool)

	var bytes int64 = 0
	var start_time = time.Now()

	for {
		resp, err := client.Do(req)

		if err != nil {
			(*user)[2] = "error"
			(*user)[3] = "error"
			(*user)[4] = err.Error()
			return
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
						(*user)[2] = "error"
						(*user)[3] = "error"
						(*user)[4] = err.Error()
						return
					}

					// since flex doesn't actually return a content-length header for go just use this instead of having to read the segments in ram
					filesize, _ := os.Stat("downloads/Flex/" + userId + "/" + filename)
					bytes = filesize.Size()
					elapsed_time := time.Since(start_time)

					(*user)[2] = tools.FormatBytes(bytes)
					(*user)[3] = tools.FormatTime(elapsed_time)
					(*user)[4] = filename

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

func Download(playlist string, nickname string, userId string) bool {
	tools.Exists("downloads/Flex/" + userId)

	tr := http.Transport{
		DisableCompression: true,
	}

	client := &http.Client{Transport: &tr}

	req, err := http.NewRequest("GET", playlist, nil)

	if err != nil {
		log.Fatal(err)
	}

	filename := nickname + "-" + userId + "-" + time.Now().Format("200601021504") + "-flex.ts"

	out, _ := os.Create("downloads/Flex/" + userId + "/" + filename)

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

					// since flex doesn't actually return a content-length header for go just use this instead of having to read the segments in ram
					filesize, _ := os.Stat("downloads/Flex/" + userId + "/" + filename)
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
