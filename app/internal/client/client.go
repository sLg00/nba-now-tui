package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

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
				json := InitiateClient(v)
				err := WriteToFiles(pc.LLFullPath(), json)
				if err != nil {
					return
				}
			}
		case "seasonStandingsURL":
			fileToCheck := fileChecker(pc.SSFullPath())
			if !fileToCheck {
				json := InitiateClient(v)
				err := WriteToFiles(pc.SSFullPath(), json)
				if err != nil {
					return
				}
			}
		case "dailyScoresURL":
			fileToCheck := fileChecker(pc.DSBFullPath())
			if !fileToCheck {
				json := InitiateClient(v)
				err := WriteToFiles(pc.DSBFullPath(), json)
				if err != nil {
					return
				}
			}
		}
	}
}
