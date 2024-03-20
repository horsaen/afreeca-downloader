package concurrent

import (
	"bufio"
	"horsaen/afreeca-downloader/plugins/afreeca"
	"horsaen/afreeca-downloader/plugins/bigo"
	"horsaen/afreeca-downloader/plugins/flex"
	"horsaen/afreeca-downloader/plugins/kick"
	"horsaen/afreeca-downloader/plugins/panda"
	"horsaen/afreeca-downloader/tools"
	"os"
	"strings"
	"time"
)

func GetUsers() [][]string {
	var users [][]string

	home, _ := os.UserHomeDir()

	file, _ := os.Open(home + "/.afreeca-downloader/users")

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, ", ")

		init := []string{"Offline", "Offline", "Offline"}

		split = append(split, init...)

		users = append(users, split)
	}

	return users
}

func Start() {
	users := GetUsers()

	for i := range users {
		switch users[i][1] {
		case "afreeca":
			go afreeca.Concurrent(&users[i])
		case "bigo":
			go bigo.Concurrent(&users[i])
		case "kick":
			go kick.Concurrent(&users[i])
		case "flex":
			go flex.Concurrent(&users[i])
		case "panda":
			go panda.Concurrent(&users[i])
		}
	}

	for {
		Table(users)
		time.Sleep(1 * time.Second)
		tools.ClearCli()
	}
}
