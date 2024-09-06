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

	cookies := tools.LoadCookies("panda")

	req, _ := http.NewRequest("GET", playlist, nil)

	req.Header.Add("Cookie", "sessKey="+cookies[0])
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Origin", "https://www.pandalive.co.kr")
	req.Header.Add("DNT", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Referer", "https://www.pandalive.co.kr/")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-site")
	req.Header.Add("Sec-GPC", "1")
	req.Header.Add("TE", "trailers")

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

					req, _ := http.NewRequest("GET", line, nil)

					req.Header.Add("Cookie", "sessKey="+cookies[0])
					req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0")
					req.Header.Add("Accept", "application/json, text/plain, */*")
					req.Header.Add("Accept-Language", "en-US,en;q=0.5")
					req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					req.Header.Add("Origin", "https://www.pandalive.co.kr")
					req.Header.Add("DNT", "1")
					req.Header.Add("Connection", "keep-alive")
					req.Header.Add("Referer", "https://www.pandalive.co.kr/")
					req.Header.Add("Sec-Fetch-Dest", "empty")
					req.Header.Add("Sec-Fetch-Mode", "cors")
					req.Header.Add("Sec-Fetch-Site", "same-site")
					req.Header.Add("Sec-GPC", "1")
					req.Header.Add("TE", "trailers")

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

func Download(userId string, playlist string) {
	tools.Exists("downloads/Panda/" + userId)

	client := &http.Client{}

	cookies := tools.LoadCookies("panda")

	req, _ := http.NewRequest("GET", playlist, nil)

	req.Header.Add("Cookie", "sessKey="+cookies[0])
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Origin", "https://www.pandalive.co.kr")
	req.Header.Add("DNT", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Referer", "https://www.pandalive.co.kr/")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-site")
	req.Header.Add("Sec-GPC", "1")
	req.Header.Add("TE", "trailers")

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
					req, _ := http.NewRequest("GET", line, nil)

					req.Header.Add("Cookie", "sessKey="+cookies[0])
					req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0")
					req.Header.Add("Accept", "application/json, text/plain, */*")
					req.Header.Add("Accept-Language", "en-US,en;q=0.5")
					req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					req.Header.Add("Origin", "https://www.pandalive.co.kr")
					req.Header.Add("DNT", "1")
					req.Header.Add("Connection", "keep-alive")
					req.Header.Add("Referer", "https://www.pandalive.co.kr/")
					req.Header.Add("Sec-Fetch-Dest", "empty")
					req.Header.Add("Sec-Fetch-Mode", "cors")
					req.Header.Add("Sec-Fetch-Site", "same-site")
					req.Header.Add("Sec-GPC", "1")
					req.Header.Add("TE", "trailers")

					resp, err := client.Do(req)

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
