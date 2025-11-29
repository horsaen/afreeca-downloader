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
	"slices"
	"strings"
	"time"
)

func Download(bjid, broad_no, base, stream string) {
	path := fmt.Sprintf("downloads/soop/%s", bjid)
	file := fmt.Sprintf("%s-%s-%s-soop.ts", bjid, broad_no, time.Now().Format("200601021504"))
	os.MkdirAll(path, os.ModePerm)
	out, err := os.Create(path + "/" + file)

	if err != nil {
		log.Fatal(err)
	}

	var prevSegments []string
	var bytes int64 = 0
	var start_time = time.Now()
	for {
		res, err := http.Get(stream)

		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(res.Body)

		currentSegments := []string{}

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasSuffix(line, ".TS") {
				currentSegments = append(currentSegments, line)
				if !slices.Contains(prevSegments, line) {
					res, err = http.Get(base + line)

					if err != nil {
						log.Fatal(err)
					}

					bytes += res.ContentLength
					elapsed_time := time.Since(start_time)

					fmt.Printf("\rDownloading to %s || %s @ %s      \x1b[?25l", file, tools.FormatTime(elapsed_time), tools.FormatBytes(bytes))

					io.Copy(out, res.Body)
				}
			}
		}

		if len(currentSegments) == 0 {
			fmt.Println("Stream has ended, resuming polling.")
			Start(bjid)
		}

		prevSegments = currentSegments
	}
}

func DownloadVod(manifest string) {
	m3uURL := strings.Replace(manifest, ".smil", ".mp4", -1)
	os.MkdirAll("downloads/soop", os.ModePerm)
	res, err := http.Get(m3uURL)

	if err != nil {
		log.Fatal(err)
	}

	base, _ := url.Parse(m3uURL)
	sc := bufio.NewScanner(res.Body)

	var initURL string
	var segments []string

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())

		if strings.HasPrefix(line, "#EXT-X-MAP:") {
			uri := extractQuotedValue(line)
			full, _ := base.Parse(uri)
			initURL = full.String()
			continue
		}

		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		full, _ := base.Parse(line)
		segments = append(segments, full.String())
	}

	if initURL == "" {
		fmt.Println("playlist missing EXT-X-MAP (init.m4s)")
		return
	}

	out, err := os.Create("downloads/soop/" + time.Now().String() + ".mp4")

	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	fmt.Println("Downloading init.m4s")
	res, _ = http.Get(initURL)
	io.Copy(out, res.Body)

	for i, seg := range segments {
		res, err = http.Get(seg)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\rDownloading segment %d/%d", i+1, len(segments))
		io.Copy(out, res.Body)
	}

}

// helper slop
func extractQuotedValue(s string) string {
	start := strings.Index(s, "\"")
	end := strings.LastIndex(s, "\"")
	if start == -1 || end == -1 || end <= start {
		return ""
	}
	return s[start+1 : end]
}
