package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var LLJson []byte
var SSJson []byte
var DSBJson []byte

// InitiateClient initializes client instances with the appropriate request URLs and headers
func InitiateClient(url requestURL) []byte {
	client := http.Client{Timeout: time.Duration(5) * time.Second}
	req, _ := http.NewRequest("GET", string(url), nil)

	req.Header = HTTPHeaderSet()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err:", err)
		return nil
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	return body
}

// MakeRequests triggers the HTTP calls towards the NBA API
//func MakeRequests() {
//	urlMap := BuildRequests()
//
//	url, ok := urlMap["leagueLeadersURL"]
//	if !ok {
//		log.Fatal("URL not found")
//	}
//	//check whether today's LL file already exists before making HTTP call towards API
//	fileToCheck := fileChecker(LLFULLPATH)
//	if !fileToCheck {
//		LLJson = InitiateClient(url)
//		err := WriteToFiles()
//		if err != nil {
//			return
//		}
//	}
//}

// MakeRequests queries the NBA APIs and populates the respective files with the returned JSON
func MakeRequests() {
	urlMap := BuildRequests()

	for k, v := range urlMap {
		switch k {
		case "leagueLeadersURL":
			fileToCheck := fileChecker(LLFULLPATH)
			if !fileToCheck {
				LLJson = InitiateClient(v)
				err := WriteToFiles(LLFULLPATH, LLJson)
				if err != nil {
					return
				}
			}
		case "seasonStandingsURL":
			fileToCheck := fileChecker(SSFULLPATH)
			if !fileToCheck {
				SSJson = InitiateClient(v)
				err := WriteToFiles(SSFULLPATH, SSJson)
				if err != nil {
					return
				}
			}
		case "dailyScoresURL":
			fileToCheck := fileChecker(DSBFULLPATH)
			if !fileToCheck {
				DSBJson = InitiateClient(v)
				err := WriteToFiles(DSBFULLPATH, DSBJson)
				if err != nil {
					return
				}
			}
		}
	}
}
