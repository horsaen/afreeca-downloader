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

func ConcurrentDownload(user *[]string, nickname string, broad_no string, playlist string) {
	bjId := (*user)[0]
	tools.Exists("downloads/Afreeca/" + bjId)

	client := &http.Client{}

	req, _ := http.NewRequest("GET", playlist, nil)

	parsedUrl, _ := url.Parse(playlist)

	pathSegments := strings.Split(parsedUrl.Path, "/")

	newPath := strings.Join(pathSegments[:len(pathSegments)-1], "/")

	filename := nickname + "-" + bjId + "-" + broad_no + "-" + time.Now().Format("200601021504") + "-afreeca.ts"

	out, _ := os.Create("downloads/Afreeca/" + bjId + "/" + filename)

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
						(*user)[2] = "error"
						(*user)[3] = "error"
						(*user)[4] = err.Error()
						return
					}

					bytes += resp.ContentLength
					elapsed_time := time.Since(start_time)

					(*user)[2] = tools.FormatBytes(bytes)
					(*user)[3] = tools.FormatTime(elapsed_time)
					(*user)[4] = filename

					_, err = io.Copy(out, resp.Body)

					if err != nil {
						(*user)[2] = "error"
						(*user)[3] = "error"
						(*user)[4] = err.Error()
						return
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

func Download(bjId string, nickname string, broad_no string, playlist string) bool {
	tools.Exists("downloads/Afreeca/" + bjId)

	client := &http.Client{}

	req, err := http.NewRequest("GET", playlist, nil)

	if err != nil {
		log.Fatal(err)
	}

	parsedUrl, _ := url.Parse(playlist)

	pathSegments := strings.Split(parsedUrl.Path, "/")

	newPath := strings.Join(pathSegments[:len(pathSegments)-1], "/")

	filename := nickname + "-" + bjId + "-" + broad_no + "-" + time.Now().Format("200601021504") + "-afreeca.ts"

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
	tools.Exists("downloads/Afreeca/Playlists/" + now)

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

		filename := "downloads/Afreeca/Playlists/" + now + "/" + strconv.Itoa(i+1) + ".ts"
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

func DownloadVods(TitleNo string, files []string) {
	tools.Exists("downloads/Afreeca/Vods/" + TitleNo)

	vodBase := "https://vod-archive-global-cdn-z02.afreecatv.com"
	length := len(files)

	for i, file := range files {
		var downloaded = 0
		segments := []string{}

		parseUrl, _ := url.Parse(file)
		resp, _ := http.Get(file)

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasSuffix(line, ".ts") {
				segments = append(segments, vodBase+path.Dir(parseUrl.Path)+"/"+line)
			}
		}

		filename := "downloads/Afreeca/Vods/" + TitleNo + "/" + strconv.Itoa(i+1) + ".ts"
		out, err := os.Create(filename)

		if err != nil {
			log.Fatal(err)
		}

		for _, segment := range segments {
			resp, _ := http.Get(segment)
			downloaded += 1
			output := fmt.Sprintf("\rDownloaded %d segments out of %d || File %d of %d", downloaded, len(segments), i+1, length)
			fmt.Print(output)
			io.Copy(out, resp.Body)
		}
	}

}
