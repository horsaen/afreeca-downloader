package tools

import (
	"os"
)

func Exists(dir string) {
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
}

func InitConfDir() {
	home, _ := os.UserHomeDir()

	configBase := home + "/.afreeca-downloader"
	cookies := configBase + "/cookies/"

	_, err := os.Stat(configBase)

	if os.IsNotExist(err) {
		os.Mkdir(configBase, 0755)
		os.Create(configBase + "/users")

		os.Mkdir(configBase+"/cookies", 0755)
		os.Create(cookies + "afreeca")
		os.Create(cookies + "flex")
		os.Create(cookies + "kick")
		os.Create(cookies + "panda")
		os.Create(cookies + "soop")
	}
}
