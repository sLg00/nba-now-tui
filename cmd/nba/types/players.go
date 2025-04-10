package types

type PlayerStatistics struct {
	Minutes                 string  `json:"minutes" isVisible:"true" display:"Minutes"`
	FieldGoalsMade          int     `json:"fieldGoalsMade" isVisible:"true" display:"FG Made"`
	FieldGoalsAttempted     int     `json:"fieldGoalsAttempted" isVisible:"true" display:"FG Attempted"`
	FieldGoalsPercentage    float64 `json:"fieldGoalsPercentage" percentage:"true" display:"FG%"`
	ThreePointersMade       int     `json:"threePointersMade" isVisible:"true" display:"3PT Made"`
	ThreePointersAttempted  int     `json:"threePointersAttempted" isVisible:"true" display:"3PT Attempted"`
	ThreePointersPercentage float64 `json:"threePointersPercentage" percentage:"true" display:"3PT%"`
	FreeThrowsMade          int     `json:"freeThrowsMade" isVisible:"true" display:"FT Made"`
	FreeThrowsAttempted     int     `json:"freeThrowsAttempted" isVisible:"true" display:"FT Attempted"`
	FreeThrowsPercentage    float64 `json:"freeThrowsPercentage" percentage:"true" display:"FT %"`
	ReboundsOffensive       int     `json:"reboundsOffensive" isVisible:"true" display:"Off. Reb"`
	ReboundsDefensive       int     `json:"reboundsDefensive" isVisible:"true" display:"Def. Reb"`
	ReboundsTotal           int     `json:"reboundsTotal" isVisible:"true" display:"Rebounds"`
	Assists                 int     `json:"assists" isVisible:"true" display:"Assists"`
	Steals                  int     `json:"steals" isVisible:"true" display:"Steals"`
	Blocks                  int     `json:"blocks" isVisible:"true" display:"Blocks"`
	Turnovers               int     `json:"turnovers" isVisible:"true" display:"Turnovers"`
	FoulsPersonal           int     `json:"foulsPersonal" isVisible:"true" display:"Fouls"`
	Points                  int     `json:"points" isVisible:"true" display:"Points"`
	PlusMinusPoints         float64 `json:"plusMinusPoints" isVisible:"true" display:"+/- Points"`
}

type BoxScorePlayer struct {
	PersonId   int              `json:"personId" isVisible:"true" isID:"true" display:"ID"`
	FirstName  string           `json:"firstName" isVisible:"false"`
	FamilyName string           `json:"familyName" isVisible:"false"`
	NameI      string           `json:"nameI" isVisible:"true" display:"Name"`
	PlayerSlug string           `json:"playerSlug" isVisible:"false"`
	Position   string           `json:"position" isVisible:"true" display:"Position"`
	Comment    string           `json:"comment" isVisible:"false"`
	JerseyNum  string           `json:"jerseyNum" isVisible:"true" display:"No."`
	Statistics PlayerStatistics `json:"statistics" isVisible:"true"`
}

type BoxScorePlayers []BoxScorePlayer

// Player struct represents a player row with their current statistical averages based on the input parameters
// Can be totals, per game averages, per 48 minutes etc.
type Player struct {
	PlayerID    int     `json:"PLAYER_ID" isVisible:"false" isID:"true"`
	Rank        int     `json:"RANK" isVisible:"true" display:"Rank" width:"8"`
	PlayerName  string  `json:"PLAYER" isVisible:"true" display:"Player" width:"25"`
	TeamID      int     `json:"TEAM_ID" isVisible:"false" isID:"true"`
	TeamAbbr    string  `json:"TEAM" isVisible:"true" display:"Team" width:"8"`
	GamesPlayed int     `json:"GP" isVisible:"true" display:"GP" width:"8"`
	Minutes     float64 `json:"MIN" isVisible:"true" display:"Minutes"`
	FGM         float64 `json:"FGM" isVisible:"true" display:"FG Made"`
	FGA         float64 `json:"FGA" isVisible:"true" display:"FG Attempted"`
	FGPCT       float64 `json:"FG_PCT" percentage:"true" isVisible:"true" display:"FG%"`
	FG3PTM      float64 `json:"FG3PTM" isVisible:"true" display:"3PT Made"`
	FG3PTA      float64 `json:"FG3PTA" isVisible:"true" display:"3PT Attempted"`
	FG3PTPCT    float64 `json:"FG3PT_PCT" percentage:"true" isVisible:"true" display:"3PT%"`
	FTM         float64 `json:"FTM" isVisible:"true" display:"FT Made"`
	FTA         float64 `json:"FTA" isVisible:"true" display:"FT Attempted"`
	FTPCT       float64 `json:"FT_PCT" percentage:"true" isVisible:"true" display:"FT%"`
	OREB        float64 `json:"OREB" isVisible:"true" display:"Off. Reb"`
	DREB        float64 `json:"DREB" isVisible:"true" display:"Def. Reb"`
	REB         float64 `json:"REB" isVisible:"true" display:"Rebounds"`
	AST         float64 `json:"AST" isVisible:"true" display:"Assists"`
	STL         float64 `json:"STL" isVisible:"true" display:"Steals"`
	BLK         float64 `json:"BLK" isVisible:"true" display:"Blocks"`
	TOV         float64 `json:"TOV" isVisible:"true" display:"Turnovers"`
	PTS         float64 `json:"PTS" isVisible:"true" display:"Points"`
	EFF         float64 `json:"EFF" isVisible:"true" display:"Efficiency"`
}

type Players []Player

type IndexPlayer struct {
	PlayerID        int     `json:"PERSON_ID" isVisible:"true" isID:"true"`
	PlayerLastName  string  `json:"PLAYER_LAST_NAME" isVisible:"true" display:"Last Name"`
	PlayerFirstName string  `json:"PLAYER_FIRST_NAME" isVisible:"true" display:"First Name"`
	PlayerSlug      string  `json:"PLAYER_SLUG" isVisible:"false"`
	TeamSlug        string  `json:"TEAM_SLUG" isVisible:"false"`
	TeamID          int     `json:"TEAM_ID" isVisible:"false"`
	TeamCity        string  `json:"TEAM_CITY" isVisible:"false"`
	TeamName        string  `json:"TEAM" isVisible:"false"`
	TeamAbbr        string  `json:"TEAM_ABBREVIATION" isVisible:"false"`
	JerseyNumber    string  `json:"JERSEY_NUMBER" isVisible:"true" display:"Number" width:"10"`
	Position        string  `json:"POSITION" isVisible:"true" display:"Position" width:"10"`
	Height          string  `json:"HEIGHT" isVisible:"true" display:"Height" width:"10"`
	Weight          string  `json:"WEIGHT" isVisible:"true" display:"Weight" width:"10"`
	College         string  `json:"COLLEGE" isVisible:"true" display:"College" width:"25"`
	Country         string  `json:"COUNTRY" isVisible:"true" display:"Country"`
	DraftYear       int     `json:"DRAFT_YEAR" isVisible:"true" display:"Draft Year"`
	DraftRound      int     `json:"DRAFT_ROUND" isVisible:"true" display:"Draft Round"`
	DraftNumber     int     `json:"DRAFT_NUMBER" isVisible:"true" display:"Draft Number"`
	RosterStatus    float64 `json:"ROSTER_STATUS" isVisible:"false"`
	FromYear        string  `json:"FROM_YEAR" isVisible:"false"`
	ToYear          string  `json:"TO_YEAR" isVisible:"false"`
	Points          float64 `json:"PTS" isVisible:"true" display:"Points"`
	Rebounds        float64 `json:"REB" isVisible:"true" display:"Rebounds"`
	Assists         float64 `json:"AST" isVisible:"true" display:"Assists"`
	StatsTimeframe  string  `json:"STATS_TIMEFRAME" isVisible:"false"`
}

type IndexPlayers []IndexPlayer

func (pst PlayerStatistics) ToStringSlice() []string {
	return structToStringSlice(pst)
}

func (bsp BoxScorePlayer) ToStringSlice() []string {
	return structToStringSlice(bsp)
}

// ToStringSlice is a method on the Player type that enables the attributes of the type to be converted to strings
func (p Player) ToStringSlice() []string {
	return structToStringSlice(p)
}

// ToStringSlice is a method on the Players type that enables the attributes of the type to be converted to strings
func (ps Players) ToStringSlice() []string {
	return structToStringSlice(ps)
}

func (ip IndexPlayer) ToStringSlice() []string {
	return structToStringSlice(ip)
}

func (ips IndexPlayers) ToStringSlice() []string {
	return structToStringSlice(ips)
}
