package afreeca

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
	"strconv"
	"strings"
	"time"
)

func Download(bjId string, nickname string, playlist string) bool {
	tools.Exists("downloads/Afreeca/" + bjId)

	client := &http.Client{}

	req, err := http.NewRequest("GET", playlist, nil)

	if err != nil {
		log.Fatal(err)
	}

	parsedUrl, _ := url.Parse(playlist)

	pathSegments := strings.Split(parsedUrl.Path, "/")

	newPath := strings.Join(pathSegments[:len(pathSegments)-1], "/")

	filename := nickname + "-" + bjId + "-" + time.Now().Format("200601021504") + "-afreeca.ts"

	out, _ := os.Create("downloads/Afreeca/" + bjId + "/" + filename)

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

func DownloadPlaylists(playlists []string) {
	now := time.Now().Format("200601021504")
	tools.Exists("downloads/Afreeca/" + now)

	vodBase := "https://vod-archive-kr-cdn-z01.afreecatv.com"
	length := len(playlists)

	for i, playlist := range playlists {
		var downloaded = 0
		segments := []string{}

		parseUrl, _ := url.Parse(playlist)
		resp, _ := http.Get(playlist)

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasSuffix(line, ".ts") {
				segments = append(segments, vodBase+path.Dir(parseUrl.Path)+"/"+line)
			}
		}

		filename := "downloads/Afreeca/" + now + "/" + strconv.Itoa(i+1) + ".ts"
		out, _ := os.Create(filename)

		for _, segment := range segments {
			resp, _ := http.Get(segment)
			downloaded += 1
			output := fmt.Sprintf("\rDownloaded %d segments out of %d || Playlist %d of %d", downloaded, len(segments), i+1, length)
			fmt.Print(output)
			io.Copy(out, resp.Body)
		}

	}
}
