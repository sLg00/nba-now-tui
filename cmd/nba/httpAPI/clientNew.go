package httpAPI

import (
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/client"
	"github.com/sLg00/nba-now-tui/cmd/nba/pathManager"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type HTTPRequester interface {
	Get(url RequestURL) ([]byte, error)
	SetHeaders() http.Header
}

type RequestBuilder interface {
	BuildRequests(param string) map[string]RequestURL
	BuildLeagueLeadersRequest() RequestURL
	BuildSeasonStandingsRequest() RequestURL
	BuildDailyScoresRequest() RequestURL
	//BuildBoxScoresRequest(gameID string) RequestURL
	//BuildTeamInfoRequest(teamID string) RequestURL
}

type nbaRequestBuilder struct {
	baseURL string
	dates   types.DateProvider
}

func NewRequestBuilder(baseURL string, dates types.DateProvider) RequestBuilder {
	return &nbaRequestBuilder{
		baseURL: baseURL,
		dates:   dates,
	}
}

func (rb *nbaRequestBuilder) BuildRequests(param string) map[string]RequestURL {
	return map[string]RequestURL{
		"leagueLeaders":   rb.BuildLeagueLeadersRequest(),
		"seasonStandings": rb.BuildSeasonStandingsRequest(),
		"dailyScores":     rb.BuildDailyScoresRequest(),
		//"boxScores:" rb.BuildBoxScoresRequests(),
		//"teamInfo":     rb.BuildTeamInfoRequest(),
	}
}

func (rb *nbaRequestBuilder) buildURL(params RequestParams) RequestURL {
	if err := params.Validate(); err != nil {
		log.Printf("Invalid request parameters for %s: %v", params.Endpoint(), err)
		return ""
	}
	u, err := url.Parse(rb.baseURL + params.Endpoint())
	if err != nil {
		log.Printf("failed to parse endpoint: %v", err)
		return ""
	}
	u.RawQuery = params.ToValues().Encode()
	return RequestURL(u.String())
}

func (rb *nbaRequestBuilder) BuildLeagueLeadersRequest() RequestURL {
	params := LeagueLeadersParams{
		LeagueID:     LeagueID,
		PerMode:      "PerGame",
		Scope:        "S",
		Season:       rb.dates.GetCurrentSeason(),
		SeasonType:   "Regular Season",
		StatCategory: "PTS",
	}
	return rb.buildURL(params)
}

func (rb *nbaRequestBuilder) BuildSeasonStandingsRequest() RequestURL {
	params := SeasonStandingsParams{
		LeagueID:   LeagueID,
		Season:     rb.dates.GetCurrentSeason(),
		SeasonType: "Regular Season",
	}
	return rb.buildURL(params)
}

func (rb *nbaRequestBuilder) BuildDailyScoresRequest() RequestURL {
	date, err := rb.dates.GetCurrentDate()
	if err != nil {
		log.Printf("failed to get current date: %v", err)
	}
	params := DailyScoresParams{
		DayOffset: "0",
		GameDate:  date,
		LeagueID:  LeagueID,
	}
	return rb.buildURL(params)
}

type NewClient struct {
	http             HTTPRequester
	requests         RequestBuilder
	Paths            pathManager.PathManager
	InstantiatePaths func(string) *client.PathComponents
	FileChecker      func(string) bool
	WriteToFiles     func(string, []byte) error
}

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{Timeout: time.Duration(8) * time.Second},
	}
}

func (h *HTTPClient) Get(url RequestURL) ([]byte, error) {
	req, _ := http.NewRequest("GET", string(url), nil)
	req.Header = h.SetHeaders()

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned status %d", resp.StatusCode)
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func NewNewClient() *NewClient {
	dateProvider := NewDateProvider(os.Args)
	return &NewClient{
		http:         NewHTTPClient(),
		requests:     NewRequestBuilder(BaseURL, dateProvider),
		Paths:        pathManager.PathFactory(dateProvider, ""),
		FileChecker:  client.FileChecker,
		WriteToFiles: client.WriteToFiles,
	}
}

func (c *NewClient) NewMakeDefaultRequests() error {
	urls := c.requests.BuildRequests("")

	dChan := make(chan struct{}, len(urls))
	eChan := make(chan error, len(urls))

	for name, reqURL := range urls {
		go func(name string, reqURL RequestURL) {
			defer func() { dChan <- struct{}{} }()

			path := c.Paths.GetFullPath(name, "")

			if name != "dailyScores" && c.FileChecker(path) {
				return
			}

			data, err := c.http.Get(reqURL)
			if err != nil {
				eChan <- fmt.Errorf("api error: %w", err)
			}

			if err = c.WriteToFiles(path, data); err != nil {
				eChan <- fmt.Errorf("write error for %s: %w", name, err)
			}
		}(name, reqURL)
	}

	for i := 0; i < len(urls); i++ {
		<-dChan
	}
	close(eChan)

	var errs []error

	for err := range eChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("encountered %d errors during API requests", len(errs))
	}
	return nil
}
