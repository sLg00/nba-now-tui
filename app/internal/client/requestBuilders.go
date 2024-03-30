package client

import (
	"net/http"
)

// HTTPHeaderSet returns a http header
func HTTPHeaderSet() http.Header {
	var headerSet = http.Header{
		"User-Agent":         {"Mozilla/5.0, (Windows NT 10.0; Win64; x64; rv:72.0), Gecko/20100101, Firefox/72.0"},
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

type URL struct {
	Domain string
	Path   string
}

type GeneralComponents struct {
	LeagueID   string
	Season     string
	SeasonType string
}

type PlayerStatsComponents struct {
	PerMode      string
	StatCategory string
	Scope        string
	ActiveFlag   *string
}

type RequestComponents struct {
	URL                   URL
	Headers               http.Header
	GeneralComponents     GeneralComponents
	PlayerStatsComponents PlayerStatsComponents
}

func buildRequest(i int) string {
	var requestSignature string
	if i == 1 {
		requestSignature = "https://stats.nba.com/stats/leagueleaders?ActiveFlag=&LeagueID=00&PerMode=PerGame&Scope=S&Season=2023-24&SeasonType=Regular+Season&StatCategory=PTS"
	}
	return requestSignature
}
