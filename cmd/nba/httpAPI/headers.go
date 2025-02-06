package httpAPI

import "net/http"

func (h *HTTPClient) SetHeaders() http.Header {
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
