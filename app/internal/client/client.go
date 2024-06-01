package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var LLJson []byte

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

func MakeRequests() {
	urlMap := BuildRequests()
	url, ok := urlMap["leagueLeadersURL"]
	if !ok {
		log.Fatal("URL not found")
	}
	LLJson = InitiateClient(url)
}
