package soop

import (
	"bufio"
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"io"
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
	out, _ := os.Create(path + "/" + file)

	var prevSegments []string
	var bytes int64 = 0
	var start_time = time.Now()
	for {
		res, _ := http.Get(stream)

		scanner := bufio.NewScanner(res.Body)

		currentSegments := []string{}

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasSuffix(line, ".TS") {
				currentSegments = append(currentSegments, line)
				if !slices.Contains(prevSegments, line) {
					res, _ = http.Get(base + line)

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
	res, _ := http.Get(m3uURL)

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

	out, _ := os.Create("downloads/soop/" + time.Now().String() + ".mp4")
	defer out.Close()

	fmt.Println("Downloading init.m4s")
	res, _ = http.Get(initURL)
	io.Copy(out, res.Body)

	for i, seg := range segments {
		res, _ = http.Get(seg)
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
