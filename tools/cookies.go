package tools

import (
	"bufio"
	"fmt"
	"os"
)

func LoadCookies(platform string) []string {
	home, _ := os.UserHomeDir()

	configBase := home + "/.afreeca-downloader"
	cookieDir := configBase + "/cookies/"

	file, _ := os.Open(cookieDir + "/" + platform)

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)
	var cookies []string

	for scanner.Scan() {
		cookies = append(cookies, scanner.Text())
	}

	file.Close()

	for _, line := range cookies {
		fmt.Println(line)
	}

	fmt.Println(cookies)

	return cookies

}
