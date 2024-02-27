package tools

import (
	"fmt"
	"time"
)

func FormatTime(elapsed_time time.Duration) string {
	hours, minutes, seconds := int(elapsed_time.Hours()), int(elapsed_time.Minutes())%60, int(elapsed_time.Seconds())%60
	time := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

	return time
}
