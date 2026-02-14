package nbaAPI

import (
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	BaseURL  = "https://stats.nba.com/stats/"
	LeagueID = "00"
)

type (
	RequestURL  string
	RequestType string
	PerMode     string
	SeasonType  string
	Scope       string
)

type RequestParams interface {
	ToValues() url.Values
	Endpoint() string
	Validate() error
}

type nbaDateProvider struct {
	date string
}

type LeagueLeadersParams struct {
	LeagueID     string
	PerMode      PerMode
	Scope        Scope
	Season       string
	SeasonType   SeasonType
	StatCategory string
}

type SeasonStandingsParams struct {
	LeagueID   string
	Season     string
	SeasonType SeasonType
}

type DailyScoresParams struct {
	DayOffset string
	GameDate  string
	LeagueID  string
}

type BoxScoreParams struct {
	EndPeriod   string
	EndRange    string
	GameID      string
	RangeType   string
	StartPeriod string
	StartRange  string
}

type TeamProfileParams struct {
	LeagueID string
	Season   string
	TeamID   string
}

type PlayerIndexParams struct {
	LeagueID string
	Season   string
	TeamID   string
}

func NewDateProvider() types.DateProvider {
	eastern, _ := time.LoadLocation("America/New_York")
	today := time.Now().In(eastern).Format("2006-01-02")
	return &nbaDateProvider{date: today}
}

func (dp *nbaDateProvider) GetCurrentDate() (string, error) {
	return dp.date, nil
}

// GetCurrentSeason calculates and provides the season date string in the format of YYYY-YY.
func (dp *nbaDateProvider) GetCurrentSeason() string {
	dateSplit := strings.Split(dp.date, "-")
	year, _ := strconv.Atoi(dateSplit[0])
	month, _ := strconv.Atoi(dateSplit[1])

	if month < 10 {
		return fmt.Sprintf("%d-%02d", year-1, year%100)
	}
	return fmt.Sprintf("%d-%02d", year, (year+1)%100)
}

func (p LeagueLeadersParams) ToValues() url.Values {
	values := url.Values{}
	values.Set("LeagueID", p.LeagueID)
	values.Set("PerMode", string(p.PerMode))
	values.Set("Scope", string(p.Scope))
	values.Set("Season", p.Season)
	values.Set("SeasonType", string(p.SeasonType))
	values.Set("StatCategory", p.StatCategory)
	return values
}

func (p LeagueLeadersParams) Endpoint() string {
	return "leagueleaders"
}

func (p LeagueLeadersParams) Validate() error {
	if p.LeagueID == "" {
		return fmt.Errorf("leagueID is required")
	}
	if p.Season == "" {
		return fmt.Errorf("season is required")
	}
	return nil
}

func (p SeasonStandingsParams) ToValues() url.Values {
	values := url.Values{}
	values.Set("LeagueID", p.LeagueID)
	values.Set("Season", p.Season)
	values.Set("SeasonType", string(p.SeasonType))
	return values
}

func (p SeasonStandingsParams) Endpoint() string {
	return "leaguestandingsv3"
}

func (p SeasonStandingsParams) Validate() error {
	if p.LeagueID == "" {
		return fmt.Errorf("leagueID is required")
	}
	if p.Season == "" {
		return fmt.Errorf("season is required")
	}
	return nil
}

func (p DailyScoresParams) ToValues() url.Values {
	values := url.Values{}
	values.Set("DayOffset", p.DayOffset)
	values.Set("GameDate", p.GameDate)
	values.Set("LeagueID", p.LeagueID)
	return values
}

func (p DailyScoresParams) Endpoint() string {
	return "scoreboardv3"
}

func (p DailyScoresParams) Validate() error {
	if p.LeagueID == "" {
		return fmt.Errorf("leagueID is required")
	}
	if p.GameDate == "" {
		return fmt.Errorf("gameDate is required")
	}
	return nil
}

func (p BoxScoreParams) ToValues() url.Values {
	values := url.Values{}
	values.Set("EndPeriod", p.EndPeriod)
	values.Set("EndRange", p.EndRange)
	values.Set("GameID", p.GameID)
	values.Set("RangeType", p.RangeType)
	values.Set("StartPeriod", p.StartPeriod)
	values.Set("StartRange", p.StartRange)
	return values
}

func (p BoxScoreParams) Endpoint() string { return "boxscoretraditionalv3" }

func (p BoxScoreParams) Validate() error {
	if p.GameID == "" {
		return fmt.Errorf("gameID is required")
	}
	if p.EndPeriod == "" {
		return fmt.Errorf("endPeriod is required")
	}
	return nil
}

func (p TeamProfileParams) ToValues() url.Values {
	values := url.Values{}
	values.Set("LeagueID", p.LeagueID)
	values.Set("Season", p.Season)
	values.Set("TeamID", p.TeamID)
	return values
}

func (p TeamProfileParams) Endpoint() string { return "teaminfocommon" }

func (p TeamProfileParams) Validate() error {
	if p.LeagueID == "" {
		return fmt.Errorf("leagueID is required")
	}
	if p.TeamID == "" {
		return fmt.Errorf("teamID is required")
	}
	return nil
}

func (pi PlayerIndexParams) ToValues() url.Values {
	values := url.Values{}
	values.Set("LeagueID", pi.LeagueID)
	values.Set("Season", pi.Season)
	values.Set("TeamID", pi.TeamID)
	return values
}

func (pi PlayerIndexParams) Endpoint() string {
	return "playerindex"
}

func (pi PlayerIndexParams) Validate() error {
	if pi.LeagueID == "" {
		return fmt.Errorf("leagueID is required")
	}
	if pi.Season == "" {
		return fmt.Errorf("season is required")
	}
	if pi.TeamID == "" {
		return fmt.Errorf("teamID is required")
	}
	return nil
}
