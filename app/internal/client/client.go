package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	filepath2 "path/filepath"
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

	path := HOME + "/.config/nba-tui/"
	fmt.Println(path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) == true {
		err = os.Mkdir(path, 0777)
		if err != nil {
			fmt.Println("Error creating directory", err)
			return
		}
	}

	llFile := "ll_" + Today
	filePath := filepath2.Join(path, llFile)

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(LLJson)
	if err != nil {
		log.Println("could not write to file")
	}
}
