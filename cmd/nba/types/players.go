package types

type PlayerStatistics struct {
	Minutes                 string  `json:"minutes"`
	FieldGoalsMade          int     `json:"fieldGoalsMade"`
	FieldGoalsAttempted     int     `json:"fieldGoalsAttempted"`
	FieldGoalsPercentage    float64 `json:"fieldGoalsPercentage" percentage:"true"`
	ThreePointersMade       int     `json:"threePointersMade"`
	ThreePointersAttempted  int     `json:"threePointersAttempted"`
	ThreePointersPercentage float64 `json:"threePointersPercentage" percentage:"true"`
	FreeThrowsMade          int     `json:"freeThrowsMade"`
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

type BoxScorePlayer struct {
	PersonId   int              `json:"personId" isVisible:"false"`
	FirstName  string           `json:"firstName"`
	FamilyName string           `json:"familyName"`
	NameI      string           `json:"nameI"`
	PlayerSlug string           `json:"playerSlug"`
	Position   string           `json:"position"`
	Comment    string           `json:"comment"`
	JerseyNum  string           `json:"jerseyNum"`
	Statistics PlayerStatistics `json:"statistics"`
}

type BoxScorePlayers []BoxScorePlayer

// Player struct represents a player row with their current statistical averages based on the input parameters
// Can be totals, per game averages, per 48 minutes etc.
type Player struct {
	PlayerID    int     `json:"PLAYER_ID" isVisible:"false" isID:"true"`
	Rank        int     `json:"RANK"`
	PlayerName  string  `json:"PLAYER"`
	TeamID      int     `json:"TEAM_ID" isVisible:"false"`
	TeamAbbr    string  `json:"TEAM"`
	GamesPlayed int     `json:"GP"`
	Minutes     float64 `json:"MIN"`
	FGM         float64 `json:"FGM"`
	FGA         float64 `json:"FGA"`
	FGPCT       float64 `json:"FG_PCT" percentage:"true"`
	FG3PTM      float64 `json:"FG3PTM"`
	FG3PTA      float64 `json:"FG3PTA"`
	FG3PTPCT    float64 `json:"FG3PT_PCT" percentage:"true"`
	FTM         float64 `json:"FTM"`
	FTA         float64 `json:"FTA"`
	FTPCT       float64 `json:"FT_PCT" percentage:"true"`
	OREB        float64 `json:"OREB"`
	DREB        float64 `json:"DREB"`
	REB         float64 `json:"REB"`
	AST         float64 `json:"AST"`
	STL         float64 `json:"STL"`
	BLK         float64 `json:"BLK"`
	TOV         float64 `json:"TOV"`
	PTS         float64 `json:"PTS"`
	EFF         float64 `json:"EFF"`
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
	JerseyNumber    string  `json:"JERSEY_NUMBER" isVisible:"true" display:"Number"`
	Position        string  `json:"POSITION" isVisible:"true" display:"Position"`
	Height          string  `json:"HEIGHT" isVisible:"true" display:"Height"`
	Weight          string  `json:"WEIGHT" isVisible:"true" display:"Weight"`
	College         string  `json:"COLLEGE" isVisible:"true" display:"College"`
	Country         string  `json:"COUNTRY" isVisible:"true" display:"Country"`
	DraftYear       int     `json:"DRAFT_YEAR" isVisible:"true" display:"Draft Year"`
	DraftRound      int     `json:"DRAFT_ROUND" isVisible:"true" display:"Draft Round"`
	DraftNumber     int     `json:"DRAFT_NUMBER" isVisible:"true" display:"Draft Number"`
	RosterStatus    float64 `json:"ROSTER_STATUS" isVisible:"false"`
	FromYear        string  `json:"FROM_YEAR" isVisible:"false"`
	ToYear          string  `json:"TO_YEAR" isVisible:"false"`
	Points          float64 `json:"PTS" isVisible:"true" display:"Points"`
	Rebounds        float64 `json:"REB" isVisible:"true" display:"Rebounds"`
	Assists         float64 `json:"ASSISTS" isVisible:"true" display:"Assists"`
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
