package nbaAPI

import (
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"log"
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
	cmdArgs []string
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

func NewDateProvider(args []string) types.DateProvider {
	return &nbaDateProvider{cmdArgs: args}
}

func (dp *nbaDateProvider) GetCurrentDate() (string, error) {
	if len(dp.cmdArgs) != 3 {
		log.Println("Cannot invoke program, date not provided in command line arguments.")
		err := fmt.Errorf("please use %s -d \"YYYY-MM-DD\"", dp.cmdArgs[0])
		return "", err
	}
	if dp.cmdArgs[1] == "-h" {
		err := fmt.Errorf("to invoke the program, please use %s -d \"YYYY-DD-MM\" with any date", dp.cmdArgs[0])
		return "", err
	}
	dateStr := dp.cmdArgs[2]

	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		return "", fmt.Errorf("date must be in YYYY-MM-DD format")
	}

	return dateStr, nil
}

func (dp *nbaDateProvider) GetCurrentSeason() string {
	date, err := dp.GetCurrentDate()
	if err != nil {
		return ""
	}
	dateSplit := strings.Split(date, "-")
	year, _ := strconv.Atoi(dateSplit[0])
	month, _ := strconv.Atoi(dateSplit[1])

	//last season
	if month < 10 {
		return fmt.Sprintf("%d-%02d", year-1, year%100)
	}

	//current season
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
	return "scoreboardv2"
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
