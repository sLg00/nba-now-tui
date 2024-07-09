package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// Client struct encompasses all the functions required to do api calls and filesystem ops
type Client struct {
	HTTPClient       *http.Client
	HeaderSet        func() http.Header
	BuildRequests    func() map[string]requestURL
	InstantiatePaths func() PathComponents
	FileChecker      func(string) bool
	WriteToFiles     func(string, []byte) error
}

// NewClient instantiates a pointer to a Client struct with the appropriate values
func NewClient() *Client {
	return &Client{
		HTTPClient:       &http.Client{Timeout: time.Duration(5) * time.Second},
		HeaderSet:        HTTPHeaderSet,
		BuildRequests:    BuildRequests,
		InstantiatePaths: InstantiatePaths,
		FileChecker:      fileChecker,
		WriteToFiles:     WriteToFiles,
	}
}

// InitiateClient initializes client instances with the appropriate request URLs and headers
func (c *Client) InitiateClient(url requestURL) []byte {
	//client := http.Client{Timeout: time.Duration(5) * time.Second}
	req, _ := http.NewRequest("GET", string(url), nil)
	req.Header = c.HeaderSet()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		fmt.Println("err:", err)
		return nil
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Could not close file.")
		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	return body
}

// MakeRequests queries the NBA APIs and populates the respective files with the returned JSON IF the files do not already
// exist for the given day
func (c *Client) MakeRequests() {
	urlMap := c.BuildRequests()
	pc := c.InstantiatePaths()

	var wg sync.WaitGroup

	for k, v := range urlMap {
		wg.Add(1)

		go func(key string, url requestURL) {
			defer wg.Done()

			switch key {
			case "leagueLeadersURL":
				fileToCheck := c.FileChecker(pc.LLFullPath())
				if !fileToCheck {
					json := c.InitiateClient(url)
					err := c.WriteToFiles(pc.LLFullPath(), json)
					if err != nil {
						return
					}
				}
			case "seasonStandingsURL":
				fileToCheck := fileChecker(pc.SSFullPath())
				if !fileToCheck {
					json := c.InitiateClient(url)
					err := c.WriteToFiles(pc.SSFullPath(), json)
					if err != nil {
						return
					}
				}
			case "dailyScoresURL":
				fileToCheck := fileChecker(pc.DSBFullPath())
				if !fileToCheck {
					json := c.InitiateClient(url)
					err := c.WriteToFiles(pc.DSBFullPath(), json)
					if err != nil {
						return
					}
				}
			}
		}(k, v)
	}
	wg.Wait()
}
