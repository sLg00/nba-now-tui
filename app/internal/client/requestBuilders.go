package client

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

var today = time.Now().Format("2006-01-31")

// HTTPHeaderSet returns a http header
func HTTPHeaderSet() http.Header {
	var headerSet = http.Header{
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
	return headerSet
}

// LeagueID is always 00 for the requests going against NBA APIs
const (
	LeagueID = "00"
	URL      = "https://stats.nba.com/stats/"
)

// identifySeason determines what season is currently ongoing and formats it in a way that is needed to query NBA APIs
func identifySeason() string {
	cyms := strings.Split(time.Now().Format("2006-01"), "-")
	year := cyms[0]
	month := cyms[1]
	yearInt, _ := strconv.Atoi(year)
	monthInt, _ := strconv.Atoi(month)

	lastYear := yearInt - 1
	nextYear := yearInt + 1

	if monthInt >= 7 {
		p1 := strconv.Itoa(yearInt)
		p2 := strconv.Itoa(nextYear)[2:]
		s := p1 + "-" + p2
		return s

	}
	if monthInt < 7 {
		p1 := strconv.Itoa(lastYear)
		p2 := strconv.Itoa(yearInt)[2:]
		s := p1 + "-" + p2
		return s
	}
	return ""
}

type (
	requestURL string
)

// leagueLeadersAPIRequestBuilder creates the URL for the API request from dynamic- and hardcoded building blocks
// TODO: SeasonType key can have 2 values, need to add identification for regular season/playoffs
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
	return requestURL(URL + "scoreboardv2?DayOffset=0&GameDate=" + today + "&LeagueID=" + LeagueID)
}

func BuildRequests() map[string]requestURL {
	urlMap := map[string]requestURL{
		"leagueLeadersURL":   leagueLeadersAPIRequestBuilder(),
		"seasonStandingsURL": seasonStandingsAPIRequestBuilder(),
		"dailyScoresURL":     dailyScoreboardAPIRequestBuilder(),
	}
	return urlMap
}
