package soop

import (
	"encoding/json"
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type UserVod struct {
	Data []struct {
		TitleNo   int `json:"title_no"`
		StationNo int `json:"station_no"`
	} `json:"data"`
	Links struct {
		Next string `json:"next"`
	} `json:"links"`
}

type Vod struct {
	Result int `json:"result"`
	Data   struct {
		Files []struct {
			File string `json:"file"`
		} `json:"files"`
	} `json:"data"`
}

func UserVods(bjid string) {
	os.MkdirAll("downloads/soop/"+bjid+"/vods", os.ModePerm)
	fmt.Println("Grabbing user vods")
	titles := GetUserVods(bjid)

	fmt.Println(strconv.Itoa(len(titles)) + " vods found")
	fmt.Println("Grabbing user's playlist urls, this might take a while..")

	var playlists []string
	for i, t := range titles {
		playlists = append(playlists, GetVideoManifest(t)...)
		fmt.Printf("\r%d playlist files found | Searched %d/%d vods", len(playlists), i+1, len(titles))
	}

	tools.WriteToFile("downloads/soop/"+bjid+"/vods/vod.txt", playlists)
	fmt.Println("\nPlaylists archived to downloads/" + bjid + "vod.txt")
	fmt.Println("Downloading vods")

	for _, p := range playlists {
		DownloadPlaylistVod(p, bjid)
	}
}

func GetUserVods(bjid string) []string {
	url := "https://chapi.sooplive.co.kr/api/" + bjid + "/vods/all/streamer"

	var playlists []string

	for url != "" {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		var userVod UserVod
		if err := json.Unmarshal(body, &userVod); err != nil {
			log.Fatal(err)
		}

		for _, v := range userVod.Data {
			playlists = append(playlists, strconv.Itoa(v.TitleNo))
		}

		url = userVod.Links.Next
	}

	return playlists
}

func GetVideoManifest(titleNo string) []string {
	url := "https://api.m.sooplive.co.kr/station/video/a/view"

	payload := fmt.Sprintf(
		"nTitleNo=%s&nApiLevel=11&nPlaylistIdx=0",
		titleNo,
	)

	req, _ := http.NewRequest("POST", url, strings.NewReader(payload))

	cookies := tools.LoadCookies("soop")

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:144.0) Gecko/20100101 Firefox/144.0")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "AuthTicket="+cookies[0])

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var vod Vod
	json.Unmarshal(body, &vod)

	var files []string
	for _, f := range vod.Data.Files {
		files = append(files, f.File)
	}

	return files
}
