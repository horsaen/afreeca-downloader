package soop

import (
	"bufio"
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func ConcurrentDownload(user *[]string, broad_no string, playlist string) {
	bjId := (*user)[0]
	tools.Exists("downloads/Soop/" + bjId)

	client := &http.Client{}

	req, err := http.NewRequest("GET", playlist, nil)

	if err != nil {
		log.Fatal(err)
	}

	parsedUrl, _ := url.Parse(playlist)

	pathSegments := strings.Split(parsedUrl.Path, "/")

	newPath := strings.Join(pathSegments[:len(pathSegments)-1], "/")

	filename := bjId + "-" + broad_no + "-" + time.Now().Format("200601021504") + "-soop.ts"

	out, _ := os.Create("downloads/Soop/" + bjId + "/" + filename)

	playlistUrls := make(map[string]bool)

	var bytes int64 = 0
	var start_time = time.Now()

	for {
		resp, err := client.Do(req)

		if err != nil {
			(*user)[2] = "ERROR"
			(*user)[3] = "RETRYING"
			(*user)[4] = err.Error()
			ConcurrentDownload(user, broad_no, playlist)
		}

		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			fmt.Println(err)
		}

		bodyText := string(bodyBytes)

		if !strings.Contains(bodyText, ".TS") {
			(*user)[2] = "Offline"
			(*user)[3] = "Offline"
			(*user)[4] = "Offline"
			Concurrent(user)
		}

		scanner := bufio.NewScanner(strings.NewReader(bodyText))

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasSuffix(line, ".TS") {

				playlistUrl := fmt.Sprintf("%s://%s%s/%s", parsedUrl.Scheme, parsedUrl.Host, newPath, line)

				if !playlistUrls[playlistUrl] {

					resp, err := http.Get(playlistUrl)

					if err != nil {
						(*user)[2] = "ERROR"
						(*user)[3] = "RETRYING"
						(*user)[4] = err.Error()
						ConcurrentDownload(user, broad_no, playlist)
					}

					bytes += resp.ContentLength
					elapsed_time := time.Since(start_time)

					fmt.Printf("\rDownloading to %s || %s @ %s      \x1b[?25l", filename, tools.FormatTime(elapsed_time), tools.FormatBytes(bytes))

					_, err = io.Copy(out, resp.Body)

					if err != nil {
						(*user)[2] = "ERROR"
						(*user)[3] = "RETRYING"
						(*user)[4] = err.Error()
						ConcurrentDownload(user, broad_no, playlist)
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

func Download(bjId string, broad_no string, playlist string) {
	tools.Exists("downloads/Soop/" + bjId)

	client := &http.Client{}

	req, err := http.NewRequest("GET", playlist, nil)

	if err != nil {
		log.Fatal(err)
	}

	parsedUrl, _ := url.Parse(playlist)

	pathSegments := strings.Split(parsedUrl.Path, "/")

	newPath := strings.Join(pathSegments[:len(pathSegments)-1], "/")

	filename := bjId + "-" + broad_no + "-" + time.Now().Format("200601021504") + "-soop.ts"

	out, _ := os.Create("downloads/Soop/" + bjId + "/" + filename)

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

		if !strings.Contains(bodyText, ".TS") {
			Start(bjId)
		}

		scanner := bufio.NewScanner(strings.NewReader(bodyText))

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasSuffix(line, ".TS") {

				playlistUrl := fmt.Sprintf("%s://%s%s/%s", parsedUrl.Scheme, parsedUrl.Host, newPath, line)

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
