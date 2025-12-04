package tools

import (
	"log"
	"os"
	"strings"
)

func WriteToFile(filename string, lines []string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(strings.Join(lines, "\n"))
	if err != nil {
		log.Fatal(err)
	}
}
