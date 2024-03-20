package concurrent

import (
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func Table(users [][]string) {
	headerFmt := color.New(color.FgMagenta, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgMagenta).SprintfFunc()

	tbl := table.New("Platform", "Name", "Size", "Duration", "Path")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, user := range users {
		tbl.AddRow(user[1], user[0], user[2], user[3], user[4])
	}

	tbl.Print()
}
