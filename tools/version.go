package tools

import (
	"fmt"
	"os"
)

func Version() {
	fmt.Println("https://github.com/horsaen/afreeca-downloader")
	fmt.Println("prerelease")
	os.Exit(0)
}
