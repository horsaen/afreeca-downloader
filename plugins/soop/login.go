package soop

import (
	"encoding/json"
	"fmt"
	"horsaen/afreeca-downloader/tools"
	"io"
	"log"
	"net/http"
	"strings"
)

type LoginStruct struct {
	Result int `json:"RESULT"`
}

func UserLogin(username, password string) bool {
	url := "https://login.sooplive.co.kr/app/LoginAction.php"

	payload := fmt.Sprintf("szWork=login&szType=json&szUid=%s&szPassword=%s&isSaveId=false&szScriptVar=oLoginRet&szAction=&isLoginRetain=Y", username, password)

	req, _ := http.NewRequest("POST", url, strings.NewReader(payload))

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var login LoginStruct
	json.Unmarshal(body, &login)

	if login.Result == 1 {
		cookies := resp.Cookies()

		for _, cookie := range cookies {
			if cookie.Name == "AuthTicket" {
				tools.WriteCookies(cookie.Value, "soop")
				return true
			}
		}
	}

	return false
}
