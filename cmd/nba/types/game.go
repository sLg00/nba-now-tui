package types

type GameResult struct {
	GameID               string `isVisible:"false"`
	HomeTeamID           int    `isVisible:"false"`
	HomeTeamName         string
	HomeTeamPts          int
	HomeTeamAbbreviation string
	AwayTeamID           int `isVisible:"false"`
	AwayTeamName         string
	AwayTeamPts          int
	AwayTeamAbbreviation string
	GameStatusID         int `json:"GAME_STATUS_ID" isVisible:"false"`
}

type DailyGameResults []GameResult

type TeamGameStatistics struct {
	Minutes                 string  `json:"minutes"`
	FieldGoalsMade          int     `json:"fieldGoalsMade"`
	FieldGoalsAttempted     int     `json:"fieldGoalsAttempted"`
	FieldGoalsPercentage    float64 `json:"fieldGoalsPercentage" percentage:"true"`
	ThreePointersMage       int     `json:"threePointersMage"`
	ThreePointersAttempted  int     `json:"threePointersAttempted"`
	ThreePointersPercentage float64 `json:"threePointersPercentage" percentage:"true"`
	FreeThrowsMage          int     `json:"freeThrowsMage"`
	FreeThrowsAttempted     int     `json:"freeThrowsAttempted"`
	FreeThrowsPercentage    float64 `json:"freeThrowsPercentage" percentage:"true"`
	ReboundsOffensive       int     `json:"reboundsOffensive"`
	ReboundsDefensive       int     `json:"reboundsDefensive"`
	ReboundsTotal           int     `json:"reboundsTotal"`
	Assists                 int     `json:"assists"`
	Steals                  int     `json:"steals"`
	Blocks                  int     `json:"blocks"`
	Turnovers               int     `json:"turnovers"`
	FoulsPersonal           int     `json:"foulsPersonal"`
	Points                  int     `json:"points"`
	PlusMinusPoints         float64 `json:"plusMinusPoints"`
}

type BoxScoreTeam struct {
	TeamID             int                `json:"teamId" isVisible:"false"`
	TeamCity           string             `json:"teamCity"`
	TeamName           string             `json:"teamName"`
	TeamTriCode        string             `json:"teamTriCode"`
	TeamSlug           string             `json:"teamSlug"`
	BoxScorePlayers    BoxScorePlayers    `json:"players"`
	TeamGameStatistics TeamGameStatistics `json:"statistics"`
}

type BoxScore struct {
	GameID     string `json:"gameId" isVisible:"false"`
	AwayTeamId int    `json:"awayTeamId" isVisible:"false"`
	HomeTeamId int    `json:"homeTeamId" isVisible:"false"`
	HomeTeam   BoxScoreTeam
	AwayTeam   BoxScoreTeam
}

// ScoreboardV3Period represents a single period's score in the V3 API
type ScoreboardV3Period struct {
	Period     int    `json:"period"`
	PeriodType string `json:"periodType"`
	Score      int    `json:"score"`
}

// ScoreboardV3Team represents team data in the V3 scoreboard API
type ScoreboardV3Team struct {
	TeamID            int                  `json:"teamId"`
	TeamName          string               `json:"teamName"`
	TeamCity          string               `json:"teamCity"`
	TeamTricode       string               `json:"teamTricode"`
	TeamSlug          string               `json:"teamSlug"`
	Wins              int                  `json:"wins"`
	Losses            int                  `json:"losses"`
	Score             int                  `json:"score"`
	Seed              int                  `json:"seed"`
	InBonus           *string              `json:"inBonus"`
	TimeoutsRemaining int                  `json:"timeoutsRemaining"`
	Periods           []ScoreboardV3Period `json:"periods"`
}

// ScoreboardV3Game represents a single game in the V3 scoreboard API
type ScoreboardV3Game struct {
	GameID            string           `json:"gameId"`
	GameCode          string           `json:"gameCode"`
	GameStatus        int              `json:"gameStatus"`
	GameStatusText    string           `json:"gameStatusText"`
	Period            int              `json:"period"`
	GameClock         string           `json:"gameClock"`
	GameTimeUTC       string           `json:"gameTimeUTC"`
	GameEt            string           `json:"gameEt"`
	RegulationPeriods int              `json:"regulationPeriods"`
	SeriesGameNumber  string           `json:"seriesGameNumber"`
	GameLabel         string           `json:"gameLabel"`
	GameSubLabel      string           `json:"gameSubLabel"`
	SeriesText        string           `json:"seriesText"`
	IfNecessary       bool             `json:"ifNecessary"`
	SeriesConference  string           `json:"seriesConference"`
	PoRoundDesc       string           `json:"poRoundDesc"`
	GameSubtype       string           `json:"gameSubtype"`
	IsNeutral         bool             `json:"isNeutral"`
	HomeTeam          ScoreboardV3Team `json:"homeTeam"`
	AwayTeam          ScoreboardV3Team `json:"awayTeam"`
}

// ScoreboardV3Data represents the scoreboard data in the V3 API
type ScoreboardV3Data struct {
	GameDate   string             `json:"gameDate"`
	LeagueID   string             `json:"leagueId"`
	LeagueName string             `json:"leagueName"`
	Games      []ScoreboardV3Game `json:"games"`
}

func (g GameResult) ToStringSlice() []string {
	return structToStringSlice(g)
}

func (d DailyGameResults) ToStringSlice() []string {
	return structToStringSlice(d)
}

func (bst BoxScoreTeam) ToStringSlice() []string {
	return structToStringSlice(bst)
}
