package soop

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetPlaylist(manifest string) {
	path := strings.Replace(manifest, ".smil", ".mp4", -1)

	res, _ := http.Get(path)

	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
}
