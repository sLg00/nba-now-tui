package types

type GameHeader struct {
	GameDateEst                   string `json:"GAME_DATE_EST"`
	GameSequence                  int    `json:"GAME_SEQUENCE"`
	GameID                        string `json:"GAME_ID" isVisible:"false"`
	GameStatusID                  int    `json:"GAME_STATUS_ID" isVisible:"false"`
	GameStatusText                string `json:"GAME_STATUS_TEXT"`
	GameCode                      string `json:"GAMECODE" isVisible:"false"`
	HomeTeamID                    int    `json:"HOME_TEAM_ID" isVisible:"false"`
	VisitorTeamID                 int    `json:"VISITOR_TEAM_ID" isVisible:"false"`
	Season                        string `json:"SEASON"`
	LivePeriod                    int    `json:"LIVE_PERIOD"`
	LivePcTime                    string `json:"LIVE_PC_TIME"`
	NatlTvBroadcasterAbbreviation string `json:"NATL_TV_BROADCASTER_ABBREVIATION"`
	HomeTvBroadcasterAbbreviation string `json:"HOME_TV_BROADCASTER_ABBREVIATION"`
	AwayTvBroadcasterAbbreviation string `json:"AWAY_TV_BROADCASTER_ABBREVIATION"`
	LivePeriodTimeBcast           string `json:"LIVE_PERIOD_TIME_BCAST"`
	ArenaName                     string `json:"ARENA_NAME"`
	WhStatus                      int    `json:"WH_STATUS"`
	WNBACommissionerFlag          int    `json:"WNBA_COMMISSIONER_FLAG"`
}

type LineScore struct {
	GameDateEst      string  `json:"GAME_DATE_EST"`
	GameSequence     int     `json:"GAME_SEQUENCE"`
	GameID           string  `json:"GAME_ID" isVisible:"false"`
	TeamID           int     `json:"TEAM_ID" isVisible:"false"`
	TeamAbbreviation string  `json:"TEAM_ABBREVIATION"`
	TeamCityName     string  `json:"TEAM_CITY_NAME"`
	TeamName         string  `json:"TEAM_NAME"`
	TeamWinsLosses   string  `json:"TEAM_WINS_LOSSES"`
	PtsQtr1          int     `json:"PTS_QTR1"`
	PtsQtr2          int     `json:"PTS_QTR2"`
	PtsQtr3          int     `json:"PTS_QTR3"`
	PtsQtr4          int     `json:"PTS_QTR4"`
	PtsOt1           int     `json:"PTS_OT1"`
	PtsOt2           int     `json:"PTS_OT2"`
	PtsOt3           int     `json:"PTS_OT3"`
	PtsOt4           int     `json:"PTS_OT4"`
	PtsOt5           int     `json:"PTS_OT5"`
	PtsOt6           int     `json:"PTS_OT6"`
	PtsOt7           int     `json:"PTS_OT7"`
	PtsOt8           int     `json:"PTS_OT8"`
	PtsOt9           int     `json:"PTS_OT9"`
	PtsOt10          int     `json:"PTS_OT10"`
	Pts              int     `json:"PTS"`
	FgPct            float64 `json:"FG_PCT" percentage:"true"`
	FtPct            float64 `json:"FT_PCT" percentage:"true"`
	Fg3Pct           float64 `json:"FG3_PCT" percentage:"true"`
	Ast              int     `json:"AST"`
	Reb              int     `json:"REB"`
	Tov              int     `json:"TOV"`
}

type SeriesStandings struct {
	GameID         string `json:"GAME_ID" isVisible:"false"`
	HomeTeamID     int    `json:"HOME_TEAM_ID" isVisible:"false"`
	VisitorTeamID  int    `json:"VISITOR_TEAM_ID" isVisible:"false"`
	GameDateEst    string `json:"GAME_DATE_EST"`
	HomeTeamWins   int    `json:"HOME_TEAM_WINS"`
	HomeTeamLosses int    `json:"HOME_TEAM_LOSSES"`
	SeriesLeader   string `json:"SERIES_LEADER"`
}

type LastMeeting struct {
	GameID                       string `json:"GAME_ID" isVisible:"false"`
	LastGameID                   string `json:"LAST_GAME_ID" isVisible:"false"`
	LastGameDateEst              string `json:"LAST_GAME_DATE_EST"`
	LastGameHomeTeamID           int    `json:"LAST_GAME_HOME_TEAM_ID" isVisible:"false"`
	LastGameHomeTeamCity         string `json:"LAST_GAME_HOME_TEAM_CITY"`
	LastGameHomeTeamName         string `json:"LAST_GAME_HOME_TEAM_NAME"`
	LastGameHomeTeamAbbreviation string `json:"LAST_GAME_HOME_TEAM_ABBREVIATION"`
	LastGameHomeTeamPoints       int    `json:"LAST_GAME_HOME_TEAM_POINTS"`
	LastGameVisitorTeamID        int    `json:"LAST_GAME_VISITOR_TEAM_ID" isVisible:"false"`
	LastGameVisitorTeamCity      string `json:"LAST_GAME_VISITOR_TEAM_CITY"`
	LastGameVisitorTeamName      string `json:"LAST_GAME_VISITOR_TEAM_NAME"`
	LastGameVisitorTeamCity1     string `json:"LAST_GAME_VISITOR_TEAM_CITY1"`
	LastGameVisitorTeamPoints    int    `json:"LAST_GAME_VISITOR_TEAM_POINTS"`
}

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

func (g GameResult) ToStringSlice() []string {
	return structToStringSlice(g)
}

func (d DailyGameResults) ToStringSlice() []string {
	return structToStringSlice(d)
}

func (bst BoxScoreTeam) ToStringSlice() []string {
	return structToStringSlice(bst)
}
