package afreeca

import (
	"bufio"
	"net/http"
	"strings"
)

func ParseVods(files []string) []string {
	vodBase := "https://vod-archive-global-cdn-z02.afreecatv.com"

	fileArr := []string{}

	for _, file := range files {
		qualities := []string{}

		resp, _ := http.Get(file)

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasSuffix(line, ".m3u8") {
				qualities = append(qualities, line)
			}
		}

		fileArr = append(fileArr, vodBase+qualities[len(qualities)-1])
	}

	return fileArr
}
