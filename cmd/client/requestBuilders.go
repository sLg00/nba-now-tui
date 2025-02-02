package client

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	//RequestURL is a custom type that represent the specific URLs required to make requests towards the NBA APIs
	RequestURL  string
	RequestType string
	PerMode     string
	SeasonType  string
	Scope       string
)

const (
	LeagueID                    = "00"
	URL                         = "https://stats.nba.com/stats/"
	LeagueLeaders   RequestType = "leagueleaders"
	SeasonStandings RequestType = "leaguestandingsv3"
	DailyScoreboard RequestType = "scoreboardv2"
	BoxScore        RequestType = "boxscoretraditionalv3"
	TeamInfo        RequestType = "teaminfocommon"
	PerGame         PerMode     = "PerGame"
	RegularSeason   SeasonType  = "Regular Season"
	ScopeS          Scope       = "S"
)

// HTTPHeaderSet returns a http header required for the NBA API
func HTTPHeaderSet() http.Header {
	return http.Header{
		"User-Agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:72.0), Gecko/20100101, Firefox/72.0"},
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

// leagueLeadersAPIRequestBuilder is a builder function to construct the appropriate URL to query leagueLeaders
func leagueLeadersAPIRequestBuilder() RequestURL {
	params := url.Values{}
	params.Set("LeagueID", LeagueID)
	params.Set("PerMode", string(PerGame))
	params.Set("Scope", string(ScopeS))
	params.Set("Season", identifySeason())
	params.Set("SeasonType", string(RegularSeason))
	params.Set("StatCategory", "PTS")

	u, err := url.Parse(URL + string(LeagueLeaders))
	if err != nil {
		return ""
	}
	u.RawQuery = params.Encode()
	return RequestURL(u.String())
}

func seasonStandingsAPIRequestBuilder() RequestURL {
	params := url.Values{}
	params.Set("LeagueID", LeagueID)
	params.Set("Season", identifySeason())
	params.Set("SeasonType", string(RegularSeason))

	u, err := url.Parse(URL + string(SeasonStandings))
	if err != nil {
		return ""
	}
	u.RawQuery = params.Encode()
	return RequestURL(u.String())

}

func dailyScoreboardAPIRequestBuilder() RequestURL {
	today, err := GetDateArg()
	if err != nil {
		log.Fatal(err)
	}

	params := url.Values{}
	params.Set("DayOffset", "0")
	params.Set("GameDate", today)
	params.Set("LeagueID", LeagueID)

	u, err := url.Parse(URL + string(DailyScoreboard))
	if err != nil {
		return ""
	}

	u.RawQuery = params.Encode()
	return RequestURL(u.String())

}

// boxScoreRequestBuilder creates the URL for the API call to query a specific game's box score.
func boxScoreRequestBuilder(s string) RequestURL {
	params := url.Values{}
	params.Set("EndPeriod", "4")
	params.Set("EngRange", "0")
	params.Set("GameID", s)
	params.Set("RangeType", "0")
	params.Set("StartPeriod", "1")
	params.Set("StartRange", "0")

	u, err := url.Parse(URL + string(BoxScore))
	if err != nil {
		return ""
	}

	u.RawQuery = params.Encode()
	return RequestURL(u.String())

}

func teamInfoCommonRequestBuilder(s string) RequestURL {
	params := url.Values{}
	params.Set("LeagueID", LeagueID)
	params.Set("Season", identifySeason())
	params.Set("TeamID", s)

	u, err := url.Parse(URL + string(TeamInfo))
	if err != nil {
		return ""
	}

	u.RawQuery = params.Encode()
	return RequestURL(u.String())
}

func BuildRequests(s string) map[string]RequestURL {
	urlMap := map[string]RequestURL{
		"leagueLeadersURL":   leagueLeadersAPIRequestBuilder(),
		"seasonStandingsURL": seasonStandingsAPIRequestBuilder(),
		"dailyScoresURL":     dailyScoreboardAPIRequestBuilder(),
		"boxScoreURL":        boxScoreRequestBuilder(s),
		"teamInfoCommonURL":  teamInfoCommonRequestBuilder(s),
	}
	return urlMap
}

func GetDateArg() (string, error) {
	if len(os.Args) != 3 {
		log.Println("Cannot invoke program, date not provided in command line arguments.")
		err := fmt.Errorf("please use %s -d \"YYYY-MM-DD\"", os.Args[0])
		return "", err
	}
	if os.Args[1] == "-h" {
		err := fmt.Errorf("to invoke the program, please use %s -d \"YYYY-DD-MM\" with any date", os.Args[0])
		return "", err
	}
	dateStr := os.Args[2]

	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		return "", fmt.Errorf("date must be in YYYY-MM-DD format")
	}

	return dateStr, nil
}
