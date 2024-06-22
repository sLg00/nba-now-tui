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

// MakeRequests queries the NBA APIs and populates the respective files with the returned JSON IF the files do not already
// exist for the given day
func MakeRequests() {
	urlMap := BuildRequests()
	pc := InstantiatePaths()
	for k, v := range urlMap {
		switch k {
		case "leagueLeadersURL":
			fileToCheck := fileChecker(pc.LLFullPath())
			if !fileToCheck {
				LLJson = InitiateClient(v)
				err := WriteToFiles(pc.LLFullPath(), LLJson)
				if err != nil {
					return
				}
			}
		case "seasonStandingsURL":
			fileToCheck := fileChecker(pc.SSFullPath())
			if !fileToCheck {
				SSJson = InitiateClient(v)
				err := WriteToFiles(pc.SSFullPath(), SSJson)
				if err != nil {
					return
				}
			}
		case "dailyScoresURL":
			fileToCheck := fileChecker(pc.DSBFullPath())
			if !fileToCheck {
				DSBJson = InitiateClient(v)
				err := WriteToFiles(pc.DSBFullPath(), DSBJson)
				if err != nil {
					return
				}
			}
		}
	}
}
