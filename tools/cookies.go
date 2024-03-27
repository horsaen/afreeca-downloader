package tools

import (
	"bufio"
	"os"
)

func LoadCookies(platform string) [3]string {
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
