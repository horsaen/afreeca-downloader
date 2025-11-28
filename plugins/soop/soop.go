package soop

import (
	"fmt"
)

func Start(bjid string) {
	if DvrCheck(bjid) {
		fmt.Printf("User %s online.\n", bjid)
		broad_no := GetBroadNo(bjid)

		aid := GetStreamAid(bjid, broad_no)

		server := GetStreamServer(broad_no)

		master := GetMasterPlist(server, aid)

		playlist := GetStreamQualities(master, 0)

		base, stream := GetStream(server, playlist)

		Download(bjid, broad_no, base, stream)
	}
}
