package soop

import (
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"time"
)

func Concurrent(index int, user []string, updates chan<- tools.ConcurrentRow) {
	bjid := user[0]
	tools.Exists("downloads/soop")

	for {
		broadNo := GetBroadNo(bjid)

		if broadNo == "" || broadNo == "0" {
			user[2] = "Offline"
			user[3] = "Offline"
			user[4] = "Offline"
			updates <- tools.SnapshotConcurrentRow(index, user)
			time.Sleep(3 * time.Second)
			continue
		}

		aid := GetStreamAid(bjid, broadNo, 0)
		server := GetStreamServer(broadNo)

		if server == "" {
			user[2] = "Offline"
			user[3] = "Offline"
			user[4] = "Offline"
			updates <- tools.SnapshotConcurrentRow(index, user)
			time.Sleep(3 * time.Second)
			continue
		}

		master := GetMasterPlist(server, aid)
		playlist := GetStreamQualities(master, 0)

		if playlist == "" {
			user[2] = "Offline"
			user[3] = "Offline"
			user[4] = "Offline"
			updates <- tools.SnapshotConcurrentRow(index, user)
			time.Sleep(3 * time.Second)
			continue
		}

		base, stream := GetStream(server, playlist)

		if !ConcurrentDownload(index, user, bjid, broadNo, base, stream, updates) {
			user[2] = "Offline"
			user[3] = "Offline"
			user[4] = "Offline"
			updates <- tools.SnapshotConcurrentRow(index, user)
			time.Sleep(3 * time.Second)
			continue
		}
	}
}

func Start(bjid string) {
	if DvrCheck(bjid) {
		fmt.Printf("User %s online.\n", bjid)
		broad_no := GetBroadNo(bjid)

		if broad_no == "" || broad_no == "0" {
			return
		}

		aid := GetStreamAid(bjid, broad_no, 0)

		server := GetStreamServer(broad_no)

		if server == "" {
			return
		}

		master := GetMasterPlist(server, aid)

		playlist := GetStreamQualities(master, 0)

		if playlist == "" {
			return
		}

		base, stream := GetStream(server, playlist)

		Download(bjid, broad_no, base, stream)
	}
}
