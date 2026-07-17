package concurrent

import (
	"bufio"
	"fmt"
	"horsaen/afreeca-downloader/plugins/bigo"
	"horsaen/afreeca-downloader/plugins/flex"
	"horsaen/afreeca-downloader/plugins/panda"
	"horsaen/afreeca-downloader/plugins/soop"
	"horsaen/afreeca-downloader/tools"
	"os"
	"strings"
	"time"
)

func GetUsers() [][]string {
	users := make([][]string, 0)
	home, err := os.UserHomeDir()
	if err != nil {
		return users
	}

	file, err := os.Open(home + "/.afreeca-downloader/users")
	if err != nil {
		return users
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, ", ")
		if len(parts) < 2 {
			continue
		}

		row := append(parts, "Offline", "Offline", "Offline")
		users = append(users, row)
	}

	if err := scanner.Err(); err != nil {
		return users
	}

	return users
}

func CloneUsers(users [][]string) [][]string {
	cloned := make([][]string, len(users))
	for index, user := range users {
		row := make([]string, len(user))
		copy(row, user)
		cloned[index] = row
	}

	return cloned
}

func Start() {
	tools.ClearCli()

	users := GetUsers()
	if len(users) == 0 {
		fmt.Println("No users found in ~/.afreeca-downloader/users")
		return
	}

	updates := make(chan tools.ConcurrentRow, len(users)*4)
	state := CloneUsers(users)

	for index, user := range users {
		current := make([]string, len(user))
		copy(current, user)

		switch current[1] {
		case "bigo":
			go bigo.Concurrent(index, current, updates)
		case "soop":
			go soop.Concurrent(index, current, updates)
		case "flex":
			go flex.Concurrent(index, current, updates)
		case "panda":
			go panda.Concurrent(index, current, updates)
		default:
			current[2] = "Unsupported"
			current[3] = "Unsupported"
			current[4] = "Unsupported"
			state[index] = current
		}
	}

	Table(state)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case update := <-updates:
			if update.Index >= 0 && update.Index < len(state) {
				row := make([]string, len(update.Values))
				copy(row, update.Values)
				state[update.Index] = row
				Table(state)
			}
		case <-ticker.C:
			Table(state)
		}
	}
}
