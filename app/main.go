package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var headerSet = http.Header{
	"User-Agent":         {"Mozilla/5.0, (Windows NT 10.0; Win64; x64; rv:72.0), Gecko/20100101, Firefox/72.0"},
	"Accept":             {"application/json , text/plain, */*"},
	"Accept-Language":    {"en-US, en;q=0.5firefox-125.0b3.tar.bz2"},
	"Accept-Encoding":    {"gzip, deflate, br"},
	"x-nba-stats-origin": {"stats"},
	"x-nba-stats-token":  {"true"},
	"Connection":         {"keep-alive"},
	"Referer":            {"https://stats.nba.com/"},
	"Pragma":             {"no-cache"},
	"Cache-Control":      {"no-cache"},
}

type ParameterSet struct {
	LeagueID   string
	Season     string
	SeasonType string
}

type RequestSet struct {
	baselineUrl  string
	parameterSet ParameterSet
	headers      http.Header
}

func buildParameterSets() (ParameterSet, error) {
	seasonSet := ParameterSet{
		LeagueID:   "00",
		Season:     "2023-24",
		SeasonType: "Regular+Season",
	}
	return seasonSet, nil
}

func buildRequestSet() (RequestSet, error) {
	seasonSetImpl, _ := buildParameterSets()

	testRequest := RequestSet{
		baselineUrl:  "https://stats.nba.com/stats/leaguestandingsv3",
		parameterSet: seasonSetImpl,
		headers:      headerSet,
	}
	return testRequest, nil
}

func main() {
	rs, _ := buildRequestSet()
	parameterSetString := fmt.Sprintf("?LeagueID=%s&Season=%s&SeasonType=%s", rs.parameterSet.LeagueID,
		rs.parameterSet.Season, rs.parameterSet.SeasonType)

	c := http.Client{Timeout: time.Duration(5) * time.Second}
	req, err := http.NewRequest("GET", rs.baselineUrl+parameterSetString, nil)

	if err != nil {
		fmt.Println("err:&s", err)
		return
	}

	log.Println(rs.baselineUrl + parameterSetString)
	req.Header = headerSet

	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("err: %s", err)
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(body)
}
