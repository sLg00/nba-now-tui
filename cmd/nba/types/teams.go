package types

type Team struct {
	LeagueID                string  `json:"LeagueID" isVisible:"false"`
	SeasonID                string  `json:"SeasonID" isVisible:"false"`
	TeamID                  int     `json:"TeamID" isVisible:"false" isID:"true"`
	TeamCity                string  `json:"TeamCity" isVisible:"true" display:"City" width:"15"`
	TeamName                string  `json:"TeamName" isVisible:"true" display:"Name" width:"15"`
	TeamSlug                string  `json:"TeamSlug" isVisible:"false"`
	Conference              string  `json:"Conference" isVisible:"false"`
	ConferenceRecord        string  `json:"ConferenceRecord" isVisible:"true" display:"Conf. Record"`
	PlayoffRank             int     `json:"PlayoffRank" isVisible:"true" display:"PO Rank"`
	ClinchIndicator         string  `json:"ClinchIndicator" isVisible:"true" display:"Clinch Indicator"`
	Division                string  `json:"Division" isVisible:"true" display:"Division"`
	DivisionRecord          string  `json:"DivisionRecord" isVisible:"true" display:"Division Record"`
	DivisionRank            int     `json:"DivisionRank" isVisible:"true" display:"Division Rank"`
	Wins                    int     `json:"WINS" isVisible:"true" display:"Wins"`
	Losses                  int     `json:"LOSSES" isVisible:"true" display:"Losses"`
	WinPCT                  float64 `json:"WinPCT" percentage:"true" display:"Win%"`
	LeagueRank              int     `json:"LeagueRank" isVisible:"true" display:"League Rank"`
	Record                  string  `json:"Record" isVisible:"true" display:"Record"`
	Home                    string  `json:"HOME" isVisible:"true" display:"Home"`
	Road                    string  `json:"ROAD" isVisible:"true" display:"Road"`
	L10                     string  `json:"L10" isVisible:"true" display:"Last10"`
	Last10Home              string  `json:"Last10Home" isVisible:"true" display:"Last10Home"`
	Last10Road              string  `json:"Last10Road" isVisible:"true" display:"Last10Road"`
	OT                      string  `json:"OT" isVisible:"true" display:"OT"`
	ThreePTSOrLess          string  `json:"ThreePTSOrLess" isVisible:"false"`
	TenPTSOrMore            string  `json:"TenPTSOrMore" isVisible:"false"`
	LongHomeStreak          int     `json:"LongHomeStreak" isVisible:"false"`
	StrLongHomeStreak       string  `json:"strLongHomeStreak" isVisible:"false"`
	LongRoadStreak          int     `json:"LongRoadStreak" isVisible:"false"`
	StrLongRoadStreak       string  `json:"strLongRoadStreak" isVisible:"false"` //
	LongWinStreak           int     `json:"LongWinStreak" isVisible:"true" display:"LongWinStreak"`
	LongLossStreak          int     `json:"LongLossStreak" isVisible:"true" display:"LongLossStreak"`
	CurrentHomeStreak       int     `json:"CurrentHomeStreak" isVisible:"false"`
	StrCurrentHomeStreak    string  `json:"strCurrentHomeStreak" isVisible:"false"` //
	CurrentRoadStreak       int     `json:"CurrentRoadStreak" isVisible:"false"`
	StrCurrentRoadStreak    string  `json:"strCurrentRoadStreak" isVisible:"false"`
	CurrentStreak           int     `json:"CurrentStreak" isVisible:"true" display:"CurrentStreak"`
	StrCurrentStreak        string  `json:"strCurrentStreak" isVisible:"false"`
	ConferenceGamesBack     float64 `json:"ConferenceGamesBack" isVisible:"false"`
	DivisionGamesBack       float64 `json:"DivisionGamesBack" isVisible:"false"`
	ClinchedConferenceTitle int     `json:"ClinchedConferenceTitle" isVisible:"true" display:"Clinched Conf.""`
	ClinchedDivisionTitle   int     `json:"ClinchedDivisionTitle" isVisible:"true" display:"Clinched Div."`
	ClinchedPlayoffBirth    int     `json:"ClinchedPlayoffBirth" isVisible:"true" display:"Clinched PO"`
	ClinchedPlayIn          int     `json:"ClinchedPlayIn" isVisible:"true" display:"Clinched PlayIn"`
	EliminatedConference    int     `json:"EliminatedConference" isVisible:"false"`
	EliminatedDivision      int     `json:"EliminatedDivision" isVisible:"false"`
	AheadAtHalf             string  `json:"AheadAtHalf" isVisible:"false"`
	BehindAtHalf            string  `json:"BehindAtHalf" isVisible:"false"`
	TiedAtHalf              string  `json:"TiedAtHalf" isVisible:"false"`
	AheadAtThird            string  `json:"AheadAtThird" isVisible:"false"`
	BehindAtThird           string  `json:"BehindAtThird" isVisible:"false"`
	TiedAtThird             string  `json:"TiedAtThird" isVisible:"false"`
	Score100PTS             string  `json:"Score100PTS" isVisible:"false"`
	OppScore100PTS          string  `json:"OppScore100PTS" isVisible:"false"`
	OppOver500              string  `json:"OppOver500" isVisible:"false"`
	LeadInFGPCT             string  `json:"LeadInFGPCT" isVisible:"false"`
	LeadInReb               string  `json:"LeadInReb" isVisible:"false"`
	FewerTurnovers          string  `json:"FewerTurnovers" isVisible:"false"`
	PointsPG                float64 `json:"PointsPG" isVisible:"true" display:"PPG"`
	OppPointsPG             float64 `json:"OppPointsPG" isVisible:"true" display:"Opp PPG"`
	DiffPointsPG            float64 `json:"DiffPointsPG" isVisible:"true" display:"Diff PPG"`
	VsEast                  string  `json:"vsEast" isVisible:"false"`
	VsAtlantic              string  `json:"vsAtlantic" isVisible:"false"`
	VsCentral               string  `json:"vsCentral" isVisible:"false"`
	VsSoutheast             string  `json:"vsSoutheast" isVisible:"false"`
	VsWest                  string  `json:"vsWest" isVisible:"false"`
	VsNorthwest             string  `json:"vsNorthwest" isVisible:"false"`
	VsPacific               string  `json:"vsPacific" isVisible:"false"`
	VsSouthwest             string  `json:"vsSouthwest" isVisible:"false"`
	Jan                     string  `json:"Jan" isVisible:"false"`
	Feb                     string  `json:"Feb" isVisible:"false"`
	Mar                     string  `json:"Mar" isVisible:"false"`
	Apr                     string  `json:"Apr" isVisible:"false"`
	May                     string  `json:"May" isVisible:"false"`
	Jun                     string  `json:"Jun" isVisible:"false"`
	Jul                     string  `json:"Jul" isVisible:"false"`
	Aug                     string  `json:"Aug" isVisible:"false"`
	Sep                     string  `json:"Sep" isVisible:"false"`
	Oct                     string  `json:"Oct" isVisible:"false"`
	Nov                     string  `json:"Nov" isVisible:"false"`
	Dec                     string  `json:"Dec" isVisible:"false"`
	Score80Plus             string  `json:"Score_80_Plus" isVisible:"false"`
	OppScore80Plus          string  `json:"Opp_Score_80_Plus" isVisible:"false"`
	ScoreBelow80            string  `json:"Score_Below_80" isVisible:"false"`
	OppScoreBelow80         string  `json:"Opp_Score_Below_80" isVisible:"false"`
	TotalPoints             int     `json:"TotalPoints" isVisible:"false"`
	OppTotalPoints          int     `json:"OppTotalPoints" isVisible:"false"`
	DiffTotalPoints         int     `json:"DiffTotalPoints" isVisible:"false"`
	LeagueGamesBack         float64 `json:"LeagueGamesBack" isVisible:"false"`
	PlayoffSeeding          int     `json:"PlayoffSeeding" isVisible:"true" display:"PlayOff Seeding"`
	ClinchedPostSeason      int     `json:"ClinchedPostSeason" isVisible:"false"`
}

type Teams []Team

type TeamCommonInfo struct {
	TeamID         int     `json:"TEAM_ID" isVisible:"false" isID:"true"`
	SeasonYear     string  `json:"SEASON_YEAR" isVisible:"false" isID:"false"`
	TeamCity       string  `json:"TEAM_CITY" isVisible:"false" isID:"false"`
	TeamName       string  `json:"TEAM_NAME" isVisible:"false" isID:"false"`
	TeamAbbrev     string  `json:"TEAM_ABBREV" isVisible:"false" isID:"false"`
	TeamConference string  `json:"TEAM_CONFERENCE" isVisible:"false" isID:"false"`
	TeamDivision   string  `json:"TEAM_DIVISION" isVisible:"false" isID:"false"`
	TeamCode       string  `json:"TEAM_CODE" isVisible:"false" isID:"false"`
	TeamSlug       string  `json:"TEAM_SLUG" isVisible:"false" isID:"false"`
	Wins           int     `json:"W" isVisible:"true" isID:"false" display:"Wins"`
	Losses         int     `json:"L" isVisible:"true" isID:"false" display:"Losses"`
	WinPct         float64 `json:"PCT" percentage:"true" isVisible:"true" isID:"false" display:"Win %"`
	ConfRank       int     `json:"CONF_RANK" isVisible:"true" isID:"false" display:"Conf rank"`
	DivRank        int     `json:"DIV_RANK" isVisible:"true" isID:"false" display:"Div rank"`
	MinYear        string  `json:"MIN_YEAR" isVisible:"false" isID:"false"`
	MaxYear        string  `json:"MAX_YEAR" isVisible:"false" isID:"false"`
}

// ToStringSlice is a method on the TeamCommonInfo type that enables the attributes of the type to be converted to strings
func (ti TeamCommonInfo) ToStringSlice() []string {
	return structToStringSlice(ti)
}

// ToStringSlice is a method on the Team type that enables the attributes of the type to be converted to strings
func (t Team) ToStringSlice() []string {
	return structToStringSlice(t)
}

// ToStringSlice is a method on the Teams type that enables the attributes of type to be converted to strings
func (ts Teams) ToStringSlice() []string {
	return structToStringSlice(ts)
}

// SplitStandingsPerConference is a method on the Teams type that separates the Western and Easter conference teams into
// separate instances of the Teams type
func (ts Teams) SplitStandingsPerConference() (Teams, Teams) {
	teams := ts
	var westTeams Teams
	var eastTeams Teams

	for _, team := range teams {
		if team.Conference == "East" {
			eastTeams = append(eastTeams, team)
		} else if team.Conference == "West" {
			westTeams = append(westTeams, team)
		}
	}
	return eastTeams, westTeams
}
