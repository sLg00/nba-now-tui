package client

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestInitiateClient(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected 'GET' request, got '%s'", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"key": "value"}`))
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
		HeaderSet: func() http.Header {
			headers := make(http.Header)
			headers.Set("Content-Type", "application/json")
			return headers
		},
	}

	url := requestURL(server.URL)
	response, err := client.InitiateClient(url)
	if err != nil {
		log.Println(err)
	}
	expectedResponse := []byte(`{"key": "value"}`)
	if !bytes.Equal(response, expectedResponse) {
		t.Errorf("Expected '%v', got '%v'", expectedResponse, response)
	}
}

func TestMakeDefaultRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"key": "value"}`))
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
		HeaderSet: func() http.Header {
			headers := make(http.Header)
			headers.Set("Content-Type", "application/json")
			return headers
		},
		BuildRequests: func(string) map[string]requestURL {
			return map[string]requestURL{
				"leagueLeadersURL":   requestURL(server.URL),
				"seasonStandingsURL": requestURL(server.URL),
				"dailyScoreboardURL": requestURL(server.URL),
			}
		},
		InstantiatePaths: func(string) *PathComponents {
			home, _ := os.UserHomeDir()
			return &PathComponents{
				Home:         home,
				Path:         "/.config/nba-tui/",
				LLFile:       "2024-01-01" + "_ll",
				SSFile:       "2024-01-01" + "_ss",
				DSBFile:      "2024-01-01" + "_dsb",
				BoxScorePath: "boxscores",
				BoxScoreFile: "2024-01-01" + "123123123",
			}
		},
		FileChecker: func(filePath string) bool {
			return false
		},
		WriteToFiles: func(filePath string, data []byte) error {
			return nil
		},
	}
	client.MakeDefaultRequests()
}
