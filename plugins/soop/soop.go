package soop

import (
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"os"
	"time"
)

func Concurrent(user *[]string) {
	exists, pwd, online, broad := GetUserData((*user)[0])

	if !exists {
		(*user)[2] = "Not found"
		(*user)[3] = "Not found"
		(*user)[4] = "EXITED"
		return
	}

	for !online {
		// fmt.Println("Stream is password protected, trying again in 3 minutes.")
		time.Sleep(3 * time.Minute)
		exists, pwd, online, broad = GetUserData((*user)[0])
	}

	for pwd {
		(*user)[2] = "PWD"
		(*user)[3] = "PWD"
		(*user)[4] = "CHECKING"
		time.Sleep(3 * time.Minute)
		exists, pwd, online, broad = GetUserData((*user)[0])
	}

	base := GetBase(broad)
	master := GetMaster((*user)[0])

	stream := GetStream(base, master)

	ConcurrentDownload(user, broad, stream)
}

func Start(bjId string) {
	tools.Exists("downloads/Soop")
	exists, pwd, online, broad := GetUserData(bjId)

	if !exists {
		fmt.Println("User not found.")
		os.Exit(1)
	}

	for !online {
		fmt.Printf("User %s offline, rechecking in 3 minutes.\n", bjId)
		time.Sleep(3 * time.Minute)
		exists, pwd, online, broad = GetUserData(bjId)
	}

	for pwd {
		fmt.Println("Stream is password protected, trying again in 3 minutes.")
		time.Sleep(3 * time.Minute)
		exists, pwd, online, broad = GetUserData(bjId)
	}

	base := GetBase(broad)
	master := GetMaster(bjId)

	stream := GetStream(base, master)

	Download(bjId, broad, stream)
}
