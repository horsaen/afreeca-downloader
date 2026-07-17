package tools

import (
	"bufio"
	"os"
	"sync"
)

var cookieMu sync.Mutex

func WriteCookies(cookie, platform string) {
	cookieMu.Lock()
	defer cookieMu.Unlock()

	home, _ := os.UserHomeDir()

	configBase := home + "/.afreeca-downloader"
	cookieDir := configBase + "/cookies/"

	if err := os.MkdirAll(cookieDir, os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.WriteFile(cookieDir+platform, []byte(cookie), 0644); err != nil {
		panic(err)
	}
}

func LoadCookies(platform string) [3]string {
	cookieMu.Lock()
	defer cookieMu.Unlock()

	home, _ := os.UserHomeDir()

	configBase := home + "/.afreeca-downloader"
	cookieDir := configBase + "/cookies/"

	file, _ := os.Open(cookieDir + "/" + platform)

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)
	cookies := [3]string{" ", " ", " "}

	index := 0
	for scanner.Scan() {
		cookies[index] = scanner.Text()
		index++
	}

	file.Close()

	return cookies
}
