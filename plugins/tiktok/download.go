package tiktok

import (
	"bufio"
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

func Download(playlist string, userId string) bool {
	tools.Exists("downloads/Tiktok/" + userId)

	client := &http.Client{}

	req, err := http.NewRequest("GET", playlist, nil)

	if err != nil {
		log.Fatal(err)
	}

	parsedUrl, _ := url.Parse(playlist)

	filename := userId + "-" + time.Now().Format("200601021504") + "-tiktok.ts"

	out, _ := os.Create("downloads/Tiktok/" + userId + "/" + filename)

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
			if strings.HasSuffix(line, ".ts") {

				playlistUrl := fmt.Sprintf("%s://%s%s/%s", parsedUrl.Scheme, parsedUrl.Host, path.Dir(parsedUrl.Path), line)

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

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}

		time.Sleep(3 * time.Second)
	}
}

func DownloadFlv(url string, userId string) {
	tools.Exists("downloads/Tiktok/" + userId)

	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)

	var start_time = time.Now()
	filename := userId + "-" + time.Now().Format("200601021504") + "-tiktok.flv"

	out, _ := os.Create("downloads/Tiktok/" + userId + "/" + filename)

	go func() {
		for {
			elapsed_time := time.Since(start_time)
			filePath := "downloads/Tiktok/" + userId + "/" + filename
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				continue
			}
			filesize, _ := os.Stat(filePath)

			fmt.Printf("\rDownloading to %s || %s @ %s      \x1b[?25l", filename, tools.FormatTime(elapsed_time), tools.FormatBytes(filesize.Size()))
			time.Sleep(1 * time.Second)
		}
	}()

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("FLV download complete.")

	Start(userId)
}
