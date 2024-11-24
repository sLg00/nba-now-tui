package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Client struct encompasses all the functions required to do api calls and filesystem ops
type Client struct {
	HTTPClient       *http.Client
	HeaderSet        func() http.Header
	BuildRequests    func(string) map[string]requestURL
	InstantiatePaths func(string) *PathComponents
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
		err = fmt.Errorf("HTTP request failed: %v", err)
		log.Println(err)
		return nil
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			err = fmt.Errorf("closing response body failed: %v", err)
			log.Println(err)
		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	return body
}

// MakeDefaultRequests queries the NBA APIs and populates the respective files with the returned JSON IF the files do not already
// exist for the given day
func (c *Client) MakeDefaultRequests() error {
	defaultString := ""
	urlMap := c.BuildRequests(defaultString)
	pc := c.InstantiatePaths(defaultString)
	err := CleanOldFiles(pc)
	if err != nil {
		log.Printf("could not clean up old files: %v\n", err)
	}

	dChan := make(chan struct{}, len(urlMap))
	eChan := make(chan error, len(urlMap))

	for k, v := range urlMap {

		go func(key string, url requestURL) {
			defer func() { dChan <- struct{}{} }()

			switch key {
			case "leagueLeadersURL":
				fileToCheck := c.FileChecker(pc.LLFullPath())
				if !fileToCheck {
					json := c.InitiateClient(url)
					err := c.WriteToFiles(pc.LLFullPath(), json)
					if err != nil {
						eChan <- fmt.Errorf("could not write to files %v\n", err)

					}
				}
			case "seasonStandingsURL":
				fileToCheck := fileChecker(pc.SSFullPath())
				if !fileToCheck {
					json := c.InitiateClient(url)
					err := c.WriteToFiles(pc.SSFullPath(), json)
					if err != nil {
						eChan <- fmt.Errorf("could not write to files %v\n", err)
					}
				}
			case "dailyScoresURL":
				fileToCheck := fileChecker(pc.DSBFullPath())
				if !fileToCheck {
					json := c.InitiateClient(url)
					err := c.WriteToFiles(pc.DSBFullPath(), json)
					if err != nil {
						eChan <- fmt.Errorf("could not write to files %v\n", err)
					}
				}
			}
		}(k, v)
	}
	for i := 0; i < len(urlMap); i++ {
		<-dChan
	}
	close(eChan)

	var errs []error
	for e := range eChan {
		errs = append(errs, e)
	}

	if len(errs) > 0 {
		for _, apiErr := range errs {
			log.Printf("API error: %v", apiErr)
		}
		return fmt.Errorf("encountered errors during API requests")
	}

	return nil
}

// MakeOnDemandRequests takes a string (a gameId, a playerID etc) and queries the NBA API on-demand
func (c *Client) MakeOnDemandRequests(s string) error {
	urlMap := c.BuildRequests(s)
	path := c.InstantiatePaths(s)

	for k, v := range urlMap {
		switch k {
		case "boxScoreURL":
			fileToCheck := c.FileChecker(path.BoxScoreFullPath())
			if !fileToCheck {
				json := c.InitiateClient(v)
				err := c.WriteToFiles(path.BoxScoreFullPath(), json)
				if err != nil {
					err = fmt.Errorf("couldn't write to files: %v", err)
					log.Println(err)
					return err
				}
			}
		}
	}
	return nil
}
