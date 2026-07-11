package concurrent

import (
	"fmt"
	"strings"
)

func Table(users [][]string) {
	widths := []int{8, 20, 14, 10, 28}
	headers := []string{"Platform", "Name", "Size", "Duration", "Path"}

	for _, user := range users {
		for index, value := range user {
			if index >= len(widths) {
				continue
			}

			if len(value) > widths[index] {
				widths[index] = len(value)
			}
		}
	}

	fmt.Print("\033[H\033[2J")
	for index, header := range headers {
		fmt.Printf("%-*s  ", widths[index], header)
	}
	fmt.Println()

	for index := range headers {
		fmt.Print(strings.Repeat("-", widths[index]))
		if index < len(headers)-1 {
			fmt.Print("  ")
		}
	}
	fmt.Println()

	for _, user := range users {
		for index := 0; index < len(headers); index++ {
			value := ""
			if index < len(user) {
				value = user[index]
			}
			fmt.Printf("%-*s  ", widths[index], value)
		}
		fmt.Println()
	}
}
