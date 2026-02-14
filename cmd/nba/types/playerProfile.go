package types

type PlayerBio struct {
	FirstName        string `json:"FIRST_NAME" isVisible:"false"`
	LastName         string `json:"LAST_NAME" isVisible:"false"`
	DisplayName      string `json:"DISPLAY_FIRST_LAST" isVisible:"false"`
	TeamName         string `json:"TEAM_NAME" isVisible:"false"`
	TeamAbbreviation string `json:"TEAM_ABBREVIATION" isVisible:"false"`
	TeamCity         string `json:"TEAM_CITY" isVisible:"false"`
	JerseyNumber     string `json:"JERSEY" isVisible:"false"`
	Position         string `json:"POSITION" isVisible:"false"`
	Height           string `json:"HEIGHT" isVisible:"false"`
	Weight           string `json:"WEIGHT" isVisible:"false"`
	Country          string `json:"COUNTRY" isVisible:"false"`
	School           string `json:"SCHOOL" isVisible:"false"`
	DraftYear        string `json:"DRAFT_YEAR" isVisible:"false"`
	DraftRound       string `json:"DRAFT_ROUND" isVisible:"false"`
	DraftNumber      string `json:"DRAFT_NUMBER" isVisible:"false"`
	SeasonExp        int    `json:"SEASON_EXP" isVisible:"false"`
	Birthdate        string `json:"BIRTHDATE" isVisible:"false"`
	TeamID           int    `json:"TEAM_ID" isVisible:"false"`
}

type SeasonStats struct {
	SeasonID string  `json:"SEASON_ID" isVisible:"true" display:"Season" width:"10"`
	TeamAbbr string  `json:"TEAM_ABBREVIATION" isVisible:"true" display:"Team" width:"8"`
	GP       int     `json:"GP" isVisible:"true" display:"GP" width:"6"`
	MIN      float64 `json:"MIN" isVisible:"true" display:"MPG" width:"8"`
	PTS      float64 `json:"PTS" isVisible:"true" display:"PPG" width:"8"`
	REB      float64 `json:"REB" isVisible:"true" display:"RPG" width:"8"`
	AST      float64 `json:"AST" isVisible:"true" display:"APG" width:"8"`
	STL      float64 `json:"STL" isVisible:"true" display:"SPG" width:"8"`
	BLK      float64 `json:"BLK" isVisible:"true" display:"BPG" width:"8"`
	FGPCT    float64 `json:"FG_PCT" percentage:"true" isVisible:"true" display:"FG%" width:"8"`
	FG3PCT   float64 `json:"FG3_PCT" percentage:"true" isVisible:"true" display:"3PT%" width:"8"`
	FTPCT    float64 `json:"FT_PCT" percentage:"true" isVisible:"true" display:"FT%" width:"8"`
}

type GameLogEntry struct {
	GameDate  string  `json:"GAME_DATE" isVisible:"true" display:"Date" width:"12"`
	Matchup   string  `json:"MATCHUP" isVisible:"true" display:"Matchup" width:"15"`
	WL        string  `json:"WL" isVisible:"true" display:"W/L" width:"6"`
	MIN       int     `json:"MIN" isVisible:"true" display:"MIN" width:"6"`
	PTS       int     `json:"PTS" isVisible:"true" display:"PTS" width:"6"`
	REB       int     `json:"REB" isVisible:"true" display:"REB" width:"6"`
	AST       int     `json:"AST" isVisible:"true" display:"AST" width:"6"`
	STL       int     `json:"STL" isVisible:"true" display:"STL" width:"6"`
	BLK       int     `json:"BLK" isVisible:"true" display:"BLK" width:"6"`
	FGPCT     float64 `json:"FG_PCT" percentage:"true" isVisible:"true" display:"FG%" width:"8"`
	PlusMinus float64 `json:"PLUS_MINUS" isVisible:"true" display:"+/-" width:"8"`
}

type SeasonStatsList []SeasonStats
type GameLog []GameLogEntry

func (pb PlayerBio) ToStringSlice() []string {
	return structToStringSlice(pb)
}

func (ss SeasonStats) ToStringSlice() []string {
	return structToStringSlice(ss)
}

func (gl GameLogEntry) ToStringSlice() []string {
	return structToStringSlice(gl)
}
