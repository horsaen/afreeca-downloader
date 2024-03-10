package afreeca

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func GetPlaylists() []string {
	inputPlaylist := []string{}

	fmt.Println("Paste playlist URLs:")

	scn := bufio.NewScanner(os.Stdin)
	for {
		scn.Scan()
		userInput := scn.Text()

		if userInput == "" {
			return inputPlaylist
		} else {
			inputPlaylist = append(inputPlaylist, userInput)
		}
	}
}

func ParsePlaylists(streams []string) []string {
	vodBase := "https://vod-archive-kr-cdn-z01.afreecatv.com"

	segments := []string{}

	for _, playlist := range streams {
		tempSegs := []string{}

		resp, _ := http.Get(playlist)

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasSuffix(line, ".m3u8") {
				tempSegs = append(tempSegs, line)
			}
		}

		segments = append(segments, vodBase+tempSegs[len(tempSegs)-1])
	}

	return segments
}
