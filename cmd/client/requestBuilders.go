package client

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	// requestURL is a custom type that represent the specific URLs required to make requests towards the NBA APIs
	requestURL string
)

// HTTPHeaderSet returns a http header required for the NBA API
func HTTPHeaderSet() http.Header {
	return http.Header{
		"User-Agent":         {"Mozilhttps://stats.nba.com/stats/la/5.0, (Windows NT 10.0; Win64; x64; rv:72.0), Gecko/20100101, Firefox/72.0"},
		"Accept":             {"application/json; charset=utf-8 , text/plain, */*"},
		"Accept-Language":    {"en-US, en;q=0.5firefox-125.0b3.tar.bz2"},
		"Accept-Encoding":    {"deflate, br"},
		"x-nba-stats-origin": {"stats"},
		"x-nba-stats-token":  {"true"},
		"Connection":         {"keep-alive"},
		"Referer":            {"https://stats.nba.com/"},
		"Pragma":             {"no-cache"},
		"Cache-Control":      {"no-cache"},
	}
}

// LeagueID is always 00 for the requests going against NBA APIs
const (
	LeagueID = "00"
	URL      = "https://stats.nba.com/stats/"
)

// identifySeason determines what season is currently ongoing and formats it in a way that is needed to query NBA APIs
func identifySeason() string {
	var seasonString string
	date, _ := GetDateArg()
	dateSplit := strings.Split(date, "-")
	year := dateSplit[0]
	month := dateSplit[1]

	monthInt, _ := strconv.Atoi(month)
	yearInt, _ := strconv.Atoi(year)
	previousYear := yearInt - 1
	nextYear := yearInt + 1

	if monthInt < 10 {
		dateStringPartOne := strconv.Itoa(previousYear)
		dateStringPartTwo := strconv.Itoa(yearInt)[2:]
		seasonString = dateStringPartOne + "-" + dateStringPartTwo

	}

	if monthInt >= 10 {
		dateStringPartOne := strconv.Itoa(yearInt)
		dateStringPartTwo := strconv.Itoa(nextYear)[2:]
		seasonString = dateStringPartOne + "-" + dateStringPartTwo
	}
	return seasonString
}

// leagueLeadersAPIRequestBuilder creates the URL for the API request from dynamic- and hardcoded building blocks
// TODO: SeasonType key can have 3 values, need to add identification for regular season/playoffs/preseason
func leagueLeadersAPIRequestBuilder() requestURL {
	return requestURL(URL + "leagueleaders?ActiveFlag=&LeagueID=" +
		LeagueID + "&PerMode=PerGame&Scope=S&Season=" +
		identifySeason() + "&SeasonType=Regular+Season&StatCategory=PTS")
}

func seasonStandingsAPIRequestBuilder() requestURL {
	return requestURL(URL + "leaguestandingsv3?LeagueID=" +
		LeagueID + "&Season=" + identifySeason() + "&SeasonType=Regular+Season")

}

func dailyScoreboardAPIRequestBuilder() requestURL {
	today, err := GetDateArg()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	return requestURL(URL + "scoreboardv2?DayOffset=0&GameDate=" + today + "&LeagueID=" + LeagueID)
}

// BoxScoreRequestBuilder creates the URL for the API call to query a specific game's box score.
// This URL does not go into the urlMap returned by the BuildRequests function.
// This is due to the fact,that box scores are fetched on-demand, not up front. For now.
func boxScoreRequestBuilder(s string) requestURL {
	return requestURL(URL + "boxscoretraditionalv3?EndPeriod=1&EndRange=0&GameID=" +
		s + "&RangeType=0&StartPeriod=1&StartRange=0")
}

func BuildRequests(s string) map[string]requestURL {
	urlMap := map[string]requestURL{
		"leagueLeadersURL":   leagueLeadersAPIRequestBuilder(),
		"seasonStandingsURL": seasonStandingsAPIRequestBuilder(),
		"dailyScoresURL":     dailyScoreboardAPIRequestBuilder(),
		"boxScoreURL":        boxScoreRequestBuilder(s),
	}
	return urlMap
}

func GetDateArg() (string, error) {
	if len(os.Args) != 3 || os.Args[1] != "-d" {
		log.Println("Cannot invoke program, date not provided in command line arguments.")
		err := fmt.Errorf("Please use %s -d \"YYYY-MM-DD\"", os.Args[0])

		return "", err
	}
	if os.Args[1] == "-h" {
		fmt.Printf("Please use %s -d \"YYYY-DD-MM\"", os.Args[0])
	}
	dateStr := os.Args[2]

	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		return "", fmt.Errorf("date must be in YYYY-MM-DD format")
	}

	return dateStr, nil
}
